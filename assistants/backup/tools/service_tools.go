package tools

import ai2 "middleman/internal/ai"

func createServiceTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "service_create",
				Description: "Create a new service",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Service name",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Service description",
						},
						"price": map[string]interface{}{
							"type":        "integer",
							"description": "Price in cents",
						},
					},
					"required": []string{"name", "description", "price"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "service_search",
				Description: "Search services",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "service_get_by_id",
				Description: "Get service by ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"service_id": map[string]interface{}{
							"type":        "string",
							"description": "Service ID",
						},
					},
					"required": []string{"service_id"},
				},
			},
		},
	}
}