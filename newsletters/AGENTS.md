# AGENTS Guide: newsletters

## Business Purpose
Manages newsletter subscriptions, campaign operations, and delivery state.

## Service Runtime Shape
- Entry point: `sfx_markt/newsletters/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/newsletters/module.go`
- Codegen entrypoint: `sfx_markt/newsletters/generate.go`
- Application facade: `sfx_markt/newsletters/internal/application/application.go`
- gRPC server implementation: `sfx_markt/newsletters/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/newsletters/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/newsletters/newsletterspb/api.proto`
- Event contract (proto): `sfx_markt/newsletters/newsletterspb/events.proto`
- Message contract (proto): `(not found)`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/newsletters/newsletterspb/api.proto`
- `sfx_markt/newsletters/newsletterspb/events.proto`

### gRPC transport layer files
- `sfx_markt/newsletters/internal/grpc/server.go`
- `sfx_markt/newsletters/internal/grpc/server_transaction.go`

### Application command handlers
- `sfx_markt/newsletters/internal/application/commands/accept_offer.go`
- `sfx_markt/newsletters/internal/application/commands/cancel_offer.go`
- `sfx_markt/newsletters/internal/application/commands/create_edition.go`
- `sfx_markt/newsletters/internal/application/commands/create_newsletter.go`
- `sfx_markt/newsletters/internal/application/commands/create_template.go`
- `sfx_markt/newsletters/internal/application/commands/delete_newsletter.go`
- `sfx_markt/newsletters/internal/application/commands/delete_template.go`
- `sfx_markt/newsletters/internal/application/commands/reject_offer.go`
- `sfx_markt/newsletters/internal/application/commands/respond_to_offer.go`
- `sfx_markt/newsletters/internal/application/commands/schedule_edition.go`
- `sfx_markt/newsletters/internal/application/commands/send_edition.go`
- `sfx_markt/newsletters/internal/application/commands/send_offer_to_buyer.go`
- `sfx_markt/newsletters/internal/application/commands/subscribe.go`
- `sfx_markt/newsletters/internal/application/commands/subscribe_newsletter.go`
- `sfx_markt/newsletters/internal/application/commands/unsubscribe.go`
- `sfx_markt/newsletters/internal/application/commands/update_edition.go`
- `sfx_markt/newsletters/internal/application/commands/update_newsletter.go`
- `sfx_markt/newsletters/internal/application/commands/update_subscription.go`
- `sfx_markt/newsletters/internal/application/commands/update_template.go`

### Application query handlers
- `sfx_markt/newsletters/internal/application/queries/get_edition.go`
- `sfx_markt/newsletters/internal/application/queries/get_edition_stats.go`
- `sfx_markt/newsletters/internal/application/queries/get_newsletter.go`
- `sfx_markt/newsletters/internal/application/queries/get_newsletter_stats.go`
- `sfx_markt/newsletters/internal/application/queries/get_offers_by_id.go`
- `sfx_markt/newsletters/internal/application/queries/get_offers_by_item.go`
- `sfx_markt/newsletters/internal/application/queries/get_offers_for_buyer.go`
- `sfx_markt/newsletters/internal/application/queries/get_offers_for_seller.go`
- `sfx_markt/newsletters/internal/application/queries/get_subscription.go`
- `sfx_markt/newsletters/internal/application/queries/get_template.go`
- `sfx_markt/newsletters/internal/application/queries/get_visited_items_for_buyer.go`
- `sfx_markt/newsletters/internal/application/queries/get_wishlist_items_for_buyer.go`
- `sfx_markt/newsletters/internal/application/queries/list_editions.go`
- `sfx_markt/newsletters/internal/application/queries/list_newsletters.go`
- `sfx_markt/newsletters/internal/application/queries/list_subscriptions.go`
- `sfx_markt/newsletters/internal/application/queries/list_templates.go`

### Domain model files
- `sfx_markt/newsletters/internal/domain/catalog_repository.go`
- `sfx_markt/newsletters/internal/domain/edition.go`
- `sfx_markt/newsletters/internal/domain/middleman_repository.go`
- `sfx_markt/newsletters/internal/domain/newsletter.go`
- `sfx_markt/newsletters/internal/domain/newsletter_events.go`
- `sfx_markt/newsletters/internal/domain/newsletter_repository.go`
- `sfx_markt/newsletters/internal/domain/newsletter_snapshots.go`
- `sfx_markt/newsletters/internal/domain/subscription.go`
- `sfx_markt/newsletters/internal/domain/subscription_preferences.go`
- `sfx_markt/newsletters/internal/domain/template.go`

### Event/integration/projection handlers
- `sfx_markt/newsletters/internal/handlers/domain_events.go`
- `sfx_markt/newsletters/internal/handlers/domain_events_transaction.go`

### Postgres/repository files
- `sfx_markt/newsletters/internal/postgres/edition_catalog_repository.go`
- `sfx_markt/newsletters/internal/postgres/middleman_repository.go`
- `sfx_markt/newsletters/internal/postgres/newsletter_catalog_repository.go`
- `sfx_markt/newsletters/internal/postgres/subscription_catalog_repository.go`
- `sfx_markt/newsletters/internal/postgres/template_catalog_repository.go`

### SQL migrations
- `sfx_markt/newsletters/migrations/001_create_tables.sql`

### RPC methods from api.proto
- `CreateNewsletter`
- `UpdateNewsletter`
- `GetNewsletter`
- `ListNewsletters`
- `DeleteNewsletter`
- `Subscribe`
- `Unsubscribe`
- `UpdateSubscription`
- `GetSubscription`
- `ListSubscriptions`
- `CreateEdition`
- `UpdateEdition`
- `GetEdition`
- `ListEditions`
- `ScheduleEdition`
- `SendEdition`
- `PreviewEdition`
- `CreateTemplate`
- `UpdateTemplate`
- `GetTemplate`
- `ListTemplates`
- `DeleteTemplate`
- `GetNewsletterStats`
- `GetEditionStats`

### Event/message constants (channel names + event/command keys)
- `NewsletterAggregateChannel`
- `SubscriptionAggregateChannel`
- `EditionAggregateChannel`
- `TemplateAggregateChannel`
- `NewsletterCreatedEvent`
- `NewsletterUpdatedEvent`
- `NewsletterActivatedEvent`
- `NewsletterDeactivatedEvent`
- `NewsletterDeletedEvent`
- `SubscribedEvent`
- `UnsubscribedEvent`
- `PreferencesUpdatedEvent`
- `SubscriptionPausedEvent`
- `SubscriptionResumedEvent`
- `EditionCreatedEvent`
- `EditionUpdatedEvent`
- `EditionScheduledEvent`
- `EditionSendingEvent`
- `EditionSentEvent`
- `TemplateCreatedEvent`
- `TemplateUpdatedEvent`
- `TemplateDeletedEvent`
- (none found)

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./newsletters/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./newsletters/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile newsletters up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/newsletters && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/newsletters && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./newsletters/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/newsletters/newsletterspb/api.proto`.
2. Implement command object + handler in `internal/application/commands`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/newsletters/newsletterspb/api.proto`.
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
   - Use `sfx_markt/newsletters/newsletterspb/events.proto` for event payload contract updates.
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
   - `cd <workspace-root>/sfx_markt && go test ./newsletters/...`
   - `cd <workspace-root>/sfx_markt && go run ./newsletters/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.
