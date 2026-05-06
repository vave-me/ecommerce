package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/geocoding/internal/domain"
	"middleman/internal/ddd"
)

type BatchGeocodeAddress struct {
	ID      string
	Address string
}

type BatchGeocodeAddressHandler struct {
	addresses domain.AddressRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewBatchGeocodeAddressHandler(addresses domain.AddressRepository, publisher ddd.EventPublisher[ddd.Event]) BatchGeocodeAddressHandler {
	return BatchGeocodeAddressHandler{
		addresses: addresses,
		publisher: publisher,
	}
}

func (g BatchGeocodeAddressHandler) BatchGeocodeAddress(ctx context.Context, cmd BatchGeocodeAddress) error {
	address, err := g.addresses.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error geocoding address")
	}

	//TODO change to geocodeAddress
	event, err := address.BatchDecodeAddress(cmd.ID, cmd.Address)
	if err != nil {
		return errors.Wrap(err, "initializing image")
	}
	return g.publisher.Publish(ctx, event)
}
