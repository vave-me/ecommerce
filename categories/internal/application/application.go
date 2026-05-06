package application

import (
	"context"
	"middleman/categories/internal/application/commands"
	"middleman/categories/internal/application/queries"
	"middleman/categories/internal/domain"
	"middleman/internal/ddd"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		AddCategory(ctx context.Context, cmd commands.AddCategory) error
		RebrandCategory(ctx context.Context, cmd commands.RebrandCategory) error
		UpdateCategory(ctx context.Context, cmd commands.UpdateCategory) error
		RemoveCategory(ctx context.Context, cmd commands.RemoveCategory) error
		ArchiveCategory(ctx context.Context, cmd commands.ArchiveCategory) error
		AddFilter(ctx context.Context, cmd commands.AddFilter) error
		RebrandFilter(ctx context.Context, cmd commands.RebrandFilter) error
		ArchiveFilter(ctx context.Context, cmd commands.ArchiveFilter) error
		RemoveFilter(ctx context.Context, cmd commands.RemoveFilter) error
	}

	Queries interface {
		GetCatalog(ctx context.Context, query queries.GetCatalog) ([]*domain.CatalogCategory, int64, error)
		GetCategory(ctx context.Context, query queries.GetCategory) (*domain.CatalogCategory, error)
		GetCategoryBySlug(ctx context.Context, query queries.GetCategoryBySlug) (*domain.CatalogCategory, error)
		GetCategories(ctx context.Context, query queries.GetCategories) ([]*domain.CatalogCategory, int64, error)
		GetMainCategories(ctx context.Context, query queries.GetMainCategories) ([]*domain.CatalogCategory, int64, error)
		GetAllMainCategories(ctx context.Context, query queries.GetAllMainCategories) ([]*domain.CatalogCategory, int64, error)
		GetSubCategories(ctx context.Context, query queries.GetSubCategories) ([]*domain.CatalogCategory, int64, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.AddCategoryHandler
		commands.RebrandCategoryHandler
		commands.UpdateCategoryHandler
		commands.RemoveCategoryHandler
		commands.ArchiveCategoryHandler
		commands.AddFilterHandler
		commands.RebrandFilterHandler
		commands.ArchiveFilterHandler
		commands.RemoveFilterHandler
	}

	appQueries struct {
		queries.GetCatalogHandler
		queries.GetCategoryHandler
		queries.GetCategoryBySlugHandler
		queries.GetCategoriesHandler
		queries.GetMainCategoriesHandler
		queries.GetAllMainCategoriesHandler
		queries.GetSubCategoriesHandler
	}
)

var _ App = (*Application)(nil)

func New(categories domain.CategoryRepository,
	catalog domain.CatalogRepository,
	cacheCatalog domain.CatalogRepository,
	filters domain.FilterRepository,
	filterCatalog domain.CatalogFilterRepository,
	publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{
			AddCategoryHandler:     commands.NewAddCategoryHandler(categories, publisher),
			RebrandCategoryHandler: commands.NewRebrandCategoryHandler(categories, publisher),
			UpdateCategoryHandler:  commands.NewUpdateCategoryHandler(categories, publisher),
			RemoveCategoryHandler:  commands.NewRemoveCategoryHandler(categories, publisher),
			ArchiveCategoryHandler: commands.NewArchiveCategoryHandler(categories, publisher),
			AddFilterHandler:       commands.NewAddFilterHandler(filters, publisher),
			RebrandFilterHandler:   commands.NewRebrandFilterHandler(filters, publisher),
			ArchiveFilterHandler:   commands.NewArchiveFilterHandler(filters, publisher),
			RemoveFilterHandler:    commands.NewRemoveFilterHandler(filters, publisher),
		},
		appQueries: appQueries{
			GetCatalogHandler:           queries.NewGetCatalogHandler(catalog),
			GetCategoryHandler:          queries.NewGetCategoryHandler(catalog),
			GetCategoryBySlugHandler:    queries.NewGetCategoryBySlugHandler(catalog),
			GetCategoriesHandler:        queries.NewGetCategoriesHandler(cacheCatalog),
			GetMainCategoriesHandler:    queries.NewGetMainCategoriesHandler(cacheCatalog),
			GetAllMainCategoriesHandler: queries.NewGetAllMainCategoriesHandler(cacheCatalog),
			GetSubCategoriesHandler:     queries.NewGetSubCategoriesHandler(cacheCatalog),
		},
	}
}
