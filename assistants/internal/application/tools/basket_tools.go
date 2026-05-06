package tools

import ai2 "middleman/internal/ai"

func createBasketTools() []ai2.Tool {
	return []ai2.Tool{
		// CRUCIAL: Add item to basket
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "basket_add_item",
				Description: "Add a product to the shopping basket",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "User ID who owns the basket",
						},
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "Product ID to add",
						},
						"quantity": map[string]interface{}{
							"type":        "integer",
							"description": "Quantity to add",
							"minimum":     1,
						},
					},
					"required": []string{"user_id", "product_id", "quantity"},
				},
			},
		},
		// CRUCIAL: Remove item from basket
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "basket_remove_product_from_basket",
				Description: "Remove a product from the basket",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"basket_id": map[string]interface{}{
							"type":        "string",
							"description": "Basket ID",
						},
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "Product ID to remove",
						},
					},
					"required": []string{"basket_id", "product_id"},
				},
			},
		},
		// CRUCIAL: Get basket
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "basket_get_user_current_basket",
				Description: "Get the current basket for a user",
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
		// CRUCIAL: Clear basket
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "basket_clear_all_basket_items",
				Description: "Clear all items from a basket",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"basket_id": map[string]interface{}{
							"type":        "string",
							"description": "Basket ID to clear",
						},
					},
					"required": []string{"basket_id"},
				},
			},
		},
		// CRUCIAL: Calculate total
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "basket_calculate_basket_total",
				Description: "Calculate the total price of items in a basket",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"basket_id": map[string]interface{}{
							"type":        "string",
							"description": "Basket ID",
						},
					},
					"required": []string{"basket_id"},
				},
			},
		},
	}
}