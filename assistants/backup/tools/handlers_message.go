package tools

import (
	"context"
	"fmt"
)

// ==================== MESSAGE HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeMessageHandlers() {
	// Conversation Management
	r.handlers["messages_create_new_conversation"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		senderID := getStringParam(params, "sender_id")
		recipientID := getStringParam(params, "recipient_id")
		itemID := getStringParam(params, "item_id")

		// Validate required parameters
		if err := ValidateIDParam("sender_id", senderID); err != nil {
			return nil, fmt.Errorf("invalid sender_id: %w", err)
		}
		if err := ValidateIDParam("recipient_id", recipientID); err != nil {
			return nil, fmt.Errorf("invalid recipient_id: %w", err)
		}
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, fmt.Errorf("invalid item_id: %w", err)
		}

		return reg.messageRepo.CreateNewConversation(ctx, senderID, recipientID, itemID)
	}

	r.handlers["messages_restore_archived_conversation"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		conversationID := getStringParam(params, "conversation_id")
		if err := ValidateIDParam("conversation_id", conversationID); err != nil {
			return nil, fmt.Errorf("invalid conversation_id: %w", err)
		}
		return reg.messageRepo.RestoreArchivedConversation(ctx, conversationID)
	}

	r.handlers["messages_archive_existing_conversation"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		conversationID := getStringParam(params, "conversation_id")
		if err := ValidateIDParam("conversation_id", conversationID); err != nil {
			return nil, fmt.Errorf("invalid conversation_id: %w", err)
		}
		return reg.messageRepo.ArchiveExistingConversation(ctx, conversationID)
	}

	r.handlers["messages_get_conversation_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		conversationID := getStringParam(params, "conversation_id")
		if err := ValidateIDParam("conversation_id", conversationID); err != nil {
			return nil, fmt.Errorf("invalid conversation_id: %w", err)
		}
		return reg.messageRepo.GetConversationByID(ctx, conversationID)
	}

	r.handlers["messages_find_conversation_by_recipient_and_item"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		recipientID := getStringParam(params, "recipient_id")
		itemID := getStringParam(params, "item_id")
		if err := ValidateIDParam("recipient_id", recipientID); err != nil {
			return nil, fmt.Errorf("invalid recipient_id: %w", err)
		}
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, fmt.Errorf("invalid item_id: %w", err)
		}
		return reg.messageRepo.FindConversationByRecipientAndItem(ctx, recipientID, itemID)
	}

	r.handlers["messages_get_user_conversation_list"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		page := getInt64Param(params, "page", 1)
		limit := getInt64Param(params, "limit", 50)
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if err := ValidatePaginationParams(page, limit); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.messageRepo.GetUserConversationList(ctx, userID, page, limit)
	}

	r.handlers["messages_get_user_active_conversations"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		page := getInt64Param(params, "page", 1)
		limit := getInt64Param(params, "limit", 50)
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if err := ValidatePaginationParams(page, limit); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.messageRepo.GetUserActiveConversations(ctx, userID, page, limit)
	}

	// Message Management
	r.handlers["messages_send_new_message"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		messageID := getStringParam(params, "message_id")
		conversationID := getStringParam(params, "conversation_id")
		senderID := getStringParam(params, "sender_id")
		recipientID := getStringParam(params, "recipient_id")
		itemID := getStringParam(params, "item_id")
		body := getStringParam(params, "body")
		isRead := getBoolParam(params, "is_read", false)

		// Validate required parameters
		if err := ValidateIDParam("message_id", messageID); err != nil {
			return nil, fmt.Errorf("invalid message_id: %w", err)
		}
		if err := ValidateIDParam("conversation_id", conversationID); err != nil {
			return nil, fmt.Errorf("invalid conversation_id: %w", err)
		}
		if err := ValidateIDParam("sender_id", senderID); err != nil {
			return nil, fmt.Errorf("invalid sender_id: %w", err)
		}
		if err := ValidateIDParam("recipient_id", recipientID); err != nil {
			return nil, fmt.Errorf("invalid recipient_id: %w", err)
		}
		if err := ValidateIDParam("item_id", itemID); err != nil {
			return nil, fmt.Errorf("invalid item_id: %w", err)
		}
		if body == "" {
			return nil, fmt.Errorf("message body is required")
		}

		// Sanitize message body
		body = SanitizeString(body)

		return reg.messageRepo.SendNewMessage(ctx, messageID, conversationID, senderID, recipientID, itemID, body, isRead)
	}

	r.handlers["messages_delete_message_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		messageID := getStringParam(params, "message_id")
		if err := ValidateIDParam("message_id", messageID); err != nil {
			return nil, fmt.Errorf("invalid message_id: %w", err)
		}
		return nil, reg.messageRepo.DeleteMessageByID(ctx, messageID)
	}

	r.handlers["messages_get_message_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		messageID := getStringParam(params, "message_id")
		if err := ValidateIDParam("message_id", messageID); err != nil {
			return nil, fmt.Errorf("invalid message_id: %w", err)
		}
		return reg.messageRepo.GetMessageByID(ctx, messageID)
	}

	r.handlers["messages_get_conversation_messages"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		conversationID := getStringParam(params, "conversation_id")
		page := getInt64Param(params, "page", 1)
		limit := getInt64Param(params, "limit", 50)
		if err := ValidateIDParam("conversation_id", conversationID); err != nil {
			return nil, fmt.Errorf("invalid conversation_id: %w", err)
		}
		if err := ValidatePaginationParams(page, limit); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.messageRepo.GetConversationMessages(ctx, conversationID, page, limit)
	}

	r.handlers["messages_mark_message_as_read_by_user"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		messageID := getStringParam(params, "message_id")
		userID := getStringParam(params, "user_id")
		if err := ValidateIDParam("message_id", messageID); err != nil {
			return nil, fmt.Errorf("invalid message_id: %w", err)
		}
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return nil, reg.messageRepo.MarkMessageAsReadByUser(ctx, messageID, userID)
	}

	r.handlers["messages_get_user_unread_message_count"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return reg.messageRepo.GetUserUnreadMessageCount(ctx, userID)
	}
}