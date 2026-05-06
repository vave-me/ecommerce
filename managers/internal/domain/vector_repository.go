package domain

import (
	"context"
	"time"
)

// VectorRepository defines the interface for vector database operations
// This abstraction allows graceful degradation when vector service is unavailable
type VectorRepository interface {
	// SearchSimilar finds entities similar to a given entity using vector similarity
	// Returns empty results if vector service is unavailable
	SearchSimilar(ctx context.Context, entityID string, entityType string, options SearchSimilarOptions) (*VectorSearchResults, error)

	// SearchByVector performs semantic search using raw vector embeddings
	// Returns empty results if vector service is unavailable
	SearchByVector(ctx context.Context, vector []float32, options VectorSearchOptions) (*VectorSearchResults, error)

	// GetRecommendations provides personalized recommendations based on user vector
	// Returns empty results if vector service is unavailable
	GetRecommendations(ctx context.Context, userVector []float32, options RecommendationOptions) (*VectorSearchResults, error)

	// GetEntityContext retrieves contextual information about an entity
	// Returns basic entity info if vector service is unavailable
	GetEntityContext(ctx context.Context, entityID string, entityType string, options ContextOptions) (*EntityContext, error)

	// Health checks if vector service is available
	Health(ctx context.Context) bool
}

// SearchSimilarOptions configures similar entity search
type SearchSimilarOptions struct {
	TopK              int64          `json:"top_k"`
	TargetEntityTypes []string       `json:"target_entity_types"`
	ExcludeSelf       bool           `json:"exclude_self"`
	ScoreThreshold    float32        `json:"score_threshold"`
	Filters           *VectorFilters `json:"filters"`
}

// VectorSearchOptions configures vector search operations
type VectorSearchOptions struct {
	EntityTypes    []string       `json:"entity_types"`
	TopK           int64          `json:"top_k"`
	ScoreThreshold float32        `json:"score_threshold"`
	WithVectors    bool           `json:"with_vectors"`
	Filters        *VectorFilters `json:"filters"`
}

// RecommendationOptions configures recommendation requests
type RecommendationOptions struct {
	EntityTypes    []string       `json:"entity_types"`
	TopK           int64          `json:"top_k"`
	DiversityLevel float32        `json:"diversity_level"`
	ExcludeIDs     []string       `json:"exclude_ids"`
	Filters        *VectorFilters `json:"filters"`
}

// ContextOptions configures entity context requests
type ContextOptions struct {
	IncludeSimilar bool  `json:"include_similar"`
	ContextSize    int64 `json:"context_size"`
}

// VectorFilters represents filtering criteria for vector operations
type VectorFilters struct {
	EntityTypes     []string          `json:"entity_types"`
	PriceRange      *PriceRange       `json:"price_range"`
	GeoFilter       *GeoFilter        `json:"geo_filter"`
	Statuses        []string          `json:"statuses"`
	Categories      []string          `json:"categories"`
	NegotiableOnly  bool              `json:"negotiable_only"`
	TimeRange       *TimeRange        `json:"time_range"`
	UserType        string            `json:"user_type"`
	MetadataFilters map[string]string `json:"metadata_filters"`
}

// PriceRange for price-based filtering
type PriceRange struct {
	MinPrice int64 `json:"min_price"`
	MaxPrice int64 `json:"max_price"`
}

// GeoFilter for location-based filtering
type GeoFilter struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	RadiusKm  float64 `json:"radius_km"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
}

// TimeRange for time-based filtering
type TimeRange struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// VectorSearchResults represents search results from vector operations
type VectorSearchResults struct {
	Results      []*VectorSearchResult `json:"results"`
	TotalFound   int64                 `json:"total_found"`
	MaxScore     float32               `json:"max_score"`
	MinScore     float32               `json:"min_score"`
	QueryTime    time.Duration         `json:"query_time"`
	SourceFailed bool                  `json:"source_failed"` // Indicates if vector service failed
}

// VectorSearchResult represents a single search result
type VectorSearchResult struct {
	Entity        *VectorEntity      `json:"entity"`
	Score         float32            `json:"score"`
	Vector        []float32          `json:"vector,omitempty"`
	FeatureScores map[string]float32 `json:"feature_scores"`
	Rank          int64              `json:"rank"`
}

// VectorEntity represents an entity in vector search results
type VectorEntity struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Location    *EntityLocation   `json:"location,omitempty"`
	Pricing     *EntityPricing    `json:"pricing,omitempty"`
	Tags        []string          `json:"tags"`
	Status      string            `json:"status"`
	Thumbnail   string            `json:"thumbnail"`
}

// EntityLocation represents geographic information
type EntityLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
}

// EntityPricing represents pricing information
type EntityPricing struct {
	BasePrice  int64  `json:"base_price"`
	SalePrice  int64  `json:"sale_price"`
	Currency   string `json:"currency"`
	Negotiable bool   `json:"negotiable"`
}

// EntityContext provides contextual information about an entity
type EntityContext struct {
	Entity  *VectorEntity         `json:"entity"`
	Similar []*VectorSearchResult `json:"similar"`
	Metrics *EntityMetrics        `json:"metrics"`
	Failed  bool                  `json:"failed"` // Indicates if vector service failed
}

// EntityMetrics represents entity statistics
type EntityMetrics struct {
	ViewCount        int64            `json:"view_count"`
	InteractionCount int64            `json:"interaction_count"`
	PopularityScore  float32          `json:"popularity_score"`
	LastUpdated      time.Time        `json:"last_updated"`
	CustomMetrics    map[string]int64 `json:"custom_metrics"`
}
