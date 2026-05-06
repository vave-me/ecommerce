package models

type MetricTypeCount string

const (
	MetricTypeCountComment   MetricTypeCount = "comment"
	MetricTypeCountMessage   MetricTypeCount = "message"
	MetricTypeCountLike      MetricTypeCount = "like"
	MetricTypeCountDislike   MetricTypeCount = "dislike"
	MetricTypeCountWishlist  MetricTypeCount = "wishlist"
	MetricTypeCountBasket    MetricTypeCount = "basket"
	MetricTypeCountFollowing MetricTypeCount = "following"
	MetricTypeCountReview    MetricTypeCount = "review"
	MetricTypeCountShare     MetricTypeCount = "share"
	MetricTypeCountView      MetricTypeCount = "view"
	MetricTypeCountVisit     MetricTypeCount = "visit"
	MetricTypeCountReport    MetricTypeCount = "report"
	MetricTypeCountFollow    MetricTypeCount = "follow"

	// User specific metric types
	MetricTypeUserMediaAdd    MetricTypeCount = "media_add"
	MetricTypeUserCommentAdd  MetricTypeCount = "comment_add"
	MetricTypeUserLikeAdd     MetricTypeCount = "like_add"
	MetricTypeUserProductAdd  MetricTypeCount = "product_add"
	MetricTypeUserVideoAdd    MetricTypeCount = "video_add"
	MetricTypeUserServiceAdd  MetricTypeCount = "service_add"
	MetricTypeUserJobAdd      MetricTypeCount = "job_add"
	MetricTypeUserPostAdd     MetricTypeCount = "post_add"
	MetricTypeUserVehicleAdd  MetricTypeCount = "vehicle_add"
	MetricTypeUserPropertyAdd MetricTypeCount = "property_add"

	MetricTypeCountUnknown MetricTypeCount = ""
)

func (s MetricTypeCount) String() string {
	switch s {
	case MetricTypeCountComment, MetricTypeCountMessage, MetricTypeCountLike, MetricTypeCountDislike,
		MetricTypeCountWishlist, MetricTypeCountBasket, MetricTypeCountFollowing, MetricTypeCountReview,
		MetricTypeCountShare, MetricTypeCountView, MetricTypeCountVisit, MetricTypeCountReport,
		MetricTypeCountFollow,
		MetricTypeUserMediaAdd, MetricTypeUserCommentAdd, MetricTypeUserLikeAdd,
		MetricTypeUserProductAdd, MetricTypeUserVideoAdd, MetricTypeUserServiceAdd,
		MetricTypeUserJobAdd, MetricTypeUserPostAdd, MetricTypeUserVehicleAdd, MetricTypeUserPropertyAdd:
		return string(s)
	default:
		return ""
	}
}

func ToMetricTypeCount(s string) MetricTypeCount {

	switch s {
	case MetricTypeCountComment.String():
		return MetricTypeCountComment
	case MetricTypeCountMessage.String():
		return MetricTypeCountMessage
	case MetricTypeCountLike.String():
		return MetricTypeCountLike
	case MetricTypeCountDislike.String():
		return MetricTypeCountDislike
	case MetricTypeCountWishlist.String():
		return MetricTypeCountWishlist
	case MetricTypeCountBasket.String():
		return MetricTypeCountBasket
	case MetricTypeCountFollowing.String():
		return MetricTypeCountFollowing
	case MetricTypeCountReview.String():
		return MetricTypeCountReview
	case MetricTypeCountShare.String():
		return MetricTypeCountShare
	case MetricTypeCountView.String():
		return MetricTypeCountView
	case MetricTypeCountVisit.String():
		return MetricTypeCountVisit
	case MetricTypeCountReport.String():
		return MetricTypeCountReport
	case MetricTypeCountFollow.String():
		return MetricTypeCountFollow
	case MetricTypeUserMediaAdd.String():
		return MetricTypeUserMediaAdd
	case MetricTypeUserCommentAdd.String():
		return MetricTypeUserCommentAdd
	case MetricTypeUserLikeAdd.String():
		return MetricTypeUserLikeAdd
	case MetricTypeUserProductAdd.String():
		return MetricTypeUserProductAdd
	case MetricTypeUserVideoAdd.String():
		return MetricTypeUserVideoAdd
	case MetricTypeUserServiceAdd.String():
		return MetricTypeUserServiceAdd
	case MetricTypeUserJobAdd.String():
		return MetricTypeUserJobAdd
	case MetricTypeUserPostAdd.String():
		return MetricTypeUserPostAdd
	case MetricTypeUserVehicleAdd.String():
		return MetricTypeUserVehicleAdd
	case MetricTypeUserPropertyAdd.String():
		return MetricTypeUserPropertyAdd
	default:
		return MetricTypeCountUnknown
	}
}
