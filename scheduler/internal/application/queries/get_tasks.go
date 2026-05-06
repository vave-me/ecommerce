package queries

import (
	"context"

	"github.com/stackus/errors"
	"middleman/scheduler/internal/domain"
)

type GetTasks struct {
	ManagerID string
	Filter    domain.TaskFilter
}

type GetTasksHandler struct {
	catalogTasks domain.CatalogTaskRepository
}

func NewGetTasksHandler(catalogTasks domain.CatalogTaskRepository) GetTasksHandler {
	return GetTasksHandler{catalogTasks: catalogTasks}
}

func (h GetTasksHandler) GetTasks(ctx context.Context, query GetTasks) ([]*domain.CatalogTask, error) {
	tasks, err := h.catalogTasks.FindByManagerID(ctx, query.ManagerID, query.Filter)
	if err != nil {
		return nil, errors.Wrap(err, "getting tasks")
	}

	return tasks, nil
}