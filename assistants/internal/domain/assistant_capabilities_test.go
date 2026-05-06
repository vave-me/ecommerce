package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"middleman/assistants/internal/domain"
)

func TestAssistantCapabilities_Deduplication(t *testing.T) {
	tests := []struct {
		name         string
		capabilities []domain.AssistantCapability
		expected     int
	}{
		{
			name: "no duplicates",
			capabilities: []domain.AssistantCapability{
				domain.CapabilityUserInteraction,
				domain.CapabilityDataAnalysis,
				domain.CapabilityLocationServices,
			},
			expected: 3,
		},
		{
			name: "with duplicates",
			capabilities: []domain.AssistantCapability{
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
			capabilities: []domain.AssistantCapability{},
			expected:     0,
		},
		{
			name: "all duplicates",
			capabilities: []domain.AssistantCapability{
				domain.CapabilityUserInteraction,
				domain.CapabilityUserInteraction,
				domain.CapabilityUserInteraction,
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create assistant
			assistant := domain.NewAssistant("test-id")
			
			// Create assistant with capabilities
			event, err := assistant.CreateAssistant(
				"test-id",
				"Test Assistant",
				"Test Description",
				"user-123",
				domain.AssistantTypeStandard,
				tt.capabilities,
				0.7,
				1000,
				"Test prompt",
			)
			
			assert.NoError(t, err)
			assert.NotNil(t, event)
			
			// Check that capabilities were deduplicated
			assert.Equal(t, tt.expected, len(assistant.Capabilities))
			
			// Verify no duplicates
			seen := make(map[domain.AssistantCapability]bool)
			for _, cap := range assistant.Capabilities {
				assert.False(t, seen[cap], "Found duplicate capability: %s", cap)
				seen[cap] = true
			}
		})
	}
}

func TestAssistantCapabilities_UpdateWithDeduplication(t *testing.T) {
	// Create assistant
	assistant := domain.NewAssistant("test-id")
	
	// Create assistant with initial capabilities
	_, err := assistant.CreateAssistant(
		"test-id",
		"Test Assistant",
		"Test Description",
		"user-123",
		domain.AssistantTypeStandard,
		[]domain.AssistantCapability{
			domain.CapabilityUserInteraction,
			domain.CapabilityDataAnalysis,
		},
		0.7,
		1000,
		"Test prompt",
	)
	assert.NoError(t, err)
	
	// Activate the assistant
	_, err = assistant.Activate()
	assert.NoError(t, err)
	
	// Update with duplicate capabilities
	event, err := assistant.UpdateConfigurationWithCapabilities(
		0.8,
		2000,
		"Updated prompt",
		[]domain.AssistantCapability{
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
	assert.Equal(t, 4, len(assistant.Capabilities))
	
	// Verify the correct capabilities are present
	expectedCaps := map[domain.AssistantCapability]bool{
		domain.CapabilityUserInteraction:  true,
		domain.CapabilityDataAnalysis:     true,
		domain.CapabilityLocationServices: true,
		domain.CapabilityAuthentication:   true,
	}
	
	for _, cap := range assistant.Capabilities {
		assert.True(t, expectedCaps[cap], "Unexpected capability: %s", cap)
	}
}

func TestAssistantCapabilities_AllValidCapabilities(t *testing.T) {
	// Test that all defined capabilities are valid
	allCapabilities := []domain.AssistantCapability{
		domain.CapabilityAssistantManagement,
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
	
	// Create assistant with all capabilities
	assistant := domain.NewAssistant("test-id")
	event, err := assistant.CreateAssistant(
		"test-id",
		"Test Assistant",
		"Test Description",
		"user-123",
		domain.AssistantTypeAdmin, // Admin can have all capabilities
		allCapabilities,
		0.7,
		1000,
		"Test prompt",
	)
	
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, len(allCapabilities), len(assistant.Capabilities))
}

func TestAssistantCapabilities_InvalidCapabilityFiltering(t *testing.T) {
	// Create assistant with some invalid capabilities (simulated by empty string)
	assistant := domain.NewAssistant("test-id")
	
	capabilities := []domain.AssistantCapability{
		domain.CapabilityUserInteraction,
		"", // Invalid capability
		domain.CapabilityDataAnalysis,
		"invalid_capability", // Invalid capability
		domain.CapabilityLocationServices,
	}
	
	event, err := assistant.CreateAssistant(
		"test-id",
		"Test Assistant",
		"Test Description",
		"user-123",
		domain.AssistantTypeStandard,
		capabilities,
		0.7,
		1000,
		"Test prompt",
	)
	
	assert.NoError(t, err)
	assert.NotNil(t, event)
	
	// Should only have the 3 valid capabilities
	assert.Equal(t, 3, len(assistant.Capabilities))
	
	// Verify only valid capabilities are present
	validCaps := map[domain.AssistantCapability]bool{
		domain.CapabilityUserInteraction:  true,
		domain.CapabilityDataAnalysis:     true,
		domain.CapabilityLocationServices: true,
	}
	
	for _, cap := range assistant.Capabilities {
		assert.True(t, validCaps[cap], "Unexpected capability: %s", cap)
	}
}

func TestAssistantCapabilities_HasCapability(t *testing.T) {
	// Create assistant
	assistant := domain.NewAssistant("test-id")
	
	// Create assistant with specific capabilities
	_, err := assistant.CreateAssistant(
		"test-id",
		"Test Assistant",
		"Test Description",
		"user-123",
		domain.AssistantTypeStandard,
		[]domain.AssistantCapability{
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
	assert.True(t, assistant.HasCapability(domain.CapabilityUserInteraction))
	assert.True(t, assistant.HasCapability(domain.CapabilityDataAnalysis))
	assert.True(t, assistant.HasCapability(domain.CapabilityLocationServices))
	
	// Test capabilities not present
	assert.False(t, assistant.HasCapability(domain.CapabilityAuthentication))
	assert.False(t, assistant.HasCapability(domain.CapabilityPrivateAPIAccess))
	assert.False(t, assistant.HasCapability(domain.CapabilityWebSearch))
}