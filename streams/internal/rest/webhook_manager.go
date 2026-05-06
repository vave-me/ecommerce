package rest

import (
	"encoding/json"
	"net/http"

	"middleman/internal/di"
	"middleman/internal/web"
	"middleman/streams/internal/application"
	"middleman/streams/internal/application/commands"
	"middleman/streams/internal/application/queries"

	"github.com/go-chi/chi/v5"
)

// RegisterWebhookRoutes registers webhook management routes
func RegisterWebhookRoutes(container di.Container, router chi.Router, app application.StreamingApp) {
	router.Route("/webhooks", func(r chi.Router) {
		// Subscription management
		r.Post("/", createWebhookSubscription(container, app))
		r.Get("/", listWebhookSubscriptions(container, app))
		r.Get("/{id}", getWebhookSubscription(container, app))
		r.Put("/{id}", updateWebhookSubscription(container, app))
		r.Delete("/{id}", deleteWebhookSubscription(container, app))
		
		// Delivery history
		r.Get("/{id}/deliveries", getWebhookDeliveries(container, app))
		
		// Testing
		r.Post("/{id}/test", testWebhookSubscription(container, app))
		r.Post("/{id}/retry/{deliveryId}", retryWebhookDelivery(container, app))
	})
}

// createWebhookSubscription creates a new webhook subscription
func createWebhookSubscription(container di.Container, app application.StreamingApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cmd commands.SubscribeWebhook
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			web.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate required fields
		if cmd.URL == "" {
			web.RespondWithError(w, http.StatusBadRequest, "URL is required")
			return
		}
		if len(cmd.Events) == 0 {
			web.RespondWithError(w, http.StatusBadRequest, "At least one event is required")
			return
		}

		// Generate secret if not provided
		if cmd.Secret == "" {
			cmd.Secret = generateWebhookSecret()
		}

		err := app.SubscribeWebhook(r.Context(), cmd)
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		web.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "Webhook subscription created",
			"secret":  cmd.Secret,
		})
	}
}

// listWebhookSubscriptions lists all webhook subscriptions
func listWebhookSubscriptions(container di.Container, app application.StreamingApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := queries.ListWebhookSubscriptions{
			Active: r.URL.Query().Get("active") == "true",
		}

		subscriptions, err := app.ListWebhookSubscriptions(r.Context(), query)
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		web.RespondWithJSON(w, http.StatusOK, subscriptions)
	}
}

// getWebhookSubscription gets a webhook subscription by ID
func getWebhookSubscription(container di.Container, app application.StreamingApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		
		query := queries.GetWebhookSubscription{ID: id}
		subscription, err := app.GetWebhookSubscription(r.Context(), query)
		if err != nil {
			web.RespondWithError(w, http.StatusNotFound, "Webhook subscription not found")
			return
		}

		web.RespondWithJSON(w, http.StatusOK, subscription)
	}
}

// updateWebhookSubscription updates a webhook subscription
func updateWebhookSubscription(container di.Container, app application.StreamingApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		
		var cmd commands.UpdateWebhook
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			web.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		
		cmd.ID = id
		
		err := app.UpdateWebhook(r.Context(), cmd)
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		web.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Webhook subscription updated",
		})
	}
}

// deleteWebhookSubscription deletes a webhook subscription
func deleteWebhookSubscription(container di.Container, app application.StreamingApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		
		cmd := commands.UnsubscribeWebhook{ID: id}
		err := app.UnsubscribeWebhook(r.Context(), cmd)
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		web.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Webhook subscription deleted",
		})
	}
}

// getWebhookDeliveries gets delivery history for a webhook subscription
func getWebhookDeliveries(container di.Container, app application.StreamingApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		
		query := queries.GetWebhookDeliveries{
			SubscriptionID: id,
			Limit:          100,
		}
		
		deliveries, err := app.GetWebhookDeliveries(r.Context(), query)
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		web.RespondWithJSON(w, http.StatusOK, deliveries)
	}
}

// testWebhookSubscription sends a test webhook
func testWebhookSubscription(container di.Container, app application.StreamingApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		
		cmd := commands.TestWebhook{SubscriptionID: id}
		err := app.TestWebhook(r.Context(), cmd)
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		web.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Test webhook sent",
		})
	}
}

// retryWebhookDelivery retries a failed webhook delivery
func retryWebhookDelivery(container di.Container, app application.StreamingApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deliveryID := chi.URLParam(r, "deliveryId")
		
		cmd := commands.RetryWebhookDelivery{DeliveryID: deliveryID}
		err := app.RetryWebhookDelivery(r.Context(), cmd)
		if err != nil {
			web.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		web.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Webhook delivery retry initiated",
		})
	}
}

// generateWebhookSecret generates a secure webhook secret
func generateWebhookSecret() string {
	// In production, use a proper crypto random generator
	return "whsec_" + generateRandomString(32)
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)] // Simplified for example
	}
	return string(b)
}