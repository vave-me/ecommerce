package commands

import (
	"context"
	"time"
	"github.com/stackus/errors"
	"middleman/notifications/internal/domain"
)

type UpdatePreferences struct {
	UserID        string
	GlobalEnabled *bool
	EmailEnabled  *bool
	PushEnabled   *bool
	Preferences   []PreferenceUpdate
}

type PreferenceUpdate struct {
	Type         string
	Enabled      bool
	EmailEnabled bool
	PushEnabled  bool
}

type UpdatePreferencesHandler struct {
	prefsRepo domain.PreferencesRepository
}

func NewUpdatePreferencesHandler(prefsRepo domain.PreferencesRepository) UpdatePreferencesHandler {
	return UpdatePreferencesHandler{prefsRepo: prefsRepo}
}

func (h UpdatePreferencesHandler) UpdatePreferences(ctx context.Context, cmd UpdatePreferences) error {
	// Get existing preferences or create new ones
	prefs, err := h.prefsRepo.GetUserPreferences(ctx, cmd.UserID)
	if err != nil && !errors.Is(err, domain.ErrAlertNotFound) {
		return errors.Wrap(err, "getting user preferences")
	}
	
	if prefs == nil {
		// Create new preferences
		prefs = &domain.UserPreferences{
			UserID:        cmd.UserID,
			GlobalEnabled: true,
			EmailEnabled:  true,
			PushEnabled:   true,
			Preferences:   make(map[domain.AlertType]*domain.NotificationPreference),
			UpdatedAt:     time.Now(),
		}
	}
	
	// Update global settings if provided
	if cmd.GlobalEnabled != nil {
		prefs.GlobalEnabled = *cmd.GlobalEnabled
	}
	if cmd.EmailEnabled != nil {
		prefs.EmailEnabled = *cmd.EmailEnabled
	}
	if cmd.PushEnabled != nil {
		prefs.PushEnabled = *cmd.PushEnabled
	}
	
	// Update individual preferences
	for _, prefUpdate := range cmd.Preferences {
		alertType, err := domain.ToAlertType(prefUpdate.Type)
		if err != nil {
			return errors.Wrap(err, "invalid alert type")
		}
		
		pref, exists := prefs.Preferences[alertType]
		if !exists {
			pref = &domain.NotificationPreference{
				UserID: cmd.UserID,
				Type:   alertType,
			}
			prefs.Preferences[alertType] = pref
		}
		
		pref.Enabled = prefUpdate.Enabled
		pref.EmailEnabled = prefUpdate.EmailEnabled
		pref.PushEnabled = prefUpdate.PushEnabled
		pref.UpdatedAt = time.Now()
	}
	
	prefs.UpdatedAt = time.Now()
	
	// Save preferences
	if err := h.prefsRepo.SaveUserPreferences(ctx, prefs); err != nil {
		return errors.Wrap(err, "saving user preferences")
	}
	
	return nil
}