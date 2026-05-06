package domain

import (
	"middleman/internal/es"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	ActivityCreatedEvent    = "activity.ActivityCreated"
	InteractionAddedEvent   = "activity.InteractionAdded"
	InteractionRemovedEvent = "activity.InteractionRemoved"
	InteractionUpdatedEvent = "activity.InteractionUpdated"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	if err := serde.Register(Activity{}, func(v interface{}) error {
		activity := v.(*Activity)
		activity.Aggregate = es.NewAggregate("", ActivityAggregate)
		return nil
	}); err != nil {
		return err
	}

	if err := serde.Register(Interaction{}, func(v interface{}) error {
		activity := v.(*Interaction)
		activity.Aggregate = es.NewAggregate("", InteractionAggregate)
		return nil
	}); err != nil {
		return err
	}

	// interaction events
	if err := serde.Register(InteractionAdded{}); err != nil {
		return err
	}
	if err := serde.Register(ActivityCreated{}); err != nil {
		return err
	}
	if err := serde.Register(InteractionRemoved{}); err != nil {
		return err
	}

	if err := serde.Register(InteractionUpdated{}); err != nil {
		return err
	}

	// interaction snapshots
	if err := serde.RegisterKey(InteractionVi{}.SnapshotName(), InteractionVi{}); err != nil {
		return err
	}
	// activity snapshots
	if err := serde.RegisterKey(ActivityVi{}.SnapshotName(), ActivityVi{}); err != nil {
		return err
	}

	return nil
}

func (Interaction) Key() string { return InteractionAggregate }
func (Activity) Key() string    { return ActivityAggregate }

func (ActivityCreated) Key() string    { return ActivityCreatedEvent }
func (InteractionAdded) Key() string   { return InteractionAddedEvent }
func (InteractionRemoved) Key() string { return InteractionRemovedEvent }
func (InteractionUpdated) Key() string { return InteractionUpdatedEvent }
