package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type AuthorizeUser struct {
	ID    string
	Token string
}

type AuthorizeUserHandler struct {
	users     domain.UserRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewAuthorizeUserHandler(users domain.UserRepository, publisher ddd.EventPublisher[ddd.Event]) AuthorizeUserHandler {
	return AuthorizeUserHandler{
		users:     users,
		publisher: publisher,
	}
}

func (h AuthorizeUserHandler) AuthorizeUser(ctx context.Context, cmd AuthorizeUser) error {
	user, err := h.users.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	event, err := user.Authorize(cmd.ID, cmd.Token)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event) // Return both token and user ID
}
