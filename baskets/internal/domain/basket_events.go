package domain

type BasketStarted struct {
	UserCustomerID string
	Status         BasketStatus
}

type BasketItemAdded struct {
	Item Item
}

type BasketItemRemoved struct {
	ProductID string
	Quantity  int64
}

type BasketCanceled struct {
	Status BasketStatus
}

type BasketCheckedOut struct {
	Status          BasketStatus
	PaymentIntentID string
}

type BasketReopened struct {
	Status BasketStatus
	Reason string
}
