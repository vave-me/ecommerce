package tools

import (
	"middleman/managers/internal/domain"
	"time"
)

// ToolExecutionContext provides context for tool execution including security
type ToolExecutionContext struct {
	UserID        string
	SessionID     string
	RequestID     string
	AssistantID   string
	AssistantType domain.AssistantType
	Metadata      map[string]interface{}
	Permissions   []string
}

// ToolExecutionResult represents the result of a tool execution from the application layer
type ToolExecutionResult struct {
	ToolName string                 `json:"tool_name"`
	Status   string                 `json:"status"`
	Result   interface{}            `json:"result,omitempty"`
	Error    string                 `json:"error,omitempty"`
	Duration int64                  `json:"duration_ms"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ToolExecutionStream represents a streaming update from tool execution
type ToolExecutionStream struct {
	ID        string                 `json:"id"`
	ToolName  string                 `json:"tool_name"`
	Status    string                 `json:"status"`
	Progress  float64                `json:"progress"`
	Result    interface{}            `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp int64                  `json:"timestamp"`
}

// ToolOperationResult represents the result of a specific tool operation (used in tools package)
type ToolOperationResult struct {
	EntityType string                 `json:"entity_type"`
	Operation  string                 `json:"operation"`
	Success    bool                   `json:"success"`
	Result     interface{}            `json:"result,omitempty"`
	Error      string                 `json:"error,omitempty"`
	Metadata   map[string]interface{} `json:"metadata"`
	Duration   time.Duration          `json:"duration"`
}

// ToolStreamConfig configures streaming behavior for tool operations
type ToolStreamConfig struct {
	BufferSize       int           `json:"buffer_size"`
	ProgressInterval time.Duration `json:"progress_interval"`
	EnableMetrics    bool          `json:"enable_metrics"`
	MaxRetries       int           `json:"max_retries"`
}
