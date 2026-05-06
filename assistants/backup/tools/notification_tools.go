package tools

import ai2 "middleman/internal/ai"

func createNotificationTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "notification_send",
				Description: "Send a notification",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "User ID",
						},
						"message": map[string]interface{}{
							"type":        "string",
							"description": "Notification message",
						},
					},
					"required": []string{"user_id", "message"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "notification_get_user_notifications",
				Description: "Get user notifications",
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