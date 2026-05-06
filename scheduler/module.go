package scheduler

import (
	"context"
	"database/sql"
	"os"
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
	"middleman/scheduler/internal/application"
	"middleman/scheduler/internal/constants"
	"middleman/scheduler/internal/domain"
	"middleman/scheduler/internal/grpc"
	"middleman/scheduler/internal/handlers"
	"middleman/scheduler/internal/postgres"
	"middleman/scheduler/internal/redis_repository"
	"middleman/scheduler/internal/rest"
	"middleman/scheduler/internal/workers"
	"middleman/scheduler/schedulerspb"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono system.SchedulerService) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.SchedulerService) (err error) {
	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := domain.Registrations(reg); err != nil {
			return nil, err
		}
		if err := schedulerspb.Registrations(reg); err != nil {
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
	container.AddScoped(constants.RedisPoolKey, func(c di.Container) (any, error) { return svc.RedisPoolScheduler(), nil })

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

	container.AddScoped(constants.CommandPublisherKey, func(c di.Container) (any, error) {
		return am.NewCommandPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
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
	container.AddScoped(constants.AggregateStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return es.AggregateStoreWithMiddleware(
			pg.NewEventStore(constants.EventsTableName, tx, reg),
			pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg),
		), nil
	})
	container.AddScoped(constants.SchedulerRepoKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return es.NewAggregateRepository[*domain.Scheduler](
			domain.SchedulerAggregate,
			reg,
			es.AggregateStoreWithMiddleware(
				pg.NewEventStore(constants.EventsTableName, tx, reg),
				pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg),
			),
		), nil
	})

	container.AddScoped(constants.ActionsRepoKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return es.NewAggregateRepository[*domain.Action](
			domain.ActionAggregate,
			reg,
			es.AggregateStoreWithMiddleware(
				pg.NewEventStore(constants.EventsTableName, tx, reg),
				pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg),
			),
		), nil
	})
	container.AddSingleton(constants.AssistantRepoKey, func(c di.Container) (any, error) {
		// Get assistant ID from environment or use default
		assistantID := constants.DefaultAssistantID
		if envID := os.Getenv("SCHEDULER_ASSISTANT_ID"); envID != "" {
			assistantID = envID
		}
		return grpc.NewAssistantRepository(
			svc.Config().Rpc.Service(constants.AssistantsServiceName),
			assistantID,
		), nil
	})

	container.AddScoped(constants.MiddlemanCacheActionRepoKey, func(c di.Container) (any, error) {
		return redis_repository.NewMiddlemanCacheActionRepository(
			constants.MiddlemanActionTableName,
			c.Get(constants.MiddlemanActionRepoKey).(domain.MiddlemanActionRepository),
		), nil
	})

	container.AddScoped(constants.MiddlemanCacheRepoKey, func(c di.Container) (any, error) {
		return redis_repository.NewMiddlemanCacheRepository(
			constants.MiddlemanCacheTableName,
			c.Get(constants.MiddlemanRepoKey).(domain.MiddlemanRepository),
		), nil
	})

	container.AddScoped(constants.MiddlemanRepoKey, func(c di.Container) (any, error) {
		return postgres.NewMiddlemanRepository(
			constants.MiddlemanTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	container.AddScoped(constants.MiddlemanActionRepoKey, func(c di.Container) (any, error) {
		return postgres.NewMiddlemanActionRepository(
			constants.MiddlemanActionTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	
	// Task repositories
	container.AddScoped(constants.TaskRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Task](
			domain.TaskAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})
	
	container.AddScoped(constants.CatalogTaskRepoKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		return postgres.NewCatalogTaskRepository(constants.CatalogTaskTableName, tx), nil
	})
	// TODO implement prometheus counters

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.SchedulerRepoKey).(domain.SchedulerRepository),
			c.Get(constants.ActionsRepoKey).(domain.ActionRepository),
			c.Get(constants.MiddlemanRepoKey).(domain.MiddlemanRepository),
			c.Get(constants.MiddlemanActionRepoKey).(domain.MiddlemanActionRepository),
			c.Get(constants.TaskRepoKey).(domain.TaskRepository),
			c.Get(constants.CatalogTaskRepoKey).(domain.CatalogTaskRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event])), nil
	})
	container.AddScoped(constants.MiddlemanHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewMiddlemanHandlers(c.Get(constants.MiddlemanCacheRepoKey).(domain.MiddlemanCacheRepository)), nil
	})
	container.AddScoped(constants.MiddlemanActionHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewMiddlemanActionHandlers(
			c.Get(constants.MiddlemanCacheActionRepoKey).(domain.MiddlemanCacheActionRepository)), nil
	})
	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})
	
	container.AddScoped(constants.TaskHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewTaskHandlers(
			c.Get(constants.CatalogTaskRepoKey).(domain.CatalogTaskRepository),
			svc.Logger(),
		), nil
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
	handlers.RegisterMiddlemanHandlersTx(container)
	handlers.RegisterMiddlemanActionHandlersTx(container)
	handlers.RegisterTaskHandlersTx(container)
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())
	
	// Start the scheduler worker
	schedulerWorker := workers.NewSchedulerWorker(
		container.Get(constants.ApplicationKey).(application.App),
		svc.Logger(),
		container.Get(constants.AssistantRepoKey).(domain.AssistantRepository),
	)
	schedulerWorker.Start(ctx)
	
	// Start the task worker
	taskWorker := workers.NewTaskWorker(
		container.Get(constants.ApplicationKey).(application.App),
		svc.Logger(),
	)
	taskWorker.Start(ctx)
	
	// Add cleanup function to stop the workers
	svc.Waiter().Add(func(ctx context.Context) error {
		schedulerWorker.Stop()
		taskWorker.Stop()
		return nil
	})
	
	return
}
func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Store
	if err = serde.Register(domain.Scheduler{}, func(v any) error {
		scheduler := v.(*domain.Scheduler)
		scheduler.Aggregate = es.NewAggregate("", domain.SchedulerAggregate)
		return nil
	}); err != nil {
		return
	}

	if err = serde.Register(domain.Action{}, func(v any) error {
		interaction := v.(*domain.Action)
		interaction.Aggregate = es.NewAggregate("", domain.ActionAggregate)
		return nil
	}); err != nil {
		return
	}
	
	if err = serde.Register(domain.Task{}, func(v any) error {
		task := v.(*domain.Task)
		task.Aggregate = es.NewAggregate("", domain.TaskAggregate)
		return nil
	}); err != nil {
		return
	}
	// store events
	if err = serde.Register(domain.SchedulerCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.ActionAdded{}); err != nil {
		return
	}
	if err = serde.Register(domain.ActionUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.ActionRemoved{}); err != nil {
		return
	}
	
	// Task events
	if err = serde.Register(domain.TaskScheduled{}); err != nil {
		return
	}
	if err = serde.Register(domain.TaskCancelled{}); err != nil {
		return
	}
	if err = serde.Register(domain.TaskUpdated{}); err != nil {
		return
	}
	if err = serde.Register(domain.TaskExecutionStarted{}); err != nil {
		return
	}
	if err = serde.Register(domain.TaskExecutionCompleted{}); err != nil {
		return
	}
	if err = serde.Register(domain.TaskExecutionFailed{}); err != nil {
		return
	}

	// store snapshots
	if err = serde.RegisterKey(domain.SchedulerVi{}.SnapshotName(), domain.SchedulerVi{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.ActionVi{}.SnapshotName(), domain.ActionVi{}); err != nil {
		return
	}
	if err = serde.RegisterKey(domain.TaskSnapshot{}.SnapshotName(), domain.TaskSnapshot{}); err != nil {
		return
	}
	return
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("scheduler outbox processor encountered an error")
		}
	}()
}
