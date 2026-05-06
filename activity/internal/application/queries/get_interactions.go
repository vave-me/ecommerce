package queries

import (
	"context"
	"middleman/activity/internal/domain"
)

type GetInteractions struct {
	ActivityID string
}

type GetInteractionsHandler struct {
	interactions domain.MiddlemanInteractionRepository
}

func NewGetInteractionsHandler(interactions domain.MiddlemanInteractionRepository) GetInteractionsHandler {
	return GetInteractionsHandler{interactions: interactions}
}

func (h GetInteractionHandler) GetInteractions(ctx context.Context, query GetInteractions) ([]*domain.MiddlemanInteraction, error) {
	return h.interactions.All(ctx, query.ActivityID)
}
