package tools

import (
	"context"
	"fmt"
	"time"

	"middleman/managers/internal/domain"
)

// GeocodingToolService handles location and geocoding operations
type GeocodingToolService struct {
	geocodingRepo domain.GeocodingRepository
	config        *ServiceConfig
}

// NewGeocodingToolService creates a new geocoding tool service
func NewGeocodingToolService(geocodingRepo domain.GeocodingRepository) *GeocodingToolService {
	return &GeocodingToolService{
		geocodingRepo: geocodingRepo,
		config: &ServiceConfig{
			MaxRetries:      3,
			EnableStreaming: true,
			EnableMetrics:   true,
		},
	}
}

// ExecuteOperation routes geocoding operations to appropriate handlers
func (s *GeocodingToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Send initial progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "geocoding_operation",
			Status:   "started",
			Progress: 0,
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "GeocodingToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	var result interface{}
	var err error

	switch operation {
	case "geocode", "encode":
		result, err = s.geocode(ctx, parameters, streamChan, toolID)
	case "reverse_geocode", "decode":
		result, err = s.reverseGeocode(ctx, parameters, streamChan, toolID)
	case "get_coordinates":
		result, err = s.getCoordinates(ctx, parameters, streamChan, toolID)
	case "get_address":
		result, err = s.getAddress(ctx, parameters, streamChan, toolID)
	case "search_places":
		result, err = s.searchPlaces(ctx, parameters, streamChan, toolID)
	case "get_nearby":
		result, err = s.getNearby(ctx, parameters, streamChan, toolID)
	case "validate_address":
		result, err = s.validateAddress(ctx, parameters, streamChan, toolID)
	case "normalize_address":
		result, err = s.normalizeAddress(ctx, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported geocoding operation: %s", operation)
	}

	// Send completion status
	if streamChan != nil {
		status := "completed"
		if err != nil {
			status = "error"
		}

		errorStr := ""
		if err != nil {
			errorStr = err.Error()
		}

		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "geocoding_operation",
			Status:   status,
			Progress: 100,
			Result:   result,
			Error:    errorStr,
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "GeocodingToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return result, err
}

func (s *GeocodingToolService) geocode(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	address := getStringParam(params, "address", "")
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	// Send progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "geocoding_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "geocoding_address",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	response, err := s.geocodingRepo.GeocodeAddress(ctx, "", address)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode address: %w", err)
	}

	return map[string]interface{}{
		"address":     address,
		"coordinates": response,
	}, nil
}

func (s *GeocodingToolService) reverseGeocode(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	lat := getFloat32Param(params, "lat", 0.0)
	lng := getFloat32Param(params, "lng", 0.0)

	if lat == 0.0 || lng == 0.0 {
		return nil, fmt.Errorf("lat and lng are required")
	}

	// Send progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "geocoding_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "reverse_geocoding",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	response, err := s.geocodingRepo.ReverseGeocodeLocation(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to reverse geocode coordinates: %w", err)
	}

	return map[string]interface{}{
		"lat":     lat,
		"lng":     lng,
		"address": response,
	}, nil
}

func (s *GeocodingToolService) getCoordinates(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	address := getStringParam(params, "address", "")
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	// Send progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "geocoding_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "getting_coordinates",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	response, err := s.geocodingRepo.GetCoordinatesForAddress(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get coordinates: %w", err)
	}

	return map[string]interface{}{
		"address":  address,
		"response": response,
	}, nil
}

func (s *GeocodingToolService) getAddress(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	lat := getFloat32Param(params, "lat", 0.0)
	lng := getFloat32Param(params, "lng", 0.0)

	if lat == 0.0 || lng == 0.0 {
		return nil, fmt.Errorf("lat and lng are required")
	}

	// Send progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "geocoding_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "getting_address",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	response, err := s.geocodingRepo.GetAddressForCoordinates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}

	return map[string]interface{}{
		"lat":      lat,
		"lng":      lng,
		"response": response,
	}, nil
}

func (s *GeocodingToolService) searchPlaces(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	query := getStringParam(params, "query", "")
	if query == "" {
		query = getStringParam(params, "search_term", "")
	}
	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	limit := getInt64Param(params, "limit", 10)
	radius := getInt64Param(params, "radius", 5000)
	lat := getFloat32Param(params, "lat", 0.0)
	lng := getFloat32Param(params, "lng", 0.0)

	// Send progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "geocoding_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "searching_places",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Use SearchGeocodingWithTerm instead of SearchPlaces
	places, err := s.geocodingRepo.SearchGeocodingWithTerm(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search places: %w", err)
	}

	return map[string]interface{}{
		"query":  query,
		"places": places,
		"lat":    lat,
		"lng":    lng,
		"radius": radius,
		"limit":  limit,
	}, nil
}

func (s *GeocodingToolService) getNearby(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	lat := getFloat32Param(params, "lat", 0.0)
	lng := getFloat32Param(params, "lng", 0.0)
	radius := getInt64Param(params, "radius", 1000)
	placeType := getStringParam(params, "place_type", "")

	if lat == 0.0 || lng == 0.0 {
		return nil, fmt.Errorf("lat and lng are required")
	}

	// Send progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "geocoding_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "getting_nearby_places",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Use GetGeocodingByLocation instead of GetNearby
	places, err := s.geocodingRepo.GetGeocodingByLocation(ctx, float64(lat), float64(lng), radius)
	if err != nil {
		return nil, fmt.Errorf("failed to get nearby places: %w", err)
	}

	return map[string]interface{}{
		"lat":        lat,
		"lng":        lng,
		"radius":     radius,
		"place_type": placeType,
		"places":     places,
	}, nil
}

func (s *GeocodingToolService) validateAddress(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	address := getStringParam(params, "address", "")
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	// Send progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "geocoding_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "validating_address",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	response, err := s.geocodingRepo.ValidateAddress(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to validate address: %w", err)
	}

	return map[string]interface{}{
		"address":            address,
		"is_valid":           response.IsValid,
		"normalized_address": response.ValidatedAddress,
		"status":             response.Status,
		"message":            response.Message,
	}, nil
}

func (s *GeocodingToolService) normalizeAddress(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	address := getStringParam(params, "address", "")
	if address == "" {
		return nil, fmt.Errorf("address is required")
	}

	// Send progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "geocoding_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "normalizing_address",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Use ValidateAddress instead of NormalizeAddress since it's not available
	response, err := s.geocodingRepo.ValidateAddress(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize address: %w", err)
	}

	return map[string]interface{}{
		"address":            address,
		"normalized_address": response.ValidatedAddress,
		"status":             response.Status,
		"message":            response.Message,
	}, nil
}
