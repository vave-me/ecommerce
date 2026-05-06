package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// NominatimGeocodingClient encapsulates interactions with a Nominatim server.
type NominatimGeocodingClient struct {
	// BaseURL could be "https://nominatim.openstreetmap.org" or
	// a self-hosted e.g. "http://nominatim:8080"
	BaseURL string
	Client  *http.Client
}

// NewNominatimGeocodingClient initializes a new NominatimGeocodingClient.
func NewNominatimGeocodingClient(baseURL string) *NominatimGeocodingClient {
	// Default to the public endpoint if baseURL is empty.
	if baseURL == "" {
		baseURL = "http://nominatim:8080"
	}
	return &NominatimGeocodingClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// -----------------------------------------------------------------------------
// 1) FORWARD GEOCODING
// -----------------------------------------------------------------------------

// ForwardNominatimResponse represents a simplified forward-geocoding result.
type ForwardNominatimResponse struct {
	DisplayName string  `json:"display_name,omitempty"`
	Latitude    float64 `json:"lat,omitempty"`
	Longitude   float64 `json:"lon,omitempty"`
}

// nominatimForwardAPIResponse models the raw JSON array from Nominatim's /search endpoint.
type nominatimForwardAPIResponse []struct {
	DisplayName string `json:"display_name"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	// Additional fields omitted for brevity
}

// ForwardGeocode queries Nominatim for an address/place string
// and returns the first match's lat/lon and display_name if found.
func (c *NominatimGeocodingClient) ForwardGeocode(ctx context.Context, address string) (*ForwardNominatimResponse, error) {
	// Build URL. e.g.: GET /search?q=Berlin&format=json
	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	// Usually /search endpoint for forward geocoding
	searchURL, err := base.Parse("/search")
	if err != nil {
		return nil, fmt.Errorf("failed to build Nominatim search URL: %w", err)
	}

	q := searchURL.Query()
	q.Set("q", address)
	q.Set("format", "json")
	q.Set("limit", "1")          // Return only the top result if you wish
	q.Set("addressdetails", "1") // Optionally get more details
	searchURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL.String(), nil)
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

	var data nominatimForwardAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode forward geocode response failed: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	first := data[0]

	// Convert lat/lon from string to float64
	lat, lon, err := parseLatLon(first.Lat, first.Lon)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lat/lon from result: %w", err)
	}

	return &ForwardNominatimResponse{
		DisplayName: first.DisplayName,
		Latitude:    lat,
		Longitude:   lon,
	}, nil
}

// -----------------------------------------------------------------------------
// 2) REVERSE GEOCODING
// -----------------------------------------------------------------------------

// ReverseNominatimResponse is a simplified structure for a reverse geocoding request result.
type ReverseNominatimResponse struct {
	DisplayName string `json:"display_name,omitempty"`
}

// nominatimReverseAPIResponse models the raw JSON from Nominatim's /reverse endpoint.
type nominatimReverseAPIResponse struct {
	DisplayName string `json:"display_name"`
	// Additional fields can go here
	Error string `json:"error,omitempty"`
}

// ReverseGeocode queries Nominatim to get a best-match address for lat/lon.
func (c *NominatimGeocodingClient) ReverseGeocode(ctx context.Context, lat, lon float64) (*ReverseNominatimResponse, error) {
	// e.g.: GET /reverse?lat=...&lon=...&format=json
	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	reverseURL, err := base.Parse("/reverse")
	if err != nil {
		return nil, fmt.Errorf("failed to build Nominatim reverse URL: %w", err)
	}

	q := reverseURL.Query()
	q.Set("lat", fmt.Sprintf("%f", lat))
	q.Set("lon", fmt.Sprintf("%f", lon))
	q.Set("format", "json")
	// optionally: q.Set("addressdetails", "1")
	reverseURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reverseURL.String(), nil)
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

	var data nominatimReverseAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode reverse geocode response failed: %w", err)
	}

	// If the 'error' field is present, it typically means no result found
	if data.Error != "" {
		return nil, fmt.Errorf("no reverse geocoding result found: %s", data.Error)
	}

	return &ReverseNominatimResponse{
		DisplayName: data.DisplayName,
	}, nil
}

// Suggestion represents one potential matching location.
type Suggestion struct {
	DisplayName string  `json:"display_name"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lon"`
}

// SuggestAddresses calls `/search` with a partial query and returns multiple results
// (instead of a single top result). This can be used for 'autocomplete' style suggestions.
func (c *NominatimGeocodingClient) SuggestAddresses(ctx context.Context, partialQuery string, limit int) ([]Suggestion, error) {
	if limit <= 0 {
		limit = 5 // default to 5 suggestions if not specified
	}

	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	searchURL, err := base.Parse("/search")
	if err != nil {
		return nil, fmt.Errorf("failed to build Nominatim search URL: %w", err)
	}

	q := searchURL.Query()
	q.Set("q", partialQuery)
	q.Set("format", "json")
	q.Set("limit", fmt.Sprintf("%d", limit))
	// For partial matches, Nominatim sometimes uses partial word expansions.
	// You can also consider 'addressdetails=1' if you want more detail.
	searchURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create suggest addresses request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("suggest addresses request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("suggest addresses failed, status %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var data nominatimForwardAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode suggest addresses response failed: %w", err)
	}

	if len(data) == 0 {
		return nil, nil // or return an empty slice
	}

	// Convert the raw results to our Suggestion type
	suggestions := make([]Suggestion, 0, len(data))
	for _, item := range data {
		latVal, lonVal, parseErr := parseLatLon(item.Lat, item.Lon)
		if parseErr != nil {
			// if one fails, skip it or handle differently
			continue
		}
		suggestions = append(suggestions, Suggestion{
			DisplayName: item.DisplayName,
			Latitude:    latVal,
			Longitude:   lonVal,
		})
	}
	return suggestions, nil
}

// SuggestCities decides whether to treat partialCity as a zip code (leading digit)
// or partial city name (leading letter). Then it queries the Nominatim `/search`
// endpoint. For zip code search we set `postalcode=partialCity`; for city we set `q=partialCity`.
func (c *NominatimGeocodingClient) SuggestCities(ctx context.Context, partialCity string, limit int) ([]Suggestion, error) {
	if limit <= 0 {
		limit = 5 // default to 5 suggestions
	}

	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	searchURL, err := base.Parse("/search")
	if err != nil {
		return nil, fmt.Errorf("failed to build Nominatim search URL: %w", err)
	}

	q := searchURL.Query()

	// If the first character is digit, treat partialCity as postal code,
	// else treat as partial city name.
	if len(partialCity) > 0 && isDigit(partialCity[0]) {
		// Attempt direct postal code matching
		q.Set("postalcode", partialCity)
		// In many Nominatim setups, partial postal code might not yield partial matches,
		// but we rely on the server's behavior.
	} else {
		// Use free-text `q` param for city partial matching
		q.Set("q", partialCity)
	}

	q.Set("format", "json")
	q.Set("limit", strconv.Itoa(limit))
	// Optional: q.Set("addressdetails", "1") // if you want more detail

	searchURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create suggest cities request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("suggest cities request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("suggest cities failed, status %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var data nominatimForwardAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode suggest cities response failed: %w", err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	suggestions := make([]Suggestion, 0, len(data))
	for _, item := range data {
		latVal, lonVal, parseErr := parseLatLon(item.Lat, item.Lon)
		if parseErr != nil {
			// skip invalid lat/lon
			continue
		}
		label := formatZipCity(item.DisplayName)
		suggestions = append(suggestions, Suggestion{
			DisplayName: label,
			Latitude:    latVal,
			Longitude:   lonVal,
		})
	}
	return suggestions, nil
}

// isDigit checks if b is '0'..'9'.
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// formatZipCity tries to produce "ZIP - CITY" from typical Nominatim strings
// e.g. "Berlin, Brandenburg, 10117, Deutschland" => "10117 - Berlin".
func formatZipCity(fullDisplay string) string {
	parts := strings.Split(fullDisplay, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	if len(parts) >= 3 {
		city := parts[0]
		zip := parts[len(parts)-2]
		if looksLikeZip(zip) {
			return fmt.Sprintf("%s - %s", zip, city)
		}
	}
	return fullDisplay
}

// looksLikeZip does a naive check if s is numeric and length 3..10.
func looksLikeZip(s string) bool {
	if len(s) < 3 || len(s) > 10 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// ValidateAddress checks if the given address string is recognized by Nominatim.
// Returns true if a forward geocode result is found, false if not.
func (c *NominatimGeocodingClient) ValidateAddress(ctx context.Context, address string) (bool, error) {
	// We'll just re-use the ForwardGeocode logic, which returns an error if no results found
	_, err := c.ForwardGeocode(ctx, address)
	if err != nil {
		// You can parse the error message to differentiate "no results" vs other errors.
		return false, nil // or return false, err if you want to treat network errors differently
	}
	return true, nil
}

// parseLatLon is a small helper to parse lat/lon from string to float64.
func parseLatLon(latStr, lonStr string) (float64, float64, error) {
	lat, err := parseFloat(latStr)
	if err != nil {
		return 0, 0, err
	}
	lon, err := parseFloat(lonStr)
	if err != nil {
		return 0, 0, err
	}
	return lat, lon, nil
}

func parseFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse float from %q: %w", s, err)
	}
	return f, nil
}
