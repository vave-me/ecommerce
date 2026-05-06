package queries

import (
	"context"
	"middleman/media/internal/domain"
)

type GetAllMediaImages struct {
	MediaID string
}

type GetAllMediaImagesHandler struct {
	catalog domain.MiddlemanImageRepository
}

func NewGetAllMediaImagesHandler(catalog domain.MiddlemanImageRepository) GetAllMediaImagesHandler {
	return GetAllMediaImagesHandler{catalog: catalog}
}

func (h GetAllMediaImagesHandler) GetAllMediaImages(ctx context.Context, query GetAllMediaImages) ([]*domain.MiddlemanImage, error) {

	return h.catalog.FindAllMediaImages(ctx, query.MediaID)
}
