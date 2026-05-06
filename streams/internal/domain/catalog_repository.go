package domain

import (
	"context"
)

type CatalogStream struct {
	ID               string
	Name             string
	Description      string
	BasePrice        int64
	UserSellerID     string
	CategoryID       string
	CategorySlug     string
	Brand            string
	Condition        StreamCondition
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
	Status           StreamStatus
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
	AddStream(ctx context.Context,
		id, name, description string,
		basePrice int64,
		userSellerID, categoryID, categorySlug, brand string,
		condition StreamCondition,
		model string,
		tags []string,
		manageStock bool,
		stock int64,
		sku string,
		attributes []Attribute,
		weight, height, width, depth int64,
		status StreamStatus,
		negotiable bool,
		userType UserType,
		middlemanService bool,
		shippingCost int64,
		hasVariants bool,
		options []Option,
		thumbnail string,
		lat, long float64) error
	UpdatePrice(ctx context.Context, productID string, oldPrice, newPrice int64) error
	GetStreamsWithFilters(ctx context.Context,
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
		lat, lng float64,
		radius int64,
		page, pageSize int64,
		sortBy, sortOrder string) ([]*CatalogStream, int64, error)
	RemoveStream(ctx context.Context, productID string, userSellerID string) error
	Find(ctx context.Context, productID string) (*CatalogStream, error)
	GetCatalog(ctx context.Context, userSellerID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogStream, int64, error)
	GetPublicCatalog(ctx context.Context, userSellerID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogStream, int64, error)
	GetStreams(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogStream, int64, error)
	GetStreamsByCategory(ctx context.Context, categoryID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogStream, int64, error)
	GetStreamsByCategorySlug(ctx context.Context, categorySlug string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogStream, int64, error)
	RebrandStream(ctx context.Context, productID string, newName, newDescription, newCategoryID, newCategorySlug, newBrand, newModel, newCondition string, newTags []string) error
	AdjustStock(ctx context.Context, productID string, userSellerID string, oldStock, newStock int64) error
	ArchiveStream(ctx context.Context, productID string, userSellerID string) error
	MarkStreamSold(ctx context.Context, productID string, userSellerID string, finalPrice int64) error
	MarkStreamLeased(ctx context.Context, productID string, userSellerID string) error
	MarkStreamPawned(ctx context.Context, productID string, userSellerID string) error
	ToggleNegotiable(ctx context.Context, productID string, userSellerID string, currentValue bool) error
	UpdateStream(ctx context.Context, productID string, name string, description string, basePrice int64, categoryID, categorySlug string, brand string, condition StreamCondition, model string, tags []string, manageStock bool, stock int64, sku string, attributes []Attribute, weight, height, width, depth int64, status StreamStatus, negotiable bool, userType UserType, middlemanService bool, shippingCost int64, hasVariants bool, options []Option, thumbnail string, lat, lng float64) error
	UpdateThumbnail(ctx context.Context, productID string, thumbnail string) error
	GetStreamBySKU(ctx context.Context, sku string) (*CatalogStream, error)
	GetStreamsBySKUs(ctx context.Context, skus []string) ([]*CatalogStream, error)
}
