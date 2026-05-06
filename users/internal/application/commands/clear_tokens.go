package commands

import (
	"context"
	"errors"
	"log"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type ClearTokens struct {
	UserID       string
	TokenID      string
	RefreshToken string
	Reason       string
}

type ClearTokensHandler struct {
	users     domain.UserRepository
	auth      *auth.Auth
	publisher ddd.EventPublisher[ddd.Event]
}

func NewClearTokensHandler(
	users domain.UserRepository,
	auth *auth.Auth,
	publisher ddd.EventPublisher[ddd.Event],
) ClearTokensHandler {
	return ClearTokensHandler{
		users:     users,
		auth:      auth,
		publisher: publisher,
	}
}

// ClearTokens handles the invalidation of user tokens during logout or security events
func (h ClearTokensHandler) ClearTokens(ctx context.Context, cmd ClearTokens) error {
	// Find the user by ID
	user, err := h.users.Load(ctx, cmd.UserID)
	if err != nil {
		log.Printf("Failed to load user %s: %v", cmd.UserID, err)
		return err
	}

	// Validate the user exists
	if user == nil {
		return errors.New("user not found")
	}

	// Default reason if not provided
	reason := cmd.Reason
	if reason == "" {
		reason = "user logout"
	}

	// Use tokenID if provided, otherwise use a placeholder
	tokenID := cmd.TokenID
	if tokenID == "" {
		tokenID = "unknown"
	}

	// Generate token invalidation event
	event, err := user.InvalidateTokens(tokenID, reason)
	if err != nil {
		return err
	}

	// Save any changes to the user
	err = h.users.Save(ctx, user)
	if err != nil {
		return err
	}

	// Publish the token invalidation event
	return h.publisher.Publish(ctx, event)
}
