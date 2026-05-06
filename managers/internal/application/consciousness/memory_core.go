package consciousness

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"middleman/internal/ddd"
	"middleman/managers/internal/domain"
)

type EventMemory struct {
	ID          string
	Type        string
	AggregateID string
	Timestamp   time.Time
	Importance  float64
	Data        interface{} // Changed from Payload to Data to match usage
}

type MemoryCore struct {
	vectorRepo      domain.VectorRepository
	recentEvents    *CircularBuffer
	importantEvents map[string]EventMemory
	mu              sync.RWMutex
	logger          zerolog.Logger
}

func NewMemoryCore(vectorRepo domain.VectorRepository, logger zerolog.Logger) *MemoryCore {
	return &MemoryCore{
		vectorRepo:      vectorRepo,
		recentEvents:    NewCircularBuffer(10000),
		importantEvents: make(map[string]EventMemory),
		logger:          logger,
	}
}

func (m *MemoryCore) StoreEvent(ctx context.Context, event ddd.Event) error {
	// Try type assertion to get AggregateID if the concrete type supports it
	aggregateID := ""
	if eventWithAggregateID, ok := event.(interface{ AggregateID() string }); ok {
		aggregateID = eventWithAggregateID.AggregateID()
	}
	
	memory := EventMemory{
		ID:          event.ID(),
		Type:        event.EventName(),
		AggregateID: aggregateID,
		Timestamp:   event.OccurredAt(),
		Importance:  m.calculateImportance(event),
		Data:        event.Payload(),
	}

	m.recentEvents.Add(memory)

	if memory.Importance > 0.7 {
		m.mu.Lock()
		m.importantEvents[memory.ID] = memory
		m.mu.Unlock()
	}

	return nil
}

func (m *MemoryCore) GetRecentEvents(count int) []EventMemory {
	return m.recentEvents.GetRecent(count)
}

func (m *MemoryCore) calculateImportance(event ddd.Event) float64 {
	criticalEvents := map[string]float64{
		"OrderCompleted":   1.0,
		"PaymentCompleted": 1.0,
		"TicketCreated":    0.9,
		"OrderCanceled":    0.9,
		"ReviewAdded":      0.85,
	}

	if importance, ok := criticalEvents[event.EventName()]; ok {
		return importance
	}

	highImportanceEvents := map[string]float64{
		"UserCreated":        0.8,
		"ProductAdded":       0.75,
		"ProductSold":        0.8,
		"MessageSent":        0.7,
		"WishlistItemAdded":  0.65,
		"BasketItemAdded":    0.6,
	}

	if importance, ok := highImportanceEvents[event.EventName()]; ok {
		return importance
	}

	return 0.5
}

// CircularBuffer for recent event storage
type CircularBuffer struct {
	data     []interface{}
	size     int
	position int
	mu       sync.RWMutex
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		data: make([]interface{}, size),
		size: size,
	}
}

func (b *CircularBuffer) Add(item interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data[b.position] = item
	b.position = (b.position + 1) % b.size
}

func (b *CircularBuffer) GetRecent(count int) []EventMemory {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if count > b.size {
		count = b.size
	}

	result := make([]EventMemory, 0, count)
	start := (b.position - count + b.size) % b.size

	for i := 0; i < count; i++ {
		pos := (start + i) % b.size
		if b.data[pos] != nil {
			if memory, ok := b.data[pos].(EventMemory); ok {
				result = append(result, memory)
			}
		}
	}

	return result
}

// GetStatus returns the current consciousness status
func (m *MemoryCore) GetStatus(ctx context.Context) ConsciousnessStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	recentEvents := m.recentEvents.GetRecent(1000)
	eventsProcessed := int64(len(recentEvents))
	
	// Count decisions made (simplified - based on important events)
	decisionsMade := int64(len(m.importantEvents))
	
	// Determine health based on recent activity
	health := "healthy"
	if len(recentEvents) == 0 {
		health = "idle"
	} else if eventsProcessed > 500 {
		health = "busy"
	}
	
	lastActivity := time.Now()
	if len(recentEvents) > 0 {
		lastActivity = recentEvents[0].Timestamp
	}
	
	return ConsciousnessStatus{
		Active:          true,
		EventsProcessed: eventsProcessed,
		DecisionsMade:   decisionsMade,
		LastActivity:    lastActivity,
		Health:          health,
	}
}