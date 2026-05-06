# AGENTS Guide: k8s

## Business Purpose
Kubernetes manifests and deployment configurations.

## Key Files In This Directory
- `sfx_markt/k8s/00-namespaces.yaml`
- `sfx_markt/k8s/02-database/node-port.yaml`
- `sfx_markt/k8s/02-database/postgres.yaml`
- `sfx_markt/k8s/03-messaging/.allowTrafic`
- `sfx_markt/k8s/03-messaging/middleware-ws.yaml`
- `sfx_markt/k8s/03-messaging/nats-ws-certificate.yaml`
- `sfx_markt/k8s/03-messaging/nats-ws-ingress.yaml`
- `sfx_markt/k8s/03-messaging/nats.yaml`
- `sfx_markt/k8s/03-messaging/node-port.yaml`
- `sfx_markt/k8s/04-redis/redis.yaml`
- `sfx_markt/k8s/05-minio/minio-api.yaml`
- `sfx_markt/k8s/05-minio/minio.yaml`
- `sfx_markt/k8s/06-common/config-common.yaml`
- `sfx_markt/k8s/06-common/job-db-init.yaml`
- `sfx_markt/k8s/06-common/secret-common.yaml`
- `sfx_markt/k8s/07-reverse-proxy/nginx-deployment.yaml`
- `sfx_markt/k8s/07-reverse-proxy/nginx-proxy.yaml`
- `sfx_markt/k8s/07-reverse-proxy/reverse-proxy-configmap.yaml`
- `sfx_markt/k8s/07-reverse-proxy/reverse-proxy_service.yaml`
- `sfx_markt/k8s/08-middleware/middleware-cors.yaml`
- `sfx_markt/k8s/08-middleware/middleware-observability.yaml`
- `sfx_markt/k8s/08-middleware/observability-auth.yaml`
- `sfx_markt/k8s/09-services/activity.yaml`
- `sfx_markt/k8s/09-services/assistants.yaml`
- `sfx_markt/k8s/09-services/baskets.yaml`
- `sfx_markt/k8s/09-services/categories.yaml`
- `sfx_markt/k8s/09-services/comments.yaml`
- `sfx_markt/k8s/09-services/cosec.yaml`
- `sfx_markt/k8s/09-services/erp.yaml`
- `sfx_markt/k8s/09-services/following.yaml`
- `sfx_markt/k8s/09-services/geocoding.yaml`
- `sfx_markt/k8s/09-services/mailer.yaml`
- `sfx_markt/k8s/09-services/managers.yaml`
- `sfx_markt/k8s/09-services/media.yaml`
- `sfx_markt/k8s/09-services/merchant.yaml`
- `sfx_markt/k8s/09-services/messages.yaml`
- `sfx_markt/k8s/09-services/metrics.yaml`
- `sfx_markt/k8s/09-services/newsletters.yaml`
- `sfx_markt/k8s/09-services/notifications.yaml`
- `sfx_markt/k8s/09-services/offers.yaml`
- `sfx_markt/k8s/09-services/ordering.yaml`
- `sfx_markt/k8s/09-services/payments.yaml`
- `sfx_markt/k8s/09-services/posts.yaml`
- `sfx_markt/k8s/09-services/products.yaml`
- `sfx_markt/k8s/09-services/reviews.yaml`
- `sfx_markt/k8s/09-services/scheduler.yaml`
- `sfx_markt/k8s/09-services/search.yaml`
- `sfx_markt/k8s/09-services/services.yaml`
- `sfx_markt/k8s/09-services/shipping.yaml`
- `sfx_markt/k8s/09-services/support.yaml`
- `sfx_markt/k8s/09-services/users.yaml`
- `sfx_markt/k8s/09-services/vectors.yaml`
- `sfx_markt/k8s/09-services/wishlists.yaml`
- `sfx_markt/k8s/10-frontend/frontend.yaml`
- `sfx_markt/k8s/11-jenkins/jenkins.yaml`
- `sfx_markt/k8s/11-jenkins/secret.yaml`
- `sfx_markt/k8s/12-nominatim/nominatim.yaml`
- `sfx_markt/k8s/13-registry/k8s-registry-config.yaml`
- `sfx_markt/k8s/14-vectors/qdrant.yaml`
- `sfx_markt/k8s/15-jaeger/jaeger.yaml`
- `sfx_markt/k8s/16-grafana/apply-fixes.sh`
- `sfx_markt/k8s/16-grafana/create-dashboard-configmaps.sh`
- `sfx_markt/k8s/16-grafana/create-dashboards-configmap.sh`
- `sfx_markt/k8s/16-grafana/deploy-production.sh`
- `sfx_markt/k8s/16-grafana/grafana.yaml`
- `sfx_markt/k8s/17-otel/otel-collector-updated.yaml`
- `sfx_markt/k8s/17-otel/otel-collector.yaml`
- `sfx_markt/k8s/18-prometheus/prometheus-scrape-config.yaml`
- `sfx_markt/k8s/18-prometheus/prometheus-updated-config.yaml`
- `sfx_markt/k8s/18-prometheus/prometheus.yaml`
- `sfx_markt/k8s/assistants-grpc-web-solution.yaml`
- `sfx_markt/k8s/deployment-guide.md`

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
