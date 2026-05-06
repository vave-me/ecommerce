package queries

import (
	"context"
	"middleman/streams/internal/domain"
)

type GetCatalog struct {
	UserID   string // Optional - for personalized catalog
	Page     int
	PageSize int
	Filters  domain.CatalogFilters
}

type GetCatalogHandler struct {
	catalog domain.StreamCatalogRepository
}

func NewGetCatalogHandler(catalog domain.StreamCatalogRepository) GetCatalogHandler {
	return GetCatalogHandler{catalog: catalog}
}

func (h GetCatalogHandler) GetCatalog(ctx context.Context, query GetCatalog) (*domain.StreamCatalog, error) {
	if query.UserID != "" {
		return h.catalog.GetUserCatalog(ctx, query.UserID, query.Filters)
	}
	return h.catalog.GetCatalog(ctx, query.Filters)
}