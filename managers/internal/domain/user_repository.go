package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type UserRepository interface {
	// User retrieval methods
	Find(ctx context.Context, userID string) (*models.User, error)
	GetBaseUser(ctx context.Context, userID string) (*models.BaseUser, error)
	ListUsers(ctx context.Context, userIDs []string) ([]*models.User, error)
	ListParticipatingUsers(ctx context.Context) ([]*models.User, error)

	// User management methods
	CreateUser(ctx context.Context, email, password, username, firstName, lastName, location string, lat, lng float32, thumbnail, language string) (string, error)
	UpdateUser(ctx context.Context, id, username, firstName, lastName, bio, privacy, background, location string, lat, lng float32, thumbnail string) (string, error)
	RenameUser(ctx context.Context, id, username string) (*models.User, error)
	EnableUser(ctx context.Context, id, verificationToken string) error
	DisableUser(ctx context.Context, id string) error

	// Authentication methods
	LoginUser(ctx context.Context, email, password string) (*models.LoginResponse, error)
	WebLoginWithGoogle(ctx context.Context, idToken string) (*models.LoginResponse, error)
	MobileLoginWithGoogle(ctx context.Context, idToken string) (*models.LoginResponse, error)
	LogUserOut(ctx context.Context, id, authToken, refreshToken string) error
	RefreshAuthToken(ctx context.Context, refreshToken, userID string) (*models.TokenResponse, error)
	ClearTokens(ctx context.Context, userID, tokenID, refreshToken, reason string) (*models.ClearTokensResponse, error)

	// Password management methods
	ForgotPassword(ctx context.Context, email string) (*models.MessageResponse, error)
	ResetPassword(ctx context.Context, token, newPassword string) (*models.MessageResponse, error)
}
