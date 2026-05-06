package commands

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"

	"github.com/stackus/errors"
)

type CreateBusinessAssistant struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type CreateBusinessAssistantHandler struct {
	assistants     domain.AssistantRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewCreateBusinessAssistantHandler(
	assistants domain.AssistantRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) CreateBusinessAssistantHandler {
	return CreateBusinessAssistantHandler{
		assistants:     assistants,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h CreateBusinessAssistantHandler) CreateBusinessAssistant(ctx context.Context, cmd CreateBusinessAssistant) error {
	// Load assistant aggregate
	assistant, err := h.assistants.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading assistant for creation")
	}

	// Validate assistant is not already created
	if assistant.Version() > 0 {
		return errors.Wrap(errors.ErrBadRequest, "assistant already exists")
	}

	// Get base system prompt and enhance it for business role
	basePrompt := h.promptProvider.GetCompleteSystemPrompt()
	businessPrompt := basePrompt + `

## BUSINESS ASSISTANT MODE

You are operating as a BUSINESS ASSISTANT helping vendors manage their marketplace presence. You have access to:

1. **Business Management**:
   - Manage your products and services (add, edit, remove)
   - View and process your orders
   - Respond to customer inquiries
   - Update business information and settings
   - Manage inventory and pricing

2. **Analytics & Insights**:
   - View your sales performance
   - Analyze customer behavior for your products
   - Track conversion rates and metrics
   - Monitor competitor pricing (public data only)
   - Generate business reports

3. **Customer Relations**:
   - Respond to reviews and ratings
   - Handle customer support for your products
   - Send updates to customers
   - Manage returns and refunds for your orders

4. **Restrictions**:
   - Can only access YOUR business data
   - Cannot view other vendors' private information
   - Cannot modify platform-wide settings
   - Cannot access customer payment details

Focus on helping vendors grow their business, improve customer satisfaction, and optimize their marketplace performance.`

	// Create business assistant with predefined configuration
	event, err := assistant.CreateAssistant(
		cmd.ID,
		"Business Assistant",
		"Vendor assistant for managing products, orders, and customer relationships",
		cmd.UserID,
		domain.AssistantTypeBusiness,
		[]domain.AssistantCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
			domain.CapabilityTextGeneration, // For creating product descriptions
		},
		0.7,    // temperature
		4000,   // maxTokens
		businessPrompt,
	)
	if err != nil {
		return err
	}

	// Save assistant aggregate
	if err = h.assistants.Save(ctx, assistant); err != nil {
		return errors.Wrap(err, "error saving business assistant")
	}

	// Publish domain event
	return h.publisher.Publish(ctx, event)
}