package domain

import "context"

type MiddlemanWishlist struct {
	ID     string
	UserID string
	Name   string
}

type MiddlemanRepository interface {
	AddWishlist(ctx context.Context, wishlistID, userID, name string) error
	Find(ctx context.Context, userID, name string) (*MiddlemanWishlist, error)
	GetWishlists(ctx context.Context, userID string) ([]*MiddlemanWishlist, error)
	Remove(ctx context.Context, wishlistID string) error
}
