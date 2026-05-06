package commands

import (
	"context"
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/messages/internal/domain"
)

type SendMessage struct {
	ID             string
	ConversationID string
	SenderID       string
	RecipientID    string
	ItemID         string
	Body           string
	IsRead         bool
}

type SendMessageHandler struct {
	messages  domain.MessageRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewSendMessageHandler(messages domain.MessageRepository, publisher ddd.EventPublisher[ddd.Event]) SendMessageHandler {
	return SendMessageHandler{
		messages:  messages,
		publisher: publisher,
	}
}

func (h SendMessageHandler) SendMessage(ctx context.Context, cmd SendMessage) error {
	message, err := h.messages.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error sending message")
	}

	event, err := message.InitMessage(cmd.ID, cmd.ConversationID, cmd.SenderID, cmd.RecipientID, cmd.ItemID, cmd.Body, cmd.IsRead)
	if err != nil {
		return errors.Wrap(err, "initializing message")
	}

	err = h.messages.Save(ctx, message)
	if err != nil {
		return errors.Wrap(err, "error saving message ")
	}

	return h.publisher.Publish(ctx, event)
}
