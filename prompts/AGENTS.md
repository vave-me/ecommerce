# AGENTS Guide: prompts

## Business Purpose
Prompt templates and prompt management resources.

## Key Files In This Directory
- `sfx_markt/prompts/README.md`
- `sfx_markt/prompts/assistant_types/energy_advisor.md`
- `sfx_markt/prompts/assistant_types/energy_advisor_en.md`
- `sfx_markt/prompts/docs/issues_fixed.md`
- `sfx_markt/prompts/docs/migration_guide.md`
- `sfx_markt/prompts/system/llm_iteration_awareness.md`
- `sfx_markt/prompts/system/marketplace_ai.md`
- `sfx_markt/prompts/system/natural_marketplace_ai.md`
- `sfx_markt/prompts/system/schema_aware_fixed.md`
- `sfx_markt/prompts/system/store_consciousness.md`
- `sfx_markt/prompts/templates/energy_calculations.md`
- `sfx_markt/prompts/templates/energy_installation_guide.md`
- `sfx_markt/prompts/templates/help_generation.md`
- `sfx_markt/prompts/templates/tool_execution.md`
- `sfx_markt/prompts/tools/energy_advisor_tools.md`

## How To Work In This Directory
1. Keep changes aligned with the owning service/business module contracts.
2. Do not edit generated/build/vendor artifacts directly.
3. Validate with targeted build/test/deploy commands relevant to affected modules.

## Relationship To Commands / Queries / Proto / Events
- This directory is support/infrastructure/integration-oriented unless it contains a dedicated service module pattern.
- For adding/removing commands, queries, RPCs, domain events, and SQL read models, edit the owning service directory (for example `sfx_markt/users`, `sfx_markt/products`, `sfx_markt/ordering`, `sfx_markt/payments`, `sfx_markt/search`).
- Regeneration baseline when shared contracts change:
  - `cd <workspace-root>/sfx_markt && go generate ./...`

## SQL / Data Safety
- Prefer additive, forward-only schema/data changes.
- Never rewrite historical migrations or previously applied seed files.
