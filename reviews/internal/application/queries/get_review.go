package queries

import (
	"context"
	"middleman/reviews/internal/domain"
)

type GetReview struct {
	ID     string
	ItemID string
}

type GetReviewHandler struct {
	reviews domain.MiddlemanRepository
}

func NewGetReviewHandler(reviews domain.MiddlemanRepository) GetReviewHandler {
	return GetReviewHandler{reviews: reviews}
}

func (h GetReviewHandler) GetReview(ctx context.Context, query GetReview) (*domain.MiddlemanReview, error) {
	return h.reviews.Find(ctx, query.ID, query.ItemID)
}
