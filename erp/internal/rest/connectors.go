package rest

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"middleman/erp/internal/application"
	"middleman/erp/internal/application/commands"
	"middleman/erp/internal/application/queries"
	"middleman/erp/internal/constants"
	"middleman/internal/auth"
	"middleman/internal/di"
)

// RegisterConnectorRoutes registers additional REST routes for connector management
// These complement the auto-generated gRPC-gateway routes
func RegisterConnectorRoutes(container di.Container, mux *chi.Mux, app application.App) {
	// Health check endpoint for a specific connector
	mux.Get("/api/erp/connectors/{connectorId}/health", func(w http.ResponseWriter, r *http.Request) {
		ctx := container.Scoped(r.Context())
		var err error
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				log.Error().Interface("panic", p).Msg("panic in connector health check; tx rolled back")
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
				log.Error().Err(err).Msg("connector health check failed; tx rolled back")
			} else {
				if cerr := tx.Commit(); cerr != nil {
					log.Error().Err(cerr).Msg("commit failed in connector health check")
				}
			}
		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

		// Validate authentication
		claims, ok := auth.ClaimsFromContext(ctx)
		if !ok {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		connectorID := chi.URLParam(r, "connectorId")
		
		// Get connector status
		status, err := app.GetConnectorStatus(ctx, queries.GetConnectorStatus{
			ConnectorID: connectorID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Return health status
		health := map[string]interface{}{
			"connector_id": connectorID,
			"status":       status.Status,
			"healthy":      status.Status == "active",
			"checked_by":   claims.Subject,
			"checked_at":   time.Now(),
			"message":      status.Message,
		}

		if status.LastSync != nil {
			health["last_sync"] = status.LastSync
		}
		if status.WebhookInfo != nil {
			health["webhook_info"] = status.WebhookInfo
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
	})

	// Test connector configuration endpoint
	mux.Post("/api/erp/connectors/test", func(w http.ResponseWriter, r *http.Request) {
		ctx := container.Scoped(r.Context())
		var err error
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				log.Error().Interface("panic", p).Msg("panic in connector test; tx rolled back")
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
				log.Error().Err(err).Msg("connector test failed; tx rolled back")
			} else {
				if cerr := tx.Commit(); cerr != nil {
					log.Error().Err(cerr).Msg("commit failed in connector test")
				}
			}
		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

		// Validate authentication
		_, ok := auth.ClaimsFromContext(ctx)
		if !ok {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		// Parse request body
		var req struct {
			Type        string                 `json:"type"`
			BaseURL     string                 `json:"base_url"`
			AuthType    string                 `json:"auth_type"`
			AuthConfig  map[string]interface{} `json:"auth_config"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Test the configuration without saving
		// This would require adding a TestConnector command or using the factory directly
		result := map[string]interface{}{
			"valid":       false,
			"message":     "Test endpoint not fully implemented",
			"tested_at":   time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// Get connector sync entities configuration
	mux.Get("/api/erp/connectors/{connectorId}/sync-entities", func(w http.ResponseWriter, r *http.Request) {
		ctx := container.Scoped(r.Context())
		var err error
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				log.Error().Interface("panic", p).Msg("panic in get sync entities; tx rolled back")
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
				log.Error().Err(err).Msg("get sync entities failed; tx rolled back")
			} else {
				if cerr := tx.Commit(); cerr != nil {
					log.Error().Err(cerr).Msg("commit failed in get sync entities")
				}
			}
		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

		// Validate authentication
		_, ok := auth.ClaimsFromContext(ctx)
		if !ok {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		connectorID := chi.URLParam(r, "connectorId")
		
		// This would require adding a query to get sync entities
		// For now, return a placeholder
		entities := []map[string]interface{}{
			{
				"entity_type":    "product",
				"enabled":        true,
				"sync_direction": "bidirectional",
				"last_sync_at":   nil,
			},
			{
				"entity_type":    "stock",
				"enabled":        true,
				"sync_direction": "inbound",
				"last_sync_at":   nil,
			},
			{
				"entity_type":    "order",
				"enabled":        true,
				"sync_direction": "outbound",
				"last_sync_at":   nil,
			},
		}

		response := map[string]interface{}{
			"connector_id": connectorID,
			"entities":     entities,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Update sync entity configuration
	mux.Put("/api/erp/connectors/{connectorId}/sync-entities/{entityType}", func(w http.ResponseWriter, r *http.Request) {
		ctx := container.Scoped(r.Context())
		var err error
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				log.Error().Interface("panic", p).Msg("panic in update sync entity; tx rolled back")
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
				log.Error().Err(err).Msg("update sync entity failed; tx rolled back")
			} else {
				if cerr := tx.Commit(); cerr != nil {
					log.Error().Err(cerr).Msg("commit failed in update sync entity")
				}
			}
		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

		// Validate authentication
		_, ok := auth.ClaimsFromContext(ctx)
		if !ok {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		connectorID := chi.URLParam(r, "connectorId")
		entityType := chi.URLParam(r, "entityType")

		// Parse request body
		var req struct {
			Enabled       bool                   `json:"enabled"`
			SyncDirection string                 `json:"sync_direction"`
			Filters       map[string]interface{} `json:"filters"`
			FieldMapping  map[string]string      `json:"field_mapping"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// This would require adding an UpdateSyncEntity command
		// For now, return success
		response := map[string]interface{}{
			"connector_id":   connectorID,
			"entity_type":    entityType,
			"enabled":        req.Enabled,
			"sync_direction": req.SyncDirection,
			"updated_at":     time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Get connector audit logs
	mux.Get("/api/erp/connectors/{connectorId}/audit", func(w http.ResponseWriter, r *http.Request) {
		ctx := container.Scoped(r.Context())
		var err error
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				log.Error().Interface("panic", p).Msg("panic in get audit logs; tx rolled back")
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
				log.Error().Err(err).Msg("get audit logs failed; tx rolled back")
			} else {
				if cerr := tx.Commit(); cerr != nil {
					log.Error().Err(cerr).Msg("commit failed in get audit logs")
				}
			}
		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

		// Validate authentication
		_, ok := auth.ClaimsFromContext(ctx)
		if !ok {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		connectorID := chi.URLParam(r, "connectorId")
		
		// This would require adding a query to get audit logs
		// For now, return a placeholder
		logs := []map[string]interface{}{
			{
				"action":     "created",
				"changed_by": "user123",
				"changed_at": time.Now().Add(-24 * time.Hour),
				"details":    map[string]interface{}{"initial": "creation"},
			},
			{
				"action":     "updated",
				"changed_by": "user456",
				"changed_at": time.Now().Add(-2 * time.Hour),
				"details":    map[string]interface{}{"field": "webhook_enabled", "old": false, "new": true},
			},
		}

		response := map[string]interface{}{
			"connector_id": connectorID,
			"audit_logs":   logs,
			"total":        len(logs),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Trigger immediate sync for a connector
	mux.Post("/api/erp/connectors/{connectorId}/sync/{entityType}", func(w http.ResponseWriter, r *http.Request) {
		ctx := container.Scoped(r.Context())
		var err error
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				log.Error().Interface("panic", p).Msg("panic in trigger sync; tx rolled back")
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
				log.Error().Err(err).Msg("trigger sync failed; tx rolled back")
			} else {
				if cerr := tx.Commit(); cerr != nil {
					log.Error().Err(cerr).Msg("commit failed in trigger sync")
				}
			}
		}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

		// Validate authentication
		_, ok := auth.ClaimsFromContext(ctx)
		if !ok {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		connectorID := chi.URLParam(r, "connectorId")
		entityType := chi.URLParam(r, "entityType")

		// Trigger appropriate sync based on entity type
		switch entityType {
		case "products":
			err = app.SyncProducts(ctx, commands.SyncProducts{
				ConnectorID: connectorID,
			})
		case "stock":
			err = app.SyncStock(ctx, commands.SyncStock{
				ConnectorID: connectorID,
			})
		case "prices":
			err = app.SyncPrices(ctx, commands.SyncPrices{
				ConnectorID: connectorID,
			})
		case "customers":
			err = app.SyncCustomers(ctx, commands.SyncCustomers{
				ConnectorID: connectorID,
			})
		default:
			http.Error(w, "invalid entity type", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"connector_id": connectorID,
			"entity_type":  entityType,
			"sync_started": true,
			"started_at":   time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}