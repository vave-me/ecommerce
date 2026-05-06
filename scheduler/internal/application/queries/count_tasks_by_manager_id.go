package queries

import (
	"context"

	"github.com/stackus/errors"
	"middleman/scheduler/internal/domain"
)

type CountTasksByManagerID struct {
	ManagerID string
	Filter    domain.TaskFilter
}

type CountTasksByManagerIDHandler struct {
	catalogTasks domain.CatalogTaskRepository
}

func NewCountTasksByManagerIDHandler(catalogTasks domain.CatalogTaskRepository) CountTasksByManagerIDHandler {
	return CountTasksByManagerIDHandler{catalogTasks: catalogTasks}
}

func (h CountTasksByManagerIDHandler) CountTasksByManagerID(ctx context.Context, query CountTasksByManagerID) (int, error) {
	count, err := h.catalogTasks.CountByManagerID(ctx, query.ManagerID, query.Filter)
	if err != nil {
		return 0, errors.Wrap(err, "counting tasks")
	}

	return count, nil
}