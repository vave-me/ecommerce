package commands

import (
	"context"
	"fmt"
	"middleman/internal/erp"
	"time"

	"middleman/erp/internal/domain"
)

// SyncStock command triggers stock synchronization from ERP
type SyncStock struct {
	ConnectorID string
	ProductIDs  []string // If empty, sync all stock
	Since       time.Time
	BatchSize   int
}

// SyncStockHandler handles the SyncStock command
type SyncStockHandler struct {
	registry       erp.ConnectorRegistry
	repository     domain.SyncLogRepository
	productRepo    domain.ProductRepository
}

// NewSyncStockHandler creates a new handler
func NewSyncStockHandler(
	registry erp.ConnectorRegistry,
	repository domain.SyncLogRepository,
	productRepo domain.ProductRepository,
) SyncStockHandler {
	return SyncStockHandler{
		registry:       registry,
		repository:     repository,
		productRepo:    productRepo,
	}
}

// SyncStock synchronizes stock levels from the ERP
func (h SyncStockHandler) SyncStock(ctx context.Context, cmd SyncStock) error {
	// Get connector
	connector, err := h.registry.GetConnector(cmd.ConnectorID)
	if err != nil {
		return fmt.Errorf("getting connector: %w", err)
	}

	// Create sync log
	syncLog := &domain.SyncLog{
		ID:          generateStockSyncID(),
		ConnectorID: cmd.ConnectorID,
		EntityType:  "stock",
		StartedAt:   time.Now(),
		Status:      domain.SyncStatusInProgress,
		Metadata: map[string]interface{}{
			"product_ids": cmd.ProductIDs,
			"batch_size":  cmd.BatchSize,
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

	// Sync stock levels
	stockLevels, err := connector.SyncStock(ctx, cmd.ProductIDs, batchSize)
	if err != nil {
		syncLog.Status = domain.SyncStatusFailed
		syncLog.Error = err.Error()
		syncLog.CompletedAt = ptrTime(time.Now())
		h.repository.Update(ctx, syncLog)
		return fmt.Errorf("syncing stock: %w", err)
	}

	// Process stock updates
	processedCount := 0
	failedCount := 0

	// Process in batches
	for i := 0; i < len(stockLevels); i += batchSize {
		end := i + batchSize
		if end > len(stockLevels) {
			end = len(stockLevels)
		}

		batch := stockLevels[i:end]
		
		// Process each stock level in the batch
		for _, stockPayload := range batch {
			// Get product by SKU to get the product ID
			product, err := h.productRepo.GetProductBySKU(ctx, stockPayload.SKU)
			if err != nil {
				syncLog.Error = fmt.Sprintf("failed to get product for SKU %s: %v", stockPayload.SKU, err)
				failedCount++
				continue
			}
			
			if product == nil {
				syncLog.Error = fmt.Sprintf("product not found for SKU %s", stockPayload.SKU)
				failedCount++
				continue
			}

			// Update stock level (convert int to int64)
			err = h.productRepo.UpdateProductStock(ctx, product.ProductID, int64(stockPayload.Quantity))
			if err != nil {
				syncLog.Error = fmt.Sprintf("failed to update stock for SKU %s: %v", stockPayload.SKU, err)
				failedCount++
				continue
			}
			
			processedCount++
		}
	}

	// Update sync log
	syncLog.Status = domain.SyncStatusCompleted
	if failedCount > 0 {
		syncLog.Status = domain.SyncStatusFailed
		syncLog.Error = fmt.Sprintf("%d out of %d stock updates failed to sync", failedCount, len(stockLevels))
	}
	syncLog.CompletedAt = ptrTime(time.Now())
	syncLog.RecordsProcessed = processedCount
	syncLog.RecordsTotal = len(stockLevels)

	if err := h.repository.Update(ctx, syncLog); err != nil {
		return fmt.Errorf("updating sync log: %w", err)
	}

	return nil
}

func generateStockSyncID() string {
	return fmt.Sprintf("stock_sync_%d", time.Now().UnixNano())
}
