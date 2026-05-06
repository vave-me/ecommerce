package queries

import (
	"context"
	"middleman/comments/internal/domain"
)

type GetCommentsBySender struct {
	SenderID string
}

type GetCommentsBySenderHandler struct {
	comments domain.MiddlemanRepository
}

func NewGetCommentsBySenderHandler(comments domain.MiddlemanRepository) GetCommentsBySenderHandler {
	return GetCommentsBySenderHandler{comments: comments}
}

func (h GetCommentsBySenderHandler) GetCommentsBySender(ctx context.Context, query GetCommentsBySender) ([]*domain.MiddlemanComment, error) {
	return h.comments.FindBySenderID(ctx, query.SenderID)
}
