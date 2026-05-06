package queries

import (
	"context"
	"middleman/ordering/internal/domain"
)

type GetOrdersByStatus struct {
	Status    string
	Page      int64
	PageSize  int64
	SortBy    string
	SortOrder string
}

type GetOrdersByStatusHandler struct {
	catalog domain.CatalogRepository
}

func NewGetOrdersByStatusHandler(catalog domain.CatalogRepository) GetOrdersByStatusHandler {
	return GetOrdersByStatusHandler{catalog: catalog}
}

func (h GetOrdersByStatusHandler) GetOrdersByStatus(ctx context.Context, query GetOrdersByStatus) ([]*domain.OrderCatalog, int64, error) {
	status := domain.ToOrderStatus(query.Status)
	return h.catalog.GetOrdersByStatus(ctx, status, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
