package processor

import (
	"middleman/managers/internal/application/tools"
	"testing"
)

func TestBuildIterativeMessage_WithContextualHelp(t *testing.T) {
	p := &LLMProcessor{}
	
	// Test with schema_generate_contextual_help result
	toolResults := []*tools.ToolOperationResult{
		{
			EntityType: "schema",
			Operation:  "schema_generate_contextual_help",
			Success:    true,
			Result:     "I can help you with comprehensive marketplace operations including product management, order processing, and user management.",
		},
	}
	
	message := p.buildIterativeMessage("Help me with the marketplace", toolResults, 1)
	
	// Verify the contextual help is included
	expectedSubstring := "I can help you with comprehensive marketplace operations"
	if !contains(message, expectedSubstring) {
		t.Errorf("Expected message to contain '%s', but got: %s", expectedSubstring, message)
	}
}

func TestBuildSchemaContextFromResults_WithContextualHelp(t *testing.T) {
	p := &LLMProcessor{}
	
	// Test with schema_generate_contextual_help result
	toolResults := []*tools.ToolOperationResult{
		{
			EntityType: "schema",
			Operation:  "schema_generate_contextual_help",
			Success:    true,
			Result:     "Here are the available operations for products: search, create, update, delete.",
		},
	}
	
	context := p.buildSchemaContextFromResults(toolResults)
	
	// Verify the contextual help is extracted
	expectedSubstring := "Here are the available operations for products"
	if !contains(context, expectedSubstring) {
		t.Errorf("Expected context to contain '%s', but got: %s", expectedSubstring, context)
	}
}

func TestBuildIterativeMessage_WithMixedSchemaResults(t *testing.T) {
	p := &LLMProcessor{}
	
	// Test with mixed schema results
	toolResults := []*tools.ToolOperationResult{
		{
			EntityType: "schema",
			Operation:  "schema_generate_contextual_help",
			Success:    true,
			Result:     "Product management operations available.",
		},
		{
			EntityType: "schema",
			Operation:  "schema_get_fields",
			Success:    true,
			Result: map[string]interface{}{
				"fields": []string{"id", "name", "price", "description"},
			},
		},
	}
	
	message := p.buildIterativeMessage("Show me products", toolResults, 1)
	
	// Verify both types of results are included
	if !contains(message, "Product management operations available") {
		t.Error("Expected message to contain contextual help")
	}
	if !contains(message, "Fields available:") {
		t.Error("Expected message to contain fields information")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && len(substr) == 0 || (len(substr) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}