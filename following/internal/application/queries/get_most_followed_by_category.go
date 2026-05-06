package queries

import (
	"context"
	"middleman/following/internal/domain"
)

type GetMostFollowedByCategory struct {
	FollowedUserType domain.FollowedUserType
	CategoryID       string
	Offset           int
	Limit            int
}

type GetMostFollowedByCategoryHandler struct {
	following domain.MiddlemanRepository
}

func NewGetMostFollowedByCategoryHandler(following domain.MiddlemanRepository) GetMostFollowedByCategoryHandler {
	return GetMostFollowedByCategoryHandler{following: following}
}

func (h GetMostFollowedByCategoryHandler) GetMostFollowedByCategory(ctx context.Context, query GetMostFollowedByCategory) ([]*domain.ItemFollowCount, error) {
	return h.following.MostFollowedItemsByCategory(ctx, query.FollowedUserType, query.CategoryID, query.Offset, query.Limit)
}
