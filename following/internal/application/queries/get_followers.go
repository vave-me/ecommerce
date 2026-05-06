package queries

import (
	"context"
	"middleman/following/internal/domain"
)

type GetFollowing struct {
	FollowedUserID string
}

type GetFollowingHandler struct {
	following domain.MiddlemanRepository
}

func NewGetFollowingHandler(following domain.MiddlemanRepository) GetFollowingHandler {
	return GetFollowingHandler{following: following}
}

func (h GetFollowingHandler) GetFollowing(ctx context.Context, query GetFollowing) ([]*domain.MiddlemanFollow, error) {
	return h.following.All(ctx, query.FollowedUserID)
}
