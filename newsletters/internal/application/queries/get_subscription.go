package queries

import (
	"context"
	"middleman/newsletters/internal/domain"
)

type GetSubscription struct {
	ID string
}

type GetSubscriptionHandler struct {
	catalog           domain.SubscriptionCatalogRepository
	newsletterCatalog domain.NewsletterCatalogRepository
}

func NewGetSubscriptionHandler(
	catalog domain.SubscriptionCatalogRepository,
	newsletterCatalog domain.NewsletterCatalogRepository,
) GetSubscriptionHandler {
	return GetSubscriptionHandler{
		catalog:           catalog,
		newsletterCatalog: newsletterCatalog,
	}
}

func (h GetSubscriptionHandler) GetSubscription(ctx context.Context, query GetSubscription) (*domain.CatalogSubscription, error) {
	subscription, err := h.catalog.Find(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	// Populate newsletter details
	newsletter, err := h.newsletterCatalog.Find(ctx, subscription.NewsletterID)
	if err == nil {
		subscription.Newsletter = newsletter
	}

	return subscription, nil
}