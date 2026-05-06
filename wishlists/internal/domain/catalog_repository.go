package domain

import (
	"context"
)

type CatalogWishlistItem struct {
	ID         string
	WishlistID string
	ItemID     string
	EntityType string
}

type CatalogRepository interface {
	AddWishlistItem(ctx context.Context, wishlistItemID, wishlistID, itemID, entityType string) error
	RemoveWishlistItem(ctx context.Context, wishlistItemID string) error
	Find(ctx context.Context, wishlistItemID string) (*CatalogWishlistItem, error)
	GetWishlistItems(ctx context.Context, userID string) ([]*CatalogWishlistItem, error)
}
