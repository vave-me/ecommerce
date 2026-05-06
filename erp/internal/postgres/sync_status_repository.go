package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"middleman/erp/internal/domain"
	"middleman/internal/postgres"
)

// SyncStatusRepository implements domain.SyncStatusRepository
type SyncStatusRepository struct {
	tableName string
	db        postgres.DB
}

// NewSyncStatusRepository creates a new sync status repository
func NewSyncStatusRepository(tableName string, db postgres.DB) domain.SyncStatusRepository {
	return &SyncStatusRepository{
		tableName: tableName,
		db:        db,
	}
}

// GetLastSyncTime gets the last sync time for an entity type
func (r *SyncStatusRepository) GetLastSyncTime(ctx context.Context, erpType, entityType string) (time.Time, error) {
	query := fmt.Sprintf(`
		SELECT last_synced_at 
		FROM %s 
		WHERE erp_type = $1 AND entity_type = $2 AND entity_id = ''
		ORDER BY last_synced_at DESC
		LIMIT 1
	`, r.tableName)

	var lastSync time.Time
	err := r.db.QueryRowContext(ctx, query, erpType, entityType).Scan(&lastSync)
	if err == sql.ErrNoRows {
		// Return zero time if no sync has been performed
		return time.Time{}, nil
	}

	return lastSync, err
}

// UpdateLastSync updates the last sync time for an entity type
func (r *SyncStatusRepository) UpdateLastSync(ctx context.Context, erpType, entityType string, syncTime time.Time) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			id, erp_type, entity_type, entity_id, last_synced_at, 
			status, error_message, retry_count, created_at, updated_at
		) VALUES (
			gen_random_uuid(), $1, $2, '', $3, 
			'success', '', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
		ON CONFLICT (erp_type, entity_type, entity_id) 
		DO UPDATE SET 
			last_synced_at = EXCLUDED.last_synced_at,
			status = EXCLUDED.status,
			updated_at = CURRENT_TIMESTAMP
	`, r.tableName)

	_, err := r.db.ExecContext(ctx, query, erpType, entityType, syncTime)
	return err
}

// Save saves a sync status record
func (r *SyncStatusRepository) Save(ctx context.Context, status *domain.SyncStatusRecord) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			id, erp_type, entity_type, entity_id, last_synced_at,
			status, error_message, retry_count, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		ON CONFLICT (erp_type, entity_type, entity_id) 
		DO UPDATE SET 
			last_synced_at = EXCLUDED.last_synced_at,
			status = EXCLUDED.status,
			error_message = EXCLUDED.error_message,
			retry_count = EXCLUDED.retry_count,
			updated_at = EXCLUDED.updated_at
	`, r.tableName)

	if status.ID == "" {
		status.ID = generateUUID()
	}
	if status.CreatedAt.IsZero() {
		status.CreatedAt = time.Now()
	}
	status.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		status.ID,
		status.ERPType,
		status.EntityType,
		status.EntityID,
		status.LastSyncedAt,
		status.Status,
		status.ErrorMessage,
		status.RetryCount,
		status.CreatedAt,
		status.UpdatedAt,
	)

	return err
}

// GetByEntity gets sync status for a specific entity
func (r *SyncStatusRepository) GetByEntity(ctx context.Context, erpType, entityType, entityID string) (*domain.SyncStatusRecord, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, entity_type, entity_id, last_synced_at,
			status, error_message, retry_count, created_at, updated_at
		FROM %s
		WHERE erp_type = $1 AND entity_type = $2 AND entity_id = $3
	`, r.tableName)

	var status domain.SyncStatusRecord
	err := r.db.QueryRowContext(ctx, query, erpType, entityType, entityID).Scan(
		&status.ID,
		&status.ERPType,
		&status.EntityType,
		&status.EntityID,
		&status.LastSyncedAt,
		&status.Status,
		&status.ErrorMessage,
		&status.RetryCount,
		&status.CreatedAt,
		&status.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &status, err
}

// GetFailedSyncs gets failed sync records
func (r *SyncStatusRepository) GetFailedSyncs(ctx context.Context, limit int) ([]*domain.SyncStatusRecord, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, entity_type, entity_id, last_synced_at,
			status, error_message, retry_count, created_at, updated_at
		FROM %s
		WHERE status = 'failed' AND retry_count < 3
		ORDER BY updated_at DESC
		LIMIT $1
	`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*domain.SyncStatusRecord
	for rows.Next() {
		var status domain.SyncStatusRecord
		err := rows.Scan(
			&status.ID,
			&status.ERPType,
			&status.EntityType,
			&status.EntityID,
			&status.LastSyncedAt,
			&status.Status,
			&status.ErrorMessage,
			&status.RetryCount,
			&status.CreatedAt,
			&status.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &status)
	}

	return records, rows.Err()
}

// generateUUID generates a new UUID
