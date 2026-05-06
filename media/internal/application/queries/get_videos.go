package queries

import (
	"context"
	"middleman/media/internal/domain"
)

type GetAllVideos struct {
	Page      int64
	PageSize  int64
	SortBy    string
	SortOrder string
}

type GetAllVideosHandler struct {
	catalog domain.MiddlemanVideoRepository
}

func NewGetAllVideosHandler(catalog domain.MiddlemanVideoRepository) GetAllVideosHandler {
	return GetAllVideosHandler{catalog: catalog}
}

func (h GetAllVideosHandler) GetAllVideos(ctx context.Context, query GetAllVideos) ([]*domain.MiddlemanVideo, int64, error) {
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

	return h.catalog.FindAllVideos(ctx, query.Page, query.PageSize, query.SortBy, query.SortOrder)
}
