package commands

import (
	"context"
	"time"
	"middleman/internal/ddd"
	"middleman/newsletters/internal/domain"
)

type Unsubscribe struct {
	SubscriptionID string
	Reason         string
}

type UnsubscribeHandler struct {
	subscriptions     domain.SubscriptionRepository
	catalog           domain.SubscriptionCatalogRepository
	newsletterCatalog domain.NewsletterCatalogRepository
	publisher         ddd.EventPublisher[ddd.Event]
}

func NewUnsubscribeHandler(
	subscriptions domain.SubscriptionRepository,
	catalog domain.SubscriptionCatalogRepository,
	newsletterCatalog domain.NewsletterCatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UnsubscribeHandler {
	return UnsubscribeHandler{
		subscriptions:     subscriptions,
		catalog:           catalog,
		newsletterCatalog: newsletterCatalog,
		publisher:         publisher,
	}
}

func (h UnsubscribeHandler) Unsubscribe(ctx context.Context, cmd Unsubscribe) error {
	subscription, err := h.subscriptions.Load(ctx, cmd.SubscriptionID)
	if err != nil {
		return err
	}

	event, err := subscription.Unsubscribe(cmd.Reason)
	if err != nil {
		return err
	}

	err = h.subscriptions.Save(ctx, subscription)
	if err != nil {
		return err
	}

	// Update catalog
	catalogSub, err := h.catalog.Find(ctx, cmd.SubscriptionID)
	if err != nil {
		return err
	}

	catalogSub.Status = domain.UnsubscribedStatus.String()
	now := time.Now()
	catalogSub.UnsubscribedAt = &now

	err = h.catalog.Update(ctx, catalogSub)
	if err != nil {
		return err
	}

	// Update subscriber count
	newsletter, err := h.newsletterCatalog.Find(ctx, subscription.NewsletterID)
	if err != nil {
		return err
	}

	if newsletter.SubscriberCount > 0 {
		newsletter.SubscriberCount--
		err = h.newsletterCatalog.Update(ctx, newsletter)
		if err != nil {
			return err
		}
	}

	return h.publisher.Publish(ctx, event)
}