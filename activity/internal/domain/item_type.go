package domain

type ItemType string

const (
	UnknownType  ItemType = ""
	ProductType  ItemType = "product"
	CommentType  ItemType = "comment"
	ReplyType    ItemType = "reply"
	ReviewType   ItemType = "review"
	ImageType    ItemType = "image"
	VideoType    ItemType = "video"
	ProfileType  ItemType = "profile"
	VehicleType  ItemType = "vehicle"
	JobType      ItemType = "job"
	PostType     ItemType = "post"
	ServiceType  ItemType = "service"
	PropertyType ItemType = "property"
)

func (t ItemType) String() string {
	switch t {
	case ReplyType, ProductType, CommentType, ProfileType, VideoType, ImageType, ReviewType, VehicleType, JobType, PropertyType, PostType, ServiceType:
		return string(t)
	default:
		return ""
	}
}
func ToInteractionType(interactionType string) ItemType {
	switch interactionType {

	case ProductType.String():
		return ProductType
	case ReplyType.String():
		return ReplyType
	case CommentType.String():
		return CommentType
	case ProfileType.String():
		return ProfileType
	case VideoType.String():
		return VideoType
	case ReviewType.String():
		return ReviewType
	case ImageType.String():
		return ImageType
	case ServiceType.String():
		return ServiceType
	case PropertyType.String():
		return PropertyType
	case PostType.String():
		return PostType
	case JobType.String():
		return JobType
	case VehicleType.String():
		return VehicleType

	default:
		return UnknownType
	}
}

type ActionType string

const (
	Like    ActionType = "like"
	Dislike ActionType = "dislike"
)

func (t ActionType) String() string {
	switch t {
	case Like, Dislike:
		return string(t)
	default:
		return ""
	}
}
