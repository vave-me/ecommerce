package tools

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
	
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
)

// ParameterError represents an error in parameter extraction/conversion
type ParameterError struct {
	Key          string
	ExpectedType string
	ActualType   string
	Value        interface{}
	Message      string
}

func (e *ParameterError) Error() string {
	return fmt.Sprintf("parameter '%s': expected %s, got %s with value %v - %s",
		e.Key, e.ExpectedType, e.ActualType, e.Value, e.Message)
}

// getStringParam safely extracts a string parameter from the parameters map
func getStringParam(params map[string]interface{}, key string, defaultValues ...string) string {
	defaultValue := ""
	if len(defaultValues) > 0 {
		defaultValue = defaultValues[0]
	}
	
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case string:
			return v
		case fmt.Stringer:
			return v.String()
		default:
			// Log type mismatch for debugging
			log.Printf("[WARN] Parameter '%s' type mismatch: expected string, got %T, using default", key, val)
		}
	}
	return defaultValue
}

// getIntParam safely extracts an int parameter from the parameters map
func getIntParam(params map[string]interface{}, key string, defaultVal int) int {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return defaultVal
}

// getInt64Param safely extracts an int64 parameter from the parameters map
func getInt64Param(params map[string]interface{}, key string, defaultValue int64) int64 {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case int32:
			return int64(v)
		case float64:
			// Check for fractional part
			if v != math.Trunc(v) {
				log.Printf("[WARN] Parameter '%s': float64 value %f has fractional part, truncating to %d", key, v, int64(v))
			}
			// Check for overflow
			if v > float64(math.MaxInt64) || v < float64(math.MinInt64) {
				log.Printf("[ERROR] Parameter '%s': float64 value %f out of int64 range, using default %d", key, v, defaultValue)
				return defaultValue
			}
			return int64(v)
		case float32:
			if float64(v) != math.Trunc(float64(v)) {
				log.Printf("[WARN] Parameter '%s': float32 value %f has fractional part, truncating", key, v)
			}
			return int64(v)
		case string:
			parsed, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err != nil {
				log.Printf("[ERROR] Parameter '%s': failed to parse string '%s' as int64: %v, using default %d", key, v, err, defaultValue)
				return defaultValue
			}
			return parsed
		default:
			log.Printf("[ERROR] Parameter '%s': unsupported type %T for int64 conversion, using default %d", key, val, defaultValue)
		}
	}
	return defaultValue
}

// getFloat32Param safely extracts a float32 parameter from the parameters map
func getFloat32Param(params map[string]interface{}, key string, defaultValue float32) float32 {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case float32:
			return v
		case float64:
			// Check for range
			if v > float64(math.MaxFloat32) || (v != 0 && math.Abs(v) < float64(math.SmallestNonzeroFloat32)) {
				log.Printf("[ERROR] Parameter '%s': float64 value %f out of float32 range, using default %f", key, v, defaultValue)
				return defaultValue
			}
			return float32(v)
		case int:
			return float32(v)
		case int64:
			return float32(v)
		case string:
			parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 32)
			if err != nil {
				log.Printf("[ERROR] Parameter '%s': failed to parse string '%s' as float32: %v, using default %f", key, v, err, defaultValue)
				return defaultValue
			}
			return float32(parsed)
		default:
			log.Printf("[ERROR] Parameter '%s': unsupported type %T for float32 conversion, using default %f", key, val, defaultValue)
		}
	}
	return defaultValue
}

// getFloat64Param safely extracts a float64 parameter from the parameters map
func getFloat64Param(params map[string]interface{}, key string, defaultValue float64) float64 {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case int:
			return float64(v)
		case int32:
			return float64(v)
		case int64:
			return float64(v)
		case string:
			parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
			if err != nil {
				log.Printf("[ERROR] Parameter '%s': failed to parse string '%s' as float64: %v, using default %f", key, v, err, defaultValue)
				return defaultValue
			}
			return parsed
		default:
			log.Printf("[ERROR] Parameter '%s': unsupported type %T for float64 conversion, using default %f", key, val, defaultValue)
		}
	}
	return defaultValue
}

// getBoolParam safely extracts a boolean parameter from the parameters map
func getBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case bool:
			return v
		case string:
			if parsed, err := strconv.ParseBool(v); err == nil {
				return parsed
			}
		case int:
			return v != 0
		case int64:
			return v != 0
		case float64:
			return v != 0
		}
	}
	return defaultValue
}

// getArrayParam safely extracts an array parameter from the parameters map
func getArrayParam(params map[string]interface{}, key string) []interface{} {
	if val, ok := params[key]; ok {
		if arr, ok := val.([]interface{}); ok {
			return arr
		}
	}
	return []interface{}{}
}

// getStringArrayParam safely extracts a string array parameter from the parameters map
func getStringArrayParam(params map[string]interface{}, key string) []string {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case []string:
			return v
		case []interface{}:
			result := make([]string, 0, len(v))
			for _, item := range v {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return []string{}
}

// getStringSliceParam safely extracts a string slice parameter from the parameters map
func getStringSliceParam(params map[string]interface{}, key string) []string {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case []string:
			return v
		case []interface{}:
			result := make([]string, len(v))
			for i, item := range v {
				if str, ok := item.(string); ok {
					result[i] = str
				}
			}
			return result
		case string:
			// Single string converted to slice
			return []string{v}
		}
	}
	return []string{}
}

// getTimeParam safely extracts a time parameter from the parameters map
func getTimeParam(params map[string]interface{}, key string, defaultValue time.Time) time.Time {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case time.Time:
			return v
		case string:
			if parsed, err := time.Parse(time.RFC3339, v); err == nil {
				return parsed
			}
		case int64:
			return time.Unix(v, 0)
		}
	}
	return defaultValue
}

// getMapParam safely extracts a map[string]interface{} parameter from the parameters map
func getMapParam(params map[string]interface{}, key string) map[string]interface{} {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case map[string]interface{}:
			return v
		default:
			log.Printf("[WARN] Parameter '%s' type mismatch: expected map[string]interface{}, got %T", key, val)
		}
	}
	return make(map[string]interface{})
}

// getStringMapParam safely extracts a map[string]string parameter from the parameters map
func getStringMapParam(params map[string]interface{}, key string) map[string]string {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case map[string]string:
			return v
		case map[string]interface{}:
			result := make(map[string]string)
			for k, value := range v {
				if str, ok := value.(string); ok {
					result[k] = str
				}
			}
			return result
		default:
			log.Printf("[WARN] Parameter '%s' type mismatch: expected map[string]string, got %T", key, val)
		}
	}
	return make(map[string]string)
}

// ContainsIgnoreCase checks if substr is present in s, case-insensitive.
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// ToolConfig represents tool configuration
type ToolConfig struct {
	MaxParallelExecutions int           `json:"max_parallel_executions"`
	ExecutionTimeout      time.Duration `json:"execution_timeout"`
	RetryAttempts         int           `json:"retry_attempts"`
	RetryDelay            time.Duration `json:"retry_delay"`
}

// DefaultToolConfig returns default tool configuration
func DefaultToolConfig() *ToolConfig {
	return &ToolConfig{
		MaxParallelExecutions: 10,
		ExecutionTimeout:      30 * time.Second,
		RetryAttempts:         3,
		RetryDelay:            1 * time.Second,
	}
}

// StreamingConfig represents streaming configuration
type StreamingConfig struct {
	StreamBufferSize  int           `json:"stream_buffer_size"`
	StreamTimeout     time.Duration `json:"stream_timeout"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}

// DefaultStreamingConfig returns default streaming configuration
func DefaultStreamingConfig() *StreamingConfig {
	return &StreamingConfig{
		StreamBufferSize:  100,
		StreamTimeout:     5 * time.Minute,
		HeartbeatInterval: 30 * time.Second,
	}
}

// getFloat32ArrayParam safely extracts a []float32 parameter from the parameters map
func getFloat32ArrayParam(params map[string]interface{}, key string) []float32 {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case []float32:
			return v
		case []interface{}:
			result := make([]float32, 0, len(v))
			for _, item := range v {
				switch num := item.(type) {
				case float32:
					result = append(result, num)
				case float64:
					result = append(result, float32(num))
				case int:
					result = append(result, float32(num))
				case int64:
					result = append(result, float32(num))
				default:
					log.Printf("[WARN] Item in array '%s' type mismatch: expected numeric type, got %T", key, item)
				}
			}
			return result
		default:
			log.Printf("[WARN] Parameter '%s' type mismatch: expected []float32, got %T", key, val)
		}
	}
	return []float32{}
}

// getAttributesParam safely extracts a []models.Attribute parameter from the parameters map
func getAttributesParam(params map[string]interface{}, key string) []models.Attribute {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case []models.Attribute:
			return v
		case []interface{}:
			var attributes []models.Attribute
			for _, item := range v {
				if attrMap, ok := item.(map[string]interface{}); ok {
					attr := models.Attribute{
						Key:   getStringParam(attrMap, "key"),
						Value: getStringParam(attrMap, "value"),
					}
					attributes = append(attributes, attr)
				}
			}
			return attributes
		default:
			log.Printf("[WARN] Parameter '%s' type mismatch: expected []models.Attribute, got %T", key, val)
		}
	}
	return []models.Attribute{}
}

// getOptionsParam safely extracts a []models.Option parameter from the parameters map
func getOptionsParam(params map[string]interface{}, key string) []models.Option {
	if val, exists := params[key]; exists {
		switch v := val.(type) {
		case []models.Option:
			return v
		case []interface{}:
			var options []models.Option
			for _, item := range v {
				if optMap, ok := item.(map[string]interface{}); ok {
					opt := models.Option{
						Name:  getStringParam(optMap, "name"),
						Value: getStringParam(optMap, "value"),
						Price: getFloat64Param(optMap, "price", 0),
					}
					options = append(options, opt)
				}
			}
			return options
		default:
			log.Printf("[WARN] Parameter '%s' type mismatch: expected []models.Option, got %T", key, val)
		}
	}
	return []models.Option{}
}

// createServiceFromParams creates a Service model from parameters
func createServiceFromParams(params map[string]interface{}) *models.Service {
	return &models.Service{
		ID:               getStringParam(params, "service_id"),
		Name:             getStringParam(params, "name"),
		Description:      getStringParam(params, "description"),
		ServiceType:      getStringParam(params, "service_type"),
		BasePrice:        getInt64Param(params, "base_price", 0),
		Pricing:          getStringArrayParam(params, "pricing"),
		Availability:     getStringParam(params, "availability"),
		ProviderName:     getStringParam(params, "provider_name"),
		UserID:           getStringParam(params, "user_id"),
		CategoryID:       getStringParam(params, "category_id"),
		CategorySlug:     getStringParam(params, "category_slug"),
		DescriptionShort: getStringParam(params, "description_short"),
		DescriptionLong:  getStringParam(params, "description_long"),
		Qualifications:   getStringArrayParam(params, "qualifications"),
		Contact:          getStringParam(params, "contact"),
		Faq:              getStringParam(params, "faq"),
		Tags:             getStringArrayParam(params, "tags"),
		Status:           getStringParam(params, "status"),
		UserType:         getStringParam(params, "user_type"),
		ShippingCost:     getInt64Param(params, "shipping_cost", 0),
		HasVariants:      getBoolParam(params, "has_variants", false),
		MiddlemanService: getBoolParam(params, "middleman_service", false),
		Negotiable:       getBoolParam(params, "negotiable", false),
		Attributes:       getAttributesParam(params, "attributes"),
		Options:          getOptionsParam(params, "options"),
		Thumbnail:        getStringParam(params, "thumbnail"),
		Lat:              getFloat64Param(params, "lat", 0),
		Lng:              getFloat64Param(params, "lng", 0),
		EntityType:       models.EntityType(getStringParam(params, "entity_type")),
	}
}

// convertToOrderItems converts params to order items
func convertToOrderItems(params map[string]interface{}) []models.Item {
	if items, ok := params["items"].([]interface{}); ok {
		var orderItems []models.Item
		for _, item := range items {
			if itemMap, ok := item.(map[string]interface{}); ok {
				orderItems = append(orderItems, models.Item{
					UserSellerID:   getStringParam(itemMap, "user_seller_id"),
					ProductID:      getStringParam(itemMap, "product_id"),
					UserSellerName: getStringParam(itemMap, "user_seller_name"),
					ProductName:    getStringParam(itemMap, "product_name"),
					Price:          getInt64Param(itemMap, "price", 0),
					Quantity:       getInt64Param(itemMap, "quantity", 1),
				})
			}
		}
		return orderItems
	}
	return []models.Item{}
}

// createVectorFilters creates VectorFilters from parameters
func createVectorFilters(params map[string]interface{}) *domain.VectorFilters {
	filters := &domain.VectorFilters{
		EntityTypes:     getStringArrayParam(params, "filter_entity_types"),
		Statuses:        getStringArrayParam(params, "filter_statuses"),
		Categories:      getStringArrayParam(params, "filter_categories"),
		NegotiableOnly:  getBoolParam(params, "filter_negotiable_only", false),
		UserType:        getStringParam(params, "filter_user_type"),
		MetadataFilters: getStringMapParam(params, "filter_metadata"),
	}

	// Price range
	if minPrice := getFloat64Param(params, "filter_min_price", 0); minPrice > 0 {
		if filters.PriceRange == nil {
			filters.PriceRange = &domain.PriceRange{}
		}
		filters.PriceRange.MinPrice = int64(minPrice)
	}
	if maxPrice := getFloat64Param(params, "filter_max_price", 0); maxPrice > 0 {
		if filters.PriceRange == nil {
			filters.PriceRange = &domain.PriceRange{}
		}
		filters.PriceRange.MaxPrice = int64(maxPrice)
	}

	// Geo filter
	if lat := getFloat64Param(params, "filter_lat", 0); lat != 0 {
		if filters.GeoFilter == nil {
			filters.GeoFilter = &domain.GeoFilter{}
		}
		filters.GeoFilter.Latitude = lat
		filters.GeoFilter.Longitude = getFloat64Param(params, "filter_lng", 0)
		filters.GeoFilter.RadiusKm = getFloat64Param(params, "filter_radius_km", 10)
		filters.GeoFilter.Country = getStringParam(params, "filter_country")
		filters.GeoFilter.City = getStringParam(params, "filter_city")
	}

	// Time range
	if startTime := getTimeParam(params, "filter_start_time", time.Time{}); !startTime.IsZero() {
		if filters.TimeRange == nil {
			filters.TimeRange = &domain.TimeRange{}
		}
		filters.TimeRange.StartTime = startTime
		filters.TimeRange.EndTime = getTimeParam(params, "filter_end_time", time.Time{})
	}

	return filters
}