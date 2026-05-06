package commands

import (
	"context"
	"fmt"
	"time"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/erp"
	"middleman/internal/es"
)

// CreateInvoice command creates a new invoice in the system and syncs to ERP
type CreateInvoice struct {
	InvoiceID      string
	InvoiceNumber  string
	OrderID        string
	CustomerID     string
	Type           domain.InvoiceType
	IssueDate      time.Time
	DueDate        time.Time
	Currency       string
	Lines          []domain.InvoiceLine
	SubTotal       float64
	TaxAmount      float64
	DiscountAmount float64
	ShippingAmount float64
	TotalAmount    float64
	PaymentTerms   string
	BillingAddress erp.Address
	TaxID          string
	PONumber       string
	Notes          string
	ConnectorID    string
}

// CreateInvoiceHandler handles creating new invoices
type CreateInvoiceHandler struct {
	invoices  es.AggregateRepository[*domain.Invoice]
	publisher ddd.EventPublisher[ddd.Event]
}

// NewCreateInvoiceHandler creates a new handler
func NewCreateInvoiceHandler(
	invoices es.AggregateRepository[*domain.Invoice],
	publisher ddd.EventPublisher[ddd.Event],
) CreateInvoiceHandler {
	return CreateInvoiceHandler{
		invoices:  invoices,
		publisher: publisher,
	}
}

// CreateInvoice handles the CreateInvoice command
func (h CreateInvoiceHandler) CreateInvoice(ctx context.Context, cmd CreateInvoice) error {
	// Load the invoice aggregate (will create new if doesn't exist)
	invoice, err := h.invoices.Load(ctx, cmd.InvoiceID)
	if err != nil {
		return fmt.Errorf("loading invoice: %w", err)
	}

	// Convert erp.Address to domain.Address
	billingAddress := domain.Address{
		Street:     cmd.BillingAddress.Street,
		City:       cmd.BillingAddress.City,
		State:      cmd.BillingAddress.State,
		PostalCode: cmd.BillingAddress.PostalCode,
		Country:    cmd.BillingAddress.Country,
	}

	event, err := invoice.CreateInvoice(
		cmd.InvoiceNumber,
		cmd.OrderID,
		cmd.CustomerID,
		cmd.Type,
		cmd.IssueDate,
		cmd.DueDate,
		cmd.Currency,
		cmd.Lines,
		cmd.SubTotal,
		cmd.TaxAmount,
		cmd.DiscountAmount,
		cmd.ShippingAmount,
		cmd.TotalAmount,
		cmd.PaymentTerms,
		billingAddress,
		cmd.TaxID,
		cmd.PONumber,
		cmd.Notes,
		cmd.ConnectorID,
	)
	if err != nil {
		return fmt.Errorf("creating invoice: %w", err)
	}

	// Save the invoice aggregate
	if err := h.invoices.Save(ctx, invoice); err != nil {
		return fmt.Errorf("saving invoice: %w", err)
	}

	// Publish the domain event
	return h.publisher.Publish(ctx, event)
}
