package grpc

import (
	"context"
	"fmt"
	"middleman/erp/internal/domain"
	"middleman/internal/rpc"
	"middleman/products/productspb"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProductRepository implements domain.ProductRepository using gRPC
type ProductRepository struct {
	endpoint string
}

var _ domain.ProductRepository = (*ProductRepository)(nil)

// NewProductRepository creates a new gRPC product repository
func NewProductRepository(endpoint string) ProductRepository {
	return ProductRepository{
		endpoint: endpoint,
	}
}

// AddProduct creates a new product via gRPC
func (r ProductRepository) AddProduct(ctx context.Context, product domain.Product) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return fmt.Errorf("dialing products service: %w", err)
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)

	// Convert domain attributes to proto attributes
	attrs := make([]*productspb.Attribute, len(product.Attributes))
	for i, attr := range product.Attributes {
		attrs[i] = &productspb.Attribute{
			Key:   attr.Key,
			Value: attr.Value,
		}
	}

	_, err = client.AddProduct(ctx, &productspb.AddProductRequest{
		Name:        product.Name,
		Description: product.Description,
		BasePrice:   product.BasePrice,
		Sku:         product.SKU,
		CategoryId:  product.CategoryID,
		Brand:       product.Brand,
		Status:      product.Status.String(),
		Weight:      product.Weight,
		Height:      product.Height,
		Width:       product.Width,
		Depth:       product.Depth,
		Stock:       product.Stock,
		Attributes:  attrs,
	})

	if err != nil {
		return fmt.Errorf("adding product: %w", err)
	}

	return nil
}

// UpdateProduct updates an existing product via gRPC
func (r ProductRepository) UpdateProduct(ctx context.Context, productID string, product domain.Product) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return fmt.Errorf("dialing products service: %w", err)
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)

	// The products service UpdateProduct API only updates specific fields
	// For a full update, we might need to call multiple endpoints
	_, err = client.UpdateProduct(ctx, &productspb.UpdateProductRequest{
		Id:          productID,
		Name:        product.Name,
		Description: product.Description,
		BasePrice:   product.BasePrice,
		Stock:       product.Stock,
		Sku:         product.SKU,
		CategoryId:  product.CategoryID,
		Status:      product.Status.String(),
	})

	if err != nil {
		return fmt.Errorf("updating product: %w", err)
	}

	return nil
}

// UpdateProductStock updates product stock levels via gRPC
func (r ProductRepository) UpdateProductStock(ctx context.Context, productID string, stock int64) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return fmt.Errorf("dialing products service: %w", err)
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)

	_, err = client.AdjustProductStock(ctx, &productspb.AdjustProductStockRequest{
		ProductId: productID,
		NewStock:  stock,
	})

	if err != nil {
		return fmt.Errorf("updating product stock: %w", err)
	}

	return nil
}

// UpdateProductPrice updates product price via gRPC
func (r ProductRepository) UpdateProductPrice(ctx context.Context, productID string, price int64) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return fmt.Errorf("dialing products service: %w", err)
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)

	// Get current product to get old price
	product, err := r.GetProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("getting current product: %w", err)
	}

	_, err = client.UpdateProductPrice(ctx, &productspb.UpdateProductPriceRequest{
		Id:       productID,
		NewPrice: price,
		OldPrice: product.BasePrice,
	})

	if err != nil {
		return fmt.Errorf("updating product price: %w", err)
	}

	return nil
}

// GetProductBySKU retrieves a product by its SKU via gRPC
func (r ProductRepository) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, fmt.Errorf("dialing products service: %w", err)
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)

	// Use GetProductBySKU if available, otherwise use GetProductsWithFilters
	resp, err := client.GetProductBySKU(ctx, &productspb.GetProductBySKURequest{
		Sku: sku,
	})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil // Product not found
		}
		return nil, fmt.Errorf("getting product by SKU: %w", err)
	}

	return r.productToDomain(resp.GetProduct()), nil
}

// GetProductByID retrieves a product by its ID via gRPC
func (r ProductRepository) GetProductByID(ctx context.Context, productID string) (*domain.Product, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, fmt.Errorf("dialing products service: %w", err)
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)

	resp, err := client.GetProduct(ctx, &productspb.GetProductRequest{
		Id: productID,
	})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil // Product not found
		}
		return nil, fmt.Errorf("getting product by ID: %w", err)
	}

	return r.productToDomain(resp.GetProduct()), nil
}

// productToDomain converts a protobuf product to domain product
func (r ProductRepository) productToDomain(pb *productspb.Product) *domain.Product {
	if pb == nil {
		return nil
	}

	attrs := make([]domain.Attribute, len(pb.GetAttributes()))
	for i, a := range pb.GetAttributes() {
		attrs[i] = domain.Attribute{
			Key:   a.GetKey(),
			Value: a.GetValue(),
		}
	}

	return &domain.Product{
		ProductID:    pb.GetId(),
		Name:         pb.GetName(),
		Description:  pb.GetDescription(),
		SKU:          pb.GetSku(),
		BasePrice:    pb.GetBasePrice(),
		Stock:        pb.GetStock(),
		CategoryID:   pb.GetCategoryId(),
		Brand:        pb.GetBrand(),
		Status:       domain.ToProductStatus(pb.GetStatus()),
		Attributes:   attrs,
		Weight:       pb.GetWeight(),
		Height:       pb.GetHeight(),
		Width:        pb.GetWidth(),
		Depth:        pb.GetDepth(),
		UserSellerID: pb.GetUserSellerId(),
	}
}

// dial creates a gRPC connection
func (r ProductRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	log.Debug().Str("endpoint", r.endpoint).Msg("Dialing products service")
	return rpc.Dial(ctx, r.endpoint)
}
