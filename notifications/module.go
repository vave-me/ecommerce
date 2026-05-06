package notifications

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	"middleman/baskets/basketspb"
	"middleman/comments/commentspb"
	"middleman/following/followingpb"
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
	"middleman/messages/messagespb"
	"middleman/notifications/internal/application"
	"middleman/notifications/internal/constants"
	"middleman/notifications/internal/domain"
	"middleman/notifications/internal/grpc"
	"middleman/notifications/internal/handlers"
	"middleman/notifications/internal/infra"
	"middleman/notifications/internal/postgres"
	"middleman/notifications/internal/rest"
	"middleman/notifications/notificationspb"
	"middleman/offers/offerspb"
	"middleman/ordering/orderingpb"
	"middleman/payments/paymentspb"
	"middleman/posts/postspb"
	"middleman/products/productspb"
	"middleman/reviews/reviewspb"
	"middleman/support/supportpb"
	"middleman/users/userspb"
	"middleman/wishlists/wishlistspb"
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
		if err := notificationspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := userspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := productspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := basketspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := commentspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := messagespb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := postspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := reviewspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := followingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := offerspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := paymentspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := supportpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := wishlistspb.Registrations(reg); err != nil {
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
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher), svc.Logger()), nil
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

	container.AddScoped(constants.AlertsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Alert](
			domain.AlertAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
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
			grpc.NewProductRepository(svc.Config().Rpc.Service(constants.ProductsServiceName)),
		), nil
	})
	container.AddScoped(constants.CatalogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewCatalogRepository(
			constants.CatalogTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// Add preferences repository - using in-memory for now
	container.AddSingleton(constants.PreferencesRepoKey, func(c di.Container) (any, error) {
		return infra.NewInMemoryPreferencesRepository(), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.AlertsRepoKey).(domain.AlertRepository),
			c.Get(constants.CatalogRepoKey).(domain.CatalogRepository),
			c.Get(constants.PreferencesRepoKey).(domain.PreferencesRepository),
			c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])), nil
	})

	container.AddScoped(constants.MiddlemanHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewMiddlemanHandlers(c.Get(constants.CatalogRepoKey).(domain.CatalogRepository)), nil
	})

	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})

	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewIntegrationEventHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ApplicationKey).(application.App),
			c.Get(constants.CatalogRepoKey).(domain.CatalogRepository),
			c.Get(constants.UsersRepoKey).(domain.UserCacheRepository),
			c.Get(constants.ProductsRepoKey).(domain.ProductCacheRepository),
			svc.Logger(),
		), nil
	})

	outboxProcessor := tm.NewOutboxProcessor(
		stream,
		pg.NewOutboxStore(constants.OutboxTableName, svc.DB()),
	)

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
	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())
	return
}
func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(domain.Alert{}, func(v any) error {
		alert := v.(*domain.Alert)
		alert.Aggregate = es.NewAggregate("", domain.AlertAggregate)
		return nil
	}); err != nil {
		return
	}
	// Alert events
	if err = serde.Register(domain.ProductAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.BasketAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.OrderAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.WishlistAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.MessageAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.InteractionAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.CommentAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.OfferAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.SupportAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.ReviewAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.PaymentAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.FollowingAlertAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.AlertRead{}); err != nil {
		return
	}

	if err = serde.RegisterKey(domain.AlertV1{}.SnapshotName(), domain.AlertV1{}); err != nil {
		return
	}
	return
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("notifications outbox processor encountered an error")
		}
	}()
}
