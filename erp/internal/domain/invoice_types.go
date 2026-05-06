package domain

// InvoiceType represents the type of invoice
type InvoiceType string

const (
	InvoiceTypeStandard   InvoiceType = "standard"
	InvoiceTypeCredit     InvoiceType = "credit_note"
	InvoiceTypeDebit      InvoiceType = "debit_note"
	InvoiceTypeProforma   InvoiceType = "proforma"
	InvoiceTypeRecurring  InvoiceType = "recurring"
	InvoiceTypeCommercial InvoiceType = "commercial"
)

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusApproved  InvoiceStatus = "approved"
	InvoiceStatusSent      InvoiceStatus = "sent"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusPartial   InvoiceStatus = "partial"
	InvoiceStatusOverdue   InvoiceStatus = "overdue"
	InvoiceStatusVoided    InvoiceStatus = "voided"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

// InvoiceLine represents a line item on an invoice
type InvoiceLine struct {
	SKU         string
	ProductName string
	Description string
	Quantity    int
	UnitPrice   float64
	TaxRate     float64
	TaxAmount   float64
	Discount    float64
	LineTotal   float64
	AccountCode string
}

