// Package grpc provides gRPC server and client implementations
package grpc

// DEPRECATED: This file is kept for backward compatibility only.
// The modern architecture uses ToolServiceRegistry in internal/application/tools
// which provides the same functionality in a cleaner, more maintainable way.
//
// DO NOT ADD NEW FEATURES HERE - use ToolServiceRegistry instead.

import (
	"context"
	"fmt"
	"log"
	"middleman/managers/internal/application/processor"
	"middleman/managers/internal/application/services"
	"reflect"
	"strings"
	"time"

	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
)

// EnhancedRepositoryAdapter provides comprehensive access to all repository operations
type EnhancedRepositoryAdapter struct {
	// Repository instances
	productRepo domain.ProductRepository
	userRepo    domain.UserRepository

	orderRepo        domain.OrderRepository
	paymentRepo      domain.PaymentRepository
	offerRepo        domain.OfferRepository
	reviewRepo       domain.ReviewRepository
	commentRepo      domain.CommentRepository
	notificationRepo domain.NotificationRepository
	newsletterRepo   domain.NewsletterRepository
	basketRepo       domain.BasketRepository
	categoryRepo     domain.CategoryRepository
	metricRepo       domain.MetricRepository
	messagesRepo     domain.MessagesRepository
	wishlistRepo     domain.WishlistRepository
	followingRepo    domain.FollowingRepository
	activityRepo     domain.ActivityRepository
	mediaRepo        domain.MiddlemanMediaRepository
	shippingRepo     domain.ShippingRepository
	supportRepo      domain.SupportRepository
	geocodingRepo    domain.GeocodingRepository
	variantRepo      domain.VariantRepository
	serviceRepo      domain.ServiceRepository

	// Enhanced interface for LLM communication
	llmInterface *processor.EnhancedLLMInterface
}

// NewEnhancedRepositoryAdapter creates a new enhanced repository adapter with all repositories
func NewEnhancedRepositoryAdapter(
	productRepo domain.ProductRepository,
	userRepo domain.UserRepository,
	orderRepo domain.OrderRepository,
	paymentRepo domain.PaymentRepository,
	offerRepo domain.OfferRepository,
	reviewRepo domain.ReviewRepository,
	commentRepo domain.CommentRepository,
	notificationRepo domain.NotificationRepository,
	newsletterRepo domain.NewsletterRepository,
	basketRepo domain.BasketRepository,
	categoryRepo domain.CategoryRepository,
	metricRepo domain.MetricRepository,
	messagesRepo domain.MessagesRepository,
	wishlistRepo domain.WishlistRepository,
	followingRepo domain.FollowingRepository,
	activityRepo domain.ActivityRepository,
	mediaRepo domain.MiddlemanMediaRepository,
	shippingRepo domain.ShippingRepository,
	supportRepo domain.SupportRepository,
	geocodingRepo domain.GeocodingRepository,
	variantRepo domain.VariantRepository,
	serviceRepo domain.ServiceRepository,
) *EnhancedRepositoryAdapter {
	return &EnhancedRepositoryAdapter{
		productRepo:      productRepo,
		userRepo:         userRepo,
		orderRepo:        orderRepo,
		paymentRepo:      paymentRepo,
		offerRepo:        offerRepo,
		reviewRepo:       reviewRepo,
		commentRepo:      commentRepo,
		notificationRepo: notificationRepo,
		newsletterRepo:   newsletterRepo,
		basketRepo:       basketRepo,
		categoryRepo:     categoryRepo,
		metricRepo:       metricRepo,
		messagesRepo:     messagesRepo,
		wishlistRepo:     wishlistRepo,
		followingRepo:    followingRepo,
		activityRepo:     activityRepo,
		mediaRepo:        mediaRepo,
		shippingRepo:     shippingRepo,
		supportRepo:      supportRepo,
		geocodingRepo:    geocodingRepo,
		variantRepo:      variantRepo,
		serviceRepo:      serviceRepo,
		llmInterface:     processor.NewEnhancedLLMInterface(),
	}
}

// Execute implements the unified repository interface with comprehensive operation support
func (era *EnhancedRepositoryAdapter) Execute(ctx context.Context, query services.RepositoryQuery) (*services.RepositoryResponse, error) {
	startTime := time.Now()

	log.Printf("EnhancedRepositoryAdapter.Execute: Executing %s operation on %s", query.Operation, query.EntityType)

	// Validate the query using enhanced interface
	if err := era.llmInterface.ValidateOperationRequest(query.EntityType, string(query.Operation), era.queryParametersToMap(query.Parameters)); err != nil {
		return &services.RepositoryResponse{
			Success: false,
			Error:   fmt.Sprintf("validation failed: %v", err),
			Metadata: services.ResponseMetadata{
				EntityType:    query.EntityType,
				Operation:     query.Operation,
				ExecutionTime: time.Since(startTime).Milliseconds(),
			},
		}, err
	}

	// Route to appropriate repository
	result, err := era.routeToRepository(ctx, query)

	executionTime := time.Since(startTime).Milliseconds()

	if err != nil {
		return &services.RepositoryResponse{
			Success: false,
			Error:   err.Error(),
			Metadata: services.ResponseMetadata{
				EntityType:    query.EntityType,
				Operation:     query.Operation,
				ExecutionTime: executionTime,
			},
		}, err
	}

	// Determine result count
	resultCount := era.getResultCount(result)

	return &services.RepositoryResponse{
		Success: true,
		Data:    result,
		Metadata: services.ResponseMetadata{
			EntityType:    query.EntityType,
			Operation:     query.Operation,
			ReturnedCount: resultCount,
			ExecutionTime: executionTime,
		},
	}, nil
}

// routeToRepository routes queries to appropriate repository based on entity type
func (era *EnhancedRepositoryAdapter) routeToRepository(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	switch query.EntityType {
	case models.ProductType:
		return era.executeProductOperation(ctx, query)
	case models.UserEntityType:
		return era.executeUserOperation(ctx, query)
	case models.EntityTypeOrder:
		return era.executeOrderOperation(ctx, query)
	case models.OfferEntityType:
		return era.executeOfferOperation(ctx, query)
	case models.ReviewType:
		return era.executeReviewOperation(ctx, query)
	case models.CommentType:
		return era.executeCommentOperation(ctx, query)
	case models.NotificationEntityType:
		return era.executeNotificationOperation(ctx, query)
	case models.NewsletterEntityType:
		return era.executeNewsletterOperation(ctx, query)
	case models.BasketType:
		return era.executeBasketOperation(ctx, query)
	case models.CategoryType:
		return era.executeCategoryOperation(ctx, query)
	case models.MetricEntityType:
		return era.executeMetricOperation(ctx, query)
	case models.MessageType:
		return era.executeMessagesOperation(ctx, query)
	case models.WishlistType:
		return era.executeWishlistOperation(ctx, query)
	case models.FollowingType:
		return era.executeFollowingOperation(ctx, query)
	case models.ActivityType:
		return era.executeActivityOperation(ctx, query)
	case models.MediaType:
		return era.executeMediaOperation(ctx, query)
	case models.ServiceType:
		return era.executeServiceOperation(ctx, query)
	default:
		return nil, fmt.Errorf("unsupported entity type: %s", query.EntityType)
	}
}

// executeProductOperation handles all product-related operations with comprehensive validation
func (era *EnhancedRepositoryAdapter) executeProductOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	params := query.Parameters

	// Validate operation using enhanced interface
	if err := era.llmInterface.ValidateOperationRequest(query.EntityType, string(query.Operation), era.queryParametersToMap(params)); err != nil {
		return nil, fmt.Errorf("product operation validation failed: %v", err)
	}

	switch query.Operation {
	case services.OpFind:
		if params.ID == "" {
			return nil, fmt.Errorf("product ID is required for find operation")
		}
		return era.productRepo.Find(ctx, params.ID)

	case services.OpSearchByTerm, services.OpSearch:
		if params.SearchTerm == "" && params.Name == "" {
			return nil, fmt.Errorf("search term or name is required for search operation")
		}
		searchTerm := params.SearchTerm
		if searchTerm == "" {
			searchTerm = params.Name
		}
		return era.productRepo.SearchWithTerm(ctx, searchTerm)

	case services.OpFilter:
		return era.productRepo.SearchWithFilters(ctx,
			params.Name, params.CategoryID, params.CategorySlug, params.MinPrice, params.MaxPrice,
			params.Brand, params.Condition, params.Model, params.Tags,
			params.ManageStock, params.MinStock, params.MaxStock, params.SKU,
			params.Status, params.Negotiable, params.UserType, params.MiddlemanService,
			params.HasVariants, params.ShippingCost, params.MinWeight, params.MaxWeight,
			params.MinHeight, params.MaxHeight, params.MinWidth, params.MaxWidth,
			params.MinDepth, params.MaxDepth, params.Offset, params.Limit,
			params.Lat, params.Lng, int64(params.RadiusMeters), params.Page, params.PageSize,
			params.SortBy, params.SortOrder)

	case services.OpAdd:
		// Validate required fields for product creation
		if params.Name == "" || params.Description == "" || params.BasePrice == 0 || params.UserSellerID == "" || params.CategoryID == "" {
			return nil, fmt.Errorf("name, description, base_price, user_seller_id, and category_id are required for add operation")
		}

		return nil, era.productRepo.Add(ctx,
			params.ID, params.Name, params.Description, params.BasePrice,
			params.UserSellerID, params.CategoryID, params.CategorySlug,
			params.Brand, params.Condition, params.Model, params.Tags,
			params.ManageStock, params.NewStock, params.SKU,
			[]models.Attribute{}, // Empty attributes for now
			params.MinWeight, params.MinHeight, params.MinWidth, params.MinDepth,
			params.Status, params.Negotiable, params.UserType,
			params.MiddlemanService, params.ShippingCost, params.HasVariants,
			[]models.Option{}, // Empty options for now
			params.Lat, params.Lng, params.Thumbnail, models.ProductType)

	case services.OpUpdate:
		if params.ID == "" {
			return nil, fmt.Errorf("product ID is required for update operation")
		}
		return nil, era.productRepo.Update(ctx, params.ID, params.NewPrice)

	case services.OpRemove:
		if params.ID == "" {
			return nil, fmt.Errorf("product ID is required for remove operation")
		}
		return nil, era.productRepo.Remove(ctx, params.ID)

	case services.OpSuggest:
		if params.Name == "" {
			return nil, fmt.Errorf("name is required for suggest operation")
		}
		return era.productRepo.SuggestProducts(ctx, params.Name)

	case services.OpGetCatalog:
		return era.productRepo.GetCatalog(ctx, params.UserID, params.Page, params.PageSize, params.SortBy, params.SortOrder)

	case services.OpGetPublicCatalog:
		return era.productRepo.GetPublicCatalog(ctx, params.UserID, params.Page, params.PageSize, params.SortBy, params.SortOrder)

	default:
		return nil, fmt.Errorf("unsupported operation %s for products", query.Operation)
	}
}

// executeUserOperation handles all user-related operations
func (era *EnhancedRepositoryAdapter) executeUserOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	params := query.Parameters

	switch query.Operation {
	case services.OpFind:
		if params.ID == "" && params.UserID == "" {
			return nil, fmt.Errorf("user ID is required for find operation")
		}
		userID := params.ID
		if userID == "" {
			userID = params.UserID
		}
		return era.userRepo.Find(ctx, userID)

	case services.OpAdd:
		// User creation
		if params.Email == "" || params.Password == "" || params.FirstName == "" || params.LastName == "" {
			return nil, fmt.Errorf("email, password, first_name, and last_name are required for user creation")
		}
		return era.userRepo.CreateUser(ctx, params.Email, params.Password, params.Email, // username = email for now
			params.FirstName, params.LastName, params.Location,
			float32(params.Lat), float32(params.Lng), params.Thumbnail, params.Language)

	case services.OpUpdate:
		if params.ID == "" {
			return nil, fmt.Errorf("user ID is required for update operation")
		}
		return era.userRepo.UpdateUser(ctx, params.ID, params.Email, params.FirstName, params.LastName,
			params.Bio, params.Privacy, params.Background, params.Location,
			float32(params.Lat), float32(params.Lng), params.Thumbnail)

	default:
		return nil, fmt.Errorf("unsupported operation %s for users", query.Operation)
	}
}

func (era *EnhancedRepositoryAdapter) executeOrderOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for order operations
	return nil, fmt.Errorf("order operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeOfferOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for offer operations
	return nil, fmt.Errorf("offer operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeReviewOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for review operations
	return nil, fmt.Errorf("review operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeCommentOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for comment operations
	return nil, fmt.Errorf("comment operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeNotificationOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for notification operations
	return nil, fmt.Errorf("notification operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeNewsletterOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for newsletter operations
	return nil, fmt.Errorf("newsletter operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeBasketOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for basket operations
	return nil, fmt.Errorf("basket operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeCategoryOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for category operations
	return nil, fmt.Errorf("category operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeMetricOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for metric operations
	return nil, fmt.Errorf("metric operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeMessagesOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for messages operations
	return nil, fmt.Errorf("messages operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeWishlistOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for wishlist operations
	return nil, fmt.Errorf("wishlist operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeFollowingOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for following operations
	return nil, fmt.Errorf("following operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeActivityOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for activity operations
	return nil, fmt.Errorf("activity operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeMediaOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for media operations
	return nil, fmt.Errorf("media operations not yet implemented")
}

func (era *EnhancedRepositoryAdapter) executeServiceOperation(ctx context.Context, query services.RepositoryQuery) (interface{}, error) {
	// Implementation for service operations
	return nil, fmt.Errorf("service operations not yet implemented")
}

// Helper methods

// ValidateQuery validates if a query is properly formed
func (era *EnhancedRepositoryAdapter) ValidateQuery(query services.RepositoryQuery) error {
	if query.EntityType == "" {
		return fmt.Errorf("entity type is required")
	}
	if query.Operation == "" {
		return fmt.Errorf("operation is required")
	}
	return era.llmInterface.ValidateOperationRequest(query.EntityType, string(query.Operation), era.queryParametersToMap(query.Parameters))
}

// GetSupportedOperations returns supported operations for an entity type
func (era *EnhancedRepositoryAdapter) GetSupportedOperations(entityType models.EntityType) []services.OperationType {
	schema := era.llmInterface.SchemaRegistry.GetSchema(entityType)
	if schema == nil {
		return []services.OperationType{}
	}

	operations := make([]services.OperationType, len(schema.Operations))
	for i, op := range schema.Operations {
		operations[i] = services.OperationType(op.Name)
	}
	return operations
}

// TranslateAIRequest translates an AI request into a repository query
func (era *EnhancedRepositoryAdapter) TranslateAIRequest(aiRequest map[string]interface{}) (*services.RepositoryQuery, error) {
	// Extract entity type
	entityTypeStr, ok := aiRequest["entity_type"].(string)
	if !ok {
		return nil, fmt.Errorf("entity_type is required")
	}
	entityType := models.ToEntityType(entityTypeStr)

	// Extract operation
	operationStr, ok := aiRequest["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation is required")
	}
	operation := services.OperationType(operationStr)

	// Extract parameters
	params := services.QueryParameters{}
	if parameters, exists := aiRequest["parameters"].(map[string]interface{}); exists {
		params = era.mapToQueryParameters(parameters)
	}

	return &services.RepositoryQuery{
		EntityType: entityType,
		Operation:  operation,
		Parameters: params,
	}, nil
}

// queryParametersToMap converts QueryParameters to map for validation
func (era *EnhancedRepositoryAdapter) queryParametersToMap(params services.QueryParameters) map[string]interface{} {
	result := make(map[string]interface{})

	// Use reflection to convert struct to map
	v := reflect.ValueOf(params)
	t := reflect.TypeOf(params)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get JSON tag name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		// Remove omitempty and other flags
		jsonTag = strings.Split(jsonTag, ",")[0]

		// Only include non-zero values
		if !field.IsZero() {
			result[jsonTag] = field.Interface()
		}
	}

	return result
}

// mapToQueryParameters converts map to QueryParameters struct
func (era *EnhancedRepositoryAdapter) mapToQueryParameters(params map[string]interface{}) services.QueryParameters {
	queryParams := services.QueryParameters{}

	// Use reflection to populate struct from map
	v := reflect.ValueOf(&queryParams).Elem()
	t := reflect.TypeOf(queryParams)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get JSON tag name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		// Remove omitempty and other flags
		jsonTag = strings.Split(jsonTag, ",")[0]

		// Set field value if it exists in the map
		if value, exists := params[jsonTag]; exists && field.CanSet() {
			era.setFieldValue(field, value)
		}
	}

	return queryParams
}

// setFieldValue sets a struct field value from interface{}
func (era *EnhancedRepositoryAdapter) setFieldValue(field reflect.Value, value interface{}) {
	if value == nil {
		return
	}

	switch field.Kind() {
	case reflect.String:
		if str, ok := value.(string); ok {
			field.SetString(str)
		}
	case reflect.Int64:
		if num, ok := value.(float64); ok {
			field.SetInt(int64(num))
		} else if num, ok := value.(int64); ok {
			field.SetInt(num)
		}
	case reflect.Float64:
		if num, ok := value.(float64); ok {
			field.SetFloat(num)
		}
	case reflect.Bool:
		if b, ok := value.(bool); ok {
			field.SetBool(b)
		}
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			if arr, ok := value.([]interface{}); ok {
				strSlice := make([]string, len(arr))
				for i, v := range arr {
					if str, ok := v.(string); ok {
						strSlice[i] = str
					}
				}
				field.Set(reflect.ValueOf(strSlice))
			}
		}
	}
}

// getResultCount determines the count of results
func (era *EnhancedRepositoryAdapter) getResultCount(result interface{}) int {
	if result == nil {
		return 0
	}

	v := reflect.ValueOf(result)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return v.Len()
	case reflect.Ptr:
		if v.IsNil() {
			return 0
		}
		return 1
	default:
		return 1
	}
}
