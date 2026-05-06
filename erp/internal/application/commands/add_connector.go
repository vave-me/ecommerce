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

// AddConnector command creates a new ERP connector
type AddConnector struct {
	Name        string
	Type        string // odoo, dynamics365, netsuite, sap, erpnext, frappe
	Environment string // production, staging, development, sandbox
	BaseURL     string
	
	// Authentication
	AuthType     string                 // basic, oauth2, apikey, token
	AuthConfig   map[string]interface{} // Credentials (will be encrypted)
	
	// Webhook configuration
	WebhookEnabled bool
	WebhookSecret  string   // Will be encrypted
	WebhookEvents  []string // Event types to subscribe to
	
	// Sync configuration
	SyncEnabled         bool
	SyncIntervalSeconds int
	BatchSize           int
	SyncEntities        []SyncEntityConfig
	
	// Rate limiting
	RateLimitRequestsPerSecond int
	RateLimitBurst             int
	
	// Retry configuration
	RetryMaxAttempts    int
	RetryInitialDelayMs int
	RetryMaxDelayMs     int
	RetryMultiplier     float64
	
	// Additional settings
	CustomHeaders  map[string]string
	TimeoutSeconds int
	
	// Metadata
	CreatedBy string
}

// SyncEntityConfig defines sync configuration for an entity type
type SyncEntityConfig struct {
	EntityType    string                 // product, stock, price, customer, order, invoice, return
	Enabled       bool
	SyncDirection string                 // inbound, outbound, bidirectional
	Filters       map[string]interface{} // Entity-specific filters
	FieldMapping  map[string]string      // Field mapping configuration
}

// AddConnectorHandler handles the AddConnector command
type AddConnectorHandler struct {
	repository domain.ConnectorRepository
	factory    erp.ConnectorFactory
	registry   erp.ConnectorRegistry
	encryptor  *crypto.Encryptor
}

// NewAddConnectorHandler creates a new handler
func NewAddConnectorHandler(
	repository domain.ConnectorRepository,
	factory erp.ConnectorFactory,
	registry erp.ConnectorRegistry,
	encryptor *crypto.Encryptor,
) AddConnectorHandler {
	return AddConnectorHandler{
		repository: repository,
		factory:    factory,
		registry:   registry,
		encryptor:  encryptor,
	}
}

// AddConnector creates and registers a new connector
func (h AddConnectorHandler) AddConnector(ctx context.Context, cmd AddConnector) error {
	// Validate input
	if err := h.validateCommand(cmd); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if connector with same name already exists
	existing, _ := h.repository.GetByName(ctx, cmd.Name)
	if existing != nil {
		return fmt.Errorf("connector with name '%s' already exists", cmd.Name)
	}

	// Generate connector ID
	connectorID := fmt.Sprintf("%s_%s", cmd.Type, uuid.New().String()[:8])

	// Generate salt for encryption
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Encrypt auth config
	authConfigJSON, err := json.Marshal(cmd.AuthConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal auth config: %w", err)
	}

	authConfigEncrypted, err := h.encryptor.EncryptJSON(authConfigJSON, salt)
	if err != nil {
		return fmt.Errorf("failed to encrypt auth config: %w", err)
	}

	// Encrypt webhook secret if provided
	var webhookSecretEncrypted []byte
	if cmd.WebhookSecret != "" {
		webhookSecretEncrypted, err = h.encryptor.EncryptWithSalt([]byte(cmd.WebhookSecret), salt)
		if err != nil {
			return fmt.Errorf("failed to encrypt webhook secret: %w", err)
		}
	}

	// Create connector entity
	connector := &domain.ConnectorEntity{
		ID:                         connectorID,
		Name:                       cmd.Name,
		Type:                       cmd.Type,
		Status:                     "inactive", // Start as inactive until health check passes
		AuthConfigEncrypted:        authConfigEncrypted,
		AuthConfigSalt:             salt,
		BaseURL:                    cmd.BaseURL,
		Environment:                cmd.Environment,
		WebhookEnabled:             cmd.WebhookEnabled,
		WebhookSecretEncrypted:     webhookSecretEncrypted,
		WebhookURL:                 fmt.Sprintf("/api/erp/webhook/%s", connectorID),
		WebhookEvents:              cmd.WebhookEvents,
		SyncEnabled:                cmd.SyncEnabled,
		SyncIntervalSeconds:        cmd.SyncIntervalSeconds,
		BatchSize:                  cmd.BatchSize,
		RateLimitRequestsPerSecond: cmd.RateLimitRequestsPerSecond,
		RateLimitBurst:             cmd.RateLimitBurst,
		RetryMaxAttempts:           cmd.RetryMaxAttempts,
		RetryInitialDelayMs:        cmd.RetryInitialDelayMs,
		RetryMaxDelayMs:            cmd.RetryMaxDelayMs,
		RetryMultiplier:            cmd.RetryMultiplier,
		CustomHeaders:              cmd.CustomHeaders,
		TimeoutSeconds:             cmd.TimeoutSeconds,
		CreatedBy:                  cmd.CreatedBy,
		UpdatedBy:                  cmd.CreatedBy,
		Version:                    1,
	}

	// Set defaults if not provided
	h.setDefaults(connector)

	// Save connector to database
	if err := h.repository.Create(ctx, connector); err != nil {
		return fmt.Errorf("failed to save connector: %w", err)
	}

	// Create sync entity configurations
	for _, syncEntity := range cmd.SyncEntities {
		entity := &domain.ConnectorSyncEntity{
			ID:            uuid.New().String(),
			ConnectorID:   connectorID,
			EntityType:    syncEntity.EntityType,
			Enabled:       syncEntity.Enabled,
			SyncDirection: syncEntity.SyncDirection,
			Filters:       syncEntity.Filters,
			FieldMapping:  syncEntity.FieldMapping,
		}

		if err := h.repository.CreateSyncEntity(ctx, entity); err != nil {
			log.Error().Err(err).
				Str("connector_id", connectorID).
				Str("entity_type", syncEntity.EntityType).
				Msg("failed to create sync entity configuration")
		}
	}

	// Create audit log entry
	auditLog := &domain.ConnectorAuditLog{
		ID:          uuid.New().String(),
		ConnectorID: connectorID,
		Action:      "created",
		ChangedBy:   cmd.CreatedBy,
		ChangedAt:   time.Now(),
		NewValues: map[string]interface{}{
			"name":        cmd.Name,
			"type":        cmd.Type,
			"environment": cmd.Environment,
			"base_url":    cmd.BaseURL,
		},
	}

	if err := h.repository.CreateAuditLog(ctx, auditLog); err != nil {
		log.Error().Err(err).
			Str("connector_id", connectorID).
			Msg("failed to create audit log")
	}

	// Create ERP config from saved connector
	erpConfig := h.buildERPConfig(connector, cmd.AuthConfig, cmd.WebhookSecret)

	// Create connector instance via factory
	connectorInstance, err := h.factory.CreateConnector(*erpConfig)
	if err != nil {
		// Update status to error
		h.repository.UpdateStatus(ctx, connectorID, "error")
		return fmt.Errorf("failed to create connector instance: %w", err)
	}

	// Test connection
	healthCheck := connectorInstance.HealthCheck(ctx)
	if healthCheck.Status != erp.HealthStatusHealthy {
		// Update health check status
		h.repository.UpdateHealthCheck(ctx, connectorID, "failed", healthCheck.Message)
		log.Warn().
			Str("connector_id", connectorID).
			Str("message", healthCheck.Message).
			Msg("connector health check failed")
	} else {
		// Update to active status
		h.repository.UpdateStatus(ctx, connectorID, "active")
		h.repository.UpdateHealthCheck(ctx, connectorID, "healthy", "Connection successful")
	}

	// Register connector in runtime registry
	if err := h.registry.RegisterConnector(connectorID, connectorInstance); err != nil {
		log.Error().Err(err).
			Str("connector_id", connectorID).
			Msg("failed to register connector in runtime registry")
	}

	log.Info().
		Str("connector_id", connectorID).
		Str("name", cmd.Name).
		Str("type", cmd.Type).
		Msg("connector added successfully")

	return nil
}

// validateCommand validates the add connector command
func (h AddConnectorHandler) validateCommand(cmd AddConnector) error {
	if cmd.Name == "" {
		return fmt.Errorf("connector name is required")
	}

	validTypes := map[string]bool{
		"odoo":        true,
		"dynamics365": true,
		"netsuite":    true,
		"sap":         true,
		"erpnext":     true,
		"frappe":      true,
	}

	if !validTypes[cmd.Type] {
		return fmt.Errorf("invalid connector type: %s", cmd.Type)
	}

	if cmd.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}

	if cmd.AuthConfig == nil || len(cmd.AuthConfig) == 0 {
		return fmt.Errorf("authentication configuration is required")
	}

	validEnvironments := map[string]bool{
		"production":  true,
		"staging":     true,
		"development": true,
		"sandbox":     true,
	}

	if cmd.Environment == "" {
		cmd.Environment = "production"
	}

	if !validEnvironments[cmd.Environment] {
		return fmt.Errorf("invalid environment: %s", cmd.Environment)
	}

	return nil
}

// setDefaults sets default values for connector configuration
func (h AddConnectorHandler) setDefaults(connector *domain.ConnectorEntity) {
	if connector.SyncIntervalSeconds == 0 {
		connector.SyncIntervalSeconds = 300 // 5 minutes
	}

	if connector.BatchSize == 0 {
		connector.BatchSize = 100
	}

	if connector.RateLimitRequestsPerSecond == 0 {
		connector.RateLimitRequestsPerSecond = 10
	}

	if connector.RateLimitBurst == 0 {
		connector.RateLimitBurst = 20
	}

	if connector.RetryMaxAttempts == 0 {
		connector.RetryMaxAttempts = 3
	}

	if connector.RetryInitialDelayMs == 0 {
		connector.RetryInitialDelayMs = 1000
	}

	if connector.RetryMaxDelayMs == 0 {
		connector.RetryMaxDelayMs = 60000
	}

	if connector.RetryMultiplier == 0 {
		connector.RetryMultiplier = 2.0
	}

	if connector.TimeoutSeconds == 0 {
		connector.TimeoutSeconds = 30
	}
}

// buildERPConfig creates an ERP config from connector entity
func (h AddConnectorHandler) buildERPConfig(connector *domain.ConnectorEntity, authConfig map[string]interface{}, webhookSecret string) *erp.ERPConfig {
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