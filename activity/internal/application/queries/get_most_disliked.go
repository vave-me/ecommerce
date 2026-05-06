package queries

import (
	"context"
	"middleman/activity/internal/domain"
)

type GetMostDisliked struct {
	ItemType string
	Limit    int64
}

type GetMostDislikedHandler struct {
	interactions domain.MiddlemanInteractionRepository
}

func NewGetMostDislikedHandler(interactions domain.MiddlemanInteractionRepository) GetMostDislikedHandler {
	return GetMostDislikedHandler{interactions: interactions}
}

func (h GetMostDislikedHandler) GetMostDisliked(ctx context.Context, query GetMostDisliked) ([]*domain.MostReactionResult, error) {
	return h.interactions.GetMostDisliked(ctx, query.ItemType, query.Limit)
}
