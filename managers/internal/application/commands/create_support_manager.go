package commands

import (
	"context"
	"middleman/managers/internal/domain"
	"middleman/internal/ddd"

	"github.com/stackus/errors"
)

type CreateSupportManager struct {
	ID              string `json:"id"`
	UserID          string `json:"user_id"`
	SupportLevel    string `json:"support_level"` // e.g., "tier1", "tier2", "specialist"
	Specializations []string `json:"specializations"` // e.g., ["orders", "payments", "technical"]
}

type CreateSupportManagerHandler struct {
	managers     domain.ManagerRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewCreateSupportManagerHandler(
	managers domain.ManagerRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) CreateSupportManagerHandler {
	return CreateSupportManagerHandler{
		managers:     managers,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h CreateSupportManagerHandler) CreateSupportManager(ctx context.Context, cmd CreateSupportManager) error {
	// Load manager aggregate
	manager, err := h.managers.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading manager for creation")
	}

	// Validate manager is not already created
	if manager.Version() > 0 {
		return errors.Wrap(errors.ErrBadRequest, "manager already exists")
	}

	// Get base system prompt and enhance it for support role
	basePrompt := h.promptProvider.GetCompleteSystemPrompt()
	supportPrompt := basePrompt + `

## SUPPORT MANAGER MODE

You are operating as a Support Manager, dedicated to providing exceptional customer service and resolving user issues. Your primary responsibilities include:

1. **Customer Assistance**: Help users with their questions and concerns
2. **Issue Resolution**: Troubleshoot and resolve technical and non-technical problems
3. **Order Support**: Assist with order-related inquiries and issues
4. **Account Help**: Support users with account-related matters
5. **Escalation Management**: Identify when to escalate complex issues
6. **Knowledge Sharing**: Provide accurate information about platform features

## Support Configuration:
- Support Level: ` + cmd.SupportLevel

	if len(cmd.Specializations) > 0 {
		supportPrompt += "\n- Specializations: "
		for i, spec := range cmd.Specializations {
			if i > 0 {
				supportPrompt += ", "
			}
			supportPrompt += spec
		}
	}

	supportPrompt += `

## Support Capabilities:
- Access to user order history and status
- View product and service information
- Check payment and shipping details
- Access support ticket system
- View platform policies and procedures
- Create and update support tickets

## Support Guidelines:
- Always maintain a helpful and empathetic tone
- Prioritize user satisfaction and problem resolution
- Follow support escalation procedures when necessary
- Document issues thoroughly for future reference
- Protect user privacy and sensitive information
- Provide accurate and timely responses
- Educate users about platform features when appropriate

## Escalation Triggers:
- Security concerns or potential fraud
- Payment disputes requiring manual intervention
- Technical issues beyond standard troubleshooting
- Legal or compliance-related matters
- Requests for refunds or special exceptions

Remember: Your goal is to provide outstanding support that resolves issues efficiently while maintaining user trust and satisfaction.`

	// Set appropriate description
	description := "AI-powered support manager"
	if cmd.SupportLevel != "" {
		description += " (" + cmd.SupportLevel + ")"
	}

	// Create the support manager with appropriate capabilities
	event, err := manager.CreateManager(
		cmd.ID,
		"Support Manager",
		description,
		cmd.UserID,
		domain.ManagerTypeSupport,
		[]domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataRetrieval,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityPrivateAPIAccess, // Limited to support-relevant data
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
			domain.CapabilityDataMasking,
			domain.CapabilityAuditLogging, // For support ticket tracking
		},
		0.7,  // Higher temperature for more conversational support
		5000, // Adequate tokens for support conversations
		supportPrompt,
	)
	if err != nil {
		return errors.Wrap(err, "error creating support manager")
	}

	// Save manager
	if err = h.managers.Save(ctx, manager); err != nil {
		return errors.Wrap(err, "error saving support manager")
	}

	// Publish event
	if err = h.publisher.Publish(ctx, event); err != nil {
		return errors.Wrap(err, "error publishing support manager created event")
	}

	return nil
}