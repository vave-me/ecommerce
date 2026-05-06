package merchant

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	"middleman/internal/am"
	"middleman/internal/amotel"
	"middleman/internal/amprom"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/jetstream"
	"middleman/internal/merchant"
	pg "middleman/internal/postgres"
	"middleman/internal/postgresotel"
	"middleman/internal/registry"
	"middleman/internal/system"
	"middleman/internal/tm"
	"middleman/merchant/internal/adapters"
	"middleman/merchant/internal/application"
	"middleman/merchant/internal/constants"
	"middleman/merchant/internal/domain"
	"middleman/merchant/internal/grpc"
	"middleman/merchant/internal/handlers"
	"middleman/merchant/internal/rest"
	"middleman/products/productspb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.MerchantService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.MerchantService) (err error) {
	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := productspb.Registrations(reg); err != nil {
			return nil, err
		}

		return reg, nil
	})
	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())
	merchantClient, err := merchant.NewMerchantCenterClient(ctx, svc.MerchantConfig().MerchantId, svc.MerchantConfig().ServiceAccountJSONPath)
	if err != nil {

		return err
	}

	// Add repository
	container.AddScoped(constants.ProductsSyncStatusRepoKey, func(c di.Container) (any, error) {
		tx := c.Get(constants.DatabaseTransactionKey).(*sql.Tx)
		return adapters.NewProductSyncStatusRepository(tx), nil
	})

	// Wrap merchant client with retry logic
	container.AddSingleton(constants.MerchantClientKey, func(c di.Container) (any, error) {
		wrapper := adapters.NewMerchantClientWrapper(merchantClient, svc.Logger())
		return wrapper, nil
	})
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
	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.MerchantClientKey).(domain.MerchantClient),
			c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event]),
			c.Get(constants.ProductsSyncStatusRepoKey).(domain.ProductSyncStatusRepository),
		), nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
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
	ctx = container.Scoped(ctx)

	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	handlers.RegisterDomainEventHandlersTx(container)

	// Setup Driver adapters
	if err = grpc.RegisterServer(
		ctx,
		container.Get(constants.ApplicationKey).(application.App),
		container.Get(constants.ProductsSyncStatusRepoKey).(domain.ProductSyncStatusRepository),
		svc.RPC(),
	); err != nil {
		return err
	}
	if err = rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
		return err
	}
	if err = rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}

	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("merchant outbox processor encountered an error")
		}
	}()
}
