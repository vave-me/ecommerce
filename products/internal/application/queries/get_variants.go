package queries

import (
	"context"
	"middleman/products/internal/domain"
)

type GetVariants struct {
	ProductID string
	Page      int64
	PageSize  int64
	SortBy    string
	SortOrder string
}

type GetVariantsHandler struct {
	variantCatalog domain.CatalogVariantRepository
}

func NewGetVariantsHandler(variantCatalog domain.CatalogVariantRepository) GetVariantsHandler {
	return GetVariantsHandler{variantCatalog: variantCatalog}
}

func (h GetVariantsHandler) GetVariants(ctx context.Context, query GetVariants) ([]*domain.CatalogVariant, int64, error) {
	catalogVariants, total, err := h.variantCatalog.GetVariants(ctx, query.ProductID, query.Page, query.PageSize, query.SortBy, query.SortOrder)
	if err != nil {
		return nil, 0, err
	}

	// Convert CatalogVariant to Variant
	variants := make([]*domain.CatalogVariant, len(catalogVariants))
	for i, cv := range catalogVariants {
		variants[i] = &domain.CatalogVariant{
			ID:           cv.ID,
			ProductID:    cv.ProductID,
			Status:       cv.Status,
			SKU:          cv.SKU,
			Barcode:      cv.Barcode,
			Condition:    cv.Condition,
			VariantPrice: cv.VariantPrice,
			CurrencyCode: cv.CurrencyCode,
			Stock:        cv.Stock,
			Weight:       cv.Weight,
			Height:       cv.Height,
			Width:        cv.Width,
			Depth:        cv.Depth,
			Attributes:   cv.Attributes,
			IsAvailable:  cv.IsAvailable,
			HasOptions:   cv.HasOptions,
			Options:      cv.Options,
		}
	}

	return variants, total, nil
}
