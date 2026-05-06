package queries

import (
	"context"
	"middleman/assistants/internal/domain"
	"middleman/internal/ddd"

	"github.com/google/uuid"
	"github.com/stackus/errors"
)

type GetOrCreateUserAssistant struct {
	UserID       string
	UserRole     string // Optional: to determine assistant type
	AssistantType domain.AssistantType // Optional: explicitly request a type
}

type GetOrCreateUserAssistantHandler struct {
	assistants         domain.AssistantRepository
	assistantReadModel domain.CatalogRepository
	publisher          ddd.EventPublisher[ddd.Event]
}

func NewGetOrCreateUserAssistantHandler(
	assistants domain.AssistantRepository,
	assistantReadModel domain.CatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) GetOrCreateUserAssistantHandler {
	return GetOrCreateUserAssistantHandler{
		assistants:         assistants,
		assistantReadModel: assistantReadModel,
		publisher:          publisher,
	}
}

func (h GetOrCreateUserAssistantHandler) GetOrCreateUserAssistant(ctx context.Context, query GetOrCreateUserAssistant) (*domain.Assistant, error) {
	// Determine assistant type based on user role or explicit type
	var assistantType domain.AssistantType
	if query.AssistantType != "" {
		assistantType = query.AssistantType
	} else {
		// Determine type based on user role
		switch query.UserRole {
		case "admin", "superadmin":
			assistantType = domain.AssistantTypeAdmin
		case "vendor", "business":
			assistantType = domain.AssistantTypeBusiness
		case "support":
			assistantType = domain.AssistantTypeSupport
		case "system":
			assistantType = domain.AssistantTypeScheduler
		default:
			assistantType = domain.AssistantTypeStandard
		}
	}

	// Validate user can create this type of assistant
	if query.UserRole != "" && !canUserCreateAssistantType(query.UserRole, assistantType) {
		return nil, errors.Wrap(errors.ErrForbidden, "insufficient permissions to create this assistant type")
	}

	// First try to find existing assistants for the user of the same type
	assistants, err := h.assistantReadModel.FindActiveByUser(ctx, query.UserID)
	if err == nil && len(assistants) > 0 {
		// Look for existing assistant of the requested type
		for _, assistant := range assistants {
			// Load the full assistant aggregate to check type
			fullAssistant, loadErr := h.assistants.Load(ctx, assistant.ID)
			if loadErr == nil && fullAssistant.Type == assistantType {
				return fullAssistant, nil
			}
		}
	}

	// No existing assistant of this type found, create a new one
	assistantID := uuid.New().String()
	assistant, err := h.assistants.Load(ctx, assistantID)
	if err != nil {
		return nil, errors.Wrap(err, "error loading new assistant aggregate")
	}

	// For GetOrCreateUserAssistant, we only create standard assistants
	// Other types should be created through specific commands
	if assistantType != domain.AssistantTypeStandard {
		return nil, errors.Wrap(errors.ErrBadRequest, "GetOrCreateUserAssistant only supports standard assistants")
	}

	// Create the standard assistant
	event, err := assistant.CreateAssistant(
		assistantID,
		"Vaver",
		"AI-powered marketplace assistant for intelligent search and recommendations",
		query.UserID,
		domain.AssistantTypeStandard,
		[]domain.AssistantCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilitySearchAndFilter,
			domain.CapabilityDataRetrieval,
			domain.CapabilityPublicAPIAccess,
			domain.CapabilityJailbreakResistant,
			domain.CapabilityScopeEnforcement,
		},
		0.7,
		4000,
		"", // Will use default system prompt
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating assistant")
	}

	// Save and publish
	if err = h.assistants.Save(ctx, assistant); err != nil {
		return nil, errors.Wrap(err, "error saving assistant")
	}

	if err = h.publisher.Publish(ctx, event); err != nil {
		return nil, errors.Wrap(err, "error publishing assistant created event")
	}

	return assistant, nil
}

// canUserCreateAssistantType checks if a user with given role can create an assistant of specified type
func canUserCreateAssistantType(userRole string, assistantType domain.AssistantType) bool {
	switch assistantType {
	case domain.AssistantTypeAdmin:
		// Only admin users can create admin assistants
		return userRole == "admin" || userRole == "superadmin"
	case domain.AssistantTypeBusiness:
		// Business users and admins can create business assistants
		return userRole == "vendor" || userRole == "business" || userRole == "admin" || userRole == "superadmin"
	case domain.AssistantTypeSupport:
		// Support staff and admins can create support assistants
		return userRole == "support" || userRole == "admin" || userRole == "superadmin"
	case domain.AssistantTypeScheduler:
		// System and admins can create scheduler assistants
		return userRole == "system" || userRole == "admin" || userRole == "superadmin"
	case domain.AssistantTypeStandard:
		// Anyone can create standard assistants
		return true
	default:
		return false
	}
}
