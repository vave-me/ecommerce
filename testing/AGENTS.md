# AGENTS Guide: testing

## Business Purpose
Integration/e2e test harness and testing utilities.

## Key Files In This Directory
- `sfx_markt/testing/e2e/baskets_feature_test.go`
- `sfx_markt/testing/e2e/customers_feature_test.go`
- `sfx_markt/testing/e2e/e2e_test.go`
- `sfx_markt/testing/e2e/features/baskets/add_item.feature`
- `sfx_markt/testing/e2e/features/baskets/checkout_basket.feature`
- `sfx_markt/testing/e2e/features/baskets/start_basket.feature`
- `sfx_markt/testing/e2e/features/kiosk/shopping.feature`
- `sfx_markt/testing/e2e/features/orders/processing.feature`
- `sfx_markt/testing/e2e/features/users/create_user.feature`
- `sfx_markt/testing/e2e/suite.go`
- `sfx_markt/testing/e2e/users_context.go`
- `sfx_markt/testing/e2e/users_feature_test.go`

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
