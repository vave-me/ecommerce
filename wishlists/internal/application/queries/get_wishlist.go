package queries

import (
	"context"
	"middleman/wishlists/internal/domain"
)

type GetWishlist struct {
	UserID string
	Name   string
}

type GetWishlistHandler struct {
	middleman domain.MiddlemanRepository
}

func NewGetWishlistHandler(catalog domain.MiddlemanRepository) GetWishlistHandler {
	return GetWishlistHandler{middleman: catalog}
}

func (h GetWishlistHandler) GetWishlist(ctx context.Context, query GetWishlist) (*domain.MiddlemanWishlist, error) {
	return h.middleman.Find(ctx, query.UserID, query.Name)
}
