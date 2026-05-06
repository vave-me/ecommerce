package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type WishlistRepository interface {
	// Wishlist management
	CreateWishlist(ctx context.Context, wishlistID, name string) error
	GetWishlist(ctx context.Context, name string) (string, error) // Returns wishlist ID
	GetWishlists(ctx context.Context) ([]*models.Wishlist, error)
	RemoveWishlist(ctx context.Context, wishlistID string) error

	// Wishlist item management
	AddWishlistItem(ctx context.Context, wishlistItemID, wishlistID, itemID, entityType string) error
	RemoveWishlistItem(ctx context.Context, wishlistItemID string) error
	GetWishlistItem(ctx context.Context, wishlistItemID, wishlistID, itemID string) (*models.WishlistItem, error)
	GetWishlistItems(ctx context.Context, wishlistID string) ([]*models.WishlistItem, error)

	// Legacy methods (for backward compatibility)
	FindItemFromWishlist(ctx context.Context, wishlistItemID string) (*models.WishlistItem, error)
	GetAllWishlistItems(ctx context.Context) ([]*models.WishlistItem, error)
	FindWishlistByName(ctx context.Context, name string) (*models.Wishlist, error)
	GetAllUserWishlists(ctx context.Context) ([]*models.Wishlist, error)

	// Methods needed by wishlist tool service
	AddToWishlist(ctx context.Context, itemID, itemType string) error
	RemoveFromWishlist(ctx context.Context, itemID string) error
	GetUserWishlist(ctx context.Context) (*models.Wishlist, error)
	GetUserWishlists(ctx context.Context, limit int32) ([]*models.Wishlist, error)
	ClearWishlist(ctx context.Context) error
	IsInWishlist(ctx context.Context, itemID string) (bool, error)
	GetWishlistCount(ctx context.Context) (int, error)
}
