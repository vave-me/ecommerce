# AGENTS Guide: erp

## Business Purpose
Coordinates ERP connectors, invoice and return lifecycles, and inventory synchronization with external systems.

## Service Runtime Shape
- Entry point: `sfx_markt/erp/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/erp/module.go`
- Codegen entrypoint: `sfx_markt/erp/generate.go`
- Application facade: `sfx_markt/erp/internal/application/application.go`
- gRPC server implementation: `sfx_markt/erp/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/erp/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/erp/erppb/api.proto`
- Event contract (proto): `sfx_markt/erp/erppb/events.proto`
- Message contract (proto): `(not found)`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/erp/erppb/api.proto`
- `sfx_markt/erp/erppb/events.proto`

### gRPC transport layer files
- `sfx_markt/erp/internal/grpc/product_repository.go`
- `sfx_markt/erp/internal/grpc/server.go`
- `sfx_markt/erp/internal/grpc/server_transaction.go`

### Application command handlers
- `sfx_markt/erp/internal/application/commands/add_connector.go`
- `sfx_markt/erp/internal/application/commands/approve_invoice.go`
- `sfx_markt/erp/internal/application/commands/approve_return.go`
- `sfx_markt/erp/internal/application/commands/complete_return.go`
- `sfx_markt/erp/internal/application/commands/create_inventory_reservation.go`
- `sfx_markt/erp/internal/application/commands/create_invoice.go`
- `sfx_markt/erp/internal/application/commands/create_return.go`
- `sfx_markt/erp/internal/application/commands/fulfill_inventory_reservation.go`
- `sfx_markt/erp/internal/application/commands/helpers.go`
- `sfx_markt/erp/internal/application/commands/process_return_start.go`
- `sfx_markt/erp/internal/application/commands/process_webhook.go`
- `sfx_markt/erp/internal/application/commands/record_invoice_payment.go`
- `sfx_markt/erp/internal/application/commands/register_connector.go`
- `sfx_markt/erp/internal/application/commands/reject_return.go`
- `sfx_markt/erp/internal/application/commands/release_inventory_reservation.go`
- `sfx_markt/erp/internal/application/commands/remove_connector.go`
- `sfx_markt/erp/internal/application/commands/restock_return_items.go`
- `sfx_markt/erp/internal/application/commands/send_invoice.go`
- `sfx_markt/erp/internal/application/commands/send_order.go`
- `sfx_markt/erp/internal/application/commands/sync_customers.go`
- `sfx_markt/erp/internal/application/commands/sync_invoice_to_erp.go`
- `sfx_markt/erp/internal/application/commands/sync_prices.go`
- `sfx_markt/erp/internal/application/commands/sync_products.go`
- `sfx_markt/erp/internal/application/commands/sync_return_to_erp.go`
- `sfx_markt/erp/internal/application/commands/sync_stock.go`
- `sfx_markt/erp/internal/application/commands/toggle_connector.go`
- `sfx_markt/erp/internal/application/commands/transfer_inventory_reservation.go`
- `sfx_markt/erp/internal/application/commands/update_connector.go`
- `sfx_markt/erp/internal/application/commands/update_inventory_reservation.go`
- `sfx_markt/erp/internal/application/commands/void_invoice.go`

### Application query handlers
- `sfx_markt/erp/internal/application/queries/get_connector_status.go`
- `sfx_markt/erp/internal/application/queries/get_sync_history.go`
- `sfx_markt/erp/internal/application/queries/list_connectors.go`

### Domain model files
- `sfx_markt/erp/internal/domain/attributes.go`
- `sfx_markt/erp/internal/domain/connector_repository.go`
- `sfx_markt/erp/internal/domain/entity_type.go`
- `sfx_markt/erp/internal/domain/events.go`
- `sfx_markt/erp/internal/domain/inventory_reservation.go`
- `sfx_markt/erp/internal/domain/inventory_reservation_events.go`
- `sfx_markt/erp/internal/domain/inventory_reservation_snapshots.go`
- `sfx_markt/erp/internal/domain/invoice.go`
- `sfx_markt/erp/internal/domain/invoice_action.go`
- `sfx_markt/erp/internal/domain/invoice_events.go`
- `sfx_markt/erp/internal/domain/invoice_snapshots.go`
- `sfx_markt/erp/internal/domain/invoice_types.go`
- `sfx_markt/erp/internal/domain/options.go`
- `sfx_markt/erp/internal/domain/order.go`
- `sfx_markt/erp/internal/domain/order_repository.go`
- `sfx_markt/erp/internal/domain/payment.go`
- `sfx_markt/erp/internal/domain/payment_repository.go`
- `sfx_markt/erp/internal/domain/product.go`
- `sfx_markt/erp/internal/domain/product_condition.go`
- `sfx_markt/erp/internal/domain/product_status.go`
- `sfx_markt/erp/internal/domain/repositories.go`
- `sfx_markt/erp/internal/domain/return.go`
- `sfx_markt/erp/internal/domain/return_events.go`
- `sfx_markt/erp/internal/domain/return_snapshots.go`
- `sfx_markt/erp/internal/domain/return_types.go`
- `sfx_markt/erp/internal/domain/seller_type.go`
- `sfx_markt/erp/internal/domain/sync_events.go`
- `sfx_markt/erp/internal/domain/sync_models.go`
- `sfx_markt/erp/internal/domain/types.go`
- `sfx_markt/erp/internal/domain/variant.go`

### Event/integration/projection handlers
- `sfx_markt/erp/internal/handlers/commands.go`
- `sfx_markt/erp/internal/handlers/commands_transaction.go`
- `sfx_markt/erp/internal/handlers/domain_events.go`
- `sfx_markt/erp/internal/handlers/domain_events_handlers.go`
- `sfx_markt/erp/internal/handlers/domain_events_transaction.go`
- `sfx_markt/erp/internal/handlers/integration_events.go`
- `sfx_markt/erp/internal/handlers/integration_events_transaction.go`

### Postgres/repository files
- `sfx_markt/erp/internal/postgres/connector_repository.go`
- `sfx_markt/erp/internal/postgres/invoice_sync_repository.go`
- `sfx_markt/erp/internal/postgres/order_sync_repository.go`
- `sfx_markt/erp/internal/postgres/sync_configuration_repository.go`
- `sfx_markt/erp/internal/postgres/sync_log_repository.go`
- `sfx_markt/erp/internal/postgres/sync_status_repository.go`
- `sfx_markt/erp/internal/postgres/webhook_event_repository.go`

### SQL migrations
- `sfx_markt/erp/migrations/001_create_tables.sql`

### RPC methods from api.proto
- `RegisterConnector`
- `AddConnector`
- `UpdateConnector`
- `RemoveConnector`
- `ToggleConnector`
- `SyncProducts`
- `SyncStock`
- `SyncPrices`
- `SyncCustomers`
- `SendOrder`
- `ProcessWebhook`
- `CreateInvoice`
- `ApproveInvoice`
- `SendInvoice`
- `VoidInvoice`
- `RecordInvoicePayment`
- `SyncInvoiceToERP`
- `CreateReturn`
- `ApproveReturn`
- `RejectReturn`
- `ProcessReturnStart`
- `CompleteReturn`
- `RestockReturnItems`
- `SyncReturnToERP`
- `CreateInventoryReservation`
- `ReleaseInventoryReservation`
- `FulfillInventoryReservation`
- `TransferInventoryReservation`
- `GetConnectorStatus`
- `GetSyncHistory`
- `ListConnectors`
- `UpdateInventoryReservation`

### Event/message constants (channel names + event/command keys)
- `ConnectorAggregateChannel`
- `OrderAggregateChannel`
- `InvoiceAggregateChannel`
- `SyncAggregateChannel`
- `WebhookAggregateChannel`
- `ConnectorRegisteredEvent`
- `ConnectorUpdatedEvent`
- `ConnectorDisabledEvent`
- `ConnectorEnabledEvent`
- `ConnectorRemovedEvent`
- `WebhookReceivedEvent`
- `WebhookProcessedEvent`
- `WebhookFailedEvent`
- `OrderSentToERPEvent`
- `OrderSyncedFromERPEvent`
- `OrderSyncFailedEvent`
- `InvoiceCreatedEvent`
- `InvoiceUpdatedEvent`
- `InvoiceApprovedEvent`
- `InvoiceVoidedEvent`
- `InvoiceSentEvent`
- `InvoicePaymentReceivedEvent`
- `ProductsSyncStartedEvent`
- `ProductsSyncCompletedEvent`
- `ProductSyncedEvent`
- `StockSyncStartedEvent`
- `StockSyncCompletedEvent`
- `StockUpdatedEvent`
- `PricesSyncStartedEvent`
- `PricesSyncCompletedEvent`
- `PriceUpdatedEvent`
- `CustomersSyncStartedEvent`
- `CustomersSyncCompletedEvent`
- `CustomerSyncedEvent`
- `InventoryReservedEvent`
- `InventoryReleasedEvent`
- `InventoryConfirmedEvent`
- `ReturnProcessedEvent`
- `ReturnFailedEvent`
- `SyncStatusUpdatedEvent`
- `SyncErrorEvent`
- `CommandChannel`
- `SyncEntityCommand`
- `ProcessERPEventCommand`
- `RetryFailedSyncCommand`
- `UpdateConnectorConfigCommand`
- (none found)

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./erp/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./erp/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile erp up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/erp && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/erp && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./erp/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/erp/erppb/api.proto`.
2. Implement command object + handler in `internal/application/commands`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/erp/erppb/api.proto`.
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
   - Use `sfx_markt/erp/erppb/events.proto` for event payload contract updates.
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
   - `cd <workspace-root>/sfx_markt && go test ./erp/...`
   - `cd <workspace-root>/sfx_markt && go run ./erp/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.
