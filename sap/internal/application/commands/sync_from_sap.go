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

// SyncFromSAPHandler handles synchronization from SAP
type SyncFromSAPHandler struct {
	sapClient      *sap.EnhancedSAPClient
	syncStatusRepo domain.SyncStatusRepository
	syncLogRepo    domain.SyncLogRepository
	eventPublisher ddd.EventPublisher[ddd.Event]
	transformer    *transformer.SAPToCanonicalTransformer
}

// NewSyncFromSAPHandler creates a new sync handler
func NewSyncFromSAPHandler(
	sapClient *sap.EnhancedSAPClient,
	syncStatusRepo domain.SyncStatusRepository,
	syncLogRepo domain.SyncLogRepository,
	eventPublisher ddd.EventPublisher[ddd.Event],
) *SyncFromSAPHandler {
	return &SyncFromSAPHandler{
		sapClient:      sapClient,
		syncStatusRepo: syncStatusRepo,
		syncLogRepo:    syncLogRepo,
		eventPublisher: eventPublisher,
		transformer:    transformer.NewSAPToCanonicalTransformer("SAP-SYNC"),
	}
}

// SyncProductsFromSAP synchronizes products from SAP
func (h *SyncFromSAPHandler) SyncProductsFromSAP(ctx context.Context, since *time.Time) error {
	log.Info().
		Time("since", *since).
		Msg("Starting product sync from SAP")
	
	startTime := time.Now()
	syncLog := &domain.SyncLog{
		ID:          uuid.New().String(),
		EventType:   "ProductSync",
		Source:      "SAP",
		Destination: "EventBus",
		Status:      "processing",
		ProcessedAt: startTime,
	}
	
	// Get last sync time if not provided
	if since == nil {
		lastSync, err := h.syncStatusRepo.GetLastSyncTime(ctx, "product")
		if err != nil {
			// Default to 24 hours ago
			defaultTime := time.Now().Add(-24 * time.Hour)
			since = &defaultTime
		} else {
			since = &lastSync
		}
	}
	
	// Get product changes using enhanced client
	changes, err := h.sapClient.GetProductChangesEnhanced(ctx, *since)
	if err != nil {
		syncLog.Status = "failed"
		syncLog.ErrorMessage = err.Error()
		syncLog.Duration = time.Since(startTime)
		_ = h.syncLogRepo.Create(ctx, syncLog)
		return fmt.Errorf("getting product changes: %w", err)
	}
	
	log.Info().
		Int("productCount", len(changes)).
		Msg("Retrieved product changes from SAP")
	
	successCount := 0
	failedCount := 0
	
	// Process each change
	for _, change := range changes {
		if err := h.processProductChange(ctx, change); err != nil {
			log.Error().
				Err(err).
				Str("productId", change.ProductID).
				Msg("Failed to process product change")
			failedCount++
		} else {
			successCount++
		}
	}
	
	// Complete sync log
	syncLog.Status = "processed"
	syncLog.Duration = time.Since(startTime)
	syncLog.Payload, _ = json.Marshal(map[string]interface{}{
		"totalProducts": len(changes),
		"successful":    successCount,
		"failed":        failedCount,
	})
	
	if err := h.syncLogRepo.Create(ctx, syncLog); err != nil {
		log.Error().Err(err).Msg("Failed to create sync log")
	}
	
	return nil
}

// SyncStockLevelsFromSAP synchronizes stock levels from SAP
func (h *SyncFromSAPHandler) SyncStockLevelsFromSAP(ctx context.Context, productIDs []string) error {
	log.Info().
		Int("productCount", len(productIDs)).
		Msg("Starting stock sync from SAP")
	
	startTime := time.Now()
	syncLog := &domain.SyncLog{
		ID:          uuid.New().String(),
		EventType:   "StockSync",
		Source:      "SAP",
		Destination: "EventBus",
		Status:      "processing",
		ProcessedAt: startTime,
	}
	
	// Get stock levels using enhanced client
	stockLevels, err := h.sapClient.GetStockLevelsEnhanced(ctx, productIDs)
	if err != nil {
		syncLog.Status = "failed"
		syncLog.ErrorMessage = err.Error()
		syncLog.Duration = time.Since(startTime)
		_ = h.syncLogRepo.Create(ctx, syncLog)
		return fmt.Errorf("getting stock levels: %w", err)
	}
	
	log.Info().
		Int("stockCount", len(stockLevels)).
		Msg("Retrieved stock levels from SAP")
	
	successCount := 0
	failedCount := 0
	
	// Process each stock level
	for _, stock := range stockLevels {
		if err := h.processStockLevel(ctx, stock); err != nil {
			log.Error().
				Err(err).
				Str("productId", stock.ProductID).
				Str("warehouseId", stock.WarehouseID).
				Msg("Failed to process stock level")
			failedCount++
		} else {
			successCount++
		}
	}
	
	// Complete sync log
	syncLog.Status = "processed"
	syncLog.Duration = time.Since(startTime)
	syncLog.Payload, _ = json.Marshal(map[string]interface{}{
		"totalStock":  len(stockLevels),
		"successful":  successCount,
		"failed":      failedCount,
		"productIDs":  productIDs,
	})
	
	if err := h.syncLogRepo.Create(ctx, syncLog); err != nil {
		log.Error().Err(err).Msg("Failed to create sync log")
	}
	
	return nil
}

// SyncPricesFromSAP synchronizes prices from SAP
func (h *SyncFromSAPHandler) SyncPricesFromSAP(ctx context.Context, productIDs []string, priceListID string) error {
	log.Info().
		Int("productCount", len(productIDs)).
		Str("priceListID", priceListID).
		Msg("Starting price sync from SAP")
	
	startTime := time.Now()
	syncLog := &domain.SyncLog{
		ID:          uuid.New().String(),
		EventType:   "PriceSync",
		Source:      "SAP",
		Destination: "EventBus",
		Status:      "processing",
		ProcessedAt: startTime,
	}
	
	// Get prices from SAP
	prices, err := h.sapClient.GetPrices(ctx, productIDs, priceListID)
	if err != nil {
		syncLog.Status = "failed"
		syncLog.ErrorMessage = err.Error()
		syncLog.Duration = time.Since(startTime)
		_ = h.syncLogRepo.Create(ctx, syncLog)
		return fmt.Errorf("getting prices: %w", err)
	}
	
	log.Info().
		Int("priceCount", len(prices)).
		Msg("Retrieved prices from SAP")
	
	successCount := 0
	failedCount := 0
	
	// Process each price
	for _, price := range prices {
		if err := h.processPrice(ctx, price); err != nil {
			log.Error().
				Err(err).
				Str("productId", price.ProductID).
				Str("priceListId", price.PriceListID).
				Msg("Failed to process price")
			failedCount++
		} else {
			successCount++
		}
	}
	
	// Complete sync log
	syncLog.Status = "processed"
	syncLog.Duration = time.Since(startTime)
	syncLog.Payload, _ = json.Marshal(map[string]interface{}{
		"totalPrices": len(prices),
		"successful":  successCount,
		"failed":      failedCount,
		"priceListID": priceListID,
		"productIDs":  productIDs,
	})
	
	if err := h.syncLogRepo.Create(ctx, syncLog); err != nil {
		log.Error().Err(err).Msg("Failed to create sync log")
	}
	
	return nil
}

// Helper methods

func (h *SyncFromSAPHandler) processProductChange(ctx context.Context, change *sap.ProductChange) error {
	// Create canonical event
	event := &transformer.CanonicalEvent{
		EventID:        uuid.New(),
		EventType:      transformer.EventTypeProductMasterUpdated,
		EventTimestamp: time.Now(),
		Source:         h.transformer.Source,
		Payload: transformer.ProductMasterPayload{
			SKU:         change.SKU,
			Name:        change.Name,
			Description: change.Description,
			Category:    change.Category,
			Weight:      change.Weight,
			Dimensions: transformer.Dimensions{
				Length: change.Dimensions.Length,
				Width:  change.Dimensions.Width,
				Height: change.Dimensions.Height,
				Unit:   change.Dimensions.Unit,
			},
			Attributes: change.Attributes,
		},
	}
	
	// Publish event
	domainEvent := ddd.NewEvent(
		string(event.EventType),
		event.Payload,
	)
	
	if err := h.eventPublisher.Publish(ctx, domainEvent); err != nil {
		return fmt.Errorf("publishing product event: %w", err)
	}
	
	// Update sync status
	syncStatus := &domain.SyncStatus{
		ID:           uuid.New().String(),
		EntityType:   "product",
		EntityID:     change.ProductID,
		LastSyncedAt: time.Now(),
		Status:       "success",
		UpdatedAt:    time.Now(),
	}
	
	return h.syncStatusRepo.Save(ctx, syncStatus)
}

func (h *SyncFromSAPHandler) processStockLevel(ctx context.Context, stock *sap.StockLevel) error {
	// Create canonical event
	event := &transformer.CanonicalEvent{
		EventID:        uuid.New(),
		EventType:      transformer.EventTypeStockLevelUpdated,
		EventTimestamp: time.Now(),
		Source:         h.transformer.Source,
		Payload: transformer.StockLevelPayload{
			SKU:         stock.SKU,
			WarehouseID: stock.WarehouseID,
			Quantity:    stock.Quantity,
			StockType:   stock.StockType,
			UpdatedAt:   stock.UpdatedAt,
		},
	}
	
	// Publish event
	domainEvent := ddd.NewEvent(
		string(event.EventType),
		event.Payload,
	)
	
	if err := h.eventPublisher.Publish(ctx, domainEvent); err != nil {
		return fmt.Errorf("publishing stock event: %w", err)
	}
	
	// Update sync status
	syncStatus := &domain.SyncStatus{
		ID:           uuid.New().String(),
		EntityType:   "stock",
		EntityID:     fmt.Sprintf("%s-%s", stock.ProductID, stock.WarehouseID),
		LastSyncedAt: time.Now(),
		Status:       "success",
		UpdatedAt:    time.Now(),
	}
	
	return h.syncStatusRepo.Save(ctx, syncStatus)
}

func (h *SyncFromSAPHandler) processPrice(ctx context.Context, price *sap.Price) error {
	// Create canonical event
	event := &transformer.CanonicalEvent{
		EventID:        uuid.New(),
		EventType:      transformer.EventTypePriceUpdated,
		EventTimestamp: time.Now(),
		Source:         h.transformer.Source,
		Payload: transformer.PricePayload{
			SKU:         price.SKU,
			PriceListID: price.PriceListID,
			Currency:    price.Currency,
			Price:       price.Amount,
			ValidFrom:   price.ValidFrom,
			ValidTo:     price.ValidTo,
		},
	}
	
	// Publish event
	domainEvent := ddd.NewEvent(
		string(event.EventType),
		event.Payload,
	)
	
	if err := h.eventPublisher.Publish(ctx, domainEvent); err != nil {
		return fmt.Errorf("publishing price event: %w", err)
	}
	
	// Update sync status
	syncStatus := &domain.SyncStatus{
		ID:           uuid.New().String(),
		EntityType:   "price",
		EntityID:     fmt.Sprintf("%s-%s", price.ProductID, price.PriceListID),
		LastSyncedAt: time.Now(),
		Status:       "success",
		UpdatedAt:    time.Now(),
	}
	
	return h.syncStatusRepo.Save(ctx, syncStatus)
}