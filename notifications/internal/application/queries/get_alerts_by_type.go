package queries

import (
	"context"
	"middleman/notifications/internal/domain"
)

type GetAlertsByType struct {
	UserID string
	Type   string
	IsRead *bool  // Optional filter
	Limit  int
	Offset int
}

type GetAlertsByTypeHandler struct {
	repo domain.CatalogRepository
}

func NewGetAlertsByTypeHandler(repo domain.CatalogRepository) GetAlertsByTypeHandler {
	return GetAlertsByTypeHandler{repo: repo}
}

func (h GetAlertsByTypeHandler) GetAlertsByType(ctx context.Context, query GetAlertsByType) ([]*domain.MiddlemanAlert, error) {
	// TODO: Update repository to support optional is_read filter and pagination
	// For now, get all and filter if needed
	isRead := false
	if query.IsRead != nil {
		isRead = *query.IsRead
	}
	return h.repo.GetAlertsByType(ctx, query.UserID, query.Type, isRead)
}
