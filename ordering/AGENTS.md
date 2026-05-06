# AGENTS Guide: ordering

## Business Purpose
Manages order lifecycle, order status transitions, and order orchestration.

## Service Runtime Shape
- Entry point: `sfx_markt/ordering/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/ordering/module.go`
- Codegen entrypoint: `sfx_markt/ordering/generate.go`
- Application facade: `sfx_markt/ordering/internal/application/application.go`
- gRPC server implementation: `sfx_markt/ordering/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/ordering/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/ordering/orderingpb/api.proto`
- Event contract (proto): `(not found)`
- Message contract (proto): `sfx_markt/ordering/orderingpb/messages.proto`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/ordering/orderingpb/api.proto`
- `sfx_markt/ordering/orderingpb/messages.proto`

### gRPC transport layer files
- `sfx_markt/ordering/internal/grpc/server.go`
- `sfx_markt/ordering/internal/grpc/server_transaction.go`

### Application command handlers
- `sfx_markt/ordering/internal/application/commands/approve_order.go`
- `sfx_markt/ordering/internal/application/commands/cancel_order.go`
- `sfx_markt/ordering/internal/application/commands/complete_order.go`
- `sfx_markt/ordering/internal/application/commands/create_order.go`
- `sfx_markt/ordering/internal/application/commands/deliver_order.go`
- `sfx_markt/ordering/internal/application/commands/ready_order.go`
- `sfx_markt/ordering/internal/application/commands/reject_order.go`
- `sfx_markt/ordering/internal/application/commands/ship_order.go`
- `sfx_markt/ordering/internal/application/commands/update_order_status.go`

### Application query handlers
- `sfx_markt/ordering/internal/application/queries/get_order.go`
- `sfx_markt/ordering/internal/application/queries/get_orders_by_customer.go`
- `sfx_markt/ordering/internal/application/queries/get_orders_by_status.go`
- `sfx_markt/ordering/internal/application/queries/list_orders.go`

### Domain model files
- `sfx_markt/ordering/internal/domain/catalog_repository.go`
- `sfx_markt/ordering/internal/domain/item.go`
- `sfx_markt/ordering/internal/domain/order.go`
- `sfx_markt/ordering/internal/domain/order_events.go`
- `sfx_markt/ordering/internal/domain/order_repository.go`
- `sfx_markt/ordering/internal/domain/order_snapshots.go`
- `sfx_markt/ordering/internal/domain/order_status.go`
- `sfx_markt/ordering/internal/domain/payment_repository.go`
- `sfx_markt/ordering/internal/domain/shopping_repository.go`
- `sfx_markt/ordering/internal/domain/user_repository.go`

### Event/integration/projection handlers
- `sfx_markt/ordering/internal/handlers/catalog_sync.go`
- `sfx_markt/ordering/internal/handlers/commands.go`
- `sfx_markt/ordering/internal/handlers/commands_transaction.go`
- `sfx_markt/ordering/internal/handlers/domain_events.go`
- `sfx_markt/ordering/internal/handlers/domain_events_transaction.go`

### Postgres/repository files
- `sfx_markt/ordering/internal/postgres/catalog_repository.go`

### SQL migrations
- `sfx_markt/ordering/migrations/001_create_tables.sql`

### RPC methods from api.proto
- `CreateOrder`
- `GetOrder`
- `ListOrders`
- `GetOrdersByCustomer`
- `GetOrdersByStatus`
- `CancelOrder`
- `ReadyOrder`
- `CompleteOrder`
- `ApproveOrder`
- `RejectOrder`
- `ShipOrder`
- `DeliverOrder`
- `UpdateOrderStatus`

### Event/message constants (channel names + event/command keys)
- (none found)
- `OrderAggregateChannel`
- `OrderCreatedEvent`
- `OrderRejectedEvent`
- `OrderApprovedEvent`
- `OrderReadiedEvent`
- `OrderCanceledEvent`
- `OrderShippedEvent`
- `OrderDeliveredEvent`
- `OrderCompletedEvent`
- `CommandChannel`
- `RejectOrderCommand`
- `ApproveOrderCommand`
- `CreateOrderCommand`

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./ordering/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./ordering/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile ordering up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/ordering && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/ordering && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./ordering/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/ordering/orderingpb/api.proto`.
2. Implement command object + handler in `internal/application/commands`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/ordering/orderingpb/api.proto`.
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
   - Use `(not found)` for event payload contract updates.
   - Use `sfx_markt/ordering/orderingpb/messages.proto` for message-bus command/event envelope updates.
4. Regenerate protobuf/generated artifacts (`go generate ./...`).
5. Update producer-side event publishing in `internal/handlers/domain_events*.go`.
6. Update consumer-side subscriptions/projections in `internal/handlers/integration_events*.go` and command-message handlers.
7. If removing events, remove constants, subscriptions, switch cases, and serde registrations together.

## SQL / Postgres Change Workflow
1. Add a new numbered SQL migration in `migrations/`; do not edit historical migrations.
2. Update read/write repositories under `internal/postgres` (or equivalent repository folders) to match schema updates.
3. Verify migration startup path in service bootstrap/main wiring.
4. Run smoke and tests:
   - `cd <workspace-root>/sfx_markt && go test ./ordering/...`
   - `cd <workspace-root>/sfx_markt && go run ./ordering/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.
