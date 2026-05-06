package commands

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"

	"github.com/stackus/errors"
)

type CreateUserAssistant struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type CreateUserAssistantHandler struct {
	assistants     domain.AssistantRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewCreateUserAssistantHandler(
	assistants domain.AssistantRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) CreateUserAssistantHandler {
	return CreateUserAssistantHandler{
		assistants:     assistants,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h CreateUserAssistantHandler) CreateUserAssistant(ctx context.Context, cmd CreateUserAssistant) error {
	// Load assistant aggregate
	assistant, err := h.assistants.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading assistant for creation")
	}

	// Validate assistant is not already created
	if assistant.Version() > 0 {
		return errors.Wrap(errors.ErrBadRequest, "assistant already exists")
	}

	// Get system prompt (standard users use the base prompt without enhancements)
	systemPrompt := h.promptProvider.GetCompleteSystemPrompt()

	// Create user assistant with predefined configuration
	event, err := assistant.CreateAssistant(
		cmd.ID,
		"Vaver",
		"AI-powered marketplace assistant for intelligent search and recommendations",
		cmd.UserID,
		domain.AssistantTypeStandard,
		[]domain.AssistantCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
		},
		0.7,    // temperature
		4000,   // maxTokens
		systemPrompt,
	)
	if err != nil {
		return err
	}

	// Save assistant aggregate
	if err = h.assistants.Save(ctx, assistant); err != nil {
		return errors.Wrap(err, "error saving user assistant")
	}

	// Publish domain event
	return h.publisher.Publish(ctx, event)
}