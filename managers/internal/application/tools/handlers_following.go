package tools

import (
	"context"
	"fmt"
)

// ==================== FOLLOWING HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeFollowingHandlers() {
	r.handlers["following_create_new_relationship"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		followedUserID := getStringParam(params, "followed_user_id")
		followedUserType := getStringParam(params, "followed_user_type")
		content := getStringParam(params, "content")
		categoryID := getStringParam(params, "category_id")
		parentID := getStringParam(params, "parent_id")

		// Validate required parameters
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if err := ValidateIDParam("followed_user_id", followedUserID); err != nil {
			return nil, fmt.Errorf("invalid followed_user_id: %w", err)
		}
		if followedUserType == "" {
			return nil, fmt.Errorf("followed_user_type is required")
		}

		// Sanitize string inputs
		content = SanitizeString(content)

		return reg.followingRepo.CreateNewFollowRelationship(ctx, userID, followedUserID, followedUserType, content, categoryID, parentID)
	}

	r.handlers["following_approve_request_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid follow request id: %w", err)
		}
		return reg.followingRepo.ApproveFollowRequestByID(ctx, id)
	}

	r.handlers["following_reject_request_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid follow request id: %w", err)
		}
		return reg.followingRepo.RejectFollowRequestByID(ctx, id)
	}

	r.handlers["following_flag_as_inappropriate"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid follow id: %w", err)
		}
		return reg.followingRepo.FlagFollowAsInappropriate(ctx, id)
	}

	r.handlers["following_edit_description"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		content := getStringParam(params, "content")
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid follow id: %w", err)
		}
		content = SanitizeString(content)
		return reg.followingRepo.EditFollowDescription(ctx, id, content)
	}

	r.handlers["following_delete_relationship"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid follow id: %w", err)
		}
		return reg.followingRepo.DeleteFollowRelationship(ctx, id)
	}

	r.handlers["following_get_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid follow id: %w", err)
		}
		return reg.followingRepo.GetFollowByID(ctx, id)
	}

	r.handlers["following_get_all_followers_for_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		followedUserID := getStringParam(params, "followed_user_id")
		if err := ValidateIDParam("followed_user_id", followedUserID); err != nil {
			return nil, fmt.Errorf("invalid followed_user_id: %w", err)
		}
		return reg.followingRepo.GetAllFollowersForUser(ctx, followedUserID)
	}

	r.handlers["following_get_most_followed_users"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.followingRepo.GetMostFollowedUsers(ctx)
	}

	r.handlers["following_get_most_followed_by_category"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		categoryID := getStringParam(params, "category_id")
		offset := getInt64Param(params, "offset", 0)
		limit := getInt64Param(params, "limit", 20)
		if err := ValidateIDParam("category_id", categoryID); err != nil {
			return nil, fmt.Errorf("invalid category_id: %w", err)
		}
		if err := ValidatePaginationParams(offset+1, limit); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.followingRepo.GetMostFollowedUsersByCategory(ctx, categoryID, offset, limit)
	}

	r.handlers["following_get_user_following_list"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return reg.followingRepo.GetUserFollowingList(ctx, userID)
	}

	r.handlers["following_get_all_approved_relationships"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.followingRepo.GetAllApprovedFollowRelationships(ctx)
	}

	r.handlers["following_check_if_user_is_following"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		followedUserID := getStringParam(params, "followed_user_id")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if err := ValidateIDParam("followed_user_id", followedUserID); err != nil {
			return nil, fmt.Errorf("invalid followed_user_id: %w", err)
		}
		return reg.followingRepo.CheckIfUserIsFollowing(ctx, userID, followedUserID)
	}

	r.handlers["following_get_followers_with_limit"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := int32(getInt64Param(params, "limit", 20))
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if err := ValidatePaginationParams(1, int64(limit)); err != nil {
			return nil, fmt.Errorf("invalid limit: %w", err)
		}
		return reg.followingRepo.GetFollowersWithLimit(ctx, userID, limit)
	}

	r.handlers["following_get_total_follower_count"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return reg.followingRepo.GetTotalFollowerCount(ctx, userID)
	}

	r.handlers["following_get_total_following_count"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return reg.followingRepo.GetTotalFollowingCount(ctx, userID)
	}

	r.handlers["following_get_mutual_between_users"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID1 := getStringParam(params, "user_id_1")
		userID2 := getStringParam(params, "user_id_2")
		if err := ValidateIDParam("user_id_1", userID1); err != nil {
			return nil, fmt.Errorf("invalid user_id_1: %w", err)
		}
		if err := ValidateIDParam("user_id_2", userID2); err != nil {
			return nil, fmt.Errorf("invalid user_id_2: %w", err)
		}
		return reg.followingRepo.GetMutualFollowingBetweenUsers(ctx, userID1, userID2)
	}

	r.handlers["following_search_by_keyword"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		keyword := getStringParam(params, "keyword")
		if keyword == "" {
			return nil, fmt.Errorf("keyword is required for search")
		}
		keyword = SanitizeString(keyword)
		return reg.followingRepo.SearchFollowsByKeyword(ctx, keyword)
	}

	r.handlers["following_get_paginated_list"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		sortBy := getStringParam(params, "sort_by")
		sortOrder := getStringParam(params, "sort_order")
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.followingRepo.GetPaginatedFollowList(ctx, page, pageSize, sortBy, sortOrder)
	}

	r.handlers["following_get_statistics"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return reg.followingRepo.GetFollowingStatistics(ctx, userID)
	}

	r.handlers["following_remove_relationship"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		followID := getStringParam(params, "follow_id")
		if err := ValidateIDParam("follow_id", followID); err != nil {
			return nil, fmt.Errorf("invalid follow_id: %w", err)
		}
		return nil, reg.followingRepo.RemoveFollowRelationship(ctx, followID)
	}
}