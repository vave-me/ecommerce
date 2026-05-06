package tools

import ai2 "middleman/internal/ai"

func createCategoryTools() []ai2.Tool {
	return []ai2.Tool{
		// CRUCIAL: Find category by ID
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "category_find",
				Description: "Find a category by its ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "string",
							"description": "The unique identifier of the category",
							"minLength":   1,
						},
					},
					"required": []string{"category_id"},
				},
			},
		},
		// CRUCIAL: Get all categories
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "category_get_all",
				Description: "Get all categories in the system",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
		// CRUCIAL: Get category by slug
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "category_get_by_slug",
				Description: "Get a category by its URL slug",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"slug": map[string]interface{}{
							"type":        "string",
							"description": "The URL-friendly slug of the category",
							"minLength":   1,
						},
					},
					"required": []string{"slug"},
				},
			},
		},
		// CRUCIAL: Search categories
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "category_search",
				Description: "Search categories by name",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query for category names",
							"minLength":   1,
						},
					},
					"required": []string{"query"},
				},
			},
		},
	}
}