package tools

import (
	"context"
)

// ==================== COMMENT HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeCommentHandlers() {
	r.handlers["comment_create_new"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		itemType := getStringParam(params, "item_type")
		content := getStringParam(params, "content")
		categoryID := getStringParam(params, "category_id")
		parentID := getStringParam(params, "parent_id")
		
		v := NewValidator()
		v.ValidateRequired("item_id", itemID)
		v.ValidateRequired("item_type", itemType)
		v.ValidateRequired("content", content).ValidateMinLength("content", content, 1).ValidateMaxLength("content", content, 5000)
		
		// Validate item type
		if itemType != "" {
			v.ValidateEnum("item_type", itemType, []string{"product", "service", "post", "order", "user"})
		}
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		content = SanitizeString(content)
		return reg.commentRepo.CreateNewComment(ctx, itemID, itemType, content, categoryID, parentID)
	}

	r.handlers["comment_get_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		commentID := getStringParam(params, "comment_id")
		
		if err := ValidateIDParam("comment_id", commentID); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.GetCommentByID(ctx, commentID)
	}

	r.handlers["comment_get_all_for_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.GetAllCommentsForItem(ctx, itemID)
	}

	r.handlers["comment_edit_content"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		commentID := getStringParam(params, "comment_id")
		content := getStringParam(params, "content")
		
		v := NewValidator()
		v.ValidateRequired("comment_id", commentID)
		v.ValidateRequired("content", content).ValidateMinLength("content", content, 1).ValidateMaxLength("content", content, 5000)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		content = SanitizeString(content)
		return reg.commentRepo.EditCommentContent(ctx, commentID, content)
	}

	r.handlers["comment_delete_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		commentID := getStringParam(params, "comment_id")
		
		if err := ValidateIDParam("comment_id", commentID); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.DeleteCommentByID(ctx, commentID)
	}

	r.handlers["comment_approve_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		commentID := getStringParam(params, "comment_id")
		
		if err := ValidateIDParam("comment_id", commentID); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.ApproveCommentByID(ctx, commentID)
	}

	r.handlers["comment_reject_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		commentID := getStringParam(params, "comment_id")
		
		if err := ValidateIDParam("comment_id", commentID); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.RejectCommentByID(ctx, commentID)
	}

	r.handlers["comment_flag_as_inappropriate"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		commentID := getStringParam(params, "comment_id")
		
		if err := ValidateIDParam("comment_id", commentID); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.FlagCommentAsInappropriate(ctx, commentID)
	}

	r.handlers["comment_get_by_sender_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		// Note: This method doesn't take parameters in the interface
		return reg.commentRepo.GetCommentsBySenderUser(ctx)
	}

	r.handlers["comment_get_all_approved"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.commentRepo.GetAllApprovedComments(ctx)
	}

	r.handlers["comment_get_most_commented_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.commentRepo.GetMostCommentedItems(ctx)
	}

	r.handlers["comment_get_most_commented_by_category"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categoryID := getStringParam(params, "category_id")
		offset := getInt64Param(params, "offset", 0)
		limit := getInt64Param(params, "limit", 20)
		
		v := NewValidator()
		v.ValidateRequired("category_id", categoryID)
		v.ValidateMinimum("offset", float64(offset), 0)
		v.ValidateRange("limit", float64(limit), 1, 100)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.GetMostCommentedItemsByCategory(ctx, categoryID, offset, limit)
	}

	r.handlers["comment_search_by_keyword"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		keyword := getStringParam(params, "keyword")
		
		v := NewValidator()
		v.ValidateRequired("keyword", keyword).ValidateMinLength("keyword", keyword, 2)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		keyword = SanitizeString(keyword)
		return reg.commentRepo.SearchCommentsByKeyword(ctx, keyword)
	}

	r.handlers["comment_get_paginated_list"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by", "created_at")
		sortOrder := getStringParam(params, "sort_order", "desc")
		
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, err
		}
		
		v := NewValidator()
		if sortBy != "" {
			v.ValidateEnum("sort_by", sortBy, []string{"created_at", "updated_at", "likes", "replies"})
		}
		if sortOrder != "" {
			v.ValidateEnum("sort_order", sortOrder, []string{"asc", "desc"})
		}
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.GetPaginatedCommentList(ctx, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["comment_get_by_approval_status"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		status := getStringParam(params, "status")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		
		v := NewValidator()
		v.ValidateRequired("status", status)
		v.ValidateEnum("status", status, []string{"pending", "approved", "rejected", "flagged"})
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.GetCommentsByApprovalStatus(ctx, status, page, pageSize)
	}

	r.handlers["comment_get_by_category_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categoryID := getStringParam(params, "category_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		
		if err := ValidateIDParam("category_id", categoryID); err != nil {
			return nil, err
		}
		
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.GetCommentsByCategoryID(ctx, categoryID, page, pageSize)
	}

	r.handlers["comment_get_statistics"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.GetCommentStatistics(ctx, itemID)
	}

	r.handlers["comment_get_recently_posted"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.GetRecentlyPostedComments(ctx, page, pageSize)
	}

	r.handlers["comment_get_child_comments"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		parentID := getStringParam(params, "parent_id")
		
		if err := ValidateIDParam("parent_id", parentID); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.GetChildCommentsForParent(ctx, parentID)
	}

	// Simplified methods for backward compatibility
	r.handlers["comment_find_by_item_and_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		commentID := getStringParam(params, "comment_id")
		itemID := getStringParam(params, "item_id")
		
		v := NewValidator()
		v.ValidateRequired("comment_id", commentID)
		v.ValidateRequired("item_id", itemID)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.FindCommentByItemAndID(ctx, commentID, itemID)
	}

	r.handlers["comment_get_all_by_item_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemID := getStringParam(params, "item_id")
		
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, err
		}
		
		return reg.commentRepo.GetAllCommentsByItemID(ctx, itemID)
	}

	r.handlers["comment_find_all_by_sender_user_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		// Note: This method doesn't take parameters in the interface
		return reg.commentRepo.FindAllCommentsBySenderUserID(ctx)
	}
}