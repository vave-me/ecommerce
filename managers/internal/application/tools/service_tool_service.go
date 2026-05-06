package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
)

// ServiceToolService handles all service-related operations
type ServiceToolService struct {
	serviceRepo domain.ServiceRepository
}

// NewServiceToolService creates a new service tool service
func NewServiceToolService(serviceRepo domain.ServiceRepository) *ServiceToolService {
	return &ServiceToolService{
		serviceRepo: serviceRepo,
	}
}

// ExecuteOperation handles service-related operations with streaming support
func (s *ServiceToolService) ExecuteOperation(ctx context.Context, operation string, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	log.Printf("ServiceToolService.ExecuteOperation: Executing service operation: %s", operation)

	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "service_operation",
		Status:   "progress",
		Progress: 10,
		Metadata: map[string]interface{}{
			"operation": operation,
			"message":   fmt.Sprintf("Processing service operation: %s", operation),
		},
		Timestamp: time.Now().Unix(),
	}

	// Extract common parameters
	id := getStringParam(parameters, "id", "")
	userID := getStringParam(parameters, "user_id", "")
	name := getStringParam(parameters, "name", "")
	searchTerm := getStringParam(parameters, "search_term", "")
	categoryID := getStringParam(parameters, "category_id", "")
	categorySlug := getStringParam(parameters, "category_slug", "")
	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)
	sortBy := getStringParam(parameters, "sort_by", "")
	sortOrder := getStringParam(parameters, "sort_order", "")

	var result interface{}
	var err error

	switch operation {
	case "find", "get":
		if id == "" {
			return nil, fmt.Errorf("service ID is required for find operation")
		}
		result, err = s.serviceRepo.Find(ctx, id)
	case "search":
		if searchTerm == "" {
			return nil, fmt.Errorf("search term is required for search operation")
		}
		result, err = s.serviceRepo.SearchWithTerm(ctx, searchTerm)
	case "suggest":
		if name == "" {
			return nil, fmt.Errorf("name is required for suggest operation")
		}
		result, err = s.serviceRepo.SuggestServices(ctx, name)
	case "get_services", "list":
		services, count, serviceErr := s.serviceRepo.GetServices(ctx, page, pageSize, sortBy, sortOrder)
		if serviceErr != nil {
			err = serviceErr
		} else {
			result = map[string]interface{}{"services": services, "total_count": count}
		}
	case "get_catalog", "catalog":
		if userID == "" {
			return nil, fmt.Errorf("user ID is required for catalog operation")
		}
		services, count, serviceErr := s.serviceRepo.GetCatalog(ctx, page, pageSize, sortBy, sortOrder)
		if serviceErr != nil {
			err = serviceErr
		} else {
			result = map[string]interface{}{"services": services, "total_count": count}
		}
	case "get_public_catalog", "public_catalog":
		if userID == "" {
			return nil, fmt.Errorf("user ID is required for public catalog operation")
		}
		services, count, serviceErr := s.serviceRepo.GetPublicCatalog(ctx, userID, page, pageSize, sortBy, sortOrder)
		if serviceErr != nil {
			err = serviceErr
		} else {
			result = map[string]interface{}{"services": services, "total_count": count}
		}
	case "search_by_category":
		if categoryID == "" {
			return nil, fmt.Errorf("category ID is required for category search")
		}
		result, err = s.serviceRepo.SearchServicesWithCategory(ctx, categoryID, page, pageSize, sortBy, sortOrder)
	case "search_by_category_slug":
		if categorySlug == "" {
			return nil, fmt.Errorf("category slug is required for category slug search")
		}
		result, err = s.serviceRepo.SearchServicesWithCategorySlug(ctx, categorySlug, page, pageSize, sortBy, sortOrder)
	case "add", "create":
		result, err = s.handleAddService(ctx, parameters, streamChan, toolID)
	case "update":
		result, err = s.handleUpdateService(ctx, parameters, streamChan, toolID)
	case "remove", "delete":
		result, err = s.handleRemoveService(ctx, parameters, streamChan, toolID)
	case "archive":
		result, err = s.handleArchiveService(ctx, parameters, streamChan, toolID)
	case "mark_sold":
		result, err = s.handleMarkServiceSold(ctx, parameters, streamChan, toolID)
	case "mark_leased":
		result, err = s.handleMarkServiceLeased(ctx, parameters, streamChan, toolID)
	case "filter":
		result, err = s.handleFilterServices(ctx, parameters, streamChan, toolID)
	default:
		streamChan <- ToolExecutionStream{
			ID:        toolID,
			ToolName:  "service_operation",
			Status:    "error",
			Progress:  100,
			Error:     fmt.Sprintf("Unsupported service operation: %s", operation),
			Timestamp: time.Now().Unix(),
		}
		return nil, fmt.Errorf("unsupported service operation: %s", operation)
	}

	// Handle result or error
	if err != nil {
		streamChan <- ToolExecutionStream{
			ID:        toolID,
			ToolName:  "service_operation",
			Status:    "error",
			Progress:  100,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
		}
		return nil, err
	}

	// Send success
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "service_operation",
		Status:   "completed",
		Progress: 100,
		Result:   result,
		Metadata: map[string]interface{}{
			"operation": operation,
			"message":   fmt.Sprintf("Service operation %s completed successfully", operation),
		},
		Timestamp: time.Now().Unix(),
	}

	return result, nil
}

func (s *ServiceToolService) handleAddService(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	// Extract parameters with intelligent defaults
	name := getStringParam(parameters, "name", "")
	description := getStringParam(parameters, "description", "")
	// If description is empty, try description_short or description_long
	if description == "" {
		description = getStringParam(parameters, "description_short", "")
		if description == "" {
			description = getStringParam(parameters, "description_long", "")
		}
	}
	basePrice := getInt64Param(parameters, "base_price", 0)
	userSellerID := getStringParam(parameters, "user_seller_id", "")
	if userSellerID == "" {
		userSellerID = getStringParam(parameters, "user_id", "")
	}
	categoryID := getStringParam(parameters, "category_id", "")

	// Generate name from description or service type if not provided
	if name == "" {
		serviceType := getStringParam(parameters, "service_type", "")
		if serviceType != "" {
			name = serviceType
		} else if description != "" {
			// Take first 50 chars of description as name
			if len(description) > 50 {
				name = description[:50] + "..."
			} else {
				name = description
			}
		} else {
			// Check if there's any other identifying info
			providerName := getStringParam(parameters, "provider_name", "")
			if providerName != "" {
				name = providerName + " Service"
			} else {
				return nil, fmt.Errorf("NAME_REQUIRED: unable to determine service name. Please provide at least one of: name, service_type, description, or provider_name")
			}
		}
		log.Printf("ServiceToolService: Generated name '%s' from available data", name)
	}

	// If category_id is missing, return an error suggesting to search for category first
	if categoryID == "" {
		// Check if we have category-related info to suggest a search
		categorySlug := getStringParam(parameters, "category_slug", "")
		categoryHint := ""
		if categorySlug != "" {
			categoryHint = fmt.Sprintf(" (category_slug: %s)", categorySlug)
		} else if serviceType := getStringParam(parameters, "service_type", ""); serviceType != "" {
			categoryHint = fmt.Sprintf(" (service_type: %s)", serviceType)
		}

		return nil, fmt.Errorf("CATEGORY_REQUIRED: category_id is missing. Please use category_operation tool to search for appropriate category%s, then retry service creation with the found category_id", categoryHint)
	}

	// Extract optional parameters
	categorySlug := getStringParam(parameters, "category_slug", "")
	serviceType := getStringParam(parameters, "service_type", "")
	_ = getInt64Param(parameters, "duration", 0) // duration not used in current interface
	availability := getStringParam(parameters, "availability", "")
	_ = getStringParam(parameters, "location", "") // location not used in current interface
	tags := getStringSliceParam(parameters, "tags")
	status := getStringParam(parameters, "status", "active")
	negotiable := getBoolParam(parameters, "negotiable", false)
	userType := getStringParam(parameters, "user_type", "")
	hasVariants := getBoolParam(parameters, "has_variants", false)
	thumbnail := getStringParam(parameters, "thumbnail", "")
	lat := getFloat32Param(parameters, "lat", 0.0)
	lng := getFloat32Param(parameters, "lng", 0.0)

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "service_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step": "adding_service",
			"name": name,
		},
		Timestamp: time.Now().Unix(),
	}

	// Convert parameters to match domain interface
	pricing := getStringSliceParam(parameters, "pricing")
	descriptionShort := getStringParam(parameters, "description_short", description)
	descriptionLong := getStringParam(parameters, "description_long", description)
	qualifications := getStringSliceParam(parameters, "qualifications")
	contact := getStringParam(parameters, "contact", "")
	faq := getStringParam(parameters, "faq", "")
	attributes := getStringSliceParam(parameters, "attributes")
	options := getStringSliceParam(parameters, "options")

	// Convert string status to models.Status
	var statusEnum models.Status
	switch status {
	case "active":
		statusEnum = models.StatusActive
	case "inactive":
		statusEnum = models.StatusPaused
	default:
		statusEnum = models.StatusActive
	}

	// Convert string userType to models.UserType
	var userTypeEnum models.UserType
	switch userType {
	case "business":
		userTypeEnum = models.UserTypeBusiness
	case "individual", "private":
		userTypeEnum = models.UserTypePrivate
	default:
		userTypeEnum = models.UserTypePrivate
	}

	shippingCost := getInt64Param(parameters, "shipping_cost", 0)
	middlemanService := getBoolParam(parameters, "middleman_service", false)

	// Extract provider name from parameters
	providerName := getStringParam(parameters, "provider_name", "")

	err := s.serviceRepo.Add(
		ctx, "", name, description, serviceType,
		basePrice, pricing, availability,
		providerName, categoryID, categorySlug,
		descriptionShort, descriptionLong,
		qualifications,
		contact, faq,
		tags,
		statusEnum,
		userTypeEnum,
		shippingCost,
		negotiable, hasVariants, middlemanService,
		attributes,
		options,
		thumbnail,
		float64(lat), float64(lng),
	)
	if err != nil {
		return nil, fmt.Errorf("add service failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "services",
		"operation":   "add",
		"result":      "Service added successfully",
		"name":        name,
	}, nil
}

func (s *ServiceToolService) handleUpdateService(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	serviceID := getStringParam(parameters, "id", "")
	if serviceID == "" {
		serviceID = getStringParam(parameters, "service_id", "")
	}
	userID := getStringParam(parameters, "user_id", "")

	if serviceID == "" || userID == "" {
		return nil, fmt.Errorf("service_id and user_id parameters required")
	}

	// Extract update parameters
	name := getStringParam(parameters, "name", "")
	description := getStringParam(parameters, "description", "")
	basePrice := getInt64Param(parameters, "base_price", 0)
	categoryID := getStringParam(parameters, "category_id", "")
	serviceType := getStringParam(parameters, "service_type", "")
	_ = getInt64Param(parameters, "duration", 0)       // duration not used in current interface
	_ = getStringParam(parameters, "availability", "") // availability not used in current interface
	_ = getStringParam(parameters, "location", "")     // location not used in current interface
	_ = getStringSliceParam(parameters, "tags")        // tags not used in current interface
	_ = getStringParam(parameters, "status", "")       // status not used in current interface
	_ = getBoolParam(parameters, "negotiable", false)  // negotiable not used in current interface

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "service_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":       "updating_service",
			"service_id": serviceID,
		},
		Timestamp: time.Now().Unix(),
	}

	// Build service model for update
	service := &models.Service{
		ID:          serviceID,
		Name:        name,
		Description: description,
		BasePrice:   basePrice,
		CategoryID:  categoryID,
		ServiceType: serviceType,
		// Add other fields as needed
	}
	_, err := s.serviceRepo.Update(ctx, serviceID, service)
	if err != nil {
		return nil, fmt.Errorf("update service failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "services",
		"operation":   "update",
		"result":      "Service updated successfully",
		"service_id":  serviceID,
	}, nil
}

func (s *ServiceToolService) handleRemoveService(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	serviceID := getStringParam(parameters, "id", "")
	if serviceID == "" {
		serviceID = getStringParam(parameters, "service_id", "")
	}
	userID := getStringParam(parameters, "user_id", "")

	if serviceID == "" || userID == "" {
		return nil, fmt.Errorf("service_id and user_id parameters required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "service_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":       "removing_service",
			"service_id": serviceID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := s.serviceRepo.Remove(ctx, serviceID, userID)
	if err != nil {
		return nil, fmt.Errorf("remove service failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "services",
		"operation":   "remove",
		"result":      "Service removed successfully",
		"service_id":  serviceID,
	}, nil
}

func (s *ServiceToolService) handleArchiveService(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	serviceID := getStringParam(parameters, "id", "")
	if serviceID == "" {
		serviceID = getStringParam(parameters, "service_id", "")
	}

	if userID == "" || serviceID == "" {
		return nil, fmt.Errorf("user_id and service_id parameters required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "service_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":       "archiving_service",
			"service_id": serviceID,
		},
		Timestamp: time.Now().Unix(),
	}

	result, err := s.serviceRepo.ArchiveService(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("archive service failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "services",
		"operation":   "archive",
		"result":      result,
		"service_id":  serviceID,
	}, nil
}

func (s *ServiceToolService) handleMarkServiceSold(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	serviceID := getStringParam(parameters, "id", "")
	if serviceID == "" {
		serviceID = getStringParam(parameters, "service_id", "")
	}

	if userID == "" || serviceID == "" {
		return nil, fmt.Errorf("user_id and service_id parameters required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "service_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":       "marking_service_sold",
			"service_id": serviceID,
		},
		Timestamp: time.Now().Unix(),
	}

	result, err := s.serviceRepo.MarkServiceSold(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("mark service sold failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "services",
		"operation":   "mark_sold",
		"result":      result,
		"service_id":  serviceID,
	}, nil
}

func (s *ServiceToolService) handleMarkServiceLeased(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	serviceID := getStringParam(parameters, "id", "")
	if serviceID == "" {
		serviceID = getStringParam(parameters, "service_id", "")
	}
	monthlyPrice := getInt64Param(parameters, "monthly_price", 0)
	leaseTermMonths := getInt64Param(parameters, "lease_term_months", 0)

	if userID == "" || serviceID == "" || monthlyPrice == 0 || leaseTermMonths == 0 {
		return nil, fmt.Errorf("user_id, service_id, monthly_price, and lease_term_months parameters required")
	}

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "service_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":              "marking_service_leased",
			"service_id":        serviceID,
			"monthly_price":     monthlyPrice,
			"lease_term_months": leaseTermMonths,
		},
		Timestamp: time.Now().Unix(),
	}

	result, err := s.serviceRepo.MarkServiceLeased(ctx, serviceID, monthlyPrice, leaseTermMonths)
	if err != nil {
		return nil, fmt.Errorf("mark service leased failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type":       "services",
		"operation":         "mark_leased",
		"result":            result,
		"service_id":        serviceID,
		"monthly_price":     monthlyPrice,
		"lease_term_months": leaseTermMonths,
	}, nil
}

func (s *ServiceToolService) handleFilterServices(ctx context.Context, parameters map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	// Extract filter parameters
	name := getStringParam(parameters, "name", "")
	description := getStringParam(parameters, "description", "")
	categoryID := getStringParam(parameters, "category_id", "")
	serviceType := getStringParam(parameters, "service_type", "")
	minPrice := getInt64Param(parameters, "min_price", 0)
	maxPrice := getInt64Param(parameters, "max_price", 0)
	minDuration := getInt64Param(parameters, "min_duration", 0)
	maxDuration := getInt64Param(parameters, "max_duration", 0)
	availability := getStringParam(parameters, "availability", "")
	location := getStringParam(parameters, "location", "")
	tags := getStringSliceParam(parameters, "tags")
	status := getStringParam(parameters, "status", "")
	negotiable := getBoolParam(parameters, "negotiable", false)
	page := getInt64Param(parameters, "page", 1)
	pageSize := getInt64Param(parameters, "page_size", 20)
	sortBy := getStringParam(parameters, "sort_by", "created_at")
	sortOrder := getStringParam(parameters, "sort_order", "desc")

	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "service_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":         "filtering_services",
			"service_type": serviceType,
			"location":     location,
		},
		Timestamp: time.Now().Unix(),
	}

	// Use existing SearchServicesWithFilter method with proper parameters
	services, err := s.serviceRepo.SearchServicesWithFilter(
		ctx,
		categoryID, "", serviceType,
		"",                  // userID not available in this context
		models.StatusActive, // use default status
		name,                // searchText
		minPrice, maxPrice,
		time.Time{}, time.Time{}, // availableFrom, availableTo not available
		false, negotiable, false, // hasVariants, negotiable, middlemanService
		models.UserTypePrivate, // default userType
		tags, []string{},       // tags, qualifications
		(page-1)*pageSize, // offset
		pageSize,          // limit
		0.0, 0.0,          // lat, lng not available
		0, // radius not available
		page, pageSize,
		sortBy, sortOrder,
	)
	totalCount := int64(len(services)) // Approximate count since exact count not available from this method
	if err != nil {
		return nil, fmt.Errorf("filter services failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "services",
		"operation":   "filter",
		"results":     services,
		"total_count": totalCount,
		"filters": map[string]interface{}{
			"name":         name,
			"description":  description,
			"category_id":  categoryID,
			"service_type": serviceType,
			"min_price":    minPrice,
			"max_price":    maxPrice,
			"min_duration": minDuration,
			"max_duration": maxDuration,
			"availability": availability,
			"location":     location,
			"tags":         tags,
			"status":       status,
			"negotiable":   negotiable,
		},
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	}, nil
}
