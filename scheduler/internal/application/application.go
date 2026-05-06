package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/scheduler/internal/application/commands"
	"middleman/scheduler/internal/application/queries"
	"middleman/scheduler/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		CreateScheduler(ctx context.Context, cmd commands.CreateScheduler) error
		AddAction(ctx context.Context, cmd commands.AddAction) error
		RemoveAction(ctx context.Context, cmd commands.RemoveAction) error
		UpdateActionStatus(ctx context.Context, cmd commands.UpdateActionStatus) error
		// Task commands
		ScheduleTask(ctx context.Context, cmd commands.ScheduleTask) error
		CancelTask(ctx context.Context, cmd commands.CancelTask) error
		UpdateTask(ctx context.Context, cmd commands.UpdateTask) error
		ExecuteTask(ctx context.Context, cmd commands.ExecuteTask) error
		CompleteTask(ctx context.Context, cmd commands.CompleteTask) error
		FailTask(ctx context.Context, cmd commands.FailTask) error
	}
	Queries interface {
		GetAction(ctx context.Context, query queries.GetAction) (*domain.MiddlemanAction, error)
		GetActions(ctx context.Context, query queries.GetActions) ([]*domain.MiddlemanAction, error)
		GetSchedulers(ctx context.Context, query queries.GetSchedulers) ([]*domain.MiddlemanScheduler, error)
		GetScheduler(ctx context.Context, query queries.GetScheduler) (*domain.MiddlemanScheduler, error)
		GetPendingActions(ctx context.Context, query queries.GetPendingActions) ([]*domain.MiddlemanAction, error)
		// Task queries
		GetTask(ctx context.Context, query queries.GetTask) (*domain.CatalogTask, error)
		GetTasks(ctx context.Context, query queries.GetTasks) ([]*domain.CatalogTask, error)
		GetPendingTasks(ctx context.Context, query queries.GetPendingTasks) ([]*domain.CatalogTask, error)
		CountTasksByManagerID(ctx context.Context, query queries.CountTasksByManagerID) (int, error)
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.CreateSchedulerHandler
		commands.AddActionHandler
		commands.RemoveActionHandler
		commands.UpdateActionStatusHandler
		// Task handlers
		commands.ScheduleTaskHandler
		commands.CancelTaskHandler
		commands.UpdateTaskHandler
		commands.ExecuteTaskHandler
		commands.CompleteTaskHandler
		commands.FailTaskHandler
	}
	appQueries struct {
		queries.GetActionHandler
		queries.GetActionsHandler
		queries.GetSchedulerHandler
		queries.GetSchedulersHandler
		queries.GetPendingActionsHandler
		// Task handlers
		queries.GetTaskHandler
		queries.GetTasksHandler
		queries.GetPendingTasksHandler
		queries.CountTasksByManagerIDHandler
	}
)

var _ App = (*Application)(nil)

func New(scheduler domain.SchedulerRepository, actions domain.ActionRepository,
	middleman domain.MiddlemanRepository, itemActions domain.MiddlemanActionRepository,
	taskRepo domain.TaskRepository, catalogTaskRepo domain.CatalogTaskRepository,
	publisher ddd.EventPublisher[ddd.Event],
) *Application {
	return &Application{
		appCommands: appCommands{
			CreateSchedulerHandler:     commands.NewCreateSchedulerHandler(scheduler, publisher),
			AddActionHandler:           commands.NewAddActionHandler(actions, publisher),
			RemoveActionHandler:        commands.NewRemoveActionHandler(actions, publisher),
			UpdateActionStatusHandler:  commands.NewUpdateActionStatusHandler(actions, publisher),
			// Task handlers
			ScheduleTaskHandler:        commands.NewScheduleTaskHandler(taskRepo, publisher),
			CancelTaskHandler:          commands.NewCancelTaskHandler(taskRepo, publisher),
			UpdateTaskHandler:          commands.NewUpdateTaskHandler(taskRepo, publisher),
			ExecuteTaskHandler:         commands.NewExecuteTaskHandler(taskRepo, publisher),
			CompleteTaskHandler:        commands.NewCompleteTaskHandler(taskRepo, publisher),
			FailTaskHandler:            commands.NewFailTaskHandler(taskRepo, publisher),
		},
		appQueries: appQueries{
			GetSchedulersHandler:     queries.NewGetSchedulersHandler(middleman),
			GetSchedulerHandler:      queries.NewGetSchedulerHandler(middleman),
			GetActionHandler:         queries.NewGetActionHandler(itemActions),
			GetActionsHandler:        queries.NewGetActionsHandler(itemActions),
			GetPendingActionsHandler: queries.NewGetPendingActionsHandler(itemActions),
			// Task handlers
			GetTaskHandler:               queries.NewGetTaskHandler(catalogTaskRepo),
			GetTasksHandler:              queries.NewGetTasksHandler(catalogTaskRepo),
			GetPendingTasksHandler:       queries.NewGetPendingTasksHandler(catalogTaskRepo),
			CountTasksByManagerIDHandler: queries.NewCountTasksByManagerIDHandler(catalogTaskRepo),
		},
	}
}
