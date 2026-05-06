package tools

import ai2 "middleman/internal/ai"

func createVectorTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "vector_search_similar",
				Description: "Find entities similar to a given entity using vector similarity",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"entity_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the entity to find similar items for",
						},
						"entity_type": map[string]interface{}{
							"type":        "string",
							"description": "The type of entity",
							"enum":        []string{"product", "service", "user", "post"},
						},
						"options": map[string]interface{}{
							"type":        "object",
							"description": "Search options",
							"properties": map[string]interface{}{
								"top_k": map[string]interface{}{
									"type":        "integer",
									"description": "Number of similar items to return",
									"minimum":     1,
									"maximum":     100,
									"default":     10,
								},
								"exclude_self": map[string]interface{}{
									"type":        "boolean",
									"description": "Whether to exclude the source entity from results",
									"default":     true,
								},
								"target_entity_types": map[string]interface{}{
									"type":        "array",
									"description": "Limit results to specific entity types",
									"items":       map[string]interface{}{"type": "string"},
								},
							},
						},
					},
					"required": []string{"entity_id", "entity_type"},
				},
			},
		},
	}
}
