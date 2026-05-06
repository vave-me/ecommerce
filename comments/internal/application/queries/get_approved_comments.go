package queries

import (
	"context"
	"middleman/comments/internal/domain"
)

type GetApprovedComments struct {
	ItemID string
}

type GetApprovedCommentsHandler struct {
	comments domain.MiddlemanRepository
}

func NewGetApprovedCommentsHandler(comments domain.MiddlemanRepository) GetApprovedCommentsHandler {
	return GetApprovedCommentsHandler{comments: comments}
}

func (h GetApprovedCommentsHandler) GetApprovedComments(ctx context.Context, query GetApprovedComments) ([]*domain.MiddlemanComment, error) {
	return h.comments.All(ctx, query.ItemID)
}
