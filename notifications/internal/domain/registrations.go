package domain

import (
	"middleman/internal/es"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	UserAlertAddedEvent        = "notifications.UserAlertAdded"
	ProductAlertAddedEvent     = "notifications.ProductAlertAdded"
	BasketAlertAddedEvent      = "notifications.BasketAlertAdded"
	OrderAlertAddedEvent       = "notifications.OrderAlertAdded"
	WishlistAlertAddedEvent    = "notifications.WishlistAlertAdded"
	MessageAlertAddedEvent     = "notifications.MessageAlertAdded"
	InteractionAlertAddedEvent = "notifications.InteractionAlertAdded"
	CommentAlertAddedEvent     = "notifications.CommentAlertAdded"
	OfferAlertAddedEvent       = "notifications.OfferAlertAdded"
	SupportAlertAddedEvent     = "notifications.SupportAlertAdded"
	ReviewAlertAddedEvent      = "notifications.ReviewAlertAdded"
	PaymentAlertAddedEvent     = "notifications.PaymentAlertAdded"
	FollowingAlertAddedEvent   = "notifications.FollowingAlertAdded"
	AlertReadEvent             = "notifications.AlertRead"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	if err := serde.Register(Alert{}, func(v interface{}) error {
		alert := v.(*Alert)
		alert.Aggregate = es.NewAggregate("", AlertAggregate)
		return nil
	}); err != nil {
		return err
	}

	if err := serde.Register(ProductAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(BasketAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(OrderAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(WishlistAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(MessageAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(InteractionAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(CommentAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(OfferAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(SupportAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(AlertRead{}); err != nil {
		return err
	}
	if err := serde.Register(ReviewAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(PaymentAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(FollowingAlertAdded{}); err != nil {
		return err
	}
	// basket snapshots

	return nil
}

func (Alert) Key() string                 { return AlertAggregate }
func (BasketAlertAdded) Key() string      { return BasketAlertAddedEvent }
func (ProductAlertAdded) Key() string     { return ProductAlertAddedEvent }
func (OrderAlertAdded) Key() string       { return OrderAlertAddedEvent }
func (WishlistAlertAdded) Key() string    { return WishlistAlertAddedEvent }
func (MessageAlertAdded) Key() string     { return MessageAlertAddedEvent }
func (InteractionAlertAdded) Key() string { return InteractionAlertAddedEvent }
func (CommentAlertAdded) Key() string     { return CommentAlertAddedEvent }
func (OfferAlertAdded) Key() string       { return OfferAlertAddedEvent }
func (SupportAlertAdded) Key() string     { return SupportAlertAddedEvent }
func (AlertRead) Key() string             { return AlertReadEvent }
func (ReviewAlertAdded) Key() string      { return ReviewAlertAddedEvent }
func (PaymentAlertAdded) Key() string     { return PaymentAlertAddedEvent }
func (FollowingAlertAdded) Key() string   { return FollowingAlertAddedEvent }
