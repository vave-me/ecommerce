package connectors

import (
	"fmt"
	"middleman/internal/erp"
	"middleman/internal/netsuite"
	"middleman/internal/odoo"
)

// RegisterConnectorBuilders registers all available connector builders
func RegisterConnectorBuilders(factory erp.ConnectorFactory) error {
	// Register Odoo connector
	factory.RegisterBuilder(erp.ERPTypeOdoo, erp.ConnectorBuilderFunc(func(config erp.ERPConfig) (erp.Connector, error) {
		return odoo.NewConnector(config)
	}))

	// Register NetSuite connector
	factory.RegisterBuilder(erp.ERPTypeNetSuite, erp.ConnectorBuilderFunc(func(config erp.ERPConfig) (erp.Connector, error) {
		return netsuite.NewConnector(config)
	}))

	// Register Dynamics 365 connector (placeholder)
	factory.RegisterBuilder(erp.ERPTypeDynamics365, erp.ConnectorBuilderFunc(func(config erp.ERPConfig) (erp.Connector, error) {
		return nil, fmt.Errorf("Dynamics 365 connector not implemented yet")
	}))

	// Register SAP connector (placeholder)
	factory.RegisterBuilder(erp.ERPTypeSAP, erp.ConnectorBuilderFunc(func(config erp.ERPConfig) (erp.Connector, error) {
		return nil, fmt.Errorf("SAP connector not implemented yet")
	}))

	// Register Oracle connector (placeholder)
	factory.RegisterBuilder(erp.ERPTypeOracle, erp.ConnectorBuilderFunc(func(config erp.ERPConfig) (erp.Connector, error) {
		return nil, fmt.Errorf("Oracle connector not implemented yet")
	}))

	return nil
}

