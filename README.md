# SFX Markt Platform

## What This Project Is
`sfx_markt` is a marketplace microservices distribution aligned with SFX Markt deployment needs.

It includes core marketplace services plus enterprise-oriented modules like:
`erp`, `merchant`, `managers`, `streams`, and `tickets`.

## Support and Contributions
This project is free to use and maintained through community collaboration and commercial engagements handled outside repository docs.

Donations and sponsorships fund infrastructure costs, maintenance time, release hardening, and long-term documentation work.

Donation channels:
- Bank transfer (SEPA/SWIFT):
```text
Recipient: vaveme
IBAN: DE04 1001 0178 8926 8156 27
BIC/SWIFT: REVODEB2
Bank name and address: Revolut Bank UAB, Zweigniederlassung Deutschland
FORA Linden Palais, Unter den Linden 40, 10117, Berlin, Germany
Correspondent bank BIC: CHASDEFX
```
- Bitcoin: `bitcoin:bc1qkwyadqjxdxpr5r27qj7epl7xqgp8xsdd20lece`
- PayPal: `@szymongol`

Support channels:
- `GitHub Issues`: bug reports and feature requests.
- `GitHub Discussions`: architecture questions, rollout guidance, and sponsorship/support coordination.
- `Consulting contact`: `simongol@proton.me`
- `Security policy channel`: private vulnerability reporting.

- Use the issue tracker for bug reports and feature requests.
- Use pull requests for code contributions and fixes.
- For commercial support, contact maintainers through repository discussions or organization channels.
- Keep personal payment details, private emails, and wallet/account information out of committed documentation.

## Architecture Summary
- Go microservices with modular composition (`module.go`, `cmd/service/main.go`)
- CQRS + event-driven communication
- gRPC/protobuf contracts across service boundaries
- Postgres read/write models with migration-first workflow
- NATS + Redis + Qdrant + MinIO + observability stack

Use service-level `AGENTS.md` for exact implementation workflows.

## Documentation Quick Access
Swagger/OpenAPI docs are available without running application containers.

Documentation locations:
- Proto contracts: `<service>/*pb/api.proto` (plus `events.proto`, `messages.proto` where used)
- Generated OpenAPI spec: `<service>/internal/rest/api.swagger.json`

Offline docs workflow:
1. List all generated specs:
```bash
find . -type f -path "*/internal/rest/api.swagger.json" | sort
```
2. Install Swagger CLI (once):
```bash
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```
3. Serve one spec locally:
```bash
SPEC=./streams/internal/rest/api.swagger.json
"$(go env GOPATH)/bin/swagger" serve --no-open "$SPEC"
```
4. Open the URL printed in terminal (default: `http://localhost:8080/docs`).

Regenerate docs after API changes:
```bash
make generate
```

Validation smoke test (auto-stops after 5 seconds):
```bash
timeout 5s "$(go env GOPATH)/bin/swagger" serve --no-open --port 9214 ./streams/internal/rest/api.swagger.json
```

Troubleshooting:
- `swagger: command not found`: re-run the install command above.
- Port already in use: add `--port 8081` (or another free port).
- Missing spec file: run `make generate` and list files again.

## Repository Layout
Service modules include:
`activity`, `assistants`, `baskets`, `categories`, `comments`, `erp`, `following`, `geocoding`, `mailer`, `managers`, `media`, `merchant`, `messages`, `metrics`, `newsletters`, `notifications`, `offers`, `ordering`, `payments`, `pim`, `posts`, `products`, `reviews`, `scheduler`, `search`, `services`, `shipping`, `streams`, `support`, `tickets`, `users`, `vectors`, `wishlists`

Support directories include:
`internal`, `docker`, `k8s`, `postgres`, `migrations`, `frontend`, `rust-streams`, `react-nats-messaging`.

## Quick Start
1. Install tools:
```bash
make install-tools
```
2. Generate code:
```bash
make generate
```
3. Run tests:
```bash
go test ./...
```
4. Start main stack:
```bash
docker compose --env-file ./docker/.env --profile microservices up
```
5. Stop stack:
```bash
docker compose --env-file ./docker/.env --profile microservices down
```

## Command Reference

### Backend lifecycle
```bash
make install-tools
make generate
go test ./...
```

### Docker runtime
```bash
cp docker/.env.example docker/.env
docker compose --env-file ./docker/.env --profile microservices up -d
docker compose --env-file ./docker/.env --profile microservices down
docker compose --env-file ./docker/.env --profile ci up -d --force-recreate
docker compose --env-file ./docker/.env logs -f nominatim
```

### Frequently used Make targets
```bash
make run-micro
make down-micro
make run-ci
make down-ci
make build-micro
make build-websocket
make build-merchant
make build-erp
make build-users
make build-users-no-cache
```

### Frontend commands
```bash
cd frontend
cp .env.example .env.local
npm install
npm run dev
npm run lint
npm run build
npm run test -- --runInBand
```

### Pre-push checks
```bash
go test ./...
cd frontend && npm run lint && npm run build && npm run test -- --runInBand
git status
```

## Required API and Infra Setup
Before end-user testing or publish, connect required providers and infra:

1. Create local runtime env from template:
```bash
cp docker/.env.example docker/.env
```
2. Configure `docker/.env`:
   - AI: `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, `DEEPSEEK_API_KEY`, `GOOGLE_AI_API_KEY`
   - Geocoding: `GEOCODING_API_KEY`, `NOMINATIM_PASSWORD`, `PBF_URL`, `REPLICATION_URL`
   - Google OAuth: `WEB_GOOGLE_OAUTH_CLIENT_ID`, `MOBILE_GOOGLE_OAUTH_CLIENT_ID`, `NEXT_PUBLIC_GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`
   - Storage: `MINIO_ROOT_USER`, `MINIO_ROOT_PASSWORD`, `MINIO_BUCKET`, `MINIO_ACCESS_KEY_ID`, `MINIO_SECRET_ACCESS_KEY`
   - Commerce: `STRIPE_KEY`, `STRIPE_SECRET`, `STRIPE_WEBHOOK_SECRET`, `DHL_API_KEY`
3. Download required map file for Nominatim:
```bash
mkdir -p docker/nominatim
curl -L -o docker/nominatim/germany-latest.osm.pbf \
  https://download.geofabrik.de/europe/germany-latest.osm.pbf
```
4. Set `PBF_URL=file:///data/germany-latest.osm.pbf`.
5. Start profile `microservices` and monitor import:
```bash
docker compose logs -f nominatim
```
6. Open `http://localhost:9099` and ensure MinIO bucket `MINIO_BUCKET` exists.

Full cross-project checklist:
- `SETUP_AND_EXTERNAL_INTEGRATIONS.md`

## Docker Compose Analysis
Compose file: `docker-compose.yaml`

Profiles available:
`activity`, `assistants`, `baskets`, `categories`, `ci`, `comments`, `erp`, `following`, `geocoding`, `mailer`, `managers`, `media`, `merchant`, `messages`, `metrics`, `microservices`, `newsletters`, `notifications`, `offers`, `ordering`, `payments`, `posts`, `products`, `reviews`, `search`, `services`, `shipping`, `support`, `users`, `vectors`, `websocket`, `wishlists`

Compose service groups:
- Domain services for marketplace and enterprise modules
- Infra: `postgres`, `nats`, `redis`, `qdrant`, `minio`, `nominatim`, `registry`
- Observability: `collector`, `prometheus`, `grafana`, `jaeger`
- Edge: `reverse-proxy`

Useful compose commands:
```bash
# Full runtime

docker compose --env-file ./docker/.env --profile microservices up

# CI profile

docker compose --env-file ./docker/.env --profile ci up -d --force-recreate
```

## Makefile Analysis
Primary targets:
- Setup/codegen: `install-tools`, `generate`
- Runtime: `run-micro`, `run-streams`, `run-ci`, `down-micro`, `down-ci`
- Build bundles: `build-micro`, `build-websocket`
- Service image builds: `build-<service>`, `build-<service>-no-cache`
- Enterprise builds: `build-erp`, `build-merchant`

Important notes:
- `run-wishlists` uses `wishilist` (typo); compose profile is `wishlists`.
- `run-frontend` exists in Makefile, but compose currently has no `frontend` profile entry.
- `run-streams` exists in Makefile, but compose currently has no `streams` profile entry.

Frontend local run:
```bash
cd frontend
npm install
npm run dev
```

## Docker Directory Analysis
`docker/` contains:
- Service image definitions: `Dockerfile*` (`microservices`, `frontend`, `mailer`, `Geo`, `nginx`)
- Runtime configs: `docker/.env`, NATS/Redis/Postgres config, nginx config
- Observability configs: OTEL/Prometheus/Grafana
- DB init scripts and email templates

## Security and Operations Notes
- Keep runtime secrets in secured environment management, not plaintext repo files.
- Validate registry/tag settings before production pushes.
- Preserve forward-only migration discipline.

## Development Workflow (Feature Changes)
1. Update protobuf contracts.
2. Implement command/query logic.
3. Wire gRPC handlers and tx wrappers.
4. Update domain/integration event handlers.
5. Add migrations and repository updates.
6. Regenerate and run tests.

Reference details: `AGENTS.md` files.

## Repository Governance Files
- `LICENSE`
- `CONTRIBUTING.md`
- `SECURITY.md`
- `CODE_OF_CONDUCT.md`
- `STANDALONE_REPOSITORY_CHECKLIST.md`
