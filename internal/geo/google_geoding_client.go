package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GoogleGeocodingClient encapsulates interactions with the Google Geocoding API.
type GoogleGeocodingClient struct {
	APIKey  string // Your Google API key
	BaseURL string // Typically "https://maps.googleapis.com/maps/api"
	Client  *http.Client
}

// NewGoogleGeocodingClient initializes a new GoogleGeocodingClient.
func NewGoogleGeocodingClient(apiKey string) *GoogleGeocodingClient {
	return &GoogleGeocodingClient{
		APIKey:  apiKey,
		BaseURL: "https://maps.googleapis.com/maps/api",
		Client:  &http.Client{},
	}
}

// ForwardGeocodeResponse is a simplified response structure for a forward geocoding request.
type ForwardGeocodeResponse struct {
	FormattedAddress string  `json:"formatted_address,omitempty"`
	Latitude         float64 `json:"lat,omitempty"`
	Longitude        float64 `json:"lng,omitempty"`
}

// geocodeAPIResponse models the raw JSON from the Google Geocoding API.
type geocodeAPIResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// ForwardGeocode takes an address string (e.g. "Berlin, Germany") and returns lat/lng plus
// a formatted address if successful.
func (c *GoogleGeocodingClient) ForwardGeocode(ctx context.Context, address string) (*ForwardGeocodeResponse, error) {
	// Build the request URL
	endpoint := fmt.Sprintf("%s/geocode/json?address=%s&key=%s",
		c.BaseURL,
		url.QueryEscape(address),
		url.QueryEscape(c.APIKey),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create forward geocode request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("forward geocode request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("forward geocode failed, status %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var data geocodeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode forward geocode response failed: %w", err)
	}

	// Check Google's status field
	if data.Status != "OK" {
		return nil, fmt.Errorf("forward geocode error: status=%s, message=%s", data.Status, data.ErrorMessage)
	}

	// Extract the first result, if any
	if len(data.Results) == 0 {
		return nil, fmt.Errorf("no geocoding results found for address: %s", address)
	}

	first := data.Results[0]
	return &ForwardGeocodeResponse{
		FormattedAddress: first.FormattedAddress,
		Latitude:         first.Geometry.Location.Lat,
		Longitude:        first.Geometry.Location.Lng,
	}, nil
}

// -----------------------------------------------------------------------------
// 2) REVERSE GEOCODING
// -----------------------------------------------------------------------------

// ReverseGeocodeResponse is a simplified response structure for a reverse geocoding request.
type ReverseGeocodeResponse struct {
	FormattedAddress string `json:"formatted_address,omitempty"`
}

// ReverseGeocode takes lat/lng and returns a best-match formatted address from Google.
func (c *GoogleGeocodingClient) ReverseGeocode(ctx context.Context, lat, lng float64) (*ReverseGeocodeResponse, error) {
	// Build the request URL
	endpoint := fmt.Sprintf("%s/geocode/json?latlng=%f,%f&key=%s",
		c.BaseURL,
		lat,
		lng,
		url.QueryEscape(c.APIKey),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create reverse geocode request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("reverse geocode request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("reverse geocode failed, status %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var data geocodeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode reverse geocode response failed: %w", err)
	}

	// Check Google's status field
	if data.Status != "OK" {
		return nil, fmt.Errorf("reverse geocode error: status=%s, message=%s", data.Status, data.ErrorMessage)
	}

	if len(data.Results) == 0 {
		return nil, fmt.Errorf("no reverse geocoding results found for coordinates: %f, %f", lat, lng)
	}

	first := data.Results[0]
	return &ReverseGeocodeResponse{
		FormattedAddress: first.FormattedAddress,
	}, nil
}
