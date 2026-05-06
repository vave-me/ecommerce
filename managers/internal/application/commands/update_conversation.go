package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/managers/internal/domain"
)

// UpdateConversation represents the command to update conversation metadata
type UpdateConversation struct {
	ConversationID string                 `json:"conversation_id"`
	UserID         string                 `json:"user_id"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type UpdateConversationHandler struct {
	conversations domain.ConversationRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewUpdateConversationHandler(
	conversations domain.ConversationRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdateConversationHandler {
	return UpdateConversationHandler{
		conversations: conversations,
		publisher:     publisher,
	}
}

func (h UpdateConversationHandler) UpdateConversation(ctx context.Context, cmd UpdateConversation) error {
	// Load the conversation
	conversation, err := h.conversations.Load(ctx, cmd.ConversationID)
	if err != nil {
		return err
	}

	// Update the metadata
	event, err := conversation.UpdateMetadata(cmd.Metadata)
	if err != nil {
		return err
	}

	// Save the conversation
	if err := h.conversations.Save(ctx, conversation); err != nil {
		return err
	}

	// Publish the event
	return h.publisher.Publish(ctx, event)
}