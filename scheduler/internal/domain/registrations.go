package domain

import (
	"middleman/internal/es"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	SchedulerCreatedEvent = "scheduler.SchedulerCreated"
	ActionAddedEvent      = "scheduler.ActionAdded"
	ActionRemovedEvent    = "scheduler.ActionRemoved"
	ActionUpdatedEvent    = "scheduler.ActionUpdated"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	if err := serde.Register(Scheduler{}, func(v interface{}) error {
		scheduler := v.(*Scheduler)
		scheduler.Aggregate = es.NewAggregate("", SchedulerAggregate)
		return nil
	}); err != nil {
		return err
	}

	if err := serde.Register(Action{}, func(v interface{}) error {
		scheduler := v.(*Action)
		scheduler.Aggregate = es.NewAggregate("", ActionAggregate)
		return nil
	}); err != nil {
		return err
	}
	
	if err := serde.Register(Task{}, func(v interface{}) error {
		task := v.(*Task)
		task.Aggregate = es.NewAggregate("", TaskAggregate)
		return nil
	}); err != nil {
		return err
	}

	// interaction events
	if err := serde.Register(ActionAdded{}); err != nil {
		return err
	}
	if err := serde.Register(SchedulerCreated{}); err != nil {
		return err
	}
	if err := serde.Register(ActionRemoved{}); err != nil {
		return err
	}

	if err := serde.Register(ActionUpdated{}); err != nil {
		return err
	}
	
	// Task events
	if err := serde.Register(TaskScheduled{}); err != nil {
		return err
	}
	if err := serde.Register(TaskCancelled{}); err != nil {
		return err
	}
	if err := serde.Register(TaskUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(TaskExecutionStarted{}); err != nil {
		return err
	}
	if err := serde.Register(TaskExecutionCompleted{}); err != nil {
		return err
	}
	if err := serde.Register(TaskExecutionFailed{}); err != nil {
		return err
	}

	// interaction snapshots
	if err := serde.RegisterKey(ActionVi{}.SnapshotName(), ActionVi{}); err != nil {
		return err
	}
	// scheduler snapshots
	if err := serde.RegisterKey(SchedulerVi{}.SnapshotName(), SchedulerVi{}); err != nil {
		return err
	}
	// task snapshots
	if err := serde.RegisterKey(TaskSnapshot{}.SnapshotName(), TaskSnapshot{}); err != nil {
		return err
	}

	return nil
}

func (Action) Key() string    { return ActionAggregate }
func (Scheduler) Key() string { return SchedulerAggregate }
func (Task) Key() string      { return TaskAggregate }

func (SchedulerCreated) Key() string { return SchedulerCreatedEvent }
func (ActionAdded) Key() string      { return ActionAddedEvent }
func (ActionRemoved) Key() string    { return ActionRemovedEvent }
func (ActionUpdated) Key() string    { return ActionUpdatedEvent }
