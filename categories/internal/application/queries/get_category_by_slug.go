package queries

import (
	"context"
	"middleman/categories/internal/domain"
)

type GetCategoryBySlug struct {
	Slug string
}

type GetCategoryBySlugHandler struct {
	catalog domain.CatalogRepository
}

func NewGetCategoryBySlugHandler(catalog domain.CatalogRepository) GetCategoryBySlugHandler {
	return GetCategoryBySlugHandler{catalog: catalog}
}

func (h GetCategoryBySlugHandler) GetCategoryBySlug(ctx context.Context, query GetCategoryBySlug) (*domain.CatalogCategory, error) {
	return h.catalog.Find(ctx, query.Slug)
}
