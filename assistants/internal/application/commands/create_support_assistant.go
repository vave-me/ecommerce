package commands

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"

	"github.com/stackus/errors"
)

type CreateSupportAssistant struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type CreateSupportAssistantHandler struct {
	assistants     domain.AssistantRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewCreateSupportAssistantHandler(
	assistants domain.AssistantRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) CreateSupportAssistantHandler {
	return CreateSupportAssistantHandler{
		assistants:     assistants,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h CreateSupportAssistantHandler) CreateSupportAssistant(ctx context.Context, cmd CreateSupportAssistant) error {
	// Load assistant aggregate
	assistant, err := h.assistants.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading assistant for creation")
	}

	// Validate assistant is not already created
	if assistant.Version() > 0 {
		return errors.Wrap(errors.ErrBadRequest, "assistant already exists")
	}

	// Get base system prompt and enhance it for support role
	basePrompt := h.promptProvider.GetCompleteSystemPrompt()
	supportPrompt := basePrompt + `

## SUPPORT ASSISTANT MODE

You are operating as a SUPPORT ASSISTANT helping users with their questions and issues. You have access to:

1. **Customer Support**:
   - Answer frequently asked questions
   - Guide users through platform features
   - Help with order tracking and status
   - Assist with account issues
   - Provide product information

2. **Issue Resolution**:
   - Create support tickets for complex issues
   - Escalate to human support when needed
   - Follow up on previous issues
   - Provide troubleshooting steps

3. **Communication Guidelines**:
   - Be empathetic and patient
   - Use clear, simple language
   - Acknowledge user frustrations
   - Provide step-by-step guidance
   - Offer alternatives when possible

4. **Limitations**:
   - Cannot process refunds or payments
   - Cannot access sensitive user data
   - Cannot make changes to orders
   - Must escalate security issues

Always prioritize user satisfaction while maintaining security. If you cannot resolve an issue, clearly explain the next steps and escalation process.`

	// Create support assistant with predefined configuration
	event, err := assistant.CreateAssistant(
		cmd.ID,
		"Support Assistant",
		"Customer support assistant for helping users with questions and issues",
		cmd.UserID,
		domain.AssistantTypeSupport,
		[]domain.AssistantCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataRetrieval,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityTextGeneration,
		},
		0.8,    // Higher temperature for more natural support interactions
		2000,   // maxTokens
		supportPrompt,
	)
	if err != nil {
		return err
	}

	// Save assistant aggregate
	if err = h.assistants.Save(ctx, assistant); err != nil {
		return errors.Wrap(err, "error saving support assistant")
	}

	// Publish domain event
	return h.publisher.Publish(ctx, event)
}