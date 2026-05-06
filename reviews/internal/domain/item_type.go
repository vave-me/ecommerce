package domain

type ItemType string

const (
	UnknownType  ItemType = ""
	ProductType  ItemType = "product"
	PostType     ItemType = "post"
	VehicleType  ItemType = "vehicle"
	JobType      ItemType = "job"
	ServiceType  ItemType = "service"
	PropertyType ItemType = "property"
	ReviewType   ItemType = "review"
	ImageType    ItemType = "image"
	VideoType    ItemType = "video"
	ProfileType  ItemType = "profile"
)

func (t ItemType) String() string {
	switch t {
	case ProductType, ProfileType, VideoType, ImageType, ReviewType, VehicleType, JobType, PropertyType, PostType, ServiceType:
		return string(t)
	default:
		return ""
	}
}
func ToItemType(interactionType string) ItemType {
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
