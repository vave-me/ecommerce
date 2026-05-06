package geocoding

import (
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	"middleman/geocoding/geocodingpb"
	"middleman/geocoding/internal/application"
	"middleman/geocoding/internal/constants"
	"middleman/geocoding/internal/domain"
	"middleman/geocoding/internal/grpc"
	"middleman/geocoding/internal/handlers"
	"middleman/geocoding/internal/postgres"
	"middleman/geocoding/internal/rest"
	"middleman/internal/am"
	"middleman/internal/amotel"
	"middleman/internal/amprom"
	"middleman/internal/ddd"
	"middleman/internal/di"
	"middleman/internal/es"
	"middleman/internal/geo"
	"middleman/internal/jetstream"
	pg "middleman/internal/postgres"
	"middleman/internal/postgresotel"
	"middleman/internal/registry"
	"middleman/internal/registry/serdes"
	"middleman/internal/system"
	"middleman/internal/tm"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.GeocodingService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.GeocodingService) (err error) {

	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := geocodingpb.Registrations(reg); err != nil {
			return nil, err
		}

		return reg, nil
	})
	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())
	container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
	})

	geocode := geo.NewGoogleGeocodingClient(svc.GeocodingConfig().GoogleAPIKey)
	nominatim := geo.NewNominatimGeocodingClient("https://geo.sfx-markt.de")
	container.AddSingleton(constants.GoogleGeocode, func(c di.Container) (any, error) {
		return geocode, nil
	})
	container.AddSingleton(constants.NominatimGeocode, func(c di.Container) (any, error) {
		return nominatim, nil
	})

	container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DB().Begin()
	})

	container.AddSingleton(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream, svc.Logger(),
			amotel.OtelMessageContextExtractor(),
			amprom.ReceivedMessagesCounter(constants.ServiceName),
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

	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		return pg.NewInboxStore(constants.InboxTableName, tx), nil
	})

	container.AddScoped(constants.AddressesRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Address](
			domain.AddressAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})
	container.AddScoped(constants.CatalogRepoKey, func(c di.Container) (any, error) {
		return postgres.NewCatalogRepository(
			constants.CatalogTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	container.AddScoped(constants.LocationsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Location](
			domain.LocationAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.AddressesRepoKey).(domain.AddressRepository),
			c.Get(constants.LocationsRepoKey).(domain.LocationRepository),
			c.Get(constants.CatalogRepoKey).(domain.CatalogRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})
	//container.AddScoped(constants.CatalogHandlersKey, func(c di.Container) (any, error) {
	//	return handlers.NewCatalogHandlers(c.Get(constants.CatalogRepoKey).(domain.CatalogRepository)), nil
	//})
	//container.AddScoped(constants.CatalogLocationHandlersKey, func(c di.Container) (any, error) {
	//	return handlers.NewCatalogLocationHandlers(c.Get(constants.CatalogLocationRepoKey).(domain.CatalogLocationRepository)), nil
	//})

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

	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())
	return nil
}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(domain.Address{}, func(v any) error {
		address := v.(*domain.Address)
		address.Aggregate = es.NewAggregate("", domain.AddressAggregate)
		return nil
	}); err != nil {
		return
	}
	// product events
	if err = serde.Register(domain.AddressCreated{}); err != nil {
		return
	}

	if err = serde.Register(domain.Location{}, func(v any) error {
		location := v.(*domain.Location)
		location.Aggregate = es.NewAggregate("", domain.LocationAggregate)
		return nil
	}); err != nil {
		return
	}
	// product events
	if err = serde.Register(domain.LocationAdded{}); err != nil {
		return
	}

	if err = serde.RegisterKey(domain.LocationV1{}.SnapshotName(), domain.LocationV1{}); err != nil {
		return
	}

	return
}
func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("products outbox processor encountered an error")
		}
	}()
}
