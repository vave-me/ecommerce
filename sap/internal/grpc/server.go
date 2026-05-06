package grpc

import (
	"context"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"middleman/internal/auth"
	"middleman/internal/errorsotel"
	"middleman/products/internal/application"
	"middleman/products/internal/application/commands"
	"middleman/products/internal/application/queries"
	"middleman/products/internal/domain"
	"middleman/products/productspb"
)

type server struct {
	app application.App
	productspb.UnimplementedProductsServiceServer
}

var _ productspb.ProductsServiceServer = (*server)(nil)

// RegisterServer registers the gRPC server implementation
func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	productspb.RegisterProductsServiceServer(registrar, server{app: app})
	return nil
}

// -----------------------------------------------------------------------------
// 1) PRODUCT METHODS
// -----------------------------------------------------------------------------
func (s server) AddProduct(ctx context.Context, req *productspb.AddProductRequest) (*productspb.AddProductResponse, error) {

	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	domainAttrs := make([]domain.Attribute, len(req.GetAttributes()))
	for i, attr := range req.GetAttributes() {
		domainAttrs[i] = domain.Attribute{
			Key:   attr.GetKey(),
			Value: attr.GetValue(),
		}
	}
	domainOptions := make([]domain.Option, len(req.GetOptions()))
	for i, opt := range req.GetOptions() {
		domainOptions[i] = domain.Option{
			Name:  opt.GetName(),
			Value: opt.GetValue(),
			Price: float64(opt.GetPrice()), // or int64 => domain int64 if you prefer
		}
	}

	cmd := commands.AddProduct{
		ID:               uuid.New().String(),
		Name:             req.GetName(),
		Description:      req.GetDescription(),
		BasePrice:        req.GetBasePrice(),
		UserSellerID:     userID,
		CategoryID:       req.GetCategoryId(),
		CategorySlug:     req.GetCategorySlug(),
		Brand:            req.GetBrand(),
		Condition:        domain.ToProductCondition(req.GetCondition()),
		Model:            req.GetModel(),
		Tags:             req.GetTags(), // if using CSV
		ManageStock:      req.GetManageStocks(),
		Stock:            req.GetStock(),
		SKU:              req.GetSku(),
		Attributes:       domainAttrs,
		Weight:           req.GetWeight(),
		Height:           req.GetHeight(),
		Width:            req.GetWidth(),
		Depth:            req.GetDepth(),
		Status:           domain.ToProductStatus(req.GetStatus()),
		Negotiable:       req.GetNegotiable(),
		UserType:         domain.ToUserType(req.GetUserType()),
		MiddlemanService: req.GetMiddlemanService(),
		ShippingCost:     req.GetShippingCost(),
		HasVariants:      req.GetHasVariants(),
		Options:          domainOptions,
		Thumbnail:        req.GetThumbnail(),
		Lat:              float64(req.GetLat()),
		Lng:              float64(req.GetLng()),
	}
	// Then call your application service:
	if err := s.app.AddProduct(ctx, cmd); err != nil {
		return nil, err
	}
	// Return the new ID in the response:
	return &productspb.AddProductResponse{Id: cmd.ID}, nil
}

func (s server) RebrandProduct(ctx context.Context, req *productspb.RebrandProductRequest) (*productspb.RebrandProductResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", req.GetId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	domainAttrs := make([]domain.Attribute, len(req.GetAttributes()))
	for i, attr := range req.GetAttributes() {
		domainAttrs[i] = domain.Attribute{
			Key:   attr.GetKey(),
			Value: attr.GetValue(),
		}
	}
	domainOptions := make([]domain.Option, len(req.GetOptions()))
	for i, opt := range req.GetOptions() {
		domainOptions[i] = domain.Option{
			Name:  opt.GetName(),
			Value: opt.GetValue(),
			Price: float64(opt.GetPrice()), // or int64 => domain int64 if you prefer
		}
	}

	err := s.app.RebrandProduct(ctx, commands.RebrandProduct{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &productspb.RebrandProductResponse{}, err
}

func (s server) UpdateProduct(ctx context.Context, req *productspb.UpdateProductRequest) (*productspb.UpdateProductResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", req.GetId()))

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	domainAttrs := make([]domain.Attribute, len(req.GetAttributes()))
	for i, attr := range req.GetAttributes() {
		domainAttrs[i] = domain.Attribute{
			Key:   attr.GetKey(),
			Value: attr.GetValue(),
		}
	}
	domainOptions := make([]domain.Option, len(req.GetOptions()))
	for i, opt := range req.GetOptions() {
		domainOptions[i] = domain.Option{
			Name:  opt.GetName(),
			Value: opt.GetValue(),
			Price: float64(opt.GetPrice()), // or int64 => domain int64 if you prefer
		}
	}

	err := s.app.UpdateProduct(ctx, commands.UpdateProduct{
		ID:               req.GetId(),
		Name:             req.GetName(),
		Description:      req.GetDescription(),
		BasePrice:        req.GetBasePrice(),
		UserSellerID:     userID,
		CategoryID:       req.GetCategoryId(),
		CategorySlug:     req.GetCategorySlug(),
		Brand:            req.GetBrand(),
		Condition:        domain.ToProductCondition(req.GetCondition()),
		Model:            req.GetModel(),
		Tags:             req.GetTags(), // if using CSV
		ManageStock:      req.GetManageStocks(),
		Stock:            req.GetStock(),
		Sku:              req.GetSku(),
		Attributes:       domainAttrs,
		Weight:           req.GetWeight(),
		Height:           req.GetHeight(),
		Width:            req.GetWidth(),
		Depth:            req.GetDepth(),
		Status:           domain.ToProductStatus(req.GetStatus()),
		Negotiable:       req.GetNegotiable(),
		UserType:         domain.ToUserType(req.GetUserType()),
		MiddlemanService: req.GetMiddlemanService(),
		ShippingCost:     req.GetShippingCost(),
		HasVariants:      req.GetHasVariants(),
		Options:          domainOptions,
		Thumbnail:        req.GetThumbnail(),
		Lat:              float64(req.GetLat()),
		Lng:              float64(req.GetLng()),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &productspb.UpdateProductResponse{}, err
}

// For Product Price Increase
func (s server) IncreaseProductPrice(ctx context.Context, request *productspb.IncreaseProductPriceRequest) (*productspb.IncreaseProductPriceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", request.GetProductId()))

	err := s.app.IncreaseProductPrice(ctx, commands.IncreaseProductPrice{
		ID:    request.GetProductId(),
		Price: request.GetPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Suppose your domain service can return oldPrice and newPrice. You could fill them here if needed:
	return &productspb.IncreaseProductPriceResponse{
		ProductId: request.GetProductId(),
		NewPrice:  request.GetPrice(),
	}, nil
}

// For Product Price Decrease
func (s server) DecreaseProductPrice(ctx context.Context, request *productspb.DecreaseProductPriceRequest) (*productspb.DecreaseProductPriceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", request.GetProductId()))

	err := s.app.DecreaseProductPrice(ctx, commands.DecreaseProductPrice{
		ID:    request.GetProductId(),
		Price: request.GetNewPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &productspb.DecreaseProductPriceResponse{
		ProductId: request.GetProductId(),
		NewPrice:  request.GetNewPrice(),
	}, nil
}

// For listing products in pages
func (s server) GetProducts(ctx context.Context, request *productspb.GetProductsRequest) (*productspb.GetProductsResponse, error) {

	span := trace.SpanFromContext(ctx)
	// 1) Guard for zero/negative PageSize
	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10 // <-- your default page size
	}

	// 2) Guard for zero/negative Page
	page := request.GetPage()
	if page <= 0 {
		page = 1 // <-- your default page number
	}

	products, totalCount, err := s.app.GetProducts(ctx, queries.GetProducts{
		Page:      request.GetPage(),
		PageSize:  request.GetPageSize(),
		SortBy:    request.GetSortBy(),
		SortOrder: request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	protoProducts := make([]*productspb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}

	return &productspb.GetProductsResponse{
		Products:    protoProducts,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

// For listing products in a category
func (s server) GetProductsByCategory(ctx context.Context, request *productspb.GetProductsByCategoryRequest) (*productspb.GetProductsByCategoryResponse, error) {
	span := trace.SpanFromContext(ctx)

	// 1) Guard for zero/negative PageSize
	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10 // <-- your default page size
	}

	// 2) Guard for zero/negative Page
	page := request.GetPage()
	if page <= 0 {
		page = 1 // <-- your default page number
	}

	products, totalCount, err := s.app.GetProductsByCategory(ctx, queries.GetProductsByCategory{
		CategoryID: request.CategoryId,
		Page:       request.GetPage(),
		PageSize:   request.GetPageSize(),
		SortBy:     request.GetSortBy(),
		SortOrder:  request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	protoProducts := make([]*productspb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}

	return &productspb.GetProductsByCategoryResponse{
		Products:    protoProducts,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

func (s server) GetProductsByCategorySlug(ctx context.Context, request *productspb.GetProductsByCategorySlugRequest) (*productspb.GetProductsByCategorySlugResponse, error) {
	span := trace.SpanFromContext(ctx)

	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}

	page := request.GetPage()
	if page <= 0 {
		page = 1
	}

	products, totalCount, err := s.app.GetProductsByCategorySlug(ctx, queries.GetProductsByCategorySlug{
		Slug:      request.GetCategorySlug(),
		Page:      request.GetPage(),
		PageSize:  request.GetPageSize(),
		SortBy:    request.GetSortBy(),
		SortOrder: request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	protoProducts := make([]*productspb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}

	return &productspb.GetProductsByCategorySlugResponse{
		Products:    protoProducts,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

func (s server) GetProductsWithFilters(ctx context.Context, request *productspb.GetProductsWithFiltersRequest) (*productspb.GetProductsWithFiltersResponse, error) {
	span := trace.SpanFromContext(ctx)

	// 1) Guard for zero/negative PageSize
	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}

	page := request.GetPage()
	if page <= 0 {
		page = 1
	}

	products, totalCount, err := s.app.GetProductsWithFilters(ctx, queries.GetProductsWithFilters{
		Name:             request.GetName(),
		Category:         request.GetCategory(),
		MinPrice:         request.GetMinPrice(),
		MaxPrice:         request.GetMaxPrice(),
		Brand:            request.GetBrand(),
		Condition:        request.GetCondition(),
		Model:            request.GetModel(),
		Tags:             request.GetTags(),
		ManageStock:      request.GetManageStock(),
		MinStock:         request.GetMinStock(),
		MaxStock:         request.GetMaxStock(),
		SKU:              request.GetSku(),
		Status:           request.GetStatus(),
		Negotiable:       request.GetNegotiable(),
		UserType:         request.GetUserType(),
		MiddlemanService: request.GetMiddlemanService(),
		HasVariants:      request.GetHasVariants(),
		ShippingCost:     request.GetShippingCost(),
		MinWeight:        request.GetMinWeight(),
		MaxWeight:        request.GetMaxWeight(),
		MinHeight:        request.GetMinHeight(),
		MaxHeight:        request.GetMaxHeight(),
		MinWidth:         request.GetMinWidth(),
		MaxWidth:         request.GetMaxWidth(),
		MinDepth:         request.GetMinDepth(),
		MaxDepth:         request.GetMaxDepth(),
		Offset:           request.GetOffset(),
		Limit:            request.GetLimit(),
		Lat:              float64(request.GetLat()),
		Lng:              float64(request.GetLng()),
		Radius:           int64(request.GetRadius()),
		Page:             request.GetPage(),
		PageSize:         request.GetPageSize(),
		SortBy:           request.GetSortBy(),
		SortOrder:        request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	protoProducts := make([]*productspb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}

	return &productspb.GetProductsWithFiltersResponse{
		Products:    protoProducts,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

// Get a single product
func (s server) GetProduct(ctx context.Context, request *productspb.GetProductRequest) (*productspb.GetProductResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", request.GetId()))

	product, err := s.app.GetProduct(ctx, queries.GetProduct{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &productspb.GetProductResponse{Product: s.productFromDomain(product)}, nil
}

// Return a user's product catalog
func (s server) GetCatalog(ctx context.Context, request *productspb.GetCatalogRequest) (*productspb.GetCatalogResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("UserID", request.GetUserId()))

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10 // <-- your default page size
	}

	page := request.GetPage()
	if page <= 0 {
		page = 1 // <-- your default page number
	}

	products, totalCount, err := s.app.GetCatalog(ctx, queries.GetCatalog{
		UserSellerID: request.GetUserId(),
		Page:         request.GetPage(),
		PageSize:     request.GetPageSize(),
		SortBy:       request.GetSortBy(),
		SortOrder:    request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	protoProducts := make([]*productspb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}

	return &productspb.GetCatalogResponse{
		Products:    protoProducts,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

func (s server) GetPublicCatalog(ctx context.Context, request *productspb.GetPublicCatalogRequest) (*productspb.GetPublicCatalogResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("UserID", request.GetUserId()))

	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10 // <-- your default page size
	}

	page := request.GetPage()
	if page <= 0 {
		page = 1 // <-- your default page number
	}

	products, totalCount, err := s.app.GetPublicCatalog(ctx, queries.GetPublicCatalog{
		UserSellerID: request.GetUserId(),
		Page:         request.GetPage(),
		PageSize:     request.GetPageSize(),
		SortBy:       request.GetSortBy(),
		SortOrder:    request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	totalPages := (totalCount + pageSize - 1) / pageSize

	protoProducts := make([]*productspb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}

	return &productspb.GetPublicCatalogResponse{
		Products:    protoProducts,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

// Remove a product from listing
func (s server) RemoveProduct(ctx context.Context, request *productspb.RemoveProductRequest) (*productspb.RemoveProductResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", request.GetId()))

	err := s.app.RemoveProduct(ctx, commands.RemoveProduct{
		ID:     request.GetId(),
		UserID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &productspb.RemoveProductResponse{}, err
}

// -----------------------------------------------------------------------------
// MISSING PRODUCT METHODS (some commented out for brevity in your snippet)
// -----------------------------------------------------------------------------

// (Fully uncomment & implement if you need the UpdateProduct method from the snippet)

// 2) AdjustProductStock
func (s server) AdjustProductStock(ctx context.Context, request *productspb.AdjustProductStockRequest) (*productspb.AdjustProductStockResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", request.GetProductId()))

	err := s.app.AdjustProductStock(ctx, commands.AdjustProductStock{
		ID:       request.GetProductId(),
		NewStock: request.GetNewStock(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// If your domain returns oldStock, you can populate it here. We'll omit that for now:
	return &productspb.AdjustProductStockResponse{
		ProductId: request.GetProductId(),
		NewStock:  request.GetNewStock(),
	}, nil
}

// 3) ArchiveProduct
func (s server) ArchiveProduct(ctx context.Context, request *productspb.ArchiveProductRequest) (*productspb.ArchiveProductResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", request.GetProductId()))

	err := s.app.ArchiveProduct(ctx, commands.ArchiveProduct{
		ID:           request.GetProductId(),
		UserSellerID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &productspb.ArchiveProductResponse{
		ProductId: request.GetProductId(),
		Archived:  true,
	}, nil
}

func (s server) AddProductThumbnail(ctx context.Context, request *productspb.AddProductThumbnailRequest) (*productspb.AddProductThumbnailResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", request.GetProductId()))

	err := s.app.AddProductThumbnail(ctx, commands.AddProductThumbnail{
		ID:        request.GetProductId(),
		Thumbnail: request.GetThumbnail(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &productspb.AddProductThumbnailResponse{}, nil
}
func (s server) UpdateProductThumbnail(ctx context.Context, request *productspb.UpdateProductThumbnailRequest) (*productspb.UpdateProductThumbnailResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", request.GetProductId()))

	err := s.app.UpdateProductThumbnail(ctx, commands.UpdateProductThumbnail{
		ID:        request.GetProductId(),
		Thumbnail: request.GetThumbnail(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &productspb.UpdateProductThumbnailResponse{}, nil
}

// 4) MarkProductSold
func (s server) MarkProductSold(ctx context.Context, request *productspb.MarkProductSoldRequest) (*productspb.MarkProductSoldResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", request.GetProductId()))

	err := s.app.MarkProductSold(ctx, commands.MarkProductSold{
		ID:           request.GetProductId(),
		UserSellerID: userID,
		// If finalPrice is in the proto, set it here: FinalPrice: request.GetFinalPrice()
		FinalPrice: 0,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &productspb.MarkProductSoldResponse{
		ProductId: request.GetProductId(),
		Status:    "sold",
	}, nil
}

// 5) MarkProductLeased
func (s server) MarkProductLeased(ctx context.Context, request *productspb.MarkProductLeasedRequest) (*productspb.MarkProductLeasedResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", request.GetProductId()))

	err := s.app.MarkProductLeased(ctx, commands.MarkProductLeased{
		ID:              request.GetProductId(),
		UserSellerID:    userID,
		MonthlyPrice:    request.GetMonthlyPrice(),
		LeaseTermMonths: request.GetLeaseTermMonths(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &productspb.MarkProductLeasedResponse{
		ProductId: request.GetProductId(),
		Status:    "leased",
	}, nil
}

// 6) MarkProductPawned
func (s server) MarkProductPawned(ctx context.Context, request *productspb.MarkProductPawnedRequest) (*productspb.MarkProductPawnedResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ProductID", request.GetProductId()))

	err := s.app.MarkProductPawned(ctx, commands.MarkProductPawned{
		ID:            request.GetProductId(),
		UserSellerID:  userID,
		LockedPrice:   request.GetLockedPrice(),
		RedemptionFee: request.GetRedemptionFee(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &productspb.MarkProductPawnedResponse{
		ProductId: request.GetProductId(),
		Status:    "pawned",
	}, nil
}

// -----------------------------------------------------------------------------
// 2) VARIANT METHODS
// -----------------------------------------------------------------------------
func (s server) AddVariant(ctx context.Context, req *productspb.AddVariantRequest) (*productspb.AddVariantResponse, error) {
	varID := uuid.New().String()

	// Convert repeated Attribute to domain
	domainAttrs := make([]domain.Attribute, len(req.GetAttributes()))
	for i, attr := range req.GetAttributes() {
		domainAttrs[i] = domain.Attribute{
			Key:   attr.GetKey(),
			Value: attr.GetValue(),
		}
	}
	// Convert repeated Option
	domainOpts := make([]domain.Option, len(req.GetOptions()))
	for i, opt := range req.GetOptions() {
		domainOpts[i] = domain.Option{
			Name:  opt.GetName(),
			Value: opt.GetValue(),
			Price: float64(opt.GetPrice()),
		}
	}

	cmd := commands.AddVariant{
		ID:           varID,
		ProductID:    req.GetProductId(),
		SKU:          req.GetSku(),
		Barcode:      req.GetBarcode(),
		VariantPrice: req.GetVariantPrice(),
		CurrencyCode: req.GetCurrencyCode(),
		Stock:        req.GetStock(),
		Weight:       req.GetWeight(),
		Height:       req.GetHeight(),
		Width:        req.GetWidth(),
		Depth:        req.GetDepth(),
		Attributes:   domainAttrs,
		IsAvailable:  req.GetIsAvailable(),
		HasOptions:   req.GetHasOptions(),
		Options:      domainOpts,
	}

	if err := s.app.AddVariant(ctx, cmd); err != nil {
		return nil, err
	}
	return &productspb.AddVariantResponse{VariantId: varID}, nil
}

func (s server) IncreaseVariantPrice(ctx context.Context, request *productspb.IncreaseVariantPriceRequest) (*productspb.IncreaseVariantPriceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("VariantID", request.GetVariantId()))

	err := s.app.IncreaseVariantPrice(ctx, commands.IncreaseVariantPrice{
		ID:    request.GetVariantId(),
		Price: request.GetPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// If your domain returns old/new price, set them here:
	return &productspb.IncreaseVariantPriceResponse{
		VariantId: request.GetVariantId(),
		NewPrice:  request.GetPrice(),
	}, nil
}

func (s server) DecreaseVariantPrice(ctx context.Context, request *productspb.DecreaseVariantPriceRequest) (*productspb.DecreaseVariantPriceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("VariantID", request.GetVariantId()))

	err := s.app.DecreaseVariantPrice(ctx, commands.DecreaseVariantPrice{
		ID:    request.GetVariantId(),
		Price: request.GetPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &productspb.DecreaseVariantPriceResponse{
		VariantId: request.GetVariantId(),
		NewPrice:  request.GetPrice(),
	}, nil
}

func (s server) AdjustVariantStock(ctx context.Context, request *productspb.AdjustVariantStockRequest) (*productspb.AdjustVariantStockResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("VariantID", request.GetVariantId()))

	err := s.app.AdjustVariantStock(ctx, commands.AdjustVariantStock{
		ID:       request.GetVariantId(),
		NewStock: request.GetNewStock(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &productspb.AdjustVariantStockResponse{
		VariantId: request.GetVariantId(),
		NewStock:  request.GetNewStock(),
	}, nil
}

func (s server) ArchiveVariant(ctx context.Context, request *productspb.ArchiveVariantRequest) (*productspb.ArchiveVariantResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("VariantID", request.GetVariantId()))

	err := s.app.ArchiveVariant(ctx, commands.ArchiveVariant{
		ID: request.GetVariantId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &productspb.ArchiveVariantResponse{
		VariantId: request.GetVariantId(),
		Archived:  true,
	}, nil
}

func (s server) RemoveVariant(ctx context.Context, request *productspb.RemoveVariantRequest) (*productspb.RemoveVariantResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("VariantID", request.GetVariantId()))

	err := s.app.RemoveVariant(ctx, commands.RemoveVariant{
		ID: request.GetVariantId(),
		// reason if needed
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &productspb.RemoveVariantResponse{
		VariantId: request.GetVariantId(),
	}, nil
}

func (s server) productFromDomain(p *domain.CatalogProduct) *productspb.Product {

	// (2) Convert domain Attributes => repeated productspb.Attribute
	pbAttributes := make([]*productspb.Attribute, len(p.Attributes))
	for i, attr := range p.Attributes {
		pbAttributes[i] = &productspb.Attribute{
			Key:   attr.Key,
			Value: attr.Value,
		}
	}

	// (3) Convert domain Options => repeated productspb.Option
	pbOptions := make([]*productspb.Option, len(p.Options))
	for i, opt := range p.Options {
		pbOptions[i] = &productspb.Option{
			Name:  opt.Name,
			Value: opt.Value,
			Price: int64(opt.Price),
		}
	}

	return &productspb.Product{
		Id:               p.ID,
		Name:             p.Name,
		Description:      p.Description,
		BasePrice:        p.BasePrice,
		UserSellerId:     p.UserSellerID,
		CategoryId:       p.CategoryID,
		CategorySlug:     p.CategorySlug,
		Brand:            p.Brand,
		Condition:        p.Condition.String(), // if Condition is a custom type
		Model:            p.Model,
		Tags:             p.Tags,
		ManageStocks:     p.ManageStock,
		Stock:            p.Stock,
		Sku:              p.SKU,
		Attributes:       pbAttributes,
		Weight:           p.Weight,
		Height:           p.Height,
		Width:            p.Width,
		Depth:            p.Depth,
		Status:           p.Status.String(),
		Negotiable:       p.Negotiable,
		MiddlemanService: p.MiddlemanService,
		UserType:         p.UserType.String(),
		ShippingCost:     p.ShippingCost,
		HasVariants:      p.HasVariants,
		Options:          pbOptions,
		Thumbnail:        p.Thumbnail,
		Lat:              float32(p.Lat),
		Lng:              float32(p.Lng),
	}
}

func (s server) variantFromDomain(v *domain.Variant) *productspb.Variant {
	// Convert []Attribute
	pbAttrs := make([]*productspb.Attribute, len(v.Attributes))
	for i, a := range v.Attributes {
		pbAttrs[i] = &productspb.Attribute{
			Key:   a.Key,
			Value: a.Value,
		}
	}
	// Convert []Option
	pbOpts := make([]*productspb.Option, len(v.Options))
	for i, o := range v.Options {
		pbOpts[i] = &productspb.Option{
			Name:  o.Name,
			Value: o.Value,
			Price: int64(o.Price),
		}
	}

	return &productspb.Variant{
		Id:           v.ID(),
		ProductId:    v.ProductID,
		Sku:          v.SKU,
		Barcode:      v.Barcode,
		VariantPrice: v.VariantPrice,
		CurrencyCode: v.CurrencyCode,
		Stock:        v.Stock,
		Weight:       v.Weight,
		Height:       v.Height,
		Width:        v.Width,
		Depth:        v.Depth,
		Attributes:   pbAttrs,
		IsAvailable:  v.IsAvailable,
		HasOptions:   v.HasOptions,
		Options:      pbOpts,
	}
}
