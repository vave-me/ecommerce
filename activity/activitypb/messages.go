package activitypb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	ActivityAggregateChannel    = "middleman.activity.events.Activity"
	InteractionAggregateChannel = "middleman.activity.events.Interaction"
	ActivityCreatedEvent        = "activityapi.ActivityCreated"
	InteractionAddedEvent       = "activityapi.InteractionAdded"
	InteractionUpdatedEvent     = "activityapi.InteractionUpdated"
	InteractionRemovedEvent     = "activityapi.InteractionRemoved"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {

	// Interaction events
	if err := serde.Register(&ActivityCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&InteractionAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&InteractionUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&InteractionRemoved{}); err != nil {
		return err
	}

	return nil

}

func (*ActivityCreated) Key() string    { return ActivityCreatedEvent }
func (*InteractionAdded) Key() string   { return InteractionAddedEvent }
func (*InteractionUpdated) Key() string { return InteractionUpdatedEvent }
func (*InteractionRemoved) Key() string { return InteractionRemovedEvent }
