package tools

import (
	"context"
	"fmt"
	"strings"
)

// ==================== NEWSLETTER HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeNewsletterHandlers() {
	r.handlers["newsletter_subscribe"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		newsletterID := getStringParam(params, "newsletter_id")
		preferences := getStringArrayParam(params, "preferences")
		if userID == "" || newsletterID == "" {
			return nil, fmt.Errorf("user_id and newsletter_id are required")
		}
		// Convert array to comma-separated string
		preferencesStr := "{}"
		if len(preferences) > 0 {
			preferencesStr = strings.Join(preferences, ",")
		}
		return reg.newsletterRepo.SubscribeNewsletter(ctx, userID, newsletterID, preferencesStr)
	}

	r.handlers["newsletter_unsubscribe"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		subscriptionID := getStringParam(params, "subscription_id")
		if subscriptionID == "" {
			return nil, fmt.Errorf("subscription_id is required")
		}
		return nil, reg.newsletterRepo.UnsubscribeNewsletter(ctx, subscriptionID)
	}

	r.handlers["newsletter_get_subscription"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		subscriptionID := getStringParam(params, "subscription_id")
		if subscriptionID == "" {
			return nil, fmt.Errorf("subscription_id is required")
		}
		return reg.newsletterRepo.GetSubscription(ctx, subscriptionID)
	}

	r.handlers["newsletter_list_subscriptions"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		subscriptionStatus := getStringParam(params, "subscription_status", "")
		page := int32(getInt64Param(params, "page", 1))
		limit := int32(getInt64Param(params, "limit", 20))
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.newsletterRepo.ListSubscriptions(ctx, userID, subscriptionStatus, page, limit)
	}

	r.handlers["newsletter_update_subscription"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		subscriptionID := getStringParam(params, "subscription_id")
		subscriptionPreferences := getStringParam(params, "subscription_preferences")
		subscriptionStatus := getStringParam(params, "subscription_status")
		if subscriptionID == "" {
			return nil, fmt.Errorf("subscription_id is required")
		}
		return reg.newsletterRepo.UpdateSubscription(ctx, subscriptionID, subscriptionPreferences, subscriptionStatus)
	}

	r.handlers["newsletter_send"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		newsletterID := getStringParam(params, "newsletter_id")
		content := getStringParam(params, "content")
		if newsletterID == "" || content == "" {
			return nil, fmt.Errorf("newsletter_id and content are required")
		}
		return reg.newsletterRepo.SendNewsletter(ctx, newsletterID, content)
	}

	// Additional handlers using methods from the repository interface
	r.handlers["newsletter_get_subscription_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		subscriptionID := getStringParam(params, "subscription_id")
		if subscriptionID == "" {
			return nil, fmt.Errorf("subscription_id is required")
		}
		return reg.newsletterRepo.GetSubscriptionByID(ctx, subscriptionID)
	}

	r.handlers["newsletter_get_subscriptions_by_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := int32(getInt64Param(params, "limit", 50))
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.newsletterRepo.GetSubscriptionsByUser(ctx, userID, limit)
	}

	r.handlers["newsletter_get_subscriptions_by_newsletter"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		newsletterID := getStringParam(params, "newsletter_id")
		limit := int32(getInt64Param(params, "limit", 50))
		if newsletterID == "" {
			return nil, fmt.Errorf("newsletter_id is required")
		}
		return reg.newsletterRepo.GetSubscriptionsByNewsletter(ctx, newsletterID, limit)
	}

	r.handlers["newsletter_get_active_subscriptions"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := int32(getInt64Param(params, "limit", 50))
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.newsletterRepo.GetActiveSubscriptions(ctx, userID, limit)
	}

	r.handlers["newsletter_get_inactive_subscriptions"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		limit := int32(getInt64Param(params, "limit", 50))
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		return reg.newsletterRepo.GetInactiveSubscriptions(ctx, userID, limit)
	}

	r.handlers["newsletter_search_subscriptions"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		query := getStringParam(params, "query")
		limit := int32(getInt64Param(params, "limit", 50))
		if userID == "" || query == "" {
			return nil, fmt.Errorf("user_id and query are required")
		}
		return reg.newsletterRepo.SearchSubscriptions(ctx, userID, query, limit)
	}

	r.handlers["newsletter_count_subscriptions"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		subscriptionStatus := getStringParam(params, "subscription_status", "")
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		count, err := reg.newsletterRepo.CountSubscriptions(ctx, userID, subscriptionStatus)
		if err != nil {
			return nil, err
		}
		return map[string]int32{"count": count}, nil
	}

	r.handlers["newsletter_count_active_subscriptions"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if userID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
		count, err := reg.newsletterRepo.CountActiveSubscriptions(ctx, userID)
		if err != nil {
			return nil, err
		}
		return map[string]int32{"count": count}, nil
	}

	r.handlers["newsletter_count_subscriptions_by_newsletter"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		newsletterID := getStringParam(params, "newsletter_id")
		if newsletterID == "" {
			return nil, fmt.Errorf("newsletter_id is required")
		}
		count, err := reg.newsletterRepo.CountSubscriptionsByNewsletter(ctx, newsletterID)
		if err != nil {
			return nil, err
		}
		return map[string]int32{"count": count}, nil
	}

	r.handlers["newsletter_get"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		newsletterID := getStringParam(params, "newsletter_id")
		if newsletterID == "" {
			return nil, fmt.Errorf("newsletter_id is required")
		}
		return reg.newsletterRepo.GetNewsletter(ctx, newsletterID)
	}

	r.handlers["newsletter_get_newsletters"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		status := getStringParam(params, "status", "")
		limit := int32(getInt64Param(params, "limit", 50))
		return reg.newsletterRepo.GetNewsletters(ctx, status, limit)
	}

	r.handlers["newsletter_search_newsletters"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		query := getStringParam(params, "query")
		limit := int32(getInt64Param(params, "limit", 50))
		if query == "" {
			return nil, fmt.Errorf("query is required")
		}
		return reg.newsletterRepo.SearchNewsletters(ctx, query, limit)
	}
}