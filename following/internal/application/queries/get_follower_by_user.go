package queries

import (
	"context"
	"middleman/following/internal/domain"
)

type GetFollowingBySender struct {
	UserID string
}

type GetFollowingBySenderHandler struct {
	following domain.MiddlemanRepository
}

func NewGetFollowingBySenderHandler(following domain.MiddlemanRepository) GetFollowingBySenderHandler {
	return GetFollowingBySenderHandler{following: following}
}

func (h GetFollowingBySenderHandler) GetFollowingBySender(ctx context.Context, query GetFollowingBySender) ([]*domain.MiddlemanFollow, error) {
	return h.following.FindByUserID(ctx, query.UserID)
}
