package queries

import (
	"context"
	"middleman/media/internal/domain"
)

type GetAllMediaVideos struct {
	MediaID string
}

type GetAllMediaVideosHandler struct {
	catalog domain.MiddlemanVideoRepository
}

func NewGetAllMediaVideosHandler(catalog domain.MiddlemanVideoRepository) GetAllMediaVideosHandler {
	return GetAllMediaVideosHandler{catalog: catalog}
}

func (h GetAllMediaVideosHandler) GetAllMediaVideos(ctx context.Context, query GetAllMediaVideos) ([]*domain.MiddlemanVideo, error) {

	return h.catalog.FindAllMediaVideos(ctx, query.MediaID)
}
