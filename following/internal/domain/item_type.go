package domain

type FollowedUserType string

const (
	UnknownType  FollowedUserType = ""
	ProductType  FollowedUserType = "product"
	PostType     FollowedUserType = "post"
	VehicleType  FollowedUserType = "vehicle"
	JobType      FollowedUserType = "job"
	ServiceType  FollowedUserType = "service"
	PropertyType FollowedUserType = "property"
	ReviewType   FollowedUserType = "review"
	ImageType    FollowedUserType = "image"
	VideoType    FollowedUserType = "video"
	ProfileType  FollowedUserType = "profile"
)

func (t FollowedUserType) String() string {
	switch t {
	case ProductType, ProfileType, VideoType, ImageType, ReviewType, VehicleType, JobType, PropertyType, PostType, ServiceType:
		return string(t)
	default:
		return ""
	}
}
func ToFollowedUserType(interactionType string) FollowedUserType {
	switch interactionType {

	case ProductType.String():
		return ProductType
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
