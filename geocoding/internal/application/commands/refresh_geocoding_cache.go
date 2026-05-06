package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/geocoding/internal/domain"
	"middleman/internal/ddd"
)

type RefreshGeocodingCache struct {
	ID      string
	Address string
}

type RefreshGeocodingCacheHandler struct {
	addresses domain.AddressRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRefreshGeocodingCacheHandler(addresses domain.AddressRepository, publisher ddd.EventPublisher[ddd.Event]) RefreshGeocodingCacheHandler {
	return RefreshGeocodingCacheHandler{
		addresses: addresses,
		publisher: publisher,
	}
}

func (g RefreshGeocodingCacheHandler) RefreshGeocodingCache(ctx context.Context, cmd RefreshGeocodingCache) error {
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
