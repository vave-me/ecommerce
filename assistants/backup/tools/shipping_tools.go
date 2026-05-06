package tools

import ai2 "middleman/internal/ai"

func createShippingTools() []ai2.Tool {
	return []ai2.Tool{
		// CRUCIAL: Calculate shipping cost
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "shipping_calculate_cost",
				Description: "Calculate shipping cost",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"weight": map[string]interface{}{
							"type":        "number",
							"description": "Weight in kg",
						},
						"destination": map[string]interface{}{
							"type":        "string",
							"description": "Destination address",
						},
					},
					"required": []string{"weight", "destination"},
				},
			},
		},
		// CRUCIAL: Track shipment
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "shipping_track",
				Description: "Track a shipment",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"tracking_number": map[string]interface{}{
							"type":        "string",
							"description": "Tracking number",
						},
					},
					"required": []string{"tracking_number"},
				},
			},
		},
	}
}