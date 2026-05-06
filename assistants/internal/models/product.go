package models

type Product struct {
	ProductID        string
	Name             string
	Description      string
	BasePrice        int64 // or int64 if you store cents
	UserSellerID     string
	CategoryID       string
	CategorySlug     string
	Brand            string
	Condition        string
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
	Status           string
	Negotiable       bool
	UserType         string
	MiddlemanService bool // was string, but domain is bool
	ShippingCost     int64
	HasVariants      bool
	Options          []Option
	Lat              float64
	Lng              float64
	Thumbnail        string
	EntityType       EntityType
	Metrics          *ItemMetric
}
