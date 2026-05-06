package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"middleman/managers/internal/domain"
)

func TestManagerCapabilities_Deduplication(t *testing.T) {
	tests := []struct {
		name         string
		capabilities []domain.ManagerCapability
		expected     int
	}{
		{
			name: "no duplicates",
			capabilities: []domain.ManagerCapability{
				domain.CapabilityUserInteraction,
				domain.CapabilityDataAnalysis,
				domain.CapabilityLocationServices,
			},
			expected: 3,
		},
		{
			name: "with duplicates",
			capabilities: []domain.ManagerCapability{
				domain.CapabilityUserInteraction,
				domain.CapabilityDataAnalysis,
				domain.CapabilityUserInteraction,
				domain.CapabilityDataAnalysis,
				domain.CapabilityLocationServices,
				domain.CapabilityUserInteraction,
			},
			expected: 3,
		},
		{
			name:         "empty capabilities",
			capabilities: []domain.ManagerCapability{},
			expected:     0,
		},
		{
			name: "all duplicates",
			capabilities: []domain.ManagerCapability{
				domain.CapabilityUserInteraction,
				domain.CapabilityUserInteraction,
				domain.CapabilityUserInteraction,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create manager
			manager := domain.NewManager("test-id")

			// Create manager with capabilities
			event, err := manager.CreateManager(
				"test-id",
				"Test Manager",
				"Test Description",
				"user-123",
				domain.ManagerTypeStandard,
				tt.capabilities,
				0.7,
				1000,
				"Test prompt",
			)

			assert.NoError(t, err)
			assert.NotNil(t, event)

			// Check that capabilities were deduplicated
			assert.Equal(t, tt.expected, len(manager.Capabilities))

			// Verify no duplicates
			seen := make(map[domain.ManagerCapability]bool)
			for _, cap := range manager.Capabilities {
				assert.False(t, seen[cap], "Found duplicate capability: %s", cap)
				seen[cap] = true
			}
		})
	}
}

func TestManagerCapabilities_UpdateWithDeduplication(t *testing.T) {
	// Create manager
	manager := domain.NewManager("test-id")

	// Create manager with initial capabilities
	_, err := manager.CreateManager(
		"test-id",
		"Test Manager",
		"Test Description",
		"user-123",
		domain.ManagerTypeStandard,
		[]domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
		},
		0.7,
		1000,
		"Test prompt",
	)
	assert.NoError(t, err)

	// Activate the manager
	_, err = manager.Activate()
	assert.NoError(t, err)

	// Update with duplicate capabilities
	event, err := manager.UpdateConfigurationWithCapabilities(
		0.8,
		2000,
		"Updated prompt",
		[]domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilityLocationServices,
			domain.CapabilityUserInteraction, // Duplicate
			domain.CapabilityDataAnalysis,    // Duplicate
			domain.CapabilityAuthentication,
			domain.CapabilityLocationServices, // Duplicate
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, event)

	// Check that capabilities were deduplicated
	assert.Equal(t, 4, len(manager.Capabilities))

	// Verify the correct capabilities are present
	expectedCaps := map[domain.ManagerCapability]bool{
		domain.CapabilityUserInteraction:  true,
		domain.CapabilityDataAnalysis:     true,
		domain.CapabilityLocationServices: true,
		domain.CapabilityAuthentication:   true,
	}

	for _, cap := range manager.Capabilities {
		assert.True(t, expectedCaps[cap], "Unexpected capability: %s", cap)
	}
}

func TestManagerCapabilities_AllValidCapabilities(t *testing.T) {
	// Test that all defined capabilities are valid
	allCapabilities := []domain.ManagerCapability{
		domain.CapabilityManagerManagement,
		domain.CapabilityUserInteraction,
		domain.CapabilityDataAnalysis,
		domain.CapabilityLocationServices,
		domain.CapabilityAuthentication,
		domain.CapabilityPublicAPIAccess,
		domain.CapabilityJailbreakResistant,
		domain.CapabilityScopeEnforcement,
		domain.CapabilityDataRetrieval,
		domain.CapabilitySearchAndFilter,
		domain.CapabilityPrivateAPIAccess,
		domain.CapabilityUserDataAccess,
		domain.CapabilityTokenManagement,
		domain.CapabilityDataMasking,
		domain.CapabilityAuditLogging,
		domain.CapabilityTextGeneration,
		domain.CapabilityCodeGeneration,
		domain.CapabilityWebSearch,
	}

	// Create manager with all capabilities
	manager := domain.NewManager("test-id")
	event, err := manager.CreateManager(
		"test-id",
		"Test Manager",
		"Test Description",
		"user-123",
		domain.ManagerTypeAdmin, // Admin can have all capabilities
		allCapabilities,
		0.7,
		1000,
		"Test prompt",
	)

	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, len(allCapabilities), len(manager.Capabilities))
}

func TestManagerCapabilities_InvalidCapabilityFiltering(t *testing.T) {
	// Create manager with some invalid capabilities (simulated by empty string)
	manager := domain.NewManager("test-id")

	capabilities := []domain.ManagerCapability{
		domain.CapabilityUserInteraction,
		"", // Invalid capability
		domain.CapabilityDataAnalysis,
		"invalid_capability", // Invalid capability
		domain.CapabilityLocationServices,
	}

	event, err := manager.CreateManager(
		"test-id",
		"Test Manager",
		"Test Description",
		"user-123",
		domain.ManagerTypeStandard,
		capabilities,
		0.7,
		1000,
		"Test prompt",
	)

	assert.NoError(t, err)
	assert.NotNil(t, event)

	// Should only have the 3 valid capabilities
	assert.Equal(t, 3, len(manager.Capabilities))

	// Verify only valid capabilities are present
	validCaps := map[domain.ManagerCapability]bool{
		domain.CapabilityUserInteraction:  true,
		domain.CapabilityDataAnalysis:     true,
		domain.CapabilityLocationServices: true,
	}

	for _, cap := range manager.Capabilities {
		assert.True(t, validCaps[cap], "Unexpected capability: %s", cap)
	}
}

func TestManagerCapabilities_HasCapability(t *testing.T) {
	// Create manager
	manager := domain.NewManager("test-id")

	// Create manager with specific capabilities
	_, err := manager.CreateManager(
		"test-id",
		"Test Manager",
		"Test Description",
		"user-123",
		domain.ManagerTypeStandard,
		[]domain.ManagerCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
			domain.CapabilityLocationServices,
		},
		0.7,
		1000,
		"Test prompt",
	)
	assert.NoError(t, err)

	// Test HasCapability
	assert.True(t, manager.HasCapability(domain.CapabilityUserInteraction))
	assert.True(t, manager.HasCapability(domain.CapabilityDataAnalysis))
	assert.True(t, manager.HasCapability(domain.CapabilityLocationServices))

	// Test capabilities not present
	assert.False(t, manager.HasCapability(domain.CapabilityAuthentication))
	assert.False(t, manager.HasCapability(domain.CapabilityPrivateAPIAccess))
	assert.False(t, manager.HasCapability(domain.CapabilityWebSearch))
}
