package models

import "time"

// ItemMetric represents metrics for items based on the protobuf specification
type ItemMetric struct {
	ItemID               string    `json:"item_id"`
	EntityType           string    `json:"entity_type"`
	LikesCount           int64     `json:"likes_count"`
	DislikesCount        int64     `json:"dislikes_count"`
	CommentsCount        int64     `json:"comments_count"`
	MessagesCount        int64     `json:"messages_count"`
	SharedCount          int64     `json:"shared_count"`
	AddedToWishlistCount int64     `json:"added_to_wishlist_count"`
	AddedToBasketCount   int64     `json:"added_to_basket_count"`
	VisitedCount         int64     `json:"visited_count"`
	ReportedCount        int64     `json:"reported_count"`
	FollowerCount        int64     `json:"follower_count"`
	ReviewCount          int64     `json:"review_count"`
	RatingCount          int64     `json:"rating_count"`
	VideosCount          int64     `json:"videos_count"`
	ImagesCount          int64     `json:"images_count"`
	Rating               int64     `json:"rating"`
	Category             string    `json:"category"`
	CategoryID           string    `json:"category_id"`
	CategorySlug         string    `json:"category_slug"`
	MediaCount           int64     `json:"media_count"`
	Price                int64     `json:"price"`
	Lat                  float32   `json:"lat"`
	Lng                  float32   `json:"lng"`
	CreatedAt            time.Time `json:"created_at,omitempty"`
	UpdatedAt            time.Time `json:"updated_at,omitempty"`
}

// GetItemMetricResponse represents the response for getting an item metric
type GetItemMetricResponse struct {
	Metric *ItemMetric `json:"metric"`
}

// GetItemsMetricResponse represents the response for getting multiple item metrics
type GetItemsMetricResponse struct {
	Metrics []*ItemMetric `json:"metrics"`
}

// UpdateItemMetricResponse represents the response for updating an item metric
type UpdateItemMetricResponse struct {
	ItemID string `json:"item_id"`
}

// ShareItemResponse represents the response for sharing an item
type ShareItemResponse struct {
	Success bool   `json:"success"`
	ItemID  string `json:"item_id,omitempty"`
}

// VisitItemResponse represents the response for visiting an item
type VisitItemResponse struct {
	Success bool   `json:"success"`
	ItemID  string `json:"item_id,omitempty"`
}

// Metric type constants
const (
	MetricTypeLikes     = "likes"
	MetricTypeDislikes  = "dislikes"
	MetricTypeComments  = "comments"
	MetricTypeMessages  = "messages"
	MetricTypeShares    = "shares"
	MetricTypeWishlist  = "wishlist"
	MetricTypeBasket    = "basket"
	MetricTypeVisits    = "visits"
	MetricTypeReports   = "reports"
	MetricTypeFollowers = "followers"
	MetricTypeReviews   = "reviews"
	MetricTypeRating    = "rating"
	MetricTypeVideos    = "videos"
	MetricTypeImages    = "images"
	MetricTypeMedia     = "media"
)

// Metric type action constants
const (
	MetricActionIncrement = "increment"
	MetricActionDecrement = "decrement"
	MetricActionSet       = "set"
	MetricActionReset     = "reset"
)
