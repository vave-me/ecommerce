package commands

import (
	"context"
	"fmt"

	"middleman/erp/internal/domain"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

// ApproveInvoice command approves an invoice for sending
type ApproveInvoice struct {
	InvoiceID  string
	ApprovedBy string
}

// ApproveInvoiceHandler handles approving invoices
type ApproveInvoiceHandler struct {
	invoices  es.AggregateRepository[*domain.Invoice]
	publisher ddd.EventPublisher[ddd.Event]
}

// NewApproveInvoiceHandler creates a new handler
func NewApproveInvoiceHandler(
	invoices es.AggregateRepository[*domain.Invoice],
	publisher ddd.EventPublisher[ddd.Event],
) ApproveInvoiceHandler {
	return ApproveInvoiceHandler{
		invoices:  invoices,
		publisher: publisher,
	}
}

// ApproveInvoice handles the ApproveInvoice command
func (h ApproveInvoiceHandler) ApproveInvoice(ctx context.Context, cmd ApproveInvoice) error {
	// Load the invoice aggregate
	invoice, err := h.invoices.Load(ctx, cmd.InvoiceID)
	if err != nil {
		return fmt.Errorf("loading invoice: %w", err)
	}

	// Approve the invoice
	event, err := invoice.ApproveInvoice(cmd.ApprovedBy)
	if err != nil {
		return fmt.Errorf("approving invoice: %w", err)
	}

	// Save the updated invoice
	if err := h.invoices.Save(ctx, invoice); err != nil {
		return fmt.Errorf("saving invoice: %w", err)
	}

	// Publish the domain event
	return h.publisher.Publish(ctx, event)
}
