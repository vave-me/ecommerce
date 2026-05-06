package queries

import (
	"context"
	"middleman/posts/internal/domain"
)

type GetPostsByCategorySlug struct {
	CategorySlug string // Category ID to filter posts
	Page         int64  // Page number for pagination
	PageSize     int64  // Number of items per page
	SortBy       string // Field to sort by (e.g., "price", "name", "created_at")
	SortOrder    string // "asc" for ascending, "desc" for descending
}

type GetPostsByCategorySlugHandler struct {
	catalog domain.CatalogRepository
}

func NewGetPostsByCategorySlugHandler(catalog domain.CatalogRepository) GetPostsByCategorySlugHandler {
	return GetPostsByCategorySlugHandler{catalog: catalog}
}

func (h GetPostsByCategorySlugHandler) GetPostsByCategorySlug(ctx context.Context, query GetPostsByCategorySlug) ([]*domain.CatalogPost, int64, error) {
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

	return h.catalog.GetPostsByCategorySlug(ctx, query.CategorySlug, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
