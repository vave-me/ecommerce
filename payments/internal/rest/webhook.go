package rest

import (
	"database/sql"
	"io"
	"net/http"

	"middleman/internal/di"
	"middleman/internal/stripe"
	"middleman/payments/internal/application"
	"middleman/payments/internal/constants"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	stripelib "github.com/stripe/stripe-go/v81"
)

// RegisterWebhookRoute sets up a plain HTTP route for /api/payments/webhook
// It scopes the DI container for each request to ensure transactional context.
func RegisterWebhookRoute(container di.Container, mux *chi.Mux, app application.PaymentDomain, stripeClient *stripe.StripeClient) {
	mux.Post("/api/payments/webhook", func(w http.ResponseWriter, r *http.Request) {
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

		// 1) read raw body
		bodyBytes, errRead := io.ReadAll(r.Body)
		if errRead != nil {
			err = errRead
			log.Error().Err(err).Msg("failed to read request body")
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		// Omit raw body debug dump in production.

		// 2) get signature
		sig := r.Header.Get("Stripe-Signature")

		// 3) validate signature
		if errVal := stripeClient.ValidateSign(bodyBytes, sig, stripeClient.WebhookSecret); errVal != nil {
			err = errVal
			log.Error().Err(err).Msg("[RegisterWebhookRoute] webhook signature invalid")
			// return 2xx if you want Stripe to stop retrying, or 4xx if you want retries
			http.Error(w, "invalid signature", http.StatusBadRequest)
			return
		}

		// 4) parse event
		event, errParse := stripeClient.ParseWebhookEvent(bodyBytes, sig, stripeClient.WebhookSecret)
		if errParse != nil {
			err = errParse
			log.Error().Err(err).Msg("[RegisterWebhookRoute] failed to parse webhook event")
			http.Error(w, "failed to parse event", http.StatusBadRequest)
			return
		}
		// Parsed event; further processing below.

		// 5) handle event type
		switch event.Type {
		case "payment_intent.created":
			pi, parseErr := stripeClient.ExtractPaymentIntent(event)
			if parseErr != nil {
				err = parseErr
				log.Error().Err(err).Msg("[RegisterWebhookRoute] failed to parse PaymentIntent from payment_intent.created")
				break
			}
			// Payment intent created externally (not from gRPC) - authorize local payment

			if errAuth := app.AuthorizePayment(ctx, application.AuthorizePaymentCommand{
				PaymentID:      pi.ID,
				UserCustomerID: extractUserCustomerID(pi), // Helper to extract user ID from metadata
				Amount:         pi.Amount,
			}); errAuth != nil {
				err = errAuth
				log.Error().Err(err).
					Str("payment_intent_id", pi.ID).
					Msg("[RegisterWebhookRoute] failed to authorize payment from payment_intent.created")
			} else {
				log.Info().
					Str("payment_intent_id", pi.ID).
					Int64("amount", pi.Amount).
					Msg("[RegisterWebhookRoute] payment authorized from webhook payment_intent.created")
			}

		case "payment_intent.succeeded":
			pi, parseErr := stripeClient.ExtractPaymentIntent(event)
			if parseErr != nil {
				err = parseErr
				log.Error().Err(err).Msg("[RegisterWebhookRoute] failed to parse PaymentIntent")
				break
			}
			// Intent succeeded – confirm local payment.

			if errConf := app.ConfirmPayment(ctx, application.ConfirmPaymentCommand{
				PaymentID: pi.ID,
			}); errConf != nil {
				err = errConf
				log.Error().Err(err).
					Str("payment_intent_id", pi.ID).
					Msg("[RegisterWebhookRoute] failed to confirm payment in local DB")
			}
			// etc...

		case "charge.succeeded":
			pi, parseErr := stripeClient.ExtractPaymentIntent(event)
			if parseErr != nil {
				err = parseErr
				log.Error().Err(err).Msg("[RegisterWebhookRoute] failed to parse PaymentIntent from charge.succeeded")
				break
			}
			// Charge succeeded – confirm local payment using payment intent ID.

			if errConf := app.ConfirmPayment(ctx, application.ConfirmPaymentCommand{
				PaymentID: pi.ID,
			}); errConf != nil {
				err = errConf
				log.Error().Err(err).
					Str("payment_intent_id", pi.ID).
					Msg("[RegisterWebhookRoute] failed to confirm payment from charge.succeeded")
			}
		default:
			// Unhandled event types are ignored.
		}

		// 6) respond 200 so Stripe won't keep retrying
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

// extractUserCustomerID extracts the user customer ID from payment intent metadata
// This assumes the user ID is stored in metadata when creating the payment intent
func extractUserCustomerID(pi *stripelib.PaymentIntent) string {
	if pi.Metadata != nil {
		if userID, exists := pi.Metadata["user_customer_id"]; exists {
			return userID
		}
		// Fallback to customer_id if user_customer_id is not present
		if customerID, exists := pi.Metadata["customer_id"]; exists {
			return customerID
		}
	}
	// Return empty string if no user ID found in metadata
	return ""
}
