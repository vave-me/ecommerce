package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const MessageAggregate = "messages.Message"

var (
	ErrMessageNameIsBlank     = errors.Wrap(errors.ErrBadRequest, "the item name cannot be blank")
	ErrMessagePriceIsNegative = errors.Wrap(errors.ErrBadRequest, "the item price cannot be negative")
	ErrNotAPriceIncrease      = errors.Wrap(errors.ErrBadRequest, "the price change would be a decrease")
	ErrNotAPriceDecrease      = errors.Wrap(errors.ErrBadRequest, "the price change would be an increase")
)

type Message struct {
	es.Aggregate
	MessageID      string
	ConversationID string
	SenderID       string
	RecipientID    string
	ItemID         string
	Body           string
	IsRead         bool
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Message)(nil)

func NewMessage(id string) *Message {
	return &Message{
		Aggregate: es.NewAggregate(id, MessageAggregate),
	}
}

// Key implements registry.Registerable
func (Message) Key() string { return MessageAggregate }

func (m *Message) InitMessage(id, conversationID, senderID, recipientID, itemID, body string, isRead bool) (ddd.Event, error) {

	m.AddEvent(MessageSentEvent, &MessageSent{
		ID:             id,
		ConversationID: conversationID,
		SenderID:       senderID,
		RecipientID:    recipientID,
		ItemID:         itemID,
		Body:           body,
		IsRead:         isRead,
	})

	return ddd.NewEvent(MessageSentEvent, m), nil
}

func (m *Message) Delete() (ddd.Event, error) {
	m.AddEvent(MessageDeletedEvent, &MessageDeleted{})

	return ddd.NewEvent(MessageDeletedEvent, m), nil
}

//// TODO MarkAsRead updates the message as read and records MessageReceivedEvent.
//func (m *Message) MarkAsRead(receiverID string) {
//	m.IsRead = true
//	m.recordEvent(MessageReceivedEvent{MessageID: m.ID, ReceiverID: receiverID})
//}

func (m *Message) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *MessageSent:
		m.MessageID = payload.ID
		m.ConversationID = payload.ConversationID
		m.SenderID = payload.SenderID
		m.RecipientID = payload.RecipientID
		m.ItemID = payload.ItemID
		m.Body = payload.Body
		m.IsRead = payload.IsRead

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", m, event.EventName(), payload)
	}

	return nil
}

func (m *Message) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *MessageV1:
		m.MessageID = ss.ID
		m.ConversationID = ss.ConversationID
		m.SenderID = ss.SenderID
		m.RecipientID = ss.RecipientID
		m.Body = ss.Body
		m.IsRead = ss.IsRead

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", m, snapshot)
	}

	return nil
}

func (m Message) ToSnapshot() es.Snapshot {
	return MessageV1{
		ID:             m.MessageID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		RecipientID:    m.RecipientID,
		ItemID:         m.ItemID,
		Body:           m.Body,
		IsRead:         m.IsRead,
	}
}
