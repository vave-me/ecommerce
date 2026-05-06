package notificationspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	AlertAggregateChannel      = "middleman.notifications.events.Alert"
	BasketAlertAddedEvent      = "notificationsapi.BasketAlertAdded"
	CommentAlertAddedEvent     = "notificationsapi.CommentAlertAdded"
	InteractionAlertAddedEvent = "notificationsapi.InteractionAlertAdded"
	MessageAlertAddedEvent     = "notificationsapi.MessageAlertAdded"
	OfferAlertAddedEvent       = "notificationsapi.OfferAlertAdded"
	ProductAlertAddedEvent     = "notificationsapi.ProductAlertAdded"
	SupportAlertAddedEvent     = "notificationsapi.SupportAlertAdded"
	UserAlertAddedEvent        = "notificationsapi.UserAlertAdded"
	WishlistAlertAddedEvent    = "notificationsapi.WishlistAlertAdded"
	OrderAlertAddedEvent       = "notificationsapi.OrderAlertAdded"
	ReviewAlertAddedEvent      = "notificationsapi.ReviewAlertAdded"
	PaymentAlertAddedEvent     = "notificationsapi.PaymentAlertAdded"
	FollowingAlertAddedEvent   = "notificationsapi.FollowingAlertAdded"
	AlertReadEvent             = "notificationsapi.AlertRead"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}
func RegistrationsWithSerde(serde registry.Serde) error {

	// Notifications events
	if err := serde.Register(&BasketAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&CommentAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&InteractionAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&MessageAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&OfferAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&SupportAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&UserAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&WishlistAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&OrderAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&ReviewAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&PaymentAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&FollowingAlertAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&AlertRead{}); err != nil {
		return err
	}
	// commands

	return nil
}

func (*BasketAlertAdded) Key() string      { return BasketAlertAddedEvent }
func (*CommentAlertAdded) Key() string     { return CommentAlertAddedEvent }
func (*InteractionAlertAdded) Key() string { return InteractionAlertAddedEvent }
func (*MessageAlertAdded) Key() string     { return MessageAlertAddedEvent }
func (*OfferAlertAdded) Key() string       { return OfferAlertAddedEvent }
func (*ProductAlertAdded) Key() string     { return ProductAlertAddedEvent }
func (*SupportAlertAdded) Key() string     { return SupportAlertAddedEvent }
func (*UserAlertAdded) Key() string        { return UserAlertAddedEvent }
func (*WishlistAlertAdded) Key() string    { return WishlistAlertAddedEvent }
func (*OrderAlertAdded) Key() string       { return OrderAlertAddedEvent }
func (*ReviewAlertAdded) Key() string      { return ReviewAlertAddedEvent }
func (*PaymentAlertAdded) Key() string     { return PaymentAlertAddedEvent }
func (*FollowingAlertAdded) Key() string   { return FollowingAlertAddedEvent }
func (*AlertRead) Key() string             { return AlertReadEvent }
