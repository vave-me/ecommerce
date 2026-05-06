package domain

import (
	"context"
)

// StreamRepository interface for stream persistence
type StreamRepository interface {
	// Find retrieves a stream by ID
	Find(ctx context.Context, streamID string) (*Stream, error)
	
	// FindByUserAccess retrieves streams accessible by a specific user
	FindByUserAccess(ctx context.Context, userID string) ([]*Stream, error)
	
	// FindByType retrieves streams by type
	FindByType(ctx context.Context, streamType StreamType) ([]*Stream, error)
	
	// FindByStatus retrieves streams by status
	FindByStatus(ctx context.Context, status StreamStatus) ([]*Stream, error)
	
	// FindByGenre retrieves streams by genre
	FindByGenre(ctx context.Context, genre string) ([]*Stream, error)
	
	// FindByAccessType retrieves streams by access type
	FindByAccessType(ctx context.Context, accessType AccessType) ([]*Stream, error)
	
	// FindBySeries retrieves all episodes for a series
	FindBySeries(ctx context.Context, seriesID string) ([]*Stream, error)
	
	// FindBySeriesAndSeason retrieves episodes for a specific season
	FindBySeriesAndSeason(ctx context.Context, seriesID string, seasonNumber int) ([]*Stream, error)
	
	// Search performs a full-text search on streams
	Search(ctx context.Context, query string, filters StreamFilters) ([]*Stream, error)
	
	// Save persists a stream
	Save(ctx context.Context, stream *Stream) error
	
	// Update updates an existing stream
	Update(ctx context.Context, stream *Stream) error
	
	// Delete removes a stream
	Delete(ctx context.Context, streamID string) error
	
	// GetPopular retrieves popular streams based on view count
	GetPopular(ctx context.Context, limit int) ([]*Stream, error)
	
	// GetRecommended retrieves recommended streams for a user
	GetRecommended(ctx context.Context, userID string, limit int) ([]*Stream, error)
	
	// GetRecentlyAdded retrieves recently added streams
	GetRecentlyAdded(ctx context.Context, limit int) ([]*Stream, error)
	
	// GetContinueWatching retrieves streams the user has started but not finished
	GetContinueWatching(ctx context.Context, userID string) ([]*Stream, error)
}

// SeriesRepository interface for series persistence
type SeriesRepository interface {
	// Find retrieves a series by ID
	Find(ctx context.Context, seriesID string) (*Series, error)
	
	// FindAll retrieves all series
	FindAll(ctx context.Context) ([]*Series, error)
	
	// FindByGenre retrieves series by genre
	FindByGenre(ctx context.Context, genre string) ([]*Series, error)
	
	// FindByStudio retrieves series by studio
	FindByStudio(ctx context.Context, studio string) ([]*Series, error)
	
	// Search performs a full-text search on series
	Search(ctx context.Context, query string) ([]*Series, error)
	
	// Save persists a series
	Save(ctx context.Context, series *Series) error
	
	// Update updates an existing series
	Update(ctx context.Context, series *Series) error
	
	// Delete removes a series
	Delete(ctx context.Context, seriesID string) error
	
	// GetPopular retrieves popular series
	GetPopular(ctx context.Context, limit int) ([]*Series, error)
}

// StreamFilters for advanced stream searching
type StreamFilters struct {
	StreamType     StreamType
	MinDuration    int64
	MaxDuration    int64
	ContentRating  []ContentRating
	Genre          []string
	Language       string
	Country        string
	Studio         string
	AccessType     AccessType
	MinRating      float64
	ReleasedAfter  string
	ReleasedBefore string
	SortBy         string // "relevance", "date", "rating", "views"
	SortOrder      string // "asc", "desc"
}

// StreamCatalogRepository interface for read-optimized queries
type StreamCatalogRepository interface {
	// GetCatalog retrieves the full catalog with filters
	GetCatalog(ctx context.Context, filters CatalogFilters) (*StreamCatalog, error)
	
	// GetUserCatalog retrieves catalog personalized for a user
	GetUserCatalog(ctx context.Context, userID string, filters CatalogFilters) (*StreamCatalog, error)
	
	// GetCategories retrieves all available categories/genres
	GetCategories(ctx context.Context) ([]Category, error)
	
	// GetStreamDetails retrieves detailed information about a stream
	GetStreamDetails(ctx context.Context, streamID string, userID string) (*StreamDetails, error)
}

// CatalogFilters for catalog queries
type CatalogFilters struct {
	Page          int
	PageSize      int
	StreamType    StreamType
	Genre         string
	AccessType    AccessType
	SortBy        string
	SearchQuery   string
}

// StreamCatalog represents a catalog response
type StreamCatalog struct {
	Streams      []*StreamSummary
	TotalCount   int
	Page         int
	PageSize     int
	HasMore      bool
	Categories   []Category
	FeaturedList []string // Stream IDs
}

// StreamSummary represents a summary view of a stream
type StreamSummary struct {
	StreamID      string
	Title         string
	ThumbnailURL  string
	Duration      int64
	ContentRating ContentRating
	AccessType    AccessType
	Rating        float64
	ViewCount     int64
	ReleaseYear   int
	IsNew         bool
	HasAccess     bool // For user-specific catalogs
}

// StreamDetails represents detailed stream information
type StreamDetails struct {
	*Stream
	HasAccess       bool
	WatchProgress   int64
	IsInWatchlist   bool
	UserRating      int
	RelatedStreams  []*StreamSummary
	NextEpisode     *StreamSummary // For series
	PreviousEpisode *StreamSummary // For series
}

// Category represents a content category
type Category struct {
	ID          string
	Name        string
	Slug        string
	StreamCount int
}