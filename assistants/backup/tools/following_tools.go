package tools

import ai2 "middleman/internal/ai"

func createFollowingTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "following_create_new_relationship",
				Description: "Create a new follow relationship between users",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user who is following",
						},
						"followed_user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user being followed",
						},
						"followed_user_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of the user being followed",
							"enum":        []string{"user", "business", "creator"},
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "Optional message or note for the follow request",
						},
						"category_id": map[string]interface{}{
							"type":        "string",
							"description": "Category ID for categorizing the follow relationship",
						},
						"parent_id": map[string]interface{}{
							"type":        "string",
							"description": "Parent ID if this is a nested follow relationship",
						},
					},
					"required": []string{"user_id", "followed_user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "following_check_if_user_is_following",
				Description: "Check if one user is following another user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the potential follower",
						},
						"followed_user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the potentially followed user",
						},
					},
					"required": []string{"user_id", "followed_user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "following_get_all_followers_for_user",
				Description: "Get all followers for a specific user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"followed_user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user whose followers to retrieve",
						},
					},
					"required": []string{"followed_user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "following_get_user_following_list",
				Description: "Get list of users that a specific user is following",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user whose following list to retrieve",
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "following_get_total_follower_count",
				Description: "Get the total number of followers for a user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user",
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "following_get_total_following_count",
				Description: "Get the total number of users that a user is following",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user",
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "following_remove_relationship",
				Description: "Remove a follow relationship (unfollow)",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"follow_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the follow relationship to remove",
						},
					},
					"required": []string{"follow_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "following_get_most_followed_users",
				Description: "Get the most followed users on the platform",
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
				Name:        "following_get_mutual_between_users",
				Description: "Get mutual followers between two users",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id_1": map[string]interface{}{
							"type":        "string",
							"description": "ID of the first user",
						},
						"user_id_2": map[string]interface{}{
							"type":        "string",
							"description": "ID of the second user",
						},
					},
					"required": []string{"user_id_1", "user_id_2"},
				},
			},
		},
	}
}