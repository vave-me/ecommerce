package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type GeocodingRepository interface {
	// Core geocoding operations from protobuf
	BatchGeocodeAddress(ctx context.Context, request *models.BatchGeocodeAddressRequest) (*models.BatchGeocodeAddressResponse, error)
	GeocodeAddress(ctx context.Context, id string, address string) (*models.GeocodeAddressResponse, error)
	RefreshGeocodingCache(ctx context.Context) (*models.RefreshGeocodingCacheResponse, error)
	ReverseGeocodeLocation(ctx context.Context) (*models.ReverseGeocodeLocationResponse, error)
	ValidateAddress(ctx context.Context) (*models.ValidateAddressResponse, error)
	GetAddressForCoordinates(ctx context.Context) (*models.GetAddressForCoordinatesResponse, error)
	GetCoordinatesForAddress(ctx context.Context) (*models.GetGeocodingDetailsResponse, error)
	GetGeocodingCache(ctx context.Context) (*models.GetGeocodingDetailsResponse, error)
	GetGeocodingDetails(ctx context.Context) (*models.GetGeocodingDetailsResponse, error)
	SuggestAddress(ctx context.Context, address string) (*models.SuggestAddressResponse, error)
	SuggestCity(ctx context.Context, name string) (*models.SuggestCityResponse, error)

	// Additional query methods for AI tooling and repository pattern compatibility
	// These would require additional RPC methods to be added to the protobuf
	GetGeocodingRequests(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*models.GeocodingRequest, error)
	SearchGeocodingWithTerm(ctx context.Context, term string) ([]*models.GeocodingRequest, error)
	FindGeocodingRequest(ctx context.Context, requestID string) (*models.GeocodingRequest, error)
	GetUserGeocodingHistory(ctx context.Context, userID string, page, pageSize int64) ([]*models.GeocodingRequest, error)
	GetGeocodingByLocation(ctx context.Context, lat, lng float64, radius int64) ([]*models.GeocodingRequest, error)
	ClearUserGeocodingCache(ctx context.Context, userID string) error
	GetGeocodingStats(ctx context.Context) (*models.GeocodingStatsResponse, error)
}
