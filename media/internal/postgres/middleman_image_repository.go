package postgres

import (
	"context"
	"fmt"

	"github.com/stackus/errors" // or any error package you prefer
	"middleman/internal/postgres"
	"middleman/media/internal/domain"
)

type MiddlemanImageRepository struct {
	tableName string
	db        postgres.DB
}

// Ensure MiddlemanImageRepository implements domain.MiddlemanImageRepository
var _ domain.MiddlemanImageRepository = (*MiddlemanImageRepository)(nil)

// Constructor
func NewMiddlemanImageRepository(tableName string, db postgres.DB) MiddlemanImageRepository {
	return MiddlemanImageRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r MiddlemanImageRepository) AddImage(
	ctx context.Context,
	id, mediaID string,
	displayOrder int,
	isMain bool,
	url,
	metadata,
	thumbnail, userID string,
) error {
	// Adjust the column list to match your actual DB schema.
	// Example columns: id, media_id, display_order, is_main, url, metadata, thumbnail, created_at, updated_at
	const query = `
      INSERT INTO %s (
        id,
        media_id,
        display_order,
        is_main,
        url,
        metadata,
        thumbnail,
     	user_id,
        created_at,
        updated_at
      )
      VALUES (
        $1, $2, $3, $4, $5, $6, $7,$8, NOW(), NOW()
      )
    `
	_, err := r.db.ExecContext(ctx, r.table(query),
		id,
		mediaID,
		displayOrder,
		isMain,
		url,
		metadata,
		thumbnail,
		userID,
	)
	if err != nil {
		return errors.Wrap(err, "inserting image")
	}
	return nil
}

// -----------------------------------------------------------------------------
// 2) FindImage
// -----------------------------------------------------------------------------
func (r MiddlemanImageRepository) FindImage(
	ctx context.Context,
	id string,
) (*domain.MiddlemanImage, error) {
	const query = `
      SELECT
        id,
        media_id,
        display_order,
        is_main,
        url,
        metadata,
        thumbnail,
        user_id
      FROM %s
      WHERE id = $1
      LIMIT 1
    `
	row := r.db.QueryRowContext(ctx, r.table(query), id)

	var img domain.MiddlemanImage
	if err := row.Scan(
		&img.ID,
		&img.MediaID,
		&img.DisplayOrder,
		&img.IsMain,
		&img.URL,
		&img.Metadata,
		&img.Thumbnail,
		&img.UserID,
	); err != nil {
		return nil, errors.Wrap(err, "scanning image")
	}
	return &img, nil
}

// -----------------------------------------------------------------------------
// 3) FindAllItemImages
// -----------------------------------------------------------------------------

func (r MiddlemanImageRepository) FindAllItemImages(
	ctx context.Context,
	itemId string,
) ([]*domain.MiddlemanImage, error) {
	const query = `
      SELECT
        id,
        media_id,
        display_order,
        is_main,
        url,
        metadata,
        thumbnail,
        user_id
      FROM %s
      WHERE item_id = $1
      ORDER BY display_order ASC
    `
	rows, err := r.db.QueryContext(ctx, r.table(query), itemId)
	if err != nil {
		return nil, errors.Wrap(err, "querying item images")
	}
	defer rows.Close()

	var images []*domain.MiddlemanImage
	for rows.Next() {
		img := &domain.MiddlemanImage{}
		if err := rows.Scan(
			&img.ID,
			&img.MediaID,
			&img.DisplayOrder,
			&img.IsMain,
			&img.URL,
			&img.Metadata,
			&img.Thumbnail,
			&img.UserID,
		); err != nil {
			return nil, errors.Wrap(err, "scanning item image row")
		}
		images = append(images, img)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing item image rows")
	}
	return images, nil
}
func (r MiddlemanImageRepository) FindAllImagesByAuthor(
	ctx context.Context, userID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.MiddlemanImage, int64, error) {
	// Calculate offset from page and pageSize
	offset := (page - 1) * pageSize
	validSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
	}

	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = $1", r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQ, userID).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting images")
	}
	query := fmt.Sprintf(`
		SELECT
			id,
			media_id,
			display_order,
			is_main,
			url,
			metadata,
			thumbnail,
			user_id,
			created_at,
			updated_at
		FROM %s
		WHERE user_id = $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying author images")
	}
	defer rows.Close()

	var images []*domain.MiddlemanImage
	for rows.Next() {
		vid := &domain.MiddlemanImage{}
		if err := rows.Scan(
			&vid.ID,
			&vid.MediaID,
			&vid.DisplayOrder,
			&vid.IsMain,
			&vid.URL,
			&vid.Metadata,
			&vid.Thumbnail,
			&vid.UserID,
			&vid.CreatedAt,
			&vid.UpdatedAt,
		); err != nil {
			return nil, 0, errors.Wrap(err, "scanning author video row")
		}
		images = append(images, vid)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "finalizing author  video rows")
	}
	return images, totalCount, nil
}

// -----------------------------------------------------------------------------
// 4) FindAllMediaImages
// -----------------------------------------------------------------------------
// Fetch all images that share the same MediaID.
func (r MiddlemanImageRepository) FindAllMediaImages(
	ctx context.Context,
	mediaId string,
) ([]*domain.MiddlemanImage, error) {
	const query = `
      SELECT
        id,
        media_id,
        display_order,
        is_main,
        url,
        metadata,
        thumbnail,
        user_id,
        created_at,
        updated_at
      FROM %s
      WHERE media_id = $1
      ORDER BY display_order ASC
    `
	rows, err := r.db.QueryContext(ctx, r.table(query), mediaId)
	if err != nil {
		return nil, errors.Wrap(err, "querying media images")
	}
	defer rows.Close()

	var images []*domain.MiddlemanImage
	for rows.Next() {
		img := &domain.MiddlemanImage{}
		if err := rows.Scan(
			&img.ID,
			&img.MediaID,
			&img.DisplayOrder,
			&img.IsMain,
			&img.URL,
			&img.Metadata,
			&img.Thumbnail,
			&img.UserID,
			&img.CreatedAt,
			&img.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning media image row")
		}
		images = append(images, img)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing media image rows")
	}
	return images, nil
}

// -----------------------------------------------------------------------------
// 5) FindAllItemTypeImages
// -----------------------------------------------------------------------------
// Similar to `FindAllItemImages`, but the domain interface includes an itemId
// for "type" filtering. If your schema doesn't have "item_type", you may
// replicate the same logic as in `FindAllItemImages` or adapt as needed.
func (r MiddlemanImageRepository) FindAllItemTypeImages(
	ctx context.Context,
	itemId string,
) ([]*domain.MiddlemanImage, error) {
	// Placeholder usage of the same logic:
	const query = `
      SELECT
        id,
        media_id,
        display_order,
        is_main,
        url,
        metadata,
        thumbnail,
        user_id,
        created_at,
        updated_at
      FROM %s
      WHERE item_id = $1
      ORDER BY display_order ASC
    `
	rows, err := r.db.QueryContext(ctx, r.table(query), itemId)
	if err != nil {
		return nil, errors.Wrap(err, "querying item type images")
	}
	defer rows.Close()

	var images []*domain.MiddlemanImage
	for rows.Next() {
		img := &domain.MiddlemanImage{}
		if err := rows.Scan(
			&img.ID,
			&img.MediaID,
			&img.DisplayOrder,
			&img.IsMain,
			&img.URL,
			&img.Metadata,
			&img.Thumbnail,
			&img.UserID,
			&img.CreatedAt,
			&img.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning item type image row")
		}
		images = append(images, img)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing item type image rows")
	}
	return images, nil
}

// -----------------------------------------------------------------------------
// 6) RemoveImage
// -----------------------------------------------------------------------------
func (r MiddlemanImageRepository) RemoveImage(
	ctx context.Context,
	id string,
) error {
	const query = `
      DELETE FROM %s
      WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query), id)
	if err != nil {
		return errors.Wrap(err, "removing image")
	}
	return nil
}

// -----------------------------------------------------------------------------
// Helper to format the table name in queries
// -----------------------------------------------------------------------------
func (r MiddlemanImageRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
