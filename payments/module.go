package payments

import (
	"context"
	"database/sql"
	"middleman/internal/am"
	"middleman/internal/amotel"
	"middleman/internal/amprom"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/jetstream"
	pg "middleman/internal/postgres"
	"middleman/internal/postgresotel"
	"middleman/internal/registry"
	"middleman/internal/stripe"
	"middleman/internal/system"
	"middleman/internal/tm"
	"middleman/ordering/orderingpb"
	"middleman/payments/internal/adapters"
	"middleman/payments/internal/application"
	"middleman/payments/internal/constants"
	"middleman/payments/internal/grpc"
	"middleman/payments/internal/handlers"
	"middleman/payments/internal/postgres"
	"middleman/payments/internal/rest"
	"middleman/payments/paymentspb"

	"github.com/rs/zerolog"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.PaymentService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.PaymentService) (err error) {
	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := paymentspb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})
	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())
	stripeClient := stripe.NewStripeClient(svc.PaymentConfig().StripeWebhookSecret, svc.PaymentConfig().StripeAPIKey)

	container.AddSingleton(constants.Stripe, func(c di.Container) (any, error) {
		return stripeClient, nil
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
	container.AddScoped(constants.InvoicesRepoKey, func(c di.Container) (any, error) {
		return postgres.NewInvoiceRepository(
			constants.InvoicesTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	container.AddScoped(constants.RecurringRepoKey, func(c di.Container) (any, error) {
		return postgres.NewRecurringRepository(
			constants.RecurringTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	container.AddScoped(constants.PaymentsRepoKey, func(c di.Container) (any, error) {
		return postgres.NewPaymentRepository(
			constants.PaymentsTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.InvoicesRepoKey).(application.InvoiceRepository),
			c.Get(constants.PaymentsRepoKey).(application.PaymentRepository),
			c.Get(constants.RecurringRepoKey).(application.RecurringRepository),
			c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event]),
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
	ctx = container.Scoped(ctx)
	appAny := di.Get(ctx, constants.ApplicationKey) // mus
	realApp := appAny.(application.App)
	// 4) Prepare a minimal domain wrapper that matches rest.PaymentDomain
	domainWrapper := adapters.NewPaymentDomainWrapper(realApp)
	// 5) Retrieve your *stripe.StripeClient from container
	rest.RegisterWebhookRoute(container, svc.Mux(), domainWrapper, stripeClient)

	if err = grpc.RegisterServerTx(container, svc.RPC(), svc.Stripe()); err != nil {
		return err
	}

	if err = rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
		return err
	}
	if err = rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}
	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	handlers.RegisterDomainEventHandlersTx(container)
	if err = handlers.RegisterCommandHandlersTx(container); err != nil {
		return err
	}
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("payments outbox processor encountered an error")
		}
	}()
}
