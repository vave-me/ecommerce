package queries

import (
	"context"
	"middleman/wishlists/internal/domain"
)

type GetWishlists struct {
	UserID string
}

type GetWishlistsHandler struct {
	middleman domain.MiddlemanRepository
}

func NewGetWishlistsHandler(catalog domain.MiddlemanRepository) GetWishlistsHandler {
	return GetWishlistsHandler{middleman: catalog}
}

func (h GetWishlistsHandler) GetWishlists(ctx context.Context, query GetWishlists) ([]*domain.MiddlemanWishlist, error) {
	return h.middleman.GetWishlists(ctx, query.UserID)
}
