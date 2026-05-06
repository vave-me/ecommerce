package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"middleman/erp/internal/domain"
	"middleman/internal/erp"
)

// RemoveConnector command removes an ERP connector
type RemoveConnector struct {
	ConnectorID string
	RemoveBy    string
	Force       bool // Force removal even if there are active syncs
}

// RemoveConnectorHandler handles the RemoveConnector command
type RemoveConnectorHandler struct {
	repository       domain.ConnectorRepository
	registry         erp.ConnectorRegistry
	syncLogRepo      domain.SyncLogRepository
	orderSyncRepo    domain.OrderSyncRepository
	invoiceSyncRepo  domain.InvoiceSyncRepository
}

// NewRemoveConnectorHandler creates a new handler
func NewRemoveConnectorHandler(
	repository domain.ConnectorRepository,
	registry erp.ConnectorRegistry,
	syncLogRepo domain.SyncLogRepository,
	orderSyncRepo domain.OrderSyncRepository,
	invoiceSyncRepo domain.InvoiceSyncRepository,
) RemoveConnectorHandler {
	return RemoveConnectorHandler{
		repository:      repository,
		registry:        registry,
		syncLogRepo:     syncLogRepo,
		orderSyncRepo:   orderSyncRepo,
		invoiceSyncRepo: invoiceSyncRepo,
	}
}

// RemoveConnector removes a connector and all associated data
func (h RemoveConnectorHandler) RemoveConnector(ctx context.Context, cmd RemoveConnector) error {
	// Get connector
	connector, err := h.repository.GetByID(ctx, cmd.ConnectorID)
	if err != nil {
		return fmt.Errorf("connector not found: %w", err)
	}

	// Check for active syncs if not forcing
	if !cmd.Force {
		// TODO: Implement these repository methods when needed
		// - GetByConnectorAndStatus for sync logs
		// - GetPendingByConnector for order and invoice syncs
		
		// For now, we'll skip these checks
		log.Warn().
			Str("connector_id", cmd.ConnectorID).
			Msg("skipping active sync checks - methods not yet implemented")
	}

	// Store connector info for audit log
	oldValues := map[string]interface{}{
		"name":        connector.Name,
		"type":        connector.Type,
		"environment": connector.Environment,
		"base_url":    connector.BaseURL,
		"status":      connector.Status,
	}

	// Remove from runtime registry first
	h.registry.RemoveConnector(cmd.ConnectorID)

	// Delete API keys (cascade will handle this, but doing it explicitly for cleanup)
	apiKeys, err := h.repository.GetAPIKeys(ctx, cmd.ConnectorID)
	if err == nil {
		for _, key := range apiKeys {
			if err := h.repository.DeleteAPIKey(ctx, key.ID); err != nil {
				log.Warn().Err(err).
					Str("key_id", key.ID).
					Msg("failed to delete API key")
			}
		}
	}

	// Delete sync entities (cascade will handle this too)
	syncEntities, err := h.repository.GetSyncEntities(ctx, cmd.ConnectorID)
	if err == nil {
		for _, entity := range syncEntities {
			if err := h.repository.DeleteSyncEntity(ctx, entity.ID); err != nil {
				log.Warn().Err(err).
					Str("entity_id", entity.ID).
					Msg("failed to delete sync entity")
			}
		}
	}

	// Create audit log entry before deletion
	auditLog := &domain.ConnectorAuditLog{
		ID:          uuid.New().String(),
		ConnectorID: cmd.ConnectorID,
		Action:      "deleted",
		ChangedBy:   cmd.RemoveBy,
		ChangedAt:   time.Now(),
		OldValues:   oldValues,
		NewValues:   map[string]interface{}{},
	}

	if err := h.repository.CreateAuditLog(ctx, auditLog); err != nil {
		log.Error().Err(err).
			Str("connector_id", cmd.ConnectorID).
			Msg("failed to create deletion audit log")
	}

	// Delete the connector
	if err := h.repository.Delete(ctx, cmd.ConnectorID); err != nil {
		return fmt.Errorf("failed to delete connector: %w", err)
	}

	// Clean up orphaned sync logs (optional, depending on requirements)
	if cmd.Force {
		// TODO: Implement FailInProgressByConnector method when needed
		log.Debug().
			Str("connector_id", cmd.ConnectorID).
			Msg("skipping in-progress sync cleanup - method not yet implemented")

		// TODO: Implement CancelPendingByConnector for orders when needed

		// TODO: Implement CancelPendingByConnector for invoices when needed
	}

	log.Info().
		Str("connector_id", cmd.ConnectorID).
		Str("name", connector.Name).
		Str("type", connector.Type).
		Bool("forced", cmd.Force).
		Msg("connector removed successfully")

	return nil
}