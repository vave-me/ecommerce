package commands

import (
	"context"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type (
	CreateUser struct {
		ID        string
		Username  string
		Email     string
		Password  string
		FirstName string
		LastName  string
		Latitude  float64
		Longitude float64
		Thumbnail string
		Role      string
	}
	CreateUserHandler struct {
		users      domain.UserRepository
		middleman  domain.MiddlemanRepository
		auth       *auth.Auth
		publisher  ddd.EventPublisher[ddd.Event]
	}
)

func NewCreateUserHandler(users domain.UserRepository, middleman domain.MiddlemanRepository, auth *auth.Auth, publisher ddd.EventPublisher[ddd.Event]) CreateUserHandler {
	return CreateUserHandler{
		users:     users,
		middleman: middleman,
		auth:      auth,
		publisher: publisher,
	}
}

func (h CreateUserHandler) CreateUser(ctx context.Context, cmd CreateUser) error {
	// Check if email already exists
	existingUser, err := h.middleman.FindByEmail(ctx, cmd.Email)
	if err == nil && existingUser != nil {
		return domain.ErrDuplicateEmail
	}

	user, err := h.users.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}
	token, err := h.auth.GeneratePasswordResetToken(user.ID(), cmd.Email)

	event, err := user.InitUser(cmd.Email, cmd.Password, cmd.Username, cmd.FirstName, cmd.LastName, cmd.Latitude, cmd.Longitude, cmd.Thumbnail, token, domain.UserRoleCustomer)
	if err != nil {
		return err
	}
	err = h.users.Save(ctx, user)
	if err != nil {
		return err
	} // Log success
	return h.publisher.Publish(ctx, event)
}
