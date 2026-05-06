package queries

import (
	"context"
	"middleman/comments/internal/domain"
)

type GetComments struct {
	ItemID string
}

type GetCommentsHandler struct {
	comments domain.MiddlemanRepository
}

func NewGetCommentsHandler(comments domain.MiddlemanRepository) GetCommentsHandler {
	return GetCommentsHandler{comments: comments}
}

func (h GetCommentsHandler) GetComments(ctx context.Context, query GetComments) ([]*domain.MiddlemanComment, error) {
	return h.comments.All(ctx, query.ItemID)
}
