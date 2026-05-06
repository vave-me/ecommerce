package queries

import (
	"context"
	"middleman/reviews/internal/domain"
)

type GetMostReviewedByCategory struct {
	ItemType   domain.ItemType
	CategoryID string
	Offset     int
	Limit      int
}

type GetMostReviewedByCategoryHandler struct {
	reviews domain.MiddlemanRepository
}

func NewGetMostReviewedByCategoryHandler(reviews domain.MiddlemanRepository) GetMostReviewedByCategoryHandler {
	return GetMostReviewedByCategoryHandler{reviews: reviews}
}

func (h GetMostReviewedByCategoryHandler) GetMostReviewedByCategory(ctx context.Context, query GetMostReviewedByCategory) ([]*domain.ItemReviewCount, error) {
	return h.reviews.MostReviewedItemsByCategory(ctx, query.ItemType, query.CategoryID, query.Offset, query.Limit)
}
