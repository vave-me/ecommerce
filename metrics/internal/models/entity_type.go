package models

// Enumeration of possible comment statuses.
type EntityType string

const (
	ProductType  EntityType = "product"
	DealType     EntityType = "deal"
	PostType     EntityType = "post"
	VehicleType  EntityType = "vehicle"
	PropertyType EntityType = "property"
	ServiceType  EntityType = "service"
	VideoType    EntityType = "video"
	MediaType    EntityType = "media"
	JobType      EntityType = "job"
	UnknownType  EntityType = "unknown"
)

func (s EntityType) String() string {
	switch s {
	case ProductType, PostType, VehicleType, PropertyType, DealType, JobType, ServiceType, VideoType, MediaType:
		return string(s)
	default:
		return ""
	}
}
func ToEntityType(status string) EntityType {
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
	default:
		return UnknownType
	}
}
