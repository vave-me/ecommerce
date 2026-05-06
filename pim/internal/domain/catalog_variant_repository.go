package domain

import "context"

// CatalogVariant is the "query model" struct that the CatalogVariantRepository provides.
// It may differ from the "Variant" aggregate, but usually shares many fields.
type CatalogVariant struct {
	ID           string
	ProductID    string
	Status       ProductStatus
	SKU          string
	Barcode      string
	Condition    ProductCondition
	VariantPrice int64
	CurrencyCode string
	Stock        int64
	Weight       int64
	Height       int64
	Width        int64
	Depth        int64
	Attributes   []Attribute
	IsAvailable  bool
	HasOptions   bool
	Options      []Option
}

// CatalogVariantRepository defines the methods to handle reading/writing "CatalogVariant" records.
type CatalogVariantRepository interface {
	// Creates a variant row in your "variants" table
	AddVariant(ctx context.Context,
		variantID string,
		productID string,
		status ProductStatus,
		sku string,
		barcode string,
		condition ProductCondition,
		variantPrice int64,
		currencyCode string,
		stock int64,
		weight int64,
		height int64,
		width int64,
		depth int64,
		attributes []Attribute,
		isAvailable bool,
		hasOptions bool,
		options []Option,
	) error

	UpdateVariantPrice(ctx context.Context, variantID string, oldPrice, newPrice int64) error
	RemoveVariant(ctx context.Context, variantID string, userSellerID string) error
	FindVariant(ctx context.Context, variantID string) (*CatalogVariant, error)
	GetVariants(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogVariant, int64, error)
	GetVariantsByProduct(ctx context.Context, productID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogVariant, int64, error)
	GetVariantsByCategory(ctx context.Context, categoryID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogVariant, int64, error)
	RebrandVariant(ctx context.Context, variantID string, newName, newDescription string) error
	AdjustVariantStock(ctx context.Context, variantID string, oldStock, newStock int64) error
	ArchiveVariant(ctx context.Context, variantID string) error
}
