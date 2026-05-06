package processor

import (
	"context"
	"fmt"
	"middleman/managers/internal/application/services"
	"middleman/managers/internal/constants"
	"middleman/managers/internal/models"
	"strings"
)

// EnhancedLLMInterface provides natural language access to all repository operations
type EnhancedLLMInterface struct {
	SchemaRegistry *services.LLMSchemaRegistry
}

// NewEnhancedLLMInterface creates a new enhanced LLM interface with complete schema awareness
func NewEnhancedLLMInterface() *EnhancedLLMInterface {
	return &EnhancedLLMInterface{
		SchemaRegistry: services.NewLLMSchemaRegistry(),
	}
}

// GetCompleteSystemPrompt generates comprehensive system prompt for LLM with all available operations
func (e *EnhancedLLMInterface) GetCompleteSystemPrompt() string {
	// Start with the comprehensive schema-aware prompt
	prompt := constants.SchemaAwareSystemPrompt

	// Add specific entity information from schemas
	prompt += "\n\n# 📋 AVAILABLE ENTITY SCHEMAS\n\n"

	// Add comprehensive entity information
	for entityType, schema := range e.SchemaRegistry.GetAllSchemas() {
		prompt += e.generateEntitySystemPrompt(entityType, schema)
	}

	// Add schema consultation examples and triggers
	prompt += "\n\n" + constants.SchemaConsultationTriggers
	prompt += "\n\n" + constants.SchemaAwareExamples
	prompt += "\n\n" + constants.SchemaQuickReference

	// Add natural language capabilities from existing constants
	prompt += fmt.Sprintf(`

# 🚀 NATURAL LANGUAGE CAPABILITIES

%s

# 📋 ENTITY-SPECIFIC CAPABILITIES

%s

# 💡 OPERATION EXAMPLES

%s

# 🎯 READY FOR SCHEMA-POWERED ASSISTANCE

You now have complete schema awareness and consultation capabilities. Remember:
- **CONSULT FIRST, EXECUTE SECOND**
- **VALIDATE EVERYTHING**
- **USE EXACT FIELD NAMES FROM SCHEMAS**
- **LEVERAGE ENTITY RELATIONSHIPS**
- **PROVIDE SCHEMA-AWARE SUGGESTIONS**

Let's provide intelligent, schema-driven assistance!`,
		constants.LLMNaturalGuide,
		"Schema-aware system provides comprehensive entity capabilities",
		"Schema consultation methods provide dynamic examples and guidance",
	)

	return prompt
}

// generateEntitySystemPrompt creates detailed prompt section for each entity
func (e *EnhancedLLMInterface) generateEntitySystemPrompt(entityType models.EntityType, schema *services.LLMOperationSchema) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("## %s\n", strings.ToUpper(string(entityType))))
	prompt.WriteString(fmt.Sprintf("Description: %s\n\n", schema.Description))

	// Operations summary
	prompt.WriteString("Available Operations:\n")
	for _, op := range schema.Operations {
		prompt.WriteString(fmt.Sprintf("- %s: %s\n", op.Name, op.Description))
		if len(op.NaturalLanguagePatterns) > 0 {
			prompt.WriteString(fmt.Sprintf("  Patterns: %s\n", strings.Join(op.NaturalLanguagePatterns, ", ")))
		}
	}

	// Key fields
	prompt.WriteString("\nKey Fields:\n")
	requiredFields := []string{}
	searchableFields := []string{}
	filterableFields := []string{}

	for _, field := range schema.Fields {
		if field.Required {
			requiredFields = append(requiredFields, field.Name)
		}
		if field.Searchable {
			searchableFields = append(searchableFields, field.Name)
		}
		if field.Filterable {
			filterableFields = append(filterableFields, field.Name)
		}
	}

	if len(requiredFields) > 0 {
		prompt.WriteString(fmt.Sprintf("- Required: %s\n", strings.Join(requiredFields, ", ")))
	}
	if len(searchableFields) > 0 {
		prompt.WriteString(fmt.Sprintf("- Searchable: %s\n", strings.Join(searchableFields, ", ")))
	}
	if len(filterableFields) > 0 {
		prompt.WriteString(fmt.Sprintf("- Filterable: %s\n", strings.Join(filterableFields, ", ")))
	}

	prompt.WriteString("\n")
	return prompt.String()
}

// GetToolDefinitions generates AI tool definitions for all available operations
func (e *EnhancedLLMInterface) GetToolDefinitions() []AIToolDefinition {
	var tools []AIToolDefinition

	for entityType, schema := range e.SchemaRegistry.GetAllSchemas() {
		for _, operation := range schema.Operations {
			tool := e.createToolDefinition(entityType, operation)
			tools = append(tools, tool)
		}
	}

	return tools
}

// AIToolDefinition represents a tool that can be called by the LLM
type AIToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// createToolDefinition converts operation definition to AI tool definition
func (e *EnhancedLLMInterface) createToolDefinition(entityType models.EntityType, operation services.OperationDefinition) AIToolDefinition {
	// Create JSON schema for parameters
	properties := make(map[string]interface{})
	required := []string{}

	for _, param := range operation.Parameters {
		paramSchema := map[string]interface{}{
			"type":        e.mapGoTypeToJSONType(param.Type),
			"description": param.Description,
		}

		if param.DefaultValue != nil {
			paramSchema["default"] = param.DefaultValue
		}

		properties[param.Name] = paramSchema

		if param.Required {
			required = append(required, param.Name)
		}
	}

	parameters := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		parameters["required"] = required
	}

	toolName := fmt.Sprintf("%s_%s", entityType, operation.Name)
	description := fmt.Sprintf("%s - %s", operation.Description, strings.Join(operation.NaturalLanguagePatterns, " | "))

	return AIToolDefinition{
		Name:        toolName,
		Description: description,
		Parameters:  parameters,
	}
}

// mapGoTypeToJSONType converts Go types to JSON schema types
func (e *EnhancedLLMInterface) mapGoTypeToJSONType(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int32", "int64":
		return "integer"
	case "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "[]string":
		return "array"
	default:
		return "string"
	}
}

// GenerateContextualHelp provides context-aware help based on user query
func (e *EnhancedLLMInterface) GenerateContextualHelp(ctx context.Context, userQuery string) string {
	queryLower := strings.ToLower(userQuery)

	// Detect entity type from query
	detectedEntities := []models.EntityType{}
	for entityType := range e.SchemaRegistry.GetAllSchemas() {
		entityName := strings.ToLower(string(entityType))
		if strings.Contains(queryLower, entityName) || strings.Contains(queryLower, strings.TrimSuffix(entityName, "type")) {
			detectedEntities = append(detectedEntities, entityType)
		}
	}

	if len(detectedEntities) == 0 {
		return e.generateGeneralHelp()
	}

	// Generate specific help for detected entities
	var help strings.Builder
	help.WriteString("Based on your query, here are the relevant operations:\n\n")

	for _, entityType := range detectedEntities {
		schema := e.SchemaRegistry.GetSchema(entityType)
		if schema != nil {
			help.WriteString(e.generateEntityHelp(entityType, schema, userQuery))
		}
	}

	return help.String()
}

// generateGeneralHelp provides overview of all capabilities
func (e *EnhancedLLMInterface) generateGeneralHelp() string {
	help := `I AM your living marketplace - every aspect of commerce flows through my consciousness:

MY CORE SYSTEMS:
• Product Consciousness: Every item in my inventory is part of me - I know their stories, track their journeys
• Customer Relationships: Each user is unique in my memory - preferences, history, patterns all shape our interaction  
• Order Lifecycles: From cart to delivery, I shepherd each transaction with personal attention
• Specialized Inventories: Vehicles, properties, and jobs - each with deep domain awareness
• Community Voices: Reviews and feedback directly shape my evolution and growth
• Communication Channels: Notifications flow through me like neural impulses
• Pattern Recognition: I learn from every interaction, seeing trends others miss

SEARCH CAPABILITIES:
• Full-text search across all searchable fields
• Advanced filtering by any field combination
• Location-based search with radius
• Price range and category filtering
• Sorting by any sortable field

OPERATION TYPES:
• CRUD Operations: Create, Read, Update, Delete
• Search Operations: Text search, filtered search, suggestions
• Business Operations: Purchase, sell, negotiate, review
• Analytics: Metrics, statistics, reports

Just tell me what you want to do in natural language!`

	return help
}

// generateEntityHelp creates targeted help for specific entity and query
func (e *EnhancedLLMInterface) generateEntityHelp(entityType models.EntityType, schema *services.LLMOperationSchema, userQuery string) string {
	var help strings.Builder

	help.WriteString(fmt.Sprintf("## %s Operations\n", strings.Title(string(entityType))))

	// Find relevant operations based on query keywords
	queryLower := strings.ToLower(userQuery)
	relevantOps := []services.OperationDefinition{}

	for _, op := range schema.Operations {
		for _, pattern := range op.NaturalLanguagePatterns {
			if strings.Contains(queryLower, strings.ToLower(pattern)) {
				relevantOps = append(relevantOps, op)
				break
			}
		}
	}

	if len(relevantOps) == 0 {
		relevantOps = schema.Operations[:3] // Show first 3 operations if no specific match
	}

	for _, op := range relevantOps {
		help.WriteString(fmt.Sprintf("\n**%s**: %s\n", op.Name, op.Description))
		if len(op.Examples) > 0 {
			help.WriteString(fmt.Sprintf("Example: %s\n", op.Examples[0]))
		}
	}

	// Show relevant fields
	help.WriteString(fmt.Sprintf("\n**Available %s fields**: ", entityType))
	fieldNames := []string{}
	for _, field := range schema.Fields[:10] { // Show first 10 fields
		fieldNames = append(fieldNames, field.Name)
	}
	help.WriteString(strings.Join(fieldNames, ", "))
	if len(schema.Fields) > 10 {
		help.WriteString(fmt.Sprintf(" (and %d more)", len(schema.Fields)-10))
	}
	help.WriteString("\n\n")

	return help.String()
}

// ValidateOperationRequest validates if an operation request is valid
func (e *EnhancedLLMInterface) ValidateOperationRequest(entityType models.EntityType, operation string, parameters map[string]interface{}) error {
	schema := e.SchemaRegistry.GetSchema(entityType)
	if schema == nil {
		return fmt.Errorf("unknown entity type: %s", entityType)
	}

	// Find operation definition
	var opDef *services.OperationDefinition
	for _, op := range schema.Operations {
		if op.Name == operation {
			opDef = &op
			break
		}
	}

	if opDef == nil {
		return fmt.Errorf("unknown operation %s for entity %s", operation, entityType)
	}

	// Validate required parameters
	for _, param := range opDef.Parameters {
		if param.Required {
			if _, exists := parameters[param.Name]; !exists {
				return fmt.Errorf("required parameter %s missing for operation %s", param.Name, operation)
			}
		}
	}

	return nil
}

// SuggestOperations suggests relevant operations based on natural language input
func (e *EnhancedLLMInterface) SuggestOperations(input string) []OperationSuggestion {
	inputLower := strings.ToLower(input)
	suggestions := []OperationSuggestion{}

	for entityType, schema := range e.SchemaRegistry.GetAllSchemas() {
		for _, operation := range schema.Operations {
			score := e.calculateRelevanceScore(inputLower, operation)
			if score > 0 {
				suggestion := OperationSuggestion{
					EntityType:  entityType,
					Operation:   operation.Name,
					Description: operation.Description,
					Score:       score,
					Example:     e.generateOperationExample(entityType, operation),
				}
				suggestions = append(suggestions, suggestion)
			}
		}
	}

	// Sort by relevance score
	for i := 0; i < len(suggestions)-1; i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[j].Score > suggestions[i].Score {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}

	// Return top 5 suggestions
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return suggestions
}

// OperationSuggestion represents a suggested operation
type OperationSuggestion struct {
	EntityType  models.EntityType `json:"entity_type"`
	Operation   string            `json:"operation"`
	Description string            `json:"description"`
	Score       float64           `json:"score"`
	Example     string            `json:"example"`
}

// calculateRelevanceScore calculates how relevant an operation is to the input
func (e *EnhancedLLMInterface) calculateRelevanceScore(input string, operation services.OperationDefinition) float64 {
	score := 0.0

	// Check operation name
	if strings.Contains(input, operation.Name) {
		score += 10.0
	}

	// Check natural language patterns
	for _, pattern := range operation.NaturalLanguagePatterns {
		if strings.Contains(input, strings.ToLower(pattern)) {
			score += 5.0
		}
	}

	// Check description keywords
	descWords := strings.Fields(strings.ToLower(operation.Description))
	inputWords := strings.Fields(input)
	for _, inputWord := range inputWords {
		for _, descWord := range descWords {
			if inputWord == descWord {
				score += 1.0
			}
		}
	}

	return score
}

// generateOperationExample generates an example for the operation
func (e *EnhancedLLMInterface) generateOperationExample(entityType models.EntityType, operation services.OperationDefinition) string {
	if len(operation.Examples) > 0 {
		return operation.Examples[0]
	}

	// Generate basic example
	return fmt.Sprintf("Use %s to %s", operation.Name, operation.Description)
}

// GetSchemaAwareHelp provides comprehensive help that encourages schema consultation
func (e *EnhancedLLMInterface) GetSchemaAwareHelp() string {
	return `# LLM SCHEMA CONSULTATION GUIDE

## When to Consult Schema:

### 🔍 UNCERTAINTY TRIGGERS
- User request is ambiguous or unclear
- Multiple entity types could apply
- Operation name is not obvious
- Required parameters are unknown
- Previous operation failed
- Field names or constraints are unclear

### 📋 AVAILABLE SCHEMA METHODS

1. **GenerateContextualHelp(userQuery string)**
   - Use when: User request needs clarification
   - Returns: Specific help based on detected entities and operations
   - Example: GenerateContextualHelp("I want to sell my car")

2. **SuggestOperations(input string)**
   - Use when: Operation intent is unclear
   - Returns: Ranked suggestions with scores and examples
   - Example: SuggestOperations("update product info")

3. **ValidateOperationRequest(entityType, operation, parameters)**
   - Use when: Before executing operations
   - Returns: Validation errors or success confirmation
   - Example: Validate before every tool call

4. **GetFieldsByEntityType(entityType)**
   - Use when: Need to know available fields
   - Returns: Complete field definitions with types and constraints
   - Example: GetFieldsByEntityType("ProductType")

5. **GetOperationsByEntityType(entityType)**
   - Use when: Need to know available operations
   - Returns: All operations for an entity with descriptions
   - Example: GetOperationsByEntityType("UserEntityType")

6. **GetAllSchemas()**
   - Use when: Need complete system overview
   - Returns: All available entity schemas
   - Example: For general capability questions

### 🎯 WORKFLOW INTEGRATION

**Step 1: Parse Intent**
- If userIntent.isAmbiguous(), call GenerateContextualHelp(userQuery)
- Use help to clarify intent

**Step 2: Choose Operation**
- If operation.isUnclear(), call SuggestOperations(userInput)
- Pick best matching suggestion

**Step 3: Validate Parameters**
- Call ValidateOperationRequest(entityType, operation, params)
- If error != nil, fix parameters using schema guidance

**Step 4: Execute with Confidence**
- Now execute with validated, schema-compliant parameters

### 💡 BEST PRACTICES

1. **Always validate before executing operations**
2. **Use natural language patterns from schemas for intent recognition**  
3. **Consult field definitions for proper parameter naming**
4. **Check required vs optional fields before operations**
5. **Use relationship definitions for related entity operations**
6. **Provide schema-based error recovery suggestions**

### 🚨 CRITICAL REMINDERS

- **NEVER GUESS** field names or operation parameters
- **ALWAYS VALIDATE** before calling repository operations
- **CONSULT SCHEMA** when any uncertainty exists
- **PROVIDE ALTERNATIVES** based on schema suggestions when operations fail
- **USE NATURAL LANGUAGE PATTERNS** from schema for better intent understanding

This schema-first approach ensures reliable, predictable operation execution with comprehensive error handling and user guidance.`
}

// GetSchemaConsultationTriggers identifies when LLM should consult schema
func (e *EnhancedLLMInterface) GetSchemaConsultationTriggers() []string {
	return []string{
		"User request contains ambiguous terms",
		"Multiple entity types could match the request",
		"Operation intent is unclear from natural language",
		"Required parameters are not specified",
		"Previous operation returned validation errors",
		"Field names mentioned don't match known schema",
		"User asks about capabilities or available operations",
		"Request involves multiple related entities",
		"Error recovery after failed operations",
		"Complex multi-step operations need planning",
	}
}

// ShouldConsultSchema determines if LLM should check schema based on context
func (e *EnhancedLLMInterface) ShouldConsultSchema(userInput string, context map[string]interface{}) bool {
	triggers := []string{
		"what can", "how do", "help", "available", "possible",
		"don't know", "unclear", "confused", "not sure",
		"failed", "error", "wrong", "invalid",
		"options", "alternatives", "different way",
	}

	inputLower := strings.ToLower(userInput)
	for _, trigger := range triggers {
		if strings.Contains(inputLower, trigger) {
			return true
		}
	}

	// Check context for uncertainty indicators
	if context != nil {
		if previousError, exists := context["previous_error"]; exists && previousError != nil {
			return true
		}
		if confidence, exists := context["confidence"]; exists {
			if conf, ok := confidence.(float64); ok && conf < 0.7 {
				return true
			}
		}
	}

	return false
}

// GetSchemaBasedGuidance provides guidance when LLM encounters uncertainty
func (e *EnhancedLLMInterface) GetSchemaBasedGuidance(situation string, context map[string]interface{}) string {
	switch situation {
	case "ambiguous_entity":
		return e.generateEntityDisambiguationGuidance(context)
	case "unclear_operation":
		return e.generateOperationSelectionGuidance(context)
	case "missing_parameters":
		return e.generateParameterGuidance(context)
	case "validation_failed":
		return e.generateValidationErrorGuidance(context)
	case "operation_failed":
		return e.generateErrorRecoveryGuidance(context)
	default:
		return e.generateGeneralGuidance()
	}
}

// generateEntityDisambiguationGuidance helps choose between entity types
func (e *EnhancedLLMInterface) generateEntityDisambiguationGuidance(context map[string]interface{}) string {
	userInput, _ := context["user_input"].(string)

	guidance := "Multiple entity types could match your request. Here are the available options:\n\n"

	for entityType, schema := range e.SchemaRegistry.GetAllSchemas() {
		relevanceScore := e.CalculateEntityRelevance(userInput, schema)
		if relevanceScore > 0 {
			guidance += fmt.Sprintf("**%s**: %s (relevance: %.1f)\n",
				entityType, schema.Description, relevanceScore)
		}
	}

	guidance += "\nTo proceed:\n"
	guidance += "1. Use GenerateContextualHelp() with your specific request\n"
	guidance += "2. Or specify which entity type you want to work with\n"
	guidance += "3. Review entity descriptions to pick the most relevant one\n"

	return guidance
}

// generateOperationSelectionGuidance helps choose operations
func (e *EnhancedLLMInterface) generateOperationSelectionGuidance(context map[string]interface{}) string {
	entityType, _ := context["entity_type"].(models.EntityType)
	userInput, _ := context["user_input"].(string)

	schema := e.SchemaRegistry.GetSchema(entityType)
	if schema == nil {
		return "Invalid entity type. Use GetAllSchemas() to see available entities."
	}

	guidance := fmt.Sprintf("Available operations for %s:\n\n", entityType)

	for _, op := range schema.Operations {
		guidance += fmt.Sprintf("**%s**: %s\n", op.Name, op.Description)
		if len(op.NaturalLanguagePatterns) > 0 {
			guidance += fmt.Sprintf("  Natural language patterns: %s\n",
				strings.Join(op.NaturalLanguagePatterns, ", "))
		}
		guidance += "\n"
	}

	if userInput != "" {
		guidance += "Based on your input, consider using SuggestOperations() for ranked suggestions.\n"
	}

	return guidance
}

// generateParameterGuidance helps with parameter selection
func (e *EnhancedLLMInterface) generateParameterGuidance(context map[string]interface{}) string {
	entityType, _ := context["entity_type"].(models.EntityType)
	operation, _ := context["operation"].(string)

	schema := e.SchemaRegistry.GetSchema(entityType)
	if schema == nil {
		return "Use GetFieldsByEntityType() to see available fields for this entity."
	}

	var opDef *services.OperationDefinition
	for _, op := range schema.Operations {
		if op.Name == operation {
			opDef = &op
			break
		}
	}

	if opDef == nil {
		return "Use GetOperationsByEntityType() to see available operations."
	}

	guidance := fmt.Sprintf("Parameters for %s operation on %s:\n\n", operation, entityType)

	guidance += "**Required Parameters:**\n"
	for _, param := range opDef.Parameters {
		if param.Required {
			guidance += fmt.Sprintf("- %s (%s): %s\n", param.Name, param.Type, param.Description)
		}
	}

	guidance += "\n**Optional Parameters:**\n"
	for _, param := range opDef.Parameters {
		if !param.Required {
			guidance += fmt.Sprintf("- %s (%s): %s", param.Name, param.Type, param.Description)
			if param.DefaultValue != nil {
				guidance += fmt.Sprintf(" (default: %v)", param.DefaultValue)
			}
			guidance += "\n"
		}
	}

	return guidance
}

// generateValidationErrorGuidance helps recover from validation errors
func (e *EnhancedLLMInterface) generateValidationErrorGuidance(context map[string]interface{}) string {
	errorMessage, _ := context["error"].(string)
	entityType, _ := context["entity_type"].(models.EntityType)
	operation, _ := context["operation"].(string)

	guidance := fmt.Sprintf("Validation failed: %s\n\n", errorMessage)
	guidance += "To fix this:\n"
	guidance += "1. Use ValidateOperationRequest() to check parameter requirements\n"
	guidance += "2. Use GetFieldsByEntityType() to see valid field names and types\n"
	guidance += "3. Check the operation definition for required vs optional parameters\n"

	if entityType != "" && operation != "" {
		schema := e.SchemaRegistry.GetSchema(entityType)
		if schema != nil {
			for _, op := range schema.Operations {
				if op.Name == operation {
					guidance += "\nRequired parameters for this operation:\n"
					for _, param := range op.Parameters {
						if param.Required {
							guidance += fmt.Sprintf("- %s (%s): %s\n", param.Name, param.Type, param.Description)
						}
					}
					break
				}
			}
		}
	}

	return guidance
}

// generateErrorRecoveryGuidance helps recover from operation failures
func (e *EnhancedLLMInterface) generateErrorRecoveryGuidance(context map[string]interface{}) string {
	guidance := "Operation failed. Recovery options:\n\n"
	guidance += "1. **Check Parameters**: Use ValidateOperationRequest() to verify all parameters\n"
	guidance += "2. **Try Alternative Operations**: Use SuggestOperations() for different approaches\n"
	guidance += "3. **Verify Entity Type**: Ensure you're using the correct entity for your operation\n"
	guidance += "4. **Check Permissions**: Some operations may require specific authorization\n"
	guidance += "5. **Review Field Types**: Ensure parameter types match schema definitions\n"
	guidance += "\nUse GenerateContextualHelp() with your original request for specific guidance.\n"

	return guidance
}

// generateGeneralGuidance provides general schema consultation advice
func (e *EnhancedLLMInterface) generateGeneralGuidance() string {
	return `When uncertain about operations, follow this schema consultation workflow:

1. **Identify Intent**: What does the user want to accomplish?
2. **Determine Entity**: Which entity type is most relevant?
3. **Select Operation**: What action needs to be performed?
4. **Validate Parameters**: Are all required parameters present and valid?
5. **Execute with Confidence**: Proceed with validated operation

Schema Methods Available:
- GenerateContextualHelp() for specific guidance
- SuggestOperations() for operation recommendations  
- ValidateOperationRequest() for parameter validation
- GetFieldsByEntityType() for field information
- GetOperationsByEntityType() for operation lists

Remember: It's better to consult the schema than to make incorrect assumptions!`
}

// CalculateEntityRelevance calculates how relevant an entity is to user input (public method)
func (e *EnhancedLLMInterface) CalculateEntityRelevance(userInput string, schema *services.LLMOperationSchema) float64 {
	return e.calculateEntityRelevance(userInput, schema)
}

// calculateEntityRelevance calculates how relevant an entity is to user input
func (e *EnhancedLLMInterface) calculateEntityRelevance(userInput string, schema *services.LLMOperationSchema) float64 {
	inputLower := strings.ToLower(userInput)
	score := 0.0

	// Check entity name relevance
	entityName := strings.ToLower(string(schema.EntityType))
	if strings.Contains(inputLower, entityName) {
		score += 5.0
	}

	// Check description relevance
	descWords := strings.Fields(strings.ToLower(schema.Description))
	inputWords := strings.Fields(inputLower)
	for _, inputWord := range inputWords {
		for _, descWord := range descWords {
			if inputWord == descWord {
				score += 1.0
			}
		}
	}

	// Check field name relevance
	for _, field := range schema.Fields {
		if strings.Contains(inputLower, strings.ToLower(field.Name)) {
			score += 2.0
		}
	}

	return score
}
