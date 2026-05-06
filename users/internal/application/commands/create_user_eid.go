package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type CreateUserEID struct {
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

type CreateUserEIDHandler struct {
	users     domain.UserRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewCreateUserEIDHandler(users domain.UserRepository, publisher ddd.EventPublisher[ddd.Event]) CreateUserEIDHandler {
	return CreateUserEIDHandler{
		users:     users,
		publisher: publisher,
	}
}

// TODO implement proper handling for it
func (h CreateUserEIDHandler) CreateUserEID(ctx context.Context, cmd CreateUserEID) error {
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
