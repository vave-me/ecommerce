package domain

import (
	"context"
	"time"
)

type MiddlemanAlert struct {
	ID        string
	UserID    string
	AlertType AlertType
	Message   string
	Payload   map[string]interface{}
	IsRead    bool
	CreatedAt time.Time
}

type CatalogRepository interface {
	Add(ctx context.Context, alertID, userID, alertType, message string, payload map[string]interface{}, isRead bool) error
	Find(ctx context.Context, alertID string) (*MiddlemanAlert, error)
	Read(ctx context.Context, alertID string, isRead bool) error
	Remove(ctx context.Context, alertID string) error
	GetAlerts(ctx context.Context, userID string, isRead bool) ([]*MiddlemanAlert, error)
	GetAlertsByType(ctx context.Context, userID, alertType string, isRead bool) ([]*MiddlemanAlert, error)
}
