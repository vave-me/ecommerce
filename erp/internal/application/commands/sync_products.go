package commands

import (
	"context"
	"fmt"
	"middleman/internal/erp"
	"time"

	"middleman/erp/internal/domain"
)

// SyncProducts command triggers product synchronization from ERP
type SyncProducts struct {
	ConnectorID string
	Since       time.Time
	BatchSize   int
}

// SyncProductsHandler handles the SyncProducts command
type SyncProductsHandler struct {
	registry       erp.ConnectorRegistry
	repository     domain.SyncLogRepository
	productRepo    domain.ProductRepository
}

// NewSyncProductsHandler creates a new handler
func NewSyncProductsHandler(
	registry erp.ConnectorRegistry,
	repository domain.SyncLogRepository,
	productRepo domain.ProductRepository,
) SyncProductsHandler {
	return SyncProductsHandler{
		registry:       registry,
		repository:     repository,
		productRepo:    productRepo,
	}
}

// SyncProducts synchronizes products from the ERP
func (h SyncProductsHandler) SyncProducts(ctx context.Context, cmd SyncProducts) error {
	// Get connector
	connector, err := h.registry.GetConnector(cmd.ConnectorID)
	if err != nil {
		return fmt.Errorf("getting connector: %w", err)
	}

	// Create sync log
	syncLog := &domain.SyncLog{
		ID:          generateProductSyncID(),
		ConnectorID: cmd.ConnectorID,
		EntityType:  "product",
		StartedAt:   time.Now(),
		Status:      domain.SyncStatusInProgress,
	}

	if err := h.repository.Create(ctx, syncLog); err != nil {
		return fmt.Errorf("creating sync log: %w", err)
	}

	// Sync products
	products, err := connector.SyncProducts(ctx, cmd.Since, cmd.BatchSize)
	if err != nil {
		syncLog.Status = domain.SyncStatusFailed
		syncLog.Error = err.Error()
		syncLog.CompletedAt = ptrTime(time.Now())
		h.repository.Update(ctx, syncLog)
		return fmt.Errorf("fetching products: %w", err)
	}

	// Process products in batches
	batchSize := cmd.BatchSize
	if batchSize == 0 {
		batchSize = 100
	}

	processedCount := 0
	failedCount := 0
	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}

		batch := products[i:end]
		
		// Process each product in the batch
		for _, productPayload := range batch {
			// Check if product exists by SKU
			existingProduct, err := h.productRepo.GetProductBySKU(ctx, productPayload.SKU)
			if err != nil {
				syncLog.Error = fmt.Sprintf("failed to check product SKU %s: %v", productPayload.SKU, err)
				failedCount++
				continue
			}

			// Convert ERP product payload to domain product
			domainProduct := domain.Product{
				Name:         productPayload.Name,
				Description:  productPayload.Description,
				SKU:          productPayload.SKU,
				BasePrice:    0, // Price will be synced separately via sync_prices
				Stock:        0, // Stock will be synced separately via sync_stock
				CategorySlug: productPayload.Category, // Use category as slug
				Brand:        productPayload.Brand,
				Status:       domain.ProductStatusActive,
				Weight:       int64(productPayload.Weight * 1000), // Convert kg to grams
				UserSellerID: "system", // ERP sync uses system user
			}

			if existingProduct != nil {
				// Update existing product
				err = h.productRepo.UpdateProduct(ctx, existingProduct.ProductID, domainProduct)
				if err != nil {
					syncLog.Error = fmt.Sprintf("failed to update product SKU %s: %v", productPayload.SKU, err)
					failedCount++
					continue
				}
			} else {
				// Create new product
				err = h.productRepo.AddProduct(ctx, domainProduct)
				if err != nil {
					syncLog.Error = fmt.Sprintf("failed to add product SKU %s: %v", productPayload.SKU, err)
					failedCount++
					continue
				}
			}
			
			processedCount++
		}
	}

	// Update sync log
	syncLog.Status = domain.SyncStatusCompleted
	if failedCount > 0 {
		syncLog.Status = domain.SyncStatusFailed
		syncLog.Error = fmt.Sprintf("%d out of %d products failed to sync", failedCount, len(products))
	}
	syncLog.CompletedAt = ptrTime(time.Now())
	syncLog.RecordsProcessed = processedCount
	syncLog.RecordsTotal = len(products)

	if err := h.repository.Update(ctx, syncLog); err != nil {
		return fmt.Errorf("updating sync log: %w", err)
	}

	return nil
}

func generateProductSyncID() string {
	return fmt.Sprintf("product_sync_%d", time.Now().UnixNano())
}
