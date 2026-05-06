package tools

import (
	"context"
)

// ==================== BASKET HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeBasketHandlers() {
	// Core basket operations
	r.handlers["basket_create_new"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		return reg.basketRepo.CreateNewBasket(ctx, userID)
	}

	r.handlers["basket_get_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		return reg.basketRepo.GetBasketByID(ctx, basketID)
	}

	r.handlers["basket_find_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		return reg.basketRepo.FindBasketByID(ctx, basketID)
	}

	r.handlers["basket_calculate_total"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		return reg.basketRepo.CalculateBasketTotal(ctx, basketID)
	}

	r.handlers["basket_get_current"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userCustomerID := getStringParam(params, "user_customer_id")
		return reg.basketRepo.GetUserCurrentBasket(ctx, userCustomerID)
	}

	r.handlers["basket_cancel"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		reason := getStringParam(params, "reason")
		return reg.basketRepo.CancelBasketWithReason(ctx, basketID, reason)
	}

	r.handlers["basket_checkout"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		userCustomerID := getStringParam(params, "user_customer_id")
		return reg.basketRepo.CheckoutUserBasket(ctx, basketID, userCustomerID)
	}

	// Item management
	r.handlers["basket_add_product"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		productID := getStringParam(params, "product_id")
		quantity := getInt64Param(params, "quantity", 1)
		return reg.basketRepo.AddProductToBasket(ctx, basketID, productID, quantity)
	}

	r.handlers["basket_remove_product"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		itemID := getStringParam(params, "item_id")
		quantity := getInt64Param(params, "quantity", 1)
		return reg.basketRepo.RemoveProductFromBasket(ctx, basketID, itemID, quantity)
	}

	r.handlers["basket_update_item_quantity"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		itemID := getStringParam(params, "item_id")
		newQuantity := getInt64Param(params, "new_quantity", 1)
		return reg.basketRepo.UpdateBasketItemQuantity(ctx, basketID, itemID, newQuantity)
	}

	r.handlers["basket_get_all_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		return reg.basketRepo.GetAllBasketItems(ctx, basketID)
	}

	r.handlers["basket_get_item_count"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		count, err := reg.basketRepo.GetBasketItemCount(ctx, basketID)
		if err != nil {
			return nil, err
		}
		return map[string]int64{"count": count}, nil
	}

	r.handlers["basket_clear_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		return nil, reg.basketRepo.ClearAllBasketItems(ctx, basketID)
	}

	// Basket listing and search
	r.handlers["basket_list_user_baskets"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		basketStatus := getStringParam(params, "basket_status", "")
		page := getInt64Param(params, "page", 1)
		limit := getInt64Param(params, "limit", 20)
		return reg.basketRepo.ListUserBaskets(ctx, userID, basketStatus, page, limit)
	}

	r.handlers["basket_get_all_active"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.basketRepo.GetAllActiveBaskets(ctx)
	}

	r.handlers["basket_get_paginated"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by", "created_at")
		sortOrder := getStringParam(params, "sort_order", "desc")
		return reg.basketRepo.GetBasketsPaginated(ctx, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["basket_search"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		term := getStringParam(params, "term")
		return reg.basketRepo.SearchBasketsByTerm(ctx, term)
	}

	r.handlers["basket_get_by_status"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		status := getStringParam(params, "status")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		return reg.basketRepo.GetBasketsByStatusFilter(ctx, status, page, pageSize)
	}

	r.handlers["basket_get_all_user_baskets"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		return reg.basketRepo.GetAllUserBaskets(ctx, userID, page, pageSize)
	}

	r.handlers["basket_get_user_history"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		return reg.basketRepo.GetUserBasketHistory(ctx, userID)
	}

	// Analytics and statistics
	r.handlers["basket_get_user_statistics"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		return reg.basketRepo.GetUserBasketStatistics(ctx, userID)
	}

	r.handlers["basket_get_user_analytics"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		timeRange := getStringParam(params, "time_range", "30d")
		return reg.basketRepo.GetUserBasketAnalytics(ctx, userID, timeRange)
	}

	r.handlers["basket_get_conversion_statistics"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		dateRange := getStringParam(params, "date_range", "30d")
		return reg.basketRepo.GetBasketConversionStatistics(ctx, dateRange)
	}

	// Maintenance and recovery
	r.handlers["basket_get_expired"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		return reg.basketRepo.GetExpiredBasketsToCleanup(ctx, page, pageSize)
	}

	r.handlers["basket_get_abandoned"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		return reg.basketRepo.GetAbandonedBasketsToRecover(ctx, page, pageSize)
	}

	r.handlers["basket_get_in_date_range"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		startDate := getStringParam(params, "start_date")
		endDate := getStringParam(params, "end_date")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		return reg.basketRepo.GetBasketsInDateRange(ctx, startDate, endDate, page, pageSize)
	}

	// Validation
	r.handlers["basket_validate_contents"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		basketID := getStringParam(params, "basket_id")
		return reg.basketRepo.ValidateBasketContents(ctx, basketID)
	}
}