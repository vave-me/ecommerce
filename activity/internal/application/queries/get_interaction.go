package queries

import (
	"context"
	"middleman/activity/internal/domain"
)

type GetInteraction struct {
	ID string
}

type GetInteractionHandler struct {
	interactions domain.MiddlemanInteractionRepository
}

func NewGetInteractionHandler(interactions domain.MiddlemanInteractionRepository) GetInteractionHandler {
	return GetInteractionHandler{interactions: interactions}
}

func (h GetInteractionHandler) GetInteraction(ctx context.Context, query GetInteraction) (*domain.MiddlemanInteraction, error) {
	return h.interactions.Find(ctx, query.ID)
}
