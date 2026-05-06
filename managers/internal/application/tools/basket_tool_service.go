package tools

import (
	"context"
	"fmt"
	"time"

	"middleman/managers/internal/domain"
)

// BasketToolService handles basket/shopping cart operations
type BasketToolService struct {
	basketRepo domain.BasketRepository
	config     *ToolStreamConfig
}

// NewBasketToolService creates a new basket tool service
func NewBasketToolService(basketRepo domain.BasketRepository) *BasketToolService {
	return &BasketToolService{
		basketRepo: basketRepo,
		config: &ToolStreamConfig{
			BufferSize:       100,
			ProgressInterval: 500 * time.Millisecond,
			EnableMetrics:    true,
			MaxRetries:       3,
		},
	}
}

// ExecuteOperation routes basket operations to appropriate handlers
func (s *BasketToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Send initial progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "basket_operation",
			Status:   "started",
			Progress: 0,
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "BasketToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	var result interface{}
	var err error

	switch operation {
	case "start_basket", "create":
		result, err = s.startBasket(ctx, parameters, streamChan, toolID)
	case "get_basket", "find":
		result, err = s.getBasket(ctx, parameters, streamChan, toolID)
	case "get_current_basket":
		result, err = s.getCurrentBasket(ctx, parameters, streamChan, toolID)
	case "add_item":
		result, err = s.addItem(ctx, parameters, streamChan, toolID)
	case "remove_item":
		result, err = s.removeItem(ctx, parameters, streamChan, toolID)
	case "update_quantity":
		result, err = s.updateItemQuantity(ctx, parameters, streamChan, toolID)
	case "get_total":
		result, err = s.getBasketTotal(ctx, parameters, streamChan, toolID)
	case "checkout":
		result, err = s.checkoutBasket(ctx, parameters, streamChan, toolID)
	case "cancel":
		result, err = s.cancelBasket(ctx, parameters, streamChan, toolID)
	case "list_baskets":
		result, err = s.listBaskets(ctx, parameters, streamChan, toolID)
	case "clear_basket":
		result, err = s.clearBasket(ctx, parameters, streamChan, toolID)
	case "get_analytics":
		result, err = s.getBasketAnalytics(ctx, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported basket operation: %s", operation)
	}

	// Send completion status
	if streamChan != nil {
		status := "completed"
		if err != nil {
			status = "error"
		}

		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "basket_operation",
			Status:   status,
			Progress: 100,
			Result:   result,
			Error:    s.getErrorString(err),
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "BasketToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return result, err
}

func (s *BasketToolService) startBasket(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := s.getStringParam(params, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	s.sendProgress(streamChan, toolID, "starting_basket", 50)

	response, err := s.basketRepo.CreateNewBasket(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to start basket: %w", err)
	}

	return response, nil
}

func (s *BasketToolService) getBasket(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	basketID := s.getStringParam(params, "basket_id", "")
	if basketID == "" {
		basketID = s.getStringParam(params, "id", "")
	}
	if basketID == "" {
		return nil, fmt.Errorf("basket_id is required")
	}

	s.sendProgress(streamChan, toolID, "getting_basket", 50)

	response, err := s.basketRepo.GetBasketByID(ctx, basketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	return response, nil
}

func (s *BasketToolService) getCurrentBasket(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userCustomerID := s.getStringParam(params, "user_customer_id", "")
	if userCustomerID == "" {
		userCustomerID = s.getStringParam(params, "user_id", "")
	}
	if userCustomerID == "" {
		return nil, fmt.Errorf("user_customer_id is required")
	}

	s.sendProgress(streamChan, toolID, "getting_current_basket", 50)

	response, err := s.basketRepo.GetUserCurrentBasket(ctx, userCustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current basket: %w", err)
	}

	return response, nil
}

func (s *BasketToolService) addItem(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	basketID := s.getStringParam(params, "basket_id", "")
	productID := s.getStringParam(params, "product_id", "")
	quantity := s.getInt64Param(params, "quantity", 1)

	if basketID == "" || productID == "" {
		return nil, fmt.Errorf("basket_id and product_id are required")
	}

	s.sendProgress(streamChan, toolID, "adding_item", 50)

	response, err := s.basketRepo.AddProductToBasket(ctx, basketID, productID, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to add item: %w", err)
	}

	return response, nil
}

func (s *BasketToolService) removeItem(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	basketID := s.getStringParam(params, "basket_id", "")
	itemID := s.getStringParam(params, "item_id", "")
	quantity := s.getInt64Param(params, "quantity", 1)

	if basketID == "" || itemID == "" {
		return nil, fmt.Errorf("basket_id and item_id are required")
	}

	s.sendProgress(streamChan, toolID, "removing_item", 50)

	response, err := s.basketRepo.RemoveProductFromBasket(ctx, basketID, itemID, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to remove item: %w", err)
	}

	return response, nil
}

func (s *BasketToolService) updateItemQuantity(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	basketID := s.getStringParam(params, "basket_id", "")
	itemID := s.getStringParam(params, "item_id", "")
	newQuantity := s.getInt64Param(params, "new_quantity", 1)

	if basketID == "" || itemID == "" {
		return nil, fmt.Errorf("basket_id and item_id are required")
	}

	s.sendProgress(streamChan, toolID, "updating_quantity", 50)

	response, err := s.basketRepo.UpdateBasketItemQuantity(ctx, basketID, itemID, newQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to update quantity: %w", err)
	}

	return response, nil
}

func (s *BasketToolService) getBasketTotal(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	basketID := s.getStringParam(params, "basket_id", "")
	if basketID == "" {
		return nil, fmt.Errorf("basket_id is required")
	}

	s.sendProgress(streamChan, toolID, "calculating_total", 50)

	response, err := s.basketRepo.CalculateBasketTotal(ctx, basketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get basket total: %w", err)
	}

	return response, nil
}

func (s *BasketToolService) checkoutBasket(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	basketID := s.getStringParam(params, "basket_id", "")
	userCustomerID := s.getStringParam(params, "user_customer_id", "")

	if basketID == "" || userCustomerID == "" {
		return nil, fmt.Errorf("basket_id and user_customer_id are required")
	}

	s.sendProgress(streamChan, toolID, "checking_out", 50)

	response, err := s.basketRepo.CheckoutUserBasket(ctx, basketID, userCustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to checkout basket: %w", err)
	}

	return response, nil
}

func (s *BasketToolService) cancelBasket(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	basketID := s.getStringParam(params, "basket_id", "")
	reason := s.getStringParam(params, "reason", "User requested cancellation")

	if basketID == "" {
		return nil, fmt.Errorf("basket_id is required")
	}

	s.sendProgress(streamChan, toolID, "cancelling_basket", 50)

	response, err := s.basketRepo.CancelBasketWithReason(ctx, basketID, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel basket: %w", err)
	}

	return response, nil
}

func (s *BasketToolService) listBaskets(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := s.getStringParam(params, "user_id", "")
	basketStatus := s.getStringParam(params, "status", "")
	page := s.getInt64Param(params, "page", 1)
	limit := s.getInt64Param(params, "limit", 20)

	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	s.sendProgress(streamChan, toolID, "listing_baskets", 50)

	response, err := s.basketRepo.ListUserBaskets(ctx, userID, basketStatus, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list baskets: %w", err)
	}

	return response, nil
}

func (s *BasketToolService) clearBasket(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	basketID := s.getStringParam(params, "basket_id", "")
	if basketID == "" {
		return nil, fmt.Errorf("basket_id is required")
	}

	s.sendProgress(streamChan, toolID, "clearing_basket", 50)

	err := s.basketRepo.ClearAllBasketItems(ctx, basketID)
	if err != nil {
		return nil, fmt.Errorf("failed to clear basket: %w", err)
	}

	return map[string]interface{}{
		"basket_id": basketID,
		"cleared":   true,
	}, nil
}

func (s *BasketToolService) getBasketAnalytics(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := s.getStringParam(params, "user_id", "")
	timeRange := s.getStringParam(params, "time_range", "30d")

	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	s.sendProgress(streamChan, toolID, "getting_analytics", 50)

	response, err := s.basketRepo.GetUserBasketAnalytics(ctx, userID, timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics: %w", err)
	}

	return response, nil
}

// Helper functions
func (s *BasketToolService) getErrorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func (s *BasketToolService) sendProgress(streamChan chan<- ToolExecutionStream, toolID string, step string, progress float64) {
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "basket_operation",
			Status:   "progress",
			Progress: progress,
			Metadata: map[string]interface{}{
				"step": step,
			},
			Timestamp: time.Now().Unix(),
		}
	}
}

func (s *BasketToolService) getStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func (s *BasketToolService) getInt64Param(params map[string]interface{}, key string, defaultValue int64) int64 {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case int:
			return int64(v)
		case int64:
			return v
		case float64:
			return int64(v)
		}
	}
	return defaultValue
}
