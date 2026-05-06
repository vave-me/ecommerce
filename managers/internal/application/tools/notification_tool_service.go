package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// NotificationToolService handles all notification-related operations with streaming execution
type NotificationToolService struct {
	notificationRepository domain.NotificationRepository
}

// NewNotificationToolService creates a new notification tool service instance
func NewNotificationToolService(notificationRepository domain.NotificationRepository) *NotificationToolService {
	return &NotificationToolService{
		notificationRepository: notificationRepository,
	}
}

// ExecuteOperation executes notification operations with streaming progress updates
func (s *NotificationToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (*ToolOperationResult, error) {
	startTime := time.Now()

	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "notification_operation",
		Status:   "progress",
		Progress: 25.0,
		Metadata: map[string]interface{}{
			"step":      "executing_notification_operation",
			"operation": operation,
		},
		Timestamp: time.Now().Unix(),
	}

	var result interface{}
	var err error

	switch operation {
	case "list_alerts":
		result, err = s.handleListAlerts(ctx, parameters, streamChan, toolID)
	case "get_alerts_by_type":
		result, err = s.handleGetAlertsByType(ctx, parameters, streamChan, toolID)
	case "get_alert", "find":
		result, err = s.handleGetAlert(ctx, parameters, streamChan, toolID)
	case "get_alerts_by_user", "search":
		result, err = s.handleGetAlertsByUser(ctx, parameters, streamChan, toolID)
	case "get_unread_alerts":
		result, err = s.handleGetUnreadAlerts(ctx, parameters, streamChan, toolID)
	case "get_read_alerts":
		result, err = s.handleGetReadAlerts(ctx, parameters, streamChan, toolID)
	case "search_alerts":
		result, err = s.handleSearchAlerts(ctx, parameters, streamChan, toolID)
	case "count_alerts":
		result, err = s.handleCountAlerts(ctx, parameters, streamChan, toolID)
	case "count_unread_alerts":
		result, err = s.handleCountUnreadAlerts(ctx, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported notification operation: %s", operation)
	}

	duration := time.Since(startTime)

	if err != nil {
		// Send error update
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "notification_operation",
			Status:   "error",
			Error:    err.Error(),
			Metadata: map[string]interface{}{
				"operation": operation,
			},
			Timestamp: time.Now().Unix(),
		}
		return &ToolOperationResult{
			EntityType: "notifications",
			Operation:  operation,
			Success:    false,
			Error:      err.Error(),
			Duration:   duration,
		}, err
	}

	// Send completion update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "notification_operation",
		Status:   "completed",
		Progress: 100,
		Result:   result,
		Metadata: map[string]interface{}{
			"operation": operation,
			"success":   true,
		},
		Timestamp: time.Now().Unix(),
	}

	return &ToolOperationResult{
		EntityType: "notifications",
		Operation:  operation,
		Success:    true,
		Result:     result,
		Duration:   duration,
	}, nil
}

// handleListAlerts lists alerts for a user with optional filtering
func (s *NotificationToolService) handleListAlerts(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getNotificationStringParam(parameters, "user_id", "")
	alertType := getNotificationStringParam(parameters, "alert_type", "")
	if alertType == "" {
		alertType = getNotificationStringParam(parameters, "type", "")
	}
	isRead := getNotificationBoolParam(parameters, "is_read", false)

	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "notification_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "listing_alerts",
			"user_id":    userID,
			"alert_type": alertType,
			"is_read":    isRead,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("NotificationToolService: Listing alerts for user: %s, type: %s, isRead: %t", userID, alertType, isRead)
	result, err := s.notificationRepository.ListAlerts(ctx, userID, alertType, isRead)
	if err != nil {
		return nil, fmt.Errorf("list alerts failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "notifications",
		"operation":   "list_alerts",
		"result":      result,
		"user_id":     userID,
	}, nil
}

// handleGetAlertsByType gets alerts by type for a user
func (s *NotificationToolService) handleGetAlertsByType(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getNotificationStringParam(parameters, "user_id", "")
	alertType := getNotificationStringParam(parameters, "alert_type", "")
	if alertType == "" {
		alertType = getNotificationStringParam(parameters, "type", "")
	}
	isRead := getNotificationBoolParam(parameters, "is_read", false)

	if userID == "" || alertType == "" {
		return nil, fmt.Errorf("user_id and alert_type parameters required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "notification_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "getting_alerts_by_type",
			"user_id":    userID,
			"alert_type": alertType,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("NotificationToolService: Getting alerts by type for user: %s, type: %s, isRead: %t", userID, alertType, isRead)
	result, err := s.notificationRepository.GetAlertsByType(ctx, userID, alertType, isRead)
	if err != nil {
		return nil, fmt.Errorf("get alerts by type failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "notifications",
		"operation":   "get_alerts_by_type",
		"result":      result,
		"user_id":     userID,
		"alert_type":  alertType,
	}, nil
}

// handleGetAlert retrieves a specific alert
func (s *NotificationToolService) handleGetAlert(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	alertID := getNotificationStringParam(parameters, "alert_id", "")
	if alertID == "" {
		alertID = getNotificationStringParam(parameters, "id", "")
	}

	if alertID == "" {
		return nil, fmt.Errorf("alert_id or id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "notification_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":     "getting_alert",
			"alert_id": alertID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("NotificationToolService: Getting alert: %s", alertID)
	result, err := s.notificationRepository.GetAlert(ctx, alertID)
	if err != nil {
		return nil, fmt.Errorf("get alert failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "notifications",
		"operation":   "get_alert",
		"result":      result,
		"alert_id":    alertID,
	}, nil
}

// handleGetAlertsByUser retrieves alerts for a user with pagination
func (s *NotificationToolService) handleGetAlertsByUser(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getNotificationStringParam(parameters, "user_id", "")
	limit := getNotificationInt64Param(parameters, "limit", 20)

	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "notification_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":    "getting_alerts_by_user",
			"user_id": userID,
			"limit":   limit,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("NotificationToolService: Getting alerts for user: %s, limit: %d", userID, limit)
	result, err := s.notificationRepository.GetAlertsByUser(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get alerts by user failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "notifications",
		"operation":   "get_alerts_by_user",
		"result":      result,
		"user_id":     userID,
	}, nil
}

// handleGetUnreadAlerts retrieves unread alerts for a user
func (s *NotificationToolService) handleGetUnreadAlerts(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getNotificationStringParam(parameters, "user_id", "")
	limit := getNotificationInt64Param(parameters, "limit", 20)

	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "notification_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":    "getting_unread_alerts",
			"user_id": userID,
			"limit":   limit,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("NotificationToolService: Getting unread alerts for user: %s, limit: %d", userID, limit)
	result, err := s.notificationRepository.GetUnreadAlerts(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get unread alerts failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "notifications",
		"operation":   "get_unread_alerts",
		"result":      result,
		"user_id":     userID,
	}, nil
}

// handleGetReadAlerts retrieves read alerts for a user
func (s *NotificationToolService) handleGetReadAlerts(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getNotificationStringParam(parameters, "user_id", "")
	limit := getNotificationInt64Param(parameters, "limit", 20)

	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "notification_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":    "getting_read_alerts",
			"user_id": userID,
			"limit":   limit,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("NotificationToolService: Getting read alerts for user: %s, limit: %d", userID, limit)
	result, err := s.notificationRepository.GetReadAlerts(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get read alerts failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "notifications",
		"operation":   "get_read_alerts",
		"result":      result,
		"user_id":     userID,
	}, nil
}

// handleSearchAlerts searches alerts for a user
func (s *NotificationToolService) handleSearchAlerts(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getNotificationStringParam(parameters, "user_id", "")
	query := getNotificationStringParam(parameters, "query", "")
	if query == "" {
		query = getNotificationStringParam(parameters, "search_term", "")
	}
	limit := getNotificationInt64Param(parameters, "limit", 20)

	if userID == "" || query == "" {
		return nil, fmt.Errorf("user_id and query parameters required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "notification_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":    "searching_alerts",
			"user_id": userID,
			"query":   query,
			"limit":   limit,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("NotificationToolService: Searching alerts for user: %s, query: %s, limit: %d", userID, query, limit)
	result, err := s.notificationRepository.SearchAlerts(ctx, userID, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search alerts failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "notifications",
		"operation":   "search_alerts",
		"result":      result,
		"user_id":     userID,
		"query":       query,
	}, nil
}

// handleCountAlerts counts alerts for a user
func (s *NotificationToolService) handleCountAlerts(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getNotificationStringParam(parameters, "user_id", "")
	alertType := getNotificationStringParam(parameters, "alert_type", "")
	if alertType == "" {
		alertType = getNotificationStringParam(parameters, "type", "")
	}
	isRead := getNotificationBoolParam(parameters, "is_read", false)

	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "notification_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":       "counting_alerts",
			"user_id":    userID,
			"alert_type": alertType,
			"is_read":    isRead,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("NotificationToolService: Counting alerts for user: %s, type: %s, isRead: %t", userID, alertType, isRead)
	result, err := s.notificationRepository.CountAlerts(ctx, userID, alertType, isRead)
	if err != nil {
		return nil, fmt.Errorf("count alerts failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "notifications",
		"operation":   "count_alerts",
		"result":      result,
		"user_id":     userID,
	}, nil
}

// handleCountUnreadAlerts counts unread alerts for a user
func (s *NotificationToolService) handleCountUnreadAlerts(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getNotificationStringParam(parameters, "user_id", "")

	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "notification_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":    "counting_unread_alerts",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("NotificationToolService: Counting unread alerts for user: %s", userID)
	result, err := s.notificationRepository.CountUnreadAlerts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("count unread alerts failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "notifications",
		"operation":   "count_unread_alerts",
		"result":      result,
		"user_id":     userID,
	}, nil
}

// Helper functions for parameter extraction
func getNotificationStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getNotificationInt64Param(params map[string]interface{}, key string, defaultValue int64) int64 {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case float64:
			return int64(v)
		}
	}
	return defaultValue
}

func getNotificationBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := params[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}
