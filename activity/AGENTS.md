# AGENTS Guide: activity

## Business Purpose
Captures and serves user activity streams, interactions, and engagement projections.

## Service Runtime Shape
- Entry point: `sfx_markt/activity/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/activity/module.go`
- Codegen entrypoint: `sfx_markt/activity/generate.go`
- Application facade: `sfx_markt/activity/internal/application/application.go`
- gRPC server implementation: `sfx_markt/activity/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/activity/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/activity/activitypb/api.proto`
- Event contract (proto): `(not found)`
- Message contract (proto): `sfx_markt/activity/activitypb/messages.proto`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/activity/activitypb/api.proto`
- `sfx_markt/activity/activitypb/messages.proto`

### gRPC transport layer files
- `sfx_markt/activity/internal/grpc/product_repository.go`
- `sfx_markt/activity/internal/grpc/server.go`
- `sfx_markt/activity/internal/grpc/server_transaction.go`
- `sfx_markt/activity/internal/grpc/user_repository.go`

### Application command handlers
- `sfx_markt/activity/internal/application/commands/add_interaction.go`
- `sfx_markt/activity/internal/application/commands/create_activity.go`
- `sfx_markt/activity/internal/application/commands/remove_interaction.go`
- `sfx_markt/activity/internal/application/commands/update_interaction.go`

### Application query handlers
- `sfx_markt/activity/internal/application/queries/get_activities.go`
- `sfx_markt/activity/internal/application/queries/get_activity.go`
- `sfx_markt/activity/internal/application/queries/get_interaction.go`
- `sfx_markt/activity/internal/application/queries/get_interactions.go`
- `sfx_markt/activity/internal/application/queries/get_metrics.go`
- `sfx_markt/activity/internal/application/queries/get_most_disliked.go`
- `sfx_markt/activity/internal/application/queries/get_most_liked.go`

### Domain model files
- `sfx_markt/activity/internal/domain/activity.go`
- `sfx_markt/activity/internal/domain/activity_events.go`
- `sfx_markt/activity/internal/domain/activity_repostitory.go`
- `sfx_markt/activity/internal/domain/activity_snapshots.go`
- `sfx_markt/activity/internal/domain/interaction.go`
- `sfx_markt/activity/internal/domain/interaction_events.go`
- `sfx_markt/activity/internal/domain/interaction_repostitory.go`
- `sfx_markt/activity/internal/domain/interaction_snapshots.go`
- `sfx_markt/activity/internal/domain/item_type.go`
- `sfx_markt/activity/internal/domain/middleman_cache_interaction_repository.go`
- `sfx_markt/activity/internal/domain/middleman_cache_repository.go`
- `sfx_markt/activity/internal/domain/middleman_interaction_repository.go`
- `sfx_markt/activity/internal/domain/middleman_repository.go`
- `sfx_markt/activity/internal/domain/product.go`
- `sfx_markt/activity/internal/domain/product_cache_repository.go`
- `sfx_markt/activity/internal/domain/product_repository.go`
- `sfx_markt/activity/internal/domain/registrations.go`
- `sfx_markt/activity/internal/domain/user.go`
- `sfx_markt/activity/internal/domain/user_cache_repository.go`
- `sfx_markt/activity/internal/domain/user_repository.go`

### Event/integration/projection handlers
- `sfx_markt/activity/internal/handlers/domain_events.go`
- `sfx_markt/activity/internal/handlers/domain_events_transaction.go`
- `sfx_markt/activity/internal/handlers/integration_events.go`
- `sfx_markt/activity/internal/handlers/integration_events_transaction.go`
- `sfx_markt/activity/internal/handlers/middleman.go`
- `sfx_markt/activity/internal/handlers/middleman_interactions.go`

### Postgres/repository files
- `sfx_markt/activity/internal/postgres/middleman_interaction_repository.go`
- `sfx_markt/activity/internal/postgres/middleman_repository.go`
- `sfx_markt/activity/internal/postgres/product_cache_repository.go`
- `sfx_markt/activity/internal/postgres/user_cache_repository.go`

### SQL migrations
- `sfx_markt/activity/migrations/001_create_tables.sql`

### RPC methods from api.proto
- `CreateActivity`
- `GetActivity`
- `GetActivities`
- `RestoreActivity`
- `ArchiveActivity`
- `AddInteraction`
- `GetInteraction`
- `GetInteractions`
- `RemoveInteraction`
- `GetMostLiked`
- `GetMostDisliked`

### Event/message constants (channel names + event/command keys)
- (none found)
- `ActivityAggregateChannel`
- `InteractionAggregateChannel`
- `ActivityCreatedEvent`
- `InteractionAddedEvent`
- `InteractionUpdatedEvent`
- `InteractionRemovedEvent`

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./activity/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./activity/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile activity up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/activity && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/activity && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./activity/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/activity/activitypb/api.proto`.
2. Implement command object + handler in `internal/application/commands`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/activity/activitypb/api.proto`.
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
   - Use `sfx_markt/activity/activitypb/messages.proto` for message-bus command/event envelope updates.
4. Regenerate protobuf/generated artifacts (`go generate ./...`).
5. Update producer-side event publishing in `internal/handlers/domain_events*.go`.
6. Update consumer-side subscriptions/projections in `internal/handlers/integration_events*.go` and command-message handlers.
7. If removing events, remove constants, subscriptions, switch cases, and serde registrations together.

## SQL / Postgres Change Workflow
1. Add a new numbered SQL migration in `migrations/`; do not edit historical migrations.
2. Update read/write repositories under `internal/postgres` (or equivalent repository folders) to match schema updates.
3. Verify migration startup path in service bootstrap/main wiring.
4. Run smoke and tests:
   - `cd <workspace-root>/sfx_markt && go test ./activity/...`
   - `cd <workspace-root>/sfx_markt && go run ./activity/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.
