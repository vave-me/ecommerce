package tools

import ai2 "middleman/internal/ai"

func createNewsletterTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "newsletter_subscribe",
				Description: "Subscribe to newsletter",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type":        "string",
							"description": "Email address",
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
				Description: "Unsubscribe from newsletter",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"email": map[string]interface{}{
							"type":        "string",
							"description": "Email address",
						},
					},
					"required": []string{"email"},
				},
			},
		},
	}
}