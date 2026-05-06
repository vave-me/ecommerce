package adapters

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/api/content/v2.1"
	"google.golang.org/api/googleapi"
	"middleman/internal/merchant"
	"middleman/merchant/internal/domain"
)

// MerchantClientWrapper wraps the merchant client with retry logic and monitoring
type MerchantClientWrapper struct {
	client *merchant.MerchantCenterClient
	logger zerolog.Logger
	
	// Retry configuration
	maxRetries int
	baseDelay  time.Duration
}

// Ensure MerchantClientWrapper implements domain.MerchantClient
var _ domain.MerchantClient = (*MerchantClientWrapper)(nil)

// NewMerchantClientWrapper creates a new wrapper with retry logic
func NewMerchantClientWrapper(client *merchant.MerchantCenterClient, logger zerolog.Logger) *MerchantClientWrapper {
	return &MerchantClientWrapper{
		client:     client,
		logger:     logger,
		maxRetries: 3,
		baseDelay:  1 * time.Second,
	}
}

// InsertProduct inserts a product with retry logic
func (w *MerchantClientWrapper) InsertProduct(ctx context.Context, product *content.Product) error {
	return w.withRetry(ctx, "InsertProduct", func() error {
		return w.client.InsertProduct(ctx, product)
	})
}

// UpdateProduct updates a product with retry logic
func (w *MerchantClientWrapper) UpdateProduct(ctx context.Context, product *content.Product) error {
	return w.withRetry(ctx, "UpdateProduct", func() error {
		return w.client.UpdateProduct(ctx, product)
	})
}

// GetProduct gets a product with retry logic
func (w *MerchantClientWrapper) GetProduct(ctx context.Context, productID string) (*content.Product, error) {
	var result *content.Product
	err := w.withRetry(ctx, "GetProduct", func() error {
		var err error
		result, err = w.client.GetProduct(ctx, productID)
		return err
	})
	return result, err
}

// DeleteProduct deletes a product with retry logic
func (w *MerchantClientWrapper) DeleteProduct(ctx context.Context, productID string) error {
	return w.withRetry(ctx, "DeleteProduct", func() error {
		return w.client.DeleteProduct(ctx, productID)
	})
}

// ListProducts lists products with retry logic
func (w *MerchantClientWrapper) ListProducts(ctx context.Context, pageSize int64, pageToken string) ([]*content.Product, string, error) {
	var products []*content.Product
	var nextToken string
	
	err := w.withRetry(ctx, "ListProducts", func() error {
		var err error
		products, nextToken, err = w.client.ListProducts(ctx, pageSize, pageToken)
		return err
	})
	
	return products, nextToken, err
}

// IsNotFoundErr checks if the error is a not found error
func (w *MerchantClientWrapper) IsNotFoundErr(err error) bool {
	return w.client.IsNotFoundErr(err)
}

// MerchantID returns the merchant ID
func (w *MerchantClientWrapper) MerchantID() uint64 {
	return w.client.MerchantID()
}

// withRetry executes a function with exponential backoff retry logic
func (w *MerchantClientWrapper) withRetry(ctx context.Context, operation string, fn func() error) error {
	var lastErr error
	
	for attempt := 0; attempt <= w.maxRetries; attempt++ {
		// Check context before attempting
		if ctx.Err() != nil {
			return ctx.Err()
		}
		
		startTime := time.Now()
		err := fn()
		duration := time.Since(startTime)
		
		// Log the operation
		logger := w.logger.With().
			Str("operation", operation).
			Int("attempt", attempt+1).
			Dur("duration", duration).
			Logger()
		
		if err == nil {
			if attempt > 0 {
				logger.Info().Msg("operation succeeded after retry")
			} else {
				logger.Debug().Msg("operation succeeded")
			}
			return nil
		}
		
		lastErr = err
		
		// Check if error is retryable
		if !w.isRetryableError(err) {
			logger.Error().Err(err).Msg("non-retryable error")
			return err
		}
		
		// Don't retry if we've exhausted attempts
		if attempt == w.maxRetries {
			logger.Error().Err(err).Msg("max retries exceeded")
			return fmt.Errorf("operation %s failed after %d attempts: %w", operation, w.maxRetries+1, err)
		}
		
		// Calculate delay with exponential backoff
		delay := w.baseDelay * time.Duration(1<<uint(attempt))
		
		logger.Warn().
			Err(err).
			Dur("retry_in", delay).
			Msg("retryable error, will retry")
		
		// Wait before retry
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	
	return lastErr
}

// isRetryableError determines if an error should trigger a retry
func (w *MerchantClientWrapper) isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check for context errors (not retryable)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	
	// Check for Google API errors
	var googleErr *googleapi.Error
	if errors.As(err, &googleErr) {
		switch googleErr.Code {
		case http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return true
		case http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
			http.StatusConflict:
			return false
		default:
			// Retry on other 5xx errors
			return googleErr.Code >= 500 && googleErr.Code < 600
		}
	}
	
	// Default to not retrying unknown errors
	return false
}