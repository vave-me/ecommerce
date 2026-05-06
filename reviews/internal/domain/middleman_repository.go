package domain

import (
	"context"
	"time"
)

// MiddlemanReview is your main reviews model
type MiddlemanReview struct {
	ID         string
	SenderID   string
	ItemID     string
	ItemType   string
	ParentID   string
	CategoryID string
	Content    string
	Approved   bool
	Flagged    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ItemReviewCount is used for aggregator results
type ItemReviewCount struct {
	ItemID       string
	ItemType     string
	CategoryID   string
	ReviewsCount int64
}

type MiddlemanRepository interface {
	Add(ctx context.Context, reviewID, senderID, itemID string, itemType ItemType, content, categoryID, parentID string) error
	Find(ctx context.Context, reviewID, itemID string) (*MiddlemanReview, error)
	All(ctx context.Context, itemID string) ([]*MiddlemanReview, error)
	FindBySenderID(ctx context.Context, senderID string) ([]*MiddlemanReview, error)
	MostReviewedItems(ctx context.Context, limit, offset int) ([]*ItemReviewCount, error)
	MostReviewedItemsByCategory(ctx context.Context, itemType ItemType, categoryID string, limit, offset int) ([]*ItemReviewCount, error)
	Remove(ctx context.Context, reviewID string) error
}
