// File: search/internal/grpc/post_repository.go
package grpc

import (
	"context"
	"fmt"
	"middleman/assistants/internal/domain"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"middleman/assistants/internal/models"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/posts/postspb"
)

// PostRepository calls the remote posts service (gRPC) as a fallback.
type PostRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.PostRepository = (*PostRepository)(nil)

// NewPostRepositoryWithAuth creates a new PostRepository with JWT authentication support
func NewPostRepository(endpoint string, authInstance *auth.Auth) PostRepository {
	return PostRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

func (r PostRepository) CreatePost(ctx context.Context, postID string, name string, description string, tags []string, status string, lat float64, lng float64, thumbnail string) error {
	log.Printf("[POST_GRPC] CreatePost called for postID=%s, name=%s", postID, name)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Printf("[POST_GRPC] Failed to connect to posts service: %v", err)
		return err
	}
	defer conn.Close()

	log.Printf("[POST_GRPC] Successfully connected to posts service, calling AddPost RPC...")

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.AddPost(ctx, &postspb.AddPostRequest{
		Name:        name,
		Description: description,
		Tags:        tags,
		Status:      status,
		Thumbnail:   thumbnail,
		Lat:         float32(lat),
		Lng:         float32(lng),
		// Note: postID, typeOfPost, userType, categoryId, categorySlug not set here
		// as they might be set by the service or in different contexts
	})
	if err != nil {
		log.Printf("[POST_GRPC] AddPost RPC failed: %v", err)
		return fmt.Errorf("AddPost RPC failed: %w", err)
	}

	log.Printf("[POST_GRPC] AddPost RPC successful, created post with ID: %s", resp.GetId())
	return nil
}

func (r PostRepository) UpdatePostDetails(ctx context.Context, postID, name, description string, tags []string, status, thumbnail string) error {
	log.Printf("[POST_GRPC] UpdatePostDetails: updating post ID=%s via gRPC", postID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Printf("[POST_GRPC] UpdatePostDetails: failed to dial gRPC: %v", err)
		return err
	}
	defer conn.Close()

	log.Printf("[POST_GRPC] Successfully connected to posts service, calling UpdatePost RPC...")

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.UpdatePost(ctx, &postspb.UpdatePostRequest{
		Id:          postID,
		Name:        name,
		Description: description,
		Tags:        tags,
		Status:      status,
		Thumbnail:   thumbnail,
		// Note: typeOfPost, categoryId, categorySlug can be added if needed
	})
	if err != nil {
		log.Printf("[POST_GRPC] UpdatePost RPC failed: %v", err)
		return fmt.Errorf("UpdatePost RPC failed: %w", err)
	}

	log.Printf("[POST_GRPC] UpdatePost RPC successful for post ID: %s", resp.GetId())
	return nil
}

func (r PostRepository) SearchPostsByTerm(ctx context.Context, searchTerm string) ([]*models.Post, error) {
	log.Printf("SearchPostsByTerm: searching posts with term=%s via gRPC", searchTerm)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("SearchPostsByTerm: failed to dial gRPC: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.GetPostsWithFilters(ctx, &postspb.GetPostsWithFiltersRequest{
		Name:     searchTerm,
		Page:     1,
		PageSize: 20,
	})
	if err != nil {
		return nil, fmt.Errorf("GetPostsWithFilters RPC failed: %w", err)
	}

	var results []*models.Post
	for _, p := range resp.GetPosts() {
		domainPost := r.postToDomain(p)
		if domainPost != nil {
			results = append(results, domainPost)
		}
	}
	return results, nil
}

func (r PostRepository) GetPostSuggestions(ctx context.Context, searchTerm string) ([]*models.Post, error) {
	log.Printf("GetPostSuggestions: suggesting posts for searchTerm=%s via gRPC", searchTerm)

	// Use SearchPostsByTerm as a fallback for suggestions, same pattern as other repositories
	return r.SearchPostsByTerm(ctx, searchTerm)
}

// GetPostByID retrieves a post by ID from the posts microservice (via gRPC).
func (r PostRepository) GetPostByID(ctx context.Context, postID string) (*models.Post, error) {
	log.Printf("GetPostByID: retrieving post with ID=%s via gRPC fallback", postID)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("GetPostByID: failed to dial gRPC: %v", err)
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

// SearchPostsAdvanced calls the remote microservice's GetPostsWithFilters (or similar).
func (r PostRepository) SearchPostsAdvanced(
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
	log.Printf("[POST_GRPC] SearchPostsAdvanced called with: name='%s', status='%s', userType='%s', page=%d, pageSize=%d", name, status, userType, page, pageSize)

	// Check if no meaningful filters provided
	noTextFilters := name == "" && description == "" && len(tags) == 0 && status == "" && userType == "" && thumbnail == ""
	noLocationFilters := lat == 0 && lng == 0 && radius == 0

	if noTextFilters && noLocationFilters {
		log.Printf("[POST_GRPC] No filters provided, using GetPosts fallback")
		return r.GetAllPosts(ctx, page, pageSize, sortBy, sortOrder)
	}

	log.Printf("[POST_GRPC] Using filtered search with GetPostsWithFilters")
	log.Printf("[POST_GRPC] Attempting to dial posts service at endpoint: %s", r.endpoint)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[POST_GRPC] Failed to connect to posts service: %v", err)
		return nil, err
	}
	defer conn.Close()

	log.Printf("[POST_GRPC] Successfully connected to posts service at %s", r.endpoint)
	log.Printf("[POST_GRPC] Calling GetPostsWithFilters RPC...")

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.GetPostsWithFilters(ctx, &postspb.GetPostsWithFiltersRequest{
		Name:        name,
		Description: description,
		Tags:        tags,
		Status:      status,
		UserType:    userType,
		Thumbnail:   thumbnail,
		Lat:         float32(lat),
		Lng:         float32(lng),
		Radius:      float32(radius),
		Page:        page,
		PageSize:    pageSize,
		SortBy:      sortBy,
		SortOrder:   sortOrder,
	})
	if err != nil {
		log.Printf("[POST_GRPC] GetPostsWithFilters RPC failed: %v", err)
		return nil, fmt.Errorf("GetPostsWithFilters RPC failed: %w", err)
	}

	log.Printf("[POST_GRPC] GetPostsWithFilters RPC successful, received %d posts", len(resp.GetPosts()))

	var results []*models.Post
	for i, p := range resp.GetPosts() {
		log.Printf("[POST_GRPC] Converting post %d: ID=%s, Name=%s", i, p.GetId(), p.GetName())
		domainPost := r.postToDomain(p)
		if domainPost != nil {
			results = append(results, domainPost)
		}
	}

	log.Printf("[POST_GRPC] GetPostsWithFilters returning %d converted posts", len(results))
	return results, nil
}

// GetPostsByCategorySlug calls the remote microservice's GetPostsByCategorySlug.
func (r PostRepository) GetPostsByCategorySlug(
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

// GetPostsByCategoryID calls the remote microservice's GetPostsByCategory.
func (r PostRepository) GetPostsByCategoryID(
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

// Update calls UpdatePost in the remote gRPC service (or partial, etc.).
func (r PostRepository) Update(ctx context.Context, postID string, price int64) error {
	conn, err := r.dialWithAuth(ctx)
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

// DeletePost calls RemovePost in the remote gRPC service.
func (r PostRepository) DeletePost(ctx context.Context, postID string) error {
	conn, err := r.dialWithAuth(ctx)
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
		PostType:     pb.GetTypeOfPost(),
		CategoryID:   pb.GetCategoryId(),
		CategorySlug: pb.GetCategorySlug(),
		Tags:         pb.GetTags(),
		Status:       pb.GetStatus(),
		UserType:     pb.GetUserType(),
		Lat:          float64(pb.GetLat()),
		Lng:          float64(pb.GetLng()),
		Thumbnail:    pb.GetThumbnail(),
		EntityType:   models.PostType,
	}
}

// dial sets up a gRPC connection with the microservice endpoint.
// dial sets up a gRPC connection with the microservice endpoint
func (r PostRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r PostRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}

// Additional methods available in gRPC service

// GetAllPosts retrieves all posts with pagination
func (r PostRepository) GetAllPosts(ctx context.Context, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error) {
	log.Printf("[POST_GRPC] GetAllPosts called with page=%d, pageSize=%d", page, pageSize)

	conn, err := r.dial(ctx)
	if err != nil {
		log.Printf("[POST_GRPC] Failed to connect to posts service: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.GetPosts(ctx, &postspb.GetPostsRequest{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		log.Printf("[POST_GRPC] GetPosts RPC failed: %v", err)
		return nil, fmt.Errorf("GetPosts RPC failed: %w", err)
	}

	log.Printf("[POST_GRPC] GetPosts RPC successful, received %d posts", len(resp.GetPosts()))

	var results []*models.Post
	for i, p := range resp.GetPosts() {
		log.Printf("[POST_GRPC] Converting post %d: ID=%s, Name=%s", i, p.GetId(), p.GetName())
		if post := r.postToDomain(p); post != nil {
			results = append(results, post)
		}
	}

	log.Printf("[POST_GRPC] GetPosts returning %d converted posts", len(results))
	return results, nil
}

// GetPostsByUserID retrieves posts for a specific user
func (r PostRepository) GetPostsByUserID(ctx context.Context, userID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error) {
	log.Printf("[POST_GRPC] GetPostsByUserID called for userID=%s", userID)

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.GetUserPosts(ctx, &postspb.GetUserPostsRequest{
		UserId:    userID,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetUserPosts RPC failed: %w", err)
	}

	var results []*models.Post
	for _, p := range resp.GetPosts() {
		if post := r.postToDomain(p); post != nil {
			results = append(results, post)
		}
	}
	return results, nil
}

// GetPublicPostCatalog retrieves user's public catalog of posts
func (r PostRepository) GetPublicPostCatalog(ctx context.Context, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error) {
	log.Printf("[POST_GRPC] GetPublicPostCatalog called for page=%d, pageSize=%d", page, pageSize)

	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	resp, err := client.GetPublicCatalog(ctx, &postspb.GetPublicCatalogRequest{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("GetPublicCatalog RPC failed: %w", err)
	}

	var results []*models.Post
	for _, p := range resp.GetPosts() {
		if post := r.postToDomain(p); post != nil {
			results = append(results, post)
		}
	}
	return results, nil
}

// ArchiveUserPost archives a post
func (r PostRepository) ArchiveUserPost(ctx context.Context, postID string) error {
	log.Printf("[POST_GRPC] ArchiveUserPost called for postID=%s", postID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	_, err = client.ArchivePost(ctx, &postspb.ArchivePostRequest{
		PostId: postID,
	})
	return err
}

// AddThumbnailToPost adds a thumbnail to a post
func (r PostRepository) AddThumbnailToPost(ctx context.Context, postID string, thumbnail string) error {
	log.Printf("[POST_GRPC] AddThumbnailToPost called for postID=%s", postID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	_, err = client.AddPostThumbnail(ctx, &postspb.AddPostThumbnailRequest{
		PostId:    postID,
		Thumbnail: thumbnail,
	})
	return err
}

// UpdatePostThumbnail updates a post's thumbnail
func (r PostRepository) UpdatePostThumbnail(ctx context.Context, postID string, thumbnail string) error {
	log.Printf("[POST_GRPC] UpdatePostThumbnail called for postID=%s", postID)

	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := postspb.NewPostsServiceClient(conn)
	_, err = client.UpdatePostThumbnail(ctx, &postspb.UpdatePostThumbnailRequest{
		PostId:    postID,
		Thumbnail: thumbnail,
	})
	return err
}

// Additional query methods needed by tool service

// GetPostsByLocation retrieves posts by location with radius filtering
func (r PostRepository) GetPostsByLocation(ctx context.Context, lat, lng float32, radius, page, pageSize int64) ([]*models.Post, int64, error) {
	// Use the existing SearchPostsAdvanced method with location parameters
	posts, err := r.SearchPostsAdvanced(ctx, "", "", []string{}, "", "", "", 0, pageSize, float64(lat), float64(lng), radius, page, pageSize, "created_at", "desc")
	if err != nil {
		return nil, 0, err
	}

	return posts, int64(len(posts)), nil
}

// GetPostsByTags retrieves posts filtered by tags
func (r PostRepository) GetPostsByTags(ctx context.Context, tags []string, page, pageSize int64) ([]*models.Post, int64, error) {
	posts, err := r.SearchPostsAdvanced(ctx, "", "", tags, "", "", "", 0, pageSize, 0, 0, 0, page, pageSize, "created_at", "desc")
	if err != nil {
		return nil, 0, err
	}

	return posts, int64(len(posts)), nil
}

// GetPopularPosts retrieves popular posts (sorted by engagement metrics)
func (r PostRepository) GetPopularPosts(ctx context.Context, page, pageSize int64) ([]*models.Post, int64, error) {
	posts, err := r.GetAllPosts(ctx, page, pageSize, "popularity", "desc")
	if err != nil {
		return nil, 0, err
	}

	return posts, int64(len(posts)), nil
}

// GetRecentPosts retrieves recent posts (sorted by creation date)
func (r PostRepository) GetRecentPosts(ctx context.Context, page, pageSize int64) ([]*models.Post, int64, error) {
	posts, err := r.GetAllPosts(ctx, page, pageSize, "created_at", "desc")
	if err != nil {
		return nil, 0, err
	}

	return posts, int64(len(posts)), nil
}

// GetTrendingPosts retrieves trending posts (sorted by recent engagement)
func (r PostRepository) GetTrendingPosts(ctx context.Context, page, pageSize int64) ([]*models.Post, int64, error) {
	posts, err := r.GetAllPosts(ctx, page, pageSize, "trending", "desc")
	if err != nil {
		return nil, 0, err
	}

	return posts, int64(len(posts)), nil
}

// AddLikeToPost likes a post for a user
func (r PostRepository) AddLikeToPost(ctx context.Context, postID string) error {
	// For now, return nil as this would need to be implemented in the posts service
	// This is a placeholder implementation
	return nil
}

// RemoveLikeFromPost unlikes a post for a user
func (r PostRepository) RemoveLikeFromPost(ctx context.Context, postID string) error {
	// For now, return nil as this would need to be implemented in the posts service
	// This is a placeholder implementation
	return nil
}

// SharePostWithUsers shares a post for a user
func (r PostRepository) SharePostWithUsers(ctx context.Context, postID string) error {
	// For now, return nil as this would need to be implemented in the posts service
	// This is a placeholder implementation
	return nil
}

// ReportInappropriatePost reports a post with a reason
func (r PostRepository) ReportInappropriatePost(ctx context.Context, postID, reason string) error {
	// For now, return nil as this would need to be implemented in the posts service
	// This is a placeholder implementation
	return nil
}
