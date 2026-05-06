package domain

type VariantV1 struct {
	ID           string
	ProductID    string
	SKU          string
	Barcode      string
	Condition    ProductCondition
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
	Thumbnail    string
	Lat          float64
	Lng          float64
}

func (VariantV1) SnapshotName() string { return "products.VariantV1" }
