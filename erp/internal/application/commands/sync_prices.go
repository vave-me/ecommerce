package commands

import (
	"context"
	"fmt"
	"middleman/internal/erp"
	"time"

	"middleman/erp/internal/domain"
)

// SyncPrices command triggers price synchronization from ERP
type SyncPrices struct {
	ConnectorID  string
	ProductIDs   []string // If empty, sync all prices
	PriceListIDs []string // Specific price lists to sync
	Since        time.Time
	BatchSize    int // Batch size for syncing
}

// SyncPricesHandler handles the SyncPrices command
type SyncPricesHandler struct {
	registry       erp.ConnectorRegistry
	repository     domain.SyncLogRepository
	productRepo    domain.ProductRepository
}

// NewSyncPricesHandler creates a new handler
func NewSyncPricesHandler(
	registry erp.ConnectorRegistry,
	repository domain.SyncLogRepository,
	productRepo domain.ProductRepository,
) SyncPricesHandler {
	return SyncPricesHandler{
		registry:       registry,
		repository:     repository,
		productRepo:    productRepo,
	}
}

// SyncPrices synchronizes prices from the ERP
func (h SyncPricesHandler) SyncPrices(ctx context.Context, cmd SyncPrices) error {
	// Get connector
	connector, err := h.registry.GetConnector(cmd.ConnectorID)
	if err != nil {
		return fmt.Errorf("getting connector: %w", err)
	}

	// Create sync log
	syncLog := &domain.SyncLog{
		ID:          generatePriceSyncID(),
		ConnectorID: cmd.ConnectorID,
		EntityType:  "price",
		StartedAt:   time.Now(),
		Status:      domain.SyncStatusInProgress,
		Metadata: map[string]interface{}{
			"product_ids": cmd.ProductIDs,
			"priceLists":  cmd.PriceListIDs,
		},
	}

	if err := h.repository.Create(ctx, syncLog); err != nil {
		return fmt.Errorf("creating sync log: %w", err)
	}

	// Set default batch size
	batchSize := cmd.BatchSize
	if batchSize == 0 {
		batchSize = 100
	}

	// Sync prices
	prices, err := connector.SyncPrices(ctx, cmd.ProductIDs, batchSize)
	if err != nil {
		syncLog.Status = domain.SyncStatusFailed
		syncLog.Error = err.Error()
		syncLog.CompletedAt = ptrTime(time.Now())
		h.repository.Update(ctx, syncLog)
		return fmt.Errorf("syncing prices: %w", err)
	}

	processedCount := 0

	// Process price updates
	for _, price := range prices {
		// Validate price data
		if price.SKU == "" || price.Price < 0 {
			syncLog.Error = fmt.Sprintf("invalid price data for SKU %s", price.SKU)
			continue
		}

		// Get product by SKU to get the product ID
		product, err := h.productRepo.GetProductBySKU(ctx, price.SKU)
		if err != nil {
			syncLog.Error = fmt.Sprintf("failed to get product for SKU %s: %v", price.SKU, err)
			continue
		}
		
		if product == nil {
			syncLog.Error = fmt.Sprintf("product not found for SKU %s", price.SKU)
			continue
		}

		// Update product price (convert float64 to int64 cents)
		priceInCents := int64(price.Price * 100)
		err = h.productRepo.UpdateProductPrice(ctx, product.ProductID, priceInCents)
		if err != nil {
			syncLog.Error = fmt.Sprintf("failed to update price for SKU %s: %v", price.SKU, err)
			continue
		}
		
		processedCount++
	}

	// Update sync log
	syncLog.Status = domain.SyncStatusCompleted
	syncLog.CompletedAt = ptrTime(time.Now())
	syncLog.RecordsProcessed = processedCount
	syncLog.RecordsTotal = len(prices)

	if err := h.repository.Update(ctx, syncLog); err != nil {
		return fmt.Errorf("updating sync log: %w", err)
	}

	return nil
}

func generatePriceSyncID() string {
	return fmt.Sprintf("price_sync_%d", time.Now().UnixNano())
}
