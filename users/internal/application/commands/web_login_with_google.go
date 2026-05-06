package commands

import (
	"context"
	"log"

	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"

	"github.com/stackus/errors"
)

// WebLoginWithGoogle command for handling Google OAuth authentication
type WebLoginWithGoogle struct {
	UserID        string // User ID in our system
	GoogleID      string // Google's subject ID (sub claim)
	Email         string // User's email from Google
	EmailVerified bool   // Whether Google has verified the email
	Enabled       bool
	FirstName     string // User's first name
	LastName      string // User's last name
	Picture       string // User's profile picture URL
}

type WebLoginWithGoogleHandler struct {
	users     domain.UserRepository
	auth      auth.Authenticator
	publisher ddd.EventPublisher[ddd.Event]
}

func NewWebLoginWithGoogleHandler(
	middleman domain.MiddlemanRepository,
	users domain.UserRepository,
	auth auth.Authenticator,
	publisher ddd.EventPublisher[ddd.Event],
) WebLoginWithGoogleHandler {
	return WebLoginWithGoogleHandler{
		users:     users,
		auth:      auth,
		publisher: publisher,
	}
}

func (h WebLoginWithGoogleHandler) WebLoginWithGoogle(ctx context.Context, cmd WebLoginWithGoogle) (string, string, string, error) {
	// Validate inputs
	if cmd.UserID == "" {
		return "", "", "", errors.Wrap(errors.ErrInvalidArgument, "user ID cannot be empty")
	}

	if cmd.GoogleID == "" {
		return "", "", "", errors.Wrap(errors.ErrInvalidArgument, "Google ID cannot be empty")
	}

	if cmd.Email == "" {
		return "", "", "", errors.Wrap(errors.ErrInvalidArgument, "email cannot be empty")
	}

	if !cmd.EmailVerified {
		return "", "", "", errors.Wrap(errors.ErrPermissionDenied, "email must be verified by Google")
	}

	// Load the user by UserID
	user, err := h.users.Load(ctx, cmd.UserID)
	if err != nil {
		log.Printf("Load user error: %v", err)
		return "", "", "", errors.Wrap(errors.ErrNotFound, "user not found")
	}

	// Verify Google ID matches or update it
	// The domain.User has GoogleId (lowercase 'd') while domain.MiddlemanUser has GoogleID (uppercase 'D')
	// We need to handle this difference properly

	// First time Google login or GoogleId not set, link the account
	if user.GoogleID == "" {
		log.Printf("Linking Google ID %s to user %s", cmd.GoogleID, user.ID())
		// Update the user's Google ID in the domain model
		user.GoogleID = cmd.GoogleID

		// Also update profile information if available
		if cmd.FirstName != "" {
			user.FirstName = cmd.FirstName
		}

		if cmd.LastName != "" {
			user.LastName = cmd.LastName
		}

		if cmd.Picture != "" {
			user.Thumbnail = cmd.Picture
		}

		// Save changes
		err = h.users.Save(ctx, user)
		if err != nil {
			log.Printf("Warning: Failed to save Google profile updates for user %s: %v", user.ID(), err)
			// Continue despite error - don't block login
		}
	} else if user.GoogleID != cmd.GoogleID {
		// Instead of blocking login, log a warning and update the GoogleId
		log.Printf("Warning: User %s logging in with different Google ID: existing %s, new %s",
			user.ID(), user.GoogleID, cmd.GoogleID)

		// Update to the new Google ID
		user.GoogleID = cmd.GoogleID
		err = h.users.Save(ctx, user)
		if err != nil {
			log.Printf("Warning: Failed to update Google ID for user %s: %v", user.ID(), err)
			// Continue despite error - don't block login
		}
	}

	// Generate JWT tokens for the user with claims from our system
	tokens, err := h.auth.GenerateTokenPair(&auth.JwtUser{
		ID:        user.ID(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Lat:       user.Lat,
		Long:      user.Lang, // Fixed: Use Lang field instead of duplicating Lat
		Role:      user.Role.String(),
	})
	if err != nil {
		log.Printf("Failed to generate tokens for user %s: %v", user.ID(), err)
		return "", "", "", errors.Wrap(err, "failed to generate authentication tokens")
	}

	// Create and publish login event
	event, err := user.Login()
	if err != nil {
		log.Printf("Failed to create login event for user %s: %v", user.ID(), err)
		// Continue despite error - don't block login
	} else {
		err = h.publisher.Publish(ctx, event)
		if err != nil {
			log.Printf("Failed to publish login event for user %s: %v", user.ID(), err)
			// Continue despite error - don't block login
		}
	}

	log.Printf("Google login successful for user: %s", user.ID())
	return tokens.AccessToken, tokens.RefreshToken, user.Username, nil
}
