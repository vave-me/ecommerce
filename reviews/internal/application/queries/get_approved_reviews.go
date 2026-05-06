package queries

import (
	"context"
	"middleman/reviews/internal/domain"
)

type GetApprovedReviews struct {
	ItemID string
}

type GetApprovedReviewsHandler struct {
	reviews domain.MiddlemanRepository
}

func NewGetApprovedReviewsHandler(reviews domain.MiddlemanRepository) GetApprovedReviewsHandler {
	return GetApprovedReviewsHandler{reviews: reviews}
}

func (h GetApprovedReviewsHandler) GetApprovedReviews(ctx context.Context, query GetApprovedReviews) ([]*domain.MiddlemanReview, error) {
	return h.reviews.All(ctx, query.ItemID)
}
