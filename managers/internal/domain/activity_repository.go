package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type ActivityRepository interface {
	// Core activity management
	CreateActivity(ctx context.Context, userID string) (string, error)
	FindActivityId(ctx context.Context, userID string) (*models.Activity, error)
	RemoveAllActivities(ctx context.Context, activityID string) error
	ArchiveActivity(ctx context.Context, activityID, reason string) error
	RestoreActivity(ctx context.Context, activityID, reason string) error

	// Interaction management
	AddLikeOrDislike(ctx context.Context, interactionID, activityID, itemID, itemType, actionType string) error
	UpdateLikeOrDislike(ctx context.Context, interactionID string, actionType string) error
	RemoveInteraction(ctx context.Context, interactionID string) error
	Find(ctx context.Context, interactionID string) (*models.Interaction, error)
	AllInteraction(ctx context.Context, activityID string) ([]*models.Interaction, error)

	// Analytics and insights
	GetMostLiked(ctx context.Context, itemType string, limit int64) ([]*models.MostReactionResult, error)
	GetMostDisliked(ctx context.Context, itemType string, limit int64) ([]*models.MostReactionResult, error)
}
