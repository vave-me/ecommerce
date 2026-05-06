# Centralized Prompts Directory

This directory contains all system prompts, templates, and AI instructions used across the application.

## Directory Structure

```
prompts/
├── system/                    # Core system prompts
│   ├── marketplace_ai.md      # Main marketplace AI prompt (fixed)
│   ├── vaver_system.md        # Vaver system prompt (fixed)
│   └── schema_aware.md        # Schema-aware system prompt (fixed)
├── assistant_types/           # Assistant type-specific prompts
│   ├── admin.md               # Admin assistant prompt
│   ├── business.md            # Business assistant prompt
│   ├── support.md             # Support assistant prompt
│   ├── scheduler.md           # Scheduler assistant prompt
│   └── standard.md            # Standard assistant prompt
├── processors/                # Content processor prompts
│   ├── document.md            # Document processing prompts
│   ├── data.md                # Data analysis prompts
│   ├── vision.md              # Vision processing prompts
│   └── audio.md               # Audio processing prompts
├── templates/                 # Reusable prompt templates
│   ├── help_generation.md     # Help and capability templates
│   ├── error_handling.md      # Error response templates
│   ├── tool_execution.md      # Tool execution guidance
│   └── nlp_patterns.md        # Natural language patterns
└── docs/                      # Documentation
    ├── issues_fixed.md        # Documentation of issues found and fixed
    ├── best_practices.md      # Prompt engineering best practices
    └── migration_guide.md     # Guide for migrating to centralized prompts

```

## Issues Fixed

### 1. Schema Consultation Loop
- **Problem**: Infinite loop when AI calls schema consultation tools repeatedly
- **Fix**: Added explicit limits and better exit conditions

### 2. Empty Response Handling
- **Problem**: AI returns empty content when confused
- **Fix**: Added explicit fallback instructions and response requirements

### 3. Language Support
- **Problem**: English-only assumptions throughout
- **Fix**: Added multilingual support templates

### 4. Contradictory Instructions
- **Problem**: Conflicting directives about planning vs executing
- **Fix**: Clarified execution flow and requirements

### 5. Scattered Prompts
- **Problem**: Prompts hardcoded across multiple files
- **Fix**: Centralized all prompts with consistent formatting

## Usage

All prompts should be loaded from this directory at runtime. Update references in code to point to these centralized files instead of hardcoded strings.

### Example Usage in Go:

```go
import (
    "embed"
    "path/filepath"
)

//go:embed prompts
var promptsFS embed.FS

func LoadPrompt(category, name string) (string, error) {
    path := filepath.Join("prompts", category, name + ".md")
    content, err := promptsFS.ReadFile(path)
    return string(content), err
}
```

## Versioning

Each prompt file includes a version header for tracking changes:
```
<!-- Version: 1.0.0 -->
<!-- Last Updated: 2025-07-21 -->
<!-- Status: Active -->
```

## Contributing

When adding or modifying prompts:
1. Follow the existing structure and formatting
2. Update the version header
3. Document changes in issues_fixed.md
4. Test thoroughly with the AI system
5. Consider multilingual implications