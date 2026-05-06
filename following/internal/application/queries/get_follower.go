package queries

import (
	"context"
	"middleman/following/internal/domain"
)

type GetFollow struct {
	ID             string
	FollowedUserID string
}

type GetFollowHandler struct {
	following domain.MiddlemanRepository
}

func NewGetFollowHandler(following domain.MiddlemanRepository) GetFollowHandler {
	return GetFollowHandler{following: following}
}

func (h GetFollowHandler) GetFollow(ctx context.Context, query GetFollow) (*domain.MiddlemanFollow, error) {
	return h.following.Find(ctx, query.ID, query.FollowedUserID)
}
