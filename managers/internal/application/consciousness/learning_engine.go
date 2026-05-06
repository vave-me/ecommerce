package consciousness

import (
	"context"
	"sync"
	"time"
	
	"github.com/rs/zerolog"
)

// LearningEngine wraps the existing LearningProcessor and adds outcome tracking
type LearningEngine struct {
	processor    *LearningProcessor
	memoryCore   *MemoryCore
	logger       zerolog.Logger
	outcomes     []OutcomeRecord
	outcomeMutex sync.RWMutex
	stats        LearningStats
}

// OutcomeRecord tracks the result of a decision
type OutcomeRecord struct {
	DecisionID string
	Decision   *Decision
	Success    bool
	Error      error
	Timestamp  time.Time
}

// LearningStats tracks overall statistics
type LearningStats struct {
	EventsProcessed int64
	DecisionsMade   int64
	SuccessfulActions int64
	FailedActions    int64
}

// NewLearningEngine creates a new learning engine
func NewLearningEngine(memoryCore *MemoryCore, processor *LearningProcessor, logger zerolog.Logger) *LearningEngine {
	return &LearningEngine{
		processor:  processor,
		memoryCore: memoryCore,
		logger:     logger,
		outcomes:   make([]OutcomeRecord, 0, 1000),
		stats:      LearningStats{},
	}
}

// RecordOutcome records the outcome of a decision
func (le *LearningEngine) RecordOutcome(ctx context.Context, decision *Decision, err error) {
	le.outcomeMutex.Lock()
	defer le.outcomeMutex.Unlock()
	
	outcome := OutcomeRecord{
		DecisionID: decision.ID,
		Decision:   decision,
		Success:    err == nil,
		Error:      err,
		Timestamp:  time.Now(),
	}
	
	le.outcomes = append(le.outcomes, outcome)
	le.stats.DecisionsMade++
	
	if err == nil {
		le.stats.SuccessfulActions += int64(len(decision.Actions))
	} else {
		le.stats.FailedActions += int64(len(decision.Actions))
	}
	
	// Keep only last 1000 outcomes to prevent memory growth
	if len(le.outcomes) > 1000 {
		le.outcomes = le.outcomes[len(le.outcomes)-1000:]
	}
	
	// Use existing processor to learn from outcome
	insight := Insight{
		ID:             generateID(),
		Type:           "decision_outcome",
		Confidence:     le.calculateConfidence(outcome),
		Insight:        le.generateInsight(outcome),
		ActionRequired: err != nil,
		Action:         le.suggestCorrectiveAction(outcome),
	}
	
	le.processor.Learn(ctx, insight)
	
	// Log learning progress
	if le.stats.DecisionsMade%10 == 0 {
		le.logger.Info().
			Int64("decisions_made", le.stats.DecisionsMade).
			Float64("success_rate", le.getSuccessRate()).
			Msg("Learning progress")
	}
}

// GetProcessedCount returns the number of events processed
func (le *LearningEngine) GetProcessedCount() int64 {
	return le.stats.EventsProcessed
}

// GetDecisionCount returns the number of decisions made
func (le *LearningEngine) GetDecisionCount() int64 {
	return le.stats.DecisionsMade
}

// GetSuccessRate calculates the overall success rate
func (le *LearningEngine) getSuccessRate() float64 {
	total := le.stats.SuccessfulActions + le.stats.FailedActions
	if total == 0 {
		return 0.0
	}
	return float64(le.stats.SuccessfulActions) / float64(total)
}

// GetSuccessRateForType calculates success rate for a specific decision type
func (le *LearningEngine) GetSuccessRateForType(decisionType string) float64 {
	le.outcomeMutex.RLock()
	defer le.outcomeMutex.RUnlock()
	
	total := 0
	successful := 0
	
	for _, outcome := range le.outcomes {
		if outcome.Decision.Type == decisionType {
			total++
			if outcome.Success {
				successful++
			}
		}
	}
	
	if total == 0 {
		return 0.0
	}
	
	return float64(successful) / float64(total)
}

// GetRecentOutcomes returns the most recent outcomes
func (le *LearningEngine) GetRecentOutcomes(limit int) []OutcomeRecord {
	le.outcomeMutex.RLock()
	defer le.outcomeMutex.RUnlock()
	
	if limit > len(le.outcomes) {
		limit = len(le.outcomes)
	}
	
	result := make([]OutcomeRecord, limit)
	copy(result, le.outcomes[len(le.outcomes)-limit:])
	return result
}

func (le *LearningEngine) calculateConfidence(outcome OutcomeRecord) float64 {
	if outcome.Success {
		return 0.9
	}
	return 0.3
}

func (le *LearningEngine) generateInsight(outcome OutcomeRecord) string {
	if outcome.Success {
		return "Decision executed successfully"
	}
	if outcome.Error != nil {
		return "Decision failed: " + outcome.Error.Error()
	}
	return "Decision outcome unknown"
}

func (le *LearningEngine) suggestCorrectiveAction(outcome OutcomeRecord) string {
	if outcome.Success {
		return ""
	}
	
	// Suggest corrective actions based on error patterns
	if outcome.Error != nil {
		errorMsg := outcome.Error.Error()
		switch {
		case contains(errorMsg, "timeout"):
			return "increase_timeout"
		case contains(errorMsg, "permission"):
			return "check_permissions"
		case contains(errorMsg, "not found"):
			return "verify_resource_exists"
		default:
			return "review_action_parameters"
		}
	}
	
	return "investigate_failure"
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr)))
}