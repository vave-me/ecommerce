package queries

import (
	"context"
	"middleman/media/internal/domain"
)

type GetMedia struct {
	MediaID string
}

type GetMediaHandler struct {
	catalog domain.MiddlemanMediaRepository
}

func NewGetMediaHandler(catalog domain.MiddlemanMediaRepository) GetMediaHandler {
	return GetMediaHandler{catalog: catalog}
}

func (h GetMediaHandler) GetMedia(ctx context.Context, query GetMedia) (*domain.MiddlemanMedia, error) {

	return h.catalog.GetMedia(ctx, query.MediaID)
}
