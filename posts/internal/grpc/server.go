package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"middleman/posts/postspb"

	"middleman/internal/auth"
	"middleman/internal/errorsotel"
	"middleman/posts/internal/application"
	"middleman/posts/internal/application/commands"
	"middleman/posts/internal/application/queries"
	"middleman/posts/internal/domain"
)

type server struct {
	app application.App
	postspb.UnimplementedPostsServiceServer
}

var _ postspb.PostsServiceServer = (*server)(nil)

// RegisterServer registers the gRPC server implementation
func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	postspb.RegisterPostsServiceServer(registrar, server{app: app})
	return nil
}

// -----------------------------------------------------------------------------
// 1) PRODUCT METHODS
// -----------------------------------------------------------------------------
func (s server) AddPost(ctx context.Context, req *postspb.AddPostRequest) (*postspb.AddPostResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject

	cmd := commands.AddPost{
		ID:           uuid.New().String(),
		Name:         req.GetName(),
		Description:  req.GetDescription(),
		TypeOfPost:   domain.ToTypeOfPost(req.GetTypeOfPost()),
		UserID:       userID,
		UserType:     domain.ToUserType(req.GetUserType()),
		CategoryID:   req.GetCategoryId(),
		CategorySlug: req.GetCategorySlug(),
		Status:       domain.ToPostStatus(req.GetStatus()),
		Thumbnail:    req.GetThumbnail(),
		Lat:          float64(req.GetLat()),
		Lng:          float64(req.GetLng()),
	}
	// Then call your application service:
	if err := s.app.AddPost(ctx, cmd); err != nil {
		return nil, err
	}

	return &postspb.AddPostResponse{Id: cmd.ID}, nil
}
func (s server) UpdatePost(ctx context.Context, req *postspb.UpdatePostRequest) (*postspb.UpdatePostResponse, error) {
	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	cmd := commands.UpdatePost{
		ID:           req.GetId(),
		Name:         req.GetName(),
		Description:  req.GetDescription(),
		TypeOfPost:   domain.ToTypeOfPost(req.GetTypeOfPost()),
		CategoryID:   req.GetCategoryId(),
		CategorySlug: req.GetCategorySlug(),
		Status:       domain.ToPostStatus(req.GetStatus()),
		Thumbnail:    req.GetThumbnail(),
	}
	// Then call your application service:
	if err := s.app.UpdatePost(ctx, cmd); err != nil {
		return nil, err
	}

	return &postspb.UpdatePostResponse{Id: cmd.ID}, nil
}

// For listing posts in pages
func (s server) GetPosts(ctx context.Context, request *postspb.GetPostsRequest) (*postspb.GetPostsResponse, error) {
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
	posts, totalCount, err := s.app.GetPosts(ctx, queries.GetPosts{
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

	protoPosts := make([]*postspb.Post, len(posts))
	for i, post := range posts {
		protoPosts[i] = s.postFromDomain(post)
	}

	return &postspb.GetPostsResponse{
		Posts:       protoPosts,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

func (s server) GetUserPosts(ctx context.Context, request *postspb.GetUserPostsRequest) (*postspb.GetUserPostsResponse, error) {
	span := trace.SpanFromContext(ctx)
	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10 // <-- your default page size
	}

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	// 2) Guard for zero/negative Page
	page := request.GetPage()
	if page <= 0 {
		page = 1 // <-- your default page number
	}
	posts, totalCount, err := s.app.GetUserPosts(ctx, queries.GetUserPosts{
		UserId:    request.GetUserId(),
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

	protoPosts := make([]*postspb.Post, len(posts))
	for i, post := range posts {
		protoPosts[i] = s.postFromDomain(post)
	}

	return &postspb.GetUserPostsResponse{
		Posts:       protoPosts,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}
func (s server) GetPublicCatalog(ctx context.Context, request *postspb.GetPublicCatalogRequest) (*postspb.GetPublicCatalogResponse, error) {
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
	posts, totalCount, err := s.app.GetPublicCatalog(ctx, queries.GetPublicCatalog{
		UserId:    request.GetUserId(),
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

	protoPosts := make([]*postspb.Post, len(posts))
	for i, post := range posts {
		protoPosts[i] = s.postFromDomain(post)
	}

	return &postspb.GetPublicCatalogResponse{
		Posts:       protoPosts,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

// -----------------------------------------------------------------------------
// 4) GetPostsWithFilters (the method we are fixing)
// -----------------------------------------------------------------------------
func (s server) GetPostsWithFilters(ctx context.Context, req *postspb.GetPostsWithFiltersRequest) (*postspb.GetPostsWithFiltersResponse, error) {
	span := trace.SpanFromContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			span.RecordError(fmt.Errorf("%v", r))
			span.SetStatus(codes.Error, "panic in GetPostsWithFilters")
		}
	}()

	page := req.GetPage()
	if page < 1 {
		page = 1
	}
	pageSize := req.GetPageSize()
	if pageSize < 1 {
		pageSize = 10
	}

	q := queries.GetPostsWithFilters{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		TypeOfPost:  req.GetTypeOfPost(),
		Tags:        req.GetTags(),
		Lat:         float64(req.GetLat()),
		Lng:         float64(req.GetLng()),
		Radius:      int64(req.GetRadius()),
		Page:        page,
		PageSize:    pageSize,
		SortBy:      req.GetSortBy(),
		SortOrder:   req.GetSortOrder(),
	}
	domainPosts, totalCount, err := s.app.GetPostsWithFilters(ctx, q)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	totalPages := (totalCount + pageSize - 1) / pageSize
	protoPosts := make([]*postspb.Post, len(domainPosts))
	for i, post := range domainPosts {
		protoPosts[i] = s.postFromDomain(post)
	}
	return &postspb.GetPostsWithFiltersResponse{
		Posts:       protoPosts,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}

// Get a single post
func (s server) GetPost(ctx context.Context, request *postspb.GetPostRequest) (*postspb.GetPostResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("PostID", request.GetId()))

	post, err := s.app.GetPost(ctx, queries.GetPost{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &postspb.GetPostResponse{Post: s.postFromDomain(post)}, nil
}

// Remove a post from listing
func (s server) RemovePost(ctx context.Context, request *postspb.RemovePostRequest) (*postspb.RemovePostResponse, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("PostID", request.GetId()))

	err := s.app.RemovePost(ctx, commands.RemovePost{
		ID:     request.GetId(),
		UserID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &postspb.RemovePostResponse{}, err
}

// 3) ArchivePost
func (s server) ArchivePost(ctx context.Context, request *postspb.ArchivePostRequest) (*postspb.ArchivePostResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("PostID", request.GetPostId()))
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	userID := claims.Subject
	err := s.app.ArchivePost(ctx, commands.ArchivePost{
		ID:     request.GetPostId(),
		UserID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &postspb.ArchivePostResponse{
		PostId:   request.GetPostId(),
		Archived: true,
	}, nil
}

func (s server) AddPostThumbnail(ctx context.Context, request *postspb.AddPostThumbnailRequest) (*postspb.AddPostThumbnailResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("PostID", request.GetPostId()))

	err := s.app.AddPostThumbnail(ctx, commands.AddPostThumbnail{
		ID:        request.GetPostId(),
		Thumbnail: request.GetThumbnail(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &postspb.AddPostThumbnailResponse{}, nil
}

// For listing posts in a category
func (s server) GetPostsByCategory(ctx context.Context, request *postspb.GetPostsByCategoryRequest) (*postspb.GetPostsByCategoryResponse, error) {
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

	posts, totalCount, err := s.app.GetPostsByCategory(ctx, queries.GetPostsByCategory{
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

	protoPosts := make([]*postspb.Post, len(posts))
	for i, post := range posts {
		protoPosts[i] = s.postFromDomain(post)
	}

	return &postspb.GetPostsByCategoryResponse{
		Posts:       protoPosts,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

func (s server) GetPostsByCategorySlug(ctx context.Context, request *postspb.GetPostsByCategorySlugRequest) (*postspb.GetPostsByCategorySlugResponse, error) {
	span := trace.SpanFromContext(ctx)

	pageSize := request.GetPageSize()
	if pageSize <= 0 {
		pageSize = 10
	}

	page := request.GetPage()
	if page <= 0 {
		page = 1
	}

	posts, totalCount, err := s.app.GetPostsByCategorySlug(ctx, queries.GetPostsByCategorySlug{
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

	protoPosts := make([]*postspb.Post, len(posts))
	for i, post := range posts {
		protoPosts[i] = s.postFromDomain(post)
	}

	return &postspb.GetPostsByCategorySlugResponse{
		Posts:       protoPosts,
		TotalCount:  totalCount,
		CurrentPage: request.GetPage(),
		TotalPages:  totalPages,
	}, nil
}

func (s server) UpdatePostThumbnail(ctx context.Context, request *postspb.UpdatePostThumbnailRequest) (*postspb.UpdatePostThumbnailResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("PostID", request.GetPostId()))

	err := s.app.UpdatePostThumbnail(ctx, commands.UpdatePostThumbnail{
		ID:        request.GetPostId(),
		Thumbnail: request.GetThumbnail(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &postspb.UpdatePostThumbnailResponse{}, nil
}
func (s server) postFromDomain(p *domain.CatalogPost) *postspb.Post {

	// (2) Convert domain Attributes => repeated postspb.Attribu

	return &postspb.Post{
		Id:          p.ID,
		UserId:      p.UserID,
		Name:        p.Name,
		Description: p.Description,
		Tags:        p.Tags,
		Status:      p.Status.String(),
		Thumbnail:   p.Thumbnail,
		Lat:         float32(p.Lat),
		Lng:         float32(p.Lng),
	}
}
