package categories

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	"middleman/categories/categoriespb"
	"middleman/categories/internal/application"
	"middleman/categories/internal/constants"
	"middleman/categories/internal/domain"
	"middleman/categories/internal/grpc"
	"middleman/categories/internal/handlers"
	"middleman/categories/internal/postgres"
	"middleman/categories/internal/rest"
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
		if err := categoriespb.Registrations(reg); err != nil {
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
	container.AddScoped(constants.CategoriesRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Category](
			domain.CategoryAggregate,
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

	container.AddScoped(constants.FiltersRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Filter](
			domain.FilterAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.CatalogFilterRepoKey, func(c di.Container) (any, error) {
		return postgres.NewCatalogFilterRepository(
			constants.CatalogFilterTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.CategoriesRepoKey).(domain.CategoryRepository),
			c.Get(constants.CatalogRepoKey).(domain.CatalogRepository),
			c.Get(constants.CatalogRepoKey).(domain.CatalogCacheRepository),
			c.Get(constants.FiltersRepoKey).(domain.FilterRepository),
			c.Get(constants.CatalogFilterRepoKey).(domain.CatalogFilterRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})
	container.AddScoped(constants.CatalogHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCatalogHandlers(c.Get(constants.CatalogRepoKey).(domain.CatalogRepository)), nil
	})
	container.AddScoped(constants.CatalogFilterHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCatalogFilterHandlers(c.Get(constants.CatalogFilterRepoKey).(domain.CatalogFilterRepository)), nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
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
	//if err = grpc.RegisterServer(container.Get(constants.ApplicationKey).(application.App), svc.RPC()); err != nil {
	//	return err
	//}
	if err = rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}
	handlers.RegisterCatalogHandlersTx(container)
	handlers.RegisterCatalogFilterHandlersTx(container)
	handlers.RegisterDomainEventHandlersTx(container)

	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil
}
func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(domain.Category{}, func(v any) error {
		category := v.(*domain.Category)
		category.Aggregate = es.NewAggregate("", domain.CategoryAggregate)
		return nil
	}); err != nil {
		return
	}
	// category events
	if err = serde.Register(domain.CategoryAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.CategoryRebranded{}); err != nil {
		return
	}
	if err = serde.Register(domain.CategoryRemoved{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.CategoryV1{}.SnapshotName(), domain.CategoryV1{}); err != nil {
		return
	}

	if err = serde.Register(domain.Filter{}, func(v any) error {
		filter := v.(*domain.Filter)
		filter.Aggregate = es.NewAggregate("", domain.FilterAggregate)
		return nil
	}); err != nil {
		return
	}
	// category events
	if err = serde.Register(domain.FilterAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.FilterRebranded{}); err != nil {
		return
	}
	if err = serde.Register(domain.FilterRemoved{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.FilterV1{}.SnapshotName(), domain.FilterV1{}); err != nil {
		return
	}

	return
}
func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("categories outbox processor encountered an error")
		}
	}()
}
