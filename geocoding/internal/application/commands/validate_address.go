package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/geocoding/internal/domain"
	"middleman/internal/ddd"
)

type ValidateAddress struct {
	ID      string
	Address string
}

type ValidateAddressHandler struct {
	addresses domain.AddressRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewValidateAddressHandler(addresses domain.AddressRepository, publisher ddd.EventPublisher[ddd.Event]) ValidateAddressHandler {
	return ValidateAddressHandler{
		addresses: addresses,
		publisher: publisher,
	}
}

func (g ValidateAddressHandler) ValidateAddress(ctx context.Context, cmd ValidateAddress) error {
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
