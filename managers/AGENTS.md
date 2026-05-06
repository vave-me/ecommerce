# AGENTS Guide: managers

## Business Purpose
Runs manager agents, conversations, and assistant workflows across product modules.

## Service Runtime Shape
- Entry point: `sfx_markt/managers/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/managers/module.go`
- Codegen entrypoint: `sfx_markt/managers/generate.go`
- Application facade: `sfx_markt/managers/internal/application/application.go`
- gRPC server implementation: `sfx_markt/managers/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/managers/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/managers/managerspb/api.proto`
- Event contract (proto): `sfx_markt/managers/managerspb/events.proto`
- Message contract (proto): `(not found)`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/managers/managerspb/api.proto`
- `sfx_markt/managers/managerspb/common_manager.proto`
- `sfx_markt/managers/managerspb/events.proto`

### gRPC transport layer files
- `sfx_markt/managers/internal/grpc/activity_repository.go`
- `sfx_markt/managers/internal/grpc/basket_repository.go`
- `sfx_markt/managers/internal/grpc/category_repository.go`
- `sfx_markt/managers/internal/grpc/comment_repository.go`
- `sfx_markt/managers/internal/grpc/enhanced_repository_adapter.go`
- `sfx_markt/managers/internal/grpc/error_interceptor.go`
- `sfx_markt/managers/internal/grpc/following_repository.go`
- `sfx_markt/managers/internal/grpc/geocoding_repository.go`
- `sfx_markt/managers/internal/grpc/mailer_repository.go`
- `sfx_markt/managers/internal/grpc/media_repository.go`
- `sfx_markt/managers/internal/grpc/messages_repository.go`
- `sfx_markt/managers/internal/grpc/metric_repository.go`
- `sfx_markt/managers/internal/grpc/newsletter_repository.go`
- `sfx_markt/managers/internal/grpc/notification_repository.go`
- `sfx_markt/managers/internal/grpc/offer_repository.go`
- `sfx_markt/managers/internal/grpc/order_repository.go`
- `sfx_markt/managers/internal/grpc/payment_repository.go`
- `sfx_markt/managers/internal/grpc/post_repository.go`
- `sfx_markt/managers/internal/grpc/product_repository.go`
- `sfx_markt/managers/internal/grpc/review_repository.go`
- `sfx_markt/managers/internal/grpc/scheduler_repository.go`
- `sfx_markt/managers/internal/grpc/server.go`
- `sfx_markt/managers/internal/grpc/server_transaction.go`
- `sfx_markt/managers/internal/grpc/service_repository.go`
- `sfx_markt/managers/internal/grpc/shipping_repository.go`
- `sfx_markt/managers/internal/grpc/support_repository.go`
- `sfx_markt/managers/internal/grpc/user_repository.go`
- `sfx_markt/managers/internal/grpc/variant_repository.go`
- `sfx_markt/managers/internal/grpc/vector_repository.go`
- `sfx_markt/managers/internal/grpc/wishlist_repository.go`

### Application command handlers
- `sfx_markt/managers/internal/application/commands/activate_manager.go`
- `sfx_markt/managers/internal/application/commands/add_manager_to_conversation.go`
- `sfx_markt/managers/internal/application/commands/chat_with_conversation.go`
- `sfx_markt/managers/internal/application/commands/create_admin_manager.go`
- `sfx_markt/managers/internal/application/commands/create_business_manager.go`
- `sfx_markt/managers/internal/application/commands/create_conversation.go`
- `sfx_markt/managers/internal/application/commands/create_manager.go`
- `sfx_markt/managers/internal/application/commands/create_scheduler_manager.go`
- `sfx_markt/managers/internal/application/commands/create_support_manager.go`
- `sfx_markt/managers/internal/application/commands/create_user_manager.go`
- `sfx_markt/managers/internal/application/commands/deactivate_manager.go`
- `sfx_markt/managers/internal/application/commands/delete_conversation.go`
- `sfx_markt/managers/internal/application/commands/ensure_consciousness_manager.go`
- `sfx_markt/managers/internal/application/commands/ensure_scheduler_manager.go`
- `sfx_markt/managers/internal/application/commands/process_document_input.go`
- `sfx_markt/managers/internal/application/commands/process_image_input.go`
- `sfx_markt/managers/internal/application/commands/process_speech_input.go`
- `sfx_markt/managers/internal/application/commands/process_user_input.go`
- `sfx_markt/managers/internal/application/commands/update_conversation.go`
- `sfx_markt/managers/internal/application/commands/update_conversation_context.go`
- `sfx_markt/managers/internal/application/commands/update_manager_configuration.go`
- `sfx_markt/managers/internal/application/commands/validation.go`

### Application query handlers
- `sfx_markt/managers/internal/application/queries/get_conversation.go`
- `sfx_markt/managers/internal/application/queries/get_conversation_messages.go`
- `sfx_markt/managers/internal/application/queries/get_conversation_stats.go`
- `sfx_markt/managers/internal/application/queries/get_manager.go`
- `sfx_markt/managers/internal/application/queries/get_managers.go`
- `sfx_markt/managers/internal/application/queries/get_or_create_user_manager.go`
- `sfx_markt/managers/internal/application/queries/get_user_conversations.go`

### Domain model files
- `sfx_markt/managers/internal/domain/activity_repository.go`
- `sfx_markt/managers/internal/domain/analysis_types.go`
- `sfx_markt/managers/internal/domain/basket_repository.go`
- `sfx_markt/managers/internal/domain/catalog_repository.go`
- `sfx_markt/managers/internal/domain/category_repository.go`
- `sfx_markt/managers/internal/domain/comment_repository.go`
- `sfx_markt/managers/internal/domain/conversation.go`
- `sfx_markt/managers/internal/domain/conversation_message.go`
- `sfx_markt/managers/internal/domain/conversation_repository.go`
- `sfx_markt/managers/internal/domain/conversation_snapshots.go`
- `sfx_markt/managers/internal/domain/conversations_events.go`
- `sfx_markt/managers/internal/domain/following_repository.go`
- `sfx_markt/managers/internal/domain/geocoding_repository.go`
- `sfx_markt/managers/internal/domain/interfaces.go`
- `sfx_markt/managers/internal/domain/item_metric_repository.go`
- `sfx_markt/managers/internal/domain/llm_journal.go`
- `sfx_markt/managers/internal/domain/mailer_repository.go`
- `sfx_markt/managers/internal/domain/manager.go`
- `sfx_markt/managers/internal/domain/manager_events.go`
- `sfx_markt/managers/internal/domain/manager_repository.go`
- `sfx_markt/managers/internal/domain/manager_snapshots.go`
- `sfx_markt/managers/internal/domain/media_repository.go`
- `sfx_markt/managers/internal/domain/message_role.go`
- `sfx_markt/managers/internal/domain/messages_repository.go`
- `sfx_markt/managers/internal/domain/metric_repository.go`
- `sfx_markt/managers/internal/domain/natural_repository_interfaces.go`
- `sfx_markt/managers/internal/domain/newsletter_repository.go`
- `sfx_markt/managers/internal/domain/notification_repository.go`
- `sfx_markt/managers/internal/domain/offer_repository.go`
- `sfx_markt/managers/internal/domain/order_repository.go`
- `sfx_markt/managers/internal/domain/payment_repository.go`
- `sfx_markt/managers/internal/domain/post_repository.go`
- `sfx_markt/managers/internal/domain/product_repository.go`
- `sfx_markt/managers/internal/domain/read_conversation_repository.go`
- `sfx_markt/managers/internal/domain/read_messages_repository.go`
- `sfx_markt/managers/internal/domain/review_repository.go`
- `sfx_markt/managers/internal/domain/scheduler_repository.go`
- `sfx_markt/managers/internal/domain/service_repository.go`
- `sfx_markt/managers/internal/domain/shipping_repository.go`
- `sfx_markt/managers/internal/domain/support_repository.go`
- `sfx_markt/managers/internal/domain/user_repository.go`
- `sfx_markt/managers/internal/domain/variant_repository.go`
- `sfx_markt/managers/internal/domain/vector_repository.go`
- `sfx_markt/managers/internal/domain/wishlist_repository.go`

### Event/integration/projection handlers
- `sfx_markt/managers/internal/handlers/domain_events.go`
- `sfx_markt/managers/internal/handlers/domain_events_transaction.go`
- `sfx_markt/managers/internal/handlers/integration_events.go`
- `sfx_markt/managers/internal/handlers/integration_events_transaction.go`
- `sfx_markt/managers/internal/handlers/manager_catalog.go`
- `sfx_markt/managers/internal/handlers/read_conversation.go`
- `sfx_markt/managers/internal/handlers/read_messages.go`

### Postgres/repository files
- `sfx_markt/managers/internal/postgres/catalog_repository.go`
- `sfx_markt/managers/internal/postgres/llm_journal_repository.go`
- `sfx_markt/managers/internal/postgres/read_conversation_repository.go`
- `sfx_markt/managers/internal/postgres/read_messages_repository.go`

### SQL migrations
- `sfx_markt/managers/migrations/001_create_tables.sql`

### RPC methods from api.proto
- `CreateManager`
- `ActivateManager`
- `DeactivateManager`
- `UpdateManagerConfiguration`
- `GetManager`
- `GetManagers`
- `ProcessUserInput`
- `ProcessSpeechInput`
- `ProcessImageInput`
- `ProcessDocumentInput`
- `ProcessManagerRequest`
- `CreateConversation`
- `GetConversation`
- `GetUserConversations`
- `GetConversationMessages`
- `GetConversationStats`
- `AddMessageToConversation`
- `UpdateConversation`
- `UpdateConversationContext`
- `DeleteConversation`
- `ArchiveConversation`
- `ChatWithConversation`

### Event/message constants (channel names + event/command keys)
- `ManagerAggregateChannel`
- `ConversationAggregateChannel`
- `ManagerAuthorizedEvent`
- `ManagerCreatedEvent`
- `ManagerActivatedEvent`
- `ManagerDeactivatedEvent`
- `ManagerConfigurationUpdatedEvent`
- `ManagerRequestProcessedEvent`
- `ManagerParticipatingToggledEvent`
- `ManagerRenamedEvent`
- `ManagerLoggedInEvent`
- `ManagerLoggedOutEvent`
- `ConversationCreatedEvent`
- `MessageAddedEvent`
- `ConversationContextUpdatedEvent`
- `ConversationArchivedEvent`
- (none found)

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./managers/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./managers/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile managers up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/managers && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/managers && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./managers/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/managers/managerspb/api.proto`.
2. Implement command object + handler in `internal/application/commands`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/managers/managerspb/api.proto`.
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
   - Use `sfx_markt/managers/managerspb/events.proto` for event payload contract updates.
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
   - `cd <workspace-root>/sfx_markt && go test ./managers/...`
   - `cd <workspace-root>/sfx_markt && go run ./managers/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.
