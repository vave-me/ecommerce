package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/media/internal/domain"
)

type (
	AddImportBatch struct {
		SessionID string
		Items     []ImportItem
	}

	ImportItem struct {
		ExternalID   string
		SKU          string
		ImageURL     string
		Metadata     map[string]string
		DisplayOrder int32
	}

	AddImportBatchHandler struct {
		imports          domain.ImportSessionRepository
		items            domain.ImportItemRepository
		productRepo      ProductRepository
		publisher        ddd.EventPublisher[ddd.Event]
	}
	
	ProductRepository interface {
		ValidateSKUs(ctx context.Context, skus []string) (map[string]string, []string, error)
	}
)

func NewAddImportBatchHandler(imports domain.ImportSessionRepository, items domain.ImportItemRepository, productRepo ProductRepository, publisher ddd.EventPublisher[ddd.Event]) AddImportBatchHandler {
	return AddImportBatchHandler{
		imports:     imports,
		items:       items,
		productRepo: productRepo,
		publisher:   publisher,
	}
}

func (h AddImportBatchHandler) AddImportBatch(ctx context.Context, cmd AddImportBatch) error {
	// Verify session exists and is active
	session, err := h.imports.Get(ctx, cmd.SessionID)
	if err != nil {
		return err
	}

	if session.Status != domain.ImportStatusPending && session.Status != domain.ImportStatusProcessing {
		return domain.ErrImportSessionNotActive
	}

	// Extract SKUs from items for validation
	skusToValidate := make([]string, 0, len(cmd.Items))
	skuToItemIndex := make(map[string][]int)
	
	for i, item := range cmd.Items {
		if item.SKU != "" {
			skusToValidate = append(skusToValidate, item.SKU)
			skuToItemIndex[item.SKU] = append(skuToItemIndex[item.SKU], i)
		}
	}

	// Validate SKUs with products service
	var invalidSKUs []string
	skuToProductID := make(map[string]string)
	
	if len(skusToValidate) > 0 && h.productRepo != nil {
		validSKUs, notFoundSKUs, err := h.productRepo.ValidateSKUs(ctx, skusToValidate)
		if err != nil {
			return errors.Wrap(err, "validating SKUs with product service")
		}
		skuToProductID = validSKUs
		invalidSKUs = notFoundSKUs
	}

	// Check if any SKUs are invalid
	if len(invalidSKUs) > 0 {
		return errors.ErrBadRequest.Msgf("invalid SKUs found: %v", invalidSKUs)
	}

	// Convert and save import items with product IDs
	var domainItems []*domain.ImportItem
	for _, item := range cmd.Items {
		domainItem := &domain.ImportItem{
			SessionID:    cmd.SessionID,
			ExternalID:   item.ExternalID,
			SKU:          item.SKU,
			ImageURL:     item.ImageURL,
			Status:       domain.ImportItemStatusPending,
			RetryCount:   0,
			Metadata:     item.Metadata,
			DisplayOrder: item.DisplayOrder,
		}
		
		// Add product ID if SKU was validated
		if item.SKU != "" && skuToProductID != nil {
			if productID, ok := skuToProductID[item.SKU]; ok {
				domainItem.ProductID = productID
			}
		}
		
		domainItems = append(domainItems, domainItem)
	}

	if err := h.items.CreateBatch(ctx, domainItems); err != nil {
		return err
	}

	// Update session total count
	session.TotalImages += int32(len(cmd.Items))
	if err := h.imports.Update(ctx, session); err != nil {
		return err
	}

	// Publish batch added event
	event := &domain.ImportBatchAdded{
		SessionID:   cmd.SessionID,
		BatchSize:   int32(len(cmd.Items)),
		BatchNumber: session.BatchCount + 1,
	}

	return h.publisher.Publish(ctx, ddd.NewEvent(event))
}