package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type CreateUserVideoIdent struct {
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

type CreateUserVideoIdentHandler struct {
	users     domain.UserRepository
	publisher ddd.EventPublisher[ddd.Event]
}

// TODO need to be validated for what this method will be used
func NewCreateUserVideoIdentHandler(users domain.UserRepository, publisher ddd.EventPublisher[ddd.Event]) CreateUserVideoIdentHandler {
	return CreateUserVideoIdentHandler{
		users:     users,
		publisher: publisher,
	}
}

// TODO implement proper handling for it
func (h CreateUserVideoIdentHandler) CreateUserVideoIdent(ctx context.Context, cmd CreateUserVideoIdent) error {
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
