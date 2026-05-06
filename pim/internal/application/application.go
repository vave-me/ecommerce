package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/products/internal/application/commands"
	"middleman/products/internal/application/queries"
	"middleman/products/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}

	Commands interface {
		AddProduct(ctx context.Context, cmd commands.AddProduct) error
		RebrandProduct(ctx context.Context, cmd commands.RebrandProduct) error
		UpdateProduct(ctx context.Context, cmd commands.UpdateProduct) error
		IncreaseProductPrice(ctx context.Context, cmd commands.IncreaseProductPrice) error
		DecreaseProductPrice(ctx context.Context, cmd commands.DecreaseProductPrice) error
		RemoveProduct(ctx context.Context, cmd commands.RemoveProduct) error
		AdjustProductStock(ctx context.Context, cmd commands.AdjustProductStock) error
		ArchiveProduct(ctx context.Context, cmd commands.ArchiveProduct) error
		MarkProductSold(ctx context.Context, cmd commands.MarkProductSold) error
		MarkProductLeased(ctx context.Context, cmd commands.MarkProductLeased) error
		MarkProductPawned(ctx context.Context, cmd commands.MarkProductPawned) error
		ReserveProduct(ctx context.Context, cmd commands.ReserveProduct) error
		ReleaseProduct(ctx context.Context, cmd commands.ReleaseProduct) error
		AddVariant(ctx context.Context, cmd commands.AddVariant) error
		RebrandVariant(ctx context.Context, cmd commands.RebrandVariant) error
		IncreaseVariantPrice(ctx context.Context, cmd commands.IncreaseVariantPrice) error
		DecreaseVariantPrice(ctx context.Context, cmd commands.DecreaseVariantPrice) error
		AdjustVariantStock(ctx context.Context, cmd commands.AdjustVariantStock) error
		ArchiveVariant(ctx context.Context, cmd commands.ArchiveVariant) error
		RemoveVariant(ctx context.Context, cmd commands.RemoveVariant) error
		AddProductThumbnail(ctx context.Context, cmd commands.AddProductThumbnail) error
		UpdateProductThumbnail(ctx context.Context, cmd commands.UpdateProductThumbnail) error
	}

	Queries interface {
		GetCatalog(ctx context.Context, query queries.GetCatalog) ([]*domain.CatalogProduct, int64, error)
		GetPublicCatalog(ctx context.Context, query queries.GetPublicCatalog) ([]*domain.CatalogProduct, int64, error)
		GetProduct(ctx context.Context, query queries.GetProduct) (*domain.CatalogProduct, error)
		GetProducts(ctx context.Context, query queries.GetProducts) ([]*domain.CatalogProduct, int64, error)
		GetProductsByCategory(ctx context.Context, query queries.GetProductsByCategory) ([]*domain.CatalogProduct, int64, error)
		GetProductsByCategorySlug(ctx context.Context, query queries.GetProductsByCategorySlug) ([]*domain.CatalogProduct, int64, error)
		GetProductsWithFilters(ctx context.Context, query queries.GetProductsWithFilters) ([]*domain.CatalogProduct, int64, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.AddProductHandler
		commands.RebrandProductHandler
		commands.UpdateProductHandler
		commands.IncreaseProductPriceHandler
		commands.DecreaseProductPriceHandler
		commands.RemoveProductHandler

		// NEW:
		commands.AdjustProductStockHandler
		commands.ArchiveProductHandler
		commands.MarkProductSoldHandler
		commands.MarkProductLeasedHandler
		commands.MarkProductPawnedHandler
		commands.ReserveProductHandler
		commands.ReleaseProductHandler

		// VARIANTS:
		commands.AddVariantHandler
		commands.RebrandVariantHandler
		commands.IncreaseVariantPriceHandler
		commands.DecreaseVariantPriceHandler
		commands.AdjustVariantStockHandler
		commands.ArchiveVariantHandler
		commands.RemoveVariantHandler
		commands.AddProductThumbnailHandler
		commands.UpdateProductThumbnailHandler
	}

	appQueries struct {
		queries.GetCatalogHandler
		queries.GetPublicCatalogHandler
		queries.GetProductHandler
		queries.GetProductsHandler
		queries.GetProductsByCategoryHandler
		queries.GetProductsByCategorySlugHandler
		queries.GetProductsWithFiltersHandler
	}
)

var _ App = (*Application)(nil)

func New(products domain.ProductRepository,
	catalog domain.CatalogRepository,
	cacheCatalog domain.CatalogRepository,
	variants domain.VariantRepository,
	variantCatalog domain.CatalogVariantRepository,
	publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{
			AddProductHandler:             commands.NewAddProductHandler(products, publisher),
			UpdateProductHandler:          commands.NewUpdateProductHandler(products, publisher),
			RebrandProductHandler:         commands.NewRebrandProductHandler(products, publisher),
			IncreaseProductPriceHandler:   commands.NewIncreaseProductPriceHandler(products, publisher),
			DecreaseProductPriceHandler:   commands.NewDecreaseProductPriceHandler(products, publisher),
			RemoveProductHandler:          commands.NewRemoveProductHandler(products, publisher),
			AdjustProductStockHandler:     commands.NewAdjustProductStockHandler(products, publisher),
			ArchiveProductHandler:         commands.NewArchiveProductHandler(products, publisher),
			MarkProductSoldHandler:        commands.NewMarkProductSoldHandler(products, publisher),
			MarkProductLeasedHandler:      commands.NewMarkProductLeasedHandler(products, publisher),
			MarkProductPawnedHandler:      commands.NewMarkProductPawnedHandler(products, publisher),
			ReserveProductHandler:         commands.NewReserveProductHandler(products, publisher),
			ReleaseProductHandler:         commands.NewReleaseProductHandler(products, publisher),
			AddVariantHandler:             commands.NewAddVariantHandler(variants, publisher),
			RebrandVariantHandler:         commands.NewRebrandVariantHandler(variants, publisher),
			IncreaseVariantPriceHandler:   commands.NewIncreaseVariantPriceHandler(variants, publisher),
			DecreaseVariantPriceHandler:   commands.NewDecreaseVariantPriceHandler(variants, publisher),
			AdjustVariantStockHandler:     commands.NewAdjustVariantStockHandler(variants, publisher),
			ArchiveVariantHandler:         commands.NewArchiveVariantHandler(variants, publisher),
			RemoveVariantHandler:          commands.NewRemoveVariantHandler(variants, publisher),
			AddProductThumbnailHandler:    commands.NewAddProductThumbnailHandler(products, publisher),
			UpdateProductThumbnailHandler: commands.NewUpdateProductThumbnailHandler(products, publisher),
		},
		appQueries: appQueries{
			GetCatalogHandler:                queries.NewGetCatalogHandler(catalog),
			GetPublicCatalogHandler:          queries.NewGetPublicCatalogHandler(catalog),
			GetProductHandler:                queries.NewGetProductHandler(catalog),
			GetProductsHandler:               queries.NewGetProductsHandler(cacheCatalog),
			GetProductsByCategoryHandler:     queries.NewGetProductsByCategoryHandler(cacheCatalog),
			GetProductsByCategorySlugHandler: queries.NewGetProductsByCategorySlugHandler(cacheCatalog),
			GetProductsWithFiltersHandler:    queries.NewGetProductsWithFiltersHandler(catalog),
		},
	}
}
