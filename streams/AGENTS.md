# AGENTS Guide: streams

## Business Purpose
Manages streaming content catalogs, access control, pricing, watch progress, and delivery webhooks.

## Service Runtime Shape
- Entry point: `sfx_markt/streams/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/streams/module.go`
- Codegen entrypoint: `sfx_markt/streams/generate.go`
- Application facade: `sfx_markt/streams/internal/application/application.go`
- gRPC server implementation: `sfx_markt/streams/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/streams/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/streams/streamspb/api.proto`
- Event contract (proto): `sfx_markt/streams/streamspb/events.proto`
- Message contract (proto): `(not found)`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/streams/streamspb/api.proto`
- `sfx_markt/streams/streamspb/common.proto`
- `sfx_markt/streams/streamspb/events.proto`

### gRPC transport layer files
- `sfx_markt/streams/internal/grpc/server.go`
- `sfx_markt/streams/internal/grpc/server_transaction.go`

### Application command handlers
- `sfx_markt/streams/internal/application/commands/add_episode.go`
- `sfx_markt/streams/internal/application/commands/add_season.go`
- `sfx_markt/streams/internal/application/commands/configure_streaming.go`
- `sfx_markt/streams/internal/application/commands/create_live_stream.go`
- `sfx_markt/streams/internal/application/commands/create_series.go`
- `sfx_markt/streams/internal/application/commands/create_stream.go`
- `sfx_markt/streams/internal/application/commands/grant_user_access.go`
- `sfx_markt/streams/internal/application/commands/publish_stream.go`
- `sfx_markt/streams/internal/application/commands/set_stream_pricing.go`
- `sfx_markt/streams/internal/application/commands/start_live_stream.go`
- `sfx_markt/streams/internal/application/commands/subscribe_webhook.go`
- `sfx_markt/streams/internal/application/commands/update_watch_progress.go`

### Application query handlers
- `sfx_markt/streams/internal/application/queries/get_catalog.go`
- `sfx_markt/streams/internal/application/queries/get_continue_watching.go`
- `sfx_markt/streams/internal/application/queries/get_live_stream.go`
- `sfx_markt/streams/internal/application/queries/get_series.go`
- `sfx_markt/streams/internal/application/queries/get_stream.go`
- `sfx_markt/streams/internal/application/queries/get_user_streams.go`
- `sfx_markt/streams/internal/application/queries/get_webhook_subscriptions.go`
- `sfx_markt/streams/internal/application/queries/search_streams.go`

### Domain model files
- `sfx_markt/streams/internal/domain/catalog_cache_repository.go`
- `sfx_markt/streams/internal/domain/catalog_repository.go`
- `sfx_markt/streams/internal/domain/errors.go`
- `sfx_markt/streams/internal/domain/live_stream.go`
- `sfx_markt/streams/internal/domain/live_stream_events.go`
- `sfx_markt/streams/internal/domain/live_stream_repository.go`
- `sfx_markt/streams/internal/domain/series.go`
- `sfx_markt/streams/internal/domain/stream.go`
- `sfx_markt/streams/internal/domain/stream_events.go`
- `sfx_markt/streams/internal/domain/stream_repository.go`
- `sfx_markt/streams/internal/domain/stream_snapshots.go`
- `sfx_markt/streams/internal/domain/webhook_subscription.go`

### Event/integration/projection handlers
- `sfx_markt/streams/internal/handlers/catalog.go`
- `sfx_markt/streams/internal/handlers/catalog_variants.go`
- `sfx_markt/streams/internal/handlers/commands.go`
- `sfx_markt/streams/internal/handlers/commands_transaction.go`
- `sfx_markt/streams/internal/handlers/domain_events.go`
- `sfx_markt/streams/internal/handlers/domain_events_live_streaming.go`
- `sfx_markt/streams/internal/handlers/domain_events_transaction.go`
- `sfx_markt/streams/internal/handlers/webhook_dispatcher.go`

### Postgres/repository files
- `sfx_markt/streams/internal/postgres/catalog_repository.go`
- `sfx_markt/streams/internal/postgres/catalog_variant_repository.go`
- `sfx_markt/streams/internal/postgres/webhook_delivery_repository.go`
- `sfx_markt/streams/internal/postgres/webhook_subscription_repository.go`

### SQL migrations
- `sfx_markt/streams/migrations/001_create_tables.sql`
- `sfx_markt/streams/migrations/002_add_webhooks.sql`
- `sfx_markt/streams/migrations/003_add_production_indexes.sql`

### RPC methods from api.proto
- `CreateStream`
- `GetStream`
- `UpdateStream`
- `PublishStream`
- `ArchiveStream`
- `SetStreamPricing`
- `GrantUserAccess`
- `RevokeUserAccess`
- `UpdateWatchProgress`
- `RateStream`
- `GetCatalog`
- `SearchStreams`
- `GetUserStreams`
- `GetContinueWatching`
- `CreateSeries`
- `GetSeries`
- `AddSeason`
- `AddEpisode`

### Event/message constants (channel names + event/command keys)
- (none found)
- (none found)

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./streams/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./streams/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile streams up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/streams && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/streams && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./streams/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/streams/streamspb/api.proto`.
2. Implement command object + handler in `internal/application/commands`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/streams/streamspb/api.proto`.
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
   - Use `sfx_markt/streams/streamspb/events.proto` for event payload contract updates.
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
   - `cd <workspace-root>/sfx_markt && go test ./streams/...`
   - `cd <workspace-root>/sfx_markt && go run ./streams/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.
