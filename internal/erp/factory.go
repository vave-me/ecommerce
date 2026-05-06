package erp

import (
	"fmt"
	"sync"
)

// ConnectorFactory creates ERP connector instances
type connectorFactory struct {
	builders map[ERPType]ConnectorBuilder
	mu       sync.RWMutex
}

// ConnectorBuilder is a function that creates a connector from config
type ConnectorBuilder func(config ERPConfig) (Connector, error)

// ConnectorBuilderFunc is an alias for ConnectorBuilder
type ConnectorBuilderFunc = ConnectorBuilder

// NewConnectorFactory creates a new connector factory
func NewConnectorFactory() ConnectorFactory {
	return &connectorFactory{
		builders: make(map[ERPType]ConnectorBuilder),
	}
}

// RegisterBuilder registers a builder for a specific ERP type
func (f *connectorFactory) RegisterBuilder(erpType ERPType, builder ConnectorBuilder) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.builders[erpType] = builder
}

// CreateConnector creates a connector instance based on configuration
func (f *connectorFactory) CreateConnector(config ERPConfig) (Connector, error) {
	f.mu.RLock()
	builder, exists := f.builders[config.Type]
	f.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no builder registered for ERP type: %s", config.Type)
	}

	return builder(config)
}

// ListTypes returns all registered ERP types
func (f *connectorFactory) ListTypes() []ERPType {
	f.mu.RLock()
	defer f.mu.RUnlock()

	types := make([]ERPType, 0, len(f.builders))
	for t := range f.builders {
		types = append(types, t)
	}
	return types
}
