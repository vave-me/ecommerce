# AGENTS Guide: baskets

## Business Purpose
Owns basket/cart state and basket lifecycle operations for commerce flows.

## Service Runtime Shape
- Entry point: `sfx_markt/baskets/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/baskets/module.go`
- Codegen entrypoint: `sfx_markt/baskets/generate.go`
- Application facade: `sfx_markt/baskets/internal/application/application.go`
- gRPC server implementation: `sfx_markt/baskets/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/baskets/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/baskets/basketspb/api.proto`
- Event contract (proto): `sfx_markt/baskets/basketspb/events.proto`
- Message contract (proto): `(not found)`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/baskets/basketspb/api.proto`
- `sfx_markt/baskets/basketspb/events.proto`

### gRPC transport layer files
- `sfx_markt/baskets/internal/grpc/product_repository.go`
- `sfx_markt/baskets/internal/grpc/server.go`
- `sfx_markt/baskets/internal/grpc/server_transaction.go`
- `sfx_markt/baskets/internal/grpc/user_repository.go`

### Application command handlers
- (none found)

### Application query handlers
- (none found)

### Domain model files
- `sfx_markt/baskets/internal/domain/attributes.go`
- `sfx_markt/baskets/internal/domain/basket.go`
- `sfx_markt/baskets/internal/domain/basket_events.go`
- `sfx_markt/baskets/internal/domain/basket_repository.go`
- `sfx_markt/baskets/internal/domain/basket_snapshots.go`
- `sfx_markt/baskets/internal/domain/basket_status.go`
- `sfx_markt/baskets/internal/domain/catalog_repository.go`
- `sfx_markt/baskets/internal/domain/fake_basket_repository.go`
- `sfx_markt/baskets/internal/domain/fake_product_cache_repository.go`
- `sfx_markt/baskets/internal/domain/fake_user_cache_repository.go`
- `sfx_markt/baskets/internal/domain/item.go`
- `sfx_markt/baskets/internal/domain/options.go`
- `sfx_markt/baskets/internal/domain/order_repository.go`
- `sfx_markt/baskets/internal/domain/product.go`
- `sfx_markt/baskets/internal/domain/product_cache_repository.go`
- `sfx_markt/baskets/internal/domain/product_condition.go`
- `sfx_markt/baskets/internal/domain/product_repository.go`
- `sfx_markt/baskets/internal/domain/product_status.go`
- `sfx_markt/baskets/internal/domain/registrations.go`
- `sfx_markt/baskets/internal/domain/seller_type.go`
- `sfx_markt/baskets/internal/domain/user.go`
- `sfx_markt/baskets/internal/domain/user_cache_repository.go`
- `sfx_markt/baskets/internal/domain/user_repository.go`

### Event/integration/projection handlers
- `sfx_markt/baskets/internal/handlers/catalog.go`
- `sfx_markt/baskets/internal/handlers/domain_events.go`
- `sfx_markt/baskets/internal/handlers/domain_events_transaction.go`
- `sfx_markt/baskets/internal/handlers/integration_events.go`
- `sfx_markt/baskets/internal/handlers/integration_events_transaction.go`

### Postgres/repository files
- `sfx_markt/baskets/internal/postgres/catalog_repository.go`
- `sfx_markt/baskets/internal/postgres/product_cache_repository.go`
- `sfx_markt/baskets/internal/postgres/user_cache_repository.go`

### SQL migrations
- `sfx_markt/baskets/migrations/001_create_tables.sql`

### RPC methods from api.proto
- `StartBasket`
- `GetBasket`
- `GetTotalBasket`
- `GetCurrentBasket`
- `CancelBasket`
- `CheckoutBasket`
- `AddItem`
- `RemoveItem`
- `UpdateItemQuantity`
- `ListBaskets`
- `GetActiveBaskets`
- `ReopenBasket`

### Event/message constants (channel names + event/command keys)
- `BasketAggregateChannel`
- `BasketStartedEvent`
- `BasketCanceledEvent`
- `BasketCheckedOutEvent`
- `BasketItemAddedEvent`
- `BasketItemRemovedEvent`
- (none found)

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./baskets/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./baskets/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile baskets up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/baskets && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/baskets && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./baskets/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/baskets/basketspb/api.proto`.
2. Implement command object + handler in application-level write methods in `internal/application/application.go`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/baskets/basketspb/api.proto`.
2. Implement query handler in application-level read methods in `internal/application/application.go`.
3. Wire query handler into `internal/application/application.go`.
4. Map RPC to application query in `internal/grpc/server.go`.
5. Ensure transaction wrapper delegates query correctly in `internal/grpc/server_transaction.go` if this service wraps read paths.
6. Extend postgres/read-model repositories to fetch required fields and pagination/sorting behavior.
7. Add/adjust tests for handler + repository query behavior.

## Add or Remove Events (Detailed)
1. Define or remove event structs in `internal/domain` where aggregate state changes occur.
2. Update event registration/serde wiring in `module.go` (registrations and snapshot registrations).
3. Update wire contracts:
   - Use `sfx_markt/baskets/basketspb/events.proto` for event payload contract updates.
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
   - `cd <workspace-root>/sfx_markt && go test ./baskets/...`
   - `cd <workspace-root>/sfx_markt && go run ./baskets/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.

## Service Shape Note
- This service does not fully use `internal/application/commands` + `internal/application/queries` for every operation.
- Treat `internal/application/application.go` and existing gRPC handlers as the command/query boundary.
- Follow existing service patterns instead of forcing a CQRS split where it is not already used.
