package commands

import (
	"middleman/internal/ddd"
	"middleman/messages/internal/domain"
)

type ArchiveConversation struct {
	ChatID string
	UserID string
}

type ArchiveConversationHandler struct {
	conversations domain.ConversationRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewArchiveConversationHandler(conversations domain.ConversationRepository, publisher ddd.EventPublisher[ddd.Event],
) ArchiveConversationHandler {
	return ArchiveConversationHandler{
		conversations: conversations,
		publisher:     publisher,
	}
}
