package toolmappings

// No imports needed for current implementation

// EventToolMapping represents the mapping between events and tools
type EventToolMapping struct {
	Mappings map[string]ToolMapping `json:"mappings"`
	Version  string                 `json:"version"`
}

// ToolMapping defines tools for specific event patterns
type ToolMapping struct {
	Tools       []string            `json:"tools"`
	Priority    int                 `json:"priority"`
	Conditions  []MappingCondition  `json:"conditions"`
	Description string              `json:"description"`
}

// MappingCondition defines when a mapping should be applied
type MappingCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// LoadEventToolMappings loads the event-to-tool mappings
func LoadEventToolMappings() (*EventToolMapping, error) {
	// Default mappings for common event patterns
	defaultMappings := &EventToolMapping{
		Version: "1.0",
		Mappings: map[string]ToolMapping{
			"user_activity_spike": {
				Tools:       []string{"user_get_base", "activity_get_system", "metric_record"},
				Priority:    1,
				Description: "Tools for handling user activity spikes",
				Conditions: []MappingCondition{
					{Field: "activity_level", Operator: "gt", Value: 0.8},
				},
			},
			"fraud_detection": {
				Tools:       []string{"user_find", "activity_log", "notification_create", "support_create_ticket"},
				Priority:    1,
				Description: "Tools for fraud detection response",
				Conditions: []MappingCondition{
					{Field: "fraud_score", Operator: "gt", Value: 0.7},
				},
			},
			"order_cancellation_spike": {
				Tools:       []string{"order_get_by_status", "metric_record", "notification_create"},
				Priority:    2,
				Description: "Tools for order cancellation spike response",
				Conditions: []MappingCondition{
					{Field: "cancellation_rate", Operator: "gt", Value: 0.3},
				},
			},
			"support_crisis": {
				Tools:       []string{"support_get_by_status", "support_create_ticket", "notification_send_push"},
				Priority:    1,
				Description: "Tools for support crisis management",
				Conditions: []MappingCondition{
					{Field: "support_tickets", Operator: "gt", Value: 100},
				},
			},
			"abandonment_risk": {
				Tools:       []string{"user_find", "notification_create", "mailer_send_templated"},
				Priority:    2,
				Description: "Tools for cart abandonment risk mitigation",
				Conditions: []MappingCondition{
					{Field: "abandonment_score", Operator: "gt", Value: 0.6},
				},
			},
		},
	}
	
	return defaultMappings, nil
}