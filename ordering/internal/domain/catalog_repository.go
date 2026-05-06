package domain

import "context"

type OrderCatalog struct {
	ID              string
	UserCustomerID  string
	PaymentMethodID string
	Status          OrderStatus
	Total           int64
	ItemCount       int
	CreatedAt       string
	UpdatedAt       string
}

type CatalogRepository interface {
	AddOrder(ctx context.Context, order *Order) error
	UpdateOrder(ctx context.Context, order *Order) error
	RemoveOrder(ctx context.Context, orderID string) error
	Find(ctx context.Context, orderID string) (*OrderCatalog, error)
	ListOrders(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*OrderCatalog, int64, error)
	GetOrdersByCustomer(ctx context.Context, userCustomerID string, page, pageSize int64, sortBy, sortOrder string) ([]*OrderCatalog, int64, error)
	GetOrdersByStatus(ctx context.Context, status OrderStatus, page, pageSize int64, sortBy, sortOrder string) ([]*OrderCatalog, int64, error)
}
