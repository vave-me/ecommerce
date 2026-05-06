package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type OrderRepository interface {
	// Order lifecycle operations
	CreateOrder(ctx context.Context, orderID string, items []models.Item, userCustomerID string) (*models.CreateOrderResponse, error)
	GetOrder(ctx context.Context, orderID string) (*models.GetOrderResponse, error)
	CancelOrder(ctx context.Context, orderID, reason string) error

	// Order status management
	ReadyOrder(ctx context.Context, orderID string) (*models.ReadyOrderResponse, error)
	CompleteOrder(ctx context.Context, orderID, invoiceID string) (*models.CompleteOrderResponse, error)
	ApproveOrder(ctx context.Context, orderID, shoppingID string) (*models.ApproveOrderResponse, error)
	RejectOrder(ctx context.Context, orderID string) (*models.RejectOrderResponse, error)
	ShipOrder(ctx context.Context, orderID string) (*models.ShipOrderResponse, error)
	DeliverOrder(ctx context.Context, orderID string) (*models.DeliverOrderResponse, error)

	// Additional query methods for AI tooling
	GetOrdersByCustomer(ctx context.Context, userCustomerID string) ([]*models.Order, error)
	GetOrdersByStatus(ctx context.Context, status string, limit int64) ([]*models.Order, error)
	SearchOrders(ctx context.Context, query string, limit int64) ([]*models.Order, error)

	// Update operations
	UpdateOrder(ctx context.Context, orderID string, updates map[string]interface{}) (*models.UpdateOrderResponse, error)
}
