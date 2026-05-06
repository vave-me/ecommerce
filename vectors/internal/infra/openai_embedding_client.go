package infra

import (
	"context"
	"fmt"
	ai2 "middleman/internal/ai"
	"middleman/vectors/internal/ports"
	"time"
)

// OpenAIEmbeddingClient adapts the AI OpenAI client to the EmbeddingClient interface
type OpenAIEmbeddingClient struct {
	client *ai2.OpenAIClient
	config EmbeddingClientConfig
}

// EmbeddingClientConfig holds configuration for the embedding client
type EmbeddingClientConfig struct {
	Model         string
	Dimensions    int
	PromptEnabled bool
	MaxRetries    int
	Timeout       time.Duration
}

// NewOpenAIEmbeddingClient creates a new OpenAI embedding client adapter
func NewOpenAIEmbeddingClient(client *ai2.OpenAIClient, config EmbeddingClientConfig) *OpenAIEmbeddingClient {
	if config.Model == "" {
		config.Model = "text-embedding-3-small"
	}
	if config.Dimensions == 0 {
		config.Dimensions = 1536
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &OpenAIEmbeddingClient{
		client: client,
		config: config,
	}
}

// GenerateEmbedding generates an embedding for the given text
func (c *OpenAIEmbeddingClient) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("empty text provided")
	}

	return c.client.GenerateEmbedding(ctx, text)
}

// GenerateBatchEmbeddings generates embeddings for multiple texts
func (c *OpenAIEmbeddingClient) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts provided")
	}

	return c.client.GenerateBatchEmbeddings(ctx, texts)
}

// GenerateEntityEmbedding generates an embedding for entity data
func (c *OpenAIEmbeddingClient) GenerateEntityEmbedding(ctx context.Context, entityData map[string]interface{}) ([]float32, error) {
	if len(entityData) == 0 {
		return nil, fmt.Errorf("empty entity data provided")
	}

	return c.client.GenerateEntityEmbedding(ctx, entityData)
}

// GenerateEmbeddingWithPrompt generates an embedding using a custom prompt
func (c *OpenAIEmbeddingClient) GenerateEmbeddingWithPrompt(ctx context.Context, text string, prompt string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("empty text provided")
	}

	return c.client.GenerateEmbeddingWithPrompt(ctx, text, prompt)
}

// GenerateOptimizedEmbedding generates an optimized embedding for specific use cases
func (c *OpenAIEmbeddingClient) GenerateOptimizedEmbedding(ctx context.Context, entityType string, entityData map[string]interface{}, optimization string) ([]float32, error) {
	if len(entityData) == 0 {
		return nil, fmt.Errorf("empty entity data provided")
	}

	return c.client.GenerateOptimizedEmbedding(ctx, entityType, entityData, optimization)
}

// GetDimensions returns the embedding dimensions
func (c *OpenAIEmbeddingClient) GetDimensions() int {
	return c.client.GetEmbeddingDimensions()
}

// GetModel returns the embedding model name
func (c *OpenAIEmbeddingClient) GetModel() string {
	return c.client.GetEmbeddingModel()
}

// IsPromptEnabled returns whether prompt enhancement is enabled
func (c *OpenAIEmbeddingClient) IsPromptEnabled() bool {
	return c.client.IsEmbeddingPromptEnabled()
}

// HealthCheck performs a health check on the embedding client
func (c *OpenAIEmbeddingClient) HealthCheck(ctx context.Context) error {
	return c.client.HealthCheck(ctx)
}

// GetProviderInfo returns information about the embedding provider
func (c *OpenAIEmbeddingClient) GetProviderInfo() ports.EmbeddingProviderInfo {
	info := c.client.GetProviderInfo()
	return ports.EmbeddingProviderInfo{
		Provider:     info.Provider,
		Model:        c.GetModel(),
		Dimensions:   c.GetDimensions(),
		MaxTokens:    8192, // Max input tokens for text-embedding models
		RateLimit:    1000, // Requests per minute
		CostPerToken: 0.02, // $0.02 per 1M tokens for text-embedding-3-small
		ConnectedAt:  time.Now(),
	}
}
