package queries

import (
	"context"
	"middleman/newsletters/internal/domain"
)

type ListSubscriptions struct {
	UserID       string
	NewsletterID string
	Status       string
	Page         int
	Limit        int
}

type ListSubscriptionsHandler struct {
	catalog           domain.SubscriptionCatalogRepository
	newsletterCatalog domain.NewsletterCatalogRepository
}

func NewListSubscriptionsHandler(
	catalog domain.SubscriptionCatalogRepository,
	newsletterCatalog domain.NewsletterCatalogRepository,
) ListSubscriptionsHandler {
	return ListSubscriptionsHandler{
		catalog:           catalog,
		newsletterCatalog: newsletterCatalog,
	}
}

func (h ListSubscriptionsHandler) ListSubscriptions(ctx context.Context, query ListSubscriptions) ([]*domain.CatalogSubscription, int, error) {
	if query.Limit == 0 {
		query.Limit = 20
	}
	
	offset := query.Page * query.Limit
	
	var subscriptions []*domain.CatalogSubscription
	var total int
	var err error
	
	if query.UserID != "" {
		subscriptions, total, err = h.catalog.FindByUser(ctx, query.UserID, query.Status, query.Limit, offset)
	} else if query.NewsletterID != "" {
		subscriptions, total, err = h.catalog.FindByNewsletter(ctx, query.NewsletterID, query.Status, query.Limit, offset)
	} else {
		// Default to finding by user if available in context
		subscriptions, total, err = h.catalog.FindByUser(ctx, "", query.Status, query.Limit, offset)
	}
	
	if err != nil {
		return nil, 0, err
	}

	// Populate newsletter details
	newsletterMap := make(map[string]*domain.CatalogNewsletter)
	for _, sub := range subscriptions {
		if _, ok := newsletterMap[sub.NewsletterID]; !ok {
			newsletter, err := h.newsletterCatalog.Find(ctx, sub.NewsletterID)
			if err == nil {
				newsletterMap[sub.NewsletterID] = newsletter
			}
		}
		sub.Newsletter = newsletterMap[sub.NewsletterID]
	}

	return subscriptions, total, nil
}