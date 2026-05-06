// grpc/mailer_service_server.go
package grpc

import (
	"context"
	"fmt"
	"log"
	"middleman/internal/di"
	"middleman/internal/errorsotel"
	"middleman/mailer/internal/application"
	"middleman/mailer/internal/constants"
	"middleman/mailer/mailerpb"
	"os"
	"path/filepath"

	"gopkg.in/mail.v2"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// server implements mailerpb.MailerServiceServer
type server struct {
	app application.App
	mailerpb.UnimplementedMailerServiceServer
}

// Ensure server implements mailerpb.MailerServiceServer
var _ mailerpb.MailerServiceServer = (*server)(nil)

// RegisterServer registers the MailerServiceServer with the provided gRPC registrar
func RegisterServer(ctx context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	mailerpb.RegisterMailerServiceServer(registrar, &server{
		app: app,
	})
	log.Println("Mailer service server registered successfully")
	return nil
}

// contentTypeForMailerType returns the appropriate content type based on the mailer type.
func contentTypeForMailerType(mailerType string) string {
	switch mailerType {
	case "image":
		return "image/jpeg"
	case "video":
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}

func (s *server) NotifyUserCreated(ctx context.Context, request *mailerpb.NotifyUserCreatedRequest) (*mailerpb.NotifyUserCreatedResponse, error) {
	fmt.Println("[DEBUG-GRPC] NotifyUserCreated gRPC endpoint called")
	span := trace.SpanFromContext(ctx)

	// Debug environment info
	currentDir, _ := os.Getwd()
	fmt.Printf("[DEBUG-GRPC] Current working directory: %s\n", currentDir)

	// Debug template location
	templateLocations := []string{
		filepath.Join(currentDir, "templates"),
		filepath.Join(currentDir, "middleman", "templates"),
		"./templates",
		"./middleman/templates",
	}

	fmt.Println("[DEBUG-GRPC] Checking template locations:")
	for _, location := range templateLocations {
		htmlPath := filepath.Join(location, "confirm_registration.html")
		_, err := os.Stat(htmlPath)
		exists := !os.IsNotExist(err)
		fmt.Printf("[DEBUG-GRPC] - %s exists: %v\n", htmlPath, exists)
	}

	// Debug SMTP dialer
	fmt.Println("[DEBUG-GRPC] Getting SMTP dialer from container")
	smtpDialerInterface := di.Get(ctx, constants.SmtpDialer)
	if smtpDialerInterface == nil {
		fmt.Println("[ERROR-GRPC] SMTP dialer not found in container")
		return nil, fmt.Errorf("SMTP dialer not found")
	}

	smtpDialer, ok := smtpDialerInterface.(*mail.Dialer)
	if !ok {
		fmt.Println("[ERROR-GRPC] Invalid SMTP dialer type")
		return nil, fmt.Errorf("Invalid SMTP dialer type")
	}

	fmt.Printf("[DEBUG-GRPC] SMTP dialer config: Host=%s, Port=%d, Username=%s\n",
		smtpDialer.Host, smtpDialer.Port, smtpDialer.Username)

	// Test simple email sending
	from := "redacted-email@example.com"
	to := request.GetEmail()
	subject := "Test Email"
	body := "Hello from Brevo SMTP!"

	fmt.Printf("[DEBUG-GRPC] Sending test email to: %s\n", to)
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	if err := smtpDialer.DialAndSend(m); err != nil {
		fmt.Printf("[ERROR-GRPC] Failed to send test email via SMTP: %v\n", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	fmt.Println("[DEBUG-GRPC] Test email sent successfully")

	mailerID := uuid.New().String()
	fmt.Printf("[DEBUG-GRPC] Generated mailerID: %s\n", mailerID)

	fmt.Println("[DEBUG-GRPC] Calling application.NotifyUserCreated")
	fmt.Printf("[DEBUG-GRPC] Request details: UserID=%s, Username=%s, Email=%s\n",
		request.GetUserId(), request.GetUsername(), request.GetEmail())

	err := s.app.NotifyUserCreated(ctx, application.UserCreated{
		UserID:            request.GetUserId(),
		UserName:          request.GetUsername(),
		FirstName:         request.GetFirstname(),
		LastName:          request.GetLastname(),
		Email:             request.GetEmail(),
		VerificationToken: "test-token-123", // Add verification token
	})
	span.SetAttributes(
		attribute.String("MailerID", mailerID),
	)

	if err != nil {
		fmt.Printf("[ERROR-GRPC] Application.NotifyUserCreated failed: %v\n", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	fmt.Println("[DEBUG-GRPC] Application.NotifyUserCreated completed successfully")

	return &mailerpb.NotifyUserCreatedResponse{
		Id: mailerID,
	}, nil
}

// CreateEmail handles the creation of a new Mailer entity
func (s *server) CreateEmail(ctx context.Context, request *mailerpb.CreateEmailRequest) (*mailerpb.CreateEmailResponse, error) {
	span := trace.SpanFromContext(ctx)

	// SMTP dialer from DI
	smtpDialer := di.Get(ctx, constants.SmtpDialer).(*mail.Dialer)

	from := "redacted-email@example.com"
	to := "redacted-email@example.com"
	subject := "Test Email"
	body := "Hello from Brevo SMTP!"

	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	if err := smtpDialer.DialAndSend(m); err != nil {
		fmt.Printf("ERROR sending email via SMTP: %v\n", err)
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	mailerID := uuid.New().String()

	span.SetAttributes(
		attribute.String("MailerID", mailerID),
	)

	return &mailerpb.CreateEmailResponse{
		Id: mailerID,
	}, nil
}
