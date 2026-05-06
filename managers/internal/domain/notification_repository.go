package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type NotificationRepository interface {
	// Core notification operations from gRPC service
	ListAlerts(ctx context.Context, userID, alertType string, isRead bool) (*models.ListAlertsResponse, error)
	GetAlertsByType(ctx context.Context, userID, alertType string, isRead bool) (*models.GetAlertsByTypeResponse, error)

	// Additional query methods for AI tooling
	GetAlert(ctx context.Context, alertID string) (*models.Alert, error)
	GetAlertsByUser(ctx context.Context, userID string, limit int64) ([]*models.Alert, error)
	GetUnreadAlerts(ctx context.Context, userID string, limit int64) ([]*models.Alert, error)
	GetReadAlerts(ctx context.Context, userID string, limit int64) ([]*models.Alert, error)
	SearchAlerts(ctx context.Context, userID, query string, limit int64) ([]*models.Alert, error)
	CountAlerts(ctx context.Context, userID string, alertType string, isRead bool) (int64, error)
	CountUnreadAlerts(ctx context.Context, userID string) (int64, error)
}
