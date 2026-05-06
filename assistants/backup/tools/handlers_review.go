package tools

import (
	"context"
	"fmt"
)

// ==================== REVIEW HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeReviewHandlers() {
	r.handlers["review_create_new"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		senderID := getStringParam(params, "sender_id")
		itemID := getStringParam(params, "item_id")
		itemType := getStringParam(params, "item_type")
		content := getStringParam(params, "content")
		categoryID := getStringParam(params, "category_id")
		parentID := getStringParam(params, "parent_id")
		if senderID == "" || itemID == "" || itemType == "" || content == "" {
			return nil, fmt.Errorf("sender_id, item_id, item_type, and content are required")
		}
		return reg.reviewRepo.CreateNewReview(ctx, senderID, itemID, itemType, content, categoryID, parentID)
	}

	r.handlers["review_get_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		reviewID := getStringParam(params, "review_id")
		if reviewID == "" {
			return nil, fmt.Errorf("review_id is required")
		}
		return reg.reviewRepo.GetReviewByID(ctx, reviewID)
	}

	r.handlers["review_get_all_for_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		if itemID == "" {
			return nil, fmt.Errorf("item_id is required")
		}
		return reg.reviewRepo.GetAllReviewsForItem(ctx, itemID)
	}

	r.handlers["review_edit_content"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		reviewID := getStringParam(params, "review_id")
		content := getStringParam(params, "content")
		if reviewID == "" || content == "" {
			return nil, fmt.Errorf("review_id and content are required")
		}
		return nil, reg.reviewRepo.EditReviewContent(ctx, reviewID, content)
	}

	r.handlers["review_delete_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		reviewID := getStringParam(params, "review_id")
		if reviewID == "" {
			return nil, fmt.Errorf("review_id is required")
		}
		return nil, reg.reviewRepo.DeleteReviewByID(ctx, reviewID)
	}

	r.handlers["review_approve_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		reviewID := getStringParam(params, "review_id")
		if reviewID == "" {
			return nil, fmt.Errorf("review_id is required")
		}
		return reg.reviewRepo.ApproveReviewByID(ctx, reviewID)
	}

	r.handlers["review_reject_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		reviewID := getStringParam(params, "review_id")
		if reviewID == "" {
			return nil, fmt.Errorf("review_id is required")
		}
		return reg.reviewRepo.RejectReviewByID(ctx, reviewID)
	}

	r.handlers["review_flag_as_inappropriate"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		reviewID := getStringParam(params, "review_id")
		if reviewID == "" {
			return nil, fmt.Errorf("review_id is required")
		}
		return reg.reviewRepo.FlagReviewAsInappropriate(ctx, reviewID)
	}

	r.handlers["review_get_by_sender_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		senderID := getStringParam(params, "sender_id")
		if senderID == "" {
			return nil, fmt.Errorf("sender_id is required")
		}
		return reg.reviewRepo.GetReviewsBySenderID(ctx, senderID)
	}

	r.handlers["review_get_all_approved"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.reviewRepo.GetAllApprovedReviews(ctx)
	}

	r.handlers["review_get_most_reviewed_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.reviewRepo.GetMostReviewedItems(ctx)
	}

	r.handlers["review_get_most_reviewed_by_category"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categoryID := getStringParam(params, "category_id")
		offset := getInt64Param(params, "offset", 0)
		limit := getInt64Param(params, "limit", 10)
		if categoryID == "" {
			return nil, fmt.Errorf("category_id is required")
		}
		return reg.reviewRepo.GetMostReviewedItemsByCategory(ctx, categoryID, offset, limit)
	}

	r.handlers["review_update_content_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		reviewID := getStringParam(params, "review_id")
		content := getStringParam(params, "content")
		if reviewID == "" || content == "" {
			return nil, fmt.Errorf("review_id and content are required")
		}
		return nil, reg.reviewRepo.UpdateReviewContentByID(ctx, reviewID, content)
	}

	r.handlers["review_remove_permanently"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		reviewID := getStringParam(params, "review_id")
		if reviewID == "" {
			return nil, fmt.Errorf("review_id is required")
		}
		return nil, reg.reviewRepo.RemoveReviewPermanently(ctx, reviewID)
	}

	r.handlers["review_unflag_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		reviewID := getStringParam(params, "review_id")
		if reviewID == "" {
			return nil, fmt.Errorf("review_id is required")
		}
		return nil, reg.reviewRepo.UnflagReviewByID(ctx, reviewID)
	}

	// Additional handlers using methods from the repository interface
	r.handlers["review_search_by_keyword"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		query := getStringParam(params, "query")
		limit := getInt64Param(params, "limit", 50)
		if query == "" {
			return nil, fmt.Errorf("query is required")
		}
		return reg.reviewRepo.SearchReviewsByKeyword(ctx, query, limit)
	}

	r.handlers["review_find_by_id_and_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		reviewID := getStringParam(params, "review_id")
		itemID := getStringParam(params, "item_id")
		if reviewID == "" || itemID == "" {
			return nil, fmt.Errorf("review_id and item_id are required")
		}
		return reg.reviewRepo.FindReviewByIDAndItem(ctx, reviewID, itemID)
	}

	r.handlers["review_get_all_by_item_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		if itemID == "" {
			return nil, fmt.Errorf("item_id is required")
		}
		return reg.reviewRepo.GetAllReviewsByItemID(ctx, itemID)
	}

	r.handlers["review_find_by_sender_user_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		senderID := getStringParam(params, "sender_id")
		if senderID == "" {
			return nil, fmt.Errorf("sender_id is required")
		}
		return reg.reviewRepo.FindReviewsBySenderUserID(ctx, senderID)
	}
}