package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stackus/errors"
	"middleman/activity/internal/domain"
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

func (r MiddlemanRepository) Add(ctx context.Context, activityID, userID string) error {
	const query = `INSERT INTO %s (id, user_id) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, r.table(query), activityID, userID)

	return err
}

func (r MiddlemanRepository) Update(ctx context.Context, activityID string, isRead bool) error {
	const query = `UPDATE %s SET is_read = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), activityID, isRead)

	return err
}

func (r MiddlemanRepository) Remove(ctx context.Context, activityID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), activityID)

	return err
}

func (r MiddlemanRepository) Find(ctx context.Context, userID string) (*domain.MiddlemanActivity, error) {
	const query = `SELECT id FROM %s WHERE user_id = $1 LIMIT 1`

	activity := &domain.MiddlemanActivity{
		UserID: userID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), userID).Scan(&activity.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msg("activity with that ID does not exist")
		}
		return nil, errors.Wrap(err, "scanning activity")
	}

	return activity, nil
}

func (r MiddlemanRepository) All(ctx context.Context, userID string) (activities []*domain.MiddlemanActivity, err error) {
	const query = `SELECT id,user_id FROM %s WHERE id = $1 LIMIT 1`

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query), userID)
	if err != nil {
		return nil, errors.Wrap(err, "querying activities")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing product rows")
		}
	}(rows)

	for rows.Next() {
		activity := &domain.MiddlemanActivity{
			UserID: userID,
		}
		err := rows.Scan(&activity.ID, &activity.UserID)
		if err != nil {
			return nil, errors.Wrap(err, "scanning product")
		}

		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing product rows")
	}

	return activities, nil
}

func (r MiddlemanRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
