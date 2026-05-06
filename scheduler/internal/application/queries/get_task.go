package queries

import (
	"context"

	"github.com/stackus/errors"
	"middleman/scheduler/internal/domain"
)

type GetTask struct {
	TaskID string
}

type GetTaskHandler struct {
	catalogTasks domain.CatalogTaskRepository
}

func NewGetTaskHandler(catalogTasks domain.CatalogTaskRepository) GetTaskHandler {
	return GetTaskHandler{catalogTasks: catalogTasks}
}

func (h GetTaskHandler) GetTask(ctx context.Context, query GetTask) (*domain.CatalogTask, error) {
	task, err := h.catalogTasks.Find(ctx, query.TaskID)
	if err != nil {
		return nil, errors.Wrap(err, "getting task")
	}

	return task, nil
}