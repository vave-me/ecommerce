package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const SchedulerAggregate = "scheduler.Scheduler"

var (
	ErrUserIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the user ID cannot be blank")
)

type Scheduler struct {
	es.Aggregate
	UserID string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Scheduler)(nil)

func NewScheduler(id string) *Scheduler {
	return &Scheduler{
		Aggregate: es.NewAggregate(id, SchedulerAggregate),
	}
}

func (i *Scheduler) InitScheduler(userID string) (ddd.Event, error) {

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	i.AddEvent(SchedulerCreatedEvent, &SchedulerCreated{
		UserID: userID,
	})

	return ddd.NewEvent(SchedulerCreatedEvent, i), nil
}
func (b *Scheduler) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *SchedulerCreated:
		b.UserID = payload.UserID

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", b, event.EventName(), payload)
	}

	return nil
}

func (b *Scheduler) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *SchedulerVi:
		b.UserID = ss.UserID
	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", b, snapshot)
	}

	return nil
}

func (b *Scheduler) ToSnapshot() es.Snapshot {
	return &SchedulerVi{
		UserID: b.UserID,
	}
}
