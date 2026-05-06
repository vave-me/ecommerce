package queries

import (
	"context"
	"middleman/scheduler/internal/domain"
)

type GetScheduler struct {
	UserID string
}

type GetSchedulerHandler struct {
	activities domain.MiddlemanRepository
}

func NewGetSchedulerHandler(activities domain.MiddlemanRepository) GetSchedulerHandler {
	return GetSchedulerHandler{activities: activities}
}

func (h GetSchedulerHandler) GetScheduler(ctx context.Context, query GetScheduler) (*domain.MiddlemanScheduler, error) {
	return h.activities.Find(ctx, query.UserID)
}
