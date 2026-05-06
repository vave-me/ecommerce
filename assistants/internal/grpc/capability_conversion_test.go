package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"middleman/assistants/assistantspb"
	"middleman/assistants/internal/domain"
)

func TestDomainToProtoCapability(t *testing.T) {
	tests := []struct {
		name       string
		domainCap  domain.AssistantCapability
		expected   assistantspb.AssistantCapability
		shouldFail bool
	}{
		{
			name:      "assistant management",
			domainCap: domain.CapabilityAssistantManagement,
			expected:  assistantspb.AssistantCapability_ASSISTANT_MANAGEMENT,
		},
		{
			name:      "user interaction",
			domainCap: domain.CapabilityUserInteraction,
			expected:  assistantspb.AssistantCapability_USER_INTERACTION,
		},
		{
			name:      "data analysis",
			domainCap: domain.CapabilityDataAnalysis,
			expected:  assistantspb.AssistantCapability_DATA_ANALYSIS,
		},
		{
			name:      "location services",
			domainCap: domain.CapabilityLocationServices,
			expected:  assistantspb.AssistantCapability_LOCATION_SERVICES,
		},
		{
			name:      "authentication",
			domainCap: domain.CapabilityAuthentication,
			expected:  assistantspb.AssistantCapability_AUTHENTICATION,
		},
		{
			name:      "public api access",
			domainCap: domain.CapabilityPublicAPIAccess,
			expected:  assistantspb.AssistantCapability_PUBLIC_API_ACCESS,
		},
		{
			name:      "jailbreak resistant",
			domainCap: domain.CapabilityJailbreakResistant,
			expected:  assistantspb.AssistantCapability_JAILBREAK_RESISTANT,
		},
		{
			name:      "scope enforcement",
			domainCap: domain.CapabilityScopeEnforcement,
			expected:  assistantspb.AssistantCapability_SCOPE_ENFORCEMENT,
		},
		{
			name:      "data retrieval",
			domainCap: domain.CapabilityDataRetrieval,
			expected:  assistantspb.AssistantCapability_DATA_RETRIEVAL,
		},
		{
			name:      "search and filter",
			domainCap: domain.CapabilitySearchAndFilter,
			expected:  assistantspb.AssistantCapability_SEARCH_AND_FILTER,
		},
		{
			name:      "private api access",
			domainCap: domain.CapabilityPrivateAPIAccess,
			expected:  assistantspb.AssistantCapability_PRIVATE_API_ACCESS,
		},
		{
			name:      "user data access",
			domainCap: domain.CapabilityUserDataAccess,
			expected:  assistantspb.AssistantCapability_USER_DATA_ACCESS,
		},
		{
			name:      "token management",
			domainCap: domain.CapabilityTokenManagement,
			expected:  assistantspb.AssistantCapability_TOKEN_MANAGEMENT,
		},
		{
			name:      "data masking",
			domainCap: domain.CapabilityDataMasking,
			expected:  assistantspb.AssistantCapability_DATA_MASKING,
		},
		{
			name:      "audit logging",
			domainCap: domain.CapabilityAuditLogging,
			expected:  assistantspb.AssistantCapability_AUDIT_LOGGING,
		},
		{
			name:      "text generation",
			domainCap: domain.CapabilityTextGeneration,
			expected:  assistantspb.AssistantCapability_TEXT_GENERATION,
		},
		{
			name:      "code generation",
			domainCap: domain.CapabilityCodeGeneration,
			expected:  assistantspb.AssistantCapability_CODE_GENERATION,
		},
		{
			name:      "web search",
			domainCap: domain.CapabilityWebSearch,
			expected:  assistantspb.AssistantCapability_WEB_SEARCH,
		},
		{
			name:       "unknown capability",
			domainCap:  "unknown_capability",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := domainToProtoCapability(tt.domainCap)
			
			if tt.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestProtoToDomainCapability(t *testing.T) {
	tests := []struct {
		name       string
		protoCap   assistantspb.AssistantCapability
		expected   domain.AssistantCapability
		shouldFail bool
	}{
		{
			name:     "assistant management",
			protoCap: assistantspb.AssistantCapability_ASSISTANT_MANAGEMENT,
			expected: domain.CapabilityAssistantManagement,
		},
		{
			name:     "user interaction",
			protoCap: assistantspb.AssistantCapability_USER_INTERACTION,
			expected: domain.CapabilityUserInteraction,
		},
		{
			name:     "data analysis",
			protoCap: assistantspb.AssistantCapability_DATA_ANALYSIS,
			expected: domain.CapabilityDataAnalysis,
		},
		{
			name:     "location services",
			protoCap: assistantspb.AssistantCapability_LOCATION_SERVICES,
			expected: domain.CapabilityLocationServices,
		},
		{
			name:     "authentication",
			protoCap: assistantspb.AssistantCapability_AUTHENTICATION,
			expected: domain.CapabilityAuthentication,
		},
		{
			name:     "public api access",
			protoCap: assistantspb.AssistantCapability_PUBLIC_API_ACCESS,
			expected: domain.CapabilityPublicAPIAccess,
		},
		{
			name:     "jailbreak resistant",
			protoCap: assistantspb.AssistantCapability_JAILBREAK_RESISTANT,
			expected: domain.CapabilityJailbreakResistant,
		},
		{
			name:     "scope enforcement",
			protoCap: assistantspb.AssistantCapability_SCOPE_ENFORCEMENT,
			expected: domain.CapabilityScopeEnforcement,
		},
		{
			name:     "data retrieval",
			protoCap: assistantspb.AssistantCapability_DATA_RETRIEVAL,
			expected: domain.CapabilityDataRetrieval,
		},
		{
			name:     "search and filter",
			protoCap: assistantspb.AssistantCapability_SEARCH_AND_FILTER,
			expected: domain.CapabilitySearchAndFilter,
		},
		{
			name:     "private api access",
			protoCap: assistantspb.AssistantCapability_PRIVATE_API_ACCESS,
			expected: domain.CapabilityPrivateAPIAccess,
		},
		{
			name:     "user data access",
			protoCap: assistantspb.AssistantCapability_USER_DATA_ACCESS,
			expected: domain.CapabilityUserDataAccess,
		},
		{
			name:     "token management",
			protoCap: assistantspb.AssistantCapability_TOKEN_MANAGEMENT,
			expected: domain.CapabilityTokenManagement,
		},
		{
			name:     "data masking",
			protoCap: assistantspb.AssistantCapability_DATA_MASKING,
			expected: domain.CapabilityDataMasking,
		},
		{
			name:     "audit logging",
			protoCap: assistantspb.AssistantCapability_AUDIT_LOGGING,
			expected: domain.CapabilityAuditLogging,
		},
		{
			name:     "text generation",
			protoCap: assistantspb.AssistantCapability_TEXT_GENERATION,
			expected: domain.CapabilityTextGeneration,
		},
		{
			name:     "code generation",
			protoCap: assistantspb.AssistantCapability_CODE_GENERATION,
			expected: domain.CapabilityCodeGeneration,
		},
		{
			name:     "web search",
			protoCap: assistantspb.AssistantCapability_WEB_SEARCH,
			expected: domain.CapabilityWebSearch,
		},
		{
			name:       "unknown capability",
			protoCap:   assistantspb.AssistantCapability(999),
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := protoToDomainCapability(tt.protoCap)
			
			if tt.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDomainToProtoCapabilities_Deduplication(t *testing.T) {
	// Test with duplicates
	domainCaps := []domain.AssistantCapability{
		domain.CapabilityUserInteraction,
		domain.CapabilityDataAnalysis,
		domain.CapabilityUserInteraction, // Duplicate
		domain.CapabilityLocationServices,
		domain.CapabilityDataAnalysis, // Duplicate
		domain.CapabilityAuthentication,
		domain.CapabilityUserInteraction, // Duplicate
	}
	
	protoCaps := domainToProtoCapabilities(domainCaps)
	
	// Should have only 4 unique capabilities
	assert.Equal(t, 4, len(protoCaps))
	
	// Verify no duplicates
	seen := make(map[assistantspb.AssistantCapability]bool)
	for _, cap := range protoCaps {
		assert.False(t, seen[cap], "Found duplicate capability")
		seen[cap] = true
	}
}

func TestProtoToDomainCapabilities_Deduplication(t *testing.T) {
	// Test with duplicates
	protoCaps := []assistantspb.AssistantCapability{
		assistantspb.AssistantCapability_USER_INTERACTION,
		assistantspb.AssistantCapability_DATA_ANALYSIS,
		assistantspb.AssistantCapability_USER_INTERACTION, // Duplicate
		assistantspb.AssistantCapability_LOCATION_SERVICES,
		assistantspb.AssistantCapability_DATA_ANALYSIS, // Duplicate
		assistantspb.AssistantCapability_AUTHENTICATION,
		assistantspb.AssistantCapability_USER_INTERACTION, // Duplicate
	}
	
	domainCaps := protoToDomainCapabilities(protoCaps)
	
	// Should have only 4 unique capabilities
	assert.Equal(t, 4, len(domainCaps))
	
	// Verify no duplicates
	seen := make(map[domain.AssistantCapability]bool)
	for _, cap := range domainCaps {
		assert.False(t, seen[cap], "Found duplicate capability")
		seen[cap] = true
	}
}

func TestCapabilityConversion_RoundTrip(t *testing.T) {
	// Test that converting from domain to proto and back preserves all capabilities
	allDomainCaps := []domain.AssistantCapability{
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
	
	// Convert to proto
	protoCaps := domainToProtoCapabilities(allDomainCaps)
	assert.Equal(t, len(allDomainCaps), len(protoCaps))
	
	// Convert back to domain
	resultDomainCaps := protoToDomainCapabilities(protoCaps)
	assert.Equal(t, len(allDomainCaps), len(resultDomainCaps))
	
	// Verify all capabilities are preserved
	expectedCaps := make(map[domain.AssistantCapability]bool)
	for _, cap := range allDomainCaps {
		expectedCaps[cap] = true
	}
	
	for _, cap := range resultDomainCaps {
		assert.True(t, expectedCaps[cap], "Capability %s was not preserved in round trip", cap)
		delete(expectedCaps, cap)
	}
	
	// All capabilities should have been found
	assert.Empty(t, expectedCaps)
}