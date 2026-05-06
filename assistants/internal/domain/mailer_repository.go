package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type MailerRepository interface {
	// Core mailer operations from gRPC service
	CreateEmail(ctx context.Context, senderID, recipient, subject, body, status string) (*models.CreateEmailResponse, error)
	NotifyUserCreated(ctx context.Context, userID, username, firstname, lastname, email string, enabled bool, subject, body, status string) (*models.NotifyUserCreatedResponse, error)

	// Additional query methods for AI tooling
	GetEmail(ctx context.Context, emailID string) (*models.Email, error)
	GetEmailsBySender(ctx context.Context, senderID string, page, pageSize int32) ([]*models.Email, error)
	GetEmailsByRecipient(ctx context.Context, recipient string, page, pageSize int32) ([]*models.Email, error)
	GetEmailsByStatus(ctx context.Context, status string, page, pageSize int32) ([]*models.Email, error)
	SearchEmails(ctx context.Context, query string, page, pageSize int32) ([]*models.Email, error)
	CountEmails(ctx context.Context, senderID, recipient, status string) (int32, error)
	UpdateEmailStatus(ctx context.Context, emailID, status string) error
	DeleteEmail(ctx context.Context, emailID string) error

	// Notification and user creation handling
	GetUserCreationNotifications(ctx context.Context, userID string, page, pageSize int32) ([]*models.Email, error)
	GetEmailsInDateRange(ctx context.Context, startDate, endDate string, page, pageSize int32) ([]*models.Email, error)
	GetFailedEmails(ctx context.Context, page, pageSize int32) ([]*models.Email, error)
	GetPendingEmails(ctx context.Context, page, pageSize int32) ([]*models.Email, error)
	GetSentEmails(ctx context.Context, page, pageSize int32) ([]*models.Email, error)
}
