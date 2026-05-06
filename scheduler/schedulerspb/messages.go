package schedulerspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	SchedulerAggregateChannel = "middleman.scheduler.events.Scheduler"
	ActionAggregateChannel    = "middleman.scheduler.events.Action"
	SchedulerCreatedEvent     = "schedulerapi.SchedulerCreated"
	ActionAddedEvent          = "schedulerapi.ActionAdded"
	ActionUpdatedEvent        = "schedulerapi.ActionUpdated"
	ActionRemovedEvent        = "schedulerapi.ActionRemoved"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {

	// Action events
	if err := serde.Register(&SchedulerCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&ActionAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&ActionUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&ActionRemoved{}); err != nil {
		return err
	}

	return nil

}

func (*SchedulerCreated) Key() string { return SchedulerCreatedEvent }
func (*ActionAdded) Key() string      { return ActionAddedEvent }
func (*ActionUpdated) Key() string    { return ActionUpdatedEvent }
func (*ActionRemoved) Key() string    { return ActionRemovedEvent }
