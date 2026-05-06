package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// SendInvoice command sends an approved invoice to customers
type SendInvoice struct {
	InvoiceID string
	SentTo    []string // Email addresses
	SentBy    string
}

// SendInvoiceHandler handles sending invoices
type SendInvoiceHandler struct {
	invoices  es.AggregateRepository[*domain.Invoice]
	publisher ddd.EventPublisher[ddd.Event]
}

// NewSendInvoiceHandler creates a new handler
func NewSendInvoiceHandler(
	invoices es.AggregateRepository[*domain.Invoice],
	publisher ddd.EventPublisher[ddd.Event],
) SendInvoiceHandler {
	return SendInvoiceHandler{
		invoices:  invoices,
		publisher: publisher,
	}
}

// SendInvoice handles the SendInvoice command
func (h SendInvoiceHandler) SendInvoice(ctx context.Context, cmd SendInvoice) error {
	// Load the invoice aggregate
	invoice, err := h.invoices.Load(ctx, cmd.InvoiceID)
	if err != nil {
		return fmt.Errorf("loading invoice: %w", err)
	}

	// Send the invoice
	event, err := invoice.SendInvoice(cmd.SentTo, cmd.SentBy)
	if err != nil {
		return fmt.Errorf("sending invoice: %w", err)
	}

	// Save the updated invoice
	if err := h.invoices.Save(ctx, invoice); err != nil {
		return fmt.Errorf("saving invoice: %w", err)
	}

	// Publish the domain event
	return h.publisher.Publish(ctx, event)
}