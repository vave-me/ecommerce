package queries

import (
	"context"
	"middleman/posts/internal/domain"
)

type GetPostsWithFilters struct {
	Name        string
	Description string
	TypeOfPost  string
	UserID      string
	UserType    string
	Tags        []string
	Lat         float64
	Lng         float64
	Radius      int64
	Page        int64
	PageSize    int64
	Offset      int64
	Limit       int64
	SortBy      string
	SortOrder   string
}
type GetPostsWithFiltersHandler struct {
	catalog domain.CatalogRepository
}

func NewGetPostsWithFiltersHandler(catalog domain.CatalogRepository) GetPostsWithFiltersHandler {
	return GetPostsWithFiltersHandler{catalog: catalog}
}

func (h GetPostsWithFiltersHandler) GetPostsWithFilters(ctx context.Context, query GetPostsWithFilters) ([]*domain.CatalogPost, int64, error) {
	// Set sensible defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.SortBy == "" {
		query.SortBy = "name" // or "created_at", etc.
	}
	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		query.SortOrder = "asc"
	}

	// Compute offset/limit based on page/pageSize
	query.Offset = (query.Page - 1) * query.PageSize
	query.Limit = query.PageSize

	return h.catalog.GetPostsWithFilters(
		ctx,
		query.Name,
		query.Description,
		query.Tags,
		query.Offset,
		query.Limit,
		query.Lat,
		query.Lng,
		query.Radius,
		query.Page,
		query.PageSize,
		query.SortBy,
		query.SortOrder,
	)
}
