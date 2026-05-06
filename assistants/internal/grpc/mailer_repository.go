package grpc

import (
	"context"
	"fmt"
	"middleman/assistants/internal/domain"
	"middleman/assistants/internal/models"
	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/mailer/mailerpb"
	"time"

	"google.golang.org/grpc"
)

// MailerRepository calls the remote mailer service (gRPC).
type MailerRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.MailerRepository = (*MailerRepository)(nil)

// NewMailerRepositoryWithAuth creates a new MailerRepository with JWT authentication support
func NewMailerRepository(endpoint string, authInstance *auth.Auth) MailerRepository {
	return MailerRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// Core mailer operations

// CreateEmail creates a new email
func (r MailerRepository) CreateEmail(ctx context.Context, senderID, recipient, subject, body, status string) (*models.CreateEmailResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mailerpb.NewMailerServiceClient(conn)
	resp, err := client.CreateEmail(ctx, &mailerpb.CreateEmailRequest{
		SenderId:  senderID,
		Recipient: recipient,
		Subject:   subject,
		Body:      body,
		Status:    status,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateEmail RPC failed: %w", err)
	}

	return &models.CreateEmailResponse{
		ID: resp.GetId(),
	}, nil
}

// NotifyUserCreated sends a user creation notification email
func (r MailerRepository) NotifyUserCreated(ctx context.Context, userID, username, firstname, lastname, email string, enabled bool, subject, body, status string) (*models.NotifyUserCreatedResponse, error) {
	conn, err := r.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := mailerpb.NewMailerServiceClient(conn)
	resp, err := client.NotifyUserCreated(ctx, &mailerpb.NotifyUserCreatedRequest{
		UserId:    userID,
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Enabled:   enabled,
		Subject:   subject,
		Body:      body,
		Status:    status,
	})
	if err != nil {
		return nil, fmt.Errorf("NotifyUserCreated RPC failed: %w", err)
	}

	return &models.NotifyUserCreatedResponse{
		ID: resp.GetId(),
	}, nil
}

// Additional query methods for AI tooling
// Note: These methods assume additional RPC methods exist in the mailer service
// In a real implementation, these would need to be added to the protobuf definition

// GetEmail retrieves a single email by ID
func (r MailerRepository) GetEmail(ctx context.Context, emailID string) (*models.Email, error) {
	// TODO: This would require a GetEmail RPC method in the protobuf definition
	// For now, returning a placeholder implementation
	return nil, fmt.Errorf("GetEmail not implemented - requires additional RPC method in mailer service")
}

// GetEmailsBySender retrieves emails by sender ID with pagination
func (r MailerRepository) GetEmailsBySender(ctx context.Context, senderID string, page, pageSize int32) ([]*models.Email, error) {
	// TODO: This would require a GetEmailsBySender RPC method in the protobuf definition
	return nil, fmt.Errorf("GetEmailsBySender not implemented - requires additional RPC method in mailer service")
}

// GetEmailsByRecipient retrieves emails by recipient with pagination
func (r MailerRepository) GetEmailsByRecipient(ctx context.Context, recipient string, page, pageSize int32) ([]*models.Email, error) {
	// TODO: This would require a GetEmailsByRecipient RPC method in the protobuf definition
	return nil, fmt.Errorf("GetEmailsByRecipient not implemented - requires additional RPC method in mailer service")
}

// GetEmailsByStatus retrieves emails by status with pagination
func (r MailerRepository) GetEmailsByStatus(ctx context.Context, status string, page, pageSize int32) ([]*models.Email, error) {
	// TODO: This would require a GetEmailsByStatus RPC method in the protobuf definition
	return nil, fmt.Errorf("GetEmailsByStatus not implemented - requires additional RPC method in mailer service")
}

// SearchEmails searches emails by query with pagination
func (r MailerRepository) SearchEmails(ctx context.Context, query string, page, pageSize int32) ([]*models.Email, error) {
	// TODO: This would require a SearchEmails RPC method in the protobuf definition
	return nil, fmt.Errorf("SearchEmails not implemented - requires additional RPC method in mailer service")
}

// CountEmails counts emails based on filters
func (r MailerRepository) CountEmails(ctx context.Context, senderID, recipient, status string) (int32, error) {
	// TODO: This would require a CountEmails RPC method in the protobuf definition
	return 0, fmt.Errorf("CountEmails not implemented - requires additional RPC method in mailer service")
}

// UpdateEmailStatus updates the status of an email
func (r MailerRepository) UpdateEmailStatus(ctx context.Context, emailID, status string) error {
	// TODO: This would require an UpdateEmailStatus RPC method in the protobuf definition
	return fmt.Errorf("UpdateEmailStatus not implemented - requires additional RPC method in mailer service")
}

// DeleteEmail deletes an email
func (r MailerRepository) DeleteEmail(ctx context.Context, emailID string) error {
	// TODO: This would require a DeleteEmail RPC method in the protobuf definition
	return fmt.Errorf("DeleteEmail not implemented - requires additional RPC method in mailer service")
}

// Notification and user creation handling

// GetUserCreationNotifications retrieves user creation notification emails
func (r MailerRepository) GetUserCreationNotifications(ctx context.Context, userID string, page, pageSize int32) ([]*models.Email, error) {
	// TODO: This would require a GetUserCreationNotifications RPC method in the protobuf definition
	return nil, fmt.Errorf("GetUserCreationNotifications not implemented - requires additional RPC method in mailer service")
}

// GetEmailsInDateRange retrieves emails within a date range
func (r MailerRepository) GetEmailsInDateRange(ctx context.Context, startDate, endDate string, page, pageSize int32) ([]*models.Email, error) {
	// TODO: This would require a GetEmailsInDateRange RPC method in the protobuf definition
	return nil, fmt.Errorf("GetEmailsInDateRange not implemented - requires additional RPC method in mailer service")
}

// GetFailedEmails retrieves failed emails
func (r MailerRepository) GetFailedEmails(ctx context.Context, page, pageSize int32) ([]*models.Email, error) {
	return r.GetEmailsByStatus(ctx, models.EmailStatusFailed, page, pageSize)
}

// GetPendingEmails retrieves pending emails
func (r MailerRepository) GetPendingEmails(ctx context.Context, page, pageSize int32) ([]*models.Email, error) {
	return r.GetEmailsByStatus(ctx, models.EmailStatusPending, page, pageSize)
}

// GetSentEmails retrieves sent emails
func (r MailerRepository) GetSentEmails(ctx context.Context, page, pageSize int32) ([]*models.Email, error) {
	return r.GetEmailsByStatus(ctx, models.EmailStatusSent, page, pageSize)
}

// Helper methods

// convertEmailFromPb converts protobuf Email to domain model
func (r MailerRepository) convertEmailFromPb(pbEmail *mailerpb.Email) *models.Email {
	if pbEmail == nil {
		return nil
	}

	return &models.Email{
		ID:        pbEmail.GetId(),
		SenderID:  pbEmail.GetSenderId(),
		Recipient: pbEmail.GetRecipient(),
		Subject:   pbEmail.GetSubject(),
		Body:      pbEmail.GetBody(),
		Status:    pbEmail.GetStatus(),
		CreatedAt: time.Now(), // Would be set from protobuf timestamp if available
		UpdatedAt: time.Now(),
	}
}

// dial establishes a gRPC connection to the mailer service
// dial sets up a gRPC connection with the microservice endpoint
func (r MailerRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

func (r MailerRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}
