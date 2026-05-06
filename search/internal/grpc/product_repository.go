// File: search/internal/grpc/product_repository.go
package grpc

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"middleman/internal/rpc"
	"middleman/products/productspb"
	"middleman/search/internal/application"
	"middleman/search/internal/models"
)

// ProductRepository calls the remote products service (gRPC) as a fallback.
type ProductRepository struct {
	endpoint string
}

var _ application.ProductRepository = (*ProductRepository)(nil)

// NewProductRepository instantiates the gRPC-based fallback repo.
func NewProductRepository(endpoint string) ProductRepository {
	return ProductRepository{
		endpoint: endpoint,
	}
}

// Find retrieves a product by ID from the products microservice (via gRPC).
func (r ProductRepository) Find(ctx context.Context, productID string) (*models.Product, error) {
	log.Printf("Find: retrieving product with ID=%s via gRPC fallback", productID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("Find: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	resp, err := client.GetProduct(ctx, &productspb.GetProductRequest{Id: productID})
	if err != nil {
		return nil, fmt.Errorf("GetProduct RPC failed: %w", err)
	}
	return r.productToDomain(resp.GetProduct()), nil
}
func (r ProductRepository) GetCatalog(ctx context.Context, userID string) ([]*models.Product, error) {
	log.Printf("Find: retrieving deal with ID=%s via gRPC fallback", userID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("Find: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	resp, err := client.GetCatalog(ctx, &productspb.GetCatalogRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("GetDeal RPC failed: %w", err)
	}

	var results []*models.Product
	for _, p := range resp.GetProducts() {
		domainProd := r.productToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// SearchWithFilters calls the remote microservice's GetProductsWithFilters (or similar).
func (r ProductRepository) SearchWithFilters(
	ctx context.Context,
	name string,
	categoryID string,
	categorySlug string,
	minPrice int64,
	maxPrice int64,
	brand string,
	condition string,
	model string,
	tags []string,
	manageStock bool,
	minStock int64,
	maxStock int64,
	sku string,
	status string,
	negotiable bool,
	userType string,
	middlemanService bool,
	hasVariants bool,
	shippingCost int64,
	minWeight int64,
	maxWeight int64,
	minHeight int64,
	maxHeight int64,
	minWidth int64,
	maxWidth int64,
	minDepth int64,
	maxDepth int64,
	offset int64,
	limit int64,
	lat, lng float64,
	radius int64,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*models.Product, error) {

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	// Pass both categoryID and categorySlug separately
	resp, err := client.GetProductsWithFilters(ctx, &productspb.GetProductsWithFiltersRequest{
		Name:             name,
		CategorySlug:     categorySlug,
		CategoryId:       categoryID,
		MinPrice:         minPrice,
		MaxPrice:         maxPrice,
		Brand:            brand,
		Condition:        condition,
		Model:            model,
		Tags:             tags,
		ManageStock:      manageStock,
		MinStock:         minStock,
		MaxStock:         maxStock,
		Sku:              sku,
		Status:           status,
		Negotiable:       negotiable,
		UserType:         userType,
		MiddlemanService: middlemanService,
		HasVariants:      hasVariants,
		ShippingCost:     shippingCost,
		MinWeight:        minWeight,
		MaxWeight:        maxWeight,
		MinHeight:        minHeight,
		MaxHeight:        maxHeight,
		MinWidth:         minWidth,
		MaxWidth:         maxWidth,
		MinDepth:         minDepth,
		MaxDepth:         maxDepth,
		Offset:           offset,
		Limit:            limit,
		Lat:              float32(lat),
		Lng:              float32(lng),
		Radius:           float32(radius),
		Page:             page,
		PageSize:         pageSize,
		SortBy:           sortBy,
		SortOrder:        sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetProductsWithFilters RPC failed: %w", err)
	}

	var results []*models.Product
	for _, p := range resp.GetProducts() {
		domainProd := r.productToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// SearchWithFilters calls the remote microservice's GetProductsWithFilters (or similar).
func (r ProductRepository) SearchProductsWithCategorySlug(
	ctx context.Context,
	categorySlug string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*models.Product, error) {

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	resp, err := client.GetProductsByCategorySlug(ctx, &productspb.GetProductsByCategorySlugRequest{
		CategorySlug: categorySlug,
		Page:         page,
		PageSize:     pageSize,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetProductsWithFilters RPC failed: %w", err)
	}

	var results []*models.Product
	for _, p := range resp.GetProducts() {
		domainProd := r.productToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// SearchWithFilters calls the remote microservice's GetProductsWithFilters (or similar).
func (r ProductRepository) SearchProductsWithCategory(
	ctx context.Context,
	categoryId string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*models.Product, error) {

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	resp, err := client.GetProductsByCategory(ctx, &productspb.GetProductsByCategoryRequest{
		CategoryId: categoryId,
		Page:       page,
		PageSize:   pageSize,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetProductsWithFilters RPC failed: %w", err)
	}

	var results []*models.Product
	for _, p := range resp.GetProducts() {
		domainProd := r.productToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// Add calls AddProduct in the remote gRPC service.
func (r ProductRepository) Add(ctx context.Context, productID string, name string, description string, basePrice int64, userSellerID string, categoryID string, brand string, condition string, model string, tags []string) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.AddProduct(ctx, &productspb.AddProductRequest{
		Name:        name,
		Description: description,
		BasePrice:   basePrice,
		CategoryId:  categoryID,
		Brand:       brand,
		Condition:   condition,
		Model:       model,
		Tags:        tags,
	})
	return err
}

// Update calls UpdateProduct in the remote gRPC service (or partial, etc.).
func (r ProductRepository) Update(ctx context.Context, productID string, price int64) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.UpdateProduct(ctx, &productspb.UpdateProductRequest{
		Id:        productID,
		BasePrice: price,
	})
	return err
}

// Remove calls RemoveProduct in the remote gRPC service.
func (r ProductRepository) Remove(ctx context.Context, productID string) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.RemoveProduct(ctx, &productspb.RemoveProductRequest{Id: productID})
	return err
}

// productToDomain converts a productspb.Product into our internal models.Product.
func (r ProductRepository) productToDomain(pb *productspb.Product) *models.Product {
	if pb == nil {
		return nil
	}
	attrs := make([]models.Attribute, len(pb.GetAttributes()))
	for i, a := range pb.GetAttributes() {
		attrs[i] = models.Attribute{
			Key:   a.GetKey(),
			Value: a.GetValue(),
		}
	}
	opts := make([]models.Option, len(pb.GetOptions()))
	for i, o := range pb.GetOptions() {
		opts[i] = models.Option{
			Name:  o.GetName(),
			Value: o.GetValue(),
			Price: float64(o.GetPrice()),
		}
	}

	return &models.Product{
		ProductID:        pb.GetId(),
		Name:             pb.GetName(),
		Description:      pb.GetDescription(),
		BasePrice:        pb.GetBasePrice(),
		UserSellerID:     pb.GetUserSellerId(),
		CategoryID:       pb.GetCategoryId(),
		CategorySlug:     pb.GetCategorySlug(),
		Brand:            pb.GetBrand(),
		Condition:        pb.GetCondition(),
		Model:            pb.GetModel(),
		Tags:             pb.GetTags(),
		ManageStock:      pb.GetManageStocks(),
		Stock:            pb.GetStock(),
		SKU:              pb.GetSku(),
		Attributes:       attrs,
		Weight:           pb.GetWeight(),
		Height:           pb.GetHeight(),
		Width:            pb.GetWidth(),
		Depth:            pb.GetDepth(),
		Status:           pb.GetStatus(),
		Negotiable:       pb.GetNegotiable(),
		MiddlemanService: pb.GetMiddlemanService(),
		UserType:         pb.GetUserType(),
		ShippingCost:     pb.GetShippingCost(),
		HasVariants:      pb.GetHasVariants(),
		Options:          opts,
		Lat:              float64(pb.GetLat()),
		Lng:              float64(pb.GetLng()),
		Thumbnail:        pb.GetThumbnail(),
		EntityType:       models.ProductType,
	}
}

// dial sets up a gRPC connection with the microservice endpoint.
func (r ProductRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
