# AGENTS Guide: users

## Business Purpose
Owns identity/auth/profile/account lifecycle and user query operations.

## Service Runtime Shape
- Entry point: `sfx_markt/users/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/users/module.go`
- Codegen entrypoint: `sfx_markt/users/generate.go`
- Application facade: `sfx_markt/users/internal/application/application.go`
- gRPC server implementation: `sfx_markt/users/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/users/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/users/userspb/api.proto`
- Event contract (proto): `sfx_markt/users/userspb/events.proto`
- Message contract (proto): `(not found)`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/users/userspb/api.proto`
- `sfx_markt/users/userspb/events.proto`

### gRPC transport layer files
- `sfx_markt/users/internal/grpc/server.go`
- `sfx_markt/users/internal/grpc/server_transaction.go`

### Application command handlers
- `sfx_markt/users/internal/application/commands/add_admin.go`
- `sfx_markt/users/internal/application/commands/add_user_location.go`
- `sfx_markt/users/internal/application/commands/add_user_thumbnail.go`
- `sfx_markt/users/internal/application/commands/archive_user_ident.go`
- `sfx_markt/users/internal/application/commands/authorize_user.go`
- `sfx_markt/users/internal/application/commands/blacklist_user.go`
- `sfx_markt/users/internal/application/commands/block_user.go`
- `sfx_markt/users/internal/application/commands/clear_tokens.go`
- `sfx_markt/users/internal/application/commands/create_user.go`
- `sfx_markt/users/internal/application/commands/create_user_eid.go`
- `sfx_markt/users/internal/application/commands/create_user_esign.go`
- `sfx_markt/users/internal/application/commands/create_user_from_google.go`
- `sfx_markt/users/internal/application/commands/create_user_video_ident.go`
- `sfx_markt/users/internal/application/commands/disable_user.go`
- `sfx_markt/users/internal/application/commands/enable_user.go`
- `sfx_markt/users/internal/application/commands/follow_user.go`
- `sfx_markt/users/internal/application/commands/forgot_password.go`
- `sfx_markt/users/internal/application/commands/kyc_verify.go`
- `sfx_markt/users/internal/application/commands/login_user.go`
- `sfx_markt/users/internal/application/commands/logout_user.go`
- `sfx_markt/users/internal/application/commands/mobile_login_with_google.go`
- `sfx_markt/users/internal/application/commands/refresh_auth_token.go`
- `sfx_markt/users/internal/application/commands/remove_user_location.go`
- `sfx_markt/users/internal/application/commands/remove_user_thumbnail.go`
- `sfx_markt/users/internal/application/commands/rename_user.go`
- `sfx_markt/users/internal/application/commands/report_user.go`
- `sfx_markt/users/internal/application/commands/reset_password.go`
- `sfx_markt/users/internal/application/commands/unblock_user.go`
- `sfx_markt/users/internal/application/commands/unfollow_user.go`
- `sfx_markt/users/internal/application/commands/update_user.go`
- `sfx_markt/users/internal/application/commands/update_user_background.go`
- `sfx_markt/users/internal/application/commands/update_user_location.go`
- `sfx_markt/users/internal/application/commands/update_user_thumbnail.go`
- `sfx_markt/users/internal/application/commands/web_login_with_google.go`

### Application query handlers
- `sfx_markt/users/internal/application/queries/get_base_user.go`
- `sfx_markt/users/internal/application/queries/get_enabled_users.go`
- `sfx_markt/users/internal/application/queries/get_user.go`
- `sfx_markt/users/internal/application/queries/get_user_by_google_id.go`
- `sfx_markt/users/internal/application/queries/get_user_by_mail.go`
- `sfx_markt/users/internal/application/queries/get_user_location.go`
- `sfx_markt/users/internal/application/queries/get_users.go`

### Domain model files
- `sfx_markt/users/internal/domain/fake_user_repository.go`
- `sfx_markt/users/internal/domain/middleman_repository.go`
- `sfx_markt/users/internal/domain/token_repository.go`
- `sfx_markt/users/internal/domain/totp_setup_detail.go`
- `sfx_markt/users/internal/domain/user.go`
- `sfx_markt/users/internal/domain/user_events.go`
- `sfx_markt/users/internal/domain/user_language.go`
- `sfx_markt/users/internal/domain/user_privacy.go`
- `sfx_markt/users/internal/domain/user_repository.go`
- `sfx_markt/users/internal/domain/user_role.go`
- `sfx_markt/users/internal/domain/user_snapshots.go`

### Event/integration/projection handlers
- `sfx_markt/users/internal/handlers/domain_events.go`
- `sfx_markt/users/internal/handlers/domain_events_transaction.go`
- `sfx_markt/users/internal/handlers/middleman.go`

### Postgres/repository files
- `sfx_markt/users/internal/postgres/middleman_repository.go`

### SQL migrations
- `sfx_markt/users/migrations/001_create_tables.sql`

### RPC methods from api.proto
- `CreateUser`
- `AddAdmin`
- `GetUser`
- `GetBaseUser`
- `GetUsers`
- `ListEnabledUsers`
- `UpdateUser`
- `RenameUser`
- `EnableUser`
- `DisableUser`
- `LoginUser`
- `WebLoginWithGoogle`
- `MobileLoginWithGoogle`
- `LogUserOut`
- `RefreshAuthToken`
- `ClearTokens`
- `ForgotPassword`
- `ResetPassword`

### Event/message constants (channel names + event/command keys)
- `UserAggregateChannel`
- `UserAuthorizedEvent`
- `UserCreatedEvent`
- `UserUpdatedEvent`
- `UserEnabledToggledEvent`
- `UserRenamedEvent`
- `UserLoggedInEvent`
- `UserLoggedOutEvent`
- `UserPasswordResetEvent`
- `UserPasswordResetRequestedEvent`
- `UserPasswordForgottenEvent`
- `UserTokenInvalidatedEvent`
- `UserTokenRefreshedEvent`
- (none found)

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./users/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./users/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile users up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/users && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/users && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./users/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/users/userspb/api.proto`.
2. Implement command object + handler in `internal/application/commands`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/users/userspb/api.proto`.
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
   - Use `sfx_markt/users/userspb/events.proto` for event payload contract updates.
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
   - `cd <workspace-root>/sfx_markt && go test ./users/...`
   - `cd <workspace-root>/sfx_markt && go run ./users/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.
