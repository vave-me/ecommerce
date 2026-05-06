package domain

import (
	"middleman/internal/es"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	MessageSentEvent          = "messages.MessageSent"
	MessageReadEvent          = "messages.MessageRead"
	MessageDeletedEvent       = "messages.MessageDeleted"
	MessageReceivedEvent      = "messages.MessageReceived"
	ConversationStartedEvent  = "messages.ConversationStarted"
	ConversationArchivedEvent = "messages.ConversationArchived"
	ConversationDeletedEvent  = "messages.ConversationDeleted"
	ConversationMuttedEvent   = "messages.ConversationMutted"
	ConversationUnmuttedEvent = "messages.ConversationUnmutted"
	ConversationClosedEvent   = "messages.ConversationClosed"
	ConversationReopenedEvent = "messages.ConversationReopened"
)

// TODO check this and uncomment to refactor
func Registrations(reg registry.Registry) (err error) {

	serde := serdes.NewJsonSerde(reg)
	//serde := serdes.NewProtoSerde(reg)

	// Store
	if err = serde.Register(Conversation{}, func(v any) error {
		conversation := v.(*Conversation)
		conversation.Aggregate = es.NewAggregate("", ConversationAggregate)
		return nil
	}); err != nil {
		return
	}

	// conversation events
	if err = serde.Register(ConversationStarted{}); err != nil {
		return
	}

	if err = serde.RegisterKey(ConversationV1{}.SnapshotName(), ConversationV1{}); err != nil {
		return
	}
	// Message
	if err = serde.Register(Message{}, func(v any) error {
		message := v.(*Message)
		message.Aggregate = es.NewAggregate("", MessageAggregate)
		return nil
	}); err != nil {
		return
	}

	if err = serde.Register(MessageSent{}); err != nil {
		return
	}
	// reply snapshots
	if err = serde.RegisterKey(MessageV1{}.SnapshotName(), MessageV1{}); err != nil {
		return
	}

	return
}

func (MessageRead) Key() string { return MessageReadEvent }

func (MessageDeleted) Key() string { return MessageDeletedEvent }

func (MessageReceived) Key() string { return MessageReceivedEvent }

func (ConversationStarted) Key() string { return ConversationStartedEvent }

func (ConversationArchived) Key() string { return ConversationArchivedEvent }
func (ConversationDeleted) Key() string  { return ConversationDeletedEvent }

func (ConversationMutted) Key() string { return ConversationMuttedEvent }

func (ConversationUnmutted) Key() string { return ConversationUnmuttedEvent }

func (ConversationReopened) Key() string { return ConversationReopenedEvent }
func (ConversationClosed) Key() string   { return ConversationClosedEvent }
func (MessageSent) Key() string          { return MessageSentEvent }
