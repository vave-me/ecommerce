package models

import "time"

// UserMetric represents metrics for users based on the protobuf specification
type UserMetric struct {
	UserID               string    `json:"user_id"`
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
	MediaAddedCount      int64     `json:"media_added_count"`
	CommentAddedCount    int64     `json:"comment_added_count"`
	LikedCount           int64     `json:"liked_count"`
	DislikedCount        int64     `json:"disliked_count"`
	ProductsAddedCount   int64     `json:"products_added_count"`
	VideosAddedCount     int64     `json:"videos_added_count"`
	ImagesAddedCount     int64     `json:"images_added_count"`
	SeriesAddedCount     int64     `json:"series_added_count"`
	JobsAddedCount       int64     `json:"jobs_added_count"`
	PostsAddedCount      int64     `json:"posts_added_count"`
	VehiclesAddedCount   int64     `json:"vehicles_added_count"`
	PropertiesAddedCount int64     `json:"properties_added_count"`
	CreatedAt            time.Time `json:"created_at,omitempty"`
	UpdatedAt            time.Time `json:"updated_at,omitempty"`
}

// GetUserMetricResponse represents the response for getting a user metric
type GetUserMetricResponse struct {
	Metric *UserMetric `json:"metric"`
}

// UpdateUserMetricResponse represents the response for updating a user metric
type UpdateUserMetricResponse struct {
	UserID string `json:"user_id"`
}
