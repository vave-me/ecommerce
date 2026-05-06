package models

import "time"

// Core comment entities
type Comment struct {
	ID            string    `json:"id"`
	SenderID      string    `json:"sender_id"`
	ItemID        string    `json:"item_id"`
	ItemType      string    `json:"item_type"`
	Content       string    `json:"content"`
	CategoryID    string    `json:"category_id"`
	ParentID      string    `json:"parent_id"`
	CommentStatus string    `json:"comment_status"`
	Flagged       bool      `json:"flagged"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ApprovedAt    time.Time `json:"approved_at,omitempty"`

	// Legacy fields for backward compatibility
	Approved bool `json:"approved"`
}

type ItemCommentCount struct {
	ItemID        string `json:"item_id"`
	CategoryID    string `json:"category_id"`
	CommentsCount int64  `json:"comments_count"`
}

// Protobuf request/response types

// AddComment
type AddCommentRequest struct {
	SenderID   string `json:"sender_id"`
	ItemID     string `json:"item_id"`
	ItemType   string `json:"item_type"`
	Content    string `json:"content"`
	CategoryID string `json:"category_id"`
	ParentID   string `json:"parent_id"`
}

type AddCommentResponse struct {
	ID string `json:"id"`
}

// ApproveComment
type ApproveCommentRequest struct {
	ID string `json:"id"`
}

type ApproveCommentResponse struct {
	ID            string    `json:"id"`
	CommentStatus string    `json:"comment_status"`
	ApprovedAt    time.Time `json:"approved_at"`
}

// RejectComment
type RejectCommentRequest struct {
	ID string `json:"id"`
}

type RejectCommentResponse struct {
	ID            string `json:"id"`
	CommentStatus string `json:"comment_status"`
}

// FlagComment
type FlagCommentRequest struct {
	ID string `json:"id"`
}

type FlagCommentResponse struct {
	ID      string `json:"id"`
	Flagged bool   `json:"flagged"`
}

// EditComment
type EditCommentRequest struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type EditCommentResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// RemoveComment
type RemoveCommentRequest struct {
	ID string `json:"id"`
}

type RemoveCommentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetComment
type GetCommentRequest struct {
	ID string `json:"id"`
}

type GetCommentResponse struct {
	Comment *Comment `json:"comment"`
}

// GetComments
type GetCommentsRequest struct {
	ItemID string `json:"item_id"`
}

type GetCommentsResponse struct {
	Comments []*Comment `json:"comments"`
}

// GetCommentsBySender
type GetCommentsBySenderRequest struct {
	SenderID string `json:"sender_id"`
}

type GetCommentsBySenderResponse struct {
	Comments []*Comment `json:"comments"`
}

// GetMostCommented
type GetMostCommentedRequest struct{}

type GetMostCommentedResponse struct {
	ItemCommentCounts []*ItemCommentCount `json:"item_comment_counts"`
}

// GetMostCommentedByCategory
type GetMostCommentedByCategoryRequest struct {
	CategoryID string `json:"category_id"`
	Offset     int64  `json:"offset"`
	Limit      int64  `json:"limit"`
}

type GetMostCommentedByCategoryResponse struct {
	ItemCommentCounts []*ItemCommentCount `json:"item_comment_counts"`
}

// GetApprovedComments
type GetApprovedCommentsRequest struct{}

type GetApprovedCommentsResponse struct {
	Comments []*Comment `json:"comments"`
}

// Additional response types for extended functionality
type CommentStatsResponse struct {
	ItemID        string `json:"item_id"`
	TotalComments int64  `json:"total_comments"`
	ApprovedCount int64  `json:"approved_count"`
	PendingCount  int64  `json:"pending_count"`
	RejectedCount int64  `json:"rejected_count"`
	FlaggedCount  int64  `json:"flagged_count"`
	ThreadsCount  int64  `json:"threads_count"`
}

// Comment status constants
const (
	CommentStatusPending  = "pending"
	CommentStatusApproved = "approved"
	CommentStatusRejected = "rejected"
	CommentStatusFlagged  = "flagged"
	CommentStatusDeleted  = "deleted"
	CommentStatusDraft    = "draft"
)
