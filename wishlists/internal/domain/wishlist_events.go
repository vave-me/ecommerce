package domain

const (
	WishlistCreatedEvent = "wishlist.WishlistCreated"
	WishlistRemovedEvent = "wishlist.WishlistRemoved"
)

type WishlistCreated struct {
	WishlistID string
	Name       string
	UserID     string
}

// Key implements registry.Registerable
func (WishlistCreated) Key() string { return WishlistCreatedEvent }

type WishlistRemoved struct {
	WishlistID string
	UserID     string
}

// Key implements registry.Registerable
func (WishlistRemoved) Key() string { return WishlistRemovedEvent }
