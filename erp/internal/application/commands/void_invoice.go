package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// VoidInvoice command voids an invoice
type VoidInvoice struct {
	InvoiceID string
	Reason    string
	VoidedBy  string
}

// VoidInvoiceHandler handles voiding invoices
type VoidInvoiceHandler struct {
	invoices  es.AggregateRepository[*domain.Invoice]
	publisher ddd.EventPublisher[ddd.Event]
}

// NewVoidInvoiceHandler creates a new handler
func NewVoidInvoiceHandler(
	invoices es.AggregateRepository[*domain.Invoice],
	publisher ddd.EventPublisher[ddd.Event],
) VoidInvoiceHandler {
	return VoidInvoiceHandler{
		invoices:  invoices,
		publisher: publisher,
	}
}

// VoidInvoice handles the VoidInvoice command
func (h VoidInvoiceHandler) VoidInvoice(ctx context.Context, cmd VoidInvoice) error {
	// Load the invoice aggregate
	invoice, err := h.invoices.Load(ctx, cmd.InvoiceID)
	if err != nil {
		return fmt.Errorf("loading invoice: %w", err)
	}

	// Void the invoice
	event, err := invoice.VoidInvoice(cmd.Reason, cmd.VoidedBy)
	if err != nil {
		return fmt.Errorf("voiding invoice: %w", err)
	}

	// Save the updated invoice
	if err := h.invoices.Save(ctx, invoice); err != nil {
		return fmt.Errorf("saving invoice: %w", err)
	}

	// Publish the domain event
	return h.publisher.Publish(ctx, event)
}