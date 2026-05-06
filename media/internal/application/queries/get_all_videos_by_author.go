package queries

import (
	"context"
	"middleman/media/internal/domain"
)

type GetAllVideosByAuthor struct {
	UserID    string
	Page      int64
	PageSize  int64
	SortBy    string
	SortOrder string
}

type GetAllVideosByAuthorHandler struct {
	catalog domain.MiddlemanVideoRepository
}

func NewGetAllVideosByAuthorHandler(catalog domain.MiddlemanVideoRepository) GetAllVideosByAuthorHandler {
	return GetAllVideosByAuthorHandler{catalog: catalog}
}

func (h GetAllVideosByAuthorHandler) GetAllVideosByAuthor(ctx context.Context, query GetAllVideosByAuthor) ([]*domain.MiddlemanVideo, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.SortBy == "" {
		query.SortBy = "createdAt" // Default sorting field
	}
	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		query.SortOrder = "asc" // Default sorting order
	}

	return h.catalog.FindAllVideosByAuthor(ctx, query.UserID, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
