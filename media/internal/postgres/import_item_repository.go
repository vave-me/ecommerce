package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/media/internal/domain"
)

// ImportItemRepository satisfies domain.ImportItemRepository
type ImportItemRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.ImportItemRepository = (*ImportItemRepository)(nil)

// Constructor
func NewImportItemRepository(tableName string, db postgres.DB) ImportItemRepository {
	return ImportItemRepository{
		tableName: tableName,
		db:        db,
	}
}

// CreateBatch inserts multiple import items
func (r ImportItemRepository) CreateBatch(ctx context.Context, items []*domain.ImportItem) error {
	if len(items) == 0 {
		return nil
	}

	const query = `
      INSERT INTO %s (
        id, 
        session_id, 
        external_id, 
        sku, 
        product_id,
        image_url, 
        status, 
        error_message, 
        retry_count, 
        media_id, 
        image_id, 
        processed_at,
        metadata,
        display_order
      )
      VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
      )
    `

	for _, item := range items {
		if item.ID == "" {
			item.ID = uuid.New().String()
		}

		var mediaID, imageID, productID sql.NullString
		var processedAt sql.NullTime
		var metadata []byte

		if item.MediaID != "" {
			mediaID = sql.NullString{String: item.MediaID, Valid: true}
		}
		if item.ImageID != "" {
			imageID = sql.NullString{String: item.ImageID, Valid: true}
		}
		if item.ProductID != "" {
			productID = sql.NullString{String: item.ProductID, Valid: true}
		}
		if !item.ProcessedAt.IsZero() {
			processedAt = sql.NullTime{Time: item.ProcessedAt, Valid: true}
		}
		
		// Marshal metadata
		if item.Metadata != nil {
			metadata, err = json.Marshal(item.Metadata)
			if err != nil {
				return errors.Wrap(err, "marshaling metadata")
			}
		}

		_, err := r.db.ExecContext(ctx, r.table(query),
			item.ID,
			item.SessionID,
			item.ExternalID,
			item.SKU,
			productID,
			item.ImageURL,
			item.Status,
			item.ErrorMessage,
			item.RetryCount,
			mediaID,
			imageID,
			processedAt,
			metadata,
			item.DisplayOrder,
		)
		if err != nil {
			return errors.Wrap(err, "inserting import item")
		}
	}

	return nil
}

// GetBySession retrieves all import items for a session with a specific status
func (r ImportItemRepository) GetBySession(ctx context.Context, sessionID string, status domain.ImportItemStatus) ([]*domain.ImportItem, error) {
	const query = `
      SELECT 
        id, 
        session_id, 
        external_id, 
        sku, 
        product_id,
        image_url, 
        status, 
        error_message, 
        retry_count, 
        media_id, 
        image_id, 
        processed_at,
        metadata,
        display_order
      FROM %s
      WHERE session_id = $1 AND status = $2
      ORDER BY id
    `

	rows, err := r.db.QueryContext(ctx, r.table(query), sessionID, status)
	if err != nil {
		return nil, errors.Wrap(err, "querying import items")
	}
	defer rows.Close()

	var items []*domain.ImportItem
	for rows.Next() {
		item := &domain.ImportItem{}
		var mediaID, imageID, productID sql.NullString
		var processedAt sql.NullTime
		var metadata []byte

		if err := rows.Scan(
			&item.ID,
			&item.SessionID,
			&item.ExternalID,
			&item.SKU,
			&productID,
			&item.ImageURL,
			&item.Status,
			&item.ErrorMessage,
			&item.RetryCount,
			&mediaID,
			&imageID,
			&processedAt,
			&metadata,
			&item.DisplayOrder,
		); err != nil {
			return nil, errors.Wrap(err, "scanning import item")
		}

		if mediaID.Valid {
			item.MediaID = mediaID.String
		}
		if imageID.Valid {
			item.ImageID = imageID.String
		}
		if productID.Valid {
			item.ProductID = productID.String
		}
		if processedAt.Valid {
			item.ProcessedAt = processedAt.Time
		}
		
		// Unmarshal metadata
		if len(metadata) > 0 {
			if err := json.Unmarshal(metadata, &item.Metadata); err != nil {
				return nil, errors.Wrap(err, "unmarshaling metadata")
			}
		} else {
			item.Metadata = make(map[string]string)
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing import item rows")
	}

	return items, nil
}

// Update modifies an existing import item
func (r ImportItemRepository) Update(ctx context.Context, item *domain.ImportItem) error {
	const query = `
      UPDATE %s 
      SET 
        status = $2, 
        error_message = $3, 
        retry_count = $4, 
        media_id = $5, 
        image_id = $6, 
        processed_at = $7,
        product_id = $8
      WHERE id = $1
    `

	var mediaID, imageID, productID sql.NullString
	var processedAt sql.NullTime

	if item.MediaID != "" {
		mediaID = sql.NullString{String: item.MediaID, Valid: true}
	}
	if item.ImageID != "" {
		imageID = sql.NullString{String: item.ImageID, Valid: true}
	}
	if item.ProductID != "" {
		productID = sql.NullString{String: item.ProductID, Valid: true}
	}
	if !item.ProcessedAt.IsZero() {
		processedAt = sql.NullTime{Time: item.ProcessedAt, Valid: true}
	}

	result, err := r.db.ExecContext(ctx, r.table(query),
		item.ID,
		item.Status,
		item.ErrorMessage,
		item.RetryCount,
		mediaID,
		imageID,
		processedAt,
		productID,
	)
	if err != nil {
		return errors.Wrap(err, "updating import item")
	}

	// Check if any row was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "checking rows affected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(errors.ErrNotFound, "import item not found")
	}

	return nil
}

// GetPendingItems retrieves pending import items with a limit, using row locking for concurrency
func (r ImportItemRepository) GetPendingItems(ctx context.Context, limit int) ([]*domain.ImportItem, error) {
	const query = `
      SELECT 
        id, 
        session_id, 
        external_id, 
        sku, 
        product_id,
        image_url, 
        status, 
        error_message, 
        retry_count, 
        media_id, 
        image_id, 
        processed_at,
        metadata,
        display_order
      FROM %s
      WHERE status = $1
      ORDER BY id
      LIMIT $2
      FOR UPDATE SKIP LOCKED
    `

	rows, err := r.db.QueryContext(ctx, r.table(query), domain.ImportItemStatusPending, limit)
	if err != nil {
		return nil, errors.Wrap(err, "querying pending items")
	}
	defer rows.Close()

	var items []*domain.ImportItem
	for rows.Next() {
		item := &domain.ImportItem{}
		var mediaID, imageID, productID sql.NullString
		var processedAt sql.NullTime
		var metadata []byte

		if err := rows.Scan(
			&item.ID,
			&item.SessionID,
			&item.ExternalID,
			&item.SKU,
			&productID,
			&item.ImageURL,
			&item.Status,
			&item.ErrorMessage,
			&item.RetryCount,
			&mediaID,
			&imageID,
			&processedAt,
			&metadata,
			&item.DisplayOrder,
		); err != nil {
			return nil, errors.Wrap(err, "scanning pending item")
		}

		if mediaID.Valid {
			item.MediaID = mediaID.String
		}
		if imageID.Valid {
			item.ImageID = imageID.String
		}
		if productID.Valid {
			item.ProductID = productID.String
		}
		if processedAt.Valid {
			item.ProcessedAt = processedAt.Time
		}
		
		// Unmarshal metadata
		if len(metadata) > 0 {
			if err := json.Unmarshal(metadata, &item.Metadata); err != nil {
				return nil, errors.Wrap(err, "unmarshaling metadata")
			}
		} else {
			item.Metadata = make(map[string]string)
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing pending item rows")
	}

	return items, nil
}

// Helper to format queries
func (r ImportItemRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
