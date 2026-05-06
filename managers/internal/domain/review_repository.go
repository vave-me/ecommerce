package domain

import (
	"context"
	"middleman/managers/internal/models"
)

type ReviewRepository interface {
	// Basic CRUD operations
	AddReview(ctx context.Context, senderID, itemID, itemType, content, categoryID, parentID string) (string, error) // returns review ID
	GetReview(ctx context.Context, reviewID string) (*models.Review, error)
	GetReviews(ctx context.Context, itemID string) ([]*models.Review, error)
	EditReview(ctx context.Context, reviewID, content string) error
	RemoveReview(ctx context.Context, reviewID string) error

	// Review management operations
	ApproveReview(ctx context.Context, reviewID string) (*models.Review, error) // returns updated review with approval info
	RejectReview(ctx context.Context, reviewID string) (*models.Review, error)  // returns updated review with rejection info
	FlagReview(ctx context.Context, reviewID string) (*models.Review, error)    // returns updated review with flag status

	// Additional review operations needed by tool service
	UpdateReview(ctx context.Context, reviewID, content string) error
	DeleteReview(ctx context.Context, reviewID string) error
	UnflagReview(ctx context.Context, reviewID string) error
	SearchReviews(ctx context.Context, query string, limit int64) ([]*models.Review, error)

	// Query operations
	GetReviewsBySender(ctx context.Context, senderID string) ([]*models.Review, error)
	GetApprovedReviews(ctx context.Context) ([]*models.Review, error)
	GetMostReviewed(ctx context.Context) ([]*models.ItemReviewCount, error)
	GetMostReviewedByCategory(ctx context.Context, categoryID string, offset, limit int64) ([]*models.ItemReviewCount, error)

	// Legacy methods (for backward compatibility)
	FindOneReview(ctx context.Context, reviewID, itemID string) (*models.Review, error)
	AllReviews(ctx context.Context, itemID string) ([]*models.Review, error)
	FindBySenderID(ctx context.Context, senderID string) ([]*models.Review, error)
}
