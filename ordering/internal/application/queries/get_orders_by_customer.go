package queries

import (
	"context"
	"middleman/ordering/internal/domain"
)

type GetOrdersByCustomer struct {
	UserCustomerID string
	Page           int64
	PageSize       int64
	SortBy         string
	SortOrder      string
}

type GetOrdersByCustomerHandler struct {
	catalog domain.CatalogRepository
}

func NewGetOrdersByCustomerHandler(catalog domain.CatalogRepository) GetOrdersByCustomerHandler {
	return GetOrdersByCustomerHandler{catalog: catalog}
}

func (h GetOrdersByCustomerHandler) GetOrdersByCustomer(ctx context.Context, query GetOrdersByCustomer) ([]*domain.OrderCatalog, int64, error) {
	return h.catalog.GetOrdersByCustomer(ctx, query.UserCustomerID, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
