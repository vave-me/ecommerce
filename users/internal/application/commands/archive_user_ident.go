package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type ArchiveUserIdent struct {
	ID          string
	UserID      string
	FirstName   string
	LastName    string
	Email       string
	BirthDate   string
	City        string
	Country     string
	MobilePhone string
	Nationality string
	Custom1     string
}

type ArchiveUserIdentHandler struct {
	users     domain.UserRepository
	publisher ddd.EventPublisher[ddd.Event]
}

// TODO need to be validated for what this method will be used
func NewArchiveUserIdentHandler(users domain.UserRepository, publisher ddd.EventPublisher[ddd.Event]) ArchiveUserIdentHandler {
	return ArchiveUserIdentHandler{
		users:     users,
		publisher: publisher,
	}
}

// TODO implement proper handling for it
func (h ArchiveUserIdentHandler) ArchiveUserIdent(ctx context.Context, cmd ArchiveUserIdent) error {
	user, err := h.users.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	event, err := user.Authorize(cmd.ID, cmd.ID)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event) // Return both token and user ID
}
