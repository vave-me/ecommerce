package commands

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"
)

type ActivateAssistant struct {
	ID string `json:"id"`
}

type ActivateAssistantHandler struct {
	assistants domain.AssistantRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewActivateAssistantHandler(assistants domain.AssistantRepository, publisher ddd.EventPublisher[ddd.Event]) ActivateAssistantHandler {
	return ActivateAssistantHandler{
		assistants: assistants,
		publisher:  publisher,
	}
}

func (h ActivateAssistantHandler) ActivateAssistant(ctx context.Context, cmd ActivateAssistant) error {
	// Step 1: Load assistant aggregate
	assistant, err := h.assistants.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Step 2: Activate assistant in aggregate
	event, err := assistant.Activate()
	if err != nil {
		return err
	}

	// Step 3: Save assistant aggregate
	if err = h.assistants.Save(ctx, assistant); err != nil {
		return err
	}

	// Step 4: Publish domain event
	err = h.publisher.Publish(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
