package domain

import "time"

// Order represents an order placed by a customer
type Order struct {
	ID              string     `json:"id"`
	UserCustomerID  string     `json:"user_customer_id"`
	PaymentMethodID string     `json:"payment_method_id"`
	Items           []Item     `json:"items"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ReadyAt         *time.Time `json:"ready_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	RejectedAt      *time.Time `json:"rejected_at,omitempty"`
	ShippedAt       *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt     *time.Time `json:"delivered_at,omitempty"`
}

// Item represents an item within an order
type Item struct {
	UserSellerID   string `json:"user_seller_id"`
	ProductID      string `json:"product_id"`
	UserSellerName string `json:"user_seller_name"`
	ProductName    string `json:"product_name"`
	Price          int64  `json:"price"`
	Quantity       int64  `json:"quantity"`
}

// CreateOrderResponse represents the response for order creation
type CreateOrderResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// GetOrderResponse represents the response for getting an order
type GetOrderResponse struct {
	Order Order `json:"order"`
}

// ReadyOrderResponse represents the response for marking order as ready
type ReadyOrderResponse struct {
	ID      string    `json:"id"`
	Status  string    `json:"status"`
	ReadyAt time.Time `json:"ready_at"`
}

// CompleteOrderResponse represents the response for completing an order
type CompleteOrderResponse struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	CompletedAt time.Time `json:"completed_at"`
}

// ApproveOrderResponse represents the response for approving an order
type ApproveOrderResponse struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	ApprovedAt time.Time `json:"approved_at"`
}

// RejectOrderResponse represents the response for rejecting an order
type RejectOrderResponse struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	RejectedAt time.Time `json:"rejected_at"`
}

// ShipOrderResponse represents the response for shipping an order
type ShipOrderResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	ShippedAt time.Time `json:"shipped_at"`
}

// DeliverOrderResponse represents the response for delivering an order
type DeliverOrderResponse struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	DeliveredAt time.Time `json:"delivered_at"`
}

// UpdateOrderResponse represents the response for updating an order
type UpdateOrderResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Order status constants
const (
	OrderStatusPending   = "pending"
	OrderStatusCreated   = "created"
	OrderStatusReady     = "ready"
	OrderStatusApproved  = "approved"
	OrderStatusRejected  = "rejected"
	OrderStatusShipped   = "shipped"
	OrderStatusDelivered = "delivered"
	OrderStatusCompleted = "completed"
	OrderStatusCanceled  = "canceled"
)
