package grpc

import (
	"context"
	"fmt"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/wishlists/wishlistspb"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// WishlistRepository calls the remote wishlists service (gRPC) as a fallback.
type WishlistRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.WishlistRepository = (*WishlistRepository)(nil)

// NewWishlistRepositoryWithAuth creates a new WishlistRepository with JWT authentication support
func NewWishlistRepository(endpoint string, authInstance *auth.Auth) WishlistRepository {
	return WishlistRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// CreateNewWishlist creates a new wishlist
func (r WishlistRepository) CreateNewWishlist(ctx context.Context, wishlistID, name string) error {
	log.Printf("[WISHLIST_GRPC] CreateNewWishlist called for  name=%s", name)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Printf("[WISHLIST_GRPC] Failed to connect to wishlists service: %v", err)
		return err
	}
	defer conn.Close()

	log.Printf("[WISHLIST_GRPC] Successfully connected to wishlists service, calling CreateWishlist RPC...")

	client := wishlistspb.NewWishlistServiceClient(conn)
	resp, err := client.CreateWishlist(ctx, &wishlistspb.CreateWishlistRequest{

		Name: name,
	})
	if err != nil {
		log.Printf("[WISHLIST_GRPC] CreateWishlist RPC failed: %v", err)
		return fmt.Errorf("CreateWishlist RPC failed: %w", err)
	}

	log.Printf("[WISHLIST_GRPC] CreateWishlist RPC successful, created wishlist with ID: %s", resp.GetId())
	return nil
}

// GetWishlistByName retrieves a wishlist by user ID and name, returns wishlist ID
func (r WishlistRepository) GetWishlistByName(ctx context.Context, name string) (string, error) {
	log.Printf("[WISHLIST_GRPC] GetWishlistByName called for name=%s", name)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Printf("[WISHLIST_GRPC] Failed to connect to wishlists service: %v", err)
		return "", err
	}
	defer conn.Close()

	client := wishlistspb.NewWishlistServiceClient(conn)
	resp, err := client.GetWishlist(ctx, &wishlistspb.GetWishlistRequest{

		Name: name,
	})
	if err != nil {
		log.Printf("[WISHLIST_GRPC] GetWishlist RPC failed: %v", err)
		return "", fmt.Errorf("GetWishlist RPC failed: %w", err)
	}

	log.Printf("[WISHLIST_GRPC] GetWishlist RPC successful, wishlist ID: %s", resp.GetWishlistId())
	return resp.GetWishlistId(), nil
}

// GetAllWishlists retrieves all wishlists for a user
func (r WishlistRepository) GetAllWishlists(ctx context.Context) ([]*models.Wishlist, error) {

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Printf("[WISHLIST_GRPC] Failed to connect to wishlists service: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := wishlistspb.NewWishlistServiceClient(conn)
	resp, err := client.GetWishlists(ctx, &wishlistspb.GetWishlistsRequest{})
	if err != nil {
		log.Printf("[WISHLIST_GRPC] GetWishlists RPC failed: %v", err)
		return nil, fmt.Errorf("GetWishlists RPC failed: %w", err)
	}

	log.Printf("[WISHLIST_GRPC] GetWishlists RPC successful, received %d wishlists", len(resp.GetWishlists()))

	var results []*models.Wishlist
	for i, w := range resp.GetWishlists() {
		log.Printf("[WISHLIST_GRPC] Converting wishlist %d: ID=%s, Name=%s", i, w.GetId(), w.GetName())
		if wishlist := r.wishlistToDomain(w); wishlist != nil {
			results = append(results, wishlist)
		}
	}

	log.Printf("[WISHLIST_GRPC] GetWishlists returning %d converted wishlists", len(results))
	return results, nil
}

// DeleteWishlist removes a wishlist by ID
func (r WishlistRepository) DeleteWishlist(ctx context.Context, wishlistID string) error {
	log.Printf("[WISHLIST_GRPC] DeleteWishlist called for wishlistID=%s", wishlistID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := wishlistspb.NewWishlistServiceClient(conn)
	_, err = client.RemoveWishlist(ctx, &wishlistspb.RemoveWishlistRequest{
		Id: wishlistID,
	})
	if err != nil {
		log.Printf("[WISHLIST_GRPC] RemoveWishlist RPC failed: %v", err)
		return fmt.Errorf("RemoveWishlist RPC failed: %w", err)
	}

	log.Printf("[WISHLIST_GRPC] RemoveWishlist RPC successful for wishlist ID: %s", wishlistID)
	return nil
}

// AddItemToWishlist adds an item to a wishlist
func (r WishlistRepository) AddItemToWishlist(ctx context.Context, wishlistItemID, wishlistID, itemID, entityType string) error {
	log.Printf("[WISHLIST_GRPC] AddItemToWishlist called for wishlistID=%s, itemID=%s", wishlistID, itemID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Printf("[WISHLIST_GRPC] Failed to connect to wishlists service: %v", err)
		return err
	}
	defer conn.Close()

	client := wishlistspb.NewWishlistServiceClient(conn)
	resp, err := client.AddWishlistItem(ctx, &wishlistspb.AddWishlistItemRequest{
		WishlistId: wishlistID,
		ItemId:     itemID,
		EntityType: entityType,
	})
	if err != nil {
		log.Printf("[WISHLIST_GRPC] AddWishlistItem RPC failed: %v", err)
		return fmt.Errorf("AddWishlistItem RPC failed: %w", err)
	}

	log.Printf("[WISHLIST_GRPC] AddWishlistItem RPC successful, created item with ID: %s", resp.GetId())
	return nil
}

// RemoveItemFromWishlist removes an item from a wishlist
func (r WishlistRepository) RemoveItemFromWishlist(ctx context.Context, wishlistItemID string) error {
	log.Printf("[WISHLIST_GRPC] RemoveItemFromWishlist called for wishlistItemID=%s", wishlistItemID)

	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := wishlistspb.NewWishlistServiceClient(conn)
	_, err = client.RemoveWishlistItem(ctx, &wishlistspb.RemoveWishlistItemRequest{
		Id: wishlistItemID,
	})
	if err != nil {
		log.Printf("[WISHLIST_GRPC] RemoveWishlistItem RPC failed: %v", err)
		return fmt.Errorf("RemoveWishlistItem RPC failed: %w", err)
	}

	log.Printf("[WISHLIST_GRPC] RemoveWishlistItem RPC successful for item ID: %s", wishlistItemID)
	return nil
}

// GetWishlistItemDetails retrieves a specific wishlist item
func (r WishlistRepository) GetWishlistItemDetails(ctx context.Context, wishlistItemID, wishlistID, itemID string) (*models.WishlistItem, error) {
	log.Printf("[WISHLIST_GRPC] GetWishlistItemDetails called for wishlistItemID=%s", wishlistItemID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Printf("[WISHLIST_GRPC] Failed to connect to wishlists service: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := wishlistspb.NewWishlistServiceClient(conn)
	resp, err := client.GetWishlistItem(ctx, &wishlistspb.GetWishlistItemRequest{
		Id:         wishlistItemID,
		WishlistId: wishlistID,
		ItemId:     itemID,
	})
	if err != nil {
		log.Printf("[WISHLIST_GRPC] GetWishlistItem RPC failed: %v", err)
		return nil, fmt.Errorf("GetWishlistItem RPC failed: %w", err)
	}

	item := r.wishlistItemToDomain(resp.GetItem())
	log.Printf("[WISHLIST_GRPC] GetWishlistItem RPC successful, item ID: %s", item.ID)
	return item, nil
}

// GetAllItemsInWishlist retrieves all items from a wishlist
func (r WishlistRepository) GetAllItemsInWishlist(ctx context.Context, wishlistID string) ([]*models.WishlistItem, error) {
	log.Printf("[WISHLIST_GRPC] GetAllItemsInWishlist called for wishlistID=%s", wishlistID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Printf("[WISHLIST_GRPC] Failed to connect to wishlists service: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := wishlistspb.NewWishlistServiceClient(conn)
	resp, err := client.GetWishlistItems(ctx, &wishlistspb.GetWishlistItemsRequest{
		WishlistId: wishlistID,
	})
	if err != nil {
		log.Printf("[WISHLIST_GRPC] GetWishlistItems RPC failed: %v", err)
		return nil, fmt.Errorf("GetWishlistItems RPC failed: %w", err)
	}

	log.Printf("[WISHLIST_GRPC] GetWishlistItems RPC successful, received %d items", len(resp.GetItems()))

	var results []*models.WishlistItem
	for i, item := range resp.GetItems() {
		log.Printf("[WISHLIST_GRPC] Converting item %d: ID=%s, ItemID=%s", i, item.GetId(), item.GetItemId())
		if domainItem := r.wishlistItemToDomain(item); domainItem != nil {
			results = append(results, domainItem)
		}
	}

	log.Printf("[WISHLIST_GRPC] GetWishlistItems returning %d converted items", len(results))
	return results, nil
}

// Legacy methods (for backward compatibility)

// GetWishlistItemByID retrieves a wishlist item by ID
func (r WishlistRepository) GetWishlistItemByID(ctx context.Context, wishlistItemID string) (*models.WishlistItem, error) {
	// Use the new GetWishlistItemDetails method, but we need wishlistID and itemID
	// For backward compatibility, we'll just pass empty strings for the missing parameters
	return r.GetWishlistItemDetails(ctx, wishlistItemID, "", "")
}

// GetAllWishlistItemsForUser retrieves all wishlist items for a user
func (r WishlistRepository) GetAllWishlistItemsForUser(ctx context.Context) ([]*models.WishlistItem, error) {

	// First get all wishlists for the user
	wishlists, err := r.GetAllWishlists(ctx)
	if err != nil {
		return nil, err
	}

	var allItems []*models.WishlistItem
	for _, wishlist := range wishlists {
		items, err := r.GetAllItemsInWishlist(ctx, wishlist.ID)
		if err != nil {
			log.Printf("[WISHLIST_GRPC] Failed to get items for wishlist %s: %v", wishlist.ID, err)
			continue
		}
		allItems = append(allItems, items...)
	}

	log.Printf("[WISHLIST_GRPC] GetAllWishlistItems returning %d total items", len(allItems))
	return allItems, nil
}

// FindWishlistByNameDetailed finds a wishlist by name for a user
func (r WishlistRepository) FindWishlistByNameDetailed(ctx context.Context, name string) (*models.Wishlist, error) {

	// Get all wishlists and find the one with matching name
	wishlists, err := r.GetAllWishlists(ctx)
	if err != nil {
		return nil, err
	}

	for _, wishlist := range wishlists {
		if wishlist.Name == name {
			log.Printf("[WISHLIST_GRPC] Found matching wishlist: %s", wishlist.ID)
			return wishlist, nil
		}
	}

	return nil, fmt.Errorf("wishlist with name '%s' ", name)
}

// GetAllUserWishlistsDetailed retrieves all wishlists for a user
func (r WishlistRepository) GetAllUserWishlistsDetailed(ctx context.Context) ([]*models.Wishlist, error) {
	return r.GetAllWishlists(ctx)
}

// AddItemToUserDefaultWishlist adds an item to user's default wishlist
func (r WishlistRepository) AddItemToUserDefaultWishlist(ctx context.Context, itemID, itemType string) error {
	log.Printf("[WISHLIST_GRPC] AddItemToUserDefaultWishlist called for , itemID=%s, itemType=%s", itemID, itemType)

	// For simplicity, we'll use a default wishlist for the user
	// In a real implementation, you might want to create or find a default wishlist
	defaultWishlistName := "default"

	// Try to get the default wishlist ID, create if doesn't exist
	wishlistID, err := r.GetWishlistByName(ctx, defaultWishlistName)
	if err != nil {
		// Create default wishlist if it doesn't exist
		wishlistID = "default"
		err = r.CreateNewWishlist(ctx, wishlistID, defaultWishlistName)
		if err != nil {
			return fmt.Errorf("failed to create default wishlist: %w", err)
		}
	}

	// Add item to the wishlist
	itemIDGenerated := "item_" + wishlistID + "_" + itemID
	return r.AddItemToWishlist(ctx, itemIDGenerated, wishlistID, itemID, itemType)
}

func (r WishlistRepository) RemoveItemFromUserDefaultWishlist(ctx context.Context, itemID string) error {

	// For simplicity, assume we're removing from the default wishlist
	// In a real implementation, you might search across all user's wishlists
	defaultWishlistName := "default"
	wishlistID, err := r.GetWishlistByName(ctx, defaultWishlistName)
	if err != nil {
		return fmt.Errorf("failed to find default wishlist: %w", err)
	}

	// Remove the item (we need to construct the wishlist item ID)
	itemIDGenerated := "item_" + wishlistID + "_" + itemID
	return r.RemoveItemFromWishlist(ctx, itemIDGenerated)
}

func (r WishlistRepository) GetUserDefaultWishlist(ctx context.Context) (*models.Wishlist, error) {

	// Get the user's wishlists and return the first one (or default)
	wishlists, err := r.GetAllWishlists(ctx)
	if err != nil {
		return nil, err
	}

	if len(wishlists) == 0 {
		return nil, fmt.Errorf("no wishlists found for user ")
	}

	// Return the first wishlist (or you could look for a "default" one)
	return wishlists[0], nil
}

func (r WishlistRepository) GetUserWishlistsWithLimit(ctx context.Context, limit int32) ([]*models.Wishlist, error) {
	log.Printf("[WISHLIST_GRPC] GetUserWishlistsWithLimit called , limit=%d", limit)

	// Get all wishlists and apply limit
	wishlists, err := r.GetAllWishlists(ctx)
	if err != nil {
		return nil, err
	}

	// Apply limit
	if int32(len(wishlists)) > limit {
		wishlists = wishlists[:limit]
	}

	return wishlists, nil
}

func (r WishlistRepository) ClearUserDefaultWishlist(ctx context.Context) error {
	log.Printf("[WISHLIST_GRPC] ClearUserDefaultWishlist called for ")

	// Get user's default wishlist and remove all items
	defaultWishlistName := "default"
	wishlistID, err := r.GetWishlistByName(ctx, defaultWishlistName)
	if err != nil {
		return fmt.Errorf("failed to find default wishlist: %w", err)
	}

	// Get all items in the wishlist
	items, err := r.GetAllItemsInWishlist(ctx, wishlistID)
	if err != nil {
		return fmt.Errorf("failed to get wishlist items: %w", err)
	}

	// Remove each item
	for _, item := range items {
		err := r.RemoveItemFromWishlist(ctx, item.ID)
		if err != nil {
			log.Printf("[WISHLIST_GRPC] Failed to remove item %s: %v", item.ID, err)
			// Continue with other items even if one fails
		}
	}

	return nil
}

func (r WishlistRepository) CheckIfItemInWishlist(ctx context.Context, itemID string) (bool, error) {
	log.Printf("[WISHLIST_GRPC] CheckIfItemInWishlist called , itemID=%s", itemID)

	// Get user's wishlists and check if item exists in any of them
	wishlists, err := r.GetAllWishlists(ctx)
	if err != nil {
		return false, err
	}

	for _, wishlist := range wishlists {
		items, err := r.GetAllItemsInWishlist(ctx, wishlist.ID)
		if err != nil {
			continue // Skip this wishlist if we can't get items
		}

		for _, item := range items {
			if item.ItemID == itemID {
				return true, nil
			}
		}
	}

	return false, nil
}

func (r WishlistRepository) GetTotalWishlistItemCount(ctx context.Context) (int, error) {
	log.Printf("[WISHLIST_GRPC] GetTotalWishlistItemCount called for ")

	// Get all items for the user and count them
	allItems, err := r.GetAllWishlistItemsForUser(ctx)
	if err != nil {
		return 0, err
	}

	return len(allItems), nil
}

// Helper methods for domain conversion

// wishlistToDomain converts a wishlistspb.Wishlist into our internal models.Wishlist
func (r WishlistRepository) wishlistToDomain(pb *wishlistspb.Wishlist) *models.Wishlist {
	if pb == nil {
		return nil
	}

	return &models.Wishlist{
		ID:          pb.GetId(),
		UserID:      pb.GetUserId(),
		Name:        pb.GetName(),
		Description: pb.GetDescription(),
	}
}

// wishlistItemToDomain converts a wishlistspb.WishlistItem into our internal models.WishlistItem
func (r WishlistRepository) wishlistItemToDomain(pb *wishlistspb.WishlistItem) *models.WishlistItem {
	if pb == nil {
		return nil
	}

	return &models.WishlistItem{
		ID:         pb.GetId(),
		WishlistID: pb.GetWishlistId(),
		ItemID:     pb.GetItemId(),
		EntityType: pb.GetEntityType(),
		Notes:      pb.GetNotes(),
	}
}

// dial sets up a gRPC connection with the microservice endpoint
// dial sets up a gRPC connection with the microservice endpoint
func (r WishlistRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r WishlistRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
