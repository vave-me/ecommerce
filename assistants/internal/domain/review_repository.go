package domain

import (
	"context"
	"middleman/assistants/internal/models"
)

type ReviewRepository interface {
	// Basic CRUD operations
	CreateNewReview(ctx context.Context, senderID, itemID, itemType, content, categoryID, parentID string) (string, error) // returns review ID
	GetReviewByID(ctx context.Context, reviewID string) (*models.Review, error)
	GetAllReviewsForItem(ctx context.Context, itemID string) ([]*models.Review, error)
	EditReviewContent(ctx context.Context, reviewID, content string) error
	DeleteReviewByID(ctx context.Context, reviewID string) error

	// Review management operations
	ApproveReviewByID(ctx context.Context, reviewID string) (*models.Review, error) // returns updated review with approval info
	RejectReviewByID(ctx context.Context, reviewID string) (*models.Review, error)  // returns updated review with rejection info
	FlagReviewAsInappropriate(ctx context.Context, reviewID string) (*models.Review, error)    // returns updated review with flag status

	// Additional review operations needed by tool service
	UpdateReviewContentByID(ctx context.Context, reviewID, content string) error
	RemoveReviewPermanently(ctx context.Context, reviewID string) error
	UnflagReviewByID(ctx context.Context, reviewID string) error
	SearchReviewsByKeyword(ctx context.Context, query string, limit int64) ([]*models.Review, error)

	// Query operations
	GetReviewsBySenderID(ctx context.Context, senderID string) ([]*models.Review, error)
	GetAllApprovedReviews(ctx context.Context) ([]*models.Review, error)
	GetMostReviewedItems(ctx context.Context) ([]*models.ItemReviewCount, error)
	GetMostReviewedItemsByCategory(ctx context.Context, categoryID string, offset, limit int64) ([]*models.ItemReviewCount, error)

	// Legacy methods (for backward compatibility)
	FindReviewByIDAndItem(ctx context.Context, reviewID, itemID string) (*models.Review, error)
	GetAllReviewsByItemID(ctx context.Context, itemID string) ([]*models.Review, error)
	FindReviewsBySenderUserID(ctx context.Context, senderID string) ([]*models.Review, error)
}
