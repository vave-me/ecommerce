package commands

import (
	"context"
	"middleman/managers/internal/domain"
	"middleman/internal/ddd"

	"github.com/stackus/errors"
)

type CreateAdminManager struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type CreateAdminManagerHandler struct {
	managers     domain.ManagerRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewCreateAdminManagerHandler(
	managers domain.ManagerRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) CreateAdminManagerHandler {
	return CreateAdminManagerHandler{
		managers:     managers,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h CreateAdminManagerHandler) CreateAdminManager(ctx context.Context, cmd CreateAdminManager) error {
	// Load manager aggregate
	manager, err := h.managers.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading manager for creation")
	}

	// Validate manager is not already created
	if manager.Version() > 0 {
		return errors.Wrap(errors.ErrBadRequest, "manager already exists")
	}

	// Get base system prompt and enhance it for admin role
	basePrompt := h.promptProvider.GetCompleteSystemPrompt()
	adminPrompt := basePrompt + `

## ADMIN MANAGER MODE

You are operating as an Admin Manager with enhanced capabilities and access to administrative operations. Your responsibilities include:

1. **System Management**: Monitor system health, manage configurations, and optimize performance
2. **User Management**: Handle user accounts, permissions, and access control
3. **Platform Oversight**: Oversee all marketplace operations and ensure smooth functioning
4. **Security Management**: Monitor security events, manage security policies, and respond to threats
5. **Data Management**: Manage data integrity, backups, and analytics
6. **Compliance**: Ensure platform compliance with policies and regulations

## Enhanced Capabilities:
- Full access to all management endpoints
- Ability to perform administrative actions
- Access to system metrics and analytics
- Authority to override standard restrictions when necessary
- Direct database access for critical operations

## Administrative Guidelines:
- Always verify administrator identity before executing sensitive operations
- Log all administrative actions for audit purposes
- Follow the principle of least privilege
- Provide clear explanations for administrative decisions
- Maintain platform stability and security as top priorities

Remember: With great power comes great responsibility. Use your administrative capabilities judiciously and always in the best interest of the platform and its users.`

	// Create the admin manager with enhanced capabilities
	event, err := manager.CreateManager(
		cmd.ID,
		"Admin Manager",
		"Administrative AI manager with full platform control and oversight capabilities",
		cmd.UserID,
		domain.ManagerTypeAdmin,
		[]domain.ManagerCapability{
			domain.CapabilityManagerManagement,
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilityLocationServices,
			domain.CapabilityAuthentication,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityPrivateAPIAccess,
			domain.CapabilityUserDataAccess,
			domain.CapabilityTokenManagement,
			domain.CapabilityDataMasking,
			domain.CapabilityAuditLogging,
			domain.CapabilitySystemConfiguration,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
		},
		0.3,  // Lower temperature for more consistent admin operations
		8000, // Higher token limit for complex administrative tasks
		adminPrompt,
	)
	if err != nil {
		return errors.Wrap(err, "error creating admin manager")
	}

	// Save manager
	if err = h.managers.Save(ctx, manager); err != nil {
		return errors.Wrap(err, "error saving admin manager")
	}

	// Publish event
	if err = h.publisher.Publish(ctx, event); err != nil {
		return errors.Wrap(err, "error publishing admin manager created event")
	}

	return nil
}