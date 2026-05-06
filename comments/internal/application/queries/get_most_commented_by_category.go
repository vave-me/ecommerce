package queries

import (
	"context"
	"middleman/comments/internal/domain"
)

type GetMostCommentedByCategory struct {
	ItemType   domain.ItemType
	CategoryID string
	Offset     int
	Limit      int
}

type GetMostCommentedByCategoryHandler struct {
	comments domain.MiddlemanRepository
}

func NewGetMostCommentedByCategoryHandler(comments domain.MiddlemanRepository) GetMostCommentedByCategoryHandler {
	return GetMostCommentedByCategoryHandler{comments: comments}
}

func (h GetMostCommentedByCategoryHandler) GetMostCommentedByCategory(ctx context.Context, query GetMostCommentedByCategory) ([]*domain.ItemCommentCount, error) {
	return h.comments.MostCommentedItemsByCategory(ctx, query.ItemType, query.CategoryID, query.Offset, query.Limit)
}
