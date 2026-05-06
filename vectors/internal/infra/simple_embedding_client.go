package infra

import (
	"context"
	"fmt"
	"middleman/vectors/internal/ports"
	"time"
)

// SimpleEmbeddingClient provides a basic implementation for Qdrant + Redis
type SimpleEmbeddingClient struct {
	model      string
	dimensions int
}

// NewSimpleEmbeddingClient creates a new simple embedding client
func NewSimpleEmbeddingClient(model string, dimensions int) *SimpleEmbeddingClient {
	if model == "" {
		model = "text-embedding-3-small"
	}
	if dimensions == 0 {
		dimensions = 1536
	}
	return &SimpleEmbeddingClient{
		model:      model,
		dimensions: dimensions,
	}
}

// GenerateEmbedding generates a simple embedding
func (s *SimpleEmbeddingClient) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	// Create a deterministic embedding based on text content
	embedding := make([]float32, s.dimensions)
	textBytes := []byte(text)

	for i := 0; i < s.dimensions; i++ {
		if i < len(textBytes) {
			embedding[i] = float32(textBytes[i%len(textBytes)]) / 255.0
		} else {
			embedding[i] = 0.1 // Default value
		}
	}

	return embedding, nil
}

// GenerateBatchEmbeddings generates embeddings for multiple texts
func (s *SimpleEmbeddingClient) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("texts array cannot be empty")
	}

	embeddings := make([][]float32, len(texts))
	for i, text := range texts {
		embedding, err := s.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for text %d: %w", i, err)
		}
		embeddings[i] = embedding
	}

	return embeddings, nil
}

// GetDimensions returns the embedding dimensions
func (s *SimpleEmbeddingClient) GetDimensions() int {
	return s.dimensions
}

// GetModel returns the model name
func (s *SimpleEmbeddingClient) GetModel() string {
	return s.model
}

// HealthCheck checks if the client is healthy
func (s *SimpleEmbeddingClient) HealthCheck(ctx context.Context) error {
	return nil // Always healthy for simple implementation
}

// GenerateEntityEmbedding generates embedding for entity data
func (s *SimpleEmbeddingClient) GenerateEntityEmbedding(ctx context.Context, entityData map[string]interface{}) ([]float32, error) {
	// Extract text from entity data
	text := s.extractEntityText(entityData)
	return s.GenerateEmbedding(ctx, text)
}

// GenerateEmbeddingWithPrompt generates embedding with prompt
func (s *SimpleEmbeddingClient) GenerateEmbeddingWithPrompt(ctx context.Context, text string, prompt string) ([]float32, error) {
	// For simple implementation, ignore prompt and just generate embedding
	return s.GenerateEmbedding(ctx, text)
}

// GenerateOptimizedEmbedding generates optimized embedding
func (s *SimpleEmbeddingClient) GenerateOptimizedEmbedding(ctx context.Context, entityType string, entityData map[string]interface{}, optimization string) ([]float32, error) {
	// For simple implementation, ignore optimization
	return s.GenerateEntityEmbedding(ctx, entityData)
}

// IsPromptEnabled returns whether prompt enhancement is supported
func (s *SimpleEmbeddingClient) IsPromptEnabled() bool {
	return false // Simple implementation doesn't support prompts
}

// GetProviderInfo returns provider information
func (s *SimpleEmbeddingClient) GetProviderInfo() ports.EmbeddingProviderInfo {
	return ports.EmbeddingProviderInfo{
		Provider:     "simple",
		Model:        s.model,
		Dimensions:   s.dimensions,
		MaxTokens:    8192,
		RateLimit:    1000,
		CostPerToken: 0.0,
		ConnectedAt:  time.Now(),
	}
}

// extractEntityText extracts text from entity data
func (s *SimpleEmbeddingClient) extractEntityText(entityData map[string]interface{}) string {
	priorityFields := []string{"name", "title", "description", "content", "summary"}

	for _, field := range priorityFields {
		if value, exists := entityData[field]; exists {
			if str, ok := value.(string); ok && str != "" {
				return str
			}
		}
	}

	// If no priority fields found, concatenate all string values
	var parts []string
	for _, value := range entityData {
		if str, ok := value.(string); ok && str != "" {
			parts = append(parts, str)
		}
	}

	if len(parts) > 0 {
		return parts[0] // Use first non-empty string
	}

	return "default_entity_text"
}
