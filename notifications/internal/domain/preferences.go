package domain

import (
	"context"
	"time"
)

// NotificationPreference represents user preferences for a specific notification type
type NotificationPreference struct {
	UserID       string
	Type         AlertType
	Enabled      bool
	EmailEnabled bool
	PushEnabled  bool
	UpdatedAt    time.Time
}

// UserPreferences represents global notification preferences for a user
type UserPreferences struct {
	UserID        string
	GlobalEnabled bool
	EmailEnabled  bool
	PushEnabled   bool
	Preferences   map[AlertType]*NotificationPreference
	UpdatedAt     time.Time
}

// PreferencesRepository defines the interface for storing user notification preferences
type PreferencesRepository interface {
	// GetUserPreferences retrieves all preferences for a user
	GetUserPreferences(ctx context.Context, userID string) (*UserPreferences, error)
	
	// SaveUserPreferences saves or updates user preferences
	SaveUserPreferences(ctx context.Context, prefs *UserPreferences) error
	
	// UpdatePreference updates a single notification type preference
	UpdatePreference(ctx context.Context, userID string, alertType AlertType, pref *NotificationPreference) error
}