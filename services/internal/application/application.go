package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/services/internal/application/commands"
	"middleman/services/internal/application/queries"
	"middleman/services/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		AddService(ctx context.Context, cmd commands.AddService) error
		UpdateService(ctx context.Context, cmd commands.UpdateService) error
		RebrandService(ctx context.Context, cmd commands.RebrandService) error
		IncreaseServicePrice(ctx context.Context, cmd commands.IncreaseServicePrice) error
		DecreaseServicePrice(ctx context.Context, cmd commands.DecreaseServicePrice) error
		RemoveService(ctx context.Context, cmd commands.RemoveService) error
		ArchiveService(ctx context.Context, cmd commands.ArchiveService) error
	}

	Queries interface {
		GetCatalog(ctx context.Context, query queries.GetCatalog) ([]*domain.CatalogService, int64, error)
		GetPublicCatalog(ctx context.Context, query queries.GetPublicCatalog) ([]*domain.CatalogService, int64, error)
		GetService(ctx context.Context, query queries.GetService) (*domain.CatalogService, error)
		GetServices(ctx context.Context, query queries.GetServices) ([]*domain.CatalogService, int64, error)
		GetServicesByCategory(ctx context.Context, query queries.GetServicesByCategory) ([]*domain.CatalogService, int64, error)
		GetServicesByCategorySlug(ctx context.Context, query queries.GetServicesByCategorySlug) ([]*domain.CatalogService, int64, error)
		GetServicesWithFilter(ctx context.Context, query queries.GetServicesWithFilter) ([]*domain.CatalogService, int64, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.AddServiceHandler
		commands.UpdateServiceHandler
		commands.RebrandServiceHandler
		commands.IncreaseServicePriceHandler
		commands.DecreaseServicePriceHandler
		commands.RemoveServiceHandler
		commands.ArchiveServiceHandler
	}

	appQueries struct {
		queries.GetCatalogHandler
		queries.GetPublicCatalogHandler
		queries.GetServiceHandler
		queries.GetServicesHandler
		queries.GetServicesByCategoryHandler
		queries.GetServicesByCategorySlugHandler
		queries.GetServicesWithFilterHandler
	}
)

var _ App = (*Application)(nil)

func New(services domain.ServiceRepository,
	catalog domain.CatalogRepository,
	cacheCatalog domain.CatalogRepository,
	publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{
			AddServiceHandler:           commands.NewAddServiceHandler(services, publisher),
			UpdateServiceHandler:        commands.NewUpdateServiceHandler(services, publisher),
			RebrandServiceHandler:       commands.NewRebrandServiceHandler(services, publisher),
			IncreaseServicePriceHandler: commands.NewIncreaseServicePriceHandler(services, publisher),
			DecreaseServicePriceHandler: commands.NewDecreaseServicePriceHandler(services, publisher),
			RemoveServiceHandler:        commands.NewRemoveServiceHandler(services, publisher),
			ArchiveServiceHandler:       commands.NewArchiveServiceHandler(services, publisher),
		},
		appQueries: appQueries{
			GetCatalogHandler:                queries.NewGetCatalogHandler(catalog),
			GetPublicCatalogHandler:          queries.NewGetPublicCatalogHandler(catalog),
			GetServiceHandler:                queries.NewGetServiceHandler(catalog),
			GetServicesHandler:               queries.NewGetServicesHandler(cacheCatalog),
			GetServicesByCategoryHandler:     queries.NewGetServicesByCategoryHandler(cacheCatalog),
			GetServicesByCategorySlugHandler: queries.NewGetServicesByCategorySlugHandler(cacheCatalog),
			GetServicesWithFilterHandler:     queries.NewGetServicesWithFilterHandler(cacheCatalog),
		},
	}
}
