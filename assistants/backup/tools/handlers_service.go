package tools

import (
	"context"
	"fmt"
	"middleman/assistants/internal/models"
	"time"
)

// ==================== SERVICE HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeServiceHandlers() {
	r.handlers["service_get_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		serviceID := getStringParam(params, "service_id")
		if err := ValidateIDParam("service_id", serviceID); err != nil {
			return nil, fmt.Errorf("invalid service_id: %w", err)
		}
		return reg.serviceRepo.GetServiceByID(ctx, serviceID)
	}

	r.handlers["service_create"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		err := reg.serviceRepo.CreateService(ctx,
			getStringParam(params, "name"),
			getStringParam(params, "description"),
			getStringParam(params, "service_type"),
			getInt64Param(params, "base_price", 0),
			getStringArrayParam(params, "pricing"),
			getStringParam(params, "availability"),
			getStringParam(params, "provider_name"),
			getStringParam(params, "category_id"),
			getStringParam(params, "category_slug"),
			getStringParam(params, "description_short"),
			getStringParam(params, "description_long"),
			getStringArrayParam(params, "qualifications"),
			getStringParam(params, "contact"),
			getStringParam(params, "faq"),
			getStringArrayParam(params, "tags"),
			models.Status(getStringParam(params, "status")),
			models.UserType(getStringParam(params, "user_type")),
			getInt64Param(params, "shipping_cost", 0),
			getBoolParam(params, "negotiable", false),
			getBoolParam(params, "has_variants", false),
			getBoolParam(params, "middleman_service", false),
			getStringArrayParam(params, "attributes"),
			getStringArrayParam(params, "options"),
			getStringParam(params, "thumbnail"),
			getFloat64Param(params, "lat", 0),
			getFloat64Param(params, "lng", 0))
		return nil, err
	}

	r.handlers["service_update_details"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		serviceID := getStringParam(params, "service_id")
		if err := ValidateIDParam("service_id", serviceID); err != nil {
			return nil, fmt.Errorf("invalid service_id: %w", err)
		}
		service := createServiceFromParams(params)
		return reg.serviceRepo.UpdateServiceDetails(ctx, serviceID, service)
	}

	r.handlers["service_delete"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		serviceID := getStringParam(params, "service_id")
		userID := getStringParam(params, "user_id")
		if err := ValidateIDParam("service_id", serviceID); err != nil {
			return nil, fmt.Errorf("invalid service_id: %w", err)
		}
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return nil, reg.serviceRepo.DeleteService(ctx, serviceID, userID)
	}

	r.handlers["service_search_by_name"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		name := getStringParam(params, "name")
		if name == "" {
			return nil, fmt.Errorf("name is required for search")
		}
		name = SanitizeString(name)
		return reg.serviceRepo.SearchServicesByName(ctx, name)
	}

	r.handlers["service_search_advanced"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.serviceRepo.SearchServicesAdvanced(ctx,
			getStringParam(params, "category_id"),
			getStringParam(params, "category_slug"),
			getStringParam(params, "service_type"),
			getStringParam(params, "user_id"),
			models.Status(getStringParam(params, "status")),
			getStringParam(params, "search_text"),
			getInt64Param(params, "min_price", 0),
			getInt64Param(params, "max_price", 0),
			getTimeParam(params, "available_from", time.Time{}),
			getTimeParam(params, "available_to", time.Time{}),
			getBoolParam(params, "has_variants", false),
			getBoolParam(params, "negotiable", false),
			getBoolParam(params, "middleman_service", false),
			models.UserType(getStringParam(params, "user_type")),
			getStringArrayParam(params, "tags"),
			getStringArrayParam(params, "qualifications"),
			getInt64Param(params, "offset", 0),
			getInt64Param(params, "limit", 20),
			getFloat64Param(params, "lat", 0),
			getFloat64Param(params, "lng", 0),
			getInt64Param(params, "radius", 0),
			getInt64Param(params, "page", 1),
			getInt64Param(params, "page_size", 20),
			getStringParam(params, "sort_by"),
			getStringParam(params, "sort_order"))
	}

	r.handlers["service_get_all"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		services, _, err := reg.serviceRepo.GetAllServices(ctx, page, pageSize, sortBy, sortOrder)
		return services, err
	}

	r.handlers["service_get_catalog"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		services, _, err := reg.serviceRepo.GetServiceCatalog(ctx, page, pageSize, sortBy, sortOrder)
		return services, err
	}

	r.handlers["service_get_public_catalog"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		services, _, err := reg.serviceRepo.GetPublicServiceCatalog(ctx, userID, page, pageSize, sortBy, sortOrder)
		return services, err
	}

	r.handlers["service_update_price_for_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		serviceID := getStringParam(params, "service_id")
		newPrice := getInt64Param(params, "new_price", 0)
		oldPrice := getInt64Param(params, "old_price", 0)
		if err := ValidateIDParam("service_id", serviceID); err != nil {
			return nil, fmt.Errorf("invalid service_id: %w", err)
		}
		if newPrice <= 0 {
			return nil, fmt.Errorf("new_price must be greater than zero")
		}
		if oldPrice < 0 {
			return nil, fmt.Errorf("old_price cannot be negative")
		}
		return nil, reg.serviceRepo.UpdateServicePriceForUser(ctx, serviceID, newPrice, oldPrice)
	}

	r.handlers["service_increase_price_by"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		serviceID := getStringParam(params, "service_id")
		increaseAmount := getInt64Param(params, "increase_amount", 0)
		if err := ValidateIDParam("service_id", serviceID); err != nil {
			return nil, fmt.Errorf("invalid service_id: %w", err)
		}
		if increaseAmount <= 0 {
			return nil, fmt.Errorf("increase_amount must be greater than zero")
		}
		oldPrice, newPrice, err := reg.serviceRepo.IncreaseServicePriceBy(ctx, serviceID, increaseAmount)
		if err != nil {
			return nil, err
		}
		return map[string]int64{"old_price": oldPrice, "new_price": newPrice}, nil
	}

	r.handlers["service_decrease_price_to"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		serviceID := getStringParam(params, "service_id")
		newPrice := getInt64Param(params, "new_price", 0)
		if err := ValidateIDParam("service_id", serviceID); err != nil {
			return nil, fmt.Errorf("invalid service_id: %w", err)
		}
		if newPrice <= 0 {
			return nil, fmt.Errorf("new_price must be greater than zero")
		}
		oldPrice, resultPrice, err := reg.serviceRepo.DecreaseServicePriceTo(ctx, serviceID, newPrice)
		if err != nil {
			return nil, err
		}
		return map[string]int64{"old_price": oldPrice, "new_price": resultPrice}, nil
	}

	r.handlers["service_update_branding"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		serviceID := getStringParam(params, "service_id")
		if err := ValidateIDParam("service_id", serviceID); err != nil {
			return nil, fmt.Errorf("invalid service_id: %w", err)
		}
		service := createServiceFromParams(params)
		return nil, reg.serviceRepo.UpdateServiceBranding(ctx, serviceID, service)
	}

	r.handlers["service_get_suggestions"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		name := getStringParam(params, "name")
		if name == "" {
			return nil, fmt.Errorf("name is required for suggestions")
		}
		name = SanitizeString(name)
		return reg.serviceRepo.GetServiceSuggestions(ctx, name)
	}

	r.handlers["service_get_by_category_slug"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categorySlug := getStringParam(params, "category_slug")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if categorySlug == "" {
			return nil, fmt.Errorf("category_slug is required")
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.serviceRepo.GetServicesByCategorySlug(ctx, categorySlug, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["service_get_by_category_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categoryID := getStringParam(params, "category_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if err := ValidateIDParam("category_id", categoryID); err != nil {
			return nil, fmt.Errorf("invalid category_id: %w", err)
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.serviceRepo.GetServicesByCategoryID(ctx, categoryID, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["service_adjust_inventory"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		serviceID := getStringParam(params, "service_id")
		newStock := getInt64Param(params, "new_stock", 0)
		if err := ValidateIDParam("service_id", serviceID); err != nil {
			return nil, fmt.Errorf("invalid service_id: %w", err)
		}
		if newStock < 0 {
			return nil, fmt.Errorf("new_stock cannot be negative")
		}
		oldStock, resultStock, err := reg.serviceRepo.AdjustServiceInventory(ctx, serviceID, newStock)
		if err != nil {
			return nil, err
		}
		return map[string]int64{"old_stock": oldStock, "new_stock": resultStock}, nil
	}

	r.handlers["service_archive_user_service"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		serviceID := getStringParam(params, "service_id")
		if err := ValidateIDParam("service_id", serviceID); err != nil {
			return nil, fmt.Errorf("invalid service_id: %w", err)
		}
		archived, err := reg.serviceRepo.ArchiveUserService(ctx, serviceID)
		if err != nil {
			return nil, err
		}
		return map[string]bool{"archived": archived}, nil
	}

	r.handlers["service_mark_as_sold"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		serviceID := getStringParam(params, "service_id")
		if err := ValidateIDParam("service_id", serviceID); err != nil {
			return nil, fmt.Errorf("invalid service_id: %w", err)
		}
		status, err := reg.serviceRepo.MarkServiceAsSold(ctx, serviceID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"status": status}, nil
	}

	r.handlers["service_mark_as_leased"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		serviceID := getStringParam(params, "service_id")
		monthlyPrice := getInt64Param(params, "monthly_price", 0)
		leaseTermMonths := getInt64Param(params, "lease_term_months", 0)
		if err := ValidateIDParam("service_id", serviceID); err != nil {
			return nil, fmt.Errorf("invalid service_id: %w", err)
		}
		if monthlyPrice <= 0 {
			return nil, fmt.Errorf("monthly_price must be greater than zero")
		}
		if leaseTermMonths <= 0 {
			return nil, fmt.Errorf("lease_term_months must be greater than zero")
		}
		status, err := reg.serviceRepo.MarkServiceAsLeased(ctx, serviceID, monthlyPrice, leaseTermMonths)
		if err != nil {
			return nil, err
		}
		return map[string]string{"status": status}, nil
	}
}