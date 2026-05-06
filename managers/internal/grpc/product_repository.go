// File: search/internal/grpc/product_repository.go
package grpc

import (
	"context"
	"fmt"
	"middleman/managers/internal/domain"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/models"
	"middleman/products/productspb"
)

// ProductRepository calls the remote products service (gRPC) as a fallback.
type ProductRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.ProductRepository = (*ProductRepository)(nil)

// NewProductRepositoryWithAuth creates a new ProductRepository with JWT authentication support
func NewProductRepository(endpoint string, authInstance *auth.Auth) ProductRepository {
	return ProductRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

func (r ProductRepository) Rebrand(ctx context.Context, productID, name, description string, basePrice int64, stock int64, sku string, categoryID string, status string) error {
	log.Printf("Rebrand: rebranding product with ID=%s via gRPC", productID)

	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.RebrandProduct(ctx, &productspb.RebrandProductRequest{
		Id:          productID,
		Name:        name,
		Description: description,
		BasePrice:   basePrice,
		Stock:       stock,
		Sku:         sku,
		CategoryId:  categoryID,
		Status:      status,
	})
	if err != nil {
		return fmt.Errorf("RebrandProduct RPC failed: %w", err)
	}

	log.Printf("Rebrand: successfully rebranded product with ID=%s", productID)
	return nil
}

func (r ProductRepository) SearchWithTerm(ctx context.Context, name string) ([]*models.Product, error) {
	log.Printf("SearchWithTerm: searching products with term=%s via gRPC", name)

	// Create a timeout context with 20-second deadline
	timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	log.Printf("SearchWithTerm: created 20-second timeout context for term=%s", name)

	conn, err := r.dial(timeoutCtx)
	if err != nil {
		log.Printf("SearchWithTerm: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	log.Printf("SearchWithTerm: successfully connected to gRPC for term=%s", name)

	client := productspb.NewProductsServiceClient(conn)

	log.Printf("SearchWithTerm: calling GetProductsWithFilters RPC for term=%s", name)

	// Use GetProductsWithFilters with name filter
	resp, err := client.GetProductsWithFilters(timeoutCtx, &productspb.GetProductsWithFiltersRequest{
		Name:     name,
		Page:     1,
		PageSize: 20,
	})
	if err != nil {
		log.Printf("SearchWithTerm: GetProductsWithFilters RPC FAILED for term=%s: %v", name, err)
		return nil, fmt.Errorf("GetProductsWithFilters RPC failed: %w", err)
	}

	log.Printf("SearchWithTerm: GetProductsWithFilters RPC SUCCESS for term=%s, got %d products", name, len(resp.GetProducts()))

	products := make([]*models.Product, 0, len(resp.GetProducts()))
	for _, pbProduct := range resp.GetProducts() {
		products = append(products, r.productToDomain(pbProduct))
	}

	log.Printf("SearchWithTerm: returning %d products for term=%s", len(products), name)
	return products, nil
}

func (r ProductRepository) SuggestProducts(ctx context.Context, name string) ([]*models.Product, error) {
	// Use SearchWithTerm as a fallback for suggestions
	return r.SearchWithTerm(ctx, name)
}

func (r ProductRepository) UpdateThumbnail(ctx context.Context, productID string, thumbnail string) error {
	log.Printf("UpdateThumbnail: updating thumbnail for product ID=%s via gRPC", productID)

	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.UpdateProductThumbnail(ctx, &productspb.UpdateProductThumbnailRequest{
		ProductId: productID,
		Thumbnail: thumbnail,
	})
	if err != nil {
		return fmt.Errorf("UpdateProductThumbnail RPC failed: %w", err)
	}

	log.Printf("UpdateThumbnail: successfully updated thumbnail for product ID=%s", productID)
	return nil
}

// Find retrieves a product by ID from the products microservice (via gRPC).
func (r ProductRepository) Find(ctx context.Context, productID string) (*models.Product, error) {
	log.Printf("Find: retrieving product with ID=%s via gRPC fallback", productID)

	// Create a timeout context with 20-second deadline
	timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	log.Printf("Find: created 20-second timeout context for productID=%s", productID)

	conn, err := r.dial(timeoutCtx)
	if err != nil {
		log.Printf("Find: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	log.Printf("Find: successfully connected to gRPC for productID=%s", productID)

	client := productspb.NewProductsServiceClient(conn)

	log.Printf("Find: calling GetProduct RPC for productID=%s", productID)

	resp, err := client.GetProduct(timeoutCtx, &productspb.GetProductRequest{Id: productID})
	if err != nil {
		log.Printf("Find: GetProduct RPC FAILED for productID=%s: %v", productID, err)
		return nil, fmt.Errorf("GetProduct RPC failed: %w", err)
	}

	log.Printf("Find: GetProduct RPC SUCCESS for productID=%s", productID)
	return r.productToDomain(resp.GetProduct()), nil
}

// SearchWithFilters calls the remote microservice's GetProductsWithFilters (or similar).
func (r ProductRepository) SearchWithFilters(
	ctx context.Context,
	name string,
	categoryId string,
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
	log.Printf("[PRODUCT_GRPC] ========== SearchWithFilters START ==========")
	log.Printf("[PRODUCT_GRPC] SearchWithFilters called with: name='%s', categoryId='%s', status='%s', page=%d, pageSize=%d", name, categoryId, status, page, pageSize)
	log.Printf("[PRODUCT_GRPC] Repository endpoint configured as: %s", r.endpoint)
	log.Printf("[PRODUCT_GRPC] Auth configured: %v", r.auth != nil)

	// Fallback: if no meaningful filters provided, call GetProducts to list all products.
	noTextFilters := strings.TrimSpace(name) == "" && strings.TrimSpace(categoryId) == "" && len(tags) == 0 && strings.TrimSpace(brand) == "" && strings.TrimSpace(condition) == "" && strings.TrimSpace(model) == "" && strings.TrimSpace(status) == ""
	noPriceFilters := minPrice == 0 && maxPrice == 0
	noStockFilters := !manageStock && minStock == 0 && maxStock == 0 && strings.TrimSpace(sku) == ""
	if noTextFilters && noPriceFilters && noStockFilters {
		log.Printf("[PRODUCT_GRPC] No filters provided, using GetProducts fallback")
		log.Printf("[PRODUCT_GRPC] Attempting to dial products service at endpoint: %s", r.endpoint)

		// Create a timeout context with 20-second deadline
		timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()

		log.Printf("[PRODUCT_GRPC] Created 20-second timeout context for GetProducts")

		conn, err := r.dial(timeoutCtx)
		if err != nil {
			log.Printf("[PRODUCT_GRPC] FAILED to connect to products service: %v", err)
			log.Printf("[PRODUCT_GRPC] Error details: endpoint=%s, error_type=%T", r.endpoint, err)
			return nil, err
		}
		defer func() {
			if closeErr := conn.Close(); closeErr != nil {
				log.Printf("[PRODUCT_GRPC] Warning: Failed to close connection: %v", closeErr)
			} else {
				log.Printf("[PRODUCT_GRPC] Connection closed successfully")
			}
		}()

		log.Printf("[PRODUCT_GRPC] Successfully connected to products service at %s", r.endpoint)
		log.Printf("[PRODUCT_GRPC] Connection state: %v", conn.GetState())
		log.Printf("[PRODUCT_GRPC] Calling GetProducts RPC...")

		client := productspb.NewProductsServiceClient(conn)
		log.Printf("[PRODUCT_GRPC] ProductsServiceClient created")

		request := &productspb.GetProductsRequest{
			Page:      page,
			PageSize:  pageSize,
			SortBy:    sortBy,
			SortOrder: sortOrder,
		}
		log.Printf("[PRODUCT_GRPC] Sending GetProducts request: %+v", request)

		resp, err := client.GetProducts(timeoutCtx, request)
		if err != nil {
			log.Printf("[PRODUCT_GRPC] GetProducts RPC FAILED: %v", err)
			log.Printf("[PRODUCT_GRPC] Error details: type=%T, endpoint=%s", err, r.endpoint)
			return nil, fmt.Errorf("GetProducts RPC failed: %w", err)
		}

		log.Printf("[PRODUCT_GRPC] GetProducts RPC successful!")
		log.Printf("[PRODUCT_GRPC] Response details: received %d products", len(resp.GetProducts()))
		log.Printf("[PRODUCT_GRPC] Response total count: %d", resp.GetTotalCount())
		log.Printf("[PRODUCT_GRPC] Response current page: %d", resp.GetCurrentPage())

		var results []*models.Product
		for i, p := range resp.GetProducts() {
			log.Printf("[PRODUCT_GRPC] Processing product %d: ID=%s, Name=%s, Price=%d", i, p.GetId(), p.GetName(), p.GetBasePrice())
			if prod := r.productToDomain(p); prod != nil {
				log.Printf("[PRODUCT_GRPC] Successfully converted product %d to domain model", i)
				results = append(results, prod)
			} else {
				log.Printf("[PRODUCT_GRPC] WARNING: Failed to convert product %d to domain model", i)
			}
		}

		log.Printf("[PRODUCT_GRPC] GetProducts returning %d converted products out of %d received", len(results), len(resp.GetProducts()))
		log.Printf("[PRODUCT_GRPC] ========== SearchWithFilters END (GetProducts path) ==========")
		return results, nil
	}

	log.Printf("[PRODUCT_GRPC] Using filtered search with GetProductsWithFilters")
	log.Printf("[PRODUCT_GRPC] Attempting to dial products service at endpoint: %s", r.endpoint)

	// Create a timeout context with 20-second deadline
	timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	log.Printf("[PRODUCT_GRPC] Created 20-second timeout context for GetProductsWithFilters")

	conn, err := r.dial(timeoutCtx)
	if err != nil {
		log.Printf("[PRODUCT_GRPC] FAILED to connect to products service: %v", err)
		return nil, err
	}
	defer conn.Close()

	log.Printf("[PRODUCT_GRPC] Successfully connected to products service at %s", r.endpoint)
	log.Printf("[PRODUCT_GRPC] Calling GetProductsWithFilters RPC...")

	client := productspb.NewProductsServiceClient(conn)
	resp, err := client.GetProductsWithFilters(timeoutCtx, &productspb.GetProductsWithFiltersRequest{
		Name:             name,
		CategoryId:       categoryId,
		CategorySlug:     categorySlug,
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
		log.Printf("[PRODUCT_GRPC] GetProductsWithFilters RPC FAILED: %v", err)
		return nil, fmt.Errorf("GetProductsWithFilters RPC failed: %w", err)
	}

	log.Printf("[PRODUCT_GRPC] GetProductsWithFilters RPC successful, received %d products", len(resp.GetProducts()))

	var results []*models.Product
	for i, p := range resp.GetProducts() {
		log.Printf("[PRODUCT_GRPC] Converting product %d: ID=%s, Name=%s", i, p.GetId(), p.GetName())
		domainProd := r.productToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}

	log.Printf("[PRODUCT_GRPC] GetProductsWithFilters returning %d converted products", len(results))
	log.Printf("[PRODUCT_GRPC] ========== SearchWithFilters END (filtered path) ==========")
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
func (r ProductRepository) Add(ctx context.Context, productID string, name string, description string, basePrice int64, userSellerID string, categoryID string, categorySlug string, brand string, condition string, model string, tags []string, manageStock bool, stock int64, sku string, attributes []models.Attribute, weight int64, height int64, width int64, depth int64, status string, negotiable bool, userType string, middlemanService bool, shippingCost int64, hasVariants bool, options []models.Option, lat float64, lng float64, thumbnail string, entityType models.EntityType) error {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.AddProduct(ctx, &productspb.AddProductRequest{
		Name:         name,
		Description:  description,
		BasePrice:    basePrice,
		CategoryId:   categoryID,
		CategorySlug: categorySlug,
		Brand:        brand,
		Condition:    condition,
		Model:        model,
		Tags:         tags,
		Sku:          sku,
		Weight:       weight,
		Height:       height,
		Width:        width,
		Depth:        depth,
		Status:       status,
		Negotiable:   negotiable,
		UserType:     userType,
		ShippingCost: shippingCost,
		HasVariants:  hasVariants,
	})
	return err
}

// Update calls UpdateProduct in the remote gRPC service (or partial, etc.).
func (r ProductRepository) Update(ctx context.Context, productID string, price int64) error {
	conn, err := r.dialWithAuth(ctx)
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
	conn, err := r.dialWithAuth(ctx)
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
	}
}

// dial sets up a gRPC connection with the microservice endpoint
func (r ProductRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r ProductRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}

// Additional methods available in gRPC service

// GetCatalog retrieves user's catalog of products
func (r ProductRepository) GetCatalog(ctx context.Context, userID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Product, error) {
	log.Printf("[PRODUCT_GRPC] GetCatalog called for userID=%s", userID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	resp, err := client.GetCatalog(ctx, &productspb.GetCatalogRequest{
		UserId:    userID,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetCatalog RPC failed: %w", err)
	}

	var results []*models.Product
	for _, p := range resp.GetProducts() {
		if prod := r.productToDomain(p); prod != nil {
			results = append(results, prod)
		}
	}
	return results, nil
}

// GetPublicCatalog retrieves user's public catalog of products
func (r ProductRepository) GetPublicCatalog(ctx context.Context, userID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Product, error) {
	log.Printf("[PRODUCT_GRPC] GetPublicCatalog called for userID=%s", userID)

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	resp, err := client.GetPublicCatalog(ctx, &productspb.GetPublicCatalogRequest{
		UserId:    userID,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetPublicCatalog RPC failed: %w", err)
	}

	var results []*models.Product
	for _, p := range resp.GetProducts() {
		if prod := r.productToDomain(p); prod != nil {
			results = append(results, prod)
		}
	}
	return results, nil
}

// UpdateProductPrice updates product price
func (r ProductRepository) UpdateProductPrice(ctx context.Context, userID string, productID string, newPrice int64, oldPrice int64) error {
	log.Printf("[PRODUCT_GRPC] UpdateProductPrice called for productID=%s", productID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.UpdateProductPrice(ctx, &productspb.UpdateProductPriceRequest{

		Id:       productID,
		NewPrice: newPrice,
		OldPrice: oldPrice,
	})
	return err
}

// AdjustProductStock adjusts product stock
func (r ProductRepository) AdjustProductStock(ctx context.Context, userID string, productID string, newStock int64) error {
	log.Printf("[PRODUCT_GRPC] AdjustProductStock called for productID=%s", productID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.AdjustProductStock(ctx, &productspb.AdjustProductStockRequest{
		ProductId: productID,
		NewStock:  newStock,
	})
	return err
}

// ArchiveProduct archives a product
func (r ProductRepository) ArchiveProduct(ctx context.Context, userID string, productID string) error {
	log.Printf("[PRODUCT_GRPC] ArchiveProduct called for productID=%s", productID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.ArchiveProduct(ctx, &productspb.ArchiveProductRequest{
		ProductId: productID,
	})
	return err
}

// MarkProductSold marks a product as sold
func (r ProductRepository) MarkProductSold(ctx context.Context, userID string, productID string) error {
	log.Printf("[PRODUCT_GRPC] MarkProductSold called for productID=%s", productID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.MarkProductSold(ctx, &productspb.MarkProductSoldRequest{
		ProductId: productID,
	})
	return err
}

// MarkProductLeased marks a product as leased
func (r ProductRepository) MarkProductLeased(ctx context.Context, userID string, productID string, monthlyPrice int64, leaseTermMonths int64) error {
	log.Printf("[PRODUCT_GRPC] MarkProductLeased called for productID=%s", productID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.MarkProductLeased(ctx, &productspb.MarkProductLeasedRequest{

		ProductId:       productID,
		MonthlyPrice:    monthlyPrice,
		LeaseTermMonths: leaseTermMonths,
	})
	return err
}

// MarkProductPawned marks a product as pawned
func (r ProductRepository) MarkProductPawned(ctx context.Context, userID string, productID string, lockedPrice int64, redemptionFee int64) error {
	log.Printf("[PRODUCT_GRPC] MarkProductPawned called for productID=%s", productID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.MarkProductPawned(ctx, &productspb.MarkProductPawnedRequest{
		ProductId:     productID,
		LockedPrice:   lockedPrice,
		RedemptionFee: redemptionFee,
	})
	return err
}

// IncreaseProductPrice increases product price
func (r ProductRepository) IncreaseProductPrice(ctx context.Context, userID string, productID string, price int64) error {
	log.Printf("[PRODUCT_GRPC] IncreaseProductPrice called for productID=%s", productID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.IncreaseProductPrice(ctx, &productspb.IncreaseProductPriceRequest{

		ProductId: productID,
		Price:     price,
	})
	return err
}

// DecreaseProductPrice decreases product price
func (r ProductRepository) DecreaseProductPrice(ctx context.Context, userID string, productID string, newPrice int64) error {
	log.Printf("[PRODUCT_GRPC] DecreaseProductPrice called for productID=%s", productID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.DecreaseProductPrice(ctx, &productspb.DecreaseProductPriceRequest{

		ProductId: productID,
		NewPrice:  newPrice,
	})
	return err
}

// AddProductThumbnail adds a thumbnail to a product
func (r ProductRepository) AddProductThumbnail(ctx context.Context, productID string, thumbnail string) error {
	log.Printf("[PRODUCT_GRPC] AddProductThumbnail called for productID=%s", productID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := productspb.NewProductsServiceClient(conn)
	_, err = client.AddProductThumbnail(ctx, &productspb.AddProductThumbnailRequest{
		ProductId: productID,
		Thumbnail: thumbnail,
	})
	return err
}
