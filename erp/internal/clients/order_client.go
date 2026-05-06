package clients

import (
	"context"
	"fmt"
	"middleman/erp/internal/domain"

	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/ordering/orderingpb"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// OrderClient calls the remote ordering service (gRPC).
type OrderClient struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.OrderRepository = (*OrderClient)(nil)

// NewOrderClient creates a new OrderClient with JWT authentication support
func NewOrderClient(endpoint string, authInstance *auth.Auth) OrderClient {
	return OrderClient{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// CreateOrder creates a new order
func (r OrderClient) CreateOrder(ctx context.Context, orderID string, items []domain.Item, userCustomerID string) (*domain.CreateOrderResponse, error) {
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

	return &domain.CreateOrderResponse{
		ID:        resp.GetId(),
		CreatedAt: createdAt,
	}, nil
}

// GetOrder retrieves an order by ID
func (r OrderClient) GetOrder(ctx context.Context, orderID string) (*domain.GetOrderResponse, error) {
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

	return &domain.GetOrderResponse{
		Order: *order,
	}, nil
}

// CancelOrder cancels an existing order
func (r OrderClient) CancelOrder(ctx context.Context, orderID, reason string) error {
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
func (r OrderClient) ReadyOrder(ctx context.Context, orderID string) (*domain.ReadyOrderResponse, error) {
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

	return &domain.ReadyOrderResponse{
		ID:      resp.GetId(),
		Status:  resp.GetStatus(),
		ReadyAt: readyAt,
	}, nil
}

// CompleteOrder completes an order
func (r OrderClient) CompleteOrder(ctx context.Context, orderID, invoiceID string) (*domain.CompleteOrderResponse, error) {
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

	return &domain.CompleteOrderResponse{
		ID:          resp.GetId(),
		Status:      resp.GetStatus(),
		CompletedAt: completedAt,
	}, nil
}

// ApproveOrder approves an order
func (r OrderClient) ApproveOrder(ctx context.Context, orderID, shoppingID string) (*domain.ApproveOrderResponse, error) {
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

	return &domain.ApproveOrderResponse{
		ID:         resp.GetId(),
		Status:     resp.GetStatus(),
		ApprovedAt: approvedAt,
	}, nil
}

// RejectOrder rejects an order
func (r OrderClient) RejectOrder(ctx context.Context, orderID string) (*domain.RejectOrderResponse, error) {
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

	return &domain.RejectOrderResponse{
		ID:         resp.GetId(),
		Status:     resp.GetStatus(),
		RejectedAt: rejectedAt,
	}, nil
}

// ShipOrder ships an order
func (r OrderClient) ShipOrder(ctx context.Context, orderID string) (*domain.ShipOrderResponse, error) {
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

	return &domain.ShipOrderResponse{
		ID:        resp.GetId(),
		Status:    resp.GetStatus(),
		ShippedAt: shippedAt,
	}, nil
}

// DeliverOrder delivers an order
func (r OrderClient) DeliverOrder(ctx context.Context, orderID string) (*domain.DeliverOrderResponse, error) {
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

	return &domain.DeliverOrderResponse{
		ID:          resp.GetId(),
		Status:      resp.GetStatus(),
		DeliveredAt: deliveredAt,
	}, nil
}

// GetOrdersByCustomer retrieves orders by customer ID (mock implementation for AI tooling)
func (r OrderClient) GetOrdersByCustomer(ctx context.Context, userCustomerID string) ([]*domain.Order, error) {
	// Note: This would typically require a GetOrdersByCustomer RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetOrdersByCustomer called for customer: %s (mock implementation)", userCustomerID)

	return []*domain.Order{
		{
			ID:              "order_1",
			UserCustomerID:  userCustomerID,
			PaymentMethodID: "payment_method_1",
			Items: []domain.Item{
				{
					UserSellerID:   "seller_1",
					ProductID:      "product_1",
					UserSellerName: "Mock Seller",
					ProductName:    "Mock Product",
					Price:          1000,
					Quantity:       1,
				},
			},
			Status:    domain.OrderStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

// GetOrdersByStatus retrieves orders by status (mock implementation for AI tooling)
func (r OrderClient) GetOrdersByStatus(ctx context.Context, status string, limit int64) ([]*domain.Order, error) {
	// Note: This would typically require a GetOrdersByStatus RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetOrdersByStatus called with status: %s, limit: %d (mock implementation)", status, limit)

	orders := make([]*domain.Order, 0, limit)
	for i := int64(0); i < limit && i < 5; i++ { // Mock max 5 results
		orders = append(orders, &domain.Order{
			ID:              fmt.Sprintf("order_%d", i+1),
			UserCustomerID:  "mock_customer",
			PaymentMethodID: "mock_payment_method",
			Items: []domain.Item{
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
func (r OrderClient) SearchOrders(ctx context.Context, query string, limit int64) ([]*domain.Order, error) {
	// Note: This would typically require a SearchOrders RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("SearchOrders called with query: %s, limit: %d (mock implementation)", query, limit)

	orders := make([]*domain.Order, 0, limit)
	for i := int64(0); i < limit && i < 3; i++ { // Mock max 3 results
		orders = append(orders, &domain.Order{
			ID:              fmt.Sprintf("order_%d", i+1),
			UserCustomerID:  "mock_customer",
			PaymentMethodID: "mock_payment_method",
			Items: []domain.Item{
				{
					UserSellerID:   fmt.Sprintf("seller_%d", i+1),
					ProductID:      fmt.Sprintf("product_%d", i+1),
					UserSellerName: fmt.Sprintf("Seller %d", i+1),
					ProductName:    fmt.Sprintf("Product %d", i+1),
					Price:          1000 * (i + 1),
					Quantity:       1,
				},
			},
			Status:    domain.OrderStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	return orders, nil
}

// convertOrderFromPb converts protobuf Order to domain Order
func (r OrderClient) convertOrderFromPb(pbOrder *orderingpb.Order) *domain.Order {
	if pbOrder == nil {
		return nil
	}

	// Convert items
	items := make([]domain.Item, len(pbOrder.GetItems()))
	for i, pbItem := range pbOrder.GetItems() {
		items[i] = domain.Item{
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

	return &domain.Order{
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
func (r OrderClient) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r OrderClient) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}

// UpdateOrder updates an existing order with the provided fields
func (r OrderClient) UpdateOrder(ctx context.Context, orderID string, updates map[string]interface{}) (*domain.UpdateOrderResponse, error) {
	// For now, return a basic response as this would need to be implemented in the orders service
	// This is a placeholder implementation
	return &domain.UpdateOrderResponse{
		ID:        orderID,
		Status:    "updated",
		UpdatedAt: time.Now(),
	}, nil
}
