package queries

import (
	"context"
	"middleman/categories/internal/domain"
)

type GetSubCategories struct {
	Lang             string
	ParentCategoryID string
	Page             int64  // Page number for pagination
	PageSize         int64  // Number of items per page
	SortBy           string // Field to sort by (e.g., "price", "name", "created_at")
	SortOrder        string // "asc" for ascending, "desc" for descending
}

type GetSubCategoriesHandler struct {
	catalog domain.CatalogRepository
}

func NewGetSubCategoriesHandler(catalog domain.CatalogRepository) GetSubCategoriesHandler {
	return GetSubCategoriesHandler{catalog: catalog}
}

func (h GetSubCategoriesHandler) GetSubCategories(ctx context.Context, query GetSubCategories) ([]*domain.CatalogCategory, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Lang == "" {
		query.Lang = "en" // Default sorting field
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
	return h.catalog.GetSubCategories(ctx, query.Lang, query.ParentCategoryID, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
