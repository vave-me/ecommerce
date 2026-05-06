package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/messages/internal/domain"
)

type StartConversation struct {
	ID          string
	SenderID    string
	RecipientID string
	ItemID      string
}

type StartConversationHandler struct {
	conversations domain.ConversationRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewStartConversationHandler(conversations domain.ConversationRepository, publisher ddd.EventPublisher[ddd.Event],
) StartConversationHandler {
	return StartConversationHandler{
		conversations: conversations,
		publisher:     publisher,
	}
}
func (h StartConversationHandler) StartConversation(ctx context.Context, cmd StartConversation) error {
	conversation, err := h.conversations.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := conversation.InitConversation(cmd.ID, cmd.SenderID, cmd.RecipientID, cmd.ItemID)
	if err != nil {
		return err
	}

	err = h.conversations.Save(ctx, conversation)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
