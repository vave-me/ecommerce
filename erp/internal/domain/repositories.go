package domain

import (
	"context"
	"time"
)

// SyncStatus represents the status of a sync operation
type SyncStatus string

const (
	SyncStatusPending    SyncStatus = "pending"
	SyncStatusInProgress SyncStatus = "in_progress"
	SyncStatusCompleted  SyncStatus = "completed"
	SyncStatusFailed     SyncStatus = "failed"
)

// SyncLog represents a sync operation log
type SyncLog struct {
	ID               string
	ConnectorID      string
	ERPType          string
	EntityType       string // product, stock, price, order, customer
	Status           SyncStatus
	StartedAt        time.Time
	CompletedAt      *time.Time
	Duration         time.Duration
	RecordsProcessed int
	RecordsTotal     int
	RecordsSuccess   int
	RecordsFailed    int
	LastSyncTime     *time.Time
	Error            string
	ErrorMessage     string
	Metadata         map[string]interface{}
}

// SyncLogRepository manages sync logs
type SyncLogRepository interface {
	Create(ctx context.Context, log *SyncLog) error
	Update(ctx context.Context, log *SyncLog) error
	GetByID(ctx context.Context, id string) (*SyncLog, error)
	GetByConnectorID(ctx context.Context, connectorID string) ([]*SyncLog, error)
	GetByStatus(ctx context.Context, status SyncStatus) ([]*SyncLog, error)
}

// WebhookStatus represents the status of a webhook event
type WebhookStatus string

const (
	WebhookStatusPending   WebhookStatus = "pending"
	WebhookStatusProcessed WebhookStatus = "processed"
	WebhookStatusFailed    WebhookStatus = "failed"
	WebhookStatusIgnored   WebhookStatus = "ignored"
)

// WebhookEvent represents an incoming webhook event
type WebhookEvent struct {
	ID           string
	ConnectorID  string
	EventID      string
	EventType    string
	ReceivedAt   time.Time
	ProcessedAt  *time.Time
	Status       WebhookStatus
	Payload      string
	Signature    string
	Headers      map[string]string
	Error        string
	ERPType      string // ERP system type
	Source       string // Source of the webhook
	ErrorMessage string // Detailed error message
	RetryCount   int    // Number of retry attempts
}

// WebhookEventRepository manages webhook events
type WebhookEventRepository interface {
	Create(ctx context.Context, event *WebhookEvent) error
	Update(ctx context.Context, event *WebhookEvent) error
	GetByID(ctx context.Context, id string) (*WebhookEvent, error)
	GetByConnectorID(ctx context.Context, connectorID string) ([]*WebhookEvent, error)
	GetByStatus(ctx context.Context, status WebhookStatus) ([]*WebhookEvent, error)
}

// OrderSyncDirection represents the direction of order sync
type OrderSyncDirection string

const (
	DirectionInbound  OrderSyncDirection = "inbound"
	DirectionOutbound OrderSyncDirection = "outbound"
)

// OrderSyncStatus represents the status of an order sync
type OrderSyncStatus string

const (
	OrderSyncStatusPending   OrderSyncStatus = "pending"
	OrderSyncStatusCompleted OrderSyncStatus = "completed"
	OrderSyncStatusFailed    OrderSyncStatus = "failed"
)

// OrderSync represents an order sync operation
type OrderSync struct {
	ID          string
	ConnectorID string
	OrderID     string
	Direction   OrderSyncDirection
	Status      OrderSyncStatus
	AttemptedAt time.Time
	CompletedAt *time.Time
	Error       string
	Payload     interface{}
}

// OrderSyncRepository manages order sync operations
type OrderSyncRepository interface {
	Create(ctx context.Context, sync *OrderSync) error
	Update(ctx context.Context, sync *OrderSync) error
	GetByID(ctx context.Context, id string) (*OrderSync, error)
	GetByOrderID(ctx context.Context, orderID string) ([]*OrderSync, error)
	GetByConnectorID(ctx context.Context, connectorID string) ([]*OrderSync, error)
}

// InvoiceSyncStatus represents the status of an invoice sync
type InvoiceSyncStatus string

const (
	InvoiceSyncStatusPending   InvoiceSyncStatus = "pending"
	InvoiceSyncStatusCompleted InvoiceSyncStatus = "completed"
	InvoiceSyncStatusFailed    InvoiceSyncStatus = "failed"
)

// InvoiceSync represents an invoice sync operation
type InvoiceSync struct {
	ID          string
	ConnectorID string
	InvoiceID   string
	ExternalID  string
	Action      string
	Status      InvoiceSyncStatus
	AttemptedAt time.Time
	CompletedAt *time.Time
	Error       string
	Payload     interface{}
}

// InvoiceSyncRepository manages invoice sync operations
type InvoiceSyncRepository interface {
	Create(ctx context.Context, sync *InvoiceSync) error
	Update(ctx context.Context, sync *InvoiceSync) error
	GetByID(ctx context.Context, id string) (*InvoiceSync, error)
	GetByInvoiceID(ctx context.Context, invoiceID string) ([]*InvoiceSync, error)
	GetByConnectorID(ctx context.Context, connectorID string) ([]*InvoiceSync, error)
}

// ProductRepository defines the interface for interacting with the products service
type ProductRepository interface {
	// AddProduct creates a new product
	AddProduct(ctx context.Context, product Product) error
	
	// UpdateProduct updates an existing product
	UpdateProduct(ctx context.Context, productID string, product Product) error
	
	// UpdateProductStock updates product stock levels
	UpdateProductStock(ctx context.Context, productID string, stock int64) error
	
	// UpdateProductPrice updates product price
	UpdateProductPrice(ctx context.Context, productID string, price int64) error
	
	// GetProductBySKU retrieves a product by its SKU
	GetProductBySKU(ctx context.Context, sku string) (*Product, error)
	
	// GetProductByID retrieves a product by its ID
	GetProductByID(ctx context.Context, productID string) (*Product, error)
}