package commands

import (
	"context"
	"github.com/google/uuid"
	"github.com/stackus/errors"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type SubscribeNewsletter struct {
	NewsletterID      string
	FrequencyOverride string   // Optional: override newsletter frequency
	Topics           []string  // Optional: filter by topics
	Format           string    // Optional: html or text
}

type SubscribeNewsletterHandler struct {
	subscriptions     domain.SubscriptionRepository
	catalog           domain.SubscriptionCatalogRepository
	newsletterCatalog domain.NewsletterCatalogRepository
	publisher         ddd.EventPublisher[ddd.Event]
}

func NewSubscribeNewsletterHandler(
	subscriptions domain.SubscriptionRepository,
	catalog domain.SubscriptionCatalogRepository,
	newsletterCatalog domain.NewsletterCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) SubscribeNewsletterHandler {
	return SubscribeNewsletterHandler{
		subscriptions:     subscriptions,
		catalog:           catalog,
		newsletterCatalog: newsletterCatalog,
		publisher:         publisher,
	}
}

func (h SubscribeNewsletterHandler) SubscribeNewsletter(ctx context.Context, cmd SubscribeNewsletter) (string, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return "", errors.ErrUnauthenticated.Msg("authentication required")
	}
	userID := claims.Subject
	
	// Verify newsletter exists and is active
	newsletter, err := h.newsletterCatalog.Find(ctx, cmd.NewsletterID)
	if err != nil {
		return "", err
	}
	if !newsletter.IsActive {
		return "", errors.ErrBadRequest.Msg("newsletter is not active")
	}
	
	// Check if already subscribed
	existing, err := h.catalog.FindByUserAndNewsletter(ctx, userID, cmd.NewsletterID)
	if err == nil && existing != nil && existing.Status == "active" {
		return "", errors.ErrConflict.Msg("already subscribed to this newsletter")
	}
	
	subscriptionID := uuid.New().String()
	subscription := domain.NewSubscription(subscriptionID)
	
	prefs := domain.SubscriptionPreferences{
		FrequencyOverride: domain.ToNewsletterFrequencyPtr(cmd.FrequencyOverride),
		Topics:            cmd.Topics,
		Format:            domain.ToContentFormat(cmd.Format),
	}
	event, err := subscription.Subscribe(userID, cmd.NewsletterID, prefs)
	if err != nil {
		return "", err
	}
	
	err = h.subscriptions.Save(ctx, subscription)
	if err != nil {
		return "", err
	}
	
	// Add to catalog
	catalogSub := &domain.CatalogSubscription{
		ID:                subscriptionID,
		UserID:            userID,
		NewsletterID:      cmd.NewsletterID,
		Status:            "active",
		FrequencyOverride: cmd.FrequencyOverride,
		Topics:            cmd.Topics,
		Format:            cmd.Format,
		SubscribedAt:      subscription.SubscribedAt,
	}
	
	err = h.catalog.Add(ctx, catalogSub)
	if err != nil {
		return "", err
	}
	
	err = h.publisher.Publish(ctx, event)
	if err != nil {
		return "", err
	}
	
	return subscriptionID, nil
}