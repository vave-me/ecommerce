package models

// PaymentStatus enumerates possible payment states.
const (
	PaymentStatusUnknown    = ""
	PaymentStatusAuthorized = "AUTHORIZED"
	PaymentStatusConfirmed  = "CONFIRMED"
	PaymentStatusFailed     = "FAILED"
	PaymentStatusPartial    = "PARTIAL"
)

// Payment captures a single transaction or authorization attempt.
type Payment struct {
	ID              string
	UserCustomerID  string
	Amount          int64
	PaymentMethodID string
	PaymentIntentID string // Stripe PaymentIntent identifier, not-null in DB
	OrderID         string // associated order id
	Status          string
}
