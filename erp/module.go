package erp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"middleman/erp/internal/application"
	"middleman/erp/internal/connectors"
	"middleman/erp/internal/constants"
	"middleman/erp/internal/crypto"
	"middleman/erp/internal/domain"
	"middleman/erp/internal/grpc"
	"middleman/erp/internal/handlers"
	"middleman/erp/internal/postgres"
	"middleman/erp/internal/rest"
	"middleman/internal/am"
	"middleman/internal/amotel"
	"middleman/internal/amprom"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/erp"
	"middleman/internal/es"
	"middleman/internal/jetstream"
	pg "middleman/internal/postgres"
	"middleman/internal/postgresotel"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
	"middleman/internal/system"
	"middleman/internal/tm"

	"github.com/rs/zerolog"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono system.Service) error {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) error {
	container := di.New()

	// Setup Registry
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})

	// Setup Stream
	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())

	container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
	})

	// Setup Database Transaction
	container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DB().Begin()
	})

	// Setup Message Publisher
	sentCounter := amprom.SentMessagesCounter(constants.ServiceName)
	container.AddScoped(constants.MessagePublisherKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		outboxStore := pg.NewOutboxStore(constants.OutboxTableName, tx)
		return am.NewMessagePublisher(
			stream,
			svc.Logger(),
			amotel.OtelMessageContextInjector(),
			sentCounter,
			tm.OutboxPublisher(outboxStore),
		), nil
	})

	// Setup Message Subscriber
	container.AddSingleton(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream,
			svc.Logger(),
			amotel.OtelMessageContextExtractor(),
			amprom.ReceivedMessagesCounter(constants.ServiceName),
		), nil
	})

	// Setup Event Publisher
	container.AddScoped(constants.EventPublisherKey, func(c di.Container) (any, error) {
		return am.NewEventPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
			svc.Logger(),
		), nil
	})

	// Setup Reply Publisher
	container.AddScoped(constants.ReplyPublisherKey, func(c di.Container) (any, error) {
		return am.NewReplyPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
			svc.Logger(),
		), nil
	})

	// Setup Inbox Store
	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		return pg.NewInboxStore(constants.InboxTableName, tx), nil
	})

	// Setup Domain Repositories
	container.AddScoped(constants.WebhookEventRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewWebhookEventRepository(
			constants.WebhookEventsTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.SyncStatusRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewSyncStatusRepository(
			constants.SyncStatusTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.SyncLogRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewSyncLogRepository(
			constants.SyncLogsTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.SyncConfigurationRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewSyncConfigurationRepository(
			constants.SyncConfigurationsTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.OrderSyncRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewOrderSyncRepository(
			constants.OrderSyncTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.InvoiceSyncRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewInvoiceSyncRepository(
			constants.InvoiceSyncTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// Setup Connector Repository
	container.AddScoped(constants.ConnectorRepositoryKey, func(c di.Container) (any, error) {
		return postgres.NewConnectorRepository(svc.DB()), nil
	})

	// Setup ERP Connector Factory and Registry
	container.AddSingleton(constants.ConnectorFactoryKey, func(c di.Container) (any, error) {
		factory := erp.NewConnectorFactory()
		// Register all connector builders
		if err := connectors.RegisterConnectorBuilders(factory); err != nil {
			return nil, err
		}
		return factory, nil
	})

	container.AddSingleton(constants.ConnectorRegistryKey, func(c di.Container) (any, error) {
		return erp.NewConnectorRegistry(), nil
	})

	// Setup Aggregate Store
	container.AddScoped(constants.AggregateStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return es.AggregateStoreWithMiddleware(
			pg.NewEventStore(constants.EventsTableName, tx, reg),
			pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg),
		), nil
	})

	// Setup Invoice Repository
	container.AddScoped(constants.InvoiceRepositoryKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Invoice](
			domain.InvoiceAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	// Setup Return Repository
	container.AddScoped(constants.ReturnRepositoryKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Return](
			domain.ReturnAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	// Setup InventoryReservation Repository
	container.AddScoped(constants.InventoryReservationRepositoryKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.InventoryReservation](
			domain.InventoryReservationAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	// Setup Product Repository
	container.AddScoped(constants.ProductRepositoryKey, func(c di.Container) (any, error) {
		return grpc.NewProductRepository(svc.Config().Rpc.Service(constants.ProductsServiceName)), nil
	})

	// Setup Encryptor
	container.AddSingleton(constants.EncryptorKey, func(c di.Container) (any, error) {
		// TODO: Get master key from environment or configuration
		masterKey := "default-encryption-key-change-me"
		return crypto.NewEncryptor(masterKey), nil
	})

	// Setup Application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.InvoiceRepositoryKey).(es.AggregateRepository[*domain.Invoice]),
			c.Get(constants.ReturnRepositoryKey).(es.AggregateRepository[*domain.Return]),
			c.Get(constants.InventoryReservationRepositoryKey).(es.AggregateRepository[*domain.InventoryReservation]),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
			c.Get(constants.ConnectorFactoryKey).(erp.ConnectorFactory),
			c.Get(constants.ConnectorRegistryKey).(erp.ConnectorRegistry),
			c.Get(constants.ConnectorRepositoryKey).(*postgres.ConnectorRepository),
			c.Get(constants.SyncLogRepositoryKey).(domain.SyncLogRepository),
			c.Get(constants.WebhookEventRepositoryKey).(domain.WebhookEventRepository),
			c.Get(constants.OrderSyncRepositoryKey).(domain.OrderSyncRepository),
			c.Get(constants.InvoiceSyncRepositoryKey).(domain.InvoiceSyncRepository),
			c.Get(constants.ProductRepositoryKey).(domain.ProductRepository),
			c.Get(constants.EncryptorKey).(*crypto.Encryptor),
		), nil
	})

	// Setup Command Handlers
	container.AddScoped(constants.CommandHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCommandHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ApplicationKey).(application.App),
			c.Get(constants.ReplyPublisherKey).(am.ReplyPublisher),
			tm.InboxHandler(c.Get(constants.InboxStoreKey).(tm.InboxStore)),
		), nil
	})

	// Setup Domain Event Handlers
	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})

	// Setup Outbox Processor
	outboxProcessor := tm.NewOutboxProcessor(
		stream,
		pg.NewOutboxStore(constants.OutboxTableName, svc.DB()),
	)

	// Setup Integration Event Handlers
	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewIntegrationEventHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ApplicationKey).(application.App),
			tm.InboxHandler(c.Get(constants.InboxStoreKey).(tm.InboxStore)),
		), nil
	})

	// Register Command Handlers
	if err := handlers.RegisterCommandHandlersTx(container); err != nil {
		return err
	}

	// Register Integration Event Handlers
	if err := handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}

	// Register Domain Event Handlers
	handlers.RegisterDomainEventHandlersTx(container)

	// Start Outbox Processor
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	// Register gRPC server
	if err := grpc.RegisterServerTx(container, svc.RPC()); err != nil {
		return err
	}

	// Register REST gateway and swagger
	if err := rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}

	// Initialize connectors from database
	logger := svc.Logger()
	if err := initializeConnectors(ctx, container, logger); err != nil {
		logger.Error().Err(err).Msg("Failed to initialize connectors")
		// Don't return error - service can still run without pre-loaded connectors
	}

	return nil
}

func registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// Register Invoice aggregate
	if err := serde.Register(domain.Invoice{}, func(v any) error {
		invoice := v.(*domain.Invoice)
		invoice.Aggregate = es.NewAggregate("", domain.InvoiceAggregate)
		return nil
	}); err != nil {
		return err
	}

	// Register Return aggregate
	if err := serde.Register(domain.Return{}, func(v any) error {
		ret := v.(*domain.Return)
		ret.Aggregate = es.NewAggregate("", domain.ReturnAggregate)
		return nil
	}); err != nil {
		return err
	}

	// Register InventoryReservation aggregate
	if err := serde.Register(domain.InventoryReservation{}, func(v any) error {
		res := v.(*domain.InventoryReservation)
		res.Aggregate = es.NewAggregate("", domain.InventoryReservationAggregate)
		return nil
	}); err != nil {
		return err
	}

	// Register Invoice events
	if err := serde.Register(domain.InvoiceCreated{}); err != nil {
		return err
	}
	if err := serde.Register(domain.InvoiceApproved{}); err != nil {
		return err
	}
	if err := serde.Register(domain.InvoiceSent{}); err != nil {
		return err
	}
	if err := serde.Register(domain.InvoicePaymentReceived{}); err != nil {
		return err
	}
	if err := serde.Register(domain.InvoiceVoided{}); err != nil {
		return err
	}
	if err := serde.Register(domain.InvoiceLinkedToERP{}); err != nil {
		return err
	}
	if err := serde.RegisterKey(domain.InvoiceV1{}.SnapshotName(), domain.InvoiceV1{}); err != nil {
		return err
	}

	// Register Return events
	if err := serde.Register(domain.ReturnCreated{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ReturnApproved{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ReturnProcessed{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ReturnCompleted{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ReturnRejected{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ReturnItemsRestocked{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ReturnLinkedToERP{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ReturnSyncFailed{}); err != nil {
		return err
	}
	if err := serde.RegisterKey(domain.ReturnV1{}.SnapshotName(), domain.ReturnV1{}); err != nil {
		return err
	}

	// Register InventoryReservation events
	if err := serde.Register(domain.ReservationCreated{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ReservationReleased{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ReservationFulfilled{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ReservationTransferred{}); err != nil {
		return err
	}
	if err := serde.Register(domain.ReservationExpired{}); err != nil {
		return err
	}
	if err := serde.RegisterKey(domain.InventoryReservationV1{}.SnapshotName(), domain.InventoryReservationV1{}); err != nil {
		return err
	}

	return nil
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		if err := outboxProcessor.Start(ctx); err != nil {
			logger.Error().Err(err).Msg("ERP outbox processor encountered an error")
		}
	}()
}

// initializeConnectors loads active connectors from database and registers them
func initializeConnectors(ctx context.Context, container di.Container, logger zerolog.Logger) error {
	// Get repositories
	connectorRepo := container.Get(constants.ConnectorRepositoryKey).(*postgres.ConnectorRepository)
	factory := container.Get(constants.ConnectorFactoryKey).(erp.ConnectorFactory)
	registry := container.Get(constants.ConnectorRegistryKey).(erp.ConnectorRegistry)
	encryptor := container.Get(constants.EncryptorKey).(*crypto.Encryptor)

	// Load active connectors from database
	connectors, err := connectorRepo.GetByStatus(ctx, domain.ConnectorStatusActive)
	if err != nil {
		return fmt.Errorf("failed to load active connectors: %w", err)
	}

	logger.Info().Int("count", len(connectors)).Msg("Loading active connectors")

	// Initialize each connector
	for _, connectorEntity := range connectors {
		// Decrypt auth config
		authConfig, err := decryptAuthConfig(encryptor, connectorEntity.AuthConfigEncrypted, connectorEntity.AuthConfigSalt)
		if err != nil {
			logger.Error().
				Err(err).
				Str("connector_id", connectorEntity.ID).
				Str("name", connectorEntity.Name).
				Msg("Failed to decrypt connector auth config")
			continue
		}

		// Create ERP config
		config := erp.ERPConfig{
			ID:       connectorEntity.ID,
			Name:     connectorEntity.Name,
			Type:     erp.ERPType(connectorEntity.Type),
			Endpoint: connectorEntity.BaseURL,
			Auth:     authConfig,
			Webhook: erp.WebhookConfig{
				Enabled:      connectorEntity.WebhookEnabled,
				URL:          connectorEntity.WebhookURL,
				Secret:       "", // Will be decrypted when needed
				ValidateSign: true,
				Events:       connectorEntity.WebhookEvents,
			},
			Sync: erp.SyncConfig{
				Enabled:   connectorEntity.SyncEnabled,
				Interval:  time.Duration(connectorEntity.SyncIntervalSeconds) * time.Second,
				BatchSize: connectorEntity.BatchSize,
			},
			RateLimit: &erp.RateLimitConfig{
				RequestsPerSecond: connectorEntity.RateLimitRequestsPerSecond,
				BurstSize:         connectorEntity.RateLimitBurst,
			},
			Retry: &erp.RetryConfig{
				MaxAttempts:  connectorEntity.RetryMaxAttempts,
				InitialDelay: time.Duration(connectorEntity.RetryInitialDelayMs) * time.Millisecond,
				MaxDelay:     time.Duration(connectorEntity.RetryMaxDelayMs) * time.Millisecond,
				Multiplier:   connectorEntity.RetryMultiplier,
			},
			Metadata: convertMapStringToInterface(connectorEntity.CustomHeaders),
		}

		// Create connector instance
		connector, err := factory.CreateConnector(config)
		if err != nil {
			logger.Error().
				Err(err).
				Str("connector_id", connectorEntity.ID).
				Str("name", connectorEntity.Name).
				Str("type", connectorEntity.Type).
				Msg("Failed to create connector instance")
			continue
		}

		// Register connector
		if err := registry.RegisterConnector(connectorEntity.ID, connector); err != nil {
			logger.Error().
				Err(err).
				Str("connector_id", connectorEntity.ID).
				Str("name", connectorEntity.Name).
				Msg("Failed to register connector")
			continue
		}

		logger.Info().
			Str("connector_id", connectorEntity.ID).
			Str("name", connectorEntity.Name).
			Str("type", connectorEntity.Type).
			Msg("Successfully initialized connector")
	}

	return nil
}

// decryptAuthConfig decrypts the auth configuration
func decryptAuthConfig(encryptor *crypto.Encryptor, encrypted []byte, salt string) (erp.AuthConfig, error) {
	if len(encrypted) == 0 {
		return erp.AuthConfig{}, nil
	}

	decrypted, err := encryptor.DecryptJSON(encrypted, salt)
	if err != nil {
		return erp.AuthConfig{}, fmt.Errorf("failed to decrypt auth config: %w", err)
	}

	var authConfig erp.AuthConfig
	if err := json.Unmarshal(decrypted, &authConfig); err != nil {
		return erp.AuthConfig{}, fmt.Errorf("failed to unmarshal auth config: %w", err)
	}

	return authConfig, nil
}

// convertMapStringToInterface converts map[string]string to map[string]interface{}
func convertMapStringToInterface(m map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = v
	}
	return result
}

