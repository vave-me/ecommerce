package queries

import (
	"context"
	"middleman/scheduler/internal/domain"
)

type GetAction struct {
	ID string
}

type GetActionHandler struct {
	interactions domain.MiddlemanActionRepository
}

func NewGetActionHandler(interactions domain.MiddlemanActionRepository) GetActionHandler {
	return GetActionHandler{interactions: interactions}
}

func (h GetActionHandler) GetAction(ctx context.Context, query GetAction) (*domain.MiddlemanAction, error) {
	return h.interactions.Find(ctx, query.ID)
}
