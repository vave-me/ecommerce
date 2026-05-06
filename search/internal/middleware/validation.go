package middleware

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	
	"middleman/search/searchpb"
)

// ValidationConfig defines validation rules
type ValidationConfig struct {
	MaxPageSize         int64
	MaxSearchRadius     int64
	MaxStringLength     int
	MaxArrayElements    int
	MaxEntityTypes      int
	MaxPriceRange       int64
	MaxBatchSize        int
}

// RequestValidator validates incoming requests
type RequestValidator struct {
	config ValidationConfig
}

// NewRequestValidator creates a new request validator
func NewRequestValidator(config ValidationConfig) *RequestValidator {
	return &RequestValidator{config: config}
}

// UnaryServerInterceptor creates a gRPC interceptor for request validation
func (rv *RequestValidator) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Validate based on method name
		if err := rv.validateRequest(req, info.FullMethod); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
		}
		
		return handler(ctx, req)
	}
}

// validateRequest validates a request based on its type
func (rv *RequestValidator) validateRequest(req interface{}, method string) error {
	switch r := req.(type) {
	case *searchpb.UnifiedSearchRequest:
		return rv.validateUnifiedSearch(r)
	case *searchpb.UnifiedFeedRequest:
		return rv.validateUnifiedFeed(r)
	case *searchpb.BatchUnifiedFeedRequest:
		return rv.validateBatchUnifiedFeed(r)
	case *searchpb.SearchProductsWithFiltersRequest:
		return rv.validateProductSearch(r)
	case *searchpb.SearchPostsWithFiltersRequest:
		return rv.validatePostSearch(r)
	case *searchpb.GetProductRequest:
		return rv.validateGetProduct(r)
	case *searchpb.GetPostRequest:
		return rv.validateGetPost(r)
	}
	
	// No validation for other request types
	return nil
}

// validateUnifiedSearch validates unified search requests
func (rv *RequestValidator) validateUnifiedSearch(req *searchpb.UnifiedSearchRequest) error {
	// Validate search term
	if len(req.GetSearchTerm()) > rv.config.MaxStringLength {
		return fmt.Errorf("search term too long (max %d characters)", rv.config.MaxStringLength)
	}
	
	if req.GetSearchTerm() == "" {
		return fmt.Errorf("search term cannot be empty")
	}
	
	// Validate entity types
	if len(req.GetEntityTypes()) > rv.config.MaxEntityTypes {
		return fmt.Errorf("too many entity types (max %d)", rv.config.MaxEntityTypes)
	}
	
	// Validate pagination
	if err := rv.validatePagination(req.GetPage(), req.GetPageSize()); err != nil {
		return err
	}
	
	// Validate filters
	if req.GetFilters() != nil {
		if err := rv.validateCommonFilters(req.GetFilters()); err != nil {
			return err
		}
	}
	
	// Validate geo parameters
	if err := rv.validateGeoParams(req.GetLat(), req.GetLng(), req.GetRadius()); err != nil {
		return err
	}
	
	return nil
}

// validateUnifiedFeed validates unified feed requests
func (rv *RequestValidator) validateUnifiedFeed(req *searchpb.UnifiedFeedRequest) error {
	// Validate entity types
	if len(req.GetEntityTypes()) > rv.config.MaxEntityTypes {
		return fmt.Errorf("too many entity types (max %d)", rv.config.MaxEntityTypes)
	}
	
	// Validate pagination
	if err := rv.validatePagination(req.GetPage(), req.GetPageSize()); err != nil {
		return err
	}
	
	// Validate feed type
	validFeedTypes := map[string]bool{
		"latest":      true,
		"trending":    true,
		"popular":     true,
		"recommended": true,
	}
	
	if req.GetFeedType() != "" && !validFeedTypes[req.GetFeedType()] {
		return fmt.Errorf("invalid feed type: %s", req.GetFeedType())
	}
	
	// Validate filters
	if req.GetFilters() != nil {
		if err := rv.validateCommonFilters(req.GetFilters()); err != nil {
			return err
		}
	}
	
	return nil
}

// validateBatchUnifiedFeed validates batch feed requests
func (rv *RequestValidator) validateBatchUnifiedFeed(req *searchpb.BatchUnifiedFeedRequest) error {
	if len(req.GetRequests()) == 0 {
		return fmt.Errorf("batch request cannot be empty")
	}
	
	if len(req.GetRequests()) > rv.config.MaxBatchSize {
		return fmt.Errorf("batch size too large (max %d)", rv.config.MaxBatchSize)
	}
	
	// Validate each individual request
	for i, feedReq := range req.GetRequests() {
		if err := rv.validateUnifiedFeed(feedReq); err != nil {
			return fmt.Errorf("request %d: %v", i, err)
		}
	}
	
	return nil
}

// validateProductSearch validates product search requests
func (rv *RequestValidator) validateProductSearch(req *searchpb.SearchProductsWithFiltersRequest) error {
	// Validate string fields
	if err := rv.validateStringField("name", req.GetName()); err != nil {
		return err
	}
	
	// Validate price range
	if req.GetMinPrice() < 0 {
		return fmt.Errorf("min price cannot be negative")
	}
	if req.GetMaxPrice() > 0 && req.GetMaxPrice() < req.GetMinPrice() {
		return fmt.Errorf("max price must be greater than min price")
	}
	if req.GetMaxPrice() > rv.config.MaxPriceRange {
		return fmt.Errorf("max price too high (max %d)", rv.config.MaxPriceRange)
	}
	
	// Validate stock range
	if req.GetMinStock() < 0 {
		return fmt.Errorf("min stock cannot be negative")
	}
	if req.GetMaxStock() > 0 && req.GetMaxStock() < req.GetMinStock() {
		return fmt.Errorf("max stock must be greater than min stock")
	}
	
	// Validate tags
	if len(req.GetTags()) > rv.config.MaxArrayElements {
		return fmt.Errorf("too many tags (max %d)", rv.config.MaxArrayElements)
	}
	
	// Validate pagination
	if err := rv.validatePagination(req.GetPage(), req.GetPageSize()); err != nil {
		return err
	}
	
	// Validate geo parameters
	if err := rv.validateGeoParams(req.GetLat(), req.GetLng(), req.GetRadius()); err != nil {
		return err
	}
	
	return nil
}

// validatePostSearch validates post search requests
func (rv *RequestValidator) validatePostSearch(req *searchpb.SearchPostsWithFiltersRequest) error {
	// Validate string fields
	if err := rv.validateStringField("name", req.GetName()); err != nil {
		return err
	}
	if err := rv.validateStringField("description", req.GetDescription()); err != nil {
		return err
	}
	
	// Validate tags
	if len(req.GetTags()) > rv.config.MaxArrayElements {
		return fmt.Errorf("too many tags (max %d)", rv.config.MaxArrayElements)
	}
	
	// Validate pagination
	if err := rv.validatePagination(req.GetPage(), req.GetPageSize()); err != nil {
		return err
	}
	
	// Validate geo parameters
	if err := rv.validateGeoParams(req.GetLat(), req.GetLng(), req.GetRadius()); err != nil {
		return err
	}
	
	return nil
}

// validateGetProduct validates get product requests
func (rv *RequestValidator) validateGetProduct(req *searchpb.GetProductRequest) error {
	if req.GetId() == "" {
		return fmt.Errorf("product ID cannot be empty")
	}
	
	if len(req.GetId()) > rv.config.MaxStringLength {
		return fmt.Errorf("product ID too long (max %d characters)", rv.config.MaxStringLength)
	}
	
	// Basic sanitization check for ID
	if strings.ContainsAny(req.GetId(), "<>\"';&|") {
		return fmt.Errorf("product ID contains invalid characters")
	}
	
	return nil
}

// validateGetPost validates get post requests
func (rv *RequestValidator) validateGetPost(req *searchpb.GetPostRequest) error {
	if req.GetId() == "" {
		return fmt.Errorf("post ID cannot be empty")
	}
	
	if len(req.GetId()) > rv.config.MaxStringLength {
		return fmt.Errorf("post ID too long (max %d characters)", rv.config.MaxStringLength)
	}
	
	// Basic sanitization check for ID
	if strings.ContainsAny(req.GetId(), "<>\"';&|") {
		return fmt.Errorf("post ID contains invalid characters")
	}
	
	return nil
}

// validateCommonFilters validates common filter parameters
func (rv *RequestValidator) validateCommonFilters(filters *searchpb.CommonFilters) error {
	// Validate price range
	if filters.GetMinPrice() < 0 {
		return fmt.Errorf("min price cannot be negative")
	}
	if filters.GetMaxPrice() > 0 && filters.GetMaxPrice() < filters.GetMinPrice() {
		return fmt.Errorf("max price must be greater than min price")
	}
	
	// Validate entity types
	if len(filters.GetEntityTypes()) > rv.config.MaxEntityTypes {
		return fmt.Errorf("too many entity types in filter (max %d)", rv.config.MaxEntityTypes)
	}
	
	// Validate tags
	if len(filters.GetTags()) > rv.config.MaxArrayElements {
		return fmt.Errorf("too many tags in filter (max %d)", rv.config.MaxArrayElements)
	}
	
	// Validate user IDs
	if len(filters.GetUserIds()) > rv.config.MaxArrayElements {
		return fmt.Errorf("too many user IDs in filter (max %d)", rv.config.MaxArrayElements)
	}
	
	// Validate geo filter
	if filters.GetGeoFilter() != nil {
		geo := filters.GetGeoFilter()
		if err := rv.validateGeoParams(geo.GetLat(), geo.GetLng(), geo.GetRadiusKm()); err != nil {
			return fmt.Errorf("geo filter: %v", err)
		}
	}
	
	return nil
}

// validatePagination validates pagination parameters
func (rv *RequestValidator) validatePagination(page, pageSize int64) error {
	if page < 1 {
		return fmt.Errorf("page must be >= 1")
	}
	
	if pageSize < 1 {
		return fmt.Errorf("page size must be >= 1")
	}
	
	if pageSize > rv.config.MaxPageSize {
		return fmt.Errorf("page size too large (max %d)", rv.config.MaxPageSize)
	}
	
	// Prevent deep pagination attacks
	maxOffset := page * pageSize
	if maxOffset > 10000 {
		return fmt.Errorf("pagination too deep (max offset 10000)")
	}
	
	return nil
}

// validateGeoParams validates geographic search parameters
func (rv *RequestValidator) validateGeoParams(lat, lng, radius float64) error {
	// Only validate if geo search is being used
	if lat == 0 && lng == 0 && radius == 0 {
		return nil
	}
	
	// Validate latitude
	if lat < -90 || lat > 90 {
		return fmt.Errorf("invalid latitude (must be between -90 and 90)")
	}
	
	// Validate longitude
	if lng < -180 || lng > 180 {
		return fmt.Errorf("invalid longitude (must be between -180 and 180)")
	}
	
	// Validate radius
	if radius < 0 {
		return fmt.Errorf("radius cannot be negative")
	}
	if radius > float64(rv.config.MaxSearchRadius) {
		return fmt.Errorf("search radius too large (max %d km)", rv.config.MaxSearchRadius)
	}
	
	return nil
}

// validateStringField validates a string field
func (rv *RequestValidator) validateStringField(name, value string) error {
	if len(value) > rv.config.MaxStringLength {
		return fmt.Errorf("%s too long (max %d characters)", name, rv.config.MaxStringLength)
	}
	
	// Check for potential injection attempts
	if strings.ContainsAny(value, "\x00") {
		return fmt.Errorf("%s contains null bytes", name)
	}
	
	return nil
}

// GetDefaultValidationConfig returns default validation configuration
func GetDefaultValidationConfig() ValidationConfig {
	return ValidationConfig{
		MaxPageSize:      100,
		MaxSearchRadius:  100, // km
		MaxStringLength:  1000,
		MaxArrayElements: 50,
		MaxEntityTypes:   10,
		MaxPriceRange:    1000000000, // 1 billion
		MaxBatchSize:     10,
	}
}