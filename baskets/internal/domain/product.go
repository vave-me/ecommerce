package domain

type Product struct {
	ID               string
	Name             string
	Description      string
	BasePrice        int64
	UserSellerID     string
	CategoryID       string
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
	Negotiable       bool
	UserType         UserType
	MiddlemanService bool
	ShippingCost     int64
	HasVariants      bool
	Options          []Option
	Thumbnail        string
	Lat              float64
	Lng              float64
}
