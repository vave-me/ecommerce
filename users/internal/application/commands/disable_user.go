package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type DisableUser struct {
	ID string
}

type DisableUserHandler struct {
	users     domain.UserRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDisableUserHandler(users domain.UserRepository, publisher ddd.EventPublisher[ddd.Event]) DisableUserHandler {
	return DisableUserHandler{
		users:     users,
		publisher: publisher,
	}
}

func (h DisableUserHandler) DisableUser(ctx context.Context, cmd DisableUser) error {
	user, err := h.users.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := user.Disable()
	if err != nil {
		return err
	}

	err = h.users.Save(ctx, user)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}