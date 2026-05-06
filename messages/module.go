package message

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog"
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
	"middleman/internal/system"
	"middleman/internal/tm"
	"middleman/messages/internal/application"
	"middleman/messages/internal/constants"
	"middleman/messages/internal/domain"
	"middleman/messages/internal/grpc"
	"middleman/messages/internal/handlers"
	"middleman/messages/internal/postgres"
	"middleman/messages/internal/rest"
	"middleman/messages/messagespb"
)

type Module struct {
}

func (m *Module) Startup(ctx context.Context, mono system.MessengerService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.MessengerService) (err error) {
	container := di.New()

	//setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := domain.Registrations(reg); err != nil {
			return nil, err
		}
		if err := messagespb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})

	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())
	natsWs := jetstream.NewStream(svc.MessengerConfig().StreamName, svc.JS(), svc.Logger())

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

	container.AddScoped(constants.WebSocketPublisherKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		outboxStore := pg.NewOutboxStore(constants.OutboxTableName, tx)
		return am.NewMessagePublisher(
			natsWs,
			svc.Logger(),
			amotel.OtelMessageContextInjector(),
			sentCounter,
			tm.OutboxPublisher(outboxStore),
		), nil
	})
	container.AddSingleton(constants.WebSocketSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			natsWs, svc.Logger(),
			amotel.OtelMessageContextExtractor(),
			amprom.ReceivedMessagesCounter(constants.ServiceName),
		), nil
	})

	container.AddSingleton(constants.Logger, func(c di.Container) (any, error) {
		return svc.Logger(), nil
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

	container.AddScoped(constants.ConversationsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Conversation](
			domain.ConversationAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.MessagesRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Message](
			domain.MessageAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.MessengerRepoKey, func(c di.Container) (any, error) {
		return postgres.NewMessengerRepository(
			constants.MessengerTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.MiddlemanRepoKey, func(c di.Container) (any, error) {
		return postgres.NewMiddlemanRepository(
			constants.MiddlemanTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.ConversationsRepoKey).(domain.ConversationRepository),
			c.Get(constants.MessagesRepoKey).(domain.MessageRepository),
			c.Get(constants.MiddlemanRepoKey).(domain.MiddlemanRepository),
			c.Get(constants.MessengerRepoKey).(domain.MessengerRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})

	container.AddScoped(constants.MiddlemanHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewMiddlemanHandlers(c.Get(constants.MiddlemanRepoKey).(domain.MiddlemanRepository)), nil
	})
	container.AddScoped(constants.MessengerHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewMessengerHandlers(c.Get(constants.MessengerRepoKey).(domain.MessengerRepository)), nil
	})
	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(
			c.Get(constants.EventPublisherKey).(am.EventPublisher),
		), nil
	})

	container.AddScoped(constants.WebsocketsEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewWebsocketEventHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ApplicationKey).(application.App),
			svc.Logger(),
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
	handlers.RegisterMessengerHandlersTx(container)
	handlers.RegisterMiddlemanHandlersTx(container)
	handlers.RegisterDomainEventHandlersTx(container)
	if err = handlers.RegisterWebsocketsCommandsHandlersTx(container, svc.Logger()); err != nil {
		return err
	}
	//if err = messagespb.RegisterAsyncApi(svc.Mux()); err != nil {return err}

	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())
	return nil
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("messages outbox processor encountered an error")
		}
	}()
}
