package queries

import (
	"context"
	"middleman/activity/internal/domain"
)

type GetActivities struct {
	ID string
}

type GetActivitiesHandler struct {
	activities domain.MiddlemanRepository
}

func NewGetActivitiesHandler(activities domain.MiddlemanRepository) GetActivitiesHandler {
	return GetActivitiesHandler{activities: activities}
}

func (h GetActivitiesHandler) GetActivities(ctx context.Context, query GetActivities) ([]*domain.MiddlemanActivity, error) {
	return h.activities.All(ctx, query.ID)
}
