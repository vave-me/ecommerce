package tools

import (
	"context"
	"fmt"
)

// ==================== NOTIFICATION HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeNotificationHandlers() {
	r.handlers["notification_get_user_alerts"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		alertType := getStringParam(params, "alert_type")
		isRead := getBoolParam(params, "is_read", false)
		if userID == "" || alertType == "" {
			return nil, fmt.Errorf("user_id and alert_type are required")
		}
		return reg.notificationRepo.GetUserAlertsList(ctx, userID, alertType, isRead)
	}

	r.handlers["notification_get_alerts_by_type"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		alertType := getStringParam(params, "alert_type")
		isRead := getBoolParam(params, "is_read", false)
		if userID == "" || alertType == "" {
			return nil, fmt.Errorf("user_id and alert_type are required")
		}
		return reg.notificationRepo.GetUserAlertsByNotificationType(ctx, userID, alertType, isRead)
	}

	r.handlers["notification_get_alert_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		alertID := getStringParam(params, "alert_id")
		if alertID == "" {
			return nil, fmt.Errorf("alert_id is required")
		}
		return reg.notificationRepo.GetAlertByID(ctx, alertID)
	}

	r.handlers["notification_get_all_for_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := getInt64Param(params, "limit", 50)
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.notificationRepo.GetAllAlertsForUser(ctx, userID, limit)
	}

	r.handlers["notification_get_unread_for_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := getInt64Param(params, "limit", 50)
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.notificationRepo.GetUnreadAlertsForUser(ctx, userID, limit)
	}

	r.handlers["notification_get_read_for_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := getInt64Param(params, "limit", 50)
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.notificationRepo.GetReadAlertsForUser(ctx, userID, limit)
	}

	r.handlers["notification_search_alerts"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		query := getStringParam(params, "query")
		limit := getInt64Param(params, "limit", 50)
		if userID == "" || query == "" {
			return nil, fmt.Errorf("user_id and query are required")
		}
		return reg.notificationRepo.SearchAlertsByKeyword(ctx, userID, query, limit)
	}

	r.handlers["notification_get_count_by_type_status"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		alertType := getStringParam(params, "alert_type")
		isRead := getBoolParam(params, "is_read", false)
		if userID == "" || alertType == "" {
			return nil, fmt.Errorf("user_id and alert_type are required")
		}
		count, err := reg.notificationRepo.GetAlertCountByTypeAndStatus(ctx, userID, alertType, isRead)
		if err != nil {
			return nil, err
		}
		return map[string]int64{"count": count}, nil
	}

	r.handlers["notification_get_unread_count"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		count, err := reg.notificationRepo.GetUnreadAlertCountForUser(ctx, userID)
		if err != nil {
			return nil, err
		}
		return map[string]int64{"count": count}, nil
	}
}