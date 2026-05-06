package models

const InvoicePaidEvent = "payments.InvoicePaid"

const InvoiceCanceledEvent = "payments.InvoiceCanceled"

const InvoicePartialPaidEvent = "payments.InvoicePartialPaid"

type InvoicePartialPaid struct {
	InvoiceID  string
	OrderID    string
	PaidAmount int64
	Remaining  int64
}
type InvoicePaid struct {
	InvoiceID string
	OrderID   string
}
type InvoiceCanceled struct {
	InvoiceID string
	OrderID   string
}
