package rest

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"middleman/internal/di"
	"middleman/sap/internal/application"
	"middleman/sap/internal/constants"
	"middleman/sap/internal/domain"
	"middleman/sap/internal/sap"
)

// RegisterWebhookRoute sets up a plain HTTP route for /api/sap/webhook
// It scopes the DI container for each request to ensure transactional context.
func RegisterWebhookRoute(container di.Container, mux *chi.Mux, app application.SAPConnectorDomain, sapClient *sap.EnhancedSAPClient) {
	mux.Post("/api/sap/webhook", func(w http.ResponseWriter, r *http.Request) {
		// Scope the DI container – opens a DB transaction and other scoped dependencies.
		ctx := container.Scoped(r.Context())
		var err error
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				log.Error().Interface("panic", p).Msg("panic in SAP webhook; tx rolled back")
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
				log.Error().Err(err).Msg("SAP webhook handler failed; tx rolled back")
			} else {
				if cerr := tx.Commit(); cerr != nil {
					log.Error().Err(cerr).Msg("commit failed in SAP webhook handler")
				}
			}
		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

		// Replace request context with scoped ctx for downstream readers
		r = r.WithContext(ctx)

		// 1) Read raw body
		bodyBytes, errRead := io.ReadAll(r.Body)
		if errRead != nil {
			err = errRead
			log.Error().Err(err).Msg("failed to read request body")
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		// 2) Get signature from headers (SAP might use different header name)
		sig := r.Header.Get("X-SAP-Signature")
		if sig == "" {
			sig = r.Header.Get("SAP-Signature")
		}

		// 3) Validate signature with enhanced security
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		if errVal := sapClient.ValidateWebhookWithSecurity(r); errVal != nil {
			err = errVal
			log.Error().Err(err).Msg("[RegisterWebhookRoute] SAP webhook validation failed")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// 4) Store raw webhook event
		webhookEvent := &domain.WebhookEvent{
			ID:         uuid.New().String(),
			EventID:    r.Header.Get("X-SAP-Event-ID"),
			EventType:  r.Header.Get("X-SAP-Event-Type"),
			Source:     "SAP",
			Signature:  sig,
			Payload:    bodyBytes,
			ReceivedAt: time.Now(),
			Status:     "received",
		}

		if errStore := app.StoreWebhookEvent(ctx, webhookEvent); errStore != nil {
			err = errStore
			log.Error().Err(err).Msg("[RegisterWebhookRoute] failed to store webhook event")
			// Continue processing even if storage fails
		}

		// 5) Parse event - try different formats
		event, contentType, errParse := parseEvent(bodyBytes, r.Header.Get("Content-Type"))
		if errParse != nil {
			err = errParse
			log.Error().Err(err).Msg("[RegisterWebhookRoute] failed to parse SAP webhook event")
			http.Error(w, "failed to parse event", http.StatusBadRequest)
			return
		}

		// 6) Handle event based on type
		switch contentType {
		case "idoc":
			err = handleIDocEvent(ctx, app, event, webhookEvent.ID)
		case "json":
			err = handleJSONEvent(ctx, app, event.(*sap.SAPEvent), webhookEvent.ID)
		default:
			log.Warn().Str("contentType", contentType).Msg("Unknown content type")
		}

		if err != nil {
			// Update webhook event status to failed
			_ = app.UpdateWebhookEventStatus(ctx, webhookEvent.ID, "failed", err.Error())
		} else {
			// Update webhook event status to processed
			_ = app.UpdateWebhookEventStatus(ctx, webhookEvent.ID, "processed", "")
		}

		// 7) Respond 200 so SAP won't keep retrying
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

// parseEvent attempts to parse the event from different formats
func parseEvent(bodyBytes []byte, contentType string) (interface{}, string, error) {
	// Try JSON first
	var jsonEvent sap.SAPEvent
	if err := json.Unmarshal(bodyBytes, &jsonEvent); err == nil && jsonEvent.Type != "" {
		return &jsonEvent, "json", nil
	}

	// Try XML/IDoc
	var idoc sap.IDoc
	if err := xml.Unmarshal(bodyBytes, &idoc); err == nil && idoc.EDI_DC40.IDOCTYP != "" {
		return &idoc, "idoc", nil
	}

	return nil, "", fmt.Errorf("unable to parse event as JSON or IDoc")
}

// handleIDocEvent processes IDoc format events
func handleIDocEvent(ctx context.Context, app application.SAPConnectorDomain, idocInterface interface{}, webhookEventID string) error {
	idoc := idocInterface.(*sap.IDoc)

	log.Info().
		Str("idocType", idoc.EDI_DC40.IDOCTYP).
		Str("docNum", idoc.EDI_DC40.DOCNUM).
		Msg("Processing IDoc event")

	switch sap.IDocType(idoc.EDI_DC40.IDOCTYP) {
	case sap.IDocTypeMATMAS:
		// Material Master update
		return app.ProcessMaterialMaster(ctx, application.ProcessMaterialMasterCommand{
			IDocData:       idocInterface,
			CorrelationID:  idoc.EDI_DC40.DOCNUM,
			WebhookEventID: webhookEventID,
		})

	case sap.IDocTypeINVCON:
		// Inventory update
		return app.ProcessInventoryUpdate(ctx, application.ProcessInventoryUpdateCommand{
			IDocData:       idocInterface,
			CorrelationID:  idoc.EDI_DC40.DOCNUM,
			WebhookEventID: webhookEventID,
		})

	case sap.IDocTypeCOND_A:
		// Pricing update
		return app.ProcessPricingUpdate(ctx, application.ProcessPricingUpdateCommand{
			IDocData:       idocInterface,
			CorrelationID:  idoc.EDI_DC40.DOCNUM,
			WebhookEventID: webhookEventID,
		})

	default:
		log.Warn().
			Str("idocType", idoc.EDI_DC40.IDOCTYP).
			Msg("Unhandled IDoc type")
		return nil
	}
}

// handleJSONEvent processes JSON format events
func handleJSONEvent(ctx context.Context, app application.SAPConnectorDomain, event *sap.SAPEvent, webhookEventID string) error {
	log.Info().
		Str("eventType", string(event.Type)).
		Str("eventId", event.ID).
		Msg("Processing JSON event")

	switch event.Type {
	case sap.EventTypeProductCreated, sap.EventTypeProductUpdated:
		return app.ProcessProductEvent(ctx, application.ProcessProductEventCommand{
			Event:          event,
			WebhookEventID: webhookEventID,
		})

	case sap.EventTypeStockUpdated:
		return app.ProcessStockEvent(ctx, application.ProcessStockEventCommand{
			Event:          event,
			WebhookEventID: webhookEventID,
		})

	case sap.EventTypePriceUpdated:
		return app.ProcessPriceEvent(ctx, application.ProcessPriceEventCommand{
			Event:          event,
			WebhookEventID: webhookEventID,
		})

	case sap.EventTypeProductDeleted:
		return app.ProcessProductDeletedEvent(ctx, application.ProcessProductDeletedEventCommand{
			Event:          event,
			WebhookEventID: webhookEventID,
		})

	default:
		log.Warn().
			Str("eventType", string(event.Type)).
			Msg("Unhandled event type")
		return nil
	}
}
