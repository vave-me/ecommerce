package tools

import (
	"fmt"
	"middleman/managers/internal/domain"
)

// PermissionChecker validates if an manager type can perform an operation
type PermissionChecker struct {
	// adminOnlyOperations defines operations that only admin managers can perform
	adminOnlyOperations map[string][]string

	// businessScopedOperations defines operations that business managers can only perform on their own data
	businessScopedOperations map[string][]string

	// supportReadOnlyOperations defines operations that support managers can only read
	supportReadOnlyOperations map[string][]string

	// restrictedEntities defines entities that certain manager types cannot access
	restrictedEntities map[domain.ManagerType][]string
}

// NewPermissionChecker creates a new permission checker with default rules
func NewPermissionChecker() *PermissionChecker {
	return &PermissionChecker{
		adminOnlyOperations: map[string][]string{
			"users":      {"ban", "suspend", "delete", "update_role", "deactivate", "activate"},
			"orders":     {"force_cancel", "override_payment", "bulk_update"},
			"payments":   {"refund", "reverse", "manual_adjustment"},
			"products":   {"bulk_delete", "force_remove", "moderate"},
			"posts":      {"force_delete", "moderate", "bulk_moderate"},
			"comments":   {"bulk_delete", "moderate_all"},
			"reviews":    {"override", "delete_all"},
			"support":    {"assign_agent", "escalate_to_admin", "view_all_tickets"},
			"metrics":    {"export_all", "view_platform_metrics"},
			"categories": {"create", "delete", "update"},
		},
		businessScopedOperations: map[string][]string{
			"products": {"add", "update", "remove", "adjust_stock", "mark_sold"},
			"services": {"add", "update", "remove", "activate", "deactivate"},
			"orders":   {"view", "update_status", "ship", "complete"},
			"reviews":  {"respond", "flag"},
			"messages": {"send", "view", "archive"},
			"metrics":  {"view_own", "export_own"},
		},
		supportReadOnlyOperations: map[string][]string{
			"products":   {"search", "find", "filter"},
			"services":   {"search", "find", "filter"},
			"users":      {"find", "get_profile"},
			"orders":     {"track", "get_status"},
			"support":    {"create_ticket", "update_ticket", "view_own_tickets"},
			"categories": {"get_categories", "search"},
		},
		restrictedEntities: map[domain.ManagerType][]string{
			domain.ManagerTypeSupport:  {"payments", "shipping", "wishlists"},
			domain.ManagerTypeBusiness: {"payments"}, // Can see payment status but not details
		},
	}
}

// CanPerformOperation checks if an manager type can perform a specific operation on an entity
func (pc *PermissionChecker) CanPerformOperation(
	managerType domain.ManagerType,
	entityType string,
	operation string,
	userID string,
	resourceOwnerID string,
) (bool, error) {
	// Admin managers can do everything
	if managerType == domain.ManagerTypeAdmin {
		return true, nil
	}

	// Check if the entity is restricted for this manager type
	if restrictedEntities, exists := pc.restrictedEntities[managerType]; exists {
		for _, restricted := range restrictedEntities {
			if restricted == entityType {
				return false, fmt.Errorf("manager type %s cannot access %s", managerType, entityType)
			}
		}
	}

	// Check admin-only operations
	if adminOps, exists := pc.adminOnlyOperations[entityType]; exists {
		for _, adminOp := range adminOps {
			if adminOp == operation {
				return false, fmt.Errorf("operation %s on %s requires admin manager", operation, entityType)
			}
		}
	}

	// Handle business manager permissions
	if managerType == domain.ManagerTypeBusiness {
		// Check if it's a business-scoped operation
		if businessOps, exists := pc.businessScopedOperations[entityType]; exists {
			for _, businessOp := range businessOps {
				if businessOp == operation {
					// Business managers can only operate on their own resources
					if userID != resourceOwnerID && resourceOwnerID != "" {
						return false, fmt.Errorf("business manager can only %s their own %s", operation, entityType)
					}
					return true, nil
				}
			}
		}
		// Allow read operations for business managers
		if isReadOperation(operation) {
			return true, nil
		}
	}

	// Handle support manager permissions
	if managerType == domain.ManagerTypeSupport {
		// Check if it's in the allowed read-only operations
		if supportOps, exists := pc.supportReadOnlyOperations[entityType]; exists {
			for _, supportOp := range supportOps {
				if supportOp == operation {
					return true, nil
				}
			}
		}
		return false, fmt.Errorf("support manager cannot perform %s on %s", operation, entityType)
	}

	// Scheduler managers have similar permissions to standard but with user context
	if managerType == domain.ManagerTypeScheduler {
		// Scheduler can perform most operations in user context
		if !isDestructiveOperation(operation) {
			return true, nil
		}
		return false, fmt.Errorf("scheduler manager cannot perform destructive operation %s", operation)
	}

	// Standard managers have default permissions
	if managerType == domain.ManagerTypeStandard {
		// Cannot perform admin operations
		if isAdminOperation(operation) {
			return false, fmt.Errorf("standard manager cannot perform admin operation %s", operation)
		}
		// Can perform most read and basic write operations
		return true, nil
	}

	return false, fmt.Errorf("unknown manager type: %s", managerType)
}

// GetAllowedOperations returns all operations an manager type can perform on an entity
func (pc *PermissionChecker) GetAllowedOperations(
	managerType domain.ManagerType,
	entityType string,
) []string {
	// Admin can do everything
	if managerType == domain.ManagerTypeAdmin {
		// Return all operations (this would need to be comprehensive)
		return []string{"*"} // Wildcard for all operations
	}

	var allowed []string

	// Check restricted entities first
	if restrictedEntities, exists := pc.restrictedEntities[managerType]; exists {
		for _, restricted := range restrictedEntities {
			if restricted == entityType {
				return []string{} // No operations allowed
			}
		}
	}

	switch managerType {
	case domain.ManagerTypeBusiness:
		// Add business-scoped operations
		if ops, exists := pc.businessScopedOperations[entityType]; exists {
			allowed = append(allowed, ops...)
		}
		// Add common read operations
		allowed = append(allowed, getCommonReadOperations()...)

	case domain.ManagerTypeSupport:
		// Only add support read-only operations
		if ops, exists := pc.supportReadOnlyOperations[entityType]; exists {
			allowed = append(allowed, ops...)
		}

	case domain.ManagerTypeScheduler, domain.ManagerTypeStandard:
		// Add all non-admin operations
		allowed = append(allowed, getCommonReadOperations()...)
		allowed = append(allowed, getCommonWriteOperations()...)

		// Remove admin-only operations
		if adminOps, exists := pc.adminOnlyOperations[entityType]; exists {
			allowed = filterOutOperations(allowed, adminOps)
		}
	}

	return allowed
}

// Helper functions

func isReadOperation(operation string) bool {
	readOps := []string{"search", "find", "filter", "get", "view", "list", "track"}
	for _, op := range readOps {
		if op == operation {
			return true
		}
	}
	return false
}

func isDestructiveOperation(operation string) bool {
	destructiveOps := []string{"delete", "remove", "force_delete", "bulk_delete", "purge", "destroy"}
	for _, op := range destructiveOps {
		if op == operation {
			return true
		}
	}
	return false
}

func isAdminOperation(operation string) bool {
	adminOps := []string{"ban", "suspend", "moderate", "override", "force", "bulk", "admin", "platform"}
	for _, op := range adminOps {
		if op == operation {
			return true
		}
	}
	return false
}

func getCommonReadOperations() []string {
	return []string{"search", "find", "filter", "get", "view", "list"}
}

func getCommonWriteOperations() []string {
	return []string{"add", "create", "update", "edit"}
}

func filterOutOperations(operations []string, toRemove []string) []string {
	filtered := []string{}
	for _, op := range operations {
		remove := false
		for _, removeOp := range toRemove {
			if op == removeOp {
				remove = true
				break
			}
		}
		if !remove {
			filtered = append(filtered, op)
		}
	}
	return filtered
}
