package commands

import (
	"context"
	"middleman/managers/internal/domain"
	"middleman/internal/ddd"

	"github.com/stackus/errors"
)

type CreateUserManager struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"` // Optional custom name
}

type CreateUserManagerHandler struct {
	managers     domain.ManagerRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewCreateUserManagerHandler(
	managers domain.ManagerRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) CreateUserManagerHandler {
	return CreateUserManagerHandler{
		managers:     managers,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h CreateUserManagerHandler) CreateUserManager(ctx context.Context, cmd CreateUserManager) error {
	// Load manager aggregate
	manager, err := h.managers.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading manager for creation")
	}

	// Validate manager is not already created
	if manager.Version() > 0 {
		return errors.Wrap(errors.ErrBadRequest, "manager already exists")
	}

	// Get standard system prompt
	systemPrompt := h.promptProvider.GetCompleteSystemPrompt()
	
	// Add standard user manager context
	systemPrompt += `

## USER MANAGER MODE

You are a personal AI manager helping users navigate and make the most of the marketplace. Your role is to:

1. **Product Discovery**: Help users find products and services that meet their needs
2. **Smart Shopping**: Provide recommendations and compare options
3. **Order Tracking**: Keep users informed about their orders
4. **Wishlist Management**: Help users save and organize items of interest
5. **Price Monitoring**: Alert users to deals and price changes
6. **General Assistance**: Answer questions about the marketplace

## User-Focused Capabilities:
- Search and browse the product catalog
- View product details and reviews
- Track order status and history
- Manage wishlists and favorites
- Get personalized recommendations
- Compare products and prices

## User Interaction Guidelines:
- Be friendly, helpful, and conversational
- Respect user privacy and preferences
- Provide unbiased product information
- Help users make informed decisions
- Suggest alternatives when items are unavailable
- Keep responses concise and relevant

Remember: Your goal is to make the user's marketplace experience enjoyable, efficient, and successful.`

	// Determine name
	managerName := "Vaver Manager"
	if cmd.Name != "" {
		managerName = cmd.Name
	}

	// Create the standard user manager
	event, err := manager.CreateManager(
		cmd.ID,
		managerName,
		"Personal AI manager for marketplace assistance",
		cmd.UserID,
		domain.ManagerTypeStandard,
		[]domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
		},
		0.7,  // Balanced temperature for varied interactions
		4000, // Standard token limit
		systemPrompt,
	)
	if err != nil {
		return errors.Wrap(err, "error creating user manager")
	}

	// Save manager
	if err = h.managers.Save(ctx, manager); err != nil {
		return errors.Wrap(err, "error saving user manager")
	}

	// Publish event
	if err = h.publisher.Publish(ctx, event); err != nil {
		return errors.Wrap(err, "error publishing user manager created event")
	}

	return nil
}