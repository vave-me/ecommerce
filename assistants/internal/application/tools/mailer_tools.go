package tools

import ai2 "middleman/internal/ai"

func createMailerTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "mailer_send_email",
				Description: "Send an email",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"to": map[string]interface{}{
							"type":        "string",
							"description": "Recipient email",
						},
						"subject": map[string]interface{}{
							"type":        "string",
							"description": "Email subject",
						},
						"body": map[string]interface{}{
							"type":        "string",
							"description": "Email body",
						},
					},
					"required": []string{"to", "subject", "body"},
				},
			},
		},
	}
}