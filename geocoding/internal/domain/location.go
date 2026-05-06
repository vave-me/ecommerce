package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const LocationAggregate = "geocoding.Location"

type Location struct {
	es.Aggregate
	ProductID string
	Latitude  float64
	Longitude float64
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Location)(nil)

func NewLocation(id string) *Location {
	return &Location{
		Aggregate: es.NewAggregate(id, LocationAggregate),
	}
}
func (a *Location) InitLocation(id, mediaID string) (ddd.Event, error) {

	a.AddEvent(LocationAddedEvent, &LocationAdded{})

	return ddd.NewEvent(LocationAddedEvent, a), nil
}
func (Location) Key() string { return LocationAggregate }

// ApplyEvent implements es.EventApplier
func (l *Location) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *LocationAdded:

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", l, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (l *Location) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *LocationV1:
		l.ProductID = ss.ProductID
		l.Latitude = ss.Latitude
		l.Longitude = ss.Longitude

	default:
		return errors.ErrInternal.Msgf(
			"%T received the unexpected snapshot %T", l, snapshot)
	}
	return nil
}

func (l Location) ToSnapshot() es.Snapshot {
	return LocationV1{}
}
