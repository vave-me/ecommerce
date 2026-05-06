package application

import (
	"context"
	"middleman/geocoding/internal/application/commands"
	"middleman/geocoding/internal/application/queries"
	"middleman/geocoding/internal/domain"
	"middleman/internal/ddd"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		GeocodeAddress(ctx context.Context, cmd commands.GeocodeAddress) error
		BatchGeocodeAddress(ctx context.Context, cmd commands.BatchGeocodeAddress) error
		ReverseGeocodeLocation(ctx context.Context, cmd commands.ReverseGeocodeLocation) error
		RefreshGeocodingCache(ctx context.Context, cmd commands.RefreshGeocodingCache) error
		ValidateAddress(ctx context.Context, cmd commands.ValidateAddress) error
	}

	Queries interface {
		GetAddressForCoordinates(ctx context.Context, query queries.GetAddressForCoordinates) (*domain.CatalogAddress, error)
		GetCoordinatesForAddress(ctx context.Context, query queries.GetCoordinatesForAddress) (*domain.CatalogAddress, error)
		GetGeocodingCache(ctx context.Context, query queries.GetGeocodingCache) (*domain.CatalogAddress, error)
		GetGeocodingDetails(ctx context.Context, query queries.GetGeocodingDetails) (*domain.CatalogAddress, error)
		SuggestAddress(ctx context.Context, query queries.SuggestAddress) ([]*domain.CatalogAddress, error)
	}
	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.GeocodeAddressHandler
		commands.BatchGeocodeAddressHandler
		commands.ReverseGeocodeLocationHandler
		commands.RefreshGeocodingCacheHandler
		commands.ValidateAddressHandler
	}
	appQueries struct {
		queries.GetAddressForCoordinatesHandler
		queries.GetCoordinatesForAddressHandler
		queries.GetGeocodingCacheHandler
		queries.GetGeocodingDetailsHandler
		queries.SuggestAddressHandler
	}
)

var _ App = (*Application)(nil)

func New(
	addresses domain.AddressRepository,
	locations domain.LocationRepository,
	catalog domain.CatalogRepository,
	publisher ddd.EventPublisher[ddd.Event]) *Application {

	return &Application{
		appCommands: appCommands{
			GeocodeAddressHandler:         commands.NewGeocodeAddressHandler(addresses, publisher),
			BatchGeocodeAddressHandler:    commands.NewBatchGeocodeAddressHandler(addresses, publisher),
			ReverseGeocodeLocationHandler: commands.NewReverseGeocodeLocationHandler(locations, publisher),
			RefreshGeocodingCacheHandler:  commands.NewRefreshGeocodingCacheHandler(addresses, publisher),
			ValidateAddressHandler:        commands.NewValidateAddressHandler(addresses, publisher),
		},
		appQueries: appQueries{
			GetAddressForCoordinatesHandler: queries.NewGetAddressForCoordinatesHandler(catalog),
			GetCoordinatesForAddressHandler: queries.NewGetCoordinatesForAddressHandler(catalog),
			GetGeocodingCacheHandler:        queries.NewGetGeocodingCacheHandler(catalog),
			GetGeocodingDetailsHandler:      queries.NewGetGeocodingDetailsHandler(catalog),
			SuggestAddressHandler:           queries.NewSuggestAddressHandler(catalog),
		},
	}
}
