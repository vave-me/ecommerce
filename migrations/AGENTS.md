# AGENTS Guide: migrations

## Business Purpose
Schema migration assets used outside service-local migration folders.

## Key Files In This Directory
- `sfx_markt/migrations/001_create_basket_schema.sql`
- `sfx_markt/migrations/002_create_users_schema.sql`
- `sfx_markt/migrations/003_create_notifications_schema.sql`
- `sfx_markt/migrations/004_create_ordering_schema.sql`
- `sfx_markt/migrations/005_create_payments_schema.sql`
- `sfx_markt/migrations/006_create_search_schema.sql`
- `sfx_markt/migrations/migrations.go`

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
