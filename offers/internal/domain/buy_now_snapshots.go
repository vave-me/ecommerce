package domain

type BuyNowV1 struct {
	OfferID             string
	FinalPrice          int64
	Status              BuyNowStatus
	NegotiatedPrice     int64
	NegotiationComments string
}

func (BuyNowV1) SnapshotName() string { return "offers.BuyNowV1" }
