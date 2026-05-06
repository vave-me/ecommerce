package dynamic365

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Response represents a response from the Business Central API
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Value      interface{}
	NextLink   string
	Context    string
}

// parseResponse parses an HTTP response into a Response struct
func parseResponse(resp *http.Response) (*Response, error) {
	defer resp.Body.Close()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		return nil, parseError(resp, body)
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}

	// Parse JSON response
	if len(body) > 0 && strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("parsing JSON response: %w", err)
		}

		// Extract value
		if value, ok := result["value"]; ok {
			response.Value = value
		} else {
			response.Value = result
		}

		// Extract OData metadata
		if context, ok := result["@odata.context"].(string); ok {
			response.Context = context
		}

		if nextLink, ok := result["@odata.nextLink"].(string); ok {
			response.NextLink = nextLink
		}
	}

	return response, nil
}

// Unmarshal unmarshals the response value into the provided interface
func (r *Response) Unmarshal(v interface{}) error {
	if r.Value == nil {
		return fmt.Errorf("no value in response")
	}

	// Convert back to JSON then unmarshal to handle type conversions
	data, err := json.Marshal(r.Value)
	if err != nil {
		return fmt.Errorf("marshaling value: %w", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshaling value: %w", err)
	}

	return nil
}

// GetItems returns the items from a collection response
func (r *Response) GetItems() ([]interface{}, error) {
	if r.Value == nil {
		return nil, fmt.Errorf("no value in response")
	}

	items, ok := r.Value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value is not a collection")
	}

	return items, nil
}

// HasMore returns true if there are more pages available
func (r *Response) HasMore() bool {
	return r.NextLink != ""
}

// APIError represents an error from the Business Central API
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Details    []ErrorDetail
	Target     string
	InnerError *InnerError
	RequestID  string
}

// ErrorDetail represents additional error details
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target,omitempty"`
}

// InnerError represents nested error information
type InnerError struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	InnerError *InnerError `json:"innerError,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.InnerError != nil {
		return fmt.Sprintf("API error %d: %s - %s (inner: %s)",
			e.StatusCode, e.Code, e.Message, e.InnerError.Message)
	}
	return fmt.Sprintf("API error %d: %s - %s", e.StatusCode, e.Code, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsConflict returns true if the error is a 409 Conflict
func (e *APIError) IsConflict() bool {
	return e.StatusCode == http.StatusConflict
}

// IsBadRequest returns true if the error is a 400 Bad Request
func (e *APIError) IsBadRequest() bool {
	return e.StatusCode == http.StatusBadRequest
}

// IsUnauthorized returns true if the error is a 401 Unauthorized
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsForbidden returns true if the error is a 403 Forbidden
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

// IsRateLimited returns true if the error is a 429 Too Many Requests
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// parseError parses an error response with body already read
func parseError(resp *http.Response, body []byte) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		RequestID:  resp.Header.Get("X-Request-ID"),
	}

	// Try to parse OData error format
	var odataError struct {
		Error struct {
			Code       string        `json:"code"`
			Message    string        `json:"message"`
			Target     string        `json:"target,omitempty"`
			Details    []ErrorDetail `json:"details,omitempty"`
			InnerError *InnerError   `json:"innererror,omitempty"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &odataError); err == nil && odataError.Error.Code != "" {
		apiErr.Code = odataError.Error.Code
		apiErr.Message = odataError.Error.Message
		apiErr.Target = odataError.Error.Target
		apiErr.Details = odataError.Error.Details
		apiErr.InnerError = odataError.Error.InnerError
		return apiErr
	}

	// Try simple error format
	var simpleError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &simpleError); err == nil && simpleError.Code != "" {
		apiErr.Code = simpleError.Code
		apiErr.Message = simpleError.Message
		return apiErr
	}

	// Default error
	apiErr.Code = http.StatusText(resp.StatusCode)
	apiErr.Message = string(body)
	if apiErr.Message == "" {
		apiErr.Message = "No error details provided"
	}

	return apiErr
}

// parseError with http.Response (overload for backward compatibility)
func parseErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading error response: %w", err)
	}
	return parseError(resp, body)
}

// BatchResponse represents a response from a batch request
type BatchResponse struct {
	Responses []BatchItemResponse `json:"responses"`
}

// BatchItemResponse represents a single response in a batch
type BatchItemResponse struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}

// PageIterator provides iteration over paged results
type PageIterator struct {
	client   *Client
	nextLink string
	finished bool
}

// NewPageIterator creates a new page iterator from a response
func (c *Client) NewPageIterator(resp *Response) *PageIterator {
	return &PageIterator{
		client:   c,
		nextLink: resp.NextLink,
		finished: resp.NextLink == "",
	}
}

// HasMore returns true if there are more pages
func (pi *PageIterator) HasMore() bool {
	return !pi.finished
}

// Next fetches the next page of results
func (pi *PageIterator) Next(ctx context.Context) (*Response, error) {
	if pi.finished {
		return nil, fmt.Errorf("no more pages")
	}

	// Extract path from next link
	u, err := parseURL(pi.nextLink)
	if err != nil {
		return nil, fmt.Errorf("parsing next link: %w", err)
	}

	path := u.Path
	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}

	resp, err := pi.client.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	pi.nextLink = resp.NextLink
	pi.finished = resp.NextLink == ""

	return resp, nil
}

// parseURL safely parses a URL
func parseURL(urlStr string) (*url, error) {
	parts := strings.SplitN(urlStr, "?", 2)
	result := &url{Path: parts[0]}
	if len(parts) > 1 {
		result.RawQuery = parts[1]
	}
	return result, nil
}

type url struct {
	Path     string
	RawQuery string
}
