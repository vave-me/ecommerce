package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"middleman/streams/internal/domain"
)

// WebhookDeliveryRepository is a PostgreSQL implementation of the webhook delivery repository
type WebhookDeliveryRepository struct {
	db *sql.DB
}

// NewWebhookDeliveryRepository creates a new PostgreSQL webhook delivery repository
func NewWebhookDeliveryRepository(db *sql.DB) *WebhookDeliveryRepository {
	return &WebhookDeliveryRepository{db: db}
}

// Create creates a new webhook delivery record
func (r *WebhookDeliveryRepository) Create(delivery *domain.WebhookDelivery) error {
	query := `
		INSERT INTO webhook_deliveries (
			id, subscription_id, event_id, event_type, payload,
			status, attempts, created_at, completed_at, last_attempt_at,
			next_retry_at, response_status, response_body, error
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.Exec(query,
		delivery.ID,
		delivery.SubscriptionID,
		delivery.EventID,
		delivery.EventType,
		delivery.Payload,
		delivery.Status,
		delivery.Attempts,
		delivery.CreatedAt,
		delivery.CompletedAt,
		delivery.LastAttemptAt,
		delivery.NextRetryAt,
		delivery.ResponseStatus,
		delivery.ResponseBody,
		delivery.Error,
	)

	return err
}

// Update updates an existing webhook delivery record
func (r *WebhookDeliveryRepository) Update(delivery *domain.WebhookDelivery) error {
	query := `
		UPDATE webhook_deliveries SET
			status = $2,
			attempts = $3,
			completed_at = $4,
			last_attempt_at = $5,
			next_retry_at = $6,
			response_status = $7,
			response_body = $8,
			error = $9
		WHERE id = $1
	`

	result, err := r.db.Exec(query,
		delivery.ID,
		delivery.Status,
		delivery.Attempts,
		delivery.CompletedAt,
		delivery.LastAttemptAt,
		delivery.NextRetryAt,
		delivery.ResponseStatus,
		delivery.ResponseBody,
		delivery.Error,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("webhook delivery not found: %s", delivery.ID)
	}

	return nil
}

// Find finds a webhook delivery by ID
func (r *WebhookDeliveryRepository) Find(id string) (*domain.WebhookDelivery, error) {
	query := `
		SELECT id, subscription_id, event_id, event_type, payload,
			status, attempts, created_at, completed_at, last_attempt_at,
			next_retry_at, response_status, response_body, error
		FROM webhook_deliveries
		WHERE id = $1
	`

	delivery := &domain.WebhookDelivery{}

	err := r.db.QueryRow(query, id).Scan(
		&delivery.ID,
		&delivery.SubscriptionID,
		&delivery.EventID,
		&delivery.EventType,
		&delivery.Payload,
		&delivery.Status,
		&delivery.Attempts,
		&delivery.CreatedAt,
		&delivery.CompletedAt,
		&delivery.LastAttemptAt,
		&delivery.NextRetryAt,
		&delivery.ResponseStatus,
		&delivery.ResponseBody,
		&delivery.Error,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("webhook delivery not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

// FindBySubscription finds webhook deliveries for a specific subscription
func (r *WebhookDeliveryRepository) FindBySubscription(subscriptionID string) ([]*domain.WebhookDelivery, error) {
	query := `
		SELECT id, subscription_id, event_id, event_type, payload,
			status, attempts, created_at, completed_at, last_attempt_at,
			next_retry_at, response_status, response_body, error
		FROM webhook_deliveries
		WHERE subscription_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`

	rows, err := r.db.Query(query, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []*domain.WebhookDelivery
	for rows.Next() {
		delivery := &domain.WebhookDelivery{}

		err := rows.Scan(
			&delivery.ID,
			&delivery.SubscriptionID,
			&delivery.EventID,
			&delivery.EventType,
			&delivery.Payload,
			&delivery.Status,
			&delivery.Attempts,
			&delivery.CreatedAt,
			&delivery.CompletedAt,
			&delivery.LastAttemptAt,
			&delivery.NextRetryAt,
			&delivery.ResponseStatus,
			&delivery.ResponseBody,
			&delivery.Error,
		)
		if err != nil {
			return nil, err
		}

		deliveries = append(deliveries, delivery)
	}

	return deliveries, rows.Err()
}

// FindPendingDeliveries finds webhook deliveries that need to be processed or retried
func (r *WebhookDeliveryRepository) FindPendingDeliveries(limit int) ([]*domain.WebhookDelivery, error) {
	query := `
		SELECT id, subscription_id, event_id, event_type, payload,
			status, attempts, created_at, completed_at, last_attempt_at,
			next_retry_at, response_status, response_body, error
		FROM webhook_deliveries
		WHERE status IN ($1, $2) 
			AND (next_retry_at IS NULL OR next_retry_at <= $3)
		ORDER BY created_at ASC
		LIMIT $4
	`

	rows, err := r.db.Query(query, 
		domain.DeliveryStatusPending, 
		domain.DeliveryStatusRetrying,
		time.Now(),
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []*domain.WebhookDelivery
	for rows.Next() {
		delivery := &domain.WebhookDelivery{}

		err := rows.Scan(
			&delivery.ID,
			&delivery.SubscriptionID,
			&delivery.EventID,
			&delivery.EventType,
			&delivery.Payload,
			&delivery.Status,
			&delivery.Attempts,
			&delivery.CreatedAt,
			&delivery.CompletedAt,
			&delivery.LastAttemptAt,
			&delivery.NextRetryAt,
			&delivery.ResponseStatus,
			&delivery.ResponseBody,
			&delivery.Error,
		)
		if err != nil {
			return nil, err
		}

		deliveries = append(deliveries, delivery)
	}

	return deliveries, rows.Err()
}

// CleanupOld deletes webhook deliveries older than the specified duration
func (r *WebhookDeliveryRepository) CleanupOld(olderThan time.Duration) error {
	query := `
		DELETE FROM webhook_deliveries
		WHERE created_at < $1
	`

	cutoffTime := time.Now().Add(-olderThan)
	_, err := r.db.Exec(query, cutoffTime)
	return err
}

// GetStats returns delivery statistics for a subscription
func (r *WebhookDeliveryRepository) GetStats(subscriptionID string) (*domain.WebhookDeliveryStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = $1 THEN 1 END) as successful,
			COUNT(CASE WHEN status = $2 THEN 1 END) as failed,
			COUNT(CASE WHEN status IN ($3, $4) THEN 1 END) as pending,
			AVG(CASE WHEN status = $1 THEN attempts END) as avg_attempts
		FROM webhook_deliveries
		WHERE subscription_id = $5
	`

	stats := &domain.WebhookDeliveryStats{}
	var avgAttempts sql.NullFloat64

	err := r.db.QueryRow(query,
		domain.DeliveryStatusDelivered,
		domain.DeliveryStatusFailed,
		domain.DeliveryStatusPending,
		domain.DeliveryStatusRetrying,
		subscriptionID,
	).Scan(
		&stats.Total,
		&stats.Successful,
		&stats.Failed,
		&stats.Pending,
		&avgAttempts,
	)

	if err != nil {
		return nil, err
	}

	if avgAttempts.Valid {
		stats.AverageAttempts = avgAttempts.Float64
	}

	return stats, nil
}

// Ensure interface compliance
var _ domain.WebhookDeliveryRepository = (*WebhookDeliveryRepository)(nil)