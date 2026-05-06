package commands

import (
	"context"
	"fmt"
	"time"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// RecordInvoicePayment command records a payment against an invoice
type RecordInvoicePayment struct {
	InvoiceID     string
	Amount        float64
	PaymentMethod string
	TransactionID string
	PaymentDate   time.Time
}

// RecordInvoicePaymentHandler handles recording invoice payments
type RecordInvoicePaymentHandler struct {
	invoices  es.AggregateRepository[*domain.Invoice]
	publisher ddd.EventPublisher[ddd.Event]
}

// NewRecordInvoicePaymentHandler creates a new handler
func NewRecordInvoicePaymentHandler(
	invoices es.AggregateRepository[*domain.Invoice],
	publisher ddd.EventPublisher[ddd.Event],
) RecordInvoicePaymentHandler {
	return RecordInvoicePaymentHandler{
		invoices:  invoices,
		publisher: publisher,
	}
}

// RecordInvoicePayment handles the RecordInvoicePayment command
func (h RecordInvoicePaymentHandler) RecordInvoicePayment(ctx context.Context, cmd RecordInvoicePayment) error {
	// Load the invoice aggregate
	invoice, err := h.invoices.Load(ctx, cmd.InvoiceID)
	if err != nil {
		return fmt.Errorf("loading invoice: %w", err)
	}

	// Record the payment
	event, err := invoice.RecordPayment(cmd.Amount, cmd.PaymentMethod, cmd.TransactionID, cmd.PaymentDate)
	if err != nil {
		return fmt.Errorf("recording payment: %w", err)
	}

	// Save the updated invoice
	if err := h.invoices.Save(ctx, invoice); err != nil {
		return fmt.Errorf("saving invoice: %w", err)
	}

	// Publish the domain event
	return h.publisher.Publish(ctx, event)
}