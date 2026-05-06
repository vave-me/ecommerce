package queries

import (
	"context"
	"middleman/media/internal/domain"
)

type GetAllItemVideos struct {
	ItemID string
}

type GetAllItemVideosHandler struct {
	catalog domain.MiddlemanVideoRepository
}

func NewGetAllItemVideosHandler(catalog domain.MiddlemanVideoRepository) GetAllItemVideosHandler {
	return GetAllItemVideosHandler{catalog: catalog}
}

func (h GetAllItemVideosHandler) GetAllItemVideos(ctx context.Context, query GetAllItemVideos) ([]*domain.MiddlemanVideo, error) {

	return h.catalog.FindAllItemVideos(ctx, query.ItemID)
}
