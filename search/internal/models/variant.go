package models

// Variant holds information about a specific configuration (e.g., size, color)
// or sub-product linked to a main product.
type Variant struct {
	VariantID    string
	ProductID    string
	Status       string
	Name         string
	SKU          string
	Barcode      string
	Condition    string
	VariantPrice int64
	CurrencyCode string
	Stock        int64
	Weight       int64
	Height       int64
	Width        int64
	Depth        int64
	Attributes   []Attribute
	IsAvailable  bool
	HasOptions   bool
	Options      []Option
}
