package models

type UserMetric struct {
	ID                   string
	EntityType           string
	LikesCount           int64
	DislikesCount        int64
	CommentsCount        int64
	MessagesCount        int64
	SharedCount          int64
	AddedToWishlistCount int64
	AddedToBasketCount   int64
	VisitedCount         int64
	ReportedCount        int64
	FollowerCount        int64
	ReviewsCount         int64
	RatingCount          int64
	VideosCount          int64
	ImagesCount          int64
	Rating               int64
	Review               int64
	Category             string
	CategoryID           string
	CategorySlug         string
	MediaAddedCount      int64
	CommentAddedCount    int64
	LikedAddedCount      int64
	ProductsAddedCount   int64
	VideosAddedCount     int64
	ServicesAddedCount   int64
	JobsAddedCount       int64
	PostsAddedCount      int64
	VehiclesAddedCount   int64
	PropertiesAddedCount int64
	CreatedAt            string
	UpdatedAt            string
}
