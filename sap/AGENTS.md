# AGENTS Guide: sap

## Business Purpose
SAP integration module/assets and related adapters.

## Key Files In This Directory
- `sfx_markt/sap/SAP_CONNECTOR_DOCUMENTATION.md`
- `sfx_markt/sap/SAP_CONNECTOR_IMPLEMENTATION_PLAN.md`
- `sfx_markt/sap/TO_DO_IMPORTANT.MD`
- `sfx_markt/sap/buf.gen.yaml`
- `sfx_markt/sap/buf.yaml`
- `sfx_markt/sap/cmd/service/main.go`
- `sfx_markt/sap/generate.txt`
- `sfx_markt/sap/internal/application/application.go`
- `sfx_markt/sap/internal/application/application_impl.go`
- `sfx_markt/sap/internal/application/commands/process_webhook.go`
- `sfx_markt/sap/internal/application/commands/sync_from_sap.go`
- `sfx_markt/sap/internal/application/mock_app.go`
- `sfx_markt/sap/internal/application/mock_commands.go`
- `sfx_markt/sap/internal/application/mock_queries.go`
- `sfx_markt/sap/internal/constants/constants.go`
- `sfx_markt/sap/internal/domain/sync_repository.go`
- `sfx_markt/sap/internal/domain/sync_status.go`
- `sfx_markt/sap/internal/grpc/server.go`
- `sfx_markt/sap/internal/grpc/server_transaction.go`
- `sfx_markt/sap/internal/handlers/catalog.go`
- `sfx_markt/sap/internal/handlers/catalog_variants.go`
- `sfx_markt/sap/internal/handlers/commands.go`
- `sfx_markt/sap/internal/handlers/commands_transaction.go`
- `sfx_markt/sap/internal/handlers/domain_events.go`
- `sfx_markt/sap/internal/handlers/domain_events_contract_test.go`
- `sfx_markt/sap/internal/handlers/domain_events_transaction.go`
- `sfx_markt/sap/internal/postgres/catalog_repository.go`
- `sfx_markt/sap/internal/postgres/catalog_variant_repository.go`
- `sfx_markt/sap/internal/rest/api.annotations.yaml`
- `sfx_markt/sap/internal/rest/api.openapi.yaml`
- `sfx_markt/sap/internal/rest/api.swagger.json`
- `sfx_markt/sap/internal/rest/gateway.go`
- `sfx_markt/sap/internal/rest/index.html`
- `sfx_markt/sap/internal/rest/swagger.go`
- `sfx_markt/sap/internal/rest/webhook.go`
- `sfx_markt/sap/internal/sap/client.go`
- `sfx_markt/sap/internal/sap/enhanced_client.go`
- `sfx_markt/sap/internal/sap/events.go`
- `sfx_markt/sap/internal/sap/hana_client.go`
- `sfx_markt/sap/internal/sap/idoc_parser.go`
- `sfx_markt/sap/internal/sap/security_client.go`
- `sfx_markt/sap/internal/sap/transformer/transformer.go`
- `sfx_markt/sap/migrations/001_create_tables.sql`
- `sfx_markt/sap/migrations/migrations.go`
- `sfx_markt/sap/module.go`
- `sfx_markt/sap/sappb/api.pb.go`
- `sfx_markt/sap/sappb/api.pb.gw.go`
- `sfx_markt/sap/sappb/api.proto`
- `sfx_markt/sap/sappb/api_grpc.pb.go`
- `sfx_markt/sap/sappb/asyncapi.go`
- `sfx_markt/sap/sappb/asyncapi.yaml`
- `sfx_markt/sap/sappb/common.pb.go`
- `sfx_markt/sap/sappb/common.proto`
- `sfx_markt/sap/sappb/css/asyncapi.min.css`
- `sfx_markt/sap/sappb/css/global.min.css`
- `sfx_markt/sap/sappb/events.go`
- `sfx_markt/sap/sappb/events.pb.go`
- `sfx_markt/sap/sappb/events.proto`
- `sfx_markt/sap/sappb/index.html`
- `sfx_markt/sap/sappb/js/asyncapi-ui.min.js`
- `sfx_markt/sap/sappb/mock_products_service_client.go`
- `sfx_markt/sap/sappb/mock_products_service_server.go`
- `sfx_markt/sap/sappb/mock_unsafe_products_service_server.go`

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
