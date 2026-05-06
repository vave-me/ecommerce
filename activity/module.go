package activity

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	"middleman/activity/activitypb"
	"middleman/activity/internal/application"
	"middleman/activity/internal/constants"
	"middleman/activity/internal/domain"
	"middleman/activity/internal/grpc"
	"middleman/activity/internal/handlers"
	"middleman/activity/internal/postgres"
	"middleman/activity/internal/redis_repository"
	"middleman/activity/internal/rest"
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
	"middleman/posts/postspb"
	"middleman/products/productspb"
	"middleman/users/userspb"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono system.ActivityService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.ActivityService) (err error) {
	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := domain.Registrations(reg); err != nil {
			return nil, err
		}
		if err := activitypb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := userspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := productspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := postspb.Registrations(reg); err != nil {
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
	container.AddScoped(constants.RedisPoolKey, func(c di.Container) (any, error) { return svc.RedisPoolActivity(), nil })

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

	container.AddScoped(constants.CommandPublisherKey, func(c di.Container) (any, error) {
		return am.NewCommandPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
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
	container.AddScoped(constants.ActivityRepoKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return es.NewAggregateRepository[*domain.Activity](
			domain.ActivityAggregate,
			reg,
			es.AggregateStoreWithMiddleware(
				pg.NewEventStore(constants.EventsTableName, tx, reg),
				pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg),
			),
		), nil
	})

	container.AddScoped(constants.InteractionsRepoKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return es.NewAggregateRepository[*domain.Interaction](
			domain.InteractionAggregate,
			reg,
			es.AggregateStoreWithMiddleware(
				pg.NewEventStore(constants.EventsTableName, tx, reg),
				pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg),
			),
		), nil
	})
	container.AddScoped(constants.UsersRepoKey, func(c di.Container) (any, error) {
		return postgres.NewUserCacheRepository(
			constants.UsersCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
			grpc.NewUserRepository(svc.Config().Rpc.Service(constants.UsersServiceName)),
		), nil
	})
	container.AddScoped(constants.ProductsRepoKey, func(c di.Container) (any, error) {
		return postgres.NewProductCacheRepository(
			constants.ProductsCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
			grpc.NewProductRepository(svc.Config().Rpc.Service(constants.UsersServiceName)),
		), nil
	})

	container.AddScoped(constants.MiddlemanCacheInteractionRepoKey, func(c di.Container) (any, error) {
		return redis_repository.NewMiddlemanCacheInteractionRepository(
			constants.MiddlemanInteractionTableName,
			c.Get(constants.MiddlemanInteractionRepoKey).(domain.MiddlemanInteractionRepository),
		), nil
	})

	container.AddScoped(constants.MiddlemanCacheRepoKey, func(c di.Container) (any, error) {
		return redis_repository.NewMiddlemanCacheRepository(
			constants.MiddlemanCacheTableName,
			c.Get(constants.MiddlemanRepoKey).(domain.MiddlemanRepository),
		), nil
	})

	container.AddScoped(constants.MiddlemanRepoKey, func(c di.Container) (any, error) {
		return postgres.NewMiddlemanRepository(
			constants.MiddlemanTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.MiddlemanInteractionRepoKey, func(c di.Container) (any, error) {
		return postgres.NewMiddlemanInteractionRepository(
			constants.MiddlemanInteractionTableName,
			constants.MiddlemanInteractionCountsTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	// TODO implement prometheus counters

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.ActivityRepoKey).(domain.ActivityRepository),
			c.Get(constants.InteractionsRepoKey).(domain.InteractionRepository),
			c.Get(constants.MiddlemanRepoKey).(domain.MiddlemanRepository),
			c.Get(constants.MiddlemanInteractionRepoKey).(domain.MiddlemanInteractionRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event])), nil
	})
	container.AddScoped(constants.MiddlemanHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewMiddlemanHandlers(c.Get(constants.MiddlemanCacheRepoKey).(domain.MiddlemanCacheRepository)), nil
	})
	container.AddScoped(constants.MiddlemanInteractionHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewMiddlemanInteractionHandlers(
			c.Get(constants.MiddlemanCacheInteractionRepoKey).(domain.MiddlemanCacheInteractionRepository)), nil
	})
	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})

	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewIntegrationEventHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.UsersRepoKey).(domain.UserCacheRepository),
			c.Get(constants.ProductsRepoKey).(domain.ProductCacheRepository),
			c.Get(constants.ApplicationKey).(application.App),
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
	if err = rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
		return err
	}
	if err = rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}
	handlers.RegisterDomainEventHandlersTx(container)
	handlers.RegisterMiddlemanHandlersTx(container)
	handlers.RegisterMiddlemanInteractionHandlersTx(container)
	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())
	return
}
func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Store
	if err = serde.Register(domain.Activity{}, func(v any) error {
		activity := v.(*domain.Activity)
		activity.Aggregate = es.NewAggregate("", domain.ActivityAggregate)
		return nil
	}); err != nil {
		return
	}

	if err = serde.Register(domain.Interaction{}, func(v any) error {
		interaction := v.(*domain.Interaction)
		interaction.Aggregate = es.NewAggregate("", domain.InteractionAggregate)
		return nil
	}); err != nil {
		return
	}
	// store events
	if err = serde.Register(domain.ActivityCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.InteractionAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.InteractionUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.InteractionRemoved{}); err != nil {
		return
	}

	// store snapshots
	if err = serde.RegisterKey(domain.ActivityVi{}.SnapshotName(), domain.ActivityVi{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.InteractionVi{}.SnapshotName(), domain.InteractionVi{}); err != nil {
		return
	}
	return
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("activity outbox processor encountered an error")
		}
	}()
}
