package services

import (
	"context"
	"middleman/internal/ai"
)

// AIClientProvider provides access to different AI clients with production features
type AIClientProvider interface {
	// GetClient returns a specific AI client by provider name
	GetClient(ctx context.Context, provider string) (ai.EnhancedAIService, error)
	
	// GetDefaultClient returns the default AI client, trying fallbacks if necessary
	GetDefaultClient(ctx context.Context) (ai.EnhancedAIService, error)
	
	// GetHealthyProvider returns the first healthy provider
	GetHealthyProvider(ctx context.Context) (ai.EnhancedAIService, string, error)
	
	// GetOptimalProvider returns the best provider for a specific request type
	GetOptimalProvider(ctx context.Context, requestType string) (ai.EnhancedAIService, string, error)
}