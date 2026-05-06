// merchant/application/service.go

package application

import (
	"context"
	"fmt"
	"middleman/internal/ddd"
	"middleman/merchant/internal/domain"
	"time"

	"google.golang.org/api/content/v2.1"
)

type (
	UpsertProductCommand struct {
		Product *content.Product
	}
	DeleteProductCommand struct {
		ProductID string
	}
	BatchUpsertProductsCommand struct {
		Products []*content.Product
	}
	BatchDeleteProductsCommand struct {
		ProductIDs []string
	}

	ListProductsQuery struct {
		PageSize  int64
		PageToken string
	}
	ListAllProductsQuery struct {
		PageSize int64
	}

	GetProductQuery struct {
		ProductID string
	}
)

// MerchantService defines the high-level interface for merchant domain commands/queries.
type App interface {
	// Commands
	UpsertProduct(ctx context.Context, cmd UpsertProductCommand) error
	DeleteProduct(ctx context.Context, cmd DeleteProductCommand) error
	BatchUpsertProducts(ctx context.Context, cmd BatchUpsertProductsCommand) error
	BatchDeleteProducts(ctx context.Context, cmd BatchDeleteProductsCommand) error

	// Queries
	ListProducts(ctx context.Context, q ListProductsQuery) ([]*content.Product, string, error)
	ListAllProducts(ctx context.Context, q ListAllProductsQuery) ([]*content.Product, error)
	GetProduct(ctx context.Context, q GetProductQuery) (*content.Product, error)
}

type Application struct {
	client     domain.MerchantClient
	publisher  ddd.EventPublisher[ddd.Event]
	syncRepo   domain.ProductSyncStatusRepository
	validator  *domain.ProductValidator
}

// Ensure it implements App fully
var _ App = (*Application)(nil)

// New constructor
func New(
	client domain.MerchantClient,
	publisher ddd.EventPublisher[ddd.Event],
	syncRepo domain.ProductSyncStatusRepository,
) *Application {
	return &Application{
		client:    client,
		publisher: publisher,
		syncRepo:  syncRepo,
		validator: domain.NewProductValidator(),
	}
}

// UpsertProduct checks if product exists (Get), if not found => Insert, else => Update
func (a *Application) UpsertProduct(ctx context.Context, cmd UpsertProductCommand) error {
	if cmd.Product == nil || cmd.Product.Id == "" {
		return fmt.Errorf("invalid product: Id is required")
	}

	// Validate product according to Google Merchant Center requirements
	if validationErrors := a.validator.ValidateProduct(cmd.Product); len(validationErrors) > 0 {
		return fmt.Errorf("product validation failed: %w", validationErrors)
	}

	// Track sync status
	syncStatus := &domain.ProductSyncStatus{
		ProductID:  cmd.Product.OfferId,
		MerchantID: fmt.Sprintf("%d", a.client.MerchantID()),
		SyncStatus: domain.SyncStatusPending,
	}

	// Check if sync status exists
	existingStatus, err := a.syncRepo.FindByProductID(ctx, cmd.Product.OfferId)
	if err != nil {
		return fmt.Errorf("error checking sync status: %w", err)
	}

	if existingStatus == nil {
		// Create new sync status
		if err := a.syncRepo.Create(ctx, syncStatus); err != nil {
			return fmt.Errorf("error creating sync status: %w", err)
		}
	}

	// Attempt to Get existing product from Google Merchant Center
	_, err = a.client.GetProduct(ctx, cmd.Product.Id)
	if err != nil {
		// If not found => Insert
		if a.client.IsNotFoundErr(err) {
			err = a.client.InsertProduct(ctx, cmd.Product)
		} else {
			return a.updateSyncStatus(ctx, cmd.Product.OfferId, domain.SyncStatusFailed, err.Error())
		}
	} else {
		// Otherwise => Update
		err = a.client.UpdateProduct(ctx, cmd.Product)
	}

	// Update sync status based on result
	if err != nil {
		return a.updateSyncStatus(ctx, cmd.Product.OfferId, domain.SyncStatusFailed, err.Error())
	}

	return a.updateSyncStatus(ctx, cmd.Product.OfferId, domain.SyncStatusSynced, "")
}

// DeleteProduct
func (a *Application) DeleteProduct(ctx context.Context, cmd DeleteProductCommand) error {
	if cmd.ProductID == "" {
		return fmt.Errorf("productID is required")
	}
	
	err := a.client.DeleteProduct(ctx, cmd.ProductID)
	if err != nil {
		return fmt.Errorf("error deleting product from merchant center: %w", err)
	}
	
	// Update sync status to removed
	return a.updateSyncStatus(ctx, cmd.ProductID, domain.SyncStatusRemoved, "")
}

func (a *Application) BatchUpsertProducts(ctx context.Context, cmd BatchUpsertProductsCommand) error {
	if len(cmd.Products) == 0 {
		return nil
	}

	// Google Merchant Center recommends batches of up to 1000 products
	const batchSize = 1000
	var errors []error

	for i := 0; i < len(cmd.Products); i += batchSize {
		end := i + batchSize
		if end > len(cmd.Products) {
			end = len(cmd.Products)
		}

		batch := cmd.Products[i:end]
		
		// Process batch with rate limiting
		for _, p := range batch {
			err := a.UpsertProduct(ctx, UpsertProductCommand{Product: p})
			if err != nil {
				errors = append(errors, fmt.Errorf("failed to upsert product %s: %w", p.OfferId, err))
			}
			
			// Basic rate limiting - Google allows 2500 requests per day
			// This gives us roughly 1 request per 35 seconds, but we'll be more conservative
			time.Sleep(100 * time.Millisecond)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("batch upsert completed with %d errors", len(errors))
	}

	return nil
}

func (a *Application) BatchDeleteProducts(ctx context.Context, cmd BatchDeleteProductsCommand) error {
	if len(cmd.ProductIDs) == 0 {
		return nil
	}

	var errors []error

	for _, pid := range cmd.ProductIDs {
		err := a.DeleteProduct(ctx, DeleteProductCommand{ProductID: pid})
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to delete product %s: %w", pid, err))
		}
		
		// Basic rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	if len(errors) > 0 {
		return fmt.Errorf("batch delete completed with %d errors", len(errors))
	}

	return nil
}

func (a *Application) ListProducts(ctx context.Context, q ListProductsQuery) ([]*content.Product, string, error) {
	return a.client.ListProducts(ctx, q.PageSize, q.PageToken)
}

func (a *Application) ListAllProducts(ctx context.Context, q ListAllProductsQuery) ([]*content.Product, error) {
	if q.PageSize <= 0 {
		q.PageSize = 100
	}

	var allProducts []*content.Product
	pageToken := ""

	for {
		prods, nextToken, err := a.client.ListProducts(ctx, q.PageSize, pageToken)
		if err != nil {
			return nil, err
		}
		allProducts = append(allProducts, prods...)
		if nextToken == "" {
			break
		}
		pageToken = nextToken
	}
	return allProducts, nil
}

func (a *Application) GetProduct(ctx context.Context, q GetProductQuery) (*content.Product, error) {
	if q.ProductID == "" {
		return nil, domain.ErrInvalidProduct
	}
	
	product, err := a.client.GetProduct(ctx, q.ProductID)
	if err != nil {
		if a.client.IsNotFoundErr(err) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}
	
	return product, nil
}

// updateSyncStatus is a helper method to update sync status
func (a *Application) updateSyncStatus(ctx context.Context, productID string, status string, errorMsg string) error {
	syncStatus := &domain.ProductSyncStatus{
		ProductID:  productID,
		MerchantID: fmt.Sprintf("%d", a.client.MerchantID()),
		SyncStatus: status,
		LastError:  errorMsg,
	}
	
	if status == domain.SyncStatusSynced {
		now := time.Now()
		syncStatus.LastSyncedAt = &now
	}
	
	return a.syncRepo.Update(ctx, syncStatus)
}
