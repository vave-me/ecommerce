package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stackus/errors"
	"middleman/comments/internal/domain"
	"middleman/internal/postgres"
)

type MiddlemanRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.MiddlemanRepository = (*MiddlemanRepository)(nil)

func NewMiddlemanRepository(tableName string, db postgres.DB) MiddlemanRepository {
	return MiddlemanRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r MiddlemanRepository) Add(ctx context.Context, commentID, senderID, itemID string, itemType domain.ItemType, content, categoryID, parentID string) error {
	const query = `
INSERT INTO %s (id, sender_id, item_id, item_type, content, category_id, parent_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
`
	_, err := r.db.ExecContext(ctx, r.table(query),
		commentID, // $1 -> id
		senderID,  // $2 -> sender_id
		itemID,    // $3 -> item_id
		itemType.String(),
		content,
		categoryID,
		parentID,
	)

	return err
}

func (r MiddlemanRepository) Find(ctx context.Context, commentID, itemID string) (*domain.MiddlemanComment, error) {
	// If you want to ensure the comment belongs to itemID, do "WHERE id=$1 AND item_id=$2"
	const query = `
SELECT id, sender_id, item_id, item_type, content, category_id, parent_id, approved, flagged, created_at, updated_at
  FROM %s
 WHERE id = $1
   AND item_id = $2
 LIMIT 1
`
	reply := &domain.MiddlemanComment{}

	err := r.db.QueryRowContext(ctx, r.table(query), commentID, itemID).Scan(
		&reply.ID,
		&reply.SenderID,
		&reply.ItemID,
		&reply.ItemType,
		&reply.Content,
		&reply.CategoryID,
		&reply.ParentID,
		&reply.Approved,
		&reply.Flagged,
		&reply.CreatedAt,
		&reply.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, "scanning comment")
	}

	return reply, nil
}

func (r MiddlemanRepository) All(ctx context.Context, itemID string) ([]*domain.MiddlemanComment, error) {
	const query = `
SELECT id, sender_id, item_id, item_type, content, category_id, parent_id, approved, flagged, created_at, updated_at
  FROM %s
 WHERE item_id = $1
 ORDER BY created_at
`
	rows, err := r.db.QueryContext(ctx, r.table(query), itemID)
	if err != nil {
		return nil, errors.Wrap(err, "querying comments by item_id")
	}
	defer rows.Close()

	var comments []*domain.MiddlemanComment
	for rows.Next() {
		comment := &domain.MiddlemanComment{}
		if err := rows.Scan(
			&comment.ID,
			&comment.SenderID,
			&comment.ItemID,
			&comment.ItemType,
			&comment.Content,
			&comment.CategoryID,
			&comment.ParentID,
			&comment.Approved,
			&comment.Flagged,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning all comments")
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing rows for All()")
	}

	return comments, nil
}

func (r MiddlemanRepository) FindBySenderID(ctx context.Context, senderID string) ([]*domain.MiddlemanComment, error) {
	const query = `
SELECT id, sender_id, item_id, item_type, content, category_id, parent_id, approved, flagged, created_at, updated_at
  FROM %s
 WHERE sender_id = $1
 ORDER BY created_at
`
	rows, err := r.db.QueryContext(ctx, r.table(query), senderID)
	if err != nil {
		return nil, errors.Wrap(err, "querying comments by sender_id")
	}
	defer rows.Close()

	var comments []*domain.MiddlemanComment
	for rows.Next() {
		c := &domain.MiddlemanComment{}
		if err := rows.Scan(
			&c.ID,
			&c.SenderID,
			&c.ItemID,
			&c.ItemType,
			&c.Content,
			&c.CategoryID,
			&c.ParentID,
			&c.Approved,
			&c.Flagged,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning comments for senderID")
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing rows for FindBySenderID()")
	}
	return comments, nil
}

func (r MiddlemanRepository) Remove(ctx context.Context, commentID string) error {
	// If `id` is truly unique, you only need the commentID.
	// If you require (commentID+itemID) as the key, you can add `AND item_id = $2`.
	const query = `DELETE FROM %s WHERE id = $1`
	_, err := r.db.ExecContext(ctx, r.table(query), commentID)
	return err
}

func (r MiddlemanRepository) MostCommentedItems(
	ctx context.Context,
	limit, offset int,
) ([]*domain.ItemCommentCount, error) {

	// This query only returns rows with category_id = '' (i.e. "no category")
	baseQuery := `
    SELECT item_id, item_type, category_id, total_comments
      FROM item_comment_counts
     WHERE category_id = ''
     ORDER BY total_comments DESC
    `
	finalQuery, args := buildLimitOffsetQuery(baseQuery, limit, offset)

	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying most commented items (all categories=empty)")
	}
	defer rows.Close()

	var results []*domain.ItemCommentCount
	for rows.Next() {
		p := &domain.ItemCommentCount{}
		if err := rows.Scan(
			&p.ItemID,
			&p.ItemType,
			&p.CategoryID,
			&p.CommentsCount,
		); err != nil {
			return nil, errors.Wrap(err, "scanning item_comment_counts")
		}
		results = append(results, p)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing aggregator rows")
	}
	return results, nil
}

func (r MiddlemanRepository) MostCommentedItemsByCategory(
	ctx context.Context,
	itemType domain.ItemType,
	categoryID string,
	limit, offset int,
) ([]*domain.ItemCommentCount, error) {

	baseQuery := `
    SELECT item_id, item_type, category_id, total_comments
      FROM item_comment_counts
     WHERE category_id = $1
       AND item_type = $2
     ORDER BY total_comments DESC
    `
	// We'll pass (categoryID, itemType) first, then limit/offset
	finalQuery, args := buildLimitOffsetQuery(baseQuery, limit, offset, categoryID, itemType.String())

	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying most commented items by category/type")
	}
	defer rows.Close()

	var results []*domain.ItemCommentCount
	for rows.Next() {
		p := &domain.ItemCommentCount{}
		if err := rows.Scan(
			&p.ItemID,
			&p.ItemType,
			&p.CategoryID,
			&p.CommentsCount,
		); err != nil {
			return nil, errors.Wrap(err, "scanning item_comment_counts by category/type")
		}
		results = append(results, p)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing aggregator rows by category/type")
	}
	return results, nil
}

// A small helper to append LIMIT and OFFSET if > 0
func buildLimitOffsetQuery(base string, limit, offset int, argsIn ...interface{}) (string, []interface{}) {
	argPos := len(argsIn) + 1
	finalQuery := base
	finalArgs := argsIn

	if limit > 0 {
		finalQuery += fmt.Sprintf(" LIMIT $%d", argPos)
		finalArgs = append(finalArgs, limit)
		argPos++
	}
	if offset > 0 {
		finalQuery += fmt.Sprintf(" OFFSET $%d", argPos)
		finalArgs = append(finalArgs, offset)
		argPos++
	}
	return finalQuery, finalArgs
}

func (r MiddlemanRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
