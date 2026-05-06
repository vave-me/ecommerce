package commands

import (
	"context"
	"github.com/google/uuid"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/mailer/internal/domain"
)

type CreateEmail struct {
	SenderID  string
	Recipient string
	Subject   string
	Body      string
	Status    string
}

type CreateEmailHandler struct {
	emails    domain.EmailRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCreateEmailHandler(emails domain.EmailRepository, publisher ddd.EventPublisher[ddd.Event]) CreateEmailHandler {
	return CreateEmailHandler{
		emails:    emails,
		publisher: publisher,
	}
}

func (h CreateEmailHandler) CreateEmail(ctx context.Context, cmd CreateEmail) error {

	emailID := uuid.New().String()
	emailAgg := domain.NewEmail(emailID)

	event, err := emailAgg.CreateEmail(
		emailID,       // Email ID
		cmd.SenderID,  // SenderID
		cmd.Recipient, // Recipient
		cmd.Subject,
		cmd.Body,
	)
	if err != nil {
		return errors.Wrap(err, "initializing email")
	}

	err = h.emails.Save(ctx, emailAgg)
	if err != nil {
		return errors.Wrap(err, "error saving the new email")
	}

	return errors.Wrap(h.publisher.Publish(ctx, event), "publishing domain event")
}
