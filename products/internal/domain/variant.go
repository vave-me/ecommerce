package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const VariantAggregate = "products.Variant"

type Variant struct {
	es.Aggregate
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

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Variant)(nil)

func NewVariant(id string) *Variant {
	return &Variant{
		Aggregate: es.NewAggregate(id, VariantAggregate),
	}
}

func (v *Variant) InitVariant(id, productID, sku string, attributes []Attribute, condition ProductCondition, barcode string, variantPrice int64, currencyCode string, stock, weight, height, width, depth int64, isAvailable, hasOptions bool, Options []Option) (ddd.Event, error) {

	v.AddEvent(VariantAddedEvent, &VariantAdded{
		ProductID:    productID,
		SKU:          sku,
		Barcode:      barcode,
		Condition:    condition,
		VariantPrice: variantPrice,
		CurrencyCode: currencyCode,
		Stock:        stock,
		Weight:       weight,
		Height:       height,
		Width:        width,
		Depth:        depth,
		Attributes:   attributes,
		IsAvailable:  isAvailable,
		HasOptions:   hasOptions,
		Options:      Options,
	})

	return ddd.NewEvent(VariantAddedEvent, v), nil
}

func (Variant) Key() string { return VariantAggregate }

func (v *Variant) IncreasePrice(variantID string, newPrice int64) (ddd.Event, error) {
	if newPrice < v.VariantPrice {
		return nil, ErrNotAPriceIncrease
	}
	v.AddEvent(VariantPriceIncreasedEvent, &VariantPriceIncreased{
		VariantID: variantID,
		OldPrice:  v.VariantPrice,
		NewPrice:  newPrice,
	})
	return ddd.NewEvent(VariantPriceIncreasedEvent, v), nil
}

func (v *Variant) DecreasePrice(variantID string, newPrice int64) (ddd.Event, error) {
	if newPrice > v.VariantPrice {
		return nil, ErrNotAPriceDecrease
	}
	v.AddEvent(VariantPriceDecreasedEvent, &VariantPriceDecreased{
		VariantID: variantID,
		OldPrice:  v.VariantPrice,
		NewPrice:  newPrice,
	})
	return ddd.NewEvent(VariantPriceDecreasedEvent, v), nil
}

func (v *Variant) AdjustStock(newStock int64) (ddd.Event, error) {
	v.AddEvent(VariantStockAdjustedEvent, &VariantStockAdjusted{
		VariantID: v.ID(),
		OldStock:  v.Stock,
		NewStock:  newStock,
	})
	return ddd.NewEvent(VariantStockAdjustedEvent, v), nil
}

func (v *Variant) Archive() (ddd.Event, error) {
	v.IsAvailable = false
	v.AddEvent(VariantArchivedEvent, &VariantArchived{
		VariantID: v.ID(),
	})
	return ddd.NewEvent(VariantArchivedEvent, v), nil
}
func (v *Variant) Rebrand(name, description string) (ddd.Event, error) {
	// If you keep name/desc for variants
	v.AddEvent(VariantRebrandedEvent, &VariantRebranded{
		VariantID: v.ID(),
		Name:      name,
		Desc:      description,
	})
	return ddd.NewEvent(VariantRebrandedEvent, v), nil
}

func (v *Variant) Remove(id string) (ddd.Event, error) {
	v.AddEvent(VariantRemovedEvent, &VariantRemoved{
		VariantID: id,
	})
	return ddd.NewEvent(VariantRemovedEvent, v), nil
}
func (v *Variant) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {

	case *VariantAdded:
		v.ProductID = e.ProductID
		v.SKU = e.SKU
		v.Barcode = e.Barcode
		v.Condition = e.Condition
		v.VariantPrice = e.VariantPrice
		v.CurrencyCode = e.CurrencyCode
		v.Stock = e.Stock
		v.Weight = e.Weight
		v.Height = e.Height
		v.Width = e.Width
		v.Depth = e.Depth
		v.Attributes = e.Attributes
		v.IsAvailable = e.IsAvailable
		v.HasOptions = e.HasOptions
		v.Options = e.Options

	case *VariantPriceIncreased:
		v.VariantPrice = e.NewPrice

	case *VariantPriceDecreased:
		v.VariantPrice = e.NewPrice

	case *VariantStockAdjusted:
		v.Stock = e.NewStock

	case *VariantArchived:
		v.IsAvailable = false // or set some "archived" state

	case *VariantRemoved:
		// Possibly mark a "removed" flag or do nothing
		// e.g. v.IsAvailable = false

	default:
		return errors.ErrInternal.
			Msgf("%T received the event %s with unexpected payload %T", v, event.EventName(), e)
	}
	return nil
}

func (v *Variant) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *VariantV1:
		v.ProductID = ss.ProductID
		v.SKU = ss.SKU
		v.Barcode = ss.Barcode
		v.Condition = ss.Condition
		v.VariantPrice = ss.VariantPrice
		v.CurrencyCode = ss.CurrencyCode
		v.Stock = ss.Stock
		v.Weight = ss.Weight
		v.Height = ss.Height
		v.Width = ss.Width
		v.Depth = ss.Depth
		v.Attributes = ss.Attributes
		v.IsAvailable = ss.IsAvailable
		v.HasOptions = ss.HasOptions
		v.Options = ss.Options

	default:
		return errors.ErrInternal.
			Msgf("%T received the unexpected snapshot %T", v, snapshot)
	}
	return nil
}
func (v Variant) ToSnapshot() es.Snapshot {
	return &VariantV1{
		ProductID:    v.ProductID,
		SKU:          v.SKU,
		Barcode:      v.Barcode,
		Condition:    v.Condition,
		VariantPrice: v.VariantPrice,
		CurrencyCode: v.CurrencyCode,
		Stock:        v.Stock,
		Weight:       v.Weight,
		Height:       v.Height,
		Width:        v.Width,
		Depth:        v.Depth,
		Attributes:   v.Attributes,
		IsAvailable:  v.IsAvailable,
		HasOptions:   v.HasOptions,
		Options:      v.Options,
	}
}
