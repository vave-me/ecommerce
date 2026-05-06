package domain

import (
	"context"
	"time"
)

// Event represents a domain event
type Event interface {
	GetEventID() string
	GetEventType() string
	GetTimestamp() time.Time
}

// BaseEvent contains common event fields
type BaseEvent struct {
	EventID   string
	EventType string
	Timestamp time.Time
}

func (e BaseEvent) GetEventID() string      { return e.EventID }
func (e BaseEvent) GetEventType() string    { return e.EventType }
func (e BaseEvent) GetTimestamp() time.Time { return e.Timestamp }

// ProductSyncedEvent is emitted when a product is synced from ERP
type ProductSyncedEvent struct {
	BaseEvent
	ConnectorID string
	Product     *ProductPayload
	WebhookID   string // Optional: if triggered by webhook
}

// StockUpdatedEvent is emitted when stock levels are updated
type StockUpdatedEvent struct {
	BaseEvent
	ConnectorID string
	SKU         string
	WarehouseID string
	Quantity    int
	Available   int
	Reserved    int
	WebhookID   string // Optional: if triggered by webhook
}

// PriceUpdatedEvent is emitted when prices are updated
type PriceUpdatedEvent struct {
	BaseEvent
	ConnectorID string
	SKU         string
	Price       float64
	Currency    string
	PriceList   string
	WebhookID   string // Optional: if triggered by webhook
}

// OrderSyncedEvent is emitted when an order is synced from ERP
type OrderSyncedEvent struct {
	BaseEvent
	ConnectorID string
	Order       *OrderPayload
	WebhookID   string // Optional: if triggered by webhook
}

// OrderSentEvent is emitted when an order is sent to ERP
type OrderSentEvent struct {
	BaseEvent
	ConnectorID string
	OrderID     string
	CustomerID  string
	TotalAmount float64
	Currency    string
	ItemCount   int
}

// CustomerSyncedEvent is emitted when a customer is synced from ERP
type CustomerSyncedEvent struct {
	BaseEvent
	ConnectorID    string
	Customer       *CustomerPayload
	WebhookID      string // Optional: if triggered by webhook
	IncludesCredit bool   // Whether credit info was synced
}

// InvoiceSyncedEvent is emitted when an invoice is synced from ERP
type InvoiceSyncedEvent struct {
	BaseEvent
	ConnectorID string
	Invoice     *InvoicePayload
	WebhookID   string // Optional: if triggered by webhook
}

// ConnectorRegisteredEvent is emitted when a new connector is registered
type ConnectorRegisteredEvent struct {
	BaseEvent
	ConnectorID string
	Type        ERPType
	Endpoint    string
}

// ConnectorHealthChangedEvent is emitted when connector health status changes
type ConnectorHealthChangedEvent struct {
	BaseEvent
	ConnectorID string
	Message     string
}

// InventoryReservationUpdatedEvent is emitted when inventory reservation changes
type InventoryReservationUpdatedEvent struct {
	BaseEvent
	ConnectorID   string
	OrderID       string
	SKU           string
	WarehouseID   string
	Quantity      int
	Action        string // reserve, release, fulfill, transfer
	ReservationID string
}

// InvoiceProcessedEvent is emitted when an invoice is processed
type InvoiceProcessedEvent struct {
	BaseEvent
	ConnectorID   string
	InvoiceID     string
	InvoiceNumber string
	CustomerID    string
	OrderID       string
	Action        string // create, update, approve, void, send, payment
	Status        string
	TotalAmount   float64
	Currency      string
	DueDate       time.Time
}

// PaymentReceivedEvent is emitted when a payment is received
type PaymentReceivedEvent struct {
	BaseEvent
	ConnectorID   string
	InvoiceID     string
	CustomerID    string
	Amount        float64
	Currency      string
	PaymentMethod string
}

// EventPublisher publishes domain events
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}
