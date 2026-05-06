package queries

import (
	"context"
	"fmt"
	"middleman/internal/erp"
	"time"

	"middleman/erp/internal/domain"
)

// GetConnectorStatus query retrieves the status of a connector
type GetConnectorStatus struct {
	ConnectorID string
}

// ConnectorStatus represents the status of a connector
type ConnectorStatus struct {
	ConnectorID string
	Type        erp.ERPType
	Status      erp.HealthStatus
	Message     string
	Details     map[string]interface{}
	LastSync    *SyncSummary
	WebhookInfo *WebhookInfo
}

// SyncSummary contains summary of last sync operations
type SyncSummary struct {
	LastProductSync *time.Time
	LastStockSync   *time.Time
	LastPriceSync   *time.Time
	LastOrderSync   *time.Time
	TotalSyncs      int
	FailedSyncs     int
}

// WebhookInfo contains webhook configuration info
type WebhookInfo struct {
	Enabled        bool
	URL            string
	LastWebhookAt  *time.Time
	TotalWebhooks  int
	FailedWebhooks int
}

// GetConnectorStatusHandler handles the GetConnectorStatus query
type GetConnectorStatusHandler struct {
	registry         erp.ConnectorRegistry
	syncLogRepo      domain.SyncLogRepository
	webhookEventRepo domain.WebhookEventRepository
}

// NewGetConnectorStatusHandler creates a new handler
func NewGetConnectorStatusHandler(
	registry erp.ConnectorRegistry,
	syncLogRepo domain.SyncLogRepository,
	webhookEventRepo domain.WebhookEventRepository,
) GetConnectorStatusHandler {
	return GetConnectorStatusHandler{
		registry:         registry,
		syncLogRepo:      syncLogRepo,
		webhookEventRepo: webhookEventRepo,
	}
}

// GetConnectorStatus retrieves connector status
func (h GetConnectorStatusHandler) GetConnectorStatus(ctx context.Context, query GetConnectorStatus) (*ConnectorStatus, error) {
	// Get connector
	connector, err := h.registry.GetConnector(query.ConnectorID)
	if err != nil {
		return nil, fmt.Errorf("getting connector: %w", err)
	}

	// Get health status
	health := connector.HealthCheck(ctx)
	if health.Status == erp.HealthStatusUnhealthy {
		// Connector exists but is unhealthy
		return &ConnectorStatus{
			ConnectorID: query.ConnectorID,
			Type:        erp.ERPType(connector.GetType()),
			Status:      erp.HealthStatusUnhealthy,
			Message:     health.Message,
		}, nil
	}

	// Get sync summary
	syncSummary, err := h.getSyncSummary(ctx, query.ConnectorID)
	if err != nil {
		// Log error but don't fail the query
		syncSummary = &SyncSummary{}
	}

	// Get webhook info
	webhookInfo, err := h.getWebhookInfo(ctx, query.ConnectorID)
	if err != nil {
		// Log error but don't fail the query
		webhookInfo = &WebhookInfo{}
	}

	return &ConnectorStatus{
		ConnectorID: query.ConnectorID,
		Type:        erp.ERPType(connector.GetType()),
		Status:      health.Status,
		Message:     health.Message,
		Details:     map[string]interface{}{}, // Default empty details
		LastSync:    syncSummary,
		WebhookInfo: webhookInfo,
	}, nil
}

func (h GetConnectorStatusHandler) getSyncSummary(ctx context.Context, connectorID string) (*SyncSummary, error) {
	// Get sync logs for the connector
	logs, err := h.syncLogRepo.GetByConnectorID(ctx, connectorID)
	if err != nil {
		return nil, err
	}

	summary := &SyncSummary{
		TotalSyncs: len(logs),
	}

	for _, log := range logs {
		if log.Status == domain.SyncStatusFailed {
			summary.FailedSyncs++
		}

		if log.Status == domain.SyncStatusCompleted && log.CompletedAt != nil {
			switch log.EntityType {
			case "product":
				if summary.LastProductSync == nil || log.CompletedAt.After(*summary.LastProductSync) {
					summary.LastProductSync = log.CompletedAt
				}
			case "stock":
				if summary.LastStockSync == nil || log.CompletedAt.After(*summary.LastStockSync) {
					summary.LastStockSync = log.CompletedAt
				}
			case "price":
				if summary.LastPriceSync == nil || log.CompletedAt.After(*summary.LastPriceSync) {
					summary.LastPriceSync = log.CompletedAt
				}
			case "order":
				if summary.LastOrderSync == nil || log.CompletedAt.After(*summary.LastOrderSync) {
					summary.LastOrderSync = log.CompletedAt
				}
			}
		}
	}

	return summary, nil
}

func (h GetConnectorStatusHandler) getWebhookInfo(ctx context.Context, connectorID string) (*WebhookInfo, error) {
	// Get webhook configuration from connector
	connector, _ := h.registry.GetConnector(connectorID)
	config := connector.GetConfig()

	info := &WebhookInfo{
		Enabled: config.Webhook.Enabled,
		URL:     config.Webhook.URL, // Webhook URL from config
	}

	// Get webhook events for the connector
	events, err := h.webhookEventRepo.GetByConnectorID(ctx, connectorID)
	if err != nil {
		return info, nil // Return partial info
	}

	info.TotalWebhooks = len(events)

	for _, event := range events {
		if event.Status == domain.WebhookStatusFailed {
			info.FailedWebhooks++
		}

		if info.LastWebhookAt == nil || event.ReceivedAt.After(*info.LastWebhookAt) {
			info.LastWebhookAt = &event.ReceivedAt
		}
	}

	return info, nil
}
