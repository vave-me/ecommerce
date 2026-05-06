# AGENTS Guide: scheduler

## Business Purpose
Schedules and executes time-based asynchronous workflows.

## Service Runtime Shape
- Entry point: `sfx_markt/scheduler/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/scheduler/module.go`
- Codegen entrypoint: `sfx_markt/scheduler/generate.go`
- Application facade: `sfx_markt/scheduler/internal/application/application.go`
- gRPC server implementation: `sfx_markt/scheduler/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/scheduler/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/scheduler/schedulerspb/api.proto`
- Event contract (proto): `(not found)`
- Message contract (proto): `sfx_markt/scheduler/schedulerspb/messages.proto`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/scheduler/schedulerspb/api.proto`
- `sfx_markt/scheduler/schedulerspb/messages.proto`

### gRPC transport layer files
- `sfx_markt/scheduler/internal/grpc/assistant_repository.go`
- `sfx_markt/scheduler/internal/grpc/server.go`
- `sfx_markt/scheduler/internal/grpc/server_transaction.go`

### Application command handlers
- `sfx_markt/scheduler/internal/application/commands/add_action.go`
- `sfx_markt/scheduler/internal/application/commands/cancel_task.go`
- `sfx_markt/scheduler/internal/application/commands/complete_task.go`
- `sfx_markt/scheduler/internal/application/commands/create_scheduler.go`
- `sfx_markt/scheduler/internal/application/commands/execute_task.go`
- `sfx_markt/scheduler/internal/application/commands/fail_task.go`
- `sfx_markt/scheduler/internal/application/commands/remove_iaction.go`
- `sfx_markt/scheduler/internal/application/commands/schedule_task.go`
- `sfx_markt/scheduler/internal/application/commands/update_action_status.go`
- `sfx_markt/scheduler/internal/application/commands/update_task.go`

### Application query handlers
- `sfx_markt/scheduler/internal/application/queries/count_tasks_by_manager_id.go`
- `sfx_markt/scheduler/internal/application/queries/get_action.go`
- `sfx_markt/scheduler/internal/application/queries/get_actions.go`
- `sfx_markt/scheduler/internal/application/queries/get_pending_actions.go`
- `sfx_markt/scheduler/internal/application/queries/get_pending_tasks.go`
- `sfx_markt/scheduler/internal/application/queries/get_scheduler.go`
- `sfx_markt/scheduler/internal/application/queries/get_schedulers.go`
- `sfx_markt/scheduler/internal/application/queries/get_task.go`
- `sfx_markt/scheduler/internal/application/queries/get_tasks.go`

### Domain model files
- `sfx_markt/scheduler/internal/domain/action.go`
- `sfx_markt/scheduler/internal/domain/action_events.go`
- `sfx_markt/scheduler/internal/domain/action_repository.go`
- `sfx_markt/scheduler/internal/domain/action_snapshots.go`
- `sfx_markt/scheduler/internal/domain/assistant_repository.go`
- `sfx_markt/scheduler/internal/domain/catalog_task_repository.go`
- `sfx_markt/scheduler/internal/domain/middleman_action_repository.go`
- `sfx_markt/scheduler/internal/domain/middleman_cache_action_repository.go`
- `sfx_markt/scheduler/internal/domain/middleman_cache_repository.go`
- `sfx_markt/scheduler/internal/domain/middleman_repository.go`
- `sfx_markt/scheduler/internal/domain/registrations.go`
- `sfx_markt/scheduler/internal/domain/scheduler.go`
- `sfx_markt/scheduler/internal/domain/scheduler_events.go`
- `sfx_markt/scheduler/internal/domain/scheduler_repository.go`
- `sfx_markt/scheduler/internal/domain/scheduler_snapshots.go`
- `sfx_markt/scheduler/internal/domain/task.go`
- `sfx_markt/scheduler/internal/domain/task_events.go`
- `sfx_markt/scheduler/internal/domain/task_repository.go`
- `sfx_markt/scheduler/internal/domain/task_snapshots.go`

### Event/integration/projection handlers
- `sfx_markt/scheduler/internal/handlers/domain_events.go`
- `sfx_markt/scheduler/internal/handlers/domain_events_transaction.go`
- `sfx_markt/scheduler/internal/handlers/middleman.go`
- `sfx_markt/scheduler/internal/handlers/middleman_actions.go`
- `sfx_markt/scheduler/internal/handlers/task_events.go`

### Postgres/repository files
- `sfx_markt/scheduler/internal/postgres/catalog_task_repository.go`
- `sfx_markt/scheduler/internal/postgres/middleman_action_repository.go`
- `sfx_markt/scheduler/internal/postgres/middleman_repository.go`
- `sfx_markt/scheduler/internal/postgres/task_repository.go`

### SQL migrations
- `sfx_markt/scheduler/migrations/001_create_tables.sql`

### RPC methods from api.proto
- `CreateScheduler`
- `GetScheduler`
- `GetSchedulers`
- `AddAction`
- `GetAction`
- `GetActions`
- `RemoveAction`
- `ScheduleTask`
- `CancelTask`
- `GetScheduledTasks`
- `UpdateTask`

### Event/message constants (channel names + event/command keys)
- (none found)
- `SchedulerAggregateChannel`
- `ActionAggregateChannel`
- `SchedulerCreatedEvent`
- `ActionAddedEvent`
- `ActionUpdatedEvent`
- `ActionRemovedEvent`

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./scheduler/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./scheduler/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile scheduler up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/scheduler && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/scheduler && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./scheduler/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/scheduler/schedulerspb/api.proto`.
2. Implement command object + handler in `internal/application/commands`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/scheduler/schedulerspb/api.proto`.
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
   - Use `sfx_markt/scheduler/schedulerspb/messages.proto` for message-bus command/event envelope updates.
4. Regenerate protobuf/generated artifacts (`go generate ./...`).
5. Update producer-side event publishing in `internal/handlers/domain_events*.go`.
6. Update consumer-side subscriptions/projections in `internal/handlers/integration_events*.go` and command-message handlers.
7. If removing events, remove constants, subscriptions, switch cases, and serde registrations together.

## SQL / Postgres Change Workflow
1. Add a new numbered SQL migration in `migrations/`; do not edit historical migrations.
2. Update read/write repositories under `internal/postgres` (or equivalent repository folders) to match schema updates.
3. Verify migration startup path in service bootstrap/main wiring.
4. Run smoke and tests:
   - `cd <workspace-root>/sfx_markt && go test ./scheduler/...`
   - `cd <workspace-root>/sfx_markt && go run ./scheduler/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.
