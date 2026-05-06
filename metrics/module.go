package metric

import (
	"context"
	"database/sql"
	"middleman/activity/activitypb"
	"middleman/baskets/basketspb"
	"middleman/comments/commentspb"
	"middleman/following/followingpb"
	"middleman/internal/am"
	"middleman/internal/amotel"
	"middleman/internal/amprom"
	"middleman/internal/di"
	"middleman/internal/jetstream"
	pg "middleman/internal/postgres"
	"middleman/internal/postgresotel"
	"middleman/internal/registry"
	"middleman/internal/system"
	"middleman/internal/tm"
	"middleman/messages/messagespb"
	"middleman/metrics/internal/application"
	"middleman/metrics/internal/constants"
	"middleman/metrics/internal/grpc"
	"middleman/metrics/internal/handlers"
	"middleman/metrics/internal/postgres"
	"middleman/metrics/internal/redis"
	"middleman/metrics/internal/rest"
	"middleman/ordering/orderingpb"
	"middleman/posts/postspb"
	"middleman/products/productspb"
	"middleman/reviews/reviewspb"
	"middleman/services/servicespb"
	"middleman/users/userspb"
	"middleman/wishlists/wishlistspb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.MetricsService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.MetricsService) (err error) {
	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := activitypb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := basketspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := commentspb.Registrations(reg); err != nil {
			return nil, err
		}

		if err := followingpb.Registrations(reg); err != nil {
			return nil, err
		}

		if err := messagespb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := wishlistspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := reviewspb.Registrations(reg); err != nil {
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
		if err := servicespb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})

	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())
	container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DB().Begin()
	})
	container.AddScoped(constants.RedisTransactionKey, func(c di.Container) (any, error) { return *svc.Redis(), nil })
	container.AddSingleton(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream, svc.Logger(),
			amotel.OtelMessageContextExtractor(),
			amprom.ReceivedMessagesCounter(constants.ServiceName),
		), nil
	})
	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		return pg.NewInboxStore(constants.InboxTableName, tx), nil
	})

	container.AddScoped(constants.UsersMetricsCacheRepoKey, func(c di.Container) (any, error) {
		return postgres.NewUserMetricCacheRepository(
			constants.UsersMetricsCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.UserMetricsRepoKey, func(c di.Container) (any, error) {
		return redis.NewUserMetricRepository(
			constants.UsersMetricsCacheTableName,
			c.Get(constants.UsersMetricsCacheRepoKey).(application.UserMetricCacheRepository),
		), nil
	})

	container.AddScoped(constants.ItemMetricsCacheRepoKey, func(c di.Container) (any, error) {
		return postgres.NewItemMetricCacheRepository(
			constants.ItemMetricsCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.ItemMetricsRepoKey, func(c di.Container) (any, error) {
		return redis.NewItemMetricRepository(
			constants.ItemMetricsCacheTableName,
			c.Get(constants.ItemMetricsCacheRepoKey).(application.ItemMetricCacheRepository),
		), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.ItemMetricsRepoKey).(application.ItemMetricRepository),
			c.Get(constants.UserMetricsRepoKey).(application.UserMetricRepository),
		), nil
	})
	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewIntegrationEventHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ItemMetricsRepoKey).(application.ItemMetricRepository),
			c.Get(constants.UserMetricsRepoKey).(application.UserMetricRepository),
			tm.InboxHandler(c.Get(constants.InboxStoreKey).(tm.InboxStore)),
		), nil
	})

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
	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}

	return nil
}
