package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"middleman/streams/internal/domain"
)

// WebhookSubscriptionRepository is a PostgreSQL implementation of the webhook subscription repository
type WebhookSubscriptionRepository struct {
	db *sql.DB
}

// NewWebhookSubscriptionRepository creates a new PostgreSQL webhook subscription repository
func NewWebhookSubscriptionRepository(db *sql.DB) *WebhookSubscriptionRepository {
	return &WebhookSubscriptionRepository{db: db}
}

// Create creates a new webhook subscription
func (r *WebhookSubscriptionRepository) Create(subscription *domain.WebhookSubscription) error {
	query := `
		INSERT INTO webhook_subscriptions (
			id, url, secret, events, headers, 
			retry_max_retries, retry_backoff_factor, retry_initial_delay, retry_max_backoff,
			active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	eventsJSON, err := json.Marshal(subscription.Events)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	headersJSON, err := json.Marshal(subscription.Headers)
	if err != nil {
		return fmt.Errorf("failed to marshal headers: %w", err)
	}

	_, err = r.db.Exec(query,
		subscription.ID,
		subscription.URL,
		subscription.Secret,
		eventsJSON,
		headersJSON,
		subscription.RetryPolicy.MaxRetries,
		subscription.RetryPolicy.BackoffFactor,
		subscription.RetryPolicy.InitialDelay.Milliseconds(),
		subscription.RetryPolicy.MaxBackoff.Milliseconds(),
		subscription.Active,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)

	return err
}

// Update updates an existing webhook subscription
func (r *WebhookSubscriptionRepository) Update(subscription *domain.WebhookSubscription) error {
	query := `
		UPDATE webhook_subscriptions SET
			url = $2,
			secret = $3,
			events = $4,
			headers = $5,
			retry_max_retries = $6,
			retry_backoff_factor = $7,
			retry_initial_delay = $8,
			retry_max_backoff = $9,
			active = $10,
			updated_at = $11
		WHERE id = $1
	`

	eventsJSON, err := json.Marshal(subscription.Events)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	headersJSON, err := json.Marshal(subscription.Headers)
	if err != nil {
		return fmt.Errorf("failed to marshal headers: %w", err)
	}

	result, err := r.db.Exec(query,
		subscription.ID,
		subscription.URL,
		subscription.Secret,
		eventsJSON,
		headersJSON,
		subscription.RetryPolicy.MaxRetries,
		subscription.RetryPolicy.BackoffFactor,
		subscription.RetryPolicy.InitialDelay.Milliseconds(),
		subscription.RetryPolicy.MaxBackoff.Milliseconds(),
		subscription.Active,
		subscription.UpdatedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("webhook subscription not found: %s", subscription.ID)
	}

	return nil
}

// Delete deletes a webhook subscription
func (r *WebhookSubscriptionRepository) Delete(id string) error {
	query := `DELETE FROM webhook_subscriptions WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("webhook subscription not found: %s", id)
	}

	return nil
}

// Find finds a webhook subscription by ID
func (r *WebhookSubscriptionRepository) Find(id string) (*domain.WebhookSubscription, error) {
	query := `
		SELECT id, url, secret, events, headers,
			retry_max_retries, retry_backoff_factor, retry_initial_delay, retry_max_backoff,
			active, created_at, updated_at
		FROM webhook_subscriptions
		WHERE id = $1
	`

	subscription := &domain.WebhookSubscription{}
	var eventsJSON, headersJSON []byte
	var initialDelayMs, maxBackoffMs int64

	err := r.db.QueryRow(query, id).Scan(
		&subscription.ID,
		&subscription.URL,
		&subscription.Secret,
		&eventsJSON,
		&headersJSON,
		&subscription.RetryPolicy.MaxRetries,
		&subscription.RetryPolicy.BackoffFactor,
		&initialDelayMs,
		&maxBackoffMs,
		&subscription.Active,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("webhook subscription not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(eventsJSON, &subscription.Events); err != nil {
		return nil, fmt.Errorf("failed to unmarshal events: %w", err)
	}

	if err := json.Unmarshal(headersJSON, &subscription.Headers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal headers: %w", err)
	}

	// Convert durations
	subscription.RetryPolicy.InitialDelay = time.Duration(initialDelayMs) * time.Millisecond
	subscription.RetryPolicy.MaxBackoff = time.Duration(maxBackoffMs) * time.Millisecond

	return subscription, nil
}

// FindAll finds all webhook subscriptions
func (r *WebhookSubscriptionRepository) FindAll() ([]*domain.WebhookSubscription, error) {
	query := `
		SELECT id, url, secret, events, headers,
			retry_max_retries, retry_backoff_factor, retry_initial_delay, retry_max_backoff,
			active, created_at, updated_at
		FROM webhook_subscriptions
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*domain.WebhookSubscription
	for rows.Next() {
		subscription := &domain.WebhookSubscription{}
		var eventsJSON, headersJSON []byte
		var initialDelayMs, maxBackoffMs int64

		err := rows.Scan(
			&subscription.ID,
			&subscription.URL,
			&subscription.Secret,
			&eventsJSON,
			&headersJSON,
			&subscription.RetryPolicy.MaxRetries,
			&subscription.RetryPolicy.BackoffFactor,
			&initialDelayMs,
			&maxBackoffMs,
			&subscription.Active,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(eventsJSON, &subscription.Events); err != nil {
			return nil, fmt.Errorf("failed to unmarshal events: %w", err)
		}

		if err := json.Unmarshal(headersJSON, &subscription.Headers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal headers: %w", err)
		}

		// Convert durations
		subscription.RetryPolicy.InitialDelay = time.Duration(initialDelayMs) * time.Millisecond
		subscription.RetryPolicy.MaxBackoff = time.Duration(maxBackoffMs) * time.Millisecond

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, rows.Err()
}

// FindActiveByEvent finds active webhook subscriptions for a specific event type
func (r *WebhookSubscriptionRepository) FindActiveByEvent(eventType string) ([]*domain.WebhookSubscription, error) {
	query := `
		SELECT id, url, secret, events, headers,
			retry_max_retries, retry_backoff_factor, retry_initial_delay, retry_max_backoff,
			active, created_at, updated_at
		FROM webhook_subscriptions
		WHERE active = true AND $1 = ANY(string_to_array(events::text, ','))
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, eventType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*domain.WebhookSubscription
	for rows.Next() {
		subscription := &domain.WebhookSubscription{}
		var eventsJSON, headersJSON []byte
		var initialDelayMs, maxBackoffMs int64

		err := rows.Scan(
			&subscription.ID,
			&subscription.URL,
			&subscription.Secret,
			&eventsJSON,
			&headersJSON,
			&subscription.RetryPolicy.MaxRetries,
			&subscription.RetryPolicy.BackoffFactor,
			&initialDelayMs,
			&maxBackoffMs,
			&subscription.Active,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(eventsJSON, &subscription.Events); err != nil {
			return nil, fmt.Errorf("failed to unmarshal events: %w", err)
		}

		if err := json.Unmarshal(headersJSON, &subscription.Headers); err != nil {
			return nil, fmt.Errorf("failed to unmarshal headers: %w", err)
		}

		// Convert durations
		subscription.RetryPolicy.InitialDelay = time.Duration(initialDelayMs) * time.Millisecond
		subscription.RetryPolicy.MaxBackoff = time.Duration(maxBackoffMs) * time.Millisecond

		// Double-check event type in Go (since PostgreSQL JSON query might not be perfect)
		for _, event := range subscription.Events {
			if event == eventType {
				subscriptions = append(subscriptions, subscription)
				break
			}
		}
	}

	return subscriptions, rows.Err()
}

// Ensure interface compliance
var _ domain.WebhookSubscriptionRepository = (*WebhookSubscriptionRepository)(nil)