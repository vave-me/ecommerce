package queries

import (
	"context"
	"middleman/tickets/internal/domain"
)

type GetMatchCatalog struct {
	Filters domain.MatchCatalogFilters
}

type GetMatchCatalogHandler struct {
	catalog domain.MatchCatalogRepository
}

func NewGetMatchCatalogHandler(catalog domain.MatchCatalogRepository) GetMatchCatalogHandler {
	return GetMatchCatalogHandler{catalog: catalog}
}

func (h GetMatchCatalogHandler) GetMatchCatalog(ctx context.Context, query GetMatchCatalog) (*domain.MatchCatalog, error) {
	return h.catalog.GetMatchCatalog(ctx, query.Filters)
}