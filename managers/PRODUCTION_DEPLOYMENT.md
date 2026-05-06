# Production Deployment Guide for Consciousness-Enabled Managers Service

## Overview

The managers service has been transformed into a production-ready, self-conscious service that autonomously monitors platform events and takes appropriate actions. This guide covers deployment, configuration, monitoring, and operations.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Deployment Requirements](#deployment-requirements)
3. [Configuration](#configuration)
4. [Deployment Steps](#deployment-steps)
5. [Monitoring & Observability](#monitoring--observability)
6. [Operations](#operations)
7. [Troubleshooting](#troubleshooting)
8. [Performance Tuning](#performance-tuning)

## Architecture Overview

### Core Components

1. **ConsciousnessManager**: Orchestrates event processing through pattern detection → decision making → action execution
2. **Dynamic Tool Selection**: Accesses all 415+ platform tools based on event context
3. **Production Features**:
   - Comprehensive error handling with recovery
   - Circuit breakers for fault tolerance
   - Rate limiting with quotas
   - Distributed tracing
   - Prometheus metrics
   - Health/readiness probes
   - Graceful shutdown

### Event Flow

```
Platform Event (NATS) 
    → Rate Limiter
    → Circuit Breaker
    → Pattern Detection
    → Dynamic Tool Selection (415+ tools)
    → AI-Powered Decision Making
    → Autonomous Action Execution
    → Learning & Metrics
```

## Deployment Requirements

### Infrastructure

- **Kubernetes**: 1.24+
- **PostgreSQL**: 14+
- **NATS JetStream**: 2.9+
- **Redis**: 7.0+ (optional, for distributed rate limiting)
- **Prometheus**: For metrics collection
- **Jaeger/Tempo**: For distributed tracing

### Resource Requirements

```yaml
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "2Gi"
    cpu: "2000m"
```

### API Keys Required

- At least one AI provider (OpenAI, Anthropic, or DeepSeek)
- Recommended: Multiple providers for fallback

## Configuration

### Environment Variables

Create a ConfigMap or Secret with the following:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: managers-config
data:
  # Core Service Configuration
  DATABASE_URL: "postgresql://user:password@postgres:5432/managers"
  NATS_URL: "nats://nats:4222"
  
  # Consciousness Configuration
  MANAGER_CONSCIOUSNESS_ENABLED: "true"
  MANAGER_CONFIDENCE_THRESHOLD: "0.8"
  MANAGER_MAX_ACTIONS_PER_MINUTE: "10"
  MANAGER_LEARNING_ENABLED: "true"
  
  # AI Provider Configuration
  AI_PROVIDER_DEFAULT: "deepseek"
  AI_PROVIDER_FALLBACK_ENABLED: "true"
  
  # Tool Execution
  MANAGER_TOOL_TIMEOUT: "30s"
  MANAGER_TOOL_MAX_CONCURRENT: "10"
  
  # Rate Limiting
  RATE_LIMIT_ENABLED: "true"
  RATE_LIMIT_REQUESTS_PER_MINUTE: "100"
  
  # Circuit Breaker
  CIRCUIT_BREAKER_ENABLED: "true"
  CIRCUIT_BREAKER_MAX_FAILURES: "5"
  CIRCUIT_BREAKER_RESET_TIMEOUT: "2m"
  
  # Performance
  MAX_EVENT_PROCESSING_TIME: "5s"
  EVENT_BATCH_SIZE: "100"
  
  # Feature Flags
  FEATURE_AI_DECISIONS: "true"
  FEATURE_AUTONOMOUS_ACTIONS: "true"
  FEATURE_DRY_RUN_MODE: "false"
```

### Secrets

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: managers-secrets
type: Opaque
data:
  AI_PROVIDER_OPENAI_API_KEY: REDACTED
  AI_PROVIDER_ANTHROPIC_API_KEY: REDACTED
  AI_PROVIDER_DEEPSEEK_API_KEY: REDACTED
```

## Deployment Steps

### 1. Database Migration

```bash
# Run migrations
kubectl create job --from=cronjob/managers-migrate managers-migrate-manual
```

### 2. Deploy Service

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: managers
  labels:
    app: managers
spec:
  replicas: 3
  selector:
    matchLabels:
      app: managers
  template:
    metadata:
      labels:
        app: managers
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9090"
        prometheus.io/path: "/metrics"
    spec:
      containers:
      - name: managers
        image: your-registry/managers:latest
        ports:
        - containerPort: 8080  # HTTP
        - containerPort: 8081  # gRPC
        - containerPort: 9090  # Metrics
        envFrom:
        - configMapRef:
            name: managers-config
        - secretRef:
            name: managers-secrets
        livenessProbe:
          httpGet:
            path: /health/consciousness
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready/consciousness
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "2000m"
```

### 3. Service & Ingress

```yaml
apiVersion: v1
kind: Service
metadata:
  name: managers
spec:
  selector:
    app: managers
  ports:
  - name: http
    port: 80
    targetPort: 8080
  - name: grpc
    port: 8081
    targetPort: 8081
  - name: metrics
    port: 9090
    targetPort: 9090
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: managers
  annotations:
    nginx.ingress.kubernetes.io/backend-protocol: "GRPC"
spec:
  rules:
  - host: managers.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: managers
            port:
              number: 8081
```

### 4. Horizontal Pod Autoscaler

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: managers-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: managers
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

## Monitoring & Observability

### Prometheus Metrics

Key metrics to monitor:

```yaml
# Consciousness metrics
managers_consciousness_events_processed_total
managers_consciousness_patterns_detected_total
managers_consciousness_decisions_made_total
managers_consciousness_actions_executed_total
managers_consciousness_tool_duration_seconds

# System health
managers_consciousness_system_health
managers_consciousness_circuit_breaker_state
managers_consciousness_errors_total

# Performance
managers_consciousness_memory_usage_bytes
managers_consciousness_goroutines_count
```

### Grafana Dashboard

Import the provided dashboard JSON:

```json
{
  "dashboard": {
    "title": "Managers Consciousness System",
    "panels": [
      {
        "title": "Event Processing Rate",
        "targets": [{
          "expr": "rate(managers_consciousness_events_processed_total[5m])"
        }]
      },
      {
        "title": "Decision Success Rate",
        "targets": [{
          "expr": "rate(managers_consciousness_decisions_made_total{source='rule_based'}[5m])"
        }]
      },
      {
        "title": "Tool Execution Latency",
        "targets": [{
          "expr": "histogram_quantile(0.95, rate(managers_consciousness_tool_duration_seconds_bucket[5m]))"
        }]
      }
    ]
  }
}
```

### Alerts

```yaml
groups:
- name: managers_consciousness
  rules:
  - alert: ConsciousnessHighErrorRate
    expr: rate(managers_consciousness_errors_total[5m]) > 0.1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High error rate in consciousness system"
      
  - alert: CircuitBreakerOpen
    expr: managers_consciousness_circuit_breaker_state > 0
    for: 10m
    labels:
      severity: critical
    annotations:
      summary: "Circuit breaker {{ $labels.circuit_name }} is open"
      
  - alert: ConsciousnessNotProcessingEvents
    expr: rate(managers_consciousness_events_processed_total[10m]) == 0
    for: 10m
    labels:
      severity: critical
    annotations:
      summary: "Consciousness system not processing events"
```

### Logging

Structured logs with trace context:

```json
{
  "level": "info",
  "time": "2024-01-20T10:00:00Z",
  "trace_id": "abc123",
  "span_id": "def456",
  "component": "consciousness_manager",
  "event_type": "OrderCreated",
  "pattern_type": "fraud_risk",
  "confidence": 0.85,
  "decision_id": "dec-789",
  "message": "Decision made"
}
```

## Operations

### Starting/Stopping Consciousness

```bash
# Disable consciousness (will still process events normally)
kubectl set env deployment/managers MANAGER_CONSCIOUSNESS_ENABLED=false

# Enable dry-run mode (log actions without executing)
kubectl set env deployment/managers FEATURE_DRY_RUN_MODE=true

# Re-enable full consciousness
kubectl set env deployment/managers MANAGER_CONSCIOUSNESS_ENABLED=true FEATURE_DRY_RUN_MODE=false
```

### Adjusting Rate Limits

```bash
# Increase global rate limit
kubectl set env deployment/managers RATE_LIMIT_REQUESTS_PER_MINUTE=200

# Adjust AI decision threshold
kubectl set env deployment/managers MANAGER_CONFIDENCE_THRESHOLD=0.9
```

### Manual Circuit Breaker Reset

Access the service pod and use the management endpoint:

```bash
kubectl exec -it managers-pod -- curl -X POST http://localhost:8080/admin/circuit-breaker/reset
```

## Troubleshooting

### Common Issues

1. **High Memory Usage**
   - Check for memory leaks in pattern detection
   - Reduce EVENT_BATCH_SIZE
   - Increase memory limits

2. **Circuit Breakers Opening Frequently**
   - Check downstream service health
   - Adjust timeout settings
   - Review error logs for root cause

3. **AI Decisions Failing**
   - Verify API keys are correct
   - Check provider quotas
   - Enable fallback providers

4. **Tools Not Executing**
   - Verify tool permissions
   - Check comprehensive tool registry initialization
   - Review tool execution logs

### Debug Mode

Enable detailed logging:

```bash
kubectl set env deployment/managers LOG_LEVEL=debug
```

### Performance Profiling

Enable pprof endpoints:

```bash
kubectl port-forward pod/managers-pod 6060:6060
go tool pprof http://localhost:6060/debug/pprof/profile
```

## Performance Tuning

### Event Processing

```yaml
# Optimize for high throughput
MAX_EVENT_PROCESSING_TIME: "10s"
EVENT_BATCH_SIZE: "500"
EVENT_BUFFER_SIZE: "5000"
MANAGER_TOOL_MAX_CONCURRENT: "20"
```

### Memory Optimization

```yaml
# Reduce memory usage
RULES_RELOAD_INTERVAL: "15m"
MANAGER_MAX_ACTIONS_PER_MINUTE: "5"
EVENT_BATCH_SIZE: "50"
```

### AI Cost Optimization

```yaml
# Use cost-effective default provider
AI_PROVIDER_DEFAULT: "deepseek"
# Only use expensive providers for complex decisions
MANAGER_CONFIDENCE_THRESHOLD: "0.9"
```

## Security Considerations

1. **API Key Rotation**: Implement regular rotation of AI provider keys
2. **Rate Limiting**: Protect against abuse with proper limits
3. **Audit Logging**: All autonomous actions are logged for compliance
4. **Network Policies**: Restrict egress to required services only

## Backup and Recovery

1. **Decision History**: Stored in PostgreSQL, backup regularly
2. **Learning Data**: Export metrics periodically for analysis
3. **Configuration**: Version control all ConfigMaps and Secrets

## Conclusion

The consciousness-enabled managers service provides autonomous platform management with production-grade reliability. Monitor metrics closely during initial deployment and adjust thresholds based on your platform's behavior patterns.

For additional support, consult the technical documentation or contact the platform team.