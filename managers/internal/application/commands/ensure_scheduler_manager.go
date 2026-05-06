package commands

import (
	"context"
	"middleman/managers/internal/domain"
	"middleman/internal/ddd"
	
	"github.com/stackus/errors"
)

// EnsureSchedulerManager ensures the scheduler manager exists
type EnsureSchedulerManager struct {
	ManagerID string
}

type EnsureSchedulerManagerHandler struct {
	managers     domain.ManagerRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewEnsureSchedulerManagerHandler(
	managers domain.ManagerRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) EnsureSchedulerManagerHandler {
	return EnsureSchedulerManagerHandler{
		managers:     managers,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h EnsureSchedulerManagerHandler) EnsureSchedulerManager(ctx context.Context, cmd EnsureSchedulerManager) error {
	// Try to load the manager
	manager, err := h.managers.Load(ctx, cmd.ManagerID)
	if err != nil {
		return errors.Wrap(err, "error loading scheduler manager")
	}
	
	// Check if manager already exists
	if manager.Version() > 0 {
		// Manager exists, check if it's the right type
		if manager.Type != domain.ManagerTypeScheduler {
			return errors.Wrap(errors.ErrBadRequest, "manager exists but is not a scheduler type")
		}
		return nil // Manager already exists and is correct type
	}
	
	// Create scheduler manager with predefined configuration
	systemUserID := "system" // Default system user
	
	// Get scheduler-specific system prompt
	basePrompt := h.promptProvider.GetCompleteSystemPrompt()
	schedulerPrompt := basePrompt + `

## SCHEDULER MANAGER MODE

You are the system's scheduler manager, responsible for automated task execution and maintenance.

Key responsibilities:
- Execute scheduled tasks reliably
- Perform system maintenance
- Generate reports and analytics
- Clean up expired data
- Monitor system health

Remember: You operate autonomously without user interaction.`

	// Create the scheduler manager
	event, err := manager.CreateManager(
		cmd.ManagerID,
		"System Scheduler",
		"Automated scheduler for system tasks and maintenance",
		systemUserID,
		domain.ManagerTypeScheduler,
		[]domain.ManagerCapability{
			domain.CapabilityDataAnalysis,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPrivateAPIAccess,
			domain.CapabilitySystemConfiguration,
			domain.CapabilityAuditLogging,
		},
		0.1,  // Very low temperature for consistent operations
		4000, // Standard token limit
		schedulerPrompt,
	)
	if err != nil {
		return errors.Wrap(err, "error creating scheduler manager")
	}

	// Save manager
	if err = h.managers.Save(ctx, manager); err != nil {
		return errors.Wrap(err, "error saving scheduler manager")
	}

	// Publish event
	if err = h.publisher.Publish(ctx, event); err != nil {
		return errors.Wrap(err, "error publishing scheduler manager created event")
	}

	return nil
}