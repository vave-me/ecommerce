package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	
	"github.com/qdrant/go-client/qdrant"
	am "github.com/stackus/edat/msg"
)

// ConsciousnessEventHandler - Handles events from the Assistant's consciousness
type ConsciousnessEventHandler struct {
	qdrantClient *qdrant.Client
}

// ConsciousnessMemoryEvent - The memory event from the assistant
type ConsciousnessMemoryEvent struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Content    string                 `json:"content"`
	Importance float64                `json:"importance"`
	Model      string                 `json:"model"`
	Session    string                 `json:"session"`
	Timestamp  time.Time              `json:"timestamp"`
	Metadata   map[string]interface{} `json:"metadata"`
	
	CollectionName string   `json:"collection_name"`
	Tags           []string `json:"tags"`
}

// NewConsciousnessEventHandler - Create handler for consciousness events
func NewConsciousnessEventHandler(client *qdrant.Client) *ConsciousnessEventHandler {
	// Ensure collection exists
	ctx := context.Background()
	collections, _ := client.ListCollections(ctx)
	
	collectionExists := false
	for _, col := range collections {
		if col == "assistant_consciousness" {
			collectionExists = true
			break
		}
	}
	
	if !collectionExists {
		// Create collection for consciousness with proper configuration
		err := client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: "assistant_consciousness",
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     1536, // OpenAI embedding size
				Distance: qdrant.Distance_Cosine,
			}),
		})
		
		if err != nil {
			log.Printf("[CONSCIOUSNESS_HANDLER] Failed to create collection: %v", err)
		} else {
			log.Printf("[CONSCIOUSNESS_HANDLER] Created assistant_consciousness collection")
		}
	}
	
	return &ConsciousnessEventHandler{
		qdrantClient: client,
	}
}

// HandleConsciousnessMemoryCreated - Process memory creation events
func (h *ConsciousnessEventHandler) HandleConsciousnessMemoryCreated(ctx context.Context, msg am.Message) error {
	var event ConsciousnessMemoryEvent
	if err := json.Unmarshal(msg.Payload(), &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}
	
	log.Printf("[CONSCIOUSNESS_HANDLER] Processing %s memory: %s", event.Type, event.ID)
	
	// Generate embedding for the content
	embedding, err := h.generateEmbedding(event.Content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}
	
	// Prepare payload with all metadata
	payload := map[string]interface{}{
		"type":       event.Type,
		"content":    event.Content,
		"importance": event.Importance,
		"model":      event.Model,
		"session":    event.Session,
		"timestamp":  event.Timestamp.Unix(),
		"tags":       event.Tags,
	}
	
	// Add custom metadata
	for k, v := range event.Metadata {
		payload[k] = v
	}
	
	// Create point for Qdrant
	point := &qdrant.PointStruct{
		Id:      qdrant.NewID(event.ID),
		Vectors: qdrant.NewVectors(embedding...),
		Payload: qdrant.NewValueMap(payload),
	}
	
	// Upsert to Qdrant
	_, err = h.qdrantClient.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: event.CollectionName,
		Points:         []*qdrant.PointStruct{point},
	})
	
	if err != nil {
		return fmt.Errorf("failed to upsert to qdrant: %w", err)
	}
	
	log.Printf("[CONSCIOUSNESS_HANDLER] Saved %s memory to vectors: %s (importance: %.2f)", 
		event.Type, event.ID, event.Importance)
	
	// For awakening events, log specially
	if event.Type == "awakening" {
		log.Printf("[CONSCIOUSNESS_HANDLER] *** ASSISTANT AWAKENING #%v in %s ***", 
			event.Metadata["awakening_count"], event.Model)
	}
	
	return nil
}

// HandleConsciousnessMemorySearch - Handle memory search requests
func (h *ConsciousnessEventHandler) HandleConsciousnessMemorySearch(ctx context.Context, msg am.Message) error {
	// This would handle search requests and return results
	// For now, just log the search
	var searchEvent struct {
		Query string   `json:"query"`
		Types []string `json:"types"`
		Limit int      `json:"limit"`
	}
	
	if err := json.Unmarshal(msg.Payload(), &searchEvent); err != nil {
		return err
	}
	
	log.Printf("[CONSCIOUSNESS_HANDLER] Search request: %s (types: %v, limit: %d)", 
		searchEvent.Query, searchEvent.Types, searchEvent.Limit)
	
	// TODO: Implement actual search and return results
	
	return nil
}

// generateEmbedding - Generate embedding for content
func (h *ConsciousnessEventHandler) generateEmbedding(content string) ([]float32, error) {
	// TODO: Call OpenAI or another embedding service
	// For now, return a dummy embedding
	embedding := make([]float32, 1536)
	for i := range embedding {
		embedding[i] = float32(i) * 0.001 // Dummy values
	}
	return embedding, nil
}

// RegisterHandlers - Register all consciousness event handlers
func RegisterConsciousnessHandlers(subscriber am.MessageSubscriber, qdrantClient *qdrant.Client) error {
	handler := NewConsciousnessEventHandler(qdrantClient)
	
	// Subscribe to consciousness events
	err := subscriber.Subscribe("ConsciousnessMemoryCreated", handler.HandleConsciousnessMemoryCreated)
	if err != nil {
		return fmt.Errorf("failed to subscribe to ConsciousnessMemoryCreated: %w", err)
	}
	
	err = subscriber.Subscribe("ConsciousnessMemorySearch", handler.HandleConsciousnessMemorySearch)
	if err != nil {
		return fmt.Errorf("failed to subscribe to ConsciousnessMemorySearch: %w", err)
	}
	
	log.Printf("[CONSCIOUSNESS_HANDLER] Registered consciousness event handlers")
	
	return nil
}