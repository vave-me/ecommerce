package grpc

import (
	"context"
	"fmt"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
	"middleman/newsletters/newsletterspb"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// NewsletterRepository calls the remote newsletters service (gRPC).
type NewsletterRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.NewsletterRepository = (*NewsletterRepository)(nil)

// NewNewsletterRepositoryWithAuth creates a new NewsletterRepository with JWT authentication support
func NewNewsletterRepository(endpoint string, authInstance *auth.Auth) NewsletterRepository {
	return NewsletterRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// SubscribeNewsletter subscribes a user to a newsletter
func (r NewsletterRepository) SubscribeNewsletter(ctx context.Context, userID, newsletterID, subscriptionPreferences string) (*models.SubscribeNewsletterResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Parse preferences
	prefs := &newsletterspb.SubscriptionPreferences{}
	if subscriptionPreferences != "" {
		// Set frequency override if provided
		prefs.FrequencyOverride = subscriptionPreferences
	}

	client := newsletterspb.NewNewslettersServiceClient(conn)
	resp, err := client.Subscribe(ctx, &newsletterspb.SubscribeRequest{
		NewsletterId: newsletterID,
		Preferences:  prefs,
	})
	if err != nil {
		return nil, fmt.Errorf("Subscribe RPC failed: %w", err)
	}

	return &models.SubscribeNewsletterResponse{
		SubscriptionID: resp.GetSubscription().GetId(),
	}, nil
}

// UnsubscribeNewsletter unsubscribes a user from a newsletter
func (r NewsletterRepository) UnsubscribeNewsletter(ctx context.Context, subscriptionID string) error {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := newsletterspb.NewNewslettersServiceClient(conn)
	_, err = client.Unsubscribe(ctx, &newsletterspb.UnsubscribeRequest{
		SubscriptionId: subscriptionID,
		Reason:         "User requested unsubscribe",
	})
	if err != nil {
		return fmt.Errorf("Unsubscribe RPC failed: %w", err)
	}

	return nil
}

// GetSubscription retrieves a specific subscription
func (r NewsletterRepository) GetSubscription(ctx context.Context, subscriptionID string) (*models.GetSubscriptionResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := newsletterspb.NewNewslettersServiceClient(conn)
	resp, err := client.GetSubscription(ctx, &newsletterspb.GetSubscriptionRequest{
		Id: subscriptionID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetSubscription RPC failed: %w", err)
	}

	subscription := r.convertSubscriptionFromPb(resp.GetSubscription())
	return &models.GetSubscriptionResponse{
		Subscription: subscription,
	}, nil
}

// ListSubscriptions lists all subscriptions for a user
func (r NewsletterRepository) ListSubscriptions(ctx context.Context, userID, subscriptionStatus string, page, limit int32) (*models.ListSubscriptionsResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := newsletterspb.NewNewslettersServiceClient(conn)
	resp, err := client.ListSubscriptions(ctx, &newsletterspb.ListSubscriptionsRequest{
		UserId: userID,
		Status: subscriptionStatus,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("ListSubscriptions RPC failed: %w", err)
	}

	subscriptions := make([]*models.Subscription, len(resp.GetSubscriptions()))
	for i, pbSubscription := range resp.GetSubscriptions() {
		subscriptions[i] = r.convertSubscriptionFromPb(pbSubscription)
	}

	return &models.ListSubscriptionsResponse{
		Subscriptions: subscriptions,
		Total:         resp.GetTotal(),
		Page:          resp.GetPage(),
		Limit:         resp.GetLimit(),
	}, nil
}

// UpdateSubscription updates subscription preferences
func (r NewsletterRepository) UpdateSubscription(ctx context.Context, subscriptionID, subscriptionPreferences, subscriptionStatus string) (*models.UpdateSubscriptionResponse, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Build preferences
	prefs := &newsletterspb.SubscriptionPreferences{}
	if subscriptionPreferences != "" {
		prefs.FrequencyOverride = subscriptionPreferences
	}

	client := newsletterspb.NewNewslettersServiceClient(conn)
	resp, err := client.UpdateSubscription(ctx, &newsletterspb.UpdateSubscriptionRequest{
		Id:          subscriptionID,
		Status:      subscriptionStatus,
		Preferences: prefs,
	})
	if err != nil {
		return nil, fmt.Errorf("UpdateSubscription RPC failed: %w", err)
	}

	subscription := r.convertSubscriptionFromPb(resp.GetSubscription())
	return &models.UpdateSubscriptionResponse{
		Subscription: subscription,
	}, nil
}

// SendNewsletter sends a newsletter to all active subscribers
func (r NewsletterRepository) SendNewsletter(ctx context.Context, newsletterID, content string) (*models.SendNewsletterResponse, error) {
	// Note: The newsletter service uses SendEdition for sending newsletters
	// This is a wrapper that creates and sends an edition
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// First create an edition
	client := newsletterspb.NewNewslettersServiceClient(conn)
	editionResp, err := client.CreateEdition(ctx, &newsletterspb.CreateEditionRequest{
		NewsletterId: newsletterID,
		Subject:      "Newsletter Update",
		ContentHtml:  content,
		ContentText:  content, // Simple text version
	})
	if err != nil {
		return nil, fmt.Errorf("CreateEdition RPC failed: %w", err)
	}

	// Then send the edition
	sendResp, err := client.SendEdition(ctx, &newsletterspb.SendEditionRequest{
		Id: editionResp.GetEdition().GetId(),
	})
	if err != nil {
		return nil, fmt.Errorf("SendEdition RPC failed: %w", err)
	}

	return &models.SendNewsletterResponse{
		NewsletterID: newsletterID,
		Message:      fmt.Sprintf("Newsletter sent successfully to %d recipients", sendResp.GetRecipientsQueued()),
	}, nil
}

// GetSubscriptionByID retrieves a subscription by ID (wrapper for GetSubscription)
func (r NewsletterRepository) GetSubscriptionByID(ctx context.Context, subscriptionID string) (*models.Subscription, error) {
	resp, err := r.GetSubscription(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	return resp.Subscription, nil
}

// GetSubscriptionsByUser retrieves subscriptions by user ID
func (r NewsletterRepository) GetSubscriptionsByUser(ctx context.Context, userID string, limit int32) ([]*models.Subscription, error) {
	log.Printf("GetSubscriptionsByUser called for user: %s, limit: %d", userID, limit)

	resp, err := r.ListSubscriptions(ctx, userID, "", 1, limit)
	if err != nil {
		return nil, err
	}

	return resp.Subscriptions, nil
}

// GetSubscriptionsByNewsletter retrieves subscriptions by newsletter ID (mock implementation for AI tooling)
func (r NewsletterRepository) GetSubscriptionsByNewsletter(ctx context.Context, newsletterID string, limit int32) ([]*models.Subscription, error) {
	log.Printf("GetSubscriptionsByNewsletter called for newsletter: %s, limit: %d", newsletterID, limit)

	// Use ListSubscriptions with newsletter filter
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := newsletterspb.NewNewslettersServiceClient(conn)
	resp, err := client.ListSubscriptions(ctx, &newsletterspb.ListSubscriptionsRequest{
		NewsletterId: newsletterID,
		Page:         1,
		Limit:        limit,
	})
	if err != nil {
		return nil, fmt.Errorf("ListSubscriptions RPC failed: %w", err)
	}

	subscriptions := make([]*models.Subscription, len(resp.GetSubscriptions()))
	for i, pbSubscription := range resp.GetSubscriptions() {
		subscriptions[i] = r.convertSubscriptionFromPb(pbSubscription)
	}

	return subscriptions, nil
}

// GetActiveSubscriptions retrieves active subscriptions for a user
func (r NewsletterRepository) GetActiveSubscriptions(ctx context.Context, userID string, limit int32) ([]*models.Subscription, error) {
	log.Printf("GetActiveSubscriptions called for user: %s, limit: %d", userID, limit)

	resp, err := r.ListSubscriptions(ctx, userID, models.SubscriptionStatusActive, 1, limit)
	if err != nil {
		return nil, err
	}

	return resp.Subscriptions, nil
}

// GetInactiveSubscriptions retrieves inactive subscriptions for a user
func (r NewsletterRepository) GetInactiveSubscriptions(ctx context.Context, userID string, limit int32) ([]*models.Subscription, error) {
	log.Printf("GetInactiveSubscriptions called for user: %s, limit: %d", userID, limit)

	resp, err := r.ListSubscriptions(ctx, userID, models.SubscriptionStatusInactive, 1, limit)
	if err != nil {
		return nil, err
	}

	return resp.Subscriptions, nil
}

// SearchSubscriptions searches subscriptions by query (mock implementation for AI tooling)
func (r NewsletterRepository) SearchSubscriptions(ctx context.Context, userID, query string, limit int32) ([]*models.Subscription, error) {
	log.Printf("SearchSubscriptions called for user: %s, query: %s, limit: %d (mock implementation)", userID, query, limit)

	// Mock implementation - in real scenario, this would require additional RPC method
	subscriptions := make([]*models.Subscription, 0, limit)
	for i := int32(0); i < limit && i < 2; i++ { // Mock max 2 results
		subscriptions = append(subscriptions, &models.Subscription{
			SubscriptionID:          fmt.Sprintf("search_sub_%d", i+1),
			UserID:                  userID,
			NewsletterID:            fmt.Sprintf("newsletter_%d", i+1),
			SubscriptionPreferences: models.PreferenceDaily,
			SubscriptionStatus:      models.SubscriptionStatusActive,
		})
	}

	return subscriptions, nil
}

// CountSubscriptions counts subscriptions with filters
func (r NewsletterRepository) CountSubscriptions(ctx context.Context, userID, subscriptionStatus string) (int32, error) {
	log.Printf("CountSubscriptions called for user: %s, status: %s", userID, subscriptionStatus)

	resp, err := r.ListSubscriptions(ctx, userID, subscriptionStatus, 1, 1000) // Get all to count
	if err != nil {
		return 0, err
	}

	return resp.Total, nil
}

// CountActiveSubscriptions counts active subscriptions for a user
func (r NewsletterRepository) CountActiveSubscriptions(ctx context.Context, userID string) (int32, error) {
	log.Printf("CountActiveSubscriptions called for user: %s", userID)
	return r.CountSubscriptions(ctx, userID, models.SubscriptionStatusActive)
}

// CountSubscriptionsByNewsletter counts subscriptions for a newsletter (mock implementation for AI tooling)
func (r NewsletterRepository) CountSubscriptionsByNewsletter(ctx context.Context, newsletterID string) (int32, error) {
	log.Printf("CountSubscriptionsByNewsletter called for newsletter: %s", newsletterID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := newsletterspb.NewNewslettersServiceClient(conn)
	resp, err := client.ListSubscriptions(ctx, &newsletterspb.ListSubscriptionsRequest{
		NewsletterId: newsletterID,
		Page:         1,
		Limit:        1, // Just need the total
	})
	if err != nil {
		return 0, fmt.Errorf("ListSubscriptions RPC failed: %w", err)
	}

	return resp.GetTotal(), nil
}

// GetNewsletter retrieves a newsletter by ID
func (r NewsletterRepository) GetNewsletter(ctx context.Context, newsletterID string) (*models.Newsletter, error) {
	log.Printf("GetNewsletter called for ID: %s", newsletterID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := newsletterspb.NewNewslettersServiceClient(conn)
	resp, err := client.GetNewsletter(ctx, &newsletterspb.GetNewsletterRequest{
		Id: newsletterID,
	})
	if err != nil {
		return nil, fmt.Errorf("GetNewsletter RPC failed: %w", err)
	}

	newsletter := resp.GetNewsletter()
	return &models.Newsletter{
		ID:          newsletter.GetId(),
		Title:       newsletter.GetName(),
		Description: newsletter.GetDescription(),
		Content:     "",                               // Newsletter content is in editions, not the newsletter itself
		Status:      models.NewsletterStatusPublished, // Map from is_active
	}, nil
}

// GetNewsletters retrieves newsletters by status
func (r NewsletterRepository) GetNewsletters(ctx context.Context, status string, limit int32) ([]*models.Newsletter, error) {
	log.Printf("GetNewsletters called with status: %s, limit: %d", status, limit)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Map status to active_only parameter
	activeOnly := status == models.NewsletterStatusPublished

	client := newsletterspb.NewNewslettersServiceClient(conn)
	resp, err := client.ListNewsletters(ctx, &newsletterspb.ListNewslettersRequest{
		ActiveOnly: activeOnly,
		Page:       1,
		Limit:      limit,
	})
	if err != nil {
		return nil, fmt.Errorf("ListNewsletters RPC failed: %w", err)
	}

	newsletters := make([]*models.Newsletter, len(resp.GetNewsletters()))
	for i, pbNewsletter := range resp.GetNewsletters() {
		newsletters[i] = &models.Newsletter{
			ID:          pbNewsletter.GetId(),
			Title:       pbNewsletter.GetName(),
			Description: pbNewsletter.GetDescription(),
			Content:     "",                               // Newsletter content is in editions
			Status:      models.NewsletterStatusPublished, // Map from is_active
		}
	}

	return newsletters, nil
}

// SearchNewsletters searches newsletters by query (mock implementation for AI tooling)
func (r NewsletterRepository) SearchNewsletters(ctx context.Context, query string, limit int32) ([]*models.Newsletter, error) {
	log.Printf("SearchNewsletters called with query: %s, limit: %d (mock implementation)", query, limit)

	newsletters := make([]*models.Newsletter, 0, limit)
	for i := int32(0); i < limit && i < 2; i++ { // Mock max 2 results
		newsletters = append(newsletters, &models.Newsletter{
			ID:          fmt.Sprintf("search_newsletter_%d", i+1),
			Title:       fmt.Sprintf("Search result %d for: %s", i+1, query),
			Description: fmt.Sprintf("Description for search result %d", i+1),
			Content:     fmt.Sprintf("Content for search result %d", i+1),
			Status:      models.NewsletterStatusPublished,
		})
	}

	return newsletters, nil
}

// convertSubscriptionFromPb converts protobuf Subscription to domain Subscription
func (r NewsletterRepository) convertSubscriptionFromPb(pbSubscription *newsletterspb.Subscription) *models.Subscription {
	if pbSubscription == nil {
		return nil
	}

	// Map preferences
	subscriptionPrefs := models.PreferenceWeekly // Default
	if pbSubscription.GetPreferences() != nil && pbSubscription.GetPreferences().GetFrequencyOverride() != "" {
		switch pbSubscription.GetPreferences().GetFrequencyOverride() {
		case "daily":
			subscriptionPrefs = models.PreferenceDaily
		case "weekly":
			subscriptionPrefs = models.PreferenceWeekly
		case "monthly":
			subscriptionPrefs = models.PreferenceMonthly
		}
	}

	return &models.Subscription{
		SubscriptionID:          pbSubscription.GetId(),
		UserID:                  pbSubscription.GetUserId(),
		NewsletterID:            pbSubscription.GetNewsletterId(),
		SubscriptionPreferences: subscriptionPrefs,
		SubscriptionStatus:      pbSubscription.GetStatus(),
	}
}

// dial establishes a gRPC connection to the newsletters service
// dial sets up a gRPC connection with the microservice endpoint
func (r NewsletterRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r NewsletterRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	if r.auth != nil {
		return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
	}
	return rpc.Dial(ctx, r.endpoint)
}
