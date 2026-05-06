package tools

import ai2 "middleman/internal/ai"

func createOrderTools() []ai2.Tool {
	return []ai2.Tool{
		// CRUCIAL: Create order
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "order_create",
				Description: "Create a new order",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "User ID placing the order",
						},
						"basket_id": map[string]interface{}{
							"type":        "string",
							"description": "Basket ID to create order from",
						},
						"shipping_address": map[string]interface{}{
							"type":        "string",
							"description": "Shipping address",
						},
					},
					"required": []string{"user_id", "basket_id"},
				},
			},
		},
		// CRUCIAL: Get order by ID
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "order_get_by_id",
				Description: "Get order details by ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"order_id": map[string]interface{}{
							"type":        "string",
							"description": "Order ID",
						},
					},
					"required": []string{"order_id"},
				},
			},
		},
		// CRUCIAL: Get user orders
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "order_get_user_orders",
				Description: "Get all orders for a user",
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