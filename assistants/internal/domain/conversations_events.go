package domain

import "time"

const (
	ConversationCreatedEvent        = "assistants.ConversationCreated"
	ConversationContextUpdatedEvent = "assistants.ConversationContextUpdated"
	MessageAddedEvent               = "assistants.MessageAdded"
	ConversationArchivedEvent       = "assistants.ConversationArchived"
)

type ConversationCreated struct {
	ConversationID string                 `json:"conversation_id"`
	UserID         string                 `json:"user_id"`
	AssistantID    string                 `json:"assistant_id"`
	CreatedAt      time.Time              `json:"created_at"`
	Context        map[string]interface{} `json:"context,omitempty"`
}

// Key implements registry.Registerable
func (ConversationCreated) Key() string { return ConversationCreatedEvent }

type MessageAdded struct {
	ConversationID string                 `json:"conversation_id"`
	AssistantID    string                 `json:"assistant_id"`
	ID             string                 `json:"id"`
	Role           MessageRole            `json:"role"`
	Content        string                 `json:"content"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	ActionsTaken   []AssistantAction      `json:"actions_taken,omitempty"`
	Timestamp      time.Time              `json:"timestamp"`
}

func (MessageAdded) Key() string { return MessageAddedEvent }

type ConversationContextUpdated struct {
	ConversationID string                 `json:"conversation_id"`
	Context        map[string]interface{} `json:"context"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

func (ConversationContextUpdated) Key() string { return ConversationContextUpdatedEvent }

type ConversationArchived struct {
	ConversationID string    `json:"conversation_id"`
	ArchivedAt     time.Time `json:"archived_at"`
}

func (ConversationArchived) Key() string { return ConversationArchivedEvent }
