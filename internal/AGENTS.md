# AGENTS Guide: internal

## Business Purpose
Shared platform libraries and infrastructure primitives used by services.

## Key Files In This Directory
- `sfx_markt/internal/ai/README_INTERFACES.md`
- `sfx_markt/internal/ai/anthropic_client.go`
- `sfx_markt/internal/ai/client_factory.go`
- `sfx_markt/internal/ai/contstats.go`
- `sfx_markt/internal/ai/deepseek_client.go`
- `sfx_markt/internal/ai/interfaces.go`
- `sfx_markt/internal/ai/openai_client.go`
- `sfx_markt/internal/ai/reasoning_client.go`
- `sfx_markt/internal/am/buf.gen.yaml`
- `sfx_markt/internal/am/buf.yaml`
- `sfx_markt/internal/am/command_messages.go`
- `sfx_markt/internal/am/event_messages.go`
- `sfx_markt/internal/am/fake_event_publisher.go`
- `sfx_markt/internal/am/generate.go`
- `sfx_markt/internal/am/message.go`
- `sfx_markt/internal/am/message_types.pb.go`
- `sfx_markt/internal/am/message_types.proto`
- `sfx_markt/internal/am/middleware.go`
- `sfx_markt/internal/am/mock_command_publisher.go`
- `sfx_markt/internal/am/mock_event_publisher.go`
- `sfx_markt/internal/am/mock_message_handler.go`
- `sfx_markt/internal/am/mock_message_publisher.go`
- `sfx_markt/internal/am/mock_message_subscriber.go`
- `sfx_markt/internal/am/mock_reply_publisher.go`
- `sfx_markt/internal/am/mock_websocket_publisher.go`
- `sfx_markt/internal/am/reply_messages.go`
- `sfx_markt/internal/am/subscriber_config.go`
- `sfx_markt/internal/am/subscription.go`
- `sfx_markt/internal/am/websocket_messages.go`
- `sfx_markt/internal/amotel/extractor.go`
- `sfx_markt/internal/amotel/injector.go`
- `sfx_markt/internal/amotel/metadata_carrier.go`
- `sfx_markt/internal/amotel/trace.go`
- `sfx_markt/internal/amprom/received.go`
- `sfx_markt/internal/amprom/sent.go`
- `sfx_markt/internal/auth/auth.go`
- `sfx_markt/internal/auth/sso/manager.go`
- `sfx_markt/internal/auth/sso/oidc_provider.go`
- `sfx_markt/internal/auth/sso/provider.go`
- `sfx_markt/internal/config/activity_config.go`
- `sfx_markt/internal/config/assistants_config.go`
- `sfx_markt/internal/config/comments_config.go`
- `sfx_markt/internal/config/config.go`
- `sfx_markt/internal/config/geocoding_config.go`
- `sfx_markt/internal/config/mailer_config.go`
- `sfx_markt/internal/config/managers_config.go`
- `sfx_markt/internal/config/media_config.go`
- `sfx_markt/internal/config/merchant_config.go`
- `sfx_markt/internal/config/messenger_config.go`
- `sfx_markt/internal/config/metrics_config.go`
- `sfx_markt/internal/config/payment_config.go`
- `sfx_markt/internal/config/redis_config.go`
- `sfx_markt/internal/config/scheduler_config.go`
- `sfx_markt/internal/config/search_config.go`
- `sfx_markt/internal/config/shipping_config.go`
- `sfx_markt/internal/config/users_config.go`
- `sfx_markt/internal/config/vectors_config.go`
- `sfx_markt/internal/config/webhook_config.go`
- `sfx_markt/internal/ddd/aggregate.go`
- `sfx_markt/internal/ddd/aggregate_build_options.go`
- `sfx_markt/internal/ddd/command.go`
- `sfx_markt/internal/ddd/entity.go`
- `sfx_markt/internal/ddd/entity_build_options.go`
- `sfx_markt/internal/ddd/event.go`
- `sfx_markt/internal/ddd/event_dispatcher.go`
- `sfx_markt/internal/ddd/generate.go`
- `sfx_markt/internal/ddd/metadata.go`
- `sfx_markt/internal/ddd/mock_aggregate.go`
- `sfx_markt/internal/ddd/mock_command_handler.go`
- `sfx_markt/internal/ddd/mock_entity.go`
- `sfx_markt/internal/ddd/mock_event_handler.go`
- `sfx_markt/internal/ddd/mock_event_publisher.go`
- `sfx_markt/internal/ddd/mock_event_subscriber.go`
- `sfx_markt/internal/ddd/mock_reply_handler.go`
- `sfx_markt/internal/ddd/mock_websocket_handler.go`
- `sfx_markt/internal/ddd/mock_websocket_publisher.go`
- `sfx_markt/internal/ddd/mock_websocket_subscriber.go`
- `sfx_markt/internal/ddd/reply.go`
- `sfx_markt/internal/ddd/websocket.go`
- `sfx_markt/internal/ddd/websocket_dispatcher.go`

## How To Work In This Directory
1. Keep changes aligned with the owning service/business module contracts.
2. Do not edit generated/build/vendor artifacts directly.
3. Validate with targeted build/test/deploy commands relevant to affected modules.

## Relationship To Commands / Queries / Proto / Events
- This directory is support/infrastructure/integration-oriented unless it contains a dedicated service module pattern.
- For adding/removing commands, queries, RPCs, domain events, and SQL read models, edit the owning service directory (for example `sfx_markt/users`, `sfx_markt/products`, `sfx_markt/ordering`, `sfx_markt/payments`, `sfx_markt/search`).
- Regeneration baseline when shared contracts change:
  - `cd <workspace-root>/sfx_markt && go generate ./...`

## SQL / Data Safety
- Prefer additive, forward-only schema/data changes.
- Never rewrite historical migrations or previously applied seed files.
