package tools

import ai2 "middleman/internal/ai"

func createUserTools() []ai2.Tool {
	return []ai2.Tool{
		// CRUCIAL: Get user by ID
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "user_get_by_id",
				Description: "Retrieve user information by user ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "The unique identifier of the user",
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		// CRUCIAL: Search users
		{
			Type: "function", 
			Function: ai2.FunctionDef{
				Name:        "user_search",
				Description: "Search for users by name or email",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query for user name or email",
							"minLength":   1,
						},
					},
					"required": []string{"query"},
				},
			},
		},
		// CRUCIAL: Update user profile
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "user_update_profile",
				Description: "Update user profile information",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "User ID to update",
						},
						"full_name": map[string]interface{}{
							"type":        "string",
							"description": "User's full name",
						},
						"email": map[string]interface{}{
							"type":        "string",
							"description": "User's email address",
						},
						"phone": map[string]interface{}{
							"type":        "string",
							"description": "User's phone number",
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		// CRUCIAL: Authenticate user
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "user_authenticate",
				Description: "Authenticate user with email and password",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type":        "string",
							"description": "User's email address",
						},
						"password": map[string]interface{}{
							"type":        "string",
							"description": "User's password",
						},
					},
					"required": []string{"email", "password"},
				},
			},
		},
	}
}