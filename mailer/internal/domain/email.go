package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const EmailAggregate = "mailer.Email"

// Domain-specific errors
var (
	ErrEmailIDBlank   = errors.Wrap(errors.ErrBadRequest, "the email ID cannot be blank")
	ErrSenderIDBlank  = errors.Wrap(errors.ErrBadRequest, "the sender ID cannot be blank")
	ErrRecipientBlank = errors.Wrap(errors.ErrBadRequest, "the recipient cannot be blank")
)

// EmailStatus can represent states like Draft, Queued, Sent, Failed, etc.
type EmailStatus string

const (
	StatusDraft  EmailStatus = "Draft"
	StatusQueued EmailStatus = "Queued"
	StatusSent   EmailStatus = "Sent"
	StatusFailed EmailStatus = "Failed"
)

// Email aggregate represents an email in the system, tracked via event sourcing.
type Email struct {
	es.Aggregate

	EmailID   string
	SenderID  string
	Recipient string
	Subject   string
	Body      string
	Status    EmailStatus
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Email)(nil)

// NewEmail creates a new Email aggregate instance with the given ID.
func NewEmail(id string) *Email {
	return &Email{
		Aggregate: es.NewAggregate(id, EmailAggregate),
	}
}
func (e *Email) CreateEmail(emailID, senderID, recipient, subject, body string) (ddd.Event, error) {
	if emailID == "" {
		return nil, ErrEmailIDBlank
	}
	if senderID == "" {
		return nil, ErrSenderIDBlank
	}
	if recipient == "" {
		return nil, ErrRecipientBlank
	}

	// Default to Draft status upon creation
	status := StatusDraft

	e.AddEvent(EmailCreatedEvent, &EmailCreated{
		EmailID:   emailID,
		SenderID:  senderID,
		Recipient: recipient,
		Subject:   subject,
		Body:      body,
		Status:    status,
	})

	return ddd.NewEvent(EmailCreatedEvent, e), nil
}

// Key implements registry.Registerable
func (Email) Key() string { return EmailAggregate }

// ApplyEvent implements es.EventApplier
func (e *Email) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *EmailCreated:
		e.EmailID = payload.EmailID
		e.SenderID = payload.SenderID
		e.Recipient = payload.Recipient
		e.Subject = payload.Subject
		e.Body = payload.Body
		e.Status = payload.Status

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", e, event.EventName(), payload)
	}
	return nil
}

// ApplySnapshot applies a snapshot to rebuild the aggregate state.
func (e *Email) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *EmailV1:
		e.EmailID = ss.EmailID
		e.SenderID = ss.SenderID
		e.Recipient = ss.Recipient
		e.Subject = ss.Subject
		e.Body = ss.Body
		e.Status = ss.Status
	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", e, snapshot)
	}
	return nil
}

func (e Email) ToSnapshot() es.Snapshot {
	return &EmailV1{
		EmailID:   e.EmailID,
		SenderID:  e.SenderID,
		Recipient: e.Recipient,
		Subject:   e.Subject,
		Body:      e.Body,
		Status:    e.Status,
	}
}
