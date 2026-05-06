package commands

import (
	"context"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
	"errors"
)

type (
	AddAdmin struct {
		ID        string
		Username  string
		Email     string
		Password  string
		FirstName string
		LastName  string
		Latitude  float64
		Longitude float64
		Thumbnail string
	}
	AddAdminHandler struct {
		users     domain.UserRepository
		middleman domain.MiddlemanRepository
		auth      *auth.Auth
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewAddAdminHandler(users domain.UserRepository, middleman domain.MiddlemanRepository, auth *auth.Auth, publisher ddd.EventPublisher[ddd.Event]) AddAdminHandler {
	return AddAdminHandler{
		users:     users,
		middleman: middleman,
		auth:      auth,
		publisher: publisher,
	}
}

func (h AddAdminHandler) AddAdmin(ctx context.Context, cmd AddAdmin) error {
	// Check if the requester has permission to create admins
	// This should be done at the service layer, but we can add an additional check here
	
	// Validate admin creation requirements
	if cmd.Email == "" || cmd.Password == "" || cmd.FirstName == "" || cmd.LastName == "" {
		return errors.New("email, password, first name, and last name are required for admin creation")
	}

	// Check if email already exists
	existingUser, err := h.middleman.FindByEmail(ctx, cmd.Email)
	if err == nil && existingUser != nil {
		return domain.ErrDuplicateEmail
	}

	user, err := h.users.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Generate verification token for admin
	token, err := h.auth.GeneratePasswordResetToken(user.ID(), user.Email)
	if err != nil {
		return err
	}

	// Create admin user with admin role
	event, err := user.InitAdminUser(cmd.Email, cmd.Password, cmd.Username, cmd.FirstName, cmd.LastName, cmd.Latitude, cmd.Longitude, cmd.Thumbnail, token)
	if err != nil {
		return err
	}

	err = h.users.Save(ctx, user)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}