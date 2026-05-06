package domain

import (
	"context"
	"time"
)

// MiddlemanComment is your main comments model
type MiddlemanComment struct {
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

// ItemCommentCount is used for aggregator results
type ItemCommentCount struct {
	ItemID        string
	ItemType      string
	CategoryID    string
	CommentsCount int64
}

type MiddlemanRepository interface {
	Add(ctx context.Context, commentID, senderID, itemID string, itemType ItemType, content, categoryID, parentID string) error
	Find(ctx context.Context, commentID, itemID string) (*MiddlemanComment, error)
	All(ctx context.Context, itemID string) ([]*MiddlemanComment, error)
	FindBySenderID(ctx context.Context, senderID string) ([]*MiddlemanComment, error)
	MostCommentedItems(ctx context.Context, limit, offset int) ([]*ItemCommentCount, error)
	MostCommentedItemsByCategory(ctx context.Context, itemType ItemType, categoryID string, limit, offset int) ([]*ItemCommentCount, error)
	Remove(ctx context.Context, commentID string) error
}
