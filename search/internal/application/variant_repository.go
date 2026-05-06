package application

import (
	"context"
	"middleman/search/internal/models"
)

type VariantRepository interface {
	Find(ctx context.Context, variantID string) (*models.Variant, error)
	Update(
		ctx context.Context,
		variantID string,
		newVariantPrice int64,
		newStock int64,
		newName string,
		newAttributes []models.Attribute) error
	Remove(ctx context.Context, variantID string) error
}

type VariantCacheRepository interface {
	Add(
		ctx context.Context,
		variantID, productID, name, sku, barcode string,
		variantPrice int64,
		currencyCode string,
		stock, weight, height, width, depth int64,
		attributes []models.Attribute,
		isAvailable bool,
		hasOptions bool,
		options []models.Option) error
	Rebrand(
		ctx context.Context,
		variantID, name string,
		variantPrice, stock int64,
		attributes []models.Attribute,
		isAvailable bool) error
	SearchWithFilters(
		ctx context.Context,
		name string,
		minPrice int64,
		maxPrice int64,
		offset int64,
		limit int64,
	) ([]*models.Variant, error)
	SearchWithTerm(ctx context.Context, name string) ([]*models.Variant, error)
	SuggestVariants(ctx context.Context, partialName string) ([]*models.Variant, error)
	VariantRepository
}
