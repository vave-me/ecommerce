package queries

import (
	"context"
	"middleman/support/internal/domain"
)

type GetSupportChannel struct {
	ID string
}

type GetSupportChannelHandler struct {
	catalog domain.SupportChannelCatalogRepository
}

func NewGetSupportChannelHandler(catalog domain.SupportChannelCatalogRepository) GetSupportChannelHandler {
	return GetSupportChannelHandler{
		catalog: catalog,
	}
}

func (h GetSupportChannelHandler) GetSupportChannel(ctx context.Context, query GetSupportChannel) (*domain.SupportChannelCatalog, error) {
	return h.catalog.Find(ctx, query.ID)
}