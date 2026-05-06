package commands

import (
	"context"
	"log"
	"middleman/internal/ddd"
	"middleman/users/internal/domain"
)

type (
	CreateUserFromGoogle struct {
		ID            string
		Email         string
		FirstName     string
		LastName      string
		Enabled       bool
		Latitude      float64
		Longitude     float64
		GoogleID      string // CRITICAL
		Thumbnail     string
		Locale        string
		EmailVerified bool
		Role          string // User role to assign
	}
	CreateUserFromGoogleHandler struct {
		users     domain.UserRepository
		middleman domain.MiddlemanRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewCreateUserFromGoogleHandler(users domain.UserRepository, middleman domain.MiddlemanRepository, publisher ddd.EventPublisher[ddd.Event]) CreateUserFromGoogleHandler {
	return CreateUserFromGoogleHandler{
		users:     users,
		middleman: middleman,
		publisher: publisher,
	}
}

func (h CreateUserFromGoogleHandler) CreateUserFromGoogle(ctx context.Context, cmd CreateUserFromGoogle) error {

	log.Printf("Error in application.CreateUserFromGoogle: %v", cmd.ID)

	// Check if email already exists
	existingUser, err := h.middleman.FindByEmail(ctx, cmd.Email)
	if err == nil && existingUser != nil {
		// Check if the existing user has a Google ID
		if existingUser.GoogleID == "" {
			// User exists with this email but hasn't linked Google account yet
			// We should link the Google account to the existing user instead of creating a duplicate
			return domain.ErrDuplicateEmail
		} else if existingUser.GoogleID == cmd.GoogleID {
			// User already exists with this Google ID, this shouldn't happen but handle gracefully
			return nil
		} else {
			// User exists with a different Google ID - this is a conflict
			return domain.ErrDuplicateEmail
		}
	}

	user, err := h.users.Load(ctx, cmd.ID)
	if err != nil {

		return err
	}
	// Convert string role to UserRole type
	userRole := domain.ToUserRole(cmd.Role)
	
	event, err := user.InitUserWithGoogle(cmd.Email, cmd.FirstName, cmd.LastName, cmd.Enabled, cmd.GoogleID, cmd.Thumbnail, cmd.Latitude, cmd.Longitude, userRole)
	if err != nil {
		return err
	}

	err = h.users.Save(ctx, user)

	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
