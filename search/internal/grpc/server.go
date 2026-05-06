// File: search/internal/grpc/server.go
package grpc

import (
	"context"
	"fmt"
	"log"
	"middleman/internal/errorsotel"
	"middleman/search/internal/application"
	"middleman/search/internal/models"
	"middleman/search/searchpb"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	app application.Application
	searchpb.UnimplementedSearchServiceServer
}

func RegisterServer(
	ctx context.Context,
	app application.Application,
	registrar grpc.ServiceRegistrar,
) error {
	searchpb.RegisterSearchServiceServer(registrar, server{app: app})
	return nil
}

func handlePanic(span trace.Span, methodName string) {
	if r := recover(); r != nil {
		err := fmt.Errorf("panic recovered in %s: %v", methodName, r)
		if span != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, "panic")
		}
		log.Printf("Panic recovered in %s: %v", methodName, r)
	}
}

// --- Order Methods ---

//func (s server) SearchOrders(ctx context.Context, req *searchpb.SearchOrdersRequest) (*searchpb.SearchOrdersResponse, error) {
//	span := trace.SpanFromContext(ctx)
//	defer handlePanic(span, "SearchOrders")
//
//	params := application.Search{ // Assumed application struct
//		UserCustomerID: req.GetUserCustomerId(),
//		After:          req.GetAfter().AsTime(),
//		Before:         req.GetBefore().AsTime(),
//		UserSellerIDs:  req.GetUserSellerIds(),
//		ProductIDs:     req.GetProductIds(),
//		MinTotal:       req.GetMinTotal(),
//		MaxTotal:       req.GetMaxTotal(),
//		Status:         req.GetStatus(),
//		Next:           req.GetNext(),
//		Limit:          req.GetLimit(),
//	}
//
//	orders, next, err := s.app.SearchOrders(ctx, params) // Assumed app method signature
//	if err != nil {
//		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
//		span.SetStatus(codes.Error, err.Error())
//		return nil, err
//	}
//
//	protoOrders := make([]*searchpb.Order, len(orders))
//	for i, order := range orders {
//		protoOrders[i] = s.orderFromDomain(order)
//	}
//
//	return &searchpb.SearchOrdersResponse{
//		Orders: protoOrders,
//		Next:   next,
//	}, nil
//}

func (s server) GetOrder(ctx context.Context, req *searchpb.GetOrderRequest) (*searchpb.GetOrderResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "GetOrder")

	span.SetAttributes(attribute.String("OrderID", req.GetId()))
	log.Printf("[SearchService] GetOrder: id=%s", req.GetId())

	order, err := s.app.GetOrder(ctx, application.GetOrder{
		OrderID: req.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &searchpb.GetOrderResponse{
		Order: s.orderFromDomain(order),
	}, nil
}

// --- User Methods ---

//func (s server) SearchUsers(ctx context.Context, req *searchpb.SearchUsersRequest) (*searchpb.SearchUsersResponse, error) {
//	span := trace.SpanFromContext(ctx)
//	defer handlePanic(span, "SearchUsers")
//
//	params := application.SearchUsers{ // Assumed application struct
//		Name: req.GetName(),
//	}
//
//	users, next, err := s.app.SearchUser(ctx, params) // Assumed app method signature
//	if err != nil {
//		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
//		span.SetStatus(codes.Error, err.Error())
//		return nil, err
//	}
//
//	protoUsers := make([]*searchpb.User, len(users))
//	for i, user := range users {
//		protoUsers[i] = s.userFromDomain(user)
//	}
//
//	return &searchpb.SearchUsersResponse{
//		Users: protoUsers,
//		Next:  next,
//	}, nil
//}

func (s server) GetUser(ctx context.Context, req *searchpb.GetUserRequest) (*searchpb.GetUserResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "GetUser")

	params := application.GetUser{ // Assumed application struct
		UserID: req.GetId(),
	}

	user, err := s.app.GetUser(ctx, params) // Assumed app method signature
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &searchpb.GetUserResponse{
		User: s.userFromDomain(user),
	}, nil
}

// --- Product Methods ---

func (s server) SearchProductsWithFilters(
	ctx context.Context,
	req *searchpb.SearchProductsWithFiltersRequest,
) (*searchpb.SearchProductsWithFiltersResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchProductsWithFilters")

	log.Printf("[SearchService] SearchProductsWithFilters: name=%s category=%s min=%d max=%d page=%d pageSize=%d",
		req.GetName(), req.GetCategoryId(), req.GetMinPrice(), req.GetMaxPrice(), req.GetPage(), req.GetPageSize())

	searchParams := application.SearchProductsWithFilters{
		Name:             req.GetName(),
		CategoryID:       req.GetCategoryId(),
		CategorySlug:     req.GetCategorySlug(),
		MinPrice:         req.GetMinPrice(),
		MaxPrice:         req.GetMaxPrice(),
		Brand:            req.GetBrand(),
		Condition:        req.GetCondition(),
		Model:            req.GetModel(),
		Tags:             req.GetTags(),
		ManageStock:      req.GetManageStock(),
		MinStock:         req.GetMinStock(),
		MaxStock:         req.GetMaxStock(),
		SKU:              req.GetSku(),
		Status:           req.GetStatus(),
		Negotiable:       req.GetNegotiable(),
		UserType:         req.GetUserType(),
		MiddlemanService: req.GetMiddlemanService(),
		HasVariants:      req.GetHasVariants(),
		ShippingCost:     req.GetShippingCost(),
		MinWeight:        req.GetMinWeight(),
		MaxWeight:        req.GetMaxWeight(),
		MinHeight:        req.GetMinHeight(),
		MaxHeight:        req.GetMaxHeight(),
		MinWidth:         req.GetMinWidth(),
		MaxWidth:         req.GetMaxWidth(),
		MinDepth:         req.GetMinDepth(),
		MaxDepth:         req.GetMaxDepth(),
		Offset:           req.GetOffset(),
		Limit:            req.GetLimit(),
		Lat:              float64(req.GetLat()),
		Lng:              float64(req.GetLng()),
		Radius:           int64(req.GetRadius()),
		Page:             req.GetPage(),
		PageSize:         req.GetPageSize(),
		SortBy:           req.GetSortBy(),
		SortOrder:        req.GetSortOrder(),
	}

	products, err := s.app.SearchProductsWithFilters(ctx, searchParams)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoProducts := make([]*searchpb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}
	return &searchpb.SearchProductsWithFiltersResponse{
		Products:    protoProducts,
		CurrentPage: req.GetPage(),
	}, nil
}
func (s server) SearchProductsWithCategorySlug(
	ctx context.Context,
	req *searchpb.SearchProductsWithCategorySlugRequest,
) (*searchpb.SearchProductsWithCategorySlugResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchProductsWithCategorySlug")

	categorySlug := req.GetCategorySlug()
	offset := req.GetOffset()
	if offset < 0 {
		offset = 0
	}
	limit := req.GetLimit()
	if limit <= 0 {
		limit = 20
	}
	lat := req.GetLat()
	lng := req.GetLng()
	radius := req.GetRadius()
	if radius < 0 {
		radius = 0
	}
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	sortBy := req.GetSortBy()
	sortOrder := req.GetSortOrder()
	if sortOrder == "" {
		sortOrder = "ASC"
	}

	searchParams := application.SearchProductsWithCategorySlug{
		CategorySlug: categorySlug,
		Offset:       offset,
		Limit:        limit,
		Lat:          float64(lat),
		Lng:          float64(lng),
		Radius:       int64(radius),
		Page:         page,
		PageSize:     pageSize,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	}

	products, err := s.app.SearchProductsWithCategorySlug(ctx, searchParams)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoProducts := make([]*searchpb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}

	return &searchpb.SearchProductsWithCategorySlugResponse{
		Products:    protoProducts,
		CurrentPage: page,
	}, nil
}

func (s server) SearchProductsWithCategory(
	ctx context.Context,
	req *searchpb.SearchProductsWithCategoryRequest,
) (*searchpb.SearchProductsWithCategoryResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchProductsWithCategory")

	categoryID := req.GetCategoryId()
	offset := req.GetOffset()
	if offset < 0 {
		offset = 0
	}
	limit := req.GetLimit()
	if limit <= 0 {
		limit = 20
	}
	lat := req.GetLat()
	lng := req.GetLng()
	radius := req.GetRadius()
	if radius < 0 {
		radius = 0
	}
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	sortBy := req.GetSortBy()
	sortOrder := req.GetSortOrder()
	if sortOrder == "" {
		sortOrder = "ASC"
	}

	searchParams := application.SearchProductsWithCategory{
		CategoryID: categoryID,
		Offset:     offset,
		Limit:      limit,
		Lat:        float64(lat),
		Lng:        float64(lng),
		Radius:     int64(radius),
		Page:       page,
		PageSize:   pageSize,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	products, err := s.app.SearchProductsWithCategory(ctx, searchParams)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoProducts := make([]*searchpb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}

	return &searchpb.SearchProductsWithCategoryResponse{
		Products:    protoProducts,
		CurrentPage: page,
	}, nil
}

func (s server) SearchProductsWithTerm(ctx context.Context, req *searchpb.SearchProductsWithTermRequest) (*searchpb.SearchProductsWithTermResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchProductsWithTerm")

	log.Printf("[SearchService] SearchProductsWithTerm: name=%s", req.GetName())

	products, err := s.app.SearchProductsWithTerm(ctx, application.SearchProductsWithTerm{
		Name: req.GetName(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoProducts := make([]*searchpb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = s.productFromDomain(product)
	}
	return &searchpb.SearchProductsWithTermResponse{
		Products: protoProducts,
	}, nil
}
func (s server) GetProduct(ctx context.Context, req *searchpb.GetProductRequest) (*searchpb.GetProductResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "GetProduct")

	span.SetAttributes(attribute.String("ProductID", req.GetId()))
	log.Printf("[SearchService] GetProduct: id=%s", req.GetId())

	product, err := s.app.GetProduct(ctx, application.GetProduct{
		ProductID: req.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &searchpb.GetProductResponse{
		Product: s.productFromDomain(product),
	}, nil
}

func (s server) SuggestProducts(ctx context.Context, req *searchpb.SuggestProductsRequest) (*searchpb.SuggestProductsResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SuggestProducts")

	log.Printf("[SearchService] SuggestProducts: name=%s", req.GetName())

	suggestions, err := s.app.SuggestProducts(ctx, application.SuggestProducts{
		Name: req.GetName(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoProducts := make([]*searchpb.Product, len(suggestions))
	for i, product := range suggestions {
		protoProducts[i] = s.productFromDomain(product)
	}
	return &searchpb.SuggestProductsResponse{
		Suggestions: protoProducts,
	}, nil
}

// --- Post Methods ---

func (s server) SearchPostsWithFilters(
	ctx context.Context,
	req *searchpb.SearchPostsWithFiltersRequest,
) (*searchpb.SearchPostsWithFiltersResponse, error) {

	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchPostsWithFilters")

	log.Printf("[SearchService] SearchPostsWithFilters: name=%s page=%d page_size=%d",
		req.GetName(), req.GetPage(), req.GetPageSize())

	searchParams := application.SearchPostsWithFilters{
		Name:         req.GetName(),
		Description:  req.GetDescription(),
		PostType:     req.GetPostType(),
		UserType:     req.GetUserType(),
		CategoryID:   req.GetCategoryId(),
		CategorySlug: req.GetCategorySlug(),
		Tags:         req.GetTags(),
		Status:       req.GetStatus(),
		Thumbnail:    req.GetThumbnail(),
		Lat:          float64(req.GetLat()),
		Lng:          float64(req.GetLng()),
		Radius:       int64(req.GetRadius()),
		Page:         req.GetPage(),
		PageSize:     req.GetPageSize(),
		SortBy:       req.GetSortBy(),
		SortOrder:    req.GetSortOrder(),
	}

	posts, err := s.app.SearchPostsWithFilters(ctx, searchParams)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoPosts := make([]*searchpb.Post, len(posts))
	for i, p := range posts {
		protoPosts[i] = s.postFromDomain(p)
	}

	return &searchpb.SearchPostsWithFiltersResponse{
		Posts:       protoPosts,
		CurrentPage: req.GetPage(),
	}, nil
}

func (s server) SearchPostsWithTerm(
	ctx context.Context,
	req *searchpb.SearchPostsWithTermRequest,
) (*searchpb.SearchPostsWithTermResponse, error) {

	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchPostsWithTerm")

	log.Printf("[SearchService] SearchPostsWithTerm: name=%s", req.GetName())

	posts, err := s.app.SearchPostsWithTerm(ctx, application.SearchPostsWithTerm{
		Name: req.GetName(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoPosts := make([]*searchpb.Post, len(posts))
	for i, p := range posts {
		protoPosts[i] = s.postFromDomain(p)
	}

	return &searchpb.SearchPostsWithTermResponse{
		Posts: protoPosts,
	}, nil
}

func (s server) SuggestPosts(
	ctx context.Context,
	req *searchpb.SuggestPostsRequest,
) (*searchpb.SuggestPostsResponse, error) {

	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SuggestPosts")

	log.Printf("[SearchService] SuggestPosts: name=%s", req.GetName())

	suggestions, err := s.app.SuggestPosts(ctx, application.SuggestPosts{
		Name: req.GetName(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoPosts := make([]*searchpb.Post, len(suggestions))
	for i, p := range suggestions {
		protoPosts[i] = s.postFromDomain(p)
	}
	return &searchpb.SuggestPostsResponse{
		Suggestions: protoPosts,
	}, nil
}

func (s server) GetPost(
	ctx context.Context,
	req *searchpb.GetPostRequest,
) (*searchpb.GetPostResponse, error) {

	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "GetPost")

	span.SetAttributes(attribute.String("PostID", req.GetId()))
	log.Printf("[SearchService] GetPost: id=%s", req.GetId())

	post, err := s.app.GetPost(ctx, application.GetPost{
		PostID: req.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &searchpb.GetPostResponse{
		Post: s.postFromDomain(post),
	}, nil
}

func (s server) SearchPostsWithCategorySlug(
	ctx context.Context,
	req *searchpb.SearchPostsWithCategorySlugRequest,
) (*searchpb.SearchPostsWithCategorySlugResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchPostsWithCategorySlug")
	// Assumed Implementation:
	categorySlug := req.GetCategorySlug()
	offset := req.GetOffset()
	if offset < 0 {
		offset = 0
	}
	limit := req.GetLimit()
	if limit <= 0 {
		limit = 20
	}
	lat := req.GetLat()
	lng := req.GetLng()
	radius := req.GetRadius()
	if radius < 0 {
		radius = 0
	}
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	sortBy := req.GetSortBy()
	sortOrder := req.GetSortOrder()
	if sortOrder == "" {
		sortOrder = "ASC"
	}

	params := application.SearchPostsWithCategorySlug{ // Assumed struct
		CategorySlug: categorySlug,
		Offset:       offset,
		Limit:        limit,
		Lat:          float64(lat),
		Lng:          float64(lng),
		Radius:       int64(radius),
		Page:         page,
		PageSize:     pageSize,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	}
	posts, err := s.app.SearchPostsWithCategorySlug(ctx, params) // Assumed method
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	protoPosts := make([]*searchpb.Post, len(posts))
	for i, p := range posts {
		protoPosts[i] = s.postFromDomain(p)
	}
	return &searchpb.SearchPostsWithCategorySlugResponse{Posts: protoPosts, CurrentPage: page}, nil
}

func (s server) SearchPostsWithCategory(
	ctx context.Context,
	req *searchpb.SearchPostsWithCategoryRequest,
) (*searchpb.SearchPostsWithCategoryResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchPostsWithCategory")
	// Assumed Implementation:
	categoryID := req.GetCategoryId()
	offset := req.GetOffset()
	if offset < 0 {
		offset = 0
	}
	limit := req.GetLimit()
	if limit <= 0 {
		limit = 20
	}
	lat := req.GetLat()
	lng := req.GetLng()
	radius := req.GetRadius()
	if radius < 0 {
		radius = 0
	}
	page := req.GetPage()
	if page <= 0 {
		page = 1
	}
	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20
	}
	sortBy := req.GetSortBy()
	sortOrder := req.GetSortOrder()
	if sortOrder == "" {
		sortOrder = "ASC"
	}

	params := application.SearchPostsWithCategory{ // Assumed struct
		CategoryID: categoryID,
		Offset:     offset,
		Limit:      limit,
		Lat:        float64(lat),
		Lng:        float64(lng),
		Radius:     int64(radius),
		Page:       page,
		PageSize:   pageSize,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}
	posts, err := s.app.SearchPostsWithCategory(ctx, params) // Assumed method
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	protoPosts := make([]*searchpb.Post, len(posts))
	for i, p := range posts {
		protoPosts[i] = s.postFromDomain(p)
	}
	return &searchpb.SearchPostsWithCategoryResponse{Posts: protoPosts, CurrentPage: page}, nil
}

// --- Service Methods (Stubs - Assuming Application methods exist) ---

func (s server) GetService(ctx context.Context, req *searchpb.GetServiceRequest) (*searchpb.GetServiceResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "GetService")
	// Assumed Implementation:
	params := application.GetService{ServiceID: req.GetId()} // Assumed struct
	service, err := s.app.GetService(ctx, params)            // Assumed method
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &searchpb.GetServiceResponse{Service: s.serviceFromDomain(service)}, nil
}

func (s server) SearchServicesWithFilters(ctx context.Context, req *searchpb.SearchServicesWithFiltersRequest) (*searchpb.SearchServicesWithFiltersResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchServicesWithFilters")
	// Assumed Implementation:
	params := application.SearchServicesWithFilters{
		CategoryID:     req.GetCategoryId(),
		CategorySlug:   req.GetCategorySlug(),
		ServiceType:    req.GetServiceType(),
		UserID:         req.GetUserId(),
		UserType:       models.UserType(req.GetUserType()),
		Status:         models.ToStatus(req.GetStatus()),
		SearchText:     req.GetSearchText(),
		MinPrice:       req.GetMinPrice(),
		MaxPrice:       req.GetMaxPrice(),
		Tags:           req.GetTags(),
		Qualifications: req.GetQualifications(),
		Negotiable:     req.GetNegotiable(),
		Lat:            float64(req.GetLat()),
		Lng:            float64(req.GetLng()),
		Radius:         req.GetRadius(),
		Page:           req.GetPage(),
		PageSize:       req.GetPageSize(),
		SortBy:         req.GetSortBy(),
		SortOrder:      req.GetSortOrder(),
	} // Assumed struct
	services, err := s.app.SearchServicesWithFilters(ctx, params) // Assumed method
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	protoServices := make([]*searchpb.Service, len(services))
	for i, sv := range services {
		protoServices[i] = s.serviceFromDomain(sv)
	}
	return &searchpb.SearchServicesWithFiltersResponse{Services: protoServices, CurrentPage: req.GetPage()}, nil
}

func (s server) SearchServicesWithTerm(ctx context.Context, req *searchpb.SearchServicesWithTermRequest) (*searchpb.SearchServicesWithTermResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchServicesWithTerm")
	// Assumed Implementation:
	params := application.SearchServicesWithTerm{Name: req.GetName()} // Assumed struct
	services, err := s.app.SearchServicesWithTerm(ctx, params)        // Assumed method
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	protoServices := make([]*searchpb.Service, len(services))
	for i, sv := range services {
		protoServices[i] = s.serviceFromDomain(sv)
	}
	return &searchpb.SearchServicesWithTermResponse{Services: protoServices}, nil
}

func (s server) SearchServicesWithCategorySlug(ctx context.Context, req *searchpb.SearchServicesWithCategorySlugRequest) (*searchpb.SearchServicesWithCategorySlugResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchServicesWithCategorySlug")
	// Assumed Implementation:
	params := application.SearchServicesWithCategorySlug{CategorySlug: req.GetCategorySlug(), Page: req.GetPage(), PageSize: req.GetPageSize() /* Map other fields */} // Assumed struct
	services, err := s.app.SearchServicesWithCategorySlug(ctx, params)                                                                                                 // Assumed method
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	protoServices := make([]*searchpb.Service, len(services))
	for i, sv := range services {
		protoServices[i] = s.serviceFromDomain(sv)
	}
	return &searchpb.SearchServicesWithCategorySlugResponse{Services: protoServices, CurrentPage: req.GetPage()}, nil
}

func (s server) SearchServicesWithCategory(ctx context.Context, req *searchpb.SearchServicesWithCategoryRequest) (*searchpb.SearchServicesWithCategoryResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SearchServicesWithCategory")
	// Assumed Implementation:
	params := application.SearchServicesWithCategory{CategoryID: req.GetCategoryId(), Page: req.GetPage(), PageSize: req.GetPageSize() /* Map other fields */} // Assumed struct
	services, err := s.app.SearchServicesWithCategory(ctx, params)                                                                                             // Assumed method
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	protoServices := make([]*searchpb.Service, len(services))
	for i, sv := range services {
		protoServices[i] = s.serviceFromDomain(sv)
	}
	return &searchpb.SearchServicesWithCategoryResponse{Services: protoServices, CurrentPage: req.GetPage()}, nil
}

func (s server) SuggestServices(ctx context.Context, req *searchpb.SuggestServicesRequest) (*searchpb.SuggestServicesResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "SuggestServices")
	// Assumed Implementation:
	params := application.SuggestServices{Name: req.GetName()} // Assumed struct
	services, err := s.app.SuggestServices(ctx, params)        // Assumed method
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	protoServices := make([]*searchpb.Service, len(services))
	for i, sv := range services {
		protoServices[i] = s.serviceFromDomain(sv)
	}
	return &searchpb.SuggestServicesResponse{Suggestions: protoServices}, nil
}

func (s server) UnifiedSearch(ctx context.Context, req *searchpb.UnifiedSearchRequest) (*searchpb.UnifiedSearchResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "UnifiedSearch")

	log.Printf("[SearchService] UnifiedSearch: term=%s entityTypes=%v page=%d pageSize=%d",
		req.GetSearchTerm(), req.GetEntityTypes(), req.GetPage(), req.GetPageSize())

	params := application.UnifiedSearchParams{
		SearchTerm:       req.GetSearchTerm(),
		EntityTypes:      req.GetEntityTypes(),
		Page:             req.GetPage(),
		PageSize:         req.GetPageSize(),
		Lat:              float64(req.GetLat()),
		Lng:              float64(req.GetLng()),
		Radius:           req.GetRadius(),
		SortBy:           req.GetSortBy(),
		SortOrder:        req.GetSortOrder(),
		MinPrice:         req.GetMinPrice(),
		MaxPrice:         req.GetMaxPrice(),
		CategoryID:       req.GetCategoryId(),
		CategorySlug:     req.GetCategorySlug(),
		UserType:         req.GetUserType(),
		Negotiable:       req.GetNegotiable(),
		Brand:            req.GetBrand(),
		Condition:        req.GetCondition(),
		Model:            req.GetModel(),
		Tags:             req.GetTags(),
		HasVariants:      req.GetHasVariants(),
		MiddlemanService: req.GetMiddlemanService(),
		Status:           req.GetStatus(),
		ServiceType:      req.GetServiceType(),
	}

	results, err := s.app.UnifiedSearch(ctx, params)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Calculate pagination metadata
	totalPages := int64(0)
	if results.PageSize > 0 {
		totalPages = (results.TotalCount + results.PageSize - 1) / results.PageSize
	}
	hasMore := results.Page < totalPages

	// Generate next cursor (could be enhanced with timestamp/ID for stable pagination)
	nextCursor := ""
	if hasMore {
		nextCursor = fmt.Sprintf("page:%d", results.Page+1)
	}

	response := &searchpb.UnifiedSearchResponse{
		Results:      make([]*searchpb.UnifiedSearchResult, len(results.Results)),
		TotalCount:   results.TotalCount,
		Page:         results.Page,
		PageSize:     results.PageSize,
		CountsByType: results.CountsByType,
		HasMore:      hasMore,
		TotalPages:   totalPages,
		NextCursor:   nextCursor,
	}

	for i, result := range results.Results {
		response.Results[i] = s.unifiedResultFromDomain(result)
	}

	return response, nil
}

func (s server) UnifiedFeed(ctx context.Context, req *searchpb.UnifiedFeedRequest) (*searchpb.UnifiedFeedResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "UnifiedFeed")

	// Apply default pagination
	page := req.GetPage()
	pageSize := req.GetPageSize()

	// Ensure reasonable defaults
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20 // Default page size
	}

	log.Printf("[SearchService] UnifiedFeed: feed_type=%s entityTypes=%v page=%d pageSize=%d",
		req.GetFeedType(), req.GetEntityTypes(), page, pageSize)

	params := application.UnifiedFeedParams{
		EntityTypes: req.GetEntityTypes(),
		FeedType:    req.GetFeedType(),
		Page:        page,
		PageSize:    pageSize,
		Lat:         float64(req.GetLat()),
		Lng:         float64(req.GetLng()),
		Radius:      req.GetRadius(),
		UserID:      req.GetUserId(),
	}

	results, err := s.app.UnifiedFeed(ctx, params)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Calculate pagination metadata
	totalPages := int64(0)
	if pageSize > 0 && results.TotalCount > 0 {
		totalPages = (results.TotalCount + pageSize - 1) / pageSize
	}
	hasMore := page < totalPages

	// Generate next cursor
	nextCursor := ""
	if hasMore {
		nextCursor = fmt.Sprintf("page:%d", page+1)
	}

	// Create response with pagination metadata
	response := &searchpb.UnifiedFeedResponse{
		Items:       make([]*searchpb.UnifiedSearchResult, len(results.Items)),
		HasMore:     hasMore,
		TotalCount:  results.TotalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
		NextCursor:  nextCursor,
	}

	for i, item := range results.Items {
		response.Items[i] = s.unifiedResultFromDomain(item)
	}

	return response, nil
}

func (s server) GetCatalog(ctx context.Context, req *searchpb.GetCatalogRequest) (*searchpb.UnifiedFeedResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer handlePanic(span, "UnifiedFeed")

	// Apply default pagination
	page := req.GetPage()
	pageSize := req.GetPageSize()

	// Ensure reasonable defaults
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20 // Default page size
	}

	log.Printf("[SearchService] UnifiedFeed: feed_type=%s entityTypes=%v page=%d pageSize=%d",
		req.GetFeedType(), req.GetEntityTypes(), page, pageSize)

	params := application.UnifiedFeedParams{
		EntityTypes: req.GetEntityTypes(),
		FeedType:    req.GetFeedType(),
		Page:        page,
		PageSize:    pageSize,
		Lat:         float64(req.GetLat()),
		Lng:         float64(req.GetLng()),
		Radius:      req.GetRadius(),
		UserID:      req.GetUserId(),
	}

	results, err := s.app.UnifiedFeed(ctx, params)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Calculate pagination metadata
	totalPages := int64(0)
	if pageSize > 0 && results.TotalCount > 0 {
		totalPages = (results.TotalCount + pageSize - 1) / pageSize
	}
	hasMore := page < totalPages

	// Generate next cursor
	nextCursor := ""
	if hasMore {
		nextCursor = fmt.Sprintf("page:%d", page+1)
	}

	// Create response with pagination metadata
	response := &searchpb.UnifiedFeedResponse{
		Items:       make([]*searchpb.UnifiedSearchResult, len(results.Items)),
		HasMore:     hasMore,
		TotalCount:  results.TotalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
		NextCursor:  nextCursor,
	}

	for i, item := range results.Items {
		response.Items[i] = s.unifiedResultFromDomain(item)
	}

	return response, nil
}

// unifiedResultFromDomain converts the domain model of UnifiedSearchResult to the protobuf version
func (s server) unifiedResultFromDomain(result application.UnifiedSearchResult) *searchpb.UnifiedSearchResult {
	pbResult := &searchpb.UnifiedSearchResult{
		EntityType:     result.EntityType,
		RelevanceScore: float32(result.RelevanceScore),
	}

	if result.CreatedAt.IsZero() == false {
		pbResult.CreatedAt = timestamppb.New(result.CreatedAt)
	}
	if result.UpdatedAt.IsZero() == false {
		pbResult.UpdatedAt = timestamppb.New(result.UpdatedAt)
	}

	// Set the appropriate entity based on the entity type
	switch result.EntityType {
	case "product":
		if result.Product != nil {
			pbResult.Item = &searchpb.UnifiedSearchResult_Product{
				Product: s.productFromDomain(result.Product),
			}
		}
	case "service":
		if result.Service != nil {
			pbResult.Item = &searchpb.UnifiedSearchResult_Service{
				Service: s.serviceFromDomain(result.Service),
			}
		}
	case "post":
		if result.Post != nil {
			pbResult.Item = &searchpb.UnifiedSearchResult_Post{
				Post: s.postFromDomain(result.Post),
			}
		}
	}

	return pbResult
}

// --- Helper Functions ---

func (s server) orderFromDomain(order *models.Order) *searchpb.Order {
	if order == nil {
		return nil
	}
	protoItems := make([]*searchpb.Order_Item, len(order.Items))
	for i, it := range order.Items {
		protoItems[i] = &searchpb.Order_Item{
			ProductId:      it.ProductID,
			UserSellerId:   it.UserSellerID,
			ProductName:    it.ProductName,
			UserSellerName: it.UserSellerName,
			Price:          it.Price,
			Quantity:       it.Quantity,
		}
	}
	return &searchpb.Order{
		OrderId:          order.OrderID,
		UserCustomerId:   order.UserCustomerID,
		UserCustomerName: order.UserCustomerName,
		Total:            order.Total,
		Status:           order.Status,
		Items:            protoItems,
	}
}

func (s server) productFromDomain(product *models.Product) *searchpb.Product {
	if product == nil {
		return nil
	}
	pbAttrs := make([]*searchpb.Attribute, len(product.Attributes))
	for i, a := range product.Attributes {
		pbAttrs[i] = &searchpb.Attribute{Key: a.Key, Value: a.Value}
	}
	pbOpts := make([]*searchpb.Option, len(product.Options))
	for i, o := range product.Options {
		pbOpts[i] = &searchpb.Option{Name: o.Name, Value: o.Value, Price: int64(o.Price)}
	}

	// Create ItemMetric protobuf if it exists
	var pbMetrics *searchpb.ItemMetric
	if product.Metrics != nil {
		pbMetrics = &searchpb.ItemMetric{
			Id:                   product.Metrics.ID,
			EntityType:           product.Metrics.EntityType,
			LikesCount:           product.Metrics.LikesCount,
			DislikesCount:        product.Metrics.DislikesCount,
			CommentsCount:        product.Metrics.CommentsCount,
			SharedCount:          product.Metrics.SharedCount,
			AddedToWishlistCount: product.Metrics.AddedToWishlistCount,
			AddedToBasketCount:   product.Metrics.AddedToBasketCount,
			VisitedCount:         product.Metrics.VisitedCount,
			ReportedCount:        product.Metrics.ReportedCount,
			FollowerCount:        product.Metrics.FollowerCount,
			ReviewsCount:         product.Metrics.ReviewsCount,
			RatingCount:          product.Metrics.RatingCount,
			VideosCount:          product.Metrics.VideosCount,
			ImagesCount:          product.Metrics.ImagesCount,
			Rating:               product.Metrics.Rating,
			Category:             product.Metrics.Category,
			CategoryId:           product.Metrics.CategoryID,
			CategorySlug:         product.Metrics.CategorySlug,
		}
	}

	return &searchpb.Product{
		Id:               product.ProductID,
		Name:             product.Name,
		Description:      product.Description,
		BasePrice:        product.BasePrice,
		UserSellerId:     product.UserSellerID,
		CategoryId:       product.CategoryID,
		CategorySlug:     product.CategorySlug,
		Brand:            product.Brand,
		Condition:        product.Condition,
		Model:            product.Model,
		Tags:             product.Tags,
		ManageStocks:     product.ManageStock,
		Stock:            product.Stock,
		Sku:              product.SKU,
		Attributes:       pbAttrs,
		Weight:           product.Weight,
		Height:           product.Height,
		Width:            product.Width,
		Depth:            product.Depth,
		Status:           product.Status,
		Negotiable:       product.Negotiable,
		MiddlemanService: product.MiddlemanService,
		UserType:         product.UserType,
		ShippingCost:     product.ShippingCost,
		HasVariants:      product.HasVariants,
		Options:          pbOpts,
		Lat:              float32(product.Lat),
		Lng:              float32(product.Lng),
		Thumbnail:        product.Thumbnail,
		EntityType:       "product",
		Metrics:          pbMetrics,
		CreatedAt:        timestamppb.New(product.CreatedAt),
		UpdatedAt:        timestamppb.New(product.UpdatedAt),
	}
}

func (s server) postFromDomain(post *models.Post) *searchpb.Post {
	if post == nil {
		return nil
	}

	// Create ItemMetric protobuf if it exists
	var pbMetrics *searchpb.ItemMetric
	if post.Metrics != nil {
		pbMetrics = &searchpb.ItemMetric{
			Id:                   post.Metrics.ID,
			EntityType:           post.Metrics.EntityType,
			LikesCount:           post.Metrics.LikesCount,
			DislikesCount:        post.Metrics.DislikesCount,
			CommentsCount:        post.Metrics.CommentsCount,
			SharedCount:          post.Metrics.SharedCount,
			AddedToWishlistCount: post.Metrics.AddedToWishlistCount,
			AddedToBasketCount:   post.Metrics.AddedToBasketCount,
			VisitedCount:         post.Metrics.VisitedCount,
			ReportedCount:        post.Metrics.ReportedCount,
			FollowerCount:        post.Metrics.FollowerCount,
			ReviewsCount:         post.Metrics.ReviewsCount,
			RatingCount:          post.Metrics.RatingCount,
			VideosCount:          post.Metrics.VideosCount,
			ImagesCount:          post.Metrics.ImagesCount,
			Rating:               post.Metrics.Rating,
			Category:             post.Metrics.Category,
			CategoryId:           post.Metrics.CategoryID,
			CategorySlug:         post.Metrics.CategorySlug,
		}
	}

	return &searchpb.Post{
		Id:           post.PostID,
		UserId:       post.UserID,
		Name:         post.Name,
		Description:  post.Description,
		TypeOfPost:   post.TypeOfPost,
		UserType:     post.UserType,
		CategoryId:   post.CategoryID,
		CategorySlug: post.CategorySlug,
		Tags:         post.Tags,
		Status:       post.Status,
		Thumbnail:    post.Thumbnail,
		Lat:          float32(post.Lat),
		Lng:          float32(post.Lng),
		EntityType:   "post",
		Metrics:      pbMetrics,
		CreatedAt:    timestamppb.New(post.CreatedAt),
		UpdatedAt:    timestamppb.New(post.UpdatedAt),
	}
}
func (s server) serviceFromDomain(service *models.Service) *searchpb.Service {
	if service == nil {
		return nil
	}

	// Create attributes
	pbAttrs := make([]*searchpb.Attribute, len(service.Attributes))
	for i, a := range service.Attributes {
		pbAttrs[i] = &searchpb.Attribute{Key: a.Key, Value: a.Value}
	}

	// Create options
	pbOpts := make([]*searchpb.Option, len(service.Options))
	for i, o := range service.Options {
		pbOpts[i] = &searchpb.Option{Name: o.Name, Value: o.Value, Price: int64(o.Price)}
	}

	// Create ItemMetric protobuf if it exists
	var pbMetrics *searchpb.ItemMetric
	if service.Metrics != nil {
		pbMetrics = &searchpb.ItemMetric{
			Id:                   service.Metrics.ID,
			EntityType:           service.Metrics.EntityType,
			LikesCount:           service.Metrics.LikesCount,
			DislikesCount:        service.Metrics.DislikesCount,
			CommentsCount:        service.Metrics.CommentsCount,
			SharedCount:          service.Metrics.SharedCount,
			AddedToWishlistCount: service.Metrics.AddedToWishlistCount,
			AddedToBasketCount:   service.Metrics.AddedToBasketCount,
			VisitedCount:         service.Metrics.VisitedCount,
			ReportedCount:        service.Metrics.ReportedCount,
			FollowerCount:        service.Metrics.FollowerCount,
			ReviewsCount:         service.Metrics.ReviewsCount,
			RatingCount:          service.Metrics.RatingCount,
			VideosCount:          service.Metrics.VideosCount,
			ImagesCount:          service.Metrics.ImagesCount,
			Rating:               service.Metrics.Rating,
			Category:             service.Metrics.Category,
			CategoryId:           service.Metrics.CategoryID,
			CategorySlug:         service.Metrics.CategorySlug,
		}
	}

	return &searchpb.Service{
		Id:               service.ID,
		Name:             service.Name,
		Description:      service.Description,
		ServiceType:      service.ServiceType,
		BasePrice:        service.BasePrice,
		Pricing:          service.Pricing,
		Availability:     service.Availability,
		ProviderName:     service.ProviderName,
		UserId:           service.UserID,
		CategoryId:       service.CategoryID,
		CategorySlug:     service.CategorySlug,
		DescriptionShort: service.DescriptionShort,
		DescriptionLong:  service.DescriptionLong,
		Qualifications:   service.Qualifications,
		Contact:          service.Contact,
		Faq:              service.Faq,
		Tags:             service.Tags,
		Status:           service.Status,
		UserType:         service.UserType,
		ShippingCost:     service.ShippingCost,
		HasVariants:      service.HasVariants,
		MiddlemanService: service.MiddlemanService,
		Negotiable:       service.Negotiable,
		Attributes:       pbAttrs,
		Options:          pbOpts,
		Thumbnail:        service.Thumbnail,
		Lat:              float32(service.Lat),
		Lng:              float32(service.Lng),
		EntityType:       "service",
		Metrics:          pbMetrics,
		CreatedAt:        timestamppb.New(service.CreatedAt),
		UpdatedAt:        timestamppb.New(service.UpdatedAt),
	}
}

// --- Stub *FromDomain helpers for missing types ---

func (s server) userFromDomain(user *models.User) *searchpb.User {
	if user == nil {
		return nil
	}
	return &searchpb.User{
		Name: user.Username,
	}
}
