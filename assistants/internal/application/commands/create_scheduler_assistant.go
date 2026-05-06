package commands

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"

	"github.com/stackus/errors"
)

type CreateSchedulerAssistant struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type CreateSchedulerAssistantHandler struct {
	assistants     domain.AssistantRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewCreateSchedulerAssistantHandler(
	assistants domain.AssistantRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) CreateSchedulerAssistantHandler {
	return CreateSchedulerAssistantHandler{
		assistants:     assistants,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h CreateSchedulerAssistantHandler) CreateSchedulerAssistant(ctx context.Context, cmd CreateSchedulerAssistant) error {
	// Load assistant aggregate
	assistant, err := h.assistants.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading assistant for creation")
	}

	// Validate assistant is not already created
	if assistant.Version() > 0 {
		return errors.Wrap(errors.ErrBadRequest, "assistant already exists")
	}

	// Get base system prompt and enhance it for scheduler role
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

	// Create scheduler assistant with predefined configuration
	event, err := assistant.CreateAssistant(
		cmd.ID,
		"Scheduler Assistant",
		"Automated assistant for executing scheduled tasks and workflows",
		cmd.UserID,
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
		return err
	}

	// Save assistant aggregate
	if err = h.assistants.Save(ctx, assistant); err != nil {
		return errors.Wrap(err, "error saving scheduler assistant")
	}

	// Publish domain event
	return h.publisher.Publish(ctx, event)
}