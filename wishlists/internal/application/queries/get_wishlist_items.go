package queries

import (
	"context"
	"middleman/wishlists/internal/domain"
)

type GetWishlistItems struct {
	WishlistID string
}

type GetWishlistItemsHandler struct {
	catalog domain.CatalogRepository
}

func NewGetWishlistItemsHandler(catalog domain.CatalogRepository) GetWishlistItemsHandler {
	return GetWishlistItemsHandler{catalog: catalog}
}

func (h GetWishlistItemsHandler) GetWishlistItems(ctx context.Context, query GetWishlistItems) ([]*domain.CatalogWishlistItem, error) {
	return h.catalog.GetWishlistItems(ctx, query.WishlistID)
}
