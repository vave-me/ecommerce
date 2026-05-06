package domain

import (
	"time"
)

// Alias types for clarity
type WebhookEventLog = WebhookEvent
type SyncConfigLog = SyncConfiguration

// ConnectorStatus represents the status of a connector
type ConnectorStatus string

const (
	ConnectorStatusPending  ConnectorStatus = "pending"
	ConnectorStatusActive   ConnectorStatus = "active"
	ConnectorStatusInactive ConnectorStatus = "inactive"
	ConnectorStatusError    ConnectorStatus = "error"
)

// ERPType represents the type of ERP system
type ERPType string

const (
	ERPTypeSAP         ERPType = "sap"
	ERPTypeOdoo        ERPType = "odoo"
	ERPTypeDynamics365 ERPType = "dynamics365"
	ERPTypeNetSuite    ERPType = "netsuite"
	ERPTypeERPNext     ERPType = "erpnext"
)

// EventType represents canonical event types
type EventType string

const (
	EventTypeProductMasterUpdated EventType = "product.master.updated"
	EventTypeStockLevelUpdated    EventType = "stock.level.updated"
	EventTypePriceUpdated         EventType = "price.updated"
	EventTypeProductCreated       EventType = "product.created"
	EventTypeProductDeleted       EventType = "product.deleted"
	EventTypeOrderCreated         EventType = "order.created"
	EventTypeOrderUpdated         EventType = "order.updated"
	EventTypeCustomerCreated      EventType = "customer.created"
	EventTypeCustomerUpdated      EventType = "customer.updated"
)

// CanonicalEvent represents a standardized event format
type CanonicalEvent struct {
	EventID        string      `json:"eventId"`
	EventType      EventType   `json:"eventType"`
	EventTimestamp time.Time   `json:"eventTimestamp"`
	Source         string      `json:"source"`
	CorrelationID  string      `json:"correlationId,omitempty"`
	Payload        interface{} `json:"payload"`
}

// ProductPayload represents product master data
type ProductPayload struct {
	SKU         string                 `json:"sku"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Brand       string                 `json:"brand,omitempty"`
	Weight      float64                `json:"weight,omitempty"`
	Dimensions  *Dimensions            `json:"dimensions,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// StockPayload represents stock level data
type StockPayload struct {
	SKU        string    `json:"sku"`
	LocationID string    `json:"locationId"`
	Quantity   int       `json:"quantity"`
	Available  int       `json:"available"`
	Reserved   int       `json:"reserved"`
	StockType  string    `json:"stockType"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// PricePayload represents pricing data
type PricePayload struct {
	SKU         string                 `json:"sku"`
	PriceListID string                 `json:"priceListId"`
	Currency    string                 `json:"currency"`
	Price       float64                `json:"price"`
	ValidFrom   time.Time              `json:"validFrom"`
	ValidTo     *time.Time             `json:"validTo,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// OrderPayload represents order data
type OrderPayload struct {
	OrderID     string                 `json:"orderId"`
	CustomerID  string                 `json:"customerId"`
	Items       []OrderItem            `json:"items"`
	TotalAmount float64                `json:"totalAmount"`
	Currency    string                 `json:"currency"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"createdAt"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	Name        string  `json:"name"`
	SKU         string  `json:"sku"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	TotalAmount float64 `json:"amount"`
	TaxRate     float64 `json:"tax_rate"`
}

// Address represents a customer address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// CustomerPayload represents customer data
type CustomerPayload struct {
	CustomerID string                 `json:"customerId"`
	Email      string                 `json:"email"`
	Name       string                 `json:"name"`
	Phone      string                 `json:"phone,omitempty"`
	Address    *Address               `json:"address,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// Dimensions represents physical dimensions
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"`
}

// WebhookPayload represents an incoming webhook
type WebhookPayload struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Source    string      `json:"source"`
	Data      interface{} `json:"data"`
	Signature string      `json:"signature,omitempty"`
}

// SyncResult represents the result of a sync operation
type SyncResult struct {
	EntityType   string    `json:"entityType"`
	TotalRecords int       `json:"totalRecords"`
	Successful   int       `json:"successful"`
	Failed       int       `json:"failed"`
	LastSyncTime time.Time `json:"lastSyncTime"`
	Errors       []string  `json:"errors,omitempty"`
}

// ERPConfig represents configuration for an ERP connector
type ERPConfig struct {
	Type      ERPType                `json:"type"`
	Endpoint  string                 `json:"endpoint"`
	Auth      AuthConfig             `json:"auth"`
	Webhook   WebhookConfig          `json:"webhook"`
	Sync      SyncConfig             `json:"sync"`
	RateLimit *RateLimitConfig       `json:"rateLimit,omitempty"`
	Retry     *RetryConfig           `json:"retry,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Type         string `json:"type"` // "oauth2", "api_key", "basic"
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	APIKey       string `json:"apiKey,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	TokenURL     string `json:"tokenUrl,omitempty"`
}

// WebhookConfig represents webhook configuration
type WebhookConfig struct {
	Enabled      bool   `json:"enabled"`
	URL          string `json:"url,omitempty"`
	Secret       string `json:"secret,omitempty"`
	ValidateSign bool   `json:"validateSign"`
}

// SyncConfig represents sync configuration
type SyncConfig struct {
	Enabled   bool          `json:"enabled"`
	Interval  time.Duration `json:"interval"`
	Entities  []string      `json:"entities"`
	BatchSize int           `json:"batchSize"`
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond int `json:"requestsPerSecond"`
	BurstSize         int `json:"burstSize"`
}

// RetryConfig represents retry configuration
type RetryConfig struct {
	MaxAttempts  int           `json:"maxAttempts"`
	InitialDelay time.Duration `json:"initialDelay"`
	MaxDelay     time.Duration `json:"maxDelay"`
	Multiplier   float64       `json:"multiplier"`
}

// InventoryAdjustment represents an inventory adjustment
type InventoryAdjustment struct {
	ReferenceID         string         `json:"referenceId"`
	ReferenceType       string         `json:"referenceType"` // order, return, adjustment, transfer
	SKU                 string         `json:"sku"`
	WarehouseID         string         `json:"warehouseId"`
	Type                AdjustmentType `json:"type"`
	QuantityDelta       int            `json:"quantityDelta"`
	ReservedDelta       int            `json:"reservedDelta"`
	Reason              string         `json:"reason"`
	TransferToWarehouse string         `json:"transferToWarehouse,omitempty"`
	Timestamp           time.Time      `json:"timestamp"`
}

// AdjustmentType represents the type of inventory adjustment
type AdjustmentType string

const (
	AdjustmentTypeReservation AdjustmentType = "reservation"
	AdjustmentTypeFulfillment AdjustmentType = "fulfillment"
	AdjustmentTypeReturn      AdjustmentType = "return"
	AdjustmentTypeTransfer    AdjustmentType = "transfer"
	AdjustmentTypeManual      AdjustmentType = "manual"
	AdjustmentTypeDamage      AdjustmentType = "damage"
)

// ReturnPayload represents a return/RMA
type ReturnPayload struct {
	ReturnID        string       `json:"returnId"`
	OriginalOrderID string       `json:"originalOrderId"`
	CustomerID      string       `json:"customerId"`
	Reason          string       `json:"reason"`
	Status          string       `json:"status"`
	CreatedAt       time.Time    `json:"createdAt"`
	Items           []ReturnItem `json:"items"`
	TotalRefund     float64      `json:"totalRefund"`
	RefundMethod    string       `json:"refundMethod"`
	WarehouseID     string       `json:"warehouseId"`
	Notes           string       `json:"notes,omitempty"`
}

// InvoicePayload represents an invoice
type InvoicePayload struct {
	InvoiceID     string        `json:"invoiceId"`
	InvoiceNumber string        `json:"invoiceNumber"`
	CustomerID    string        `json:"customerId"`
	OrderID       string        `json:"orderId,omitempty"`
	IssueDate     time.Time     `json:"issueDate"`
	DueDate       time.Time     `json:"dueDate"`
	Currency      string        `json:"currency"`
	Lines         []InvoiceLine `json:"lines"`
	SubTotal      float64       `json:"subTotal"`
	TaxAmount     float64       `json:"taxAmount"`
	TotalAmount   float64       `json:"totalAmount"`
	Status        string        `json:"status"`
	PaymentTerms  string        `json:"paymentTerms,omitempty"`
	Notes         string        `json:"notes,omitempty"`
}

// InvoiceLine is defined in invoice.go

// PaymentPayload represents a payment
type PaymentPayload struct {
	PaymentID     string    `json:"paymentId"`
	InvoiceID     string    `json:"invoiceId"`
	CustomerID    string    `json:"customerId"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"paymentMethod"`
	PaymentDate   time.Time `json:"paymentDate"`
	Reference     string    `json:"reference,omitempty"`
	Notes         string    `json:"notes,omitempty"`
}
