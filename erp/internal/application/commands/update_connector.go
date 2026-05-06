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

// UpdateConnector command updates an existing ERP connector
type UpdateConnector struct {
	ConnectorID string
	Name        string
	Environment string
	BaseURL     string
	
	// Authentication - only update if provided
	AuthConfig *map[string]interface{}
	
	// Webhook configuration
	WebhookEnabled *bool
	WebhookSecret  *string // Will be encrypted if provided
	WebhookEvents  []string
	
	// Sync configuration
	SyncEnabled         *bool
	SyncIntervalSeconds *int
	BatchSize           *int
	
	// Rate limiting
	RateLimitRequestsPerSecond *int
	RateLimitBurst             *int
	
	// Retry configuration
	RetryMaxAttempts    *int
	RetryInitialDelayMs *int
	RetryMaxDelayMs     *int
	RetryMultiplier     *float64
	
	// Additional settings
	CustomHeaders  map[string]string
	TimeoutSeconds *int
	
	// Metadata
	UpdatedBy string
}

// UpdateConnectorHandler handles the UpdateConnector command
type UpdateConnectorHandler struct {
	repository domain.ConnectorRepository
	factory    erp.ConnectorFactory
	registry   erp.ConnectorRegistry
	encryptor  *crypto.Encryptor
}

// NewUpdateConnectorHandler creates a new handler
func NewUpdateConnectorHandler(
	repository domain.ConnectorRepository,
	factory erp.ConnectorFactory,
	registry erp.ConnectorRegistry,
	encryptor *crypto.Encryptor,
) UpdateConnectorHandler {
	return UpdateConnectorHandler{
		repository: repository,
		factory:    factory,
		registry:   registry,
		encryptor:  encryptor,
	}
}

// UpdateConnector updates an existing connector
func (h UpdateConnectorHandler) UpdateConnector(ctx context.Context, cmd UpdateConnector) error {
	// Get existing connector
	connector, err := h.repository.GetByID(ctx, cmd.ConnectorID)
	if err != nil {
		return fmt.Errorf("connector not found: %w", err)
	}

	// Store old values for audit log
	oldValues := map[string]interface{}{
		"name":        connector.Name,
		"environment": connector.Environment,
		"base_url":    connector.BaseURL,
		"status":      connector.Status,
	}

	// Update fields if provided
	if cmd.Name != "" && cmd.Name != connector.Name {
		// Check if new name already exists
		existing, _ := h.repository.GetByName(ctx, cmd.Name)
		if existing != nil && existing.ID != connector.ID {
			return fmt.Errorf("connector with name '%s' already exists", cmd.Name)
		}
		connector.Name = cmd.Name
	}

	if cmd.Environment != "" {
		connector.Environment = cmd.Environment
	}

	if cmd.BaseURL != "" {
		connector.BaseURL = cmd.BaseURL
	}

	// Update auth config if provided
	var authConfig map[string]interface{}
	if cmd.AuthConfig != nil {
		authConfig = *cmd.AuthConfig
		
		// Generate new salt for re-encryption
		salt, err := crypto.GenerateSalt()
		if err != nil {
			return fmt.Errorf("failed to generate salt: %w", err)
		}

		// Encrypt new auth config
		authConfigJSON, err := json.Marshal(authConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal auth config: %w", err)
		}

		authConfigEncrypted, err := h.encryptor.EncryptJSON(authConfigJSON, salt)
		if err != nil {
			return fmt.Errorf("failed to encrypt auth config: %w", err)
		}

		connector.AuthConfigEncrypted = authConfigEncrypted
		connector.AuthConfigSalt = salt
	} else {
		// Decrypt existing auth config for connector recreation
		authConfigJSON, err := h.encryptor.DecryptJSON(connector.AuthConfigEncrypted, connector.AuthConfigSalt)
		if err != nil {
			return fmt.Errorf("failed to decrypt existing auth config: %w", err)
		}

		if err := json.Unmarshal(authConfigJSON, &authConfig); err != nil {
			return fmt.Errorf("failed to unmarshal auth config: %w", err)
		}
	}

	// Update webhook configuration
	var webhookSecret string
	if cmd.WebhookEnabled != nil {
		connector.WebhookEnabled = *cmd.WebhookEnabled
	}

	if cmd.WebhookSecret != nil {
		webhookSecret = *cmd.WebhookSecret
		webhookSecretEncrypted, err := h.encryptor.EncryptWithSalt([]byte(webhookSecret), connector.AuthConfigSalt)
		if err != nil {
			return fmt.Errorf("failed to encrypt webhook secret: %w", err)
		}
		connector.WebhookSecretEncrypted = webhookSecretEncrypted
	} else if len(connector.WebhookSecretEncrypted) > 0 {
		// Decrypt existing webhook secret
		decrypted, err := h.encryptor.DecryptWithSalt(connector.WebhookSecretEncrypted, connector.AuthConfigSalt)
		if err != nil {
			log.Warn().Err(err).Msg("failed to decrypt webhook secret")
		} else {
			webhookSecret = string(decrypted)
		}
	}

	if len(cmd.WebhookEvents) > 0 {
		connector.WebhookEvents = cmd.WebhookEvents
	}

	// Update sync configuration
	if cmd.SyncEnabled != nil {
		connector.SyncEnabled = *cmd.SyncEnabled
	}

	if cmd.SyncIntervalSeconds != nil {
		connector.SyncIntervalSeconds = *cmd.SyncIntervalSeconds
	}

	if cmd.BatchSize != nil {
		connector.BatchSize = *cmd.BatchSize
	}

	// Update rate limiting
	if cmd.RateLimitRequestsPerSecond != nil {
		connector.RateLimitRequestsPerSecond = *cmd.RateLimitRequestsPerSecond
	}

	if cmd.RateLimitBurst != nil {
		connector.RateLimitBurst = *cmd.RateLimitBurst
	}

	// Update retry configuration
	if cmd.RetryMaxAttempts != nil {
		connector.RetryMaxAttempts = *cmd.RetryMaxAttempts
	}

	if cmd.RetryInitialDelayMs != nil {
		connector.RetryInitialDelayMs = *cmd.RetryInitialDelayMs
	}

	if cmd.RetryMaxDelayMs != nil {
		connector.RetryMaxDelayMs = *cmd.RetryMaxDelayMs
	}

	if cmd.RetryMultiplier != nil {
		connector.RetryMultiplier = *cmd.RetryMultiplier
	}

	// Update additional settings
	if cmd.CustomHeaders != nil {
		connector.CustomHeaders = cmd.CustomHeaders
	}

	if cmd.TimeoutSeconds != nil {
		connector.TimeoutSeconds = *cmd.TimeoutSeconds
	}

	connector.UpdatedBy = cmd.UpdatedBy

	// Update in database
	if err := h.repository.Update(ctx, connector); err != nil {
		return fmt.Errorf("failed to update connector: %w", err)
	}

	// Create audit log entry
	newValues := map[string]interface{}{
		"name":        connector.Name,
		"environment": connector.Environment,
		"base_url":    connector.BaseURL,
		"status":      connector.Status,
	}

	auditLog := &domain.ConnectorAuditLog{
		ID:          uuid.New().String(),
		ConnectorID: connector.ID,
		Action:      "updated",
		ChangedBy:   cmd.UpdatedBy,
		ChangedAt:   time.Now(),
		OldValues:   oldValues,
		NewValues:   newValues,
	}

	if err := h.repository.CreateAuditLog(ctx, auditLog); err != nil {
		log.Error().Err(err).
			Str("connector_id", connector.ID).
			Msg("failed to create audit log")
	}

	// Recreate connector instance with updated config
	erpConfig := h.buildERPConfig(connector, authConfig, webhookSecret)

	// Remove old instance from registry
	h.registry.RemoveConnector(connector.ID)

	// Create new connector instance
	connectorInstance, err := h.factory.CreateConnector(*erpConfig)
	if err != nil {
		// Update status to error
		h.repository.UpdateStatus(ctx, connector.ID, "error")
		return fmt.Errorf("failed to recreate connector instance: %w", err)
	}

	// Test connection
	healthCheck := connectorInstance.HealthCheck(ctx)
	if healthCheck.Status != erp.HealthStatusHealthy {
		// Update health check status
		h.repository.UpdateHealthCheck(ctx, connector.ID, "failed", healthCheck.Message)
		h.repository.UpdateStatus(ctx, connector.ID, "error")
		log.Warn().
			Str("connector_id", connector.ID).
			Str("message", healthCheck.Message).
			Msg("connector health check failed after update")
	} else {
		h.repository.UpdateHealthCheck(ctx, connector.ID, "healthy", "")
		// Only set to active if it wasn't already in maintenance
		if connector.Status != "maintenance" {
			h.repository.UpdateStatus(ctx, connector.ID, "active")
		}
	}

	// Register updated connector in runtime registry
	if err := h.registry.RegisterConnector(connector.ID, connectorInstance); err != nil {
		log.Error().Err(err).
			Str("connector_id", connector.ID).
			Msg("failed to register updated connector in runtime registry")
	}

	log.Info().
		Str("connector_id", connector.ID).
		Str("name", connector.Name).
		Msg("connector updated successfully")

	return nil
}

// buildERPConfig creates an ERP config from connector entity
func (h UpdateConnectorHandler) buildERPConfig(connector *domain.ConnectorEntity, authConfig map[string]interface{}, webhookSecret string) *erp.ERPConfig {
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