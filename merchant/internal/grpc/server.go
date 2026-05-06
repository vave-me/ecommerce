package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/api/content/v2.1"
	"fmt"

	"middleman/merchant/internal/application"
	"middleman/merchant/internal/domain"
	"middleman/merchant/merchantpb"
)

// server implements the MerchantServiceServer gRPC interface.
type server struct {
	app      application.App
	syncRepo domain.ProductSyncStatusRepository
	merchantpb.UnimplementedMerchantServiceServer
}

var _ merchantpb.MerchantServiceServer = (*server)(nil)

// RegisterServer registers this gRPC server with the provided registrar.
func RegisterServer(_ context.Context, app application.App, syncRepo domain.ProductSyncStatusRepository, registrar grpc.ServiceRegistrar) error {
	merchantpb.RegisterMerchantServiceServer(registrar, server{
		app:      app,
		syncRepo: syncRepo,
	})
	return nil
}

// SyncProduct syncs a single product to Google Merchant Center
func (s server) SyncProduct(ctx context.Context, req *merchantpb.SyncProductRequest) (*merchantpb.SyncProductResponse, error) {
	if req.Product == nil {
		return nil, status.Error(codes.InvalidArgument, "product is required")
	}

	// Convert merchantpb.Product to content.Product
	contentProduct := s.convertToContentProduct(req.Product)
	
	err := s.app.UpsertProduct(ctx, application.UpsertProductCommand{
		Product: contentProduct,
	})
	
	if err != nil {
		return &merchantpb.SyncProductResponse{
			MerchantProductId: req.ProductId,
			Success:           false,
			ErrorMessage:      err.Error(),
		}, nil
	}

	return &merchantpb.SyncProductResponse{
		MerchantProductId: contentProduct.Id,
		Success:           true,
	}, nil
}

// BatchSyncProducts syncs multiple products to Google Merchant Center
func (s server) BatchSyncProducts(ctx context.Context, req *merchantpb.BatchSyncProductsRequest) (*merchantpb.BatchSyncProductsResponse, error) {
	if len(req.Products) == 0 {
		return &merchantpb.BatchSyncProductsResponse{}, nil
	}

	var contentProducts []*content.Product
	for _, p := range req.Products {
		contentProducts = append(contentProducts, s.convertToContentProduct(p))
	}

	_ = s.app.BatchUpsertProducts(ctx, application.BatchUpsertProductsCommand{
		Products: contentProducts,
	})

	// Get sync results
	var results []*merchantpb.SyncResult
	var successCount, failedCount int32

	for _, p := range req.Products {
		status, err := s.syncRepo.FindByProductID(ctx, p.Id)
		if err != nil {
			continue
		}

		result := &merchantpb.SyncResult{
			ProductId: p.Id,
			Success:   status.SyncStatus == domain.SyncStatusSynced,
		}

		if !result.Success {
			result.ErrorMessage = status.LastError
			failedCount++
		} else {
			successCount++
		}

		results = append(results, result)
	}

	return &merchantpb.BatchSyncProductsResponse{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Results:      results,
	}, nil
}

// RemoveProduct removes a product from Google Merchant Center
func (s server) RemoveProduct(ctx context.Context, req *merchantpb.RemoveProductRequest) (*merchantpb.RemoveProductResponse, error) {
	err := s.app.DeleteProduct(ctx, application.DeleteProductCommand{
		ProductID: req.ProductId,
	})
	if err != nil {
		return &merchantpb.RemoveProductResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}
	return &merchantpb.RemoveProductResponse{
		Success: true,
	}, nil
}

// GetProductStatus gets the sync status of a product
func (s server) GetProductStatus(ctx context.Context, req *merchantpb.GetProductStatusRequest) (*merchantpb.GetProductStatusResponse, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product ID is required")
	}

	syncStatus, err := s.syncRepo.FindByProductID(ctx, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get sync status: %v", err)
	}

	if syncStatus == nil {
		return nil, status.Error(codes.NotFound, "product sync status not found")
	}

	response := &merchantpb.GetProductStatusResponse{
		Status: s.convertSyncStatus(syncStatus.SyncStatus),
	}

	if syncStatus.LastSyncedAt != nil {
		response.LastSyncedAt = timestamppb.New(*syncStatus.LastSyncedAt)
	}

	if syncStatus.LastError != "" {
		response.Warnings = []string{syncStatus.LastError}
	}

	return response, nil
}

// ListProducts lists all synced products
func (s server) ListProducts(ctx context.Context, req *merchantpb.ListProductsRequest) (*merchantpb.ListProductsResponse, error) {
	pageSize := int64(req.PageSize)
	if pageSize <= 0 {
		pageSize = 50
	}

	products, nextToken, err := s.app.ListProducts(ctx, application.ListProductsQuery{
		PageSize:  pageSize,
		PageToken: req.PageToken,
	})
	if err != nil {
		return nil, err
	}

	// Convert content.Product to merchantpb.Product
	var pbProducts []*merchantpb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, s.convertFromContentProduct(p))
	}

	return &merchantpb.ListProductsResponse{
		Products:      pbProducts,
		NextPageToken: nextToken,
		TotalCount:    int32(len(pbProducts)),
	}, nil
}

// convertToContentProduct converts merchantpb.Product to content.Product
func (s server) convertToContentProduct(p *merchantpb.Product) *content.Product {
	product := &content.Product{
		OfferId:               p.Id,
		Title:                 p.Name,
		Description:           p.Description,
		Link:                  p.Link,
		ImageLink:             p.ImageUrl,
		ContentLanguage:       "en",
		TargetCountry:         "US",
		Availability:          p.Availability,
		Condition:             p.Condition,
		Brand:                 p.Brand,
		GoogleProductCategory: p.GoogleProductCategory,
	}

	// Convert price
	if p.Price > 0 {
		priceValue := fmt.Sprintf("%.2f", float64(p.Price)/100.0)
		product.Price = &content.Price{
			Value:    priceValue,
			Currency: p.Currency,
		}
	}

	// Set SKU if provided
	if p.Sku != "" {
		product.Gtin = p.Sku
	}

	// Generate ID for Google Merchant Center
	product.Id = fmt.Sprintf("online:en:%s", p.Id)

	return product
}

// convertFromContentProduct converts content.Product to merchantpb.Product
func (s server) convertFromContentProduct(p *content.Product) *merchantpb.Product {
	product := &merchantpb.Product{
		Id:                    p.OfferId,
		Name:                  p.Title,
		Description:           p.Description,
		Link:                  p.Link,
		ImageUrl:              p.ImageLink,
		Availability:          p.Availability,
		Condition:             p.Condition,
		Brand:                 p.Brand,
		GoogleProductCategory: p.GoogleProductCategory,
	}

	// Convert price
	if p.Price != nil {
		// Parse price value and convert to cents
		var priceFloat float64
		fmt.Sscanf(p.Price.Value, "%f", &priceFloat)
		product.Price = int64(priceFloat * 100)
		product.Currency = p.Price.Currency
	}

	// Set SKU from GTIN if available
	if p.Gtin != "" {
		product.Sku = p.Gtin
	}

	return product
}

// convertSyncStatus converts domain sync status to protobuf enum
func (s server) convertSyncStatus(status string) merchantpb.ProductStatus {
	switch status {
	case domain.SyncStatusSynced:
		return merchantpb.ProductStatus_PRODUCT_STATUS_SYNCED
	case domain.SyncStatusPending:
		return merchantpb.ProductStatus_PRODUCT_STATUS_PENDING
	case domain.SyncStatusFailed:
		return merchantpb.ProductStatus_PRODUCT_STATUS_FAILED
	case domain.SyncStatusRemoved:
		return merchantpb.ProductStatus_PRODUCT_STATUS_REMOVED
	default:
		return merchantpb.ProductStatus_PRODUCT_STATUS_UNSPECIFIED
	}
}