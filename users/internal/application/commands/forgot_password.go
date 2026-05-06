package commands

import (
	"context"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
	"time"

	"github.com/stackus/errors"
)

type ForgotPassword struct {
	MiddlemanUserID string
	Email           string
}

type ForgotPasswordHandler struct {
	users           domain.UserRepository
	middlemanRepo   domain.MiddlemanRepository
	publisher       ddd.EventPublisher[ddd.Event]
	auth            *auth.Auth
	tokenExpiration time.Duration // How long the reset token should be valid
}

func NewForgotPasswordHandler(
	users domain.UserRepository,
	middlemanRepo domain.MiddlemanRepository,
	publisher ddd.EventPublisher[ddd.Event],
	auth *auth.Auth,
) ForgotPasswordHandler {
	return ForgotPasswordHandler{
		users:           users,
		middlemanRepo:   middlemanRepo,
		publisher:       publisher,
		auth:            auth,
		tokenExpiration: 24 * time.Hour, // Default 24-hour expiration
	}
}

func (h ForgotPasswordHandler) ForgotPassword(ctx context.Context, cmd ForgotPassword) error {
	// We should use the MiddlemanUserID which is already loaded in the server handler
	user, err := h.users.Load(ctx, cmd.MiddlemanUserID)
	if err != nil {
		return err
	}

	// Safety check: make sure auth service is available
	if h.auth == nil {
		return errors.Wrap(errors.ErrInternal, "auth service not available")
	}

	// Generate a secure token
	token, err := h.auth.GeneratePasswordResetToken(user.ID(), user.Email)
	if err != nil {
		return err
	}

	// Set the token expiration time
	expirationTime := time.Now().Add(h.tokenExpiration)

	// Request a password reset with the token and expiration time
	event, err := user.RequestPasswordReset(token, cmd.Email, expirationTime)
	if err != nil {
		return err
	}

	// Save the updated user with the reset token
	if err = h.users.Save(ctx, user); err != nil {
		return err
	}

	// Publish the event that can trigger email sending
	return h.publisher.Publish(ctx, event)
}
