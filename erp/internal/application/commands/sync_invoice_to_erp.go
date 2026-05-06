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

// SyncInvoiceToERP command syncs an invoice to the external ERP system
type SyncInvoiceToERP struct {
	InvoiceID   string
	ConnectorID string
}

// SyncInvoiceToERPHandler handles syncing invoices to ERP
type SyncInvoiceToERPHandler struct {
	invoices   es.AggregateRepository[*domain.Invoice]
	registry   erp.ConnectorRegistry
	repository domain.InvoiceSyncRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

// NewSyncInvoiceToERPHandler creates a new handler
func NewSyncInvoiceToERPHandler(
	invoices es.AggregateRepository[*domain.Invoice],
	registry erp.ConnectorRegistry,
	repository domain.InvoiceSyncRepository,
	publisher ddd.EventPublisher[ddd.Event],
) SyncInvoiceToERPHandler {
	return SyncInvoiceToERPHandler{
		invoices:   invoices,
		registry:   registry,
		repository: repository,
		publisher:  publisher,
	}
}

// SyncInvoiceToERP handles the SyncInvoiceToERP command
func (h SyncInvoiceToERPHandler) SyncInvoiceToERP(ctx context.Context, cmd SyncInvoiceToERP) error {
	// Load the invoice aggregate
	invoice, err := h.invoices.Load(ctx, cmd.InvoiceID)
	if err != nil {
		return fmt.Errorf("loading invoice: %w", err)
	}

	// Get connector
	connector, err := h.registry.GetConnector(cmd.ConnectorID)
	if err != nil {
		return fmt.Errorf("getting connector: %w", err)
	}

	// Create invoice sync record
	invoiceSync := &domain.InvoiceSync{
		ID:          generateID("inv_sync"),
		ConnectorID: cmd.ConnectorID,
		InvoiceID:   cmd.InvoiceID,
		Action:      "sync",
		Status:      domain.InvoiceSyncStatusPending,
		AttemptedAt: time.Now(),
		Payload: map[string]interface{}{
			"invoice_number": invoice.InvoiceNumber,
			"customer_id":    invoice.CustomerID,
			"total_amount":   invoice.TotalAmount,
			"status":         invoice.Status,
		},
	}

	if err := h.repository.Create(ctx, invoiceSync); err != nil {
		return fmt.Errorf("creating invoice sync record: %w", err)
	}

	// Convert to ERP invoice format
	erpInvoice := convertToERPInvoice(invoice)

	// Create or update invoice in ERP based on whether it has an external ID
	var externalID string
	if invoice.ExternalID == "" {
		// Create new invoice in ERP
		externalID, err = connector.CreateInvoice(ctx, erpInvoice)
		if err != nil {
			updateSyncError(ctx, h.repository, invoiceSync, err)
			return fmt.Errorf("creating invoice in ERP: %w", err)
		}
	} else {
		// Update existing invoice in ERP
		erpInvoice.InvoiceID = invoice.ExternalID
		erpInvoice.Status = string(invoice.Status)
		err = connector.UpdateInvoice(ctx, invoice.ExternalID, erpInvoice)
		if err != nil {
			updateSyncError(ctx, h.repository, invoiceSync, err)
			return fmt.Errorf("updating invoice in ERP: %w", err)
		}
		externalID = invoice.ExternalID
	}

	// Update sync status to completed
	invoiceSync.Status = domain.InvoiceSyncStatusCompleted
	invoiceSync.ExternalID = externalID
	invoiceSync.CompletedAt = ptrTime(time.Now())
	if err := h.repository.Update(ctx, invoiceSync); err != nil {
		return fmt.Errorf("updating invoice sync record: %w", err)
	}

	// If this is a new sync, link the invoice to the external ID
	if invoice.ExternalID == "" && externalID != "" {
		event, err := invoice.LinkToERP(externalID)
		if err != nil {
			return fmt.Errorf("linking invoice to ERP: %w", err)
		}

		// Save the updated invoice
		if err := h.invoices.Save(ctx, invoice); err != nil {
			return fmt.Errorf("saving invoice: %w", err)
		}

		// Publish the domain event
		h.publisher.Publish(ctx, event)
	}

	return nil
}

func convertToERPInvoice(inv *domain.Invoice) *erp.InvoicePayload {
	erpLines := make([]erp.InvoiceLine, len(inv.Lines))
	for i, line := range inv.Lines {
		erpLines[i] = erp.InvoiceLine{
			SKU:         line.SKU,
			Description: line.Description,
			Quantity:    float64(line.Quantity),
			UnitPrice:   line.UnitPrice,
			TaxRate:     line.TaxRate,
			TaxAmount:   line.TaxAmount,
			LineTotal:   line.LineTotal,
		}
	}

	return &erp.InvoicePayload{
		InvoiceID:     inv.ID(),
		InvoiceNumber: inv.InvoiceNumber,
		CustomerID:    inv.CustomerID,
		OrderID:       inv.OrderID,
		IssueDate:     inv.IssueDate,
		DueDate:       inv.DueDate,
		Currency:      inv.Currency,
		Lines:         erpLines,
		SubTotal:      inv.SubTotal,
		TaxAmount:     inv.TaxAmount,
		TotalAmount:   inv.TotalAmount,
		Status:        string(inv.Status),
		PaymentTerms:  inv.PaymentTerms,
		Notes:         inv.Notes,
	}
}

func updateSyncError(ctx context.Context, repo domain.InvoiceSyncRepository, sync *domain.InvoiceSync, err error) {
	sync.Status = domain.InvoiceSyncStatusFailed
	sync.Error = err.Error()
	sync.CompletedAt = ptrTime(time.Now())
	repo.Update(ctx, sync)
}

