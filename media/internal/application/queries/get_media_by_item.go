package queries

import (
	"context"
	"middleman/media/internal/domain"
)

type GetMediaByItem struct {
	ItemID string
}

type GetMediaByItemHandler struct {
	catalog domain.MiddlemanMediaRepository
}

func NewGetMediaByItemHandler(catalog domain.MiddlemanMediaRepository) GetMediaByItemHandler {
	return GetMediaByItemHandler{catalog: catalog}
}

func (h GetMediaByItemHandler) GetMediaByItem(ctx context.Context, query GetMediaByItem) (*domain.MiddlemanMedia, error) {

	return h.catalog.GetMediaByItem(ctx, query.ItemID)
}
