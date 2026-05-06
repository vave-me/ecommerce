package queries

import (
	"context"
	"middleman/managers/internal/domain"
)

type GetManager struct {
	ID string
}

type GetManagerHandler struct {
	managers domain.CatalogRepository
}

func NewGetManagerHandler(managers domain.CatalogRepository) GetManagerHandler {
	return GetManagerHandler{
		managers: managers,
	}
}

func (h GetManagerHandler) GetManager(ctx context.Context, query GetManager) (*domain.CatalogManager, error) {
	return h.managers.Find(ctx, query.ID)
}
