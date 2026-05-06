package queries

import (
	"context"
	"middleman/reviews/internal/domain"
)

type GetMostReviewed struct {
	Offset int
	Limit  int
}

type GetMostReviewedHandler struct {
	reviews domain.MiddlemanRepository
}

func NewGetMostReviewedHandler(reviews domain.MiddlemanRepository) GetMostReviewedHandler {
	return GetMostReviewedHandler{reviews: reviews}
}

func (h GetMostReviewedHandler) GetMostReviewed(ctx context.Context, query GetMostReviewed) ([]*domain.ItemReviewCount, error) {
	return h.reviews.MostReviewedItems(ctx, query.Offset, query.Limit)
}
