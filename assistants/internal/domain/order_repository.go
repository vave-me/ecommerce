package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type OrderRepository interface {
	// Order lifecycle operations
	CreateNewCustomerOrder(ctx context.Context, orderID string, items []models.Item, userCustomerID string) (*models.CreateOrderResponse, error)
	GetOrderDetailsByID(ctx context.Context, orderID string) (*models.GetOrderResponse, error)
	CancelOrderWithReason(ctx context.Context, orderID, reason string) error

	// Order status management
	MarkOrderAsReadyForProcessing(ctx context.Context, orderID string) (*models.ReadyOrderResponse, error)
	CompleteOrderWithInvoice(ctx context.Context, orderID, invoiceID string) (*models.CompleteOrderResponse, error)
	ApproveOrderForShopping(ctx context.Context, orderID, shoppingID string) (*models.ApproveOrderResponse, error)
	RejectOrderRequest(ctx context.Context, orderID string) (*models.RejectOrderResponse, error)
	MarkOrderAsShipped(ctx context.Context, orderID string) (*models.ShipOrderResponse, error)
	MarkOrderAsDelivered(ctx context.Context, orderID string) (*models.DeliverOrderResponse, error)

	// Additional query methods for AI tooling
	GetCustomerOrderHistory(ctx context.Context, userCustomerID string) ([]*models.Order, error)
	FilterOrdersByCurrentStatus(ctx context.Context, status string, limit int64) ([]*models.Order, error)
	SearchOrdersByKeyword(ctx context.Context, query string, limit int64) ([]*models.Order, error)

	// Update operations
	UpdateOrderInformation(ctx context.Context, orderID string, updates map[string]interface{}) (*models.UpdateOrderResponse, error)
}
