package domain

type WishlistItemV1 struct {
	WishlistID string
	ItemID     string
}

func (WishlistItemV1) SnapshotName() string { return "wishlist.WishlistItemV1" }
