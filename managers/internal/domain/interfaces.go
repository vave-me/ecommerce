package domain

// SystemPromptProvider provides system prompts for managers
type SystemPromptProvider interface {
	GetCompleteSystemPrompt() string
}
