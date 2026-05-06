package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type ActivityRepository interface {
	// Core activity management
	CreateUserActivity(ctx context.Context, userID string) (string, error)
	GetActivityByUserID(ctx context.Context, userID string) (*models.Activity, error)
	DeleteAllUserActivities(ctx context.Context, activityID string) error
	ArchiveUserActivity(ctx context.Context, activityID, reason string) error
	RestoreUserActivity(ctx context.Context, activityID, reason string) error
	AddUserInteraction(ctx context.Context, interactionID, activityID, itemID, itemType, actionType string) error
	UpdateUserInteraction(ctx context.Context, interactionID string, actionType string) error
	DeleteUserInteraction(ctx context.Context, interactionID string) error
	GetInteractionByID(ctx context.Context, interactionID string) (*models.Interaction, error)
	GetAllActivityInteractions(ctx context.Context, activityID string) ([]*models.Interaction, error)
	GetMostLikedItems(ctx context.Context, itemType string, limit int64) ([]*models.MostReactionResult, error)
	GetMostDislikedItems(ctx context.Context, itemType string, limit int64) ([]*models.MostReactionResult, error)
}
