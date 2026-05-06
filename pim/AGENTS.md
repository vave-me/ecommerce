# AGENTS Guide: pim

## Business Purpose
Owns product catalog, variants, pricing, stock, and public catalog read models.

## Service Runtime Shape
- Entry point: `sfx_markt/pim/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/pim/module.go`
- Codegen entrypoint: `sfx_markt/pim/generate.go`
- Application facade: `sfx_markt/pim/internal/application/application.go`
- gRPC server implementation: `sfx_markt/pim/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/pim/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/pim/productspb/api.proto`
- Event contract (proto): `sfx_markt/pim/productspb/events.proto`
- Message contract (proto): `(not found)`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/pim/productspb/api.proto`
- `sfx_markt/pim/productspb/common.proto`
- `sfx_markt/pim/productspb/events.proto`

### gRPC transport layer files
- `sfx_markt/pim/internal/grpc/server.go`
- `sfx_markt/pim/internal/grpc/server_transaction.go`

### Application command handlers
- `sfx_markt/pim/internal/application/commands/add_product.go`
- `sfx_markt/pim/internal/application/commands/add_product_location.go`
- `sfx_markt/pim/internal/application/commands/add_product_thumbnail.go`
- `sfx_markt/pim/internal/application/commands/add_variant.go`
- `sfx_markt/pim/internal/application/commands/add_variant_option.go`
- `sfx_markt/pim/internal/application/commands/adjust_product_stock.go`
- `sfx_markt/pim/internal/application/commands/adjust_variant_stock.go`
- `sfx_markt/pim/internal/application/commands/archive_product.go`
- `sfx_markt/pim/internal/application/commands/archive_variant.go`
- `sfx_markt/pim/internal/application/commands/change_product_location.go`
- `sfx_markt/pim/internal/application/commands/decrease_product_price.go`
- `sfx_markt/pim/internal/application/commands/decrease_variant_price.go`
- `sfx_markt/pim/internal/application/commands/increase_product_price.go`
- `sfx_markt/pim/internal/application/commands/increase_variant_price.go`
- `sfx_markt/pim/internal/application/commands/mark_product_leased.go`
- `sfx_markt/pim/internal/application/commands/mark_product_pawned.go`
- `sfx_markt/pim/internal/application/commands/mark_product_sold.go`
- `sfx_markt/pim/internal/application/commands/rebrand_product.go`
- `sfx_markt/pim/internal/application/commands/rebrand_variant.go`
- `sfx_markt/pim/internal/application/commands/release_product.go`
- `sfx_markt/pim/internal/application/commands/remove_product.go`
- `sfx_markt/pim/internal/application/commands/remove_variant.go`
- `sfx_markt/pim/internal/application/commands/remove_variant_option.go`
- `sfx_markt/pim/internal/application/commands/reserve_product.go`
- `sfx_markt/pim/internal/application/commands/toggle_product_negotiable.go`
- `sfx_markt/pim/internal/application/commands/update_product.go`
- `sfx_markt/pim/internal/application/commands/update_product_thumbnail.go`
- `sfx_markt/pim/internal/application/commands/update_variant_attributes.go`

### Application query handlers
- `sfx_markt/pim/internal/application/queries/get_available_variants.go`
- `sfx_markt/pim/internal/application/queries/get_catalog.go`
- `sfx_markt/pim/internal/application/queries/get_negotiable_products.go`
- `sfx_markt/pim/internal/application/queries/get_product.go`
- `sfx_markt/pim/internal/application/queries/get_products.go`
- `sfx_markt/pim/internal/application/queries/get_products_by_category.go`
- `sfx_markt/pim/internal/application/queries/get_products_by_category_slug.go`
- `sfx_markt/pim/internal/application/queries/get_products_with_filters.go`
- `sfx_markt/pim/internal/application/queries/get_products_with_term.go`
- `sfx_markt/pim/internal/application/queries/get_public_catalog.go`
- `sfx_markt/pim/internal/application/queries/get_variants.go`

### Domain model files
- `sfx_markt/pim/internal/domain/attributes.go`
- `sfx_markt/pim/internal/domain/catalog_cache_repository.go`
- `sfx_markt/pim/internal/domain/catalog_repository.go`
- `sfx_markt/pim/internal/domain/catalog_variant_cache_repository.go`
- `sfx_markt/pim/internal/domain/catalog_variant_repository.go`
- `sfx_markt/pim/internal/domain/fake_product_repository.go`
- `sfx_markt/pim/internal/domain/options.go`
- `sfx_markt/pim/internal/domain/product.go`
- `sfx_markt/pim/internal/domain/product_condition.go`
- `sfx_markt/pim/internal/domain/product_events.go`
- `sfx_markt/pim/internal/domain/product_repository.go`
- `sfx_markt/pim/internal/domain/product_snapshots.go`
- `sfx_markt/pim/internal/domain/product_status.go`
- `sfx_markt/pim/internal/domain/seller_type.go`
- `sfx_markt/pim/internal/domain/variant.go`
- `sfx_markt/pim/internal/domain/variant_events.go`
- `sfx_markt/pim/internal/domain/variant_repository.go`
- `sfx_markt/pim/internal/domain/variant_snapshots.go`

### Event/integration/projection handlers
- `sfx_markt/pim/internal/handlers/catalog.go`
- `sfx_markt/pim/internal/handlers/catalog_variants.go`
- `sfx_markt/pim/internal/handlers/commands.go`
- `sfx_markt/pim/internal/handlers/commands_transaction.go`
- `sfx_markt/pim/internal/handlers/domain_events.go`
- `sfx_markt/pim/internal/handlers/domain_events_transaction.go`

### Postgres/repository files
- `sfx_markt/pim/internal/postgres/catalog_repository.go`
- `sfx_markt/pim/internal/postgres/catalog_variant_repository.go`

### SQL migrations
- `sfx_markt/pim/migrations/001_create_tables.sql`

### RPC methods from api.proto
- `AddProduct`
- `GetProduct`
- `GetProducts`
- `GetProductsByCategory`
- `GetProductsByCategorySlug`
- `GetProductsWithFilters`
- `GetCatalog`
- `GetPublicCatalog`
- `RemoveProduct`
- `UpdateProductPrice`
- `RebrandProduct`
- `UpdateProduct`
- `AdjustProductStock`
- `ArchiveProduct`
- `MarkProductSold`
- `MarkProductLeased`
- `MarkProductPawned`
- `IncreaseProductPrice`
- `DecreaseProductPrice`
- `AddVariant`
- `GetVariant`
- `GetVariants`
- `IncreaseVariantPrice`
- `DecreaseVariantPrice`
- `AdjustVariantStock`
- `ArchiveVariant`
- `RemoveVariant`
- `AddProductThumbnail`
- `UpdateProductThumbnail`

### Event/message constants (channel names + event/command keys)
- `ProductAggregateChannel`
- `VariantAggregateChannel`
- `ProductAddedEvent`
- `ProductRebrandedEvent`
- `ProductUpdatedEvent`
- `ProductPriceIncreasedEvent`
- `ProductPriceDecreasedEvent`
- `ProductStockAdjustedEvent`
- `ProductRemovedEvent`
- `ProductThumbnailAddedEvent`
- `ProductThumbnailUpdatedEvent`
- `ProductArchivedEvent`
- `ProductNegotiableToggledEvent`
- `ProductSoldEvent`
- `ProductLeasedEvent`
- `ProductPawnedEvent`
- `VariantAddedEvent`
- `VariantPriceIncreasedEvent`
- `VariantPriceDecreasedEvent`
- `VariantStockAdjustedEvent`
- `VariantArchivedEvent`
- `VariantRemovedEvent`
- `CommandChannel`
- `ReserveProductCommand`
- `ReleaseProductCommand`
- `ReserveProductsCommand`
- `ReleaseProductsCommand`
- (none found)

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./pim/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./pim/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile pim up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/pim && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/pim && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./pim/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/pim/productspb/api.proto`.
2. Implement command object + handler in `internal/application/commands`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/pim/productspb/api.proto`.
2. Implement query handler in `internal/application/queries`.
3. Wire query handler into `internal/application/application.go`.
4. Map RPC to application query in `internal/grpc/server.go`.
5. Ensure transaction wrapper delegates query correctly in `internal/grpc/server_transaction.go` if this service wraps read paths.
6. Extend postgres/read-model repositories to fetch required fields and pagination/sorting behavior.
7. Add/adjust tests for handler + repository query behavior.

## Add or Remove Events (Detailed)
1. Define or remove event structs in `internal/domain` where aggregate state changes occur.
2. Update event registration/serde wiring in `module.go` (registrations and snapshot registrations).
3. Update wire contracts:
   - Use `sfx_markt/pim/productspb/events.proto` for event payload contract updates.
   - Use `(not found)` for message-bus command/event envelope updates.
4. Regenerate protobuf/generated artifacts (`go generate ./...`).
5. Update producer-side event publishing in `internal/handlers/domain_events*.go`.
6. Update consumer-side subscriptions/projections in `internal/handlers/integration_events*.go` and command-message handlers.
7. If removing events, remove constants, subscriptions, switch cases, and serde registrations together.

## SQL / Postgres Change Workflow
1. Add a new numbered SQL migration in `migrations/`; do not edit historical migrations.
2. Update read/write repositories under `internal/postgres` (or equivalent repository folders) to match schema updates.
3. Verify migration startup path in service bootstrap/main wiring.
4. Run smoke and tests:
   - `cd <workspace-root>/sfx_markt && go test ./pim/...`
   - `cd <workspace-root>/sfx_markt && go run ./pim/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.
