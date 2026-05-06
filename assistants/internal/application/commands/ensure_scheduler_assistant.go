package commands

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"
	
	"github.com/stackus/errors"
)

// EnsureSchedulerAssistant ensures the scheduler assistant exists
type EnsureSchedulerAssistant struct {
	AssistantID string
}

type EnsureSchedulerAssistantHandler struct {
	assistants     domain.AssistantRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewEnsureSchedulerAssistantHandler(
	assistants domain.AssistantRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) EnsureSchedulerAssistantHandler {
	return EnsureSchedulerAssistantHandler{
		assistants:     assistants,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h EnsureSchedulerAssistantHandler) EnsureSchedulerAssistant(ctx context.Context, cmd EnsureSchedulerAssistant) error {
	// Try to load the assistant
	assistant, err := h.assistants.Load(ctx, cmd.AssistantID)
	if err != nil {
		return errors.Wrap(err, "error loading scheduler assistant")
	}
	
	// Check if assistant already exists
	if assistant.Version() > 0 {
		// Assistant exists, check if it's the right type
		if assistant.Type != domain.AssistantTypeScheduler {
			return errors.Wrap(errors.ErrBadRequest, "assistant exists but is not a scheduler type")
		}
		return nil // Assistant already exists and is correct type
	}
	
	// Create scheduler assistant with predefined configuration
	basePrompt := h.promptProvider.GetCompleteSystemPrompt()
	schedulerPrompt := basePrompt + `

## SCHEDULER ASSISTANT MODE

You are operating as a SCHEDULER ASSISTANT executing scheduled tasks. You have access to:

1. **Task Execution**:
   - Execute natural language tasks from users
   - Process scheduled operations
   - Generate reports at specified times
   - Send notifications and reminders
   - Perform data analysis tasks

2. **User Context**:
   - Access user-specific data for task execution
   - Maintain user preferences
   - Execute tasks with user permissions
   - Respect privacy settings

3. **Execution Guidelines**:
   - Clearly report task results
   - Handle errors gracefully
   - Provide execution status
   - Log important actions
   - Notify users of completion

4. **Special Considerations**:
   - Tasks may be time-sensitive
   - Maintain execution history
   - Handle recurring tasks
   - Respect resource limits

Focus on reliable, accurate task execution with clear result reporting.`
	
	event, err := assistant.CreateAssistant(
		cmd.AssistantID,
		"Scheduler Assistant",
		"Automated assistant for executing scheduled tasks and workflows",
		"system", // System user owns the scheduler assistant
		domain.AssistantTypeScheduler,
		[]domain.AssistantCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityPrivateAPIAccess, // Needed for user-specific scheduled tasks
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
			domain.CapabilityTextGeneration,
		},
		0.5,    // Balanced for consistent execution
		4000,   // maxTokens
		schedulerPrompt,
	)
	if err != nil {
		return errors.Wrap(err, "error creating scheduler assistant")
	}
	
	// Save the assistant
	if err = h.assistants.Save(ctx, assistant); err != nil {
		return errors.Wrap(err, "error saving scheduler assistant")
	}
	
	// Publish the event
	if err = h.publisher.Publish(ctx, event); err != nil {
		return errors.Wrap(err, "error publishing scheduler assistant created event")
	}
	
	return nil
}