package tools

import (
	ai2 "middleman/internal/ai"
)

// CreateGenericUserTools creates tool definitions for generic users with read-only access to most resources
// and write access only to their own data (wishlists, newsletters, notifications, etc.)
func CreateGenericUserTools() []ai2.Tool {
	tools := []ai2.Tool{}

	// READ-ONLY TOOLS FOR PUBLIC DATA
	
	// Product tools - read only
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "product_get_by_id",
				Description: "View detailed information about a specific product",
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
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "product_search_advanced",
				Description: "Advanced product search with filtering options",
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
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "product_get_by_category_id",
				Description: "Browse products by category",
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
	}...)

	// Service tools - read only
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "service_get_by_id",
				Description: "View details about a specific service",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"service_id": map[string]interface{}{
							"type":        "string",
							"description": "The unique identifier of the service",
							"minLength":   1,
						},
					},
					"required": []string{"service_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "service_search",
				Description: "Search for available services",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type":        "string",
							"description": "Search keyword",
						},
						"category": map[string]interface{}{
							"type":        "string",
							"description": "Service category",
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "Page number",
							"minimum":     1,
							"default":     1,
						},
					},
					"required": []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "service_get_available",
				Description: "Get list of all available services",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
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
	}...)

	// Post tools - read only
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "post_get_by_id",
				Description: "View a specific post",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"post_id": map[string]interface{}{
							"type":        "string",
							"description": "The unique identifier of the post",
							"minLength":   1,
						},
					},
					"required": []string{"post_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "post_list_public",
				Description: "Browse public posts",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
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
							"maximum":     50,
							"default":     10,
						},
					},
					"required": []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "post_search",
				Description: "Search posts by title or content",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query",
							"minLength":   1,
						},
					},
					"required": []string{"query"},
				},
			},
		},
	}...)

	// Category tools - read only
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "category_get_by_id",
				Description: "Get category details",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "string",
							"description": "The category ID",
							"minLength":   1,
						},
					},
					"required": []string{"category_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "category_list_all",
				Description: "Get all available categories",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "category_get_children",
				Description: "Get subcategories of a parent category",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"parent_id": map[string]interface{}{
							"type":        "string",
							"description": "Parent category ID",
							"minLength":   1,
						},
					},
					"required": []string{"parent_id"},
				},
			},
		},
	}...)

	// USER-SPECIFIC TOOLS WITH WRITE ACCESS

	// Wishlist tools - user can manage their own wishlist
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "wishlist_add_item",
				Description: "Add a product to your wishlist",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "Product ID to add",
							"minLength":   1,
						},
					},
					"required": []string{"product_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "wishlist_remove_item",
				Description: "Remove a product from your wishlist",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "Product ID to remove",
							"minLength":   1,
						},
					},
					"required": []string{"product_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "wishlist_get_my_items",
				Description: "View your wishlist items",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
	}...)

	// Newsletter tools - user can manage their subscriptions
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "newsletter_subscribe",
				Description: "Subscribe to newsletter updates",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type":        "string",
							"description": "Email address",
							"format":      "email",
						},
						"preferences": map[string]interface{}{
							"type":        "array",
							"description": "Newsletter preferences (e.g., 'deals', 'new-products', 'updates')",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"email"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "newsletter_unsubscribe",
				Description: "Unsubscribe from newsletters",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type":        "string",
							"description": "Email address",
							"format":      "email",
						},
					},
					"required": []string{"email"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "newsletter_update_preferences",
				Description: "Update your newsletter preferences",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"preferences": map[string]interface{}{
							"type":        "array",
							"description": "Updated preferences",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"preferences"},
				},
			},
		},
	}...)

	// Basket tools - user can manage their own basket
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "basket_add_item",
				Description: "Add a product to your shopping basket",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "Product ID to add",
							"minLength":   1,
						},
						"quantity": map[string]interface{}{
							"type":        "integer",
							"description": "Quantity to add",
							"minimum":     1,
							"default":     1,
						},
					},
					"required": []string{"product_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "basket_remove_item",
				Description: "Remove a product from your basket",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "Product ID to remove",
							"minLength":   1,
						},
					},
					"required": []string{"product_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "basket_update_quantity",
				Description: "Update quantity of a product in your basket",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "Product ID",
							"minLength":   1,
						},
						"quantity": map[string]interface{}{
							"type":        "integer",
							"description": "New quantity",
							"minimum":     0,
						},
					},
					"required": []string{"product_id", "quantity"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "basket_get_current",
				Description: "View your current basket contents",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "basket_clear",
				Description: "Clear all items from your basket",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
	}...)

	// Order tools - user can view their own orders
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "order_get_my_orders",
				Description: "View your order history",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Filter by status (pending, processing, shipped, delivered, cancelled)",
							"enum":        []string{"pending", "processing", "shipped", "delivered", "cancelled"},
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "Page number",
							"minimum":     1,
							"default":     1,
						},
					},
					"required": []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "order_get_details",
				Description: "View details of your specific order",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"order_id": map[string]interface{}{
							"type":        "string",
							"description": "Order ID",
							"minLength":   1,
						},
					},
					"required": []string{"order_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "order_track_shipment",
				Description: "Track shipment for your order",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"order_id": map[string]interface{}{
							"type":        "string",
							"description": "Order ID",
							"minLength":   1,
						},
					},
					"required": []string{"order_id"},
				},
			},
		},
	}...)

	// User profile tools - limited to own profile
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "user_get_my_profile",
				Description: "View your profile information",
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "user_update_profile",
				Description: "Update your profile information",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"display_name": map[string]interface{}{
							"type":        "string",
							"description": "Display name",
							"maxLength":   100,
						},
						"phone": map[string]interface{}{
							"type":        "string",
							"description": "Phone number",
						},
						"bio": map[string]interface{}{
							"type":        "string",
							"description": "Profile bio",
							"maxLength":   500,
						},
					},
					"required": []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "user_update_preferences",
				Description: "Update your account preferences",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"language": map[string]interface{}{
							"type":        "string",
							"description": "Preferred language",
							"enum":        []string{"en", "es", "fr", "de", "it", "pt", "zh", "ja"},
						},
						"currency": map[string]interface{}{
							"type":        "string",
							"description": "Preferred currency",
							"enum":        []string{"USD", "EUR", "GBP", "CAD", "AUD"},
						},
						"notifications": map[string]interface{}{
							"type":        "boolean",
							"description": "Enable notifications",
						},
					},
					"required": []string{},
				},
			},
		},
	}...)

	// Comment tools - user can add comments
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "comment_add",
				Description: "Add a comment to a product or post",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"entity_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of entity to comment on",
							"enum":        []string{"product", "post", "service"},
						},
						"entity_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the entity",
							"minLength":   1,
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "Comment content",
							"minLength":   1,
							"maxLength":   1000,
						},
					},
					"required": []string{"entity_type", "entity_id", "content"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "comment_get_for_entity",
				Description: "Get comments for a product or post",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"entity_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of entity",
							"enum":        []string{"product", "post", "service"},
						},
						"entity_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the entity",
							"minLength":   1,
						},
					},
					"required": []string{"entity_type", "entity_id"},
				},
			},
		},
	}...)

	// Notification tools - user's own notifications
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "notification_get_my_notifications",
				Description: "Get your notifications",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"unread_only": map[string]interface{}{
							"type":        "boolean",
							"description": "Show only unread notifications",
							"default":     false,
						},
					},
					"required": []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "notification_mark_as_read",
				Description: "Mark notification as read",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"notification_id": map[string]interface{}{
							"type":        "string",
							"description": "Notification ID",
							"minLength":   1,
						},
					},
					"required": []string{"notification_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "notification_update_preferences",
				Description: "Update notification preferences",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"email_notifications": map[string]interface{}{
							"type":        "boolean",
							"description": "Receive email notifications",
						},
						"push_notifications": map[string]interface{}{
							"type":        "boolean",
							"description": "Receive push notifications",
						},
						"order_updates": map[string]interface{}{
							"type":        "boolean",
							"description": "Receive order updates",
						},
						"promotions": map[string]interface{}{
							"type":        "boolean",
							"description": "Receive promotional notifications",
						},
					},
					"required": []string{},
				},
			},
		},
	}...)

	// Activity tools - read only for viewing activities
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "activity_get_for_product",
				Description: "View activity for a product (likes, views)",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "Product ID",
							"minLength":   1,
						},
					},
					"required": []string{"product_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "activity_like_product",
				Description: "Like a product",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"product_id": map[string]interface{}{
							"type":        "string",
							"description": "Product ID to like",
							"minLength":   1,
						},
					},
					"required": []string{"product_id"},
				},
			},
		},
	}...)

	// Media tools - read only
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "media_get_for_entity",
				Description: "Get media files for a product or post",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"entity_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of entity",
							"enum":        []string{"product", "post", "service"},
						},
						"entity_id": map[string]interface{}{
							"type":        "string",
							"description": "Entity ID",
							"minLength":   1,
						},
					},
					"required": []string{"entity_type", "entity_id"},
				},
			},
		},
	}...)

	// Shipping tools - calculate shipping for user's orders
	tools = append(tools, []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "shipping_calculate_rate",
				Description: "Calculate shipping rate for your order",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"from_zip": map[string]interface{}{
							"type":        "string",
							"description": "Origin ZIP code",
							"pattern":     "^[0-9]{5}$",
						},
						"to_zip": map[string]interface{}{
							"type":        "string",
							"description": "Destination ZIP code",
							"pattern":     "^[0-9]{5}$",
						},
						"weight": map[string]interface{}{
							"type":        "integer",
							"description": "Package weight in grams",
							"minimum":     1,
						},
					},
					"required": []string{"from_zip", "to_zip", "weight"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "shipping_get_options",
				Description: "Get available shipping options",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"weight": map[string]interface{}{
							"type":        "integer",
							"description": "Package weight in grams",
							"minimum":     1,
						},
					},
					"required": []string{"weight"},
				},
			},
		},
	}...)

	return tools
}