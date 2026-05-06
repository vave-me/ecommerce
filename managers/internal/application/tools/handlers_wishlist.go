package tools

import (
	"context"
	"fmt"
)

// ==================== WISHLIST HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeWishlistHandlers() {
	// Wishlist management
	r.handlers["wishlist_create_new"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		wishlistID := getStringParam(params, "wishlist_id")
		name := getStringParam(params, "name")
		if wishlistID == "" || name == "" {
			return nil, fmt.Errorf("wishlist_id and name are required")
		}
		return nil, reg.wishlistRepo.CreateNewWishlist(ctx, wishlistID, name)
	}

	r.handlers["wishlist_get_by_name"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		name := getStringParam(params, "name")
		if name == "" {
			return nil, fmt.Errorf("name is required")
		}
		return reg.wishlistRepo.GetWishlistByName(ctx, name)
	}

	r.handlers["wishlist_get_all"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.wishlistRepo.GetAllWishlists(ctx)
	}

	r.handlers["wishlist_delete"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		wishlistID := getStringParam(params, "wishlist_id")
		if wishlistID == "" {
			return nil, fmt.Errorf("wishlist_id is required")
		}
		return nil, reg.wishlistRepo.DeleteWishlist(ctx, wishlistID)
	}

	// Wishlist item management
	r.handlers["wishlist_add_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		wishlistItemID := getStringParam(params, "wishlist_item_id")
		wishlistID := getStringParam(params, "wishlist_id")
		itemID := getStringParam(params, "item_id")
		entityType := getStringParam(params, "entity_type")
		if wishlistItemID == "" || wishlistID == "" || itemID == "" || entityType == "" {
			return nil, fmt.Errorf("wishlist_item_id, wishlist_id, item_id, and entity_type are required")
		}
		return nil, reg.wishlistRepo.AddItemToWishlist(ctx, wishlistItemID, wishlistID, itemID, entityType)
	}

	r.handlers["wishlist_remove_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		wishlistItemID := getStringParam(params, "wishlist_item_id")
		if wishlistItemID == "" {
			return nil, fmt.Errorf("wishlist_item_id is required")
		}
		return nil, reg.wishlistRepo.RemoveItemFromWishlist(ctx, wishlistItemID)
	}

	r.handlers["wishlist_get_item_details"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		wishlistItemID := getStringParam(params, "wishlist_item_id")
		wishlistID := getStringParam(params, "wishlist_id")
		itemID := getStringParam(params, "item_id")
		if wishlistItemID == "" || wishlistID == "" || itemID == "" {
			return nil, fmt.Errorf("wishlist_item_id, wishlist_id, and item_id are required")
		}
		return reg.wishlistRepo.GetWishlistItemDetails(ctx, wishlistItemID, wishlistID, itemID)
	}

	r.handlers["wishlist_get_all_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		wishlistID := getStringParam(params, "wishlist_id")
		if wishlistID == "" {
			return nil, fmt.Errorf("wishlist_id is required")
		}
		return reg.wishlistRepo.GetAllItemsInWishlist(ctx, wishlistID)
	}

	r.handlers["wishlist_get_item_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		wishlistItemID := getStringParam(params, "wishlist_item_id")
		if wishlistItemID == "" {
			return nil, fmt.Errorf("wishlist_item_id is required")
		}
		return reg.wishlistRepo.GetWishlistItemByID(ctx, wishlistItemID)
	}

	r.handlers["wishlist_get_all_user_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.wishlistRepo.GetAllWishlistItemsForUser(ctx)
	}

	r.handlers["wishlist_find_by_name_detailed"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		name := getStringParam(params, "name")
		if name == "" {
			return nil, fmt.Errorf("name is required")
		}
		return reg.wishlistRepo.FindWishlistByNameDetailed(ctx, name)
	}

	r.handlers["wishlist_get_all_user_detailed"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.wishlistRepo.GetAllUserWishlistsDetailed(ctx)
	}

	// User convenience methods
	r.handlers["wishlist_add_to_user_default"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		itemType := getStringParam(params, "item_type")
		if itemID == "" || itemType == "" {
			return nil, fmt.Errorf("item_id and item_type are required")
		}
		return nil, reg.wishlistRepo.AddItemToUserDefaultWishlist(ctx, itemID, itemType)
	}

	r.handlers["wishlist_remove_from_user_default"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		if itemID == "" {
			return nil, fmt.Errorf("item_id is required")
		}
		return nil, reg.wishlistRepo.RemoveItemFromUserDefaultWishlist(ctx, itemID)
	}

	r.handlers["wishlist_get_user_default"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.wishlistRepo.GetUserDefaultWishlist(ctx)
	}

	r.handlers["wishlist_get_user_with_limit"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		limit := int32(getInt64Param(params, "limit", 10))
		return reg.wishlistRepo.GetUserWishlistsWithLimit(ctx, limit)
	}

	r.handlers["wishlist_clear_user_default"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return nil, reg.wishlistRepo.ClearUserDefaultWishlist(ctx)
	}

	r.handlers["wishlist_check_item_in_wishlist"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		if itemID == "" {
			return nil, fmt.Errorf("item_id is required")
		}
		exists, err := reg.wishlistRepo.CheckIfItemInWishlist(ctx, itemID)
		if err != nil {
			return nil, err
		}
		return map[string]bool{"exists": exists}, nil
	}

	r.handlers["wishlist_get_total_item_count"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		count, err := reg.wishlistRepo.GetTotalWishlistItemCount(ctx)
		if err != nil {
			return nil, err
		}
		return map[string]int{"count": count}, nil
	}
}