package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/media/internal/domain"
)

// MiddlemanMediaRepository satisfies domain.MiddlemanMediaRepository
type MiddlemanMediaRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.MiddlemanMediaRepository = (*MiddlemanMediaRepository)(nil)

// Constructor
func NewMiddlemanMediaRepository(tableName string, db postgres.DB) MiddlemanMediaRepository {
	return MiddlemanMediaRepository{
		tableName: tableName,
		db:        db,
	}
}

// 1) AddMedia
func (r MiddlemanMediaRepository) AddMedia(
	ctx context.Context,
	id, itemID string,
	itemType domain.ItemType,
	userID string,
	status domain.MediaStatus,
) error {
	const query = `
      INSERT INTO %s (
        id,
        item_id,
        item_type,
        user_id,
        status,
        created_at,
        updated_at
      )
      VALUES (
        $1, $2, $3, $4, $5, NOW(), NOW()
      )
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		id,
		itemID,
		itemType,
		userID,
		status,
	)
	if err != nil {
		return errors.Wrap(err, "inserting media")
	}
	return nil
}

// 2) GetMedia
func (r MiddlemanMediaRepository) GetMedia(
	ctx context.Context,
	id string,
) (*domain.MiddlemanMedia, error) {
	const query = `
      SELECT
        id,
        item_id,
        item_type,
        user_id,
        status,
        media_order
      FROM %s
      WHERE id = $1
      LIMIT 1
    `
	row := r.db.QueryRowContext(ctx, r.table(query), id)

	m := &domain.MiddlemanMedia{}
	var rawOrder sql.NullString

	if err := row.Scan(
		&m.ID,
		&m.ItemID,
		&m.ItemType,
		&m.UserID,
		&m.Status,
		&rawOrder,
	); err != nil {
		return m, errors.Wrap(err, "scanning media")
	}

	// If 'media_order' is not empty, parse it as map[int]MediaOrderItem
	if rawOrder.Valid && rawOrder.String != "" {
		var parsed map[int]domain.MediaOrderItem
		if err := json.Unmarshal([]byte(rawOrder.String), &parsed); err != nil {
			return m, errors.Wrap(err, "unmarshal media_order JSON")
		}
		m.MediaOrder = parsed
	} else {
		// If no JSON, just assign an empty map
		m.MediaOrder = make(map[int]domain.MediaOrderItem)
	}

	return m, nil
}

// 3) GetMediaByItem
func (r MiddlemanMediaRepository) GetMediaByItem(
	ctx context.Context,
	itemID string,
) (*domain.MiddlemanMedia, error) {
	const query = `
      SELECT
        id,
        item_id,
        item_type,
        user_id,
        status,
        media_order
      FROM %s
      WHERE item_id = $1
      LIMIT 1
    `
	row := r.db.QueryRowContext(ctx, r.table(query), itemID)

	m := &domain.MiddlemanMedia{}
	var rawOrder sql.NullString

	if err := row.Scan(
		&m.ID,
		&m.ItemID,
		&m.ItemType,
		&m.UserID,
		&m.Status,
		&rawOrder,
	); err != nil {
		return m, errors.Wrap(err, "scanning media")
	}

	if rawOrder.Valid && rawOrder.String != "" {
		var parsed map[int]domain.MediaOrderItem
		if err := json.Unmarshal([]byte(rawOrder.String), &parsed); err != nil {
			return m, errors.Wrap(err, "unmarshal media_order JSON")
		}
		m.MediaOrder = parsed
	} else {
		m.MediaOrder = make(map[int]domain.MediaOrderItem)
	}

	return m, nil
}

// 4) GetAllUserMedia
func (r MiddlemanMediaRepository) GetAllUserMedia(
	ctx context.Context,
	userID string,
) ([]*domain.MiddlemanMedia, error) {
	const query = `
      SELECT
        id,
        item_id,
        item_type,
        user_id,
        status,
        media_order
      FROM %s
      WHERE user_id = $1
    `
	rows, err := r.db.QueryContext(ctx, r.table(query), userID)
	if err != nil {
		return nil, errors.Wrap(err, "querying user media")
	}
	defer rows.Close()

	var medias []*domain.MiddlemanMedia
	for rows.Next() {
		m := &domain.MiddlemanMedia{}
		var rawOrder sql.NullString

		if err := rows.Scan(
			&m.ID,
			&m.ItemID,
			&m.ItemType,
			&m.UserID,
			&m.Status,
			&rawOrder,
		); err != nil {
			return nil, errors.Wrap(err, "scanning user media row")
		}

		if rawOrder.Valid && rawOrder.String != "" {
			var parsed map[int]domain.MediaOrderItem
			if err := json.Unmarshal([]byte(rawOrder.String), &parsed); err != nil {
				return nil, errors.Wrap(err, "unmarshal media_order JSON in loop")
			}
			m.MediaOrder = parsed
		} else {
			m.MediaOrder = make(map[int]domain.MediaOrderItem)
		}

		medias = append(medias, m)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing user media rows")
	}

	return medias, nil
}
func (r MiddlemanMediaRepository) UpdateMedia(
	ctx context.Context,
	id, itemID string,
	itemType domain.ItemType,
	userID string,
	status domain.MediaStatus,
) error {
	const query = `
      UPDATE %s
      SET 
        item_id = $2,
        item_type = $3,
        user_id = $4,
        status = $5,
        updated_at = NOW()
      WHERE id = $1
    `
	result, err := r.db.ExecContext(ctx, r.table(query),
		id,
		itemID,
		itemType,
		userID,
		status,
	)
	if err != nil {
		return errors.Wrap(err, "updating media")
	}

	// Check if any row was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "checking rows affected")
	}

	if rowsAffected == 0 {
		return errors.ErrOK
	}

	return nil
}

// 5) GetItemMedia
func (r MiddlemanMediaRepository) GetItemMedia(
	ctx context.Context,
	id string,
) (*domain.MiddlemanMedia, error) {
	const query = `
      SELECT
        id,
        item_id,
        item_type,
        user_id,
        status,
        media_order
      FROM %s
      WHERE item_id = $1
      LIMIT 1
    `
	row := r.db.QueryRowContext(ctx, r.table(query), id)

	m := &domain.MiddlemanMedia{}
	var rawOrder sql.NullString

	if err := row.Scan(
		&m.ID,
		&m.ItemID,
		&m.ItemType,
		&m.UserID,
		&m.Status,
		&rawOrder,
	); err != nil {
		return m, errors.Wrap(err, "scanning item media")
	}

	if rawOrder.Valid && rawOrder.String != "" {
		var parsed map[int]domain.MediaOrderItem
		if err := json.Unmarshal([]byte(rawOrder.String), &parsed); err != nil {
			return m, errors.Wrap(err, "unmarshal media_order JSON for item media")
		}
		m.MediaOrder = parsed
	} else {
		m.MediaOrder = make(map[int]domain.MediaOrderItem)
	}

	return m, nil
}

// 6) RemoveMedia
func (r MiddlemanMediaRepository) RemoveMedia(
	ctx context.Context,
	id string,
) error {
	const query = `
      DELETE FROM %s
      WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query), id)
	if err != nil {
		return errors.Wrap(err, "removing media")
	}
	return nil
}

// 7) AddMediaItemOrder
func (r MiddlemanMediaRepository) AddMediaItemOrder(
	ctx context.Context,
	id, mediaItemId, mediaItemUrl string,
	displayOrder int,
) error {
	// 1) Load existing JSON
	const selectQuery = `
      SELECT media_order 
      FROM %s
      WHERE id = $1
      LIMIT 1
    `
	var rawOrder sql.NullString
	err := r.db.QueryRowContext(ctx, r.table(selectQuery), id).Scan(&rawOrder)
	if err != nil {
		return errors.Wrap(err, "retrieving media_order")
	}

	// 2) Parse into map[int]MediaOrderItem
	existing := make(map[int]domain.MediaOrderItem)
	if rawOrder.Valid && rawOrder.String != "" {
		if unmarshalErr := json.Unmarshal([]byte(rawOrder.String), &existing); unmarshalErr != nil {
			return errors.Wrap(unmarshalErr, "unmarshalling media_order JSON")
		}
	}

	// Remove if the item already exists
	for k, v := range existing {
		if v.MediaItemID == mediaItemId {
			delete(existing, k)
			break
		}
	}

	// SHIFT occupant if needed
	shiftOccupant(existing, displayOrder)

	// Insert new occupant
	existing[displayOrder] = domain.MediaOrderItem{
		MediaItemID: mediaItemId,
		URL:         mediaItemUrl,
	}

	// 3) Convert back to JSON
	newJSON, marshalErr := json.Marshal(existing)
	if marshalErr != nil {
		return errors.Wrap(marshalErr, "marshalling updated media_order JSON")
	}

	// 4) Update DB
	const updateQuery = `
      UPDATE %s
      SET media_order = $2,
          updated_at = NOW()
      WHERE id = $1
    `
	_, err = r.db.ExecContext(ctx, r.table(updateQuery), id, string(newJSON))
	if err != nil {
		return errors.Wrap(err, "updating media_order")
	}

	return nil
}

// shiftOccupant is a top-level named function to handle occupant shifting
func shiftOccupant(
	orders map[int]domain.MediaOrderItem,
	order int,
) {
	occupant, found := orders[order]
	if !found {
		return
	}
	next := order + 1
	// Recursively shift occupant at next
	shiftOccupant(orders, next)

	orders[next] = occupant
	delete(orders, order)
}

// Helper to format queries
func (r MiddlemanMediaRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
