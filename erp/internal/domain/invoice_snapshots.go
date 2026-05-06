package domain

import (
	"time"
)

type InvoiceV1 struct {
	InvoiceNumber  string
	OrderID        string
	CustomerID     string
	Type           InvoiceType
	Status         InvoiceStatus
	IssueDate      time.Time
	DueDate        time.Time
	Currency       string
	Lines          []InvoiceLine
	SubTotal       float64
	TaxAmount      float64
	DiscountAmount float64
	ShippingAmount float64
	TotalAmount    float64
	PaidAmount     float64
	BalanceDue     float64
	PaymentTerms   string
	PaymentMethod  string
	BillingAddress Address
	TaxID          string
	PONumber       string
	Notes          string
	Attachments    []string
	ConnectorID    string
	ExternalID     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ApprovedAt     *time.Time
	SentAt         *time.Time
	VoidedAt       *time.Time
	PaidAt         *time.Time
}

func (InvoiceV1) SnapshotName() string { return "erp.InvoiceV1" }