package queries

import (
	"context"
	"middleman/reviews/internal/domain"
)

type GetReviewsBySender struct {
	SenderID string
}

type GetReviewsBySenderHandler struct {
	reviews domain.MiddlemanRepository
}

func NewGetReviewsBySenderHandler(reviews domain.MiddlemanRepository) GetReviewsBySenderHandler {
	return GetReviewsBySenderHandler{reviews: reviews}
}

func (h GetReviewsBySenderHandler) GetReviewsBySender(ctx context.Context, query GetReviewsBySender) ([]*domain.MiddlemanReview, error) {
	return h.reviews.FindBySenderID(ctx, query.SenderID)
}
