package tickets

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
	"middleman/tickets/internal/application"
	"middleman/tickets/internal/constants"
	"middleman/tickets/internal/domain"
	"middleman/tickets/internal/grpc"
	"middleman/tickets/internal/handlers"
	"middleman/tickets/internal/postgres"
	"middleman/tickets/internal/rest"
	"middleman/tickets/ticketspb"

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
		if err := ticketspb.Registrations(reg); err != nil {
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

	container.AddScoped(constants.TicketsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Ticket](
			domain.TicketAggregate,
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

	container.AddScoped(constants.VariantsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Variant](
			domain.VariantAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.CatalogVariantRepoKey, func(c di.Container) (any, error) {
		return postgres.NewCatalogVariantRepository(
			constants.CatalogVariantTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.TicketsRepoKey).(domain.TicketRepository),
			c.Get(constants.CatalogRepoKey).(domain.CatalogRepository),
			c.Get(constants.CatalogRepoKey).(domain.CatalogCacheRepository),
			c.Get(constants.VariantsRepoKey).(domain.VariantRepository),
			c.Get(constants.CatalogVariantRepoKey).(domain.CatalogVariantRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})
	container.AddScoped(constants.CatalogHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCatalogHandlers(c.Get(constants.CatalogRepoKey).(domain.CatalogRepository)), nil
	})
	container.AddScoped(constants.CatalogVariantHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCatalogVariantHandlers(c.Get(constants.CatalogVariantRepoKey).(domain.CatalogVariantRepository)), nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
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
	handlers.RegisterCatalogVariantHandlersTx(container)
	handlers.RegisterDomainEventHandlersTx(container)
	if err = handlers.RegisterCommandHandlersTx(container); err != nil {
		return err
	}

	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil
}
func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(domain.Ticket{}, func(v any) error {
		ticket := v.(*domain.Ticket)
		ticket.Aggregate = es.NewAggregate("", domain.TicketAggregate)
		return nil
	}); err != nil {
		return
	}
	// ticket events
	if err = serde.Register(domain.TicketAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketRebranded{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketPriceIncreased{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketPriceDecreased{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketRemoved{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.TicketV1{}.SnapshotName(), domain.TicketV1{}); err != nil {
		return
	}

	if err = serde.Register(domain.Variant{}, func(v any) error {
		variant := v.(*domain.Variant)
		variant.Aggregate = es.NewAggregate("", domain.VariantAggregate)
		return nil
	}); err != nil {
		return
	}
	// ticket events
	if err = serde.Register(domain.VariantAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.VariantRebranded{}); err != nil {
		return
	}
	if err = serde.Register(domain.VariantPriceIncreased{}); err != nil {
		return
	}
	if err = serde.Register(domain.VariantPriceDecreased{}); err != nil {
		return
	}
	if err = serde.Register(domain.VariantRemoved{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.VariantV1{}.SnapshotName(), domain.VariantV1{}); err != nil {
		return
	}

	return
}
func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("tickets outbox processor encountered an error")
		}
	}()
}
