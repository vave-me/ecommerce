package tools

import (
	"context"
	"middleman/assistants/internal/models"
)

// ==================== PRODUCT HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeProductHandlers() {
	r.handlers["product_get_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		
		// Validate required parameter
		if err := ValidateIDParam("product_id", productID); err != nil {
			return nil, err
		}
		
		return reg.productRepo.GetProductByID(ctx, productID)
	}

	r.handlers["product_create"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		// Extract parameters
		name := getStringParam(params, "name", "")
		description := getStringParam(params, "description", "")
		basePrice := getInt64Param(params, "base_price", 0)
		categoryID := getStringParam(params, "category_id", "")
		categorySlug := getStringParam(params, "category_slug", "")
		brand := getStringParam(params, "brand", "")
		condition := getStringParam(params, "condition", "")
		model := getStringParam(params, "model", "")
		tags := getStringArrayParam(params, "tags")
		manageStock := getBoolParam(params, "manage_stock", true)
		stock := getInt64Param(params, "stock", 0)
		sku := getStringParam(params, "sku", "")
		attributes := getAttributesParam(params, "attributes")
		weight := getInt64Param(params, "weight", 0)
		
		// Validate parameters
		v := NewValidator()
		v.ValidateRequired("name", name).ValidateMinLength("name", name, 1).ValidateMaxLength("name", name, 200)
		v.ValidateRequired("description", description)
		v.ValidateMinimum("base_price", float64(basePrice), 0)
		
		if condition != "" {
			v.ValidateEnum("condition", condition, []string{"new", "used", "refurbished", "damaged"})
		}
		
		if stock < 0 {
			v.errors = append(v.errors, ValidationError{
				Field:   "stock",
				Message: "must be non-negative",
			})
		}
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		// Sanitize string inputs
		name = SanitizeString(name)
		description = SanitizeString(description)
		brand = SanitizeString(brand)
		model = SanitizeString(model)
		
		err := reg.productRepo.CreateProduct(ctx,
			name,
			description,
			basePrice,
			categoryID,
			categorySlug,
			brand,
			condition,
			model,
			tags,
			manageStock,
			stock,
			sku,
			attributes,
			weight,
			getInt64Param(params, "height", 0),
			getInt64Param(params, "width", 0),
			getInt64Param(params, "depth", 0),
			getStringParam(params, "status", "active"),
			getBoolParam(params, "negotiable", false),
			getStringParam(params, "user_type", ""),
			getBoolParam(params, "middleman_service", false),
			getInt64Param(params, "shipping_cost", 0),
			getBoolParam(params, "has_variants", false),
			getOptionsParam(params, "options"),
			getFloat64Param(params, "lat", 0),
			getFloat64Param(params, "lng", 0),
			getStringParam(params, "thumbnail", ""),
			models.EntityType(getStringParam(params, "entity_type", "product")))
		return nil, err
	}

	r.handlers["product_update_base_price"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		newBasePrice := getInt64Param(params, "new_base_price", 0)
		return nil, reg.productRepo.UpdateProductPrice(ctx, productID, newBasePrice)
	}

	r.handlers["product_delete"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		return nil, reg.productRepo.DeleteProduct(ctx, productID)
	}

	r.handlers["product_search_by_name"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		name := getStringParam(params, "name")
		return reg.productRepo.SearchProductsByName(ctx, name)
	}

	r.handlers["product_search_advanced"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		// Validate search parameters
		if err := ValidateProductSearchParams(params); err != nil {
			return nil, err
		}
		
		return reg.productRepo.SearchProductsAdvanced(ctx,
			getStringParam(params, "name"),
			getStringParam(params, "category_id"),
			getStringParam(params, "category_slug"),
			getInt64Param(params, "min_price", 0),
			getInt64Param(params, "max_price", 0),
			getStringParam(params, "brand"),
			getStringParam(params, "condition"),
			getStringParam(params, "model"),
			getStringArrayParam(params, "tags"),
			getBoolParam(params, "manage_stock", false),
			getInt64Param(params, "min_stock", 0),
			getInt64Param(params, "max_stock", 0),
			getStringParam(params, "sku"),
			getStringParam(params, "status"),
			getBoolParam(params, "negotiable", false),
			getStringParam(params, "user_type"),
			getBoolParam(params, "middleman_service", false),
			getBoolParam(params, "has_variants", false),
			getInt64Param(params, "shipping_cost", 0),
			getInt64Param(params, "min_weight", 0),
			getInt64Param(params, "max_weight", 0),
			getInt64Param(params, "min_height", 0),
			getInt64Param(params, "max_height", 0),
			getInt64Param(params, "min_width", 0),
			getInt64Param(params, "max_width", 0),
			getInt64Param(params, "min_depth", 0),
			getInt64Param(params, "max_depth", 0),
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

	r.handlers["product_get_suggestions"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		name := getStringParam(params, "name")
		return reg.productRepo.GetProductSuggestions(ctx, name)
	}

	r.handlers["product_update_thumbnail"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		thumbnail := getStringParam(params, "thumbnail")
		return nil, reg.productRepo.UpdateProductThumbnail(ctx, productID, thumbnail)
	}

	r.handlers["product_get_by_category_slug"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categorySlug := getStringParam(params, "category_slug")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		return reg.productRepo.GetProductsByCategorySlug(ctx, categorySlug, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["product_get_by_category_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categoryID := getStringParam(params, "category_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		return reg.productRepo.GetProductsByCategoryID(ctx, categoryID, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["product_get_user_catalog"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		return reg.productRepo.GetUserProductCatalog(ctx, userID, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["product_get_public_catalog"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		return reg.productRepo.GetPublicProductCatalog(ctx, userID, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["product_update_price_for_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		newPrice := getInt64Param(params, "new_price", 0)
		oldPrice := getInt64Param(params, "old_price", 0)
		return nil, reg.productRepo.UpdateProductPriceForUser(ctx, productID, newPrice, oldPrice)
	}

	r.handlers["product_adjust_inventory"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		newStock := getInt64Param(params, "new_stock", 0)
		return nil, reg.productRepo.AdjustProductInventory(ctx, productID, newStock)
	}

	r.handlers["product_archive_user_product"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		productID := getStringParam(params, "product_id")
		return nil, reg.productRepo.ArchiveUserProduct(ctx, userID, productID)
	}

	r.handlers["product_mark_as_sold"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		return nil, reg.productRepo.MarkProductAsSold(ctx, productID)
	}

	r.handlers["product_mark_as_leased"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		monthlyPrice := getInt64Param(params, "monthly_price", 0)
		leaseTermMonths := getInt64Param(params, "lease_term_months", 0)
		return nil, reg.productRepo.MarkProductAsLeased(ctx, productID, monthlyPrice, leaseTermMonths)
	}

	r.handlers["product_mark_as_pawned"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		productID := getStringParam(params, "product_id")
		lockedPrice := getInt64Param(params, "locked_price", 0)
		redemptionFee := getInt64Param(params, "redemption_fee", 0)
		return nil, reg.productRepo.MarkProductAsPawned(ctx, userID, productID, lockedPrice, redemptionFee)
	}

	r.handlers["product_increase_price_by"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		increaseAmount := getInt64Param(params, "increase_amount", 0)
		return nil, reg.productRepo.IncreaseProductPriceBy(ctx, productID, increaseAmount)
	}

	r.handlers["product_decrease_price_to"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		newPrice := getInt64Param(params, "new_price", 0)
		return nil, reg.productRepo.DecreaseProductPriceTo(ctx, productID, newPrice)
	}

	r.handlers["product_add_thumbnail"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		productID := getStringParam(params, "product_id")
		thumbnail := getStringParam(params, "thumbnail")
		return nil, reg.productRepo.AddThumbnailToProduct(ctx, productID, thumbnail)
	}

	r.handlers["product_update_details"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return nil, reg.productRepo.UpdateProductDetails(ctx,
			getStringParam(params, "product_id"),
			getStringParam(params, "name"),
			getStringParam(params, "description"),
			getInt64Param(params, "base_price", 0),
			getInt64Param(params, "stock", 0),
			getStringParam(params, "sku"),
			getStringParam(params, "category_id"),
			getStringParam(params, "status"))
	}
}