package commands

import (
	"context"
	"middleman/internal/ddd"
	"middleman/managers/internal/domain"
)

type UpdateManagerConfiguration struct {
	ID           string                     `json:"id"`
	UserID       string                     `json:"user_id"`
	Name         string                     `json:"name,omitempty"`
	Description  string                     `json:"description,omitempty"`
	Temperature  float64                    `json:"temperature,omitempty"`
	MaxTokens    int                        `json:"max_tokens,omitempty"`
	SystemPrompt string                     `json:"system_prompt,omitempty"`
	Capabilities []domain.ManagerCapability `json:"capabilities,omitempty"`
}

type UpdateManagerConfigurationHandler struct {
	managers  domain.ManagerRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewUpdateManagerConfigurationHandler(managers domain.ManagerRepository, publisher ddd.EventPublisher[ddd.Event]) UpdateManagerConfigurationHandler {
	return UpdateManagerConfigurationHandler{
		managers:  managers,
		publisher: publisher,
	}
}

func (h UpdateManagerConfigurationHandler) UpdateManagerConfiguration(ctx context.Context, cmd UpdateManagerConfiguration) error {
	// Step 1: Load manager aggregate
	manager, err := h.managers.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Step 2: Prepare configuration values with defaults
	temperature := cmd.Temperature
	if temperature == 0 {
		temperature = manager.Temperature
	}

	maxTokens := cmd.MaxTokens
	if maxTokens == 0 {
		maxTokens = manager.MaxTokens
	}

	systemPrompt := cmd.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = manager.SystemPrompt
	}

	// Step 3: Update configuration in aggregate
	var event ddd.Event
	// Check if capabilities are being updated
	if len(cmd.Capabilities) > 0 {
		event, err = manager.UpdateConfigurationWithCapabilities(temperature, maxTokens, systemPrompt, cmd.Capabilities)
	} else {
		event, err = manager.UpdateConfiguration(temperature, maxTokens, systemPrompt)
	}

	if err != nil {
		return err
	}

	// Step 4: Save manager aggregate
	if err = h.managers.Save(ctx, manager); err != nil {
		return err
	}

	// Step 5: Publish domain event
	err = h.publisher.Publish(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
