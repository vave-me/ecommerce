package commands

import (
	"context"
	"errors"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type EnableUser struct {
	ID                string
	VerificationToken string // Optional verification token
}

type EnableUserHandler struct {
	users     domain.UserRepository
	publisher ddd.EventPublisher[ddd.Event]
	auth      *auth.Auth // Added auth for token validation
}

func NewEnableUserHandler(users domain.UserRepository, publisher ddd.EventPublisher[ddd.Event], auth *auth.Auth) EnableUserHandler {
	return EnableUserHandler{
		users:     users,
		publisher: publisher,
		auth:      auth,
	}
}

func (h EnableUserHandler) EnableUser(ctx context.Context, cmd EnableUser) error {
	user, err := h.users.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// If a verification token is provided, validate it
	if cmd.VerificationToken != "" {
		// Validate the token
		userID, err := h.auth.ValidatePasswordResetToken(cmd.VerificationToken)
		if err != nil {
			return errors.New("invalid verification token")
		}

		// Ensure the token belongs to this user
		if userID != cmd.ID {
			return errors.New("verification token does not match user")
		}
	}

	event, err := user.Enable()
	if err != nil {
		return err
	}

	err = h.users.Save(ctx, user)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}