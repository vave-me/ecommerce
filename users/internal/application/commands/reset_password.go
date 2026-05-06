package commands

import (
	"context"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type ResetPassword struct {
	Token       string
	NewPassword string
}

type ResetPasswordHandler struct {
	users     domain.UserRepository
	publisher ddd.EventPublisher[ddd.Event]
	auth      *auth.Auth
}

func NewResetPasswordHandler(users domain.UserRepository, publisher ddd.EventPublisher[ddd.Event], auth *auth.Auth) ResetPasswordHandler {
	return ResetPasswordHandler{
		users:     users,
		publisher: publisher,
		auth:      auth,
	}
}

func (h ResetPasswordHandler) ResetPassword(ctx context.Context, cmd ResetPassword) error {
	// Validate the token and get the user ID
	userID, err := h.auth.ValidatePasswordResetToken(cmd.Token)

	if err != nil {
		return err
	}

	// Find the user by ID
	user, err := h.users.Load(ctx, userID)
	if err != nil {
		return err
	}

	event, err := user.ResetPassword(cmd.NewPassword)
	if err != nil {
		return err
	}

	if err = h.users.Save(ctx, user); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
