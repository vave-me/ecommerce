package geocodingpb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

// Channels
const (
	AddressAggregateChannel  = "middleman.geocoding.events.Address"
	LocationAggregateChannel = "middleman.products.events.Location"
	AddressCreatedEvent      = "geocodingapi.AddressCreated"
	LocationAddedEvent       = "geocodingapi.LocationAdded"
)

// Registrations and Serde
func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Product
	if err := serde.Register(&AddressCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&LocationAdded{}); err != nil {
		return err
	}

	return nil
}

func (*AddressCreated) Key() string { return AddressCreatedEvent }
func (*LocationAdded) Key() string  { return LocationAddedEvent }
