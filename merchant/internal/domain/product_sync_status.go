package domain

import "time"

// ProductSyncStatus represents the synchronization status of a product with Google Merchant Center
type ProductSyncStatus struct {
	ProductID      string
	MerchantID     string
	SyncStatus     string
	LastSyncedAt   *time.Time
	LastError      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// SyncStatus constants
const (
	SyncStatusPending = "PENDING"
	SyncStatusSynced  = "SYNCED"
	SyncStatusFailed  = "FAILED"
	SyncStatusRemoved = "REMOVED"
)