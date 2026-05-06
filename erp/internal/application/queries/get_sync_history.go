package queries

import (
	"context"
	"time"

	"middleman/erp/internal/domain"
)

// GetSyncHistory query retrieves sync history for a connector
type GetSyncHistory struct {
	ConnectorID string
	EntityType  string    // Optional: filter by entity type (product, stock, price, order)
	Status      string    // Optional: filter by status (completed, failed, in_progress)
	Since       time.Time // Optional: filter by date
	Limit       int       // Optional: limit results (default 100)
}

// SyncHistoryItem represents a sync history entry
type SyncHistoryItem struct {
	ID               string
	ConnectorID      string
	EntityType       string
	Status           domain.SyncStatus
	StartedAt        time.Time
	CompletedAt      *time.Time
	RecordsProcessed int
	RecordsTotal     int
	Error            string
	Metadata         map[string]interface{}
}

// GetSyncHistoryHandler handles the GetSyncHistory query
type GetSyncHistoryHandler struct {
	repository domain.SyncLogRepository
}

// NewGetSyncHistoryHandler creates a new handler
func NewGetSyncHistoryHandler(repository domain.SyncLogRepository) GetSyncHistoryHandler {
	return GetSyncHistoryHandler{
		repository: repository,
	}
}

// GetSyncHistory retrieves sync history
func (h GetSyncHistoryHandler) GetSyncHistory(ctx context.Context, query GetSyncHistory) ([]SyncHistoryItem, error) {
	// Set default limit
	limit := query.Limit
	if limit == 0 {
		limit = 100
	}

	// Get sync logs
	logs, err := h.repository.GetByConnectorID(ctx, query.ConnectorID)
	if err != nil {
		return nil, err
	}

	// Apply filters and convert to result items
	result := make([]SyncHistoryItem, 0)
	for _, log := range logs {
		// Apply entity type filter
		if query.EntityType != "" && log.EntityType != query.EntityType {
			continue
		}

		// Apply status filter
		if query.Status != "" && string(log.Status) != query.Status {
			continue
		}

		// Apply date filter
		if !query.Since.IsZero() && log.StartedAt.Before(query.Since) {
			continue
		}

		// Convert to result item
		item := SyncHistoryItem{
			ID:               log.ID,
			ConnectorID:      log.ConnectorID,
			EntityType:       log.EntityType,
			Status:           log.Status,
			StartedAt:        log.StartedAt,
			CompletedAt:      log.CompletedAt,
			RecordsProcessed: log.RecordsProcessed,
			RecordsTotal:     log.RecordsTotal,
			Error:            log.Error,
			Metadata:         log.Metadata,
		}

		result = append(result, item)

		// Apply limit
		if len(result) >= limit {
			break
		}
	}

	return result, nil
}