package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"middleman/merchant/internal/domain"
)

type productSyncStatusRepository struct {
	db interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}
}

// NewProductSyncStatusRepository creates a new PostgreSQL-based repository
// Can accept either *sql.DB or *sql.Tx
func NewProductSyncStatusRepository(db interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}) domain.ProductSyncStatusRepository {
	return &productSyncStatusRepository{db: db}
}

func (r *productSyncStatusRepository) Create(ctx context.Context, status *domain.ProductSyncStatus) error {
	query := `
		INSERT INTO products_sync_status (product_id, merchant_id, sync_status, last_synced_at, last_error, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	now := time.Now()
	status.CreatedAt = now
	status.UpdatedAt = now
	
	_, err := r.db.ExecContext(ctx, query,
		status.ProductID,
		status.MerchantID,
		status.SyncStatus,
		status.LastSyncedAt,
		status.LastError,
		status.CreatedAt,
		status.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create sync status: %w", err)
	}
	
	return nil
}

func (r *productSyncStatusRepository) Update(ctx context.Context, status *domain.ProductSyncStatus) error {
	query := `
		UPDATE products_sync_status 
		SET sync_status = $2, last_synced_at = $3, last_error = $4, updated_at = $5
		WHERE product_id = $1
	`
	
	status.UpdatedAt = time.Now()
	
	result, err := r.db.ExecContext(ctx, query,
		status.ProductID,
		status.SyncStatus,
		status.LastSyncedAt,
		status.LastError,
		status.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update sync status: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("product sync status not found: %s", status.ProductID)
	}
	
	return nil
}

func (r *productSyncStatusRepository) FindByProductID(ctx context.Context, productID string) (*domain.ProductSyncStatus, error) {
	query := `
		SELECT product_id, merchant_id, sync_status, last_synced_at, last_error, created_at, updated_at
		FROM products_sync_status
		WHERE product_id = $1
	`
	
	var status domain.ProductSyncStatus
	var lastSyncedAt sql.NullTime
	var lastError sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&status.ProductID,
		&status.MerchantID,
		&status.SyncStatus,
		&lastSyncedAt,
		&lastError,
		&status.CreatedAt,
		&status.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to find sync status: %w", err)
	}
	
	if lastSyncedAt.Valid {
		status.LastSyncedAt = &lastSyncedAt.Time
	}
	
	status.LastError = lastError.String
	
	return &status, nil
}

func (r *productSyncStatusRepository) FindByStatus(ctx context.Context, status string, limit int, offset int) ([]*domain.ProductSyncStatus, error) {
	query := `
		SELECT product_id, merchant_id, sync_status, last_synced_at, last_error, created_at, updated_at
		FROM products_sync_status
		WHERE sync_status = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find by status: %w", err)
	}
	defer rows.Close()
	
	var results []*domain.ProductSyncStatus
	
	for rows.Next() {
		var s domain.ProductSyncStatus
		var lastSyncedAt sql.NullTime
		var lastError sql.NullString
		
		err := rows.Scan(
			&s.ProductID,
			&s.MerchantID,
			&s.SyncStatus,
			&lastSyncedAt,
			&lastError,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		
		if lastSyncedAt.Valid {
			s.LastSyncedAt = &lastSyncedAt.Time
		}
		
		s.LastError = lastError.String
		
		results = append(results, &s)
	}
	
	return results, nil
}

func (r *productSyncStatusRepository) FindFailedSyncs(ctx context.Context, since time.Time, limit int) ([]*domain.ProductSyncStatus, error) {
	query := `
		SELECT product_id, merchant_id, sync_status, last_synced_at, last_error, created_at, updated_at
		FROM products_sync_status
		WHERE sync_status = $1 AND updated_at >= $2
		ORDER BY updated_at DESC
		LIMIT $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, domain.SyncStatusFailed, since, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find failed syncs: %w", err)
	}
	defer rows.Close()
	
	var results []*domain.ProductSyncStatus
	
	for rows.Next() {
		var s domain.ProductSyncStatus
		var lastSyncedAt sql.NullTime
		var lastError sql.NullString
		
		err := rows.Scan(
			&s.ProductID,
			&s.MerchantID,
			&s.SyncStatus,
			&lastSyncedAt,
			&lastError,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		
		if lastSyncedAt.Valid {
			s.LastSyncedAt = &lastSyncedAt.Time
		}
		
		s.LastError = lastError.String
		
		results = append(results, &s)
	}
	
	return results, nil
}

func (r *productSyncStatusRepository) DeleteByProductID(ctx context.Context, productID string) error {
	query := `DELETE FROM products_sync_status WHERE product_id = $1`
	
	_, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete sync status: %w", err)
	}
	
	return nil
}

func (r *productSyncStatusRepository) GetSyncStats(ctx context.Context) (*domain.SyncStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN sync_status = $1 THEN 1 END) as synced,
			COUNT(CASE WHEN sync_status = $2 THEN 1 END) as pending,
			COUNT(CASE WHEN sync_status = $3 THEN 1 END) as failed,
			MAX(last_synced_at) as last_sync
		FROM products_sync_status
	`
	
	var stats domain.SyncStats
	var lastSync sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query,
		domain.SyncStatusSynced,
		domain.SyncStatusPending,
		domain.SyncStatusFailed,
	).Scan(
		&stats.TotalProducts,
		&stats.SyncedProducts,
		&stats.PendingProducts,
		&stats.FailedProducts,
		&lastSync,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get sync stats: %w", err)
	}
	
	if lastSync.Valid {
		stats.LastSyncTime = &lastSync.Time
	}
	
	return &stats, nil
}