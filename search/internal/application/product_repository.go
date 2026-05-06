package application

import (
	"context"
	"middleman/search/internal/models"
)

type ProductRepository interface {
	Find(ctx context.Context, productID string) (*models.Product, error)
	SearchWithFilters(
		ctx context.Context,
		name string,
		categoryID string,
		categorySlug string,
		minPrice int64,
		maxPrice int64,
		brand string,
		condition string,
		model string,
		tags []string,
		manageStock bool,
		minStock int64,
		maxStock int64,
		sku string,
		status string,
		negotiable bool,
		userType string,
		middlemanService bool,
		hasVariants bool,
		shippingCost int64,
		minWeight int64,
		maxWeight int64,
		minHeight int64,
		maxHeight int64,
		minWidth int64,
		maxWidth int64,
		minDepth int64,
		maxDepth int64,
		offset int64,
		limit int64,
		lat float64,
		lng float64,
		radius int64,
		page int64,
		pageSize int64,
		sortBy string,
		sortOrder string,
	) ([]*models.Product, error)

	SearchProductsWithCategorySlug(ctx context.Context, categorySlug string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Product, error)
	SearchProductsWithCategory(ctx context.Context, categoryId string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Product, error)
	GetCatalog(ctx context.Context, userId string) ([]*models.Product, error)
}

type ProductCacheRepository interface {
	ProductRepository
	Remove(ctx context.Context, productID string) error
	Update(ctx context.Context,
		productID, name, description string,
		basePrice int64,
		categoryID, categorySlug, brand, condition, model string,
		tags []string,
		manageStock bool,
		stock int64,
		sku string,
		attributes []models.Attribute,
		weight, height, width, depth int64,
		status string,
		negotiable, middlemanService bool,
		userType string,
		shippingCost int64,
		hasVariants bool,
		option []models.Option,
		thumbnail string, lat, lng float64) error
	Add(
		ctx context.Context,
		productID, name, description string,
		basePrice int64,
		userSellerID, categoryID, categorySlug string,
		brand string,
		condition string,
		model string, tags []string,
		manageStock bool,
		stock int64,
		sku string,
		attributes []models.Attribute,
		weight int64,
		height int64,
		width int64,
		depth int64,
		status string,
		negotiable bool,
		userType string,
		middlemanService bool,
		shippingCost int64,
		hasVariants bool,
		options []models.Option,
		lat, lng float64,
		thumbnail string,
		entityType models.EntityType,
	) error

	Rebrand(
		ctx context.Context,
		productID, name, description string,
		basePrice int64,
		stock int64,
		sku string,
		categoryID string,
		status string,
	) error

	SearchWithTerm(ctx context.Context, name string) ([]*models.Product, error)
	SuggestProducts(ctx context.Context, name string) ([]*models.Product, error)
	UpdateThumbnail(ctx context.Context, productID string, thumbnail string) error

	// Dedicated methods for each event type to avoid bottlenecks
	IncreasePrice(ctx context.Context, productID string, newPrice int64) error
	DecreasePrice(ctx context.Context, productID string, newPrice int64) error
	MarkAsLeased(ctx context.Context, productID string) error
	MarkAsSold(ctx context.Context, productID string) error
	MarkAsPawned(ctx context.Context, productID string) error
	AdjustStock(ctx context.Context, productID string, newStock int64) error
	ToggleNegotiable(ctx context.Context, productID string, negotiable bool) error
	ArchiveProduct(ctx context.Context, productID string) error
}
