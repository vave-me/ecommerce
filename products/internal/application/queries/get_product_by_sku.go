package queries

import (
	"context"
	"database/sql"

	"github.com/stackus/errors"
	"middleman/products/internal/domain"
)

type GetProductBySKU struct {
	SKU string
}

type GetProductBySKUHandler struct {
	catalog domain.CatalogRepository
}

func NewGetProductBySKUHandler(catalog domain.CatalogRepository) GetProductBySKUHandler {
	return GetProductBySKUHandler{catalog: catalog}
}

func (h GetProductBySKUHandler) GetProductBySKU(ctx context.Context, query GetProductBySKU) (*domain.CatalogProduct, error) {
	product, err := h.catalog.GetProductBySKU(ctx, query.SKU)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "product not found with SKU")
		}
		return nil, err
	}

	return product, nil
}