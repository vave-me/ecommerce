package queries

import (
	"context"
	"middleman/managers/internal/domain"
	"middleman/internal/ddd"

	"github.com/google/uuid"
	"github.com/stackus/errors"
)

type GetOrCreateUserManager struct {
	UserID       string
	UserRole     string // Optional: to determine manager type
	ManagerType domain.ManagerType // Optional: explicitly request a type
}

type GetOrCreateUserManagerHandler struct {
	managers         domain.ManagerRepository
	managerReadModel domain.CatalogRepository
	publisher          ddd.EventPublisher[ddd.Event]
}

func NewGetOrCreateUserManagerHandler(
	managers domain.ManagerRepository,
	managerReadModel domain.CatalogRepository,
	publisher ddd.EventPublisher[ddd.Event],
) GetOrCreateUserManagerHandler {
	return GetOrCreateUserManagerHandler{
		managers:         managers,
		managerReadModel: managerReadModel,
		publisher:          publisher,
	}
}

func (h GetOrCreateUserManagerHandler) GetOrCreateUserManager(ctx context.Context, query GetOrCreateUserManager) (*domain.Manager, error) {
	// Determine manager type based on user role or explicit type
	var managerType domain.ManagerType
	if query.ManagerType != "" {
		managerType = query.ManagerType
	} else {
		// Determine type based on user role
		switch query.UserRole {
		case "admin", "superadmin":
			managerType = domain.ManagerTypeAdmin
		case "vendor", "business":
			managerType = domain.ManagerTypeBusiness
		case "support":
			managerType = domain.ManagerTypeSupport
		case "system":
			managerType = domain.ManagerTypeScheduler
		default:
			managerType = domain.ManagerTypeStandard
		}
	}

	// Validate user can create this type of manager
	if query.UserRole != "" && !canUserCreateManagerType(query.UserRole, managerType) {
		return nil, errors.Wrap(errors.ErrForbidden, "insufficient permissions to create this manager type")
	}

	// First try to find existing managers for the user of the same type
	managers, err := h.managerReadModel.FindActiveByUser(ctx, query.UserID)
	if err == nil && len(managers) > 0 {
		// Look for existing manager of the requested type
		for _, manager := range managers {
			// Load the full manager aggregate to check type
			fullManager, loadErr := h.managers.Load(ctx, manager.ID)
			if loadErr == nil && fullManager.Type == managerType {
				return fullManager, nil
			}
		}
	}

	// No existing manager of this type found, create a new one
	managerID := uuid.New().String()
	manager, err := h.managers.Load(ctx, managerID)
	if err != nil {
		return nil, errors.Wrap(err, "error loading new manager aggregate")
	}

	// For GetOrCreateUserManager, we only create standard managers
	// Other types should be created through specific commands
	if managerType != domain.ManagerTypeStandard {
		return nil, errors.Wrap(errors.ErrBadRequest, "GetOrCreateUserManager only supports standard managers")
	}

	// Create the standard manager
	event, err := manager.CreateManager(
		managerID,
		"Vaver Manager",
		"AI-powered marketplace manager for intelligent operations and management",
		query.UserID,
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
		0.7,
		4000,
		"", // Will use default system prompt
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating manager")
	}

	// Save and publish
	if err = h.managers.Save(ctx, manager); err != nil {
		return nil, errors.Wrap(err, "error saving manager")
	}

	if err = h.publisher.Publish(ctx, event); err != nil {
		return nil, errors.Wrap(err, "error publishing manager created event")
	}

	return manager, nil
}

// canUserCreateManagerType checks if a user with given role can create a manager of specified type
func canUserCreateManagerType(userRole string, managerType domain.ManagerType) bool {
	switch managerType {
	case domain.ManagerTypeAdmin:
		// Only admin users can create admin managers
		return userRole == "admin" || userRole == "superadmin"
	case domain.ManagerTypeBusiness:
		// Business users and admins can create business managers
		return userRole == "vendor" || userRole == "business" || userRole == "admin" || userRole == "superadmin"
	case domain.ManagerTypeSupport:
		// Support staff and admins can create support managers
		return userRole == "support" || userRole == "admin" || userRole == "superadmin"
	case domain.ManagerTypeScheduler:
		// System and admins can create scheduler managers
		return userRole == "system" || userRole == "admin" || userRole == "superadmin"
	case domain.ManagerTypeStandard:
		// Anyone can create standard managers
		return true
	default:
		return false
	}
}