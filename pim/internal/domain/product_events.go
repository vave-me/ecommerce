package domain

// ------------------------------------------------------------------------------------
// 1. Product Event Names
// ------------------------------------------------------------------------------------

const (
	ProductAddedEvent             = "products.ProductAdded"             // new product created
	ProductUpdatedEvent           = "products.ProductUpdated"           // new product created
	ProductRebrandedEvent         = "products.ProductRebranded"         // changes name, brand, model, or description
	ProductPriceIncreasedEvent    = "products.ProductPriceIncreased"    // product base price goes up
	ProductPriceDecreasedEvent    = "products.ProductPriceDecreased"    // product base price goes down
	ProductStockAdjustedEvent     = "products.ProductStockAdjusted"     // product stock changes (if manageStock is true)
	ProductNegotiableToggledEvent = "products.ProductNegotiableToggled" // toggles the negotiable field
	ProductRemovedEvent           = "products.ProductRemoved"           // product is removed from listing
	ProductSoldEvent              = "products.ProductSold"              // product is sold
	ProductLeasedEvent            = "products.ProductLeased"            // product is leased
	ProductPawnedEvent            = "products.ProductPawned"
	ProductArchivedEvent          = "products.ProductArchived"         // product is pawned
	ProductThumbnailAddedEvent    = "products.ProductThumbnailAdded"   // product is pawned
	ProductThumbnailRemovedEvent  = "products.ProductThumbnailRemoved" // product is pawned
	ProductThumbnailUpdatedEvent  = "products.ProductThumbnailUpdated" // product is pawned
)

// ------------------------------------------------------------------------------------
// 2. ProductAdded
// This event sets the essential fields for a newly created product
// ------------------------------------------------------------------------------------

type ProductAdded struct {
	Name             string
	Description      string
	BasePrice        int64
	UserSellerID     string
	CategoryID       string
	CategorySlug     string
	Brand            string
	Condition        ProductCondition
	Model            string
	Tags             []string
	ManageStock      bool
	Stock            int64
	SKU              string
	Attributes       []Attribute
	Weight           int64
	Height           int64
	Width            int64
	Depth            int64
	Status           ProductStatus // e.g. "draft" or "active"
	Negotiable       bool
	UserType         UserType
	MiddlemanService bool // e.g. "premium", "base"
	ShippingCost     int64
	HasVariants      bool
	Options          []Option
	Thumbnail        string
	Lat              float64
	Lng              float64
	// If you want to store top-level product options, you could do:
	// Options []Option
}

func (ProductAdded) Key() string { return ProductAddedEvent }

type ProductUpdated struct {
	Name             string
	Description      string
	BasePrice        int64
	UserSellerID     string
	CategoryID       string
	CategorySlug     string
	Brand            string
	Condition        ProductCondition
	Model            string
	Tags             []string
	ManageStock      bool
	Stock            int64
	SKU              string
	Attributes       []Attribute
	Weight           int64
	Height           int64
	Width            int64
	Depth            int64
	Status           ProductStatus // e.g. "draft" or "active"
	Negotiable       bool
	UserType         UserType
	MiddlemanService bool // e.g. "premium", "base"
	ShippingCost     int64
	HasVariants      bool
	Options          []Option
	Thumbnail        string
	Lat              float64
	Lng              float64
	// If you want to store top-level product options, you could do:
	// Options []Option
}

func (ProductUpdated) Key() string { return ProductUpdatedEvent }

type ProductRebranded struct {
	Name        string
	Description string
	Brand       string
	Model       string
	// Possibly condition if rebranding includes changing condition
}

func (ProductRebranded) Key() string { return ProductRebrandedEvent }

// If you prefer storing a delta:
type ProductPriceIncreased struct {
	ProductID string
	OldPrice  int64
	NewPrice  int64
}

func (ProductPriceIncreased) Key() string { return ProductPriceIncreasedEvent }

type ProductPriceDecreased struct {
	ProductID string
	OldPrice  int64
	NewPrice  int64
}

func (ProductPriceDecreased) Key() string { return ProductPriceDecreasedEvent }

type ProductStockAdjusted struct {
	ProductID string
	OldStock  int64
	NewStock  int64
}

func (ProductStockAdjusted) Key() string { return ProductStockAdjustedEvent }

// ------------------------------------------------------------------------------------
// 6. ProductNegotiableToggled
// Switches the 'Negotiable' field from true->false or false->true
// ------------------------------------------------------------------------------------

type ProductNegotiableToggled struct {
	ProductID string
	OldValue  bool
	NewValue  bool
}

func (ProductNegotiableToggled) Key() string { return ProductNegotiableToggledEvent }

// ------------------------------------------------------------------------------------
// 7. ProductRemoved
// Product is removed from listing or from system
// ------------------------------------------------------------------------------------

type ProductRemoved struct {
	ProductID    string
	UserSellerID string // or reason
}

func (ProductRemoved) Key() string { return ProductRemovedEvent }

type ProductSold struct {
	ProductID    string
	UserSellerID string
	BuyerID      string // if needed
	FinalPrice   int64
	SoldAt       string // or a timestamp
}

func (ProductSold) Key() string { return ProductSoldEvent }

type ProductLeased struct {
	ProductID    string
	UserSellerID string
	LesseeID     string
}

func (ProductLeased) Key() string { return ProductLeasedEvent }

type ProductPawned struct {
	ProductID    string
	UserSellerID string
	PawnID       string
}

func (ProductPawned) Key() string { return ProductPawnedEvent }

type ProductArchived struct {
	ProductID    string
	UserSellerID string
}

func (ProductArchived) Key() string { return ProductArchivedEvent }

type ProductThumbnailAdded struct {
	ProductID string
	Thumbnail string
}

func (ProductThumbnailAdded) Key() string { return ProductThumbnailAddedEvent }

type ProductThumbnailRemoved struct {
	ProductID string
	Thumbnail string
}

func (ProductThumbnailRemoved) Key() string { return ProductThumbnailRemovedEvent }

type ProductThumbnailUpdated struct {
	ProductID string
	Thumbnail string
}

func (ProductThumbnailUpdated) Key() string { return ProductThumbnailUpdatedEvent }
