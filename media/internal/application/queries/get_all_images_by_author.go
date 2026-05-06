package queries

import (
	"context"
	"middleman/media/internal/domain"
)

type GetAllImagesByAuthor struct {
	UserID    string
	Page      int64
	PageSize  int64
	SortBy    string
	SortOrder string
}

type GetAllImagesByAuthorHandler struct {
	catalog domain.MiddlemanImageRepository
}

func NewGetAllImagesByAuthorHandler(catalog domain.MiddlemanImageRepository) GetAllImagesByAuthorHandler {
	return GetAllImagesByAuthorHandler{catalog: catalog}
}

func (h GetAllImagesByAuthorHandler) GetAllImagesByAuthor(ctx context.Context, query GetAllImagesByAuthor) ([]*domain.MiddlemanImage, int64, error) {
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

	return h.catalog.FindAllImagesByAuthor(ctx, query.UserID, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
