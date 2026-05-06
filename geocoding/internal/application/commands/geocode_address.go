package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/geocoding/internal/domain"
	"middleman/internal/ddd"
)

type GeocodeAddress struct {
	ID        string
	Address   string
	Latitude  float64
	Longitude float64
}

type GeocodeAddressHandler struct {
	addresses domain.AddressRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewGeocodeAddressHandler(addresses domain.AddressRepository, publisher ddd.EventPublisher[ddd.Event]) GeocodeAddressHandler {
	return GeocodeAddressHandler{
		addresses: addresses,
		publisher: publisher,
	}
}

func (g GeocodeAddressHandler) GeocodeAddress(ctx context.Context, cmd GeocodeAddress) error {
	address, err := g.addresses.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error geocoding address")
	}

	//TODO change to geocodeAddress
	event, err := address.InitAddress(cmd.ID, cmd.Address, cmd.Latitude, cmd.Longitude)
	if err != nil {
		return errors.Wrap(err, "initializing image")
	}
	return g.publisher.Publish(ctx, event)
}
