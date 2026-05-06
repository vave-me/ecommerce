package domain

import "time"

const (
	ConversationCreatedEvent        = "managers.ConversationCreated"
	ConversationContextUpdatedEvent = "managers.ConversationContextUpdated"
	MessageAddedEvent               = "managers.MessageAdded"
	ConversationArchivedEvent       = "managers.ConversationArchived"
	ConversationMetadataUpdatedEvent = "managers.ConversationMetadataUpdated"
	ConversationDeletedEvent        = "managers.ConversationDeleted"
)

type ConversationCreated struct {
	ConversationID string                 `json:"conversation_id"`
	UserID         string                 `json:"user_id"`
	ManagerID      string                 `json:"manager_id"`
	CreatedAt      time.Time              `json:"created_at"`
	Context        map[string]interface{} `json:"context,omitempty"`
}

// Key implements registry.Registerable
func (ConversationCreated) Key() string { return ConversationCreatedEvent }

type MessageAdded struct {
	ConversationID string                 `json:"conversation_id"`
	ManagerID      string                 `json:"manager_id"`
	ID             string                 `json:"id"`
	Role           MessageRole            `json:"role"`
	Content        string                 `json:"content"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	ActionsTaken   []ManagerAction        `json:"actions_taken,omitempty"`
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

type ConversationMetadataUpdated struct {
	ConversationID string                 `json:"conversation_id"`
	Metadata       map[string]interface{} `json:"metadata"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

func (ConversationMetadataUpdated) Key() string { return ConversationMetadataUpdatedEvent }

type ConversationDeleted struct {
	ConversationID string    `json:"conversation_id"`
	DeletedAt      time.Time `json:"deleted_at"`
}

func (ConversationDeleted) Key() string { return ConversationDeletedEvent }
