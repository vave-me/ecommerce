package shipping

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
	"middleman/shipping/internal/application"
	"middleman/shipping/internal/constants"
	"middleman/shipping/internal/dhl"
	"middleman/shipping/internal/domain"
	"middleman/shipping/internal/grpc"
	"middleman/shipping/internal/handlers"
	postgres "middleman/shipping/internal/postgres"
	"middleman/shipping/internal/rest"
	"middleman/shipping/shippingpb"
)

type Module struct {
}

func (m *Module) Startup(ctx context.Context, mono system.ShippingService) (err error) {
	return Root(ctx, mono)
}
func Root(ctx context.Context, svc system.ShippingService) (err error) {

	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := shippingpb.Registrations(reg); err != nil {
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
	// Use the existing DHLClient from config package for now
	// In production, you would update ShippingConfig to include the needed fields
	dhlClient := dhl.NewClient(
		svc.ShippingConfig().ClientID,      // Using ClientID as username
		svc.ShippingConfig().ClientSecret,  // Using ClientSecret as password
		svc.ShippingConfig().APIEndpoint,
		"", // Account number would need to be added to config
		true, // IsTest - would need to be added to config
	)

	container.AddSingleton(constants.DHLClient, func(c di.Container) (any, error) {
		return dhlClient, nil
	})

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

	container.AddScoped(constants.ShippingRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Shipment](
			domain.ShipmentAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.ShippingCatalogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewShippingCatalogRepository(
			constants.CatalogTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.ShippingRepoKey).(domain.ShippingRepository),
			c.Get(constants.ShippingCatalogRepoKey).(domain.ShippingCatalogRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})
	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})
	container.AddScoped(constants.CatalogHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCatalogHandlers(c.Get(constants.ShippingCatalogRepoKey).(domain.ShippingCatalogRepository)), nil
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
	
	// Register command handlers
	if err = handlers.RegisterCommandHandlersTx(container); err != nil {
		return err
	}

	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil

}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Store
	if err = serde.Register(domain.Shipment{}, func(v any) error {
		shipping := v.(*domain.Shipment)
		shipping.Aggregate = es.NewAggregate("", domain.ShipmentAggregate)
		return nil
	}); err != nil {
		return
	}
	// store events
	if err = serde.Register(domain.ShipmentCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.CarrierAssigned{}); err != nil {
		return
	}
	if err = serde.Register(domain.ShipmentStarted{}); err != nil {
		return
	}
	if err = serde.Register(domain.ShipmentStatusUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.ShipmentCancelled{}); err != nil {
		return
	}
	if err = serde.Register(domain.ShipmentDelivered{}); err != nil {
		return
	}
	if err = serde.Register(domain.PickupScheduled{}); err != nil {
		return
	}
	if err = serde.Register(domain.ShipmentReturned{}); err != nil {
		return
	}

	// store snapshots
	if err = serde.RegisterKey(domain.ShipmentV1{}.SnapshotName(), domain.ShipmentV1{}); err != nil {
		return
	}

	return
}
func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("shipping outbox processor encountered an error")
		}
	}()
}
