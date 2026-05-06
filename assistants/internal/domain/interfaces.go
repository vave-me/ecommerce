package domain

// SystemPromptProvider provides system prompts for assistants
type SystemPromptProvider interface {
	GetCompleteSystemPrompt() string
}