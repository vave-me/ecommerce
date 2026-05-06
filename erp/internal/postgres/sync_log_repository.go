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

// SyncLogRepository implements domain.SyncLogRepository
type SyncLogRepository struct {
	tableName string
	db        postgres.DB
}

func (r *SyncLogRepository) GetByConnectorID(ctx context.Context, connectorID string) ([]*domain.SyncLog, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, entity_type, started_at, completed_at,
			duration, status, records_processed, records_success,
			records_failed, last_sync_time, error_message, metadata
		FROM %s
		WHERE metadata->>'connector_id' = $1
		ORDER BY started_at DESC
	`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, connectorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.SyncLog
	for rows.Next() {
		var log domain.SyncLog
		var metadataJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.ERPType,
			&log.EntityType,
			&log.StartedAt,
			&log.CompletedAt,
			&log.Duration,
			&log.Status,
			&log.RecordsProcessed,
			&log.RecordsSuccess,
			&log.RecordsFailed,
			&log.LastSyncTime,
			&log.ErrorMessage,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}

		if len(metadataJSON) > 0 {
			_ = json.Unmarshal(metadataJSON, &log.Metadata)
		}

		logs = append(logs, &log)
	}

	return logs, rows.Err()
}

func (r *SyncLogRepository) GetByStatus(ctx context.Context, status domain.SyncStatus) ([]*domain.SyncLog, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, entity_type, started_at, completed_at,
			duration, status, records_processed, records_success,
			records_failed, last_sync_time, error_message, metadata
		FROM %s
		WHERE status = $1
		ORDER BY started_at DESC
	`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.SyncLog
	for rows.Next() {
		var log domain.SyncLog
		var metadataJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.ERPType,
			&log.EntityType,
			&log.StartedAt,
			&log.CompletedAt,
			&log.Duration,
			&log.Status,
			&log.RecordsProcessed,
			&log.RecordsSuccess,
			&log.RecordsFailed,
			&log.LastSyncTime,
			&log.ErrorMessage,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}

		if len(metadataJSON) > 0 {
			_ = json.Unmarshal(metadataJSON, &log.Metadata)
		}

		logs = append(logs, &log)
	}

	return logs, rows.Err()
}

// NewSyncLogRepository creates a new sync log repository
func NewSyncLogRepository(tableName string, db postgres.DB) domain.SyncLogRepository {
	return &SyncLogRepository{
		tableName: tableName,
		db:        db,
	}
}

// Create creates a new sync log entry
func (r *SyncLogRepository) Create(ctx context.Context, log *domain.SyncLog) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			id, erp_type, entity_type, started_at, completed_at, 
			duration, status, records_processed, records_success, 
			records_failed, last_sync_time, error_message, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, r.tableName)

	metadataJSON, _ := json.Marshal(log.Metadata)

	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.ERPType,
		log.EntityType,
		log.StartedAt,
		log.CompletedAt,
		log.Duration,
		log.Status,
		log.RecordsProcessed,
		log.RecordsSuccess,
		log.RecordsFailed,
		log.LastSyncTime,
		log.ErrorMessage,
		metadataJSON,
	)

	return err
}

// Update updates a sync log entry
func (r *SyncLogRepository) Update(ctx context.Context, log *domain.SyncLog) error {
	query := fmt.Sprintf(`
		UPDATE %s SET
			completed_at = $2,
			duration = $3,
			status = $4,
			records_processed = $5,
			records_success = $6,
			records_failed = $7,
			error_message = $8,
			metadata = $9
		WHERE id = $1
	`, r.tableName)

	metadataJSON, _ := json.Marshal(log.Metadata)

	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.CompletedAt,
		log.Duration,
		log.Status,
		log.RecordsProcessed,
		log.RecordsSuccess,
		log.RecordsFailed,
		log.ErrorMessage,
		metadataJSON,
	)

	return err
}

// GetByID gets a sync log by ID
func (r *SyncLogRepository) GetByID(ctx context.Context, id string) (*domain.SyncLog, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, entity_type, started_at, completed_at,
			duration, status, records_processed, records_success,
			records_failed, last_sync_time, error_message, metadata
		FROM %s
		WHERE id = $1
	`, r.tableName)

	var log domain.SyncLog
	var metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID,
		&log.ERPType,
		&log.EntityType,
		&log.StartedAt,
		&log.CompletedAt,
		&log.Duration,
		&log.Status,
		&log.RecordsProcessed,
		&log.RecordsSuccess,
		&log.RecordsFailed,
		&log.LastSyncTime,
		&log.ErrorMessage,
		&metadataJSON,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err == nil && len(metadataJSON) > 0 {
		_ = json.Unmarshal(metadataJSON, &log.Metadata)
	}

	return &log, err
}

// GetRecent gets recent sync logs for an ERP type
func (r *SyncLogRepository) GetRecent(ctx context.Context, erpType string, limit int) ([]*domain.SyncLog, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, entity_type, started_at, completed_at,
			duration, status, records_processed, records_success,
			records_failed, last_sync_time, error_message, metadata
		FROM %s
		WHERE erp_type = $1
		ORDER BY started_at DESC
		LIMIT $2
	`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, erpType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.SyncLog
	for rows.Next() {
		var log domain.SyncLog
		var metadataJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.ERPType,
			&log.EntityType,
			&log.StartedAt,
			&log.CompletedAt,
			&log.Duration,
			&log.Status,
			&log.RecordsProcessed,
			&log.RecordsSuccess,
			&log.RecordsFailed,
			&log.LastSyncTime,
			&log.ErrorMessage,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}

		if len(metadataJSON) > 0 {
			_ = json.Unmarshal(metadataJSON, &log.Metadata)
		}

		logs = append(logs, &log)
	}

	return logs, rows.Err()
}

// GetByDateRange gets sync logs within a date range
func (r *SyncLogRepository) GetByDateRange(ctx context.Context, erpType string, start, end time.Time) ([]*domain.SyncLog, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, entity_type, started_at, completed_at,
			duration, status, records_processed, records_success,
			records_failed, last_sync_time, error_message, metadata
		FROM %s
		WHERE erp_type = $1 AND started_at >= $2 AND started_at <= $3
		ORDER BY started_at DESC
	`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, erpType, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.SyncLog
	for rows.Next() {
		var log domain.SyncLog
		var metadataJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.ERPType,
			&log.EntityType,
			&log.StartedAt,
			&log.CompletedAt,
			&log.Duration,
			&log.Status,
			&log.RecordsProcessed,
			&log.RecordsSuccess,
			&log.RecordsFailed,
			&log.LastSyncTime,
			&log.ErrorMessage,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}

		if len(metadataJSON) > 0 {
			_ = json.Unmarshal(metadataJSON, &log.Metadata)
		}

		logs = append(logs, &log)
	}

	return logs, rows.Err()
}
