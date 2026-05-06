package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"middleman/erp/internal/crypto"
	"middleman/erp/internal/domain"
	"middleman/internal/erp"
)

// ToggleConnector command activates or deactivates a connector
type ToggleConnector struct {
	ConnectorID string
	Activate    bool   // true to activate, false to deactivate
	ChangedBy   string
	Reason      string // Optional reason for the change
}

// ToggleConnectorHandler handles the ToggleConnector command
type ToggleConnectorHandler struct {
	repository domain.ConnectorRepository
	registry   erp.ConnectorRegistry
	factory    erp.ConnectorFactory
	encryptor  *crypto.Encryptor
}

// NewToggleConnectorHandler creates a new handler
func NewToggleConnectorHandler(
	repository domain.ConnectorRepository,
	registry erp.ConnectorRegistry,
	factory erp.ConnectorFactory,
	encryptor *crypto.Encryptor,
) ToggleConnectorHandler {
	return ToggleConnectorHandler{
		repository: repository,
		registry:   registry,
		factory:    factory,
		encryptor:  encryptor,
	}
}

// ToggleConnector activates or deactivates a connector
func (h ToggleConnectorHandler) ToggleConnector(ctx context.Context, cmd ToggleConnector) error {
	// Get connector
	connector, err := h.repository.GetByID(ctx, cmd.ConnectorID)
	if err != nil {
		return fmt.Errorf("connector not found: %w", err)
	}

	// Check current status
	currentStatus := connector.Status
	targetStatus := "inactive"
	action := "deactivated"

	if cmd.Activate {
		targetStatus = "active"
		action = "activated"

		// If activating, check if connector is healthy
		if currentStatus == "error" {
			// Try to create and test the connector
			authConfig, webhookSecret, err := h.decryptCredentials(connector)
			if err != nil {
				return fmt.Errorf("failed to decrypt credentials: %w", err)
			}

			erpConfig := h.buildERPConfig(connector, authConfig, webhookSecret)

			// Create connector instance
			connectorInstance, err := h.factory.CreateConnector(*erpConfig)
			if err != nil {
				return fmt.Errorf("cannot activate connector: failed to create instance: %w", err)
			}

			// Test connection
			healthCheck := connectorInstance.HealthCheck(ctx)
			if healthCheck.Status != erp.HealthStatusHealthy {
				h.repository.UpdateHealthCheck(ctx, connector.ID, "failed", healthCheck.Message)
				return fmt.Errorf("cannot activate connector: health check failed: %s", healthCheck.Message)
			}

			// Register in runtime registry
			if err := h.registry.RegisterConnector(connector.ID, connectorInstance); err != nil {
				return fmt.Errorf("failed to register connector: %w", err)
			}

			h.repository.UpdateHealthCheck(ctx, connector.ID, "healthy", "")
		}
	} else {
		// If deactivating, remove from runtime registry
		h.registry.RemoveConnector(connector.ID)
	}

	// Update status
	if err := h.repository.UpdateStatus(ctx, cmd.ConnectorID, targetStatus); err != nil {
		return fmt.Errorf("failed to update connector status: %w", err)
	}

	// Create audit log
	auditLog := &domain.ConnectorAuditLog{
		ID:          uuid.New().String(),
		ConnectorID: cmd.ConnectorID,
		Action:      action,
		ChangedBy:   cmd.ChangedBy,
		ChangedAt:   time.Now(),
		OldValues: map[string]interface{}{
			"status": currentStatus,
		},
		NewValues: map[string]interface{}{
			"status": targetStatus,
			"reason": cmd.Reason,
		},
	}

	if err := h.repository.CreateAuditLog(ctx, auditLog); err != nil {
		log.Error().Err(err).
			Str("connector_id", cmd.ConnectorID).
			Msg("failed to create audit log")
	}

	log.Info().
		Str("connector_id", cmd.ConnectorID).
		Str("name", connector.Name).
		Str("old_status", currentStatus).
		Str("new_status", targetStatus).
		Str("action", action).
		Str("reason", cmd.Reason).
		Msg("connector status changed")

	return nil
}

// decryptCredentials decrypts the auth config and webhook secret
func (h ToggleConnectorHandler) decryptCredentials(connector *domain.ConnectorEntity) (map[string]interface{}, string, error) {
	// Decrypt auth config
	authConfigJSON, err := h.encryptor.DecryptJSON(connector.AuthConfigEncrypted, connector.AuthConfigSalt)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decrypt auth config: %w", err)
	}

	var authConfig map[string]interface{}
	if err := json.Unmarshal(authConfigJSON, &authConfig); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal auth config: %w", err)
	}

	// Decrypt webhook secret if present
	var webhookSecret string
	if len(connector.WebhookSecretEncrypted) > 0 {
		decrypted, err := h.encryptor.DecryptWithSalt(connector.WebhookSecretEncrypted, connector.AuthConfigSalt)
		if err != nil {
			log.Warn().Err(err).Msg("failed to decrypt webhook secret")
		} else {
			webhookSecret = string(decrypted)
		}
	}

	return authConfig, webhookSecret, nil
}

// buildERPConfig creates an ERP config from connector entity
func (h ToggleConnectorHandler) buildERPConfig(connector *domain.ConnectorEntity, authConfig map[string]interface{}, webhookSecret string) *erp.ERPConfig {
	// Build auth config
	auth := erp.AuthConfig{
		Type: authConfig["type"].(string),
	}
	
	// Map auth fields based on type
	if v, ok := authConfig["client_id"].(string); ok {
		auth.ClientID = v
	}
	if v, ok := authConfig["client_secret"].(string); ok {
		auth.ClientSecret = v
	}
	if v, ok := authConfig["api_key"].(string); ok {
		auth.APIKey = v
	}
	if v, ok := authConfig["username"].(string); ok {
		auth.Username = v
	}
	if v, ok := authConfig["password"].(string); ok {
		auth.Password = v
	}
	if v, ok := authConfig["token_url"].(string); ok {
		auth.TokenURL = v
	}
	// OAuth 1.0a fields for NetSuite
	if v, ok := authConfig["consumer_key"].(string); ok {
		auth.ConsumerKey = v
	}
	if v, ok := authConfig["consumer_secret"].(string); ok {
		auth.ConsumerSecret = v
	}
	if v, ok := authConfig["token_id"].(string); ok {
		auth.TokenID = v
	}
	if v, ok := authConfig["token_secret"].(string); ok {
		auth.TokenSecret = v
	}

	return &erp.ERPConfig{
		ID:          connector.ID,
		Name:        connector.Name,
		Type:        erp.ERPType(connector.Type),
		Endpoint:    connector.BaseURL,
		URL:         connector.BaseURL,
		Auth:        auth,
		Webhook: erp.WebhookConfig{
			Enabled:      connector.WebhookEnabled,
			Secret:       webhookSecret,
			ValidateSign: true,
			URL:          connector.WebhookURL,
			Events:       connector.WebhookEvents,
		},
		Sync: erp.SyncConfig{
			Enabled:   connector.SyncEnabled,
			Interval:  time.Duration(connector.SyncIntervalSeconds) * time.Second,
			BatchSize: connector.BatchSize,
		},
		RateLimit: &erp.RateLimitConfig{
			RequestsPerSecond: connector.RateLimitRequestsPerSecond,
			BurstSize:         connector.RateLimitBurst,
		},
		Retry: &erp.RetryConfig{
			MaxAttempts:  connector.RetryMaxAttempts,
			InitialDelay: time.Duration(connector.RetryInitialDelayMs) * time.Millisecond,
			MaxDelay:     time.Duration(connector.RetryMaxDelayMs) * time.Millisecond,
			Multiplier:   connector.RetryMultiplier,
		},
		Metadata: authConfig,
	}
}