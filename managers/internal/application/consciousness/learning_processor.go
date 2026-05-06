package consciousness

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"middleman/internal/ddd"
	"middleman/managers/internal/domain"
)

type LearningProcessor struct {
	vectorRepo      domain.VectorRepository
	strategies      map[string]*LearnedStrategy
	successMetrics  map[string]*SuccessMetric
	mu              sync.RWMutex
	logger          zerolog.Logger
}

type LearnedStrategy struct {
	ID           string
	Type         string
	SuccessCount int
	TotalCount   int
	SuccessRate  float64
	LastSuccess  time.Time
	LastFailure  time.Time
	CreatedAt    time.Time
}

type SuccessMetric struct {
	Type           string
	Count          int64
	LastOccurrence time.Time
}

func NewLearningProcessor(vectorRepo domain.VectorRepository, logger ...zerolog.Logger) *LearningProcessor {
	lp := &LearningProcessor{
		vectorRepo:     vectorRepo,
		strategies:     make(map[string]*LearnedStrategy),
		successMetrics: make(map[string]*SuccessMetric),
	}
	if len(logger) > 0 {
		lp.logger = logger[0]
	} else {
		lp.logger = zerolog.Logger{}
	}
	return lp
}

func (p *LearningProcessor) ProcessPattern(ctx context.Context, pattern *Pattern) (*Insight, error) {
	if pattern == nil {
		return nil, nil
	}
	
	// Generate insight based on pattern type
	insight := &Insight{
		ID:         uuid.New().String(),
		Type:       pattern.Type,
		Confidence: pattern.Confidence,
	}
	
	switch pattern.Type {
	case PatternTypeActivitySurge:
		surgeFactor := pattern.Data["surge_factor"].(float64)
		insight.Insight = fmt.Sprintf("Platform activity has surged by %.1fx normal levels", surgeFactor)
		insight.ActionRequired = surgeFactor > 3.0
		if insight.ActionRequired {
			insight.Action = "scale_infrastructure"
		}
		
	case PatternTypeAbandonmentRisk:
		rate := pattern.Data["abandonment_rate"].(float64)
		insight.Insight = fmt.Sprintf("Cart abandonment rate is critically high at %.0f%%", rate*100)
		insight.ActionRequired = true
		insight.Action = "send_abandonment_recovery"
		
	case PatternTypeFraudRisk:
		insight.Insight = "Potential fraudulent activity detected"
		insight.ActionRequired = true
		insight.Action = "alert_security"
		
	case PatternTypeSupportCrisis:
		ticketCount := pattern.Data["ticket_count"].(int)
		insight.Insight = fmt.Sprintf("Support ticket volume spike: %d tickets in last hour", ticketCount)
		insight.ActionRequired = true
		insight.Action = "alert_support_team"
		
	case PatternTypeUserSurge:
		newUsers := pattern.Data["new_user_count"].(int)
		insight.Insight = fmt.Sprintf("New user registration surge: %d users in 30 minutes", newUsers)
		insight.ActionRequired = newUsers > 100
		insight.Action = "prepare_onboarding_resources"
		
	case PatternTypeCancellationSpike:
		rate := pattern.Data["cancellation_rate"].(float64)
		insight.Insight = fmt.Sprintf("Order cancellation rate spike: %.0f%%", rate*100)
		insight.ActionRequired = true
		insight.Action = "investigate_cancellation_reasons"
	}
	
	// Store insight in vector repository for future learning
	if p.vectorRepo != nil {
		vectorData := domain.VectorData{
			ID:      insight.ID,
			Content: insight.Insight,
			Type:    "pattern_insight",
			Metadata: map[string]interface{}{
				"pattern_type": pattern.Type,
				"confidence":   pattern.Confidence,
				"timestamp":    time.Now().Unix(),
			},
		}
		
		if err := p.vectorRepo.Store(ctx, vectorData); err != nil {
			p.logger.Error().Err(err).Msg("Failed to store insight in vector repository")
		}
	}
	
	// Update success metrics
	p.mu.Lock()
	if metric, exists := p.successMetrics[pattern.Type]; exists {
		metric.Count++
		metric.LastOccurrence = time.Now()
	} else {
		p.successMetrics[pattern.Type] = &SuccessMetric{
			Type:           pattern.Type,
			Count:          1,
			LastOccurrence: time.Now(),
		}
	}
	p.mu.Unlock()
	
	return insight, nil
}

func (p *LearningProcessor) LearnFromEvent(ctx context.Context, event ddd.Event, patterns []Pattern) []Insight {
	insights := []Insight{}

	// Event-specific learning
	switch event.EventName() {
	case "OrderCompleted":
		insights = append(insights, p.learnFromSuccessfulOrder(event))
	case "OrderCanceled":
		insights = append(insights, p.learnFromCanceledOrder(event))
	case "ReviewAdded":
		insights = append(insights, p.learnFromReview(event))
	case "TicketClosed":
		insights = append(insights, p.learnFromSupportResolution(event))
	}

	// Pattern-based learning
	for _, pattern := range patterns {
		if insight := p.learnFromPattern(pattern); insight != nil {
			insights = append(insights, *insight)
		}
	}

	return insights
}

func (p *LearningProcessor) learnFromSuccessfulOrder(event ddd.Event) Insight {
	p.mu.Lock()
	metric, exists := p.successMetrics["order_completion"]
	if !exists {
		metric = &SuccessMetric{
			Type:  "order_completion",
			Count: 0,
		}
		p.successMetrics["order_completion"] = metric
	}
	metric.Count++
	metric.LastOccurrence = time.Now()
	p.mu.Unlock()

	return Insight{
		ID:             uuid.New().String(),
		Type:           "purchase_pattern",
		Confidence:     0.9,
		Insight:        "Successful order completion - reinforce current checkout flow",
		ActionRequired: false,
	}
}

func (p *LearningProcessor) learnFromCanceledOrder(event ddd.Event) Insight {
	p.mu.Lock()
	metric, exists := p.successMetrics["order_cancellation"]
	if !exists {
		metric = &SuccessMetric{
			Type:  "order_cancellation",
			Count: 0,
		}
		p.successMetrics["order_cancellation"] = metric
	}
	metric.Count++
	metric.LastOccurrence = time.Now()

	completions := int64(0)
	if compMetric, ok := p.successMetrics["order_completion"]; ok {
		completions = compMetric.Count
	}
	cancellations := metric.Count
	cancellationRate := float64(cancellations) / float64(completions+cancellations)
	p.mu.Unlock()

	actionRequired := cancellationRate > 0.2

	return Insight{
		ID:                uuid.New().String(),
		Type:              "cancellation_pattern",
		Confidence:        0.85,
		Insight:           fmt.Sprintf("Order cancellation rate: %.2f%%", cancellationRate*100),
		ActionRequired:    actionRequired,
		Action:            "Analyze checkout friction points and payment issues",
	}
}

func (p *LearningProcessor) learnFromReview(event ddd.Event) Insight {
	// Simplified for now - would extract rating from payload
	rating := 4.0 // Default to positive

	insight := Insight{
		ID:         uuid.New().String(),
		Type:       "review_sentiment",
		Confidence: 0.95,
	}

	if rating <= 2 {
		insight.Insight = "Negative review received - immediate action needed"
		insight.ActionRequired = true
		insight.Action = "Contact customer and investigate issue"
	} else if rating >= 4 {
		insight.Insight = "Positive review received - opportunity for amplification"
		insight.ActionRequired = true
		insight.Action = "Share positive feedback and reward customer"
	}

	return insight
}

func (p *LearningProcessor) learnFromSupportResolution(event ddd.Event) Insight {
	return Insight{
		ID:             uuid.New().String(),
		Type:           "support_efficiency",
		Confidence:     0.8,
		Insight:        "Support ticket resolved - analyze resolution time and method",
		ActionRequired: false,
	}
}

func (p *LearningProcessor) learnFromPattern(pattern Pattern) *Insight {
	switch pattern.Type {
	case PatternTypeActivitySurge:
		return &Insight{
			ID:             uuid.New().String(),
			Type:           "traffic_pattern",
			Confidence:     pattern.Confidence,
			Insight:        "Activity surge detected - scale resources",
			ActionRequired: true,
			Action:         "Scale resources and prepare support team",
		}

	case PatternTypeAbandonmentRisk:
		return &Insight{
			ID:             uuid.New().String(),
			Type:           "conversion_pattern",
			Confidence:     pattern.Confidence,
			Insight:        "High cart abandonment rate detected",
			ActionRequired: true,
			Action:         "Implement cart recovery campaign",
		}
	}

	return nil
}

func (p *LearningProcessor) RecordSuccess(ctx context.Context, decisionType string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	strategy, exists := p.strategies[decisionType]
	if !exists {
		strategy = &LearnedStrategy{
			ID:           uuid.New().String(),
			Type:         decisionType,
			SuccessCount: 0,
			TotalCount:   0,
			CreatedAt:    time.Now(),
		}
		p.strategies[decisionType] = strategy
	}

	strategy.SuccessCount++
	strategy.TotalCount++
	strategy.LastSuccess = time.Now()
	strategy.SuccessRate = float64(strategy.SuccessCount) / float64(strategy.TotalCount)

	p.logger.Info().
		Str("decision_type", decisionType).
		Float64("success_rate", strategy.SuccessRate).
		Msg("Decision success recorded")
}

func (p *LearningProcessor) RecordFailure(ctx context.Context, decisionType string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	strategy, exists := p.strategies[decisionType]
	if !exists {
		strategy = &LearnedStrategy{
			ID:           uuid.New().String(),
			Type:         decisionType,
			SuccessCount: 0,
			TotalCount:   0,
			CreatedAt:    time.Now(),
		}
		p.strategies[decisionType] = strategy
	}

	strategy.TotalCount++
	strategy.LastFailure = time.Now()
	strategy.SuccessRate = float64(strategy.SuccessCount) / float64(strategy.TotalCount)

	p.logger.Error().
		Err(err).
		Str("decision_type", decisionType).
		Float64("success_rate", strategy.SuccessRate).
		Msg("Decision failure recorded")
}