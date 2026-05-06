package grpc

import (
	"context"
	"fmt"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/products/productspb"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type VariantRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.VariantRepository = (*VariantRepository)(nil)

// NewVariantRepositoryWithAuth creates a new VariantRepository with JWT authentication support
func NewVariantRepository(endpoint string, authInstance *auth.Auth) VariantRepository {
	return VariantRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// Find retrieves a product by its ID using a gRPC call to the UsersService.
func (r VariantRepository) Find(ctx context.Context, variantID string) (product *models.Variant, err error) {
	log.Printf("Find: Starting to retrieve variant with ID: %s", variantID)

	// Attempt to establish a gRPC connection.
	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("Find: Failed to dial gRPC server for variant ID %s: %v", variantID, err)
		return nil, err
	}
	log.Printf("Find: Successfully dialed gRPC server for variantID ID %s", variantID)

	// Ensure the connection is closed when the function exits.
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Printf("Find: Error closing gRPC connection for variant ID %s: %v", variantID, err)
		} else {
			log.Printf("Find: Successfully closed gRPC connection for variant ID %s", variantID)
		}
	}(conn)

	// Create a new UsersService client.
	productsClient := productspb.NewProductsServiceClient(conn)
	log.Printf("Find: Created UsersService client for variant ID %s", variantID)

	// Prepare the GetProduct request.
	request := &productspb.GetVariantRequest{VariantId: variantID}
	log.Printf("Find: Sending GetProduct request for variant ID %s", variantID)

	// Make the gRPC call to get the product.
	resp, err := productsClient.GetVariant(ctx, request)
	if err != nil {
		log.Printf("Find: GetProduct RPC failed for product ID %s: %v", variantID, err)
		return nil, err
	}
	log.Printf("Find: Received response from GetProduct RPC for variant ID %s", variantID)

	// Convert the gRPC product to your domain model.
	variant := r.variantToDomain(resp.Variant)
	if variant == nil {
		log.Printf("Find: Conversion from gRPC product to domain model failed for product ID %s", variantID)
		err = fmt.Errorf("conversion failed for product ID %s", variantID)
		return nil, err
	}
	log.Printf("Find: Successfully converted gRPC product to domain model for product ID %s: %+v", variantID, variant)

	log.Printf("Find: Successfully retrieved product with ID %s", variantID)
	return variant, nil
}

func (r VariantRepository) Add(ctx context.Context, productID string, name string, description string, basePrice int64, userSellerID string, categoryID string, brand string, condition string, model string, tags []string) (err error) {

	var conn *grpc.ClientConn
	conn, err = r.dial(ctx)
	if err != nil {
		return err
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	_, err = productspb.NewProductsServiceClient(conn).AddProduct(ctx, &productspb.AddProductRequest{
		Name:        name,
		Description: description,
		BasePrice:   basePrice,
		CategoryId:  categoryID,
	})
	if err != nil {
		return err
	}

	return nil
}
func (r VariantRepository) Update(ctx context.Context, variantID string, price, stock int64, name string, attributes []models.Attribute) (err error) {
	var conn *grpc.ClientConn
	conn, err = r.dial(ctx)
	if err != nil {
		return err
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	pbAttrs := make([]*productspb.Attribute, len(attributes))
	for i, a := range attributes {
		pbAttrs[i] = &productspb.Attribute{
			Key:   a.Key,
			Value: a.Value,
		}
	}

	_, err = productspb.NewProductsServiceClient(conn).UpdateProduct(ctx, &productspb.UpdateProductRequest{
		Id:         variantID,
		Name:       name,
		BasePrice:  price,
		Stock:      stock,
		Attributes: pbAttrs,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r VariantRepository) Remove(ctx context.Context, productID string) (err error) {
	var conn *grpc.ClientConn
	conn, err = r.dial(ctx)
	if err != nil {
		return err
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	_, err = productspb.NewProductsServiceClient(conn).RemoveProduct(ctx, &productspb.RemoveProductRequest{Id: productID})
	if err != nil {
		return err
	}

	return nil
}

// Additional methods needed by variant tool service
func (r VariantRepository) CreateVariant(ctx context.Context, productID, name, sku string, price int64) (string, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// For now, delegate to existing Add method or implement with placeholder
	// TODO: Replace with actual CreateVariant RPC call when available
	return "variant_" + productID + "_" + sku, nil
}

func (r VariantRepository) UpdateVariant(ctx context.Context, variantID, name, sku string, price int64) error {
	// Delegate to existing Update method with appropriate parameters
	return r.Update(ctx, variantID, price, 0, name, []models.Attribute{})
}

func (r VariantRepository) DeleteVariant(ctx context.Context, variantID string) error {
	// Delegate to existing Remove method
	return r.Remove(ctx, variantID)
}

func (r VariantRepository) GetVariants(ctx context.Context, page, limit int64, sortBy, sortOrder string) ([]*models.Variant, int64, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close()

	// TODO: Implement with actual GetVariants RPC call when available
	// For now, return empty results
	return []*models.Variant{}, 0, nil
}

func (r VariantRepository) GetProductVariants(ctx context.Context, productID string) ([]*models.Variant, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// TODO: Implement with actual GetProductVariants RPC call when available
	// For now, return empty results
	return []*models.Variant{}, nil
}

func (r VariantRepository) SearchVariants(ctx context.Context, query string, page, limit int64) ([]*models.Variant, int64, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close()

	// TODO: Implement with actual SearchVariants RPC call when available
	// For now, return empty results
	return []*models.Variant{}, 0, nil
}

func (r VariantRepository) UpdateInventory(ctx context.Context, variantID string, quantity int) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// TODO: Implement with actual UpdateInventory RPC call when available
	// For now, return success
	return nil
}

func (r VariantRepository) GetInventory(ctx context.Context, variantID string) (int, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// TODO: Implement with actual GetInventory RPC call when available
	// For now, return placeholder inventory
	return 100, nil
}

func (r VariantRepository) variantToDomain(variant *productspb.Variant) *models.Variant {
	options := r.optionToDomain(variant)
	attributes := r.attributesToDomain(variant)
	return &models.Variant{
		VariantID:    variant.GetId(),
		ProductID:    variant.GetProductId(),
		Status:       variant.GetStatus(),
		SKU:          variant.GetSku(),
		Barcode:      variant.GetBarcode(),
		VariantPrice: variant.GetVariantPrice(),
		CurrencyCode: variant.GetCurrencyCode(),
		Stock:        variant.GetStock(),
		Weight:       variant.GetWeight(),
		Height:       variant.GetHeight(),
		Width:        variant.GetWidth(),
		Depth:        variant.GetDepth(),
		Attributes:   attributes,
		IsAvailable:  variant.GetIsAvailable(),
		HasOptions:   variant.GetHasOptions(),
		Options:      options,
	}
}

func (r VariantRepository) optionToDomain(product *productspb.Variant) []models.Option {

	// Convert []Option
	opts := make([]models.Option, len(product.GetOptions()))
	for i, o := range product.Options {
		opts[i] = models.Option{
			Name:  o.Name,
			Value: o.Value,
			Price: float64(o.Price),
		}
	}
	return opts
}

func (r VariantRepository) attributesToDomain(product *productspb.Variant) []models.Attribute {

	// Convert []Option
	opts := make([]models.Attribute, len(product.GetAttributes()))
	for i, o := range product.Attributes {
		opts[i] = models.Attribute{
			Key:   o.Key,
			Value: o.Value,
		}
	}
	return opts
}

func (r VariantRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r VariantRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
