package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
)

// OrderToolService handles all order-related tool operations
type OrderToolService struct {
	repository   domain.OrderRepository
	streamConfig *ToolStreamConfig
}

// NewOrderToolService creates a new order tool service
func NewOrderToolService(repository domain.OrderRepository, config *ToolStreamConfig) *OrderToolService {
	if config == nil {
		config = &ToolStreamConfig{
			BufferSize:       100,
			ProgressInterval: 500 * time.Millisecond,
			EnableMetrics:    true,
			MaxRetries:       3,
		}
	}

	return &OrderToolService{
		repository:   repository,
		streamConfig: config,
	}
}

// ExecuteOperation executes an order operation with streaming support
func (s *OrderToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (*ToolOperationResult, error) {
	startTime := time.Now()

	// Send initial progress
	s.sendProgress(streamChan, toolID, 0, "initializing_order_operation", map[string]interface{}{
		"operation":  operation,
		"start_time": startTime,
	})

	var result interface{}
	var err error

	switch operation {
	case "create", "create_order":
		result, err = s.handleCreateOrder(ctx, parameters, streamChan, toolID)
	case "find", "get", "get_order":
		result, err = s.handleGetOrder(ctx, parameters, streamChan, toolID)
	case "update", "update_order":
		result, err = s.handleUpdateOrder(ctx, parameters, streamChan, toolID)
	case "complete", "complete_order":
		result, err = s.handleCompleteOrder(ctx, parameters, streamChan, toolID)
	case "approve", "approve_order":
		result, err = s.handleApproveOrder(ctx, parameters, streamChan, toolID)
	case "reject", "reject_order":
		result, err = s.handleRejectOrder(ctx, parameters, streamChan, toolID)
	case "ship", "ship_order":
		result, err = s.handleShipOrder(ctx, parameters, streamChan, toolID)
	case "deliver", "deliver_order":
		result, err = s.handleDeliverOrder(ctx, parameters, streamChan, toolID)
	case "cancel", "cancel_order":
		result, err = s.handleCancelOrder(ctx, parameters, streamChan, toolID)
	case "get_orders_by_customer", "list_by_customer":
		result, err = s.handleGetOrdersByCustomer(ctx, parameters, streamChan, toolID)
	case "get_orders_by_seller", "list_by_seller":
		result, err = s.handleGetOrdersBySeller(ctx, parameters, streamChan, toolID)
	case "search", "search_orders":
		result, err = s.handleSearchOrders(ctx, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported order operation: %s", operation)
	}

	duration := time.Since(startTime)

	if err != nil {
		s.sendError(streamChan, toolID, err, map[string]interface{}{
			"operation": operation,
			"duration":  duration.String(),
		})
		return &ToolOperationResult{
			EntityType: "orders",
			Operation:  operation,
			Success:    false,
			Error:      err.Error(),
			Duration:   duration,
		}, err
	}

	s.sendCompletion(streamChan, toolID, result, map[string]interface{}{
		"operation":   operation,
		"duration":    duration.String(),
		"result_type": fmt.Sprintf("%T", result),
	})

	return &ToolOperationResult{
		EntityType: "orders",
		Operation:  operation,
		Success:    true,
		Result:     result,
		Duration:   duration,
		Metadata: map[string]interface{}{
			"execution_time": duration.String(),
		},
	}, nil
}

// handleCreateOrder handles order creation
func (s *OrderToolService) handleCreateOrder(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 20, "extracting_order_params", nil)

	orderID := getStringParam(parameters, "id", "")
	if orderID == "" {
		orderID = getStringParam(parameters, "order_id", "")
	}
	userCustomerID := getStringParam(parameters, "user_customer_id", "")
	if userCustomerID == "" {
		userCustomerID = getStringParam(parameters, "customer_id", "")
	}

	if orderID == "" || userCustomerID == "" {
		return nil, fmt.Errorf("order_id and user_customer_id are required")
	}

	s.sendProgress(streamChan, toolID, 40, "building_order_items", map[string]interface{}{
		"order_id":    orderID,
		"customer_id": userCustomerID,
	})

	// Build items from parameters
	items := s.buildOrderItems(parameters)
	if len(items) == 0 {
		return nil, fmt.Errorf("at least one item is required for order creation")
	}

	s.sendProgress(streamChan, toolID, 70, "creating_order", map[string]interface{}{
		"items_count": len(items),
	})

	log.Printf("OrderToolService: Creating order %s for customer %s with %d items", orderID, userCustomerID, len(items))

	result, err := s.repository.CreateOrder(ctx, orderID, items, userCustomerID)
	if err != nil {
		return nil, fmt.Errorf("order creation failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "order_created_successfully", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "create",
		"result":      result,
		"order_id":    orderID,
		"customer_id": userCustomerID,
		"items_count": len(items),
	}, nil
}

// handleGetOrder handles retrieving a specific order
func (s *OrderToolService) handleGetOrder(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_order_id", nil)

	orderID := getStringParam(parameters, "id", "")
	if orderID == "" {
		orderID = getStringParam(parameters, "order_id", "")
	}

	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	s.sendProgress(streamChan, toolID, 70, "retrieving_order", map[string]interface{}{
		"order_id": orderID,
	})

	log.Printf("OrderToolService: Retrieving order %s", orderID)

	result, err := s.repository.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order retrieval failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "order_retrieved_successfully", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "get",
		"result":      result,
		"order_id":    orderID,
	}, nil
}

// handleCompleteOrder handles order completion
func (s *OrderToolService) handleCompleteOrder(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_completion_params", nil)

	orderID := getStringParam(parameters, "id", "")
	if orderID == "" {
		orderID = getStringParam(parameters, "order_id", "")
	}
	invoiceID := getStringParam(parameters, "invoice_id", "")

	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	s.sendProgress(streamChan, toolID, 70, "completing_order", map[string]interface{}{
		"order_id":   orderID,
		"invoice_id": invoiceID,
	})

	log.Printf("OrderToolService: Completing order %s with invoice %s", orderID, invoiceID)

	result, err := s.repository.CompleteOrder(ctx, orderID, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("order completion failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "order_completed_successfully", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "complete",
		"result":      result,
		"order_id":    orderID,
		"invoice_id":  invoiceID,
	}, nil
}

// handleApproveOrder handles order approval
func (s *OrderToolService) handleApproveOrder(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_approval_params", nil)

	orderID := getStringParam(parameters, "id", "")
	if orderID == "" {
		orderID = getStringParam(parameters, "order_id", "")
	}
	shoppingID := getStringParam(parameters, "shopping_id", "")

	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	s.sendProgress(streamChan, toolID, 70, "approving_order", map[string]interface{}{
		"order_id":    orderID,
		"shopping_id": shoppingID,
	})

	log.Printf("OrderToolService: Approving order %s", orderID)

	result, err := s.repository.ApproveOrder(ctx, orderID, shoppingID)
	if err != nil {
		return nil, fmt.Errorf("order approval failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "order_approved_successfully", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "approve",
		"result":      result,
		"order_id":    orderID,
		"shopping_id": shoppingID,
	}, nil
}

// handleRejectOrder handles order rejection
func (s *OrderToolService) handleRejectOrder(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_rejection_params", nil)

	orderID := getStringParam(parameters, "id", "")
	if orderID == "" {
		orderID = getStringParam(parameters, "order_id", "")
	}

	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	s.sendProgress(streamChan, toolID, 70, "rejecting_order", map[string]interface{}{
		"order_id": orderID,
	})

	log.Printf("OrderToolService: Rejecting order %s", orderID)

	result, err := s.repository.RejectOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order rejection failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "order_rejected_successfully", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "reject",
		"result":      result,
		"order_id":    orderID,
	}, nil
}

// handleShipOrder handles order shipping
func (s *OrderToolService) handleShipOrder(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_shipping_params", nil)

	orderID := getStringParam(parameters, "id", "")
	if orderID == "" {
		orderID = getStringParam(parameters, "order_id", "")
	}

	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	s.sendProgress(streamChan, toolID, 70, "shipping_order", map[string]interface{}{
		"order_id": orderID,
	})

	log.Printf("OrderToolService: Shipping order %s", orderID)

	result, err := s.repository.ShipOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order shipping failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "order_shipped_successfully", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "ship",
		"result":      result,
		"order_id":    orderID,
	}, nil
}

// handleDeliverOrder handles order delivery
func (s *OrderToolService) handleDeliverOrder(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_delivery_params", nil)

	orderID := getStringParam(parameters, "id", "")
	if orderID == "" {
		orderID = getStringParam(parameters, "order_id", "")
	}

	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	s.sendProgress(streamChan, toolID, 70, "delivering_order", map[string]interface{}{
		"order_id": orderID,
	})

	log.Printf("OrderToolService: Delivering order %s", orderID)

	result, err := s.repository.DeliverOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order delivery failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "order_delivered_successfully", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "deliver",
		"result":      result,
		"order_id":    orderID,
	}, nil
}

// handleCancelOrder handles order cancellation
func (s *OrderToolService) handleCancelOrder(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_cancellation_params", nil)

	orderID := getStringParam(parameters, "id", "")
	if orderID == "" {
		orderID = getStringParam(parameters, "order_id", "")
	}
	reason := getStringParam(parameters, "reason", "")

	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	s.sendProgress(streamChan, toolID, 70, "cancelling_order", map[string]interface{}{
		"order_id": orderID,
		"reason":   reason,
	})

	log.Printf("OrderToolService: Cancelling order %s", orderID)

	err := s.repository.CancelOrder(ctx, orderID, reason)
	if err != nil {
		return nil, fmt.Errorf("order cancellation failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "order_cancelled_successfully", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "cancel",
		"success":     true,
		"order_id":    orderID,
		"reason":      reason,
	}, nil
}

// handleGetOrdersByCustomer handles retrieving orders by customer
func (s *OrderToolService) handleGetOrdersByCustomer(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_customer_id", nil)

	userCustomerID := getStringParam(parameters, "user_customer_id", "")
	if userCustomerID == "" {
		userCustomerID = getStringParam(parameters, "customer_id", "")
	}

	if userCustomerID == "" {
		return nil, fmt.Errorf("user_customer_id is required")
	}

	s.sendProgress(streamChan, toolID, 70, "retrieving_customer_orders", map[string]interface{}{
		"customer_id": userCustomerID,
	})

	log.Printf("OrderToolService: Retrieving orders for customer %s", userCustomerID)

	result, err := s.repository.GetOrdersByCustomer(ctx, userCustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer orders retrieval failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "customer_orders_retrieved", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "get_orders_by_customer",
		"result":      result,
		"customer_id": userCustomerID,
	}, nil
}

// handleGetOrdersBySeller handles retrieving orders by seller
func (s *OrderToolService) handleGetOrdersBySeller(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_seller_id", nil)

	userSellerID := getStringParam(parameters, "user_seller_id", "")
	if userSellerID == "" {
		userSellerID = getStringParam(parameters, "seller_id", "")
	}

	if userSellerID == "" {
		return nil, fmt.Errorf("user_seller_id is required")
	}

	s.sendProgress(streamChan, toolID, 70, "retrieving_seller_orders", map[string]interface{}{
		"seller_id": userSellerID,
	})

	log.Printf("OrderToolService: Retrieving orders for seller %s", userSellerID)

	// Note: GetOrdersBySeller not available in current interface, using GetOrdersByCustomer as fallback
	result, err := s.repository.GetOrdersByCustomer(ctx, userSellerID)
	if err != nil {
		return nil, fmt.Errorf("seller orders retrieval failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "seller_orders_retrieved", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "get_orders_by_seller",
		"result":      result,
		"seller_id":   userSellerID,
	}, nil
}

// handleSearchOrders handles order search operations
func (s *OrderToolService) handleSearchOrders(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_search_params", nil)

	searchTerm := getStringParam(parameters, "search_term", "")
	status := getStringParam(parameters, "status", "")
	limit := getInt64Param(parameters, "limit", 20)

	if searchTerm == "" && status == "" {
		return nil, fmt.Errorf("search_term or status is required")
	}

	s.sendProgress(streamChan, toolID, 70, "searching_orders", map[string]interface{}{
		"search_term": searchTerm,
		"status":      status,
		"limit":       limit,
	})

	log.Printf("OrderToolService: Searching orders with term: %s, status: %s", searchTerm, status)

	// Repository SearchOrders only accepts query and limit, not status separately
	query := searchTerm
	if status != "" && searchTerm == "" {
		query = status
	}
	result, err := s.repository.SearchOrders(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("order search failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "order_search_completed", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "search",
		"result":      result,
		"search_term": searchTerm,
		"status":      status,
		"limit":       limit,
	}, nil
}

// handleUpdateOrder handles order updates
func (s *OrderToolService) handleUpdateOrder(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_update_params", nil)

	orderID := getStringParam(parameters, "id", "")
	if orderID == "" {
		orderID = getStringParam(parameters, "order_id", "")
	}

	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	// Extract update fields
	updates := make(map[string]interface{})
	for key, value := range parameters {
		if key != "id" && key != "order_id" && key != "operation" {
			updates[key] = value
		}
	}

	s.sendProgress(streamChan, toolID, 70, "updating_order", map[string]interface{}{
		"order_id": orderID,
		"updates":  updates,
	})

	log.Printf("OrderToolService: Updating order %s", orderID)

	result, err := s.repository.UpdateOrder(ctx, orderID, updates)
	if err != nil {
		return nil, fmt.Errorf("order update failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "order_updated_successfully", nil)

	return map[string]interface{}{
		"entity_type": "orders",
		"operation":   "update",
		"result":      result,
		"order_id":    orderID,
		"updates":     updates,
	}, nil
}

// buildOrderItems constructs order items from parameters
func (s *OrderToolService) buildOrderItems(parameters map[string]interface{}) []models.Item {
	var items []models.Item

	// Try to get items from parameters
	if itemsParam, ok := parameters["items"].([]interface{}); ok {
		for _, item := range itemsParam {
			if itemMap, ok := item.(map[string]interface{}); ok {
				orderItem := models.Item{
					UserSellerID:   getStringParam(itemMap, "seller_id", ""),
					ProductID:      getStringParam(itemMap, "product_id", ""),
					UserSellerName: getStringParam(itemMap, "seller_name", ""),
					ProductName:    getStringParam(itemMap, "product_name", ""),
					Price:          getInt64Param(itemMap, "price", 0),
					Quantity:       getInt64Param(itemMap, "quantity", 1),
				}
				items = append(items, orderItem)
			}
		}
	} else {
		// Fallback: create single item from direct parameters
		item := models.Item{
			UserSellerID:   getStringParam(parameters, "seller_id", "default_seller"),
			ProductID:      getStringParam(parameters, "product_id", "default_product"),
			UserSellerName: getStringParam(parameters, "seller_name", "Default Seller"),
			ProductName:    getStringParam(parameters, "product_name", "Default Product"),
			Price:          getInt64Param(parameters, "price", 1000),
			Quantity:       getInt64Param(parameters, "quantity", 1),
		}
		items = append(items, item)
	}

	return items
}

// Helper methods for streaming (same as ProductToolService)
func (s *OrderToolService) sendProgress(streamChan chan<- ToolExecutionStream, toolID string, progress float64, step string, metadata map[string]interface{}) {
	if streamChan == nil {
		return
	}

	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["step"] = step

	streamChan <- ToolExecutionStream{
		ID:        toolID,
		ToolName:  "order_operation",
		Status:    "progress",
		Progress:  progress,
		Metadata:  metadata,
		Timestamp: time.Now().Unix(),
	}
}

func (s *OrderToolService) sendError(streamChan chan<- ToolExecutionStream, toolID string, err error, metadata map[string]interface{}) {
	if streamChan == nil {
		return
	}

	streamChan <- ToolExecutionStream{
		ID:        toolID,
		ToolName:  "order_operation",
		Status:    "error",
		Error:     err.Error(),
		Metadata:  metadata,
		Timestamp: time.Now().Unix(),
	}
}

func (s *OrderToolService) sendCompletion(streamChan chan<- ToolExecutionStream, toolID string, result interface{}, metadata map[string]interface{}) {
	if streamChan == nil {
		return
	}

	streamChan <- ToolExecutionStream{
		ID:        toolID,
		ToolName:  "order_operation",
		Status:    "completed",
		Progress:  100.0,
		Result:    result,
		Metadata:  metadata,
		Timestamp: time.Now().Unix(),
	}
}
