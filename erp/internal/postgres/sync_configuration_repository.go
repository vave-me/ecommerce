package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"middleman/erp/internal/domain"
	"middleman/internal/postgres"
)

// SyncConfigurationRepository implements domain.SyncConfigurationRepository
type SyncConfigurationRepository struct {
	tableName string
	db        postgres.DB
}

// NewSyncConfigurationRepository creates a new sync configuration repository
func NewSyncConfigurationRepository(tableName string, db postgres.DB) domain.SyncConfigurationRepository {
	return &SyncConfigurationRepository{
		tableName: tableName,
		db:        db,
	}
}

// GetByEntity gets sync configuration for an entity type
func (r *SyncConfigurationRepository) GetByEntity(ctx context.Context, erpType, entityType string) (*domain.SyncConfiguration, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, entity_type, enabled, sync_interval,
			batch_size, retry_attempts, retry_delay, filters,
			created_at, updated_at
		FROM %s
		WHERE erp_type = $1 AND entity_type = $2
	`, r.tableName)

	var config domain.SyncConfiguration
	var filtersJSON []byte

	err := r.db.QueryRowContext(ctx, query, erpType, entityType).Scan(
		&config.ID,
		&config.ERPType,
		&config.EntityType,
		&config.Enabled,
		&config.SyncInterval,
		&config.BatchSize,
		&config.RetryAttempts,
		&config.RetryDelay,
		&filtersJSON,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err == nil && len(filtersJSON) > 0 {
		_ = json.Unmarshal(filtersJSON, &config.Filters)
	}

	return &config, err
}

// Save saves a sync configuration
func (r *SyncConfigurationRepository) Save(ctx context.Context, config *domain.SyncConfiguration) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			id, erp_type, entity_type, enabled, sync_interval,
			batch_size, retry_attempts, retry_delay, filters,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (erp_type, entity_type) 
		DO UPDATE SET
			enabled = EXCLUDED.enabled,
			sync_interval = EXCLUDED.sync_interval,
			batch_size = EXCLUDED.batch_size,
			retry_attempts = EXCLUDED.retry_attempts,
			retry_delay = EXCLUDED.retry_delay,
			filters = EXCLUDED.filters,
			updated_at = EXCLUDED.updated_at
	`, r.tableName)

	if config.ID == "" {
		config.ID = generateUUID()
	}
	if config.CreatedAt.IsZero() {
		config.CreatedAt = time.Now()
	}
	config.UpdatedAt = time.Now()

	filtersJSON, _ := json.Marshal(config.Filters)

	_, err := r.db.ExecContext(ctx, query,
		config.ID,
		config.ERPType,
		config.EntityType,
		config.Enabled,
		config.SyncInterval,
		config.BatchSize,
		config.RetryAttempts,
		config.RetryDelay,
		filtersJSON,
		config.CreatedAt,
		config.UpdatedAt,
	)

	return err
}

// GetEnabled gets all enabled sync configurations
func (r *SyncConfigurationRepository) GetEnabled(ctx context.Context) ([]*domain.SyncConfiguration, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, entity_type, enabled, sync_interval,
			batch_size, retry_attempts, retry_delay, filters,
			created_at, updated_at
		FROM %s
		WHERE enabled = true
		ORDER BY erp_type, entity_type
	`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*domain.SyncConfiguration
	for rows.Next() {
		var config domain.SyncConfiguration
		var filtersJSON []byte

		err := rows.Scan(
			&config.ID,
			&config.ERPType,
			&config.EntityType,
			&config.Enabled,
			&config.SyncInterval,
			&config.BatchSize,
			&config.RetryAttempts,
			&config.RetryDelay,
			&filtersJSON,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if len(filtersJSON) > 0 {
			_ = json.Unmarshal(filtersJSON, &config.Filters)
		}

		configs = append(configs, &config)
	}

	return configs, rows.Err()
}

// Delete deletes a sync configuration
func (r *SyncConfigurationRepository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE id = $1
	`, r.tableName)

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// generateUUID generates a new UUID
func generateUUID() string {
	// In production, use a proper UUID generator
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

