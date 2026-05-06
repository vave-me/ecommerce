package offers

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
	"middleman/offers/internal/application"
	"middleman/offers/internal/constants"
	"middleman/offers/internal/domain"
	"middleman/offers/internal/grpc"
	"middleman/offers/internal/handlers"
	"middleman/offers/internal/postgres"
	"middleman/offers/internal/rest"
	"middleman/offers/offerspb"
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
		if err := offerspb.Registrations(reg); err != nil {
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

			c.Get(constants.MessagePublisherKey).(am.MessagePublisher), svc.Logger(),
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

	container.AddScoped(constants.OffersRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Offer](
			domain.OfferAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.LeasingRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Lease](
			domain.LeaseAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.BuyBacksRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.BuyBack](
			domain.BuyBackAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.BuyNowRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.BuyNow](
			domain.BuyNowAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.ReservationRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Reservation](
			domain.ReservationAggregate,
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
			c.Get(constants.OffersRepoKey).(domain.OfferRepository),
			c.Get(constants.LeasingRepoKey).(domain.LeaseRepository),
			c.Get(constants.BuyBacksRepoKey).(domain.BuyBackRepository),
			c.Get(constants.BuyNowRepoKey).(domain.BuyNowRepository),
			c.Get(constants.ReservationRepoKey).(domain.ReservationRepository),
			c.Get(constants.MiddlemanRepoKey).(domain.MiddlemanRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
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
	handlers.RegisterMiddlemanHandlersTx(container)
	handlers.RegisterDomainEventHandlersTx(container)
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil
}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Register Offer aggregate
	if err = serde.Register(domain.Offer{}, func(v any) error {
		offer := v.(*domain.Offer)
		offer.Aggregate = es.NewAggregate("", domain.OfferAggregate)
		return nil
	}); err != nil {
		return
	}

	// Register Lease aggregate
	if err = serde.Register(domain.Lease{}, func(v any) error {
		lease := v.(*domain.Lease)
		lease.Aggregate = es.NewAggregate("", domain.LeaseAggregate)
		return nil
	}); err != nil {
		return
	}

	// Register BuyBack aggregate
	if err = serde.Register(domain.BuyBack{}, func(v any) error {
		buyBack := v.(*domain.BuyBack)
		buyBack.Aggregate = es.NewAggregate("", domain.BuyBackAggregate)
		return nil
	}); err != nil {
		return
	}

	// Register BuyNow aggregate
	if err = serde.Register(domain.BuyNow{}, func(v any) error {
		buyNow := v.(*domain.BuyNow)
		buyNow.Aggregate = es.NewAggregate("", domain.BuyNowAggregate)
		return nil
	}); err != nil {
		return
	}

	// Register Reservation aggregate
	if err = serde.Register(domain.Reservation{}, func(v any) error {
		reservation := v.(*domain.Reservation)
		reservation.Aggregate = es.NewAggregate("", domain.ReservationAggregate)
		return nil
	}); err != nil {
		return
	}

	// Register Offer events
	if err = serde.Register(domain.OfferCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.OfferActivated{}); err != nil {
		return
	}
	if err = serde.Register(domain.OfferAccepted{}); err != nil {
		return
	}
	if err = serde.Register(domain.OfferClosed{}); err != nil {
		return
	}

	// Register Lease events
	if err = serde.Register(domain.LeaseCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.LeaseStarted{}); err != nil {
		return
	}
	if err = serde.Register(domain.LeasePaymentMade{}); err != nil {
		return
	}
	if err = serde.Register(domain.LeaseBuyoutExecuted{}); err != nil {
		return
	}
	if err = serde.Register(domain.LeaseEnded{}); err != nil {
		return
	}
	if err = serde.Register(domain.LeaseCanceled{}); err != nil {
		return
	}
	if err = serde.Register(domain.LeaseDefaulted{}); err != nil {
		return
	}
	if err = serde.Register(domain.LeaseDeclined{}); err != nil {
		return
	}
	if err = serde.Register(domain.LeaseRejected{}); err != nil {
		return
	}
	if err = serde.Register(domain.LeaseNegotiationRequested{}); err != nil {
		return
	}
	if err = serde.Register(domain.LeaseNegotiationAccepted{}); err != nil {
		return
	}
	if err = serde.Register(domain.LeaseNegotiationUpdated{}); err != nil {
		return
	}

	// Register BuyNow events
	if err = serde.Register(domain.BuyNowCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.BuyNowConfirmed{}); err != nil {
		return
	}
	if err = serde.Register(domain.BuyNowCanceled{}); err != nil {
		return
	}
	if err = serde.Register(domain.BuyNowNegotiationRequested{}); err != nil {
		return
	}
	if err = serde.Register(domain.BuyNowNegotiationAccepted{}); err != nil {
		return
	}
	if err = serde.Register(domain.BuyNowNegotiationDeclined{}); err != nil {
		return
	}

	// Register BuyBack events
	if err = serde.Register(domain.BuyBackCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.BuyBackRedeemed{}); err != nil {
		return
	}
	if err = serde.Register(domain.BuyBackExpired{}); err != nil {
		return
	}
	if err = serde.Register(domain.BuyBackCanceled{}); err != nil {
		return
	}
	if err = serde.Register(domain.BuyBackNegotiationRequested{}); err != nil {
		return
	}
	if err = serde.Register(domain.BuyBackNegotiationAccepted{}); err != nil {
		return
	}
	if err = serde.Register(domain.BuyBackNegotiationDeclined{}); err != nil {
		return
	}

	// Register Reservation events
	if err = serde.Register(domain.ReservationCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.ReservationRedeemed{}); err != nil {
		return
	}
	if err = serde.Register(domain.ReservationExpired{}); err != nil {
		return
	}
	if err = serde.Register(domain.ReservationCanceled{}); err != nil {
		return
	}
	if err = serde.Register(domain.ReservationNegotiationRequested{}); err != nil {
		return
	}
	if err = serde.Register(domain.ReservationNegotiationAccepted{}); err != nil {
		return
	}
	if err = serde.Register(domain.ReservationNegotiationDeclined{}); err != nil {
		return
	}

	// Register snapshots
	if err = serde.RegisterKey(domain.OfferV1{}.SnapshotName(), domain.OfferV1{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.LeaseV1{}.SnapshotName(), domain.LeaseV1{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.BuyBackV1{}.SnapshotName(), domain.BuyBackV1{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.BuyNowV1{}.SnapshotName(), domain.BuyNowV1{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.ReservationV1{}.SnapshotName(), domain.ReservationV1{}); err != nil {
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
