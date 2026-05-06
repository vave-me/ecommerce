package queries

import (
	"context"
	"middleman/activity/internal/domain"
)

type GetActivity struct {
	UserID string
}

type GetActivityHandler struct {
	activities domain.MiddlemanRepository
}

func NewGetActivityHandler(activities domain.MiddlemanRepository) GetActivityHandler {
	return GetActivityHandler{activities: activities}
}

func (h GetActivityHandler) GetActivity(ctx context.Context, query GetActivity) (*domain.MiddlemanActivity, error) {
	return h.activities.Find(ctx, query.UserID)
}
