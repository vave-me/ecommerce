package tools

import (
	"context"
)

// ==================== ACTIVITY HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeActivityHandlers() {
	r.handlers["activity_create_user_activity"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, err
		}
		
		return reg.activityRepo.CreateUserActivity(ctx, userID)
	}

	r.handlers["activity_get_by_user_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, err
		}
		
		return reg.activityRepo.GetActivityByUserID(ctx, userID)
	}

	r.handlers["activity_delete_all_user_activities"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		activityID := getStringParam(params, "activity_id")
		
		if err := ValidateIDParam("activity_id", activityID); err != nil {
			return nil, err
		}
		
		return nil, reg.activityRepo.DeleteAllUserActivities(ctx, activityID)
	}

	r.handlers["activity_archive_user_activity"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		activityID := getStringParam(params, "activity_id")
		reason := getStringParam(params, "reason")
		
		v := NewValidator()
		v.ValidateRequired("activity_id", activityID)
		v.ValidateRequired("reason", reason).ValidateMaxLength("reason", reason, 500)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		reason = SanitizeString(reason)
		return nil, reg.activityRepo.ArchiveUserActivity(ctx, activityID, reason)
	}

	r.handlers["activity_restore_user_activity"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		activityID := getStringParam(params, "activity_id")
		reason := getStringParam(params, "reason")
		
		v := NewValidator()
		v.ValidateRequired("activity_id", activityID)
		v.ValidateRequired("reason", reason).ValidateMaxLength("reason", reason, 500)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		reason = SanitizeString(reason)
		return nil, reg.activityRepo.RestoreUserActivity(ctx, activityID, reason)
	}

	r.handlers["activity_add_user_interaction"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		interactionID := getStringParam(params, "interaction_id")
		activityID := getStringParam(params, "activity_id")
		itemID := getStringParam(params, "item_id")
		itemType := getStringParam(params, "item_type")
		actionType := getStringParam(params, "action_type")
		
		v := NewValidator()
		v.ValidateRequired("interaction_id", interactionID)
		v.ValidateRequired("activity_id", activityID)
		v.ValidateRequired("item_id", itemID)
		v.ValidateRequired("item_type", itemType)
		v.ValidateRequired("action_type", actionType)
		
		// Validate item type
		if itemType != "" {
			v.ValidateEnum("item_type", itemType, []string{"product", "service", "post", "comment", "review", "user"})
		}
		
		// Validate action type
		if actionType != "" {
			v.ValidateEnum("action_type", actionType, []string{"like", "dislike", "view", "share", "save", "report", "comment"})
		}
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return nil, reg.activityRepo.AddUserInteraction(ctx, interactionID, activityID, itemID, itemType, actionType)
	}

	r.handlers["activity_update_user_interaction"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		interactionID := getStringParam(params, "interaction_id")
		actionType := getStringParam(params, "action_type")
		
		v := NewValidator()
		v.ValidateRequired("interaction_id", interactionID)
		v.ValidateRequired("action_type", actionType)
		
		// Validate action type
		if actionType != "" {
			v.ValidateEnum("action_type", actionType, []string{"like", "dislike", "view", "share", "save", "report", "comment"})
		}
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return nil, reg.activityRepo.UpdateUserInteraction(ctx, interactionID, actionType)
	}

	r.handlers["activity_delete_user_interaction"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		interactionID := getStringParam(params, "interaction_id")
		
		if err := ValidateIDParam("interaction_id", interactionID); err != nil {
			return nil, err
		}
		
		return nil, reg.activityRepo.DeleteUserInteraction(ctx, interactionID)
	}

	r.handlers["activity_get_interaction_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		interactionID := getStringParam(params, "interaction_id")
		
		if err := ValidateIDParam("interaction_id", interactionID); err != nil {
			return nil, err
		}
		
		return reg.activityRepo.GetInteractionByID(ctx, interactionID)
	}

	r.handlers["activity_get_all_activity_interactions"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		activityID := getStringParam(params, "activity_id")
		
		if err := ValidateIDParam("activity_id", activityID); err != nil {
			return nil, err
		}
		
		return reg.activityRepo.GetAllActivityInteractions(ctx, activityID)
	}

	r.handlers["activity_get_most_liked_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemType := getStringParam(params, "item_type")
		limit := getInt64Param(params, "limit", 10)
		
		v := NewValidator()
		if itemType != "" {
			v.ValidateEnum("item_type", itemType, []string{"product", "service", "post", "comment", "review", "user"})
		}
		v.ValidateRange("limit", float64(limit), 1, 100)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.activityRepo.GetMostLikedItems(ctx, itemType, limit)
	}

	r.handlers["activity_get_most_disliked_items"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		itemType := getStringParam(params, "item_type")
		limit := getInt64Param(params, "limit", 10)
		
		v := NewValidator()
		if itemType != "" {
			v.ValidateEnum("item_type", itemType, []string{"product", "service", "post", "comment", "review", "user"})
		}
		v.ValidateRange("limit", float64(limit), 1, 100)
		
		if err := v.GetError(); err != nil {
			return nil, err
		}
		
		return reg.activityRepo.GetMostDislikedItems(ctx, itemType, limit)
	}
}