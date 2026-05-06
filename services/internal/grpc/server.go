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
	"middleman/services/internal/application"
	"middleman/services/internal/application/commands"
	"middleman/services/internal/application/queries"
	"middleman/services/internal/domain"
	"middleman/services/servicespb"
	"time"
)

type server struct {
	app application.App
	servicespb.UnimplementedServicesServiceServer
}

var _ servicespb.ServicesServiceServer = (*server)(nil)

// RegisterServer registers the gRPC server implementation
func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	servicespb.RegisterServicesServiceServer(registrar, server{app: app})
	return nil
}

// -----------------------------------------------------------------------------
// 1) PRODUCT METHODS
// -----------------------------------------------------------------------------
func (s server) AddService(ctx context.Context, req *servicespb.AddServiceRequest) (*servicespb.AddServiceResponse, error) {

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

	cmd := commands.AddService{
		ID:               uuid.New().String(),
		Name:             req.GetName(),
		Description:      req.GetDescription(),
		BasePrice:        req.GetBasePrice(),
		ServiceType:      req.GetServiceType(),
		Pricing:          req.GetPricing(),
		Availability:     req.GetAvailability(),
		ProviderName:     req.GetProviderName(),
		UserID:           userID,
		CategoryID:       req.GetCategoryId(),
		CategorySlug:     req.GetCategorySlug(),
		DescriptionShort: req.GetDescriptionShort(),
		DescriptionLong:  req.GetDescriptionLong(),
		Qualifications:   req.GetQualifications(),
		Contact:          req.Contact,
		Faq:              req.Faq,
		Tags:             req.GetTags(), // if using CSV
		Attributes:       domainAttrs,
		Status:           domain.ToServiceStatus(req.GetStatus()),
		Negotiable:       req.GetNegotiable(),
		MiddlemanService: req.GetMiddlemanService(),
		ShippingCost:     req.GetShippingCost(),
		HasVariants:      req.GetHasVariants(),
		Options:          domainOptions,
		Thumbnail:        req.GetThumbnail(),
		Lat:              float64(req.GetLat()),
		Lng:              float64(req.GetLng()),
	}
	// Then call your application service:
	if err := s.app.AddService(ctx, cmd); err != nil {
		return nil, err
	}
	// Return the new ID in the response:
	return &servicespb.AddServiceResponse{Id: cmd.ID}, nil
}
func (s server) UpdateService(ctx context.Context, req *servicespb.UpdateServiceRequest) (*servicespb.UpdateServiceResponse, error) {

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

	cmd := commands.UpdateService{
		ID:               req.GetId(),
		Name:             req.GetName(),
		Description:      req.GetDescription(),
		ServiceType:      req.GetServiceType(),
		BasePrice:        req.GetBasePrice(),
		Pricing:          req.GetPricing(),
		Availability:     req.GetAvailability(),
		ProviderName:     req.GetProviderName(),
		UserID:           userID,
		CategoryID:       req.GetCategoryId(),
		CategorySlug:     req.GetCategorySlug(),
		DescriptionShort: req.GetDescriptionShort(),
		DescriptionLong:  req.GetDescriptionLong(),
		Qualifications:   req.GetQualifications(),
		Contact:          req.Contact,
		Faq:              req.Faq,
		Tags:             req.GetTags(), // if using CSV
		Attributes:       domainAttrs,
		Status:           domain.ToServiceStatus(req.GetStatus()),
		Negotiable:       req.GetNegotiable(),
		MiddlemanService: req.GetMiddlemanService(),
		ShippingCost:     req.GetShippingCost(),
		HasVariants:      req.GetHasVariants(),
		Options:          domainOptions,
		Thumbnail:        req.GetThumbnail(),
		Lat:              float64(req.GetLat()),
		Lng:              float64(req.GetLng()),
	}
	// Then call your application service:
	if err := s.app.UpdateService(ctx, cmd); err != nil {
		return nil, err
	}
	// Return the new ID in the response:
	return &servicespb.UpdateServiceResponse{Id: cmd.ID}, nil
}
func (s server) RebrandService(ctx context.Context, request *servicespb.RebrandServiceRequest) (*servicespb.RebrandServiceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ServiceID", request.GetId()))

	err := s.app.RebrandService(ctx, commands.RebrandService{
		ID:          request.GetId(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		//Brand:       request.GetBrand(),
		//Model:       request.GetModel(),
		//Condition:   request.GetCondition(),
		// Possibly other fields if your domain command supports them
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &servicespb.RebrandServiceResponse{}, err
}

// For Service Price Increase
func (s server) IncreaseServicePrice(ctx context.Context, request *servicespb.IncreaseServicePriceRequest) (*servicespb.IncreaseServicePriceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ServiceID", request.GetServiceId()))

	err := s.app.IncreaseServicePrice(ctx, commands.IncreaseServicePrice{
		ID:    request.GetServiceId(),
		Price: request.GetPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Suppose your domain service can return oldPrice and newPrice. You could fill them here if needed:
	return &servicespb.IncreaseServicePriceResponse{
		ServiceId: request.GetServiceId(),
		NewPrice:  request.GetPrice(),
	}, nil
}

// For Service Price Decrease
func (s server) DecreaseServicePrice(ctx context.Context, request *servicespb.DecreaseServicePriceRequest) (*servicespb.DecreaseServicePriceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ServiceID", request.GetServiceId()))

	err := s.app.DecreaseServicePrice(ctx, commands.DecreaseServicePrice{
		ID:    request.GetServiceId(),
		Price: request.GetNewPrice(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &servicespb.DecreaseServicePriceResponse{
		ServiceId: request.GetServiceId(),
		NewPrice:  request.GetNewPrice(),
	}, nil
}

// For listing services in pages
func (s server) GetServices(ctx context.Context, request *servicespb.GetServicesRequest) (*servicespb.GetServicesResponse, error) {
	span := trace.SpanFromContext(ctx)
	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10 // <-- your default page size
	}

	// 2) Guard for zero/negative Page
	page := request.GetPage()
	if page <= 0 {
		page = 1 // <-- your default page number
	}

	services, totalCount, err := s.app.GetServices(ctx, queries.GetServices{
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

	protoServices := make([]*servicespb.Service, len(services))
	for i, service := range services {
		protoServices[i] = s.serviceFromDomain(service)
	}

	return &servicespb.GetServicesResponse{
		Services:    protoServices,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

// For listing services in a category
func (s server) GetServicesByCategory(ctx context.Context, request *servicespb.GetServicesByCategoryRequest) (*servicespb.GetServicesByCategoryResponse, error) {
	span := trace.SpanFromContext(ctx)

	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10 // <-- your default page size
	}

	// 2) Guard for zero/negative Page
	page := request.GetPage()
	if page <= 0 {
		page = 1 // <-- your default page number
	}

	services, totalCount, err := s.app.GetServicesByCategory(ctx, queries.GetServicesByCategory{
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

	protoServices := make([]*servicespb.Service, len(services))
	for i, service := range services {
		protoServices[i] = s.serviceFromDomain(service)
	}

	return &servicespb.GetServicesByCategoryResponse{
		Services:    protoServices,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}
func (s server) GetServicesByCategorySlug(ctx context.Context, request *servicespb.GetServicesByCategorySlugRequest) (*servicespb.GetServicesByCategorySlugResponse, error) {
	span := trace.SpanFromContext(ctx)

	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10 // <-- your default page size
	}

	// 2) Guard for zero/negative Page
	page := request.GetPage()
	if page <= 0 {
		page = 1 // <-- your default page number
	}

	services, totalCount, err := s.app.GetServicesByCategorySlug(ctx, queries.GetServicesByCategorySlug{
		CategorySlug: request.GetCategorySlug(),
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

	protoServices := make([]*servicespb.Service, len(services))
	for i, service := range services {
		protoServices[i] = s.serviceFromDomain(service)
	}

	return &servicespb.GetServicesByCategorySlugResponse{
		Services:    protoServices,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}
func (s server) GetServicesWithFilter(ctx context.Context, request *servicespb.GetServicesWithFilterRequest) (*servicespb.GetServicesWithFilterResponse, error) {
	span := trace.SpanFromContext(ctx)

	// Validate and set default values for pagination
	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}
	page := request.GetPage()
	if page <= 0 {
		page = 1
	}

	// Convert proto fields to domain types
	var availableFrom, availableTo time.Time
	if request.GetAvailableFrom() > 0 {
		availableFrom = time.Unix(request.GetAvailableFrom(), 0)
	}
	if request.GetAvailableTo() > 0 {
		availableTo = time.Unix(request.GetAvailableTo(), 0)
	}

	// Execute the filtered query
	services, totalCount, err := s.app.GetServicesWithFilter(ctx, queries.GetServicesWithFilter{
		CategoryID:       request.GetCategoryId(),
		CategorySlug:     request.GetCategorySlug(),
		ServiceType:      request.GetServiceType(),
		Status:           domain.ToServiceStatus(request.GetStatus()),
		SearchText:       request.GetSearchText(),
		MinPrice:         request.GetMinPrice(),
		MaxPrice:         request.GetMaxPrice(),
		Latitude:         float64(request.GetLat()),
		Longitude:        float64(request.GetLng()),
		Radius:           float64(request.GetRadius()),
		AvailableFrom:    availableFrom,
		AvailableTo:      availableTo,
		HasVariants:      request.GetHasVariants(),
		Negotiable:       request.GetNegotiable(),
		MiddlemanService: request.GetMiddlemanService(),
		UserType:         domain.ToUserType(request.GetUserType()),
		Tags:             request.GetTags(),
		Qualifications:   request.GetQualifications(),
		Page:             page,
		PageSize:         pageSize,
		SortBy:           request.GetSortBy(),
		SortOrder:        request.GetSortOrder(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Calculate total pages
	totalPages := (totalCount + pageSize - 1) / pageSize

	// Convert domain services to proto services
	protoServices := make([]*servicespb.Service, len(services))
	for i, service := range services {
		protoServices[i] = s.serviceFromDomain(service)
	}

	return &servicespb.GetServicesWithFilterResponse{
		Services:    protoServices,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}

// Get a single service
func (s server) GetService(ctx context.Context, request *servicespb.GetServiceRequest) (*servicespb.GetServiceResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ServiceID", request.GetId()))

	service, err := s.app.GetService(ctx, queries.GetService{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &servicespb.GetServiceResponse{Service: s.serviceFromDomain(service)}, nil
}

// Return a user's service catalog
func (s server) GetCatalog(ctx context.Context, request *servicespb.GetCatalogRequest) (*servicespb.GetCatalogResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("UserID", request.GetUserId()))

	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10 // <-- your default page size
	}

	// 2) Guard for zero/negative Page
	page := request.GetPage()
	if page <= 0 {
		page = 1 // <-- your default page number
	}

	services, totalCount, err := s.app.GetCatalog(ctx, queries.GetCatalog{
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

	protoServices := make([]*servicespb.Service, len(services))
	for i, service := range services {
		protoServices[i] = s.serviceFromDomain(service)
	}

	return &servicespb.GetCatalogResponse{
		Services:    protoServices,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}
func (s server) GetPublicCatalog(ctx context.Context, request *servicespb.GetPublicCatalogRequest) (*servicespb.GetPublicCatalogResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("UserID", request.GetUserId()))

	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10 // <-- your default page size
	}

	// 2) Guard for zero/negative Page
	page := request.GetPage()
	if page <= 0 {
		page = 1 // <-- your default page number
	}

	services, totalCount, err := s.app.GetPublicCatalog(ctx, queries.GetPublicCatalog{
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

	protoServices := make([]*servicespb.Service, len(services))
	for i, service := range services {
		protoServices[i] = s.serviceFromDomain(service)
	}

	return &servicespb.GetPublicCatalogResponse{
		Services:    protoServices,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

// Remove a service from listing
func (s server) RemoveService(ctx context.Context, request *servicespb.RemoveServiceRequest) (*servicespb.RemoveServiceResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ServiceID", request.GetId()))

	err := s.app.RemoveService(ctx, commands.RemoveService{
		ID:     request.GetId(),
		UserID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &servicespb.RemoveServiceResponse{}, err
}

// 3) ArchiveService
func (s server) ArchiveService(ctx context.Context, request *servicespb.ArchiveServiceRequest) (*servicespb.ArchiveServiceResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("ServiceID", request.GetServiceId()))

	err := s.app.ArchiveService(ctx, commands.ArchiveService{
		ID:           request.GetServiceId(),
		UserSellerID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &servicespb.ArchiveServiceResponse{
		ServiceId: request.GetServiceId(),
		Archived:  true,
	}, nil
}

func (s server) serviceFromDomain(p *domain.CatalogService) *servicespb.Service {

	// (2) Convert domain Attributes => repeated servicespb.Attribute
	pbAttributes := make([]*servicespb.Attribute, len(p.Attributes))
	for i, attr := range p.Attributes {
		pbAttributes[i] = &servicespb.Attribute{
			Key:   attr.Key,
			Value: attr.Value,
		}
	}

	// (3) Convert domain Options => repeated servicespb.Option
	pbOptions := make([]*servicespb.Option, len(p.Options))
	for i, opt := range p.Options {
		pbOptions[i] = &servicespb.Option{
			Name:  opt.Name,
			Value: opt.Value,
			// domain has Price as int64 or float64? If float64, cast to int64:
			Price: int64(opt.Price),
		}
	}

	return &servicespb.Service{
		Id:               p.ID,
		Name:             p.Name,
		Description:      p.Description,
		ServiceType:      p.ServiceType,
		BasePrice:        p.BasePrice,
		Pricing:          p.Pricing,
		Availability:     p.Availability,
		ProviderName:     p.ProviderName,
		UserId:           p.UserID,
		CategoryId:       p.CategoryID,
		CategorySlug:     p.CategorySlug,
		DescriptionShort: p.DescriptionShort,
		DescriptionLong:  p.DescriptionLong,
		Qualifications:   p.Qualifications,
		Contact:          p.Contact,
		Faq:              p.Faq,
		Tags:             p.Tags,
		Attributes:       pbAttributes,
		Status:           p.Status.String(),
		Negotiable:       p.Negotiable,
		MiddlemanService: p.MiddlemanService,
		HasVariants:      p.HasVariants,
		UserType:         p.UserType.String(),
		ShippingCost:     p.ShippingCost,
		Options:          pbOptions,
		Thumbnail:        p.Thumbnail,
		Lat:              float32(p.Lat),
		Lng:              float32(p.Lng),
	}
}
