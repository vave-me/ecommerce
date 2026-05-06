package domain

// ---------------------
// SNAPSHOT
// ---------------------

type OfferV1 struct {
	UserSellerID   string
	UserCustomerID string
	ProductID      string
	Price          int64
	Status         OfferStatus
}

func (OfferV1) SnapshotName() string { return "offers.OfferV1" }
