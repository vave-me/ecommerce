package queries

import (
	"context"
	"middleman/products/internal/domain"
)

type GetProductsByCategorySlug struct {
	Slug      string // Category ID to filter products
	Page      int64  // Page number for pagination
	PageSize  int64  // Number of items per page
	SortBy    string // Field to sort by (e.g., "price", "name", "created_at")
	SortOrder string // "asc" for ascending, "desc" for descending
}

type GetProductsByCategorySlugHandler struct {
	catalog domain.CatalogRepository
}

func NewGetProductsByCategorySlugHandler(catalog domain.CatalogRepository) GetProductsByCategorySlugHandler {
	return GetProductsByCategorySlugHandler{catalog: catalog}
}

func (h GetProductsByCategorySlugHandler) GetProductsByCategorySlug(ctx context.Context, query GetProductsByCategorySlug) ([]*domain.CatalogProduct, int64, error) {
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

	return h.catalog.GetProductsByCategorySlug(ctx, query.Slug, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
