package queries

import (
	"context"
	"middleman/following/internal/domain"
)

type GetMostFollowed struct {
	Offset int
	Limit  int
}

type GetMostFollowedHandler struct {
	following domain.MiddlemanRepository
}

func NewGetMostFollowedHandler(following domain.MiddlemanRepository) GetMostFollowedHandler {
	return GetMostFollowedHandler{following: following}
}

func (h GetMostFollowedHandler) GetMostFollowed(ctx context.Context, query GetMostFollowed) ([]*domain.ItemFollowCount, error) {
	return h.following.MostFollowedItems(ctx, query.Offset, query.Limit)
}
