package commands

import (
	"context"
	"log"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type LogoutUser struct {
	UserID       string
	AuthToken    string
	RefreshToken string
}

type LogoutUserHandler struct {
	users     domain.UserRepository
	auth      *auth.Auth
	publisher ddd.EventPublisher[ddd.Event]
}

func NewLogoutUserHandler(
	users domain.UserRepository,
	auth *auth.Auth,
	publisher ddd.EventPublisher[ddd.Event],
) LogoutUserHandler {
	return LogoutUserHandler{
		users:     users,
		auth:      auth,
		publisher: publisher,
	}
}

func (h LogoutUserHandler) LogoutUser(ctx context.Context, cmd LogoutUser) error {
	// Find the user by ID
	user, err := h.users.Load(ctx, cmd.UserID)
	if err != nil {
		log.Printf("Failed to load user %s during logout: %v", cmd.UserID, err)
		return err
	}

	// Generate domain event for logout
	logoutEvent, err := user.Logout()
	if err != nil {
		return err
	}

	// Also generate a token invalidation event if tokens are provided
	var tokenEvent ddd.Event
	if cmd.AuthToken != "" || cmd.RefreshToken != "" {
		// Use the first few characters of the refresh token as a token ID for tracking
		tokenID := "unknown"
		if cmd.RefreshToken != "" {
			// Use safer substring operation to avoid potential panic
			maxLen := 8
			if len(cmd.RefreshToken) < maxLen {
				maxLen = len(cmd.RefreshToken)
			}
			tokenID = cmd.RefreshToken[:maxLen]
		}

		tokenEvent, err = user.InvalidateTokens(tokenID, "user logout")
		if err != nil {
			log.Printf("Failed to generate token invalidation event: %v", err)
			// Continue with logout even if token invalidation fails
		}
	}

	// Save the updated user state
	err = h.users.Save(ctx, user)
	if err != nil {
		return err
	}

	// Publish the logout event
	err = h.publisher.Publish(ctx, logoutEvent)
	if err != nil {
		return err
	}

	// Also publish the token invalidation event if it was created
	if tokenEvent != nil {
		err = h.publisher.Publish(ctx, tokenEvent)
		if err != nil {
			log.Printf("Failed to publish token invalidation event: %v", err)
			// Don't fail the logout if token event publishing fails
		}
	}

	return nil
}
