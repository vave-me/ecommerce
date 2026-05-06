package domain

import (
	"context"
	"time"
)

// SyncStatusRecord represents the sync status for an entity type
type SyncStatusRecord struct {
	ID           string    `json:"id" db:"id"`
	ERPType      string    `json:"erpType" db:"erp_type"`
	EntityType   string    `json:"entityType" db:"entity_type"`
	EntityID     string    `json:"entityId,omitempty" db:"entity_id"`
	LastSyncedAt time.Time `json:"lastSyncedAt" db:"last_synced_at"`
	Status       string    `json:"status" db:"status"`
	ErrorMessage string    `json:"errorMessage,omitempty" db:"error_message"`
	RetryCount   int       `json:"retryCount" db:"retry_count"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// SyncConfiguration represents sync configuration for an entity
type SyncConfiguration struct {
	ID            string                 `json:"id" db:"id"`
	ERPType       string                 `json:"erpType" db:"erp_type"`
	EntityType    string                 `json:"entityType" db:"entity_type"`
	Enabled       bool                   `json:"enabled" db:"enabled"`
	SyncInterval  time.Duration          `json:"syncInterval" db:"sync_interval"`
	BatchSize     int                    `json:"batchSize" db:"batch_size"`
	RetryAttempts int                    `json:"retryAttempts" db:"retry_attempts"`
	RetryDelay    time.Duration          `json:"retryDelay" db:"retry_delay"`
	Filters       map[string]interface{} `json:"filters,omitempty" db:"filters"`
	CreatedAt     time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time              `json:"updatedAt" db:"updated_at"`
}

// Repositories

// SyncStatusRepository defines sync status operations
type SyncStatusRepository interface {
	GetLastSyncTime(ctx context.Context, erpType, entityType string) (time.Time, error)
	UpdateLastSync(ctx context.Context, erpType, entityType string, syncTime time.Time) error
	Save(ctx context.Context, status *SyncStatusRecord) error
	GetByEntity(ctx context.Context, erpType, entityType, entityID string) (*SyncStatusRecord, error)
	GetFailedSyncs(ctx context.Context, limit int) ([]*SyncStatusRecord, error)
}

// SyncConfigurationRepository defines sync configuration operations
type SyncConfigurationRepository interface {
	GetByEntity(ctx context.Context, erpType, entityType string) (*SyncConfiguration, error)
	Save(ctx context.Context, config *SyncConfiguration) error
	GetEnabled(ctx context.Context) ([]*SyncConfiguration, error)
	Delete(ctx context.Context, id string) error
}
