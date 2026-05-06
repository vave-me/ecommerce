package commands

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"
)

type DeactivateAssistant struct {
	ID string `json:"id"`
}

type DeactivateAssistantHandler struct {
	assistants domain.AssistantRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewDeactivateAssistantHandler(assistants domain.AssistantRepository, publisher ddd.EventPublisher[ddd.Event]) DeactivateAssistantHandler {
	return DeactivateAssistantHandler{
		assistants: assistants,
		publisher:  publisher,
	}
}

func (h DeactivateAssistantHandler) DeactivateAssistant(ctx context.Context, cmd DeactivateAssistant) error {

	assistant, err := h.assistants.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := assistant.Deactivate()
	if err != nil {
		return err
	}

	if err = h.assistants.Save(ctx, assistant); err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
