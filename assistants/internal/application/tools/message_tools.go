package tools

import ai2 "middleman/internal/ai"

func createMessageTools() []ai2.Tool {
	return []ai2.Tool{
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_create_new_conversation",
				Description: "Create a new conversation between users about a specific item",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"sender_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user initiating the conversation",
						},
						"recipient_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user receiving the message",
						},
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the item being discussed",
						},
					},
					"required": []string{"sender_id", "recipient_id", "item_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_restore_archived_conversation",
				Description: "Restore an archived conversation to active status",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"conversation_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the conversation to restore",
						},
					},
					"required": []string{"conversation_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_archive_existing_conversation",
				Description: "Archive an active conversation",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"conversation_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the conversation to archive",
						},
					},
					"required": []string{"conversation_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_get_conversation_by_id",
				Description: "Get a specific conversation by its ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"conversation_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the conversation to retrieve",
						},
					},
					"required": []string{"conversation_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_find_conversation_by_recipient_and_item",
				Description: "Find a conversation by recipient and item",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"recipient_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the recipient user",
						},
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the item being discussed",
						},
					},
					"required": []string{"recipient_id", "item_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_get_user_conversation_list",
				Description: "Get all conversations for a user with pagination",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user",
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "Page number for pagination",
							"minimum":     1,
							"default":     1,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Number of conversations per page",
							"minimum":     1,
							"maximum":     100,
							"default":     50,
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_get_user_active_conversations",
				Description: "Get only active (non-archived) conversations for a user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user",
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum number of conversations to return",
							"minimum":     1,
							"maximum":     100,
							"default":     50,
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_send_new_message",
				Description: "Send a message in an existing conversation",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"message_id": map[string]interface{}{
							"type":        "string",
							"description": "Unique ID for this message",
						},
						"conversation_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the conversation",
						},
						"sender_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the message sender",
						},
						"recipient_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the message recipient",
						},
						"item_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the item being discussed",
						},
						"body": map[string]interface{}{
							"type":        "string",
							"description": "Message content",
						},
						"is_read": map[string]interface{}{
							"type":        "boolean",
							"description": "Whether the message has been read",
							"default":     false,
						},
					},
					"required": []string{"message_id", "conversation_id", "sender_id", "recipient_id", "item_id", "body"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_delete_message_by_id",
				Description: "Delete a specific message",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"message_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the message to delete",
						},
					},
					"required": []string{"message_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_get_message_by_id",
				Description: "Get a specific message by its ID",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"message_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the message to retrieve",
						},
					},
					"required": []string{"message_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_get_conversation_messages",
				Description: "Get all messages in a conversation with pagination",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"conversation_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the conversation",
						},
						"page": map[string]interface{}{
							"type":        "integer",
							"description": "Page number for pagination",
							"minimum":     1,
							"default":     1,
						},
						"limit": map[string]interface{}{
							"type":        "integer",
							"description": "Number of messages per page",
							"minimum":     1,
							"maximum":     100,
							"default":     50,
						},
					},
					"required": []string{"conversation_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_mark_message_as_read_by_user",
				Description: "Mark a message as read by a specific user",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"message_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the message to mark as read",
						},
						"user_id": map[string]interface{}{
							"type":        "string",
							"description": "ID of the user who read the message",
						},
					},
					"required": []string{"message_id", "user_id"},
				},
			},
		},
		{
			Type: "function",
			Function: ai2.FunctionDef{
				Name:        "messages_get_user_unread_message_count",
				Description: "Get count of unread messages for a user",
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
	}
}