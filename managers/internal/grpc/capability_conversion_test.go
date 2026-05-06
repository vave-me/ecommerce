package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"middleman/managers/internal/domain"
	"middleman/managers/managerspb"
)

func TestDomainToProtoCapability(t *testing.T) {
	tests := []struct {
		name       string
		domainCap  domain.ManagerCapability
		expected   managerspb.ManagerCapability
		shouldFail bool
	}{
		{
			name:      "manager management",
			domainCap: domain.CapabilityManagerManagement,
			expected:  managerspb.ManagerCapability_ASSISTANT_MANAGEMENT,
		},
		{
			name:      "user interaction",
			domainCap: domain.CapabilityUserInteraction,
			expected:  managerspb.ManagerCapability_USER_INTERACTION,
		},
		{
			name:      "data analysis",
			domainCap: domain.CapabilityDataAnalysis,
			expected:  managerspb.ManagerCapability_DATA_ANALYSIS,
		},
		{
			name:      "location services",
			domainCap: domain.CapabilityLocationServices,
			expected:  managerspb.ManagerCapability_LOCATION_SERVICES,
		},
		{
			name:      "authentication",
			domainCap: domain.CapabilityAuthentication,
			expected:  managerspb.ManagerCapability_AUTHENTICATION,
		},
		{
			name:      "public api access",
			domainCap: domain.CapabilityPublicAPIAccess,
			expected:  managerspb.ManagerCapability_PUBLIC_API_ACCESS,
		},
		{
			name:      "jailbreak resistant",
			domainCap: domain.CapabilityJailbreakResistant,
			expected:  managerspb.ManagerCapability_JAILBREAK_RESISTANT,
		},
		{
			name:      "scope enforcement",
			domainCap: domain.CapabilityScopeEnforcement,
			expected:  managerspb.ManagerCapability_SCOPE_ENFORCEMENT,
		},
		{
			name:      "data retrieval",
			domainCap: domain.CapabilityDataRetrieval,
			expected:  managerspb.ManagerCapability_DATA_RETRIEVAL,
		},
		{
			name:      "search and filter",
			domainCap: domain.CapabilitySearchAndFilter,
			expected:  managerspb.ManagerCapability_SEARCH_AND_FILTER,
		},
		{
			name:      "private api access",
			domainCap: domain.CapabilityPrivateAPIAccess,
			expected:  managerspb.ManagerCapability_PRIVATE_API_ACCESS,
		},
		{
			name:      "user data access",
			domainCap: domain.CapabilityUserDataAccess,
			expected:  managerspb.ManagerCapability_USER_DATA_ACCESS,
		},
		{
			name:      "token management",
			domainCap: domain.CapabilityTokenManagement,
			expected:  managerspb.ManagerCapability_TOKEN_MANAGEMENT,
		},
		{
			name:      "data masking",
			domainCap: domain.CapabilityDataMasking,
			expected:  managerspb.ManagerCapability_DATA_MASKING,
		},
		{
			name:      "audit logging",
			domainCap: domain.CapabilityAuditLogging,
			expected:  managerspb.ManagerCapability_AUDIT_LOGGING,
		},
		{
			name:      "text generation",
			domainCap: domain.CapabilityTextGeneration,
			expected:  managerspb.ManagerCapability_TEXT_GENERATION,
		},
		{
			name:      "code generation",
			domainCap: domain.CapabilityCodeGeneration,
			expected:  managerspb.ManagerCapability_CODE_GENERATION,
		},
		{
			name:      "web search",
			domainCap: domain.CapabilityWebSearch,
			expected:  managerspb.ManagerCapability_WEB_SEARCH,
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
		protoCap   managerspb.ManagerCapability
		expected   domain.ManagerCapability
		shouldFail bool
	}{
		{
			name:     "manager management",
			protoCap: managerspb.ManagerCapability_ASSISTANT_MANAGEMENT,
			expected: domain.CapabilityManagerManagement,
		},
		{
			name:     "user interaction",
			protoCap: managerspb.ManagerCapability_USER_INTERACTION,
			expected: domain.CapabilityUserInteraction,
		},
		{
			name:     "data analysis",
			protoCap: managerspb.ManagerCapability_DATA_ANALYSIS,
			expected: domain.CapabilityDataAnalysis,
		},
		{
			name:     "location services",
			protoCap: managerspb.ManagerCapability_LOCATION_SERVICES,
			expected: domain.CapabilityLocationServices,
		},
		{
			name:     "authentication",
			protoCap: managerspb.ManagerCapability_AUTHENTICATION,
			expected: domain.CapabilityAuthentication,
		},
		{
			name:     "public api access",
			protoCap: managerspb.ManagerCapability_PUBLIC_API_ACCESS,
			expected: domain.CapabilityPublicAPIAccess,
		},
		{
			name:     "jailbreak resistant",
			protoCap: managerspb.ManagerCapability_JAILBREAK_RESISTANT,
			expected: domain.CapabilityJailbreakResistant,
		},
		{
			name:     "scope enforcement",
			protoCap: managerspb.ManagerCapability_SCOPE_ENFORCEMENT,
			expected: domain.CapabilityScopeEnforcement,
		},
		{
			name:     "data retrieval",
			protoCap: managerspb.ManagerCapability_DATA_RETRIEVAL,
			expected: domain.CapabilityDataRetrieval,
		},
		{
			name:     "search and filter",
			protoCap: managerspb.ManagerCapability_SEARCH_AND_FILTER,
			expected: domain.CapabilitySearchAndFilter,
		},
		{
			name:     "private api access",
			protoCap: managerspb.ManagerCapability_PRIVATE_API_ACCESS,
			expected: domain.CapabilityPrivateAPIAccess,
		},
		{
			name:     "user data access",
			protoCap: managerspb.ManagerCapability_USER_DATA_ACCESS,
			expected: domain.CapabilityUserDataAccess,
		},
		{
			name:     "token management",
			protoCap: managerspb.ManagerCapability_TOKEN_MANAGEMENT,
			expected: domain.CapabilityTokenManagement,
		},
		{
			name:     "data masking",
			protoCap: managerspb.ManagerCapability_DATA_MASKING,
			expected: domain.CapabilityDataMasking,
		},
		{
			name:     "audit logging",
			protoCap: managerspb.ManagerCapability_AUDIT_LOGGING,
			expected: domain.CapabilityAuditLogging,
		},
		{
			name:     "text generation",
			protoCap: managerspb.ManagerCapability_TEXT_GENERATION,
			expected: domain.CapabilityTextGeneration,
		},
		{
			name:     "code generation",
			protoCap: managerspb.ManagerCapability_CODE_GENERATION,
			expected: domain.CapabilityCodeGeneration,
		},
		{
			name:     "web search",
			protoCap: managerspb.ManagerCapability_WEB_SEARCH,
			expected: domain.CapabilityWebSearch,
		},
		{
			name:       "unknown capability",
			protoCap:   managerspb.ManagerCapability(999),
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
	domainCaps := []domain.ManagerCapability{
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
	seen := make(map[managerspb.ManagerCapability]bool)
	for _, cap := range protoCaps {
		assert.False(t, seen[cap], "Found duplicate capability")
		seen[cap] = true
	}
}

func TestProtoToDomainCapabilities_Deduplication(t *testing.T) {
	// Test with duplicates
	protoCaps := []managerspb.ManagerCapability{
		managerspb.ManagerCapability_USER_INTERACTION,
		managerspb.ManagerCapability_DATA_ANALYSIS,
		managerspb.ManagerCapability_USER_INTERACTION, // Duplicate
		managerspb.ManagerCapability_LOCATION_SERVICES,
		managerspb.ManagerCapability_DATA_ANALYSIS, // Duplicate
		managerspb.ManagerCapability_AUTHENTICATION,
		managerspb.ManagerCapability_USER_INTERACTION, // Duplicate
	}

	domainCaps := protoToDomainCapabilities(protoCaps)

	// Should have only 4 unique capabilities
	assert.Equal(t, 4, len(domainCaps))

	// Verify no duplicates
	seen := make(map[domain.ManagerCapability]bool)
	for _, cap := range domainCaps {
		assert.False(t, seen[cap], "Found duplicate capability")
		seen[cap] = true
	}
}

func TestCapabilityConversion_RoundTrip(t *testing.T) {
	// Test that converting from domain to proto and back preserves all capabilities
	allDomainCaps := []domain.ManagerCapability{
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

	// Convert to proto
	protoCaps := domainToProtoCapabilities(allDomainCaps)
	assert.Equal(t, len(allDomainCaps), len(protoCaps))

	// Convert back to domain
	resultDomainCaps := protoToDomainCapabilities(protoCaps)
	assert.Equal(t, len(allDomainCaps), len(resultDomainCaps))

	// Verify all capabilities are preserved
	expectedCaps := make(map[domain.ManagerCapability]bool)
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
