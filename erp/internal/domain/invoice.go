package domain

import (
	"middleman/internal/ddd"
	"middleman/internal/es"
	"time"

	"github.com/stackus/errors"
)

const InvoiceAggregate = "erp.Invoice"

var (
	ErrInvoiceNumberIsBlank   = errors.Wrap(errors.ErrBadRequest, "invoice number cannot be blank")
	ErrInvoiceCustomerIsBlank = errors.Wrap(errors.ErrBadRequest, "customer ID cannot be blank")
	ErrInvoiceOrderIsBlank    = errors.Wrap(errors.ErrBadRequest, "order ID cannot be blank")
	ErrInvoiceAmountNegative  = errors.Wrap(errors.ErrBadRequest, "invoice amount cannot be negative")
	ErrInvoiceNoLines         = errors.Wrap(errors.ErrBadRequest, "invoice must have at least one line item")
	ErrInvoiceAlreadyApproved = errors.Wrap(errors.ErrBadRequest, "invoice is already approved")
	ErrInvoiceAlreadyVoided   = errors.Wrap(errors.ErrBadRequest, "invoice is already voided")
	ErrInvoiceNotApproved     = errors.Wrap(errors.ErrBadRequest, "invoice must be approved before sending")
	ErrInvoiceAlreadySent     = errors.Wrap(errors.ErrBadRequest, "invoice is already sent")
	ErrInvoicePaymentExceeds  = errors.Wrap(errors.ErrBadRequest, "payment amount exceeds balance due")
	ErrInvoiceDueDateInvalid  = errors.Wrap(errors.ErrBadRequest, "due date must be after issue date")
)

type Invoice struct {
	es.Aggregate
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

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Invoice)(nil)

func NewInvoice(id string) *Invoice {
	return &Invoice{
		Aggregate: es.NewAggregate(id, InvoiceAggregate),
	}
}

// Key implements registry.Registerable
func (Invoice) Key() string { return InvoiceAggregate }

// CreateInvoice initializes a new invoice
func (i *Invoice) CreateInvoice(
	invoiceNumber, orderID, customerID string,
	invoiceType InvoiceType,
	issueDate, dueDate time.Time,
	currency string,
	lines []InvoiceLine,
	subTotal, taxAmount, discountAmount, shippingAmount, totalAmount float64,
	paymentTerms string,
	billingAddress Address,
	taxID, poNumber, notes string,
	connectorID string,
) (ddd.Event, error) {
	// Validations
	if invoiceNumber == "" {
		return nil, ErrInvoiceNumberIsBlank
	}
	if customerID == "" {
		return nil, ErrInvoiceCustomerIsBlank
	}
	if orderID == "" {
		return nil, ErrInvoiceOrderIsBlank
	}
	if len(lines) == 0 {
		return nil, ErrInvoiceNoLines
	}
	if totalAmount < 0 {
		return nil, ErrInvoiceAmountNegative
	}
	if dueDate.Before(issueDate) {
		return nil, ErrInvoiceDueDateInvalid
	}

	i.AddEvent(InvoiceCreatedEvent, &InvoiceCreated{
		InvoiceNumber:  invoiceNumber,
		OrderID:        orderID,
		CustomerID:     customerID,
		Type:           invoiceType,
		IssueDate:      issueDate,
		DueDate:        dueDate,
		Currency:       currency,
		Lines:          lines,
		SubTotal:       subTotal,
		TaxAmount:      taxAmount,
		DiscountAmount: discountAmount,
		ShippingAmount: shippingAmount,
		TotalAmount:    totalAmount,
		PaymentTerms:   paymentTerms,
		BillingAddress: billingAddress,
		TaxID:          taxID,
		PONumber:       poNumber,
		Notes:          notes,
		ConnectorID:    connectorID,
	})

	return ddd.NewEvent(InvoiceCreatedEvent, i), nil
}

// ApproveInvoice approves the invoice
func (i *Invoice) ApproveInvoice(approvedBy string) (ddd.Event, error) {
	if i.Status == InvoiceStatusApproved {
		return nil, ErrInvoiceAlreadyApproved
	}
	if i.Status == InvoiceStatusVoided {
		return nil, ErrInvoiceAlreadyVoided
	}

	i.AddEvent(InvoiceApprovedEvent, &InvoiceApproved{
		ApprovedBy: approvedBy,
		ApprovedAt: time.Now(),
	})

	return ddd.NewEvent(InvoiceApprovedEvent, i), nil
}

// SendInvoice sends the invoice to the customer
func (i *Invoice) SendInvoice(sentTo []string, sentBy string) (ddd.Event, error) {
	if i.Status != InvoiceStatusApproved {
		return nil, ErrInvoiceNotApproved
	}
	if i.Status == InvoiceStatusSent {
		return nil, ErrInvoiceAlreadySent
	}

	i.AddEvent(InvoiceSentEvent, &InvoiceSent{
		SentTo: sentTo,
		SentBy: sentBy,
		SentAt: time.Now(),
	})

	return ddd.NewEvent(InvoiceSentEvent, i), nil
}

// VoidInvoice voids the invoice
func (i *Invoice) VoidInvoice(reason, voidedBy string) (ddd.Event, error) {
	if i.Status == InvoiceStatusVoided {
		return nil, ErrInvoiceAlreadyVoided
	}

	i.AddEvent(InvoiceVoidedEvent, &InvoiceVoided{
		Reason:   reason,
		VoidedBy: voidedBy,
		VoidedAt: time.Now(),
	})

	return ddd.NewEvent(InvoiceVoidedEvent, i), nil
}

// RecordPayment records a payment against the invoice
func (i *Invoice) RecordPayment(amount float64, paymentMethod, transactionID string, paymentDate time.Time) (ddd.Event, error) {
	if amount > i.BalanceDue {
		return nil, ErrInvoicePaymentExceeds
	}

	i.AddEvent(InvoicePaymentReceivedEvent, &InvoicePaymentReceived{
		Amount:        amount,
		PaymentMethod: paymentMethod,
		TransactionID: transactionID,
		PaymentDate:   paymentDate,
	})

	return ddd.NewEvent(InvoicePaymentReceivedEvent, i), nil
}

// LinkToERP links the invoice to an external ERP system
func (i *Invoice) LinkToERP(externalID string) (ddd.Event, error) {
	i.AddEvent(InvoiceLinkedToERPEvent, &InvoiceLinkedToERP{
		ExternalID: externalID,
		LinkedAt:   time.Now(),
	})

	return ddd.NewEvent(InvoiceLinkedToERPEvent, i), nil
}

// ApplyEvent implements es.EventApplier
func (i *Invoice) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {
	case *InvoiceCreated:
		i.InvoiceNumber = e.InvoiceNumber
		i.OrderID = e.OrderID
		i.CustomerID = e.CustomerID
		i.Type = e.Type
		i.Status = InvoiceStatusDraft
		i.IssueDate = e.IssueDate
		i.DueDate = e.DueDate
		i.Currency = e.Currency
		i.Lines = e.Lines
		i.SubTotal = e.SubTotal
		i.TaxAmount = e.TaxAmount
		i.DiscountAmount = e.DiscountAmount
		i.ShippingAmount = e.ShippingAmount
		i.TotalAmount = e.TotalAmount
		i.BalanceDue = e.TotalAmount
		i.PaymentTerms = e.PaymentTerms
		i.BillingAddress = e.BillingAddress
		i.TaxID = e.TaxID
		i.PONumber = e.PONumber
		i.Notes = e.Notes
		i.ConnectorID = e.ConnectorID
		i.CreatedAt = time.Now()
		i.UpdatedAt = time.Now()

	case *InvoiceApproved:
		i.Status = InvoiceStatusApproved
		i.ApprovedAt = &e.ApprovedAt
		i.UpdatedAt = e.ApprovedAt

	case *InvoiceSent:
		i.Status = InvoiceStatusSent
		i.SentAt = &e.SentAt
		i.UpdatedAt = e.SentAt

	case *InvoiceVoided:
		i.Status = InvoiceStatusVoided
		i.VoidedAt = &e.VoidedAt
		i.UpdatedAt = e.VoidedAt

	case *InvoicePaymentReceived:
		i.PaidAmount += e.Amount
		i.BalanceDue = i.TotalAmount - i.PaidAmount
		if i.BalanceDue <= 0 {
			i.Status = InvoiceStatusPaid
			now := e.PaymentDate
			i.PaidAt = &now
		}
		i.UpdatedAt = e.PaymentDate

	case *InvoiceLinkedToERP:
		i.ExternalID = e.ExternalID
		i.UpdatedAt = e.LinkedAt

	default:
		return errors.ErrInternal.Msgf(
			"%T received the event %s with unexpected payload %T",
			i, event.EventName(), e)
	}
	return nil
}

// ApplySnapshot implements es.Snapshotter
func (i *Invoice) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *InvoiceV1:
		i.InvoiceNumber = ss.InvoiceNumber
		i.OrderID = ss.OrderID
		i.CustomerID = ss.CustomerID
		i.Type = ss.Type
		i.Status = ss.Status
		i.IssueDate = ss.IssueDate
		i.DueDate = ss.DueDate
		i.Currency = ss.Currency
		i.Lines = ss.Lines
		i.SubTotal = ss.SubTotal
		i.TaxAmount = ss.TaxAmount
		i.DiscountAmount = ss.DiscountAmount
		i.ShippingAmount = ss.ShippingAmount
		i.TotalAmount = ss.TotalAmount
		i.PaidAmount = ss.PaidAmount
		i.BalanceDue = ss.BalanceDue
		i.PaymentTerms = ss.PaymentTerms
		i.PaymentMethod = ss.PaymentMethod
		i.BillingAddress = ss.BillingAddress
		i.TaxID = ss.TaxID
		i.PONumber = ss.PONumber
		i.Notes = ss.Notes
		i.Attachments = ss.Attachments
		i.ConnectorID = ss.ConnectorID
		i.ExternalID = ss.ExternalID
		i.CreatedAt = ss.CreatedAt
		i.UpdatedAt = ss.UpdatedAt
		i.ApprovedAt = ss.ApprovedAt
		i.SentAt = ss.SentAt
		i.VoidedAt = ss.VoidedAt
		i.PaidAt = ss.PaidAt

	default:
		return errors.ErrInternal.Msgf(
			"%T received the unexpected snapshot %T", i, snapshot)
	}
	return nil
}

// ToSnapshot implements es.Snapshotter
func (i Invoice) ToSnapshot() es.Snapshot {
	return InvoiceV1{
		InvoiceNumber:  i.InvoiceNumber,
		OrderID:        i.OrderID,
		CustomerID:     i.CustomerID,
		Type:           i.Type,
		Status:         i.Status,
		IssueDate:      i.IssueDate,
		DueDate:        i.DueDate,
		Currency:       i.Currency,
		Lines:          i.Lines,
		SubTotal:       i.SubTotal,
		TaxAmount:      i.TaxAmount,
		DiscountAmount: i.DiscountAmount,
		ShippingAmount: i.ShippingAmount,
		TotalAmount:    i.TotalAmount,
		PaidAmount:     i.PaidAmount,
		BalanceDue:     i.BalanceDue,
		PaymentTerms:   i.PaymentTerms,
		PaymentMethod:  i.PaymentMethod,
		BillingAddress: i.BillingAddress,
		TaxID:          i.TaxID,
		PONumber:       i.PONumber,
		Notes:          i.Notes,
		Attachments:    i.Attachments,
		ConnectorID:    i.ConnectorID,
		ExternalID:     i.ExternalID,
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
		ApprovedAt:     i.ApprovedAt,
		SentAt:         i.SentAt,
		VoidedAt:       i.VoidedAt,
		PaidAt:         i.PaidAt,
	}
}