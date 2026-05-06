package handlers

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"middleman/erp/internal/application"
	"middleman/erp/internal/application/commands"
	"middleman/erp/internal/domain"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/erp"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/ordering/orderingpb"
	"middleman/payments/paymentspb"
)

// integrationHandlers implements ddd.EventHandler
type integrationHandlers[T ddd.Event] struct {
	app application.App
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(
	reg registry.Registry,
	app application.App,
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(
		reg,
		integrationHandlers[ddd.Event]{
			app: app,
		},
		zerolog.Logger{},
		mws...,
	)
}

func RegisterIntegrationEventHandlers(
	subscriber am.MessageSubscriber,
	handlers am.MessageHandler,
) (err error) {
	// Subscribe to ordering events
	_, err = subscriber.Subscribe(
		orderingpb.OrderAggregateChannel,
		handlers,
		am.MessageFilter{
			orderingpb.OrderReadiedEvent,
			orderingpb.OrderCanceledEvent,
		},
		am.GroupName("erp-service"),
	)
	if err != nil {
		return err
	}

	// Subscribe to payment events
	_, err = subscriber.Subscribe(
		paymentspb.PaymentAggregateChannel,
		handlers,
		am.MessageFilter{},
		am.GroupName("erp-service"),
	)
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling integration event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled integration event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling integration event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	// Order events
	case orderingpb.OrderReadiedEvent:
		return h.onOrderReadied(ctx, event)
	case orderingpb.OrderCanceledEvent:
		return h.onOrderCanceled(ctx, event)
	// Payment events
	case paymentspb.InvoicePaidEvent:
		return h.onInvoicePaid(ctx, event)
	}

	return nil
}

// onOrderReadied creates an invoice when an order is readied for payment
func (h integrationHandlers[T]) onOrderReadied(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*orderingpb.OrderReadied)
	// For OrderReadied, we only have the total amount, not individual items
	// We'll create a single line item for the order
	totalAmount := float64(payload.Total) / 100.0 // Convert cents to dollars
	subTotal := totalAmount / 1.1                 // Assuming 10% tax included
	taxAmount := totalAmount - subTotal

	lines := []domain.InvoiceLine{{
		SKU:         "ORDER-" + payload.Id,
		ProductName: "Order Total",
		Description: "Order #" + payload.Id,
		Quantity:    1,
		UnitPrice:   subTotal,
		LineTotal:   subTotal,
	}}

	// Create the invoice
	return h.app.CreateInvoice(ctx, commands.CreateInvoice{
		InvoiceID:      uuid.New().String(),
		InvoiceNumber:  "INV-" + payload.Id,
		OrderID:        payload.Id,
		CustomerID:     payload.UserCustomerId,
		Type:           domain.InvoiceTypeStandard,
		IssueDate:      time.Now(),
		DueDate:        time.Now().AddDate(0, 0, 30), // 30 days payment terms
		Currency:       "USD",                        // Should come from order
		Lines:          lines,
		SubTotal:       subTotal,
		TaxAmount:      taxAmount,
		ShippingAmount: 0, // Not available in OrderReadied
		TotalAmount:    totalAmount,
		PaymentTerms:   "Net 30",
		BillingAddress: erp.Address{
			// Address fields not available in OrderReadied event
		},
		ConnectorID: h.getConnectorID(ctx),
	})
}

// onOrderCanceled voids the invoice when an order is canceled
func (h integrationHandlers[T]) onOrderCanceled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*orderingpb.OrderCanceled)
	// Find invoice by order ID (would need a query for this)
	// For now, assuming we have the invoice ID
	invoiceID := h.findInvoiceByOrderID(ctx, payload.Id)
	if invoiceID == "" {
		return nil // No invoice found
	}

	return h.app.VoidInvoice(ctx, commands.VoidInvoice{
		InvoiceID: invoiceID,
		Reason:    "Order canceled",
		VoidedBy:  "system",
	})
}

// Helper methods
func (h integrationHandlers[T]) getConnectorID(ctx context.Context) string {
	// Try to get from context (e.g., tenant-specific)
	if connectorID, ok := ctx.Value("erp_connector_id").(string); ok {
		return connectorID
	}
	// Default connector ID (could be from configuration)
	return "default_connector"
}

func (h integrationHandlers[T]) findInvoiceByOrderID(ctx context.Context, orderID string) string {
	// In a real implementation, this would query the invoice repository
	// You might want to add a query to the application layer for this
	// For now, return empty string
	return ""
}

// onInvoicePaid records payment when an invoice is paid
func (h integrationHandlers[T]) onInvoicePaid(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*paymentspb.InvoicePaid)
	// Find invoice by order ID since InvoicePaid only has order ID
	invoiceID := h.findInvoiceByOrderID(ctx, payload.OrderId)
	if invoiceID == "" {
		return nil // No invoice found for this order
	}
	return h.app.RecordInvoicePayment(ctx, commands.RecordInvoicePayment{
		InvoiceID:     invoiceID,
		Amount:        100.0,  // Default amount, would need to query invoice for actual amount
		PaymentMethod: "card", // Default method
		TransactionID: payload.Id,
		PaymentDate:   time.Now(),
	})
}
