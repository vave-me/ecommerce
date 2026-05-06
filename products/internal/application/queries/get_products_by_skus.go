package queries

import (
	"context"
	"middleman/products/internal/domain"
)

type GetProductsBySKUs struct {
	SKUs []string
}

type GetProductsBySKUsHandler struct {
	catalog domain.CatalogRepository
}

func NewGetProductsBySKUsHandler(catalog domain.CatalogRepository) GetProductsBySKUsHandler {
	return GetProductsBySKUsHandler{catalog: catalog}
}

func (h GetProductsBySKUsHandler) GetProductsBySKUs(ctx context.Context, query GetProductsBySKUs) ([]*domain.CatalogProduct, []string, error) {
	if len(query.SKUs) == 0 {
		return []*domain.CatalogProduct{}, []string{}, nil
	}

	products, err := h.catalog.GetProductsBySKUs(ctx, query.SKUs)
	if err != nil {
		return nil, nil, err
	}

	// Find which SKUs were not found
	foundSKUs := make(map[string]bool)
	for _, product := range products {
		foundSKUs[product.SKU] = true
	}

	var notFoundSKUs []string
	for _, sku := range query.SKUs {
		if !foundSKUs[sku] {
			notFoundSKUs = append(notFoundSKUs, sku)
		}
	}

	return products, notFoundSKUs, nil
}