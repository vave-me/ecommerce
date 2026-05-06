// File: search/internal/grpc/post_repository.go
package grpc

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"middleman/internal/rpc"
	"middleman/posts/postspb"
	"middleman/search/internal/application"
	"middleman/search/internal/models"
)

// PostRepository calls the remote posts service (gRPC) as a fallback.
type PostRepository struct {
	endpoint string
}

var _ application.PostRepository = (*PostRepository)(nil)

// NewPostRepository instantiates the gRPC-based fallback repo.
func NewPostRepository(endpoint string) PostRepository {
	return PostRepository{
		endpoint: endpoint,
	}
}

// Find retrieves a post by ID from the posts microservice (via gRPC).
func (r PostRepository) Find(ctx context.Context, postID string) (*models.Post, error) {
	log.Printf("Find: retrieving post with ID=%s via gRPC fallback", postID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("Find: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.GetPost(ctx, &postspb.GetPostRequest{Id: postID})
	if err != nil {
		return nil, fmt.Errorf("GetPost RPC failed: %w", err)
	}
	return r.postToDomain(resp.GetPost()), nil
}

func (r PostRepository) GetCatalog(ctx context.Context, userID string) ([]*models.Post, error) {
	log.Printf("Find: retrieving deal with ID=%s via gRPC fallback", userID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("Find: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.GetPublicCatalog(ctx, &postspb.GetPublicCatalogRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("GetDeal RPC failed: %w", err)
	}

	var results []*models.Post
	for _, p := range resp.GetPosts() {
		domainProd := r.postToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// SearchWithFilters calls the remote microservice's GetPostsWithFilters (or similar).
func (r PostRepository) SearchPostsWithFilters(
	ctx context.Context,
	name string,
	description string,
	tags []string,
	status string,
	userType string,
	thumbnail string,
	offset int64,
	limit int64,
	lat float64,
	lng float64,
	radius int64,
	page int64,
	pageSize int64,
	sortBy string,
	sortOrder string,
) ([]*models.Post, error) {

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.GetPostsWithFilters(ctx, &postspb.GetPostsWithFiltersRequest{
		Name:      name,
		Tags:      tags,
		Status:    status,
		UserType:  userType,
		Lat:       float32(lat),
		Lng:       float32(lng),
		Radius:    float32(radius),
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetPostsWithFilters RPC failed: %w", err)
	}

	var results []*models.Post
	for _, p := range resp.GetPosts() {
		domainProd := r.postToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// SearchWithFilters calls the remote microservice's GetPostsWithFilters (or similar).
func (r PostRepository) SearchPostsWithCategorySlug(
	ctx context.Context,
	categorySlug string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*models.Post, error) {

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.GetPostsByCategorySlug(ctx, &postspb.GetPostsByCategorySlugRequest{
		CategorySlug: categorySlug,
		Page:         page,
		PageSize:     pageSize,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetPostsWithFilters RPC failed: %w", err)
	}

	var results []*models.Post
	for _, p := range resp.GetPosts() {
		domainProd := r.postToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// SearchWithFilters calls the remote microservice's GetPostsWithFilters (or similar).
func (r PostRepository) SearchPostsWithCategory(
	ctx context.Context,
	categoryId string,
	page, pageSize int64,
	sortBy, sortOrder string,
) ([]*models.Post, error) {

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.GetPostsByCategory(ctx, &postspb.GetPostsByCategoryRequest{
		CategoryId: categoryId,
		Page:       page,
		PageSize:   pageSize,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetPostsWithFilters RPC failed: %w", err)
	}

	var results []*models.Post
	for _, p := range resp.GetPosts() {
		domainProd := r.postToDomain(p)
		if domainProd != nil {
			results = append(results, domainProd)
		}
	}
	return results, nil
}

// Add calls AddPost in the remote gRPC service.
func (r PostRepository) Add(ctx context.Context, postID string, name string, description string, basePrice int64, userSellerID string, categoryID string, brand string, condition string, model string, tags []string) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	_, err = client.AddPost(ctx, &postspb.AddPostRequest{
		Name:        name,
		Description: description,
		CategoryId:  categoryID,

		Tags: tags,
	})
	return err
}

// Update calls UpdatePost in the remote gRPC service (or partial, etc.).
func (r PostRepository) Update(ctx context.Context, postID string, price int64) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	_, err = client.UpdatePost(ctx, &postspb.UpdatePostRequest{
		Id: postID,
	})
	return err
}

// Remove calls RemovePost in the remote gRPC service.
func (r PostRepository) Remove(ctx context.Context, postID string) error {
	conn, err := r.dial(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	_, err = client.RemovePost(ctx, &postspb.RemovePostRequest{Id: postID})
	return err
}

// postToDomain converts a postspb.Post into our internal models.Post.
func (r PostRepository) postToDomain(pb *postspb.Post) *models.Post {
	if pb == nil {
		return nil
	}

	return &models.Post{
		PostID:       pb.GetId(),
		Name:         pb.GetName(),
		Description:  pb.GetDescription(),
		UserID:       pb.GetUserId(),
		TypeOfPost:   pb.GetTypeOfPost(),
		CategoryID:   pb.GetCategoryId(),
		CategorySlug: pb.GetCategorySlug(),
		Tags:         pb.GetTags(),
		Status:       pb.GetStatus(),
		UserType:     pb.GetUserType(),
		Lat:          float64(pb.GetLat()),
		Lng:          float64(pb.GetLng()),
		Thumbnail:    pb.GetThumbnail(),
	}
}

// dial sets up a gRPC connection with the microservice endpoint.
func (r PostRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}
