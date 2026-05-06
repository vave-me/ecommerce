package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/messages/internal/domain"
)

type DeleteMessage struct {
	ID string
}

type DeleteMessageHandler struct {
	messages  domain.MessageRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDeleteMessageHandler(messages domain.MessageRepository, publisher ddd.EventPublisher[ddd.Event]) DeleteMessageHandler {
	return DeleteMessageHandler{
		messages:  messages,
		publisher: publisher,
	}
}

func (h DeleteMessageHandler) DeleteMessage(ctx context.Context, cmd DeleteMessage) error {
	message, err := h.messages.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := message.Delete()
	if err != nil {
		return err
	}

	err = h.messages.Save(ctx, message)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
