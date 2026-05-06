package domain

type BasketV1 struct {
	UserCustomerID  string
	PaymentMethodID string
	Items           map[string]Item
	Status          BasketStatus
	PaymentIntentID string
}

func (BasketV1) SnapshotName() string { return "baskets.BasketV1" }
