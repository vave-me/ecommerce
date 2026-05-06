package queries

import (
	"context"
	"middleman/notifications/internal/domain"
)

type ListAlerts struct {
	UserID string
	IsRead *bool  // Optional filter
	Limit  int
	Offset int
}

type ListAlertsHandler struct {
	repo domain.CatalogRepository
}

func NewListAlertsHandler(repo domain.CatalogRepository) ListAlertsHandler {
	return ListAlertsHandler{repo: repo}
}

func (h ListAlertsHandler) ListAlerts(ctx context.Context, query ListAlerts) ([]*domain.MiddlemanAlert, error) {
	// TODO: Update repository to support optional is_read filter and pagination
	// For now, get all and filter if needed
	isRead := false
	if query.IsRead != nil {
		isRead = *query.IsRead
	}
	return h.repo.GetAlerts(ctx, query.UserID, isRead)
}
