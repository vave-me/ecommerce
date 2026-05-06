package domain

import "context"

// CatalogPost is the minimal representation of a post in your catalog.
type CatalogPost struct {
	ID           string
	Name         string
	Description  string
	TypeOfPost   TypeOfPost
	UserID       string
	UserType     UserType
	CategoryID   string
	CategorySlug string
	Tags         []string
	Status       PostStatus
	Thumbnail    string
	Lat          float64
	Lng          float64
}

// CatalogRepository defines the essential CRUD operations.
// Additional methods (e.g. pagination) can be added as needed.
type CatalogRepository interface {
	AddPost(ctx context.Context, id, name, description string, typeOfPost TypeOfPost, userId string, userType UserType, categoryID, categorySlug string, tags []string, status PostStatus, thumbnail string, lat, lng float64) error
	UpdatePost(ctx context.Context, postID, name, description string, typeOfPost TypeOfPost, userID string, tags []string, status PostStatus, thumbnail string) error
	ArchivePost(ctx context.Context, postID, userID string) error
	RemovePost(ctx context.Context, postID, userID string) error
	Find(ctx context.Context, postID string) (*CatalogPost, error)
	GetPosts(ctx context.Context, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogPost, int64, error)
	GetUserPosts(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogPost, int64, error)
	GetPublicCatalog(ctx context.Context, userID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogPost, int64, error)
	FindByLocation(ctx context.Context, lat, lng float64, radiusMeters float64, limit int) ([]*CatalogPost, error)
	GetPostsByCategory(ctx context.Context, categoryID string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogPost, int64, error)
	GetPostsByCategorySlug(ctx context.Context, categorySlug string, page, pageSize int64, sortBy, sortOrder string) ([]*CatalogPost, int64, error)
	GetPostsWithFilters(ctx context.Context,
		name, description string,
		tags []string,
		offset, limit int64,
		lat, lng float64,
		radius, page, pageSize int64,
		sortBy, sortOrder string,
	) ([]*CatalogPost, int64, error)
	UpdatePostThumbnail(ctx context.Context, postID string, thumbnail string) error
}
