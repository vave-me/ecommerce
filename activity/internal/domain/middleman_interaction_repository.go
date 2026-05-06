package domain

import "context"

type MiddlemanInteraction struct {
	ID         string
	ActivityID string
	ItemID     string
	ItemType   string
	ActionType string
}

type MostReactionResult struct {
	ItemID   string
	ItemType string
	Action   string // "like" or "dislike"
	Count    int64
}

type MiddlemanInteractionRepository interface {
	Add(ctx context.Context, interactionID, activityID, itemID, itemType, actionType string) error
	Update(ctx context.Context, interactionID string, actionType string) error
	Remove(ctx context.Context, interactionID string) error
	Find(ctx context.Context, interactionID string) (*MiddlemanInteraction, error)
	All(ctx context.Context, activityID string) ([]*MiddlemanInteraction, error)
	GetMostLiked(ctx context.Context, itemType string, limit int64) ([]*MostReactionResult, error)
	GetMostDisliked(ctx context.Context, itemType string, limit int64) ([]*MostReactionResult, error)
}
