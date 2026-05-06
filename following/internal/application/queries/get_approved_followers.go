package queries

import (
	"context"
	"middleman/following/internal/domain"
)

type GetApprovedFollowing struct {
	FollowedUserID string
}

type GetApprovedFollowingHandler struct {
	following domain.MiddlemanRepository
}

func NewGetApprovedFollowingHandler(following domain.MiddlemanRepository) GetApprovedFollowingHandler {
	return GetApprovedFollowingHandler{following: following}
}

func (h GetApprovedFollowingHandler) GetApprovedFollowing(ctx context.Context, query GetApprovedFollowing) ([]*domain.MiddlemanFollow, error) {
	return h.following.All(ctx, query.FollowedUserID)
}
