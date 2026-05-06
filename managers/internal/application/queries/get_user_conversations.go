package queries

import (
	"context"
	"middleman/managers/internal/domain"
)

type GetUserConversations struct {
	UserID     string
	ActiveOnly bool
	Limit      int
	Offset     int
}

type GetUserConversationsResult struct {
	Conversations []*domain.ReadConversation
	TotalCount    int
}

type GetUserConversationsHandler struct {
	readRepo domain.ReadConversationRepository
}

func NewGetUserConversationsHandler(readRepo domain.ReadConversationRepository) GetUserConversationsHandler {
	return GetUserConversationsHandler{
		readRepo: readRepo,
	}
}

func (h GetUserConversationsHandler) GetUserConversations(ctx context.Context, query GetUserConversations) (*GetUserConversationsResult, error) {
	// Set default pagination if not provided
	limit := query.Limit
	if limit <= 0 {
		limit = 20 // Default page size
	}
	if limit > 100 {
		limit = 100 // Max page size
	}

	offset := query.Offset
	if offset < 0 {
		offset = 0
	}

	conversations, totalCount, err := h.readRepo.GetUserConversations(ctx, query.UserID, query.ActiveOnly, limit, offset)
	if err != nil {
		return nil, err
	}

	return &GetUserConversationsResult{
		Conversations: conversations,
		TotalCount:    totalCount,
	}, nil
}
