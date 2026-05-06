package users

import (
	"context"
	"database/sql"
	"middleman/internal/am"
	"middleman/internal/amotel"
	"middleman/internal/amprom"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/es"
	"middleman/internal/jetstream"
	oidcclient "middleman/internal/oid"
	pg "middleman/internal/postgres"
	"middleman/internal/postgresotel"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
	"middleman/internal/system"
	"middleman/internal/tm"
	"middleman/users/internal/application"
	"middleman/users/internal/constants"
	"middleman/users/internal/domain"
	"middleman/users/internal/grpc"
	"middleman/users/internal/handlers"
	"middleman/users/internal/postgres"
	"middleman/users/internal/rest"
	"middleman/users/userspb"

	"github.com/rs/zerolog"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono system.UsersService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.UsersService) (err error) {

	cfg := svc.Config()

	// Initialize the Auth struct
	authJWT := &auth.Auth{
		Issuer:        cfg.JWTIssuer,
		Audience:      cfg.JWTAudience,
		Secret:        cfg.JWTSecret,
		TokenExpiry:   cfg.JWTTokenExpiry,
		RefreshExpiry: cfg.JWTRefreshExpiry,
		CookieDomain:  cfg.CookieDomain,
		CookiePath:    cfg.CookiePath,
		CookieName:    cfg.CookieName,
	}

	container := di.New()

	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := userspb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})

	mobileAuthClient := oidcclient.NewGoogleOIDCClient(ctx, svc.UsersConfig().MobileGoogleOAuthClientID, svc.UsersConfig().Issuer)
	webAuthClient := oidcclient.NewGoogleOIDCClient(ctx, svc.UsersConfig().WebGoogleOAuthClientID, svc.UsersConfig().Issuer)

	if mobileAuthClient == nil { // Defensive check
		return err
	}

	// Also ensure configurations are not empty before calling NewGoogleOIDCClient
	if svc.UsersConfig().MobileGoogleOAuthClientID == "" {
		return err
	}
	if svc.UsersConfig().WebGoogleOAuthClientID == "" {
		return err
	}
	if svc.UsersConfig().Issuer == "" { // Depending on whether your NewGoogleOIDCClient handles default issuer or requires it
		return err
	}

	container.AddSingleton(constants.MobileGoogleVerifierKey, func(c di.Container) (any, error) {
		// By this point, authClient should be non-nil if the error check above is implemented.
		// If you want an extra layer of safety within the DI factory:
		if mobileAuthClient == nil {
			return nil, err
		}
		return mobileAuthClient, nil
	})

	container.AddSingleton(constants.WebGoogleVerifierKey, func(c di.Container) (any, error) {
		// By this point, authClient should be non-nil if the error check above is implemented.
		// If you want an extra layer of safety within the DI factory:
		if webAuthClient == nil {
			return nil, err
		}
		return webAuthClient, nil
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
	container.AddScoped(constants.UsersRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.User](
			domain.UserAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
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
			c.Get(constants.UsersRepoKey).(domain.UserRepository),
			c.Get(constants.MiddlemanRepoKey).(domain.MiddlemanRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
			authJWT,
		), nil
	})

	container.AddScoped(constants.MiddlemanHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewMiddlemanHandlers(c.Get(constants.MiddlemanRepoKey).(domain.MiddlemanRepository)), nil
	})
	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})
	outboxProcessor := tm.NewOutboxProcessor(
		stream,
		pg.NewOutboxStore(constants.OutboxTableName, svc.DB()),
	)

	// setup Driver adapters with Google verifier
	if err = grpc.RegisterServerTx(container, svc.RPC()); err != nil {
		return err
	}
	if err = rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
		return err
	}

	//}
	if err = rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}
	handlers.RegisterMiddlemanHandlersTx(container)
	handlers.RegisterDomainEventHandlersTx(container)

	if err = userspb.RegisterAsyncAPI(svc.Mux()); err != nil {
		return err
	}
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil
}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Store
	if err = serde.Register(domain.User{}, func(v any) error {
		user := v.(*domain.User)
		user.Aggregate = es.NewAggregate("", domain.UserAggregate)
		return nil
	}); err != nil {
		return
	}
	// store events
	if err = serde.Register(domain.UserCreated{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.UserEnabledEvent, domain.UserEnabledToggled{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.UserDisabledEvent, domain.UserEnabledToggled{}); err != nil {
		return
	}
	if err = serde.Register(domain.UserLoggedIn{}); err != nil {
		return
	}
	if err = serde.Register(domain.UserPasswordResetRequested{}); err != nil {
		return
	}
	if err = serde.Register(domain.UserPasswordReset{}); err != nil {
		return
	}
	if err = serde.Register(domain.UserRenamed{}); err != nil {
		return
	}
	// store snapshots
	if err = serde.RegisterKey(domain.UserV1{}.SnapshotName(), domain.UserV1{}); err != nil {
		return
	}

	return
}
func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("stores outbox processor encountered an error")
		}
	}()
}
