package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/scheduler/internal/domain"
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

func (r MiddlemanRepository) Add(ctx context.Context, schedulerID, userID string) error {
	const query = `INSERT INTO %s (id, user_id) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, r.table(query), schedulerID, userID)

	return err
}

func (r MiddlemanRepository) Update(ctx context.Context, schedulerID string, isRead bool) error {
	const query = `UPDATE %s SET is_read = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), schedulerID, isRead)

	return err
}

func (r MiddlemanRepository) Remove(ctx context.Context, schedulerID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), schedulerID)

	return err
}

func (r MiddlemanRepository) Find(ctx context.Context, userID string) (*domain.MiddlemanScheduler, error) {
	const query = `SELECT id FROM %s WHERE user_id = $1 LIMIT 1`

	scheduler := &domain.MiddlemanScheduler{
		UserID: userID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), userID).Scan(&scheduler.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msg("scheduler with that ID does not exist")
		}
		return nil, errors.Wrap(err, "scanning scheduler")
	}

	return scheduler, nil
}

func (r MiddlemanRepository) All(ctx context.Context, userID string) (schedulers []*domain.MiddlemanScheduler, err error) {
	const query = `SELECT id, user_id FROM %s WHERE user_id = $1`

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query), userID)
	if err != nil {
		return nil, errors.Wrap(err, "querying schedulers")
	}
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = errors.Wrap(closeErr, "closing scheduler rows")
		}
	}(rows)

	for rows.Next() {
		scheduler := &domain.MiddlemanScheduler{}
		err := rows.Scan(&scheduler.ID, &scheduler.UserID)
		if err != nil {
			return nil, errors.Wrap(err, "scanning scheduler")
		}

		schedulers = append(schedulers, scheduler)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing scheduler rows")
	}

	return schedulers, nil
}

func (r MiddlemanRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
