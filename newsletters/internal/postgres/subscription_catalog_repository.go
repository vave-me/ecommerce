package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	
	"middleman/internal/postgres"
	"middleman/newsletters/internal/domain"
	
	"github.com/stackus/errors"
)

type SubscriptionCatalogRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.SubscriptionCatalogRepository = (*SubscriptionCatalogRepository)(nil)

func NewSubscriptionCatalogRepository(tableName string, db postgres.DB) SubscriptionCatalogRepository {
	return SubscriptionCatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r SubscriptionCatalogRepository) Find(ctx context.Context, id string) (*domain.CatalogSubscription, error) {
	const query = `
		SELECT id, user_id, newsletter_id, status, preferences, subscribed_at, unsubscribed_at
		FROM %s
		WHERE id = $1`

	subscription := &domain.CatalogSubscription{}
	var preferences []byte
	
	err := r.db.QueryRowContext(ctx, r.table(query), id).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.NewsletterID,
		&subscription.Status,
		&preferences,
		&subscription.SubscribedAt,
		&subscription.UnsubscribedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "subscription not found")
		}
		return nil, errors.Wrap(err, "scanning subscription")
	}

	// Parse preferences JSON
	if len(preferences) > 0 {
		var prefs map[string]interface{}
		if err := json.Unmarshal(preferences, &prefs); err == nil {
			if freq, ok := prefs["frequency_override"].(string); ok {
				subscription.FrequencyOverride = freq
			}
			if topics, ok := prefs["topics"].([]interface{}); ok {
				for _, topic := range topics {
					if t, ok := topic.(string); ok {
						subscription.Topics = append(subscription.Topics, t)
					}
				}
			}
			if format, ok := prefs["format"].(string); ok {
				subscription.Format = format
			}
		}
	}
	
	return subscription, nil
}

func (r SubscriptionCatalogRepository) FindByUser(ctx context.Context, userID string, status string, limit, offset int) ([]*domain.CatalogSubscription, int, error) {
	query := `
		SELECT id, user_id, newsletter_id, status, preferences, subscribed_at, unsubscribed_at
		FROM %s
		WHERE user_id = $1`
	
	countQuery := `SELECT COUNT(*) FROM %s WHERE user_id = $1`
	
	args := []interface{}{userID}
	
	if status != "" {
		query += " AND status = $2"
		countQuery += " AND status = $2"
		args = append(args, status)
	}
	
	// Add placeholders for limit and offset
	placeholderOffset := len(args) + 1
	query += fmt.Sprintf(" ORDER BY subscribed_at DESC LIMIT $%d OFFSET $%d", placeholderOffset, placeholderOffset+1)

	var total int
	err := r.db.QueryRowContext(ctx, r.table(countQuery), args...).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting subscriptions")
	}

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, r.table(query), args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying subscriptions")
	}
	defer rows.Close()

	var subscriptions []*domain.CatalogSubscription
	for rows.Next() {
		subscription := &domain.CatalogSubscription{}
		var preferences []byte
		
		err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.NewsletterID,
			&subscription.Status,
			&preferences,
			&subscription.SubscribedAt,
			&subscription.UnsubscribedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning subscription")
		}

		// Parse preferences JSON
		if len(preferences) > 0 {
			var prefs map[string]interface{}
			if err := json.Unmarshal(preferences, &prefs); err == nil {
				if freq, ok := prefs["frequency_override"].(string); ok {
					subscription.FrequencyOverride = freq
				}
				if topics, ok := prefs["topics"].([]interface{}); ok {
					for _, topic := range topics {
						if t, ok := topic.(string); ok {
							subscription.Topics = append(subscription.Topics, t)
						}
					}
				}
				if format, ok := prefs["format"].(string); ok {
					subscription.Format = format
				}
			}
		}
		
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, total, nil
}

func (r SubscriptionCatalogRepository) FindByNewsletter(ctx context.Context, newsletterID string, status string, limit, offset int) ([]*domain.CatalogSubscription, int, error) {
	query := `
		SELECT id, user_id, newsletter_id, status, preferences, subscribed_at, unsubscribed_at
		FROM %s
		WHERE newsletter_id = $1`
	
	countQuery := `SELECT COUNT(*) FROM %s WHERE newsletter_id = $1`
	
	args := []interface{}{newsletterID}
	
	if status != "" {
		query += " AND status = $2"
		countQuery += " AND status = $2"
		args = append(args, status)
	}
	
	// Add placeholders for limit and offset
	placeholderOffset := len(args) + 1
	query += fmt.Sprintf(" ORDER BY subscribed_at DESC LIMIT $%d OFFSET $%d", placeholderOffset, placeholderOffset+1)

	var total int
	err := r.db.QueryRowContext(ctx, r.table(countQuery), args...).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting subscriptions")
	}

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, r.table(query), args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying subscriptions")
	}
	defer rows.Close()

	var subscriptions []*domain.CatalogSubscription
	for rows.Next() {
		subscription := &domain.CatalogSubscription{}
		var preferences []byte
		
		err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.NewsletterID,
			&subscription.Status,
			&preferences,
			&subscription.SubscribedAt,
			&subscription.UnsubscribedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(err, "scanning subscription")
		}

		// Parse preferences JSON
		if len(preferences) > 0 {
			var prefs map[string]interface{}
			if err := json.Unmarshal(preferences, &prefs); err == nil {
				if freq, ok := prefs["frequency_override"].(string); ok {
					subscription.FrequencyOverride = freq
				}
				if topics, ok := prefs["topics"].([]interface{}); ok {
					for _, topic := range topics {
						if t, ok := topic.(string); ok {
							subscription.Topics = append(subscription.Topics, t)
						}
					}
				}
				if format, ok := prefs["format"].(string); ok {
					subscription.Format = format
				}
			}
		}
		
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, total, nil
}

func (r SubscriptionCatalogRepository) FindByUserAndNewsletter(ctx context.Context, userID, newsletterID string) (*domain.CatalogSubscription, error) {
	const query = `
		SELECT id, user_id, newsletter_id, status, preferences, subscribed_at, unsubscribed_at
		FROM %s
		WHERE user_id = $1 AND newsletter_id = $2`

	subscription := &domain.CatalogSubscription{}
	var preferences []byte
	
	err := r.db.QueryRowContext(ctx, r.table(query), userID, newsletterID).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.NewsletterID,
		&subscription.Status,
		&preferences,
		&subscription.SubscribedAt,
		&subscription.UnsubscribedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "subscription not found")
		}
		return nil, errors.Wrap(err, "scanning subscription")
	}

	// Parse preferences JSON
	if len(preferences) > 0 {
		var prefs map[string]interface{}
		if err := json.Unmarshal(preferences, &prefs); err == nil {
			if freq, ok := prefs["frequency_override"].(string); ok {
				subscription.FrequencyOverride = freq
			}
			if topics, ok := prefs["topics"].([]interface{}); ok {
				for _, topic := range topics {
					if t, ok := topic.(string); ok {
						subscription.Topics = append(subscription.Topics, t)
					}
				}
			}
			if format, ok := prefs["format"].(string); ok {
				subscription.Format = format
			}
		}
	}
	
	return subscription, nil
}

func (r SubscriptionCatalogRepository) Add(ctx context.Context, subscription *domain.CatalogSubscription) error {
	// Build preferences JSON
	prefs := map[string]interface{}{
		"frequency_override": subscription.FrequencyOverride,
		"topics":            subscription.Topics,
		"format":            subscription.Format,
	}
	preferences, err := json.Marshal(prefs)
	if err != nil {
		return errors.Wrap(err, "marshaling preferences")
	}

	const query = `
		INSERT INTO %s (id, user_id, newsletter_id, status, preferences, subscribed_at, unsubscribed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = r.db.ExecContext(ctx, r.table(query),
		subscription.ID,
		subscription.UserID,
		subscription.NewsletterID,
		subscription.Status,
		preferences,
		subscription.SubscribedAt,
		subscription.UnsubscribedAt,
	)

	return errors.Wrap(err, "inserting subscription")
}

func (r SubscriptionCatalogRepository) Update(ctx context.Context, subscription *domain.CatalogSubscription) error {
	// Build preferences JSON
	prefs := map[string]interface{}{
		"frequency_override": subscription.FrequencyOverride,
		"topics":            subscription.Topics,
		"format":            subscription.Format,
	}
	preferences, err := json.Marshal(prefs)
	if err != nil {
		return errors.Wrap(err, "marshaling preferences")
	}

	const query = `
		UPDATE %s 
		SET status = $2, preferences = $3, unsubscribed_at = $4
		WHERE id = $1`

	_, err = r.db.ExecContext(ctx, r.table(query),
		subscription.ID,
		subscription.Status,
		preferences,
		subscription.UnsubscribedAt,
	)

	return errors.Wrap(err, "updating subscription")
}

func (r SubscriptionCatalogRepository) CountActiveByNewsletter(ctx context.Context, newsletterID string) (int, error) {
	const query = `
		SELECT COUNT(*) 
		FROM %s 
		WHERE newsletter_id = $1 AND status = 'active'`

	var count int
	err := r.db.QueryRowContext(ctx, r.table(query), newsletterID).Scan(&count)
	
	return count, errors.Wrap(err, "counting active subscriptions")
}

func (r SubscriptionCatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}