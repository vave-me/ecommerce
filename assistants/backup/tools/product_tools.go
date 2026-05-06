package tools

import ai2 "middleman/internal/ai"

func createProductTools() []ai2.Tool {
	return []ai2.Tool{
		// CRUCIAL: Get product details
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "product_get_by_id",
				Description: "Retrieve detailed information about a specific product",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "The unique identifier of the product",
							"minLength":   1,
						},
					},
					"required": []string{"product_id"},
				},
			},
		},
		// CRUCIAL: Search products
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "product_search_by_name",
				Description: "Search for products by name",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Product name to search for",
							"minLength":   1,
						},
					},
					"required": []string{"name"},
				},
			},
		},
		// CRUCIAL: Advanced search
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "product_search_advanced",
				Description: "Advanced product search with comprehensive filtering options",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Search by product name",
						},
						"category_id": map[string]interface{}{
							"type":        "string",
							"description": "Filter by category ID",
						},
						"min_price": map[string]interface{}{
							"type":        "integer",
							"description": "Minimum price in cents",
							"minimum":     0,
						},
						"max_price": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum price in cents",
							"minimum":     0,
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "Page number",
							"minimum":     1,
							"default":     1,
						},
						"page_size": map[string]interface{}{
							"type":        "integer",
							"description": "Results per page",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{},
				},
			},
		},
		// CRUCIAL: Get products by category
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "product_get_by_category_id",
				Description: "Get products by category ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "string",
							"description": "Category ID to filter by",
							"minLength":   1,
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "Page number",
							"minimum":     1,
							"default":     1,
						},
						"page_size": map[string]interface{}{
							"type":        "integer",
							"description": "Results per page",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"category_id"},
				},
			},
		},
		// CRUCIAL: Create product
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "product_create",
				Description: "Create a new product listing",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Product name",
							"minLength":   1,
							"maxLength":   200,
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Product description",
							"minLength":   1,
							"maxLength":   5000,
						},
						"base_price": map[string]interface{}{
							"type":        "integer",
							"description": "Base price in cents",
							"minimum":     0,
						},
						"category_id": map[string]interface{}{
							"type":        "string",
							"description": "Category ID for the product",
							"minLength":   1,
						},
					},
					"required": []string{"name", "description", "base_price", "category_id"},
				},
			},
		},
		// CRUCIAL: Update product price
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "product_update_base_price",
				Description: "Update the base price of a product",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "The product ID to update",
							"minLength":   1,
						},
						"new_base_price": map[string]interface{}{
							"type":        "integer",
							"description": "New base price in cents",
							"minimum":     0,
						},
					},
					"required": []string{"product_id", "new_base_price"},
				},
			},
		},
	}
}