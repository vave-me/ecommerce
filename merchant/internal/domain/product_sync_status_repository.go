package domain

import (
	"context"
	"time"
)

// ProductSyncStatusRepository defines the interface for product sync status persistence
type ProductSyncStatusRepository interface {
	// Create creates a new sync status record
	Create(ctx context.Context, status *ProductSyncStatus) error
	
	// Update updates an existing sync status record
	Update(ctx context.Context, status *ProductSyncStatus) error
	
	// FindByProductID finds sync status by product ID
	FindByProductID(ctx context.Context, productID string) (*ProductSyncStatus, error)
	
	// FindByStatus finds all products with a specific sync status
	FindByStatus(ctx context.Context, status string, limit int, offset int) ([]*ProductSyncStatus, error)
	
	// FindFailedSyncs finds products that failed to sync
	FindFailedSyncs(ctx context.Context, since time.Time, limit int) ([]*ProductSyncStatus, error)
	
	// DeleteByProductID deletes sync status for a product
	DeleteByProductID(ctx context.Context, productID string) error
	
	// GetSyncStats gets sync statistics
	GetSyncStats(ctx context.Context) (*SyncStats, error)
}

// SyncStats represents synchronization statistics
type SyncStats struct {
	TotalProducts   int
	SyncedProducts  int
	PendingProducts int
	FailedProducts  int
	LastSyncTime    *time.Time
}