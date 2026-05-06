package domain

import (
	"errors"
	"fmt"
)

// Common domain errors
var (
	// ErrProductNotFound is returned when a product is not found
	ErrProductNotFound = errors.New("product not found")
	
	// ErrSyncStatusNotFound is returned when sync status is not found
	ErrSyncStatusNotFound = errors.New("sync status not found")
	
	// ErrInvalidProduct is returned when product data is invalid
	ErrInvalidProduct = errors.New("invalid product")
	
	// ErrSyncInProgress is returned when a sync is already in progress
	ErrSyncInProgress = errors.New("sync already in progress")
	
	// ErrMerchantAPIUnavailable is returned when Google Merchant Center API is unavailable
	ErrMerchantAPIUnavailable = errors.New("merchant center API unavailable")
	
	// ErrQuotaExceeded is returned when API quota is exceeded
	ErrQuotaExceeded = errors.New("API quota exceeded")
)

// SyncError represents an error during synchronization
type SyncError struct {
	ProductID string
	Operation string
	Err       error
}

func (e *SyncError) Error() string {
	return fmt.Sprintf("sync error for product %s during %s: %v", e.ProductID, e.Operation, e.Err)
}

func (e *SyncError) Unwrap() error {
	return e.Err
}

// BatchSyncError represents errors during batch synchronization
type BatchSyncError struct {
	Errors       []error
	SuccessCount int
	FailedCount  int
}

func (e *BatchSyncError) Error() string {
	return fmt.Sprintf("batch sync completed with %d successes and %d failures", e.SuccessCount, e.FailedCount)
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}