package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const AddressAggregate = "geocoding.Address"

type Address struct {
	es.Aggregate
	FormattedAddress string
	UserID           string
	StreetName       string
	StreetNumber     int
	StreetSuffix     string
	City             string
	State            string
	Zip              string
	Latitude         float64
	Longitude        float64
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Address)(nil)

func NewAddress(id string) *Address {
	return &Address{
		Aggregate: es.NewAggregate(id, AddressAggregate),
	}
}
func (a *Address) InitAddress(id, address string, lat, lng float64) (ddd.Event, error) {

	a.AddEvent(AddressCreatedEvent, &AddressCreated{
		Address:   address,
		Latitude:  lat,
		Longitude: lng,
	})

	return ddd.NewEvent(AddressCreatedEvent, a), nil
}

func (a *Address) BatchDecodeAddress(id, address string) (ddd.Event, error) {

	a.AddEvent(AddressBatchDecodedEvent, &AddressBatchDecoded{
		Address: address,
	})

	return ddd.NewEvent(AddressCreatedEvent, a), nil
}
func (Address) Key() string { return AddressAggregate }

// ApplyEvent implements es.EventApplier
func (m *Address) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *AddressCreated:

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", m, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (a *Address) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *AddressV1:
		a.FormattedAddress = ss.FormattedAddress
		a.UserID = ss.UserID
		a.StreetName = ss.StreetName
		a.StreetNumber = ss.StreetNumber
		a.StreetSuffix = ss.StreetSuffix
		a.City = ss.City
		a.State = ss.State
		a.Zip = ss.Zip
		a.Latitude = ss.Latitude
		a.Longitude = ss.Longitude
	default:
		return errors.ErrInternal.Msgf(
			"%T received the unexpected snapshot %T", a, snapshot)
	}
	return nil
}

func (a Address) ToSnapshot() es.Snapshot {
	return AddressV1{}
}
