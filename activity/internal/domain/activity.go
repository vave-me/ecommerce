package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const ActivityAggregate = "activity.Activity"

var (
	ErrUserIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the user ID cannot be blank")
)

type Activity struct {
	es.Aggregate
	UserID string
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Activity)(nil)

func NewActivity(id string) *Activity {
	return &Activity{
		Aggregate: es.NewAggregate(id, ActivityAggregate),
	}
}

func (i *Activity) InitActivity(userID string) (ddd.Event, error) {

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	i.AddEvent(ActivityCreatedEvent, &ActivityCreated{
		UserID: userID,
	})

	return ddd.NewEvent(ActivityCreatedEvent, i), nil
}
func (b *Activity) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *ActivityCreated:
		b.UserID = payload.UserID

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", b, event.EventName(), payload)
	}

	return nil
}

func (b *Activity) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *ActivityVi:
		b.UserID = ss.UserID
	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", b, snapshot)
	}

	return nil
}

func (b *Activity) ToSnapshot() es.Snapshot {
	return &ActivityVi{
		UserID: b.UserID,
	}
}
