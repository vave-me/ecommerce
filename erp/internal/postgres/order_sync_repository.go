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

type orderSyncRepository struct {
	tableName string
	db        postgres.DB
}

// NewOrderSyncRepository creates a new OrderSyncRepository
func NewOrderSyncRepository(tableName string, db postgres.DB) domain.OrderSyncRepository {
	return &orderSyncRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r *orderSyncRepository) Create(ctx context.Context, sync *domain.OrderSync) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			id, connector_id, order_id, direction, status, 
			attempted_at, completed_at, error, payload
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, r.tableName)

	payloadJSON, err := json.Marshal(sync.Payload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		sync.ID,
		sync.ConnectorID,
		sync.OrderID,
		sync.Direction,
		sync.Status,
		sync.AttemptedAt,
		sync.CompletedAt,
		sync.Error,
		payloadJSON,
	)

	if err != nil {
		return fmt.Errorf("creating order sync: %w", err)
	}

	return nil
}

func (r *orderSyncRepository) Update(ctx context.Context, sync *domain.OrderSync) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET status = $2, completed_at = $3, error = $4, payload = $5, updated_at = $6
		WHERE id = $1
	`, r.tableName)

	payloadJSON, err := json.Marshal(sync.Payload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query,
		sync.ID,
		sync.Status,
		sync.CompletedAt,
		sync.Error,
		payloadJSON,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("updating order sync: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("order sync not found: %s", sync.ID)
	}

	return nil
}

func (r *orderSyncRepository) GetByID(ctx context.Context, id string) (*domain.OrderSync, error) {
	query := fmt.Sprintf(`
		SELECT id, connector_id, order_id, direction, status, 
		       attempted_at, completed_at, error, payload, created_at, updated_at
		FROM %s
		WHERE id = $1
	`, r.tableName)

	row := r.db.QueryRowContext(ctx, query, id)

	var sync domain.OrderSync
	var payloadJSON []byte
	var createdAt, updatedAt time.Time

	err := row.Scan(
		&sync.ID,
		&sync.ConnectorID,
		&sync.OrderID,
		&sync.Direction,
		&sync.Status,
		&sync.AttemptedAt,
		&sync.CompletedAt,
		&sync.Error,
		&payloadJSON,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order sync not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("getting order sync: %w", err)
	}

	if len(payloadJSON) > 0 {
		if err := json.Unmarshal(payloadJSON, &sync.Payload); err != nil {
			return nil, fmt.Errorf("unmarshaling payload: %w", err)
		}
	}

	return &sync, nil
}

func (r *orderSyncRepository) GetByOrderID(ctx context.Context, orderID string) ([]*domain.OrderSync, error) {
	query := fmt.Sprintf(`
		SELECT id, connector_id, order_id, direction, status, 
		       attempted_at, completed_at, error, payload, created_at, updated_at
		FROM %s
		WHERE order_id = $1
		ORDER BY attempted_at DESC
	`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("querying order syncs: %w", err)
	}
	defer rows.Close()

	var syncs []*domain.OrderSync
	for rows.Next() {
		var sync domain.OrderSync
		var payloadJSON []byte
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&sync.ID,
			&sync.ConnectorID,
			&sync.OrderID,
			&sync.Direction,
			&sync.Status,
			&sync.AttemptedAt,
			&sync.CompletedAt,
			&sync.Error,
			&payloadJSON,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning order sync: %w", err)
		}

		if len(payloadJSON) > 0 {
			if err := json.Unmarshal(payloadJSON, &sync.Payload); err != nil {
				return nil, fmt.Errorf("unmarshaling payload: %w", err)
			}
		}

		syncs = append(syncs, &sync)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating order syncs: %w", err)
	}

	return syncs, nil
}

func (r *orderSyncRepository) GetByConnectorID(ctx context.Context, connectorID string) ([]*domain.OrderSync, error) {
	query := fmt.Sprintf(`
		SELECT id, connector_id, order_id, direction, status, 
		       attempted_at, completed_at, error, payload, created_at, updated_at
		FROM %s
		WHERE connector_id = $1
		ORDER BY attempted_at DESC
		LIMIT 100
	`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, connectorID)
	if err != nil {
		return nil, fmt.Errorf("querying order syncs: %w", err)
	}
	defer rows.Close()

	var syncs []*domain.OrderSync
	for rows.Next() {
		var sync domain.OrderSync
		var payloadJSON []byte
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&sync.ID,
			&sync.ConnectorID,
			&sync.OrderID,
			&sync.Direction,
			&sync.Status,
			&sync.AttemptedAt,
			&sync.CompletedAt,
			&sync.Error,
			&payloadJSON,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning order sync: %w", err)
		}

		if len(payloadJSON) > 0 {
			if err := json.Unmarshal(payloadJSON, &sync.Payload); err != nil {
				return nil, fmt.Errorf("unmarshaling payload: %w", err)
			}
		}

		syncs = append(syncs, &sync)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating order syncs: %w", err)
	}

	return syncs, nil
}
