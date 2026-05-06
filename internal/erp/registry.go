package erp

import (
	"context"
	"fmt"
	"sync"
)

// connectorRegistry manages active ERP connector instances
type connectorRegistry struct {
	connectors map[string]Connector
	mu         sync.RWMutex
}

// NewConnectorRegistry creates a new connector registry
func NewConnectorRegistry() ConnectorRegistry {
	return &connectorRegistry{
		connectors: make(map[string]Connector),
	}
}

// RegisterConnector registers a connector instance
func (r *connectorRegistry) RegisterConnector(id string, connector Connector) error {
	if id == "" {
		return fmt.Errorf("connector ID cannot be empty")
	}
	if connector == nil {
		return fmt.Errorf("connector cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.connectors[id]; exists {
		return fmt.Errorf("connector with ID %s already registered", id)
	}

	r.connectors[id] = connector
	return nil
}

// GetConnector retrieves a connector by ID
func (r *connectorRegistry) GetConnector(id string) (Connector, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	connector, exists := r.connectors[id]
	if !exists {
		return nil, fmt.Errorf("connector with ID %s not found", id)
	}

	return connector, nil
}

// RemoveConnector removes a connector from the registry
func (r *connectorRegistry) RemoveConnector(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.connectors[id]; !exists {
		return fmt.Errorf("connector with ID %s not found", id)
	}

	delete(r.connectors, id)
	return nil
}

// ListConnectors returns all registered connector IDs
func (r *connectorRegistry) ListConnectors() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.connectors))
	for id := range r.connectors {
		ids = append(ids, id)
	}
	return ids
}

// HealthCheckAll performs health checks on all registered connectors
func (r *connectorRegistry) HealthCheckAll(ctx context.Context) map[string]HealthCheck {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make(map[string]HealthCheck)
	for id, connector := range r.connectors {
		results[id] = connector.HealthCheck(ctx)
	}
	return results
}

// GetAllConnectors returns all registered connectors
func (r *connectorRegistry) GetAllConnectors() []Connector {
	r.mu.RLock()
	defer r.mu.RUnlock()

	connectors := make([]Connector, 0, len(r.connectors))
	for _, connector := range r.connectors {
		connectors = append(connectors, connector)
	}
	return connectors
}

// GetConnectorsByType returns all connectors of a specific type
func (r *connectorRegistry) GetConnectorsByType(erpType string) []Connector {
	r.mu.RLock()
	defer r.mu.RUnlock()

	connectors := make([]Connector, 0)
	for _, connector := range r.connectors {
		if connector.GetType() == erpType {
			connectors = append(connectors, connector)
		}
	}
	return connectors
}
