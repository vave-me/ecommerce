package workers

import (
	"context"
	"time"
	
	"github.com/rs/zerolog"
	"middleman/scheduler/internal/application"
	"middleman/scheduler/internal/application/commands"
	"middleman/scheduler/internal/application/queries"
	"middleman/scheduler/internal/domain"
)

type SchedulerWorker struct {
	app                 application.App
	logger              zerolog.Logger
	assistantRepository domain.AssistantRepository
	ticker              *time.Ticker
	done                chan bool
}

func NewSchedulerWorker(app application.App, logger zerolog.Logger, assistantRepository domain.AssistantRepository) *SchedulerWorker {
	return &SchedulerWorker{
		app:                 app,
		logger:              logger,
		assistantRepository: assistantRepository,
		done:                make(chan bool),
	}
}

func (w *SchedulerWorker) Start(ctx context.Context) {
	// Check for pending tasks every 30 seconds
	w.ticker = time.NewTicker(30 * time.Second)
	
	w.logger.Info().Msg("Scheduler worker started")
	
	// Run once immediately
	w.processPendingActions(ctx)
	
	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.processPendingActions(ctx)
			case <-w.done:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (w *SchedulerWorker) Stop() {
	w.logger.Info().Msg("Stopping scheduler worker")
	if w.ticker != nil {
		w.ticker.Stop()
	}
	close(w.done)
}

func (w *SchedulerWorker) processPendingActions(ctx context.Context) {
	// Get all pending actions that should be executed by now
	actions, err := w.app.GetPendingActions(ctx, queries.GetPendingActions{
		BeforeTime: time.Now(),
	})
	if err != nil {
		w.logger.Error().Err(err).Msg("Failed to get pending actions")
		return
	}
	
	if len(actions) == 0 {
		return
	}
	
	w.logger.Info().Int("count", len(actions)).Msg("Processing pending actions")
	
	// Process each action
	for _, action := range actions {
		// Create a new context with timeout for each action
		actionCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		
		w.processAction(actionCtx, action.ID, action.NaturalLanguageTask)
		
		cancel()
	}
}

func (w *SchedulerWorker) processAction(ctx context.Context, actionID, task string) {
	logger := w.logger.With().
		Str("action_id", actionID).
		Str("task", task).
		Logger()
		
	logger.Info().Msg("Processing action")
	
	// Update status to executing
	err := w.app.UpdateActionStatus(ctx, commands.UpdateActionStatus{
		ID:     actionID,
		Status: "executing",
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update action status to executing")
		return
	}
	
	// Execute the task through the assistant service
	result, err := w.assistantRepository.ProcessUserInput(ctx, task, map[string]string{
		"action_id": actionID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to execute task")
		
		// Update status to failed
		updateErr := w.app.UpdateActionStatus(ctx, commands.UpdateActionStatus{
			ID:           actionID,
			Status:       "failed",
			ErrorMessage: err.Error(),
		})
		if updateErr != nil {
			logger.Error().Err(updateErr).Msg("Failed to update action status to failed")
		}
		return
	}
	
	// Update status to completed
	err = w.app.UpdateActionStatus(ctx, commands.UpdateActionStatus{
		ID:     actionID,
		Status: "completed",
		Result: result,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update action status to completed")
		return
	}
	
	logger.Info().Msg("Action processed successfully")
}