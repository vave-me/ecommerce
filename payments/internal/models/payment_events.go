package models

const PaymentAuthorizedEvent = "payments.PaymentAuthorized"
const PaymentFailedEvent = "payments.PaymentFailed"
const PaymentConfirmedEvent = "payments.PaymentConfirmed"

type PaymentAuthorized struct {
	PaymentID string
	BasketID  string // correlation with customer basket / order
	OrderID   string // filled once correlation exists (optional)
	Amount    int64
	UserID    string
}

type PaymentFailed struct {
	PaymentID string
	Reason    string
}

type PaymentConfirmed struct {
	PaymentID  string
	OrderID    string // the associated Order aggregate ID
	ShoppingID string // reservation / fulfillment reference (optional)
	Amount     int64
}
