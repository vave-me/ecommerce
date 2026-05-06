package queries

import (
	"context"
	"middleman/scheduler/internal/domain"
)

type GetActions struct {
	SchedulerID string
}

type GetActionsHandler struct {
	interactions domain.MiddlemanActionRepository
}

func NewGetActionsHandler(interactions domain.MiddlemanActionRepository) GetActionsHandler {
	return GetActionsHandler{interactions: interactions}
}

func (h GetActionsHandler) GetActions(ctx context.Context, query GetActions) ([]*domain.MiddlemanAction, error) {
	return h.interactions.All(ctx, query.SchedulerID)
}
