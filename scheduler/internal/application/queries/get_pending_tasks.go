package queries

import (
	"context"
	"time"

	"github.com/stackus/errors"
	"middleman/scheduler/internal/domain"
)

type GetPendingTasks struct {
	BeforeTime time.Time
	Limit      int
}

type GetPendingTasksHandler struct {
	catalogTasks domain.CatalogTaskRepository
}

func NewGetPendingTasksHandler(catalogTasks domain.CatalogTaskRepository) GetPendingTasksHandler {
	return GetPendingTasksHandler{catalogTasks: catalogTasks}
}

func (h GetPendingTasksHandler) GetPendingTasks(ctx context.Context, query GetPendingTasks) ([]*domain.CatalogTask, error) {
	tasks, err := h.catalogTasks.FindPendingTasks(ctx, query.BeforeTime, query.Limit)
	if err != nil {
		return nil, errors.Wrap(err, "getting pending tasks")
	}

	return tasks, nil
}