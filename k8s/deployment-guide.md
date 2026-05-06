# Deployment Guide for Integrated Observability Stack

This guide outlines the deployment steps for the integrated Docker configurations into your Kubernetes cluster.

## Prerequisites

Ensure the observability namespace and auth middleware are already deployed:
```bash
kubectl apply -f 00-namespaces.yaml
kubectl apply -f 08-middleware/observability-auth.yaml
```

## Deployment Order

### 1. Update Prometheus Configuration

Apply the updated Prometheus configuration with comprehensive service scraping:

```bash
# Replace the existing Prometheus ConfigMap
kubectl apply -f 18-prometheus/prometheus-updated-config.yaml

# Restart Prometheus to pick up new configuration
kubectl rollout restart statefulset/prometheus -n observability
```

### 2. Deploy Grafana with Dashboards

Apply the Grafana dashboard ConfigMap and updated deployment:

```bash
# Apply the dashboard ConfigMap (contains 9 dashboards from Docker)
kubectl apply -f 16-grafana/grafana-dashboards-configmap.yaml

# Apply the updated Grafana deployment
kubectl apply -f 16-grafana/grafana-updated.yaml

# Ensure subdomain ingress is applied
kubectl apply -f 16-grafana/grafana-subdomain-ingress.yaml

# Restart Grafana to ensure all configs are loaded
kubectl rollout restart deployment/grafana -n observability
```

### 3. Update OpenTelemetry Collector

Apply the enhanced OTel Collector configuration:

```bash
# Apply the updated OTel configuration
kubectl apply -f 17-otel/otel-collector-updated.yaml

# The deployment will automatically restart due to ConfigMap changes
```

## Verification Steps

### 1. Check Pod Status
```bash
kubectl get pods -n observability
```

All pods should be in Running state.

### 2. Verify Prometheus Targets
Access Prometheus at https://prometheus.sfx-markt.de (use credentials: admin / ObservabilityP@ss2024!)

Navigate to Status → Targets to verify:
- All service endpoints are being scraped
- OTel Collector metrics are visible
- Infrastructure components are healthy

### 3. Verify Grafana Dashboards
Access Grafana at https://grafana.sfx-markt.de

Check:
- All 9 dashboards are imported under "Docker Import" folder
- Data sources (Prometheus and Jaeger) are connected
- Metrics are flowing from Prometheus

### 4. Verify OTel Collector
```bash
# Check collector logs
kubectl logs -n observability -l app=otel-collector --tail=50

# Access zpages for debugging
kubectl port-forward -n observability svc/otel-collector 55679:55679
# Open http://localhost:55679/debug/tracez
```

### 5. Test Service Discovery

For services to be discovered by Prometheus, ensure they have the required annotations:

```yaml
metadata:
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "80"    # or your metrics port
    prometheus.io/path: "/metrics"
```

## Integrated Features

### From Docker Configurations:

1. **Prometheus**:
   - Service discovery for all my-project services
   - Individual scrape configs for 28 services
   - Infrastructure monitoring (Grafana, Jaeger)
   - Enhanced relabeling configurations

2. **Grafana**:
   - 9 pre-configured dashboards
   - Business KPIs, system metrics, trace analytics
   - Automatic dashboard provisioning
   - Exemplar support for trace correlation

3. **OTel Collector**:
   - Enhanced resource detection (CPU, OS details)
   - Optimized batching and memory limits
   - Advanced metric filtering
   - Span metrics with custom buckets
   - Compression and retry for Jaeger exports

## Rollback Procedure

If issues occur, rollback to original configurations:

```bash
# Rollback Prometheus
kubectl apply -f 18-prometheus/prometheus.yaml

# Rollback Grafana
kubectl apply -f 16-grafana/grafana.yaml

# Rollback OTel
kubectl apply -f 17-otel/otel-collector.yaml

# Restart all components
kubectl rollout restart statefulset/prometheus -n observability
kubectl rollout restart deployment/grafana -n observability
kubectl rollout restart deployment/otel-collector -n observability
```

## Notes

- Dashboard provisioning warnings in Grafana logs are expected and non-critical
- Services need proper annotations for Prometheus discovery
- The integrated configurations maintain all existing functionality while adding enhanced features from Docker setup