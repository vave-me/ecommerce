package processor

import (
	"context"
	"fmt"
	"log"
	"middleman/managers/internal/application/services"
	"strings"

	"middleman/managers/internal/constants"
	"middleman/managers/internal/models"
)

// SchemaConsultationService provides centralized schema consultation for LLM
type SchemaConsultationService struct {
	enhancedInterface *EnhancedLLMInterface
	logger            *log.Logger
}

// NewSchemaConsultationService creates a new schema consultation service
func NewSchemaConsultationService() *SchemaConsultationService {
	return &SchemaConsultationService{
		enhancedInterface: NewEnhancedLLMInterface(),
		logger:            log.New(log.Writer(), "[SchemaConsultation] ", log.LstdFlags),
	}
}

// ConsultationRequest represents a request for schema consultation
type ConsultationRequest struct {
	UserInput       string                 `json:"user_input"`
	Context         map[string]interface{} `json:"context"`
	UncertaintyType string                 `json:"uncertainty_type"`
	PreviousError   string                 `json:"previous_error,omitempty"`
	Confidence      float64                `json:"confidence,omitempty"`
	EntityType      string                 `json:"entity_type,omitempty"`
	Operation       string                 `json:"operation,omitempty"`
	Parameters      map[string]interface{} `json:"parameters,omitempty"`
}

// ConsultationResponse provides guidance and suggestions
type ConsultationResponse struct {
	ShouldConsult     bool                           `json:"should_consult"`
	Guidance          string                         `json:"guidance"`
	Suggestions       []OperationSuggestion          `json:"suggestions,omitempty"`
	ValidFields       []services.FieldDefinition     `json:"valid_fields,omitempty"`
	ValidOperations   []services.OperationDefinition `json:"valid_operations,omitempty"`
	RecommendedAction string                         `json:"recommended_action"`
	Confidence        float64                        `json:"confidence"`
	NextSteps         []string                       `json:"next_steps"`
}

// ConsultSchema is the main entry point for LLM schema consultation
func (scs *SchemaConsultationService) ConsultSchema(ctx context.Context, request ConsultationRequest) (*ConsultationResponse, error) {
	scs.logger.Printf("Schema consultation requested: type=%s, input='%s'", request.UncertaintyType, request.UserInput)

	response := &ConsultationResponse{
		ShouldConsult: true,
		Confidence:    0.5,
	}

	// Determine if consultation is needed
	shouldConsult := scs.enhancedInterface.ShouldConsultSchema(request.UserInput, request.Context)
	response.ShouldConsult = shouldConsult

	if !shouldConsult && request.UncertaintyType == "" {
		response.ShouldConsult = false
		response.Guidance = "No schema consultation needed. Proceed with your current understanding."
		response.Confidence = 0.9
		return response, nil
	}

	// Provide specific guidance based on uncertainty type
	switch request.UncertaintyType {
	case "ambiguous_entity":
		return scs.handleAmbiguousEntity(request)
	case "unclear_operation":
		return scs.handleUnclearOperation(request)
	case "missing_parameters":
		return scs.handleMissingParameters(request)
	case "validation_failed":
		return scs.handleValidationFailure(request)
	case "operation_failed":
		return scs.handleOperationFailure(request)
	case "general_uncertainty":
		return scs.handleGeneralUncertainty(request)
	default:
		return scs.handleAutoDetectedUncertainty(request)
	}
}

// handleAmbiguousEntity helps when multiple entities could match
func (scs *SchemaConsultationService) handleAmbiguousEntity(request ConsultationRequest) (*ConsultationResponse, error) {
	scs.logger.Printf("Handling ambiguous entity for input: %s", request.UserInput)

	guidance := scs.enhancedInterface.GetSchemaBasedGuidance("ambiguous_entity", map[string]interface{}{
		"user_input": request.UserInput,
	})

	// Get contextual help to suggest entities
	contextualHelp := scs.enhancedInterface.GenerateContextualHelp(context.Background(), request.UserInput)

	response := &ConsultationResponse{
		ShouldConsult:     true,
		Guidance:          guidance,
		RecommendedAction: "Use GenerateContextualHelp() to identify the most relevant entity type",
		Confidence:        0.7,
		NextSteps: []string{
			"Review entity descriptions and select the most relevant one",
			"Use natural language patterns to guide entity selection",
			"Consider the context of previous operations",
			"Call GenerateContextualHelp() with specific entity mentions",
		},
	}

	// Extract entity suggestions from contextual help
	response.Guidance = fmt.Sprintf("%s\n\n=== CONTEXTUAL HELP ===\n%s", guidance, contextualHelp)

	return response, nil
}

// handleUnclearOperation helps when operation intent is unclear
func (scs *SchemaConsultationService) handleUnclearOperation(request ConsultationRequest) (*ConsultationResponse, error) {
	scs.logger.Printf("Handling unclear operation for entity: %s, input: %s", request.EntityType, request.UserInput)

	entityType := models.ToEntityType(request.EntityType)
	guidance := scs.enhancedInterface.GetSchemaBasedGuidance("unclear_operation", map[string]interface{}{
		"entity_type": entityType,
		"user_input":  request.UserInput,
	})

	// Get operation suggestions
	suggestions := scs.enhancedInterface.SuggestOperations(request.UserInput)

	// Get valid operations for the entity
	var validOperations []services.OperationDefinition
	if request.EntityType != "" {
		validOperations = scs.enhancedInterface.SchemaRegistry.GetOperationsByEntityType(entityType)
	}

	response := &ConsultationResponse{
		ShouldConsult:     true,
		Guidance:          guidance,
		Suggestions:       suggestions,
		ValidOperations:   validOperations,
		RecommendedAction: "Use SuggestOperations() to get ranked operation suggestions",
		Confidence:        0.8,
		NextSteps: []string{
			"Review natural language patterns for operations",
			"Consider the user's apparent intent (search, create, update, etc.)",
			"Check operation descriptions for best match",
			"Validate selected operation supports required parameters",
		},
	}

	return response, nil
}

// handleMissingParameters helps when parameters are missing or unclear
func (scs *SchemaConsultationService) handleMissingParameters(request ConsultationRequest) (*ConsultationResponse, error) {
	scs.logger.Printf("Handling missing parameters for %s.%s", request.EntityType, request.Operation)

	entityType := models.ToEntityType(request.EntityType)
	guidance := scs.enhancedInterface.GetSchemaBasedGuidance("missing_parameters", map[string]interface{}{
		"entity_type": entityType,
		"operation":   request.Operation,
	})

	// Get valid fields for the entity
	var validFields []services.FieldDefinition
	if request.EntityType != "" {
		validFields = scs.enhancedInterface.SchemaRegistry.GetFieldsByEntityType(entityType)
	}

	response := &ConsultationResponse{
		ShouldConsult:     true,
		Guidance:          guidance,
		ValidFields:       validFields,
		RecommendedAction: "Use GetFieldsByEntityType() to see available fields and their requirements",
		Confidence:        0.9,
		NextSteps: []string{
			"Check required vs optional parameters for the operation",
			"Review field types and constraints",
			"Validate parameter names match schema field names",
			"Use default values where appropriate",
		},
	}

	return response, nil
}

// handleValidationFailure helps recover from validation errors
func (scs *SchemaConsultationService) handleValidationFailure(request ConsultationRequest) (*ConsultationResponse, error) {
	scs.logger.Printf("Handling validation failure: %s", request.PreviousError)

	guidance := scs.enhancedInterface.GetSchemaBasedGuidance("validation_failed", map[string]interface{}{
		"error":       request.PreviousError,
		"entity_type": models.ToEntityType(request.EntityType),
		"operation":   request.Operation,
	})

	response := &ConsultationResponse{
		ShouldConsult:     true,
		Guidance:          guidance,
		RecommendedAction: "Use ValidateOperationRequest() to check parameter requirements",
		Confidence:        0.9,
		NextSteps: []string{
			"Fix parameter types to match schema requirements",
			"Add any missing required parameters",
			"Remove invalid or unsupported parameters",
			"Retry validation before executing operation",
		},
	}

	return response, nil
}

// handleOperationFailure helps recover from operation execution failures
func (scs *SchemaConsultationService) handleOperationFailure(request ConsultationRequest) (*ConsultationResponse, error) {
	scs.logger.Printf("Handling operation failure: %s", request.PreviousError)

	guidance := scs.enhancedInterface.GetSchemaBasedGuidance("operation_failed", request.Context)

	// Get alternative suggestions
	suggestions := scs.enhancedInterface.SuggestOperations(request.UserInput)

	response := &ConsultationResponse{
		ShouldConsult:     true,
		Guidance:          guidance,
		Suggestions:       suggestions,
		RecommendedAction: "Use SuggestOperations() to find alternative approaches",
		Confidence:        0.8,
		NextSteps: []string{
			"Try alternative operations that might achieve the same goal",
			"Verify entity type is correct for the operation",
			"Check if authorization/permissions are required",
			"Review error message for specific guidance",
		},
	}

	return response, nil
}

// handleGeneralUncertainty provides general guidance
func (scs *SchemaConsultationService) handleGeneralUncertainty(request ConsultationRequest) (*ConsultationResponse, error) {
	scs.logger.Printf("Handling general uncertainty for input: %s", request.UserInput)

	guidance := scs.enhancedInterface.GetSchemaBasedGuidance("", request.Context)
	contextualHelp := scs.enhancedInterface.GenerateContextualHelp(context.Background(), request.UserInput)

	response := &ConsultationResponse{
		ShouldConsult:     true,
		Guidance:          fmt.Sprintf("%s\n\n=== CONTEXTUAL HELP ===\n%s", guidance, contextualHelp),
		RecommendedAction: "Use GenerateContextualHelp() for specific guidance on your request",
		Confidence:        0.6,
		NextSteps: []string{
			"Clarify your intent with more specific language",
			"Identify the primary entity you want to work with",
			"Specify the action you want to perform",
			"Use schema methods to explore available options",
		},
	}

	return response, nil
}

// handleAutoDetectedUncertainty automatically detects and handles uncertainty
func (scs *SchemaConsultationService) handleAutoDetectedUncertainty(request ConsultationRequest) (*ConsultationResponse, error) {
	scs.logger.Printf("Auto-detecting uncertainty for input: %s", request.UserInput)

	// Analyze input to determine uncertainty type
	inputLower := strings.ToLower(request.UserInput)

	// Check for entity ambiguity
	entityMatches := scs.countEntityMatches(request.UserInput)
	if entityMatches > 1 {
		request.UncertaintyType = "ambiguous_entity"
		return scs.handleAmbiguousEntity(request)
	}

	// Check for operation uncertainty
	if scs.hasOperationUncertainty(inputLower) {
		request.UncertaintyType = "unclear_operation"
		return scs.handleUnclearOperation(request)
	}

	// Check for parameter uncertainty
	if scs.hasParameterUncertainty(inputLower) {
		request.UncertaintyType = "missing_parameters"
		return scs.handleMissingParameters(request)
	}

	// Default to general uncertainty
	request.UncertaintyType = "general_uncertainty"
	return scs.handleGeneralUncertainty(request)
}

// countEntityMatches counts how many entities could match the input
func (scs *SchemaConsultationService) countEntityMatches(userInput string) int {
	count := 0
	schemas := scs.enhancedInterface.SchemaRegistry.GetAllSchemas()

	for _, schema := range schemas {
		if scs.enhancedInterface.CalculateEntityRelevance(userInput, schema) > 0 {
			count++
		}
	}

	return count
}

// hasOperationUncertainty checks if the input suggests operation uncertainty
func (scs *SchemaConsultationService) hasOperationUncertainty(input string) bool {
	uncertaintyIndicators := []string{
		"how do i", "how can i", "what operation", "which operation",
		"how to", "what action", "what should i do",
	}

	for _, indicator := range uncertaintyIndicators {
		if strings.Contains(input, indicator) {
			return true
		}
	}

	return false
}

// hasParameterUncertainty checks if the input suggests parameter uncertainty
func (scs *SchemaConsultationService) hasParameterUncertainty(input string) bool {
	uncertaintyIndicators := []string{
		"what fields", "what parameters", "required fields", "missing parameter",
		"what data", "what information", "need to provide",
	}

	for _, indicator := range uncertaintyIndicators {
		if strings.Contains(input, indicator) {
			return true
		}
	}

	return false
}

// GetSchemaOverview provides a comprehensive overview of available schemas
func (scs *SchemaConsultationService) GetSchemaOverview() map[string]interface{} {
	schemas := scs.enhancedInterface.SchemaRegistry.GetAllSchemas()
	overview := make(map[string]interface{})

	overview["total_entities"] = len(schemas)
	overview["entities"] = make(map[string]interface{})

	totalOperations := 0
	for entityType, schema := range schemas {
		entityInfo := map[string]interface{}{
			"description":       schema.Description,
			"total_fields":      len(schema.Fields),
			"total_operations":  len(schema.Operations),
			"searchable_fields": len(schema.SearchableFields),
			"filterable_fields": len(schema.FilterableFields),
		}
		overview["entities"].(map[string]interface{})[string(entityType)] = entityInfo
		totalOperations += len(schema.Operations)
	}

	overview["total_operations"] = totalOperations
	overview["schema_consultation_available"] = true

	return overview
}

// ValidateWithSchema validates an operation against schema requirements
func (scs *SchemaConsultationService) ValidateWithSchema(entityType, operation string, parameters map[string]interface{}) error {
	return scs.enhancedInterface.ValidateOperationRequest(models.ToEntityType(entityType), operation, parameters)
}

// GetContextualGuidance provides contextual guidance for specific scenarios
func (scs *SchemaConsultationService) GetContextualGuidance(userQuery string) string {
	return scs.enhancedInterface.GenerateContextualHelp(context.Background(), userQuery)
}

// GetOperationSuggestions provides operation suggestions for user input
func (scs *SchemaConsultationService) GetOperationSuggestions(userInput string) []OperationSuggestion {
	return scs.enhancedInterface.SuggestOperations(userInput)
}

// GetSchemaAwarePrompt generates a context-specific prompt that includes schema consultation guidance
func (s *SchemaConsultationService) GetSchemaAwarePrompt(contextData map[string]interface{}) string {
	basePrompt := constants.SchemaAwareSystemPrompt

	// Add context-specific schema guidance
	if userQuery, exists := contextData["user_query"].(string); exists {
		ctx := context.Background()
		guidance := s.enhancedInterface.GenerateContextualHelp(ctx, userQuery)
		basePrompt += fmt.Sprintf("\n\n# 🎯 CONTEXT-SPECIFIC GUIDANCE\n\n%s", guidance)
	}

	// Add quick reference for immediate use
	basePrompt += "\n\n" + constants.SchemaQuickReference

	return basePrompt
}

// GetUncertaintyGuidance provides specific guidance when LLM encounters uncertainty
func (s *SchemaConsultationService) GetUncertaintyGuidance(uncertaintyType string, context map[string]interface{}) string {
	switch uncertaintyType {
	case "entity_ambiguous":
		return s.generateEntityDisambiguationHelp(context)
	case "operation_unclear":
		return s.generateOperationSelectionHelp(context)
	case "fields_unknown":
		return s.generateFieldDefinitionHelp(context)
	case "validation_failed":
		return s.generateValidationRecoveryHelp(context)
	case "multiple_options":
		return s.generateMultiOptionHelp(context)
	default:
		return s.generateGeneralUncertaintyHelp()
	}
}

// generateEntityDisambiguationHelp helps when entity type is unclear
func (s *SchemaConsultationService) generateEntityDisambiguationHelp(context map[string]interface{}) string {
	userInput, _ := context["user_input"].(string)

	help := "🔍 **ENTITY TYPE DISAMBIGUATION**\n\n"
	help += "Your request could apply to multiple entity types. Here are the relevant options:\n\n"

	// Get suggestions from enhanced interface
	suggestions := s.enhancedInterface.SuggestOperations(userInput)

	for _, suggestion := range suggestions {
		help += fmt.Sprintf("- **%s**: %s\n", suggestion.EntityType, suggestion.Description)
		if suggestion.Example != "" {
			help += fmt.Sprintf("  Example: %s\n", suggestion.Example)
		}
	}

	help += "\n**Recommendation**: Use GenerateContextualHelp() with your specific intent for targeted guidance."

	return help
}

// generateOperationSelectionHelp helps when operation intent is unclear
func (s *SchemaConsultationService) generateOperationSelectionHelp(context map[string]interface{}) string {
	entityType, _ := context["entity_type"].(string)

	help := "🛠️ **OPERATION SELECTION GUIDANCE**\n\n"

	if entityType != "" {
		help += fmt.Sprintf("For %s entity, available operations include:\n\n", entityType)
		// Here you would get operations from schema
		help += "**Recommendation**: Use GetOperationsByEntityType() to see all available operations.\n"
	} else {
		help += "**Recommendation**: Use SuggestOperations() to get ranked operation suggestions.\n"
	}

	return help
}

// generateFieldDefinitionHelp helps when field names/types are unknown
func (s *SchemaConsultationService) generateFieldDefinitionHelp(context map[string]interface{}) string {
	entityType, _ := context["entity_type"].(string)

	help := "📋 **FIELD DEFINITION GUIDANCE**\n\n"

	if entityType != "" {
		help += fmt.Sprintf("To get complete field definitions for %s:\n", entityType)
		help += fmt.Sprintf("Use: GetFieldsByEntityType(\"%s\")\n\n", entityType)
		help += "This will provide:\n"
		help += "- Field names and types\n"
		help += "- Required vs optional fields\n"
		help += "- Field constraints and validation rules\n"
		help += "- Default values\n"
	} else {
		help += "**Recommendation**: First identify the entity type, then use GetFieldsByEntityType().\n"
	}

	return help
}

// generateValidationRecoveryHelp helps when validation fails
func (s *SchemaConsultationService) generateValidationRecoveryHelp(context map[string]interface{}) string {
	validationError, _ := context["validation_error"].(string)

	help := "🔧 **VALIDATION RECOVERY GUIDANCE**\n\n"
	help += "Validation failed. Here's how to recover:\n\n"

	if validationError != "" {
		help += fmt.Sprintf("Error: %s\n\n", validationError)
	}

	help += "**Recovery Steps**:\n"
	help += "1. Use ValidateOperationRequest() to identify specific issues\n"
	help += "2. Use GetFieldsByEntityType() to check field requirements\n"
	help += "3. Fix parameters based on schema definitions\n"
	help += "4. Re-validate before execution\n"

	return help
}

// generateMultiOptionHelp helps when multiple approaches are possible
func (s *SchemaConsultationService) generateMultiOptionHelp(context map[string]interface{}) string {
	options, _ := context["options"].([]string)

	help := "🎯 **MULTIPLE OPTIONS GUIDANCE**\n\n"
	help += "Several approaches are possible. Consider:\n\n"

	for i, option := range options {
		help += fmt.Sprintf("%d. %s\n", i+1, option)
	}

	help += "\n**Recommendation**: Use the schema consultation methods to evaluate each option's feasibility."

	return help
}

// generateGeneralUncertaintyHelp provides general guidance for uncertainty
func (s *SchemaConsultationService) generateGeneralUncertaintyHelp() string {
	return constants.SchemaConsultationTriggers + "\n\n" + constants.SchemaQuickReference
}

// ProvideLearningGuidance helps users understand how to use schema consultation effectively
func (s *SchemaConsultationService) ProvideLearningGuidance() string {
	return `# 🎓 SCHEMA CONSULTATION LEARNING GUIDE

## 🚀 Quick Start

1. **When unclear about anything**: Start with GenerateContextualHelp(userQuery)
2. **Need to see all options**: Use GetAllSchemas() for system overview
3. **Before every operation**: Use ValidateOperationRequest() to ensure success

## 📚 Method Usage Patterns

### GenerateContextualHelp(userQuery)
- **Use for**: "I want to sell something", "find cheap items", "update my info"
- **Returns**: Specific entity types, operations, and parameter guidance
- **Example**: GenerateContextualHelp("I want to sell my car")

### SuggestOperations(input)
- **Use for**: "do something with my order", "handle my listing"  
- **Returns**: Ranked list of possible operations with scores
- **Example**: SuggestOperations("update my product")

### ValidateOperationRequest(entity, operation, parameters)
- **Use for**: Before every single operation execution
- **Returns**: Validation success/failure with specific error details
- **Example**: ValidateOperationRequest("product", "add", {...})

### GetFieldsByEntityType(entityType)
- **Use for**: Unknown field names, types, or requirements
- **Returns**: Complete field definitions with constraints
- **Example**: GetFieldsByEntityType("ProductType")

### GetOperationsByEntityType(entityType)
- **Use for**: "What can I do with products?"
- **Returns**: All available operations for that entity
- **Example**: GetOperationsByEntityType("VehicleType")

## 🎯 Best Practices

1. **Consult first, execute second** - Always validate before running operations
2. **Use specific entity types** - "ProductType" not "product"
3. **Leverage relationships** - Schema shows entity connections
4. **Provide alternatives** - When one approach fails, suggest others
5. **Explain schema reasoning** - Help users understand why certain approaches work

## 🚨 Common Pitfalls to Avoid

- ❌ Guessing field names instead of consulting schema
- ❌ Assuming parameter types without validation
- ❌ Ignoring required fields
- ❌ Not checking operation availability
- ❌ Missing relationship opportunities

## ✅ Success Indicators

- ✅ Validating parameters before execution
- ✅ Using exact field names from schemas
- ✅ Converting types correctly (e.g., prices to cents)
- ✅ Suggesting schema-aware next steps
- ✅ Providing helpful error recovery

Remember: The schema is your comprehensive guide to the entire system!`
}
