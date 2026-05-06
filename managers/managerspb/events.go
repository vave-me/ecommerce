package managerspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	ManagerAggregateChannel          = "middleman.managers.events.Manager"
	ConversationAggregateChannel     = "middleman.managers.events.Conversation"
	ManagerAuthorizedEvent           = "managersapi.ManagerAuthorized"
	ManagerCreatedEvent              = "managersapi.ManagerCreated"
	ManagerActivatedEvent            = "managersapi.ManagerActivated"
	ManagerDeactivatedEvent          = "managersapi.ManagerDeactivated"
	ManagerConfigurationUpdatedEvent = "managersapi.ManagerConfigurationUpdated"
	ManagerRequestProcessedEvent     = "managersapi.ManagerRequestProcessed"
	ManagerParticipatingToggledEvent = "managersapi.ManagerParticipatingToggled"
	ManagerRenamedEvent              = "managersapi.ManagerRenamed"
	ManagerLoggedInEvent             = "managersapi.ManagerLoggedIn"
	ManagerLoggedOutEvent            = "managersapi.ManagerLoggedOut"
	ConversationCreatedEvent         = "managersapi.ConversationCreated"
	MessageAddedEvent                = "managersapi.MessageAdded"
	ConversationContextUpdatedEvent  = "managersapi.ConversationContextUpdated"
	ConversationArchivedEvent        = "managersapi.ConversationArchivedEvent"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {

	//manager events
	if err := serde.Register(&ManagerCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&ManagerAuthorized{}); err != nil {
		return err
	}
	if err := serde.Register(&ManagerActivated{}); err != nil {
		return err
	}
	if err := serde.Register(&ManagerDeactivated{}); err != nil {
		return err
	}
	if err := serde.Register(&ManagerConfigurationUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&ManagerRequestProcessed{}); err != nil {
		return err
	}
	if err := serde.Register(&ManagerLoggedIn{}); err != nil {
		return err
	}
	if err := serde.Register(&ManagerLoggedOut{}); err != nil {
		return err
	}
	if err := serde.Register(&ManagerParticipationToggled{}); err != nil {
		return err
	}
	if err := serde.Register(&ManagerRenamed{}); err != nil {
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

func (*ManagerAuthorized) Key() string           { return ManagerAuthorizedEvent }
func (*ManagerCreated) Key() string              { return ManagerCreatedEvent }
func (*ManagerActivated) Key() string            { return ManagerActivatedEvent }
func (*ManagerDeactivated) Key() string          { return ManagerDeactivatedEvent }
func (*ManagerConfigurationUpdated) Key() string { return ManagerConfigurationUpdatedEvent }
func (*ManagerRequestProcessed) Key() string     { return ManagerRequestProcessedEvent }
func (*ManagerParticipationToggled) Key() string { return ManagerParticipatingToggledEvent }
func (*ManagerRenamed) Key() string              { return ManagerRenamedEvent }
func (*ManagerLoggedIn) Key() string             { return ManagerLoggedInEvent }
func (*ManagerLoggedOut) Key() string            { return ManagerLoggedOutEvent }
func (*ConversationCreated) Key() string         { return ConversationCreatedEvent }
func (*MessageAdded) Key() string                { return MessageAddedEvent }
func (*ConversationContextUpdated) Key() string  { return ConversationContextUpdatedEvent }
func (*ConversationArchived) Key() string        { return ConversationArchivedEvent }
