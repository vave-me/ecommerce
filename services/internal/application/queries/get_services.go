package queries

import (
	"context"
	"middleman/services/internal/domain"
)

type GetServices struct {
	Page      int64  // Page number for pagination
	PageSize  int64  // Number of items per page
	SortBy    string // Field to sort by (e.g., "price", "name", "created_at")
	SortOrder string // "asc" for ascending, "desc" for descending
}

type GetServicesHandler struct {
	catalog domain.CatalogRepository
}

func NewGetServicesHandler(catalog domain.CatalogRepository) GetServicesHandler {
	return GetServicesHandler{catalog: catalog}
}

func (h GetServicesHandler) GetServices(ctx context.Context, query GetServices) ([]*domain.CatalogService, int64, error) {
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
	return h.catalog.GetServices(ctx, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
