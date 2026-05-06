package tools

import (
	"context"
	"fmt"
	"middleman/assistants/internal/models"
)

// ==================== GEOCODING HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeGeocodingHandlers() {
	// Core geocoding operations from protobuf
	r.handlers["geocoding_batch_geocode_address"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		addressStrings := getStringArrayParam(params, "addresses")
		if len(addressStrings) == 0 {
			return nil, fmt.Errorf("addresses array is required")
		}
		
		// Convert string array to Address array
		addresses := make([]models.Address, len(addressStrings))
		for i, addr := range addressStrings {
			addresses[i] = models.Address{
				ID:      fmt.Sprintf("addr_%d", i),
				Address: addr,
			}
		}
		
		req := &models.BatchGeocodeAddressRequest{
			Addresses: addresses,
		}
		return reg.geocodingRepo.BatchGeocodeAddress(ctx, req)
	}

	r.handlers["geocoding_geocode_address"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		address := getStringParam(params, "address")
		if id == "" || address == "" {
			return nil, fmt.Errorf("id and address are required")
		}
		return reg.geocodingRepo.GeocodeAddress(ctx, id, address)
	}

	r.handlers["geocoding_refresh_cache"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.geocodingRepo.RefreshGeocodingCache(ctx)
	}

	r.handlers["geocoding_reverse_geocode"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.geocodingRepo.ReverseGeocodeLocation(ctx)
	}

	r.handlers["geocoding_validate_address"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.geocodingRepo.ValidateAddress(ctx)
	}

	r.handlers["geocoding_get_address_for_coordinates"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.geocodingRepo.GetAddressForCoordinates(ctx)
	}

	r.handlers["geocoding_get_coordinates_for_address"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.geocodingRepo.GetCoordinatesForAddress(ctx)
	}

	r.handlers["geocoding_get_cache"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.geocodingRepo.GetGeocodingCache(ctx)
	}

	r.handlers["geocoding_get_details"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.geocodingRepo.GetGeocodingDetails(ctx)
	}

	r.handlers["geocoding_suggest_address"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		address := getStringParam(params, "address")
		if address == "" {
			return nil, fmt.Errorf("address is required")
		}
		return reg.geocodingRepo.SuggestAddress(ctx, address)
	}

	r.handlers["geocoding_suggest_city"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		name := getStringParam(params, "name")
		if name == "" {
			return nil, fmt.Errorf("name is required")
		}
		return reg.geocodingRepo.SuggestCity(ctx, name)
	}

	// Additional query methods for AI tooling
	r.handlers["geocoding_get_requests"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by", "created_at")
		sortOrder := getStringParam(params, "sort_order", "desc")
		return reg.geocodingRepo.GetGeocodingRequests(ctx, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["geocoding_search_with_term"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		term := getStringParam(params, "term")
		if term == "" {
			return nil, fmt.Errorf("term is required")
		}
		return reg.geocodingRepo.SearchGeocodingWithTerm(ctx, term)
	}

	r.handlers["geocoding_find_request"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		requestID := getStringParam(params, "request_id")
		if requestID == "" {
			return nil, fmt.Errorf("request_id is required")
		}
		return reg.geocodingRepo.FindGeocodingRequest(ctx, requestID)
	}

	r.handlers["geocoding_get_user_history"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.geocodingRepo.GetUserGeocodingHistory(ctx, userID, page, pageSize)
	}

	r.handlers["geocoding_get_by_location"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		lat := getFloat64Param(params, "lat", 0)
		lng := getFloat64Param(params, "lng", 0)
		radius := getInt64Param(params, "radius", 1000)
		if lat == 0 || lng == 0 {
			return nil, fmt.Errorf("lat and lng are required")
		}
		return reg.geocodingRepo.GetGeocodingByLocation(ctx, lat, lng, radius)
	}

	r.handlers["geocoding_clear_user_cache"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return nil, reg.geocodingRepo.ClearUserGeocodingCache(ctx, userID)
	}

	r.handlers["geocoding_get_stats"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.geocodingRepo.GetGeocodingStats(ctx)
	}
}