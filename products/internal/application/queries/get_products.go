package queries

import (
	"context"
	"middleman/products/internal/domain"
)

type GetProducts struct {
	Page      int64  // Page number for pagination
	PageSize  int64  // Number of items per page
	SortBy    string // Field to sort by (e.g., "price", "name", "created_at")
	SortOrder string // "asc" for ascending, "desc" for descending
}

type GetProductsHandler struct {
	catalog domain.CatalogRepository
}

func NewGetProductsHandler(catalog domain.CatalogRepository) GetProductsHandler {
	return GetProductsHandler{catalog: catalog}
}

func (h GetProductsHandler) GetProducts(ctx context.Context, query GetProducts) ([]*domain.CatalogProduct, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.SortBy == "" {
		query.SortBy = "name" // Default sorting field
	}
	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		query.SortOrder = "asc" // Default sorting order
	}

	return h.catalog.GetProducts(ctx, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
