package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"
)

type CreateConversation struct {
	ID             string                 `json:"id"`
	UserID         string                 `json:"user_id"`
	AssistantID    string                 `json:"assistant_id"`
	InitialContext map[string]interface{} `json:"initial_context,omitempty"`
}

type CreateConversationHandler struct {
	conversations domain.ConversationRepository
	publisher     ddd.EventPublisher[ddd.Event]
}

func NewCreateConversationHandler(conversations domain.ConversationRepository,
	publisher ddd.EventPublisher[ddd.Event]) CreateConversationHandler {
	return CreateConversationHandler{
		conversations: conversations,
		publisher:     publisher,
	}
}

func (h CreateConversationHandler) CreateConversation(ctx context.Context, cmd CreateConversation) error {

	conversation, err := h.conversations.Load(ctx, cmd.ID)
	if err != nil {

		return errors.Wrap(err, "error creating conversation")
	}

	event, err := conversation.CreateConversation(
		cmd.ID,
		cmd.UserID,
		cmd.AssistantID,
		cmd.InitialContext,
	)
	if err != nil {
		return err
	}

	// Save the conversation aggregate to persist the event
	if err = h.conversations.Save(ctx, conversation); err != nil {
		return errors.Wrap(err, "error saving conversation")
	}

	return h.publisher.Publish(ctx, event)
}
