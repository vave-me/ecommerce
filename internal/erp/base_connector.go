package erp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

// BaseConnector provides common functionality for all ERP connectors
type BaseConnector struct {
	Config      ERPConfig
	HTTPClient  *http.Client
	RateLimiter *rate.Limiter
	Logger      zerolog.Logger

	// Metrics
	requestCount    int64
	errorCount      int64
	lastHealthCheck time.Time
}

// NewBaseConnector creates a new base connector with production-ready configuration
func NewBaseConnector(config ERPConfig) *BaseConnector {
	// Configure HTTP transport with production settings
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   getTimeout(config),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	// Configure rate limiter with sensible defaults
	var rateLimiter *rate.Limiter
	if config.RateLimit != nil && config.RateLimit.RequestsPerSecond > 0 {
		rateLimiter = rate.NewLimiter(
			rate.Limit(config.RateLimit.RequestsPerSecond),
			config.RateLimit.BurstSize,
		)
	} else {
		// Default rate limit to prevent accidental overload
		rateLimiter = rate.NewLimiter(rate.Limit(10), 20)
	}

	// Configure logger with context
	logger := log.With().
		Str("erp_type", string(config.Type)).
		Str("endpoint", config.Endpoint).
		Str("component", "connector").
		Logger()

	return &BaseConnector{
		Config:      config,
		HTTPClient:  httpClient,
		RateLimiter: rateLimiter,
		Logger:      logger,
	}
}

// ExecuteWithRetry executes a function with exponential backoff retry logic
func (b *BaseConnector) ExecuteWithRetry(ctx context.Context, fn func() error) error {
	// Set default retry config if not provided
	if b.Config.Retry == nil {
		b.Config.Retry = &RetryConfig{
			MaxAttempts:  3,
			InitialDelay: 1 * time.Second,
			MaxDelay:     30 * time.Second,
			Multiplier:   2.0,
		}
	}

	var err error
	delay := b.Config.Retry.InitialDelay

	for attempt := 0; attempt < b.Config.Retry.MaxAttempts; attempt++ {
		// Wait before retry (except first attempt)
		if attempt > 0 {
			// Add jitter to prevent thundering herd
			jitter := time.Duration(rand.Float64() * float64(delay) * 0.1)
			actualDelay := delay + jitter

			b.Logger.Debug().
				Int("attempt", attempt).
				Dur("delay", actualDelay).
				Msg("Waiting before retry")

			select {
			case <-time.After(actualDelay):
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			}
		}

		// Execute the operation
		startTime := time.Now()
		err = fn()
		duration := time.Since(startTime)

		if err == nil {
			if attempt > 0 {
				b.Logger.Info().
					Int("attempt", attempt+1).
					Dur("duration", duration).
					Msg("Operation succeeded after retry")
			}
			return nil
		}

		// Check if error is retryable
		if !isRetryable(err) {
			b.Logger.Warn().
				Err(err).
				Int("attempt", attempt+1).
				Msg("Non-retryable error occurred")
			return err
		}

		b.Logger.Warn().
			Err(err).
			Int("attempt", attempt+1).
			Int("max_attempts", b.Config.Retry.MaxAttempts).
			Dur("duration", duration).
			Msg("Retryable error occurred")

		// Calculate next delay
		delay = time.Duration(float64(delay) * b.Config.Retry.Multiplier)
		if delay > b.Config.Retry.MaxDelay {
			delay = b.Config.Retry.MaxDelay
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w", b.Config.Retry.MaxAttempts, err)
}

// WaitForRateLimit waits for rate limit if configured
func (b *BaseConnector) WaitForRateLimit(ctx context.Context) error {
	if b.RateLimiter == nil {
		return nil
	}

	return b.RateLimiter.Wait(ctx)
}

// DoRequest executes an HTTP request with rate limiting and retries
func (b *BaseConnector) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Add standard headers
	req.Header.Set("User-Agent", fmt.Sprintf("ERP-Connector/%s", b.Config.Type))
	if req.Header.Get("Content-Type") == "" && req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add auth headers
	for k, v := range b.GetAuthHeaders() {
		req.Header.Set(k, v)
	}

	// Apply rate limiting
	if err := b.WaitForRateLimit(ctx); err != nil {
		return nil, err
	}

	var resp *http.Response
	operation := fmt.Sprintf("%s %s", req.Method, req.URL.Path)

	err := b.ExecuteWithRetry(ctx, func() error {
		var err error
		resp, err = b.HTTPClient.Do(req.WithContext(ctx))
		b.requestCount++

		if err != nil {
			return err
		}

		// Check for error status codes
		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			resp.Body.Close()

			// Rate limit error
			if resp.StatusCode == 429 {
				return &RateLimitError{
					StatusCode: resp.StatusCode,
					Message:    string(body),
					RetryAfter: resp.Header.Get("Retry-After"),
				}
			}

			// Server errors are retryable
			if resp.StatusCode >= 500 {
				return &ServerError{
					StatusCode: resp.StatusCode,
					Message:    string(body),
				}
			}

			// Client errors are not retryable
			return &ClientError{
				StatusCode: resp.StatusCode,
				Message:    string(body),
			}
		}

		return nil
	})

	if err != nil {
		b.errorCount++
		b.Logger.Error().
			Err(err).
			Str("operation", operation).
			Msg("Request failed")
	}

	return resp, err
}

// GetAuthHeaders returns authentication headers based on config
func (b *BaseConnector) GetAuthHeaders() map[string]string {
	headers := make(map[string]string)

	switch b.Config.Auth.Type {
	case "api_key":
		headers["X-API-Key"] = b.Config.Auth.APIKey
	case "basic":
		// Basic auth is handled by http.Request.SetBasicAuth
	case "oauth2":
		// OAuth2 token should be managed separately and added as Bearer token
	}

	return headers
}

// ParseJSONResponse parses JSON response with error handling
func (b *BaseConnector) ParseJSONResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	// Read body with size limit
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024)) // 10MB limit
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("parsing JSON response: %w", err)
	}

	return nil
}

// BuildURL constructs a URL with proper escaping
func (b *BaseConnector) BuildURL(path string, params map[string]string) (string, error) {
	baseURL, err := url.Parse(b.Config.Endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Parse the path
	pathURL, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	// Resolve the path against base URL
	fullURL := baseURL.ResolveReference(pathURL)

	// Add query parameters
	if len(params) > 0 {
		q := fullURL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		fullURL.RawQuery = q.Encode()
	}

	return fullURL.String(), nil
}

// Helper functions

func getTimeout(config ERPConfig) time.Duration {
	if config.Metadata != nil {
		if timeout, ok := config.Metadata["timeout"].(int); ok && timeout > 0 {
			return time.Duration(timeout) * time.Second
		}
	}
	return 30 * time.Second
}

// isRetryable determines if an error is retryable
func isRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Context errors are not retryable
	if err == context.Canceled || err == context.DeadlineExceeded {
		return false
	}

	// Check for specific error types
	switch e := err.(type) {
	case *ClientError:
		return false
	case *RateLimitError, *ServerError:
		return true
	case net.Error:
		return e.Temporary() || e.Timeout()
	}

	// Check for specific HTTP status codes in error message
	errStr := err.Error()
	if strings.Contains(errStr, "HTTP 4") && !strings.Contains(errStr, "HTTP 429") {
		return false // 4xx errors (except 429) are not retryable
	}

	return true
}

// Error types for better error handling

// ClientError represents a client-side error (4xx)
type ClientError struct {
	StatusCode int
	Message    string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("client error %d: %s", e.StatusCode, e.Message)
}

// ServerError represents a server-side error (5xx)
type ServerError struct {
	StatusCode int
	Message    string
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("server error %d: %s", e.StatusCode, e.Message)
}

// RateLimitError represents a rate limit error
type RateLimitError struct {
	StatusCode int
	Message    string
	RetryAfter string
}

func (e *RateLimitError) Error() string {
	if e.RetryAfter != "" {
		return fmt.Sprintf("rate limited %d: %s (retry after: %s)", e.StatusCode, e.Message, e.RetryAfter)
	}
	return fmt.Sprintf("rate limited %d: %s", e.StatusCode, e.Message)
}
