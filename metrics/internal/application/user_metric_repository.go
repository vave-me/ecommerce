package application

import (
	"context"
	"middleman/metrics/internal/models"
)

type UserMetricRepository interface {
	AddUserMetric(ctx context.Context, userID, userType string) error
	GetUserMetric(ctx context.Context, userID string) (*models.UserMetric, error)
	RemoveUserMetric(ctx context.Context, userID string) error
	UpdateUserMetric(ctx context.Context, userID string, metricType, metricTypeAction string) error
}

type UserMetricCacheRepository interface {
	UserMetricRepository
}
