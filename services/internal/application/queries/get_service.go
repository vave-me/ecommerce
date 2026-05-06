package queries

import (
	"context"
	"middleman/services/internal/domain"
)

type GetService struct {
	ID string
}

type GetServiceHandler struct {
	catalog domain.CatalogRepository
}

func NewGetServiceHandler(catalog domain.CatalogRepository) GetServiceHandler {
	return GetServiceHandler{catalog: catalog}
}

func (h GetServiceHandler) GetService(ctx context.Context, query GetService) (*domain.CatalogService, error) {
	return h.catalog.Find(ctx, query.ID)
}
