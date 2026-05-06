package queries

import (
	"context"
	"middleman/posts/internal/domain"
)

type GetUserPosts struct {
	UserId    string
	Page      int64  // Page number for pagination
	PageSize  int64  // Number of items per page
	SortBy    string // Field to sort by (e.g., "price", "name", "created_at")
	SortOrder string // "asc" for ascending, "desc" for descending
}

type GetUserPostsHandler struct {
	catalog domain.CatalogRepository
}

func NewGetUserPostsHandler(catalog domain.CatalogRepository) GetUserPostsHandler {
	return GetUserPostsHandler{catalog: catalog}
}

func (h GetUserPostsHandler) GetUserPosts(ctx context.Context, query GetUserPosts) ([]*domain.CatalogPost, int64, error) {
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
	return h.catalog.GetUserPosts(ctx, query.UserId, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
