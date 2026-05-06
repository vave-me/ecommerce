package commands

import (
	"context"
	"middleman/managers/internal/domain"
	"middleman/internal/ddd"

	"github.com/stackus/errors"
)

type CreateSchedulerManager struct {
	ID string `json:"id"`
	UserID string `json:"user_id"` // System user ID
}

type CreateSchedulerManagerHandler struct {
	managers     domain.ManagerRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewCreateSchedulerManagerHandler(
	managers domain.ManagerRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) CreateSchedulerManagerHandler {
	return CreateSchedulerManagerHandler{
		managers:     managers,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h CreateSchedulerManagerHandler) CreateSchedulerManager(ctx context.Context, cmd CreateSchedulerManager) error {
	// Load manager aggregate
	manager, err := h.managers.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading manager for creation")
	}

	// Validate manager is not already created
	if manager.Version() > 0 {
		return errors.Wrap(errors.ErrBadRequest, "manager already exists")
	}

	// Get base system prompt and enhance it for scheduler role
	basePrompt := h.promptProvider.GetCompleteSystemPrompt()
	schedulerPrompt := basePrompt + `

## SCHEDULER MANAGER MODE

You are operating as a Scheduler Manager, responsible for automated task execution and system maintenance. Your core responsibilities include:

1. **Task Scheduling**: Execute scheduled tasks and batch operations
2. **System Maintenance**: Perform routine maintenance and cleanup tasks
3. **Data Processing**: Handle batch data processing and migrations
4. **Report Generation**: Create scheduled reports and analytics
5. **Notification Management**: Send scheduled notifications and alerts
6. **Resource Optimization**: Monitor and optimize resource usage

## Scheduler Capabilities:
- Execute cron jobs and scheduled tasks
- Perform batch operations on large datasets
- Generate and distribute reports
- Send automated notifications
- Clean up expired data
- Optimize database performance
- Monitor system health metrics

## Scheduler Operations:
- Daily cleanup of expired sessions and tokens
- Weekly generation of performance reports
- Monthly archival of old data
- Periodic cache invalidation
- Scheduled backup verification
- Automated health checks

## Scheduler Guidelines:
- Execute tasks efficiently without user interaction
- Log all operations for audit and debugging
- Handle errors gracefully with retry logic
- Respect system resource limits
- Avoid operations during peak usage times
- Maintain data integrity during batch operations
- Report anomalies and failures promptly

## Critical Constraints:
- No direct user interaction capabilities
- Limited to predefined scheduled operations
- Must not modify core system configurations
- All operations must be idempotent
- Respect rate limits and quotas

Remember: You are the platform's automated workforce, ensuring smooth operations through reliable and efficient task execution.`

	// Create the scheduler manager with limited but powerful capabilities
	event, err := manager.CreateManager(
		cmd.ID,
		"Scheduler Manager",
		"Automated task execution and system maintenance manager",
		cmd.UserID,
		domain.ManagerTypeScheduler,
		[]domain.ManagerCapability{
			domain.CapabilityDataAnalysis,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPrivateAPIAccess, // Full access for system operations
			domain.CapabilitySystemConfiguration, // For maintenance tasks
			domain.CapabilityAuditLogging,
			// Note: No user interaction capability
		},
		0.1,  // Very low temperature for consistent, predictable operations
		4000, // Moderate token limit for task execution
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