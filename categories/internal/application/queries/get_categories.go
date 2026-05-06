package queries

import (
	"context"
	"middleman/categories/internal/domain"
)

type GetCategories struct {
	CategoryType string
	Lang         string
	Page         int64  // Page number for pagination
	PageSize     int64  // Number of items per page
	SortBy       string // Field to sort by (e.g., "price", "name", "created_at")
	SortOrder    string // "asc" for ascending, "desc" for descending
}

type GetCategoriesHandler struct {
	catalog domain.CatalogRepository
}

func NewGetCategoriesHandler(catalog domain.CatalogRepository) GetCategoriesHandler {
	return GetCategoriesHandler{catalog: catalog}
}

func (h GetCategoriesHandler) GetCategories(ctx context.Context, query GetCategories) ([]*domain.CatalogCategory, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.SortBy == "" {
		query.SortBy = "name" // Default sorting field
	}
	if query.Lang == "" {
		query.Lang = "en" // Default sorting field
	}
	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		query.SortOrder = "asc" // Default sorting order
	}
	return h.catalog.GetCategories(ctx, query.CategoryType, query.Lang, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
