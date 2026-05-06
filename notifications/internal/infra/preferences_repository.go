package infra

import (
	"context"
	"sync"
	"time"
	"middleman/notifications/internal/domain"
)

// InMemoryPreferencesRepository is a simple in-memory implementation
// TODO: Replace with PostgreSQL implementation
type InMemoryPreferencesRepository struct {
	mu    sync.RWMutex
	prefs map[string]*domain.UserPreferences
}

func NewInMemoryPreferencesRepository() *InMemoryPreferencesRepository {
	return &InMemoryPreferencesRepository{
		prefs: make(map[string]*domain.UserPreferences),
	}
}

func (r *InMemoryPreferencesRepository) GetUserPreferences(ctx context.Context, userID string) (*domain.UserPreferences, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	prefs, exists := r.prefs[userID]
	if !exists {
		return nil, nil // Return nil if not found, not an error
	}
	
	// Deep copy to prevent external modifications
	return r.copyPreferences(prefs), nil
}

func (r *InMemoryPreferencesRepository) SaveUserPreferences(ctx context.Context, prefs *domain.UserPreferences) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Deep copy to prevent external modifications
	r.prefs[prefs.UserID] = r.copyPreferences(prefs)
	return nil
}

func (r *InMemoryPreferencesRepository) UpdatePreference(ctx context.Context, userID string, alertType domain.AlertType, pref *domain.NotificationPreference) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	userPrefs, exists := r.prefs[userID]
	if !exists {
		// Create new preferences if they don't exist
		userPrefs = &domain.UserPreferences{
			UserID:        userID,
			GlobalEnabled: true,
			EmailEnabled:  true,
			PushEnabled:   true,
			Preferences:   make(map[domain.AlertType]*domain.NotificationPreference),
			UpdatedAt:     time.Now(),
		}
		r.prefs[userID] = userPrefs
	}
	
	// Update the specific preference
	userPrefs.Preferences[alertType] = &domain.NotificationPreference{
		UserID:       pref.UserID,
		Type:         pref.Type,
		Enabled:      pref.Enabled,
		EmailEnabled: pref.EmailEnabled,
		PushEnabled:  pref.PushEnabled,
		UpdatedAt:    time.Now(),
	}
	userPrefs.UpdatedAt = time.Now()
	
	return nil
}

func (r *InMemoryPreferencesRepository) copyPreferences(prefs *domain.UserPreferences) *domain.UserPreferences {
	copied := &domain.UserPreferences{
		UserID:        prefs.UserID,
		GlobalEnabled: prefs.GlobalEnabled,
		EmailEnabled:  prefs.EmailEnabled,
		PushEnabled:   prefs.PushEnabled,
		Preferences:   make(map[domain.AlertType]*domain.NotificationPreference),
		UpdatedAt:     prefs.UpdatedAt,
	}
	
	for k, v := range prefs.Preferences {
		copied.Preferences[k] = &domain.NotificationPreference{
			UserID:       v.UserID,
			Type:         v.Type,
			Enabled:      v.Enabled,
			EmailEnabled: v.EmailEnabled,
			PushEnabled:  v.PushEnabled,
			UpdatedAt:    v.UpdatedAt,
		}
	}
	
	return copied
}