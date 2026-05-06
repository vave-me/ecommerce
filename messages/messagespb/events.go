package messagespb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	WebSocketChannel               = "messenger.>"
	SendMessageCommand             = "messenger.SendMessage"
	ConversationAggregateChannel   = "middleman.messages.events.Conversation"
	ConversationStartedEvent       = "messenger.ConversationStarted"
	ConversationActiveToggledEvent = "messenger.ConversationActiveToggled"
	ConversationRebrandedEvent     = "messenger.ConversationRebranded"
	MessageAggregateChannel        = "middleman.messages.events.Message"
	MessageSentEvent               = "messenger.MessageSent"
	MessageDeletedEvent            = "messenger.MessageDeleted"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Conversation events
	if err := serde.Register(&ConversationStarted{}); err != nil {
		return err
	}
	if err := serde.Register(&ConversationActiveToggled{}); err != nil {
		return err
	}

	if err := serde.Register(&MessageSent{}); err != nil {
		return err
	}

	if err := serde.Register(&MessageDeleted{}); err != nil {
		return err
	}

	if err := serde.Register(&SendMessage{}); err != nil {
		return err
	}
	return nil
}

func (*ConversationStarted) Key() string       { return ConversationStartedEvent }
func (*ConversationActiveToggled) Key() string { return ConversationActiveToggledEvent }

func (*MessageSent) Key() string    { return MessageSentEvent }
func (*MessageDeleted) Key() string { return MessageDeletedEvent }

func (*SendMessage) Key() string { return SendMessageCommand }
