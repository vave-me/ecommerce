package sap

import (
	"context"
	"database/sql"
	"os"
	
	"github.com/rs/zerolog"
	
	"middleman/internal/am"
	"middleman/internal/amotel"
	"middleman/internal/amprom"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/jetstream"
	pg "middleman/internal/postgres"
	"middleman/internal/postgresotel"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
	"middleman/internal/system"
	"middleman/internal/tm"
	
	"middleman/sap/internal/application"
	"middleman/sap/internal/constants"
	"middleman/sap/internal/domain"
	"middleman/sap/internal/grpc"
	"middleman/sap/internal/handlers"
	"middleman/sap/internal/postgres"
	"middleman/sap/internal/rest"
	"middleman/sap/internal/sap"
	"middleman/sap/sappb"
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
		if err := sappb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})
	
	// Setup SAP Client
	container.AddSingleton(constants.SAPClientKey, func(c di.Container) (any, error) {
		config := &sap.Config{
			BaseURL:       getEnvOrDefault("SAP_BASE_URL", "https://api.sap.example.com"),
			APIKey:        os.Getenv("SAP_API_KEY"),
			WebhookSecret: os.Getenv("SAP_WEBHOOK_SECRET"),
			ClientID:      os.Getenv("SAP_CLIENT_ID"),
			ClientSecret:  os.Getenv("SAP_CLIENT_SECRET"),
			TokenURL:      os.Getenv("SAP_TOKEN_URL"),
			
			// HANA Configuration
			HANAHost:     os.Getenv("SAP_HANA_HOST"),
			HANAPort:     getEnvOrDefault("SAP_HANA_PORT", "443"),
			HANAUser:     os.Getenv("SAP_HANA_USER"),
			HANAPassword: os.Getenv("SAP_HANA_PASSWORD"),
			HANAUseTLS:   getEnvBool("SAP_HANA_USE_TLS", true),
			
			// Security Configuration
			IASInstanceName: os.Getenv("SAP_IAS_INSTANCE_NAME"),
			Issuer:          os.Getenv("SAP_ISSUER"),
			Audience:        []string{getEnvOrDefault("SAP_AUDIENCE", "sap-connector")},
			
			// Features
			UseDirectHANA:  getEnvBool("SAP_USE_DIRECT_HANA", false),
			EnableSecurity: getEnvBool("SAP_ENABLE_SECURITY", false),
		}
		
		return sap.NewEnhancedSAPClient(config)
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
	container.AddScoped(constants.SyncStatusRepoKey, func(c di.Container) (any, error) {
		return postgres.NewSyncStatusRepository(
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	
	container.AddScoped(constants.SyncLogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewSyncLogRepository(
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	
	container.AddScoped(constants.SyncConfigRepoKey, func(c di.Container) (any, error) {
		return postgres.NewSyncConfigurationRepository(
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	
	container.AddScoped(constants.WebhookEventRepoKey, func(c di.Container) (any, error) {
		return postgres.NewWebhookEventRepository(
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	
	// Setup Application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.NewApplication(
			c.Get(constants.SAPClientKey).(*sap.EnhancedSAPClient),
			c.Get(constants.EventPublisherKey).(ddd.EventPublisher[ddd.Event]),
			c.Get(constants.SyncStatusRepoKey).(domain.SyncStatusRepository),
			c.Get(constants.SyncLogRepoKey).(domain.SyncLogRepository),
			c.Get(constants.SyncConfigRepoKey).(domain.SyncConfigurationRepository),
			c.Get(constants.WebhookEventRepoKey).(domain.WebhookEventRepository),
		), nil
	})
	
	// Setup Handlers
	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})
	
	container.AddScoped(constants.CommandHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCommandHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ApplicationKey).(application.SAPConnectorDomain),
			c.Get(constants.ReplyPublisherKey).(am.ReplyPublisher),
			tm.InboxHandler(c.Get(constants.InboxStoreKey).(tm.InboxStore)),
		), nil
	})
	
	// Setup Outbox Processor
	outboxProcessor := tm.NewOutboxProcessor(
		stream,
		pg.NewOutboxStore(constants.OutboxTableName, svc.DB()),
	)
	
	// Register gRPC Server
	if err := grpc.RegisterServerTx(container, svc.RPC()); err != nil {
		return err
	}
	
	// Register REST Gateway
	if err := rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
		return err
	}
	
	// Register Webhook Route
	rest.RegisterWebhookRoute(
		container,
		svc.Mux(),
		container.Get(constants.ApplicationKey).(application.SAPConnectorDomain),
		container.Get(constants.SAPClientKey).(*sap.EnhancedSAPClient),
	)
	
	// Register Swagger
	if err := rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}
	
	// Register Domain Event Handlers
	handlers.RegisterDomainEventHandlersTx(container)
	
	// Register Command Handlers
	if err := handlers.RegisterCommandHandlersTx(container); err != nil {
		return err
	}
	
	// Start Outbox Processor
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())
	
	// Start Scheduled Sync (if enabled)
	if getEnvBool("SAP_ENABLE_SCHEDULED_SYNC", false) {
		startScheduledSync(ctx, container, svc.Logger())
	}
	
	return nil
}

func registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)
	
	// Register domain events
	if err := serde.Register(domain.SyncStatus{}); err != nil {
		return err
	}
	if err := serde.Register(domain.SyncLog{}); err != nil {
		return err
	}
	if err := serde.Register(domain.SyncConfiguration{}); err != nil {
		return err
	}
	if err := serde.Register(domain.WebhookEvent{}); err != nil {
		return err
	}
	
	return nil
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		if err := outboxProcessor.Start(ctx); err != nil {
			logger.Error().Err(err).Msg("SAP connector outbox processor encountered an error")
		}
	}()
}

func startScheduledSync(ctx context.Context, container di.Container, logger zerolog.Logger) {
	go func() {
		// Implementation for scheduled sync
		// This would periodically check for sync configurations and execute them
		logger.Info().Msg("Scheduled sync started for SAP connector")
	}()
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}