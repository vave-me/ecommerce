package queries

import (
	"context"
	"middleman/comments/internal/domain"
)

type GetMostCommented struct {
	Offset int
	Limit  int
}

type GetMostCommentedHandler struct {
	comments domain.MiddlemanRepository
}

func NewGetMostCommentedHandler(comments domain.MiddlemanRepository) GetMostCommentedHandler {
	return GetMostCommentedHandler{comments: comments}
}

func (h GetMostCommentedHandler) GetMostCommented(ctx context.Context, query GetMostCommented) ([]*domain.ItemCommentCount, error) {
	return h.comments.MostCommentedItems(ctx, query.Offset, query.Limit)
}
