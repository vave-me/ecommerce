package tools

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
)

// ProductToolService handles all product-related tool operations
type ProductToolService struct {
	repository   domain.ProductRepository
	streamConfig *ToolStreamConfig
}

// NewProductToolService creates a new product tool service
func NewProductToolService(repository domain.ProductRepository, config *ToolStreamConfig) *ProductToolService {
	if config == nil {
		config = &ToolStreamConfig{
			BufferSize:       100,
			ProgressInterval: 500 * time.Millisecond,
			EnableMetrics:    true,
			MaxRetries:       3,
		}
	}

	return &ProductToolService{
		repository:   repository,
		streamConfig: config,
	}
}

// ExecuteOperation executes a product operation with streaming support
func (s *ProductToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (*ToolOperationResult, error) {
	startTime := time.Now()

	// Enhanced logging for operation start
	log.Printf("[PRODUCT_TOOL] ========== EXECUTE_OPERATION START ==========")
	log.Printf("[PRODUCT_TOOL] Operation: %s", operation)
	log.Printf("[PRODUCT_TOOL] ToolID: %s", toolID)
	log.Printf("[PRODUCT_TOOL] Parameters: %+v", parameters)

	// Send initial progress
	s.sendProgress(streamChan, toolID, 0, "initializing", map[string]interface{}{
		"operation":  operation,
		"start_time": startTime,
	})

	var result interface{}
	var err error

	log.Printf("[PRODUCT_TOOL] Routing to handler for operation: %s", operation)

	switch operation {
	case "search", "search_by_term":
		log.Printf("[PRODUCT_TOOL] Handling as SEARCH operation")
		result, err = s.handleProductSearch(ctx, parameters, streamChan, toolID)
	case "find", "get":
		log.Printf("[PRODUCT_TOOL] Handling as FIND operation")
		result, err = s.handleProductFind(ctx, parameters, streamChan, toolID)
	case "filter", "search_with_filters":
		log.Printf("[PRODUCT_TOOL] Handling as FILTER operation")
		result, err = s.handleProductFilters(ctx, parameters, streamChan, toolID)
	case "suggest":
		log.Printf("[PRODUCT_TOOL] Handling as SUGGEST operation")
		result, err = s.handleProductSuggestions(ctx, parameters, streamChan, toolID)
	case "category_search":
		log.Printf("[PRODUCT_TOOL] Handling as CATEGORY_SEARCH operation")
		result, err = s.handleCategorySearch(ctx, parameters, streamChan, toolID)
	case "add", "create":
		log.Printf("[PRODUCT_TOOL] Handling as ADD operation")
		result, err = s.handleProductAdd(ctx, parameters, streamChan, toolID)
	case "update":
		log.Printf("[PRODUCT_TOOL] Handling as UPDATE operation")
		result, err = s.handleProductUpdate(ctx, parameters, streamChan, toolID)
	case "remove", "delete":
		log.Printf("[PRODUCT_TOOL] Handling as REMOVE operation")
		result, err = s.handleProductRemove(ctx, parameters, streamChan, toolID)
	default:
		log.Printf("[PRODUCT_TOOL] ERROR: Unsupported operation: %s", operation)
		err = fmt.Errorf("unsupported product operation: %s", operation)
	}

	duration := time.Since(startTime)

	if err != nil {
		log.Printf("[PRODUCT_TOOL] ❌ Operation FAILED: %v", err)
		log.Printf("[PRODUCT_TOOL] Duration: %v", duration)
		s.sendError(streamChan, toolID, err, map[string]interface{}{
			"operation": operation,
			"duration":  duration.String(),
		})
		return &ToolOperationResult{
			EntityType: "products",
			Operation:  operation,
			Success:    false,
			Error:      err.Error(),
			Duration:   duration,
		}, err
	}

	// Log successful result details
	log.Printf("[PRODUCT_TOOL] ✅ Operation SUCCESS")
	log.Printf("[PRODUCT_TOOL] Duration: %v", duration)
	log.Printf("[PRODUCT_TOOL] Result type: %T", result)

	// Log result content
	if resultMap, ok := result.(map[string]interface{}); ok {
		log.Printf("[PRODUCT_TOOL] Result map keys: %v", func() []string {
			keys := make([]string, 0, len(resultMap))
			for k := range resultMap {
				keys = append(keys, k)
			}
			return keys
		}())

		if count, exists := resultMap["count"]; exists {
			log.Printf("[PRODUCT_TOOL] Result count: %v", count)
		}

		if results, exists := resultMap["results"]; exists {
			log.Printf("[PRODUCT_TOOL] Results type: %T", results)
			if resultsSlice, ok := results.([]map[string]interface{}); ok {
				log.Printf("[PRODUCT_TOOL] Results array length: %d", len(resultsSlice))
				if len(resultsSlice) > 0 {
					log.Printf("[PRODUCT_TOOL] First result sample: %+v", resultsSlice[0])
				}
			}
		}
	}

	s.sendCompletion(streamChan, toolID, result, map[string]interface{}{
		"operation":   operation,
		"duration":    duration.String(),
		"result_type": fmt.Sprintf("%T", result),
	})

	finalResult := &ToolOperationResult{
		EntityType: "products",
		Operation:  operation,
		Success:    true,
		Result:     result,
		Duration:   duration,
		Metadata: map[string]interface{}{
			"records_count":  s.getResultCount(result),
			"execution_time": duration.String(),
		},
	}

	log.Printf("[PRODUCT_TOOL] Final ToolOperationResult: %+v", finalResult)
	log.Printf("[PRODUCT_TOOL] ========== EXECUTE_OPERATION END ==========")

	return finalResult, nil
}

// handleProductSearch handles product search operations
func (s *ProductToolService) handleProductSearch(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	log.Printf("[PRODUCT_TOOL_SEARCH] === SEARCH HANDLER START ===")
	log.Printf("[PRODUCT_TOOL_SEARCH] Parameters received: %+v", parameters)

	s.sendProgress(streamChan, toolID, 25, "extracting_search_term", nil)

	// Extract search term - make it optional for "list all" functionality
	searchTerm := getStringParam(parameters, "search_term", "")
	log.Printf("[PRODUCT_TOOL_SEARCH] search_term parameter: '%s'", searchTerm)

	if searchTerm == "" {
		searchTerm = getStringParam(parameters, "query", "")
		log.Printf("[PRODUCT_TOOL_SEARCH] query parameter: '%s'", searchTerm)
	}
	if searchTerm == "" {
		searchTerm = getStringParam(parameters, "name", "")
		log.Printf("[PRODUCT_TOOL_SEARCH] name parameter: '%s'", searchTerm)
	}

	log.Printf("[PRODUCT_TOOL_SEARCH] Final search term: '%s'", searchTerm)

	s.sendProgress(streamChan, toolID, 50, "executing_search", map[string]interface{}{
		"search_term": searchTerm,
	})

	var products []*models.Product
	var err error

	if searchTerm != "" {
		// Search with specific term
		log.Printf("[PRODUCT_TOOL_SEARCH] Executing SearchWithTerm for: '%s'", searchTerm)
		products, err = s.repository.SearchWithTerm(ctx, searchTerm)
	} else {
		// Default search with basic filters to return all products (like deal service)
		log.Printf("[PRODUCT_TOOL_SEARCH] Executing SearchWithFilters for ALL products")
		products, err = s.repository.SearchWithFilters(
			ctx, "", "", "", 0, 0, "", "", "", []string{}, false, 0, 0, "", "", false, "", false, false, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0.0, 0.0, 0, 1, 20, "created_at", "desc",
		)
	}

	if err != nil {
		log.Printf("[PRODUCT_TOOL_SEARCH] ❌ Search ERROR: %v", err)
		return nil, fmt.Errorf("product search failed: %w", err)
	}

	log.Printf("[PRODUCT_TOOL_SEARCH] ✅ Search returned %d products", len(products))

	// Log first few products for debugging
	for i, product := range products {
		if i < 3 && product != nil {
			log.Printf("[PRODUCT_TOOL_SEARCH] Product[%d]: ID=%s, Name=%s, Price=%d, Status=%s",
				i, product.ProductID, product.Name, product.BasePrice, product.Status)
		}
	}

	s.sendProgress(streamChan, toolID, 90, "formatting_results", map[string]interface{}{
		"products_found": len(products),
	})

	// Convert domain products to serializable maps
	log.Printf("[PRODUCT_TOOL_SEARCH] Converting %d products to maps", len(products))
	productMaps := make([]map[string]interface{}, len(products))
	for i, product := range products {
		productMap := s.productToMap(product)
		productMaps[i] = productMap

		// Log the first converted product for verification
		if i == 0 && productMap != nil {
			log.Printf("[PRODUCT_TOOL_SEARCH] First product map sample: %+v", productMap)
		}
	}

	result := map[string]interface{}{
		"entity_type": "products",
		"operation":   "search",
		"results":     productMaps,
		"count":       len(products),
		"search_term": searchTerm,
	}

	log.Printf("[PRODUCT_TOOL_SEARCH] Final result structure: entity_type=%s, operation=%s, count=%d",
		result["entity_type"], result["operation"], result["count"])
	log.Printf("[PRODUCT_TOOL_SEARCH] === SEARCH HANDLER END ===")

	return result, nil
}

// handleProductFind handles finding a specific product
func (s *ProductToolService) handleProductFind(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_product_id", nil)

	productID, ok := parameters["id"].(string)
	if !ok {
		productID, ok = parameters["product_id"].(string)
		if !ok {
			return nil, fmt.Errorf("id or product_id parameter required")
		}
	}

	s.sendProgress(streamChan, toolID, 50, "finding_product", map[string]interface{}{
		"product_id": productID,
	})

	log.Printf("ProductToolService: Finding product with ID: %s", productID)
	product, err := s.repository.Find(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product find failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "formatting_result", nil)

	return map[string]interface{}{
		"entity_type": "products",
		"operation":   "find",
		"result":      s.productToMap(product),
		"product_id":  productID,
	}, nil
}

// handleProductFilters handles complex product filtering
func (s *ProductToolService) handleProductFilters(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	log.Printf("[PRODUCT_TOOL_FILTER] === FILTER HANDLER START ===")
	log.Printf("[PRODUCT_TOOL_FILTER] Parameters received: %+v", parameters)

	s.sendProgress(streamChan, toolID, 20, "parsing_filters", nil)

	// Extract filter parameters with defaults
	name := getStringParam(parameters, "name", "")
	categoryId := getStringParam(parameters, "categoryId", "")
	categorySlug := getStringParam(parameters, "categorySlug", "")
	minPrice := getInt64Param(parameters, "min_price", 0)
	maxPrice := getInt64Param(parameters, "max_price", 0)
	brand := getStringParam(parameters, "brand", "")
	condition := getStringParam(parameters, "condition", "")
	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)
	sortBy := getStringParam(parameters, "sort_by", "")
	sortOrder := getStringParam(parameters, "sort_order", "desc")

	log.Printf("[PRODUCT_TOOL_FILTER] Parsed filters:")
	log.Printf("[PRODUCT_TOOL_FILTER]   - name: '%s'", name)
	log.Printf("[PRODUCT_TOOL_FILTER]   - categoryId: '%s'", categoryId)
	log.Printf("[PRODUCT_TOOL_FILTER]   - categorySlug: '%s'", categorySlug)
	log.Printf("[PRODUCT_TOOL_FILTER]   - price range: %d-%d", minPrice, maxPrice)
	log.Printf("[PRODUCT_TOOL_FILTER]   - brand: '%s'", brand)
	log.Printf("[PRODUCT_TOOL_FILTER]   - condition: '%s'", condition)
	log.Printf("[PRODUCT_TOOL_FILTER]   - pagination: page=%d, size=%d", page, pageSize)
	log.Printf("[PRODUCT_TOOL_FILTER]   - sort: by=%s, order=%s", sortBy, sortOrder)

	s.sendProgress(streamChan, toolID, 50, "applying_filters", map[string]interface{}{
		"filters": map[string]interface{}{
			"name":         name,
			"categoryId":   categoryId,
			"categorySlug": categorySlug,
			"min_price":    minPrice,
			"max_price":    maxPrice,
			"brand":        brand,
			"condition":    condition,
		},
	})

	log.Printf("[PRODUCT_TOOL_FILTER] Calling repository.SearchWithFilters")
	products, err := s.repository.SearchWithFilters(
		ctx, name, categoryId, categorySlug, minPrice, maxPrice, brand, condition,
		"", []string{}, false, 0, 0, "", "", false, "", false, false, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0.0, 0.0, 0, page, pageSize, sortBy, sortOrder,
	)

	if err != nil {
		log.Printf("[PRODUCT_TOOL_FILTER] ❌ Filter search ERROR: %v", err)
		return nil, fmt.Errorf("product filter search failed: %w", err)
	}

	log.Printf("[PRODUCT_TOOL_FILTER] ✅ Filter search returned %d products", len(products))

	// Log first few filtered products
	for i, product := range products {
		if i < 3 && product != nil {
			log.Printf("[PRODUCT_TOOL_FILTER] Product[%d]: ID=%s, Name=%s, Price=%d, Category=%s",
				i, product.ProductID, product.Name, product.BasePrice, product.CategoryID)
		}
	}

	s.sendProgress(streamChan, toolID, 90, "formatting_filtered_results", map[string]interface{}{
		"products_found": len(products),
	})

	// Convert domain products to serializable maps
	log.Printf("[PRODUCT_TOOL_FILTER] Converting %d products to maps", len(products))
	productMaps := make([]map[string]interface{}, len(products))
	for i, product := range products {
		productMaps[i] = s.productToMap(product)
	}

	result := map[string]interface{}{
		"entity_type": "products",
		"operation":   "filter",
		"results":     productMaps,
		"count":       len(products),
		"filters": map[string]interface{}{
			"name":         name,
			"categoryId":   categoryId,
			"categorySlug": categorySlug,
			"min_price":    minPrice,
			"max_price":    maxPrice,
			"brand":        brand,
			"condition":    condition,
		},
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	}

	log.Printf("[PRODUCT_TOOL_FILTER] Final result: count=%d, has filters=%v, has pagination=%v",
		result["count"], result["filters"] != nil, result["pagination"] != nil)
	log.Printf("[PRODUCT_TOOL_FILTER] === FILTER HANDLER END ===")

	return result, nil
}

// handleProductSuggestions handles product suggestion operations
func (s *ProductToolService) handleProductSuggestions(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_suggestion_params", nil)

	query := getStringParam(parameters, "query", "")
	limit := getInt64Param(parameters, "limit", 10)

	if query == "" {
		return nil, fmt.Errorf("query parameter required for suggestions")
	}

	s.sendProgress(streamChan, toolID, 50, "generating_suggestions", map[string]interface{}{
		"query": query,
		"limit": limit,
	})

	// Use search for suggestions (could be enhanced with dedicated suggestion logic)
	products, err := s.repository.SearchWithTerm(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("product suggestions failed: %w", err)
	}

	// Limit results for suggestions
	if int64(len(products)) > limit {
		products = products[:limit]
	}

	s.sendProgress(streamChan, toolID, 90, "formatting_suggestions", map[string]interface{}{
		"suggestions_count": len(products),
	})

	// Convert domain products to serializable maps
	productMaps := make([]map[string]interface{}, len(products))
	for i, product := range products {
		productMaps[i] = s.productToMap(product)
	}

	return map[string]interface{}{
		"entity_type": "products",
		"operation":   "suggest",
		"results":     productMaps,
		"count":       len(products),
		"query":       query,
		"limit":       limit,
	}, nil
}

// handleCategorySearch handles category-based product search
func (s *ProductToolService) handleCategorySearch(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_category", nil)

	categoryId := getStringParam(parameters, "categoryId", "")
	categorySlug := getStringParam(parameters, "categorySlug", "")
	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)
	sortBy := getStringParam(parameters, "sort_by", "")
	sortOrder := getStringParam(parameters, "sort_order", "desc")

	if categoryId == "" || categorySlug == "" {
		return nil, fmt.Errorf("category parameter required")
	}

	s.sendProgress(streamChan, toolID, 50, "searching_by_category", map[string]interface{}{
		"categoryId": categoryId,
	})

	// Use filter search with category
	products, err := s.repository.SearchWithFilters(
		ctx, "", categoryId, categorySlug, 0, 0, "", "",
		"", []string{}, false, 0, 0, "", "", false, "", false, false, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0.0, 0.0, 0, page, pageSize, sortBy, sortOrder,
	)

	if err != nil {
		return nil, fmt.Errorf("category search failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "formatting_category_results", map[string]interface{}{
		"products_found": len(products),
	})

	// Convert domain products to serializable maps
	productMaps := make([]map[string]interface{}, len(products))
	for i, product := range products {
		productMaps[i] = s.productToMap(product)
	}

	return map[string]interface{}{
		"entity_type":  "products",
		"operation":    "category_search",
		"results":      productMaps,
		"count":        len(products),
		"categoryId":   categoryId,
		"categorySlug": categorySlug,
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	}, nil
}

// handleProductAdd handles product creation operations
func (s *ProductToolService) handleProductAdd(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 10, "extracting_product_data", nil)

	// Extract required parameters
	name := getStringParam(parameters, "name", "")
	description := getStringParam(parameters, "description", "")
	basePrice := getInt64Param(parameters, "base_price", 0)
	if basePrice == 0 {
		basePrice = getInt64Param(parameters, "price", 0)
	}
	userSellerID := getStringParam(parameters, "user_seller_id", "")
	if userSellerID == "" {
		userSellerID = getStringParam(parameters, "seller_id", "")
	}
	categoryID := getStringParam(parameters, "category_id", "")

	// Basic validation
	if name == "" {
		return nil, fmt.Errorf("name parameter is required")
	}
	if basePrice <= 0 {
		return nil, fmt.Errorf("base_price or price parameter is required and must be positive")
	}
	if userSellerID == "" {
		return nil, fmt.Errorf("user_seller_id or seller_id parameter is required")
	}
	if categoryID == "" {
		return nil, fmt.Errorf("category_id parameter is required")
	}

	s.sendProgress(streamChan, toolID, 30, "extracting_optional_parameters", nil)

	// Extract optional parameters with defaults
	categorySlug := getStringParam(parameters, "category_slug", "")
	brand := getStringParam(parameters, "brand", "")
	condition := getStringParam(parameters, "condition", "new")
	model := getStringParam(parameters, "model", "")
	tags := getStringSliceParam(parameters, "tags")
	manageStock := getBoolParam(parameters, "manage_stock", false)
	stock := getInt64Param(parameters, "stock", 0)
	sku := getStringParam(parameters, "sku", "")
	weight := getInt64Param(parameters, "weight", 0)
	height := getInt64Param(parameters, "height", 0)
	width := getInt64Param(parameters, "width", 0)
	depth := getInt64Param(parameters, "depth", 0)
	status := getStringParam(parameters, "status", "active")
	negotiable := getBoolParam(parameters, "negotiable", false)
	userType := getStringParam(parameters, "user_type", "individual")
	middlemanService := getBoolParam(parameters, "middleman_service", false)
	shippingCost := getInt64Param(parameters, "shipping_cost", 0)
	hasVariants := getBoolParam(parameters, "has_variants", false)
	lat := float64(getFloat32Param(parameters, "lat", 0.0))
	lng := float64(getFloat32Param(parameters, "lng", 0.0))
	thumbnail := getStringParam(parameters, "thumbnail", "")

	// Generate product ID
	productID := fmt.Sprintf("product_%d", time.Now().UnixNano())

	s.sendProgress(streamChan, toolID, 60, "creating_product", map[string]interface{}{
		"product_id": productID,
		"name":       name,
		"price":      basePrice,
	})

	log.Printf("ProductToolService: Adding product - Name: %s, Price: %d, Seller: %s", name, basePrice, userSellerID)

	// Call repository Add method
	err := s.repository.Add(
		ctx,
		productID,
		name,
		description,
		basePrice,
		userSellerID,
		categoryID,
		categorySlug,
		brand,
		condition,
		model,
		tags,
		manageStock,
		stock,
		sku,
		[]models.Attribute{}, // Empty attributes for now
		weight,
		height,
		width,
		depth,
		status,
		negotiable,
		userType,
		middlemanService,
		shippingCost,
		hasVariants,
		[]models.Option{}, // Empty options for now
		lat,
		lng,
		thumbnail,
		models.ProductType,
	)

	if err != nil {
		return nil, fmt.Errorf("product creation failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "formatting_creation_result", nil)

	return map[string]interface{}{
		"entity_type": "products",
		"operation":   "add",
		"success":     true,
		"product_id":  productID,
		"name":        name,
		"price":       basePrice,
		"seller_id":   userSellerID,
		"category_id": categoryID,
		"status":      status,
		"message":     fmt.Sprintf("Product '%s' created successfully with ID: %s", name, productID),
	}, nil
}

// handleProductUpdate handles product update operations
func (s *ProductToolService) handleProductUpdate(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_update_data", nil)

	// Extract required parameters
	productID := getStringParam(parameters, "product_id", "")
	if productID == "" {
		productID = getStringParam(parameters, "id", "")
	}
	if productID == "" {
		return nil, fmt.Errorf("product_id or id parameter is required")
	}

	newPrice := getInt64Param(parameters, "new_price", 0)
	if newPrice == 0 {
		newPrice = getInt64Param(parameters, "price", 0)
	}
	if newPrice == 0 {
		newPrice = getInt64Param(parameters, "base_price", 0)
	}
	if newPrice <= 0 {
		return nil, fmt.Errorf("new_price, price, or base_price parameter is required and must be positive")
	}

	s.sendProgress(streamChan, toolID, 60, "updating_product", map[string]interface{}{
		"product_id": productID,
		"new_price":  newPrice,
	})

	log.Printf("ProductToolService: Updating product %s with new price: %d", productID, newPrice)

	// Call repository Update method
	err := s.repository.Update(ctx, productID, newPrice)
	if err != nil {
		return nil, fmt.Errorf("product update failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "formatting_update_result", nil)

	return map[string]interface{}{
		"entity_type": "products",
		"operation":   "update",
		"success":     true,
		"product_id":  productID,
		"new_price":   newPrice,
		"message":     fmt.Sprintf("Product %s updated successfully with new price: %d", productID, newPrice),
	}, nil
}

// handleProductRemove handles product removal operations
func (s *ProductToolService) handleProductRemove(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	s.sendProgress(streamChan, toolID, 25, "extracting_removal_data", nil)

	// Extract required parameters
	productID := getStringParam(parameters, "product_id", "")
	if productID == "" {
		productID = getStringParam(parameters, "id", "")
	}
	if productID == "" {
		return nil, fmt.Errorf("product_id or id parameter is required")
	}

	s.sendProgress(streamChan, toolID, 60, "removing_product", map[string]interface{}{
		"product_id": productID,
	})

	log.Printf("ProductToolService: Removing product with ID: %s", productID)

	// Call repository Remove method
	err := s.repository.Remove(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product removal failed: %w", err)
	}

	s.sendProgress(streamChan, toolID, 90, "formatting_removal_result", nil)

	return map[string]interface{}{
		"entity_type": "products",
		"operation":   "remove",
		"success":     true,
		"product_id":  productID,
		"message":     fmt.Sprintf("Product %s removed successfully", productID),
	}, nil
}

// Helper functions are defined in helpers.go

// Helper methods for streaming
func (s *ProductToolService) sendProgress(streamChan chan<- ToolExecutionStream, toolID string, progress float64, step string, metadata map[string]interface{}) {
	if streamChan == nil {
		return
	}

	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["step"] = step

	streamChan <- ToolExecutionStream{
		ID:        toolID,
		ToolName:  "product_operation",
		Status:    "progress",
		Progress:  progress,
		Metadata:  metadata,
		Timestamp: time.Now().Unix(),
	}
}

func (s *ProductToolService) sendError(streamChan chan<- ToolExecutionStream, toolID string, err error, metadata map[string]interface{}) {
	if streamChan == nil {
		return
	}

	streamChan <- ToolExecutionStream{
		ID:        toolID,
		ToolName:  "product_operation",
		Status:    "error",
		Error:     err.Error(),
		Metadata:  metadata,
		Timestamp: time.Now().Unix(),
	}
}

func (s *ProductToolService) sendCompletion(streamChan chan<- ToolExecutionStream, toolID string, result interface{}, metadata map[string]interface{}) {
	if streamChan == nil {
		return
	}

	streamChan <- ToolExecutionStream{
		ID:        toolID,
		ToolName:  "product_operation",
		Status:    "completed",
		Progress:  100.0,
		Result:    result,
		Metadata:  metadata,
		Timestamp: time.Now().Unix(),
	}
}

func (s *ProductToolService) getResultCount(result interface{}) int {
	if resultMap, ok := result.(map[string]interface{}); ok {
		if count, exists := resultMap["count"].(int); exists {
			return count
		}
	}
	return 0
}

// productToMap converts a domain.Product to a serializable map
func (s *ProductToolService) productToMap(product *models.Product) map[string]interface{} {
	if product == nil {
		log.Printf("[PRODUCT_TO_MAP] WARNING: Received nil product")
		return nil
	}

	log.Printf("[PRODUCT_TO_MAP] Converting product: ID=%s, Name=%s", product.ProductID, product.Name)

	productMap := map[string]interface{}{
		"id":                product.ProductID,
		"name":              product.Name,
		"description":       product.Description,
		"base_price":        product.BasePrice,
		"user_seller_id":    product.UserSellerID,
		"category_id":       product.CategoryID,
		"category_slug":     product.CategorySlug,
		"brand":             product.Brand,
		"condition":         product.Condition,
		"model":             product.Model,
		"sku":               product.SKU,
		"manage_stock":      product.ManageStock,
		"stock":             product.Stock,
		"status":            product.Status,
		"negotiable":        product.Negotiable,
		"shipping_cost":     product.ShippingCost,
		"middleman_service": product.MiddlemanService,
		"has_variants":      product.HasVariants,
		"user_type":         product.UserType,
	}

	// Log core fields
	log.Printf("[PRODUCT_TO_MAP] Core fields - ID: %s, Name: %s, Price: %d, Status: %s",
		productMap["id"], productMap["name"], productMap["base_price"], productMap["status"])

	// Add optional fields if they exist
	if product.Thumbnail != "" {
		productMap["thumbnail"] = product.Thumbnail
		log.Printf("[PRODUCT_TO_MAP] Added thumbnail: %s", product.Thumbnail)
	}
	if product.Weight > 0 {
		productMap["weight"] = product.Weight
	}
	if product.Height > 0 {
		productMap["height"] = product.Height
	}
	if product.Width > 0 {
		productMap["width"] = product.Width
	}
	if product.Depth > 0 {
		productMap["depth"] = product.Depth
	}
	if product.Lat != 0 {
		productMap["lat"] = product.Lat
	}
	if product.Lng != 0 {
		productMap["lng"] = product.Lng
	}
	if len(product.Tags) > 0 {
		productMap["tags"] = product.Tags
		log.Printf("[PRODUCT_TO_MAP] Added %d tags", len(product.Tags))
	}
	if len(product.Options) > 0 {
		productMap["options"] = product.Options
		log.Printf("[PRODUCT_TO_MAP] Added %d options", len(product.Options))
	}
	if len(product.Attributes) > 0 {
		productMap["attributes"] = product.Attributes
		log.Printf("[PRODUCT_TO_MAP] Added %d attributes", len(product.Attributes))
	}

	// Verify the map is properly serializable
	mapKeys := make([]string, 0, len(productMap))
	for k := range productMap {
		mapKeys = append(mapKeys, k)
	}
	log.Printf("[PRODUCT_TO_MAP] Final map has %d keys: %v", len(productMap), mapKeys)

	// Verify no pointers are being returned
	for key, value := range productMap {
		valueType := fmt.Sprintf("%T", value)
		if strings.Contains(valueType, "*") {
			log.Printf("[PRODUCT_TO_MAP] WARNING: Field '%s' contains pointer type: %s", key, valueType)
		}
	}

	return productMap
}
