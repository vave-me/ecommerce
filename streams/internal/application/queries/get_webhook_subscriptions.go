package queries

import (
	"context"

	"middleman/streams/internal/domain"
)

// ListWebhookSubscriptions lists webhook subscriptions
type ListWebhookSubscriptions struct {
	Active bool
}

// ListWebhookSubscriptionsHandler handles listing webhook subscriptions
type ListWebhookSubscriptionsHandler struct {
	webhookRepo domain.WebhookSubscriptionRepository
}

// NewListWebhookSubscriptionsHandler creates a new handler
func NewListWebhookSubscriptionsHandler(webhookRepo domain.WebhookSubscriptionRepository) ListWebhookSubscriptionsHandler {
	return ListWebhookSubscriptionsHandler{
		webhookRepo: webhookRepo,
	}
}

// ListWebhookSubscriptions returns all webhook subscriptions
func (h ListWebhookSubscriptionsHandler) ListWebhookSubscriptions(ctx context.Context, query ListWebhookSubscriptions) ([]*domain.WebhookSubscription, error) {
	subscriptions, err := h.webhookRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Filter by active status if requested
	if query.Active {
		var activeSubscriptions []*domain.WebhookSubscription
		for _, sub := range subscriptions {
			if sub.Active {
				activeSubscriptions = append(activeSubscriptions, sub)
			}
		}
		return activeSubscriptions, nil
	}

	return subscriptions, nil
}

// GetWebhookSubscription gets a specific webhook subscription
type GetWebhookSubscription struct {
	ID string
}

// GetWebhookSubscriptionHandler handles getting a webhook subscription
type GetWebhookSubscriptionHandler struct {
	webhookRepo domain.WebhookSubscriptionRepository
}

// NewGetWebhookSubscriptionHandler creates a new handler
func NewGetWebhookSubscriptionHandler(webhookRepo domain.WebhookSubscriptionRepository) GetWebhookSubscriptionHandler {
	return GetWebhookSubscriptionHandler{
		webhookRepo: webhookRepo,
	}
}

// GetWebhookSubscription returns a webhook subscription by ID
func (h GetWebhookSubscriptionHandler) GetWebhookSubscription(ctx context.Context, query GetWebhookSubscription) (*domain.WebhookSubscription, error) {
	return h.webhookRepo.Find(query.ID)
}

// GetWebhookDeliveries gets webhook deliveries for a subscription
type GetWebhookDeliveries struct {
	SubscriptionID string
	Limit          int
}

// GetWebhookDeliveriesHandler handles getting webhook deliveries
type GetWebhookDeliveriesHandler struct {
	deliveryRepo domain.WebhookDeliveryRepository
}

// NewGetWebhookDeliveriesHandler creates a new handler
func NewGetWebhookDeliveriesHandler(deliveryRepo domain.WebhookDeliveryRepository) GetWebhookDeliveriesHandler {
	return GetWebhookDeliveriesHandler{
		deliveryRepo: deliveryRepo,
	}
}

// GetWebhookDeliveries returns webhook deliveries for a subscription
func (h GetWebhookDeliveriesHandler) GetWebhookDeliveries(ctx context.Context, query GetWebhookDeliveries) ([]*domain.WebhookDelivery, error) {
	return h.deliveryRepo.FindBySubscription(query.SubscriptionID)
}