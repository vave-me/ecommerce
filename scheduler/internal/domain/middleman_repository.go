package domain

import "context"

type MiddlemanScheduler struct {
	ID     string
	UserID string
}

type MiddlemanRepository interface {
	Add(ctx context.Context, schedulerID, userID string) error
	Remove(ctx context.Context, schedulerID string) error
	Find(ctx context.Context, userID string) (*MiddlemanScheduler, error)
	All(ctx context.Context, schedulerID string) ([]*MiddlemanScheduler, error)
}
