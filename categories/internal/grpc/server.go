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

	"middleman/categories/categoriespb"
	"middleman/categories/internal/application"
	"middleman/categories/internal/application/commands"
	"middleman/categories/internal/application/queries"
	"middleman/categories/internal/domain"
	"middleman/internal/auth"
	"middleman/internal/errorsotel"
)

type server struct {
	app application.App
	categoriespb.UnimplementedCategoriesServiceServer
}

var _ categoriespb.CategoriesServiceServer = (*server)(nil)

// RegisterServer registers the gRPC server implementation
func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	categoriespb.RegisterCategoriesServiceServer(registrar, server{app: app})
	return nil
}

func (s server) AddCategory(ctx context.Context, req *categoriespb.AddCategoryRequest) (*categoriespb.AddCategoryResponse, error) {

	cmd := commands.AddCategory{
		ID:               uuid.New().String(),
		Description:      req.GetDescription(),
		ParentID:         req.GetParentId(),
		GoogleCategoryID: req.GetGoogleCategoryId(),
		Tags:             req.GetTags(),
		IsActive:         req.GetIsActive(),
		Slug:             req.GetSlug(),
		SeoTitle:         req.GetSeoTitle(),
		SeoKeywords:      req.GetSeoKeywords(),
		SeoDesc:          req.GetSeoDesc(),
		CategoryType:     req.GetCategoryType(),
		Lang:             req.GetLang(),
	}
	// Then call your application service:
	if err := s.app.AddCategory(ctx, cmd); err != nil {
		return nil, err
	}
	// Return the new ID in the response:
	return &categoriespb.AddCategoryResponse{Id: cmd.ID}, nil
}

func (s server) RebrandCategory(ctx context.Context, request *categoriespb.RebrandCategoryRequest) (*categoriespb.RebrandCategoryResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("CategoryID", request.GetId()))

	err := s.app.RebrandCategory(ctx, commands.RebrandCategory{
		ID:          request.GetId(),
		Slug:        request.GetNewSlug(),
		Description: request.GetNewDesc(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &categoriespb.RebrandCategoryResponse{}, err
}
func (s server) UpdateCategory(ctx context.Context, request *categoriespb.UpdateCategoryRequest) (*categoriespb.UpdateCategoryResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("CategoryID", request.GetId()))

	err := s.app.UpdateCategory(ctx, commands.UpdateCategory{
		ID:               request.GetId(),
		Description:      request.GetDescription(),
		ParentID:         request.GetParentId(),
		GoogleCategoryID: request.GetGoogleCategoryId(),
		Tags:             request.GetTags(),
		IsActive:         request.GetIsActive(),
		Slug:             request.GetSlug(),
		SeoTitle:         request.GetSeoTitle(),
		SeoKeywords:      request.GetSeoKeywords(),
		SeoDesc:          request.GetSeoDesc(),
		CategoryType:     request.GetCategoryType(),
		Lang:             request.GetLang(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &categoriespb.UpdateCategoryResponse{}, err
}

func (s server) GetCategories(ctx context.Context, request *categoriespb.GetCategoriesRequest) (*categoriespb.GetCategoriesResponse, error) {
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

	categories, totalCount, err := s.app.GetCategories(ctx, queries.GetCategories{
		CategoryType: request.GetCategoryType(),
		Lang:         request.GetLang(),
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

	protoCategories := make([]*categoriespb.Category, len(categories))
	for i, category := range categories {
		protoCategories[i] = s.categoryFromDomain(category)
	}

	return &categoriespb.GetCategoriesResponse{
		Categories:  protoCategories,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}
func (s server) GetMainCategories(ctx context.Context, request *categoriespb.GetMainCategoriesRequest) (*categoriespb.GetMainCategoriesResponse, error) {
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

	categories, totalCount, err := s.app.GetMainCategories(ctx, queries.GetMainCategories{
		CategoryType: request.GetCategoryType(),
		Lang:         request.GetLang(),
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

	protoCategories := make([]*categoriespb.Category, len(categories))
	for i, category := range categories {
		protoCategories[i] = s.categoryFromDomain(category)
	}

	return &categoriespb.GetMainCategoriesResponse{
		Categories:  protoCategories,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}
func (s server) GetAllMainCategories(ctx context.Context, request *categoriespb.GetAllMainCategoriesRequest) (*categoriespb.GetAllMainCategoriesResponse, error) {
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

	categories, totalCount, err := s.app.GetAllMainCategories(ctx, queries.GetAllMainCategories{
		Lang:      request.GetLang(),
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

	protoCategories := make([]*categoriespb.Category, len(categories))
	for i, category := range categories {
		protoCategories[i] = s.categoryFromDomain(category)
	}

	return &categoriespb.GetAllMainCategoriesResponse{
		Categories:  protoCategories,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}
func (s server) GetSubCategories(ctx context.Context, request *categoriespb.GetSubCategoriesRequest) (*categoriespb.GetSubCategoriesResponse, error) {
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

	categories, totalCount, err := s.app.GetSubCategories(ctx, queries.GetSubCategories{
		Lang:             request.GetLang(),
		ParentCategoryID: request.GetParentCategoryId(),
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

	protoCategories := make([]*categoriespb.Category, len(categories))
	for i, category := range categories {
		protoCategories[i] = s.categoryFromDomain(category)
	}

	return &categoriespb.GetSubCategoriesResponse{
		Categories:  protoCategories,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

func (s server) GetCategory(ctx context.Context, request *categoriespb.GetCategoryRequest) (*categoriespb.GetCategoryResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("CategoryID", request.GetId()))

	category, err := s.app.GetCategory(ctx, queries.GetCategory{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &categoriespb.GetCategoryResponse{Category: s.categoryFromDomain(category)}, nil
}

func (s server) GetCategoryBySlug(ctx context.Context, request *categoriespb.GetCategoryBySlugRequest) (*categoriespb.GetCategoryBySlugResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("Category slug", request.GetSlug()))

	category, err := s.app.GetCategoryBySlug(ctx, queries.GetCategoryBySlug{
		Slug: request.GetSlug(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &categoriespb.GetCategoryBySlugResponse{Category: s.categoryFromDomain(category)}, nil
}

func (s server) RemoveCategory(ctx context.Context, request *categoriespb.RemoveCategoryRequest) (*categoriespb.RemoveCategoryResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("CategoryID", request.GetId()))

	err := s.app.RemoveCategory(ctx, commands.RemoveCategory{
		ID:     request.GetId(),
		UserID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &categoriespb.RemoveCategoryResponse{}, err
}

func (s server) ArchiveCategory(ctx context.Context, request *categoriespb.ArchiveCategoryRequest) (*categoriespb.ArchiveCategoryResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("CategoryID", request.GetId()))

	err := s.app.ArchiveCategory(ctx, commands.ArchiveCategory{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &categoriespb.ArchiveCategoryResponse{
		CategoryId: request.GetId(),
		Archived:   true,
	}, nil
}

func (s server) AddFilter(ctx context.Context, req *categoriespb.AddFilterRequest) (*categoriespb.AddFilterResponse, error) {
	varID := uuid.New().String()

	cmd := commands.AddFilter{
		ID:         varID,
		CategoryID: req.GetCategoryId(),
		Name:       req.GetName(),
	}

	if err := s.app.AddFilter(ctx, cmd); err != nil {
		return nil, err
	}
	return &categoriespb.AddFilterResponse{FilterId: varID}, nil
}

func (s server) ArchiveFilter(ctx context.Context, request *categoriespb.ArchiveFilterRequest) (*categoriespb.ArchiveFilterResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("FilterID", request.GetFilterId()))

	err := s.app.ArchiveFilter(ctx, commands.ArchiveFilter{
		ID: request.GetFilterId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &categoriespb.ArchiveFilterResponse{
		FilterId: request.GetFilterId(),
		Archived: true,
	}, nil
}

func (s server) RemoveFilter(ctx context.Context, request *categoriespb.RemoveFilterRequest) (*categoriespb.RemoveFilterResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("FilterID", request.GetFilterId()))

	err := s.app.RemoveFilter(ctx, commands.RemoveFilter{
		ID: request.GetFilterId(),
		// reason if needed
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &categoriespb.RemoveFilterResponse{
		FilterId: request.GetFilterId(),
	}, nil
}

func (s server) categoryFromDomain(p *domain.CatalogCategory) *categoriespb.Category {

	return &categoriespb.Category{
		Id:               p.ID,
		Description:      p.Description,
		Slug:             p.Slug,
		GoogleCategoryId: p.GoogleCategoryID,
		IsActive:         p.IsActive,
		SeoTitle:         p.SeoTitle,
		SeoKeywords:      p.SeoKeywords,
		SeoDesc:          p.SeoDesc,
		CategoryType:     p.CategoryType,
		Lang:             p.Lang,
	}
}

func (s server) filterFromDomain(v *domain.Filter) *categoriespb.Filter {

	return &categoriespb.Filter{
		Id:         v.ID(),
		CategoryId: v.CategoryID,
		Name:       v.Name,
	}
}
