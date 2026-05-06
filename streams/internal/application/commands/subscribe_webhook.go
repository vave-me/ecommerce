package commands

import (
	"context"
	"time"

	"middleman/streams/internal/domain"
	"middleman/streams/internal/handlers"
	"middleman/streams/internal/infrastructure"
	
	"github.com/google/uuid"
)

// SubscribeWebhook creates a new webhook subscription
type SubscribeWebhook struct {
	URL         string                    `json:"url" validate:"required,url"`
	Events      []string                  `json:"events" validate:"required,min=1"`
	Secret      string                    `json:"secret"`
	Headers     map[string]string         `json:"headers"`
	RetryPolicy domain.RetryPolicy        `json:"retry_policy"`
}

// SubscribeWebhookHandler handles webhook subscription commands
type SubscribeWebhookHandler struct {
	webhookRepo domain.WebhookSubscriptionRepository
}

// NewSubscribeWebhookHandler creates a new webhook subscription handler
func NewSubscribeWebhookHandler(webhookRepo domain.WebhookSubscriptionRepository) SubscribeWebhookHandler {
	return SubscribeWebhookHandler{
		webhookRepo: webhookRepo,
	}
}

// SubscribeWebhook creates a new webhook subscription
func (h SubscribeWebhookHandler) SubscribeWebhook(ctx context.Context, cmd SubscribeWebhook) error {
	// Set default retry policy if not provided
	if cmd.RetryPolicy.MaxRetries == 0 {
		cmd.RetryPolicy = domain.RetryPolicy{
			MaxRetries:    3,
			BackoffFactor: 2.0,
			InitialDelay:  1 * time.Second,
			MaxBackoff:    5 * time.Minute,
		}
	}

	subscription := &domain.WebhookSubscription{
		ID:          uuid.New().String(),
		URL:         cmd.URL,
		Secret:      cmd.Secret,
		Events:      cmd.Events,
		Headers:     cmd.Headers,
		RetryPolicy: cmd.RetryPolicy,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return h.webhookRepo.Create(subscription)
}

// UpdateWebhook updates an existing webhook subscription
type UpdateWebhook struct {
	ID          string                    `json:"id"`
	URL         string                    `json:"url"`
	Events      []string                  `json:"events"`
	Headers     map[string]string         `json:"headers"`
	RetryPolicy *domain.RetryPolicy       `json:"retry_policy"`
	Active      *bool                     `json:"active"`
}

// UpdateWebhookHandler handles webhook update commands
type UpdateWebhookHandler struct {
	webhookRepo domain.WebhookSubscriptionRepository
}

// NewUpdateWebhookHandler creates a new webhook update handler
func NewUpdateWebhookHandler(webhookRepo domain.WebhookSubscriptionRepository) UpdateWebhookHandler {
	return UpdateWebhookHandler{
		webhookRepo: webhookRepo,
	}
}

// UpdateWebhook updates a webhook subscription
func (h UpdateWebhookHandler) UpdateWebhook(ctx context.Context, cmd UpdateWebhook) error {
	subscription, err := h.webhookRepo.Find(cmd.ID)
	if err != nil {
		return err
	}

	// Update fields if provided
	if cmd.URL != "" {
		subscription.URL = cmd.URL
	}
	if len(cmd.Events) > 0 {
		subscription.Events = cmd.Events
	}
	if cmd.Headers != nil {
		subscription.Headers = cmd.Headers
	}
	if cmd.RetryPolicy != nil {
		subscription.RetryPolicy = *cmd.RetryPolicy
	}
	if cmd.Active != nil {
		subscription.Active = *cmd.Active
	}
	
	subscription.UpdatedAt = time.Now()

	return h.webhookRepo.Update(subscription)
}

// UnsubscribeWebhook deletes a webhook subscription
type UnsubscribeWebhook struct {
	ID string `json:"id"`
}

// UnsubscribeWebhookHandler handles webhook unsubscribe commands
type UnsubscribeWebhookHandler struct {
	webhookRepo domain.WebhookSubscriptionRepository
}

// NewUnsubscribeWebhookHandler creates a new webhook unsubscribe handler
func NewUnsubscribeWebhookHandler(webhookRepo domain.WebhookSubscriptionRepository) UnsubscribeWebhookHandler {
	return UnsubscribeWebhookHandler{
		webhookRepo: webhookRepo,
	}
}

// UnsubscribeWebhook deletes a webhook subscription
func (h UnsubscribeWebhookHandler) UnsubscribeWebhook(ctx context.Context, cmd UnsubscribeWebhook) error {
	return h.webhookRepo.Delete(cmd.ID)
}

// TestWebhook sends a test webhook
type TestWebhook struct {
	SubscriptionID string `json:"subscription_id"`
}

// TestWebhookHandler handles test webhook commands
type TestWebhookHandler struct {
	webhookRepo   domain.WebhookSubscriptionRepository
	webhookClient *infrastructure.WebhookClient
}

// NewTestWebhookHandler creates a new test webhook handler
func NewTestWebhookHandler(
	webhookRepo domain.WebhookSubscriptionRepository,
	webhookClient *infrastructure.WebhookClient,
) TestWebhookHandler {
	return TestWebhookHandler{
		webhookRepo:   webhookRepo,
		webhookClient: webhookClient,
	}
}

// TestWebhook sends a test webhook
func (h TestWebhookHandler) TestWebhook(ctx context.Context, cmd TestWebhook) error {
	subscription, err := h.webhookRepo.Find(cmd.SubscriptionID)
	if err != nil {
		return err
	}

	return h.webhookClient.TestWebhook(ctx, subscription)
}

// RetryWebhookDelivery retries a failed webhook delivery
type RetryWebhookDelivery struct {
	DeliveryID string `json:"delivery_id"`
}

// RetryWebhookDeliveryHandler handles webhook retry commands
type RetryWebhookDeliveryHandler struct {
	deliveryRepo  domain.WebhookDeliveryRepository
	dispatcher    *handlers.WebhookDispatcher
}

// NewRetryWebhookDeliveryHandler creates a new retry handler
func NewRetryWebhookDeliveryHandler(
	deliveryRepo domain.WebhookDeliveryRepository,
	dispatcher *handlers.WebhookDispatcher,
) RetryWebhookDeliveryHandler {
	return RetryWebhookDeliveryHandler{
		deliveryRepo: deliveryRepo,
		dispatcher:   dispatcher,
	}
}

// RetryWebhookDelivery retries a webhook delivery
func (h RetryWebhookDeliveryHandler) RetryWebhookDelivery(ctx context.Context, cmd RetryWebhookDelivery) error {
	delivery, err := h.deliveryRepo.Find(cmd.DeliveryID)
	if err != nil {
		return err
	}

	// Reset retry time
	delivery.NextRetryAt = nil
	delivery.Status = domain.DeliveryStatusPending
	
	return h.deliveryRepo.Update(delivery)
}