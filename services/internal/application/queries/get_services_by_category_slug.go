package queries

import (
	"context"
	"middleman/services/internal/domain"
)

type GetServicesByCategorySlug struct {
	CategorySlug string // Category ID to filter services
	Page         int64  // Page number for pagination
	PageSize     int64  // Number of items per page
	SortBy       string // Field to sort by (e.g., "price", "name", "created_at")
	SortOrder    string // "asc" for ascending, "desc" for descending
}

type GetServicesByCategorySlugHandler struct {
	catalog domain.CatalogRepository
}

func NewGetServicesByCategorySlugHandler(catalog domain.CatalogRepository) GetServicesByCategorySlugHandler {
	return GetServicesByCategorySlugHandler{catalog: catalog}
}

func (h GetServicesByCategorySlugHandler) GetServicesByCategorySlug(ctx context.Context, query GetServicesByCategorySlug) ([]*domain.CatalogService, int64, error) {
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

	return h.catalog.GetServicesByCategorySlug(ctx, query.CategorySlug, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
