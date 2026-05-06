package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"middleman/managers/internal/domain"
	"middleman/managers/internal/models"
)

// MessageToolService handles messaging operations
type MessageToolService struct {
	messageRepo domain.MessagesRepository
	config      *ServiceConfig
}

// NewMessageToolService creates a new message tool service
func NewMessageToolService(messageRepo domain.MessagesRepository) *MessageToolService {
	return &MessageToolService{
		messageRepo: messageRepo,
		config: &ServiceConfig{
			MaxRetries:      3,
			EnableStreaming: true,
			EnableMetrics:   true,
		},
	}
}

// ExecuteOperation routes message operations to appropriate handlers
func (s *MessageToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {

	// Send initial progress
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "message_operation",
			Status:   "started",
			Progress: 0,
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "MessageToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	var result interface{}
	var err error

	switch operation {
	case "send_message", "create":
		result, err = s.sendMessage(ctx, parameters, streamChan, toolID)
	case "get_message", "find":
		result, err = s.getMessage(ctx, parameters, streamChan, toolID)
	case "list_messages", "list":
		result, err = s.listMessages(ctx, parameters, streamChan, toolID)
	case "get_conversation":
		result, err = s.getConversation(ctx, parameters, streamChan, toolID)
	case "mark_as_read":
		result, err = s.markAsRead(ctx, parameters, streamChan, toolID)
	case "delete_message":
		result, err = s.deleteMessage(ctx, parameters, streamChan, toolID)
	case "search_messages":
		result, err = s.searchMessages(ctx, parameters, streamChan, toolID)
	case "get_unread_count":
		result, err = s.getUnreadCount(ctx, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported message operation: %s", operation)
	}

	// Send completion status
	if streamChan != nil {
		status := "completed"
		errorStr := ""
		if err != nil {
			status = "error"
			errorStr = err.Error()
		}

		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "message_operation",
			Status:   status,
			Progress: 100,
			Result:   result,
			Error:    errorStr,
			Metadata: map[string]interface{}{
				"operation": operation,
				"service":   "MessageToolService",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	return result, err
}

func (s *MessageToolService) sendMessage(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	// Required parameters
	senderID := getStringParam(params, "sender_id", "")
	recipientID := getStringParam(params, "recipient_id", getStringParam(params, "receiver_id", ""))
	body := getStringParam(params, "content", getStringParam(params, "body", ""))

	// Optional parameters
	conversationID := getStringParam(params, "conversation_id", "")
	itemID := getStringParam(params, "item_id", "")
	isRead := getBoolParam(params, "is_read", false)

	if senderID == "" || recipientID == "" || body == "" {
		return nil, fmt.Errorf("sender_id, recipient_id, and content/body are required")
	}

	// Generate message ID if not provided
	messageID := getStringParam(params, "message_id", "")
	if messageID == "" {
		messageID = uuid.New().String()
	}

	// Progress event
	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "message_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "sending_message",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	msgResp, err := s.messageRepo.SendMessage(ctx, messageID, conversationID, senderID, recipientID, itemID, body, isRead)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return map[string]interface{}{
		"message_id":      msgResp.ID,
		"conversation_id": conversationID,
		"sender_id":       senderID,
		"recipient_id":    recipientID,
		"item_id":         itemID,
		"content":         body,
		"sent_at":         msgResp.SentAt.Format(time.RFC3339),
		"is_read":         isRead,
	}, nil
}

func (s *MessageToolService) getMessage(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	messageID := getStringParam(params, "id", "")
	if messageID == "" {
		messageID = getStringParam(params, "message_id", "")
	}
	if messageID == "" {
		return nil, fmt.Errorf("message ID is required")
	}

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "message_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "getting_message",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	message, err := s.messageRepo.GetMessage(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to find message: %w", err)
	}

	return map[string]interface{}{
		"message": message,
		"id":      messageID,
	}, nil
}

func (s *MessageToolService) listMessages(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	conversationID := getStringParam(params, "conversation_id", getStringParam(params, "id", ""))
	if conversationID == "" {
		return nil, fmt.Errorf("conversation_id is required to list messages")
	}

	page := getInt64Param(params, "page", 1)
	limit := getInt64Param(params, "limit", 20)

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "message_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "listing_messages",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	messagesResp, err := s.messageRepo.GetMessages(ctx, conversationID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	return map[string]interface{}{
		"messages":        messagesResp.Messages,
		"conversation_id": conversationID,
		"page":            page,
		"limit":           limit,
		"total":           messagesResp.Total,
	}, nil
}

func (s *MessageToolService) getConversation(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	conversationID := getStringParam(params, "conversation_id", "")
	if conversationID == "" {
		conversationID = getStringParam(params, "id", "")
	}
	if conversationID == "" {
		return nil, fmt.Errorf("conversation ID is required")
	}

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "message_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "getting_conversation",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	conversation, err := s.messageRepo.GetConversation(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return map[string]interface{}{
		"conversation":    conversation,
		"conversation_id": conversationID,
	}, nil
}

func (s *MessageToolService) markAsRead(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	messageID := getStringParam(params, "message_id", "")
	userID := getStringParam(params, "user_id", "")

	if messageID == "" || userID == "" {
		return nil, fmt.Errorf("message_id and user_id are required")
	}

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "message_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "marking_message_as_read",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	err := s.messageRepo.MarkMessageAsRead(ctx, messageID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark message as read: %w", err)
	}

	return map[string]interface{}{
		"message_id": messageID,
		"user_id":    userID,
		"marked":     true,
	}, nil
}

func (s *MessageToolService) deleteMessage(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	messageID := getStringParam(params, "message_id", "")

	if messageID == "" {
		return nil, fmt.Errorf("message ID is required")
	}

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "message_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "deleting_message",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	err := s.messageRepo.DeleteMessage(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete message: %w", err)
	}

	return map[string]interface{}{
		"message_id": messageID,
		"deleted":    true,
	}, nil
}

func (s *MessageToolService) searchMessages(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	conversationID := getStringParam(params, "conversation_id", "")
	query := getStringParam(params, "query", "")

	if conversationID == "" || query == "" {
		return nil, fmt.Errorf("conversation_id and query are required")
	}

	page := getInt64Param(params, "page", 1)
	limit := getInt64Param(params, "limit", 20)

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "message_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "searching_messages",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	// Retrieve messages and perform basic client-side filtering until backend supports search
	messagesResp, err := s.messageRepo.GetMessages(ctx, conversationID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}

	filtered := []*models.Message{}
	for _, m := range messagesResp.Messages {
		if ContainsIgnoreCase(m.Body, query) {
			filtered = append(filtered, m)
		}
	}

	return map[string]interface{}{
		"messages":        filtered,
		"query":           query,
		"conversation_id": conversationID,
		"page":            page,
		"limit":           limit,
		"total":           int64(len(filtered)),
	}, nil
}

func (s *MessageToolService) getUnreadCount(ctx context.Context, params map[string]interface{}, streamChan chan<- ToolExecutionStream, toolID string) (interface{}, error) {
	userID := getStringParam(params, "user_id", "")

	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	if streamChan != nil {
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "message_operation",
			Status:   "progress",
			Progress: 30,
			Metadata: map[string]interface{}{
				"step": "getting_unread_count",
			},
			Timestamp: time.Now().Unix(),
		}
	}

	count, err := s.messageRepo.GetUnreadMessagesCount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread count: %w", err)
	}

	return map[string]interface{}{
		"user_id":      userID,
		"unread_count": count,
	}, nil
}
