package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/managers/internal/domain"
)

type UpdateConversationContext struct {
	ConversationID string                 `json:"conversation_id"`
	UserID         string                 `json:"user_id"`
	Context        map[string]interface{} `json:"context"`
}

type UpdateConversationContextHandler struct {
	conversations     domain.ConversationRepository
	readConversations domain.ReadConversationRepository
	publisher         ddd.EventPublisher[ddd.Event]
}

func NewUpdateConversationContextHandler(
	conversations domain.ConversationRepository,
	readConversations domain.ReadConversationRepository,
	publisher ddd.EventPublisher[ddd.Event],
) UpdateConversationContextHandler {
	return UpdateConversationContextHandler{
		conversations:     conversations,
		readConversations: readConversations,
		publisher:         publisher,
	}
}

func (h UpdateConversationContextHandler) UpdateConversationContext(ctx context.Context, cmd UpdateConversationContext) error {

	_, err := h.readConversations.GetConversation(ctx, cmd.ConversationID, cmd.UserID)
	if err != nil {
		return errors.Wrap(err, "error verifying conversation access")
	}

	// Access verified

	conversation, err := h.conversations.Load(ctx, cmd.ConversationID)
	if err != nil {
		return errors.Wrap(err, "error loading conversation")
	}

	event, err := conversation.UpdateContext(cmd.Context)
	if err != nil {

		return errors.Wrap(err, "error updating conversation context")
	}

	if err = h.conversations.Save(ctx, conversation); err != nil {
		return errors.Wrap(err, "error saving conversation")
	}

	return h.publisher.Publish(ctx, event)
}
