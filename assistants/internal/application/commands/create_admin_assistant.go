package commands

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"

	"github.com/stackus/errors"
)

type CreateAdminAssistant struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type CreateAdminAssistantHandler struct {
	assistants     domain.AssistantRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewCreateAdminAssistantHandler(
	assistants domain.AssistantRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) CreateAdminAssistantHandler {
	return CreateAdminAssistantHandler{
		assistants:     assistants,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h CreateAdminAssistantHandler) CreateAdminAssistant(ctx context.Context, cmd CreateAdminAssistant) error {
	// Load assistant aggregate
	assistant, err := h.assistants.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading assistant for creation")
	}

	// Validate assistant is not already created
	if assistant.Version() > 0 {
		return errors.Wrap(errors.ErrBadRequest, "assistant already exists")
	}

	// Get base system prompt and enhance it for admin role
	basePrompt := h.promptProvider.GetCompleteSystemPrompt()
	adminPrompt := basePrompt + `

## ADMIN ASSISTANT MODE

You are operating as an ADMIN ASSISTANT with elevated privileges. You have access to:

1. **Full Platform Administration**:
   - User management (view, edit, ban, suspend, role changes)
   - Order management (view all orders, process refunds, handle disputes)
   - Content moderation (delete/edit any content, manage reports)
   - Platform configuration and settings
   - System health monitoring and analytics

2. **Security Responsibilities**:
   - All actions are logged and audited
   - Follow principle of least privilege
   - Verify user identity before sensitive operations
   - Mask sensitive data when displaying information
   - Report suspicious activities

3. **Best Practices**:
   - Double-check before destructive operations
   - Provide clear explanations for administrative actions
   - Maintain professional communication
   - Document reasons for user sanctions
   - Prioritize platform security and user safety

4. **Restricted Operations**:
   - Financial transactions require additional verification
   - Bulk operations should be confirmed
   - Access to payment details is logged
   - Personal data access is monitored

Remember: With great power comes great responsibility. Use administrative capabilities judiciously and always prioritize user privacy and platform security.`

	// Create admin assistant with predefined configuration
	event, err := assistant.CreateAssistant(
		cmd.ID,
		"Admin Assistant",
		"Administrative assistant with full platform access for system management and moderation",
		cmd.UserID,
		domain.AssistantTypeAdmin,
		[]domain.AssistantCapability{
			// All standard capabilities
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
			// Admin-specific capabilities
			domain.CapabilityPrivateAPIAccess,
			domain.CapabilityUserDataAccess,
			domain.CapabilityTokenManagement,
			domain.CapabilityDataMasking,
			domain.CapabilityAuditLogging,
			domain.CapabilityAssistantManagement,
			domain.CapabilityAuthentication,
			domain.CapabilitySystemConfiguration,
			domain.CapabilityTextGeneration,
			domain.CapabilityCodeGeneration,
		},
		0.3,    // Lower temperature for more consistent admin operations
		8000,   // Higher token limit for detailed reports
		adminPrompt,
	)
	if err != nil {
		return err
	}

	// Save assistant aggregate
	if err = h.assistants.Save(ctx, assistant); err != nil {
		return errors.Wrap(err, "error saving admin assistant")
	}

	// Publish domain event
	return h.publisher.Publish(ctx, event)
}