package consciousness

import (
	"time"
)

// Pattern represents a detected pattern
type Pattern struct {
	ID         string
	Type       string
	Confidence float64
	Data       map[string]interface{}
}

// Insight represents a learning insight
type Insight struct {
	ID             string
	Type           string
	Confidence     float64
	Insight        string
	ActionRequired bool
	Action         string
}

// Decision represents a consciousness decision
type Decision struct {
	ID       string
	Type     string
	Priority string
	Actions  []Action
}

// Action represents an action to take
type Action struct {
	Type       string
	Parameters map[string]interface{}
}

// ConsciousnessStatus represents the current state of the store consciousness
type ConsciousnessStatus struct {
	Active          bool
	EventsProcessed int64
	DecisionsMade   int64
	LastActivity    time.Time
	Health          string
}