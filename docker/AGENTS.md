# AGENTS Guide: docker

## Business Purpose
Container build/runtime assets and docker compose wiring.

## Key Files In This Directory
- `sfx_markt/docker/.env`
- `sfx_markt/docker/Dockerfile`
- `sfx_markt/docker/Dockerfile.Geo`
- `sfx_markt/docker/Dockerfile.frontend`
- `sfx_markt/docker/Dockerfile.mailer`
- `sfx_markt/docker/Dockerfile.microservices`
- `sfx_markt/docker/Dockerfile.nginx`
- `sfx_markt/docker/database/001_create_service_dbs.sh`
- `sfx_markt/docker/grafana/BUSINESS_DASHBOARDS_SUMMARY.md`
- `sfx_markt/docker/grafana/PROFESSIONAL_DASHBOARDS_SUMMARY.md`
- `sfx_markt/docker/grafana/grafana.ini`
- `sfx_markt/docker/grafana/provisioning/dashboards/business-kpi-dashboard.json`
- `sfx_markt/docker/grafana/provisioning/dashboards/default.yaml`
- `sfx_markt/docker/grafana/provisioning/dashboards/executive-analytics.json`
- `sfx_markt/docker/grafana/provisioning/dashboards/middleman.json`
- `sfx_markt/docker/grafana/provisioning/dashboards/opentelemetry-collector.json`
- `sfx_markt/docker/grafana/provisioning/dashboards/red-metrics.json`
- `sfx_markt/docker/grafana/provisioning/dashboards/sla-compliance.json`
- `sfx_markt/docker/grafana/provisioning/dashboards/sla-slo-monitoring.json`
- `sfx_markt/docker/grafana/provisioning/dashboards/spanmetrics.json`
- `sfx_markt/docker/grafana/provisioning/dashboards/technical-operations.json`
- `sfx_markt/docker/grafana/provisioning/datasources/default.yaml`
- `sfx_markt/docker/grafana/provisioning/datasources/jaeger.yaml`
- `sfx_markt/docker/jenkins/Dockerfile`
- `sfx_markt/docker/jenkins/casc_configs/jenkins.yaml`
- `sfx_markt/docker/jenkins/plugins.txt`
- `sfx_markt/docker/nats-allow-policy.yaml`
- `sfx_markt/docker/nats-ws-middleware.yaml`
- `sfx_markt/docker/nats/server.conf`
- `sfx_markt/docker/nginx.conf`
- `sfx_markt/docker/nginx_local.conf`
- `sfx_markt/docker/nominatim/germany-latest.osm.pbf`
- `sfx_markt/docker/otel/check_resources.sh`
- `sfx_markt/docker/otel/otel-config.yml`
- `sfx_markt/docker/otel/otel-config_back.yml`
- `sfx_markt/docker/otel/test-trace.sh`
- `sfx_markt/docker/otel/update-k8s-otel.sh`
- `sfx_markt/docker/postgres/postgres-custom.conf`
- `sfx_markt/docker/prometheus/prometheus-config.yml`
- `sfx_markt/docker/redis/redis.conf`
- `sfx_markt/docker/templates/account_locked.html`
- `sfx_markt/docker/templates/account_locked.txt`
- `sfx_markt/docker/templates/confirm_registration.html`
- `sfx_markt/docker/templates/confirm_registration.txt`
- `sfx_markt/docker/templates/de/account_locked.html`
- `sfx_markt/docker/templates/de/account_locked.txt`
- `sfx_markt/docker/templates/de/confirm_registration.html`
- `sfx_markt/docker/templates/de/confirm_registration.txt`
- `sfx_markt/docker/templates/de/forgot_password.html`
- `sfx_markt/docker/templates/de/forgot_password.txt`
- `sfx_markt/docker/templates/de/new_comment.html`
- `sfx_markt/docker/templates/de/new_comment.txt`
- `sfx_markt/docker/templates/de/new_message.html`
- `sfx_markt/docker/templates/de/new_message.txt`
- `sfx_markt/docker/templates/de/newsletter.html`
- `sfx_markt/docker/templates/de/newsletter.txt`
- `sfx_markt/docker/templates/de/welcome_verified.html`
- `sfx_markt/docker/templates/de/welcome_verified.txt`
- `sfx_markt/docker/templates/en/account_locked.html`
- `sfx_markt/docker/templates/en/account_locked.txt`
- `sfx_markt/docker/templates/en/confirm_registration.html`
- `sfx_markt/docker/templates/en/confirm_registration.txt`
- `sfx_markt/docker/templates/en/forgot_password.html`
- `sfx_markt/docker/templates/en/forgot_password.txt`
- `sfx_markt/docker/templates/en/new_comment.html`
- `sfx_markt/docker/templates/en/new_comment.txt`
- `sfx_markt/docker/templates/en/new_message.html`
- `sfx_markt/docker/templates/en/new_message.txt`
- `sfx_markt/docker/templates/en/newsletter.html`
- `sfx_markt/docker/templates/en/newsletter.txt`
- `sfx_markt/docker/templates/en/welcome_verified.html`
- `sfx_markt/docker/templates/en/welcome_verified.txt`
- `sfx_markt/docker/templates/forgot_password.html`
- `sfx_markt/docker/templates/forgot_password.txt`
- `sfx_markt/docker/templates/new_comment.html`
- `sfx_markt/docker/templates/new_comment.txt`
- `sfx_markt/docker/templates/new_message.html`
- `sfx_markt/docker/templates/new_message.txt`
- `sfx_markt/docker/templates/newsletter.html`
- `sfx_markt/docker/templates/newsletter.txt`

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
