package grpc

import (
	"context"
	"fmt"
	"log"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
	"time"

	"google.golang.org/grpc"
)

// BasketRepository calls the remote basket service (gRPC) as a fallback.
type BasketRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.BasketRepository = (*BasketRepository)(nil)

// NewBasketRepository instantiates the gRPC-based fallback repo.
func NewBasketRepository(endpoint string) BasketRepository {
	return BasketRepository{
		endpoint: endpoint,
		auth:     nil, // No auth by default for backwards compatibility
	}
}

// Core protobuf methods

func (r BasketRepository) StartBasket(ctx context.Context, userID string) (*models.StartBasketResponse, error) {
	log.Printf("[BASKET_GRPC] StartBasket: starting basket for user=%s via gRPC", userID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Printf("[BASKET_GRPC] StartBasket: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	// Mock implementation since basketspb is not available yet
	// TODO: Replace with actual protobuf client when basketspb is generated
	log.Printf("[BASKET_GRPC] StartBasket: Mock implementation - would call basketspb.NewBasketServiceClient")

	// Mock response
	mockBasketID := fmt.Sprintf("basket_%s_%d", userID, time.Now().Unix())

	return &models.StartBasketResponse{
		ID: mockBasketID,
	}, nil
}

func (r BasketRepository) GetBasket(ctx context.Context, basketID string) (*models.GetBasketResponse, error) {
	log.Printf("[BASKET_GRPC] GetBasket: retrieving basket ID=%s via gRPC", basketID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Mock implementation
	log.Printf("[BASKET_GRPC] GetBasket: Mock implementation - would call basketspb.GetBasket")

	mockBasket := &models.Basket{
		ID:           basketID,
		UserID:       "mock_user",
		Items:        []*models.BasketItem{},
		BasketStatus: models.BasketStatusActive,
		TotalAmount:  0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return &models.GetBasketResponse{
		Basket: mockBasket,
	}, nil
}

func (r BasketRepository) GetTotalBasket(ctx context.Context, basketID string) (*models.GetTotalBasketAmountResponse, error) {
	log.Printf("[BASKET_GRPC] GetTotalBasket: getting total for basket ID=%s via gRPC", basketID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Mock implementation
	log.Printf("[BASKET_GRPC] GetTotalBasket: Mock implementation - would call basketspb.GetTotalBasket")

	return &models.GetTotalBasketAmountResponse{
		Amount: 0, // Mock total amount
	}, nil
}

func (r BasketRepository) GetCurrentBasket(ctx context.Context, userCustomerID string) (*models.GetCurrentBasketResponse, error) {
	log.Printf("[BASKET_GRPC] GetCurrentBasket: getting current basket for user=%s via gRPC", userCustomerID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Mock implementation
	log.Printf("[BASKET_GRPC] GetCurrentBasket: Mock implementation - would call basketspb.GetCurrentBasket")

	mockBasketID := fmt.Sprintf("current_basket_%s", userCustomerID)

	return &models.GetCurrentBasketResponse{
		BasketID:     mockBasketID,
		BasketStatus: models.BasketStatusActive,
	}, nil
}

func (r BasketRepository) CancelBasket(ctx context.Context, basketID, reason string) (*models.CancelBasketResponse, error) {
	log.Printf("[BASKET_GRPC] CancelBasket: canceling basket ID=%s with reason=%s via gRPC", basketID, reason)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Mock implementation
	log.Printf("[BASKET_GRPC] CancelBasket: Mock implementation - would call basketspb.CancelBasket")

	return &models.CancelBasketResponse{
		Success: true,
		Message: "Basket canceled successfully (mock implementation)",
	}, nil
}

func (r BasketRepository) CheckoutBasket(ctx context.Context, basketID, userCustomerID string) (*models.CheckoutBasketResponse, error) {
	log.Printf("[BASKET_GRPC] CheckoutBasket: checking out basket ID=%s for user=%s via gRPC", basketID, userCustomerID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Mock implementation
	log.Printf("[BASKET_GRPC] CheckoutBasket: Mock implementation - would call basketspb.CheckoutBasket")

	mockOrderID := fmt.Sprintf("order_%s_%d", basketID, time.Now().Unix())

	return &models.CheckoutBasketResponse{
		BasketID: basketID,
		OrderID:  mockOrderID,
		Success:  true,
		Message:  "Basket checked out successfully (mock implementation)",
	}, nil
}

func (r BasketRepository) AddItem(ctx context.Context, basketID, productID string, quantity int64) (*models.AddItemResponse, error) {
	log.Printf("[BASKET_GRPC] AddItem: adding item productID=%s (qty=%d) to basket=%s via gRPC", productID, quantity, basketID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Mock implementation
	log.Printf("[BASKET_GRPC] AddItem: Mock implementation - would call basketspb.AddItem")

	mockItemID := fmt.Sprintf("item_%s_%s", basketID, productID)

	return &models.AddItemResponse{
		BasketID: basketID,
		ItemID:   mockItemID,
		Success:  true,
		Message:  "Item added successfully (mock implementation)",
	}, nil
}

func (r BasketRepository) RemoveItem(ctx context.Context, basketID, itemID string, quantity int64) (*models.RemoveItemResponse, error) {
	log.Printf("[BASKET_GRPC] RemoveItem: removing item ID=%s (qty=%d) from basket=%s via gRPC", itemID, quantity, basketID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Mock implementation
	log.Printf("[BASKET_GRPC] RemoveItem: Mock implementation - would call basketspb.RemoveItem")

	return &models.RemoveItemResponse{
		BasketID: basketID,
		ItemID:   itemID,
		Success:  true,
		Message:  "Item removed successfully (mock implementation)",
	}, nil
}

func (r BasketRepository) UpdateItemQuantity(ctx context.Context, basketID, itemID string, newQuantity int64) (*models.UpdateItemQuantityResponse, error) {
	log.Printf("[BASKET_GRPC] UpdateItemQuantity: updating item ID=%s to qty=%d in basket=%s via gRPC", itemID, newQuantity, basketID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Mock implementation
	log.Printf("[BASKET_GRPC] UpdateItemQuantity: Mock implementation - would call basketspb.UpdateItemQuantity")

	return &models.UpdateItemQuantityResponse{
		BasketID:    basketID,
		ItemID:      itemID,
		NewQuantity: newQuantity,
		Success:     true,
		Message:     "Item quantity updated successfully (mock implementation)",
	}, nil
}

func (r BasketRepository) ListBaskets(ctx context.Context, userID, basketStatus string, page, limit int64) (*models.ListBasketsResponse, error) {
	log.Printf("[BASKET_GRPC] ListBaskets: listing baskets for user=%s, status=%s, page=%d, limit=%d via gRPC", userID, basketStatus, page, limit)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Mock implementation
	log.Printf("[BASKET_GRPC] ListBaskets: Mock implementation - would call basketspb.ListBaskets")

	// Mock baskets
	mockBaskets := []*models.Basket{
		{
			ID:           "basket_1",
			UserID:       userID,
			Items:        []*models.BasketItem{},
			BasketStatus: models.BasketStatusActive,
			TotalAmount:  0,
			CreatedAt:    time.Now().Add(-time.Hour),
			UpdatedAt:    time.Now().Add(-time.Hour),
		},
	}

	return &models.ListBasketsResponse{
		Baskets: mockBaskets,
		Total:   int64(len(mockBaskets)),
		Page:    page,
		Limit:   limit,
	}, nil
}

func (r BasketRepository) GetActiveBaskets(ctx context.Context) (*models.GetActiveBasketsResponse, error) {
	log.Printf("[BASKET_GRPC] GetActiveBaskets: retrieving all active baskets via gRPC")

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Mock implementation
	log.Printf("[BASKET_GRPC] GetActiveBaskets: Mock implementation - would call basketspb.GetActiveBaskets")

	// Mock active baskets
	mockBaskets := []*models.Basket{
		{
			ID:           "active_basket_1",
			UserID:       "user_1",
			Items:        []*models.BasketItem{},
			BasketStatus: models.BasketStatusActive,
			TotalAmount:  0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	return &models.GetActiveBasketsResponse{
		Baskets: mockBaskets,
	}, nil
}

// Additional AI tooling methods (mock implementations)

func (r BasketRepository) GetBasketsWithPagination(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Basket, error) {
	log.Printf("[BASKET_GRPC] GetBasketsWithPagination: Mock implementation")
	return []*models.Basket{}, nil
}

func (r BasketRepository) SearchBasketsWithTerm(ctx context.Context, term string) ([]*models.Basket, error) {
	log.Printf("[BASKET_GRPC] SearchBasketsWithTerm: Mock implementation for term=%s", term)
	return []*models.Basket{}, nil
}

func (r BasketRepository) GetBasketsByStatus(ctx context.Context, status string, page, pageSize int64) ([]*models.Basket, error) {
	log.Printf("[BASKET_GRPC] GetBasketsByStatus: Mock implementation for status=%s", status)
	return []*models.Basket{}, nil
}

func (r BasketRepository) GetUserBaskets(ctx context.Context, userID string, page, pageSize int64) ([]*models.Basket, error) {
	log.Printf("[BASKET_GRPC] GetUserBaskets: Mock implementation for user=%s", userID)
	return []*models.Basket{}, nil
}

func (r BasketRepository) GetBasketStats(ctx context.Context, userID string) (*models.BasketStatsResponse, error) {
	log.Printf("[BASKET_GRPC] GetBasketStats: Mock implementation for user=%s", userID)
	return &models.BasketStatsResponse{
		TotalBaskets:       0,
		ActiveBaskets:      0,
		CanceledBaskets:    0,
		CheckedOutBaskets:  0,
		TotalValue:         0,
		AverageBasketValue: 0.0,
	}, nil
}

func (r BasketRepository) GetExpiredBaskets(ctx context.Context, page, pageSize int64) ([]*models.Basket, error) {
	log.Printf("[BASKET_GRPC] GetExpiredBaskets: Mock implementation")
	return []*models.Basket{}, nil
}

func (r BasketRepository) GetBasketsByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int64) ([]*models.Basket, error) {
	log.Printf("[BASKET_GRPC] GetBasketsByDateRange: Mock implementation")
	return []*models.Basket{}, nil
}

func (r BasketRepository) GetBasketItems(ctx context.Context, basketID string) ([]*models.BasketItem, error) {
	log.Printf("[BASKET_GRPC] GetBasketItems: Mock implementation for basket=%s", basketID)
	return []*models.BasketItem{}, nil
}

func (r BasketRepository) GetItemsCount(ctx context.Context, basketID string) (int64, error) {
	log.Printf("[BASKET_GRPC] GetItemsCount: Mock implementation for basket=%s", basketID)
	return 0, nil
}

func (r BasketRepository) ValidateBasket(ctx context.Context, basketID string) (*models.BasketValidationResponse, error) {
	log.Printf("[BASKET_GRPC] ValidateBasket: Mock implementation for basket=%s", basketID)
	return &models.BasketValidationResponse{
		IsValid:     true,
		Errors:      []string{},
		Warnings:    []string{},
		TotalAmount: 0,
		ItemsCount:  0,
	}, nil
}

func (r BasketRepository) ClearBasket(ctx context.Context, basketID string) error {
	log.Printf("[BASKET_GRPC] ClearBasket: Mock implementation for basket=%s", basketID)
	return nil
}

func (r BasketRepository) GetBasketAnalytics(ctx context.Context, userID string, timeRange string) (*models.BasketAnalyticsResponse, error) {
	log.Printf("[BASKET_GRPC] GetBasketAnalytics: Mock implementation for user=%s", userID)
	return &models.BasketAnalyticsResponse{
		Period:            timeRange,
		TotalBaskets:      0,
		ConversionRate:    0.0,
		AbandonmentRate:   0.0,
		AverageValue:      0.0,
		PopularProducts:   []*models.ProductStats{},
		UserBehaviorStats: make(map[string]interface{}),
	}, nil
}

func (r BasketRepository) GetAbandonedBaskets(ctx context.Context, page, pageSize int64) ([]*models.Basket, error) {
	log.Printf("[BASKET_GRPC] GetAbandonedBaskets: Mock implementation")
	return []*models.Basket{}, nil
}

func (r BasketRepository) GetBasketConversionStats(ctx context.Context, dateRange string) (*models.BasketConversionStatsResponse, error) {
	log.Printf("[BASKET_GRPC] GetBasketConversionStats: Mock implementation")
	return &models.BasketConversionStatsResponse{
		Period:                dateRange,
		TotalBaskets:          0,
		CheckedOutBaskets:     0,
		AbandonedBaskets:      0,
		ConversionRate:        0.0,
		AbandonmentRate:       0.0,
		AverageTimeToCheckout: 0.0,
	}, nil
}

// Legacy methods for backward compatibility

func (r BasketRepository) FindBasket(ctx context.Context, basketID string) (*models.Basket, error) {
	log.Printf("[BASKET_GRPC] FindBasket: Mock implementation for basket=%s", basketID)
	return &models.Basket{
		ID:           basketID,
		UserID:       "mock_user",
		Items:        []*models.BasketItem{},
		BasketStatus: models.BasketStatusActive,
		TotalAmount:  0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (r BasketRepository) GetAllBaskets(ctx context.Context, userID string) ([]*models.Basket, error) {
	log.Printf("[BASKET_GRPC] GetAllBaskets: Mock implementation for user=%s", userID)
	return []*models.Basket{}, nil
}

// Helper methods

// basketToDomain converts a protobuf Basket to domain model (placeholder)
func (r BasketRepository) basketToDomain(pb interface{}) *models.Basket {
	// TODO: Implement conversion when basketspb is available
	log.Printf("[BASKET_GRPC] basketToDomain: Mock conversion")
	return &models.Basket{
		ID:           "mock_basket",
		UserID:       "mock_user",
		Items:        []*models.BasketItem{},
		BasketStatus: models.BasketStatusActive,
		TotalAmount:  0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// basketItemToDomain converts a protobuf BasketItem to domain model (placeholder)
func (r BasketRepository) basketItemToDomain(pb interface{}) *models.BasketItem {
	// TODO: Implement conversion when basketspb is available
	log.Printf("[BASKET_GRPC] basketItemToDomain: Mock conversion")
	return &models.BasketItem{
		ID:             "mock_item",
		UserSellerID:   "mock_seller",
		ProductID:      "mock_product",
		UserSellerName: "Mock Seller",
		ProductName:    "Mock Product",
		ProductPrice:   1000,
		Quantity:       1,
		BasketID:       "mock_basket",
		AddedAt:        time.Now(),
	}
}

// dial sets up a gRPC connection with the microservice endpoint
func (r BasketRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r BasketRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
