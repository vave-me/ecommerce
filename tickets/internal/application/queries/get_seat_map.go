package queries

import (
	"context"
	"middleman/tickets/internal/domain"
)

type GetSeatMap struct {
	MatchID  string
	SectorID string
}

type GetSeatMapHandler struct {
	catalog domain.MatchCatalogRepository
}

func NewGetSeatMapHandler(catalog domain.MatchCatalogRepository) GetSeatMapHandler {
	return GetSeatMapHandler{catalog: catalog}
}

func (h GetSeatMapHandler) GetSeatMap(ctx context.Context, query GetSeatMap) (*domain.SectorSeatMap, error) {
	return h.catalog.GetSeatMap(ctx, query.MatchID, query.SectorID)
}