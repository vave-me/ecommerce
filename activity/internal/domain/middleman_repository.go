package domain

import "context"

type MiddlemanActivity struct {
	ID     string
	UserID string
}

type MiddlemanRepository interface {
	Add(ctx context.Context, activityID, userID string) error
	Remove(ctx context.Context, activityID string) error
	Find(ctx context.Context, userID string) (*MiddlemanActivity, error)
	All(ctx context.Context, activityID string) ([]*MiddlemanActivity, error)
}

