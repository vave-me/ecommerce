# AGENTS Guide: offers

## Business Purpose
Owns offer lifecycle and offer negotiation state transitions.

## Service Runtime Shape
- Entry point: `sfx_markt/offers/cmd/service/main.go`
- Module composition/wiring: `sfx_markt/offers/module.go`
- Codegen entrypoint: `sfx_markt/offers/generate.go`
- Application facade: `sfx_markt/offers/internal/application/application.go`
- gRPC server implementation: `sfx_markt/offers/internal/grpc/server.go`
- gRPC transaction wrapper: `sfx_markt/offers/internal/grpc/server_transaction.go`
- API contract (proto): `sfx_markt/offers/offerspb/api.proto`
- Event contract (proto): `sfx_markt/offers/offerspb/events.proto`
- Message contract (proto): `(not found)`

## Concrete Files To Read First
### Proto contracts
- `sfx_markt/offers/offerspb/api.proto`
- `sfx_markt/offers/offerspb/events.proto`

### gRPC transport layer files
- `sfx_markt/offers/internal/grpc/server.go`
- `sfx_markt/offers/internal/grpc/server_transaction.go`

### Application command handlers
- `sfx_markt/offers/internal/application/commands/accept_buy_back_negotiation.go`
- `sfx_markt/offers/internal/application/commands/accept_buy_now_negotiation.go`
- `sfx_markt/offers/internal/application/commands/accept_lease_negotiation.go`
- `sfx_markt/offers/internal/application/commands/accept_offer.go`
- `sfx_markt/offers/internal/application/commands/accept_reservation_negotiation.go`
- `sfx_markt/offers/internal/application/commands/activate_offer.go`
- `sfx_markt/offers/internal/application/commands/cancel_buy_back.go`
- `sfx_markt/offers/internal/application/commands/cancel_buy_now.go`
- `sfx_markt/offers/internal/application/commands/cancel_lease.go`
- `sfx_markt/offers/internal/application/commands/cancel_reservation.go`
- `sfx_markt/offers/internal/application/commands/close_offer.go`
- `sfx_markt/offers/internal/application/commands/confirm_buy_now.go`
- `sfx_markt/offers/internal/application/commands/create_buy_back.go`
- `sfx_markt/offers/internal/application/commands/create_buy_now.go`
- `sfx_markt/offers/internal/application/commands/create_lease.go`
- `sfx_markt/offers/internal/application/commands/create_offer.go`
- `sfx_markt/offers/internal/application/commands/create_reservation.go`
- `sfx_markt/offers/internal/application/commands/decline_buy_back_negotiation.go`
- `sfx_markt/offers/internal/application/commands/decline_buy_now_negotiation.go`
- `sfx_markt/offers/internal/application/commands/decline_lease_negotiation.go`
- `sfx_markt/offers/internal/application/commands/decline_reservation_negotiation.go`
- `sfx_markt/offers/internal/application/commands/default_lease.go`
- `sfx_markt/offers/internal/application/commands/end_lease.go`
- `sfx_markt/offers/internal/application/commands/execute_lease_buyout.go`
- `sfx_markt/offers/internal/application/commands/expire_buy_back.go`
- `sfx_markt/offers/internal/application/commands/expire_reservation.go`
- `sfx_markt/offers/internal/application/commands/make_lease_payment.go`
- `sfx_markt/offers/internal/application/commands/mark_buy_now_paid.go`
- `sfx_markt/offers/internal/application/commands/redeem_buy_back.go`
- `sfx_markt/offers/internal/application/commands/redeem_reservation.go`
- `sfx_markt/offers/internal/application/commands/requeast_lease_negotiation.go`
- `sfx_markt/offers/internal/application/commands/request_buy_back_negotiation.go`
- `sfx_markt/offers/internal/application/commands/request_buy_now_negotiation.go`
- `sfx_markt/offers/internal/application/commands/request_reservation_negotiation.go`
- `sfx_markt/offers/internal/application/commands/start_lease.go`

### Application query handlers
- `sfx_markt/offers/internal/application/queries/get_offer.go`
- `sfx_markt/offers/internal/application/queries/get_offers_by_id.go`
- `sfx_markt/offers/internal/application/queries/get_offers_by_item.go`
- `sfx_markt/offers/internal/application/queries/get_offers_for_buyer.go`
- `sfx_markt/offers/internal/application/queries/get_offers_for_seller.go`
- `sfx_markt/offers/internal/application/queries/get_visited_items_for_buyer.go`
- `sfx_markt/offers/internal/application/queries/get_wishlist_items_for_buyer.go`
- `sfx_markt/offers/internal/application/queries/list_offers.go`

### Domain model files
- `sfx_markt/offers/internal/domain/buy_back.go`
- `sfx_markt/offers/internal/domain/buy_back_events.go`
- `sfx_markt/offers/internal/domain/buy_back_repository.go`
- `sfx_markt/offers/internal/domain/buy_back_snapshots.go`
- `sfx_markt/offers/internal/domain/buy_back_status.go`
- `sfx_markt/offers/internal/domain/buy_now.go`
- `sfx_markt/offers/internal/domain/buy_now_events.go`
- `sfx_markt/offers/internal/domain/buy_now_repository.go`
- `sfx_markt/offers/internal/domain/buy_now_snapshots.go`
- `sfx_markt/offers/internal/domain/buy_now_status.go`
- `sfx_markt/offers/internal/domain/kyc_status.go`
- `sfx_markt/offers/internal/domain/lease.go`
- `sfx_markt/offers/internal/domain/lease_events.go`
- `sfx_markt/offers/internal/domain/lease_repository.go`
- `sfx_markt/offers/internal/domain/lease_snapshots.go`
- `sfx_markt/offers/internal/domain/lease_status.go`
- `sfx_markt/offers/internal/domain/middleman_buy_back_repository.go`
- `sfx_markt/offers/internal/domain/middleman_buy_now_repository.go`
- `sfx_markt/offers/internal/domain/middleman_lease_repository.go`
- `sfx_markt/offers/internal/domain/middleman_repository.go`
- `sfx_markt/offers/internal/domain/middleman_reservation_repository.go`
- `sfx_markt/offers/internal/domain/offer.go`
- `sfx_markt/offers/internal/domain/offer_events.go`
- `sfx_markt/offers/internal/domain/offer_repository.go`
- `sfx_markt/offers/internal/domain/offer_snapshots.go`
- `sfx_markt/offers/internal/domain/offer_status.go`
- `sfx_markt/offers/internal/domain/offer_type.go`
- `sfx_markt/offers/internal/domain/reservation.go`
- `sfx_markt/offers/internal/domain/reservation_events.go`
- `sfx_markt/offers/internal/domain/reservation_repository.go`
- `sfx_markt/offers/internal/domain/reservation_snapshots.go`
- `sfx_markt/offers/internal/domain/reservation_status.go`

### Event/integration/projection handlers
- `sfx_markt/offers/internal/handlers/domain_events.go`
- `sfx_markt/offers/internal/handlers/domain_events_transaction.go`
- `sfx_markt/offers/internal/handlers/middleman.go`

### Postgres/repository files
- `sfx_markt/offers/internal/postgres/middleman_buy_back_repository.go`
- `sfx_markt/offers/internal/postgres/middleman_buy_now_repository.go`
- `sfx_markt/offers/internal/postgres/middleman_lease_repository.go`
- `sfx_markt/offers/internal/postgres/middleman_repository.go`
- `sfx_markt/offers/internal/postgres/middleman_reservation_repository.go`

### SQL migrations
- `sfx_markt/offers/migrations/001_create_tables.sql`

### RPC methods from api.proto
- `CreateOffer`
- `ActivateOffer`
- `CloseOffer`
- `AcceptOffer`
- `GetOffer`
- `ListOffers`
- `CreateBuyNow`
- `ConfirmBuyNow`
- `CancelBuyNow`
- `RequestBuyNowNegotiation`
- `AcceptBuyNowNegotiation`
- `DeclineBuyNowNegotiation`
- `CreateLease`
- `StartLease`
- `MakeLeasePayment`
- `ExecuteLeaseBuyout`
- `EndLease`
- `CancelLease`
- `DefaultLease`
- `RequestLeaseNegotiation`
- `AcceptLeaseNegotiation`
- `DeclineLeaseNegotiation`
- `CreateBuyBack`
- `RedeemBuyBack`
- `ExpireBuyBack`
- `CancelBuyBack`
- `RequestBuyBackNegotiation`
- `AcceptBuyBackNegotiation`
- `DeclineBuyBackNegotiation`
- `CreateReservation`
- `RedeemReservation`
- `ExpireReservation`
- `CancelReservation`
- `RequestReservationNegotiation`
- `AcceptReservationNegotiation`
- `DeclineReservationNegotiation`

### Event/message constants (channel names + event/command keys)
- `OfferAggregateChannel`
- `LeaseAggregateChannel`
- `BuyBackAggregateChannel`
- `BuyNowAggregateChannel`
- `ReservationAggregateChannel`
- `OfferCreatedEvent`
- `OfferActivatedEvent`
- `OfferClosedEvent`
- `OfferAcceptedEvent`
- `BuyNowCreatedEvent`
- `BuyNowConfirmedEvent`
- `LeaseAddedEvent`
- `LeaseCreatedEvent`
- `LeaseDefaultedEvent`
- `LeaseEndedEvent`
- `BuyBackCreatedEvent`
- `BuyBackCanceledEvent`
- `BuyBackExpiredEvent`
- `BuyBackRedeemedEvent`
- `ReservationCreatedEvent`
- `ReservationCanceledEvent`
- `ReservationExpiredEvent`
- `ReservationRedeemedEvent`
- (none found)

## Local Development Workflow (Service)
1. Run tests for this service:
   - `cd <workspace-root>/sfx_markt && go test ./offers/...`
2. Start the service directly:
   - `cd <workspace-root>/sfx_markt && go run ./offers/cmd/service`
3. Optional container run (if profile exists):
   - `cd <workspace-root>/sfx_markt && docker compose --env-file ./docker/.env --profile offers up`

## Regenerate Proto / Codegen Workflow
1. Install toolchain (project root):
   - `cd <workspace-root>/sfx_markt && make install-tools`
2. Regenerate this service:
   - `cd <workspace-root>/sfx_markt/offers && go generate ./...`
3. If only proto files changed:
   - `cd <workspace-root>/sfx_markt/offers && buf generate`
4. Re-run tests:
   - `cd <workspace-root>/sfx_markt && go test ./offers/...`

## Add or Remove a Command (Detailed)
1. Update command request/response messages and RPC (if externally callable) in `sfx_markt/offers/offerspb/api.proto`.
2. Implement command object + handler in `internal/application/commands`.
3. Register handler in `internal/application/application.go` (interface fields + constructor wiring).
4. Map RPC request to command execution in `internal/grpc/server.go`.
5. Add transaction-wrapped counterpart in `internal/grpc/server_transaction.go` when write-path should run in a DB transaction.
6. Update domain aggregate behavior under `internal/domain` so command mutations emit correct domain events.
7. If this command is async/bus-driven, update channel/command constants in generated pb files and handler dispatch code under `internal/handlers`.
8. Validate repositories and migrations are aligned with new write/read needs.

## Add or Remove a Query (Detailed)
1. Add query request/response messages + RPC in `sfx_markt/offers/offerspb/api.proto`.
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
   - Use `sfx_markt/offers/offerspb/events.proto` for event payload contract updates.
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
   - `cd <workspace-root>/sfx_markt && go test ./offers/...`
   - `cd <workspace-root>/sfx_markt && go run ./offers/cmd/service`
5. For breaking schema updates, use staged additive migration sequence.

## Agent Definition of Done
1. Proto contracts and generated artifacts are synchronized.
2. gRPC server, app handlers, transaction wrapper, and domain/repository wiring are consistent.
3. Event producer/consumer mappings and channel constants are consistent.
4. SQL migrations are forward-only and repository SQL is aligned.
5. Service tests pass, or failing tests are explicitly documented.
