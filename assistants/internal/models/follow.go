package models

import "time"

// Core follow entities
type Follow struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	FollowedUserID   string    `json:"followed_user_id"`
	FollowedUserType string    `json:"followed_user_type"`
	Content          string    `json:"content"`
	CategoryID       string    `json:"category_id"`
	ParentID         string    `json:"parent_id"`
	FollowStatus     string    `json:"follow_status"`
	Flagged          bool      `json:"flagged"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ApprovedAt       time.Time `json:"approved_at,omitempty"`

	// Legacy fields for backward compatibility
	Approved bool `json:"approved"`
}

type ItemFollowCount struct {
	FollowedUserID string `json:"followed_user_id"`
	CategoryID     string `json:"category_id"`
	FollowingCount int64  `json:"following_count"`
}

// Protobuf request/response types

// AddFollow
type AddFollowRequest struct {
	UserID           string `json:"user_id"`
	FollowedUserID   string `json:"followed_user_id"`
	FollowedUserType string `json:"followed_user_type"`
	Content          string `json:"content"`
	CategoryID       string `json:"category_id"`
	ParentID         string `json:"parent_id"`
}

type AddFollowResponse struct {
	ID string `json:"id"`
}

// ApproveFollow
type ApproveFollowRequest struct {
	ID string `json:"id"`
}

type ApproveFollowResponse struct {
	ID           string    `json:"id"`
	FollowStatus string    `json:"follow_status"`
	ApprovedAt   time.Time `json:"approved_at"`
}

// RejectFollow
type RejectFollowRequest struct {
	ID string `json:"id"`
}

type RejectFollowResponse struct {
	ID           string `json:"id"`
	FollowStatus string `json:"follow_status"`
}

// FlagFollow
type FlagFollowRequest struct {
	ID string `json:"id"`
}

type FlagFollowResponse struct {
	ID      string `json:"id"`
	Flagged bool   `json:"flagged"`
}

// EditFollow
type EditFollowRequest struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type EditFollowResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// RemoveFollow
type RemoveFollowRequest struct {
	ID string `json:"id"`
}

type RemoveFollowResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetFollow
type GetFollowRequest struct {
	ID string `json:"id"`
}

type GetFollowResponse struct {
	Follow *Follow `json:"follow"`
}

// GetFollowing
type GetFollowingRequest struct {
	FollowedUserID string `json:"followed_user_id"`
}

type GetFollowingResponse struct {
	Following []*Follow `json:"following"`
	Count     int64     `json:"count"`
}

// GetFollowingBySender
type GetFollowingBySenderRequest struct {
	UserID string `json:"user_id"`
}

type GetFollowingBySenderResponse struct {
	Following []*Follow `json:"following"`
	Count     int64     `json:"count"`
}

// GetMostFollowed
type GetMostFollowedRequest struct {
	// Empty request
}

type GetMostFollowedResponse struct {
	ItemFollowCounts []*ItemFollowCount `json:"item_follow_counts"`
	Count            int64              `json:"count"`
}

// GetMostFollowedByCategory
type GetMostFollowedByCategoryRequest struct {
	CategoryID string `json:"category_id"`
	Offset     int64  `json:"offset"`
	Limit      int64  `json:"limit"`
}

type GetMostFollowedByCategoryResponse struct {
	ItemFollowCounts []*ItemFollowCount `json:"item_follow_counts"`
	Count            int64              `json:"count"`
}

// GetApprovedFollowing
type GetApprovedFollowingRequest struct {
	// Empty request
}

type GetApprovedFollowingResponse struct {
	Following []*Follow `json:"following"`
	Count     int64     `json:"count"`
}

// Additional response types for extended functionality
type FollowingStatsResponse struct {
	UserID          string `json:"user_id"`
	TotalFollowing  int64  `json:"total_following"`
	TotalFollowers  int64  `json:"total_followers"`
	PendingFollows  int64  `json:"pending_follows"`
	ApprovedFollows int64  `json:"approved_follows"`
	RejectedFollows int64  `json:"rejected_follows"`
	FlaggedFollows  int64  `json:"flagged_follows"`
}

// Follow status constants
const (
	FollowStatusPending  = "pending"
	FollowStatusApproved = "approved"
	FollowStatusRejected = "rejected"
	FollowStatusActive   = "active"
	FollowStatusBlocked  = "blocked"
)
