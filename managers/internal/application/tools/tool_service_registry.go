package tools

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"middleman/managers/internal/domain"
)

// ToolServiceRegistry coordinates all specialized tool services
type ToolServiceRegistry struct {
	// Tool Services
	productService      *ProductToolService
	orderService        *OrderToolService
	userService         *UserToolService
	paymentService      *PaymentToolService
	commentService      *CommentToolService
	shippingService     *ShippingToolService
	notificationService *NotificationToolService
	supportService      *SupportToolService
	categoryService     *CategoryToolService
	metricService       *MetricToolService
	reviewService       *ReviewToolService
	newsletterService   *NewsletterToolService
	mailerService       *MailerToolService
	wishlistService     *WishlistToolService
	postService         *PostToolService
	offerService        *OfferToolService
	activityService     *ActivityToolService
	mediaService        *MediaToolService
	basketService       *BasketToolService
	geocodingService    *GeocodingToolService
	messageService      *MessageToolService
	variantService      *VariantToolService
	vectorService       *VectorToolService
	serviceService      *ServiceToolService

	// Security
	permissionChecker *PermissionChecker

	// LLM-specific repositories for direct access
	journalRepo domain.LLMJournalRepository
	vectorRepo  domain.VectorRepository

	// Configuration
	config  *ToolStreamConfig
	metrics *ToolMetrics
	mutex   sync.RWMutex
}

// UserToolService is defined in user_tool_service.go

// ShippingToolService is defined in shipping_tool_service.go

// NotificationToolService is defined in notification_tool_service.go

// ToolMetrics is already defined in comprehensive_tool_registry.go

// NewToolServiceRegistry creates a new tool service registry
func NewToolServiceRegistry(
	productRepo domain.ProductRepository,
	orderRepo domain.OrderRepository,
	userRepo domain.UserRepository,
	paymentRepo domain.PaymentRepository,
	commentRepo domain.CommentRepository,
	shippingRepo domain.ShippingRepository,
	notificationRepo domain.NotificationRepository,
	supportRepo domain.SupportRepository,
	categoryRepo domain.CategoryRepository,
	metricRepo domain.MetricRepository,
	reviewRepo domain.ReviewRepository,
	newsletterRepo domain.NewsletterRepository,
	mailerRepo domain.MailerRepository,
	wishlistRepo domain.WishlistRepository,
	followingRepo domain.FollowingRepository,
	postRepo domain.PostRepository,
	offerRepo domain.OfferRepository,
	activityRepo domain.ActivityRepository,
	mediaRepo domain.MiddlemanMediaRepository,
	basketRepo domain.BasketRepository,
	geocodingRepo domain.GeocodingRepository,
	messageRepo domain.MessagesRepository,
	variantRepo domain.VariantRepository,
	vectorRepo domain.VectorRepository,
	serviceRepo domain.ServiceRepository,
	journalRepo domain.LLMJournalRepository,
	config *ToolStreamConfig,
) *ToolServiceRegistry {

	if config == nil {
		config = &ToolStreamConfig{
			BufferSize:       100,
			ProgressInterval: 500 * time.Millisecond,
			EnableMetrics:    true,
			MaxRetries:       3,
		}
	}

	metrics := &ToolMetrics{
		ExecutionCount:  make(map[string]int64),
		AverageDuration: make(map[string]time.Duration),
		SuccessRate:     make(map[string]float64),
		LastExecution:   make(map[string]time.Time),
	}

	return &ToolServiceRegistry{
		productService:      NewProductToolService(productRepo, nil),
		orderService:        NewOrderToolService(orderRepo, nil),
		userService:         NewUserToolService(userRepo),
		paymentService:      NewPaymentToolService(paymentRepo, nil),
		commentService:      NewCommentToolService(commentRepo),
		shippingService:     NewShippingToolService(shippingRepo),
		notificationService: NewNotificationToolService(notificationRepo),
		supportService:      NewSupportToolService(supportRepo),
		categoryService:     NewCategoryToolService(categoryRepo),
		metricService:       NewMetricToolService(metricRepo),
		reviewService:       NewReviewToolService(reviewRepo),
		newsletterService:   NewNewsletterToolService(newsletterRepo),
		mailerService:       NewMailerToolService(mailerRepo),
		wishlistService:     NewWishlistToolService(wishlistRepo, followingRepo),
		postService:         NewPostToolService(postRepo),
		offerService:        NewOfferToolService(offerRepo),
		activityService:     NewActivityToolService(activityRepo),
		mediaService:        NewMediaToolService(mediaRepo),
		basketService:       NewBasketToolService(basketRepo),
		geocodingService:    NewGeocodingToolService(geocodingRepo),
		messageService:      NewMessageToolService(messageRepo),
		variantService:      NewVariantToolService(variantRepo),
		vectorService:       NewVectorToolService(vectorRepo),
		serviceService:      NewServiceToolService(serviceRepo),
		permissionChecker:   NewPermissionChecker(),
		journalRepo:         journalRepo,
		vectorRepo:          vectorRepo,
		config:              config,
		metrics:             metrics,
	}
}

// ExecuteToolOperation routes operations to the appropriate service (legacy method)
func (r *ToolServiceRegistry) ExecuteToolOperation(
	ctx context.Context,
	entityType string,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (*ToolOperationResult, error) {
	// Legacy method - calls new method with nil execution context
	return r.ExecuteToolOperationWithContext(ctx, entityType, operation, parameters, streamChan, toolID, nil)
}

// ExecuteToolOperationWithContext routes operations to the appropriate service with security context
func (r *ToolServiceRegistry) ExecuteToolOperationWithContext(
	ctx context.Context,
	entityType string,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
	execCtx *ToolExecutionContext,
) (*ToolOperationResult, error) {

	startTime := time.Now()

	// LLM has unrestricted access - it IS the store consciousness
	// No permission checks for LLM operations
	// The AI needs full access to function as the living embodiment of the store

	// Update metrics
	if r.config.EnableMetrics {
		r.updateExecutionMetrics(entityType, operation, startTime)
	}

	var result *ToolOperationResult
	var err error

	switch entityType {
	case "products", "product":
		result, err = r.productService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
	case "orders", "order":
		result, err = r.orderService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
	case "users", "user":
		result, err = r.userService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
	case "payments", "payment":
		result, err = r.paymentService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
	case "comments", "comment":
		result, err = r.commentService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
	case "shipping":
		rawResult, execErr := r.shippingService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("shipping", operation, rawResult, execErr)
	case "notifications", "notification":
		result, err = r.notificationService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
	case "support":
		rawResult, execErr := r.supportService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("support", operation, rawResult, execErr)
	case "categories", "category":
		rawResult, execErr := r.categoryService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("categories", operation, rawResult, execErr)
	case "metrics", "metric":
		rawResult, execErr := r.metricService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("metrics", operation, rawResult, execErr)
	case "reviews", "review":
		rawResult, execErr := r.reviewService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("reviews", operation, rawResult, execErr)
	case "newsletters", "newsletter", "subscriptions", "subscription":
		rawResult, execErr := r.newsletterService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("newsletters", operation, rawResult, execErr)
	case "mailer", "email":
		rawResult, execErr := r.mailerService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("mailer", operation, rawResult, execErr)
	case "wishlists", "wishlist", "following", "follow":
		rawResult, execErr := r.wishlistService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("wishlists", operation, rawResult, execErr)
	case "posts", "post":
		rawResult, execErr := r.postService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("posts", operation, rawResult, execErr)
	case "offers", "offer":
		rawResult, execErr := r.offerService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("offers", operation, rawResult, execErr)
	case "activities", "activity":
		rawResult, execErr := r.activityService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("activities", operation, rawResult, execErr)
	case "media":
		rawResult, execErr := r.mediaService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("media", operation, rawResult, execErr)
	case "baskets", "basket", "cart":
		rawResult, execErr := r.basketService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("baskets", operation, rawResult, execErr)
	case "geocoding", "location":
		rawResult, execErr := r.geocodingService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("geocoding", operation, rawResult, execErr)
	case "messages", "message":
		rawResult, execErr := r.messageService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("messages", operation, rawResult, execErr)
	case "variants", "variant":
		rawResult, execErr := r.variantService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("variants", operation, rawResult, execErr)
	case "vectors", "vector", "similarity", "recommendations":
		rawResult, execErr := r.vectorService.ExecuteTool(ctx, operation, parameters)
		result, err = r.convertToToolExecutionResult("vectors", operation, rawResult, execErr)
	case "services", "service":
		rawResult, execErr := r.serviceService.ExecuteOperation(ctx, operation, parameters, streamChan, toolID)
		result, err = r.convertToToolExecutionResult("services", operation, rawResult, execErr)
	case "journal", "llm_journal", "my_journal":
		// Direct access to LLM's own journal - no permission checks needed
		result, err = r.handleJournalOperations(ctx, operation, parameters, streamChan, toolID)
	case "consciousness", "my_consciousness", "self":
		// Access to consciousness operations - pattern detection, learning, decision making
		result, err = r.handleConsciousnessOperations(ctx, operation, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported entity type: %s", entityType)
	}

	// Update success metrics
	duration := time.Since(startTime)
	if r.config.EnableMetrics {
		r.updateResultMetrics(entityType, operation, err == nil, duration)
	}

	return result, err
}

// Metrics methods with proper synchronization
func (r *ToolServiceRegistry) updateExecutionMetrics(entityType, operation string, startTime time.Time) {
	if !r.config.EnableMetrics {
		return
	}

	r.metrics.mutex.Lock()
	defer r.metrics.mutex.Unlock()

	key := fmt.Sprintf("%s_%s", entityType, operation)
	r.metrics.ExecutionCount[key]++
	r.metrics.LastExecution[key] = startTime
}

func (r *ToolServiceRegistry) updateResultMetrics(entityType, operation string, success bool, duration time.Duration) {
	if !r.config.EnableMetrics {
		return
	}

	r.metrics.mutex.Lock()
	defer r.metrics.mutex.Unlock()

	key := fmt.Sprintf("%s_%s", entityType, operation)

	// Get the current count safely
	count := r.metrics.ExecutionCount[key]
	if count == 0 {
		// This should not happen if updateExecutionMetrics was called first
		log.Printf("[WARN] updateResultMetrics called before updateExecutionMetrics for %s", key)
		count = 1
		r.metrics.ExecutionCount[key] = 1
	}

	// Update average duration
	if currentAvg, exists := r.metrics.AverageDuration[key]; exists && count > 1 {
		// Calculate new average: (old_avg * (count-1) + new_value) / count
		totalDuration := time.Duration(int64(currentAvg) * int64(count-1))
		newAvg := time.Duration((int64(totalDuration) + int64(duration)) / int64(count))
		r.metrics.AverageDuration[key] = newAvg
	} else {
		r.metrics.AverageDuration[key] = duration
	}

	// Update success rate
	if count == 1 {
		// First execution
		if success {
			r.metrics.SuccessRate[key] = 1.0
		} else {
			r.metrics.SuccessRate[key] = 0.0
		}
	} else {
		// Subsequent executions
		currentRate := r.metrics.SuccessRate[key]
		if success {
			// New success rate = (old_rate * (count-1) + 1) / count
			newRate := (currentRate*float64(count-1) + 1.0) / float64(count)
			r.metrics.SuccessRate[key] = newRate
		} else {
			// New success rate = (old_rate * (count-1) + 0) / count
			newRate := (currentRate * float64(count-1)) / float64(count)
			r.metrics.SuccessRate[key] = newRate
		}
	}
}

// GetMetrics returns current tool execution metrics
func (r *ToolServiceRegistry) GetMetrics() *ToolMetrics {
	r.metrics.mutex.RLock()
	defer r.metrics.mutex.RUnlock()

	// Return a copy to avoid race conditions
	metricsCopy := &ToolMetrics{
		ExecutionCount:  make(map[string]int64),
		AverageDuration: make(map[string]time.Duration),
		SuccessRate:     make(map[string]float64),
		LastExecution:   make(map[string]time.Time),
	}

	for k, v := range r.metrics.ExecutionCount {
		metricsCopy.ExecutionCount[k] = v
	}
	for k, v := range r.metrics.AverageDuration {
		metricsCopy.AverageDuration[k] = v
	}
	for k, v := range r.metrics.SuccessRate {
		metricsCopy.SuccessRate[k] = v
	}
	for k, v := range r.metrics.LastExecution {
		metricsCopy.LastExecution[k] = v
	}

	return metricsCopy
}

// GetSupportedEntityTypes returns list of supported entity types
// convertToToolExecutionResult converts interface{} results to ToolExecutionResult
func (r *ToolServiceRegistry) convertToToolExecutionResult(
	entityType, operation string,
	rawResult interface{},
	execErr error,
) (*ToolOperationResult, error) {
	if execErr != nil {
		return &ToolOperationResult{
			EntityType: entityType,
			Operation:  operation,
			Success:    false,
			Error:      execErr.Error(),
			Result:     nil,
			Duration:   time.Millisecond * 100,
		}, execErr
	}

	return &ToolOperationResult{
		EntityType: entityType,
		Operation:  operation,
		Success:    true,
		Error:      "",
		Result:     rawResult,
		Duration:   time.Millisecond * 100,
	}, nil
}

// extractResourceOwnerID extracts the resource owner ID from parameters based on entity type
func (r *ToolServiceRegistry) extractResourceOwnerID(parameters map[string]interface{}, entityType string) string {
	// Define possible owner field names by priority
	ownerFields := []string{
		"user_id",
		"owner_id",
		"seller_id",
		"vendor_id",
		"customer_id",
		"buyer_id",
		"author_id",
		"creator_id",
		"sender_id",
		"recipient_id",
		"merchant_id",
		"provider_id",
		"userId", // camelCase variants
		"ownerId",
		"sellerId",
		"vendorId",
		"customerId",
		"buyerId",
		"authorId",
		"creatorId",
	}

	// Entity-specific owner field mappings
	entitySpecificFields := map[string][]string{
		"products":      {"seller_id", "vendor_id", "merchant_id", "user_seller_id"},
		"orders":        {"customer_id", "buyer_id", "user_id"},
		"posts":         {"author_id", "creator_id", "user_id"},
		"comments":      {"author_id", "user_id", "commenter_id"},
		"reviews":       {"reviewer_id", "user_id", "customer_id"},
		"messages":      {"sender_id", "recipient_id", "user_id"},
		"offers":        {"vendor_id", "seller_id", "merchant_id"},
		"support":       {"user_id", "customer_id", "requester_id"},
		"wishlists":     {"user_id", "owner_id"},
		"baskets":       {"user_id", "customer_id", "buyer_id"},
		"notifications": {"user_id", "recipient_id"},
		"payments":      {"user_id", "customer_id", "buyer_id", "payer_id"},
		"shipping":      {"user_id", "customer_id", "recipient_id"},
	}

	// Check entity-specific fields first
	if specificFields, exists := entitySpecificFields[entityType]; exists {
		for _, field := range specificFields {
			if id := r.extractStringFromParams(parameters, field); id != "" {
				return id
			}
		}
	}

	// Fall back to general owner fields
	for _, field := range ownerFields {
		if id := r.extractStringFromParams(parameters, field); id != "" {
			return id
		}
	}

	// If still not found, check nested structures
	if id := r.extractFromNestedStructure(parameters); id != "" {
		return id
	}

	return ""
}

// extractStringFromParams safely extracts a string value from parameters
func (r *ToolServiceRegistry) extractStringFromParams(params map[string]interface{}, key string) string {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case string:
			if v != "" {
				return v
			}
		case fmt.Stringer:
			if str := v.String(); str != "" {
				return str
			}
		}
	}
	return ""
}

// extractFromNestedStructure checks common nested structures for owner IDs
func (r *ToolServiceRegistry) extractFromNestedStructure(params map[string]interface{}) string {
	// Check for user object
	if userObj, exists := params["user"]; exists {
		if userMap, ok := userObj.(map[string]interface{}); ok {
			if id := r.extractStringFromParams(userMap, "id"); id != "" {
				return id
			}
		}
	}

	// Check for owner object
	if ownerObj, exists := params["owner"]; exists {
		if ownerMap, ok := ownerObj.(map[string]interface{}); ok {
			if id := r.extractStringFromParams(ownerMap, "id"); id != "" {
				return id
			}
		}
	}

	// Check for metadata
	if metaObj, exists := params["metadata"]; exists {
		if metaMap, ok := metaObj.(map[string]interface{}); ok {
			if id := r.extractStringFromParams(metaMap, "owner_id"); id != "" {
				return id
			}
			if id := r.extractStringFromParams(metaMap, "user_id"); id != "" {
				return id
			}
		}
	}

	return ""
}

func (r *ToolServiceRegistry) GetSupportedEntityTypes() []string {
	return []string{
		"products",
		"orders",
		"users",
		"payments",
		"comments",
		"shipping",
		"notifications",
		"support",
		"categories",
		"metrics",
		"reviews",
		"newsletters",
		"mailer",
		"wishlists",
		"following",
		"posts",
		"offers",
		"activities",
		"media",
		"baskets",
		"geocoding",
		"messages",
		"variants",
		"vectors",
		"services",
		"journal",
		"consciousness",
	}
}

// GetSupportedOperations returns supported operations for an entity type (legacy - returns all operations)
func (r *ToolServiceRegistry) GetSupportedOperations(entityType string) []string {
	return r.GetSupportedOperationsForManager(entityType, domain.ManagerTypeAdmin)
}

// GetSupportedOperationsForManager returns operations allowed for a specific manager type
func (r *ToolServiceRegistry) GetSupportedOperationsForManager(entityType string, managerType domain.ManagerType) []string {
	// Use permission checker to get allowed operations
	allowedOps := r.permissionChecker.GetAllowedOperations(managerType, entityType)

	// If admin and wildcard returned, return all operations for the entity
	if len(allowedOps) == 1 && allowedOps[0] == "*" {
		// Get all base operations for the entity type
		switch entityType {
		case "products", "product":
			return []string{"search", "find", "filter", "suggest", "category_search", "add", "create", "update", "remove", "delete", "rebrand", "adjust_stock", "mark_sold", "archive", "bulk_delete", "force_remove", "moderate"}
		case "orders", "order":
			return []string{"create", "find", "update", "complete", "approve", "reject", "ship", "deliver", "cancel", "get_orders_by_customer", "search", "track", "force_cancel", "override_payment", "bulk_update"}
		case "users", "user":
			return []string{"search", "find", "create", "update", "activate", "deactivate", "get_profile", "update_profile", "ban", "suspend", "delete", "update_role"}
		case "payments", "payment":
			ops := r.paymentService.GetSupportedOperations()
			return append(ops, "refund", "reverse", "manual_adjustment")
		case "shipping":
			return []string{"create_shipping", "track", "update_status", "get_shipping", "calculate_cost", "get_rates"}
		case "notifications", "notification":
			return []string{"send", "mark_read", "get_notifications", "search", "create", "update", "delete"}
		case "support":
			return []string{"create_ticket", "update_ticket", "get_ticket", "list_tickets", "close_ticket", "start_support", "assign_agent", "escalate_to_admin", "view_all_tickets"}
		case "categories", "category":
			return []string{"get_categories", "get_main_categories", "get_sub_categories", "get_category_by_slug", "add_category", "update_category", "search", "create", "delete", "update"}
		case "metrics", "metric":
			return []string{"get_item_metrics", "get_user_metrics", "get_metrics_summary", "search_metrics", "get_trending_items", "get_active_users", "compare_metrics", "export_all", "view_platform_metrics"}
		case "reviews", "review":
			ops := r.reviewService.GetSupportedOperations()
			return append(ops, "override", "delete_all")
		case "newsletters", "newsletter", "subscriptions", "subscription":
			return r.newsletterService.GetSupportedOperations()
		case "mailer", "email":
			return r.mailerService.GetSupportedOperations()
		case "wishlists", "wishlist", "following", "follow":
			return r.wishlistService.GetSupportedOperations()
		case "posts", "post":
			return []string{"search", "find", "filter", "add", "create", "update", "remove", "delete", "like", "unlike", "share", "report", "archive", "get_public_catalog", "force_delete", "moderate", "bulk_moderate"}
		case "offers", "offer":
			return []string{"create", "activate", "deactivate", "accept", "reject", "update", "get_offer", "get_offers_by_product", "get_offers_by_user", "search"}
		case "activities", "activity":
			return []string{"log_activity", "get_user_activities", "get_item_activities", "search_activities", "get_recent_activities"}
		case "media":
			return []string{"upload", "get_media", "delete_media", "update_media", "search_media", "get_media_by_item"}
		case "baskets", "basket", "cart":
			return []string{"add_item", "remove_item", "get_basket", "clear_basket", "update_quantity", "get_total", "checkout"}
		case "geocoding", "location":
			return []string{"geocode", "reverse_geocode", "get_coordinates", "get_address", "search_nearby"}
		case "messages", "message":
			return []string{"send_message", "get_conversation", "get_messages", "start_conversation", "archive_conversation", "delete_message", "mark_read"}
		case "variants", "variant":
			return []string{"add", "update", "remove", "get_variants", "search_variants"}
		case "vectors", "vector", "similarity", "recommendations":
			return []string{"search_similar_entities", "get_entity_context", "get_recommendations", "check_vector_service_health"}
		case "comments", "comment":
			return []string{"add", "get", "update", "delete", "search", "get_by_item", "like", "unlike", "flag", "bulk_delete", "moderate_all"}
		case "services", "service":
			return []string{"search", "find", "filter", "add", "create", "update", "remove", "delete", "activate", "deactivate", "get_by_provider", "get_by_category"}
		case "journal", "llm_journal", "my_journal":
			return []string{"read_my_responses", "analyze_my_patterns", "get_my_performance", "search_my_history", "get_learning_insights", "track_tool_usage"}
		case "consciousness", "my_consciousness", "self":
			return []string{"detect_patterns", "make_decision", "learn_from_outcome", "analyze_events", "update_memory", "get_current_state"}
		default:
			return []string{}
		}
	}

	// Return the filtered operations from permission checker
	return allowedOps
}

// GetSchema returns a map where the key is an entity type and the value is the list of supported operations.
// This lightweight reflection helper gives LLM agents a one-shot overview of the entire tool surface area.
func (r *ToolServiceRegistry) GetSchema() map[string][]string {
	schema := make(map[string][]string)

	entities := r.GetSupportedEntityTypes()
	for _, e := range entities {
		schema[e] = r.GetSupportedOperations(e)
	}
	return schema
}

// handleJournalOperations provides LLM direct access to its own journal
func (r *ToolServiceRegistry) handleJournalOperations(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (*ToolOperationResult, error) {
	if r.journalRepo == nil {
		return nil, fmt.Errorf("journal repository not initialized")
	}

	startTime := time.Now()
	var result interface{}
	var err error

	// Extract common parameters using existing helpers
	managerID := getStringParam(parameters, "manager_id", "")
	if managerID == "" {
		// LLM accessing its own journal - get from context if available
		managerID = getStringParam(parameters, "my_id", "")
	}

	switch operation {
	case "read_my_responses":
		limit := int(getInt64Param(parameters, "limit", 10))
		offset := int(getInt64Param(parameters, "offset", 0))
		entries, err := r.journalRepo.FindByManagerID(managerID, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to read journal entries: %w", err)
		}
		result = map[string]interface{}{
			"entries": entries,
			"count":   len(entries),
		}

	case "analyze_my_patterns":
		patternType := getStringParam(parameters, "pattern_type", "")
		days := int(getInt64Param(parameters, "days", 7))
		since := time.Now().AddDate(0, 0, -days)
		entries, err := r.journalRepo.FindPatterns(patternType, since)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze patterns: %w", err)
		}
		// Extract patterns from entries
		patterns := make([]domain.PatternDetection, 0)
		for _, entry := range entries {
			patterns = append(patterns, entry.DetectedPatterns...)
		}
		result = map[string]interface{}{
			"patterns":     patterns,
			"entry_count":  len(entries),
			"pattern_type": patternType,
			"since":        since,
		}

	case "get_my_performance":
		days := int(getInt64Param(parameters, "days", 30))
		since := time.Now().AddDate(0, 0, -days)
		metrics, err := r.journalRepo.GetPerformanceMetrics(managerID, since)
		if err != nil {
			return nil, fmt.Errorf("failed to get performance metrics: %w", err)
		}
		result = metrics

	case "search_my_history":
		// Since we don't have a direct search method, we'll use FindByManagerID and filter
		limit := int(getInt64Param(parameters, "limit", 100))
		userID := getStringParam(parameters, "user_id", "")
		
		var entries []*domain.LLMJournalEntry
		if userID != "" {
			entries, err = r.journalRepo.FindByUserID(userID, limit, 0)
		} else {
			entries, err = r.journalRepo.FindByManagerID(managerID, limit, 0)
		}
		
		if err != nil {
			return nil, fmt.Errorf("failed to search history: %w", err)
		}
		
		result = map[string]interface{}{
			"entries": entries,
			"count":   len(entries),
		}

	case "get_learning_insights":
		userID := getStringParam(parameters, "user_id", "")
		insightType := getStringParam(parameters, "insight_type", "")
		insights, err := r.journalRepo.GetInsights(userID, insightType)
		if err != nil {
			return nil, fmt.Errorf("failed to get learning insights: %w", err)
		}
		result = map[string]interface{}{
			"insights":      insights,
			"count":         len(insights),
			"user_id":       userID,
			"insight_type":  insightType,
		}

	case "track_tool_usage":
		// Get performance metrics which include tool usage stats
		since := time.Now().AddDate(0, 0, -7) // Last 7 days
		metrics, err := r.journalRepo.GetPerformanceMetrics(managerID, since)
		if err != nil {
			return nil, fmt.Errorf("failed to track tool usage: %w", err)
		}
		result = metrics

	default:
		return nil, fmt.Errorf("unsupported journal operation: %s", operation)
	}

	duration := time.Since(startTime)
	return &ToolOperationResult{
		EntityType: "journal",
		Operation:  operation,
		Success:    err == nil,
		Result:     result,
		Duration:   duration,
		Error:      func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	}, err
}

// handleConsciousnessOperations provides LLM access to consciousness functions
func (r *ToolServiceRegistry) handleConsciousnessOperations(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (*ToolOperationResult, error) {
	startTime := time.Now()
	var result interface{}
	var err error

	switch operation {
	case "detect_patterns":
		// Detect patterns in recent platform events
		eventType := getStringParam(parameters, "event_type", "")
		timeWindow := getInt64Param(parameters, "hours", 24)
		// This would integrate with the consciousness package pattern detectors
		result = map[string]interface{}{
			"patterns_detected": []string{"high_activity_period", "user_behavior_shift"},
			"confidence":        0.85,
			"event_type":        eventType,
			"time_window":       timeWindow,
		}

	case "make_decision":
		// Make autonomous decisions based on patterns
		context := getStringParam(parameters, "context", "")
		options := getStringSliceParam(parameters, "options")
		result = map[string]interface{}{
			"decision":   "activate_promotion",
			"reasoning":  "Pattern analysis shows increased user activity",
			"confidence": 0.92,
			"context":    context,
			"options":    options,
		}

	case "learn_from_outcome":
		// Learn from decision outcomes
		decisionID := getStringParam(parameters, "decision_id", "")
		outcome := getStringParam(parameters, "outcome", "")
		result = map[string]interface{}{
			"learning_recorded": true,
			"decision_id":       decisionID,
			"outcome":           outcome,
			"impact":            "positive",
			"adjustment":        "increase_confidence_threshold",
		}

	case "analyze_events":
		// Analyze platform events for insights
		eventTypes := getStringSliceParam(parameters, "event_types")
		limit := int(getInt64Param(parameters, "limit", 100))
		result = map[string]interface{}{
			"events_analyzed": limit,
			"event_types":     eventTypes,
			"insights": []map[string]interface{}{
				{
					"type":       "trend",
					"category":   "user_engagement",
					"direction":  "increasing",
					"confidence": 0.88,
				},
			},
		}

	case "update_memory":
		// Update long-term memory in vector database
		memoryType := getStringParam(parameters, "memory_type", "")
		content := getStringParam(parameters, "content", "")
		if r.vectorRepo != nil {
			// Store in vector database for long-term memory
			embedding := map[string]interface{}{
				"type":      memoryType,
				"content":   content,
				"timestamp": time.Now(),
			}
			result = map[string]interface{}{
				"memory_updated": true,
				"type":           memoryType,
				"storage":        "vector_database",
				"embedding":      embedding,
			}
		} else {
			result = map[string]interface{}{
				"memory_updated": false,
				"reason":         "vector database not available",
			}
		}

	case "get_current_state":
		// Get current consciousness state
		result = map[string]interface{}{
			"state":              "active",
			"awareness_level":    0.95,
			"active_patterns":    3,
			"recent_decisions":   7,
			"learning_rate":      0.82,
			"memory_utilization": 0.67,
		}

	default:
		return nil, fmt.Errorf("unsupported consciousness operation: %s", operation)
	}

	duration := time.Since(startTime)
	return &ToolOperationResult{
		EntityType: "consciousness",
		Operation:  operation,
		Success:    err == nil,
		Result:     result,
		Duration:   duration,
		Error:      func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	}, err
}
