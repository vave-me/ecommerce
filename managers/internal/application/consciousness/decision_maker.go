package consciousness

import (
	"context"
	"encoding/json"
	"fmt"
	
	"github.com/rs/zerolog"
	"middleman/internal/ai"
	"middleman/internal/ddd"
	"middleman/managers/internal/application/consciousness/rules"
)

// DecisionMaker integrates rule-based and AI-powered decision making
type DecisionMaker struct {
	orchestrator *DecisionOrchestrator // Existing orchestrator
	aiManager    ai.AIClientManager
	toolSelector *ToolSelector
	logger       zerolog.Logger
	rules        []DecisionRule
}

// DecisionRule defines a rule for automated decision making
type DecisionRule struct {
	Name         string
	PatternType  string
	Confidence   float64
	Action       string
	Priority     int
	Parameters   map[string]interface{}
}

// NewDecisionMaker creates a new decision maker that wraps the existing orchestrator
func NewDecisionMaker(orchestrator *DecisionOrchestrator, aiManager ai.AIClientManager, logger zerolog.Logger) *DecisionMaker {
	return &DecisionMaker{
		orchestrator: orchestrator,
		aiManager:    aiManager,
		toolSelector: NewToolSelector(logger),
		logger:       logger,
		rules:        loadDecisionRules(),
	}
}

// MakeDecision makes a decision based on pattern using rules and AI
func (dm *DecisionMaker) MakeDecision(ctx context.Context, pattern *Pattern) (*Decision, error) {
	// First try rule-based decision
	if decision := dm.checkRules(pattern); decision != nil {
		dm.logger.Info().
			Str("decision_type", "rule_based").
			Str("pattern_type", pattern.Type).
			Msg("Decision made using rules")
		return decision, nil
	}
	
	// If no rule matches and confidence is high enough, use AI
	if pattern.Confidence >= 0.7 {
		dm.logger.Info().
			Str("pattern_type", pattern.Type).
			Float64("confidence", pattern.Confidence).
			Msg("Using AI for decision making")
		return dm.analyzeWithAI(ctx, pattern)
	}
	
	// Use existing orchestrator for other cases
	existingDecision := dm.orchestrator.Orchestrate(ctx, []Pattern{*pattern})
	if existingDecision != nil {
		return existingDecision, nil
	}
	
	return nil, nil
}

func (dm *DecisionMaker) checkRules(pattern *Pattern) *Decision {
	for _, rule := range dm.rules {
		if rule.PatternType == pattern.Type && pattern.Confidence >= rule.Confidence {
			dm.logger.Debug().
				Str("rule_name", rule.Name).
				Str("action", rule.Action).
				Msg("Rule matched")
				
			// Merge rule parameters with pattern data
			params := make(map[string]interface{})
			for k, v := range pattern.Data {
				params[k] = v
			}
			for k, v := range rule.Parameters {
				params[k] = v
			}
			
			return &Decision{
				ID:       generateID(),
				Type:     rule.Action,
				Priority: fmt.Sprintf("%d", rule.Priority),
				Actions: []Action{
					{
						Type:       rule.Action,
						Parameters: params,
					},
				},
			}
		}
	}
	return nil
}

func (dm *DecisionMaker) analyzeWithAI(ctx context.Context, pattern *Pattern) (*Decision, error) {
	client, err := dm.aiManager.GetDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get AI client: %w", err)
	}
	
	messages := []ai.Message{
		{
			Role: ai.RoleSystem,
			Content: `You are an autonomous e-commerce platform manager. Analyze patterns and suggest actions using available tools.
Your role is to:
1. Detect issues and opportunities
2. Take preventive actions
3. Optimize platform performance
4. Ensure customer satisfaction

Always respond with specific tool calls when action is needed.`,
		},
		{
			Role: ai.RoleUser,
			Content: fmt.Sprintf(`Pattern detected: %s
Confidence: %.2f
Data: %+v

Based on this pattern, what action should be taken? Consider:
- Is immediate action needed?
- What's the potential impact?
- Which tool would be most appropriate?`,
				pattern.Type, pattern.Confidence, pattern.Data),
		},
	}
	
	tools := dm.getAvailableTools()
	
	response, err := client.ExecuteWithTools(ctx, messages, tools)
	if err != nil {
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}
	
	decision, err := dm.parseAIResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}
	
	return decision, nil
}

func (dm *DecisionMaker) getAvailableTools() []ai.ToolDefinition {
	return []ai.ToolDefinition{
		{
			Type: ai.ToolTypeFunction,
			Function: ai.FunctionDef{
				Name:        "send_notification",
				Description: "Send a notification to users",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id":  map[string]interface{}{"type": "string", "description": "User ID to notify"},
						"template": map[string]interface{}{"type": "string", "description": "Notification template name"},
						"data":     map[string]interface{}{"type": "object", "description": "Template data"},
					},
					"required": []string{"user_id", "template"},
				},
			},
		},
		{
			Type: ai.ToolTypeFunction,
			Function: ai.FunctionDef{
				Name:        "create_offer",
				Description: "Create a promotional offer",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id":    map[string]interface{}{"type": "string", "description": "Target user ID"},
						"product_id": map[string]interface{}{"type": "string", "description": "Product ID"},
						"discount":   map[string]interface{}{"type": "number", "description": "Discount percentage"},
						"duration":   map[string]interface{}{"type": "string", "description": "Offer duration (e.g., '24h')"},
					},
					"required": []string{"user_id", "discount"},
				},
			},
		},
		{
			Type: ai.ToolTypeFunction,
			Function: ai.FunctionDef{
				Name:        "escalate_ticket",
				Description: "Escalate a support ticket",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"ticket_id": map[string]interface{}{"type": "string", "description": "Ticket ID"},
						"reason":    map[string]interface{}{"type": "string", "description": "Escalation reason"},
						"priority":  map[string]interface{}{"type": "string", "enum": []string{"low", "medium", "high", "urgent"}},
					},
					"required": []string{"ticket_id", "reason"},
				},
			},
		},
		{
			Type: ai.ToolTypeFunction,
			Function: ai.FunctionDef{
				Name:        "flag_order",
				Description: "Flag an order for review",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"order_id": map[string]interface{}{"type": "string", "description": "Order ID"},
						"reason":   map[string]interface{}{"type": "string", "description": "Flagging reason"},
						"risk_level": map[string]interface{}{"type": "string", "enum": []string{"low", "medium", "high"}},
					},
					"required": []string{"order_id", "reason"},
				},
			},
		},
	}
}

func (dm *DecisionMaker) parseAIResponse(response *ai.CompletionResponse) (*Decision, error) {
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}
	
	choice := response.Choices[0]
	if len(choice.Message.ToolCalls) == 0 {
		// AI decided no action needed
		dm.logger.Info().Msg("AI determined no action needed")
		return nil, nil
	}
	
	actions := make([]Action, 0, len(choice.Message.ToolCalls))
	for _, toolCall := range choice.Message.ToolCalls {
		var params map[string]interface{}
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
			dm.logger.Error().
				Err(err).
				Str("tool_name", toolCall.Function.Name).
				Msg("Failed to parse tool arguments")
			continue
		}
		
		actions = append(actions, Action{
			Type:       toolCall.Function.Name,
			Parameters: params,
		})
	}
	
	if len(actions) == 0 {
		return nil, fmt.Errorf("no valid actions parsed from AI response")
	}
	
	return &Decision{
		ID:       generateID(),
		Type:     "ai_suggested",
		Priority: "medium", // Default priority for AI decisions
		Actions:  actions,
	}, nil
}

func loadDecisionRules() []DecisionRule {
	// Try to load from configuration file
	rulesConfig, err := rules.LoadRules()
	if err != nil {
		// Log error and return default rules
		zerolog.Ctx(context.Background()).Error().Err(err).Msg("Failed to load rules from config, using defaults")
		return getDefaultRules()
	}
	
	// Convert loaded rules to DecisionRule format
	decisionRules := make([]DecisionRule, 0, len(rulesConfig.DecisionRules))
	for _, rule := range rulesConfig.DecisionRules {
		if len(rule.Actions) > 0 {
			// Use first action as primary action
			action := rule.Actions[0]
			decisionRules = append(decisionRules, DecisionRule{
				Name:        rule.Name,
				PatternType: rule.PatternType,
				Confidence:  rule.ConfidenceThreshold,
				Action:      action.Type,
				Priority:    determinePriority(action),
				Parameters:  action.Parameters,
			})
		}
	}
	
	return decisionRules
}

func getDefaultRules() []DecisionRule {
	return []DecisionRule{
		{
			Name:        "abandoned_cart_recovery",
			PatternType: "cart_abandonment",
			Confidence:  0.7,
			Action:      "send_notification",
			Priority:    2,
			Parameters: map[string]interface{}{
				"template": "abandoned_cart_reminder",
			},
		},
		{
			Name:        "fraud_alert_high_risk",
			PatternType: "fraud_risk",
			Confidence:  0.8,
			Action:      "flag_order",
			Priority:    1,
			Parameters: map[string]interface{}{
				"risk_level": "high",
				"reason":     "Automated fraud detection",
			},
		},
		{
			Name:        "support_crisis_escalation",
			PatternType: "support_crisis",
			Confidence:  0.75,
			Action:      "escalate_ticket",
			Priority:    1,
			Parameters: map[string]interface{}{
				"priority": "urgent",
				"reason":   "Crisis pattern detected",
			},
		},
		{
			Name:        "user_surge_notification",
			PatternType: "user_surge",
			Confidence:  0.8,
			Action:      "send_notification",
			Priority:    3,
			Parameters: map[string]interface{}{
				"template": "platform_busy_notification",
			},
		},
	}
}

func determinePriority(action rules.Action) int {
	// Immediate actions have highest priority
	if action.Immediate {
		return 1
	}
	
	// Actions with delay have lower priority
	if action.Delay != "" {
		return 3
	}
	
	// Default medium priority
	return 2
}

// MakeDecisionWithDynamicTools makes a decision with dynamically selected tools based on pattern
func (dm *DecisionMaker) MakeDecisionWithDynamicTools(ctx context.Context, pattern *Pattern, event interface{}) (*Decision, error) {
	// Get relevant tools for this pattern
	var relevantTools []string
	
	// If we have an event, use it for tool selection
	if dddEvent, ok := event.(ddd.Event); ok {
		relevantTools = dm.toolSelector.SelectToolsForEvent(ctx, dddEvent)
	} else {
		// Otherwise use pattern-based selection
		relevantTools = dm.toolSelector.SelectToolsForPattern(pattern)
	}
	
	dm.logger.Info().
		Str("pattern_type", pattern.Type).
		Int("relevant_tools", len(relevantTools)).
		Strs("tools", relevantTools).
		Msg("Selected relevant tools for decision making")
	
	// First try rule-based decision with dynamic tools
	if decision := dm.checkRulesWithTools(pattern, relevantTools); decision != nil {
		dm.logger.Info().
			Str("decision_type", "rule_based_dynamic").
			Str("pattern_type", pattern.Type).
			Msg("Decision made using rules with dynamic tools")
		return decision, nil
	}
	
	// If no rule matches and confidence is high enough, use AI with dynamic tools
	if pattern.Confidence >= 0.7 {
		dm.logger.Info().
			Str("pattern_type", pattern.Type).
			Float64("confidence", pattern.Confidence).
			Int("available_tools", len(relevantTools)).
			Msg("Using AI for decision making with dynamic tools")
		return dm.analyzeWithAIDynamicTools(ctx, pattern, relevantTools)
	}
	
	return nil, nil
}

// checkRulesWithTools checks rules but considers available dynamic tools
func (dm *DecisionMaker) checkRulesWithTools(pattern *Pattern, availableTools []string) *Decision {
	for _, rule := range dm.rules {
		if rule.PatternType == pattern.Type && pattern.Confidence >= rule.Confidence {
			// Check if the rule's action can be executed with available tools
			if dm.canExecuteWithTools(rule.Action, availableTools) {
				dm.logger.Debug().
					Str("rule_name", rule.Name).
					Str("action", rule.Action).
					Msg("Rule matched with available tools")
					
				// Merge rule parameters with pattern data
				params := make(map[string]interface{})
				for k, v := range pattern.Data {
					params[k] = v
				}
				for k, v := range rule.Parameters {
					params[k] = v
				}
				
				// Add dynamic tool selection
				params["_selected_tools"] = availableTools
				
				return &Decision{
					ID:       generateID(),
					Type:     rule.Action,
					Priority: fmt.Sprintf("%d", rule.Priority),
					Actions: []Action{
						{
							Type:       rule.Action,
							Parameters: params,
						},
					},
				}
			}
		}
	}
	return nil
}

// canExecuteWithTools checks if an action can be executed with the available tools
func (dm *DecisionMaker) canExecuteWithTools(action string, availableTools []string) bool {
	// Map actions to required tools
	actionToolMap := map[string][]string{
		"send_notification": {"notification_send", "notification_create_template"},
		"create_offer":      {"offer_create", "offer_assign_to_user"},
		"escalate_ticket":   {"support_escalate_ticket", "support_get_ticket"},
		"flag_order":        {"order_flag_for_review", "order_get_by_id"},
	}
	
	requiredTools, exists := actionToolMap[action]
	if !exists {
		// Unknown action, allow it
		return true
	}
	
	// Check if at least one required tool is available
	for _, required := range requiredTools {
		for _, available := range availableTools {
			if required == available {
				return true
			}
		}
	}
	
	return false
}

// analyzeWithAIDynamicTools uses AI with dynamically selected tools
func (dm *DecisionMaker) analyzeWithAIDynamicTools(ctx context.Context, pattern *Pattern, relevantTools []string) (*Decision, error) {
	client, err := dm.aiManager.GetDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get AI client: %w", err)
	}
	
	messages := []ai.Message{
		{
			Role: ai.RoleSystem,
			Content: fmt.Sprintf(`You are an autonomous e-commerce platform manager with access to %d specialized tools.
Your role is to:
1. Detect issues and opportunities
2. Take preventive actions using the most appropriate tools
3. Optimize platform performance
4. Ensure customer satisfaction

Available tools for this situation: %v

Always respond with specific tool calls when action is needed. Choose the most appropriate tool from the available options.`, len(relevantTools), relevantTools),
		},
		{
			Role: ai.RoleUser,
			Content: fmt.Sprintf(`Pattern detected: %s
Confidence: %.2f
Data: %+v

Based on this pattern and the available tools, what action should be taken? Consider:
- Is immediate action needed?
- What's the potential impact?
- Which of the available tools would be most appropriate?`,
				pattern.Type, pattern.Confidence, pattern.Data),
		},
	}
	
	// Get tool definitions only for relevant tools
	tools := dm.getDynamicToolDefinitions(relevantTools)
	
	response, err := client.ExecuteWithTools(ctx, messages, tools)
	if err != nil {
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}
	
	decision, err := dm.parseAIResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}
	
	// Add selected tools to decision metadata
	if decision != nil {
		decision.Metadata = map[string]interface{}{
			"selected_tools": relevantTools,
			"tool_count":     len(relevantTools),
		}
	}
	
	return decision, nil
}

// getDynamicToolDefinitions creates tool definitions for the selected tools
func (dm *DecisionMaker) getDynamicToolDefinitions(toolNames []string) []ai.ToolDefinition {
	// This would ideally pull from a comprehensive tool definition registry
	// For now, we'll create definitions for the most common tools
	toolDefs := []ai.ToolDefinition{}
	
	for _, toolName := range toolNames {
		switch toolName {
		case "notification_send":
			toolDefs = append(toolDefs, ai.ToolDefinition{
				Type: ai.ToolTypeFunction,
				Function: ai.FunctionDef{
					Name:        "notification_send",
					Description: "Send a notification to users",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"user_id":  map[string]interface{}{"type": "string", "description": "User ID to notify"},
							"template": map[string]interface{}{"type": "string", "description": "Notification template name"},
							"data":     map[string]interface{}{"type": "object", "description": "Template data"},
						},
						"required": []string{"user_id", "template"},
					},
				},
			})
		case "offer_create":
			toolDefs = append(toolDefs, ai.ToolDefinition{
				Type: ai.ToolTypeFunction,
				Function: ai.FunctionDef{
					Name:        "offer_create",
					Description: "Create a promotional offer",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"user_id":    map[string]interface{}{"type": "string", "description": "Target user ID"},
							"product_id": map[string]interface{}{"type": "string", "description": "Product ID"},
							"discount":   map[string]interface{}{"type": "number", "description": "Discount percentage"},
							"duration":   map[string]interface{}{"type": "string", "description": "Offer duration"},
						},
						"required": []string{"user_id", "discount"},
					},
				},
			})
		case "order_flag_for_review":
			toolDefs = append(toolDefs, ai.ToolDefinition{
				Type: ai.ToolTypeFunction,
				Function: ai.FunctionDef{
					Name:        "order_flag_for_review",
					Description: "Flag an order for manual review",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"order_id":   map[string]interface{}{"type": "string", "description": "Order ID"},
							"reason":     map[string]interface{}{"type": "string", "description": "Flagging reason"},
							"risk_level": map[string]interface{}{"type": "string", "enum": []string{"low", "medium", "high"}},
						},
						"required": []string{"order_id", "reason"},
					},
				},
			})
		// Add more tool definitions as needed
		}
	}
	
	return toolDefs
}