# AGENTS Guide: media

## Business Purpose
Manages media metadata and media linkage to domain entities.

## Service Runtime Shape
- Entry point: `sfx_markt/media/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/media/module.go`
- Codegen entrypoint: `sfx_markt/media/generate.go`
- Application facade: `sfx_markt/media/internal/application/application.go`
- gRPC server implementation: `sfx_markt/media/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/media/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/media/mediapb/api.proto`
- Event contract (proto): `sfx_markt/media/mediapb/events.proto`
- Message contract (proto): `(not found)`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/media/mediapb/api.proto`
- `sfx_markt/media/mediapb/events.proto`

### gRPC transport layer files
- `sfx_markt/media/internal/grpc/product_repository.go`
- `sfx_markt/media/internal/grpc/server.go`
- `sfx_markt/media/internal/grpc/server_transaction.go`

### Application command handlers
- `sfx_markt/media/internal/application/commands/add_image.go`
- `sfx_markt/media/internal/application/commands/add_import_batch.go`
- `sfx_markt/media/internal/application/commands/add_video.go`
- `sfx_markt/media/internal/application/commands/cancel_import.go`
- `sfx_markt/media/internal/application/commands/create_media.go`
- `sfx_markt/media/internal/application/commands/has_media.go`
- `sfx_markt/media/internal/application/commands/remove_image.go`
- `sfx_markt/media/internal/application/commands/remove_media.go`
- `sfx_markt/media/internal/application/commands/remove_video.go`
- `sfx_markt/media/internal/application/commands/start_bulk_import.go`
- `sfx_markt/media/internal/application/commands/update_media.go`

### Application query handlers
- `sfx_markt/media/internal/application/queries/get_all_images_by_author.go`
- `sfx_markt/media/internal/application/queries/get_all_item_images.go`
- `sfx_markt/media/internal/application/queries/get_all_item_videos.go`
- `sfx_markt/media/internal/application/queries/get_all_media_images.go`
- `sfx_markt/media/internal/application/queries/get_all_media_videos.go`
- `sfx_markt/media/internal/application/queries/get_all_videos_by_author.go`
- `sfx_markt/media/internal/application/queries/get_import_status.go`
- `sfx_markt/media/internal/application/queries/get_media.go`
- `sfx_markt/media/internal/application/queries/get_media_by_item.go`
- `sfx_markt/media/internal/application/queries/get_videos.go`

### Domain model files
- `sfx_markt/media/internal/domain/image.go`
- `sfx_markt/media/internal/domain/image_events.go`
- `sfx_markt/media/internal/domain/image_repository.go`
- `sfx_markt/media/internal/domain/image_snapshots.go`
- `sfx_markt/media/internal/domain/import_events.go`
- `sfx_markt/media/internal/domain/import_repository.go`
- `sfx_markt/media/internal/domain/import_session.go`
- `sfx_markt/media/internal/domain/importer.go`
- `sfx_markt/media/internal/domain/importer_snapshots.go`
- `sfx_markt/media/internal/domain/media.go`
- `sfx_markt/media/internal/domain/media_events.go`
- `sfx_markt/media/internal/domain/media_repository.go`
- `sfx_markt/media/internal/domain/media_snapshots.go`
- `sfx_markt/media/internal/domain/media_type.go`
- `sfx_markt/media/internal/domain/middleman_image_repository.go`
- `sfx_markt/media/internal/domain/middleman_media_repository.go`
- `sfx_markt/media/internal/domain/middleman_video_repository.go`
- `sfx_markt/media/internal/domain/video.go`
- `sfx_markt/media/internal/domain/video_events.go`
- `sfx_markt/media/internal/domain/video_repository.go`
- `sfx_markt/media/internal/domain/video_snapshots.go`

### Event/integration/projection handlers
- `sfx_markt/media/internal/handlers/domain_events.go`
- `sfx_markt/media/internal/handlers/domain_events_transaction.go`
- `sfx_markt/media/internal/handlers/middleman_image.go`
- `sfx_markt/media/internal/handlers/middleman_media.go`
- `sfx_markt/media/internal/handlers/middleman_video.go`

### Postgres/repository files
- `sfx_markt/media/internal/postgres/import_item_repository.go`
- `sfx_markt/media/internal/postgres/import_session_repository.go`
- `sfx_markt/media/internal/postgres/middleman_image_repository.go`
- `sfx_markt/media/internal/postgres/middleman_media_repository.go`
- `sfx_markt/media/internal/postgres/middleman_video_repository.go`

### SQL migrations
- `sfx_markt/media/migrations/001_create_tables.sql`
- `sfx_markt/media/migrations/002_add_product_id_to_import_items.sql`

### RPC methods from api.proto
- `CreateMedia`
- `UpdateMedia`
- `AddImage`
- `AddVideo`
- `RemoveMedia`
- `RemoveImage`
- `RemoveVideo`
- `GetMedia`
- `GetMediaByItem`
- `GetAllItemImages`
- `GetAllItemVideos`
- `GetAllVideos`
- `GetAllVideosByAuthor`
- `GetAllImagesByAuthor`
- `GetAllMediaVideos`
- `GetAllMediaImages`
- `StartBulkImport`
- `AddImportBatch`
- `GetImportStatus`
- `CancelImport`

### Event/message constants (channel names + event/command keys)
- `MediaAggregateChannel`
- `MediaCreatedEvent`
- `MediaUpdatedEvent`
- `ImageAggregateChannel`
- `ImageAddedEvent`
- `VideoAggregateChannel`
- `VideoAddedEvent`
- `ImageRemovedEvent`
- `MediaRemovedEvent`
- `VideoRemovedEvent`
- `ImporterAggregateChannel`
- `BulkImportStartedEvent`
- `ImportBatchAddedEvent`
- `ImportItemProcessedEvent`
- `ImportItemFailedEvent`
- `BulkImportCompletedEvent`
- `BulkImportCancelledEvent`
- (none found)

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./media/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./media/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile media up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/media && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/media && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./media/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/media/mediapb/api.proto`.
2. Implement command object + handler in `internal/application/commands`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/media/mediapb/api.proto`.
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
   - Use `sfx_markt/media/mediapb/events.proto` for event payload contract updates.
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
   - `cd <workspace-root>/sfx_markt && go test ./media/...`
   - `cd <workspace-root>/sfx_markt && go run ./media/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.
