package domain

// Enumeration of possible entity types.
type ItemType string

const (
	EntityTypeUnspecified ItemType = "unspecified"
	ProductType           ItemType = "product"
	PostTypeType          ItemType = "post"
	UserType              ItemType = "user" // User entity.
	ReviewType            ItemType = "review"
	CommentType           ItemType = "comment" // Comment entity.
	MessageType           ItemType = "message"
	RatingType            ItemType = "rating"
	VehicleType           ItemType = "vehicle"
	PropertyType          ItemType = "property"
	ServiceType           ItemType = "service"
	JobType               ItemType = "job"
	AdvertiseType         ItemType = "advertise"
	ArticleType           ItemType = "article"
	AiGeneratedType       ItemType = "generated"
)

type MediaType string

const (
	MediaTypeUnspecified MediaType = "unspecified"
	ImageType            MediaType = "image"
	VideoType            MediaType = "video"
)

type MediaStatus string

const (
	MediaStatusCreated   MediaStatus = "created"
	MediaStatusUploaded  MediaStatus = "uploaded"
	MediaStatusProcessed MediaStatus = "processed"
	// Add other statuses as needed
)
