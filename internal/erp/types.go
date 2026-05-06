package erp

import (
	"time"
)

// ERPType represents the type of ERP system
type ERPType string

const (
	ERPTypeSAP         ERPType = "sap"
	ERPTypeOdoo        ERPType = "odoo"
	ERPTypeDynamics365 ERPType = "dynamics365"
	ERPTypeNetSuite    ERPType = "netsuite"
	ERPTypeERPNext     ERPType = "erpnext"
	ERPTypeOracle      ERPType = "oracle"
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
	EventTypeOrderShipped         EventType = "order.shipped"
	EventTypeCustomerCreated      EventType = "customer.created"
	EventTypeCustomerUpdated      EventType = "customer.updated"
	EventTypeInvoiceCreated       EventType = "invoice.created"
	EventTypeInvoiceUpdated       EventType = "invoice.updated"
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
	SKU        string                 `json:"sku"`
	LocationID string                 `json:"locationId"`
	Quantity   int                    `json:"quantity"`
	Available  int                    `json:"available"`
	Reserved   int                    `json:"reserved"`
	StockType  string                 `json:"stockType"`
	UpdatedAt  time.Time              `json:"updatedAt"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
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
	OrderID         string                 `json:"orderId"`
	CustomerID      string                 `json:"customerId"`
	Items           []OrderItem            `json:"items"`
	TotalAmount     float64                `json:"totalAmount"`
	Currency        string                 `json:"currency"`
	Status          string                 `json:"status"`
	CreatedAt       time.Time              `json:"createdAt"`
	ShippingAddress *Address               `json:"shippingAddress,omitempty"`
	Attributes      map[string]interface{} `json:"attributes,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	Name        string  `json:"name"`
	SKU         string  `json:"sku"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`       // Line total
	TotalAmount float64 `json:"amount"`
	Discount    float64 `json:"discount"`
	Description string  `json:"description"`
	TaxRate     float64 `json:"tax_rate"`
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

// Address represents a physical address
type Address struct {
	Street     string `json:"street"`
	Street2    string `json:"street2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
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
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      ERPType                `json:"type"`
	Endpoint  string                 `json:"endpoint"`
	Auth      AuthConfig             `json:"auth"`
	Webhook   WebhookConfig          `json:"webhook"`
	Sync      SyncConfig             `json:"sync"`
	RateLimit *RateLimitConfig       `json:"rateLimit,omitempty"`
	Retry     *RetryConfig           `json:"retry,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	URL       string                 `json:"url"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Type         string `json:"type"` // "oauth2", "oauth1", "api_key", "basic"
	ClientID     string `json:"clientId,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	APIKey       string `json:"apiKey,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	TokenURL     string `json:"tokenUrl,omitempty"`
	// OAuth 1.0a fields for NetSuite
	ConsumerKey    string `json:"consumerKey,omitempty"`
	ConsumerSecret string `json:"consumerSecret,omitempty"`
	TokenID        string `json:"tokenId,omitempty"`
	TokenSecret    string `json:"tokenSecret,omitempty"`
}

// WebhookConfig represents webhook configuration
type WebhookConfig struct {
	Enabled      bool     `json:"enabled"`
	Secret       string   `json:"secret,omitempty"`
	ValidateSign bool     `json:"validateSign"`
	URL          string   `json:"url"`
	Events       []string `json:"events,omitempty"` // List of event types to subscribe to
}

// SyncEntityConfig represents configuration for a specific entity sync
type SyncEntityConfig struct {
	Enabled   bool          `json:"enabled"`
	Interval  time.Duration `json:"interval,omitempty"`
	BatchSize int           `json:"batchSize,omitempty"`
}

// SyncConfig represents sync configuration
type SyncConfig struct {
	Enabled   bool              `json:"enabled"`
	Interval  time.Duration     `json:"interval"`
	BatchSize int               `json:"batchSize"`
	Products  SyncEntityConfig  `json:"products"`
	Stock     SyncEntityConfig  `json:"stock"`
	Prices    SyncEntityConfig  `json:"prices"`
	Orders    SyncEntityConfig  `json:"orders"`
	Customers SyncEntityConfig  `json:"customers"`
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

// NOTE: The ERPConnector interface has been moved to interfaces.go as Connector

// InventoryAdjustment represents an inventory adjustment
type InventoryAdjustment struct {
	ReferenceID         string         `json:"referenceId"`
	ReferenceType       string         `json:"referenceType"` // order, return, adjustment, transfer
	SKU                 string         `json:"sku"`
	WarehouseID         string         `json:"warehouseId"`
	LocationID          string         `json:"locationId"` // Alias for WarehouseID
	Type                AdjustmentType `json:"type"`
	Quantity            int            `json:"quantity"`      // Total quantity after adjustment
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
	ExternalID      string       `json:"externalId,omitempty"`
	OriginalOrderID string       `json:"originalOrderId"`
	OrderID         string       `json:"orderId"` // Alias for OriginalOrderID
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

// ReturnItem represents an item in a return
type ReturnItem struct {
	SKU           string  `json:"sku"`
	Quantity      int     `json:"quantity"`
	Reason        string  `json:"reason"`
	RefundAmount  float64 `json:"refundAmount"`
	RestockingFee float64 `json:"restockingFee,omitempty"`
	Notes         string  `json:"notes,omitempty"`
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

// InvoiceLine represents a line item in an invoice
type InvoiceLine struct {
	SKU         string  `json:"sku,omitempty"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	TaxRate     float64 `json:"taxRate"`
	TaxAmount   float64 `json:"taxAmount"`
	LineTotal   float64 `json:"lineTotal"`
}

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
