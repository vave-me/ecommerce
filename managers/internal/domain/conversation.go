package domain

import (
	"fmt"
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"github.com/stackus/errors"
)

const ConversationAggregate = "managers.Conversation"

var (
	ErrConversationIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the conversation id cannot be blank")
	ErrConversationNotFound        = errors.Wrap(errors.ErrNotFound, "conversation not found")
	ErrMessageCannotBeBlank        = errors.Wrap(errors.ErrBadRequest, "message cannot be blank")
	ErrUnauthorized                = errors.Wrap(errors.ErrUnauthorized, "unauthorized access")
)

// Conversation aggregate
type Conversation struct {
	es.Aggregate
	UserID    string
	ManagerID string
	Messages  []ConversationMessage
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
	Context   map[string]interface{}
}

func (Conversation) Key() string { return ConversationAggregate }

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Conversation)(nil)

func NewConversation(id string) *Conversation {
	return &Conversation{
		Aggregate: es.NewAggregate(id, ConversationAggregate),
	}
}

// CreateConversation creates a new conversation
func (c *Conversation) CreateConversation(id, userID, managerID string, initialContext map[string]interface{}) (ddd.Event, error) {
	fmt.Printf("[CreateConversation] Called with: id=%s, userID=%s, managerID=%s\n", id, userID, managerID)

	if id == "" {
		return nil, ErrConversationIDCannotBeBlank
	}

	if c.UserID != "" {
		return nil, errors.Wrap(errors.ErrBadRequest, "conversation already exists")
	}

	eventPayload := &ConversationCreated{
		ConversationID: id,
		UserID:         userID,
		ManagerID:      managerID,
		CreatedAt:      time.Now(),
		Context:        initialContext,
	}

	fmt.Printf("[CreateConversation] Event payload: ConversationID=%s, UserID=%s, ManagerID=%s\n",
		eventPayload.ConversationID, eventPayload.UserID, eventPayload.ManagerID)

	c.AddEvent(ConversationCreatedEvent, eventPayload)

	// Return the event with the actual event payload
	return ddd.NewEvent(ConversationCreatedEvent, eventPayload), nil
}

// AddMessage adds a message to the conversation
func (c *Conversation) AddMessage(conversationsID, managerID, messageID string, role MessageRole, content string, metadata map[string]interface{}, actionsTaken []ManagerAction) (ddd.Event, error) {
	if content == "" {
		return nil, ErrMessageCannotBeBlank
	}

	timestamp := time.Now()

	eventPayload := &MessageAdded{
		ConversationID: conversationsID,
		ManagerID:      managerID,
		ID:             messageID,
		Role:           role,
		Content:        content,
		Metadata:       metadata,
		ActionsTaken:   actionsTaken,
		Timestamp:      timestamp,
	}

	c.AddEvent(MessageAddedEvent, eventPayload)

	// Return the event with the actual event payload
	return ddd.NewEvent(MessageAddedEvent, eventPayload), nil
}

// UpdateContext updates conversation context
func (c *Conversation) UpdateContext(newContext map[string]interface{}) (ddd.Event, error) {
	eventPayload := &ConversationContextUpdated{
		ConversationID: c.ID(),
		Context:        newContext,
		UpdatedAt:      time.Now(),
	}

	c.AddEvent(ConversationContextUpdatedEvent, eventPayload)

	// Return the event with the actual event payload
	return ddd.NewEvent(ConversationContextUpdatedEvent, eventPayload), nil
}

// ArchiveConversation archives the conversation
func (c *Conversation) ArchiveConversation() (ddd.Event, error) {
	if !c.Active {
		return nil, errors.Wrap(errors.ErrBadRequest, "conversation is already archived")
	}

	eventPayload := &ConversationArchived{
		ConversationID: c.ID(),
		ArchivedAt:     time.Now(),
	}

	c.AddEvent(ConversationArchivedEvent, eventPayload)

	// Return the event with the actual event payload
	return ddd.NewEvent(ConversationArchivedEvent, eventPayload), nil
}

// GetHistory returns the conversation message history
func (c *Conversation) GetHistory() []ConversationMessage {
	return c.Messages
}

// GetMessageCount returns the number of messages in the conversation
func (c *Conversation) GetMessageCount() int {
	return len(c.Messages)
}

// GetLastMessage returns the last message in the conversation
func (c *Conversation) GetLastMessage() *ConversationMessage {
	if len(c.Messages) == 0 {
		return nil
	}
	return &c.Messages[len(c.Messages)-1]
}

// IsActive returns whether the conversation is active
func (c *Conversation) IsActive() bool {
	return c.Active
}

// GetContext returns the conversation context
func (c *Conversation) GetContext() map[string]interface{} {
	if c.Context == nil {
		return make(map[string]interface{})
	}
	return c.Context
}

// UpdateMetadata updates the conversation metadata
func (c *Conversation) UpdateMetadata(metadata map[string]interface{}) (ddd.Event, error) {
	eventPayload := &ConversationMetadataUpdated{
		ConversationID: c.ID(),
		Metadata:       metadata,
		UpdatedAt:      time.Now(),
	}

	c.AddEvent(ConversationMetadataUpdatedEvent, eventPayload)
	
	// Return the event with the actual event payload
	return ddd.NewEvent(ConversationMetadataUpdatedEvent, eventPayload), nil
}

// Delete marks the conversation as deleted
func (c *Conversation) Delete() (ddd.Event, error) {
	if !c.Active {
		return nil, errors.Wrap(errors.ErrBadRequest, "conversation is already inactive")
	}

	eventPayload := &ConversationDeleted{
		ConversationID: c.ID(),
		DeletedAt:      time.Now(),
	}

	c.AddEvent(ConversationDeletedEvent, eventPayload)
	
	// Return the event with the actual event payload
	return ddd.NewEvent(ConversationDeletedEvent, eventPayload), nil
}

// GetUserID returns the user ID of the conversation
func (c *Conversation) GetUserID() string {
	return c.UserID
}

// ApplyEvent applies an event to the conversation aggregate
func (c *Conversation) ApplyEvent(event ddd.Event) error {

	switch payload := event.Payload().(type) {
	case *ConversationCreated:

		c.UserID = payload.UserID
		c.ManagerID = payload.ManagerID
		c.CreatedAt = payload.CreatedAt
		c.UpdatedAt = payload.CreatedAt
		c.Active = true
		c.Context = payload.Context
		if c.Context == nil {
			c.Context = make(map[string]interface{})
		}
		c.Messages = []ConversationMessage{} // Initialize empty messages slice

	case *MessageAdded:
		message := ConversationMessage{
			ID:             payload.ID,
			Role:           payload.Role,
			ManagerID:      payload.ManagerID,
			ConversationID: payload.ConversationID,
			Content:        payload.Content,
			Timestamp:      payload.Timestamp,
			Metadata:       payload.Metadata,
			ActionsTaken:   payload.ActionsTaken,
		}
		c.Messages = append(c.Messages, message)
		c.UpdatedAt = payload.Timestamp

	case *ConversationContextUpdated:
		c.Context = payload.Context
		if c.Context == nil {
			c.Context = make(map[string]interface{})
		}
		c.UpdatedAt = payload.UpdatedAt

	case *ConversationArchived:

		c.Active = false
		c.UpdatedAt = payload.ArchivedAt

	case *ConversationMetadataUpdated:
		if c.Context == nil {
			c.Context = make(map[string]interface{})
		}
		// Merge metadata
		for k, v := range payload.Metadata {
			c.Context[k] = v
		}
		c.UpdatedAt = payload.UpdatedAt

	case *ConversationDeleted:
		c.Active = false
		c.UpdatedAt = payload.DeletedAt

	default:

		return fmt.Errorf("unknown event payload type: %T for event %s", payload, event.EventName())
	}

	return nil
}

// ApplySnapshot applies a snapshot to the conversation aggregate
func (c *Conversation) ApplySnapshot(snapshot es.Snapshot) error {
	if conversationSnapshot, ok := snapshot.(ConversationV1); ok {
		c.UserID = conversationSnapshot.UserID
		c.ManagerID = conversationSnapshot.ManagerID
		c.Messages = conversationSnapshot.Messages
		c.CreatedAt = conversationSnapshot.CreatedAt
		c.UpdatedAt = conversationSnapshot.UpdatedAt
		c.Active = conversationSnapshot.Active
		c.Context = conversationSnapshot.Context
		if c.Context == nil {
			c.Context = make(map[string]interface{})
		}
		if c.Messages == nil {
			c.Messages = []ConversationMessage{}
		}
	}
	return nil
}

// ToSnapshot creates a snapshot of the conversation aggregate
func (c *Conversation) ToSnapshot() es.Snapshot {
	return ConversationV1{
		ID:        c.ID(),
		UserID:    c.UserID,
		ManagerID: c.ManagerID,
		Messages:  c.Messages,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Active:    c.Active,
		Context:   c.Context,
	}
}

// ConversationStats represents conversation statistics for a user
type ConversationStats struct {
	TotalConversations         int64     `json:"total_conversations"`
	ActiveConversations        int64     `json:"active_conversations"`
	TotalMessages              int64     `json:"total_messages"`
	MessagesToday              int64     `json:"messages_today"`
	MessagesThisWeek           int64     `json:"messages_this_week"`
	MessagesThisMonth          int64     `json:"messages_this_month"`
	FirstConversationAt        time.Time `json:"first_conversation_at"`
	LastConversationAt         time.Time `json:"last_conversation_at"`
	AvgMessagesPerConversation float64   `json:"avg_messages_per_conversation"`
	MostUsedManagerID          string    `json:"most_used_manager_id"`
}
