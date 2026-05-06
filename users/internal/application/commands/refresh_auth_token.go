package commands

import (
	"context"
	"errors"
	"log"
	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type RefreshAuthToken struct {
	RefreshToken string
	UserID       string
}

type RefreshAuthTokenHandler struct {
	users     domain.UserRepository
	middleman domain.MiddlemanRepository
	auth      *auth.Auth
	publisher ddd.EventPublisher[ddd.Event]
}

func NewRefreshAuthTokenHandler(
	users domain.UserRepository,
	middleman domain.MiddlemanRepository,
	auth *auth.Auth,
	publisher ddd.EventPublisher[ddd.Event],
) RefreshAuthTokenHandler {
	return RefreshAuthTokenHandler{
		users:     users,
		middleman: middleman,
		auth:      auth,
		publisher: publisher,
	}
}

func (h RefreshAuthTokenHandler) RefreshAuthToken(ctx context.Context, cmd RefreshAuthToken) (string, string, error) {
	if cmd.RefreshToken == "" {
		return "", "", errors.New("refresh token is required")
	}

	// Validate the refresh token
	jwtUser, err := h.auth.ValidateRefreshToken(cmd.RefreshToken)
	if err != nil {
		log.Printf("Failed to validate refresh token: %v", err)
		return "", "", errors.New("invalid refresh token")
	}

	// Optional: Verify the user ID in the token matches the one provided
	if cmd.UserID != "" && jwtUser.ID != cmd.UserID {
		log.Printf("User ID mismatch: token has %s, request has %s", jwtUser.ID, cmd.UserID)
		return "", "", errors.New("invalid user ID in refresh token")
	}

	// Load the user from the repository to ensure it exists and is enabled
	user, err := h.middleman.Find(ctx, jwtUser.ID)
	if err != nil {
		log.Printf("Failed to find user %s: %v", jwtUser.ID, err)
		return "", "", errors.New("user not found")
	}

	if !user.Enabled {
		return "", "", errors.New("user account is disabled")
	}

	// Generate a new token pair (both access and refresh tokens)
	tokens, err := h.auth.GenerateTokenPair(&auth.JwtUser{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Lat:       user.Lat,
		Long:      user.Lng,
		Role:      user.Role,
	})
	if err != nil {
		log.Printf("Failed to generate token pair: %v", err)
		return "", "", err
	}

	// Publish a token refreshed event
	domainUser, err := h.users.Load(ctx, user.ID)
	if err == nil {
		// Create token IDs from the first few characters of the tokens for tracking
		oldTokenID := safeSubstring(cmd.RefreshToken, 8)
		newTokenID := safeSubstring(tokens.RefreshToken, 8)

		event, err := domainUser.TokenRefreshed(oldTokenID, newTokenID)
		if err == nil {
			// Save any changes to the user
			err = h.users.Save(ctx, domainUser)
			if err != nil {
				log.Printf("Failed to save user after token refresh: %v", err)
				// Don't fail the token refresh if saving fails
			}

			// Publish the token refresh event
			err = h.publisher.Publish(ctx, event)
			if err != nil {
				log.Printf("Failed to publish token refresh event: %v", err)
				// Don't fail the token refresh if publishing fails
			}
		}
	}

	log.Printf("Successfully refreshed tokens for user: %s", user.ID)
	return tokens.AccessToken, tokens.RefreshToken, nil
}

// Helper function to safely get a substring up to maxLen
func safeSubstring(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
