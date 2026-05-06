package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type UserRepository interface {
	// User retrieval methods
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	GetBaseUserByID(ctx context.Context, userID string) (*models.BaseUser, error)
	GetMultipleUsersByIDs(ctx context.Context, userIDs []string) ([]*models.User, error)
	GetAllParticipatingUsers(ctx context.Context) ([]*models.User, error)

	// User management methods
	CreateNewUser(ctx context.Context, email, password, username, firstName, lastName, location string, lat, lng float32, thumbnail, language string) (string, error)
	UpdateUserProfile(ctx context.Context, id, username, firstName, lastName, bio, privacy, background, location string, lat, lng float32, thumbnail string) (string, error)
	ChangeUsername(ctx context.Context, id, username string) (*models.User, error)
	ActivateUserAccount(ctx context.Context, id, verificationToken string) error
	DeactivateUserAccount(ctx context.Context, id string) error

	// Authentication methods
	AuthenticateUser(ctx context.Context, email, password string) (*models.LoginResponse, error)
	AuthenticateWithGoogleWeb(ctx context.Context, idToken string) (*models.LoginResponse, error)
	AuthenticateWithGoogleMobile(ctx context.Context, idToken string) (*models.LoginResponse, error)
	LogoutUser(ctx context.Context, id, authToken, refreshToken string) error
	RefreshUserAuthToken(ctx context.Context, refreshToken, userID string) (*models.TokenResponse, error)
	RevokeUserTokens(ctx context.Context, userID, tokenID, refreshToken, reason string) (*models.ClearTokensResponse, error)

	// Password management methods
	SendPasswordResetEmail(ctx context.Context, email string) (*models.MessageResponse, error)
	ResetUserPassword(ctx context.Context, token, newPassword string) (*models.MessageResponse, error)
}
