package commands

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"
)

type UpdateAssistantConfiguration struct {
	ID           string                       `json:"id"`
	UserID       string                       `json:"user_id"`
	Name         string                       `json:"name,omitempty"`
	Description  string                       `json:"description,omitempty"`
	Temperature  float64                      `json:"temperature,omitempty"`
	MaxTokens    int                          `json:"max_tokens,omitempty"`
	SystemPrompt string                       `json:"system_prompt,omitempty"`
	Capabilities []domain.AssistantCapability `json:"capabilities,omitempty"`
}

type UpdateAssistantConfigurationHandler struct {
	assistants domain.AssistantRepository
	publisher  ddd.EventPublisher[ddd.Event]
}

func NewUpdateAssistantConfigurationHandler(assistants domain.AssistantRepository, publisher ddd.EventPublisher[ddd.Event]) UpdateAssistantConfigurationHandler {
	return UpdateAssistantConfigurationHandler{
		assistants: assistants,
		publisher:  publisher,
	}
}

func (h UpdateAssistantConfigurationHandler) UpdateAssistantConfiguration(ctx context.Context, cmd UpdateAssistantConfiguration) error {
	// Step 1: Load assistant aggregate
	assistant, err := h.assistants.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Step 2: Prepare configuration values with defaults
	temperature := cmd.Temperature
	if temperature == 0 {
		temperature = assistant.Temperature
	}

	maxTokens := cmd.MaxTokens
	if maxTokens == 0 {
		maxTokens = assistant.MaxTokens
	}

	systemPrompt := cmd.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = assistant.SystemPrompt
	}

	// Step 3: Update configuration in aggregate
	var event ddd.Event
	// Check if capabilities are being updated
	if len(cmd.Capabilities) > 0 {
		event, err = assistant.UpdateConfigurationWithCapabilities(temperature, maxTokens, systemPrompt, cmd.Capabilities)
	} else {
		event, err = assistant.UpdateConfiguration(temperature, maxTokens, systemPrompt)
	}
	
	if err != nil {
		return err
	}

	// Step 4: Save assistant aggregate
	if err = h.assistants.Save(ctx, assistant); err != nil {
		return err
	}

	// Step 5: Publish domain event
	err = h.publisher.Publish(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
