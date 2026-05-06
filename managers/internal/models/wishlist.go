package models

type WishlistItem struct {
	ID         string
	WishlistID string
	ItemID     string
	EntityType string
	Notes      string
}

type Wishlist struct {
	ID          string
	UserID      string
	Name        string
	Description string
}
