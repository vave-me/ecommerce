package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type PostRepository interface {
	// Basic CRUD operations
	CreatePost(
		ctx context.Context,
		postID, name, description string,
		tags []string,
		status string,
		lat, lng float64,
		thumbnail string,
	) error
	GetPostByID(ctx context.Context, postID string) (*models.Post, error)
	UpdatePostDetails(ctx context.Context, postID, name, description string, tags []string, status, thumbnail string) error
	DeletePost(ctx context.Context, postID string) error
	
	// Search and discovery
	SearchPostsByTerm(ctx context.Context, searchTerm string) ([]*models.Post, error)
	GetPostSuggestions(ctx context.Context, searchTerm string) ([]*models.Post, error)
	SearchPostsAdvanced(
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
	) ([]*models.Post, error)
	GetPostsByCategorySlug(ctx context.Context, categorySlug string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error)
	GetPostsByCategoryID(ctx context.Context, categoryId string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error)

	// Post retrieval
	GetAllPosts(ctx context.Context, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error)
	GetPostsByUserID(ctx context.Context, userID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error)
	GetPublicPostCatalog(ctx context.Context, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error)
	
	// Post management
	ArchiveUserPost(ctx context.Context, postID string) error
	UpdatePostThumbnail(ctx context.Context, postID string, thumbnail string) error
	AddThumbnailToPost(ctx context.Context, postID string, thumbnail string) error

	// Location and tag based queries
	GetPostsByLocation(ctx context.Context, lat, lng float32, radius, page, pageSize int64) ([]*models.Post, int64, error)
	GetPostsByTags(ctx context.Context, tags []string, page, pageSize int64) ([]*models.Post, int64, error)
	
	// Trending and popular posts
	GetPopularPosts(ctx context.Context, page, pageSize int64) ([]*models.Post, int64, error)
	GetRecentPosts(ctx context.Context, page, pageSize int64) ([]*models.Post, int64, error)
	GetTrendingPosts(ctx context.Context, page, pageSize int64) ([]*models.Post, int64, error)
	
	// Post interactions
	AddLikeToPost(ctx context.Context, postID string) error
	RemoveLikeFromPost(ctx context.Context, postID string) error
	SharePostWithUsers(ctx context.Context, postID string) error
	ReportInappropriatePost(ctx context.Context, postID, reason string) error
}
