package queries

import (
	"context"
	"middleman/media/internal/domain"
)

type GetAllItemImages struct {
	ItemID string
}

type GetAllItemImagesHandler struct {
	catalog domain.MiddlemanImageRepository
}

func NewGetAllItemImagesHandler(catalog domain.MiddlemanImageRepository) GetAllItemImagesHandler {
	return GetAllItemImagesHandler{catalog: catalog}
}

func (h GetAllItemImagesHandler) GetAllItemImages(ctx context.Context, query GetAllItemImages) ([]*domain.MiddlemanImage, error) {

	return h.catalog.FindAllItemImages(ctx, query.ItemID)
}
