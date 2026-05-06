package tools

import (
	"context"
	"fmt"
)

// ==================== ORDER HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeOrderHandlers() {
	r.handlers["order_create"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		items := convertToOrderItems(params)
		userCustomerID := getStringParam(params, "user_customer_id")

		// Validate required parameters
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		if err := ValidateIDParam("user_customer_id", userCustomerID); err != nil {
			return nil, fmt.Errorf("invalid user_customer_id: %w", err)
		}
		if len(items) == 0 {
			return nil, fmt.Errorf("order must contain at least one item")
		}

		return reg.orderRepo.CreateNewCustomerOrder(ctx, orderID, items, userCustomerID)
	}

	r.handlers["order_get"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		return reg.orderRepo.GetOrderDetailsByID(ctx, orderID)
	}

	r.handlers["order_update"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		updates := getMapParam(params, "updates")
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		if len(updates) == 0 {
			return nil, fmt.Errorf("updates map cannot be empty")
		}
		return reg.orderRepo.UpdateOrderInformation(ctx, orderID, updates)
	}

	r.handlers["order_get_customer_history"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userCustomerID := getStringParam(params, "user_customer_id")
		if err := ValidateIDParam("user_customer_id", userCustomerID); err != nil {
			return nil, fmt.Errorf("invalid user_customer_id: %w", err)
		}
		return reg.orderRepo.GetCustomerOrderHistory(ctx, userCustomerID)
	}

	r.handlers["order_mark_ready"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		return reg.orderRepo.MarkOrderAsReadyForProcessing(ctx, orderID)
	}

	r.handlers["order_cancel"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		reason := getStringParam(params, "reason")
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		if reason == "" {
			return nil, fmt.Errorf("cancellation reason is required")
		}
		reason = SanitizeString(reason)
		return nil, reg.orderRepo.CancelOrderWithReason(ctx, orderID, reason)
	}

	r.handlers["order_complete"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		invoiceID := getStringParam(params, "invoice_id")
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		if err := ValidateIDParam("invoice_id", invoiceID); err != nil {
			return nil, fmt.Errorf("invalid invoice_id: %w", err)
		}
		return reg.orderRepo.CompleteOrderWithInvoice(ctx, orderID, invoiceID)
	}

	r.handlers["order_approve"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		shoppingID := getStringParam(params, "shopping_id")
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		if err := ValidateIDParam("shopping_id", shoppingID); err != nil {
			return nil, fmt.Errorf("invalid shopping_id: %w", err)
		}
		return reg.orderRepo.ApproveOrderForShopping(ctx, orderID, shoppingID)
	}

	r.handlers["order_get_by_status"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		status := getStringParam(params, "status")
		limit := getInt64Param(params, "limit", 50)
		if status == "" {
			return nil, fmt.Errorf("status is required")
		}
		if err := ValidatePaginationParams(1, limit); err != nil {
			return nil, fmt.Errorf("invalid limit: %w", err)
		}
		return reg.orderRepo.FilterOrdersByCurrentStatus(ctx, status, limit)
	}

	r.handlers["order_reject"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		return reg.orderRepo.RejectOrderRequest(ctx, orderID)
	}

	r.handlers["order_ship"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		return reg.orderRepo.MarkOrderAsShipped(ctx, orderID)
	}

	r.handlers["order_deliver"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		orderID := getStringParam(params, "order_id")
		if err := ValidateIDParam("order_id", orderID); err != nil {
			return nil, fmt.Errorf("invalid order_id: %w", err)
		}
		return reg.orderRepo.MarkOrderAsDelivered(ctx, orderID)
	}

	r.handlers["order_search"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		query := getStringParam(params, "query")
		limit := getInt64Param(params, "limit", 50)
		if query == "" {
			return nil, fmt.Errorf("search query is required")
		}
		if err := ValidatePaginationParams(1, limit); err != nil {
			return nil, fmt.Errorf("invalid limit: %w", err)
		}
		query = SanitizeString(query)
		return reg.orderRepo.SearchOrdersByKeyword(ctx, query, limit)
	}
}