package wishlistspb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

const (
	WishlistAggregateChannel = "middleman.wishlists.events.Wishlist"

	WishlistCreatedEvent = "wishlistsapi.WishlistCreated"
	WishlistRemovedEvent = "wishlistsapi.WishlistRemoved"

	WishlistItemAggregateChannel = "middleman.wishlists.events.WishlistItem"
	WishlistItemAddedEvent       = "wishlistsapi.WishlistItemAdded"
	WishlistItemRemovedEvent     = "wishlistsapi.WishlistItemRemoved"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Store events
	if err := serde.Register(&WishlistCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&WishlistRemoved{}); err != nil {
		return err
	}

	if err := serde.Register(&WishlistItemAdded{}); err != nil {
		return err
	}
	if err := serde.Register(&WishlistItemRemoved{}); err != nil {
		return err
	}

	return nil
}

func (*WishlistCreated) Key() string { return WishlistCreatedEvent }
func (*WishlistRemoved) Key() string { return WishlistRemovedEvent }

func (*WishlistItemAdded) Key() string   { return WishlistItemAddedEvent }
func (*WishlistItemRemoved) Key() string { return WishlistItemRemovedEvent }
