package tools

import ai2 "middleman/internal/ai"

// ConvertToolsToDefinitions converts []Tool to []ToolDefinition for compatibility
func ConvertToolsToDefinitions(tools []ai2.Tool) []ai2.ToolDefinition {
	definitions := make([]ai2.ToolDefinition, len(tools))
	for i, tool := range tools {
		definitions[i] = ai2.ToolDefinition{
			Type:     tool.Type,
			Function: tool.Function,
		}
	}
	return definitions
}

// ConvertDefinitionsToTools converts []ToolDefinition to []Tool for compatibility
func ConvertDefinitionsToTools(definitions []ai2.ToolDefinition) []ai2.Tool {
	tools := make([]ai2.Tool, len(definitions))
	for i, def := range definitions {
		tools[i] = ai2.Tool{
			Type:     def.Type,
			Function: def.Function,
		}
	}
	return tools
}