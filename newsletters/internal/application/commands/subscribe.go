package commands

import (
	"context"
	"github.com/google/uuid"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type Subscribe struct {
	UserID       string
	NewsletterID string
	Preferences  domain.SubscriptionPreferences
}

type SubscribeHandler struct {
	subscriptions    domain.SubscriptionRepository
	catalog          domain.SubscriptionCatalogRepository
	newsletterCatalog domain.NewsletterCatalogRepository
	publisher        ddd.EventPublisher[ddd.Event]
}

func NewSubscribeHandler(
	subscriptions domain.SubscriptionRepository,
	catalog domain.SubscriptionCatalogRepository,
	newsletterCatalog domain.NewsletterCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) SubscribeHandler {
	return SubscribeHandler{
		subscriptions:     subscriptions,
		catalog:           catalog,
		newsletterCatalog: newsletterCatalog,
		publisher:         publisher,
	}
}

func (h SubscribeHandler) Subscribe(ctx context.Context, cmd Subscribe) (string, error) {
	// Check if newsletter exists and is active
	newsletter, err := h.newsletterCatalog.Find(ctx, cmd.NewsletterID)
	if err != nil {
		return "", errors.Wrap(err, "newsletter not found")
	}
	if !newsletter.IsActive {
		return "", errors.ErrBadRequest.Msg("newsletter is not active")
	}

	// Check if already subscribed
	existing, _ := h.catalog.FindByUserAndNewsletter(ctx, cmd.UserID, cmd.NewsletterID)
	if existing != nil && existing.Status == domain.ActiveStatus.String() {
		return "", errors.ErrBadRequest.Msg("already subscribed to this newsletter")
	}

	subscriptionID := uuid.New().String()
	subscription := domain.NewSubscription(subscriptionID)

	event, err := subscription.Subscribe(cmd.UserID, cmd.NewsletterID, cmd.Preferences)
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
		UserID:            cmd.UserID,
		NewsletterID:      cmd.NewsletterID,
		Status:            domain.ActiveStatus.String(),
		FrequencyOverride: "",
		Topics:            cmd.Preferences.Topics,
		Format:            cmd.Preferences.Format.String(),
		SubscribedAt:      subscription.SubscribedAt,
	}
	
	if cmd.Preferences.FrequencyOverride != nil {
		catalogSub.FrequencyOverride = cmd.Preferences.FrequencyOverride.String()
	}

	err = h.catalog.Add(ctx, catalogSub)
	if err != nil {
		return "", err
	}

	// Update subscriber count
	newsletter.SubscriberCount++
	err = h.newsletterCatalog.Update(ctx, newsletter)
	if err != nil {
		return "", err
	}

	err = h.publisher.Publish(ctx, event)
	if err != nil {
		return "", err
	}

	return subscriptionID, nil
}