package adapters

import (
	"context"
	"database/sql"
	"fmt"
	
	"middleman/merchant/internal/domain"
)

// TransactionalProductSyncStatusRepository wraps the repository with transaction support
type TransactionalProductSyncStatusRepository struct {
	db *sql.DB
}

// NewTransactionalProductSyncStatusRepository creates a new transactional repository
func NewTransactionalProductSyncStatusRepository(db *sql.DB) *TransactionalProductSyncStatusRepository {
	return &TransactionalProductSyncStatusRepository{db: db}
}

// WithTransaction executes a function within a database transaction
func (r *TransactionalProductSyncStatusRepository) WithTransaction(ctx context.Context, fn func(repo domain.ProductSyncStatusRepository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	// Create repository with transaction
	repo := NewProductSyncStatusRepository(tx)
	
	// Execute the function
	if err := fn(repo); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// GetRepository returns a non-transactional repository
func (r *TransactionalProductSyncStatusRepository) GetRepository() domain.ProductSyncStatusRepository {
	return NewProductSyncStatusRepository(r.db)
}