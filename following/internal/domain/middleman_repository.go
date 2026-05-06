package domain

import (
	"context"
	"time"
)

// MiddlemanFollow is your main following model
type MiddlemanFollow struct {
	ID               string
	UserID           string
	FollowedUserID   string
	FollowedUserType string
	ParentID         string
	CategoryID       string
	Content          string
	Approved         bool
	Flagged          bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ItemFollowCount is used for aggregator results
type ItemFollowCount struct {
	FollowedUserID   string
	FollowedUserType string
	CategoryID       string
	FollowingCount   int64
}

type MiddlemanRepository interface {
	Add(ctx context.Context, followID, userID, followedUserID string, followedUserType FollowedUserType, content, categoryID, parentID string) error
	Find(ctx context.Context, followID, followedUserID string) (*MiddlemanFollow, error)
	All(ctx context.Context, followedUserID string) ([]*MiddlemanFollow, error)
	FindByUserID(ctx context.Context, userID string) ([]*MiddlemanFollow, error)
	MostFollowedItems(ctx context.Context, limit, offset int) ([]*ItemFollowCount, error)
	MostFollowedItemsByCategory(ctx context.Context, followedUserType FollowedUserType, categoryID string, limit, offset int) ([]*ItemFollowCount, error)
	Remove(ctx context.Context, followID string) error
}
