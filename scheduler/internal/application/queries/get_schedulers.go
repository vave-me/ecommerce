package queries

import (
	"context"
	"middleman/scheduler/internal/domain"
)

type GetSchedulers struct {
	ID string
}

type GetSchedulersHandler struct {
	activities domain.MiddlemanRepository
}

func NewGetSchedulersHandler(activities domain.MiddlemanRepository) GetSchedulersHandler {
	return GetSchedulersHandler{activities: activities}
}

func (h GetSchedulersHandler) GetSchedulers(ctx context.Context, query GetSchedulers) ([]*domain.MiddlemanScheduler, error) {
	return h.activities.All(ctx, query.ID)
}
