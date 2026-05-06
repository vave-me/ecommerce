package queries

import (
	"context"
	"middleman/newsletters/internal/domain"
)

type GetEdition struct {
	ID string
}

type GetEditionHandler struct {
	catalog domain.EditionCatalogRepository
}

func NewGetEditionHandler(catalog domain.EditionCatalogRepository) GetEditionHandler {
	return GetEditionHandler{catalog: catalog}
}

func (h GetEditionHandler) GetEdition(ctx context.Context, query GetEdition) (*domain.CatalogEdition, error) {
	return h.catalog.Find(ctx, query.ID)
}