# AGENTS Guide: react-nats-messaging

## Business Purpose
React + NATS messaging integration workspace.

## Key Files In This Directory
- `sfx_markt/react-nats-messaging/.eslintrc.js`
- `sfx_markt/react-nats-messaging/.gitignore`
- `sfx_markt/react-nats-messaging/README.md`
- `sfx_markt/react-nats-messaging/examples/chat-app/ChatExample.tsx`
- `sfx_markt/react-nats-messaging/examples/comments-system/CommentsExample.tsx`
- `sfx_markt/react-nats-messaging/examples/complete-integration/App.tsx`
- `sfx_markt/react-nats-messaging/jest.config.js`
- `sfx_markt/react-nats-messaging/package-lock.json`
- `sfx_markt/react-nats-messaging/package.json`
- `sfx_markt/react-nats-messaging/rollup.config.mjs`
- `sfx_markt/react-nats-messaging/src/api/comments.ts`
- `sfx_markt/react-nats-messaging/src/api/messaging.ts`
- `sfx_markt/react-nats-messaging/src/api/support.ts`
- `sfx_markt/react-nats-messaging/src/components/ConnectionIndicator.tsx`
- `sfx_markt/react-nats-messaging/src/components/MessageDebugger.tsx`
- `sfx_markt/react-nats-messaging/src/core/MessageEncoder.ts`
- `sfx_markt/react-nats-messaging/src/core/NatsConnection.ts`
- `sfx_markt/react-nats-messaging/src/core/NatsProvider.tsx`
- `sfx_markt/react-nats-messaging/src/core/types.ts`
- `sfx_markt/react-nats-messaging/src/generated_proto/comments_api_pb.d.ts`
- `sfx_markt/react-nats-messaging/src/generated_proto/comments_api_pb.js`
- `sfx_markt/react-nats-messaging/src/generated_proto/message_api_pb.d.ts`
- `sfx_markt/react-nats-messaging/src/generated_proto/message_api_pb.js`
- `sfx_markt/react-nats-messaging/src/generated_proto/message_types_pb.d.ts`
- `sfx_markt/react-nats-messaging/src/generated_proto/message_types_pb.js`
- `sfx_markt/react-nats-messaging/src/generated_proto/messages_api_events_pb.d.ts`
- `sfx_markt/react-nats-messaging/src/generated_proto/messages_api_events_pb.js`
- `sfx_markt/react-nats-messaging/src/generated_proto/messages_service_api_pb.d.ts`
- `sfx_markt/react-nats-messaging/src/generated_proto/messages_service_api_pb.js`
- `sfx_markt/react-nats-messaging/src/hooks/__tests__/useNats.test.tsx`
- `sfx_markt/react-nats-messaging/src/hooks/useChatHistory.ts`
- `sfx_markt/react-nats-messaging/src/hooks/useCommentsActions.ts`
- `sfx_markt/react-nats-messaging/src/hooks/useConnectionStatus.ts`
- `sfx_markt/react-nats-messaging/src/hooks/useMessageHandler.ts`
- `sfx_markt/react-nats-messaging/src/hooks/useNats.ts`
- `sfx_markt/react-nats-messaging/src/hooks/usePublish.ts`
- `sfx_markt/react-nats-messaging/src/hooks/useSubscription.ts`
- `sfx_markt/react-nats-messaging/src/index.ts`
- `sfx_markt/react-nats-messaging/src/proto/comments_api.proto`
- `sfx_markt/react-nats-messaging/src/proto/message_api.proto`
- `sfx_markt/react-nats-messaging/src/proto/message_types.proto`
- `sfx_markt/react-nats-messaging/src/proto/messages_api_events.proto`
- `sfx_markt/react-nats-messaging/src/proto/messages_service_api.proto`
- `sfx_markt/react-nats-messaging/src/proto/types.ts`
- `sfx_markt/react-nats-messaging/src/setupTests.ts`
- `sfx_markt/react-nats-messaging/src/utils/deduplication.ts`
- `sfx_markt/react-nats-messaging/src/utils/retry.ts`
- `sfx_markt/react-nats-messaging/tsconfig.json`

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
