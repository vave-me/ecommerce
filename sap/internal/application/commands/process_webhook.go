package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	
	"middleman/internal/ddd"
	"middleman/sap/internal/domain"
	"middleman/sap/internal/sap"
	"middleman/sap/internal/sap/transformer"
)

// ProcessWebhookHandler handles all webhook processing commands
type ProcessWebhookHandler struct {
	webhookRepo    domain.WebhookEventRepository
	syncLogRepo    domain.SyncLogRepository
	syncStatusRepo domain.SyncStatusRepository
	eventPublisher ddd.EventPublisher[ddd.Event]
	transformer    *transformer.SAPToCanonicalTransformer
}

// NewProcessWebhookHandler creates a new webhook processor
func NewProcessWebhookHandler(
	webhookRepo domain.WebhookEventRepository,
	syncLogRepo domain.SyncLogRepository,
	syncStatusRepo domain.SyncStatusRepository,
	eventPublisher ddd.EventPublisher[ddd.Event],
) *ProcessWebhookHandler {
	return &ProcessWebhookHandler{
		webhookRepo:    webhookRepo,
		syncLogRepo:    syncLogRepo,
		syncStatusRepo: syncStatusRepo,
		eventPublisher: eventPublisher,
		transformer:    transformer.NewSAPToCanonicalTransformer("SAP"),
	}
}

// StoreWebhookEvent stores a new webhook event
func (h *ProcessWebhookHandler) StoreWebhookEvent(ctx context.Context, event *domain.WebhookEvent) error {
	return h.webhookRepo.Create(ctx, event)
}

// UpdateWebhookEventStatus updates the status of a webhook event
func (h *ProcessWebhookHandler) UpdateWebhookEventStatus(ctx context.Context, id, status, errorMessage string) error {
	return h.webhookRepo.UpdateStatus(ctx, id, status, errorMessage)
}

// ProcessMaterialMaster processes MATMAS IDoc
func (h *ProcessWebhookHandler) ProcessMaterialMaster(ctx context.Context, idocData interface{}, correlationID, webhookEventID string) error {
	startTime := time.Now()
	
	// Convert to MATMAS structure
	matmas, ok := idocData.(*sap.MATMAS)
	if !ok {
		// Try parsing from raw IDoc
		idoc, ok := idocData.(*sap.IDoc)
		if !ok {
			return fmt.Errorf("invalid IDoc data type")
		}
		// Parse MATMAS from IDoc data
		// Implementation would depend on exact IDoc structure
		return fmt.Errorf("MATMAS parsing not yet implemented")
	}
	
	// Transform to canonical event
	event, err := h.transformer.TransformMATMAS(matmas, correlationID)
	if err != nil {
		return fmt.Errorf("transforming MATMAS: %w", err)
	}
	
	// Create sync log
	syncLog := &domain.SyncLog{
		ID:          uuid.New().String(),
		EventID:     event.EventID.String(),
		EventType:   string(event.EventType),
		Source:      "SAP-MATMAS",
		Destination: "EventBus",
		Status:      "processing",
		ProcessedAt: time.Now(),
	}
	
	// Marshal payload for logging
	payloadJSON, _ := json.Marshal(event.Payload)
	syncLog.Payload = payloadJSON
	
	// Publish canonical event
	if err := h.publishEvent(ctx, event); err != nil {
		syncLog.Status = "failed"
		syncLog.ErrorMessage = err.Error()
		syncLog.Duration = time.Since(startTime)
		_ = h.syncLogRepo.Create(ctx, syncLog)
		return fmt.Errorf("publishing event: %w", err)
	}
	
	// Update sync status
	productPayload := event.Payload.(transformer.ProductMasterPayload)
	syncStatus := &domain.SyncStatus{
		ID:           uuid.New().String(),
		EntityType:   "product",
		EntityID:     productPayload.SKU,
		LastSyncedAt: time.Now(),
		Status:       "success",
		UpdatedAt:    time.Now(),
	}
	
	if err := h.syncStatusRepo.Save(ctx, syncStatus); err != nil {
		log.Error().Err(err).Msg("Failed to save sync status")
	}
	
	// Complete sync log
	syncLog.Status = "processed"
	syncLog.Duration = time.Since(startTime)
	
	return h.syncLogRepo.Create(ctx, syncLog)
}

// ProcessInventoryUpdate processes INVCON IDoc
func (h *ProcessWebhookHandler) ProcessInventoryUpdate(ctx context.Context, idocData interface{}, correlationID, webhookEventID string) error {
	startTime := time.Now()
	
	// Convert to INVCON structure
	invcon, ok := idocData.(*sap.INVCON)
	if !ok {
		return fmt.Errorf("invalid INVCON data type")
	}
	
	// Transform to canonical event
	event, err := h.transformer.TransformINVCON(invcon, correlationID)
	if err != nil {
		return fmt.Errorf("transforming INVCON: %w", err)
	}
	
	// Create sync log
	syncLog := &domain.SyncLog{
		ID:          uuid.New().String(),
		EventID:     event.EventID.String(),
		EventType:   string(event.EventType),
		Source:      "SAP-INVCON",
		Destination: "EventBus",
		Status:      "processing",
		ProcessedAt: time.Now(),
	}
	
	// Marshal payload for logging
	payloadJSON, _ := json.Marshal(event.Payload)
	syncLog.Payload = payloadJSON
	
	// Publish canonical event
	if err := h.publishEvent(ctx, event); err != nil {
		syncLog.Status = "failed"
		syncLog.ErrorMessage = err.Error()
		syncLog.Duration = time.Since(startTime)
		_ = h.syncLogRepo.Create(ctx, syncLog)
		return fmt.Errorf("publishing event: %w", err)
	}
	
	// Update sync status
	stockPayload := event.Payload.(transformer.StockLevelPayload)
	syncStatus := &domain.SyncStatus{
		ID:           uuid.New().String(),
		EntityType:   "stock",
		EntityID:     fmt.Sprintf("%s-%s", stockPayload.SKU, stockPayload.WarehouseID),
		LastSyncedAt: time.Now(),
		Status:       "success",
		UpdatedAt:    time.Now(),
	}
	
	if err := h.syncStatusRepo.Save(ctx, syncStatus); err != nil {
		log.Error().Err(err).Msg("Failed to save sync status")
	}
	
	// Complete sync log
	syncLog.Status = "processed"
	syncLog.Duration = time.Since(startTime)
	
	return h.syncLogRepo.Create(ctx, syncLog)
}

// ProcessPricingUpdate processes COND_A IDoc
func (h *ProcessWebhookHandler) ProcessPricingUpdate(ctx context.Context, idocData interface{}, correlationID, webhookEventID string) error {
	startTime := time.Now()
	
	// Convert to COND_A structure
	conda, ok := idocData.(*sap.COND_A)
	if !ok {
		return fmt.Errorf("invalid COND_A data type")
	}
	
	// Transform to canonical event
	event, err := h.transformer.TransformCOND_A(conda, correlationID)
	if err != nil {
		return fmt.Errorf("transforming COND_A: %w", err)
	}
	
	// Create sync log
	syncLog := &domain.SyncLog{
		ID:          uuid.New().String(),
		EventID:     event.EventID.String(),
		EventType:   string(event.EventType),
		Source:      "SAP-COND_A",
		Destination: "EventBus",
		Status:      "processing",
		ProcessedAt: time.Now(),
	}
	
	// Marshal payload for logging
	payloadJSON, _ := json.Marshal(event.Payload)
	syncLog.Payload = payloadJSON
	
	// Publish canonical event
	if err := h.publishEvent(ctx, event); err != nil {
		syncLog.Status = "failed"
		syncLog.ErrorMessage = err.Error()
		syncLog.Duration = time.Since(startTime)
		_ = h.syncLogRepo.Create(ctx, syncLog)
		return fmt.Errorf("publishing event: %w", err)
	}
	
	// Update sync status
	pricePayload := event.Payload.(transformer.PricePayload)
	syncStatus := &domain.SyncStatus{
		ID:           uuid.New().String(),
		EntityType:   "price",
		EntityID:     fmt.Sprintf("%s-%s", pricePayload.SKU, pricePayload.PriceListID),
		LastSyncedAt: time.Now(),
		Status:       "success",
		UpdatedAt:    time.Now(),
	}
	
	if err := h.syncStatusRepo.Save(ctx, syncStatus); err != nil {
		log.Error().Err(err).Msg("Failed to save sync status")
	}
	
	// Complete sync log
	syncLog.Status = "processed"
	syncLog.Duration = time.Since(startTime)
	
	return h.syncLogRepo.Create(ctx, syncLog)
}

// ProcessProductEvent processes JSON product events
func (h *ProcessWebhookHandler) ProcessProductEvent(ctx context.Context, event *sap.SAPEvent, webhookEventID string) error {
	// Transform to canonical event
	canonicalEvent, err := h.transformer.TransformSAPEvent(event)
	if err != nil {
		return fmt.Errorf("transforming product event: %w", err)
	}
	
	// Publish and log
	return h.processAndPublishEvent(ctx, canonicalEvent, "SAP-Product")
}

// ProcessStockEvent processes JSON stock events
func (h *ProcessWebhookHandler) ProcessStockEvent(ctx context.Context, event *sap.SAPEvent, webhookEventID string) error {
	// Transform to canonical event
	canonicalEvent, err := h.transformer.TransformSAPEvent(event)
	if err != nil {
		return fmt.Errorf("transforming stock event: %w", err)
	}
	
	// Publish and log
	return h.processAndPublishEvent(ctx, canonicalEvent, "SAP-Stock")
}

// ProcessPriceEvent processes JSON price events
func (h *ProcessWebhookHandler) ProcessPriceEvent(ctx context.Context, event *sap.SAPEvent, webhookEventID string) error {
	// Transform to canonical event
	canonicalEvent, err := h.transformer.TransformSAPEvent(event)
	if err != nil {
		return fmt.Errorf("transforming price event: %w", err)
	}
	
	// Publish and log
	return h.processAndPublishEvent(ctx, canonicalEvent, "SAP-Price")
}

// ProcessProductDeletedEvent processes product deletion events
func (h *ProcessWebhookHandler) ProcessProductDeletedEvent(ctx context.Context, event *sap.SAPEvent, webhookEventID string) error {
	// Transform to canonical event
	canonicalEvent, err := h.transformer.TransformSAPEvent(event)
	if err != nil {
		return fmt.Errorf("transforming delete event: %w", err)
	}
	
	// Publish and log
	return h.processAndPublishEvent(ctx, canonicalEvent, "SAP-Delete")
}

// Helper method to process and publish events
func (h *ProcessWebhookHandler) processAndPublishEvent(ctx context.Context, event *transformer.CanonicalEvent, source string) error {
	startTime := time.Now()
	
	// Create sync log
	syncLog := &domain.SyncLog{
		ID:          uuid.New().String(),
		EventID:     event.EventID.String(),
		EventType:   string(event.EventType),
		Source:      source,
		Destination: "EventBus",
		Status:      "processing",
		ProcessedAt: time.Now(),
	}
	
	// Marshal payload for logging
	payloadJSON, _ := json.Marshal(event.Payload)
	syncLog.Payload = payloadJSON
	
	// Publish canonical event
	if err := h.publishEvent(ctx, event); err != nil {
		syncLog.Status = "failed"
		syncLog.ErrorMessage = err.Error()
		syncLog.Duration = time.Since(startTime)
		_ = h.syncLogRepo.Create(ctx, syncLog)
		return fmt.Errorf("publishing event: %w", err)
	}
	
	// Complete sync log
	syncLog.Status = "processed"
	syncLog.Duration = time.Since(startTime)
	
	return h.syncLogRepo.Create(ctx, syncLog)
}

// publishEvent publishes a canonical event
func (h *ProcessWebhookHandler) publishEvent(ctx context.Context, event *transformer.CanonicalEvent) error {
	// Convert to DDD event
	domainEvent := ddd.NewEvent(
		string(event.EventType),
		event.Payload,
	)
	
	log.Info().
		Str("eventId", event.EventID.String()).
		Str("eventType", string(event.EventType)).
		Str("source", event.Source).
		Msg("Publishing canonical event")
	
	return h.eventPublisher.Publish(ctx, domainEvent)
}