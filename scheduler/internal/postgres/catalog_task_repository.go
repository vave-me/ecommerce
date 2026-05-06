package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/stackus/errors"
	"middleman/internal/postgres"
	"middleman/scheduler/internal/domain"
)

type catalogTaskRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.CatalogTaskRepository = (*catalogTaskRepository)(nil)

// NewCatalogTaskRepository creates a new PostgreSQL catalog task repository
func NewCatalogTaskRepository(tableName string, db postgres.DB) domain.CatalogTaskRepository {
	return &catalogTaskRepository{
		tableName: tableName,
		db:        db,
	}
}

// Add adds a new task to the catalog
func (r *catalogTaskRepository) Add(ctx context.Context, task *domain.CatalogTask) error {
	const query = `
		INSERT INTO %s (
			id, manager_id, task_type, scheduled_at, payload, 
			status, created_at, updated_at, executed_at, result, error_message
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	payloadJSON, err := json.Marshal(task.Payload)
	if err != nil {
		return errors.Wrap(err, "marshaling payload")
	}

	_, err = r.db.ExecContext(ctx, fmt.Sprintf(query, r.tableName),
		task.ID, task.ManagerID, task.TaskType, task.ScheduledAt, payloadJSON,
		task.Status, task.CreatedAt, task.UpdatedAt, task.ExecutedAt, task.Result, task.ErrorMessage,
	)
	if err != nil {
		return errors.Wrap(err, "inserting task")
	}

	return nil
}

// Update updates an existing task in the catalog
func (r *catalogTaskRepository) Update(ctx context.Context, taskID string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	var setClauses []string
	var args []interface{}
	argCount := 1

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argCount))
		args = append(args, value)
		argCount++
	}

	// Always update updated_at
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())
	argCount++

	// Add task ID as last argument
	args = append(args, taskID)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d",
		r.tableName,
		strings.Join(setClauses, ", "),
		argCount,
	)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "updating task")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "getting rows affected")
	}

	if rows == 0 {
		return errors.Wrap(errors.ErrNotFound, "task not found")
	}

	return nil
}

// Find retrieves a task by ID
func (r *catalogTaskRepository) Find(ctx context.Context, taskID string) (*domain.CatalogTask, error) {
	const query = `
		SELECT id, manager_id, task_type, scheduled_at, payload,
			   status, created_at, updated_at, executed_at, result, error_message
		FROM %s
		WHERE id = $1
	`

	task := &domain.CatalogTask{}
	var payloadJSON []byte

	err := r.db.QueryRowContext(ctx, fmt.Sprintf(query, r.tableName), taskID).Scan(
		&task.ID, &task.ManagerID, &task.TaskType, &task.ScheduledAt, &payloadJSON,
		&task.Status, &task.CreatedAt, &task.UpdatedAt, &task.ExecutedAt, &task.Result, &task.ErrorMessage,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(errors.ErrNotFound, "task not found")
		}
		return nil, errors.Wrap(err, "querying task")
	}

	if err = json.Unmarshal(payloadJSON, &task.Payload); err != nil {
		return nil, errors.Wrap(err, "unmarshaling payload")
	}

	return task, nil
}

// FindByManagerID retrieves tasks for a specific manager
func (r *catalogTaskRepository) FindByManagerID(ctx context.Context, managerID string, filter domain.TaskFilter) ([]*domain.CatalogTask, error) {
	query := fmt.Sprintf(`
		SELECT id, manager_id, task_type, scheduled_at, payload,
			   status, created_at, updated_at, executed_at, result, error_message
		FROM %s
		WHERE manager_id = $1
	`, r.tableName)

	args := []interface{}{managerID}
	argCount := 2

	// Apply filters
	var conditions []string

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *filter.Status)
		argCount++
	}

	if filter.ScheduledAfter != nil {
		conditions = append(conditions, fmt.Sprintf("scheduled_at > $%d", argCount))
		args = append(args, *filter.ScheduledAfter)
		argCount++
	}

	if filter.ScheduledBefore != nil {
		conditions = append(conditions, fmt.Sprintf("scheduled_at < $%d", argCount))
		args = append(args, *filter.ScheduledBefore)
		argCount++
	}

	if len(filter.TaskTypes) > 0 {
		placeholders := make([]string, len(filter.TaskTypes))
		for i, taskType := range filter.TaskTypes {
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, taskType)
			argCount++
		}
		conditions = append(conditions, fmt.Sprintf("task_type IN (%s)", strings.Join(placeholders, ", ")))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY scheduled_at ASC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
		argCount++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
	}

	return r.queryTasks(ctx, query, args...)
}

// FindPendingTasks retrieves tasks that are due for execution
func (r *catalogTaskRepository) FindPendingTasks(ctx context.Context, beforeTime time.Time, limit int) ([]*domain.CatalogTask, error) {
	query := fmt.Sprintf(`
		SELECT id, manager_id, task_type, scheduled_at, payload,
			   status, created_at, updated_at, executed_at, result, error_message
		FROM %s
		WHERE status = $1 AND scheduled_at <= $2
		ORDER BY scheduled_at ASC
		LIMIT $3
	`, r.tableName)

	return r.queryTasks(ctx, query, domain.TaskStatusPending, beforeTime, limit)
}

// FindByStatus retrieves tasks by status
func (r *catalogTaskRepository) FindByStatus(ctx context.Context, status string, limit int) ([]*domain.CatalogTask, error) {
	query := fmt.Sprintf(`
		SELECT id, manager_id, task_type, scheduled_at, payload,
			   status, created_at, updated_at, executed_at, result, error_message
		FROM %s
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, r.tableName)

	return r.queryTasks(ctx, query, status, limit)
}

// CountByManagerID counts tasks for a specific manager
func (r *catalogTaskRepository) CountByManagerID(ctx context.Context, managerID string, filter domain.TaskFilter) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE manager_id = $1", r.tableName)
	args := []interface{}{managerID}
	argCount := 2

	// Apply same filters as FindByManagerID
	var conditions []string

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *filter.Status)
		argCount++
	}

	if filter.ScheduledAfter != nil {
		conditions = append(conditions, fmt.Sprintf("scheduled_at > $%d", argCount))
		args = append(args, *filter.ScheduledAfter)
		argCount++
	}

	if filter.ScheduledBefore != nil {
		conditions = append(conditions, fmt.Sprintf("scheduled_at < $%d", argCount))
		args = append(args, *filter.ScheduledBefore)
		argCount++
	}

	if len(filter.TaskTypes) > 0 {
		placeholders := make([]string, len(filter.TaskTypes))
		for i, taskType := range filter.TaskTypes {
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, taskType)
			argCount++
		}
		conditions = append(conditions, fmt.Sprintf("task_type IN (%s)", strings.Join(placeholders, ", ")))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "counting tasks")
	}

	return count, nil
}

// Delete removes a task from the catalog
func (r *catalogTaskRepository) Delete(ctx context.Context, taskID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.tableName)

	result, err := r.db.ExecContext(ctx, query, taskID)
	if err != nil {
		return errors.Wrap(err, "deleting task")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "getting rows affected")
	}

	if rows == 0 {
		return errors.Wrap(errors.ErrNotFound, "task not found")
	}

	return nil
}

// queryTasks is a helper method to query multiple tasks
func (r *catalogTaskRepository) queryTasks(ctx context.Context, query string, args ...interface{}) ([]*domain.CatalogTask, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "querying tasks")
	}
	defer rows.Close()

	var tasks []*domain.CatalogTask
	for rows.Next() {
		task := &domain.CatalogTask{}
		var payloadJSON []byte

		err := rows.Scan(
			&task.ID, &task.ManagerID, &task.TaskType, &task.ScheduledAt, &payloadJSON,
			&task.Status, &task.CreatedAt, &task.UpdatedAt, &task.ExecutedAt, &task.Result, &task.ErrorMessage,
		)
		if err != nil {
			return nil, errors.Wrap(err, "scanning task")
		}

		if err = json.Unmarshal(payloadJSON, &task.Payload); err != nil {
			return nil, errors.Wrap(err, "unmarshaling payload")
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "iterating rows")
	}

	return tasks, nil
}