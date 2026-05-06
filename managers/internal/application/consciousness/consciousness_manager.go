package consciousness

import (
	"context"
	"fmt"
	"time"
	
	"github.com/rs/zerolog"
	"middleman/internal/ai"
	"middleman/internal/ddd"
)

// ConsciousnessManager orchestrates all consciousness components
type ConsciousnessManager struct {
	app              App
	memoryCore       *MemoryCore
	patternDetector  *PatternDetector
	decisionMaker    *DecisionMaker
	actionExecutor   *AutonomousActionExecutor
	learningEngine   *LearningEngine
	aiManager        ai.AIClientManager
	logger           zerolog.Logger
}

// NewConsciousnessManager creates a new consciousness manager using existing components
func NewConsciousnessManager(
	app App,
	memoryCore *MemoryCore,
	patternDetector *PatternDetector,
	learningProcessor *LearningProcessor,
	decisionOrchestrator *DecisionOrchestrator,
	actionExecutor *ActionExecutor,
	aiManager ai.AIClientManager,
	logger zerolog.Logger,
) *ConsciousnessManager {
	cm := &ConsciousnessManager{
		app:             app,
		memoryCore:      memoryCore,
		patternDetector: patternDetector,
		aiManager:       aiManager,
		logger:          logger,
	}
	
	// Create new components that integrate with AI
	cm.decisionMaker = NewDecisionMaker(decisionOrchestrator, aiManager, logger)
	cm.actionExecutor = NewAutonomousActionExecutor(app, actionExecutor, aiManager, logger)
	cm.learningEngine = NewLearningEngine(memoryCore, learningProcessor, logger)
	
	return cm
}

// ProcessEvent processes incoming platform events through the consciousness pipeline
func (cm *ConsciousnessManager) ProcessEvent(ctx context.Context, event ddd.Event) error {
	cm.logger.Info().
		Str("event_type", event.EventName()).
		Str("event_id", event.ID()).
		Msg("ConsciousnessManager processing event")
	
	// 1. Store in memory (already done by ProcessPlatformEvent)
	// Skip to avoid duplicate storage
	
	// 2. Detect patterns
	pattern := cm.patternDetector.DetectPattern(ctx, event)
	if pattern == nil {
		cm.logger.Debug().Msg("No pattern detected")
		return nil // No pattern detected, no action needed
	}
	
	cm.logger.Info().
		Str("pattern_type", pattern.Type).
		Float64("confidence", pattern.Confidence).
		Interface("pattern_data", pattern.Data).
		Msg("Pattern detected")
	
	// 3. Make decision with dynamic tool selection
	decision, err := cm.decisionMaker.MakeDecisionWithDynamicTools(ctx, pattern, event)
	if err != nil {
		cm.logger.Error().Err(err).Msg("Failed to make decision")
		return fmt.Errorf("decision making failed: %w", err)
	}
	
	if decision == nil {
		cm.logger.Debug().Msg("No decision made")
		return nil // No decision made
	}
	
	cm.logger.Info().
		Str("decision_id", decision.ID).
		Str("decision_type", decision.Type).
		Str("priority", decision.Priority).
		Int("action_count", len(decision.Actions)).
		Msg("Decision made")
	
	// 4. Execute action
	if err := cm.actionExecutor.ExecuteDecision(ctx, decision); err != nil {
		cm.logger.Error().
			Err(err).
			Str("decision_id", decision.ID).
			Msg("Failed to execute decision")
		// Record failure but don't return error to avoid blocking event processing
	}
	
	// 5. Learn from outcome
	cm.learningEngine.RecordOutcome(ctx, decision, err)
	
	return nil
}

// GetStatus returns the current status of the consciousness system
func (cm *ConsciousnessManager) GetStatus() ConsciousnessStatus {
	return ConsciousnessStatus{
		Active:          true,
		EventsProcessed: cm.learningEngine.GetProcessedCount(),
		DecisionsMade:   cm.learningEngine.GetDecisionCount(),
		LastActivity:    time.Now(),
		Health:          "healthy",
	}
}