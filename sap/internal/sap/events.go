package sap

import (
	"encoding/json"
	"time"
)

// EventType represents the type of SAP event
type EventType string

const (
	EventTypeProductCreated      EventType = "product.created"
	EventTypeProductUpdated      EventType = "product.updated"
	EventTypeProductDeleted      EventType = "product.deleted"
	EventTypeStockUpdated        EventType = "stock.updated"
	EventTypePriceUpdated        EventType = "price.updated"
	EventTypeOrderCreated        EventType = "order.created"
	EventTypeOrderStatusChanged  EventType = "order.status_changed"
)

// SAPEvent represents a webhook event from SAP
type SAPEvent struct {
	ID            string          `json:"id" xml:"id"`
	Type          EventType       `json:"type" xml:"type"`
	Timestamp     time.Time       `json:"timestamp" xml:"timestamp"`
	CorrelationID string          `json:"correlationId" xml:"correlationId"`
	Source        string          `json:"source" xml:"source"`
	Data          json.RawMessage `json:"data" xml:"data"`
}

// ProductChange represents a product change from SAP
type ProductChange struct {
	ProductID   string                 `json:"productId"`
	SKU         string                 `json:"sku"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Weight      float64                `json:"weight"`
	Dimensions  Dimensions             `json:"dimensions"`
	Attributes  map[string]interface{} `json:"attributes"`
	ChangedAt   time.Time              `json:"changedAt"`
	ChangeType  string                 `json:"changeType"`
}

// StockLevel represents stock information from SAP
type StockLevel struct {
	ProductID    string    `json:"productId"`
	SKU          string    `json:"sku"`
	WarehouseID  string    `json:"warehouseId"`
	Quantity     int       `json:"quantity"`
	ReservedQty  int       `json:"reservedQuantity"`
	AvailableQty int       `json:"availableQuantity"`
	StockType    string    `json:"stockType"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Price represents pricing information from SAP
type Price struct {
	ProductID   string     `json:"productId"`
	SKU         string     `json:"sku"`
	PriceListID string     `json:"priceListId"`
	Currency    string     `json:"currency"`
	Amount      float64    `json:"amount"`
	ValidFrom   time.Time  `json:"validFrom"`
	ValidTo     *time.Time `json:"validTo,omitempty"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// OrderData represents order information to send to SAP
type OrderData struct {
	OrderID     string      `json:"orderId"`
	CustomerID  string      `json:"customerId"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"totalAmount"`
	Currency    string      `json:"currency"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"createdAt"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string  `json:"productId"`
	SKU       string  `json:"sku"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Dimensions represents product dimensions
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"`
}

// ProductEventData represents product event data from SAP webhook
type ProductEventData struct {
	ProductID   string                 `json:"productId"`
	SKU         string                 `json:"sku"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Weight      float64                `json:"weight"`
	Dimensions  Dimensions             `json:"dimensions"`
	Attributes  map[string]interface{} `json:"attributes"`
}

// StockEventData represents stock event data from SAP webhook
type StockEventData struct {
	ProductID    string `json:"productId"`
	SKU          string `json:"sku"`
	WarehouseID  string `json:"warehouseId"`
	OldQuantity  int    `json:"oldQuantity"`
	NewQuantity  int    `json:"newQuantity"`
	Adjustment   int    `json:"adjustment"`
	Reason       string `json:"reason"`
	AdjustedBy   string `json:"adjustedBy"`
}

// PriceEventData represents price event data from SAP webhook
type PriceEventData struct {
	ProductID   string    `json:"productId"`
	SKU         string    `json:"sku"`
	PriceListID string    `json:"priceListId"`
	Currency    string    `json:"currency"`
	OldPrice    float64   `json:"oldPrice"`
	NewPrice    float64   `json:"newPrice"`
	ValidFrom   time.Time `json:"validFrom"`
	ValidTo     *time.Time `json:"validTo,omitempty"`
	ChangedBy   string    `json:"changedBy"`
}