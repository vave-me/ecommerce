package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const ConversationAggregate = "messages.Conversation"

// Conversation represents a conversation aggregate with event sourcing.
type Conversation struct {
	es.Aggregate
	SenderID    string
	RecipientID string
	ItemID      string
	Messages    []Message
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Conversation)(nil)

// NewConversation creates a new Conversation aggregate.
func NewConversation(id string) *Conversation {
	return &Conversation{
		Aggregate: es.NewAggregate(id, ConversationAggregate),
	}
}

// Key implements registry.Registerable
func (Conversation) Key() string { return ConversationAggregate }

func (c *Conversation) InitConversation(id, senderID, recipientID, itemID string) (ddd.Event, error) {
	c.AddEvent(ConversationStartedEvent, &ConversationStarted{
		ConversationID: id,
		SenderID:       senderID,
		RecipientID:    recipientID,
		ItemID:         itemID,
	})
	return ddd.NewEvent(ConversationStartedEvent, c), nil
}
func (c *Conversation) Delete() (ddd.Event, error) {
	c.AddEvent(ConversationDeletedEvent, &ConversationDeleted{})

	return ddd.NewEvent(ConversationDeletedEvent, c), nil
}

//// isUserInConversation checks if a user is part of the conversation.
//func (c *Conversation) isUserInConversation(userID string) bool {
//	for _, id := range c.Users {
//		if id == userID {
//			return true
//		}
//	}
//	return false
//}

// ApplyEvent implements es.EventApplier
func (c *Conversation) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *ConversationStarted:
		c.SenderID = payload.SenderID
		c.RecipientID = payload.RecipientID
		c.ItemID = payload.ItemID
	case *ConversationDeleted:
		//nothing
	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", c, event.EventName(), payload)
	}

	return nil

}

// ApplySnapshot implements es.Snapshotter
func (c *Conversation) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ConversationV1:
		c.SenderID = ss.SenderID
		c.RecipientID = ss.RecipientID
		c.ItemID = ss.ItemID

	default:
		return errors.Wrap(errors.ErrBadRequest, "the basket has no items")
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (s Conversation) ToSnapshot() es.Snapshot {
	return ConversationV1{}
}
