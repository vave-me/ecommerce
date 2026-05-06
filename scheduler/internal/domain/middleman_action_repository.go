package domain

import (
	"context"
	"time"
)

type MiddlemanAction struct {
	ID                  string
	SchedulerID         string
	NaturalLanguageTask string
	ExecutionTime       time.Time
	Status              string
	CreatedAt           time.Time
	ExecutedAt          *time.Time
	Result              string
	ErrorMessage        string
}

type MiddlemanActionRepository interface {
	Add(ctx context.Context, actionID, schedulerID, task string, executionTime time.Time) error
	UpdateStatus(ctx context.Context, actionID, status, result, errorMessage string) error
	Remove(ctx context.Context, actionID string) error
	Find(ctx context.Context, actionID string) (*MiddlemanAction, error)
	All(ctx context.Context, schedulerID string) ([]*MiddlemanAction, error)
	GetPendingActions(ctx context.Context, beforeTime time.Time) ([]*MiddlemanAction, error)
}
