package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/geocoding/internal/domain"
	"middleman/internal/ddd"
)

type ReverseGeocodeLocation struct {
	ID      string
	Address string
}

type ReverseGeocodeLocationHandler struct {
	locations domain.LocationRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewReverseGeocodeLocationHandler(locations domain.LocationRepository, publisher ddd.EventPublisher[ddd.Event]) ReverseGeocodeLocationHandler {
	return ReverseGeocodeLocationHandler{
		locations: locations,
		publisher: publisher,
	}
}

func (g ReverseGeocodeLocationHandler) ReverseGeocodeLocation(ctx context.Context, cmd ReverseGeocodeLocation) error {
	location, err := g.locations.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error geocoding address")
	}

	//TODO change to geocodeAddress
	event, err := location.InitLocation(cmd.ID, cmd.Address)
	if err != nil {
		return errors.Wrap(err, "initializing image")
	}
	return g.publisher.Publish(ctx, event)
}
