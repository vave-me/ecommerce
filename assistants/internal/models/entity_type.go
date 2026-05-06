package models

import "strings"

// Enumeration of possible comment statuses.
type EntityType string

const (
	ProductType            EntityType = "product"
	DealType               EntityType = "deal"
	PostType               EntityType = "post"
	VehicleType            EntityType = "vehicle"
	PropertyType           EntityType = "property"
	ServiceType            EntityType = "service"
	JobType                EntityType = "job"
	UserEntityType         EntityType = "user"
	EntityTypeOrder        EntityType = "order"
	OfferEntityType        EntityType = "offer"
	LeaseEntityType        EntityType = "lease"
	BuyNowEntityType       EntityType = "buynow"
	BuyBackEntityType      EntityType = "buyback"
	ReservationEntityType  EntityType = "reservation"
	NotificationEntityType EntityType = "notification"
	AlertEntityType        EntityType = "alert"
	NewsletterEntityType   EntityType = "newsletter"
	SubscriptionEntityType EntityType = "subscription"
	MetricEntityType       EntityType = "metric"
	ItemMetricEntityType   EntityType = "item_metric"
	UserMetricEntityType   EntityType = "user_metric"
	MessageType            EntityType = "message"
	ConversationType       EntityType = "conversation"
	WishlistType           EntityType = "wishlist"
	CommentType            EntityType = "comment"
	ReviewType             EntityType = "review"
	MediaType              EntityType = "media"
	ImageType              EntityType = "image"
	VideoType              EntityType = "video"
	ActivityType           EntityType = "activity"
	InteractionType        EntityType = "interaction"
	EmailType              EntityType = "email"
	MailerType             EntityType = "mailer"
	GeocodingType          EntityType = "geocoding"
	GeocodingRequestType   EntityType = "geocoding_request"
	FollowType             EntityType = "follow"
	FollowingType          EntityType = "following"
	CategoryType           EntityType = "category"
	FilterType             EntityType = "filter"
	BasketType             EntityType = "basket"
	PaymentEntityType      EntityType = "payment"
	ShippingEntityType     EntityType = "shipping"
	SupportEntityType      EntityType = "support"
	VariantEntityType      EntityType = "variant"
	UnknownType            EntityType = "unknown"
)

func (s EntityType) String() string {
	switch s {
	case ProductType, PostType, VehicleType, PropertyType, DealType, JobType, ServiceType, EntityTypeOrder, OfferEntityType, LeaseEntityType, BuyNowEntityType, BuyBackEntityType, ReservationEntityType, NotificationEntityType, AlertEntityType, NewsletterEntityType, SubscriptionEntityType, MetricEntityType, ItemMetricEntityType, UserMetricEntityType, MessageType, ConversationType, WishlistType, CommentType, ReviewType, MediaType, ImageType, VideoType, EmailType, MailerType, GeocodingType, GeocodingRequestType, FollowType, FollowingType, CategoryType, FilterType, BasketType, PaymentEntityType, ShippingEntityType, SupportEntityType, VariantEntityType:
		return string(s)
	default:
		return ""
	}
}
func ToEntityType(status string) EntityType {
	// Normalize input: trim spaces and convert to lowercase for comparison
	normalized := strings.ToLower(strings.TrimSpace(status))
	
	// First, try direct lowercase match
	switch normalized {
	case "product", "products":
		return ProductType
	case "user", "users":
		return UserEntityType
	case "order", "orders":
		return EntityTypeOrder
	case "payment", "payments":
		return PaymentEntityType
	case "comment", "comments":
		return CommentType
	case "review", "reviews":
		return ReviewType
	case "category", "categories":
		return CategoryType
	case "wishlist", "wishlists":
		return WishlistType
	case "basket", "baskets":
		return BasketType
	case "notification", "notifications":
		return NotificationEntityType
	case "message", "messages":
		return MessageType
	case "offer", "offers":
		return OfferEntityType
	case "metric", "metrics":
		return MetricEntityType
	case "newsletter", "newsletters":
		return NewsletterEntityType
	case "vehicle", "vehicles":
		return VehicleType
	case "property", "properties":
		return PropertyType
	case "service", "services":
		return ServiceType
	case "job", "jobs":
		return JobType
	case "post", "posts":
		return PostType
	case "deal", "deals":
		return DealType
	case "variant", "variants":
		return VariantEntityType
	case "shipping":
		return ShippingEntityType
	case "support":
		return SupportEntityType
	case "media":
		return MediaType
	case "video":
		return VideoType
	case "activity", "activities":
		return ActivityType
	case "following":
		return FollowingType
	case "mailer":
		return MailerType
	case "geocoding":
		return GeocodingType
	}
	
	// Now check original case-sensitive matches
	switch status {
	case ProductType.String():
		return ProductType
	case MediaType.String():
		return MediaType
	case VideoType.String():
		return VideoType
	case PostType.String():
		return PostType
	case DealType.String():
		return DealType
	case VehicleType.String():
		return VehicleType
	case PropertyType.String():
		return PropertyType
	case ServiceType.String():
		return ServiceType
	case JobType.String():
		return JobType
	case EntityTypeOrder.String():
		return EntityTypeOrder
	case OfferEntityType.String():
		return OfferEntityType
	case LeaseEntityType.String():
		return LeaseEntityType
	case BuyNowEntityType.String():
		return BuyNowEntityType
	case BuyBackEntityType.String():
		return BuyBackEntityType
	case ReservationEntityType.String():
		return ReservationEntityType
	case NotificationEntityType.String():
		return NotificationEntityType
	case AlertEntityType.String():
		return AlertEntityType
	case NewsletterEntityType.String():
		return NewsletterEntityType
	case SubscriptionEntityType.String():
		return SubscriptionEntityType
	case MetricEntityType.String():
		return MetricEntityType
	case ItemMetricEntityType.String():
		return ItemMetricEntityType
	case UserMetricEntityType.String():
		return UserMetricEntityType
	case MessageType.String():
		return MessageType
	case ConversationType.String():
		return ConversationType
	case WishlistType.String():
		return WishlistType
	case CommentType.String():
		return CommentType
	case ReviewType.String():
		return ReviewType
	case MediaType.String():
		return MediaType
	case ImageType.String():
		return ImageType
	case VideoType.String():
		return VideoType
	case EmailType.String():
		return EmailType
	case MailerType.String():
		return MailerType
	case GeocodingType.String():
		return GeocodingType
	case GeocodingRequestType.String():
		return GeocodingRequestType
	case FollowType.String():
		return FollowType
	case FollowingType.String():
		return FollowingType
	case CategoryType.String():
		return CategoryType
	case InteractionType.String():
		return InteractionType
	case FilterType.String():
		return FilterType
	case ActivityType.String():
		return ActivityType
	case BasketType.String():
		return BasketType
	case PaymentEntityType.String():
		return PaymentEntityType
	case ShippingEntityType.String():
		return ShippingEntityType
	case SupportEntityType.String():
		return SupportEntityType
	case VariantEntityType.String():
		return VariantEntityType
	case UserEntityType.String():
		return UserEntityType
		
	// Handle constant names (for LLM compatibility)
	case "ProductType", "Product":
		return ProductType
	case "PostType", "Post":
		return PostType
	case "DealType", "Deal":
		return DealType
	case "UserEntityType", "UserType", "User":
		return UserEntityType
	case "OrderType", "EntityTypeOrder", "Order":
		return EntityTypeOrder
	case "OfferType", "OfferEntityType", "Offer":
		return OfferEntityType
	case "PaymentType", "PaymentEntityType", "Payment":
		return PaymentEntityType
	case "CommentType", "Comment":
		return CommentType
	case "ReviewType", "Review":
		return ReviewType
	case "CategoryType", "Category":
		return CategoryType
	case "WishlistType", "Wishlist":
		return WishlistType
	case "BasketType", "Basket":
		return BasketType
	case "NotificationType", "NotificationEntityType", "Notification":
		return NotificationEntityType
	case "MessageType", "Message":
		return MessageType
	case "ShippingType", "ShippingEntityType", "Shipping":
		return ShippingEntityType
	case "SupportType", "SupportEntityType", "Support":
		return SupportEntityType
	case "MetricType", "MetricEntityType", "Metric":
		return MetricEntityType
	case "VariantType", "VariantEntityType", "Variant":
		return VariantEntityType
	case "ServiceType", "Service":
		return ServiceType
	case "FollowingType", "Following":
		return FollowingType
	case "MediaType", "Media":
		return MediaType
	case "ActivityType", "Activity":
		return ActivityType
	case "NewsletterType", "NewsletterEntityType", "Newsletter":
		return NewsletterEntityType
	case "MailerType", "Mailer":
		return MailerType
	case "GeocodingType", "Geocoding":
		return GeocodingType
	case "VehicleType", "Vehicle":
		return VehicleType
	case "PropertyType", "Property":
		return PropertyType
	case "JobType", "Job":
		return JobType
		
	default:
		return UnknownType
	}
}
