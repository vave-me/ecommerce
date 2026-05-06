package domain

import (
	"time"
)

// ManagerV1 snapshot
type ManagerV1 struct {
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Type         ManagerType         `json:"type"`
	Model        string              `json:"model"`
	UserID       string              `json:"user_id"`
	Capabilities []ManagerCapability `json:"capabilities"`
	Active       bool                `json:"active"`
	Temperature  float64             `json:"temperature"`
	MaxTokens    int                 `json:"max_tokens"`
	SystemPrompt string              `json:"system_prompt"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// SnapshotName implements es.Snapshot
func (ManagerV1) SnapshotName() string { return "managers.ManagerV1" }
