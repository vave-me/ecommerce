package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type BasketRepository interface {
	// Core basket operations from protobuf specification
	CreateNewBasket(ctx context.Context, userID string) (*models.StartBasketResponse, error)
	GetBasketByID(ctx context.Context, basketID string) (*models.GetBasketResponse, error)
	CalculateBasketTotal(ctx context.Context, basketID string) (*models.GetTotalBasketAmountResponse, error)
	GetUserCurrentBasket(ctx context.Context, userCustomerID string) (*models.GetCurrentBasketResponse, error)
	CancelBasketWithReason(ctx context.Context, basketID, reason string) (*models.CancelBasketResponse, error)
	CheckoutUserBasket(ctx context.Context, basketID, userCustomerID string) (*models.CheckoutBasketResponse, error)
	AddProductToBasket(ctx context.Context, basketID, productID string, quantity int64) (*models.AddItemResponse, error)
	RemoveProductFromBasket(ctx context.Context, basketID, itemID string, quantity int64) (*models.RemoveItemResponse, error)
	UpdateBasketItemQuantity(ctx context.Context, basketID, itemID string, newQuantity int64) (*models.UpdateItemQuantityResponse, error)
	ListUserBaskets(ctx context.Context, userID, basketStatus string, page, limit int64) (*models.ListBasketsResponse, error)
	GetAllActiveBaskets(ctx context.Context) (*models.GetActiveBasketsResponse, error)
	GetBasketsPaginated(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.Basket, error)
	SearchBasketsByTerm(ctx context.Context, term string) ([]*models.Basket, error)
	GetBasketsByStatusFilter(ctx context.Context, status string, page, pageSize int64) ([]*models.Basket, error)
	GetAllUserBaskets(ctx context.Context, userID string, page, pageSize int64) ([]*models.Basket, error)
	GetUserBasketStatistics(ctx context.Context, userID string) (*models.BasketStatsResponse, error)
	GetExpiredBasketsToCleanup(ctx context.Context, page, pageSize int64) ([]*models.Basket, error)
	GetBasketsInDateRange(ctx context.Context, startDate, endDate string, page, pageSize int64) ([]*models.Basket, error)
	GetAllBasketItems(ctx context.Context, basketID string) ([]*models.BasketItem, error)
	GetBasketItemCount(ctx context.Context, basketID string) (int64, error)
	ValidateBasketContents(ctx context.Context, basketID string) (*models.BasketValidationResponse, error)
	ClearAllBasketItems(ctx context.Context, basketID string) error
	GetUserBasketAnalytics(ctx context.Context, userID string, timeRange string) (*models.BasketAnalyticsResponse, error)
	GetAbandonedBasketsToRecover(ctx context.Context, page, pageSize int64) ([]*models.Basket, error)
	GetBasketConversionStatistics(ctx context.Context, dateRange string) (*models.BasketConversionStatsResponse, error)
	FindBasketByID(ctx context.Context, basketID string) (*models.Basket, error)
	GetUserBasketHistory(ctx context.Context, userID string) ([]*models.Basket, error)
}
