package queries

import (
	"context"
	"middleman/products/internal/domain"
)

type GetProductsByCategory struct {
	CategoryID string // Category ID to filter products
	Page       int64  // Page number for pagination
	PageSize   int64  // Number of items per page
	SortBy     string // Field to sort by (e.g., "price", "name", "created_at")
	SortOrder  string // "asc" for ascending, "desc" for descending
}

type GetProductsByCategoryHandler struct {
	catalog domain.CatalogRepository
}

func NewGetProductsByCategoryHandler(catalog domain.CatalogRepository) GetProductsByCategoryHandler {
	return GetProductsByCategoryHandler{catalog: catalog}
}

func (h GetProductsByCategoryHandler) GetProductsByCategory(ctx context.Context, query GetProductsByCategory) ([]*domain.CatalogProduct, int64, error) {
	// Set defaults for pagination and sorting
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.SortBy == "" {
		query.SortBy = "name" // Default sort field
	}
	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		query.SortOrder = "asc" // Default sort order
	}

	return h.catalog.GetProductsByCategory(ctx, query.CategoryID, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
