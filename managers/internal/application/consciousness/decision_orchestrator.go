package consciousness

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"middleman/internal/ddd"
	"middleman/managers/internal/domain"
)

type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
	PriorityUrgent
)

type DecisionOrchestrator struct {
	schedulerRepo  domain.SchedulerRepository
	actionExecutor *ActionExecutor
	executionQueue chan Decision
	workers        int
	logger         zerolog.Logger
}

func NewDecisionOrchestrator(params ...interface{}) *DecisionOrchestrator {
	orchestrator := &DecisionOrchestrator{
		executionQueue: make(chan Decision, 1000),
		workers:        10,
		logger:         zerolog.Logger{},
	}
	
	// Parse optional parameters
	for _, param := range params {
		switch v := param.(type) {
		case domain.SchedulerRepository:
			orchestrator.schedulerRepo = v
		case zerolog.Logger:
			orchestrator.logger = v
		case *ActionExecutor:
			orchestrator.actionExecutor = v
		}
	}
	
	// Create action executor if not provided
	if orchestrator.actionExecutor == nil && orchestrator.schedulerRepo != nil {
		orchestrator.actionExecutor = NewActionExecutor(orchestrator.schedulerRepo, orchestrator.logger)
	}

	// Start workers only if we have an action executor
	if orchestrator.actionExecutor != nil {
		for i := 0; i < orchestrator.workers; i++ {
			go orchestrator.worker(i)
		}
	}

	return orchestrator
}

func (o *DecisionOrchestrator) MakeDecision(ctx context.Context, pattern *Pattern, insight *Insight) *Decision {
	if pattern == nil {
		return nil
	}
	
	decision := &Decision{
		ID:   uuid.New().String(),
		Type: pattern.Type,
	}
	
	// Determine priority based on pattern confidence and type
	switch pattern.Type {
	case PatternTypeFraudRisk, PatternTypeSupportCrisis:
		decision.Priority = "critical"
	case PatternTypeAbandonmentRisk, PatternTypeCancellationSpike:
		decision.Priority = "high"
	case PatternTypeActivitySurge, PatternTypeUserSurge:
		if pattern.Confidence > 0.8 {
			decision.Priority = "high"
		} else {
			decision.Priority = "medium"
		}
	default:
		decision.Priority = "low"
	}
	
	// Add actions based on insight
	if insight != nil && insight.ActionRequired {
		switch insight.Action {
		case "scale_infrastructure":
			decision.Actions = append(decision.Actions, Action{
				Type: "schedule_task",
				Parameters: map[string]interface{}{
					"task_type": "scale_resources",
					"run_at":    time.Now().Add(5 * time.Minute).Format(time.RFC3339),
					"data": map[string]interface{}{
						"surge_factor": pattern.Data["surge_factor"],
					},
				},
			})
			
		case "send_abandonment_recovery":
			decision.Actions = append(decision.Actions, Action{
				Type: "send_notification",
				Parameters: map[string]interface{}{
					"type":    "cart_recovery",
					"message": "Complete your purchase with a special discount!",
				},
			})
			
		case "alert_security":
			decision.Actions = append(decision.Actions, Action{
				Type: "alert_support",
				Parameters: map[string]interface{}{
					"alert_type": "security",
					"priority":   "critical",
					"details":    insight.Insight,
				},
			})
			
		case "alert_support_team":
			decision.Actions = append(decision.Actions, Action{
				Type: "alert_support",
				Parameters: map[string]interface{}{
					"alert_type": "volume_spike",
					"priority":   "high",
					"details":    insight.Insight,
				},
			})
			
		case "prepare_onboarding_resources":
			decision.Actions = append(decision.Actions, Action{
				Type: "optimize_inventory",
				Parameters: map[string]interface{}{
					"optimization": "prepare_for_new_users",
				},
			})
			
		case "investigate_cancellation_reasons":
			decision.Actions = append(decision.Actions, Action{
				Type: "analyze_trend",
				Parameters: map[string]interface{}{
					"trend_type": "cancellations",
					"data":       pattern.Data,
				},
			})
		}
	}
	
	// Only return decision if there are actions to take
	if len(decision.Actions) == 0 {
		return nil
	}
	
	return decision
}

func (o *DecisionOrchestrator) MakeDecisions(
	ctx context.Context,
	event ddd.Event,
	patterns []Pattern,
	insights []Insight,
) []Decision {
	decisions := []Decision{}

	// Event-driven decisions
	switch event.EventName() {
	case "UserCreated":
		decisions = append(decisions, o.onNewCustomerDecision(event, patterns))
	case "ProductAdded":
		decisions = append(decisions, o.onNewProductDecision(event, patterns))
	case "OrderCompleted":
		decisions = append(decisions, o.onOrderCompletedDecision(event, insights))
	case "OrderCanceled":
		decisions = append(decisions, o.onOrderCanceledDecision(event, patterns))
	case "ReviewAdded":
		decisions = append(decisions, o.onReviewDecision(event))
	case "TicketCreated":
		decisions = append(decisions, o.onSupportTicketDecision(event, patterns))
	case "BasketItemAdded":
		decisions = append(decisions, o.onBasketActivityDecision(event, patterns))
	}

	// Pattern-driven decisions
	for _, pattern := range patterns {
		if pattern.Confidence > 0.8 {
			if decision := o.createPatternDecision(pattern); decision != nil {
				decisions = append(decisions, *decision)
			}
		}
	}

	// Learning-driven decisions
	for _, insight := range insights {
		if insight.Confidence > 0.75 && insight.ActionRequired {
			if decision := o.createLearningDecision(insight); decision != nil {
				decisions = append(decisions, *decision)
			}
		}
	}

	return o.prioritizeDecisions(decisions)
}

func (o *DecisionOrchestrator) ExecuteDecision(ctx context.Context, decision Decision) error {
	select {
	case o.executionQueue <- decision:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		o.logger.Error().
			Str("decision_id", decision.ID).
			Msg("Execution queue full, decision dropped")
		return fmt.Errorf("execution queue full")
	}
}

func (o *DecisionOrchestrator) worker(id int) {
	o.logger.Info().Int("worker_id", id).Msg("Decision worker started")

	for decision := range o.executionQueue {
		ctx := context.Background()

		for _, action := range decision.Actions {
			if err := o.executeAction(ctx, decision, action); err != nil {
				o.logger.Error().
					Err(err).
					Str("decision_id", decision.ID).
					Str("action_type", action.Type).
					Msg("Action execution failed")
			}
		}
	}
}

func (o *DecisionOrchestrator) executeAction(
	ctx context.Context,
	decision Decision,
	action Action,
) error {
	actionCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	task := o.buildActionTask(decision, action)
	actionID := fmt.Sprintf("%s_%s_%s", decision.ID, action.Type, uuid.New().String())

	// Create scheduled action using the repository
	scheduledAction := &domain.ScheduledAction{
		ID:          actionID,
		Name:        fmt.Sprintf("Consciousness Decision: %s", action.Type),
		Description: task,
		EntityID:    "store_consciousness",
		EntityType:  "consciousness",
		Action:      action.Type,
		Parameters:  action.Parameters,
		ScheduledAt: time.Now(), // Execute immediately
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	var err error
	for attempt := 0; attempt < 3; attempt++ {
		err = o.schedulerRepo.ScheduleAction(actionCtx, scheduledAction)
		if err == nil {
			break
		}

		select {
		case <-time.After(time.Duration(attempt+1) * time.Second):
			continue
		case <-actionCtx.Done():
			return fmt.Errorf("context cancelled during retry: %w", actionCtx.Err())
		}
	}

	if err != nil {
		return fmt.Errorf("failed to schedule action after 3 attempts: %w", err)
	}

	o.logger.Info().
		Str("action_id", actionID).
		Str("decision_id", decision.ID).
		Str("action_type", action.Type).
		Msg("Action scheduled successfully")

	return nil
}

func (o *DecisionOrchestrator) buildActionTask(decision Decision, action Action) string {
	contextJSON, _ := json.Marshal(map[string]interface{}{
		"decision_type": decision.Type,
		"priority":      decision.Priority,
	})
	parametersJSON, _ := json.Marshal(action.Parameters)

	return fmt.Sprintf(`STORE CONSCIOUSNESS ACTION

Decision Type: %s
Priority: %s
Action: %s

Context:
%s

Parameters:
%s

Execute this action to optimize platform operations.`,
		decision.Type,
		decision.Priority,
		action.Type,
		string(contextJSON),
		string(parametersJSON))
}

// Decision creation methods

func (o *DecisionOrchestrator) onNewCustomerDecision(event ddd.Event, patterns []Pattern) Decision {
	priority := "high"
	for _, pattern := range patterns {
		if pattern.Type == PatternTypeUserSurge {
			priority = "medium"
			break
		}
	}

	return Decision{
		ID:       uuid.New().String(),
		Type:     "customer_onboarding",
		Priority: priority,
		Actions: []Action{
			{
				Type: "welcome_sequence",
				Parameters: map[string]interface{}{
					"template":                "new_user_welcome_2024",
					"include_recommendations": true,
				},
			},
			{
				Type: "engagement_monitoring",
				Parameters: map[string]interface{}{
					"metrics":             []string{"first_search", "first_view"},
					"alert_on_inactivity": true,
				},
			},
		},
	}
}

func (o *DecisionOrchestrator) onNewProductDecision(event ddd.Event, patterns []Pattern) Decision {
	return Decision{
		ID:       uuid.New().String(),
		Type:     "product_launch_optimization",
		Priority: "medium",
		Actions: []Action{
			{
				Type: "market_analysis",
				Parameters: map[string]interface{}{
					"compare_similar": true,
					"price_analysis":  true,
				},
			},
			{
				Type: "visibility_optimization",
				Parameters: map[string]interface{}{
					"boost_if_high_potential": true,
					"notify_interested_users": true,
				},
			},
		},
	}
}

func (o *DecisionOrchestrator) onOrderCompletedDecision(event ddd.Event, insights []Insight) Decision {
	return Decision{
		ID:       uuid.New().String(),
		Type:     "post_purchase_optimization",
		Priority: "high",
		Actions: []Action{
			{
				Type: "transaction_celebration",
				Parameters: map[string]interface{}{
					"include_invoice": true,
					"tracking_info":   true,
				},
			},
			{
				Type: "satisfaction_monitoring",
				Parameters: map[string]interface{}{
					"request_feedback":   true,
					"incentivize_review": true,
				},
			},
		},
	}
}

func (o *DecisionOrchestrator) onOrderCanceledDecision(event ddd.Event, patterns []Pattern) Decision {
	urgency := "medium"
	for _, pattern := range patterns {
		if pattern.Type == PatternTypeCancellationSpike {
			urgency = "urgent"
			break
		}
	}

	return Decision{
		ID:       uuid.New().String(),
		Type:     "cancellation_recovery",
		Priority: urgency,
		Actions: []Action{
			{
				Type: "root_cause_analysis",
				Parameters: map[string]interface{}{
					"check_payment_issues":    true,
					"check_shipping_concerns": true,
				},
			},
			{
				Type: "recovery_attempt",
				Parameters: map[string]interface{}{
					"offer_discount":       true,
					"alternative_products": true,
				},
			},
		},
	}
}

func (o *DecisionOrchestrator) onReviewDecision(event ddd.Event) Decision {
	// Simplified - would extract rating from payload
	return Decision{
		ID:       uuid.New().String(),
		Type:     "review_response",
		Priority: "medium",
		Actions: []Action{
			{
				Type: "positive_amplification",
				Parameters: map[string]interface{}{
					"thank_reviewer":    true,
					"share_with_seller": true,
				},
			},
		},
	}
}

func (o *DecisionOrchestrator) onSupportTicketDecision(event ddd.Event, patterns []Pattern) Decision {
	priority := "high"
	for _, pattern := range patterns {
		if pattern.Type == PatternTypeSupportCrisis {
			priority = "urgent"
			break
		}
	}

	return Decision{
		ID:       uuid.New().String(),
		Type:     "support_excellence",
		Priority: priority,
		Actions: []Action{
			{
				Type: "instant_acknowledgment",
				Parameters: map[string]interface{}{
					"personalized":     true,
					"set_expectations": true,
				},
			},
			{
				Type: "intelligent_routing",
				Parameters: map[string]interface{}{
					"check_customer_history": true,
					"match_expertise":        true,
				},
			},
		},
	}
}

func (o *DecisionOrchestrator) onBasketActivityDecision(event ddd.Event, patterns []Pattern) Decision {
	abandonmentRisk := false
	for _, pattern := range patterns {
		if pattern.Type == PatternTypeAbandonmentRisk {
			abandonmentRisk = true
			break
		}
	}

	actions := []Action{
		{
			Type: "conversion_optimization",
			Parameters: map[string]interface{}{
				"show_stock_levels":     true,
				"suggest_complementary": true,
			},
		},
	}

	if abandonmentRisk {
		actions = append(actions, Action{
			Type: "abandonment_prevention",
			Parameters: map[string]interface{}{
				"send_reminder":    true,
				"offer_assistance": true,
			},
		})
	}

	return Decision{
		ID:       uuid.New().String(),
		Type:     "conversion_assistance",
		Priority: "medium",
		Actions:  actions,
	}
}

func (o *DecisionOrchestrator) createPatternDecision(pattern Pattern) *Decision {
	switch pattern.Type {
	case PatternTypeActivitySurge:
		return &Decision{
			ID:       uuid.New().String(),
			Type:     "surge_management",
			Priority: "urgent",
			Actions: []Action{
				{
					Type: "resource_scaling",
					Parameters: map[string]interface{}{
						"scale_factor": 1.5,
						"alert_team":   true,
					},
				},
			},
		}

	case PatternTypeFraudRisk:
		return &Decision{
			ID:       uuid.New().String(),
			Type:     "fraud_prevention",
			Priority: "urgent",
			Actions: []Action{
				{
					Type: "security_response",
					Parameters: map[string]interface{}{
						"increase_verification": true,
						"flag_transactions":     true,
					},
				},
			},
		}
	}

	return nil
}

func (o *DecisionOrchestrator) createLearningDecision(insight Insight) *Decision {
	if !insight.ActionRequired {
		return nil
	}

	return &Decision{
		ID:       uuid.New().String(),
		Type:     fmt.Sprintf("learning_%s", insight.Type),
		Priority: "medium",
		Actions: []Action{
			{
				Type:       "apply_learning",
				Parameters: map[string]interface{}{},
			},
		},
	}
}

func (o *DecisionOrchestrator) prioritizeDecisions(decisions []Decision) []Decision {
	// Sort by priority
	sort.Slice(decisions, func(i, j int) bool {
		return decisions[i].Priority > decisions[j].Priority
	})

	// Remove duplicates by type
	seen := make(map[string]bool)
	filtered := []Decision{}

	for _, decision := range decisions {
		if !seen[decision.Type] {
			seen[decision.Type] = true
			filtered = append(filtered, decision)
		}
	}

	// Limit to top 10
	if len(filtered) > 10 {
		filtered = filtered[:10]
	}

	return filtered
}
