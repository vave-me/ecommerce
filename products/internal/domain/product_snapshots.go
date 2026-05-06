package domain

type ProductV1 struct {
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
	Status           ProductStatus
	Thumbnail        string
	Negotiable       bool
	UserType         UserType
	MiddlemanService bool
	ShippingCost     int64
	HasVariants      bool
	Options          []Option
	Lat              float64
	Lng              float64
}

func (ProductV1) SnapshotName() string { return "products.ProductV1" }
