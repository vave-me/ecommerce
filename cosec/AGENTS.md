# AGENTS Guide: cosec

## Business Purpose
Security/coordination support module/assets.

## Key Files In This Directory
- `sfx_markt/cosec/cmd/service/main.go`
- `sfx_markt/cosec/internal/constants/constants.go`
- `sfx_markt/cosec/internal/handlers/integration_events.go`
- `sfx_markt/cosec/internal/handlers/integration_events_transaction.go`
- `sfx_markt/cosec/internal/handlers/replies.go`
- `sfx_markt/cosec/internal/handlers/replies_transaction.go`
- `sfx_markt/cosec/internal/models/data.go`
- `sfx_markt/cosec/internal/saga.go`
- `sfx_markt/cosec/migrations/001_create_tables.sql`
- `sfx_markt/cosec/migrations/migrations.go`
- `sfx_markt/cosec/module.go`

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
