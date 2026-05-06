# Prompt Issues Fixed

## 1. Schema Consultation Loop (Critical)

### Problem
- LLM stuck in infinite loop calling `schema_generate_contextual_help`
- Occurred when user asked general questions like "what can you do?"
- Loop continued for 10 iterations with empty responses
- Total execution time exceeded 17 seconds

### Root Cause
```go
// In llm_processor.go
if hasSchemaTools {
    continue // This bypassed ALL exit conditions!
}
```

### Fix Applied
1. Added iteration limit for schema consultations (max 3)
2. Improved help generation templates
3. Added explicit response requirements in prompts
4. Better handling of general capability questions

## 2. Empty Response Issue

### Problem
- LLM returning empty content (`''`) when confused
- Filtering of thinking blocks removed entire responses
- No fallback mechanism for empty content

### Fix Applied
1. Explicit instruction: "NEVER return empty responses"
2. Required response format with examples
3. Fallback templates for common scenarios
4. Better error recovery prompts

## 3. Language Support

### Problem
- All prompts assumed English
- Polish query "co mozesz dla mnie zrobic?" caused confusion
- No multilingual templates

### Fix Applied
1. Added language detection guidance
2. Multilingual help templates
3. Common phrases in multiple languages
4. Instruction to respond in user's language

## 4. Contradictory Instructions

### Problem
- "Plan before executing" vs "Execute immediately"
- "Zero tolerance for generic responses" vs need for help responses
- Thinking block visibility confusion

### Fix Applied
1. Clarified execution flow: Analyze → Plan → Validate → Execute → Respond
2. Removed contradictory directives
3. Clear thinking block usage instructions

## 5. Tool Priority Issues

### Problem
- Schema consultation tools had highest priority
- No clear hierarchy for tool selection
- LLM defaulted to schema tools when confused

### Fix Applied
1. Reorganized tool ordering
2. Added context-specific tool selection
3. Limited schema consultation usage
4. Better tool selection guidance

## 6. Security Vulnerabilities

### Problem
- No limits on tool execution
- Potential for prompt injection
- Missing validation requirements

### Fix Applied
1. Added explicit execution limits
2. Strengthened boundary definitions
3. Required validation steps
4. Security-first instructions

## 7. Complex Technical Language

### Problem
- Overly technical architectural comments
- Dense explanations confusing LLM
- Mixed implementation details with behavior

### Fix Applied
1. Simplified language
2. Removed architectural comments
3. Focus on clear, actionable instructions
4. User-friendly explanations

## 8. Scattered Prompts

### Problem
- Prompts hardcoded across multiple files
- Inconsistent formatting
- Difficult to maintain

### Fix Applied
1. Centralized all prompts in `/prompts` directory
2. Consistent markdown formatting
3. Version tracking
4. Clear categorization

## 9. Missing Response Examples

### Problem
- No clear examples of good responses
- LLM unsure how to format output
- Inconsistent response quality

### Fix Applied
1. Added response format examples
2. Template for common scenarios
3. Clear success/error patterns
4. Step-by-step response building

## 10. Validation Gaps

### Problem
- No size limits for data inputs
- Missing parameter validation
- Unclear schema requirements

### Fix Applied
1. Added explicit validation requirements
2. Size and parameter limits
3. Schema validation steps
4. Error boundary definitions

## Performance Improvements

### Before
- 17+ seconds for simple query
- 10 iterations of tool calls
- Empty final response

### After (Expected)
- <3 seconds for simple query
- 1-3 tool calls maximum
- Clear, helpful responses
- Multilingual support