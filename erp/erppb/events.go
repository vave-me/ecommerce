package erppb

import (
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
)

// Channels
const (
	ConnectorAggregateChannel = "middleman.erp.events.Connector"
	OrderAggregateChannel     = "middleman.erp.events.Order"
	InvoiceAggregateChannel   = "middleman.erp.events.Invoice"
	SyncAggregateChannel      = "middleman.erp.events.Sync"
	WebhookAggregateChannel   = "middleman.erp.events.Webhook"
)

// Connector Event Names
const (
	ConnectorRegisteredEvent = "erpapi.ConnectorRegistered"
	ConnectorUpdatedEvent    = "erpapi.ConnectorUpdated"
	ConnectorDisabledEvent   = "erpapi.ConnectorDisabled"
	ConnectorEnabledEvent    = "erpapi.ConnectorEnabled"
	ConnectorRemovedEvent    = "erpapi.ConnectorRemoved"
)

// Webhook Event Names
const (
	WebhookReceivedEvent  = "erpapi.WebhookReceived"
	WebhookProcessedEvent = "erpapi.WebhookProcessed"
	WebhookFailedEvent    = "erpapi.WebhookFailed"
)

// Order Event Names
const (
	OrderSentToERPEvent     = "erpapi.OrderSentToERP"
	OrderSyncedFromERPEvent = "erpapi.OrderSyncedFromERP"
	OrderSyncFailedEvent    = "erpapi.OrderSyncFailed"
)

// Invoice Event Names
const (
	InvoiceCreatedEvent         = "erpapi.InvoiceCreated"
	InvoiceUpdatedEvent         = "erpapi.InvoiceUpdated"
	InvoiceApprovedEvent        = "erpapi.InvoiceApproved"
	InvoiceVoidedEvent          = "erpapi.InvoiceVoided"
	InvoiceSentEvent            = "erpapi.InvoiceSent"
	InvoicePaymentReceivedEvent = "erpapi.InvoicePaymentReceived"
)

// Product Sync Event Names
const (
	ProductsSyncStartedEvent   = "erpapi.ProductsSyncStarted"
	ProductsSyncCompletedEvent = "erpapi.ProductsSyncCompleted"
	ProductSyncedEvent         = "erpapi.ProductSynced"
)

// Stock Sync Event Names
const (
	StockSyncStartedEvent   = "erpapi.StockSyncStarted"
	StockSyncCompletedEvent = "erpapi.StockSyncCompleted"
	StockUpdatedEvent       = "erpapi.StockUpdated"
)

// Price Sync Event Names
const (
	PricesSyncStartedEvent   = "erpapi.PricesSyncStarted"
	PricesSyncCompletedEvent = "erpapi.PricesSyncCompleted"
	PriceUpdatedEvent        = "erpapi.PriceUpdated"
)

// Customer Sync Event Names
const (
	CustomersSyncStartedEvent   = "erpapi.CustomersSyncStarted"
	CustomersSyncCompletedEvent = "erpapi.CustomersSyncCompleted"
	CustomerSyncedEvent         = "erpapi.CustomerSynced"
)

// Inventory Reservation Event Names
const (
	InventoryReservedEvent  = "erpapi.InventoryReserved"
	InventoryReleasedEvent  = "erpapi.InventoryReleased"
	InventoryConfirmedEvent = "erpapi.InventoryConfirmed"
)

// Return Event Names
const (
	ReturnProcessedEvent = "erpapi.ReturnProcessed"
	ReturnFailedEvent    = "erpapi.ReturnFailed"
)

// Sync Status Event Names
const (
	SyncStatusUpdatedEvent = "erpapi.SyncStatusUpdated"
	SyncErrorEvent         = "erpapi.SyncError"
)

// Commands & Command Channel
const (
	CommandChannel               = "middleman.erp.commands"
	SyncEntityCommand            = "erpapi.SyncEntity"
	ProcessERPEventCommand       = "erpapi.ProcessERPEvent"
	RetryFailedSyncCommand       = "erpapi.RetryFailedSync"
	UpdateConnectorConfigCommand = "erpapi.UpdateConnectorConfig"
)

// Registrations and Serde
func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Connector Events
	if err := serde.Register(&ConnectorRegistered{}); err != nil {
		return err
	}
	if err := serde.Register(&ConnectorUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&ConnectorDisabled{}); err != nil {
		return err
	}
	if err := serde.Register(&ConnectorEnabled{}); err != nil {
		return err
	}
	if err := serde.Register(&ConnectorRemoved{}); err != nil {
		return err
	}

	// Webhook Events
	if err := serde.Register(&WebhookReceived{}); err != nil {
		return err
	}
	if err := serde.Register(&WebhookProcessed{}); err != nil {
		return err
	}
	if err := serde.Register(&WebhookFailed{}); err != nil {
		return err
	}

	// Order Events
	if err := serde.Register(&OrderSentToERP{}); err != nil {
		return err
	}
	if err := serde.Register(&OrderSyncedFromERP{}); err != nil {
		return err
	}
	if err := serde.Register(&OrderSyncFailed{}); err != nil {
		return err
	}

	// Invoice Events
	if err := serde.Register(&InvoiceCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&InvoiceUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&InvoiceApproved{}); err != nil {
		return err
	}
	if err := serde.Register(&InvoiceVoided{}); err != nil {
		return err
	}
	if err := serde.Register(&InvoiceSent{}); err != nil {
		return err
	}
	if err := serde.Register(&InvoicePaymentReceived{}); err != nil {
		return err
	}

	// Product Sync Events
	if err := serde.Register(&ProductsSyncStarted{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductsSyncCompleted{}); err != nil {
		return err
	}
	if err := serde.Register(&ProductSynced{}); err != nil {
		return err
	}

	// Stock Sync Events
	if err := serde.Register(&StockSyncStarted{}); err != nil {
		return err
	}
	if err := serde.Register(&StockSyncCompleted{}); err != nil {
		return err
	}
	if err := serde.Register(&StockUpdated{}); err != nil {
		return err
	}

	// Price Sync Events
	if err := serde.Register(&PricesSyncStarted{}); err != nil {
		return err
	}
	if err := serde.Register(&PricesSyncCompleted{}); err != nil {
		return err
	}
	if err := serde.Register(&PriceUpdated{}); err != nil {
		return err
	}

	// Customer Sync Events
	if err := serde.Register(&CustomersSyncStarted{}); err != nil {
		return err
	}
	if err := serde.Register(&CustomersSyncCompleted{}); err != nil {
		return err
	}
	if err := serde.Register(&CustomerSynced{}); err != nil {
		return err
	}

	// Inventory Reservation Events
	if err := serde.Register(&InventoryReserved{}); err != nil {
		return err
	}
	if err := serde.Register(&InventoryReleased{}); err != nil {
		return err
	}
	if err := serde.Register(&InventoryConfirmed{}); err != nil {
		return err
	}

	// Return Events
	if err := serde.Register(&ReturnProcessed{}); err != nil {
		return err
	}
	if err := serde.Register(&ReturnFailed{}); err != nil {
		return err
	}

	// Sync Status Events
	if err := serde.Register(&SyncStatusUpdated{}); err != nil {
		return err
	}
	if err := serde.Register(&SyncError{}); err != nil {
		return err
	}

	// Commands
	if err := serde.Register(&SyncEntity{}); err != nil {
		return err
	}
	if err := serde.Register(&ProcessERPEvent{}); err != nil {
		return err
	}
	if err := serde.Register(&RetryFailedSync{}); err != nil {
		return err
	}
	if err := serde.Register(&UpdateConnectorConfig{}); err != nil {
		return err
	}

	return nil
}

// Connector Events
func (*ConnectorRegistered) Key() string { return ConnectorRegisteredEvent }
func (*ConnectorUpdated) Key() string    { return ConnectorUpdatedEvent }
func (*ConnectorDisabled) Key() string   { return ConnectorDisabledEvent }
func (*ConnectorEnabled) Key() string    { return ConnectorEnabledEvent }
func (*ConnectorRemoved) Key() string    { return ConnectorRemovedEvent }

// Webhook Events
func (*WebhookReceived) Key() string  { return WebhookReceivedEvent }
func (*WebhookProcessed) Key() string { return WebhookProcessedEvent }
func (*WebhookFailed) Key() string    { return WebhookFailedEvent }

// Order Events
func (*OrderSentToERP) Key() string     { return OrderSentToERPEvent }
func (*OrderSyncedFromERP) Key() string { return OrderSyncedFromERPEvent }
func (*OrderSyncFailed) Key() string    { return OrderSyncFailedEvent }

// Invoice Events
func (*InvoiceCreated) Key() string         { return InvoiceCreatedEvent }
func (*InvoiceUpdated) Key() string         { return InvoiceUpdatedEvent }
func (*InvoiceApproved) Key() string        { return InvoiceApprovedEvent }
func (*InvoiceVoided) Key() string          { return InvoiceVoidedEvent }
func (*InvoiceSent) Key() string            { return InvoiceSentEvent }
func (*InvoicePaymentReceived) Key() string { return InvoicePaymentReceivedEvent }

// Product Sync Events
func (*ProductsSyncStarted) Key() string   { return ProductsSyncStartedEvent }
func (*ProductsSyncCompleted) Key() string { return ProductsSyncCompletedEvent }
func (*ProductSynced) Key() string         { return ProductSyncedEvent }

// Stock Sync Events
func (*StockSyncStarted) Key() string   { return StockSyncStartedEvent }
func (*StockSyncCompleted) Key() string { return StockSyncCompletedEvent }
func (*StockUpdated) Key() string       { return StockUpdatedEvent }

// Price Sync Events
func (*PricesSyncStarted) Key() string   { return PricesSyncStartedEvent }
func (*PricesSyncCompleted) Key() string { return PricesSyncCompletedEvent }
func (*PriceUpdated) Key() string        { return PriceUpdatedEvent }

// Customer Sync Events
func (*CustomersSyncStarted) Key() string   { return CustomersSyncStartedEvent }
func (*CustomersSyncCompleted) Key() string { return CustomersSyncCompletedEvent }
func (*CustomerSynced) Key() string         { return CustomerSyncedEvent }

// Inventory Reservation Events
func (*InventoryReserved) Key() string  { return InventoryReservedEvent }
func (*InventoryReleased) Key() string  { return InventoryReleasedEvent }
func (*InventoryConfirmed) Key() string { return InventoryConfirmedEvent }

// Return Events
func (*ReturnProcessed) Key() string { return ReturnProcessedEvent }
func (*ReturnFailed) Key() string    { return ReturnFailedEvent }

// Sync Status Events
func (*SyncStatusUpdated) Key() string { return SyncStatusUpdatedEvent }
func (*SyncError) Key() string         { return SyncErrorEvent }

// Commands implement registry.Registerable so they can travel via NATS
func (*SyncEntity) Key() string            { return SyncEntityCommand }
func (*ProcessERPEvent) Key() string       { return ProcessERPEventCommand }
func (*RetryFailedSync) Key() string       { return RetryFailedSyncCommand }
func (*UpdateConnectorConfig) Key() string { return UpdateConnectorConfigCommand }
