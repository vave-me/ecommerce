# AGENTS Guide: core-api-modules-vaveme

## Business Purpose
Shared core API modules reused across services.

## Key Files In This Directory
- `sfx_markt/core-api-modules-vaveme/.gitignore`
- `sfx_markt/core-api-modules-vaveme/MIGRATION.md`
- `sfx_markt/core-api-modules-vaveme/README.md`
- `sfx_markt/core-api-modules-vaveme/__tests__/auth.service.test.ts`
- `sfx_markt/core-api-modules-vaveme/examples/migration-example.js`
- `sfx_markt/core-api-modules-vaveme/examples/quick-start.ts`
- `sfx_markt/core-api-modules-vaveme/jest.config.js`
- `sfx_markt/core-api-modules-vaveme/package.json`
- `sfx_markt/core-api-modules-vaveme/rollup.config.mjs`
- `sfx_markt/core-api-modules-vaveme/src/clients/base-client.ts`
- `sfx_markt/core-api-modules-vaveme/src/core/axios-client.ts`
- `sfx_markt/core-api-modules-vaveme/src/core/config.ts`
- `sfx_markt/core-api-modules-vaveme/src/core/error-handler.ts`
- `sfx_markt/core-api-modules-vaveme/src/core/token-manager.ts`
- `sfx_markt/core-api-modules-vaveme/src/index.ts`
- `sfx_markt/core-api-modules-vaveme/src/services/auth/auth.service.ts`
- `sfx_markt/core-api-modules-vaveme/src/services/search/search.service.ts`
- `sfx_markt/core-api-modules-vaveme/src/services/user/user.service.ts`
- `sfx_markt/core-api-modules-vaveme/src/utils/encoders.ts`
- `sfx_markt/core-api-modules-vaveme/src/utils/mappers.ts`
- `sfx_markt/core-api-modules-vaveme/src/utils/validators.ts`
- `sfx_markt/core-api-modules-vaveme/tsconfig.json`

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
