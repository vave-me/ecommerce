package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/notifications/internal/domain"
)

// CatalogRepository implements domain.CatalogRepository.
type CatalogRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.CatalogRepository = (*CatalogRepository)(nil)

// NewCatalogRepository constructs a CatalogRepository for the given table name and DB interface.
func NewCatalogRepository(tableName string, db postgres.DB) CatalogRepository {
	return CatalogRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r CatalogRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

// Add inserts a new alert row into the 'alerts' table.
func (r CatalogRepository) Add(
	ctx context.Context,
	alertID, userID, alertType, message string,
	payload map[string]interface{},
	isRead bool,
) error {
	const query = `
        INSERT INTO %s (id, user_id, alert_type, message, payload, is_read)
        VALUES ($1, $2, $3, $4, $5, $6)`

	// Serialize the payload from map[string]interface{} into JSON
	// so we can store it in the 'payload' JSON column.
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload into JSON")
	}

	_, err = r.db.ExecContext(
		ctx,
		r.table(query),
		alertID,
		userID,
		alertType,
		message,
		payloadBytes, // store as JSON bytes
		isRead,
	)
	return err
}

// Read updates the 'is_read' status of a given alert row.
func (r CatalogRepository) Read(ctx context.Context, alertID string, isRead bool) error {
	const query = `UPDATE %s SET is_read = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), alertID, isRead)
	return err
}

// Remove deletes an alert row by its ID.
func (r CatalogRepository) Remove(ctx context.Context, alertID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), alertID)
	return err
}

// Find retrieves a single alert by its ID, including unmarshaling JSON 'payload'.
func (r CatalogRepository) Find(ctx context.Context, alertID string) (*domain.MiddlemanAlert, error) {
	const query = `
        SELECT
          id,
          user_id,
          alert_type,
          message,
          payload,
          is_read
        FROM %s
        WHERE id = $1
        LIMIT 1`

	var (
		rawPayload []byte
		id         string
		userID     string
		alertType  string
		message    string
		isRead     bool
	)

	err := r.db.QueryRowContext(ctx, r.table(query), alertID).Scan(
		&id,
		&userID,
		&alertType,
		&message,
		&rawPayload, // we'll manually unmarshal this
		&isRead,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msg(
				"alert with that ID does not exist",
			)
		}
		return nil, errors.Wrap(err, "scanning alert row")
	}

	// Unmarshal the JSON payload
	var payload map[string]interface{}
	if len(rawPayload) > 0 {
		if umErr := json.Unmarshal(rawPayload, &payload); umErr != nil {
			return nil, errors.Wrap(umErr, "unmarshaling alert payload JSON")
		}
	}

	return &domain.MiddlemanAlert{
		ID:        id,
		UserID:    userID,
		AlertType: domain.AlertType(alertType),
		Message:   message,
		Payload:   payload,
		IsRead:    isRead,
		CreatedAt: time.Now(), // TODO: get from database
	}, nil
}

// GetAlerts retrieves multiple alerts for a given userID.
func (r CatalogRepository) GetAlerts(ctx context.Context, userID string, isRead bool) ([]*domain.MiddlemanAlert, error) {
	const query = `
        SELECT
          id,
          user_id,
          alert_type,
          message,
          payload,
          is_read
        FROM %s
        WHERE user_id = $1 AND is_read = $2
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, r.table(query), userID, isRead)
	if err != nil {
		return nil, errors.Wrap(err, "querying alerts")
	}
	defer func() {
		_ = rows.Close()
	}()

	var notifications []*domain.MiddlemanAlert

	for rows.Next() {
		var (
			rawPayload []byte
			id         string
			gotUserID  string
			alertType  string
			message    string
			isRead     bool
		)

		if scErr := rows.Scan(
			&id,
			&gotUserID,
			&alertType,
			&message,
			&rawPayload,
			&isRead,
		); scErr != nil {
			return nil, errors.Wrap(scErr, "scanning alert row")
		}

		var payload map[string]interface{}
		if len(rawPayload) > 0 {
			if umErr := json.Unmarshal(rawPayload, &payload); umErr != nil {
				return nil, errors.Wrap(umErr, "unmarshaling alert payload JSON")
			}
		}

		notifications = append(notifications, &domain.MiddlemanAlert{
			ID:        id,
			UserID:    gotUserID,
			AlertType: domain.AlertType(alertType),
			Message:   message,
			Payload:   payload,
			IsRead:    isRead,
			CreatedAt: time.Now(), // TODO: get from database
		})
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing alert rows")
	}

	return notifications, nil
}

// GetAlertsByType retrieves alerts by userID and alert_type.
func (r CatalogRepository) GetAlertsByType(ctx context.Context, userID, notificationType string, isRead bool) ([]*domain.MiddlemanAlert, error) {
	const query = `
        SELECT
          id,
          user_id,
          alert_type,
          message,
          payload,
          is_read
        FROM %s
        WHERE user_id = $1
          AND alert_type = $2
        AND is_read = $3
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, r.table(query), userID, notificationType, isRead)
	if err != nil {
		return nil, errors.Wrap(err, "querying alerts by type")
	}
	defer func() {
		_ = rows.Close()
	}()

	var notifications []*domain.MiddlemanAlert

	for rows.Next() {
		var (
			rawPayload []byte
			id         string
			userID     string
			alertType  string
			message    string
			isRead     bool
		)

		if scErr := rows.Scan(
			&id,
			&userID,
			&alertType,
			&message,
			&rawPayload,
			&isRead,
		); scErr != nil {
			return nil, errors.Wrap(scErr, "scanning alert row")
		}

		var payload map[string]interface{}
		if len(rawPayload) > 0 {
			if umErr := json.Unmarshal(rawPayload, &payload); umErr != nil {
				return nil, errors.Wrap(umErr, "unmarshaling alert payload JSON")
			}
		}

		notifications = append(notifications, &domain.MiddlemanAlert{
			ID:        id,
			UserID:    userID,
			AlertType: domain.AlertType(alertType),
			Message:   message,
			Payload:   payload,
			IsRead:    isRead,
			CreatedAt: time.Now(), // TODO: get from database
		})
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing alert rows")
	}

	return notifications, nil
}
