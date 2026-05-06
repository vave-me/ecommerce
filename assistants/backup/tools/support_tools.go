package tools

import ai2 "middleman/internal/ai"

func createSupportTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "support_start",
				Description: "Initiate a support session for a user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user initiating support",
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "support_create_ticket",
				Description: "Create a new support ticket",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"support_id": map[string]interface{}{
							"type":        "string",
							"description": "Support channel ID",
						},
						"title": map[string]interface{}{
							"type":        "string",
							"description": "Ticket title/summary",
							"maxLength":   200,
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Detailed description of the issue",
							"maxLength":   5000,
						},
					},
					"required": []string{"support_id", "title", "description"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "support_list_tickets",
				Description: "List support tickets with pagination",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"support_id": map[string]interface{}{
							"type":        "string",
							"description": "Support channel ID",
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "Page number",
							"minimum":     1,
							"default":     1,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Items per page",
							"minimum":     1,
							"maximum":     100,
							"default":     20,
						},
					},
					"required": []string{"support_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "support_get_ticket",
				Description: "Get a specific support ticket by ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"ticket_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the ticket to retrieve",
						},
					},
					"required": []string{"ticket_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "support_update_ticket",
				Description: "Update support ticket status and assignment",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"ticket_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the ticket to update",
						},
						"status": map[string]interface{}{
							"type":        "string",
							"description": "New ticket status",
							"enum":        []string{"open", "in_progress", "resolved", "closed", "waiting_customer", "waiting_support"},
						},
						"assigned_to": map[string]interface{}{
							"type":        "string",
							"description": "ID of support agent to assign to",
						},
					},
					"required": []string{"ticket_id", "status"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "support_delete_ticket",
				Description: "Delete a support ticket",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"support_id": map[string]interface{}{
							"type":        "string",
							"description": "Support channel ID",
						},
						"ticket_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the ticket to delete",
						},
					},
					"required": []string{"support_id", "ticket_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "support_close_ticket",
				Description: "Close a support ticket with reason",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"ticket_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the ticket to close",
						},
						"reason": map[string]interface{}{
							"type":        "string",
							"description": "Reason for closing the ticket",
						},
					},
					"required": []string{"ticket_id", "reason"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "support_get_tickets",
				Description: "Get all tickets for a support channel",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"support_id": map[string]interface{}{
							"type":        "string",
							"description": "Support channel ID",
						},
					},
					"required": []string{"support_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "support_get_user_support",
				Description: "Get support information for a specific user",
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
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "support_get_by_status",
				Description: "Filter tickets by their current status",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Ticket status to filter by",
							"enum":        []string{"open", "in_progress", "resolved", "closed", "waiting_customer", "waiting_support"},
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     50,
						},
					},
					"required": []string{"status"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "support_search_tickets",
				Description: "Search support tickets by keyword",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query keyword",
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of results",
							"minimum":     1,
							"maximum":     100,
							"default":     50,
						},
					},
					"required": []string{"query"},
				},
			},
		},
	}
}