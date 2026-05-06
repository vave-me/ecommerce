package tools

import ai2 "middleman/internal/ai"

func createGeocodingTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_batch_geocode_address",
				Description: "Geocode multiple addresses in a single batch operation",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"addresses": map[string]interface{}{
							"type":        "array",
							"description": "List of addresses to geocode",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"addresses"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_geocode_address",
				Description: "Convert a street address to geographic coordinates (latitude/longitude)",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "string",
							"description": "Unique identifier for this geocoding request",
						},
						"address": map[string]interface{}{
							"type":        "string",
							"description": "The street address to geocode",
						},
					},
					"required": []string{"address"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_reverse_geocode",
				Description: "Convert geographic coordinates to a street address",
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
				Name:        "geocoding_validate_address",
				Description: "Validate and standardize a street address",
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
				Name:        "geocoding_suggest_address",
				Description: "Get address suggestions based on partial input",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"address": map[string]interface{}{
							"type":        "string",
							"description": "Partial address to get suggestions for",
						},
					},
					"required": []string{"address"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_suggest_city",
				Description: "Get city suggestions based on partial city name",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Partial city name to get suggestions for",
						},
					},
					"required": []string{"name"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_refresh_cache",
				Description: "Refresh the geocoding cache to clear outdated entries",
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
				Name:        "geocoding_get_requests",
				Description: "Get a list of geocoding requests with pagination",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "Page number for pagination",
							"minimum":     1,
							"default":     1,
						},
						"page_size": map[string]interface{}{
							"type":        "integer",
							"description": "Number of results per page",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
						"sort_by": map[string]interface{}{
							"type":        "string",
							"description": "Field to sort by",
							"default":     "created_at",
						},
						"sort_order": map[string]interface{}{
							"type":        "string",
							"description": "Sort order",
							"enum":        []string{"asc", "desc"},
							"default":     "desc",
						},
					},
					"required": []string{},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_search_with_term",
				Description: "Search geocoding requests by term",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"term": map[string]interface{}{
							"type":        "string",
							"description": "Search term to find in geocoding requests",
						},
					},
					"required": []string{"term"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_find_request",
				Description: "Find a specific geocoding request by ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"request_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the geocoding request",
						},
					},
					"required": []string{"request_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_get_user_history",
				Description: "Get geocoding history for a specific user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the user",
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "Page number for pagination",
							"minimum":     1,
							"default":     1,
						},
						"page_size": map[string]interface{}{
							"type":        "integer",
							"description": "Number of results per page",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_get_by_location",
				Description: "Get geocoding requests within a radius of a location",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"lat": map[string]interface{}{
							"type":        "number",
							"description": "Latitude of the center point",
						},
						"lng": map[string]interface{}{
							"type":        "number",
							"description": "Longitude of the center point",
						},
						"radius": map[string]interface{}{
							"type":        "integer",
							"description": "Search radius in meters",
							"minimum":     1,
							"default":     1000,
						},
					},
					"required": []string{"lat", "lng"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_clear_user_cache",
				Description: "Clear geocoding cache for a specific user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "The ID of the user whose cache to clear",
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_get_stats",
				Description: "Get statistics about geocoding usage",
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
				Name:        "geocoding_get_coordinates",
				Description: "Get geographic coordinates for an address (legacy method)",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"address": map[string]interface{}{
							"type":        "string",
							"description": "The address to get coordinates for",
						},
					},
					"required": []string{"address"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "geocoding_get_address",
				Description: "Get address for geographic coordinates (legacy method)",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"lat": map[string]interface{}{
							"type":        "number",
							"description": "Latitude",
						},
						"lng": map[string]interface{}{
							"type":        "number",
							"description": "Longitude",
						},
					},
					"required": []string{"lat", "lng"},
				},
			},
		},
	}
}
