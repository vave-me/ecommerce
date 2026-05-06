package domain

import (
	"time"
)

// ConversatonV1 snapshot
type ConversationV1 struct {
	ID          string
	UserID      string
	AssistantID string
	Messages    []ConversationMessage
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Active      bool
	Context     map[string]interface{}
}

// SnapshotName implements es.Snapshot
func (ConversationV1) SnapshotName() string { return "assistants.ConversationV1" }
