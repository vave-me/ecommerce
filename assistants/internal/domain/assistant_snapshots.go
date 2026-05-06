package domain

import (
	"time"
)

// AssistantV1 snapshot
type AssistantV1 struct {
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	Type         AssistantType         `json:"type"`
	Model        string                `json:"model"`
	UserID       string                `json:"user_id"`
	Capabilities []AssistantCapability `json:"capabilities"`
	Active       bool                  `json:"active"`
	Temperature  float64               `json:"temperature"`
	MaxTokens    int                   `json:"max_tokens"`
	SystemPrompt string                `json:"system_prompt"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// SnapshotName implements es.Snapshot
func (AssistantV1) SnapshotName() string { return "assistants.AssistantV1" }
