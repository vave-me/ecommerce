package models

import "time"

// Email represents an email record in the mailer system
type Email struct {
	ID        string    `json:"id" db:"id"`
	SenderID  string    `json:"sender_id" db:"sender_id"`
	Recipient string    `json:"recipient" db:"recipient"`
	Subject   string    `json:"subject" db:"subject"`
	Body      string    `json:"body" db:"body"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateEmailResponse represents the response from creating an email
type CreateEmailResponse struct {
	ID string `json:"id"`
}

// NotifyUserCreatedResponse represents the response from user creation notification
type NotifyUserCreatedResponse struct {
	ID string `json:"id"`
}

// GetEmailResponse represents the response from getting a single email
type GetEmailResponse struct {
	Email *Email `json:"email"`
}

// GetEmailsResponse represents the response from getting multiple emails
type GetEmailsResponse struct {
	Emails []*Email `json:"emails"`
	Total  int32    `json:"total"`
	Page   int32    `json:"page"`
	Size   int32    `json:"size"`
}

// SearchEmailsResponse represents the response from searching emails
type SearchEmailsResponse struct {
	Emails []*Email `json:"emails"`
	Total  int32    `json:"total"`
	Page   int32    `json:"page"`
	Size   int32    `json:"size"`
	Query  string   `json:"query"`
}

// CountEmailsResponse represents the response from counting emails
type CountEmailsResponse struct {
	Count int32 `json:"count"`
}

// UpdateEmailStatusResponse represents the response from updating email status
type UpdateEmailStatusResponse struct {
	Success bool `json:"success"`
}

// DeleteEmailResponse represents the response from deleting an email
type DeleteEmailResponse struct {
	Success bool `json:"success"`
}

// Email status constants
const (
	EmailStatusPending    = "pending"
	EmailStatusSent       = "sent"
	EmailStatusFailed     = "failed"
	EmailStatusDraftEmail = "draft"
)

// EmailStatus enumeration
type EmailStatus string

const (
	StatusPending    EmailStatus = "pending"
	StatusSent       EmailStatus = "sent"
	StatusFailed     EmailStatus = "failed"
	StatusEmailDraft EmailStatus = "draft"
)

func (s EmailStatus) String() string {
	return string(s)
}

// ToEmailStatus converts string to EmailStatus
func ToEmailStatus(status string) EmailStatus {
	switch status {
	case "pending":
		return StatusPending
	case "sent":
		return StatusSent
	case "failed":
		return StatusFailed
	case "draft":
		return StatusEmailDraft
	default:
		return StatusPending
	}
}

// IsValid returns true if the email status is valid
func (s EmailStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusSent, StatusFailed, StatusEmailDraft:
		return true
	default:
		return false
	}
}
