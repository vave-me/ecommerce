package domain

import "context"

type MiddlemanUser struct {
	ID         string
	Email      string
	Username   string
	FirstName  string
	LastName   string
	GoogleID   string
	Enabled    bool
	Lat        float64
	Lng        float64
	Thumbnail  string
	Bio        string
	Background string
	Privacy    string
	Role       string
}
type MiddlemanViewUser struct {
	ID         string
	Username   string
	Lat        float64
	Lng        float64
	Location   string
	Thumbnail  string
	Bio        string
	Background string
	Privacy    string
}

type MiddlemanRepository interface {
	// AddUser creates a new user in the system
	AddUser(ctx context.Context, userID, email, username, firstname, lastname, googleID string, enabled bool, lat, lng float64, thumbnail, role string) error
	UpdateUser(ctx context.Context, userID, username, firstname, lastname, bio, privacy, background string, lat, lng float64, thumbnail, role string) error
	// EnableUser updates the enabled status of a user
	EnableUser(ctx context.Context, userID string, enabled bool) error
	// RenameUser updates the name of the user
	RenameUser(ctx context.Context, userID string, newName string) error
	// Find retrieves a user by ID
	Find(ctx context.Context, userID string) (*MiddlemanUser, error)
	// All retrieves all users in the system
	All(ctx context.Context) ([]*MiddlemanUser, error)
	// AllEnabled retrieves all users who are enabled
	AllEnabled(ctx context.Context) ([]*MiddlemanUser, error)
	// LogUserIn logs the user in by storing the necessary session/token data (to be implemented)
	LogUserIn(ctx context.Context, userID string) error
	FindSimple(ctx context.Context, userID string) (*MiddlemanViewUser, error)
	// LogUserOut logs the user out by clearing the necessary session/token data (to be implemented)
	LogUserOut(ctx context.Context, userID string) error
	// FindByEmail user by his email
	FindByEmail(ctx context.Context, email string) (*MiddlemanUser, error)
	FindByGoogleID(ctx context.Context, googleID string) (*MiddlemanUser, error)
}
