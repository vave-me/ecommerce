package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"middleman/internal/auth"
	"middleman/internal/ddd"
	"middleman/managers/internal/application/services"
	"middleman/managers/internal/domain"
)

// ProcessUserInput represents the command to process user input for an manager.
type ProcessUserInput struct {
	ID          string                 `json:"id"`
	ManagerID   string                 `json:"manager_id"`
	UserID      string                 `json:"user_id"`
	Message     string                 `json:"message"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Timestamp   time.Time              `json:"timestamp,omitempty"`
	RequestType string                 `json:"request_type,omitempty"`
}

// ProcessUserInputResult holds the structured result of processing user input.
type ProcessUserInputResult struct {
	ResponseID         string                 `json:"response_id"`
	ResponseMessage    string                 `json:"response_message"`
	ResponseStatus     string                 `json:"response_status"`
	ResponseConfidence float64                `json:"response_confidence"`
	ResponseTimestamp  time.Time              `json:"response_timestamp"`
	ExecutedActions    []domain.ManagerAction `json:"executed_actions,omitempty"`
}

// ProcessUserInputHandler orchestrates the processing of user input.
type ProcessUserInputHandler struct {
	managers       domain.ManagerRepository
	publisher      ddd.EventPublisher[ddd.Event]
	llmProcessor   services.LLMProcessor
	promptProvider domain.SystemPromptProvider
	// Note: We don't inject individual repositories here anymore
	// The LLMProcessor (ProductionToolRegistry) handles all tool execution internally
}

// NewProcessUserInputHandler creates a new handler with the clean workflow.
func NewProcessUserInputHandler(
	managers domain.ManagerRepository,
	llmProcessor services.LLMProcessor,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) (ProcessUserInputHandler, error) {
	if managers == nil {
		return ProcessUserInputHandler{}, errors.New("managers repository cannot be nil")
	}
	if publisher == nil {
		return ProcessUserInputHandler{}, errors.New("event publisher cannot be nil")
	}
	if llmProcessor == nil {
		return ProcessUserInputHandler{}, errors.New("LLM processor cannot be nil")
	}
	return ProcessUserInputHandler{
		managers:       managers,
		publisher:      publisher,
		llmProcessor:   llmProcessor,
		promptProvider: promptProvider,
	}, nil
}

// ProcessUserInput handles the user's input command using the expected workflow.
func (h ProcessUserInputHandler) ProcessUserInput(ctx context.Context, cmd ProcessUserInput) (*ProcessUserInputResult, error) {

	// Load the manager
	manager, err := h.managers.Load(ctx, cmd.ManagerID)
	if err != nil {
		return nil, fmt.Errorf("failed to load manager: %w", err)
	}

	// Check if this is a new manager that needs creation
	if manager.Version() == 0 && cmd.RequestType == "new_manager_interaction" {

		// Determine manager type based on user role from JWT claims
		managerType := domain.ManagerTypeStandard // Default

		// Extract role from JWT claims
		if claims, ok := auth.ClaimsFromContext(ctx); ok && claims.Role != "" {
			// Map user role to manager type
			switch claims.Role {
			case "admin", "superadmin":
				managerType = domain.ManagerTypeAdmin
			case "business", "vendor":
				managerType = domain.ManagerTypeBusiness
			case "support":
				managerType = domain.ManagerTypeSupport
			default:
				managerType = domain.ManagerTypeStandard
			}
		}

		// Get capabilities for the manager type
		capabilities := getCapabilitiesForType(managerType)

		event, createErr := manager.CreateManager(
			cmd.ManagerID,
			"Vaver",
			"AI-powered marketplace manager",
			cmd.UserID,
			managerType,
			capabilities,
			0.7,
			4000,
			h.promptProvider.GetCompleteSystemPrompt(),
		)
		if createErr != nil {
			return nil, fmt.Errorf("failed to create manager: %w", createErr)
		}

		if err := h.managers.Save(ctx, manager); err != nil {
			return nil, fmt.Errorf("failed to save new manager: %w", err)
		}
		if err := h.publisher.Publish(ctx, event); err != nil {
			// Non-critical error, continue
		}
	}

	response, actions, confidence, err := h.llmProcessor.ProcessWithHistory(
		ctx,
		manager,
		cmd.Message,
		[]domain.ConversationMessage{}, // No history for simple processing
		cmd.Context,
	)
	if err != nil {
		return nil, fmt.Errorf("LLM processing failed: %w", err)
	}

	// Save manager state
	if err := h.managers.Save(ctx, manager); err != nil {
		// Non-critical error, continue
	}

	// Return the result
	result := &ProcessUserInputResult{
		ResponseID:         cmd.ID,
		ResponseMessage:    response,
		ResponseStatus:     "completed",
		ResponseConfidence: confidence,
		ResponseTimestamp:  time.Now(),
		ExecutedActions:    actions, // Actions that were executed during LLM processing
	}

	return result, nil
}

// getCapabilitiesForType returns the appropriate capabilities for an manager type
func getCapabilitiesForType(managerType domain.ManagerType) []domain.ManagerCapability {
	switch managerType {
	case domain.ManagerTypeAdmin:
		return []domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
			domain.CapabilityPrivateAPIAccess,
			domain.CapabilityUserDataAccess,
			domain.CapabilityTokenManagement,
			domain.CapabilityDataMasking,
			domain.CapabilityAuditLogging,
			domain.CapabilityManagerManagement,
			domain.CapabilityAuthentication,
			domain.CapabilitySystemConfiguration,
			domain.CapabilityTextGeneration,
			domain.CapabilityCodeGeneration,
		}
	case domain.ManagerTypeBusiness:
		return []domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
			domain.CapabilityTextGeneration,
		}
	case domain.ManagerTypeSupport:
		return []domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataRetrieval,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityTextGeneration,
		}
	case domain.ManagerTypeScheduler:
		return []domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityPrivateAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
			domain.CapabilityTextGeneration,
		}
	default: // ManagerTypeStandard
		return []domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
		}
	}
}
