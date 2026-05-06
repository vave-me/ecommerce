package commands

import (
	"context"
	"fmt"
	"middleman/internal/erp"
)

// RegisterConnector command registers a new ERP connector
type RegisterConnector struct {
	ConnectorID string
	Type        erp.ERPType
	Config      erp.ERPConfig
}

// RegisterConnectorHandler handles the RegisterConnector command
type RegisterConnectorHandler struct {
	factory  erp.ConnectorFactory
	registry erp.ConnectorRegistry
}

// NewRegisterConnectorHandler creates a new handler
func NewRegisterConnectorHandler(factory erp.ConnectorFactory, registry erp.ConnectorRegistry) RegisterConnectorHandler {
	return RegisterConnectorHandler{
		factory:  factory,
		registry: registry,
	}
}

// RegisterConnector registers a new connector
func (h RegisterConnectorHandler) RegisterConnector(ctx context.Context, cmd RegisterConnector) error {
	// Validate configuration
	if cmd.ConnectorID == "" {
		return fmt.Errorf("connector ID is required")
	}
	if cmd.Type == "" {
		return fmt.Errorf("connector type is required")
	}

	// Check if connector already exists
	if _, err := h.registry.GetConnector(cmd.ConnectorID); err == nil {
		return fmt.Errorf("connector with ID %s already exists", cmd.ConnectorID)
	}

	// Create connector instance
	connector, err := h.factory.CreateConnector(cmd.Config)
	if err != nil {
		return fmt.Errorf("creating connector: %w", err)
	}

	// Health check
	health := connector.HealthCheck(ctx)
	if health.Status != erp.HealthStatusHealthy {
		return fmt.Errorf("connector is not healthy: %s", health.Message)
	}

	// Register connector
	if err := h.registry.RegisterConnector(cmd.ConnectorID, connector); err != nil {
		return fmt.Errorf("registering connector: %w", err)
	}

	return nil
}
