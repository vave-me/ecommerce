package grpc

import (
	"context"
	"fmt"
	"github.com/stackus/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"middleman/internal/auth"
	"middleman/internal/errorsotel"
	"middleman/notifications/internal/application"
	"middleman/notifications/internal/application/commands"
	"middleman/notifications/internal/application/queries"
	"middleman/notifications/internal/domain"
	"middleman/notifications/notificationspb"
)

type server struct {
	app application.App
	notificationspb.UnimplementedNotificationsServiceServer
}

var _ notificationspb.NotificationsServiceServer = (*server)(nil)

func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	notificationspb.RegisterNotificationsServiceServer(registrar, server{app: app})
	return nil
}

func (s server) ListAlerts(ctx context.Context, request *notificationspb.ListAlertsRequest) (resp *notificationspb.ListAlertsResponse, err error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	span.SetAttributes(
		attribute.String("UserID", userID),
	)

	query := queries.ListAlerts{
		UserID: userID,
	}
	
	// Handle optional is_read filter
	if request.IsRead != nil {
		query.IsRead = request.IsRead
	}
	
	// Handle pagination
	if request.Limit > 0 {
		query.Limit = int(request.Limit)
	} else {
		query.Limit = 50 // default
	}
	query.Offset = int(request.Offset)
	
	alerts, err := s.app.ListAlerts(ctx, query)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Apply pagination in memory for now (until repository supports it)
	totalCount := len(alerts)
	
	// Apply offset and limit
	start := query.Offset
	if start > len(alerts) {
		start = len(alerts)
	}
	
	end := start + query.Limit
	if end > len(alerts) {
		end = len(alerts)
	}
	
	paginatedAlerts := alerts[start:end]
	hasMore := end < totalCount

	protoAlerts := make([]*notificationspb.Alert, len(paginatedAlerts))
	for i, alert := range paginatedAlerts {
		protoAlerts[i] = s.alertsFromDomain(alert)
	}

	return &notificationspb.ListAlertsResponse{
		Alerts:     protoAlerts,
		TotalCount: int32(totalCount),
		HasMore:    hasMore,
	}, nil
}
func (s server) GetAlertsByType(ctx context.Context, request *notificationspb.GetAlertsByTypeRequest) (resp *notificationspb.GetAlertsByTypeResponse, err error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span.SetAttributes(
		attribute.String("UserID", userID),
	)

	query := queries.GetAlertsByType{
		UserID: userID,
		Type:   request.GetType(),
	}
	
	// Handle optional is_read filter
	if request.IsRead != nil {
		query.IsRead = request.IsRead
	}
	
	// Handle pagination
	if request.Limit > 0 {
		query.Limit = int(request.Limit)
	} else {
		query.Limit = 50 // default
	}
	query.Offset = int(request.Offset)
	
	alerts, err := s.app.GetAlertsByType(ctx, query)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Apply pagination in memory for now (until repository supports it)
	totalCount := len(alerts)
	
	// Apply offset and limit
	start := query.Offset
	if start > len(alerts) {
		start = len(alerts)
	}
	
	end := start + query.Limit
	if end > len(alerts) {
		end = len(alerts)
	}
	
	paginatedAlerts := alerts[start:end]
	hasMore := end < totalCount

	protoAlerts := make([]*notificationspb.Alert, len(paginatedAlerts))
	for i, alert := range paginatedAlerts {
		protoAlerts[i] = s.alertsFromDomain(alert)
	}

	return &notificationspb.GetAlertsByTypeResponse{
		Alerts:     protoAlerts,
		TotalCount: int32(totalCount),
		HasMore:    hasMore,
	}, nil
}

func (s server) MarkAlertAsRead(ctx context.Context, request *notificationspb.MarkAlertAsReadRequest) (*emptypb.Empty, error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	alertID := request.GetAlertId()
	
	// Validate alert ID
	if alertID == "" {
		return nil, status.Error(grpc_code.InvalidArgument, "alert_id is required")
	}
	
	span.SetAttributes(
		attribute.String("UserID", userID),
		attribute.String("AlertID", alertID),
	)

	// TODO: Verify alert belongs to user before marking as read
	// This would require loading the alert and checking the userID
	
	// Use the existing ReadAlert command
	cmd := commands.ReadAlert{
		ID: alertID,
	}

	if err := s.app.ReadAlert(ctx, cmd); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		
		// Provide more specific error messages
		switch {
		case errors.Is(err, domain.ErrAlertNotFound):
			return nil, status.Error(grpc_code.NotFound, "alert not found")
		case errors.Is(err, domain.ErrUserIDCannotBeBlank):
			return nil, status.Error(grpc_code.InvalidArgument, "invalid alert")
		default:
			return nil, status.Error(grpc_code.Internal, "failed to mark alert as read")
		}
	}

	return &emptypb.Empty{}, nil
}

func (s server) MarkAllAlertsAsRead(ctx context.Context, request *notificationspb.MarkAllAlertsAsReadRequest) (*emptypb.Empty, error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span.SetAttributes(
		attribute.String("UserID", userID),
		attribute.String("Type", request.GetType()),
	)

	// Get all unread alerts for the user
	query := queries.ListAlerts{
		UserID: userID,
		IsRead: &[]bool{false}[0], // pointer to false
		Limit:  1000, // Get all unread alerts
	}

	// If type is specified, filter by type
	if request.GetType() != "" {
		// Use GetAlertsByType instead
		typeQuery := queries.GetAlertsByType{
			UserID: userID,
			Type:   request.GetType(),
			IsRead: &[]bool{false}[0],
			Limit:  1000,
		}
		alerts, err := s.app.GetAlertsByType(ctx, typeQuery)
		if err != nil {
			span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
			span.SetStatus(codes.Error, err.Error())
			return nil, status.Error(grpc_code.Internal, "failed to get alerts")
		}

		// Mark each alert as read
		for _, alert := range alerts {
			cmd := commands.ReadAlert{
				ID: alert.ID,
			}
			if err := s.app.ReadAlert(ctx, cmd); err != nil {
				// Log error but continue with other alerts
				attrs := append(errorsotel.ErrAttrs(err), attribute.String("AlertID", alert.ID))
				span.RecordError(err, trace.WithAttributes(attrs...))
			}
		}
	} else {
		// Mark all alerts as read
		alerts, err := s.app.ListAlerts(ctx, query)
		if err != nil {
			span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
			span.SetStatus(codes.Error, err.Error())
			return nil, status.Error(grpc_code.Internal, "failed to get alerts")
		}

		// Mark each alert as read
		for _, alert := range alerts {
			cmd := commands.ReadAlert{
				ID: alert.ID,
			}
			if err := s.app.ReadAlert(ctx, cmd); err != nil {
				// Log error but continue with other alerts
				attrs := append(errorsotel.ErrAttrs(err), attribute.String("AlertID", alert.ID))
				span.RecordError(err, trace.WithAttributes(attrs...))
			}
		}
	}

	return &emptypb.Empty{}, nil
}

func (s server) DeleteAlert(ctx context.Context, request *notificationspb.DeleteAlertRequest) (*emptypb.Empty, error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	alertID := request.GetAlertId()
	
	// Validate alert ID
	if alertID == "" {
		return nil, status.Error(grpc_code.InvalidArgument, "alert_id is required")
	}
	
	span.SetAttributes(
		attribute.String("UserID", userID),
		attribute.String("AlertID", alertID),
	)

	// Delete the alert
	cmd := commands.DeleteAlert{
		AlertID: alertID,
		UserID:  userID,
	}

	if err := s.app.DeleteAlert(ctx, cmd); err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		
		// Provide more specific error messages
		switch {
		case errors.Is(err, domain.ErrAlertNotFound):
			return nil, status.Error(grpc_code.NotFound, "alert not found")
		case errors.Is(err, errors.ErrUnauthorized):
			return nil, status.Error(grpc_code.PermissionDenied, "alert does not belong to user")
		default:
			return nil, status.Error(grpc_code.Internal, "failed to delete alert")
		}
	}

	return &emptypb.Empty{}, nil
}

func (s server) GetUnreadCount(ctx context.Context, request *notificationspb.GetUnreadCountRequest) (*notificationspb.GetUnreadCountResponse, error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span.SetAttributes(
		attribute.String("UserID", userID),
	)

	// Get all unread alerts for the user
	query := queries.ListAlerts{
		UserID: userID,
		IsRead: &[]bool{false}[0], // pointer to false
		Limit:  1000, // Get all for counting
	}

	alerts, err := s.app.ListAlerts(ctx, query)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(grpc_code.Internal, "failed to get unread count")
	}

	return &notificationspb.GetUnreadCountResponse{
		Count: int32(len(alerts)),
	}, nil
}

func (s server) GetNotificationStats(ctx context.Context, request *notificationspb.GetNotificationStatsRequest) (*notificationspb.GetNotificationStatsResponse, error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span.SetAttributes(
		attribute.String("UserID", userID),
	)

	// Get all alerts for the user
	query := queries.ListAlerts{
		UserID: userID,
		Limit:  10000, // Get all for stats
	}

	alerts, err := s.app.ListAlerts(ctx, query)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(grpc_code.Internal, "failed to get notification stats")
	}

	// Calculate stats
	totalCount := len(alerts)
	unreadCount := 0
	countByType := make(map[string]int32)
	var lastReadAt *timestamppb.Timestamp

	for _, alert := range alerts {
		if !alert.IsRead {
			unreadCount++
		} else if lastReadAt == nil || alert.CreatedAt.After(lastReadAt.AsTime()) {
			lastReadAt = timestamppb.New(alert.CreatedAt)
		}

		typeStr := alert.AlertType.String()
		countByType[typeStr]++
	}

	return &notificationspb.GetNotificationStatsResponse{
		TotalCount:   int32(totalCount),
		UnreadCount:  int32(unreadCount),
		CountByType:  countByType,
		LastReadAt:   lastReadAt,
	}, nil
}

func (s server) GetNotificationPreferences(ctx context.Context, request *notificationspb.GetNotificationPreferencesRequest) (*notificationspb.GetNotificationPreferencesResponse, error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span.SetAttributes(
		attribute.String("UserID", userID),
	)

	// TODO: Implement preferences storage and retrieval
	// For now, return default preferences
	allTypes := []domain.AlertType{
		domain.MessageType,
		domain.CommentType,
		domain.OfferType,
		domain.OrderType,
		domain.PaymentType,
		domain.ReviewType,
		domain.FollowingType,
		domain.ProductType,
		domain.WishlistType,
		domain.SupportType,
		domain.BasketType,
		domain.InteractionType,
	}

	preferences := make([]*notificationspb.NotificationPreference, len(allTypes))
	for i, alertType := range allTypes {
		preferences[i] = &notificationspb.NotificationPreference{
			Type:         alertType.String(),
			Enabled:      true,
			EmailEnabled: true,
			PushEnabled:  true,
		}
	}

	return &notificationspb.GetNotificationPreferencesResponse{
		Preferences:    preferences,
		GlobalEnabled:  true,
		EmailEnabled:   true,
		PushEnabled:    true,
	}, nil
}

func (s server) UpdateNotificationPreferences(ctx context.Context, request *notificationspb.UpdateNotificationPreferencesRequest) (*emptypb.Empty, error) {
	span := trace.SpanFromContext(ctx)

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span.SetAttributes(
		attribute.String("UserID", userID),
	)

	// TODO: Implement preferences storage and update
	// For now, return not implemented
	return nil, status.Error(grpc_code.Unimplemented, "update preferences not yet implemented")
}

func (s server) alertsFromDomain(alert *domain.MiddlemanAlert) *notificationspb.Alert {
	payload := s.payloadFromDomain(alert.Payload)
	return &notificationspb.Alert{
		Id:      alert.ID,
		UserId:  alert.UserID,
		Type:    alert.AlertType.String(),
		Payload: payload,
		Message: alert.Message,
		IsRead:  alert.IsRead,
		CreatedAt: timestamppb.New(alert.CreatedAt),
	}
}

func (s server) payloadFromDomain(payload map[string]interface{}) map[string]string {
	convertedPayload := make(map[string]string)
	for key, value := range payload {
		if strVal, ok := value.(string); ok {
			convertedPayload[key] = strVal
		} else {
			convertedPayload[key] = fmt.Sprintf("%v", value)
		}
	}
	return convertedPayload
}
