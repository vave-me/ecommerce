package assistantspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	AssistantAggregateChannel          = "middleman.assistants.events.Assistant"
	ConversationAggregateChannel       = "middleman.assistants.events.Conversation"
	AssistantAuthorizedEvent           = "assistantsapi.AssistantAuthorized"
	AssistantCreatedEvent              = "assistantsapi.AssistantCreated"
	AssistantActivatedEvent            = "assistantsapi.AssistantActivated"
	AssistantDeactivatedEvent          = "assistantsapi.AssistantDeactivated"
	AssistantConfigurationUpdatedEvent = "assistantsapi.AssistantConfigurationUpdated"
	AssistantRequestProcessedEvent     = "assistantsapi.AssistantRequestProcessed"
	AssistantParticipatingToggledEvent = "assistantsapi.AssistantParticipatingToggled"
	AssistantRenamedEvent              = "assistantsapi.AssistantRenamed"
	AssistantLoggedInEvent             = "assistantsapi.AssistantLoggedIn"
	AssistantLoggedOutEvent            = "assistantsapi.AssistantLoggedOut"
	ConversationCreatedEvent           = "assistantsapi.ConversationCreated"
	MessageAddedEvent                  = "assistantsapi.MessageAdded"
	ConversationContextUpdatedEvent    = "assistantsapi.ConversationContextUpdated"
	ConversationArchivedEvent          = "assistantsapi.ConversationArchivedEvent"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {

	//assistant events
	if err := serde.Register(&AssistantCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&AssistantAuthorized{}); err != nil {
		return err
	}
	if err := serde.Register(&AssistantActivated{}); err != nil {
		return err
	}
	if err := serde.Register(&AssistantDeactivated{}); err != nil {
		return err
	}
	if err := serde.Register(&AssistantConfigurationUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&AssistantRequestProcessed{}); err != nil {
		return err
	}
	if err := serde.Register(&AssistantLoggedIn{}); err != nil {
		return err
	}
	if err := serde.Register(&AssistantLoggedOut{}); err != nil {
		return err
	}
	if err := serde.Register(&AssistantParticipationToggled{}); err != nil {
		return err
	}
	if err := serde.Register(&AssistantRenamed{}); err != nil {
		return err
	}
	if err := serde.Register(&ConversationCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&MessageAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&ConversationContextUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&ConversationArchived{}); err != nil {
		return err
	}
	return nil
}

func (*AssistantAuthorized) Key() string           { return AssistantAuthorizedEvent }
func (*AssistantCreated) Key() string              { return AssistantCreatedEvent }
func (*AssistantActivated) Key() string            { return AssistantActivatedEvent }
func (*AssistantDeactivated) Key() string          { return AssistantDeactivatedEvent }
func (*AssistantConfigurationUpdated) Key() string { return AssistantConfigurationUpdatedEvent }
func (*AssistantRequestProcessed) Key() string     { return AssistantRequestProcessedEvent }
func (*AssistantParticipationToggled) Key() string { return AssistantParticipatingToggledEvent }
func (*AssistantRenamed) Key() string              { return AssistantRenamedEvent }
func (*AssistantLoggedIn) Key() string             { return AssistantLoggedInEvent }
func (*AssistantLoggedOut) Key() string            { return AssistantLoggedOutEvent }
func (*ConversationCreated) Key() string           { return ConversationCreatedEvent }
func (*MessageAdded) Key() string                  { return MessageAddedEvent }
func (*ConversationContextUpdated) Key() string    { return ConversationContextUpdatedEvent }
func (*ConversationArchived) Key() string          { return ConversationArchivedEvent }
