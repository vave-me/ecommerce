package rules

// No imports needed for current implementation

// Action represents a rule-based action
type Action struct {
	Type        string                 `json:"type"`
	Parameters  map[string]interface{} `json:"parameters"`
	Priority    int                    `json:"priority"`
	Description string                 `json:"description"`
}

// Rule represents a decision rule
type Rule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Conditions  []Condition            `json:"conditions"`
	Actions     []Action               `json:"actions"`
	Priority    int                    `json:"priority"`
	Enabled     bool                   `json:"enabled"`
}

// Condition represents a rule condition
type Condition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// RulesConfig holds all rules configuration
type RulesConfig struct {
	Rules   []Rule `json:"rules"`
	Version string `json:"version"`
}

// LoadRules loads the rules configuration
func LoadRules() (*RulesConfig, error) {
	// Default rules configuration
	defaultRules := &RulesConfig{
		Version: "1.0",
		Rules: []Rule{
			{
				ID:      "high_activity_response",
				Name:    "High Activity Response",
				Enabled: true,
				Priority: 1,
				Conditions: []Condition{
					{Field: "activity_spike", Operator: "gt", Value: 0.8},
				},
				Actions: []Action{
					{
						Type:        "scale_resources",
						Priority:    1,
						Description: "Scale up resources for high activity",
						Parameters: map[string]interface{}{
							"scale_factor": 1.5,
						},
					},
				},
			},
			{
				ID:      "fraud_detection_response",
				Name:    "Fraud Detection Response",
				Enabled: true,
				Priority: 1,
				Conditions: []Condition{
					{Field: "fraud_risk", Operator: "gt", Value: 0.7},
				},
				Actions: []Action{
					{
						Type:        "security_alert",
						Priority:    1,
						Description: "Send security alert for fraud risk",
						Parameters: map[string]interface{}{
							"alert_type": "fraud_detected",
						},
					},
				},
			},
		},
	}
	
	return defaultRules, nil
}