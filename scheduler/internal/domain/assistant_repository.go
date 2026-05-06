package domain

import "context"

// AssistantRepository provides access to the LLM/Assistant service
type AssistantRepository interface {
	// ProcessUserInput sends a natural language task to the assistant service for processing
	ProcessUserInput(ctx context.Context, task string, context map[string]string) (string, error)
}