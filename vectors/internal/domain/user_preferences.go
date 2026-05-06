package domain

import (
	"time"
)

// UserBehaviorEvent represents a user interaction with an entity
type UserBehaviorEvent struct {
	UserID       string                 `json:"user_id"`
	EntityID     string                 `json:"entity_id"`
	EntityType   string                 `json:"entity_type"`
	EventType    UserEventType          `json:"event_type"`
	Timestamp    time.Time              `json:"timestamp"`
	Duration     time.Duration          `json:"duration,omitempty"`      // For view duration
	Context      map[string]interface{} `json:"context,omitempty"`       // Additional context
	Score        float32                `json:"score,omitempty"`         // Implicit score based on action
	SessionID    string                 `json:"session_id,omitempty"`    // For session grouping
	DeviceType   string                 `json:"device_type,omitempty"`   // mobile, desktop, tablet
	Source       string                 `json:"source,omitempty"`        // search, recommendation, browse
}

// UserEventType defines types of user interactions
type UserEventType string

const (
	EventView          UserEventType = "view"
	EventClick         UserEventType = "click"
	EventAddToCart     UserEventType = "add_to_cart"
	EventRemoveFromCart UserEventType = "remove_from_cart"
	EventPurchase      UserEventType = "purchase"
	EventAddToWishlist UserEventType = "add_to_wishlist"
	EventShare         UserEventType = "share"
	EventComment       UserEventType = "comment"
	EventLike          UserEventType = "like"
	EventDislike       UserEventType = "dislike"
	EventSearch        UserEventType = "search"
	EventFilter        UserEventType = "filter"
)

// EventWeight defines the importance of different events for preference learning
var EventWeights = map[UserEventType]float32{
	EventView:           0.1,
	EventClick:          0.2,
	EventAddToCart:      0.5,
	EventRemoveFromCart: -0.3,
	EventPurchase:       1.0,
	EventAddToWishlist:  0.6,
	EventShare:          0.7,
	EventComment:        0.8,
	EventLike:           0.4,
	EventDislike:        -0.5,
	EventSearch:         0.3,
	EventFilter:         0.2,
}

// UserPreferenceVector represents a user's learned preferences
type UserPreferenceVector struct {
	UserID            string                      `json:"user_id"`
	Vector            []float32                   `json:"vector"`
	LastUpdated       time.Time                   `json:"last_updated"`
	InteractionCount  int64                       `json:"interaction_count"`
	CategoryWeights   map[string]float32          `json:"category_weights"`   // Category preferences
	PriceRangePrefs   *PriceRangePreference       `json:"price_range_prefs"`  // Price preferences
	BrandAffinities   map[string]float32          `json:"brand_affinities"`   // Brand preferences
	LocationPrefs     *LocationPreference         `json:"location_prefs"`     // Location preferences
	TemporalPatterns  map[string]float32          `json:"temporal_patterns"`  // Time-based patterns
	EntityTypePrefs   map[string]float32          `json:"entity_type_prefs"`  // Product vs Post vs Service
	QualityMetrics    *UserQualityMetrics         `json:"quality_metrics"`    // User engagement quality
}

// PriceRangePreference captures user's price sensitivity
type PriceRangePreference struct {
	PreferredMin      int64   `json:"preferred_min"`
	PreferredMax      int64   `json:"preferred_max"`
	AverageSpend      float64 `json:"average_spend"`
	PriceSensitivity  float32 `json:"price_sensitivity"`  // 0-1, higher = more sensitive
	LuxuryAffinity    float32 `json:"luxury_affinity"`    // 0-1, higher = prefers premium
	BargainHunting    float32 `json:"bargain_hunting"`    // 0-1, higher = seeks deals
}

// LocationPreference captures user's geographic preferences
type LocationPreference struct {
	PreferredRadius   float64            `json:"preferred_radius_km"`
	FrequentLocations []GeoPoint         `json:"frequent_locations"`
	LocationFlexibility float32          `json:"location_flexibility"` // 0-1, higher = more flexible
	RemoteAffinity    float32            `json:"remote_affinity"`      // For services
}

// GeoPoint represents a geographic location
type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Weight    float32 `json:"weight"`
}

// UserQualityMetrics tracks user engagement quality
type UserQualityMetrics struct {
	AverageViewDuration   time.Duration `json:"avg_view_duration"`
	PurchaseConversionRate float32      `json:"purchase_conversion_rate"`
	ReturnRate            float32       `json:"return_rate"`
	EngagementScore       float32       `json:"engagement_score"`
	TrustScore            float32       `json:"trust_score"`
}

// UserPreferenceUpdate represents an incremental update to user preferences
type UserPreferenceUpdate struct {
	UserID           string             `json:"user_id"`
	BehaviorEvents   []UserBehaviorEvent `json:"behavior_events"`
	UpdateStrategy   UpdateStrategy      `json:"update_strategy"`
	DecayFactor      float32            `json:"decay_factor"`      // For time-based decay
	LearningRate     float32            `json:"learning_rate"`     // For gradient updates
}

// UpdateStrategy defines how to update user preferences
type UpdateStrategy string

const (
	StrategyIncremental UpdateStrategy = "incremental"  // Add new signals
	StrategyDecay       UpdateStrategy = "decay"        // Apply time decay
	StrategyRecompute   UpdateStrategy = "recompute"    // Full recomputation
	StrategyHybrid      UpdateStrategy = "hybrid"       // Combine strategies
)

// UserSegment represents a cluster of similar users
type UserSegment struct {
	SegmentID       string    `json:"segment_id"`
	Name            string    `json:"name"`
	CentroidVector  []float32 `json:"centroid_vector"`
	UserCount       int64     `json:"user_count"`
	Characteristics map[string]interface{} `json:"characteristics"`
}