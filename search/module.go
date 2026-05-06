package search

import (
	"context"
	"database/sql"
	"middleman/services/servicespb"

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

	"middleman/ordering/orderingpb"
	"middleman/posts/postspb"
	"middleman/products/productspb"

	"middleman/search/internal/application"
	"middleman/search/internal/constants"
	"middleman/search/internal/grpc"
	"middleman/search/internal/handlers"
	"middleman/search/internal/postgres"
	"middleman/search/internal/redis"
	"middleman/search/internal/rest"

	"middleman/users/userspb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.SearchService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.SearchService) (err error) {
	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := orderingpb.Registrations(reg); err != nil {
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
	container.AddScoped(constants.RedisearchClientKey, func(c di.Container) (any, error) { return *svc.Redisearch(), nil })

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
	// Register Products Fallback Repository (PostgreSQL)
	container.AddScoped(constants.UsersRepoKey, func(c di.Container) (any, error) {
		return postgres.NewUserCacheRepository(
			constants.UsersCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
			grpc.NewUserRepository(svc.Config().Rpc.Service(constants.UsersServiceName)),
		), nil
	})

	container.AddScoped(constants.ProductsRepoKey, func(c di.Container) (any, error) {
		return redis.NewProductCacheRepository(
			grpc.NewProductRepository(svc.Config().Rpc.Service(constants.ProductsServiceName)),
		), nil
	})

	container.AddScoped(constants.PostsRepoKey, func(c di.Container) (any, error) {
		return redis.NewPostCacheRepository(
			grpc.NewPostRepository(svc.Config().Rpc.Service(constants.PostsServiceName)),
		), nil
	})
	container.AddScoped(constants.ServicesRepoKey, func(c di.Container) (any, error) {
		return redis.NewServiceCacheRepository(
			grpc.NewServiceRepository(svc.Config().Rpc.Service(constants.ServicesServiceName)),
		), nil
	})
	container.AddScoped(constants.VariantsRepoKey, func(c di.Container) (any, error) {
		return redis.NewVariantCacheRepository(
			grpc.NewVariantRepository(svc.Config().Rpc.Service(constants.UsersServiceName)),
		), nil
	})
	container.AddScoped(constants.OrdersRepoKey, func(c di.Container) (any, error) {
		return postgres.NewOrderRepository(
			constants.OrdersTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// Add the MetricsRepo
	container.AddScoped(constants.MetricsRepoKey, func(c di.Container) (any, error) {
		return grpc.NewMetricRepository(svc.Config().Rpc.Service(constants.MetricsServiceName)), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.OrdersRepoKey).(application.OrderRepository),
			c.Get(constants.ProductsRepoKey).(application.ProductCacheRepository),
			c.Get(constants.VariantsRepoKey).(application.VariantCacheRepository),
			c.Get(constants.PostsRepoKey).(application.PostCacheRepository),
			c.Get(constants.UsersRepoKey).(application.UserCacheRepository),
			c.Get(constants.ServicesRepoKey).(application.ServiceCacheRepository),
			c.Get(constants.MetricsRepoKey).(application.MetricRepository),
		), nil
	})
	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewIntegrationEventHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.OrdersRepoKey).(application.OrderRepository),
			c.Get(constants.UsersRepoKey).(application.UserCacheRepository),
			c.Get(constants.ProductsRepoKey).(application.ProductCacheRepository),
			c.Get(constants.PostsRepoKey).(application.PostCacheRepository),
			c.Get(constants.ServicesRepoKey).(application.ServiceCacheRepository),
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
