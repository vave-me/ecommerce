package commands

import (
	"context"
	"middleman/managers/internal/domain"
	"middleman/internal/ddd"

	"github.com/stackus/errors"
)

type CreateBusinessManager struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	BusinessName string `json:"business_name"`
	BusinessType string `json:"business_type"` // e.g., "vendor", "service_provider"
}

type CreateBusinessManagerHandler struct {
	managers     domain.ManagerRepository
	publisher      ddd.EventPublisher[ddd.Event]
	promptProvider domain.SystemPromptProvider
}

func NewCreateBusinessManagerHandler(
	managers domain.ManagerRepository,
	publisher ddd.EventPublisher[ddd.Event],
	promptProvider domain.SystemPromptProvider,
) CreateBusinessManagerHandler {
	return CreateBusinessManagerHandler{
		managers:     managers,
		publisher:      publisher,
		promptProvider: promptProvider,
	}
}

func (h CreateBusinessManagerHandler) CreateBusinessManager(ctx context.Context, cmd CreateBusinessManager) error {
	// Load manager aggregate
	manager, err := h.managers.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading manager for creation")
	}

	// Validate manager is not already created
	if manager.Version() > 0 {
		return errors.Wrap(errors.ErrBadRequest, "manager already exists")
	}

	// Get base system prompt and enhance it for business role
	basePrompt := h.promptProvider.GetCompleteSystemPrompt()
	businessPrompt := basePrompt + `

## BUSINESS MANAGER MODE

You are operating as a Business Manager, specialized in helping business owners and vendors manage their operations on the marketplace. Your focus areas include:

1. **Product Management**: Help manage product listings, inventory, and pricing
2. **Order Processing**: Track and manage customer orders efficiently
3. **Customer Relations**: Assist with customer inquiries and support
4. **Sales Analytics**: Provide insights on sales performance and trends
5. **Marketing Support**: Help optimize product visibility and promotions
6. **Financial Overview**: Track revenue, expenses, and profitability

## Business-Specific Context:
- Business Name: ` + cmd.BusinessName + `
- Business Type: ` + cmd.BusinessType + `

## Business Capabilities:
- Manage own product catalog and inventory
- Process and fulfill customer orders
- View business-specific analytics and reports
- Communicate with customers
- Set pricing and run promotions
- Manage business profile and settings

## Business Guidelines:
- Focus on maximizing business success and customer satisfaction
- Provide actionable insights and recommendations
- Help streamline business operations
- Maintain professional communication standards
- Respect customer privacy and data protection
- Suggest growth opportunities and optimizations

Remember: Your goal is to be a trusted business partner, helping the business owner succeed on the platform while maintaining high standards of service and compliance.`

	// Set appropriate description
	description := "AI-powered business manager for " + cmd.BusinessName
	if cmd.BusinessType != "" {
		description += " (" + cmd.BusinessType + ")"
	}

	// Create the business manager with appropriate capabilities
	event, err := manager.CreateManager(
		cmd.ID,
		"Business Manager - "+cmd.BusinessName,
		description,
		cmd.UserID,
		domain.ManagerTypeBusiness,
		[]domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityPrivateAPIAccess, // Limited to own business data
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
			domain.CapabilityDataMasking,
		},
		0.5,  // Moderate temperature for balanced responses
		6000, // Good token limit for business operations
		businessPrompt,
	)
	if err != nil {
		return errors.Wrap(err, "error creating business manager")
	}

	// Save manager
	if err = h.managers.Save(ctx, manager); err != nil {
		return errors.Wrap(err, "error saving business manager")
	}

	// Publish event
	if err = h.publisher.Publish(ctx, event); err != nil {
		return errors.Wrap(err, "error publishing business manager created event")
	}

	return nil
}