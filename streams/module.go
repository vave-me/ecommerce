package streams

import (
	"context"
	"database/sql"
	"middleman/internal/am"
	"middleman/internal/amotel"
	"middleman/internal/amprom"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/es"
	"middleman/internal/jetstream"
	pg "middleman/internal/postgres"
	"middleman/internal/postgresotel"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
	"middleman/internal/system"
	"middleman/internal/tm"
	"middleman/streams/internal/application"
	"middleman/streams/internal/constants"
	"middleman/streams/internal/domain"
	"middleman/streams/internal/grpc"
	"middleman/streams/internal/handlers"
	"middleman/streams/internal/infrastructure"
	"middleman/streams/internal/infrastructure/streaming"
	"middleman/streams/internal/postgres"
	"middleman/streams/internal/rest"
	"middleman/streams/streamspb"
	
	"go.uber.org/zap"

	"github.com/rs/zerolog"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {

	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := streamspb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})

	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())
	container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
	})

	container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DB().Begin()
	})

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
	container.AddSingleton(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream,
			svc.Logger(),
			amotel.OtelMessageContextExtractor(),
			amprom.ReceivedMessagesCounter(constants.ServiceName),
		), nil
	})
	container.AddScoped(constants.EventPublisherKey, func(c di.Container) (any, error) {
		return am.NewEventPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
			svc.Logger(),
		), nil
	})
	container.AddScoped(constants.ReplyPublisherKey, func(c di.Container) (any, error) {
		return am.NewReplyPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
			svc.Logger(),
		), nil
	})
	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		return pg.NewInboxStore(constants.InboxTableName, tx), nil
	})
	container.AddScoped(constants.AggregateStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return es.AggregateStoreWithMiddleware(
			pg.NewEventStore(constants.EventsTableName, tx, reg),
			pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg),
		), nil
	})

	container.AddScoped(constants.LiveStreamsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.LiveStream](
			domain.LiveStreamAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.CatalogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewCatalogRepository(
			constants.CatalogTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// Webhook repositories
	container.AddScoped(constants.WebhookSubscriptionRepoKey, func(c di.Container) (any, error) {
		return postgres.NewWebhookSubscriptionRepository(
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	
	container.AddScoped(constants.WebhookDeliveryRepoKey, func(c di.Container) (any, error) {
		return postgres.NewWebhookDeliveryRepository(
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// Infrastructure services
	container.AddSingleton(constants.WebhookClientKey, func(c di.Container) (any, error) {
		return infrastructure.NewWebhookClient(svc.Logger()), nil
	})
	
	container.AddSingleton(constants.StreamingServerKey, func(c di.Container) (any, error) {
		// Initialize with default config - should be from actual config
		return streaming.NewStreamingServer(&streaming.StreamingConfig{
			RTMPEnabled: true,
			HLSEnabled: true,
			DASHEnabled: true,
		}, svc.Logger()), nil
	})
	
	container.AddSingleton(constants.CDNManagerKey, func(c di.Container) (any, error) {
		return streaming.NewCDNManager(svc.Logger()), nil
	})
	
	container.AddSingleton(constants.DRMManagerKey, func(c di.Container) (any, error) {
		return streaming.NewDRMManager(), nil
	})
	
	container.AddSingleton(constants.WebRTCServerKey, func(c di.Container) (any, error) {
		// Create zap logger adapter from zerolog
		zapLogger, _ := zap.NewProduction()
		return streaming.NewWebRTCServer(&streaming.WebRTCConfig{
			STUNServers: []string{"stun:stun.l.google.com:19302"},
		}, zapLogger), nil
	})
	
	container.AddSingleton(constants.WebhookDispatcherKey, func(c di.Container) (any, error) {
		// Create zap logger adapter from zerolog
		zapLogger, _ := zap.NewProduction()
		// For webhook dispatcher, we'll create non-transactional repositories
		webhookSubRepo := postgres.NewWebhookSubscriptionRepository(svc.DB())
		webhookDelRepo := postgres.NewWebhookDeliveryRepository(svc.DB())
		return handlers.NewWebhookDispatcher(
			webhookSubRepo,
			webhookDelRepo,
			c.Get(constants.WebhookClientKey).(*infrastructure.WebhookClient),
			zapLogger,
			10, // max concurrent webhook deliveries
		), nil
	})
	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.LiveStreamsRepoKey).(domain.LiveStreamRepository),
			c.Get(constants.WebhookSubscriptionRepoKey).(domain.WebhookSubscriptionRepository),
			c.Get(constants.WebhookDeliveryRepoKey).(domain.WebhookDeliveryRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
			c.Get(constants.StreamingServerKey).(*streaming.StreamingServer),
			c.Get(constants.CDNManagerKey).(*streaming.CDNManager),
			c.Get(constants.DRMManagerKey).(*streaming.DRMManager),
			c.Get(constants.WebRTCServerKey).(*streaming.WebRTCServer),
			c.Get(constants.WebhookClientKey).(*infrastructure.WebhookClient),
			c.Get(constants.WebhookDispatcherKey).(*handlers.WebhookDispatcher),
		), nil
	})
	container.AddScoped(constants.CatalogHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCatalogHandlers(c.Get(constants.CatalogRepoKey).(domain.CatalogRepository)), nil
	})
	container.AddScoped(constants.CatalogVariantHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCatalogVariantHandlers(c.Get(constants.CatalogVariantRepoKey).(domain.CatalogVariantRepository)), nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewLiveStreamingDomainEventHandlers(
			c.Get(constants.EventPublisherKey).(am.EventPublisher),
			c.Get(constants.WebhookDispatcherKey).(*handlers.WebhookDispatcher),
		), nil
	})

	// Command handlers
	container.AddScoped(constants.CommandHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCommandHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ApplicationKey).(application.App),
			c.Get(constants.ReplyPublisherKey).(am.ReplyPublisher),
			tm.InboxHandler(c.Get(constants.InboxStoreKey).(tm.InboxStore)),
		), nil
	})

	outboxProcessor := tm.NewOutboxProcessor(
		stream,
		pg.NewOutboxStore(constants.OutboxTableName, svc.DB()),
	)

	// setup Driver adapters
	if err = grpc.RegisterServerTx(container, svc.RPC()); err != nil {
		return err
	}
	if err = rest.RegisterGatewayWithWebhooks(ctx, container, svc.Mux(), svc.Config().Rpc.Address(), 
		container.Get(constants.ApplicationKey).(application.StreamingApp)); err != nil {
		return err
	}
	//if err = grpc.RegisterServer(container.Get(constants.ApplicationKey).(application.App), svc.RPC()); err != nil {
	//	return err
	//}
	if err = rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}
	handlers.RegisterCatalogHandlersTx(container)
	handlers.RegisterCatalogVariantHandlersTx(container)
	handlers.RegisterLiveStreamingDomainEventHandlersTx(container)
	if err = handlers.RegisterCommandHandlersTx(container); err != nil {
		return err
	}

	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil
}
func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Register LiveStream aggregate
	if err = serde.Register(domain.LiveStream{}, func(v any) error {
		stream := v.(*domain.LiveStream)
		stream.Aggregate = es.NewAggregate("", domain.LiveStreamAggregate)
		return nil
	}); err != nil {
		return
	}
	
	// LiveStream events
	if err = serde.Register(domain.LiveStreamCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.StreamingConfigured{}); err != nil {
		return
	}
	if err = serde.Register(domain.LiveStreamStarted{}); err != nil {
		return
	}
	if err = serde.Register(domain.LiveStreamStopped{}); err != nil {
		return
	}
	if err = serde.Register(domain.ViewerJoined{}); err != nil {
		return
	}
	if err = serde.Register(domain.ViewerLeft{}); err != nil {
		return
	}
	if err = serde.Register(domain.StreamQualityChanged{}); err != nil {
		return
	}
	if err = serde.Register(domain.StreamHealthUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.CDNEndpointAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.DRMConfigured{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.LiveStreamV1{}.SnapshotName(), domain.LiveStreamV1{}); err != nil {
		return
	}

	return
}
func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("streams outbox processor encountered an error")
		}
	}()
}
