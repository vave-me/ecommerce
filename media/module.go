package media

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	"middleman/internal/am"
	"middleman/internal/amotel"
	"middleman/internal/amprom"
	"middleman/internal/config"
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
	"middleman/media/internal/application"
	"middleman/media/internal/application/commands"
	"middleman/media/internal/constants"
	"middleman/media/internal/domain"
	"middleman/media/internal/grpc"
	"middleman/media/internal/handlers"
	"middleman/media/internal/postgres"
	"middleman/media/internal/rest"
	"middleman/media/mediapb"
)

type Module struct {
}

func (m *Module) Startup(ctx context.Context, mono system.MediaService) (err error) {
	return Root(ctx, mono)
}
func Root(ctx context.Context, svc system.MediaService) (err error) {

	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := mediapb.Registrations(reg); err != nil {
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

	minioClient, err := config.NewMinioClient(ctx, svc.MediaConfig())
	if err != nil {
		return fmt.Errorf("failed to init minio client: %w", err)
	}
	// 3) Put minioClient into DI
	container.AddSingleton(constants.MinioClient, func(c di.Container) (any, error) {
		return minioClient, nil
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

	container.AddScoped(constants.MediaRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Media](
			domain.MediaAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.ImagesRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Image](
			domain.ImageAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})
	container.AddScoped(constants.VideosRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Video](
			domain.VideoAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})
	container.AddScoped(constants.ImporterRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Importer](
			domain.ImporterAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})

	container.AddScoped(constants.MiddlemanMediaRepoKey, func(c di.Container) (any, error) {
		return postgres.NewMiddlemanMediaRepository(
			constants.MiddlemanMediaTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	container.AddScoped(constants.MiddlemanImageRepoKey, func(c di.Container) (any, error) {
		return postgres.NewMiddlemanImageRepository(
			constants.MiddlemanImageTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	container.AddScoped(constants.MiddlemanVideoRepoKey, func(c di.Container) (any, error) {
		return postgres.NewMiddlemanVideoRepository(
			constants.MiddlemanVideoTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	container.AddScoped(constants.ImportSessionRepoKey, func(c di.Container) (any, error) {
		return postgres.NewImportSessionRepository(
			constants.ImportSessionTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	container.AddScoped(constants.ImportItemRepoKey, func(c di.Container) (any, error) {
		return postgres.NewImportItemRepository(
			constants.ImportItemTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	container.AddScoped(constants.ProductRepoKey, func(c di.Container) (any, error) {
		return grpc.NewProductRepository(svc.Config().Rpc.ServiceAddress(constants.ProductsServiceName)), nil
	})
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.MediaRepoKey).(domain.MediaRepository),
			c.Get(constants.VideosRepoKey).(domain.VideoRepository),
			c.Get(constants.ImagesRepoKey).(domain.ImageRepository),
			c.Get(constants.MiddlemanMediaRepoKey).(domain.MiddlemanMediaRepository),
			c.Get(constants.MiddlemanImageRepoKey).(domain.MiddlemanImageRepository),
			c.Get(constants.MiddlemanVideoRepoKey).(domain.MiddlemanVideoRepository),
			c.Get(constants.ImporterRepoKey).(domain.ImporterRepository),
			c.Get(constants.ImportSessionRepoKey).(domain.ImportSessionRepository),
			c.Get(constants.ImportItemRepoKey).(domain.ImportItemRepository),
			c.Get(constants.ProductRepoKey).(commands.ProductRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})
	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})
	container.AddScoped(constants.MiddlemanMediaHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewMiddlemanMediaHandlers(c.Get(constants.MiddlemanMediaRepoKey).(domain.MiddlemanMediaRepository)), nil
	})
	container.AddScoped(constants.MiddlemanImageHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewMiddlemanImageHandlers(c.Get(constants.MiddlemanImageRepoKey).(domain.MiddlemanImageRepository)), nil
	})

	container.AddScoped(constants.MiddlemanVideoHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewMiddlemanVideoHandlers(c.Get(constants.MiddlemanVideoRepoKey).(domain.MiddlemanVideoRepository)), nil
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
	handlers.RegisterMiddlemanMediaHandlersTx(container)
	handlers.RegisterMiddlemanImageHandlersTx(container)
	handlers.RegisterMiddlemanVideoHandlersTx(container)
	//handlers.RegisterMallHandlersTx(container)
	handlers.RegisterDomainEventHandlersTx(container)

	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil

}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Store
	if err = serde.Register(domain.Media{}, func(v any) error {
		media := v.(*domain.Media)
		media.Aggregate = es.NewAggregate("", domain.MediaAggregate)
		return nil
	}); err != nil {
		return
	}
	// store events
	if err = serde.Register(domain.MediaCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.MediaUpdated{}); err != nil {
		return
	}

	if err = serde.Register(domain.MediaStatusChanged{}); err != nil {
		return
	}
	if err = serde.Register(domain.MediaDeleted{}); err != nil {
		return
	}
	// store snapshots
	if err = serde.RegisterKey(domain.MediaV1{}.SnapshotName(), domain.MediaV1{}); err != nil {
		return
	}

	// Image
	if err = serde.Register(domain.Image{}, func(v any) error {
		image := v.(*domain.Image)
		image.Aggregate = es.NewAggregate("", domain.ImageAggregate)
		return nil
	}); err != nil {
		return
	}
	// image events
	if err = serde.Register(domain.ImageAdded{}); err != nil {
		return
	}
	//if err = serde.Register(domain.ImageOrderChanged{}); err != nil {
	//	return
	//}
	//if err = serde.Register(domain.ImageMetadataUpdated{}); err != nil {return}
	if err = serde.Register(domain.ImageRemoved{}); err != nil {
		return
	}

	if err = serde.RegisterKey(domain.ImageV1{}.SnapshotName(), domain.ImageV1{}); err != nil {
		return
	}

	// Video
	if err = serde.Register(domain.Video{}, func(v any) error {
		image := v.(*domain.Video)
		image.Aggregate = es.NewAggregate("", domain.VideoAggregate)
		return nil
	}); err != nil {
		return
	}
	// image events
	if err = serde.Register(domain.VideoAdded{}); err != nil {
		return
	}
	//if err = serde.Register(domain.VideoOrderChanged{}); err != nil {
	//	return
	//}
	//if err = serde.Register(domain.VideoMetadataUpdated{}); err != nil {return}
	if err = serde.Register(domain.VideoRemoved{}); err != nil {
		return
	}

	if err = serde.RegisterKey(domain.VideoV1{}.SnapshotName(), domain.VideoV1{}); err != nil {
		return
	}

	// Importer
	if err = serde.Register(domain.Importer{}, func(v any) error {
		importer := v.(*domain.Importer)
		importer.Aggregate = es.NewAggregate("", domain.ImporterAggregate)
		return nil
	}); err != nil {
		return
	}
	// importer events
	if err = serde.Register(domain.BulkImportStarted{}); err != nil {
		return
	}
	if err = serde.Register(domain.ImportBatchAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.ImportItemProcessed{}); err != nil {
		return
	}
	if err = serde.Register(domain.ImportItemFailed{}); err != nil {
		return
	}
	if err = serde.Register(domain.BulkImportCompleted{}); err != nil {
		return
	}
	if err = serde.Register(domain.BulkImportCancelled{}); err != nil {
		return
	}
	// importer snapshot
	if err = serde.RegisterKey(domain.ImporterV1{}.SnapshotName(), domain.ImporterV1{}); err != nil {
		return
	}

	return
}
func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("media outbox processor encountered an error")
		}
	}()
}
