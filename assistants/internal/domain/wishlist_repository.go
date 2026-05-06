package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type WishlistRepository interface {
	// Wishlist management
	CreateNewWishlist(ctx context.Context, wishlistID, name string) error
	GetWishlistByName(ctx context.Context, name string) (string, error) // Returns wishlist ID
	GetAllWishlists(ctx context.Context) ([]*models.Wishlist, error)
	DeleteWishlist(ctx context.Context, wishlistID string) error

	// Wishlist item management
	AddItemToWishlist(ctx context.Context, wishlistItemID, wishlistID, itemID, entityType string) error
	RemoveItemFromWishlist(ctx context.Context, wishlistItemID string) error
	GetWishlistItemDetails(ctx context.Context, wishlistItemID, wishlistID, itemID string) (*models.WishlistItem, error)
	GetAllItemsInWishlist(ctx context.Context, wishlistID string) ([]*models.WishlistItem, error)

	// Item operations
	GetWishlistItemByID(ctx context.Context, wishlistItemID string) (*models.WishlistItem, error)
	GetAllWishlistItemsForUser(ctx context.Context) ([]*models.WishlistItem, error)
	FindWishlistByNameDetailed(ctx context.Context, name string) (*models.Wishlist, error)
	GetAllUserWishlistsDetailed(ctx context.Context) ([]*models.Wishlist, error)

	// User convenience methods
	AddItemToUserDefaultWishlist(ctx context.Context, itemID, itemType string) error
	RemoveItemFromUserDefaultWishlist(ctx context.Context, itemID string) error
	GetUserDefaultWishlist(ctx context.Context) (*models.Wishlist, error)
	GetUserWishlistsWithLimit(ctx context.Context, limit int32) ([]*models.Wishlist, error)
	ClearUserDefaultWishlist(ctx context.Context) error
	CheckIfItemInWishlist(ctx context.Context, itemID string) (bool, error)
	GetTotalWishlistItemCount(ctx context.Context) (int, error)
}
