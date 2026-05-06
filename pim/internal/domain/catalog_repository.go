package domain

import (
	"context"
)

type CatalogProduct struct {
	ID               string
	Name             string
	Description      string
	BasePrice        int64
	UserSellerID     string
	CategoryID       string
	CategorySlug     string
	Brand            string
	Condition        ProductCondition
	Model            string
	Tags             []string
	ManageStock      bool
	Stock            int64
	SKU              string
	Attributes       []Attribute
	Weight           int64
	Height           int64
	Width            int64
	Depth            int64
	Status           ProductStatus
	Negotiable       bool
	UserType         UserType
	MiddlemanService bool
	ShippingCost     int64
	HasVariants      bool
	Options          []Option
	Thumbnail        string
	Lat              float64
	Lng              float64
}

type CatalogRepository interface {
	AddProduct(ctx context.Context,
		id, name, description string,
		basePrice int64,
		userSellerID, categoryID, categorySlug, brand string,
		condition ProductCondition,
		model string,
		tags []string,
		manageStock bool,
		stock int64,
		sku string,
		attributes []Attribute,
		weight, height, width, depth int64,
		status ProductStatus,
		negotiable bool,
		userType UserType,
		middlemanService bool,
		shippingCost int64,
		hasVariants bool,
		options []Option,
		thumbnail string,
		lat, long float64) error
	UpdatePrice(ctx context.Context, productID string, oldPrice, newPrice int64) error
	GetProductsWithFilters(ctx context.Context,
		name string,
		categoryId string,
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
		lat, lng float64,
		radius int64,
		page, pageSize int64,
		sortBy, sortOrder string) ([]*CatalogProduct, int64, error)
	RemoveProduct(ctx context.Context, productID string, userSellerID string) error
	Find(ctx context.Context, productID string) (*CatalogProduct, error)
	GetCatalog(ctx context.Context, userSellerID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogProduct, int64, error)
	GetPublicCatalog(ctx context.Context, userSellerID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogProduct, int64, error)
	GetProducts(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogProduct, int64, error)
	GetProductsByCategory(ctx context.Context, categoryID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogProduct, int64, error)
	GetProductsByCategorySlug(ctx context.Context, categorySlug string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogProduct, int64, error)
	RebrandProduct(ctx context.Context, productID string, newName, newDescription, newCategoryID, newCategorySlug, newBrand, newModel, newCondition string, newTags []string) error
	AdjustStock(ctx context.Context, productID string, userSellerID string, oldStock, newStock int64) error
	ArchiveProduct(ctx context.Context, productID string, userSellerID string) error
	MarkProductSold(ctx context.Context, productID string, userSellerID string, finalPrice int64) error
	MarkProductLeased(ctx context.Context, productID string, userSellerID string) error
	MarkProductPawned(ctx context.Context, productID string, userSellerID string) error
	ToggleNegotiable(ctx context.Context, productID string, userSellerID string, currentValue bool) error
	UpdateProduct(ctx context.Context, productID string, name string, description string, basePrice int64, categoryID, categorySlug string, brand string, condition ProductCondition, model string, tags []string, manageStock bool, stock int64, sku string, attributes []Attribute, weight, height, width, depth int64, status ProductStatus, negotiable bool, userType UserType, middlemanService bool, shippingCost int64, hasVariants bool, options []Option, thumbnail string, lat, lng float64) error
	UpdateThumbnail(ctx context.Context, productID string, thumbnail string) error
}
