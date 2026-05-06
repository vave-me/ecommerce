package models

import "time"

// Core geocoding entities
type Address struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

type SuggestionAddress struct {
	SuggestedAddress string  `json:"suggested_address"`
	Latitude         float32 `json:"latitude"`
	Longitude        float32 `json:"longitude"`
}

type SuggestionCity struct {
	SuggestedCity string  `json:"suggested_city"`
	Latitude      float32 `json:"latitude"`
	Longitude     float32 `json:"longitude"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Domain entity for tracking geocoding requests
type GeocodingRequest struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Address   string    `json:"address"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Protobuf request/response types

// BatchGeocodeAddress
type BatchGeocodeAddressRequest struct {
	Addresses []Address `json:"addresses"`
	UserID    string    `json:"user_id"`
}

type BatchGeocodeAddressResponse struct {
	Results   []GeocodeAddressResponse `json:"results"`
	RequestID string                   `json:"request_id"`
	Status    string                   `json:"status"`
}

// GeocodeAddress
type GeocodeAddressRequest struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

type GeocodeAddressResponse struct {
	Address   string  `json:"address"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Status    string  `json:"status"`
	RequestID string  `json:"request_id"`
}

// RefreshGeocodingCache
type RefreshGeocodingCacheRequest struct {
	UserID string `json:"user_id"`
}

type RefreshGeocodingCacheResponse struct {
	Status       string `json:"status"`
	CacheCleared bool   `json:"cache_cleared"`
	Message      string `json:"message"`
}

// ReverseGeocodeLocation
type ReverseGeocodeLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UserID    string  `json:"user_id"`
}

type ReverseGeocodeLocationResponse struct {
	Address   string `json:"address"`
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
}

// ValidateAddress
type ValidateAddressRequest struct {
	Address string `json:"address"`
	UserID  string `json:"user_id"`
}

type ValidateAddressResponse struct {
	IsValid          bool   `json:"is_valid"`
	ValidatedAddress string `json:"validated_address"`
	Status           string `json:"status"`
	Message          string `json:"message"`
}

// SuggestAddress
type SuggestAddressRequest struct {
	Address string `json:"address"`
}

type SuggestAddressResponse struct {
	SuggestionAddresses []SuggestionAddress `json:"suggestion_addresses"`
	Status              string              `json:"status"`
}

// SuggestCity
type SuggestCityRequest struct {
	Name string `json:"name"`
}

type SuggestCityResponse struct {
	SuggestedCities []SuggestionCity `json:"suggested_cities"`
	Status          string           `json:"status"`
}

// GetAddressForCoordinates
type GetAddressForCoordinatesRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UserID    string  `json:"user_id"`
}

type GetAddressForCoordinatesResponse struct {
	Address   string `json:"address"`
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
}

// GetGeocodingDetails
type GetGeocodingDetailsRequest struct {
	RequestID string `json:"request_id"`
	UserID    string `json:"user_id"`
}

type GetGeocodingDetailsResponse struct {
	RequestID string  `json:"request_id"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// Additional response types for extended functionality
type GeocodingStatsResponse struct {
	TotalRequests      int64 `json:"total_requests"`
	SuccessfulRequests int64 `json:"successful_requests"`
	FailedRequests     int64 `json:"failed_requests"`
	CacheHits          int64 `json:"cache_hits"`
	CacheMisses        int64 `json:"cache_misses"`
}

// Status constants
const (
	GeocodingStatusPending   = "pending"
	GeocodingStatusSuccess   = "success"
	GeocodingStatusFailed    = "failed"
	GeocodingStatusCached    = "cached"
	GeocodingStatusValidated = "validated"
	GeocodingStatusInvalid   = "invalid"
)
