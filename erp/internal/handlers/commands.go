package handlers

import (
	"context"
	"fmt"
	"time"

	"middleman/erp/erppb"
	"middleman/erp/internal/application"
	"middleman/erp/internal/application/commands"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/erp"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// commandHandlers routes incoming NATS/JetStream commands to the application.
type commandHandlers struct {
	app application.App
}

var _ ddd.CommandHandler[ddd.Command] = (*commandHandlers)(nil)

// NewCommandHandlers wires registry-based serialization, a reply publisher and middlewares.
func NewCommandHandlers(reg registry.Registry, app application.App, replyPublisher am.ReplyPublisher, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewCommandHandler(reg, replyPublisher, commandHandlers{app: app}, mws...)
}

// RegisterCommandHandlers subscribes the given subscriber to the ERP command channel.
func RegisterCommandHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(erppb.CommandChannel, handlers, am.MessageFilter{
		erppb.SyncEntityCommand,
		erppb.ProcessERPEventCommand,
		erppb.RetryFailedSyncCommand,
		erppb.UpdateConnectorConfigCommand,
	}, am.GroupName("erp-commands"))
	return err
}

// HandleCommand dispatches based on CommandName.
func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (reply ddd.Reply, err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent("error", trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		}
		span.AddEvent("done", trace.WithAttributes(attribute.Int64("ms", time.Since(started).Milliseconds())))
	}(time.Now())

	switch cmd.CommandName() {
	case erppb.SyncEntityCommand:
		return h.doSyncEntity(ctx, cmd)
	case erppb.ProcessERPEventCommand:
		return h.doProcessERPEvent(ctx, cmd)
	case erppb.RetryFailedSyncCommand:
		return h.doRetryFailedSync(ctx, cmd)
	case erppb.UpdateConnectorConfigCommand:
		return h.doUpdateConnectorConfig(ctx, cmd)
	}
	return nil, nil
}

func (h commandHandlers) doSyncEntity(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*erppb.SyncEntity)

	// Route to appropriate sync handler based on entity type
	switch payload.GetEntityType() {
	case "products":
		if err := h.app.SyncProducts(ctx, commands.SyncProducts{
			ConnectorID: payload.GetConnectorId(),
			Since:       time.Time{}, // Use zero time for now
			BatchSize:   100,         // Default batch size
		}); err != nil {
			return nil, err
		}
	case "stock":
		if err := h.app.SyncStock(ctx, commands.SyncStock{
			ConnectorID: payload.GetConnectorId(),
			ProductIDs:  []string{}, // Empty means all products
			Since:       time.Time{}, // Use zero time for now
			BatchSize:   100,         // Default batch size
		}); err != nil {
			return nil, err
		}
	case "prices":
		if err := h.app.SyncPrices(ctx, commands.SyncPrices{
			ConnectorID:  payload.GetConnectorId(),
			ProductIDs:   []string{},  // Empty means all products
			PriceListIDs: []string{},  // Empty means all price lists
			Since:        time.Time{}, // Use zero time for now
			BatchSize:    100,         // Default batch size
		}); err != nil {
			return nil, err
		}
	case "customers":
		if err := h.app.SyncCustomers(ctx, commands.SyncCustomers{
			ConnectorID:   payload.GetConnectorId(),
			CustomerIDs:   []string{},  // Empty means all customers
			Since:         time.Time{}, // Use zero time for now
			IncludeCredit: false,       // Default value
			BatchSize:     100,         // Default batch size
		}); err != nil {
			return nil, err
		}
	}

	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doProcessERPEvent(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*erppb.ProcessERPEvent)

	// Convert proto payload to JSON bytes for webhook processing
	payloadBytes, err := payload.GetPayload().MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshaling payload: %w", err)
	}

	// Process the webhook through the application
	if err := h.app.ProcessWebhook(ctx, commands.ProcessWebhook{
		ConnectorID: payload.GetConnectorId(),
		EventID:     fmt.Sprintf("erp_event_%d", time.Now().UnixNano()), // Generate event ID
		EventType:   payload.GetEventType(),
		Payload:     payloadBytes,
		Signature:   "", // No signature in ProcessERPEvent
		Headers:     make(map[string]string), // Empty headers
	}); err != nil {
		return nil, err
	}

	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doRetryFailedSync(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*erppb.RetryFailedSync)

	// Implementation would retry failed sync operations based on sync log ID
	// For now, we'll return success
	_ = payload
	return ddd.NewReply(am.SuccessReply, nil), nil
}

func (h commandHandlers) doUpdateConnectorConfig(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*erppb.UpdateConnectorConfig)

	// Convert proto Struct to map
	configMap := payload.GetConfig().AsMap()

	// Extract ERP type from config (required field)
	erpTypeStr, ok := configMap["type"].(string)
	if !ok {
		return nil, fmt.Errorf("ERP type not specified in config")
	}

	// Create ERPConfig from map
	config := erp.ERPConfig{
		Type:     erp.ERPType(erpTypeStr),
		Endpoint: getStringFromMap(configMap, "endpoint"),
		URL:      getStringFromMap(configMap, "url"),
	}

	// Parse auth config if present
	if authMap, ok := configMap["auth"].(map[string]interface{}); ok {
		config.Auth = erp.AuthConfig{
			Type:           getStringFromMap(authMap, "type"),
			ClientID:       getStringFromMap(authMap, "clientId"),
			ClientSecret:   getStringFromMap(authMap, "clientSecret"),
			APIKey:         getStringFromMap(authMap, "apiKey"),
			Username:       getStringFromMap(authMap, "username"),
			Password:       getStringFromMap(authMap, "password"),
			TokenURL:       getStringFromMap(authMap, "tokenUrl"),
			ConsumerKey:    getStringFromMap(authMap, "consumerKey"),
			ConsumerSecret: getStringFromMap(authMap, "consumerSecret"),
			TokenID:        getStringFromMap(authMap, "tokenId"),
			TokenSecret:    getStringFromMap(authMap, "tokenSecret"),
		}
	}

	// Parse webhook config if present
	if webhookMap, ok := configMap["webhook"].(map[string]interface{}); ok {
		config.Webhook = erp.WebhookConfig{
			Enabled: getBoolFromMap(webhookMap, "enabled"),
			URL:     getStringFromMap(webhookMap, "url"),
			Secret:  getStringFromMap(webhookMap, "secret"),
			Events:  getStringSliceFromMap(webhookMap, "events"),
		}
	}

	// Parse sync config if present
	if syncMap, ok := configMap["sync"].(map[string]interface{}); ok {
		config.Sync = erp.SyncConfig{
			Products:  erp.SyncEntityConfig{Enabled: getBoolFromMap(syncMap["products"], "enabled")},
			Stock:     erp.SyncEntityConfig{Enabled: getBoolFromMap(syncMap["stock"], "enabled")},
			Prices:    erp.SyncEntityConfig{Enabled: getBoolFromMap(syncMap["prices"], "enabled")},
			Orders:    erp.SyncEntityConfig{Enabled: getBoolFromMap(syncMap["orders"], "enabled")},
			Customers: erp.SyncEntityConfig{Enabled: getBoolFromMap(syncMap["customers"], "enabled")},
		}
	}

	// Store any extra metadata
	config.Metadata = configMap

	// Register or update connector
	if err := h.app.RegisterConnector(ctx, commands.RegisterConnector{
		ConnectorID: payload.GetConnectorId(),
		Type:        config.Type,
		Config:      config,
	}); err != nil {
		return nil, err
	}

	return ddd.NewReply(am.SuccessReply, nil), nil
}

// Helper functions
func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getBoolFromMap(m interface{}, key string) bool {
	if mapVal, ok := m.(map[string]interface{}); ok {
		if v, ok := mapVal[key].(bool); ok {
			return v
		}
	}
	return false
}

func getStringSliceFromMap(m map[string]interface{}, key string) []string {
	if v, ok := m[key].([]interface{}); ok {
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	return nil
}
