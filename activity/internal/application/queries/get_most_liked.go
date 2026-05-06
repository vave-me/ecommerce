package queries

import (
	"context"
	"middleman/activity/internal/domain"
)

type GetMostLiked struct {
	ItemType string
	Limit    int64
}

type GetMostLikedHandler struct {
	interactions domain.MiddlemanInteractionRepository
}

func NewGetMostLikedHandler(interactions domain.MiddlemanInteractionRepository) GetMostLikedHandler {
	return GetMostLikedHandler{interactions: interactions}
}

func (h GetMostLikedHandler) GetMostLiked(ctx context.Context, query GetMostLiked) ([]*domain.MostReactionResult, error) {
	return h.interactions.GetMostLiked(ctx, query.ItemType, query.Limit)
}
