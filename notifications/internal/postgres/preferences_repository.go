package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/stackus/errors"
	"middleman/notifications/internal/domain"
)

type PreferencesRepository struct {
	tableName string
	db        *sql.DB
}

func NewPreferencesRepository(tableName string, db *sql.DB) PreferencesRepository {
	return PreferencesRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r PreferencesRepository) GetUserPreferences(ctx context.Context, userID string) (*domain.UserPreferences, error) {
	query := `
		SELECT 
			user_id, 
			global_enabled, 
			email_enabled, 
			push_enabled, 
			preferences,
			updated_at
		FROM %s
		WHERE user_id = $1
	`
	
	var prefs domain.UserPreferences
	var prefsJSON []byte
	
	err := r.db.QueryRowContext(ctx, r.table(query), userID).Scan(
		&prefs.UserID,
		&prefs.GlobalEnabled,
		&prefs.EmailEnabled,
		&prefs.PushEnabled,
		&prefsJSON,
		&prefs.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No preferences found, return nil
		}
		return nil, errors.Wrap(err, "querying user preferences")
	}
	
	// Unmarshal preferences JSON
	prefMap := make(map[string]*domain.NotificationPreference)
	if err := json.Unmarshal(prefsJSON, &prefMap); err != nil {
		return nil, errors.Wrap(err, "unmarshaling preferences")
	}
	
	// Convert string keys back to AlertType
	prefs.Preferences = make(map[domain.AlertType]*domain.NotificationPreference)
	for typeStr, pref := range prefMap {
		alertType, err := domain.ToAlertType(typeStr)
		if err != nil {
			continue // Skip invalid types
		}
		prefs.Preferences[alertType] = pref
	}
	
	return &prefs, nil
}

func (r PreferencesRepository) SaveUserPreferences(ctx context.Context, prefs *domain.UserPreferences) error {
	// Convert AlertType keys to strings for JSON marshaling
	prefMap := make(map[string]*domain.NotificationPreference)
	for alertType, pref := range prefs.Preferences {
		prefMap[alertType.String()] = pref
	}
	
	prefsJSON, err := json.Marshal(prefMap)
	if err != nil {
		return errors.Wrap(err, "marshaling preferences")
	}
	
	query := `
		INSERT INTO %s (
			user_id, 
			global_enabled, 
			email_enabled, 
			push_enabled, 
			preferences,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) 
		DO UPDATE SET
			global_enabled = EXCLUDED.global_enabled,
			email_enabled = EXCLUDED.email_enabled,
			push_enabled = EXCLUDED.push_enabled,
			preferences = EXCLUDED.preferences,
			updated_at = EXCLUDED.updated_at
	`
	
	_, err = r.db.ExecContext(ctx, r.table(query),
		prefs.UserID,
		prefs.GlobalEnabled,
		prefs.EmailEnabled,
		prefs.PushEnabled,
		prefsJSON,
		time.Now(),
	)
	
	if err != nil {
		return errors.Wrap(err, "saving user preferences")
	}
	
	return nil
}

func (r PreferencesRepository) UpdatePreference(ctx context.Context, userID string, alertType domain.AlertType, pref *domain.NotificationPreference) error {
	// First get existing preferences
	existingPrefs, err := r.GetUserPreferences(ctx, userID)
	if err != nil {
		return err
	}
	
	// If no preferences exist, create new ones
	if existingPrefs == nil {
		existingPrefs = &domain.UserPreferences{
			UserID:        userID,
			GlobalEnabled: true,
			EmailEnabled:  true,
			PushEnabled:   true,
			Preferences:   make(map[domain.AlertType]*domain.NotificationPreference),
			UpdatedAt:     time.Now(),
		}
	}
	
	// Update the specific preference
	existingPrefs.Preferences[alertType] = pref
	existingPrefs.UpdatedAt = time.Now()
	
	// Save back
	return r.SaveUserPreferences(ctx, existingPrefs)
}

func (r PreferencesRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}