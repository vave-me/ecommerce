package domain

import "context"

type WishlistItemRepository interface {
	Load(ctx context.Context, id string) (*WishlistItem, error)
	Save(ctx context.Context, item *WishlistItem) error
}
