package postgres

import (
	"context"
	"database/sql"

	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/postgres"
)

// WebhookEventRepository implements domain.WebhookEventRepository
type WebhookEventRepository struct {
	tableName string
	db        postgres.DB
}

func (r *WebhookEventRepository) GetByConnectorID(ctx context.Context, connectorID string) ([]*domain.WebhookEvent, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, event_id, event_type, source,
			signature, payload, received_at, processed_at, status, 
			error_message, retry_count
		FROM %s
		WHERE source LIKE $1
		ORDER BY received_at DESC
	`, r.tableName)

	// Source contains connector ID in format "ERPType-ConnectorID"
	sourcePattern := fmt.Sprintf("%%-%%s", connectorID)

	rows, err := r.db.QueryContext(ctx, query, sourcePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.WebhookEvent
	for rows.Next() {
		var event domain.WebhookEvent
		err := rows.Scan(
			&event.ID,
			&event.ERPType,
			&event.EventID,
			&event.EventType,
			&event.Source,
			&event.Signature,
			&event.Payload,
			&event.ReceivedAt,
			&event.ProcessedAt,
			&event.Status,
			&event.ErrorMessage,
			&event.RetryCount,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, rows.Err()
}

func (r *WebhookEventRepository) GetByStatus(ctx context.Context, status domain.WebhookStatus) ([]*domain.WebhookEvent, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, event_id, event_type, source,
			signature, payload, received_at, processed_at, status, 
			error_message, retry_count
		FROM %s
		WHERE status = $1
		ORDER BY received_at DESC
	`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.WebhookEvent
	for rows.Next() {
		var event domain.WebhookEvent
		err := rows.Scan(
			&event.ID,
			&event.ERPType,
			&event.EventID,
			&event.EventType,
			&event.Source,
			&event.Signature,
			&event.Payload,
			&event.ReceivedAt,
			&event.ProcessedAt,
			&event.Status,
			&event.ErrorMessage,
			&event.RetryCount,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, rows.Err()
}

// NewWebhookEventRepository creates a new webhook event repository
func NewWebhookEventRepository(tableName string, db postgres.DB) domain.WebhookEventRepository {
	return &WebhookEventRepository{
		tableName: tableName,
		db:        db,
	}
}

// Create creates a new webhook event
func (r *WebhookEventRepository) Create(ctx context.Context, event *domain.WebhookEvent) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			id, erp_type, event_id, event_type, source, 
			signature, payload, received_at, status, error_message, retry_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, r.tableName)

	_, err := r.db.ExecContext(ctx, query,
		event.ID,
		event.ERPType,
		event.EventID,
		event.EventType,
		event.Source,
		event.Signature,
		event.Payload,
		event.ReceivedAt,
		event.Status,
		event.ErrorMessage,
		event.RetryCount,
	)

	return err
}

// Update updates a webhook event
func (r *WebhookEventRepository) Update(ctx context.Context, event *domain.WebhookEvent) error {
	query := fmt.Sprintf(`
		UPDATE %s SET
			processed_at = $2,
			status = $3,
			error_message = $4,
			retry_count = $5
		WHERE id = $1
	`, r.tableName)

	_, err := r.db.ExecContext(ctx, query,
		event.ID,
		event.ProcessedAt,
		event.Status,
		event.ErrorMessage,
		event.RetryCount,
	)

	return err
}

// UpdateStatus updates the status of a webhook event
func (r *WebhookEventRepository) UpdateStatus(ctx context.Context, id string, status domain.WebhookStatus, errorMessage string) error {
	query := fmt.Sprintf(`
		UPDATE %s SET
			status = $2,
			error_message = $3,
			processed_at = CASE WHEN $2 IN ('processed', 'failed') THEN CURRENT_TIMESTAMP ELSE processed_at END
		WHERE id = $1
	`, r.tableName)

	_, err := r.db.ExecContext(ctx, query, id, status, errorMessage)
	return err
}

// GetByID gets a webhook event by ID
func (r *WebhookEventRepository) GetByID(ctx context.Context, id string) (*domain.WebhookEvent, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, event_id, event_type, source,
			signature, payload, received_at, processed_at, status, 
			error_message, retry_count
		FROM %s
		WHERE id = $1
	`, r.tableName)

	var event domain.WebhookEvent
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.ERPType,
		&event.EventID,
		&event.EventType,
		&event.Source,
		&event.Signature,
		&event.Payload,
		&event.ReceivedAt,
		&event.ProcessedAt,
		&event.Status,
		&event.ErrorMessage,
		&event.RetryCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &event, err
}

// GetPending gets pending webhook events
func (r *WebhookEventRepository) GetPending(ctx context.Context, limit int) ([]*domain.WebhookEvent, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, event_id, event_type, source,
			signature, payload, received_at, processed_at, status, 
			error_message, retry_count
		FROM %s
		WHERE status IN ('received', 'processing') AND retry_count < 3
		ORDER BY received_at ASC
		LIMIT $1
	`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.WebhookEvent
	for rows.Next() {
		var event domain.WebhookEvent
		err := rows.Scan(
			&event.ID,
			&event.ERPType,
			&event.EventID,
			&event.EventType,
			&event.Source,
			&event.Signature,
			&event.Payload,
			&event.ReceivedAt,
			&event.ProcessedAt,
			&event.Status,
			&event.ErrorMessage,
			&event.RetryCount,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, rows.Err()
}

// GetByERPType gets webhook events by ERP type
func (r *WebhookEventRepository) GetByERPType(ctx context.Context, erpType string, limit int) ([]*domain.WebhookEvent, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, erp_type, event_id, event_type, source,
			signature, payload, received_at, processed_at, status, 
			error_message, retry_count
		FROM %s
		WHERE erp_type = $1
		ORDER BY received_at DESC
		LIMIT $2
	`, r.tableName)

	rows, err := r.db.QueryContext(ctx, query, erpType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.WebhookEvent
	for rows.Next() {
		var event domain.WebhookEvent
		err := rows.Scan(
			&event.ID,
			&event.ERPType,
			&event.EventID,
			&event.EventType,
			&event.Source,
			&event.Signature,
			&event.Payload,
			&event.ReceivedAt,
			&event.ProcessedAt,
			&event.Status,
			&event.ErrorMessage,
			&event.RetryCount,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, rows.Err()
}

