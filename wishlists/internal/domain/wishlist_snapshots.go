package domain

type WishlistV1 struct {
	UserID      string
	Name        string
	Description string
}

func (WishlistV1) SnapshotName() string { return "wishlists.WishlistV1" }
