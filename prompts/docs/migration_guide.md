# Prompt Migration Guide

## Overview

This guide helps migrate from scattered, hardcoded prompts to the centralized prompt system.

## File Mapping

### Old Location → New Location

1. **System Prompts**
   - `/internal/ai/contstats.go` (VaverSystemPrompt) → `/prompts/system/vaver_system.md`
   - `/assistants/internal/constants/schema_aware_prompts.go` → `/prompts/system/schema_aware_fixed.md`
   - `/assistants/internal/constants/llm_natural_interface.go` → `/prompts/templates/nlp_patterns.md`

2. **Assistant Type Prompts**
   - `/assistants/internal/domain/assistant_factory.go` (embedded) → `/prompts/assistant_types/*.md`

3. **Processor Prompts**
   - `/assistants/internal/application/processor/openai_document_processor.go` → `/prompts/processors/document.md`
   - `/assistants/internal/application/processor/openai_data_processor.go` → `/prompts/processors/data.md`

## Code Changes Required

### 1. Replace Hardcoded Constants

**Before:**
```go
const SchemaAwareSystemPrompt = `...long prompt...`
```

**After:**
```go
//go:embed prompts/system/schema_aware_fixed.md
var schemaAwarePrompt string
```

### 2. Update Assistant Factory

**Before:**
```go
adminPrompt := `You are an Admin Assistant with full platform access...`
```

**After:**
```go
adminPrompt, err := LoadPrompt("assistant_types", "admin")
```

### 3. Fix LLM Processor Loop

**Before:**
```go
if hasSchemaTools {
    continue // This caused infinite loops!
}
```

**After:**
```go
if hasSchemaTools && schemaIterations < 3 {
    schemaIterations++
    continue
} else if schemaIterations >= 3 {
    // Break loop and provide fallback response
    break
}
```

### 4. Add Response Validation

**Before:**
```go
return content, allActions, confidence, nil
```

**After:**
```go
// Ensure non-empty response
if content == "" {
    content = GenerateFallbackResponse(ctx, currentMessage, allActions)
}
return content, allActions, confidence, nil
```

## Implementation Steps

### Phase 1: Setup (Immediate)
1. ✅ Create `/prompts` directory structure
2. ✅ Copy and fix all prompts
3. ✅ Document issues and fixes
4. ⏳ Create prompt loader utility

### Phase 2: Migration (Next Sprint)
1. Update code references to use centralized prompts
2. Remove hardcoded prompt constants
3. Implement prompt versioning
4. Add prompt validation tests

### Phase 3: Enhancement (Future)
1. Add prompt hot-reloading
2. Implement A/B testing framework
3. Add analytics for prompt effectiveness
4. Create prompt editor UI

## Prompt Loader Implementation

```go
package prompts

import (
    "embed"
    "fmt"
    "path/filepath"
    "strings"
)

//go:embed all:prompts
var promptsFS embed.FS

type PromptLoader struct {
    cache map[string]string
}

func NewPromptLoader() *PromptLoader {
    return &PromptLoader{
        cache: make(map[string]string),
    }
}

func (pl *PromptLoader) Load(category, name string) (string, error) {
    key := fmt.Sprintf("%s/%s", category, name)
    
    // Check cache
    if cached, ok := pl.cache[key]; ok {
        return cached, nil
    }
    
    // Load from embedded FS
    path := filepath.Join("prompts", category, name + ".md")
    content, err := promptsFS.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("failed to load prompt %s: %w", path, err)
    }
    
    // Strip metadata comments
    lines := strings.Split(string(content), "\n")
    var promptLines []string
    inMetadata := false
    
    for _, line := range lines {
        if strings.HasPrefix(line, "<!--") {
            inMetadata = true
        } else if strings.HasPrefix(line, "-->") {
            inMetadata = false
            continue
        } else if !inMetadata {
            promptLines = append(promptLines, line)
        }
    }
    
    prompt := strings.TrimSpace(strings.Join(promptLines, "\n"))
    pl.cache[key] = prompt
    
    return prompt, nil
}
```

## Testing Checklist

- [ ] Schema consultation loop fixed
- [ ] Empty responses eliminated  
- [ ] Multilingual support working
- [ ] Response times under 3 seconds
- [ ] Tool execution limits enforced
- [ ] Error messages helpful
- [ ] Security boundaries maintained

## Rollback Plan

If issues arise:
1. Keep original prompt constants as backup
2. Use feature flags for gradual rollout
3. Monitor response quality metrics
4. Have quick revert capability

## Success Metrics

### Before Migration
- Response time: 17+ seconds
- Empty responses: Common
- Schema loops: Frequent
- Language support: English only

### After Migration (Target)
- Response time: <3 seconds
- Empty responses: None
- Schema loops: Prevented
- Language support: 5+ languages