package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type PostRepository interface {
	Add(
		ctx context.Context,
		postID, name, description string,
		tags []string,
		status string,
		lat, lng float64,
		thumbnail string,
	) error
	UpdatePost(ctx context.Context, postID, name, description string, tags []string, status, thumbnail string) error
	SearchWithTerm(ctx context.Context, name string) ([]*models.Post, error)
	SuggestPosts(ctx context.Context, name string) ([]*models.Post, error)
	UpdateThumbnail(ctx context.Context, postID string, thumbnail string) error
	Find(ctx context.Context, postID string) (*models.Post, error)
	Remove(ctx context.Context, postID string) error
	SearchPostsWithFilters(
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
	SearchPostsWithCategorySlug(ctx context.Context, categorySlug string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error)
	SearchPostsWithCategory(ctx context.Context, categoryId string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error)

	// Additional methods available in gRPC service
	GetPosts(ctx context.Context, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error)
	GetUserPosts(ctx context.Context, userID string, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error)
	GetPublicCatalog(ctx context.Context, page int64, pageSize int64, sortBy string, sortOrder string) ([]*models.Post, error)
	ArchivePost(ctx context.Context, postID string) error
	AddPostThumbnail(ctx context.Context, postID string, thumbnail string) error
	UpdatePostThumbnail(ctx context.Context, postID string, thumbnail string) error

	// Additional query methods needed by tool service
	GetPostsByLocation(ctx context.Context, lat, lng float32, radius, page, pageSize int64) ([]*models.Post, int64, error)
	GetPostsByTags(ctx context.Context, tags []string, page, pageSize int64) ([]*models.Post, int64, error)
	GetPopularPosts(ctx context.Context, page, pageSize int64) ([]*models.Post, int64, error)
	GetRecentPosts(ctx context.Context, page, pageSize int64) ([]*models.Post, int64, error)
	GetTrendingPosts(ctx context.Context, page, pageSize int64) ([]*models.Post, int64, error)
	LikePost(ctx context.Context, postID string) error
	UnlikePost(ctx context.Context, postID string) error
	SharePost(ctx context.Context, postID string) error
	ReportPost(ctx context.Context, postID, reason string) error
}
