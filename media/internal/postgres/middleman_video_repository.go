package postgres

import (
	"context"
	"fmt"

	"github.com/stackus/errors" // or your error package of choice
	"middleman/internal/postgres"
	"middleman/media/internal/domain"
)

type MiddlemanVideoRepository struct {
	tableName string
	db        postgres.DB
}

// Ensure MiddlemanVideoRepository implements the domain interface
var _ domain.MiddlemanVideoRepository = (*MiddlemanVideoRepository)(nil)

// Constructor
func NewMiddlemanVideoRepository(tableName string, db postgres.DB) MiddlemanVideoRepository {
	return MiddlemanVideoRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r MiddlemanVideoRepository) AddVideo(
	ctx context.Context,
	id, mediaID string,
	displayOrder int,
	isMain bool,
	url string,
	metadata string,
	thumbnail string,
	userID string,
) error {

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
		return errors.Wrap(err, "inserting video")
	}
	return nil
}

func (r MiddlemanVideoRepository) FindVideo(
	ctx context.Context,
	id string,
) (*domain.MiddlemanVideo, error) {
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
      WHERE id = $1
      LIMIT 1
    `
	row := r.db.QueryRowContext(ctx, r.table(query), id)

	vid := &domain.MiddlemanVideo{
		ID: id,
	}
	if err := row.Scan(
		&vid.ID,
		&vid.MediaID,
		&vid.DisplayOrder,
		&vid.IsMain,
		&vid.Url,
		&vid.Metadata,
		&vid.Thumbnail,
		&vid.UserID,
		&vid.CreatedAt,
		&vid.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "scanning video")
	}
	return vid, nil
}

func (r MiddlemanVideoRepository) FindAllItemVideos(
	ctx context.Context,
	itemId string,
) ([]*domain.MiddlemanVideo, error) {
	// Adjust if your table uses a different column name or requires a JOIN.
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
		return nil, errors.Wrap(err, "querying item videos")
	}
	defer rows.Close()

	var videos []*domain.MiddlemanVideo
	for rows.Next() {
		vid := &domain.MiddlemanVideo{}
		if err := rows.Scan(
			&vid.ID,
			&vid.MediaID,
			&vid.DisplayOrder,
			&vid.IsMain,
			&vid.Url,
			&vid.Metadata,
			&vid.Thumbnail,
			&vid.UserID,
			&vid.CreatedAt,
			&vid.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning item video row")
		}
		videos = append(videos, vid)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing item video rows")
	}
	return videos, nil
}

func (r MiddlemanVideoRepository) FindAllVideosByAuthor(
	ctx context.Context, userID string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.MiddlemanVideo, int64, error) {
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
		return nil, 0, errors.Wrap(err, "counting products")
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
		return nil, 0, errors.Wrap(err, "querying author videos")
	}
	defer rows.Close()

	var videos []*domain.MiddlemanVideo
	for rows.Next() {
		vid := &domain.MiddlemanVideo{}
		if err := rows.Scan(
			&vid.ID,
			&vid.MediaID,
			&vid.DisplayOrder,
			&vid.IsMain,
			&vid.Url,
			&vid.Metadata,
			&vid.Thumbnail,
			&vid.UserID,
			&vid.CreatedAt,
			&vid.UpdatedAt,
		); err != nil {
			return nil, 0, errors.Wrap(err, "scanning author video row")
		}
		videos = append(videos, vid)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "finalizing author  video rows")
	}
	return videos, totalCount, nil
}

// -----------------------------------------------------------------------------
// 4) FindAllMediaVideos
// -----------------------------------------------------------------------------
// Fetch all videos for a given 'mediaId'. This matches the domain struct field
// 'MediaID'.
func (r MiddlemanVideoRepository) FindAllMediaVideos(
	ctx context.Context,
	mediaId string,
) ([]*domain.MiddlemanVideo, error) {
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
		return nil, errors.Wrap(err, "querying media videos")
	}
	defer rows.Close()

	var videos []*domain.MiddlemanVideo
	for rows.Next() {
		vid := &domain.MiddlemanVideo{}
		if err := rows.Scan(
			&vid.ID,
			&vid.MediaID,
			&vid.DisplayOrder,
			&vid.IsMain,
			&vid.Url,
			&vid.Metadata,
			&vid.Thumbnail,
			&vid.UserID,
			&vid.CreatedAt,
			&vid.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning media video row")
		}
		videos = append(videos, vid)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing media video rows")
	}
	return videos, nil
}

func (r MiddlemanVideoRepository) FindAllVideos(
	ctx context.Context,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*domain.MiddlemanVideo, int64, error) {
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
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM %s", r.tableName)
	var totalCount int64
	if err := r.db.QueryRowContext(ctx, countQ).Scan(&totalCount); err != nil {
		return nil, 0, errors.Wrap(err, "counting products")
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
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, r.tableName, sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying media videos")
	}
	defer rows.Close()

	var videos []*domain.MiddlemanVideo
	for rows.Next() {
		vid := &domain.MiddlemanVideo{}
		if err := rows.Scan(
			&vid.ID,
			&vid.MediaID,
			&vid.DisplayOrder,
			&vid.IsMain,
			&vid.Url,
			&vid.Metadata,
			&vid.Thumbnail,
			&vid.UserID,
			&vid.CreatedAt,
			&vid.UpdatedAt,
		); err != nil {
			return nil, 0, errors.Wrap(err, "scanning media video row")
		}
		videos = append(videos, vid)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "finalizing media video rows")
	}
	return videos, totalCount, nil
}

// -----------------------------------------------------------------------------
// 5) FindAllItemTypeVideos
// -----------------------------------------------------------------------------
// This signature is the same as FindAllItemVideos in the domain. Possibly
// this is meant to filter by item type or do something else. The domain struct
// has no 'ItemType', so we'll assume there's a column or a relation. For now,
// this is just a placeholder that duplicates the item logic or uses a different
// filter. Adjust as needed for your real schema.
func (r MiddlemanVideoRepository) FindAllItemTypeVideos(
	ctx context.Context,
	itemId string,
) ([]*domain.MiddlemanVideo, error) {
	// Example placeholder that simply re-queries by item_id.
	// Modify if you have a separate 'item_type' or join logic.
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
		return nil, errors.Wrap(err, "querying item type videos")
	}
	defer rows.Close()

	var videos []*domain.MiddlemanVideo
	for rows.Next() {
		vid := &domain.MiddlemanVideo{}
		if err := rows.Scan(
			&vid.ID,
			&vid.MediaID,
			&vid.DisplayOrder,
			&vid.IsMain,
			&vid.Url,
			&vid.Metadata,
			&vid.Thumbnail,
			&vid.UserID,
			&vid.CreatedAt,
			&vid.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning item type video row")
		}
		videos = append(videos, vid)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finalizing item type video rows")
	}
	return videos, nil
}

// -----------------------------------------------------------------------------
// 6) RemoveVideo
// -----------------------------------------------------------------------------
func (r MiddlemanVideoRepository) RemoveVideo(
	ctx context.Context,
	id string,
) error {
	const query = `
      DELETE FROM %s
      WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, r.table(query), id)
	if err != nil {
		return errors.Wrap(err, "removing video")
	}
	return nil
}

// -----------------------------------------------------------------------------
// Helper to format the table name in queries
// -----------------------------------------------------------------------------
func (r MiddlemanVideoRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
