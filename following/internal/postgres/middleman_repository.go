package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stackus/errors"
	"middleman/following/internal/domain"
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

func (r MiddlemanRepository) Add(ctx context.Context, followID, userID, followedUserID string, followedUserType domain.FollowedUserType, content, categoryID, parentID string) error {
	const query = `
INSERT INTO %s (id, user_id, followed_user_id, followed_user_type, content, category_id, parent_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
`
	_, err := r.db.ExecContext(ctx, r.table(query),
		followID,       // $1 -> id
		userID,         // $2 -> user_id
		followedUserID, // $3 -> followed_user_id
		followedUserType.String(),
		content,
		categoryID,
		parentID,
	)

	return err
}

func (r MiddlemanRepository) Find(ctx context.Context, followID, followedUserID string) (*domain.MiddlemanFollow, error) {
	// If you want to ensure the follow belongs to followedUserID, do "WHERE id=$1 AND followed_user_id=$2"
	const query = `
SELECT id, user_id, followed_user_id, followed_user_type, content, category_id, parent_id, approved, flagged, created_at, updated_at
  FROM %s
 WHERE id = $1
   AND followed_user_id = $2
 LIMIT 1
`
	reply := &domain.MiddlemanFollow{}

	err := r.db.QueryRowContext(ctx, r.table(query), followID, followedUserID).Scan(
		&reply.ID,
		&reply.UserID,
		&reply.FollowedUserID,
		&reply.FollowedUserType,
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
		return nil, errors.Wrap(err, "scanning follow")
	}

	return reply, nil
}

func (r MiddlemanRepository) All(ctx context.Context, followedUserID string) ([]*domain.MiddlemanFollow, error) {
	const query = `
SELECT id, user_id, followed_user_id, followed_user_type, content, category_id, parent_id, approved, flagged, created_at, updated_at
  FROM %s
 WHERE followed_user_id = $1
 ORDER BY created_at
`
	rows, err := r.db.QueryContext(ctx, r.table(query), followedUserID)
	if err != nil {
		return nil, errors.Wrap(err, "querying following by followed_user_id")
	}
	defer rows.Close()

	var following []*domain.MiddlemanFollow
	for rows.Next() {
		follow := &domain.MiddlemanFollow{}
		if err := rows.Scan(
			&follow.ID,
			&follow.UserID,
			&follow.FollowedUserID,
			&follow.FollowedUserType,
			&follow.Content,
			&follow.CategoryID,
			&follow.ParentID,
			&follow.Approved,
			&follow.Flagged,
			&follow.CreatedAt,
			&follow.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning all following")
		}
		following = append(following, follow)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing rows for All()")
	}

	return following, nil
}

func (r MiddlemanRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.MiddlemanFollow, error) {
	const query = `
SELECT id, user_id, followed_user_id, followed_user_type, content, category_id, parent_id, approved, flagged, created_at, updated_at
  FROM %s
 WHERE user_id = $1
 ORDER BY created_at
`
	rows, err := r.db.QueryContext(ctx, r.table(query), userID)
	if err != nil {
		return nil, errors.Wrap(err, "querying following by user_id")
	}
	defer rows.Close()

	var following []*domain.MiddlemanFollow
	for rows.Next() {
		c := &domain.MiddlemanFollow{}
		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.FollowedUserID,
			&c.FollowedUserType,
			&c.Content,
			&c.CategoryID,
			&c.ParentID,
			&c.Approved,
			&c.Flagged,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scanning following for userID")
		}
		following = append(following, c)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing rows for FindByUserID()")
	}
	return following, nil
}

func (r MiddlemanRepository) Remove(ctx context.Context, followID string) error {
	// If `id` is truly unique, you only need the followID.
	// If you require (followID+followedUserID) as the key, you can add `AND followed_user_id = $2`.
	const query = `DELETE FROM %s WHERE id = $1`
	_, err := r.db.ExecContext(ctx, r.table(query), followID)
	return err
}

func (r MiddlemanRepository) MostFollowedItems(
	ctx context.Context,
	limit, offset int,
) ([]*domain.ItemFollowCount, error) {

	// This query only returns rows with category_id = '' (i.e. "no category")
	baseQuery := `
    SELECT followed_user_id, followed_user_type, category_id, total_following
      FROM item_follow_counts
     WHERE category_id = ''
     ORDER BY total_following DESC
    `
	finalQuery, args := buildLimitOffsetQuery(baseQuery, limit, offset)

	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying most followed items (all categories=empty)")
	}
	defer rows.Close()

	var results []*domain.ItemFollowCount
	for rows.Next() {
		p := &domain.ItemFollowCount{}
		if err := rows.Scan(
			&p.FollowedUserID,
			&p.FollowedUserType,
			&p.CategoryID,
			&p.FollowingCount,
		); err != nil {
			return nil, errors.Wrap(err, "scanning item_follow_counts")
		}
		results = append(results, p)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing aggregator rows")
	}
	return results, nil
}

func (r MiddlemanRepository) MostFollowedItemsByCategory(
	ctx context.Context,
	followedUserType domain.FollowedUserType,
	categoryID string,
	limit, offset int,
) ([]*domain.ItemFollowCount, error) {

	baseQuery := `
    SELECT followed_user_id, followed_user_type, category_id, total_following
      FROM item_follow_counts
     WHERE category_id = $1
       AND followed_user_type = $2
     ORDER BY total_following DESC
    `
	// We'll pass (categoryID, followedUserType) first, then limit/offset
	finalQuery, args := buildLimitOffsetQuery(baseQuery, limit, offset, categoryID, followedUserType.String())

	rows, err := r.db.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying most followed items by category/type")
	}
	defer rows.Close()

	var results []*domain.ItemFollowCount
	for rows.Next() {
		p := &domain.ItemFollowCount{}
		if err := rows.Scan(
			&p.FollowedUserID,
			&p.FollowedUserType,
			&p.CategoryID,
			&p.FollowingCount,
		); err != nil {
			return nil, errors.Wrap(err, "scanning item_follow_counts by category/type")
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
