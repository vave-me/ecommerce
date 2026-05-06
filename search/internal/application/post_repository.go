package application

import (
	"context"
	"middleman/search/internal/models"
)

type PostRepository interface {
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
	GetCatalog(ctx context.Context, userId string) ([]*models.Post, error)
}

type PostCacheRepository interface {
	PostRepository
	Add(
		ctx context.Context,
		postID, name, description string,
		userID string,
		tags []string,
		status string,
		lat, lng float64,
		thumbnail string,
		entityType models.EntityType,
	) error
	UpdatePost(ctx context.Context, postID, name, description string, tags []string, status, thumbnail string) error
	SearchWithTerm(ctx context.Context, name string) ([]*models.Post, error)
	SuggestPosts(ctx context.Context, name string) ([]*models.Post, error)
	UpdateThumbnail(ctx context.Context, postID string, thumbnail string) error
	FindBatch(ctx context.Context, postIDs []string) (map[string]*models.Post, error)
}
