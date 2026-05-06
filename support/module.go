package support

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	"middleman/assistants/assistantspb"
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
	"middleman/ordering/orderingpb"
	"middleman/payments/paymentspb"
	"middleman/support/internal/application"
	"middleman/support/internal/constants"
	"middleman/support/internal/domain"
	"middleman/support/internal/grpc"
	"middleman/support/internal/handlers"
	"middleman/support/internal/postgres"
	"middleman/support/internal/rest"
	"middleman/support/supportpb"
	"middleman/users/userspb"
)

type Module struct {
}

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
		if err := supportpb.Registrations(reg); err != nil {
			return nil, err
		}
		// Register external events we subscribe to
		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := paymentspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := userspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := assistantspb.Registrations(reg); err != nil {
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

	// Event-sourced repositories
	container.AddScoped(constants.SupportChannelRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.SupportChannel](
			domain.SupportChannelAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.TicketRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Ticket](
			domain.TicketAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	// Catalog repositories (using non-transactional DB for queries)
	container.AddSingleton(constants.SupportChannelCatalogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewSupportChannelCatalogRepository(
			constants.SupportChannelsCatalogTableName,
			svc.DB(),
		), nil
	})

	container.AddSingleton(constants.TicketCatalogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewTicketCatalogRepository(
			constants.TicketsCatalogTableName,
			svc.DB(),
		), nil
	})

	container.AddSingleton(constants.KnowledgeArticleCatalogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewKnowledgeArticleCatalogRepository(
			constants.KnowledgeArticlesCatalogTableName,
			svc.DB(),
		), nil
	})

	container.AddSingleton(constants.CommunicationRepoKey, func(c di.Container) (any, error) {
		return postgres.NewCommunicationRepository(
			constants.CommunicationsTableName,
			svc.DB(),
		), nil
	})

	container.AddSingleton(constants.AIConfigurationRepoKey, func(c di.Container) (any, error) {
		return postgres.NewAIConfigurationRepository(
			constants.AIConfigurationsTableName,
			svc.DB(),
		), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.SupportChannelRepoKey).(domain.SupportChannelRepository),
			c.Get(constants.TicketRepoKey).(domain.TicketRepository),
			c.Get(constants.SupportChannelCatalogRepoKey).(domain.SupportChannelCatalogRepository),
			c.Get(constants.TicketCatalogRepoKey).(domain.TicketCatalogRepository),
			c.Get(constants.KnowledgeArticleCatalogRepoKey).(domain.KnowledgeArticleCatalogRepository),
			c.Get(constants.CommunicationRepoKey).(domain.CommunicationRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})
	
	container.AddScoped(constants.CatalogHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCatalogHandlers(
			c.Get(constants.SupportChannelCatalogRepoKey).(domain.SupportChannelCatalogRepository),
			c.Get(constants.TicketCatalogRepoKey).(domain.TicketCatalogRepository),
		), nil
	})
	
	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewIntegrationEventHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
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
	handlers.RegisterCatalogHandlersTx(container)
	handlers.RegisterDomainEventHandlersTx(container)
	
	// Register integration event handlers with transaction support
	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}

	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil
}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Support Channel aggregate
	if err = serde.Register(domain.SupportChannel{}, func(v any) error {
		channel := v.(*domain.SupportChannel)
		channel.Aggregate = es.NewAggregate("", domain.SupportChannelAggregate)
		return nil
	}); err != nil {
		return
	}

	// Ticket aggregate
	if err = serde.Register(domain.Ticket{}, func(v any) error {
		ticket := v.(*domain.Ticket)
		ticket.Aggregate = es.NewAggregate("", domain.TicketAggregate)
		return nil
	}); err != nil {
		return
	}

	// Support Channel events
	if err = serde.Register(domain.SupportChannelCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.SupportChannelSettingsUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.SupportChannelClosed{}); err != nil {
		return
	}
	if err = serde.Register(domain.SupportChannelReactivated{}); err != nil {
		return
	}

	// Ticket events
	if err = serde.Register(domain.TicketCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketAssigned{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketPriorityUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketEscalated{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketResolved{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketReopened{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketClosed{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketsMerged{}); err != nil {
		return
	}
	if err = serde.Register(domain.TicketsLinked{}); err != nil {
		return
	}

	// Communication events
	if err = serde.Register(domain.TicketReplyAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.InternalNoteAdded{}); err != nil {
		return
	}

	// Snapshots
	if err = serde.RegisterKey(domain.SupportChannelV1{}.SnapshotName(), domain.SupportChannelV1{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.TicketV1{}.SnapshotName(), domain.TicketV1{}); err != nil {
		return
	}

	return
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("support outbox processor encountered an error")
		}
	}()
}