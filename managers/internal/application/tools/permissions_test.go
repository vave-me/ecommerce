package tools_test

import (
	"middleman/managers/internal/application/tools"
	"middleman/managers/internal/domain"
	"testing"
)

func TestPermissionChecker_CanPerformOperation(t *testing.T) {
	checker := tools.NewPermissionChecker()

	tests := []struct {
		name            string
		managerType     domain.ManagerType
		entityType      string
		operation       string
		userID          string
		resourceOwnerID string
		expected        bool
		description     string
	}{
		// Admin Manager Tests
		{
			"Admin can ban users",
			domain.ManagerTypeAdmin,
			"users",
			"ban",
			"admin-123",
			"",
			true,
			"Admin managers should be able to ban users",
		},
		{
			"Admin can force delete products",
			domain.ManagerTypeAdmin,
			"products",
			"force_remove",
			"admin-123",
			"vendor-456",
			true,
			"Admin managers should be able to force remove any product",
		},
		{
			"Admin can view all metrics",
			domain.ManagerTypeAdmin,
			"metrics",
			"view_platform_metrics",
			"admin-123",
			"",
			true,
			"Admin managers should be able to view platform metrics",
		},

		// Business Manager Tests
		{
			"Business can update own products",
			domain.ManagerTypeBusiness,
			"products",
			"update",
			"vendor-123",
			"vendor-123",
			true,
			"Business managers should be able to update their own products",
		},
		{
			"Business cannot update others' products",
			domain.ManagerTypeBusiness,
			"products",
			"update",
			"vendor-123",
			"vendor-456",
			false,
			"Business managers should not be able to update others' products",
		},
		{
			"Business can search products",
			domain.ManagerTypeBusiness,
			"products",
			"search",
			"vendor-123",
			"",
			true,
			"Business managers should be able to search products",
		},
		{
			"Business cannot access payments",
			domain.ManagerTypeBusiness,
			"payments",
			"view",
			"vendor-123",
			"",
			false,
			"Business managers should not access payment details",
		},

		// Support Manager Tests
		{
			"Support can search products",
			domain.ManagerTypeSupport,
			"products",
			"search",
			"support-123",
			"",
			true,
			"Support managers should be able to search products",
		},
		{
			"Support cannot delete products",
			domain.ManagerTypeSupport,
			"products",
			"delete",
			"support-123",
			"",
			false,
			"Support managers should not be able to delete products",
		},
		{
			"Support can create tickets",
			domain.ManagerTypeSupport,
			"support",
			"create_ticket",
			"support-123",
			"",
			true,
			"Support managers should be able to create support tickets",
		},
		{
			"Support cannot access payments",
			domain.ManagerTypeSupport,
			"payments",
			"view",
			"support-123",
			"",
			false,
			"Support managers should not access payments",
		},

		// Scheduler Manager Tests
		{
			"Scheduler can create orders",
			domain.ManagerTypeScheduler,
			"orders",
			"create",
			"scheduler",
			"",
			true,
			"Scheduler managers should be able to create orders",
		},
		{
			"Scheduler cannot delete users",
			domain.ManagerTypeScheduler,
			"users",
			"delete",
			"scheduler",
			"",
			false,
			"Scheduler managers should not perform destructive operations",
		},

		// Standard Manager Tests
		{
			"Standard can search products",
			domain.ManagerTypeStandard,
			"products",
			"search",
			"user-123",
			"",
			true,
			"Standard managers should be able to search products",
		},
		{
			"Standard cannot ban users",
			domain.ManagerTypeStandard,
			"users",
			"ban",
			"user-123",
			"",
			false,
			"Standard managers should not be able to ban users",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := checker.CanPerformOperation(
				test.managerType,
				test.entityType,
				test.operation,
				test.userID,
				test.resourceOwnerID,
			)

			if result != test.expected {
				if test.expected {
					t.Errorf("%s: expected allowed, got denied. Error: %v", test.description, err)
				} else {
					t.Errorf("%s: expected denied, got allowed", test.description)
				}
			}
		})
	}
}

func TestPermissionChecker_GetAllowedOperations(t *testing.T) {
	checker := tools.NewPermissionChecker()

	tests := []struct {
		name             string
		managerType      domain.ManagerType
		entityType       string
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			"Admin gets all operations",
			domain.ManagerTypeAdmin,
			"users",
			[]string{"*"},
			[]string{},
		},
		{
			"Business gets limited product operations",
			domain.ManagerTypeBusiness,
			"products",
			[]string{"add", "update", "search", "find"},
			[]string{"force_remove", "bulk_delete"},
		},
		{
			"Support gets read-only operations",
			domain.ManagerTypeSupport,
			"products",
			[]string{"search", "find", "filter"},
			[]string{"add", "update", "delete"},
		},
		{
			"Support cannot access payments",
			domain.ManagerTypeSupport,
			"payments",
			[]string{},
			[]string{"view", "refund"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			operations := checker.GetAllowedOperations(test.managerType, test.entityType)

			// Check operations that should be present
			for _, op := range test.shouldContain {
				found := false
				for _, allowedOp := range operations {
					if allowedOp == op {
						found = true
						break
					}
				}
				if !found && op != "*" { // Skip wildcard check
					t.Errorf("Expected operation '%s' to be allowed for %s on %s",
						op, test.managerType, test.entityType)
				}
			}

			// Check operations that should NOT be present
			for _, op := range test.shouldNotContain {
				for _, allowedOp := range operations {
					if allowedOp == op {
						t.Errorf("Operation '%s' should not be allowed for %s on %s",
							op, test.managerType, test.entityType)
					}
				}
			}
		})
	}
}
