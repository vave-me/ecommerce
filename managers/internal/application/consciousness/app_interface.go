package consciousness

import (
	"context"
	"time"
)

// Command types defined locally to avoid import cycles

// EnsureConsciousnessManager represents the command to ensure consciousness manager exists
type EnsureConsciousnessManager struct {
	ManagerID string `json:"manager_id"`
}

// ProcessUserInput represents the command to process user input
type ProcessUserInput struct {
	ID          string                 `json:"id"`
	ManagerID   string                 `json:"manager_id"`
	UserID      string                 `json:"user_id"`
	Message     string                 `json:"message"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Timestamp   time.Time              `json:"timestamp,omitempty"`
	RequestType string                 `json:"request_type,omitempty"`
}

// ProcessUserInputResult holds the result of processing user input
type ProcessUserInputResult struct {
	ResponseID         string      `json:"response_id"`
	ResponseMessage    string      `json:"response_message"`
	ResponseStatus     string      `json:"response_status"`
	ResponseConfidence float64     `json:"response_confidence"`
	ResponseTimestamp  time.Time   `json:"response_timestamp"`
	ExecutedActions    interface{} `json:"executed_actions,omitempty"` // Using interface{} to avoid more imports
}

// App defines the minimal application interface needed by consciousness components
type App interface {
	// Command execution needed for autonomous actions
	EnsureConsciousnessManager(ctx context.Context, cmd EnsureConsciousnessManager) error
	ProcessUserInput(ctx context.Context, cmd ProcessUserInput) (*ProcessUserInputResult, error)
}