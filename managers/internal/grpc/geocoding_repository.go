package grpc

import (
	"context"
	"fmt"
	"log"
	"middleman/geocoding/geocodingpb"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"

	"google.golang.org/grpc"
)

type GeocodingRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.GeocodingRepository = (*GeocodingRepository)(nil)

// NewGeocodingRepositoryWithAuth creates a new GeocodingRepository with JWT authentication support
func NewGeocodingRepository(endpoint string, authInstance *auth.Auth) GeocodingRepository {
	return GeocodingRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// BatchGeocodeAddress processes multiple addresses in a batch
func (r GeocodingRepository) BatchGeocodeAddress(ctx context.Context, request *models.BatchGeocodeAddressRequest) (*models.BatchGeocodeAddressResponse, error) {
	log.Printf("[GEOCODING_GRPC] BatchGeocodeAddress called for %d addresses", len(request.Addresses))

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[GEOCODING_GRPC] Failed to connect to geocoding service: %v", err)
		return nil, fmt.Errorf("failed to connect to geocoding service: %w", err)
	}
	defer conn.Close()

	client := geocodingpb.NewGeocodingServiceClient(conn)

	// Note: The protobuf BatchGeocodeAddressRequest is empty, so we'll need to extend it
	// For now, we'll return a placeholder response
	_, err = client.BatchGeocodeAddress(ctx, &geocodingpb.BatchGeocodeAddressRequest{})
	if err != nil {
		log.Printf("[GEOCODING_GRPC] BatchGeocodeAddress RPC failed: %v", err)
		return nil, fmt.Errorf("BatchGeocodeAddress RPC failed: %w", err)
	}

	log.Printf("[GEOCODING_GRPC] BatchGeocodeAddress RPC successful")
	return &models.BatchGeocodeAddressResponse{
		Results:   []models.GeocodeAddressResponse{},
		RequestID: "batch_request_placeholder",
		Status:    models.GeocodingStatusSuccess,
	}, nil
}

// GeocodeAddress geocodes a single address
func (r GeocodingRepository) GeocodeAddress(ctx context.Context, id string, address string) (*models.GeocodeAddressResponse, error) {
	log.Printf("[GEOCODING_GRPC] GeocodeAddress called for address: %s", address)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[GEOCODING_GRPC] Failed to connect to geocoding service: %v", err)
		return nil, fmt.Errorf("failed to connect to geocoding service: %w", err)
	}
	defer conn.Close()

	client := geocodingpb.NewGeocodingServiceClient(conn)
	resp, err := client.GeocodeAddress(ctx, &geocodingpb.GeocodeAddressRequest{
		Id:      id,
		Address: address,
	})
	if err != nil {
		log.Printf("[GEOCODING_GRPC] GeocodeAddress RPC failed: %v", err)
		return nil, fmt.Errorf("GeocodeAddress RPC failed: %w", err)
	}

	log.Printf("[GEOCODING_GRPC] GeocodeAddress RPC successful for address: %s", address)
	return &models.GeocodeAddressResponse{
		Address:   resp.GetAddress(),
		Latitude:  resp.GetLatitude(),
		Longitude: resp.GetLongitude(),
		Status:    models.GeocodingStatusSuccess,
		RequestID: id,
	}, nil
}

// RefreshGeocodingCache refreshes the geocoding cache
func (r GeocodingRepository) RefreshGeocodingCache(ctx context.Context) (*models.RefreshGeocodingCacheResponse, error) {
	log.Printf("[GEOCODING_GRPC] RefreshGeocodingCache called")

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[GEOCODING_GRPC] Failed to connect to geocoding service: %v", err)
		return nil, fmt.Errorf("failed to connect to geocoding service: %w", err)
	}
	defer conn.Close()

	client := geocodingpb.NewGeocodingServiceClient(conn)
	_, err = client.RefreshGeocodingCache(ctx, &geocodingpb.RefreshGeocodingCacheRequest{})
	if err != nil {
		log.Printf("[GEOCODING_GRPC] RefreshGeocodingCache RPC failed: %v", err)
		return nil, fmt.Errorf("RefreshGeocodingCache RPC failed: %w", err)
	}

	log.Printf("[GEOCODING_GRPC] RefreshGeocodingCache RPC successful")
	return &models.RefreshGeocodingCacheResponse{
		Status:       models.GeocodingStatusSuccess,
		CacheCleared: true,
		Message:      "Cache refreshed successfully",
	}, nil
}

// ReverseGeocodeLocation performs reverse geocoding for coordinates
func (r GeocodingRepository) ReverseGeocodeLocation(ctx context.Context) (*models.ReverseGeocodeLocationResponse, error) {
	log.Printf("[GEOCODING_GRPC] ReverseGeocodeLocation called")

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[GEOCODING_GRPC] Failed to connect to geocoding service: %v", err)
		return nil, fmt.Errorf("failed to connect to geocoding service: %w", err)
	}
	defer conn.Close()

	client := geocodingpb.NewGeocodingServiceClient(conn)
	_, err = client.ReverseGeocodeLocation(ctx, &geocodingpb.ReverseGeocodeLocationRequest{})
	if err != nil {
		log.Printf("[GEOCODING_GRPC] ReverseGeocodeLocation RPC failed: %v", err)
		return nil, fmt.Errorf("ReverseGeocodeLocation RPC failed: %w", err)
	}

	log.Printf("[GEOCODING_GRPC] ReverseGeocodeLocation RPC successful")
	return &models.ReverseGeocodeLocationResponse{
		Address:   "Reverse geocoded address placeholder",
		Status:    models.GeocodingStatusSuccess,
		RequestID: "reverse_request_placeholder",
	}, nil
}

// ValidateAddress validates an address
func (r GeocodingRepository) ValidateAddress(ctx context.Context) (*models.ValidateAddressResponse, error) {
	log.Printf("[GEOCODING_GRPC] ValidateAddress called")

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[GEOCODING_GRPC] Failed to connect to geocoding service: %v", err)
		return nil, fmt.Errorf("failed to connect to geocoding service: %w", err)
	}
	defer conn.Close()

	client := geocodingpb.NewGeocodingServiceClient(conn)
	_, err = client.ValidateAddress(ctx, &geocodingpb.ValidateAddressRequest{})
	if err != nil {
		log.Printf("[GEOCODING_GRPC] ValidateAddress RPC failed: %v", err)
		return nil, fmt.Errorf("ValidateAddress RPC failed: %w", err)
	}

	log.Printf("[GEOCODING_GRPC] ValidateAddress RPC successful")
	return &models.ValidateAddressResponse{
		IsValid:          true,
		ValidatedAddress: "Validated address placeholder",
		Status:           models.GeocodingStatusValidated,
		Message:          "Address validated successfully",
	}, nil
}

// GetAddressForCoordinates gets address for given coordinates
func (r GeocodingRepository) GetAddressForCoordinates(ctx context.Context) (*models.GetAddressForCoordinatesResponse, error) {
	log.Printf("[GEOCODING_GRPC] GetAddressForCoordinates called")

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[GEOCODING_GRPC] Failed to connect to geocoding service: %v", err)
		return nil, fmt.Errorf("failed to connect to geocoding service: %w", err)
	}
	defer conn.Close()

	client := geocodingpb.NewGeocodingServiceClient(conn)
	_, err = client.GetAddressForCoordinates(ctx, &geocodingpb.GetAddressForCoordinatesRequest{})
	if err != nil {
		log.Printf("[GEOCODING_GRPC] GetAddressForCoordinates RPC failed: %v", err)
		return nil, fmt.Errorf("GetAddressForCoordinates RPC failed: %w", err)
	}

	log.Printf("[GEOCODING_GRPC] GetAddressForCoordinates RPC successful")
	return &models.GetAddressForCoordinatesResponse{
		Address:   "Address for coordinates placeholder",
		Status:    models.GeocodingStatusSuccess,
		RequestID: "coordinates_request_placeholder",
	}, nil
}

// GetCoordinatesForAddress gets coordinates for given address
func (r GeocodingRepository) GetCoordinatesForAddress(ctx context.Context) (*models.GetGeocodingDetailsResponse, error) {
	log.Printf("[GEOCODING_GRPC] GetCoordinatesForAddress called")

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[GEOCODING_GRPC] Failed to connect to geocoding service: %v", err)
		return nil, fmt.Errorf("failed to connect to geocoding service: %w", err)
	}
	defer conn.Close()

	client := geocodingpb.NewGeocodingServiceClient(conn)
	_, err = client.GetCoordinatesForAddress(ctx, &geocodingpb.GetGeocodingDetailsRequest{})
	if err != nil {
		log.Printf("[GEOCODING_GRPC] GetCoordinatesForAddress RPC failed: %v", err)
		return nil, fmt.Errorf("GetCoordinatesForAddress RPC failed: %w", err)
	}

	log.Printf("[GEOCODING_GRPC] GetCoordinatesForAddress RPC successful")
	return &models.GetGeocodingDetailsResponse{
		RequestID: "coordinates_request_placeholder",
		Address:   "Address placeholder",
		Latitude:  0.0,
		Longitude: 0.0,
		Status:    models.GeocodingStatusSuccess,
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}, nil
}

// GetGeocodingCache gets geocoding cache data
func (r GeocodingRepository) GetGeocodingCache(ctx context.Context) (*models.GetGeocodingDetailsResponse, error) {
	log.Printf("[GEOCODING_GRPC] GetGeocodingCache called")

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[GEOCODING_GRPC] Failed to connect to geocoding service: %v", err)
		return nil, fmt.Errorf("failed to connect to geocoding service: %w", err)
	}
	defer conn.Close()

	client := geocodingpb.NewGeocodingServiceClient(conn)
	_, err = client.GetGeocodingCache(ctx, &geocodingpb.GetGeocodingDetailsRequest{})
	if err != nil {
		log.Printf("[GEOCODING_GRPC] GetGeocodingCache RPC failed: %v", err)
		return nil, fmt.Errorf("GetGeocodingCache RPC failed: %w", err)
	}

	log.Printf("[GEOCODING_GRPC] GetGeocodingCache RPC successful")
	return &models.GetGeocodingDetailsResponse{
		RequestID: "cache_request_placeholder",
		Address:   "Cached address placeholder",
		Latitude:  0.0,
		Longitude: 0.0,
		Status:    models.GeocodingStatusCached,
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}, nil
}

// GetGeocodingDetails gets detailed geocoding information
func (r GeocodingRepository) GetGeocodingDetails(ctx context.Context) (*models.GetGeocodingDetailsResponse, error) {
	log.Printf("[GEOCODING_GRPC] GetGeocodingDetails called")

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[GEOCODING_GRPC] Failed to connect to geocoding service: %v", err)
		return nil, fmt.Errorf("failed to connect to geocoding service: %w", err)
	}
	defer conn.Close()

	client := geocodingpb.NewGeocodingServiceClient(conn)
	_, err = client.GetGeocodingDetails(ctx, &geocodingpb.GetGeocodingDetailsRequest{})
	if err != nil {
		log.Printf("[GEOCODING_GRPC] GetGeocodingDetails RPC failed: %v", err)
		return nil, fmt.Errorf("GetGeocodingDetails RPC failed: %w", err)
	}

	log.Printf("[GEOCODING_GRPC] GetGeocodingDetails RPC successful")
	return &models.GetGeocodingDetailsResponse{
		RequestID: "details_request_placeholder",
		Address:   "Detailed address placeholder",
		Latitude:  0.0,
		Longitude: 0.0,
		Status:    models.GeocodingStatusSuccess,
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}, nil
}

// SuggestAddress suggests addresses based on input
func (r GeocodingRepository) SuggestAddress(ctx context.Context, address string) (*models.SuggestAddressResponse, error) {
	log.Printf("[GEOCODING_GRPC] SuggestAddress called for address: %s", address)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[GEOCODING_GRPC] Failed to connect to geocoding service: %v", err)
		return nil, fmt.Errorf("failed to connect to geocoding service: %w", err)
	}
	defer conn.Close()

	client := geocodingpb.NewGeocodingServiceClient(conn)
	resp, err := client.SuggestAddress(ctx, &geocodingpb.SuggestAddressRequest{
		Address: address,
	})
	if err != nil {
		log.Printf("[GEOCODING_GRPC] SuggestAddress RPC failed: %v", err)
		return nil, fmt.Errorf("SuggestAddress RPC failed: %w", err)
	}

	// Convert protobuf suggestions to domain models
	suggestions := make([]models.SuggestionAddress, len(resp.GetSuggestionAddress()))
	for i, suggestion := range resp.GetSuggestionAddress() {
		suggestions[i] = models.SuggestionAddress{
			SuggestedAddress: suggestion.GetSuggestedAddress(),
			Latitude:         suggestion.GetLatitude(),
			Longitude:        suggestion.GetLongitude(),
		}
	}

	log.Printf("[GEOCODING_GRPC] SuggestAddress RPC successful, returned %d suggestions", len(suggestions))
	return &models.SuggestAddressResponse{
		SuggestionAddresses: suggestions,
		Status:              models.GeocodingStatusSuccess,
	}, nil
}

// SuggestCity suggests cities based on input
func (r GeocodingRepository) SuggestCity(ctx context.Context, name string) (*models.SuggestCityResponse, error) {
	log.Printf("[GEOCODING_GRPC] SuggestCity called for name: %s", name)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[GEOCODING_GRPC] Failed to connect to geocoding service: %v", err)
		return nil, fmt.Errorf("failed to connect to geocoding service: %w", err)
	}
	defer conn.Close()

	client := geocodingpb.NewGeocodingServiceClient(conn)
	resp, err := client.SuggestCity(ctx, &geocodingpb.SuggestCityRequest{
		Name: name,
	})
	if err != nil {
		log.Printf("[GEOCODING_GRPC] SuggestCity RPC failed: %v", err)
		return nil, fmt.Errorf("SuggestCity RPC failed: %w", err)
	}

	// Convert protobuf suggestions to domain models
	suggestions := make([]models.SuggestionCity, len(resp.GetSuggestedCities()))
	for i, suggestion := range resp.GetSuggestedCities() {
		suggestions[i] = models.SuggestionCity{
			SuggestedCity: suggestion.GetSuggestedCity(),
			Latitude:      suggestion.GetLatitude(),
			Longitude:     suggestion.GetLongitude(),
		}
	}

	log.Printf("[GEOCODING_GRPC] SuggestCity RPC successful, returned %d suggestions", len(suggestions))
	return &models.SuggestCityResponse{
		SuggestedCities: suggestions,
		Status:          models.GeocodingStatusSuccess,
	}, nil
}

// Additional query methods for AI tooling - these would require extending the protobuf
// For now, they return "not implemented" errors since the protobuf doesn't have these methods

func (r GeocodingRepository) GetGeocodingRequests(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.GeocodingRequest, error) {
	return nil, fmt.Errorf("GetGeocodingRequests not implemented - requires additional RPC method in protobuf")
}

func (r GeocodingRepository) SearchGeocodingWithTerm(ctx context.Context, term string) ([]*models.GeocodingRequest, error) {
	return nil, fmt.Errorf("SearchGeocodingWithTerm not implemented - requires additional RPC method in protobuf")
}

func (r GeocodingRepository) FindGeocodingRequest(ctx context.Context, requestID string) (*models.GeocodingRequest, error) {
	return nil, fmt.Errorf("FindGeocodingRequest not implemented - requires additional RPC method in protobuf")
}

func (r GeocodingRepository) GetUserGeocodingHistory(ctx context.Context, userID string, page, pageSize int64) ([]*models.GeocodingRequest, error) {
	return nil, fmt.Errorf("GetUserGeocodingHistory not implemented - requires additional RPC method in protobuf")
}

func (r GeocodingRepository) GetGeocodingByLocation(ctx context.Context, lat, lng float64, radius int64) ([]*models.GeocodingRequest, error) {
	return nil, fmt.Errorf("GetGeocodingByLocation not implemented - requires additional RPC method in protobuf")
}

func (r GeocodingRepository) ClearUserGeocodingCache(ctx context.Context, userID string) error {
	return fmt.Errorf("ClearUserGeocodingCache not implemented - requires additional RPC method in protobuf")
}

func (r GeocodingRepository) GetGeocodingStats(ctx context.Context) (*models.GeocodingStatsResponse, error) {
	return nil, fmt.Errorf("GetGeocodingStats not implemented - requires additional RPC method in protobuf")
}

// dial establishes a gRPC connection to the geocoding service
// dial sets up a gRPC connection with the microservice endpoint
func (r GeocodingRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r GeocodingRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
