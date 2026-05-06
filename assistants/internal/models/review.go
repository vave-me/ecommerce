package models

import "time"

type Review struct {
	ID           string
	SenderID     string
	ItemID       string
	ItemType     string
	ParentID     string
	CategoryID   string
	Content      string
	ReviewStatus string // "pending", "approved", "rejected"
	Flagged      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ApprovedAt   *time.Time // Pointer to allow nil for non-approved reviews
}

// ItemReviewCount represents the count of reviews for an item
type ItemReviewCount struct {
	ItemID       string
	CategoryID   string
	ReviewsCount int64
}

// ReviewStatus constants
const (
	ReviewStatusPending  = "pending"
	ReviewStatusApproved = "approved"
	ReviewStatusRejected = "rejected"
)
