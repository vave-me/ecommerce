package queries

import (
	"context"
	"middleman/ordering/internal/domain"
)

type ListOrders struct {
	Page      int64
	PageSize  int64
	SortBy    string
	SortOrder string
}

type ListOrdersHandler struct {
	catalog domain.CatalogRepository
}

func NewListOrdersHandler(catalog domain.CatalogRepository) ListOrdersHandler {
	return ListOrdersHandler{catalog: catalog}
}

func (h ListOrdersHandler) ListOrders(ctx context.Context, query ListOrders) ([]*domain.OrderCatalog, int64, error) {
	return h.catalog.ListOrders(ctx, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
