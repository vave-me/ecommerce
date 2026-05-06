package application

import (
	"context"
	"middleman/managers/internal/domain"
)

// AIModelSelector selects the appropriate AI model for a task
type AIModelSelector struct {
	// Add configuration
}

// NewAIModelSelector creates a new AI model selector
func NewAIModelSelector() *AIModelSelector {
	return &AIModelSelector{}
}

// SelectModel selects the best model for a given manager and task
func (s *AIModelSelector) SelectModel(ctx context.Context, manager *domain.Manager, taskType string) string {
	// For now, return the manager's configured model
	if manager.Model != "" {
		return manager.Model
	}
	
	// Default model selection based on task type
	switch taskType {
	case "code_generation":
		return "gpt-4"
	case "simple_chat":
		return "gpt-3.5-turbo"
	default:
		return "gpt-4"
	}
}

// GetModelCapabilities returns capabilities of a model
func (s *AIModelSelector) GetModelCapabilities(model string) map[string]bool {
	capabilities := map[string]bool{
		"chat": true,
		"completion": true,
		"embeddings": false,
	}
	
	// Add model-specific capabilities
	switch model {
	case "gpt-4", "gpt-4-turbo":
		capabilities["vision"] = true
		capabilities["function_calling"] = true
		capabilities["long_context"] = true
	case "gpt-3.5-turbo":
		capabilities["function_calling"] = true
		capabilities["fast_response"] = true
	}
	
	return capabilities
}