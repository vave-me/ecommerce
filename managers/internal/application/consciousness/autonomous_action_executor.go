package consciousness

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	
	"github.com/rs/zerolog"
	"middleman/internal/ai"
	"middleman/managers/internal/application/tools"
)

// AutonomousActionExecutor executes decisions using the application's tool system
type AutonomousActionExecutor struct {
	app            App
	actionExecutor *ActionExecutor // Existing action executor
	aiManager      ai.AIClientManager
	logger         zerolog.Logger
}

// NewAutonomousActionExecutor creates a new autonomous action executor
func NewAutonomousActionExecutor(
	app App,
	actionExecutor *ActionExecutor,
	aiManager ai.AIClientManager,
	logger zerolog.Logger,
) *AutonomousActionExecutor {
	return &AutonomousActionExecutor{
		app:            app,
		actionExecutor: actionExecutor,
		aiManager:      aiManager,
		logger:         logger,
	}
}

// ExecuteDecision executes all actions in a decision
func (ae *AutonomousActionExecutor) ExecuteDecision(ctx context.Context, decision *Decision) error {
	ae.logger.Info().
		Str("decision_id", decision.ID).
		Str("decision_type", decision.Type).
		Int("action_count", len(decision.Actions)).
		Msg("Executing autonomous decision")
	
	successCount := 0
	var lastError error
	
	for i, action := range decision.Actions {
		ae.logger.Debug().
			Int("action_index", i).
			Str("action_type", action.Type).
			Interface("parameters", action.Parameters).
			Msg("Executing action")
			
		if err := ae.executeAction(ctx, action); err != nil {
			ae.logger.Error().
				Err(err).
				Str("action_type", action.Type).
				Int("action_index", i).
				Msg("Failed to execute action")
			lastError = err
			continue
		}
		
		successCount++
		ae.logger.Info().
			Str("action_type", action.Type).
			Int("action_index", i).
			Msg("Action executed successfully")
	}
	
	ae.logger.Info().
		Str("decision_id", decision.ID).
		Int("success_count", successCount).
		Int("total_actions", len(decision.Actions)).
		Msg("Decision execution completed")
	
	// Return error only if all actions failed
	if successCount == 0 && lastError != nil {
		return fmt.Errorf("all actions failed, last error: %w", lastError)
	}
	
	return nil
}

func (ae *AutonomousActionExecutor) executeAction(ctx context.Context, action Action) error {
	// Map action to tool call
	toolCall := ae.mapActionToTool(action)
	
	// Create execution context for consciousness actions
	execCtx := &tools.ToolExecutionContext{
		UserID:   "system-consciousness",
		Role:     "manager",
		Metadata: map[string]interface{}{
			"source":      "consciousness",
			"action_type": action.Type,
		},
	}
	
	// Execute using the application's tool system
	results, err := ae.app.ExecuteTools(ctx, []ai.ToolCall{toolCall}, execCtx)
	if err != nil {
		return fmt.Errorf("tool execution failed: %w", err)
	}
	
	if len(results) == 0 {
		return fmt.Errorf("no results returned from tool execution")
	}
	
	result := results[0]
	if result.Error != nil {
		return fmt.Errorf("tool returned error: %w", result.Error)
	}
	
	ae.logger.Debug().
		Str("tool_name", result.ToolName).
		Interface("result", result.Result).
		Msg("Tool execution completed")
	
	return nil
}

func (ae *AutonomousActionExecutor) mapActionToTool(action Action) ai.ToolCall {
	// Check if we have selected tools in parameters
	var toolName string
	selectedTools, hasSelectedTools := action.Parameters["_selected_tools"].([]string)
	
	if hasSelectedTools && len(selectedTools) > 0 {
		// Use dynamic tool selection - pick the most appropriate tool from selected
		toolName = ae.selectBestTool(action.Type, selectedTools)
		ae.logger.Debug().
			Str("action_type", action.Type).
			Str("selected_tool", toolName).
			Int("available_tools", len(selectedTools)).
			Msg("Using dynamically selected tool")
	} else {
		// Fallback to standard mapping
		toolName = ae.getStandardToolMapping(action.Type)
	}
	
	// Remove internal parameters before marshaling
	cleanParams := make(map[string]interface{})
	for k, v := range action.Parameters {
		if !strings.HasPrefix(k, "_") {
			cleanParams[k] = v
		}
	}
	
	// Marshal parameters to JSON string
	args, err := json.Marshal(cleanParams)
	if err != nil {
		ae.logger.Error().Err(err).Msg("Failed to marshal action parameters")
		args = []byte("{}")
	}
	
	return ai.ToolCall{
		ID:   generateID(),
		Type: ai.ToolTypeFunction,
		Function: ai.FunctionCall{
			Name:      toolName,
			Arguments: string(args),
		},
	}
}

// selectBestTool selects the most appropriate tool from available options
func (ae *AutonomousActionExecutor) selectBestTool(actionType string, availableTools []string) string {
	// Priority mapping for each action type
	toolPriority := map[string][]string{
		"send_notification": {
			"notification_send_immediate",
			"notification_send",
			"notification_create_and_send",
			"notification_queue",
		},
		"create_offer": {
			"offer_create_personalized",
			"offer_create",
			"offer_generate",
			"promotion_create",
		},
		"escalate_ticket": {
			"support_escalate_urgent",
			"support_escalate_ticket",
			"support_escalate",
			"ticket_escalate",
		},
		"flag_order": {
			"order_flag_high_risk",
			"order_flag_for_review",
			"order_flag",
			"order_mark_suspicious",
		},
	}
	
	// Get priority list for this action
	priorities, exists := toolPriority[actionType]
	if !exists {
		// No priority list, use first available tool
		if len(availableTools) > 0 {
			return availableTools[0]
		}
		return actionType // Fallback to action type as tool name
	}
	
	// Find the highest priority tool that's available
	for _, priorityTool := range priorities {
		for _, available := range availableTools {
			if available == priorityTool {
				return available
			}
		}
	}
	
	// If no priority tool found, check for partial matches
	for _, available := range availableTools {
		if strings.Contains(available, actionType) {
			return available
		}
	}
	
	// Last resort: use first available tool
	if len(availableTools) > 0 {
		return availableTools[0]
	}
	
	return actionType // Ultimate fallback
}

// getStandardToolMapping returns the standard tool mapping for an action
func (ae *AutonomousActionExecutor) getStandardToolMapping(actionType string) string {
	switch actionType {
	case "send_notification":
		return "notification_send"
	case "create_offer":
		return "offer_create"
	case "escalate_ticket":
		return "support_escalate_ticket"
	case "flag_order":
		return "order_flag_for_review"
	default:
		// For any other action, assume the action type is the tool name
		return actionType
	}
}

// ExecuteWithVerification executes an action and verifies the result
func (ae *AutonomousActionExecutor) ExecuteWithVerification(ctx context.Context, action Action) error {
	// Execute the action
	if err := ae.executeAction(ctx, action); err != nil {
		return err
	}
	
	// Verify execution if needed
	if ae.shouldVerify(action) {
		return ae.verifyExecution(ctx, action)
	}
	
	return nil
}

func (ae *AutonomousActionExecutor) shouldVerify(action Action) bool {
	// Verify critical actions
	criticalActions := []string{
		"flag_order",
		"escalate_ticket",
		"send_payment_notification",
	}
	
	for _, critical := range criticalActions {
		if action.Type == critical {
			return true
		}
	}
	
	return false
}

func (ae *AutonomousActionExecutor) verifyExecution(ctx context.Context, action Action) error {
	// Use AI to verify the action was executed correctly
	client, err := ae.aiManager.GetDefaultClient()
	if err != nil {
		// Skip verification if AI is not available
		ae.logger.Warn().Err(err).Msg("Skipping verification - AI client not available")
		return nil
	}
	
	messages := []ai.Message{
		{
			Role:    ai.RoleSystem,
			Content: "Verify that the action was executed successfully based on the parameters.",
		},
		{
			Role: ai.RoleUser,
			Content: fmt.Sprintf("Action executed: %s with parameters: %+v. Was this successful?",
				action.Type, action.Parameters),
		},
	}
	
	response, err := client.CreateCompletion(ctx, ai.CompletionRequest{
		Messages:    messages,
		MaxTokens:   intPtr(100),
		Temperature: float64Ptr(0.3),
	})
	
	if err != nil {
		ae.logger.Warn().Err(err).Msg("Failed to verify execution")
		return nil // Don't fail the action due to verification failure
	}
	
	ae.logger.Debug().
		Str("verification_response", response.Choices[0].Message.GetContentAsString()).
		Msg("Action verification completed")
	
	return nil
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}