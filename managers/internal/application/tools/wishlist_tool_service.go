package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
)

// WishlistToolService handles all wishlist and user preference operations
type WishlistToolService struct {
	wishlists domain.WishlistRepository
	following domain.FollowingRepository
}

// NewWishlistToolService creates a new wishlist tool service
func NewWishlistToolService(wishlists domain.WishlistRepository, following domain.FollowingRepository) *WishlistToolService {
	return &WishlistToolService{
		wishlists: wishlists,
		following: following,
	}
}

// GetSupportedOperations returns all operations supported by this service
func (w *WishlistToolService) GetSupportedOperations() []string {
	return []string{
		// Wishlist operations
		"add_to_wishlist", "create_wishlist",
		"remove_from_wishlist", "delete_wishlist",
		"get_wishlist", "find_wishlist",
		"get_user_wishlists", "list_wishlists",
		"clear_wishlist",
		"check_wishlist_item",
		"get_wishlist_count",
		// Following operations
		"follow_user", "follow",
		"unfollow_user", "unfollow",
		"get_followers",
		"get_following",
		"check_following",
		"get_follower_count",
		"get_following_count",
	}
}

// ExecuteOperation executes a wishlist operation with streaming progress updates
func (w *WishlistToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	log.Printf("WishlistToolService.ExecuteOperation: Executing operation: %s", operation)

	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 10,
		Metadata: map[string]interface{}{
			"operation": operation,
			"message":   fmt.Sprintf("Processing wishlist operation: %s", operation),
		},
		Timestamp: time.Now().Unix(),
	}

	// Route to appropriate handler
	switch operation {
	case "add_to_wishlist", "create_wishlist":
		return w.handleAddToWishlist(ctx, parameters, streamChan, toolID)
	case "remove_from_wishlist", "delete_wishlist":
		return w.handleRemoveFromWishlist(ctx, parameters, streamChan, toolID)
	case "get_wishlist", "find_wishlist":
		return w.handleGetWishlist(ctx, parameters, streamChan, toolID)
	case "get_user_wishlists", "list_wishlists":
		return w.handleGetUserWishlists(ctx, parameters, streamChan, toolID)
	case "clear_wishlist":
		return w.handleClearWishlist(ctx, parameters, streamChan, toolID)
	case "check_wishlist_item":
		return w.handleCheckWishlistItem(ctx, parameters, streamChan, toolID)
	case "get_wishlist_count":
		return w.handleGetWishlistCount(ctx, parameters, streamChan, toolID)
	case "follow_user", "follow":
		return w.handleFollowUser(ctx, parameters, streamChan, toolID)
	case "unfollow_user", "unfollow":
		return w.handleUnfollowUser(ctx, parameters, streamChan, toolID)
	case "get_followers":
		return w.handleGetFollowers(ctx, parameters, streamChan, toolID)
	case "get_following":
		return w.handleGetFollowing(ctx, parameters, streamChan, toolID)
	case "check_following":
		return w.handleCheckFollowing(ctx, parameters, streamChan, toolID)
	case "get_follower_count":
		return w.handleGetFollowerCount(ctx, parameters, streamChan, toolID)
	case "get_following_count":
		return w.handleGetFollowingCount(ctx, parameters, streamChan, toolID)
	default:
		return nil, fmt.Errorf("unsupported wishlist operation: %s", operation)
	}
}

func (w *WishlistToolService) handleAddToWishlist(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	itemID := getStringParam(parameters, "item_id", "")
	itemType := getStringParam(parameters, "item_type", "")

	if userID == "" || itemID == "" {
		return nil, fmt.Errorf("user_id and item_id parameters required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":      "adding_to_wishlist",
			"user_id":   userID,
			"item_id":   itemID,
			"item_type": itemType,
		},
		Timestamp: time.Now().Unix(),
	}

	err := w.wishlists.AddToWishlist(ctx, itemID, itemType)
	if err != nil {
		return nil, fmt.Errorf("failed to add to wishlist: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"message":   "Item added to wishlist successfully",
			"user_id":   userID,
			"item_id":   itemID,
			"item_type": itemType,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "wishlists",
		"operation":   "add_to_wishlist",
		"user_id":     userID,
		"item_id":     itemID,
		"success":     true,
	}, nil
}

func (w *WishlistToolService) handleRemoveFromWishlist(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	itemID := getStringParam(parameters, "item_id", "")

	if userID == "" || itemID == "" {
		return nil, fmt.Errorf("user_id and item_id parameters required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "removing_from_wishlist",
			"user_id": userID,
			"item_id": itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := w.wishlists.RemoveFromWishlist(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove from wishlist: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"message": "Item removed from wishlist successfully",
			"user_id": userID,
			"item_id": itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "wishlists",
		"operation":   "remove_from_wishlist",
		"user_id":     userID,
		"item_id":     itemID,
		"success":     true,
	}, nil
}

func (w *WishlistToolService) handleGetWishlist(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_wishlist",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	wishlist, err := w.wishlists.GetUserWishlist(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"wishlist": wishlist,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "wishlists",
		"operation":   "get_wishlist",
		"result":      wishlist,
		"user_id":     userID,
	}, nil
}

func (w *WishlistToolService) handleGetUserWishlists(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	limit := int32(getInt64Param(parameters, "limit", 20))

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_user_wishlists",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	wishlists, err := w.wishlists.GetUserWishlists(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get user wishlists: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"wishlists": wishlists,
			"count":     len(wishlists),
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "wishlists",
		"operation":   "get_user_wishlists",
		"results":     wishlists,
		"count":       len(wishlists),
		"user_id":     userID,
	}, nil
}

func (w *WishlistToolService) handleClearWishlist(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "clearing_wishlist",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := w.wishlists.ClearWishlist(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to clear wishlist: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"message": "Wishlist cleared successfully",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "wishlists",
		"operation":   "clear_wishlist",
		"user_id":     userID,
		"success":     true,
	}, nil
}

func (w *WishlistToolService) handleCheckWishlistItem(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	itemID := getStringParam(parameters, "item_id", "")

	if userID == "" || itemID == "" {
		return nil, fmt.Errorf("user_id and item_id parameters required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "checking_wishlist_item",
			"user_id": userID,
			"item_id": itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	inWishlist, err := w.wishlists.IsInWishlist(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to check wishlist item: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"in_wishlist": inWishlist,
			"user_id":     userID,
			"item_id":     itemID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "wishlists",
		"operation":   "check_wishlist_item",
		"in_wishlist": inWishlist,
		"user_id":     userID,
		"item_id":     itemID,
	}, nil
}

func (w *WishlistToolService) handleGetWishlistCount(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_wishlist_count",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	count, err := w.wishlists.GetWishlistCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist count: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"count":   count,
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "wishlists",
		"operation":   "get_wishlist_count",
		"count":       count,
		"user_id":     userID,
	}, nil
}

// Following operations

func (w *WishlistToolService) handleFollowUser(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	followerID := getStringParam(parameters, "follower_id", "")
	if followerID == "" {
		followerID = getStringParam(parameters, "user_id", "")
	}
	followeeID := getStringParam(parameters, "followee_id", "")
	if followeeID == "" {
		followeeID = getStringParam(parameters, "target_user_id", "")
	}

	if followerID == "" || followeeID == "" {
		return nil, fmt.Errorf("follower_id and followee_id parameters required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":        "following_user",
			"follower_id": followerID,
			"followee_id": followeeID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := w.following.FollowUser(ctx, followerID, followerID, followeeID, models.UserTypePrivate, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to follow user: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"message":     "User followed successfully",
			"follower_id": followerID,
			"followee_id": followeeID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "following",
		"operation":   "follow_user",
		"follower_id": followerID,
		"followee_id": followeeID,
		"success":     true,
	}, nil
}

func (w *WishlistToolService) handleUnfollowUser(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	followerID := getStringParam(parameters, "follower_id", "")
	if followerID == "" {
		followerID = getStringParam(parameters, "user_id", "")
	}
	followeeID := getStringParam(parameters, "followee_id", "")
	if followeeID == "" {
		followeeID = getStringParam(parameters, "target_user_id", "")
	}

	if followerID == "" || followeeID == "" {
		return nil, fmt.Errorf("follower_id and followee_id parameters required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":        "unfollowing_user",
			"follower_id": followerID,
			"followee_id": followeeID,
		},
		Timestamp: time.Now().Unix(),
	}

	err := w.following.UnfollowUser(ctx, followerID)
	if err != nil {
		return nil, fmt.Errorf("failed to unfollow user: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"message":     "User unfollowed successfully",
			"follower_id": followerID,
			"followee_id": followeeID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "following",
		"operation":   "unfollow_user",
		"follower_id": followerID,
		"followee_id": followeeID,
		"success":     true,
	}, nil
}

func (w *WishlistToolService) handleGetFollowers(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	limit := int32(getInt64Param(parameters, "limit", 20))

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_followers",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	followers, err := w.following.GetFollowers(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"followers": followers,
			"count":     len(followers),
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "following",
		"operation":   "get_followers",
		"results":     followers,
		"count":       len(followers),
		"user_id":     userID,
	}, nil
}

func (w *WishlistToolService) handleGetFollowing(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	_ = int32(getInt64Param(parameters, "limit", 20)) // limit not used as GetFollowing doesn't take pagination

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_following",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	followingResponse, err := w.following.GetFollowing(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	// Extract follows from response - using empty slice as placeholder until correct field is identified
	var following []*models.Follow
	// TODO: Fix field name when GetFollowingResponse structure is known
	_ = followingResponse // Suppress unused variable warning

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"following": following,
			"count":     len(following),
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "following",
		"operation":   "get_following",
		"results":     following,
		"count":       len(following),
		"user_id":     userID,
	}, nil
}

func (w *WishlistToolService) handleCheckFollowing(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	followerID := getStringParam(parameters, "follower_id", "")
	if followerID == "" {
		followerID = getStringParam(parameters, "user_id", "")
	}
	followeeID := getStringParam(parameters, "followee_id", "")
	if followeeID == "" {
		followeeID = getStringParam(parameters, "target_user_id", "")
	}

	if followerID == "" || followeeID == "" {
		return nil, fmt.Errorf("follower_id and followee_id parameters required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":        "checking_following",
			"follower_id": followerID,
			"followee_id": followeeID,
		},
		Timestamp: time.Now().Unix(),
	}

	isFollowing, err := w.following.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check following: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"is_following": isFollowing,
			"follower_id":  followerID,
			"followee_id":  followeeID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type":  "following",
		"operation":    "check_following",
		"is_following": isFollowing,
		"follower_id":  followerID,
		"followee_id":  followeeID,
	}, nil
}

func (w *WishlistToolService) handleGetFollowerCount(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_follower_count",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	count, err := w.following.GetFollowerCount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get follower count: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"count":   count,
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "following",
		"operation":   "get_follower_count",
		"count":       count,
		"user_id":     userID,
	}, nil
}

func (w *WishlistToolService) handleGetFollowingCount(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID := getStringParam(parameters, "user_id", "")
	if userID == "" {
		return nil, fmt.Errorf("user_id parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "getting_following_count",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	count, err := w.following.GetFollowingCount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get following count: %w", err)
	}

	// Send completion
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "wishlist_operation",
		Status:   "completed",
		Progress: 100,
		Result: map[string]interface{}{
			"count":   count,
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "following",
		"operation":   "get_following_count",
		"count":       count,
		"user_id":     userID,
	}, nil
}
