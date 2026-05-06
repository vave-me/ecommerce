package tools

import (
	"context"
	"fmt"
)

// ==================== MAILER HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeMailerHandlers() {
	r.handlers["mailer_create_email"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		senderID := getStringParam(params, "sender_id")
		recipient := getStringParam(params, "recipient")
		subject := getStringParam(params, "subject")
		body := getStringParam(params, "body")
		status := getStringParam(params, "status", "pending")

		// Validate required parameters
		if err := ValidateIDParam("sender_id", senderID); err != nil {
			return nil, fmt.Errorf("invalid sender_id: %w", err)
		}
		if err := ValidateEmailParam(recipient); err != nil {
			return nil, fmt.Errorf("invalid recipient email: %w", err)
		}
		if subject == "" {
			return nil, fmt.Errorf("subject is required")
		}
		if body == "" {
			return nil, fmt.Errorf("body is required")
		}

		// Sanitize string inputs
		subject = SanitizeString(subject)
		body = SanitizeString(body)

		return reg.mailerRepo.CreateEmail(ctx, senderID, recipient, subject, body, status)
	}

	r.handlers["mailer_notify_user_created"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		username := getStringParam(params, "username")
		firstname := getStringParam(params, "firstname")
		lastname := getStringParam(params, "lastname")
		email := getStringParam(params, "email")
		enabled := getBoolParam(params, "enabled", true)
		subject := getStringParam(params, "subject")
		body := getStringParam(params, "body")
		status := getStringParam(params, "status", "pending")

		// Validate required parameters
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if err := ValidateEmailParam(email); err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}
		if username == "" {
			return nil, fmt.Errorf("username is required")
		}
		if firstname == "" {
			return nil, fmt.Errorf("firstname is required")
		}
		if lastname == "" {
			return nil, fmt.Errorf("lastname is required")
		}
		if subject == "" {
			return nil, fmt.Errorf("subject is required")
		}
		if body == "" {
			return nil, fmt.Errorf("body is required")
		}

		// Sanitize string inputs
		username = SanitizeString(username)
		firstname = SanitizeString(firstname)
		lastname = SanitizeString(lastname)
		subject = SanitizeString(subject)
		body = SanitizeString(body)

		return reg.mailerRepo.NotifyUserCreated(ctx, userID, username, firstname, lastname, email, enabled, subject, body, status)
	}

	r.handlers["mailer_get_email"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		emailID := getStringParam(params, "email_id")
		if err := ValidateIDParam("email_id", emailID); err != nil {
			return nil, fmt.Errorf("invalid email_id: %w", err)
		}
		return reg.mailerRepo.GetEmail(ctx, emailID)
	}

	r.handlers["mailer_get_emails_by_sender"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		senderID := getStringParam(params, "sender_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if err := ValidateIDParam("sender_id", senderID); err != nil {
			return nil, fmt.Errorf("invalid sender_id: %w", err)
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.mailerRepo.GetEmailsBySender(ctx, senderID, int32(page), int32(pageSize))
	}

	r.handlers["mailer_get_emails_by_recipient"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		recipient := getStringParam(params, "recipient")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if err := ValidateEmailParam(recipient); err != nil {
			return nil, fmt.Errorf("invalid recipient email: %w", err)
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.mailerRepo.GetEmailsByRecipient(ctx, recipient, int32(page), int32(pageSize))
	}

	r.handlers["mailer_get_emails_by_status"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		status := getStringParam(params, "status")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if status == "" {
			return nil, fmt.Errorf("status is required")
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.mailerRepo.GetEmailsByStatus(ctx, status, int32(page), int32(pageSize))
	}

	r.handlers["mailer_search_emails"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		query := getStringParam(params, "query")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if query == "" {
			return nil, fmt.Errorf("search query is required")
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		query = SanitizeString(query)
		return reg.mailerRepo.SearchEmails(ctx, query, int32(page), int32(pageSize))
	}

	r.handlers["mailer_count_emails"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		senderID := getStringParam(params, "sender_id")
		recipient := getStringParam(params, "recipient")
		status := getStringParam(params, "status")
		count, err := reg.mailerRepo.CountEmails(ctx, senderID, recipient, status)
		if err != nil {
			return nil, err
		}
		return map[string]int32{"count": count}, nil
	}

	r.handlers["mailer_update_email_status"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		emailID := getStringParam(params, "email_id")
		status := getStringParam(params, "status")
		if err := ValidateIDParam("email_id", emailID); err != nil {
			return nil, fmt.Errorf("invalid email_id: %w", err)
		}
		if status == "" {
			return nil, fmt.Errorf("status is required")
		}
		return nil, reg.mailerRepo.UpdateEmailStatus(ctx, emailID, status)
	}

	r.handlers["mailer_delete_email"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		emailID := getStringParam(params, "email_id")
		if err := ValidateIDParam("email_id", emailID); err != nil {
			return nil, fmt.Errorf("invalid email_id: %w", err)
		}
		return nil, reg.mailerRepo.DeleteEmail(ctx, emailID)
	}

	r.handlers["mailer_get_user_creation_notifications"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.mailerRepo.GetUserCreationNotifications(ctx, userID, int32(page), int32(pageSize))
	}

	r.handlers["mailer_get_emails_in_date_range"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		startDate := getStringParam(params, "start_date")
		endDate := getStringParam(params, "end_date")
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if startDate == "" {
			return nil, fmt.Errorf("start_date is required")
		}
		if endDate == "" {
			return nil, fmt.Errorf("end_date is required")
		}
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.mailerRepo.GetEmailsInDateRange(ctx, startDate, endDate, int32(page), int32(pageSize))
	}

	r.handlers["mailer_get_failed_emails"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.mailerRepo.GetFailedEmails(ctx, int32(page), int32(pageSize))
	}

	r.handlers["mailer_get_pending_emails"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.mailerRepo.GetPendingEmails(ctx, int32(page), int32(pageSize))
	}

	r.handlers["mailer_get_sent_emails"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		page := getInt64Param(params, "page", 1)
		pageSize := getInt64Param(params, "page_size", 20)
		if err := ValidatePaginationParams(page, pageSize); err != nil {
			return nil, fmt.Errorf("invalid pagination parameters: %w", err)
		}
		return reg.mailerRepo.GetSentEmails(ctx, int32(page), int32(pageSize))
	}
}