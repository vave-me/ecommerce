package tools

import ai2 "middleman/internal/ai"

func createMetricTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_update_item",
				Description: "Update metrics for a specific item (views, likes, shares, etc.)",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the item to update metrics for",
							"minLength":   1,
						},
						"metric_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of metric to update",
							"enum":        []string{"views", "likes", "shares", "saves", "clicks", "impressions", "purchases", "ratings"},
						},
						"metric_type_action": map[string]interface{}{
							"type":        "string",
							"description": "Action to perform on the metric",
							"enum":        []string{"increment", "decrement", "set", "reset"},
						},
					},
					"required": []string{"item_id", "metric_type", "metric_type_action"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_share_item",
				Description: "Record that an item has been shared",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the item that was shared",
							"minLength":   1,
						},
					},
					"required": []string{"item_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_visit_item",
				Description: "Record that an item has been visited/viewed",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the item that was visited",
							"minLength":   1,
						},
					},
					"required": []string{"item_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_update_user",
				Description: "Update metrics for a specific user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the user to update metrics for",
							"minLength":   1,
						},
						"metric_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of metric to update",
							"enum":        []string{"followers", "following", "posts", "purchases", "sales", "reviews", "ratings", "reputation"},
						},
						"metric_type_action": map[string]interface{}{
							"type":        "string",
							"description": "Action to perform on the metric",
							"enum":        []string{"increment", "decrement", "set", "reset"},
						},
					},
					"required": []string{"user_id", "metric_type", "metric_type_action"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_user",
				Description: "Get all metrics for a specific user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the user to get metrics for",
							"minLength":   1,
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_item",
				Description: "Get all metrics for a specific item",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the item to get metrics for",
							"minLength":   1,
						},
					},
					"required": []string{"item_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_items",
				Description: "Get metrics for multiple items",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"item_ids": map[string]interface{}{
							"type":        "array",
							"description": "List of item IDs to get metrics for",
							"items": map[string]interface{}{
								"type": "string",
							},
							"minItems": 1,
							"maxItems": 100,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results to return",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"item_ids"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_highest_by_type",
				Description: "Get items with the highest values for a specific metric type",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"metric_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of metric to sort by",
							"enum":        []string{"views", "likes", "shares", "saves", "clicks", "impressions", "purchases", "ratings"},
						},
						"entity_types": map[string]interface{}{
							"type":        "array",
							"description": "Types of entities to include",
							"items": map[string]interface{}{
								"type": "string",
								"enum": []string{"product", "service", "post", "offer"},
							},
						},
						"category_id": map[string]interface{}{
							"type":        "string",
							"description": "Filter by category ID",
						},
						"min_price": map[string]interface{}{
							"type":        "integer",
							"description": "Minimum price filter",
							"minimum":     0,
						},
						"max_price": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum price filter",
							"minimum":     0,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
						"lat": map[string]interface{}{
							"type":        "number",
							"description": "Latitude for location-based filtering",
							"minimum":     -90,
							"maximum":     90,
						},
						"lng": map[string]interface{}{
							"type":        "number",
							"description": "Longitude for location-based filtering",
							"minimum":     -180,
							"maximum":     180,
						},
						"radius": map[string]interface{}{
							"type":        "number",
							"description": "Radius in kilometers for location filtering",
							"minimum":     0,
							"maximum":     1000,
						},
						"created_from": map[string]interface{}{
							"type":        "string",
							"description": "Start date for creation filter (ISO 8601)",
							"format":      "date-time",
						},
						"created_to": map[string]interface{}{
							"type":        "string",
							"description": "End date for creation filter (ISO 8601)",
							"format":      "date-time",
						},
					},
					"required": []string{"metric_type"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_lowest_by_type",
				Description: "Get items with the lowest values for a specific metric type",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"metric_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of metric to sort by",
							"enum":        []string{"views", "likes", "shares", "saves", "clicks", "impressions", "purchases", "ratings"},
						},
						"entity_types": map[string]interface{}{
							"type":        "array",
							"description": "Types of entities to include",
							"items": map[string]interface{}{
								"type": "string",
								"enum": []string{"product", "service", "post", "offer"},
							},
						},
						"category_id": map[string]interface{}{
							"type":        "string",
							"description": "Filter by category ID",
						},
						"min_price": map[string]interface{}{
							"type":        "integer",
							"description": "Minimum price filter",
							"minimum":     0,
						},
						"max_price": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum price filter",
							"minimum":     0,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
						"lat": map[string]interface{}{
							"type":        "number",
							"description": "Latitude for location-based filtering",
							"minimum":     -90,
							"maximum":     90,
						},
						"lng": map[string]interface{}{
							"type":        "number",
							"description": "Longitude for location-based filtering",
							"minimum":     -180,
							"maximum":     180,
						},
						"radius": map[string]interface{}{
							"type":        "number",
							"description": "Radius in kilometers for location filtering",
							"minimum":     0,
							"maximum":     1000,
						},
						"created_from": map[string]interface{}{
							"type":        "string",
							"description": "Start date for creation filter (ISO 8601)",
							"format":      "date-time",
						},
						"created_to": map[string]interface{}{
							"type":        "string",
							"description": "End date for creation filter (ISO 8601)",
							"format":      "date-time",
						},
					},
					"required": []string{"metric_type"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_item_by_type",
				Description: "Get specific metric type for an item",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the item",
							"minLength":   1,
						},
						"metric_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of metric to retrieve",
							"enum":        []string{"views", "likes", "shares", "saves", "clicks", "impressions", "purchases", "ratings"},
						},
					},
					"required": []string{"item_id", "metric_type"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_user_by_type",
				Description: "Get specific metric type for a user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the user",
							"minLength":   1,
						},
						"metric_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of metric to retrieve",
							"enum":        []string{"followers", "following", "posts", "purchases", "sales", "reviews", "ratings", "reputation"},
						},
					},
					"required": []string{"user_id", "metric_type"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_items_by_category",
				Description: "Get item metrics filtered by category",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"category_id": map[string]interface{}{
							"type":        "string",
							"description": "Category ID to filter by",
							"minLength":   1,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"category_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_top_items",
				Description: "Get top items by a specific metric",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"metric_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of metric to rank by",
							"enum":        []string{"views", "likes", "shares", "saves", "clicks", "impressions", "purchases", "ratings"},
						},
						"entity_types": map[string]interface{}{
							"type":        "array",
							"description": "Types of entities to include",
							"items": map[string]interface{}{
								"type": "string",
								"enum": []string{"product", "service", "post", "offer"},
							},
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     10,
						},
					},
					"required": []string{"metric_type"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_top_users",
				Description: "Get top users by a specific metric",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"metric_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of metric to rank by",
							"enum":        []string{"followers", "following", "posts", "purchases", "sales", "reviews", "ratings", "reputation"},
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     10,
						},
					},
					"required": []string{"metric_type"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_items_in_range",
				Description: "Get item metrics within a geographic range",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"lat": map[string]interface{}{
							"type":        "number",
							"description": "Latitude of the center point",
							"minimum":     -90,
							"maximum":     90,
						},
						"lng": map[string]interface{}{
							"type":        "number",
							"description": "Longitude of the center point",
							"minimum":     -180,
							"maximum":     180,
						},
						"radius": map[string]interface{}{
							"type":        "number",
							"description": "Radius in kilometers",
							"minimum":     0,
							"maximum":     1000,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"lat", "lng", "radius"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_trending_items",
				Description: "Get trending items based on recent activity",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"entity_types": map[string]interface{}{
							"type":        "array",
							"description": "Types of entities to include",
							"items": map[string]interface{}{
								"type": "string",
								"enum": []string{"product", "service", "post", "offer"},
							},
						},
						"days": map[string]interface{}{
							"type":        "integer",
							"description": "Number of days to consider for trending",
							"minimum":     1,
							"maximum":     30,
							"default":     7,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
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
				Name:        "metric_compare_items",
				Description: "Compare metrics between two items",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"item_id1": map[string]interface{}{
							"type":        "string",
							"description": "First item ID to compare",
							"minLength":   1,
						},
						"item_id2": map[string]interface{}{
							"type":        "string",
							"description": "Second item ID to compare",
							"minLength":   1,
						},
					},
					"required": []string{"item_id1", "item_id2"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_compare_users",
				Description: "Compare metrics between two users",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id1": map[string]interface{}{
							"type":        "string",
							"description": "First user ID to compare",
							"minLength":   1,
						},
						"user_id2": map[string]interface{}{
							"type":        "string",
							"description": "Second user ID to compare",
							"minLength":   1,
						},
					},
					"required": []string{"user_id1", "user_id2"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_reset_item",
				Description: "Reset specific metrics for an item",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the item to reset metrics for",
							"minLength":   1,
						},
						"metric_types": map[string]interface{}{
							"type":        "array",
							"description": "Types of metrics to reset",
							"items": map[string]interface{}{
								"type": "string",
								"enum": []string{"views", "likes", "shares", "saves", "clicks", "impressions", "purchases", "ratings"},
							},
							"minItems": 1,
						},
					},
					"required": []string{"item_id", "metric_types"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_reset_user",
				Description: "Reset specific metrics for a user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the user to reset metrics for",
							"minLength":   1,
						},
						"metric_types": map[string]interface{}{
							"type":        "array",
							"description": "Types of metrics to reset",
							"items": map[string]interface{}{
								"type": "string",
								"enum": []string{"followers", "following", "posts", "purchases", "sales", "reviews", "ratings", "reputation"},
							},
							"minItems": 1,
						},
					},
					"required": []string{"user_id", "metric_types"},
				},
			},
		},

		// ==================== SYSTEM METRICS ====================
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_record",
				Description: "Record a metric data point with tags and unit",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "The metric name (e.g., 'api.request.count', 'user.login.duration')",
							"minLength":   1,
							"pattern":     "^[a-z0-9._-]+$",
						},
						"value": map[string]interface{}{
							"type":        "number",
							"description": "The metric value to record",
						},
						"unit": map[string]interface{}{
							"type":        "string",
							"description": "The unit of measurement",
							"enum":        []string{"count", "milliseconds", "seconds", "bytes", "percent", "ratio"},
						},
						"tags": map[string]interface{}{
							"type":        "object",
							"description": "Key-value pairs for categorizing the metric",
							"additionalProperties": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"name", "value", "unit"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get",
				Description: "Retrieve metric data for a specific time range",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "The metric name to retrieve",
							"minLength":   1,
							"pattern":     "^[a-z0-9._-]+$",
						},
						"start_time": map[string]interface{}{
							"type":        "string",
							"description": "Start time in ISO 8601 format",
							"format":      "date-time",
						},
						"end_time": map[string]interface{}{
							"type":        "string",
							"description": "End time in ISO 8601 format",
							"format":      "date-time",
						},
					},
					"required": []string{"name", "start_time", "end_time"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_get_system",
				Description: "Get overall system metrics and statistics",
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
				Name:        "metric_increment_counter",
				Description: "Increment a counter metric by 1",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "The counter metric name to increment",
							"minLength":   1,
							"pattern":     "^[a-z0-9._-]+$",
						},
						"tags": map[string]interface{}{
							"type":        "object",
							"description": "Key-value pairs for categorizing the metric",
							"additionalProperties": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"name"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_set_gauge",
				Description: "Set a gauge metric to a specific value",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "The gauge metric name",
							"minLength":   1,
							"pattern":     "^[a-z0-9._-]+$",
						},
						"value": map[string]interface{}{
							"type":        "number",
							"description": "The value to set the gauge to",
						},
						"tags": map[string]interface{}{
							"type":        "object",
							"description": "Key-value pairs for categorizing the metric",
							"additionalProperties": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"name", "value"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "metric_record_histogram",
				Description: "Record a value in a histogram metric",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "The histogram metric name",
							"minLength":   1,
							"pattern":     "^[a-z0-9._-]+$",
						},
						"value": map[string]interface{}{
							"type":        "number",
							"description": "The value to record in the histogram",
						},
						"tags": map[string]interface{}{
							"type":        "object",
							"description": "Key-value pairs for categorizing the metric",
							"additionalProperties": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"name", "value"},
				},
			},
		},
	}
}
