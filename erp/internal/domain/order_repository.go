package domain

import (
	"context"
)

type OrderRepository interface {
	// Order lifecycle operations
	CreateOrder(ctx context.Context, orderID string, items []Item, userCustomerID string) (*CreateOrderResponse, error)
	GetOrder(ctx context.Context, orderID string) (*GetOrderResponse, error)
	CancelOrder(ctx context.Context, orderID, reason string) error

	// Order status management
	ReadyOrder(ctx context.Context, orderID string) (*ReadyOrderResponse, error)
	CompleteOrder(ctx context.Context, orderID, invoiceID string) (*CompleteOrderResponse, error)
	ApproveOrder(ctx context.Context, orderID, shoppingID string) (*ApproveOrderResponse, error)
	RejectOrder(ctx context.Context, orderID string) (*RejectOrderResponse, error)
	ShipOrder(ctx context.Context, orderID string) (*ShipOrderResponse, error)
	DeliverOrder(ctx context.Context, orderID string) (*DeliverOrderResponse, error)

	// Additional query methods for AI tooling
	GetOrdersByCustomer(ctx context.Context, userCustomerID string) ([]*Order, error)
	GetOrdersByStatus(ctx context.Context, status string, limit int64) ([]*Order, error)
	SearchOrders(ctx context.Context, query string, limit int64) ([]*Order, error)

	// Update operations
	UpdateOrder(ctx context.Context, orderID string, updates map[string]interface{}) (*UpdateOrderResponse, error)
}
