package queries

import (
	"context"
	"middleman/internal/erp"
	"time"
)

// ListConnectors query retrieves all registered connectors
type ListConnectors struct {
	Type   erp.ERPType // Optional filter by type
	Status string      // Optional filter by status (healthy, unhealthy, all)
}

// ConnectorListItem represents a connector in the list
type ConnectorListItem struct {
	ID          string
	Type        erp.ERPType
	Status      erp.HealthStatus
	Endpoint    string
	LastChecked *time.Time
}

// ListConnectorsHandler handles the ListConnectors query
type ListConnectorsHandler struct {
	registry erp.ConnectorRegistry
}

// NewListConnectorsHandler creates a new handler
func NewListConnectorsHandler(registry erp.ConnectorRegistry) ListConnectorsHandler {
	return ListConnectorsHandler{
		registry: registry,
	}
}

// ListConnectors retrieves the list of connectors
func (h ListConnectorsHandler) ListConnectors(ctx context.Context, query ListConnectors) ([]ConnectorListItem, error) {
	var connectors []erp.Connector

	if query.Type != "" {
		// Get connectors by type
		connectors = h.registry.GetConnectorsByType(string(query.Type))
	} else {
		// Get all connectors
		connectors = h.registry.GetAllConnectors()
	}

	// Build result list
	result := make([]ConnectorListItem, 0, len(connectors))
	for _, connector := range connectors {
		// Get health status
		health := connector.HealthCheck(ctx)
		status := health.Status

		// Apply status filter if specified
		if query.Status != "" && query.Status != "all" {
			if query.Status == "healthy" && status != erp.HealthStatusHealthy {
				continue
			}
			if query.Status == "unhealthy" && status == erp.HealthStatusHealthy {
				continue
			}
		}

		config := connector.GetConfig()
		item := ConnectorListItem{
			ID:          connector.GetID(),
			Type:        erp.ERPType(connector.GetType()),
			Status:      status,
			Endpoint:    config.Endpoint,
			LastChecked: ptrTime(time.Now()),
		}

		result = append(result, item)
	}

	return result, nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
