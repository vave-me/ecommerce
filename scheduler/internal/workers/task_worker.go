package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"middleman/scheduler/internal/application"
	"middleman/scheduler/internal/application/commands"
	"middleman/scheduler/internal/application/queries"
	"middleman/scheduler/internal/domain"
)

// TaskWorker processes scheduled tasks for managers
type TaskWorker struct {
	app       application.App
	logger    zerolog.Logger
	ticker    *time.Ticker
	done      chan bool
	wg        sync.WaitGroup
}

// NewTaskWorker creates a new task worker
func NewTaskWorker(app application.App, logger zerolog.Logger) *TaskWorker {
	return &TaskWorker{
		app:    app,
		logger: logger.With().Str("worker", "taskWorker").Logger(),
		done:   make(chan bool),
	}
}

// Start begins processing scheduled tasks
func (w *TaskWorker) Start(ctx context.Context) {
	w.ticker = time.NewTicker(30 * time.Second) // Check every 30 seconds
	w.wg.Add(1)
	
	go func() {
		defer w.wg.Done()
		
		// Process immediately on start
		w.processPendingTasks(ctx)
		
		for {
			select {
			case <-ctx.Done():
				w.logger.Info().Msg("task worker stopping due to context cancellation")
				return
			case <-w.done:
				w.logger.Info().Msg("task worker stopping")
				return
			case <-w.ticker.C:
				w.processPendingTasks(ctx)
			}
		}
	}()
	
	w.logger.Info().Msg("task worker started")
}

// Stop stops the worker
func (w *TaskWorker) Stop() {
	if w.ticker != nil {
		w.ticker.Stop()
	}
	close(w.done)
	w.wg.Wait()
	w.logger.Info().Msg("task worker stopped")
}

func (w *TaskWorker) processPendingTasks(ctx context.Context) {
	// Create a timeout context for this batch
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	
	// Get pending tasks that are due
	tasks, err := w.app.GetPendingTasks(ctx, queries.GetPendingTasks{
		BeforeTime: time.Now(),
		Limit:      100, // Process up to 100 tasks at a time
	})
	
	if err != nil {
		w.logger.Error().Err(err).Msg("failed to get pending tasks")
		return
	}
	
	if len(tasks) == 0 {
		return
	}
	
	w.logger.Info().Int("count", len(tasks)).Msg("processing pending tasks")
	
	// Process each task
	for _, task := range tasks {
		// Skip if context is done
		if ctx.Err() != nil {
			w.logger.Warn().Msg("stopping task processing due to context cancellation")
			break
		}
		
		w.processTask(ctx, task.ID)
	}
}

func (w *TaskWorker) processTask(ctx context.Context, taskID string) {
	logger := w.logger.With().Str("taskID", taskID).Logger()
	
	// Mark task as executing
	err := w.app.ExecuteTask(ctx, commands.ExecuteTask{
		TaskID: taskID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to mark task as executing")
		return
	}
	
	// Get the task details to execute it
	task, err := w.app.GetTask(ctx, queries.GetTask{
		TaskID: taskID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get task details")
		_ = w.app.FailTask(ctx, commands.FailTask{
			TaskID:       taskID,
			ErrorMessage: "failed to get task details: " + err.Error(),
		})
		return
	}
	
	// Execute task based on type
	result, err := w.executeTaskByType(ctx, task)
	if err != nil {
		logger.Error().Err(err).Msg("task execution failed")
		_ = w.app.FailTask(ctx, commands.FailTask{
			TaskID:       taskID,
			ErrorMessage: err.Error(),
		})
		return
	}
	
	// Mark task as completed
	err = w.app.CompleteTask(ctx, commands.CompleteTask{
		TaskID: taskID,
		Result: result,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to mark task as completed")
		
		// Try to mark as failed
		_ = w.app.FailTask(ctx, commands.FailTask{
			TaskID:       taskID,
			ErrorMessage: err.Error(),
		})
		return
	}
	
	logger.Info().Msg("task processed successfully")
}

func (w *TaskWorker) executeTaskByType(ctx context.Context, task *domain.CatalogTask) (string, error) {
	// This is where you would implement the actual task execution logic
	// based on the task type. For now, we'll just return a placeholder.
	
	// Example implementation:
	switch task.TaskType {
	case "analyze_patterns":
		// Call pattern analysis service
		return "Pattern analysis completed", nil
		
	case "generate_report":
		// Call report generation service
		return "Report generated successfully", nil
		
	case "sync_data":
		// Call data synchronization service
		return "Data synchronized", nil
		
	default:
		// For unknown task types, simulate execution
		time.Sleep(100 * time.Millisecond)
		return fmt.Sprintf("Task type '%s' executed with payload: %v", task.TaskType, task.Payload), nil
	}
}