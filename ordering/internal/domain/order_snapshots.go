package domain

type OrderV1 struct {
	UserCustomerID  string
	PaymentMethodID string
	InvoiceID       string
	ShoppingID      string
	Items           []Item
	Status          OrderStatus
	PaymentIntent   string
	BasketID        string
}

func (OrderV1) SnapshotName() string { return "ordering.OrderV1" }
