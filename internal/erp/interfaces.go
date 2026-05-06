package erp

import (
	"context"
	"time"
)

// HealthStatus represents the health status of a connector
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
)

// ConnectorRegistry manages registered ERP connectors
type ConnectorRegistry interface {
	// RegisterConnector registers a connector instance
	RegisterConnector(id string, connector Connector) error

	// GetConnector retrieves a connector by ID
	GetConnector(id string) (Connector, error)

	// RemoveConnector removes a connector from the registry
	RemoveConnector(id string) error

	// ListConnectors returns all registered connector IDs
	ListConnectors() []string

	// HealthCheckAll performs health checks on all registered connectors
	HealthCheckAll(ctx context.Context) map[string]HealthCheck

	// GetAllConnectors returns all registered connectors
	GetAllConnectors() []Connector

	// GetConnectorsByType returns all connectors of a specific type
	GetConnectorsByType(erpType string) []Connector
}

// ConnectorFactory creates connector instances
type ConnectorFactory interface {
	// CreateConnector creates a new connector instance
	CreateConnector(config ERPConfig) (Connector, error)

	// RegisterBuilder registers a builder for a specific ERP type
	RegisterBuilder(erpType ERPType, builder ConnectorBuilder)

	// ListTypes returns all registered ERP types
	ListTypes() []ERPType
}

// ConnectorInfo contains information about a registered connector
type ConnectorInfo struct {
	ID        string
	Name      string
	Type      string
	Status    string
	Config    map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

// HealthCheck result from a connector
type HealthCheck struct {
	Status  HealthStatus
	Message string
}

// Connector interface defines operations for ERP connectors
type Connector interface {
	// GetID returns the connector ID
	GetID() string

	// GetType returns the type of the connector
	GetType() string

	// GetConfig returns the connector configuration
	GetConfig() ERPConfig

	// HealthCheck verifies the connector is operational
	HealthCheck(ctx context.Context) HealthCheck

	// SendOrder sends an order to the ERP system
	SendOrder(ctx context.Context, order *OrderPayload) error

	// GetOrder retrieves an order from the ERP system
	GetOrder(ctx context.Context, orderID string) (*OrderPayload, error)

	// SyncProducts synchronizes products from the ERP
	SyncProducts(ctx context.Context, since time.Time, batchSize int) ([]*ProductPayload, error)

	// SyncStock synchronizes stock levels from the ERP
	SyncStock(ctx context.Context, productIDs []string, batchSize int) ([]*StockPayload, error)

	// SyncPrices synchronizes prices from the ERP
	SyncPrices(ctx context.Context, productIDs []string, batchSize int) ([]*PricePayload, error)

	// SyncCustomers synchronizes customers from the ERP
	SyncCustomers(ctx context.Context, since time.Time, batchSize int) ([]*CustomerPayload, error)

	// ProcessWebhook processes an incoming webhook
	ProcessWebhook(ctx context.Context, payload []byte, signature string) error

	// ValidateWebhook validates a webhook signature
	ValidateWebhook(payload []byte, signature string) error

	// ParseWebhook parses webhook payload into canonical event
	ParseWebhook(payload []byte) (*CanonicalEvent, error)

	// CreateInvoice creates an invoice in the ERP system
	CreateInvoice(ctx context.Context, invoice *InvoicePayload) (string, error)

	// FetchCustomer fetches a single customer from the ERP
	FetchCustomer(ctx context.Context, customerID string) (*CustomerPayload, error)

	// FetchCustomers fetches multiple customers from the ERP
	FetchCustomers(ctx context.Context, customerIDs []string) ([]*CustomerPayload, error)

	// UpdateInvoice updates an invoice in the ERP system
	UpdateInvoice(ctx context.Context, invoiceID string, invoice *InvoicePayload) error

	// FetchPrices fetches prices from the ERP
	FetchPrices(ctx context.Context, productIDs []string, priceListID string) ([]*PricePayload, error)

	// FetchProducts fetches products from the ERP
	FetchProducts(ctx context.Context, productIDs []string) ([]*ProductPayload, error)

	// ProcessReturn processes a return in the ERP
	ProcessReturn(ctx context.Context, returnPayload *ReturnPayload) (*ReturnPayload, error)

	// UpdateInventory updates inventory in the ERP
	UpdateInventory(ctx context.Context, adjustments []*InventoryAdjustment) error

	// FetchStock fetches stock for specific products
	FetchStock(ctx context.Context, productIDs []string) ([]*StockPayload, error)

	// FetchAllStock fetches all stock information
	FetchAllStock(ctx context.Context, since time.Time) ([]*StockPayload, error)
}
