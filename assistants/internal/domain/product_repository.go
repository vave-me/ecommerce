package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type ProductRepository interface {
	// Basic CRUD operations with clear names
	DeleteProduct(ctx context.Context, productID string) error
	UpdateProductPrice(ctx context.Context, productID string, newBasePrice int64) error
	CreateProduct(
		ctx context.Context,
		name, description string,
		basePrice int64,
		categoryID, categorySlug string,
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

	UpdateProductDetails(
		ctx context.Context,
		productID, name, description string,
		basePrice int64,
		stock int64,
		sku string,
		categoryID string,
		status string,
	) error

	// Search and discovery methods
	SearchProductsByName(ctx context.Context, name string) ([]*models.Product, error)
	GetProductSuggestions(ctx context.Context, name string) ([]*models.Product, error)
	UpdateProductThumbnail(ctx context.Context, productID string, thumbnail string) error
	GetProductByID(ctx context.Context, productID string) (*models.Product, error)
	SearchProductsAdvanced(
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
	GetProductsByCategorySlug(ctx context.Context, categorySlug string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Product, error)
	GetProductsByCategoryID(ctx context.Context, categoryId string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Product, error)
	GetUserProductCatalog(ctx context.Context, userID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Product, error)
	GetPublicProductCatalog(ctx context.Context, userID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Product, error)
	UpdateProductPriceForUser(ctx context.Context, productID string, newPrice int64, oldPrice int64) error
	AdjustProductInventory(ctx context.Context, productID string, newStock int64) error
	ArchiveUserProduct(ctx context.Context, userID string, productID string) error
	MarkProductAsSold(ctx context.Context, productID string) error
	MarkProductAsLeased(ctx context.Context, productID string, monthlyPrice int64, leaseTermMonths int64) error
	MarkProductAsPawned(ctx context.Context, userID string, productID string, lockedPrice int64, redemptionFee int64) error
	IncreaseProductPriceBy(ctx context.Context, productID string, increaseAmount int64) error
	DecreaseProductPriceTo(ctx context.Context, productID string, newPrice int64) error
	AddThumbnailToProduct(ctx context.Context, productID string, thumbnail string) error
}
