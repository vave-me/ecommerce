package commands

import (
	"context"
	"fmt"
	"time"
	"middleman/internal/erp"
	"middleman/erp/internal/domain"

	"github.com/rs/zerolog/log"
)

// ProcessWebhook processes an incoming webhook event
type ProcessWebhook struct {
	ConnectorID string
	EventID     string
	EventType   string
	Payload     []byte
	Signature   string
	Headers     map[string]string
}

// ProcessWebhookHandler handles webhook processing
type ProcessWebhookHandler struct {
	registry   erp.ConnectorRegistry
	repository domain.ConnectorRepository
	webhookRepo domain.WebhookEventRepository
	eventPublisher domain.EventPublisher
}

// NewProcessWebhookHandler creates a new webhook processor
func NewProcessWebhookHandler(
	registry erp.ConnectorRegistry,
	repository domain.ConnectorRepository,
	webhookRepo domain.WebhookEventRepository,
	eventPublisher domain.EventPublisher,
) ProcessWebhookHandler {
	return ProcessWebhookHandler{
		registry:   registry,
		repository: repository,
		webhookRepo: webhookRepo,
		eventPublisher: eventPublisher,
	}
}

// ProcessWebhook processes the webhook
func (h ProcessWebhookHandler) ProcessWebhook(ctx context.Context, cmd ProcessWebhook) error {
	log.Info().
		Str("connector_id", cmd.ConnectorID).
		Str("event_id", cmd.EventID).
		Str("event_type", cmd.EventType).
		Msg("Processing webhook event")

	// Create webhook event record
	webhookEvent := &domain.WebhookEvent{
		ID:          cmd.EventID,
		ConnectorID: cmd.ConnectorID,
		EventID:     cmd.EventID,
		EventType:   cmd.EventType,
		Payload:     string(cmd.Payload),
		Signature:   cmd.Signature,
		Headers:     cmd.Headers,
		Status:      domain.WebhookStatusPending,
		ReceivedAt:  time.Now(),
	}

	// Save webhook event
	if err := h.webhookRepo.Create(ctx, webhookEvent); err != nil {
		return fmt.Errorf("saving webhook event: %w", err)
	}

	// Get connector instance from registry
	connector, err := h.registry.GetConnector(cmd.ConnectorID)
	if err != nil {
		webhookEvent.Status = domain.WebhookStatusFailed
		webhookEvent.Error = fmt.Sprintf("connector not found: %v", err)
		h.webhookRepo.Update(ctx, webhookEvent)
		return fmt.Errorf("connector not found: %w", err)
	}

	// Validate webhook signature
	if err := connector.ValidateWebhook(cmd.Payload, cmd.Signature); err != nil {
		webhookEvent.Status = domain.WebhookStatusFailed
		webhookEvent.Error = fmt.Sprintf("invalid signature: %v", err)
		h.webhookRepo.Update(ctx, webhookEvent)
		return fmt.Errorf("webhook validation failed: %w", err)
	}

	// Parse webhook into canonical event
	canonicalEvent, err := connector.ParseWebhook(cmd.Payload)
	if err != nil {
		webhookEvent.Status = domain.WebhookStatusFailed
		webhookEvent.Error = fmt.Sprintf("parse error: %v", err)
		h.webhookRepo.Update(ctx, webhookEvent)
		return fmt.Errorf("parsing webhook: %w", err)
	}

	// Process the webhook through the connector
	if err := connector.ProcessWebhook(ctx, cmd.Payload, cmd.Signature); err != nil {
		webhookEvent.Status = domain.WebhookStatusFailed
		webhookEvent.Error = fmt.Sprintf("processing error: %v", err)
		h.webhookRepo.Update(ctx, webhookEvent)
		return fmt.Errorf("processing webhook: %w", err)
	}

	// Publish domain event based on webhook type
	event := domain.WebhookProcessed{
		WebhookID:   webhookEvent.ID,
		ProcessedAt: time.Now(),
	}

	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		log.Error().Err(err).
			Str("connector_id", cmd.ConnectorID).
			Str("event_id", cmd.EventID).
			Msg("failed to publish webhook event")
	}

	// Update webhook status
	now := time.Now()
	webhookEvent.Status = domain.WebhookStatusProcessed
	webhookEvent.ProcessedAt = &now
	if err := h.webhookRepo.Update(ctx, webhookEvent); err != nil {
		log.Error().Err(err).
			Str("webhook_id", webhookEvent.ID).
			Msg("failed to update webhook status")
	}

	log.Info().
		Str("connector_id", cmd.ConnectorID).
		Str("event_id", cmd.EventID).
		Str("event_type", string(canonicalEvent.EventType)).
		Msg("Successfully processed webhook event")

	return nil
}