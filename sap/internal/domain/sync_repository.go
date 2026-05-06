package domain

import (
	"context"
	"time"
)

// SyncStatusRepository manages sync status persistence
type SyncStatusRepository interface {
	// GetByEntityID retrieves sync status for a specific entity
	GetByEntityID(ctx context.Context, entityType, entityID string) (*SyncStatus, error)
	
	// Save creates or updates sync status
	Save(ctx context.Context, status *SyncStatus) error
	
	// UpdateStatus updates the sync status and error message
	UpdateStatus(ctx context.Context, id, status, errorMessage string) error
	
	// GetFailedSyncs retrieves all failed sync records for retry
	GetFailedSyncs(ctx context.Context, entityType string, maxRetries int) ([]*SyncStatus, error)
	
	// GetLastSyncTime gets the last successful sync time for an entity type
	GetLastSyncTime(ctx context.Context, entityType string) (time.Time, error)
}

// SyncLogRepository manages sync log persistence
type SyncLogRepository interface {
	// Create creates a new sync log entry
	Create(ctx context.Context, log *SyncLog) error
	
	// GetByEventID retrieves logs for a specific event
	GetByEventID(ctx context.Context, eventID string) ([]*SyncLog, error)
	
	// GetRecentLogs retrieves recent sync logs
	GetRecentLogs(ctx context.Context, limit int) ([]*SyncLog, error)
	
	// GetByDateRange retrieves logs within a date range
	GetByDateRange(ctx context.Context, from, to time.Time) ([]*SyncLog, error)
}

// SyncConfigurationRepository manages sync configuration persistence
type SyncConfigurationRepository interface {
	// GetByEntityType retrieves configuration for a specific entity type
	GetByEntityType(ctx context.Context, entityType string) (*SyncConfiguration, error)
	
	// GetEnabled retrieves all enabled sync configurations
	GetEnabled(ctx context.Context) ([]*SyncConfiguration, error)
	
	// Save creates or updates sync configuration
	Save(ctx context.Context, config *SyncConfiguration) error
	
	// UpdateLastExecuted updates the last execution time
	UpdateLastExecuted(ctx context.Context, id string, executedAt time.Time) error
	
	// GetDueForExecution retrieves configurations due for execution
	GetDueForExecution(ctx context.Context) ([]*SyncConfiguration, error)
}

// WebhookEventRepository manages webhook event persistence
type WebhookEventRepository interface {
	// Create creates a new webhook event record
	Create(ctx context.Context, event *WebhookEvent) error
	
	// GetByID retrieves a webhook event by ID
	GetByID(ctx context.Context, id string) (*WebhookEvent, error)
	
	// UpdateStatus updates the processing status of a webhook event
	UpdateStatus(ctx context.Context, id, status, errorMessage string) error
	
	// GetUnprocessed retrieves unprocessed webhook events
	GetUnprocessed(ctx context.Context, limit int) ([]*WebhookEvent, error)
	
	// GetByEventID retrieves a webhook event by its external event ID
	GetByEventID(ctx context.Context, eventID string) (*WebhookEvent, error)
}