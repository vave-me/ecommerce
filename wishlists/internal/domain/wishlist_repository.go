package domain

import (
	"context"
)

type WishlistRepository interface {
	Load(ctx context.Context, wishlistID string) (*Wishlist, error)
	Save(ctx context.Context, wishlist *Wishlist) error
}
