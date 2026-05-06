package queries

import (
	"context"
	"middleman/services/internal/domain"
)

type GetServicesByCategory struct {
	CategoryID string // Category ID to filter services
	Page       int64  // Page number for pagination
	PageSize   int64  // Number of items per page
	SortBy     string // Field to sort by (e.g., "price", "name", "created_at")
	SortOrder  string // "asc" for ascending, "desc" for descending
}

type GetServicesByCategoryHandler struct {
	catalog domain.CatalogRepository
}

func NewGetServicesByCategoryHandler(catalog domain.CatalogRepository) GetServicesByCategoryHandler {
	return GetServicesByCategoryHandler{catalog: catalog}
}

func (h GetServicesByCategoryHandler) GetServicesByCategory(ctx context.Context, query GetServicesByCategory) ([]*domain.CatalogService, int64, error) {
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

	return h.catalog.GetServicesByCategory(ctx, query.CategoryID, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
