package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type RenameUser struct {
	ID   string
	Name string
}

type RenameUserHandler struct {
	users     domain.UserRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRenameUserHandler(users domain.UserRepository, publisher ddd.EventPublisher[ddd.Event]) RenameUserHandler {
	return RenameUserHandler{
		users:     users,
		publisher: publisher,
	}
}

func (h RenameUserHandler) RenameUser(ctx context.Context, cmd RenameUser) error {
	user, err := h.users.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := user.Rename(cmd.Name)
	if err != nil {
		return err
	}

	err = h.users.Save(ctx, user)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
