package queries

import (
	"context"
	"middleman/reviews/internal/domain"
)

type GetReviews struct {
	ItemID string
}

type GetReviewsHandler struct {
	reviews domain.MiddlemanRepository
}

func NewGetReviewsHandler(reviews domain.MiddlemanRepository) GetReviewsHandler {
	return GetReviewsHandler{reviews: reviews}
}

func (h GetReviewsHandler) GetReviews(ctx context.Context, query GetReviews) ([]*domain.MiddlemanReview, error) {
	return h.reviews.All(ctx, query.ItemID)
}
