package application

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"middleman/internal/di"
	"middleman/mailer/internal/constants"
	"os"
	"path/filepath"
	"time"

	"github.com/stackus/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/mail.v2"
)

type UserCreated struct {
	UserID            string
	Email             string
	UserName          string
	FirstName         string
	LastName          string
	Enabled           bool
	VerificationToken string
}

type PasswordReset struct {
	Email          string
	ResetToken     string
	ExpirationDate time.Time
}

// ConvertProtoTimestamp converts a Protocol Buffer timestamp to Go time.Time
func ConvertProtoTimestamp(protoTimestamp *timestamppb.Timestamp) time.Time {
	if protoTimestamp == nil {
		return time.Time{} // Return zero time if nil
	}
	return protoTimestamp.AsTime()
}

type App interface {
	NotifyUserCreated(ctx context.Context, notify UserCreated) error
	ResetPassword(ctx context.Context, reset PasswordReset) error
}

// AppConfig struct removed

type Application struct {
	users UserRepository
}

var _ App = (*Application)(nil)

func New(users UserRepository) *Application {
	return &Application{
		users: users,
	}
}

type confirmationEmailData struct {
	FirstName       string
	VerificationURL string
}

type emailTemplateData struct {
	FirstName       string
	VerificationURL string
	AppName         string
	SupportEmail    string
	CompanyName     string
	CurrentYear     int
}

func (a *Application) NotifyUserCreated(ctx context.Context, notify UserCreated) error {
	fmt.Printf("[DEBUG] Entering NotifyUserCreated for user: %s, email: %s\n", notify.UserID, notify.Email)

	// --- Dummy Config Values ---
	dummyAppBaseURL := "https://sfx-markt.de"
	dummyAppName := "sfx markt"
	dummySupportEmail := "redacted-email@example.com"
	dummyCompanyName := "sfx markt"
	dummyEmailFromName := "sfx markt Team"
	dummyEmailFromAddr := "redacted-email@example.com"
	// add locale to match right language

	// --- Path to templates (EXACT ORIGINAL PATH) ---
	templateDir := "./middleman/templates" // Keep the original path exactly

	currentDir, err := os.Getwd()
	fmt.Println("Current working directory:", os.DirFS(fmt.Sprintf("..%s", currentDir)))

	verificationURL := fmt.Sprintf("%s/verify?token=%s&email=%s", dummyAppBaseURL, notify.VerificationToken, notify.Email)
	emailData := confirmationEmailData{
		FirstName:       notify.FirstName,
		VerificationURL: verificationURL,
	}

	recipient := notify.Email
	subject := fmt.Sprintf("Confirm your registration for %s", dummyAppName)

	err = a.sendConfirmationEmail(
		ctx,
		recipient,
		subject,
		emailData,
		templateDir, // Pass template dir path
		dummyAppName,
		dummySupportEmail,
		dummyCompanyName,
		dummyEmailFromName,
		dummyEmailFromAddr,
	)
	if err != nil {
		fmt.Printf("Error sending confirmation email for %s: %v\n", recipient, err)
		return errors.Wrap(err, "failed to send user confirmation email")
	}

	fmt.Printf("User confirmation email sent successfully to %s\n", recipient)
	return nil
}

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Helper function to search for a file in the directory tree (with depth limit)
func searchForFile(root, filename string, maxDepth int) {
	if maxDepth <= 0 {
		return
	}

	files, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("[DEBUG] Error reading directory %s: %v\n", root, err)
		return
	}

	for _, file := range files {
		path := filepath.Join(root, file.Name())

		if file.IsDir() {
			// Skip .git, node_modules and other common large directories
			if file.Name() == ".git" || file.Name() == "node_modules" || file.Name() == "vendor" {
				continue
			}
			searchForFile(path, filename, maxDepth-1)
		} else if file.Name() == filename {
			fmt.Printf("[DEBUG] FOUND TEMPLATE: %s\n", path)
		}
	}
}

func (a *Application) sendConfirmationEmail(
	ctx context.Context,
	recipient, subject string,
	data confirmationEmailData,
	templateDir string, // Receive template dir path
	appName string,     // Receive dummy values
	supportEmail string,
	companyName string,
	emailFromName string,
	emailFromAddr string,
) error {
	htmlPath := filepath.Join(templateDir, "confirm_registration.html")
	textPath := filepath.Join(templateDir, "confirm_registration.txt")

	fmt.Printf("DEBUG: Parsing templates on demand from: %s and %s\n", htmlPath, textPath)

	htmlConfirmTmpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse HTML template file: %s", htmlPath)
	}

	textConfirmTmpl, err := template.ParseFiles(textPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse text template file: %s", textPath)
	}

	templateData := emailTemplateData{
		FirstName:       data.FirstName,
		VerificationURL: data.VerificationURL,
		AppName:         appName,      // Use passed dummy value
		SupportEmail:    supportEmail, // Use passed dummy value
		CompanyName:     companyName,  // Use passed dummy value
		CurrentYear:     time.Now().Year(),
	}

	var htmlBody bytes.Buffer
	if err := htmlConfirmTmpl.ExecuteTemplate(&htmlBody, "confirm_registration.html", templateData); err != nil {
		return errors.Wrap(err, "failed to execute HTML confirmation template by name")
	}
	htmlBodyStr := htmlBody.String()

	var textBody bytes.Buffer
	if err := textConfirmTmpl.ExecuteTemplate(&textBody, "confirm_registration.txt", templateData); err != nil {
		return errors.Wrap(err, "failed to execute text confirmation template by name")
	}
	textBodyStr := textBody.String()

	smtpDialerInterface := di.Get(ctx, constants.SmtpDialer)
	if smtpDialerInterface == nil {
		return errors.ErrNotFound.Msg("SMTP dialer not found in DI container")
	}
	smtpDialer, ok := smtpDialerInterface.(*mail.Dialer)
	if !ok {
		return errors.ErrNotFound.Msg("Invalid SMTP dialer type found in DI container")
	}

	from := emailFromAddr // Use passed dummy value
	if emailFromName != "" {
		from = fmt.Sprintf("%s <%s>", emailFromName, emailFromAddr) // Use passed dummy values
	}

	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", textBodyStr)
	m.AddAlternative("text/html", htmlBodyStr)

	if err := smtpDialer.DialAndSend(m); err != nil {
		fmt.Printf("ERROR sending email via SMTP: %v\n", err)
		return errors.Wrap(err, "SMTP SendEmail failed")
	}

	return nil
}

// Define a new type for password reset email data
type passwordResetData struct {
	FirstName string
	ResetURL  string
}

// Add ResetPassword method to Application
func (a *Application) ResetPassword(ctx context.Context, reset PasswordReset) error {
	fmt.Printf("[DEBUG] Entering ResetPassword for email: %s\n", reset.Email)

	// --- Dummy Config Values ---
	dummyAppBaseURL := "http://localhost:3000"
	//dummyAppBaseURL := "https://sfx-markt.de"
	dummyAppName := "sfx markt"
	dummySupportEmail := "redacted-email@example.com"
	dummyCompanyName := "sfx markt"
	dummyEmailFromName := "sfx markt Team"
	dummyEmailFromAddr := "redacted-email@example.com"

	// --- Path to templates (EXACT ORIGINAL PATH) ---
	templateDir := "./middleman/templates" // Keep the original path exactly

	// Get user information (first name) if available
	user, err := a.users.Find(ctx, "") // We don't have a userID, just email
	firstName := "User"                // Default if user not found or no first name available
	if err == nil && user != nil && user.FirstName != "" {
		firstName = user.FirstName
	} else {
		fmt.Printf("Using default name 'User' as couldn't find user by email: %s\n", reset.Email)
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s&email=%s", dummyAppBaseURL, reset.ResetToken, reset.Email)
	emailData := passwordResetData{
		FirstName: firstName,
		ResetURL:  resetURL,
	}

	recipient := reset.Email
	subject := fmt.Sprintf("Reset your password for %s", dummyAppName)

	err = a.sendPasswordResetEmail(
		ctx,
		recipient,
		subject,
		emailData,
		templateDir,
		dummyAppName,
		dummySupportEmail,
		dummyCompanyName,
		dummyEmailFromName,
		dummyEmailFromAddr,
	)
	if err != nil {
		fmt.Printf("Error sending password reset email for %s: %v\n", recipient, err)
		return errors.Wrap(err, "failed to send password reset email")
	}

	fmt.Printf("Password reset email sent successfully to %s\n", recipient)
	return nil
}

func (a *Application) sendPasswordResetEmail(
	ctx context.Context,
	recipient, subject string,
	data passwordResetData,
	templateDir string,
	appName string,
	supportEmail string,
	companyName string,
	emailFromName string,
	emailFromAddr string,
) error {
	htmlPath := filepath.Join(templateDir, "forgot_password.html")
	textPath := filepath.Join(templateDir, "forgot_password.txt")

	fmt.Printf("DEBUG: Parsing password reset templates from: %s and %s\n", htmlPath, textPath)

	htmlResetTmpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse HTML reset template file: %s", htmlPath)
	}

	textResetTmpl, err := template.ParseFiles(textPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse text reset template file: %s", textPath)
	}

	templateData := emailTemplateData{
		FirstName:       data.FirstName,
		VerificationURL: data.ResetURL, // Use the reset URL field
		AppName:         appName,
		SupportEmail:    supportEmail,
		CompanyName:     companyName,
		CurrentYear:     time.Now().Year(),
	}

	var htmlBody bytes.Buffer
	if err := htmlResetTmpl.ExecuteTemplate(&htmlBody, "forgot_password.html", templateData); err != nil {
		return errors.Wrap(err, "failed to execute HTML reset template by name")
	}
	htmlBodyStr := htmlBody.String()

	var textBody bytes.Buffer
	if err := textResetTmpl.ExecuteTemplate(&textBody, "forgot_password.txt", templateData); err != nil {
		return errors.Wrap(err, "failed to execute text reset template by name")
	}
	textBodyStr := textBody.String()

	smtpDialerInterface := di.Get(ctx, constants.SmtpDialer)
	if smtpDialerInterface == nil {
		return errors.ErrNotFound.Msg("SMTP dialer not found in DI container")
	}
	smtpDialer, ok := smtpDialerInterface.(*mail.Dialer)
	if !ok {
		return errors.ErrNotFound.Msg("Invalid SMTP dialer type found in DI container")
	}

	from := emailFromAddr
	if emailFromName != "" {
		from = fmt.Sprintf("%s <%s>", emailFromName, emailFromAddr)
	}

	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", textBodyStr)
	m.AddAlternative("text/html", htmlBodyStr)

	if err := smtpDialer.DialAndSend(m); err != nil {
		fmt.Printf("ERROR sending password reset email via SMTP: %v\n", err)
		return errors.Wrap(err, "SMTP SendEmail failed")
	}

	return nil
}

func PtrString(s string) *string {
	return &s
}
