package domain

const (
	WishlistItemAddedEvent   = "wishlists.WishlistItemAdded"
	WishlistItemRemovedEvent = "wishlists.WishlistItemRemoved"
)

type WishlistItemAdded struct {
	WishlistID string
	ItemID     string
	EntityType string
}

// Key implements registry.Registerable
func (WishlistItemAdded) Key() string { return WishlistItemAddedEvent }

type WishlistItemRemoved struct {
	ID         string
	WishlistID string
	EntityType string
}

// Key implements registry.Registerable
func (WishlistItemRemoved) Key() string { return WishlistItemRemovedEvent }
