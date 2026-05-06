package queries

import (
	"context"
	"middleman/notifications/internal/domain"
)

type GetPreferences struct {
	UserID string
}

type GetPreferencesHandler struct {
	prefsRepo domain.PreferencesRepository
}

func NewGetPreferencesHandler(prefsRepo domain.PreferencesRepository) GetPreferencesHandler {
	return GetPreferencesHandler{prefsRepo: prefsRepo}
}

func (h GetPreferencesHandler) GetPreferences(ctx context.Context, query GetPreferences) (*domain.UserPreferences, error) {
	prefs, err := h.prefsRepo.GetUserPreferences(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	
	// If no preferences exist, create default ones
	if prefs == nil {
		prefs = h.createDefaultPreferences(query.UserID)
	}
	
	return prefs, nil
}

func (h GetPreferencesHandler) createDefaultPreferences(userID string) *domain.UserPreferences {
	allTypes := []domain.AlertType{
		domain.MessageType,
		domain.CommentType,
		domain.OfferType,
		domain.OrderType,
		domain.PaymentType,
		domain.ReviewType,
		domain.FollowingType,
		domain.ProductType,
		domain.WishlistType,
		domain.SupportType,
		domain.BasketType,
		domain.InteractionType,
	}
	
	prefs := &domain.UserPreferences{
		UserID:        userID,
		GlobalEnabled: true,
		EmailEnabled:  true,
		PushEnabled:   true,
		Preferences:   make(map[domain.AlertType]*domain.NotificationPreference),
	}
	
	for _, alertType := range allTypes {
		prefs.Preferences[alertType] = &domain.NotificationPreference{
			UserID:       userID,
			Type:         alertType,
			Enabled:      true,
			EmailEnabled: true,
			PushEnabled:  true,
		}
	}
	
	return prefs
}