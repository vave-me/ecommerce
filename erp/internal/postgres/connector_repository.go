package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
	"middleman/erp/internal/domain"
)

type ConnectorRepository struct {
	db *sql.DB
}

func NewConnectorRepository(db *sql.DB) *ConnectorRepository {
	return &ConnectorRepository{db: db}
}

// Create inserts a new connector
func (r *ConnectorRepository) Create(ctx context.Context, connector *domain.ConnectorEntity) error {
	query := `
		INSERT INTO connectors (
			id, name, type, status, auth_config_encrypted, auth_config_salt,
			base_url, environment, webhook_enabled, webhook_secret_encrypted,
			webhook_url, webhook_events, sync_enabled, sync_interval_seconds,
			batch_size, rate_limit_requests_per_second, rate_limit_burst,
			retry_max_attempts, retry_initial_delay_ms, retry_max_delay_ms,
			retry_multiplier, custom_headers, timeout_seconds,
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25
		)`

	customHeaders, err := json.Marshal(connector.CustomHeaders)
	if err != nil {
		return fmt.Errorf("failed to marshal custom headers: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		connector.ID, connector.Name, connector.Type, connector.Status,
		connector.AuthConfigEncrypted, connector.AuthConfigSalt,
		connector.BaseURL, connector.Environment, connector.WebhookEnabled,
		connector.WebhookSecretEncrypted, connector.WebhookURL,
		pq.Array(connector.WebhookEvents), connector.SyncEnabled,
		connector.SyncIntervalSeconds, connector.BatchSize,
		connector.RateLimitRequestsPerSecond, connector.RateLimitBurst,
		connector.RetryMaxAttempts, connector.RetryInitialDelayMs,
		connector.RetryMaxDelayMs, connector.RetryMultiplier,
		customHeaders, connector.TimeoutSeconds,
		connector.CreatedBy, connector.UpdatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create connector: %w", err)
	}

	return nil
}

// GetByID retrieves a connector by ID
func (r *ConnectorRepository) GetByID(ctx context.Context, id string) (*domain.ConnectorEntity, error) {
	query := `
		SELECT 
			id, name, type, status, auth_config_encrypted, auth_config_salt,
			base_url, environment, webhook_enabled, webhook_secret_encrypted,
			webhook_url, webhook_events, sync_enabled, sync_interval_seconds,
			batch_size, rate_limit_requests_per_second, rate_limit_burst,
			retry_max_attempts, retry_initial_delay_ms, retry_max_delay_ms,
			retry_multiplier, custom_headers, timeout_seconds,
			last_health_check_at, last_health_check_status, last_health_check_error,
			created_at, updated_at, created_by, updated_by, version
		FROM connectors
		WHERE id = $1`

	return r.scanConnector(r.db.QueryRowContext(ctx, query, id))
}

// GetByName retrieves a connector by name
func (r *ConnectorRepository) GetByName(ctx context.Context, name string) (*domain.ConnectorEntity, error) {
	query := `
		SELECT 
			id, name, type, status, auth_config_encrypted, auth_config_salt,
			base_url, environment, webhook_enabled, webhook_secret_encrypted,
			webhook_url, webhook_events, sync_enabled, sync_interval_seconds,
			batch_size, rate_limit_requests_per_second, rate_limit_burst,
			retry_max_attempts, retry_initial_delay_ms, retry_max_delay_ms,
			retry_multiplier, custom_headers, timeout_seconds,
			last_health_check_at, last_health_check_status, last_health_check_error,
			created_at, updated_at, created_by, updated_by, version
		FROM connectors
		WHERE name = $1`

	return r.scanConnector(r.db.QueryRowContext(ctx, query, name))
}

// GetByType retrieves all connectors of a specific type
func (r *ConnectorRepository) GetByType(ctx context.Context, connectorType string) ([]*domain.ConnectorEntity, error) {
	query := `
		SELECT 
			id, name, type, status, auth_config_encrypted, auth_config_salt,
			base_url, environment, webhook_enabled, webhook_secret_encrypted,
			webhook_url, webhook_events, sync_enabled, sync_interval_seconds,
			batch_size, rate_limit_requests_per_second, rate_limit_burst,
			retry_max_attempts, retry_initial_delay_ms, retry_max_delay_ms,
			retry_multiplier, custom_headers, timeout_seconds,
			last_health_check_at, last_health_check_status, last_health_check_error,
			created_at, updated_at, created_by, updated_by, version
		FROM connectors
		WHERE type = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, connectorType)
	if err != nil {
		return nil, fmt.Errorf("failed to query connectors by type: %w", err)
	}
	defer rows.Close()

	return r.scanConnectors(rows)
}

// GetAll retrieves all connectors
func (r *ConnectorRepository) GetAll(ctx context.Context) ([]*domain.ConnectorEntity, error) {
	query := `
		SELECT 
			id, name, type, status, auth_config_encrypted, auth_config_salt,
			base_url, environment, webhook_enabled, webhook_secret_encrypted,
			webhook_url, webhook_events, sync_enabled, sync_interval_seconds,
			batch_size, rate_limit_requests_per_second, rate_limit_burst,
			retry_max_attempts, retry_initial_delay_ms, retry_max_delay_ms,
			retry_multiplier, custom_headers, timeout_seconds,
			last_health_check_at, last_health_check_status, last_health_check_error,
			created_at, updated_at, created_by, updated_by, version
		FROM connectors
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all connectors: %w", err)
	}
	defer rows.Close()

	return r.scanConnectors(rows)
}

// GetByStatus retrieves all connectors with the specified status
func (r *ConnectorRepository) GetByStatus(ctx context.Context, status domain.ConnectorStatus) ([]*domain.ConnectorEntity, error) {
	query := `
		SELECT 
			id, name, type, status, auth_config_encrypted, auth_config_salt,
			base_url, environment, webhook_enabled, webhook_secret_encrypted,
			webhook_url, webhook_events, sync_enabled, sync_interval_seconds,
			batch_size, rate_limit_requests_per_second, rate_limit_burst,
			retry_max_attempts, retry_initial_delay_ms, retry_max_delay_ms,
			retry_multiplier, custom_headers, timeout_seconds,
			last_health_check_at, last_health_check_status, last_health_check_error,
			created_at, updated_at, created_by, updated_by, version
		FROM connectors
		WHERE status = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query connectors by status: %w", err)
	}
	defer rows.Close()

	return r.scanConnectors(rows)
}

// GetActive retrieves all active connectors
func (r *ConnectorRepository) GetActive(ctx context.Context) ([]*domain.ConnectorEntity, error) {
	query := `
		SELECT 
			id, name, type, status, auth_config_encrypted, auth_config_salt,
			base_url, environment, webhook_enabled, webhook_secret_encrypted,
			webhook_url, webhook_events, sync_enabled, sync_interval_seconds,
			batch_size, rate_limit_requests_per_second, rate_limit_burst,
			retry_max_attempts, retry_initial_delay_ms, retry_max_delay_ms,
			retry_multiplier, custom_headers, timeout_seconds,
			last_health_check_at, last_health_check_status, last_health_check_error,
			created_at, updated_at, created_by, updated_by, version
		FROM connectors
		WHERE status = 'active'
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active connectors: %w", err)
	}
	defer rows.Close()

	return r.scanConnectors(rows)
}

// Update updates an existing connector
func (r *ConnectorRepository) Update(ctx context.Context, connector *domain.ConnectorEntity) error {
	query := `
		UPDATE connectors SET
			name = $2, type = $3, status = $4, auth_config_encrypted = $5,
			auth_config_salt = $6, base_url = $7, environment = $8,
			webhook_enabled = $9, webhook_secret_encrypted = $10,
			webhook_url = $11, webhook_events = $12, sync_enabled = $13,
			sync_interval_seconds = $14, batch_size = $15,
			rate_limit_requests_per_second = $16, rate_limit_burst = $17,
			retry_max_attempts = $18, retry_initial_delay_ms = $19,
			retry_max_delay_ms = $20, retry_multiplier = $21,
			custom_headers = $22, timeout_seconds = $23,
			updated_by = $24
		WHERE id = $1 AND version = $25`

	customHeaders, err := json.Marshal(connector.CustomHeaders)
	if err != nil {
		return fmt.Errorf("failed to marshal custom headers: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query,
		connector.ID, connector.Name, connector.Type, connector.Status,
		connector.AuthConfigEncrypted, connector.AuthConfigSalt,
		connector.BaseURL, connector.Environment, connector.WebhookEnabled,
		connector.WebhookSecretEncrypted, connector.WebhookURL,
		pq.Array(connector.WebhookEvents), connector.SyncEnabled,
		connector.SyncIntervalSeconds, connector.BatchSize,
		connector.RateLimitRequestsPerSecond, connector.RateLimitBurst,
		connector.RetryMaxAttempts, connector.RetryInitialDelayMs,
		connector.RetryMaxDelayMs, connector.RetryMultiplier,
		customHeaders, connector.TimeoutSeconds,
		connector.UpdatedBy, connector.Version,
	)

	if err != nil {
		return fmt.Errorf("failed to update connector: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("connector not found or version mismatch")
	}

	return nil
}

// Delete removes a connector
func (r *ConnectorRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM connectors WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete connector: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("connector not found")
	}

	return nil
}

// UpdateStatus updates the status of a connector
func (r *ConnectorRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE connectors SET status = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update connector status: %w", err)
	}

	return nil
}

// UpdateHealthCheck updates the health check status
func (r *ConnectorRepository) UpdateHealthCheck(ctx context.Context, id string, status string, errorMsg string) error {
	query := `
		UPDATE connectors SET 
			last_health_check_at = $2,
			last_health_check_status = $3,
			last_health_check_error = $4
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, time.Now(), status, errorMsg)
	if err != nil {
		return fmt.Errorf("failed to update health check: %w", err)
	}

	return nil
}

// CreateSyncEntity creates a new sync entity configuration
func (r *ConnectorRepository) CreateSyncEntity(ctx context.Context, entity *domain.ConnectorSyncEntity) error {
	query := `
		INSERT INTO connector_sync_entities (
			id, connector_id, entity_type, enabled, sync_direction,
			filters, field_mapping
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	filters, err := json.Marshal(entity.Filters)
	if err != nil {
		return fmt.Errorf("failed to marshal filters: %w", err)
	}

	fieldMapping, err := json.Marshal(entity.FieldMapping)
	if err != nil {
		return fmt.Errorf("failed to marshal field mapping: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.ConnectorID, entity.EntityType,
		entity.Enabled, entity.SyncDirection, filters, fieldMapping,
	)

	if err != nil {
		return fmt.Errorf("failed to create sync entity: %w", err)
	}

	return nil
}

// GetSyncEntities retrieves all sync entities for a connector
func (r *ConnectorRepository) GetSyncEntities(ctx context.Context, connectorID string) ([]*domain.ConnectorSyncEntity, error) {
	query := `
		SELECT id, connector_id, entity_type, enabled, sync_direction,
			   last_sync_at, filters, field_mapping, created_at, updated_at
		FROM connector_sync_entities
		WHERE connector_id = $1
		ORDER BY entity_type`

	rows, err := r.db.QueryContext(ctx, query, connectorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sync entities: %w", err)
	}
	defer rows.Close()

	var entities []*domain.ConnectorSyncEntity
	for rows.Next() {
		entity, err := r.scanSyncEntity(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

// UpdateSyncEntity updates a sync entity configuration
func (r *ConnectorRepository) UpdateSyncEntity(ctx context.Context, entity *domain.ConnectorSyncEntity) error {
	query := `
		UPDATE connector_sync_entities SET
			enabled = $3, sync_direction = $4, filters = $5, field_mapping = $6
		WHERE id = $1 AND connector_id = $2`

	filters, err := json.Marshal(entity.Filters)
	if err != nil {
		return fmt.Errorf("failed to marshal filters: %w", err)
	}

	fieldMapping, err := json.Marshal(entity.FieldMapping)
	if err != nil {
		return fmt.Errorf("failed to marshal field mapping: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.ConnectorID, entity.Enabled,
		entity.SyncDirection, filters, fieldMapping,
	)

	if err != nil {
		return fmt.Errorf("failed to update sync entity: %w", err)
	}

	return nil
}

// DeleteSyncEntity deletes a sync entity configuration
func (r *ConnectorRepository) DeleteSyncEntity(ctx context.Context, id string) error {
	query := `DELETE FROM connector_sync_entities WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sync entity: %w", err)
	}

	return nil
}

// UpdateLastSyncTime updates the last sync time for an entity
func (r *ConnectorRepository) UpdateLastSyncTime(ctx context.Context, connectorID string, entityType string, syncTime time.Time) error {
	query := `
		UPDATE connector_sync_entities SET last_sync_at = $3
		WHERE connector_id = $1 AND entity_type = $2`

	_, err := r.db.ExecContext(ctx, query, connectorID, entityType, syncTime)
	if err != nil {
		return fmt.Errorf("failed to update last sync time: %w", err)
	}

	return nil
}

// CreateAPIKey creates a new API key
func (r *ConnectorRepository) CreateAPIKey(ctx context.Context, key *domain.ConnectorAPIKey) error {
	query := `
		INSERT INTO connector_api_keys (
			id, connector_id, key_name, key_value_encrypted, key_salt,
			key_type, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		key.ID, key.ConnectorID, key.KeyName,
		key.KeyValueEncrypted, key.KeySalt,
		key.KeyType, key.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create API key: %w", err)
	}

	return nil
}

// GetAPIKeys retrieves all API keys for a connector
func (r *ConnectorRepository) GetAPIKeys(ctx context.Context, connectorID string) ([]*domain.ConnectorAPIKey, error) {
	query := `
		SELECT id, connector_id, key_name, key_value_encrypted, key_salt,
			   key_type, expires_at, last_used_at, created_at
		FROM connector_api_keys
		WHERE connector_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, connectorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query API keys: %w", err)
	}
	defer rows.Close()

	var keys []*domain.ConnectorAPIKey
	for rows.Next() {
		key, err := r.scanAPIKey(rows)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// GetAPIKeyByName retrieves an API key by name
func (r *ConnectorRepository) GetAPIKeyByName(ctx context.Context, connectorID string, keyName string) (*domain.ConnectorAPIKey, error) {
	query := `
		SELECT id, connector_id, key_name, key_value_encrypted, key_salt,
			   key_type, expires_at, last_used_at, created_at
		FROM connector_api_keys
		WHERE connector_id = $1 AND key_name = $2`

	row := r.db.QueryRowContext(ctx, query, connectorID, keyName)
	return r.scanAPIKeyRow(row)
}

// UpdateAPIKeyLastUsed updates the last used timestamp of an API key
func (r *ConnectorRepository) UpdateAPIKeyLastUsed(ctx context.Context, id string) error {
	query := `UPDATE connector_api_keys SET last_used_at = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update API key last used: %w", err)
	}

	return nil
}

// DeleteAPIKey deletes an API key
func (r *ConnectorRepository) DeleteAPIKey(ctx context.Context, id string) error {
	query := `DELETE FROM connector_api_keys WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	return nil
}

// DeleteExpiredAPIKeys deletes all expired API keys
func (r *ConnectorRepository) DeleteExpiredAPIKeys(ctx context.Context) error {
	query := `DELETE FROM connector_api_keys WHERE expires_at < $1`

	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired API keys: %w", err)
	}

	return nil
}

// CreateAuditLog creates a new audit log entry
func (r *ConnectorRepository) CreateAuditLog(ctx context.Context, log *domain.ConnectorAuditLog) error {
	query := `
		INSERT INTO connector_audit_log (
			id, connector_id, action, changed_by, old_values,
			new_values, ip_address, user_agent
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	oldValues, err := json.Marshal(log.OldValues)
	if err != nil {
		return fmt.Errorf("failed to marshal old values: %w", err)
	}

	newValues, err := json.Marshal(log.NewValues)
	if err != nil {
		return fmt.Errorf("failed to marshal new values: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		log.ID, log.ConnectorID, log.Action, log.ChangedBy,
		oldValues, newValues, log.IPAddress, log.UserAgent,
	)

	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetAuditLogs retrieves audit logs for a connector
func (r *ConnectorRepository) GetAuditLogs(ctx context.Context, connectorID string, limit int) ([]*domain.ConnectorAuditLog, error) {
	query := `
		SELECT id, connector_id, action, changed_by, changed_at,
			   old_values, new_values, ip_address, user_agent
		FROM connector_audit_log
		WHERE connector_id = $1
		ORDER BY changed_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, connectorID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*domain.ConnectorAuditLog
	for rows.Next() {
		log, err := r.scanAuditLog(rows)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// GetAuditLogsByAction retrieves audit logs by action since a specific time
func (r *ConnectorRepository) GetAuditLogsByAction(ctx context.Context, action string, since time.Time) ([]*domain.ConnectorAuditLog, error) {
	query := `
		SELECT id, connector_id, action, changed_by, changed_at,
			   old_values, new_values, ip_address, user_agent
		FROM connector_audit_log
		WHERE action = $1 AND changed_at >= $2
		ORDER BY changed_at DESC`

	rows, err := r.db.QueryContext(ctx, query, action, since)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs by action: %w", err)
	}
	defer rows.Close()

	var logs []*domain.ConnectorAuditLog
	for rows.Next() {
		log, err := r.scanAuditLog(rows)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// Helper functions for scanning

func (r *ConnectorRepository) scanConnector(row *sql.Row) (*domain.ConnectorEntity, error) {
	var connector domain.ConnectorEntity
	var customHeaders []byte
	var healthCheckStatus sql.NullString
	var healthCheckError sql.NullString

	err := row.Scan(
		&connector.ID, &connector.Name, &connector.Type, &connector.Status,
		&connector.AuthConfigEncrypted, &connector.AuthConfigSalt,
		&connector.BaseURL, &connector.Environment, &connector.WebhookEnabled,
		&connector.WebhookSecretEncrypted, &connector.WebhookURL,
		pq.Array(&connector.WebhookEvents), &connector.SyncEnabled,
		&connector.SyncIntervalSeconds, &connector.BatchSize,
		&connector.RateLimitRequestsPerSecond, &connector.RateLimitBurst,
		&connector.RetryMaxAttempts, &connector.RetryInitialDelayMs,
		&connector.RetryMaxDelayMs, &connector.RetryMultiplier,
		&customHeaders, &connector.TimeoutSeconds,
		&connector.LastHealthCheckAt, &healthCheckStatus, &healthCheckError,
		&connector.CreatedAt, &connector.UpdatedAt,
		&connector.CreatedBy, &connector.UpdatedBy, &connector.Version,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("connector not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan connector: %w", err)
	}

	if healthCheckStatus.Valid {
		connector.LastHealthCheckStatus = healthCheckStatus.String
	}
	if healthCheckError.Valid {
		connector.LastHealthCheckError = healthCheckError.String
	}

	if customHeaders != nil {
		if err := json.Unmarshal(customHeaders, &connector.CustomHeaders); err != nil {
			return nil, fmt.Errorf("failed to unmarshal custom headers: %w", err)
		}
	}

	return &connector, nil
}

func (r *ConnectorRepository) scanConnectors(rows *sql.Rows) ([]*domain.ConnectorEntity, error) {
	var connectors []*domain.ConnectorEntity

	for rows.Next() {
		var connector domain.ConnectorEntity
		var customHeaders []byte
		var healthCheckStatus sql.NullString
		var healthCheckError sql.NullString

		err := rows.Scan(
			&connector.ID, &connector.Name, &connector.Type, &connector.Status,
			&connector.AuthConfigEncrypted, &connector.AuthConfigSalt,
			&connector.BaseURL, &connector.Environment, &connector.WebhookEnabled,
			&connector.WebhookSecretEncrypted, &connector.WebhookURL,
			pq.Array(&connector.WebhookEvents), &connector.SyncEnabled,
			&connector.SyncIntervalSeconds, &connector.BatchSize,
			&connector.RateLimitRequestsPerSecond, &connector.RateLimitBurst,
			&connector.RetryMaxAttempts, &connector.RetryInitialDelayMs,
			&connector.RetryMaxDelayMs, &connector.RetryMultiplier,
			&customHeaders, &connector.TimeoutSeconds,
			&connector.LastHealthCheckAt, &healthCheckStatus, &healthCheckError,
			&connector.CreatedAt, &connector.UpdatedAt,
			&connector.CreatedBy, &connector.UpdatedBy, &connector.Version,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan connector row: %w", err)
		}

		if healthCheckStatus.Valid {
			connector.LastHealthCheckStatus = healthCheckStatus.String
		}
		if healthCheckError.Valid {
			connector.LastHealthCheckError = healthCheckError.String
		}

		if customHeaders != nil {
			if err := json.Unmarshal(customHeaders, &connector.CustomHeaders); err != nil {
				return nil, fmt.Errorf("failed to unmarshal custom headers: %w", err)
			}
		}

		connectors = append(connectors, &connector)
	}

	return connectors, nil
}

func (r *ConnectorRepository) scanSyncEntity(rows *sql.Rows) (*domain.ConnectorSyncEntity, error) {
	var entity domain.ConnectorSyncEntity
	var filters []byte
	var fieldMapping []byte

	err := rows.Scan(
		&entity.ID, &entity.ConnectorID, &entity.EntityType,
		&entity.Enabled, &entity.SyncDirection, &entity.LastSyncAt,
		&filters, &fieldMapping, &entity.CreatedAt, &entity.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan sync entity: %w", err)
	}

	if filters != nil {
		if err := json.Unmarshal(filters, &entity.Filters); err != nil {
			return nil, fmt.Errorf("failed to unmarshal filters: %w", err)
		}
	}

	if fieldMapping != nil {
		if err := json.Unmarshal(fieldMapping, &entity.FieldMapping); err != nil {
			return nil, fmt.Errorf("failed to unmarshal field mapping: %w", err)
		}
	}

	return &entity, nil
}

func (r *ConnectorRepository) scanAPIKey(rows *sql.Rows) (*domain.ConnectorAPIKey, error) {
	var key domain.ConnectorAPIKey

	err := rows.Scan(
		&key.ID, &key.ConnectorID, &key.KeyName,
		&key.KeyValueEncrypted, &key.KeySalt,
		&key.KeyType, &key.ExpiresAt, &key.LastUsedAt, &key.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan API key: %w", err)
	}

	return &key, nil
}

func (r *ConnectorRepository) scanAPIKeyRow(row *sql.Row) (*domain.ConnectorAPIKey, error) {
	var key domain.ConnectorAPIKey

	err := row.Scan(
		&key.ID, &key.ConnectorID, &key.KeyName,
		&key.KeyValueEncrypted, &key.KeySalt,
		&key.KeyType, &key.ExpiresAt, &key.LastUsedAt, &key.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("API key not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan API key: %w", err)
	}

	return &key, nil
}

func (r *ConnectorRepository) scanAuditLog(rows *sql.Rows) (*domain.ConnectorAuditLog, error) {
	var log domain.ConnectorAuditLog
	var oldValues []byte
	var newValues []byte

	err := rows.Scan(
		&log.ID, &log.ConnectorID, &log.Action, &log.ChangedBy,
		&log.ChangedAt, &oldValues, &newValues,
		&log.IPAddress, &log.UserAgent,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan audit log: %w", err)
	}

	if oldValues != nil {
		if err := json.Unmarshal(oldValues, &log.OldValues); err != nil {
			return nil, fmt.Errorf("failed to unmarshal old values: %w", err)
		}
	}

	if newValues != nil {
		if err := json.Unmarshal(newValues, &log.NewValues); err != nil {
			return nil, fmt.Errorf("failed to unmarshal new values: %w", err)
		}
	}

	return &log, nil
}