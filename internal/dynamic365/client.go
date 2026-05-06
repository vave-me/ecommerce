package dynamic365

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

// Client represents a production-ready Dynamics 365 Business Central client
type Client struct {
	// Configuration
	tenantID    string
	environment string
	companyID   string
	apiVersion  string
	baseURL     string

	// Authentication
	msalClient confidential.Client
	tokenCache *tokenCache

	// HTTP client with connection pooling
	httpClient *http.Client

	// Rate limiting
	rateLimiter *rate.Limiter

	// Logging
	logger zerolog.Logger

	// Metrics
	metrics *clientMetrics

	// Request ID generator
	requestIDGen *requestIDGenerator
}

// ClientConfig holds the configuration for the D365 client
type ClientConfig struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	Environment  string // "production" or "sandbox"
	CompanyID    string // Optional: specific company ID
	APIVersion   string // Optional: defaults to "v2.0"

	// HTTP configuration
	Timeout      time.Duration
	MaxIdleConns int

	// Rate limiting
	RequestsPerSecond float64
	BurstSize         int

	// Logging
	Logger zerolog.Logger
}

// NewClient creates a new production-ready Dynamics 365 Business Central client
func NewClient(ctx context.Context, config ClientConfig) (*Client, error) {
	// Validate configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Set defaults
	if config.APIVersion == "" {
		config.APIVersion = "v2.0"
	}
	if config.Environment == "" {
		config.Environment = "production"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 100
	}
	if config.RequestsPerSecond == 0 {
		config.RequestsPerSecond = 10
	}
	if config.BurstSize == 0 {
		config.BurstSize = 20
	}

	// Create MSAL confidential client
	cred, err := confidential.NewCredFromSecret(config.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("creating credential: %w", err)
	}

	authority := fmt.Sprintf("https://login.microsoftonline.com/%s", config.TenantID)
	msalClient, err := confidential.New(
		config.ClientID,
		cred,
		confidential.WithAuthority(authority),
		confidential.WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("creating MSAL client: %w", err)
	}

	// Create HTTP client with production settings
	httpClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:          config.MaxIdleConns,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableCompression:    false,
			DisableKeepAlives:     false,
		},
	}

	// Build base URL
	baseURL := fmt.Sprintf("https://api.businesscentral.dynamics.com/v2.0/%s/%s/api/%s",
		config.TenantID,
		config.Environment,
		config.APIVersion,
	)

	client := &Client{
		tenantID:     config.TenantID,
		environment:  config.Environment,
		companyID:    config.CompanyID,
		apiVersion:   config.APIVersion,
		baseURL:      baseURL,
		msalClient:   msalClient,
		tokenCache:   newTokenCache(),
		httpClient:   httpClient,
		rateLimiter:  rate.NewLimiter(rate.Limit(config.RequestsPerSecond), config.BurstSize),
		logger:       config.Logger,
		metrics:      newClientMetrics(),
		requestIDGen: newRequestIDGenerator(),
	}

	// Test authentication
	if _, err := client.getAccessToken(ctx); err != nil {
		return nil, fmt.Errorf("initial authentication failed: %w", err)
	}

	client.logger.Info().
		Str("tenant_id", config.TenantID).
		Str("environment", config.Environment).
		Str("api_version", config.APIVersion).
		Msg("Dynamics 365 Business Central client initialized")

	return client, nil
}

// getAccessToken retrieves an access token using MSAL
func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	// Check cache first
	if token := c.tokenCache.Get(); token != "" {
		return token, nil
	}

	// Acquire new token
	scopes := []string{"https://api.businesscentral.dynamics.com/.default"}

	result, err := c.msalClient.AcquireTokenByCredential(ctx, scopes)
	if err != nil {
		c.metrics.incrementAuthErrors()
		return "", fmt.Errorf("acquiring token: %w", err)
	}

	// Cache the token
	c.tokenCache.Set(result.AccessToken, result.ExpiresOn)
	c.metrics.incrementTokenAcquisitions()

	c.logger.Debug().
		Time("expires_on", result.ExpiresOn).
		Msg("Acquired new access token")

	return result.AccessToken, nil
}

// Request performs an authenticated HTTP request to the Business Central API
func (c *Client) Request(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	// Apply rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter: %w", err)
	}

	// Get access token
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting access token: %w", err)
	}

	// Build URL
	url := c.buildURL(path)

	// Prepare request body
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "D365-BC-Go-Client/1.0")
	req.Header.Set("X-Request-ID", c.requestIDGen.Generate())

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("If-Match", "*") // For updates
	}

	// Add OData headers
	req.Header.Set("OData-MaxVersion", "4.0")
	req.Header.Set("OData-Version", "4.0")

	// Log request
	c.logger.Debug().
		Str("method", method).
		Str("url", url).
		Str("request_id", req.Header.Get("X-Request-ID")).
		Msg("Sending API request")

	// Execute request with retry
	resp, err := c.executeWithRetry(ctx, req)
	if err != nil {
		c.metrics.incrementErrors()
		return nil, err
	}

	c.metrics.incrementRequests()

	// Log response
	c.logger.Debug().
		Str("method", method).
		Str("url", url).
		Int("status_code", resp.StatusCode).
		Str("request_id", req.Header.Get("X-Request-ID")).
		Msg("Received API response")

	return resp, nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, options ...RequestOption) (*Response, error) {
	// Apply options
	opts := &requestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	// Build query parameters
	if opts.queryParams != nil {
		path = fmt.Sprintf("%s?%s", path, opts.queryParams.Encode())
	}

	resp, err := c.Request(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	return parseResponse(resp)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	resp, err := c.Request(ctx, "POST", path, body)
	if err != nil {
		return nil, err
	}

	return parseResponse(resp)
}

// Patch performs a PATCH request
func (c *Client) Patch(ctx context.Context, path string, body interface{}) (*Response, error) {
	resp, err := c.Request(ctx, "PATCH", path, body)
	if err != nil {
		return nil, err
	}

	return parseResponse(resp)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string) error {
	resp, err := c.Request(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return parseError(resp)
	}

	return nil
}

// buildURL constructs the full URL for an API endpoint
func (c *Client) buildURL(path string) string {
	// Remove leading slash if present
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// Add company ID if specified and not already in path
	if c.companyID != "" && !contains(path, "companies(") {
		return fmt.Sprintf("%s/companies(%s)/%s", c.baseURL, c.companyID, path)
	}

	return fmt.Sprintf("%s/%s", c.baseURL, path)
}

// executeWithRetry performs a request with exponential backoff retry
func (c *Client) executeWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	maxRetries := 3
	backoff := 1 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(backoff):
				backoff *= 2
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if attempt == maxRetries-1 {
				return nil, err
			}
			continue
		}

		// Check if we should retry
		if shouldRetry(resp.StatusCode) && attempt < maxRetries-1 {
			resp.Body.Close()

			// Handle rate limiting
			if resp.StatusCode == http.StatusTooManyRequests {
				if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
					if duration, err := time.ParseDuration(retryAfter + "s"); err == nil {
						backoff = duration
					}
				}
			}

			c.logger.Warn().
				Int("attempt", attempt+1).
				Int("status_code", resp.StatusCode).
				Dur("backoff", backoff).
				Msg("Retrying request")

			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("max retries exceeded")
}

// tokenCache provides thread-safe token caching
type tokenCache struct {
	mu        sync.RWMutex
	token     string
	expiresOn time.Time
}

func newTokenCache() *tokenCache {
	return &tokenCache{}
}

func (tc *tokenCache) Get() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// Check if token is still valid (with 5 minute buffer)
	if time.Now().Add(5 * time.Minute).After(tc.expiresOn) {
		return ""
	}

	return tc.token
}

func (tc *tokenCache) Set(token string, expiresOn time.Time) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.token = token
	tc.expiresOn = expiresOn
}

// clientMetrics tracks client metrics
type clientMetrics struct {
	mu                sync.RWMutex
	totalRequests     int64
	totalErrors       int64
	authErrors        int64
	tokenAcquisitions int64
}

func newClientMetrics() *clientMetrics {
	return &clientMetrics{}
}

func (m *clientMetrics) incrementRequests() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalRequests++
}

func (m *clientMetrics) incrementErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalErrors++
}

func (m *clientMetrics) incrementAuthErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.authErrors++
}

func (m *clientMetrics) incrementTokenAcquisitions() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokenAcquisitions++
}

// GetMetrics returns current metrics
func (c *Client) GetMetrics() map[string]int64 {
	c.metrics.mu.RLock()
	defer c.metrics.mu.RUnlock()

	return map[string]int64{
		"total_requests":     c.metrics.totalRequests,
		"total_errors":       c.metrics.totalErrors,
		"auth_errors":        c.metrics.authErrors,
		"token_acquisitions": c.metrics.tokenAcquisitions,
	}
}

// requestIDGenerator generates unique request IDs
type requestIDGenerator struct {
	mu      sync.Mutex
	counter int64
}

func newRequestIDGenerator() *requestIDGenerator {
	return &requestIDGenerator{}
}

func (g *requestIDGenerator) Generate() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.counter++
	return fmt.Sprintf("req-%d-%d", time.Now().Unix(), g.counter)
}

// Helper functions

func validateConfig(config ClientConfig) error {
	if config.TenantID == "" {
		return fmt.Errorf("tenant ID is required")
	}
	if config.ClientID == "" {
		return fmt.Errorf("client ID is required")
	}
	if config.ClientSecret == "" {
		return fmt.Errorf("client secret is required")
	}
	return nil
}

func shouldRetry(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusRequestTimeout ||
		statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusGatewayTimeout
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(substr) < len(s) && containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 1; i < len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// RequestOption is a function that modifies request options
type RequestOption func(*requestOptions)

type requestOptions struct {
	queryParams url.Values
}

// WithQueryParams adds query parameters to the request
func WithQueryParams(params map[string]string) RequestOption {
	return func(opts *requestOptions) {
		if opts.queryParams == nil {
			opts.queryParams = url.Values{}
		}
		for k, v := range params {
			opts.queryParams.Set(k, v)
		}
	}
}

// WithFilter adds an OData filter to the request
func WithFilter(filter string) RequestOption {
	return WithQueryParams(map[string]string{"$filter": filter})
}

// WithSelect adds an OData select to the request
func WithSelect(fields ...string) RequestOption {
	return WithQueryParams(map[string]string{"$select": joinStrings(fields, ",")})
}

// WithExpand adds an OData expand to the request
func WithExpand(fields ...string) RequestOption {
	return WithQueryParams(map[string]string{"$expand": joinStrings(fields, ",")})
}

// WithTop adds an OData top to the request
func WithTop(count int) RequestOption {
	return WithQueryParams(map[string]string{"$top": fmt.Sprint(count)})
}

// WithOrderBy adds an OData orderby to the request
func WithOrderBy(orderBy string) RequestOption {
	return WithQueryParams(map[string]string{"$orderby": orderBy})
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
