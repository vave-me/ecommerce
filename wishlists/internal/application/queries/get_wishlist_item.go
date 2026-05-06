package queries

import (
	"context"
	"middleman/wishlists/internal/domain"
)

type GetWishlistItem struct {
	WishlistItemID string
}

type GetWishlistItemHandler struct {
	catalog domain.CatalogRepository
}

func NewGetWishlistItemHandler(catalog domain.CatalogRepository) GetWishlistItemHandler {
	return GetWishlistItemHandler{catalog: catalog}
}

func (h GetWishlistItemHandler) GetWishlistItem(ctx context.Context, query GetWishlistItem) (*domain.CatalogWishlistItem, error) {
	return h.catalog.Find(ctx, query.WishlistItemID)
}
