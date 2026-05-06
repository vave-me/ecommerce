package commands

import (
	"context"
	"fmt"
	"time"

	"middleman/assistants/internal/application/services"
	"middleman/assistants/internal/domain"
	"middleman/internal/auth"
	"middleman/internal/ddd"
)


// ProcessUserInput represents the command to process user input for an assistant.
type ProcessUserInput struct {
	ID          string                 `json:"id"`
	AssistantID string                 `json:"assistant_id"`
	UserID      string                 `json:"user_id"`
	Message     string                 `json:"message"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Timestamp   time.Time              `json:"timestamp,omitempty"`
	RequestType string                 `json:"request_type,omitempty"`
}

// ProcessUserInputResult holds the structured result of processing user input.
type ProcessUserInputResult struct {
	ResponseID         string                   `json:"response_id"`
	ResponseMessage    string                   `json:"response_message"`
	ResponseStatus     string                   `json:"response_status"`
	ResponseConfidence float64                  `json:"response_confidence"`
	ResponseTimestamp  time.Time                `json:"response_timestamp"`
	ExecutedActions    []domain.AssistantAction `json:"executed_actions,omitempty"`
}

// ProcessUserInputHandler orchestrates the processing of user input.
type ProcessUserInputHandler struct {
	assistants     domain.AssistantRepository
	publisher      ddd.EventPublisher[ddd.Event]
	llmProcessor   services.LLMProcessor
	promptProvider domain.SystemPromptProvider
	// Note: We don't inject individual repositories here anymore
	// The LLMProcessor (ProductionToolRegistry) handles all tool execution internally
}

// NewProcessUserInputHandler creates a new handler with the clean workflow.
func NewProcessUserInputHandler(
	assistants domain.AssistantRepository,
	llmProcessor services.LLMProcessor,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) ProcessUserInputHandler {
	if assistants == nil || publisher == nil {
		panic("Critical dependencies (assistants, publisher) cannot be nil")
	}
	return ProcessUserInputHandler{
		assistants:     assistants,
		publisher:      publisher,
		llmProcessor:   llmProcessor,
		promptProvider: promptProvider,
	}
}

// ProcessUserInput handles the user's input command using the expected workflow.
func (h ProcessUserInputHandler) ProcessUserInput(ctx context.Context, cmd ProcessUserInput) (*ProcessUserInputResult, error) {

	// Load the assistant
	assistant, err := h.assistants.Load(ctx, cmd.AssistantID)
	if err != nil {
		return nil, fmt.Errorf("failed to load assistant: %w", err)
	}

	// Check if this is a new assistant that needs creation
	if assistant.Version() == 0 && cmd.RequestType == "new_assistant_interaction" {

		// Determine assistant type based on user role from JWT claims
		assistantType := domain.AssistantTypeStandard // Default

		// Extract role from JWT claims
		if claims, ok := auth.ClaimsFromContext(ctx); ok && claims.Role != "" {
			// Map user role to assistant type
			switch claims.Role {
			case "admin", "superadmin":
				assistantType = domain.AssistantTypeAdmin
			case "business", "vendor":
				assistantType = domain.AssistantTypeBusiness
			case "support":
				assistantType = domain.AssistantTypeSupport
			default:
				assistantType = domain.AssistantTypeStandard
			}
		}

		// Get capabilities for the assistant type
		capabilities := getCapabilitiesForType(assistantType)

		event, createErr := assistant.CreateAssistant(
			cmd.AssistantID,
			"Vaver",
			"AI-powered marketplace assistant",
			cmd.UserID,
			assistantType,
			capabilities,
			0.7,
			4000,
			h.promptProvider.GetCompleteSystemPrompt(),
		)
		if createErr != nil {
			return nil, fmt.Errorf("failed to create assistant: %w", createErr)
		}

		if err := h.assistants.Save(ctx, assistant); err != nil {
			return nil, fmt.Errorf("failed to save new assistant: %w", err)
		}
		if err := h.publisher.Publish(ctx, event); err != nil {
			// Non-critical error, continue
		}
	}

	response, actions, confidence, err := h.llmProcessor.ProcessWithHistory(
		ctx,
		assistant,
		cmd.Message,
		[]domain.ConversationMessage{}, // No history for simple processing
		cmd.Context,
	)
	if err != nil {
		return nil, fmt.Errorf("LLM processing failed: %w", err)
	}

	// Save assistant state
	if err := h.assistants.Save(ctx, assistant); err != nil {
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

// getCapabilitiesForType returns the appropriate capabilities for an assistant type
func getCapabilitiesForType(assistantType domain.AssistantType) []domain.AssistantCapability {
	switch assistantType {
	case domain.AssistantTypeAdmin:
		return []domain.AssistantCapability{
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
			domain.CapabilityAssistantManagement,
			domain.CapabilityAuthentication,
			domain.CapabilitySystemConfiguration,
			domain.CapabilityTextGeneration,
			domain.CapabilityCodeGeneration,
		}
	case domain.AssistantTypeBusiness:
		return []domain.AssistantCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
			domain.CapabilityTextGeneration,
		}
	case domain.AssistantTypeSupport:
		return []domain.AssistantCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataRetrieval,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityTextGeneration,
		}
	case domain.AssistantTypeScheduler:
		return []domain.AssistantCapability{
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
	default: // AssistantTypeStandard
		return []domain.AssistantCapability{
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
