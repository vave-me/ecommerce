package productspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

// Channels
const (
	ProductAggregateChannel = "middleman.products.events.Product"
	VariantAggregateChannel = "middleman.products.events.Variant"
)

// Product Event Names
const (
	ProductAddedEvent             = "productsapi.ProductAdded"
	ProductRebrandedEvent         = "productsapi.ProductRebranded"
	ProductUpdatedEvent           = "productsapi.ProductUpdated"
	ProductPriceIncreasedEvent    = "productsapi.ProductPriceIncreased"
	ProductPriceDecreasedEvent    = "productsapi.ProductPriceDecreased"
	ProductStockAdjustedEvent     = "productsapi.ProductStockAdjusted"
	ProductRemovedEvent           = "productsapi.ProductRemoved"
	ProductThumbnailAddedEvent    = "productsapi.ProductThumbnailAdded"
	ProductThumbnailUpdatedEvent  = "productsapi.ProductThumbnailUpdated"
	ProductArchivedEvent          = "productsapi.ProductArchived"
	ProductNegotiableToggledEvent = "productsapi.ProductNegotiableToggled"
	ProductSoldEvent              = "productsapi.ProductSold"
	ProductLeasedEvent            = "productsapi.ProductLeased"
	ProductPawnedEvent            = "productsapi.ProductPawned"
)

// Variant Event Names
const (
	VariantAddedEvent          = "productsapi.VariantAdded"
	VariantPriceIncreasedEvent = "productsapi.VariantPriceIncreased"
	VariantPriceDecreasedEvent = "productsapi.VariantPriceDecreased"
	VariantStockAdjustedEvent  = "productsapi.VariantStockAdjusted"
	VariantArchivedEvent       = "productsapi.VariantArchived"
	VariantRemovedEvent        = "productsapi.VariantRemoved"
)

// Commands & Command Channel
const (
	CommandChannel         = "middleman.products.commands"
	ReserveProductCommand  = "productsapi.ReserveProduct"
	ReleaseProductCommand  = "productsapi.ReleaseProduct"
	ReserveProductsCommand = "productsapi.ReserveProducts"
	ReleaseProductsCommand = "productsapi.ReleaseProducts"
)

// Registrations and Serde
func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Product
	if err := serde.Register(&ProductAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductRebranded{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductPriceIncreased{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductPriceDecreased{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductRemoved{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductStockAdjusted{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductArchived{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductNegotiableToggled{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductSold{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductLeased{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductPawned{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductThumbnailAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductThumbnailUpdated{}); err != nil {
		return err
	}
	// Variants
	if err := serde.Register(&VariantAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&VariantPriceIncreased{}); err != nil {
		return err
	}
	if err := serde.Register(&VariantPriceDecreased{}); err != nil {
		return err
	}
	if err := serde.Register(&VariantStockAdjusted{}); err != nil {
		return err
	}
	if err := serde.Register(&VariantArchived{}); err != nil {
		return err
	}
	if err := serde.Register(&VariantRemoved{}); err != nil {
		return err
	}

	// Commands
	if err := serde.Register(&ReserveProduct{}); err != nil {
		return err
	}
	if err := serde.Register(&ReleaseProduct{}); err != nil {
		return err
	}
	if err := serde.Register(&ReserveProducts{}); err != nil {
		return err
	}
	if err := serde.Register(&ReleaseProducts{}); err != nil {
		return err
	}

	return nil
}

func (*ProductAdded) Key() string             { return ProductAddedEvent }
func (*ProductUpdated) Key() string           { return ProductUpdatedEvent }
func (*ProductRebranded) Key() string         { return ProductRebrandedEvent }
func (*ProductPriceIncreased) Key() string    { return ProductPriceIncreasedEvent }
func (*ProductPriceDecreased) Key() string    { return ProductPriceDecreasedEvent }
func (*ProductRemoved) Key() string           { return ProductRemovedEvent }
func (*ProductThumbnailAdded) Key() string    { return ProductThumbnailAddedEvent }
func (*ProductThumbnailUpdated) Key() string  { return ProductThumbnailUpdatedEvent }
func (*ProductStockAdjusted) Key() string     { return ProductStockAdjustedEvent }
func (*ProductArchived) Key() string          { return ProductArchivedEvent }
func (*ProductNegotiableToggled) Key() string { return ProductNegotiableToggledEvent }
func (*ProductSold) Key() string              { return ProductSoldEvent }
func (*ProductLeased) Key() string            { return ProductLeasedEvent }
func (*ProductPawned) Key() string            { return ProductPawnedEvent }
func (*VariantAdded) Key() string             { return VariantAddedEvent }
func (*VariantPriceIncreased) Key() string    { return VariantPriceIncreasedEvent }
func (*VariantPriceDecreased) Key() string    { return VariantPriceDecreasedEvent }
func (*VariantStockAdjusted) Key() string     { return VariantStockAdjustedEvent }
func (*VariantArchived) Key() string          { return VariantArchivedEvent }
func (*VariantRemoved) Key() string           { return VariantRemovedEvent }

// Commands implement registry.Registerable so they can travel via NATS
func (*ReserveProduct) Key() string  { return ReserveProductCommand }
func (*ReleaseProduct) Key() string  { return ReleaseProductCommand }
func (*ReserveProducts) Key() string { return ReserveProductsCommand }
func (*ReleaseProducts) Key() string { return ReleaseProductsCommand }
