package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type (
	UpdateUser struct {
		ID         string
		Username   string
		Bio        string
		Privacy    domain.UserPrivacy
		Background string
		FirstName  string
		LastName   string
		Latitude   float64
		Longitude  float64
		Thumbnail  string
	}
	UpdateUserHandler struct {
		users     domain.UserRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewUpdateUserHandler(users domain.UserRepository, publisher ddd.EventPublisher[ddd.Event]) UpdateUserHandler {
	return UpdateUserHandler{
		users:     users,
		publisher: publisher,
	}
}

func (h UpdateUserHandler) UpdateUser(ctx context.Context, cmd UpdateUser) error {

	user, err := h.users.Load(ctx, cmd.ID)
	if err != nil {

		return err
	}

	event, err := user.UpdateUser(cmd.Username, cmd.Bio, cmd.Privacy.String(), cmd.FirstName, cmd.LastName, cmd.Background, cmd.Latitude, cmd.Longitude, cmd.Thumbnail)
	if err != nil {
		return err
	}
	err = h.users.Save(ctx, user)
	if err != nil {
		return err
	} // Log success
	return h.publisher.Publish(ctx, event)
}
