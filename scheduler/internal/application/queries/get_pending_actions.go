package queries

import (
	"context"
	"time"
	"middleman/scheduler/internal/domain"
)

type (
	GetPendingActions struct {
		BeforeTime time.Time
	}

	GetPendingActionsHandler struct {
		actions domain.MiddlemanActionRepository
	}
)

func NewGetPendingActionsHandler(actions domain.MiddlemanActionRepository) GetPendingActionsHandler {
	return GetPendingActionsHandler{actions: actions}
}

func (h GetPendingActionsHandler) GetPendingActions(ctx context.Context, query GetPendingActions) ([]*domain.MiddlemanAction, error) {
	return h.actions.GetPendingActions(ctx, query.BeforeTime)
}