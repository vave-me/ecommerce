package queries

import (
	"context"
	"middleman/products/internal/domain"
)

type GetPublicCatalog struct {
	UserSellerID string
	Page         int64  // Page number for pagination
	PageSize     int64  // Number of items per page
	SortBy       string // Field to sort by (e.g., "price", "name", "created_at")
	SortOrder    string // "asc" for ascending, "desc" for descending
}

type GetPublicCatalogHandler struct {
	catalog domain.CatalogRepository
}

func NewGetPublicCatalogHandler(catalog domain.CatalogRepository) GetPublicCatalogHandler {
	return GetPublicCatalogHandler{catalog: catalog}
}

func (h GetPublicCatalogHandler) GetPublicCatalog(ctx context.Context, query GetPublicCatalog) ([]*domain.CatalogProduct, int64, error) {
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
	return h.catalog.GetPublicCatalog(ctx, query.UserSellerID, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
