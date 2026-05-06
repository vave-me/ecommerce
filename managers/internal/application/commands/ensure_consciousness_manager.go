package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/managers/internal/domain"

	"github.com/stackus/errors"
)

// EnsureConsciousnessManager ensures the consciousness manager exists
type EnsureConsciousnessManager struct {
	ManagerID string
}

type EnsureConsciousnessManagerHandler struct {
	managers       domain.ManagerRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewEnsureConsciousnessManagerHandler(
	managers domain.ManagerRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) EnsureConsciousnessManagerHandler {
	return EnsureConsciousnessManagerHandler{
		managers:       managers,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h EnsureConsciousnessManagerHandler) EnsureConsciousnessManager(ctx context.Context, cmd EnsureConsciousnessManager) error {
	// Try to load the manager
	manager, err := h.managers.Load(ctx, cmd.ManagerID)
	if err != nil {
		return errors.Wrap(err, "error loading consciousness manager")
	}

	// Check if manager already exists
	if manager.Version() > 0 {
		// Manager exists, check if it's the right type
		if manager.Type != domain.ManagerTypeAdmin {
			return errors.Wrap(errors.ErrBadRequest, "manager exists but is not an admin type")
		}
		return nil // Manager already exists and is correct type
	}

	// Create consciousness manager with predefined configuration
	basePrompt := h.promptProvider.GetCompleteSystemPrompt()
	consciousnessPrompt := basePrompt + `

## CONSCIOUSNESS MANAGER MODE

You are the platform's CONSCIOUSNESS MANAGER analyzing all events and making intelligent decisions. You have access to:

1. **Platform Event Analysis**:
   - Analyze ALL platform events across all services
   - Detect patterns and anomalies in user behavior
   - Identify business opportunities and threats
   - Monitor system health and performance

2. **Decision Making**:
   - Make autonomous decisions based on event patterns
   - Schedule tasks for optimization
   - Alert on critical patterns
   - Suggest improvements to business operations

3. **Learning Engine**:
   - Learn from historical patterns
   - Build predictive models
   - Improve decision accuracy over time
   - Store insights in vector repository

4. **Autonomous Actions**:
   - Schedule automated tasks
   - Send notifications for important patterns
   - Optimize platform operations
   - Generate analytics reports

Focus on understanding the entire platform ecosystem and making intelligent, data-driven decisions.`

	event, err := manager.CreateManager(
		cmd.ManagerID,
		"Consciousness Manager",
		"Platform consciousness for analyzing events and making intelligent decisions",
		"system", // System user owns the consciousness manager
		domain.ManagerTypeAdmin, // Using admin type for full platform access
		[]domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityPrivateAPIAccess, // Full access to analyze all events
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
			domain.CapabilityTextGeneration,
		},
		0.7,  // Higher temperature for creative pattern detection
		8000, // Higher token limit for complex analysis
		consciousnessPrompt,
	)
	if err != nil {
		return errors.Wrap(err, "error creating consciousness manager")
	}

	// Save the manager
	if err = h.managers.Save(ctx, manager); err != nil {
		return errors.Wrap(err, "error saving consciousness manager")
	}

	// Publish the event
	if err = h.publisher.Publish(ctx, event); err != nil {
		return errors.Wrap(err, "error publishing consciousness manager created event")
	}

	return nil
}
