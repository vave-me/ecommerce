package models

import (
	"time"
)

// Basket represents a shopping basket containing multiple items
type Basket struct {
	ID           string        `json:"id"`
	UserID       string        `json:"user_id"`
	Items        []*BasketItem `json:"items"`
	BasketStatus string        `json:"basket_status"`
	TotalAmount  int64         `json:"total_amount"`
	CreatedAt    time.Time     `json:"created_at,omitempty"`
	UpdatedAt    time.Time     `json:"updated_at,omitempty"`
}

// BasketItem represents a product added to the basket
type BasketItem struct {
	ID             string    `json:"id"`
	UserSellerID   string    `json:"user_seller_id"`
	ProductID      string    `json:"product_id"`
	UserSellerName string    `json:"user_seller_name"`
	ProductName    string    `json:"product_name"`
	ProductPrice   int64     `json:"product_price"`
	Quantity       int64     `json:"quantity"`
	BasketID       string    `json:"basket_id"`
	AddedAt        time.Time `json:"added_at,omitempty"`
}

// BasketStatus constants
const (
	BasketStatusActive     = "ACTIVE"
	BasketStatusCanceled   = "CANCELED"
	BasketStatusCheckedOut = "CHECKED_OUT"
)

// StartBasketRequest represents request to start a new basket
type StartBasketRequest struct {
	UserID string `json:"user_id"`
}

// StartBasketResponse represents response after starting a new basket
type StartBasketResponse struct {
	ID string `json:"id"`
}

// GetBasketRequest represents request to retrieve a basket
type GetBasketRequest struct {
	BasketID string `json:"basket_id"`
}

// GetBasketResponse represents response for retrieving a basket
type GetBasketResponse struct {
	Basket *Basket `json:"basket"`
}

// GetTotalBasketAmountRequest represents request to get total basket amount
type GetTotalBasketAmountRequest struct {
	BasketID string `json:"basket_id"`
}

// GetTotalBasketAmountResponse represents response for total basket amount
type GetTotalBasketAmountResponse struct {
	Amount int64 `json:"amount"`
}

// GetCurrentBasketRequest represents request to get current basket
type GetCurrentBasketRequest struct {
	UserCustomerID string `json:"user_customer_id"`
}

// GetCurrentBasketResponse represents response for current basket
type GetCurrentBasketResponse struct {
	BasketID     string `json:"basket_id"`
	BasketStatus string `json:"basket_status"`
}

// CancelBasketRequest represents request to cancel a basket
type CancelBasketRequest struct {
	BasketID string `json:"basket_id"`
	Reason   string `json:"reason"`
}

// CancelBasketResponse represents response after canceling a basket
type CancelBasketResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CheckoutBasketRequest represents request to checkout a basket
type CheckoutBasketRequest struct {
	BasketID       string `json:"basket_id"`
	UserCustomerID string `json:"user_customer_id"`
}

// CheckoutBasketResponse represents response after checking out a basket
type CheckoutBasketResponse struct {
	BasketID string `json:"basket_id"`
	OrderID  string `json:"order_id"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

// AddItemRequest represents request to add an item to the basket
type AddItemRequest struct {
	BasketID  string `json:"basket_id"`
	ProductID string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
}

// AddItemResponse represents response after adding an item
type AddItemResponse struct {
	BasketID string `json:"basket_id"`
	ItemID   string `json:"item_id"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

// RemoveItemRequest represents request to remove an item from the basket
type RemoveItemRequest struct {
	BasketID string `json:"basket_id"`
	ItemID   string `json:"item_id"`
	Quantity int64  `json:"quantity"`
}

// RemoveItemResponse represents response after removing an item
type RemoveItemResponse struct {
	BasketID string `json:"basket_id"`
	ItemID   string `json:"item_id"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

// UpdateItemQuantityRequest represents request to update item quantity
type UpdateItemQuantityRequest struct {
	BasketID    string `json:"basket_id"`
	ItemID      string `json:"item_id"`
	NewQuantity int64  `json:"new_quantity"`
}

// UpdateItemQuantityResponse represents response after updating item quantity
type UpdateItemQuantityResponse struct {
	BasketID    string `json:"basket_id"`
	ItemID      string `json:"item_id"`
	NewQuantity int64  `json:"new_quantity"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
}

// ListBasketsRequest represents request to list baskets with filtering and pagination
type ListBasketsRequest struct {
	UserID       string `json:"user_id"`
	BasketStatus string `json:"basket_status"`
	Page         int64  `json:"page"`
	Limit        int64  `json:"limit"`
}

// ListBasketsResponse represents response for listing baskets
type ListBasketsResponse struct {
	Baskets []*Basket `json:"baskets"`
	Total   int64     `json:"total"`
	Page    int64     `json:"page"`
	Limit   int64     `json:"limit"`
}

// GetActiveBasketsRequest represents request to retrieve all active baskets
type GetActiveBasketsRequest struct {
	// Empty request
}

// GetActiveBasketsResponse represents response for retrieving active baskets
type GetActiveBasketsResponse struct {
	Baskets []*Basket `json:"baskets"`
}

// Additional response types for AI tooling methods

// BasketStatsResponse represents basket statistics for a user
type BasketStatsResponse struct {
	TotalBaskets       int64   `json:"total_baskets"`
	ActiveBaskets      int64   `json:"active_baskets"`
	CanceledBaskets    int64   `json:"canceled_baskets"`
	CheckedOutBaskets  int64   `json:"checked_out_baskets"`
	TotalValue         int64   `json:"total_value"`
	AverageBasketValue float64 `json:"average_basket_value"`
}

// BasketValidationResponse represents validation result for a basket
type BasketValidationResponse struct {
	IsValid     bool     `json:"is_valid"`
	Errors      []string `json:"errors"`
	Warnings    []string `json:"warnings"`
	TotalAmount int64    `json:"total_amount"`
	ItemsCount  int64    `json:"items_count"`
}

// BasketAnalyticsResponse represents analytics data for baskets
type BasketAnalyticsResponse struct {
	Period            string                 `json:"period"`
	TotalBaskets      int64                  `json:"total_baskets"`
	ConversionRate    float64                `json:"conversion_rate"`
	AbandonmentRate   float64                `json:"abandonment_rate"`
	AverageValue      float64                `json:"average_value"`
	PopularProducts   []*ProductStats        `json:"popular_products"`
	UserBehaviorStats map[string]interface{} `json:"user_behavior_stats"`
}

// BasketConversionStatsResponse represents conversion statistics
type BasketConversionStatsResponse struct {
	Period                string  `json:"period"`
	TotalBaskets          int64   `json:"total_baskets"`
	CheckedOutBaskets     int64   `json:"checked_out_baskets"`
	AbandonedBaskets      int64   `json:"abandoned_baskets"`
	ConversionRate        float64 `json:"conversion_rate"`
	AbandonmentRate       float64 `json:"abandonment_rate"`
	AverageTimeToCheckout float64 `json:"average_time_to_checkout"`
}

// ProductStats represents statistics for products in baskets
type ProductStats struct {
	ProductID      string  `json:"product_id"`
	ProductName    string  `json:"product_name"`
	TimesAdded     int64   `json:"times_added"`
	TotalQuantity  int64   `json:"total_quantity"`
	TotalValue     int64   `json:"total_value"`
	ConversionRate float64 `json:"conversion_rate"`
}

// BasketFilter represents filtering options for basket queries
type BasketFilter struct {
	UserID      string    `json:"user_id"`
	Status      string    `json:"status"`
	MinAmount   int64     `json:"min_amount"`
	MaxAmount   int64     `json:"max_amount"`
	CreatedFrom time.Time `json:"created_from"`
	CreatedTo   time.Time `json:"created_to"`
	HasItems    bool      `json:"has_items"`
	ProductIDs  []string  `json:"product_ids"`
}

// BasketSortOptions represents sorting options for basket queries
type BasketSortOptions struct {
	SortBy    string `json:"sort_by"`    // created_at, updated_at, total_amount, items_count
	SortOrder string `json:"sort_order"` // asc, desc
}

// ItemFilter represents filtering options for basket items
type ItemFilter struct {
	ProductID    string    `json:"product_id"`
	UserSellerID string    `json:"user_seller_id"`
	MinPrice     int64     `json:"min_price"`
	MaxPrice     int64     `json:"max_price"`
	MinQuantity  int64     `json:"min_quantity"`
	MaxQuantity  int64     `json:"max_quantity"`
	AddedFrom    time.Time `json:"added_from"`
	AddedTo      time.Time `json:"added_to"`
}
