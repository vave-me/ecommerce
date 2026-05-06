package tools

import ai2 "middleman/internal/ai"

func createWishlistTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "wishlist_add_item",
				Description: "Add item to wishlist",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "User ID",
						},
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "Product ID",
						},
					},
					"required": []string{"user_id", "product_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "wishlist_remove_item",
				Description: "Remove item from wishlist",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "User ID",
						},
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "Product ID",
						},
					},
					"required": []string{"user_id", "product_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "wishlist_get",
				Description: "Get user's wishlist",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "User ID",
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
	}
}