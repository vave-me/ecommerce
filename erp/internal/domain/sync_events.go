package domain

import "time"

// Sync Event Names
const (
	ProductsSyncCompletedEvent  = "erp.ProductsSyncCompleted"
	StockSyncCompletedEvent     = "erp.StockSyncCompleted"
	PricesSyncCompletedEvent    = "erp.PricesSyncCompleted"
	CustomersSyncCompletedEvent = "erp.CustomersSyncCompleted"
	OrderSentToERPEvent         = "erp.OrderSentToERP"
	OrderSyncedFromERPEvent     = "erp.OrderSyncedFromERP"
	WebhookReceivedEvent        = "erp.WebhookReceived"
	WebhookProcessedEvent       = "erp.WebhookProcessed"
	WebhookFailedEvent          = "erp.WebhookFailed"
)

// ProductsSyncCompleted event
type ProductsSyncCompleted struct {
	ERPType       string
	TotalProducts int32
	SuccessCount  int32
	FailedCount   int32
	CompletedAt   time.Time
}

func (ProductsSyncCompleted) Key() string { return ProductsSyncCompletedEvent }

// StockSyncCompleted event
type StockSyncCompleted struct {
	ERPType      string
	TotalItems   int32
	SuccessCount int32
	FailedCount  int32
	CompletedAt  time.Time
}

func (StockSyncCompleted) Key() string { return StockSyncCompletedEvent }

// PricesSyncCompleted event
type PricesSyncCompleted struct {
	ERPType      string
	TotalItems   int32
	SuccessCount int32
	FailedCount  int32
	CompletedAt  time.Time
}

func (PricesSyncCompleted) Key() string { return PricesSyncCompletedEvent }

// CustomersSyncCompleted event
type CustomersSyncCompleted struct {
	ERPType        string
	TotalCustomers int32
	SuccessCount   int32
	FailedCount    int32
	CompletedAt    time.Time
}

func (CustomersSyncCompleted) Key() string { return CustomersSyncCompletedEvent }

// OrderSentToERP event
type OrderSentToERP struct {
	OrderID     string
	ERPOrderID  string
	ERPType     string
	ConnectorID string
	SentAt      time.Time
}

func (OrderSentToERP) Key() string { return OrderSentToERPEvent }

// OrderSyncedFromERP event
type OrderSyncedFromERP struct {
	OrderID    string
	ERPOrderID string
	ERPType    string
	SyncedAt   time.Time
}

func (OrderSyncedFromERP) Key() string { return OrderSyncedFromERPEvent }

// WebhookReceived event
type WebhookReceived struct {
	WebhookID  string
	ERPType    string
	EventType  string
	ReceivedAt time.Time
}

func (WebhookReceived) Key() string { return WebhookReceivedEvent }
func (e WebhookReceived) GetEventID() string { return e.WebhookID }
func (e WebhookReceived) GetEventType() string { return WebhookReceivedEvent }
func (e WebhookReceived) GetTimestamp() time.Time { return e.ReceivedAt }

// WebhookProcessed event
type WebhookProcessed struct {
	WebhookID   string
	ProcessedAt time.Time
}

func (WebhookProcessed) Key() string { return WebhookProcessedEvent }
func (e WebhookProcessed) GetEventID() string { return e.WebhookID }
func (e WebhookProcessed) GetEventType() string { return WebhookProcessedEvent }
func (e WebhookProcessed) GetTimestamp() time.Time { return e.ProcessedAt }

// WebhookFailed event
type WebhookFailed struct {
	WebhookID string
	Error     string
	FailedAt  time.Time
}

func (WebhookFailed) Key() string { return WebhookFailedEvent }
func (e WebhookFailed) GetEventID() string { return e.WebhookID }
func (e WebhookFailed) GetEventType() string { return WebhookFailedEvent }
func (e WebhookFailed) GetTimestamp() time.Time { return e.FailedAt }