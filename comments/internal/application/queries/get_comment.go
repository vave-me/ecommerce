package queries

import (
	"context"
	"middleman/comments/internal/domain"
)

type GetComment struct {
	ID     string
	ItemID string
}

type GetCommentHandler struct {
	comments domain.MiddlemanRepository
}

func NewGetCommentHandler(comments domain.MiddlemanRepository) GetCommentHandler {
	return GetCommentHandler{comments: comments}
}

func (h GetCommentHandler) GetComment(ctx context.Context, query GetComment) (*domain.MiddlemanComment, error) {
	return h.comments.Find(ctx, query.ID, query.ItemID)
}
