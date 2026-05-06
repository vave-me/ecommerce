package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/messages/internal/domain"
)

type DeleteConversation struct {
	ConversationID string
	UserID         string
}

type DeleteConversationHandler struct {
	conversations domain.ConversationRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewDeleteConversationHandler(conversations domain.ConversationRepository, publisher ddd.EventPublisher[ddd.Event],
) DeleteConversationHandler {
	return DeleteConversationHandler{
		conversations: conversations,
		publisher:     publisher,
	}
}
func (h DeleteConversationHandler) DeleteConversation(ctx context.Context, cmd DeleteConversation) error {
	conversation, err := h.conversations.Load(ctx, cmd.ConversationID)
	if err != nil {
		return err
	}

	event, err := conversation.Delete()
	if err != nil {
		return err
	}

	err = h.conversations.Save(ctx, conversation)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
