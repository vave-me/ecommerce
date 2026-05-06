package postgres

import (
	"context"
	"fmt"

	"github.com/stackus/errors"
	"middleman/internal/es"
	"middleman/scheduler/internal/domain"
)

type taskRepository struct {
	store es.AggregateStore
}

var _ domain.TaskRepository = (*taskRepository)(nil)

// NewTaskRepository creates a new PostgreSQL task repository
func NewTaskRepository(store es.AggregateStore) domain.TaskRepository {
	return &taskRepository{
		store: store,
	}
}

// Load retrieves a task aggregate by ID
func (r *taskRepository) Load(ctx context.Context, taskID string) (*domain.Task, error) {
	task := domain.NewTask(taskID)
	
	err := r.store.Load(ctx, task)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("loading task %s", taskID))
	}

	return task, nil
}

// Save persists a task aggregate
func (r *taskRepository) Save(ctx context.Context, task *domain.Task) error {
	if err := r.store.Save(ctx, task); err != nil {
		return errors.Wrap(err, fmt.Sprintf("saving task %s", task.ID()))
	}

	return nil
}