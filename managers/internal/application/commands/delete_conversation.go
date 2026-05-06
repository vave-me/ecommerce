package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/managers/internal/domain"
)

// DeleteConversation represents the command to delete a conversation
type DeleteConversation struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id"`
}

type DeleteConversationHandler struct {
	conversations domain.ConversationRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewDeleteConversationHandler(
	conversations domain.ConversationRepository,
	publisher ddd.EventPublisher[ddd.Event],
) DeleteConversationHandler {
	return DeleteConversationHandler{
		conversations: conversations,
		publisher:     publisher,
	}
}

func (h DeleteConversationHandler) DeleteConversation(ctx context.Context, cmd DeleteConversation) error {
	// Load the conversation
	conversation, err := h.conversations.Load(ctx, cmd.ConversationID)
	if err != nil {
		return err
	}

	// Verify ownership
	if conversation.GetUserID() != cmd.UserID {
		return domain.ErrUnauthorized
	}

	// Mark conversation as deleted
	event, err := conversation.Delete()
	if err != nil {
		return err
	}

	// Save the conversation with delete status
	if err := h.conversations.Save(ctx, conversation); err != nil {
		return err
	}

	// Publish the event
	return h.publisher.Publish(ctx, event)
}