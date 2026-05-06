package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// MailerToolService handles all email and mailing operations
type MailerToolService struct {
	mailer domain.MailerRepository
}

// NewMailerToolService creates a new mailer tool service
func NewMailerToolService(mailer domain.MailerRepository) *MailerToolService {
	return &MailerToolService{
		mailer: mailer,
	}
}

// GetSupportedOperations returns all operations supported by this service
func (m *MailerToolService) GetSupportedOperations() []string {
	return []string{
		"create_email", "send_email",
		"notify_user_created",
		"get_email",
		"get_emails_by_sender",
		"search_emails",
	}
}

// ExecuteOperation executes a mailer operation with streaming progress updates
func (m *MailerToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	log.Printf("MailerToolService.ExecuteOperation: Executing mailer operation: %s", operation)

	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "mailer_operation",
		Status:   "progress",
		Progress: 10,
		Metadata: map[string]interface{}{
			"operation": operation,
			"message":   fmt.Sprintf("Processing mailer operation: %s", operation),
		},
		Timestamp: time.Now().Unix(),
	}

	// Extract common parameters
	email := getStringParam(parameters, "email", "")
	subject := getStringParam(parameters, "subject", "")
	message := getStringParam(parameters, "message", "")

	if email == "" {
		return nil, fmt.Errorf("email parameter required")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "mailer_operation",
		Status:   "progress",
		Progress: 50,
		Metadata: map[string]interface{}{
			"step":    "processing_email",
			"email":   email,
			"subject": subject,
		},
		Timestamp: time.Now().Unix(),
	}

	var result interface{}
	var err error

	switch operation {
	case "create_email", "send_email":
		result, err = m.mailer.CreateEmail(ctx,
			getStringParam(parameters, "sender_id", ""),
			email, subject, message,
			getStringParam(parameters, "status", "pending"))

	case "notify_user_created":
		result, err = m.mailer.NotifyUserCreated(ctx,
			getStringParam(parameters, "user_id", ""),
			getStringParam(parameters, "username", ""),
			getStringParam(parameters, "firstname", ""),
			getStringParam(parameters, "lastname", ""),
			email,
			getBoolParam(parameters, "enabled", true),
			subject, message,
			getStringParam(parameters, "status", "pending"))

	case "get_email":
		emailID := getStringParam(parameters, "email_id", "")
		if emailID == "" {
			return nil, fmt.Errorf("email_id parameter required")
		}
		result, err = m.mailer.GetEmail(ctx, emailID)

	case "get_emails_by_sender":
		senderID := getStringParam(parameters, "sender_id", "")
		if senderID == "" {
			return nil, fmt.Errorf("sender_id parameter required")
		}
		result, err = m.mailer.GetEmailsBySender(ctx, senderID,
			int32(getInt64Param(parameters, "page", 1)),
			int32(getInt64Param(parameters, "page_size", 20)))

	case "search_emails":
		query := getStringParam(parameters, "query", "")
		if query == "" {
			return nil, fmt.Errorf("query parameter required")
		}
		result, err = m.mailer.SearchEmails(ctx, query,
			int32(getInt64Param(parameters, "page", 1)),
			int32(getInt64Param(parameters, "page_size", 20)))

	default:
		return nil, fmt.Errorf("unsupported mailer operation: %s", operation)
	}

	if err != nil {
		// Send error update
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "mailer_operation",
			Status:   "error",
			Error:    err.Error(),
			Metadata: map[string]interface{}{
				"operation": operation,
			},
			Timestamp: time.Now().Unix(),
		}
		return nil, fmt.Errorf("mailer operation failed: %w", err)
	}

	// Send completion update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "mailer_operation",
		Status:   "completed",
		Progress: 100,
		Result:   result,
		Metadata: map[string]interface{}{
			"operation": operation,
			"success":   true,
		},
		Timestamp: time.Now().Unix(),
	}

	return map[string]interface{}{
		"entity_type": "mailer",
		"operation":   operation,
		"result":      result,
		"success":     true,
	}, nil
}
