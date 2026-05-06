# AGENTS Guide: sfx_markt

## Business Purpose
Top-level project root for sfx_markt. This repository follows modular service boundaries (CQRS/event-driven in many services) with shared infra and deployment/tooling directories.

## First Read
### Top-level service/support directories
- `sfx_markt/activity`
- `sfx_markt/assistants`
- `sfx_markt/baskets`
- `sfx_markt/categories`
- `sfx_markt/comments`
- `sfx_markt/core-api-modules-vaveme`
- `sfx_markt/cosec`
- `sfx_markt/docker`
- `sfx_markt/erp`
- `sfx_markt/following`
- `sfx_markt/frontend`
- `sfx_markt/geocoding`
- `sfx_markt/internal`
- `sfx_markt/k8s`
- `sfx_markt/mailer`
- `sfx_markt/managers`
- `sfx_markt/media`
- `sfx_markt/merchant`
- `sfx_markt/messages`
- `sfx_markt/metrics`
- `sfx_markt/migrations`
- `sfx_markt/newsletters`
- `sfx_markt/notifications`
- `sfx_markt/offers`
- `sfx_markt/ordering`
- `sfx_markt/payments`
- `sfx_markt/pim`
- `sfx_markt/postgres`
- `sfx_markt/posts`
- `sfx_markt/products`
- `sfx_markt/prompts`
- `sfx_markt/react-nats-messaging`
- `sfx_markt/reviews`
- `sfx_markt/rust-streams`
- `sfx_markt/sap`
- `sfx_markt/scheduler`
- `sfx_markt/search`
- `sfx_markt/services`
- `sfx_markt/shipping`
- `sfx_markt/streams`
- `sfx_markt/support`
- `sfx_markt/testing`
- `sfx_markt/tickets`
- `sfx_markt/users`
- `sfx_markt/vectors`
- `sfx_markt/wishlists`

### Top-level root files
- `sfx_markt/.gitignore`
- `sfx_markt/CLAUDE.md`
- `sfx_markt/Jenkinsfile`
- `sfx_markt/Makefile`
- `sfx_markt/docker-compose.yaml`
- `sfx_markt/go.mod`
- `sfx_markt/go.sum`

## Project Workflow For Agents
1. Start in service directory that owns the requested business behavior (for example: `users`, `products`, `ordering`, `payments`, `search`).
2. Read that directory's local `AGENTS.md` before modifying code.
3. Keep transport contracts, application layer, handlers, and SQL/repository changes synchronized.
4. Prefer forward-only migrations and additive schema evolution.

## Proto / Command / Query / Event Workflow
- Service-local implementation lives under `sfx_markt/<service>`.
- For command/query/event work, apply the per-service workflow documented in each service `AGENTS.md`.
- Regeneration baseline:
  - `cd <workspace-root>/sfx_markt && go generate ./...`

## SQL / Schema Safety
- Never rewrite historical migrations already applied to environments.
- Add new migrations for every schema/data evolution step.
