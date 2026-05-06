package grpc

import (
	"context"
	"fmt"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
	"middleman/notifications/notificationspb"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// NotificationRepository calls the remote notifications service (gRPC).
type NotificationRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.NotificationRepository = (*NotificationRepository)(nil)

// NewNotificationRepositoryWithAuth creates a new NotificationRepository with JWT authentication support
func NewNotificationRepository(endpoint string, authInstance *auth.Auth) NotificationRepository {
	return NotificationRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// ListAlerts lists alerts for a user with optional filters
func (r NotificationRepository) ListAlerts(ctx context.Context, userID string, alertType string, isRead bool) (*models.ListAlertsResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := notificationspb.NewNotificationsServiceClient(conn)
	resp, err := client.ListAlerts(ctx, &notificationspb.ListAlertsRequest{
		Type:   alertType,
		IsRead: &isRead,
	})
	if err != nil {
		return nil, fmt.Errorf("ListAlerts RPC failed: %w", err)
	}

	alerts := make([]models.Alert, len(resp.GetAlerts()))
	for i, pbAlert := range resp.GetAlerts() {
		alert := r.convertAlertFromPb(pbAlert)
		alerts[i] = *alert
	}

	return &models.ListAlertsResponse{
		Alerts: alerts,
		Total:  int64(len(alerts)), // Note: protobuf doesn't have total field, using length
	}, nil
}

// GetAlertsByType gets alerts by type for a user
func (r NotificationRepository) GetAlertsByType(ctx context.Context, userID, alertType string, isRead bool) (*models.GetAlertsByTypeResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := notificationspb.NewNotificationsServiceClient(conn)
	resp, err := client.GetAlertsByType(ctx, &notificationspb.GetAlertsByTypeRequest{

		Type:   alertType,
		IsRead: &isRead,
	})
	if err != nil {
		return nil, fmt.Errorf("GetAlertsByType RPC failed: %w", err)
	}

	alerts := make([]models.Alert, len(resp.GetAlerts()))
	for i, pbAlert := range resp.GetAlerts() {
		alert := r.convertAlertFromPb(pbAlert)
		alerts[i] = *alert
	}

	return &models.GetAlertsByTypeResponse{
		Alerts: alerts,
	}, nil
}

// GetAlert retrieves a single alert by ID (mock implementation for AI tooling)
func (r NotificationRepository) GetAlert(ctx context.Context, alertID string) (*models.Alert, error) {
	// Note: This would typically require a GetAlert RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("GetAlert called for ID: %s (mock implementation)", alertID)

	return &models.Alert{
		ID:      alertID,
		UserID:  "mock_user",
		Type:    models.AlertTypeInfo,
		Message: "Mock alert message",
		Payload: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		IsRead: false,
	}, nil
}

// GetAlertsByUser retrieves alerts by user ID (mock implementation for AI tooling)
func (r NotificationRepository) GetAlertsByUser(ctx context.Context, userID string, limit int64) ([]*models.Alert, error) {
	// Note: This would typically use ListAlerts with proper pagination
	// For now, we'll use ListAlerts and limit the results
	log.Printf("GetAlertsByUser called for user: %s, limit: %d", userID, limit)

	resp, err := r.ListAlerts(ctx, userID, "", false) // Get all alerts (read and unread)
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice and apply limit
	alerts := make([]*models.Alert, 0, limit)
	for i, alert := range resp.Alerts {
		if int64(i) >= limit {
			break
		}
		alertCopy := alert
		alerts = append(alerts, &alertCopy)
	}

	return alerts, nil
}

// GetUnreadAlerts retrieves unread alerts for a user (mock implementation for AI tooling)
func (r NotificationRepository) GetUnreadAlerts(ctx context.Context, userID string, limit int64) ([]*models.Alert, error) {
	log.Printf("GetUnreadAlerts called for user: %s, limit: %d", userID, limit)

	resp, err := r.ListAlerts(ctx, userID, "", false) // Get unread alerts
	if err != nil {
		return nil, err
	}

	// Filter unread alerts and apply limit
	alerts := make([]*models.Alert, 0, limit)
	for _, alert := range resp.Alerts {
		if !alert.IsRead && int64(len(alerts)) < limit {
			alertCopy := alert
			alerts = append(alerts, &alertCopy)
		}
	}

	return alerts, nil
}

// GetReadAlerts retrieves read alerts for a user (mock implementation for AI tooling)
func (r NotificationRepository) GetReadAlerts(ctx context.Context, userID string, limit int64) ([]*models.Alert, error) {
	log.Printf("GetReadAlerts called for user: %s, limit: %d", userID, limit)

	resp, err := r.ListAlerts(ctx, userID, "", true) // Get read alerts
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice and apply limit
	alerts := make([]*models.Alert, 0, limit)
	for i, alert := range resp.Alerts {
		if int64(i) >= limit {
			break
		}
		alertCopy := alert
		alerts = append(alerts, &alertCopy)
	}

	return alerts, nil
}

// SearchAlerts searches alerts by query (mock implementation for AI tooling)
func (r NotificationRepository) SearchAlerts(ctx context.Context, userID, query string, limit int64) ([]*models.Alert, error) {
	// Note: This would typically require a SearchAlerts RPC method in the protobuf
	// For now, we'll return a mock implementation
	log.Printf("SearchAlerts called for user: %s, query: %s, limit: %d (mock implementation)", userID, query, limit)

	alerts := make([]*models.Alert, 0, limit)
	for i := int64(0); i < limit && i < 3; i++ { // Mock max 3 results
		alerts = append(alerts, &models.Alert{
			ID:      fmt.Sprintf("alert_%d", i+1),
			UserID:  userID,
			Type:    models.AlertTypeInfo,
			Message: fmt.Sprintf("Search result %d for query: %s", i+1, query),
			Payload: map[string]string{
				"search_query": query,
				"result_rank":  fmt.Sprintf("%d", i+1),
			},
			IsRead: false,
		})
	}

	return alerts, nil
}

// CountAlerts counts alerts with filters (mock implementation for AI tooling)
func (r NotificationRepository) CountAlerts(ctx context.Context, userID string, alertType string, isRead bool) (int64, error) {
	log.Printf("CountAlerts called for user: %s, type: %s, isRead: %t (mock implementation)", userID, alertType, isRead)

	resp, err := r.ListAlerts(ctx, userID, alertType, isRead)
	if err != nil {
		return 0, err
	}

	return int64(len(resp.Alerts)), nil
}

// CountUnreadAlerts counts unread alerts for a user (mock implementation for AI tooling)
func (r NotificationRepository) CountUnreadAlerts(ctx context.Context, userID string) (int64, error) {
	log.Printf("CountUnreadAlerts called for user: %s (mock implementation)", userID)

	return r.CountAlerts(ctx, userID, "", false)
}

// convertAlertFromPb converts protobuf Alert to domain Alert
func (r NotificationRepository) convertAlertFromPb(pbAlert *notificationspb.Alert) *models.Alert {
	if pbAlert == nil {
		return nil
	}

	return &models.Alert{
		ID:      pbAlert.GetId(),
		UserID:  pbAlert.GetUserId(),
		Type:    pbAlert.GetType(),
		Message: pbAlert.GetMessage(),
		Payload: pbAlert.GetPayload(),
		IsRead:  false, // Note: protobuf doesn't have is_read field, defaulting to false
	}
}

// dial establishes a gRPC connection to the notifications service
// dial sets up a gRPC connection with the microservice endpoint
func (r NotificationRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r NotificationRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
