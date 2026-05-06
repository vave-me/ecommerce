package domain

import (
	"time"
)

// Invoice Event Names
const (
	InvoiceCreatedEvent         = "erp.InvoiceCreated"
	InvoiceApprovedEvent        = "erp.InvoiceApproved"
	InvoiceSentEvent            = "erp.InvoiceSent"
	InvoiceVoidedEvent          = "erp.InvoiceVoided"
	InvoicePaymentReceivedEvent = "erp.InvoicePaymentReceived"
	InvoiceLinkedToERPEvent     = "erp.InvoiceLinkedToERP"
)

// InvoiceCreated event - fired when a new invoice is created
type InvoiceCreated struct {
	InvoiceNumber  string
	OrderID        string
	CustomerID     string
	Type           InvoiceType
	IssueDate      time.Time
	DueDate        time.Time
	Currency       string
	Lines          []InvoiceLine
	SubTotal       float64
	TaxAmount      float64
	DiscountAmount float64
	ShippingAmount float64
	TotalAmount    float64
	PaymentTerms   string
	BillingAddress Address
	TaxID          string
	PONumber       string
	Notes          string
	ConnectorID    string
}

// Key implements event registry
func (InvoiceCreated) Key() string { return InvoiceCreatedEvent }

// InvoiceApproved event - fired when invoice is approved
type InvoiceApproved struct {
	ApprovedBy string
	ApprovedAt time.Time
}

// Key implements event registry
func (InvoiceApproved) Key() string { return InvoiceApprovedEvent }

// InvoiceSent event - fired when invoice is sent to customer
type InvoiceSent struct {
	SentTo []string
	SentBy string
	SentAt time.Time
}

// Key implements event registry
func (InvoiceSent) Key() string { return InvoiceSentEvent }

// InvoiceVoided event - fired when invoice is voided
type InvoiceVoided struct {
	Reason   string
	VoidedBy string
	VoidedAt time.Time
}

// Key implements event registry
func (InvoiceVoided) Key() string { return InvoiceVoidedEvent }

// InvoicePaymentReceived event - fired when payment is received
type InvoicePaymentReceived struct {
	Amount        float64
	PaymentMethod string
	TransactionID string
	PaymentDate   time.Time
}

// Key implements event registry
func (InvoicePaymentReceived) Key() string { return InvoicePaymentReceivedEvent }

// InvoiceLinkedToERP event - fired when invoice is linked to external ERP
type InvoiceLinkedToERP struct {
	ExternalID string
	LinkedAt   time.Time
}

// Key implements event registry
func (InvoiceLinkedToERP) Key() string { return InvoiceLinkedToERPEvent }