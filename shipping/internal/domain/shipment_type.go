package domain

// Enumeration of possible entity types.
type ShipmentType string

const (
	EntityTypeUnspecified ShipmentType = "unspecified"
	ProductType           ShipmentType = "product"
	UserType              ShipmentType = "user" // User entity.
	ReviewType            ShipmentType = "review"
	CommentType           ShipmentType = "comment" // Comment entity.
	MessageType           ShipmentType = "message"
	RatingType            ShipmentType = "rating"
)

type ShippingType string

const (
	ShippingTypeUnspecified ShippingType = "unspecified"
	ImageType               ShippingType = "image"
	VideoType               ShippingType = "video"
)

type ShippingStatus string

const (
	ShippingStatusCreated   ShippingStatus = "created"
	ShippingStatusUploaded  ShippingStatus = "uploaded"
	ShippingStatusProcessed ShippingStatus = "processed"
	// Add other statuses as needed
)

type ShipmentDimensions string

const (
	sizeM   ShipmentDimensions = "sizeM"
	sizeL   ShipmentDimensions = "sizeL"
	sizeXL  ShipmentDimensions = "sizeXL"
	sizeXXL ShipmentDimensions = "sizeXXL"
)
