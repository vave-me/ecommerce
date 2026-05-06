package newsletters

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
	"middleman/internal/registry/serdes"
	"middleman/internal/system"
	"middleman/internal/tm"
	"middleman/newsletters/internal/application"
	"middleman/newsletters/internal/constants"
	"middleman/newsletters/internal/domain"
	"middleman/newsletters/internal/grpc"
	"middleman/newsletters/internal/handlers"
	"middleman/newsletters/internal/postgres"
	"middleman/newsletters/internal/rest"
	"middleman/newsletters/newsletterspb"
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
		if err := newsletterspb.Registrations(reg); err != nil {
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

	// Newsletter aggregate repository
	container.AddScoped(constants.NewslettersRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Newsletter](
			domain.NewsletterAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	// Subscription aggregate repository
	container.AddScoped(constants.SubscriptionsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Subscription](
			domain.SubscriptionAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	// Edition aggregate repository
	container.AddScoped(constants.EditionsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Edition](
			domain.EditionAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	// Template aggregate repository
	container.AddScoped(constants.TemplatesRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Template](
			domain.TemplateAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	// Catalog repositories
	container.AddScoped(constants.NewsletterCatalogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewNewsletterCatalogRepository(
			constants.NewslettersTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.SubscriptionCatalogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewSubscriptionCatalogRepository(
			constants.SubscriptionsTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.EditionCatalogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewEditionCatalogRepository(
			constants.EditionsTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.TemplateCatalogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewTemplateCatalogRepository(
			constants.TemplatesTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.NewslettersRepoKey).(domain.NewsletterRepository),
			c.Get(constants.SubscriptionsRepoKey).(domain.SubscriptionRepository),
			c.Get(constants.EditionsRepoKey).(domain.EditionRepository),
			c.Get(constants.TemplatesRepoKey).(domain.TemplateRepository),
			c.Get(constants.NewsletterCatalogRepoKey).(domain.NewsletterCatalogRepository),
			c.Get(constants.SubscriptionCatalogRepoKey).(domain.SubscriptionCatalogRepository),
			c.Get(constants.EditionCatalogRepoKey).(domain.EditionCatalogRepository),
			c.Get(constants.TemplateCatalogRepoKey).(domain.TemplateCatalogRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
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
	if err = rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}
	handlers.RegisterDomainEventHandlersTx(container)

	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil
}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Newsletter aggregate
	if err = serde.Register(domain.Newsletter{}, func(v any) error {
		newsletter := v.(*domain.Newsletter)
		newsletter.Aggregate = es.NewAggregate("", domain.NewsletterAggregate)
		return nil
	}); err != nil {
		return
	}

	// Subscription aggregate
	if err = serde.Register(domain.Subscription{}, func(v any) error {
		subscription := v.(*domain.Subscription)
		subscription.Aggregate = es.NewAggregate("", domain.SubscriptionAggregate)
		return nil
	}); err != nil {
		return
	}

	// Edition aggregate
	if err = serde.Register(domain.Edition{}, func(v any) error {
		edition := v.(*domain.Edition)
		edition.Aggregate = es.NewAggregate("", domain.EditionAggregate)
		return nil
	}); err != nil {
		return
	}

	// Template aggregate
	if err = serde.Register(domain.Template{}, func(v any) error {
		template := v.(*domain.Template)
		template.Aggregate = es.NewAggregate("", domain.TemplateAggregate)
		return nil
	}); err != nil {
		return
	}

	// Newsletter events
	if err = serde.Register(domain.NewsletterCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.NewsletterUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.NewsletterActivated{}); err != nil {
		return
	}
	if err = serde.Register(domain.NewsletterDeactivated{}); err != nil {
		return
	}
	if err = serde.Register(domain.NewsletterDeleted{}); err != nil {
		return
	}

	// Subscription events
	if err = serde.Register(domain.Subscribed{}); err != nil {
		return
	}
	if err = serde.Register(domain.Unsubscribed{}); err != nil {
		return
	}
	if err = serde.Register(domain.PreferencesUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.SubscriptionPaused{}); err != nil {
		return
	}
	if err = serde.Register(domain.SubscriptionResumed{}); err != nil {
		return
	}

	// Edition events
	if err = serde.Register(domain.EditionCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.EditionUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.EditionScheduled{}); err != nil {
		return
	}
	if err = serde.Register(domain.EditionSending{}); err != nil {
		return
	}
	if err = serde.Register(domain.EditionSent{}); err != nil {
		return
	}

	// Template events
	if err = serde.Register(domain.TemplateCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.TemplateUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.TemplateDeleted{}); err != nil {
		return
	}

	// Newsletter snapshots
	if err = serde.RegisterKey(domain.NewsletterV1{}.SnapshotName(), domain.NewsletterV1{}); err != nil {
		return
	}

	// Subscription snapshots
	if err = serde.RegisterKey(domain.SubscriptionV1{}.SnapshotName(), domain.SubscriptionV1{}); err != nil {
		return
	}

	// Edition snapshots
	if err = serde.RegisterKey(domain.EditionV1{}.SnapshotName(), domain.EditionV1{}); err != nil {
		return
	}

	// Template snapshots
	if err = serde.RegisterKey(domain.TemplateV1{}.SnapshotName(), domain.TemplateV1{}); err != nil {
		return
	}

	return
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("newsletters outbox processor encountered an error")
		}
	}()
}