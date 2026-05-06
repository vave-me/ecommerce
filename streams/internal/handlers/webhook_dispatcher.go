package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"middleman/internal/ddd"
	"middleman/streams/internal/domain"
	"middleman/streams/internal/infrastructure"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// WebhookDispatcher handles dispatching domain events to webhook subscribers
type WebhookDispatcher struct {
	subscriptionRepo domain.WebhookSubscriptionRepository
	deliveryRepo     domain.WebhookDeliveryRepository
	webhookClient    *infrastructure.WebhookClient
	logger           *zap.Logger
	maxConcurrent    int
	deliveryChan     chan *webhookDeliveryTask
	workerCtx        context.Context
	workerCancel     context.CancelFunc
}

type webhookDeliveryTask struct {
	subscription *domain.WebhookSubscription
	event        *domain.WebhookEvent
	delivery     *domain.WebhookDelivery
}

// NewWebhookDispatcher creates a new webhook dispatcher
func NewWebhookDispatcher(
	subscriptionRepo domain.WebhookSubscriptionRepository,
	deliveryRepo domain.WebhookDeliveryRepository,
	webhookClient *infrastructure.WebhookClient,
	logger *zap.Logger,
	maxConcurrent int,
) *WebhookDispatcher {
	ctx, cancel := context.WithCancel(context.Background())
	
	dispatcher := &WebhookDispatcher{
		subscriptionRepo: subscriptionRepo,
		deliveryRepo:     deliveryRepo,
		webhookClient:    webhookClient,
		logger:           logger,
		maxConcurrent:    maxConcurrent,
		deliveryChan:     make(chan *webhookDeliveryTask, 1000),
		workerCtx:        ctx,
		workerCancel:     cancel,
	}

	// Start worker pool
	dispatcher.startWorkers()

	return dispatcher
}

// startWorkers starts the webhook delivery worker pool
func (d *WebhookDispatcher) startWorkers() {
	for i := 0; i < d.maxConcurrent; i++ {
		go d.deliveryWorker()
	}

	// Start retry worker
	go d.retryWorker()
}

// deliveryWorker processes webhook deliveries
func (d *WebhookDispatcher) deliveryWorker() {
	for {
		select {
		case <-d.workerCtx.Done():
			return
		case task := <-d.deliveryChan:
			d.processDelivery(task)
		}
	}
}

// retryWorker checks for failed deliveries and retries them
func (d *WebhookDispatcher) retryWorker() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-d.workerCtx.Done():
			return
		case <-ticker.C:
			d.processRetries()
		}
	}
}

// HandleEvent handles a domain event and dispatches webhooks
func (d *WebhookDispatcher) HandleEvent(ctx context.Context, event ddd.Event) error {
	// Map domain event to webhook event type
	webhookEventType := domain.MapDomainEventToWebhookEvent(event.EventName())
	if webhookEventType == "" {
		// Event not configured for webhooks
		return nil
	}

	// Get active subscriptions for this event type
	subscriptions, err := d.subscriptionRepo.FindActiveByEvent(string(webhookEventType))
	if err != nil {
		d.logger.Error("Failed to get webhook subscriptions",
			zap.String("event_type", string(webhookEventType)),
			zap.Error(err))
		return err
	}

	if len(subscriptions) == 0 {
		return nil
	}

	// Extract stream ID from event
	streamID := d.extractStreamID(event)

	// Create webhook event
	webhookEvent := &domain.WebhookEvent{
		ID:        uuid.New().String(),
		Type:      string(webhookEventType),
		StreamID:  streamID,
		Timestamp: time.Now(),
		Data:      d.extractEventData(event),
	}

	// Create delivery records and queue for processing
	for _, subscription := range subscriptions {
		delivery := &domain.WebhookDelivery{
			ID:             uuid.New().String(),
			SubscriptionID: subscription.ID,
			EventID:        webhookEvent.ID,
			EventType:      webhookEvent.Type,
			Status:         domain.DeliveryStatusPending,
			Attempts:       0,
			CreatedAt:      time.Now(),
		}

		// Marshal payload for storage
		payload, err := json.Marshal(webhookEvent)
		if err != nil {
			d.logger.Error("Failed to marshal webhook event",
				zap.String("event_id", webhookEvent.ID),
				zap.Error(err))
			continue
		}
		delivery.Payload = payload

		// Save delivery record
		if err := d.deliveryRepo.Create(delivery); err != nil {
			d.logger.Error("Failed to create webhook delivery record",
				zap.String("subscription_id", subscription.ID),
				zap.Error(err))
			continue
		}

		// Queue for delivery
		select {
		case d.deliveryChan <- &webhookDeliveryTask{
			subscription: subscription,
			event:        webhookEvent,
			delivery:     delivery,
		}:
		default:
			d.logger.Warn("Webhook delivery queue full, delivery will be retried later",
				zap.String("delivery_id", delivery.ID))
		}
	}

	return nil
}

// processDelivery processes a single webhook delivery
func (d *WebhookDispatcher) processDelivery(task *webhookDeliveryTask) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Update delivery status
	task.delivery.Status = domain.DeliveryStatusRetrying
	task.delivery.Attempts++
	task.delivery.LastAttemptAt = timePtr(time.Now())

	// Attempt delivery
	result, err := d.webhookClient.DeliverWebhook(ctx, task.subscription, task.event)

	if err == nil && result.StatusCode >= 200 && result.StatusCode < 300 {
		// Success
		task.delivery.Status = domain.DeliveryStatusDelivered
		task.delivery.CompletedAt = timePtr(time.Now())
		task.delivery.ResponseStatus = result.StatusCode
		task.delivery.ResponseBody = result.ResponseBody

		d.logger.Info("Webhook delivered successfully",
			zap.String("delivery_id", task.delivery.ID),
			zap.String("subscription_id", task.subscription.ID),
			zap.Int("attempts", task.delivery.Attempts))
	} else {
		// Failure
		if err != nil {
			task.delivery.Error = err.Error()
		}
		if result != nil {
			task.delivery.ResponseStatus = result.StatusCode
			task.delivery.ResponseBody = result.ResponseBody
		}

		// Check if we should retry
		if task.delivery.Attempts < task.subscription.RetryPolicy.MaxRetries {
			// Calculate next retry time
			delay := d.calculateRetryDelay(task.delivery.Attempts, task.subscription.RetryPolicy)
			task.delivery.NextRetryAt = timePtr(time.Now().Add(delay))
			task.delivery.Status = domain.DeliveryStatusRetrying

			d.logger.Warn("Webhook delivery failed, will retry",
				zap.String("delivery_id", task.delivery.ID),
				zap.String("subscription_id", task.subscription.ID),
				zap.Int("attempts", task.delivery.Attempts),
				zap.Duration("retry_in", delay),
				zap.Error(err))
		} else {
			// Max retries exceeded
			task.delivery.Status = domain.DeliveryStatusFailed
			task.delivery.CompletedAt = timePtr(time.Now())

			d.logger.Error("Webhook delivery failed after max retries",
				zap.String("delivery_id", task.delivery.ID),
				zap.String("subscription_id", task.subscription.ID),
				zap.Int("attempts", task.delivery.Attempts),
				zap.Error(err))
		}
	}

	// Update delivery record
	if err := d.deliveryRepo.Update(task.delivery); err != nil {
		d.logger.Error("Failed to update webhook delivery record",
			zap.String("delivery_id", task.delivery.ID),
			zap.Error(err))
	}
}

// processRetries checks for deliveries that need to be retried
func (d *WebhookDispatcher) processRetries() {
	// Get pending deliveries that are ready for retry
	deliveries, err := d.deliveryRepo.FindPendingDeliveries(100)
	if err != nil {
		d.logger.Error("Failed to find pending webhook deliveries", zap.Error(err))
		return
	}

	for _, delivery := range deliveries {
		// Skip if not ready for retry
		if delivery.NextRetryAt != nil && delivery.NextRetryAt.After(time.Now()) {
			continue
		}

		// Get subscription
		subscription, err := d.subscriptionRepo.Find(delivery.SubscriptionID)
		if err != nil {
			d.logger.Error("Failed to find webhook subscription",
				zap.String("subscription_id", delivery.SubscriptionID),
				zap.Error(err))
			continue
		}

		// Skip inactive subscriptions
		if !subscription.Active {
			delivery.Status = domain.DeliveryStatusFailed
			delivery.Error = "subscription inactive"
			d.deliveryRepo.Update(delivery)
			continue
		}

		// Unmarshal event
		var event domain.WebhookEvent
		if err := json.Unmarshal(delivery.Payload, &event); err != nil {
			d.logger.Error("Failed to unmarshal webhook event",
				zap.String("delivery_id", delivery.ID),
				zap.Error(err))
			continue
		}

		// Queue for retry
		select {
		case d.deliveryChan <- &webhookDeliveryTask{
			subscription: subscription,
			event:        &event,
			delivery:     delivery,
		}:
		default:
			// Queue full, will try again later
		}
	}
}

// calculateRetryDelay calculates the delay before next retry
func (d *WebhookDispatcher) calculateRetryDelay(attempt int, policy domain.RetryPolicy) time.Duration {
	// Exponential backoff with jitter
	baseDelay := policy.InitialDelay
	if baseDelay == 0 {
		baseDelay = 1 * time.Second
	}

	delay := time.Duration(float64(baseDelay) * policy.BackoffFactor * float64(attempt))
	
	if policy.MaxBackoff > 0 && delay > policy.MaxBackoff {
		delay = policy.MaxBackoff
	}

	// Add jitter (±10%)
	jitter := time.Duration(float64(delay) * 0.1)
	delay = delay + time.Duration(time.Now().UnixNano()%int64(jitter))

	return delay
}

// extractStreamID extracts stream ID from domain event
func (d *WebhookDispatcher) extractStreamID(event ddd.Event) string {
	// This would need to be implemented based on your event structure
	// For now, return aggregate ID
	return event.AggregateID()
}

// extractEventData extracts relevant data from domain event
func (d *WebhookDispatcher) extractEventData(event ddd.Event) map[string]interface{} {
	// Convert event payload to map
	data := make(map[string]interface{})
	
	// Marshal and unmarshal to convert to map
	if payload, err := json.Marshal(event.Payload()); err == nil {
		json.Unmarshal(payload, &data)
	}

	// Add metadata
	data["event_id"] = event.EventID()
	data["aggregate_id"] = event.AggregateID()
	data["occurred_at"] = event.OccurredAt()

	return data
}

// Stop gracefully stops the webhook dispatcher
func (d *WebhookDispatcher) Stop() {
	d.logger.Info("Stopping webhook dispatcher")
	
	// Cancel context to signal workers to stop
	d.workerCancel()
	
	// Give workers time to finish current deliveries
	time.Sleep(2 * time.Second)
	
	// Close delivery channel
	close(d.deliveryChan)
	
	d.logger.Info("Webhook dispatcher stopped")
}

// timePtr returns a pointer to a time.Time
func timePtr(t time.Time) *time.Time {
	return &t
}