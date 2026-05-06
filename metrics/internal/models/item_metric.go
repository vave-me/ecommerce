package models

type ItemMetric struct {
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
	Price                int64
	Lat                  float64
	Lng                  float64
	CreatedAt            string
	UpdatedAt            string
}
