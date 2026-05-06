package rest

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"middleman/erp/internal/domain"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"middleman/erp/internal/application"
	"middleman/erp/internal/application/commands"
	"middleman/erp/internal/constants"
	"middleman/internal/di"
	"middleman/internal/erp"
)

// RegisterWebhookRoutes sets up webhook routes for ERP systems
// It follows a similar pattern to payments, mapping webhook events directly to application commands
func RegisterWebhookRoutes(container di.Container, mux *chi.Mux, app application.App, registry erp.ConnectorRegistry) {
	// Generic webhook endpoint for any ERP connector
	mux.Post("/api/erp/webhook/{connectorId}", func(w http.ResponseWriter, r *http.Request) {
		// Scope the DI container – opens a DB transaction and other scoped dependencies.
		ctx := container.Scoped(r.Context())
		var err error
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				log.Error().Interface("panic", p).Msg("panic in webhook; tx rolled back")
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
				log.Error().Err(err).Msg("webhook handler failed; tx rolled back")
			} else {
				if cerr := tx.Commit(); cerr != nil {
					log.Error().Err(cerr).Msg("commit failed in webhook handler")
				}
			}
		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

		// replace request context with scoped ctx for downstream readers
		r = r.WithContext(ctx)

		// Get connector ID from URL
		connectorID := chi.URLParam(r, "connectorId")
		if connectorID == "" {
			err = fmt.Errorf("connector ID is required")
			http.Error(w, "connector ID is required", http.StatusBadRequest)
			return
		}

		// Verify connector exists
		_, errGet := registry.GetConnector(connectorID)
		if errGet != nil {
			err = errGet
			log.Error().Err(err).Str("connectorId", connectorID).Msg("connector not found")
			http.Error(w, "connector not found", http.StatusNotFound)
			return
		}

		// 1) read raw body
		bodyBytes, errRead := io.ReadAll(r.Body)
		if errRead != nil {
			err = errRead
			log.Error().Err(err).Msg("failed to read request body")
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		// 2) extract headers
		headers := make(map[string]string)
		for key, values := range r.Header {
			if len(values) > 0 {
				headers[key] = values[0]
			}
		}

		// 3) extract event metadata from headers
		eventID := r.Header.Get("X-Event-ID")
		if eventID == "" {
			eventID = r.Header.Get("X-Webhook-ID")
		}
		//eventType := r.Header.Get("X-Event-Type")

		// 4) extract signature (handled by ProcessWebhook command)
		signature := extractSignature(r)

		// 5) Get connector to validate and parse webhook
		connector, errConn := registry.GetConnector(connectorID)
		if errConn != nil {
			err = errConn
			log.Error().Err(err).Str("connectorId", connectorID).Msg("failed to get connector")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// 6) Validate webhook signature
		if errVal := connector.ValidateWebhook(bodyBytes, signature); errVal != nil {
			err = errVal
			log.Error().Err(err).Msg("webhook signature invalid")
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		// 7) Parse webhook into canonical event
		canonicalEvent, errParse := connector.ParseWebhook(bodyBytes)
		if errParse != nil {
			err = errParse
			log.Error().Err(err).Msg("failed to parse webhook")
			http.Error(w, "failed to parse webhook", http.StatusBadRequest)
			return
		}

		// 8) Handle event based on type - map to specific application commands
		if err = handleWebhookEvent(ctx, app, connectorID, canonicalEvent); err != nil {
			log.Error().Err(err).
				Str("connector_id", connectorID).
				Str("event_type", string(canonicalEvent.EventType)).
				Str("event_id", canonicalEvent.EventID).
				Msg("failed to handle webhook event")
			http.Error(w, "webhook processing failed", http.StatusInternalServerError)
			return
		}

		log.Info().
			Str("connector_id", connectorID).
			Str("event_id", canonicalEvent.EventID).
			Str("event_type", string(canonicalEvent.EventType)).
			Msg("webhook processed successfully")

		// 9) respond 200 so ERP won't keep retrying
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

// handleWebhookEvent maps canonical events to specific application commands
func handleWebhookEvent(ctx context.Context, app application.App, connectorID string, event *erp.CanonicalEvent) error {
	switch event.EventType {
	// Product events
	case erp.EventTypeProductCreated, erp.EventTypeProductMasterUpdated:
		// Trigger product sync for this specific product
		if product, ok := event.Payload.(*erp.ProductPayload); ok {
			// Could create a more targeted sync command in the future
			log.Info().
				Str("sku", product.SKU).
				Str("name", product.Name).
				Msg("product webhook received")
			// For now, just log - actual sync would be handled by scheduled sync
		}

	// Stock events
	case erp.EventTypeStockLevelUpdated:
		if stock, ok := event.Payload.(*erp.StockPayload); ok {
			log.Info().
				Str("sku", stock.SKU).
				Str("location", stock.LocationID).
				Int("quantity", stock.Quantity).
				Msg("stock level webhook received")
			// Could trigger immediate stock update for this SKU
		}

	// Price events
	case erp.EventTypePriceUpdated:
		if price, ok := event.Payload.(*erp.PricePayload); ok {
			log.Info().
				Str("sku", price.SKU).
				Float64("price", price.Price).
				Str("currency", price.Currency).
				Msg("price update webhook received")
			// Could trigger immediate price update for this SKU
		}

	// Order events
	case erp.EventTypeOrderCreated, erp.EventTypeOrderUpdated:
		if order, ok := event.Payload.(*erp.OrderPayload); ok {
			log.Info().
				Str("order_id", order.OrderID).
				Str("status", order.Status).
				Msg("order webhook received")
			// Could sync order status back to our system
		}

	case erp.EventTypeOrderShipped:
		if order, ok := event.Payload.(*erp.OrderPayload); ok {
			// When order is shipped in ERP, we might want to:
			// 1. Update local order status
			// 2. Send shipping notification
			// 3. Update inventory reservations
			log.Info().
				Str("order_id", order.OrderID).
				Msg("order shipped webhook received")
		}

	// Customer events
	case erp.EventTypeCustomerCreated, erp.EventTypeCustomerUpdated:
		if customer, ok := event.Payload.(*erp.CustomerPayload); ok {
			log.Info().
				Str("customer_id", customer.CustomerID).
				Str("email", customer.Email).
				Msg("customer webhook received")
			// Could sync customer data
		}

	// Invoice events
	case erp.EventTypeInvoiceCreated:
		if invoice, ok := event.Payload.(*erp.InvoicePayload); ok {
			// When invoice is created in ERP, we might want to:
			// 1. Create local invoice record
			// 2. Send invoice notification to customer
			return app.CreateInvoice(ctx, commands.CreateInvoice{
				InvoiceID:     invoice.InvoiceID,
				InvoiceNumber: invoice.InvoiceNumber,
				OrderID:       invoice.OrderID,
				CustomerID:    invoice.CustomerID,
				Type:          domain.InvoiceTypeStandard,
				IssueDate:     invoice.IssueDate,
				DueDate:       invoice.DueDate,
				Currency:      invoice.Currency,
				Lines:         convertInvoiceLines(invoice.Lines),
				SubTotal:      invoice.SubTotal,
				TaxAmount:     invoice.TaxAmount,
				TotalAmount:   invoice.TotalAmount,
				PaymentTerms:  invoice.PaymentTerms,
				ConnectorID:   connectorID,
			})
		}

	case erp.EventTypeInvoiceUpdated:
		if invoice, ok := event.Payload.(*erp.InvoicePayload); ok {
			// Handle invoice updates - might need to update status, add payments, etc.
			switch invoice.Status {
			case "paid":
				// Record payment if invoice is marked as paid in ERP
				return app.RecordInvoicePayment(ctx, commands.RecordInvoicePayment{
					InvoiceID:     invoice.InvoiceID,
					Amount:        invoice.TotalAmount,
					PaymentMethod: "erp_sync",
					TransactionID: fmt.Sprintf("erp_%s_%s", connectorID, event.EventID),
					PaymentDate:   time.Now(),
				})
			case "void", "cancelled":
				// Void invoice if cancelled in ERP
				return app.VoidInvoice(ctx, commands.VoidInvoice{
					InvoiceID: invoice.InvoiceID,
					Reason:    fmt.Sprintf("Voided in ERP: %s", invoice.Status),
					VoidedBy:  "erp_sync",
				})
			}
		}

	default:
		// Log unhandled event types
		log.Warn().
			Str("event_type", string(event.EventType)).
			Str("event_id", event.EventID).
			Msg("unhandled webhook event type")
	}

	return nil
}

// convertInvoiceLines converts ERP invoice lines to domain invoice lines
func convertInvoiceLines(erpLines []erp.InvoiceLine) []domain.InvoiceLine {
	lines := make([]domain.InvoiceLine, len(erpLines))
	for i, line := range erpLines {
		lines[i] = domain.InvoiceLine{
			SKU:         line.SKU,
			Description: line.Description,
			Quantity:    int(line.Quantity),
			UnitPrice:   line.UnitPrice,
			TaxRate:     line.TaxRate,
			TaxAmount:   line.TaxAmount,
			LineTotal:   line.LineTotal,
		}
	}
	return lines
}

// extractSignature extracts the webhook signature from headers
func extractSignature(r *http.Request) string {
	// Check common signature headers
	signatureHeaders := []string{
		"X-Webhook-Signature",
		"X-Signature",
		"X-Hub-Signature",
		"X-Hub-Signature-256",
		// ERP specific headers
		"X-SAP-Signature",
		"SAP-Signature",
		"X-Odoo-Signature",
		"X-MS-Signature",
		"X-Dynamics-Signature",
		"X-NetSuite-Signature",
		"X-ERPNext-Signature",
		"X-Frappe-Signature",
	}

	for _, header := range signatureHeaders {
		if sig := r.Header.Get(header); sig != "" {
			return sig
		}
	}

	return ""
}

// extractEventID attempts to extract event ID from various headers
func extractEventID(r *http.Request) string {
	eventIDHeaders := []string{
		"X-Event-ID",
		"X-Webhook-ID",
		"X-Request-ID",
		"X-Message-ID",
		// ERP specific
		"X-Odoo-Event-ID",
		"X-MS-Event-ID",
		"X-NetSuite-Event-ID",
		"X-SAP-Message-ID",
	}

	for _, header := range eventIDHeaders {
		if id := r.Header.Get(header); id != "" {
			return id
		}
	}

	// Generate one if not provided
	return fmt.Sprintf("webhook_%d", time.Now().UnixNano())
}

// extractEventType attempts to extract event type from headers
func extractEventType(r *http.Request) string {
	eventTypeHeaders := []string{
		"X-Event-Type",
		"X-Webhook-Event",
		"X-Event-Name",
		// ERP specific
		"X-Odoo-Event",
		"X-MS-Event-Type",
		"X-NetSuite-Event-Type",
		"X-SAP-Event-Type",
	}

	for _, header := range eventTypeHeaders {
		if eventType := r.Header.Get(header); eventType != "" {
			return eventType
		}
	}

	return ""
}
