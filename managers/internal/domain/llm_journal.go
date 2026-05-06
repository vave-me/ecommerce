package domain

import (
	"time"
)

// LLMJournalEntry represents a single LLM response journal entry
type LLMJournalEntry struct {
	ID               string
	ManagerID        string
	ConversationID   *string
	MessageID        *string
	UserID           string
	
	// Request details
	RequestType      string                 // 'chat', 'process_input', 'tool_execution', etc.
	RequestContent   string
	RequestContext   map[string]interface{}
	
	// Response details
	ResponseContent  string
	ResponseMetadata map[string]interface{} // Contains model used, tokens, latency, etc.
	ToolCalls        []ToolCall            // Tool calls made during response
	
	// Learning & patterns
	DetectedPatterns []PatternDetection    // Patterns detected in this interaction
	LearningInsights []LearningInsight     // What the AI learned from this
	ConfidenceScore  float64
	
	// Performance metrics
	ProcessingTimeMs int
	TokensUsed       int
	ModelUsed        string
	Provider         string
	
	// Timestamps
	CreatedAt        time.Time
}

// ToolCall represents a tool invocation during LLM processing
type ToolCall struct {
	ToolName   string                 `json:"tool_name"`
	Operation  string                 `json:"operation"`
	Parameters map[string]interface{} `json:"parameters"`
	Result     interface{}            `json:"result"`
	Success    bool                   `json:"success"`
	ErrorMsg   string                 `json:"error_msg,omitempty"`
	Duration   int                    `json:"duration_ms"`
}

// PatternDetection represents a pattern detected in the interaction
type PatternDetection struct {
	PatternType  string                 `json:"pattern_type"`  // 'behavioral', 'linguistic', 'temporal', etc.
	Description  string                 `json:"description"`
	Confidence   float64                `json:"confidence"`
	Evidence     map[string]interface{} `json:"evidence"`
	Timestamp    time.Time              `json:"timestamp"`
}

// LearningInsight represents what the AI learned from an interaction
type LearningInsight struct {
	InsightType  string                 `json:"insight_type"`  // 'preference', 'pattern', 'optimization', etc.
	Description  string                 `json:"description"`
	UserSpecific bool                   `json:"user_specific"`
	AppliedTo    []string               `json:"applied_to"`    // Future interactions this applies to
	Metadata     map[string]interface{} `json:"metadata"`
}

// LLMJournalRepository provides access to LLM journal entries
type LLMJournalRepository interface {
	// Save stores a new journal entry
	Save(entry *LLMJournalEntry) error
	
	// FindByID retrieves a specific journal entry
	FindByID(id string) (*LLMJournalEntry, error)
	
	// FindByManagerID retrieves all entries for a specific manager
	FindByManagerID(managerID string, limit int, offset int) ([]*LLMJournalEntry, error)
	
	// FindByUserID retrieves all entries for a specific user
	FindByUserID(userID string, limit int, offset int) ([]*LLMJournalEntry, error)
	
	// FindByConversationID retrieves all entries for a specific conversation
	FindByConversationID(conversationID string) ([]*LLMJournalEntry, error)
	
	// FindPatterns retrieves entries with specific pattern types
	FindPatterns(patternType string, since time.Time) ([]*LLMJournalEntry, error)
	
	// GetInsights retrieves learning insights for a user
	GetInsights(userID string, insightType string) ([]LearningInsight, error)
	
	// GetPerformanceMetrics retrieves performance statistics
	GetPerformanceMetrics(managerID string, since time.Time) (map[string]interface{}, error)
}