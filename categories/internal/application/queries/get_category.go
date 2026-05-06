package queries

import (
	"context"
	"middleman/categories/internal/domain"
)

type GetCategory struct {
	ID string
}

type GetCategoryHandler struct {
	catalog domain.CatalogRepository
}

func NewGetCategoryHandler(catalog domain.CatalogRepository) GetCategoryHandler {
	return GetCategoryHandler{catalog: catalog}
}

func (h GetCategoryHandler) GetCategory(ctx context.Context, query GetCategory) (*domain.CatalogCategory, error) {
	return h.catalog.Find(ctx, query.ID)
}
