package consciousness

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"middleman/managers/internal/domain"
)

// ActionExecutor handles direct action execution through the scheduler
type ActionExecutor struct {
	schedulerRepo domain.SchedulerRepository
	toolRegistry  interface{} // Will be *tools.ToolServiceRegistry
	logger        zerolog.Logger
}

// NewActionExecutor creates a new action executor
func NewActionExecutor(params ...interface{}) *ActionExecutor {
	executor := &ActionExecutor{
		logger: zerolog.Logger{},
	}
	
	// Parse parameters
	for _, param := range params {
		switch v := param.(type) {
		case domain.SchedulerRepository:
			executor.schedulerRepo = v
		case zerolog.Logger:
			executor.logger = v
		default:
			// Assume it's the tool registry
			executor.toolRegistry = v
		}
	}
	
	return executor
}

// ExecuteActionWithType executes an action through the scheduler
func (e *ActionExecutor) ExecuteActionWithType(ctx context.Context, actionType string, parameters map[string]interface{}) error {
	// Build natural language task for scheduler
	task := e.buildActionTask(actionType, parameters)
	
	// Create unique action ID
	actionID := fmt.Sprintf("consciousness_%s_%d", actionType, time.Now().UnixNano())
	
	// Create scheduled action using the repository
	scheduledAction := &domain.ScheduledAction{
		ID:          actionID,
		Name:        fmt.Sprintf("Consciousness Action: %s", actionType),
		Description: task,
		EntityID:    "store_consciousness",
		EntityType:  "consciousness",
		Action:      actionType,
		Parameters:  parameters,
		ScheduledAt: time.Now(), // Execute immediately
		Status:      "pending",
		CreatedAt:   time.Now(),
	}
	
	err := e.schedulerRepo.ScheduleAction(ctx, scheduledAction)
	if err != nil {
		e.logger.Error().
			Err(err).
			Str("action_id", actionID).
			Str("action_type", actionType).
			Msg("Failed to schedule action")
		return fmt.Errorf("failed to schedule action: %w", err)
	}
	
	e.logger.Info().
		Str("action_id", actionID).
		Str("action_type", actionType).
		Msg("Action scheduled successfully")
	
	return nil
}

// buildActionTask creates a natural language task for the scheduler
func (e *ActionExecutor) buildActionTask(actionType string, parameters map[string]interface{}) string {
	parametersJSON, _ := json.Marshal(parameters)
	
	taskTemplates := map[string]string{
		"send_notification": "Send notification to user with parameters: %s",
		"update_metrics": "Update platform metrics with data: %s",
		"trigger_campaign": "Trigger marketing campaign with configuration: %s",
		"scale_resources": "Scale platform resources based on: %s",
		"analyze_data": "Analyze data and generate insights for: %s",
	}
	
	template, exists := taskTemplates[actionType]
	if !exists {
		template = "Execute %s action with parameters: %s"
		return fmt.Sprintf(template, actionType, string(parametersJSON))
	}
	
	return fmt.Sprintf(template, string(parametersJSON))
}

// ExecuteAction executes an Action struct
func (e *ActionExecutor) ExecuteAction(ctx context.Context, action Action) error {
	return e.ExecuteActionWithType(ctx, action.Type, action.Parameters)
}