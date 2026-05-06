package grpc

import (
	"context"
	"fmt"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
	"middleman/ordering/orderingpb"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// OrderRepository calls the remote ordering service (gRPC).
type OrderRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.OrderRepository = (*OrderRepository)(nil)

// NewOrderRepositoryWithAuth creates a new OrderRepository with JWT authentication support
func NewOrderRepository(endpoint string, authInstance *auth.Auth) OrderRepository {
	return OrderRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// CreateOrder creates a new order
func (r OrderRepository) CreateOrder(ctx context.Context, orderID string, items []models.Item, userCustomerID string) (*models.CreateOrderResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Convert items to protobuf format
	pbItems := make([]*orderingpb.Item, len(items))
	for i, item := range items {
		pbItems[i] = &orderingpb.Item{
			UserSellerId:   item.UserSellerID,
			ProductId:      item.ProductID,
			UserSellerName: item.UserSellerName,
			ProductName:    item.ProductName,
			Price:          item.Price,
			Quantity:       item.Quantity,
		}
	}

	client := orderingpb.NewOrderingServiceClient(conn)
	resp, err := client.CreateOrder(ctx, &orderingpb.CreateOrderRequest{
		Id:             orderID,
		Items:          pbItems,
		UserCustomerId: userCustomerID,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateOrder RPC failed: %w", err)
	}

	var createdAt time.Time
	if resp.GetCreatedAt() != nil {
		createdAt = resp.GetCreatedAt().AsTime()
	}

	return &models.CreateOrderResponse{
		ID:        resp.GetId(),
		CreatedAt: createdAt,
	}, nil
}

// GetOrder retrieves an order by ID
func (r OrderRepository) GetOrder(ctx context.Context, orderID string) (*models.GetOrderResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := orderingpb.NewOrderingServiceClient(conn)
	resp, err := client.GetOrder(ctx, &orderingpb.GetOrderRequest{
		Id: orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetOrder RPC failed: %w", err)
	}

	order := r.convertOrderFromPb(resp.GetOrder())

	return &models.GetOrderResponse{
		Order: *order,
	}, nil
}

// CancelOrder cancels an existing order
func (r OrderRepository) CancelOrder(ctx context.Context, orderID, reason string) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := orderingpb.NewOrderingServiceClient(conn)
	_, err = client.CancelOrder(ctx, &orderingpb.CancelOrderRequest{
		Id:     orderID,
		Reason: reason,
	})
	if err != nil {
		return fmt.Errorf("CancelOrder RPC failed: %w", err)
	}

	return nil
}

// ReadyOrder marks an order as ready
func (r OrderRepository) ReadyOrder(ctx context.Context, orderID string) (*models.ReadyOrderResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := orderingpb.NewOrderingServiceClient(conn)
	resp, err := client.ReadyOrder(ctx, &orderingpb.ReadyOrderRequest{
		Id: orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("ReadyOrder RPC failed: %w", err)
	}

	var readyAt time.Time
	if resp.GetReadyAt() != nil {
		readyAt = resp.GetReadyAt().AsTime()
	}

	return &models.ReadyOrderResponse{
		ID:      resp.GetId(),
		Status:  resp.GetStatus(),
		ReadyAt: readyAt,
	}, nil
}

// CompleteOrder completes an order
func (r OrderRepository) CompleteOrder(ctx context.Context, orderID, invoiceID string) (*models.CompleteOrderResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := orderingpb.NewOrderingServiceClient(conn)
	resp, err := client.CompleteOrder(ctx, &orderingpb.CompleteOrderRequest{
		Id:        orderID,
		InvoiceId: invoiceID,
	})
	if err != nil {
		return nil, fmt.Errorf("CompleteOrder RPC failed: %w", err)
	}

	var completedAt time.Time
	if resp.GetCompletedAt() != nil {
		completedAt = resp.GetCompletedAt().AsTime()
	}

	return &models.CompleteOrderResponse{
		ID:          resp.GetId(),
		Status:      resp.GetStatus(),
		CompletedAt: completedAt,
	}, nil
}

// ApproveOrder approves an order
func (r OrderRepository) ApproveOrder(ctx context.Context, orderID, shoppingID string) (*models.ApproveOrderResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := orderingpb.NewOrderingServiceClient(conn)
	resp, err := client.ApproveOrder(ctx, &orderingpb.ApproveOrderRequest{
		Id:         orderID,
		ShoppingId: shoppingID,
	})
	if err != nil {
		return nil, fmt.Errorf("ApproveOrder RPC failed: %w", err)
	}

	var approvedAt time.Time
	if resp.GetApprovedAt() != nil {
		approvedAt = resp.GetApprovedAt().AsTime()
	}

	return &models.ApproveOrderResponse{
		ID:         resp.GetId(),
		Status:     resp.GetStatus(),
		ApprovedAt: approvedAt,
	}, nil
}

// RejectOrder rejects an order
func (r OrderRepository) RejectOrder(ctx context.Context, orderID string) (*models.RejectOrderResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := orderingpb.NewOrderingServiceClient(conn)
	resp, err := client.RejectOrder(ctx, &orderingpb.RejectOrderRequest{
		Id: orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("RejectOrder RPC failed: %w", err)
	}

	var rejectedAt time.Time
	if resp.GetRejectedAt() != nil {
		rejectedAt = resp.GetRejectedAt().AsTime()
	}

	return &models.RejectOrderResponse{
		ID:         resp.GetId(),
		Status:     resp.GetStatus(),
		RejectedAt: rejectedAt,
	}, nil
}

// ShipOrder ships an order
func (r OrderRepository) ShipOrder(ctx context.Context, orderID string) (*models.ShipOrderResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := orderingpb.NewOrderingServiceClient(conn)
	resp, err := client.ShipOrder(ctx, &orderingpb.ShipOrderRequest{
		Id: orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("ShipOrder RPC failed: %w", err)
	}

	var shippedAt time.Time
	if resp.GetShippedAt() != nil {
		shippedAt = resp.GetShippedAt().AsTime()
	}

	return &models.ShipOrderResponse{
		ID:        resp.GetId(),
		Status:    resp.GetStatus(),
		ShippedAt: shippedAt,
	}, nil
}

// DeliverOrder delivers an order
func (r OrderRepository) DeliverOrder(ctx context.Context, orderID string) (*models.DeliverOrderResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := orderingpb.NewOrderingServiceClient(conn)
	resp, err := client.DeliverOrder(ctx, &orderingpb.DeliverOrderRequest{
		Id: orderID,
	})
	if err != nil {
		return nil, fmt.Errorf("DeliverOrder RPC failed: %w", err)
	}

	var deliveredAt time.Time
	if resp.GetDeliveredAt() != nil {
		deliveredAt = resp.GetDeliveredAt().AsTime()
	}

	return &models.DeliverOrderResponse{
		ID:          resp.GetId(),
		Status:      resp.GetStatus(),
		DeliveredAt: deliveredAt,
	}, nil
}

// GetOrdersByCustomer retrieves orders by customer ID (mock implementation for AI tooling)
func (r OrderRepository) GetOrdersByCustomer(ctx context.Context, userCustomerID string) ([]*models.Order, error) {
	// Note: This would typically require a GetOrdersByCustomer RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetOrdersByCustomer called for customer: %s (mock implementation)", userCustomerID)

	return []*models.Order{
		{
			ID:              "order_1",
			UserCustomerID:  userCustomerID,
			PaymentMethodID: "payment_method_1",
			Items: []models.Item{
				{
					UserSellerID:   "seller_1",
					ProductID:      "product_1",
					UserSellerName: "Mock Seller",
					ProductName:    "Mock Product",
					Price:          1000,
					Quantity:       1,
				},
			},
			Status:    models.OrderStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

// GetOrdersByStatus retrieves orders by status (mock implementation for AI tooling)
func (r OrderRepository) GetOrdersByStatus(ctx context.Context, status string, limit int64) ([]*models.Order, error) {
	// Note: This would typically require a GetOrdersByStatus RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetOrdersByStatus called with status: %s, limit: %d (mock implementation)", status, limit)

	orders := make([]*models.Order, 0, limit)
	for i := int64(0); i < limit && i < 5; i++ { // Mock max 5 results
		orders = append(orders, &models.Order{
			ID:              fmt.Sprintf("order_%d", i+1),
			UserCustomerID:  "mock_customer",
			PaymentMethodID: "mock_payment_method",
			Items: []models.Item{
				{
					UserSellerID:   fmt.Sprintf("seller_%d", i+1),
					ProductID:      fmt.Sprintf("product_%d", i+1),
					UserSellerName: fmt.Sprintf("Seller %d", i+1),
					ProductName:    fmt.Sprintf("Product %d", i+1),
					Price:          1000 * (i + 1),
					Quantity:       1,
				},
			},
			Status:    status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return orders, nil
}

// SearchOrders searches orders by query (mock implementation for AI tooling)
func (r OrderRepository) SearchOrders(ctx context.Context, query string, limit int64) ([]*models.Order, error) {
	// Note: This would typically require a SearchOrders RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("SearchOrders called with query: %s, limit: %d (mock implementation)", query, limit)

	orders := make([]*models.Order, 0, limit)
	for i := int64(0); i < limit && i < 3; i++ { // Mock max 3 results
		orders = append(orders, &models.Order{
			ID:              fmt.Sprintf("order_%d", i+1),
			UserCustomerID:  "mock_customer",
			PaymentMethodID: "mock_payment_method",
			Items: []models.Item{
				{
					UserSellerID:   fmt.Sprintf("seller_%d", i+1),
					ProductID:      fmt.Sprintf("product_%d", i+1),
					UserSellerName: fmt.Sprintf("Seller %d", i+1),
					ProductName:    fmt.Sprintf("Product %d", i+1),
					Price:          1000 * (i + 1),
					Quantity:       1,
				},
			},
			Status:    models.OrderStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return orders, nil
}

// convertOrderFromPb converts protobuf Order to domain Order
func (r OrderRepository) convertOrderFromPb(pbOrder *orderingpb.Order) *models.Order {
	if pbOrder == nil {
		return nil
	}

	// Convert items
	items := make([]models.Item, len(pbOrder.GetItems()))
	for i, pbItem := range pbOrder.GetItems() {
		items[i] = models.Item{
			UserSellerID:   pbItem.GetUserSellerId(),
			ProductID:      pbItem.GetProductId(),
			UserSellerName: pbItem.GetUserSellerName(),
			ProductName:    pbItem.GetProductName(),
			Price:          pbItem.GetPrice(),
			Quantity:       pbItem.GetQuantity(),
		}
	}

	// Convert timestamps
	var createdAt, updatedAt time.Time
	if pbOrder.GetCreatedAt() != nil {
		createdAt = pbOrder.GetCreatedAt().AsTime()
	}
	if pbOrder.GetUpdatedAt() != nil {
		updatedAt = pbOrder.GetUpdatedAt().AsTime()
	}

	return &models.Order{
		ID:              pbOrder.GetId(),
		UserCustomerID:  pbOrder.GetUserCustomerId(),
		PaymentMethodID: pbOrder.GetPaymentMethodId(),
		Items:           items,
		Status:          pbOrder.GetStatus(),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

// dial establishes a gRPC connection to the ordering service
// dial sets up a gRPC connection with the microservice endpoint
func (r OrderRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r OrderRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}

// UpdateOrder updates an existing order with the provided fields
func (r OrderRepository) UpdateOrder(ctx context.Context, orderID string, updates map[string]interface{}) (*models.UpdateOrderResponse, error) {
	// For now, return a basic response as this would need to be implemented in the orders service
	// This is a placeholder implementation
	return &models.UpdateOrderResponse{
		ID:        orderID,
		Status:    "updated",
		UpdatedAt: time.Now(),
	}, nil
}
