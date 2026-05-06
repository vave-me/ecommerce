package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type NotificationRepository interface {
	// Core notification operations from gRPC service
	GetUserAlertsList(ctx context.Context, userID, alertType string, isRead bool) (*models.ListAlertsResponse, error)
	GetUserAlertsByNotificationType(ctx context.Context, userID, alertType string, isRead bool) (*models.GetAlertsByTypeResponse, error)

	// Additional query methods for AI tooling
	GetAlertByID(ctx context.Context, alertID string) (*models.Alert, error)
	GetAllAlertsForUser(ctx context.Context, userID string, limit int64) ([]*models.Alert, error)
	GetUnreadAlertsForUser(ctx context.Context, userID string, limit int64) ([]*models.Alert, error)
	GetReadAlertsForUser(ctx context.Context, userID string, limit int64) ([]*models.Alert, error)
	SearchAlertsByKeyword(ctx context.Context, userID, query string, limit int64) ([]*models.Alert, error)
	GetAlertCountByTypeAndStatus(ctx context.Context, userID string, alertType string, isRead bool) (int64, error)
	GetUnreadAlertCountForUser(ctx context.Context, userID string) (int64, error)
}
