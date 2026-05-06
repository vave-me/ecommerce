package ai

import (
	"context"
	"fmt"
)

// EmbeddingInterface defines the interface for generating embeddings
type EmbeddingInterface interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
	GetDimensions() int
	GetModel() string
}

// SimpleEmbeddingClient provides a basic implementation for compilation purposes
// This allows Qdrant + Redis functionality to work while OpenAI client issues are resolved
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

// GenerateEmbedding generates a simple embedding for compilation purposes
func (s *SimpleEmbeddingClient) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	// Create a deterministic embedding based on text content
	// This is simplified for compilation - real implementation would use OpenAI API
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
