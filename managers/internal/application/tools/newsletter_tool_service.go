package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// NewsletterToolService handles all newsletter and subscription operations
type NewsletterToolService struct {
	newsletters domain.NewsletterRepository
}

// NewNewsletterToolService creates a new newsletter tool service
func NewNewsletterToolService(newsletters domain.NewsletterRepository) *NewsletterToolService {
	return &NewsletterToolService{
		newsletters: newsletters,
	}
}

// GetSupportedOperations returns all operations supported by this service
func (n *NewsletterToolService) GetSupportedOperations() []string {
	return []string{
		"subscribe_newsletter",
		"unsubscribe_newsletter",
		"get_subscription", "find",
		"list_subscriptions",
		"update_subscription",
		"send_newsletter",
		"get_subscription_by_id",
		"get_subscriptions_by_user", "search",
		"get_subscriptions_by_newsletter",
		"get_active_subscriptions",
		"get_inactive_subscriptions",
		"search_subscriptions",
		"count_subscriptions",
		"count_active_subscriptions",
		"count_subscriptions_by_newsletter",
		"get_newsletter",
		"get_newsletters",
		"search_newsletters",
	}
}

// ExecuteOperation executes a newsletter operation with streaming progress updates
func (n *NewsletterToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	log.Printf("NewsletterToolService.ExecuteOperation: Executing newsletter operation: %s", operation)

	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "newsletter_operation",
		Status:   "progress",
		Progress: 10,
		Metadata: map[string]interface{}{
			"operation": operation,
			"message":   fmt.Sprintf("Processing newsletter operation: %s", operation),
		},
		Timestamp: time.Now().Unix(),
	}

	// Extract common parameters
	userID := getStringParam(parameters, "user_id", "")
	subscriptionID := getStringParam(parameters, "subscription_id", "")
	if subscriptionID == "" {
		subscriptionID = getStringParam(parameters, "id", "")
	}
	newsletterID := getStringParam(parameters, "newsletter_id", "")
	subscriptionPreferences := getStringParam(parameters, "subscription_preferences", "")
	subscriptionStatus := getStringParam(parameters, "subscription_status", "")
	content := getStringParam(parameters, "content", "")
	if content == "" {
		content = getStringParam(parameters, "newsletter_content", "")
	}
	limit := int32(getInt64Param(parameters, "limit", 20))
	page := int32(getInt64Param(parameters, "page", 1))
	query := getStringParam(parameters, "query", "")
	if query == "" {
		query = getStringParam(parameters, "search_term", "")
	}
	status := getStringParam(parameters, "status", "")
	if status == "" {
		status = getStringParam(parameters, "newsletter_status", "")
	}

	var result interface{}
	var err error

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "newsletter_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"operation": operation,
			"step":      "executing_operation",
		},
		Timestamp: time.Now().Unix(),
	}

	switch operation {
	case "subscribe_newsletter":
		if userID == "" {
			return nil, fmt.Errorf("user_id parameter required")
		}
		if newsletterID == "" {
			return nil, fmt.Errorf("newsletter_id parameter required")
		}
		log.Printf("NewsletterToolService: Subscribing user %s to newsletter %s", userID, newsletterID)
		result, err = n.newsletters.SubscribeNewsletter(ctx, userID, newsletterID, subscriptionPreferences)

	case "unsubscribe_newsletter":
		if subscriptionID == "" {
			return nil, fmt.Errorf("subscription_id or id parameter required")
		}
		log.Printf("NewsletterToolService: Unsubscribing subscription %s", subscriptionID)
		err = n.newsletters.UnsubscribeNewsletter(ctx, subscriptionID)
		result = map[string]interface{}{"success": true, "subscription_id": subscriptionID}

	case "get_subscription", "find":
		if subscriptionID == "" {
			return nil, fmt.Errorf("subscription_id or id parameter required")
		}
		log.Printf("NewsletterToolService: Getting subscription %s", subscriptionID)
		result, err = n.newsletters.GetSubscription(ctx, subscriptionID)

	case "list_subscriptions":
		if userID == "" {
			return nil, fmt.Errorf("user_id parameter required")
		}
		log.Printf("NewsletterToolService: Listing subscriptions for user %s, status: %s", userID, subscriptionStatus)
		result, err = n.newsletters.ListSubscriptions(ctx, userID, subscriptionStatus, page, limit)

	case "update_subscription":
		if subscriptionID == "" {
			return nil, fmt.Errorf("subscription_id or id parameter required")
		}
		log.Printf("NewsletterToolService: Updating subscription %s", subscriptionID)
		result, err = n.newsletters.UpdateSubscription(ctx, subscriptionID, subscriptionPreferences, subscriptionStatus)

	case "send_newsletter":
		if newsletterID == "" {
			return nil, fmt.Errorf("newsletter_id parameter required")
		}
		if content == "" {
			return nil, fmt.Errorf("content or newsletter_content parameter required")
		}
		log.Printf("NewsletterToolService: Sending newsletter %s", newsletterID)
		result, err = n.newsletters.SendNewsletter(ctx, newsletterID, content)

	case "get_subscription_by_id":
		if subscriptionID == "" {
			return nil, fmt.Errorf("subscription_id or id parameter required")
		}
		log.Printf("NewsletterToolService: Getting subscription by ID %s", subscriptionID)
		result, err = n.newsletters.GetSubscriptionByID(ctx, subscriptionID)

	case "get_subscriptions_by_user", "search":
		if userID == "" {
			return nil, fmt.Errorf("user_id parameter required")
		}
		log.Printf("NewsletterToolService: Getting subscriptions for user %s", userID)
		result, err = n.newsletters.GetSubscriptionsByUser(ctx, userID, limit)

	case "get_subscriptions_by_newsletter":
		if newsletterID == "" {
			return nil, fmt.Errorf("newsletter_id parameter required")
		}
		log.Printf("NewsletterToolService: Getting subscriptions for newsletter %s", newsletterID)
		result, err = n.newsletters.GetSubscriptionsByNewsletter(ctx, newsletterID, limit)

	case "get_active_subscriptions":
		if userID == "" {
			return nil, fmt.Errorf("user_id parameter required")
		}
		log.Printf("NewsletterToolService: Getting active subscriptions for user %s", userID)
		result, err = n.newsletters.GetActiveSubscriptions(ctx, userID, limit)

	case "get_inactive_subscriptions":
		if userID == "" {
			return nil, fmt.Errorf("user_id parameter required")
		}
		log.Printf("NewsletterToolService: Getting inactive subscriptions for user %s", userID)
		result, err = n.newsletters.GetInactiveSubscriptions(ctx, userID, limit)

	case "search_subscriptions":
		if userID == "" {
			return nil, fmt.Errorf("user_id parameter required")
		}
		if query == "" {
			return nil, fmt.Errorf("query or search_term parameter required")
		}
		log.Printf("NewsletterToolService: Searching subscriptions for user %s, query: %s", userID, query)
		result, err = n.newsletters.SearchSubscriptions(ctx, userID, query, limit)

	case "count_subscriptions":
		if userID == "" {
			return nil, fmt.Errorf("user_id parameter required")
		}
		log.Printf("NewsletterToolService: Counting subscriptions for user %s, status: %s", userID, subscriptionStatus)
		result, err = n.newsletters.CountSubscriptions(ctx, userID, subscriptionStatus)

	case "count_active_subscriptions":
		if userID == "" {
			return nil, fmt.Errorf("user_id parameter required")
		}
		log.Printf("NewsletterToolService: Counting active subscriptions for user %s", userID)
		result, err = n.newsletters.CountActiveSubscriptions(ctx, userID)

	case "count_subscriptions_by_newsletter":
		if newsletterID == "" {
			return nil, fmt.Errorf("newsletter_id parameter required")
		}
		log.Printf("NewsletterToolService: Counting subscriptions for newsletter %s", newsletterID)
		result, err = n.newsletters.CountSubscriptionsByNewsletter(ctx, newsletterID)

	case "get_newsletter":
		if newsletterID == "" {
			return nil, fmt.Errorf("newsletter_id or id parameter required")
		}
		log.Printf("NewsletterToolService: Getting newsletter %s", newsletterID)
		result, err = n.newsletters.GetNewsletter(ctx, newsletterID)

	case "get_newsletters":
		log.Printf("NewsletterToolService: Getting newsletters with status %s", status)
		result, err = n.newsletters.GetNewsletters(ctx, status, limit)

	case "search_newsletters":
		if query == "" {
			return nil, fmt.Errorf("query or search_term parameter required")
		}
		log.Printf("NewsletterToolService: Searching newsletters with query: %s", query)
		result, err = n.newsletters.SearchNewsletters(ctx, query, limit)

	default:
		return nil, fmt.Errorf("unsupported newsletter operation: %s", operation)
	}

	if err != nil {
		// Send error update
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "newsletter_operation",
			Status:   "error",
			Error:    err.Error(),
			Metadata: map[string]interface{}{
				"operation": operation,
			},
			Timestamp: time.Now().Unix(),
		}
		return nil, fmt.Errorf("newsletter operation failed: %w", err)
	}

	// Send completion update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "newsletter_operation",
		Status:   "completed",
		Progress: 100,
		Result:   result,
		Metadata: map[string]interface{}{
			"operation": operation,
			"success":   true,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "newsletters",
		"operation":   operation,
		"result":      result,
		"success":     true,
	}, nil
}
