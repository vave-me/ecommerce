package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type ProductRepository interface {
	Remove(ctx context.Context, productID string) error
	Update(ctx context.Context, productID string, newBasePrice int64) error
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
	Find(ctx context.Context, productID string) (*models.Product, error)
	SearchWithFilters(
		ctx context.Context,
		name string,
		categoryId string,
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

	// Additional methods available in gRPC service
	GetCatalog(ctx context.Context, userID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Product, error)
	GetPublicCatalog(ctx context.Context, userID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Product, error)
	UpdateProductPrice(ctx context.Context, userID string, productID string, newPrice int64, oldPrice int64) error
	AdjustProductStock(ctx context.Context, userID string, productID string, newStock int64) error
	ArchiveProduct(ctx context.Context, userID string, productID string) error
	MarkProductSold(ctx context.Context, userID string, productID string) error
	MarkProductLeased(ctx context.Context, userID string, productID string, monthlyPrice int64, leaseTermMonths int64) error
	MarkProductPawned(ctx context.Context, userID string, productID string, lockedPrice int64, redemptionFee int64) error
	IncreaseProductPrice(ctx context.Context, userID string, productID string, price int64) error
	DecreaseProductPrice(ctx context.Context, userID string, productID string, newPrice int64) error
	AddProductThumbnail(ctx context.Context, productID string, thumbnail string) error
}
