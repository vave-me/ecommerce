package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/stackus/errors"

	"middleman/internal/postgres"
	"middleman/scheduler/internal/domain"
)

type MiddlemanActionRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.MiddlemanActionRepository = (*MiddlemanActionRepository)(nil)

func NewMiddlemanActionRepository(tableName string, db postgres.DB) MiddlemanActionRepository {
	return MiddlemanActionRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r MiddlemanActionRepository) Add(
	ctx context.Context,
	actionID, schedulerID, task string,
	executionTime time.Time,
) error {
	const query = `INSERT INTO %s (id, scheduler_id, natural_language_task, execution_time, status, created_at)
                   VALUES ($1, $2, $3, $4, 'pending', NOW())`

	_, err := r.db.ExecContext(ctx, r.table(query), actionID, schedulerID, task, executionTime)
	if err != nil {
		return errors.Wrap(err, "inserting new action")
	}
	return nil
}

func (r MiddlemanActionRepository) UpdateStatus(
	ctx context.Context,
	actionID, status, result, errorMessage string,
) error {
	const query = `UPDATE %s
                   SET status = $2, result = $3, error_message = $4, executed_at = NOW(), updated_at = NOW()
                   WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), actionID, status, result, errorMessage)
	if err != nil {
		return errors.Wrap(err, "updating action status")
	}
	return nil
}

func (r MiddlemanActionRepository) Remove(
	ctx context.Context,
	actionID string,
) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), actionID)
	if err != nil {
		return errors.Wrap(err, "deleting action")
	}
	return nil
}

func (r MiddlemanActionRepository) Find(
	ctx context.Context,
	actionID string,
) (*domain.MiddlemanAction, error) {
	const query = `SELECT id, scheduler_id, natural_language_task, execution_time, status, 
                          created_at, executed_at, result, error_message
                   FROM %s
                   WHERE id = $1
                   LIMIT 1`

	action := &domain.MiddlemanAction{}

	err := r.db.QueryRowContext(ctx, r.table(query), actionID).
		Scan(&action.ID, &action.SchedulerID, &action.NaturalLanguageTask, 
		     &action.ExecutionTime, &action.Status, &action.CreatedAt,
		     &action.ExecutedAt, &action.Result, &action.ErrorMessage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Msg("action with that ID does not exist")
		}
		return nil, errors.Wrap(err, "scanning action")
	}
	return action, nil
}

func (r MiddlemanActionRepository) All(
	ctx context.Context,
	schedulerID string,
) ([]*domain.MiddlemanAction, error) {
	const query = `SELECT id, scheduler_id, natural_language_task, execution_time, status,
                          created_at, executed_at, result, error_message
                   FROM %s
                   WHERE scheduler_id = $1
                   ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, r.table(query), schedulerID)
	if err != nil {
		return nil, errors.Wrap(err, "querying actions")
	}
	defer rows.Close()

	var actions []*domain.MiddlemanAction
	for rows.Next() {
		action := &domain.MiddlemanAction{}
		if err := rows.Scan(&action.ID, &action.SchedulerID, &action.NaturalLanguageTask,
			&action.ExecutionTime, &action.Status, &action.CreatedAt,
			&action.ExecutedAt, &action.Result, &action.ErrorMessage); err != nil {
			return nil, errors.Wrap(err, "scanning action row")
		}
		actions = append(actions, action)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing actions rows")
	}
	return actions, nil
}

func (r MiddlemanActionRepository) GetPendingActions(
	ctx context.Context,
	beforeTime time.Time,
) ([]*domain.MiddlemanAction, error) {
	const query = `SELECT id, scheduler_id, natural_language_task, execution_time, status,
                          created_at, executed_at, result, error_message
                   FROM %s
                   WHERE status = 'pending' AND execution_time <= $1
                   ORDER BY execution_time ASC
                   FOR UPDATE SKIP LOCKED`

	rows, err := r.db.QueryContext(ctx, r.table(query), beforeTime)
	if err != nil {
		return nil, errors.Wrap(err, "querying pending actions")
	}
	defer rows.Close()

	var actions []*domain.MiddlemanAction
	for rows.Next() {
		action := &domain.MiddlemanAction{}
		if err := rows.Scan(&action.ID, &action.SchedulerID, &action.NaturalLanguageTask,
			&action.ExecutionTime, &action.Status, &action.CreatedAt,
			&action.ExecutedAt, &action.Result, &action.ErrorMessage); err != nil {
			return nil, errors.Wrap(err, "scanning pending action row")
		}
		actions = append(actions, action)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing pending actions rows")
	}
	return actions, nil
}

func (r MiddlemanActionRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}