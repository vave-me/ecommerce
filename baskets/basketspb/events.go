package basketspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	BasketAggregateChannel = "middleman.baskets.events.Basket"

	BasketStartedEvent     = "basketsapi.BasketStarted"
	BasketCanceledEvent    = "basketsapi.BasketCanceled"
	BasketCheckedOutEvent  = "basketsapi.BasketCheckedOut"
	BasketItemAddedEvent   = "basketsapi.BasketItemAdded"
	BasketItemRemovedEvent = "basketsapi.BasketItemRemoved"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// Basket events
	if err := serde.Register(&BasketStarted{}); err != nil {
		return err
	}
	if err := serde.Register(&BasketCanceled{}); err != nil {
		return err
	}
	if err := serde.Register(&BasketCheckedOut{}); err != nil {
		return err
	}
	if err := serde.Register(&BasketItemAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&BasketItemRemoved{}); err != nil {
		return err
	}

	return nil
}

func (*BasketStarted) Key() string     { return BasketStartedEvent }
func (*BasketCanceled) Key() string    { return BasketCanceledEvent }
func (*BasketCheckedOut) Key() string  { return BasketCheckedOutEvent }
func (*BasketItemAdded) Key() string   { return BasketItemAddedEvent }
func (*BasketItemRemoved) Key() string { return BasketItemRemovedEvent }
