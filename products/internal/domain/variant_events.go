package domain

const (
	VariantAddedEvent          = "products.VariantAdded"
	VariantPriceIncreasedEvent = "products.VariantPriceIncreased"
	VariantPriceDecreasedEvent = "products.VariantPriceDecreased"
	VariantStockAdjustedEvent  = "products.VariantStockAdjusted"
	VariantArchivedEvent       = "products.VariantArchived"
	VariantRemovedEvent        = "products.VariantRemoved"
	VariantRebrandedEvent      = "products.VariantRebranded"
	// etc.
)

// VariantCreated sets up a new variant referencing a product
type VariantAdded struct {
	ProductID    string
	Status       ProductStatus
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
}

func (VariantAdded) Key() string { return VariantAddedEvent }

type VariantRebranded struct {
	VariantID string
	Name      string
	Desc      string
}

func (VariantRebranded) Key() string { return VariantRebrandedEvent }

type VariantPriceDecreased struct {
	VariantID string // variant ID
	OldPrice  int64
	NewPrice  int64
}

func (VariantPriceDecreased) Key() string { return VariantPriceDecreasedEvent }

// VariantPriceChanged if the variant price changes
type VariantPriceIncreased struct {
	VariantID string
	OldPrice  int64
	NewPrice  int64
}

func (VariantPriceIncreased) Key() string { return VariantPriceIncreasedEvent }

// VariantStockAdjusted for changes in variant-level stock
type VariantStockAdjusted struct {
	VariantID string
	OldStock  int64
	NewStock  int64
}

func (VariantStockAdjusted) Key() string { return VariantStockAdjustedEvent }

// VariantArchived might mark the variant as no longer available
type VariantArchived struct {
	VariantID string
}

func (VariantArchived) Key() string { return VariantArchivedEvent }

// VariantRemoved if the variant is fully removed from the system
type VariantRemoved struct {
	VariantID string
}

func (VariantRemoved) Key() string { return VariantRemovedEvent }
