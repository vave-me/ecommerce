package models

type InvoiceStatus string

const (
	InvoiceIsUnknown   InvoiceStatus = ""
	InvoiceIsPending   InvoiceStatus = "pending" // default after creation
	InvoiceIsPartially InvoiceStatus = "partial" // if partial payment was made
	InvoiceIsPaid      InvoiceStatus = "paid"
	InvoiceIsCanceled  InvoiceStatus = "canceled"
)

type Invoice struct {
	ID      string        // Unique invoice ID
	OrderID string        // Link to the related Order
	Amount  int64         // Total amount due
	PaidAmt int64         // How much has been paid so far (for partial payments)
	Status  InvoiceStatus // Current invoice status
}

// String returns a string representation of the InvoiceStatus.
// If the status is unrecognized, returns "".
func (s InvoiceStatus) String() string {
	switch s {
	case InvoiceIsPending, InvoiceIsPartially, InvoiceIsPaid, InvoiceIsCanceled:
		return string(s)
	default:
		return ""
	}
}
