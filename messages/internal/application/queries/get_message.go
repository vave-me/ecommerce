package queries

import (
	"context"
	"middleman/messages/internal/domain"
)

type GetMessage struct {
	MessageID string
}
type GetMessageHandler struct {
	messages domain.MessengerRepository
}

func NewGetMessageHandler(messages domain.MessengerRepository) GetMessageHandler {
	return GetMessageHandler{messages: messages}
}

func (h GetMessageHandler) GetMessage(ctx context.Context, query GetMessage) (*domain.MiddlemanMessage, error) {
	return h.messages.Find(ctx, query.MessageID)
}
