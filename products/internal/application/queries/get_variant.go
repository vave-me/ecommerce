package queries

import (
	"context"
	"middleman/products/internal/domain"
)

type GetVariant struct {
	VariantID string
}

type GetVariantHandler struct {
	variantCatalog domain.CatalogVariantRepository
}

func NewGetVariantHandler(variantCatalog domain.CatalogVariantRepository) GetVariantHandler {
	return GetVariantHandler{variantCatalog: variantCatalog}
}

func (h GetVariantHandler) GetVariant(ctx context.Context, query GetVariant) (*domain.CatalogVariant, error) {
	return h.variantCatalog.GetVariant(ctx, query.VariantID)
}
