# SAP Connector Implementation Plan for Products Service

## Executive Summary

This document outlines the implementation plan for creating a SAP connector/listener that integrates with the existing products service. The connector will handle bidirectional synchronization of product data between SAP and the products service using event-driven architecture.

## Architecture Overview

### Current State
- **Products Service**: Event-sourced microservice with CQRS pattern
- **Messaging**: NATS JetStream for async communication  
- **Events**: Domain events for all product/variant changes
- **APIs**: gRPC primary, REST gateway secondary

### Target Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   SAP System    │────▶│  SAP Connector   │────▶│ Products Service│
│                 │◀────│   Microservice   │◀────│                 │
└─────────────────┘     └──────────────────┘     └─────────────────┘
       │                         │                         │
       │                         │                         │
       └─────── IDoc/API ────────┴────── NATS/Events ─────┘
```

## Phase 1: SAP Connector Microservice Design

### Directory Structure
```
sap-connector/
├── cmd/
│   └── service/
│       └── main.go
├── internal/
│   ├── application/
│   │   ├── commands/
│   │   │   ├── sync_product_from_sap.go
│   │   │   ├── sync_stock_from_sap.go
│   │   │   └── sync_price_from_sap.go
│   │   ├── queries/
│   │   │   └── get_sap_sync_status.go
│   │   └── services/
│   │       ├── sap_listener.go
│   │       └── event_publisher.go
│   ├── domain/
│   │   ├── models/
│   │   │   ├── sap_product.go
│   │   │   ├── sap_stock.go
│   │   │   └── sap_price.go
│   │   └── events/
│   │       ├── canonical_events.go
│   │       └── sap_events.go
│   ├── infrastructure/
│   │   ├── sap/
│   │   │   ├── client.go
│   │   │   ├── idoc_parser.go
│   │   │   └── webhook_handler.go
│   │   ├── nats/
│   │   │   ├── publisher.go
│   │   │   └── subscriber.go
│   │   └── repository/
│   │       └── sync_status_repository.go
│   └── rest/
│       ├── handlers/
│       │   └── webhook_handler.go
│       └── server.go
├── pkg/
│   └── transformer/
│       ├── sap_to_canonical.go
│       └── canonical_to_sap.go
└── module.go
```

## Phase 2: Ingestion Layer Implementation

### 2.1 Push Model (IDoc/Webhook Handler)

```go
// internal/rest/handlers/webhook_handler.go
type WebhookHandler struct {
    sapService     application.SAPService
    eventPublisher application.EventPublisher
    logger         *slog.Logger
}

func (h *WebhookHandler) HandleSAPIDoc(w http.ResponseWriter, r *http.Request) {
    // 1. Validate authentication (API key, certificate)
    if !h.validateRequest(r) {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // 2. Read IDoc payload
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }
    
    // 3. Queue for processing
    if err := h.sapService.QueueIDoc(r.Context(), body); err != nil {
        h.logger.Error("Failed to queue IDoc", "error", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    
    // 4. Acknowledge receipt immediately
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
```

### 2.2 Pull Model (Scheduled Poller)

```go
// internal/application/services/sap_poller.go
type SAPPoller struct {
    sapClient      infrastructure.SAPClient
    lastSyncRepo   repository.SyncStatusRepository
    eventPublisher EventPublisher
    interval       time.Duration
}

func (p *SAPPoller) Start(ctx context.Context) {
    ticker := time.NewTicker(p.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            p.poll(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (p *SAPPoller) poll(ctx context.Context) {
    // 1. Get last sync timestamp
    lastSync, _ := p.lastSyncRepo.GetLastSync(ctx, "products")
    
    // 2. Query SAP for changes
    changes, err := p.sapClient.GetProductChanges(ctx, lastSync)
    if err != nil {
        p.logger.Error("Failed to poll SAP", "error", err)
        return
    }
    
    // 3. Process each change
    for _, change := range changes {
        if err := p.processChange(ctx, change); err != nil {
            p.logger.Error("Failed to process change", "error", err)
            continue
        }
    }
    
    // 4. Update last sync timestamp
    p.lastSyncRepo.UpdateLastSync(ctx, "products", time.Now())
}
```

## Phase 3: Canonical Event Model

### 3.1 Event Definitions

```go
// internal/domain/events/canonical_events.go
package events

import (
    "time"
    "github.com/google/uuid"
)

type EventType string

const (
    ProductMasterUpdated EventType = "ProductMasterUpdated"
    StockLevelUpdated    EventType = "StockLevelUpdated"
    PriceUpdated         EventType = "PriceUpdated"
    ProductCreated       EventType = "ProductCreated"
    ProductDeleted       EventType = "ProductDeleted"
)

type CanonicalEvent struct {
    EventID        uuid.UUID              `json:"eventId"`
    EventType      EventType              `json:"eventType"`
    EventTimestamp time.Time              `json:"eventTimestamp"`
    Source         string                 `json:"source"`
    CorrelationID  string                 `json:"correlationId,omitempty"`
    Payload        interface{}            `json:"payload"`
}

type ProductMasterPayload struct {
    SKU         string                 `json:"sku"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Category    string                 `json:"category"`
    Weight      float64                `json:"weight"`
    Dimensions  Dimensions             `json:"dimensions"`
    Attributes  map[string]interface{} `json:"attributes"`
}

type StockLevelPayload struct {
    SKU          string    `json:"sku"`
    WarehouseID  string    `json:"warehouseId"`
    Quantity     int       `json:"quantity"`
    StockType    string    `json:"stockType"`
    UpdatedAt    time.Time `json:"updatedAt"`
}

type PricePayload struct {
    SKU         string    `json:"sku"`
    PriceListID string    `json:"priceListId"`
    Currency    string    `json:"currency"`
    Price       float64   `json:"price"`
    ValidFrom   time.Time `json:"validFrom"`
    ValidTo     *time.Time `json:"validTo,omitempty"`
}
```

## Phase 4: Transformation Engine

### 4.1 SAP to Canonical Transformer

```go
// pkg/transformer/sap_to_canonical.go
package transformer

import (
    "github.com/google/uuid"
    "time"
    "sap-connector/internal/domain/events"
)

type SAPToCanonicalTransformer struct {
    sourceSystem string
}

func (t *SAPToCanonicalTransformer) TransformMATMAS(idoc *MATMASIDoc) (*events.CanonicalEvent, error) {
    // Extract material master data
    material := idoc.E1MARAM
    
    payload := events.ProductMasterPayload{
        SKU:         material.MATNR,
        Name:        material.MAKTX,
        Description: material.MAKTX_LONG,
        Category:    material.MTART,
        Weight:      parseFloat(material.BRGEW),
        Dimensions: events.Dimensions{
            Length: parseFloat(material.LAENG),
            Width:  parseFloat(material.BREIT),
            Height: parseFloat(material.HOEHE),
        },
        Attributes: map[string]interface{}{
            "baseUnit":     material.MEINS,
            "materialType": material.MTART,
            "division":     material.SPART,
        },
    }
    
    return &events.CanonicalEvent{
        EventID:        uuid.New(),
        EventType:      events.ProductMasterUpdated,
        EventTimestamp: time.Now(),
        Source:         t.sourceSystem,
        CorrelationID:  idoc.ControlRecord.DOCNUM,
        Payload:        payload,
    }, nil
}

func (t *SAPToCanonicalTransformer) TransformINVCON(idoc *INVCONIDoc) (*events.CanonicalEvent, error) {
    // Transform inventory/stock update
    inventory := idoc.E1INVCO
    
    payload := events.StockLevelPayload{
        SKU:         inventory.MATNR,
        WarehouseID: inventory.LGORT,
        Quantity:    parseInt(inventory.MENGE),
        StockType:   mapStockType(inventory.SOBKZ),
        UpdatedAt:   parseDate(inventory.AEDAT),
    }
    
    return &events.CanonicalEvent{
        EventID:        uuid.New(),
        EventType:      events.StockLevelUpdated,
        EventTimestamp: time.Now(),
        Source:         t.sourceSystem,
        CorrelationID:  idoc.ControlRecord.DOCNUM,
        Payload:        payload,
    }, nil
}

func (t *SAPToCanonicalTransformer) TransformCOND_A(idoc *COND_AIDoc) (*events.CanonicalEvent, error) {
    // Transform pricing condition
    condition := idoc.E1KOMG
    
    payload := events.PricePayload{
        SKU:         condition.MATNR,
        PriceListID: condition.KSCHL,
        Currency:    condition.WAERS,
        Price:       parseFloat(condition.KBETR),
        ValidFrom:   parseDate(condition.DATAB),
        ValidTo:     parseDatePtr(condition.DATBI),
    }
    
    return &events.CanonicalEvent{
        EventID:        uuid.New(),
        EventType:      events.PriceUpdated,
        EventTimestamp: time.Now(),
        Source:         t.sourceSystem,
        CorrelationID:  idoc.ControlRecord.DOCNUM,
        Payload:        payload,
    }, nil
}
```

## Phase 5: Event Publishing to Products Service

### 5.1 NATS Publisher

```go
// internal/infrastructure/nats/publisher.go
package nats

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/nats-io/nats.go"
    "sap-connector/internal/domain/events"
)

type EventPublisher struct {
    js     nats.JetStreamContext
    logger *slog.Logger
}

func (p *EventPublisher) PublishEvent(ctx context.Context, event *events.CanonicalEvent) error {
    // Determine the subject based on event type
    subject := p.getSubject(event.EventType)
    
    // Serialize event
    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }
    
    // Publish to NATS JetStream
    _, err = p.js.PublishAsync(subject, data,
        nats.MsgId(event.EventID.String()),
        nats.Header("Event-Type", string(event.EventType)),
        nats.Header("Source", event.Source),
    )
    
    if err != nil {
        return fmt.Errorf("failed to publish event: %w", err)
    }
    
    p.logger.Info("Published event",
        "eventId", event.EventID,
        "eventType", event.EventType,
        "subject", subject,
    )
    
    return nil
}

func (p *EventPublisher) getSubject(eventType events.EventType) string {
    switch eventType {
    case events.ProductMasterUpdated, events.ProductCreated:
        return "sap.products.master"
    case events.StockLevelUpdated:
        return "sap.products.stock"
    case events.PriceUpdated:
        return "sap.products.price"
    default:
        return "sap.products.unknown"
    }
}
```

## Phase 6: Products Service Event Consumers

### 6.1 Create SAP Event Handler

```go
// In products service: internal/handlers/sap_events.go
package handlers

import (
    "context"
    "encoding/json"
    "github.com/nats-io/nats.go"
    "products/internal/application"
    "sap-connector/internal/domain/events"
)

type SAPEventHandler struct {
    app    application.Application
    logger *slog.Logger
}

func (h *SAPEventHandler) SubscribeToSAPEvents(js nats.JetStreamContext) error {
    // Subscribe to product master updates
    _, err := js.Subscribe("sap.products.master", h.handleProductMasterUpdate,
        nats.Durable("products-sap-master"),
        nats.AckExplicit(),
    )
    if err != nil {
        return err
    }
    
    // Subscribe to stock updates
    _, err = js.Subscribe("sap.products.stock", h.handleStockUpdate,
        nats.Durable("products-sap-stock"),
        nats.AckExplicit(),
    )
    if err != nil {
        return err
    }
    
    // Subscribe to price updates
    _, err = js.Subscribe("sap.products.price", h.handlePriceUpdate,
        nats.Durable("products-sap-price"),
        nats.AckExplicit(),
    )
    
    return err
}

func (h *SAPEventHandler) handleProductMasterUpdate(msg *nats.Msg) {
    var event events.CanonicalEvent
    if err := json.Unmarshal(msg.Data, &event); err != nil {
        h.logger.Error("Failed to unmarshal event", "error", err)
        msg.Nak()
        return
    }
    
    var payload events.ProductMasterPayload
    if err := mapstructure.Decode(event.Payload, &payload); err != nil {
        h.logger.Error("Failed to decode payload", "error", err)
        msg.Nak()
        return
    }
    
    // Check if product exists
    product, err := h.app.Queries.GetProductBySKU.Handle(context.Background(), payload.SKU)
    if err != nil {
        // Create new product
        cmd := application.AddProductCommand{
            SKU:         payload.SKU,
            Name:        payload.Name,
            Description: payload.Description,
            Weight:      payload.Weight,
            Source:      "SAP",
        }
        
        if err := h.app.Commands.AddProduct.Handle(context.Background(), cmd); err != nil {
            h.logger.Error("Failed to add product", "error", err)
            msg.Nak()
            return
        }
    } else {
        // Update existing product
        cmd := application.UpdateProductCommand{
            ProductID:   product.ID,
            Name:        &payload.Name,
            Description: &payload.Description,
            Weight:      &payload.Weight,
        }
        
        if err := h.app.Commands.UpdateProduct.Handle(context.Background(), cmd); err != nil {
            h.logger.Error("Failed to update product", "error", err)
            msg.Nak()
            return
        }
    }
    
    msg.Ack()
}

func (h *SAPEventHandler) handleStockUpdate(msg *nats.Msg) {
    var event events.CanonicalEvent
    if err := json.Unmarshal(msg.Data, &event); err != nil {
        h.logger.Error("Failed to unmarshal event", "error", err)
        msg.Nak()
        return
    }
    
    var payload events.StockLevelPayload
    if err := mapstructure.Decode(event.Payload, &payload); err != nil {
        h.logger.Error("Failed to decode payload", "error", err)
        msg.Nak()
        return
    }
    
    // Find product by SKU
    product, err := h.app.Queries.GetProductBySKU.Handle(context.Background(), payload.SKU)
    if err != nil {
        h.logger.Error("Product not found", "sku", payload.SKU)
        msg.Ack() // Acknowledge to avoid redelivery
        return
    }
    
    // Adjust stock
    cmd := application.AdjustProductStockCommand{
        ProductID:    product.ID,
        NewQuantity:  payload.Quantity,
        WarehouseID:  payload.WarehouseID,
        AdjustmentBy: "SAP-SYNC",
    }
    
    if err := h.app.Commands.AdjustProductStock.Handle(context.Background(), cmd); err != nil {
        h.logger.Error("Failed to adjust stock", "error", err)
        msg.Nak()
        return
    }
    
    msg.Ack()
}

func (h *SAPEventHandler) handlePriceUpdate(msg *nats.Msg) {
    var event events.CanonicalEvent
    if err := json.Unmarshal(msg.Data, &event); err != nil {
        h.logger.Error("Failed to unmarshal event", "error", err)
        msg.Nak()
        return
    }
    
    var payload events.PricePayload
    if err := mapstructure.Decode(event.Payload, &payload); err != nil {
        h.logger.Error("Failed to decode payload", "error", err)
        msg.Nak()
        return
    }
    
    // Find product by SKU
    product, err := h.app.Queries.GetProductBySKU.Handle(context.Background(), payload.SKU)
    if err != nil {
        h.logger.Error("Product not found", "sku", payload.SKU)
        msg.Ack()
        return
    }
    
    // Update price based on current price
    var cmd interface{}
    if payload.Price > product.Price {
        cmd = application.IncreaseProductPriceCommand{
            ProductID:  product.ID,
            NewPrice:   payload.Price,
            Currency:   payload.Currency,
            IncreasedBy: "SAP-SYNC",
        }
    } else {
        cmd = application.DecreaseProductPriceCommand{
            ProductID:  product.ID,
            NewPrice:   payload.Price,
            Currency:   payload.Currency,
            DecreasedBy: "SAP-SYNC",
        }
    }
    
    if err := h.app.Commands.UpdatePrice.Handle(context.Background(), cmd); err != nil {
        h.logger.Error("Failed to update price", "error", err)
        msg.Nak()
        return
    }
    
    msg.Ack()
}
```

## Phase 7: Error Handling & Retry Mechanisms

### 7.1 Circuit Breaker Pattern

```go
// internal/infrastructure/resilience/circuit_breaker.go
package resilience

import (
    "github.com/sony/gobreaker"
    "time"
)

func NewSAPCircuitBreaker() *gobreaker.CircuitBreaker {
    settings := gobreaker.Settings{
        Name:        "SAP-API",
        MaxRequests: 3,
        Interval:    10 * time.Second,
        Timeout:     30 * time.Second,
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
            return counts.Requests >= 3 && failureRatio >= 0.6
        },
        OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
            logger.Info("Circuit breaker state change", 
                "name", name, 
                "from", from, 
                "to", to)
        },
    }
    
    return gobreaker.NewCircuitBreaker(settings)
}
```

### 7.2 Retry with Exponential Backoff

```go
// internal/infrastructure/resilience/retry.go
package resilience

import (
    "context"
    "fmt"
    "time"
)

type RetryConfig struct {
    MaxAttempts     int
    InitialDelay    time.Duration
    MaxDelay        time.Duration
    BackoffFactor   float64
}

func RetryWithBackoff(ctx context.Context, config RetryConfig, fn func() error) error {
    var err error
    delay := config.InitialDelay
    
    for attempt := 0; attempt < config.MaxAttempts; attempt++ {
        if attempt > 0 {
            select {
            case <-time.After(delay):
            case <-ctx.Done():
                return ctx.Err()
            }
        }
        
        err = fn()
        if err == nil {
            return nil
        }
        
        // Check if error is retryable
        if !isRetryable(err) {
            return err
        }
        
        // Calculate next delay
        delay = time.Duration(float64(delay) * config.BackoffFactor)
        if delay > config.MaxDelay {
            delay = config.MaxDelay
        }
    }
    
    return fmt.Errorf("max retry attempts reached: %w", err)
}

func isRetryable(err error) bool {
    // Define retryable errors (network, timeout, 5xx status codes)
    // Non-retryable: authentication errors, validation errors, 4xx status codes
    return true // Simplified
}
```

## Phase 8: Monitoring & Observability

### 8.1 Metrics Collection

```go
// internal/infrastructure/metrics/collector.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    EventsReceived = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "sap_connector_events_received_total",
        Help: "Total number of events received from SAP",
    }, []string{"event_type", "source"})
    
    EventsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "sap_connector_events_processed_total",
        Help: "Total number of events successfully processed",
    }, []string{"event_type", "destination"})
    
    EventProcessingDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name: "sap_connector_event_processing_duration_seconds",
        Help: "Duration of event processing",
    }, []string{"event_type"})
    
    SAPAPILatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name: "sap_api_request_duration_seconds",
        Help: "SAP API request latency",
    }, []string{"operation"})
    
    ErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "sap_connector_errors_total",
        Help: "Total number of errors",
    }, []string{"error_type", "component"})
)
```

### 8.2 Structured Logging

```go
// internal/infrastructure/logging/logger.go
package logging

import (
    "context"
    "log/slog"
    "os"
)

func NewLogger() *slog.Logger {
    opts := &slog.HandlerOptions{
        Level: slog.LevelInfo,
        AddSource: true,
    }
    
    handler := slog.NewJSONHandler(os.Stdout, opts)
    logger := slog.New(handler)
    
    return logger
}

func LogEvent(logger *slog.Logger, event *events.CanonicalEvent) {
    logger.Info("Processing event",
        slog.String("eventId", event.EventID.String()),
        slog.String("eventType", string(event.EventType)),
        slog.String("source", event.Source),
        slog.String("correlationId", event.CorrelationID),
        slog.Time("timestamp", event.EventTimestamp),
    )
}
```

## Phase 9: Configuration & Deployment

### 9.1 Configuration Structure

```yaml
# config/sap-connector.yaml
sap:
  endpoint: "https://sap.company.com/api"
  auth:
    type: "oauth2"
    clientId: "${SAP_CLIENT_ID}"
    clientSecret: "${SAP_CLIENT_SECRET}"
  polling:
    enabled: true
    interval: "2m"
  webhooks:
    enabled: true
    port: 8080
    path: "/webhooks/sap"

nats:
  url: "nats://localhost:4222"
  streamName: "SAP_EVENTS"
  subjects:
    - "sap.products.>"
    - "sap.inventory.>"
    - "sap.prices.>"

products_service:
  grpcEndpoint: "localhost:50051"
  httpEndpoint: "http://localhost:8080"

monitoring:
  metricsPort: 9090
  healthPort: 8081

logging:
  level: "info"
  format: "json"
```

### 9.2 Docker Deployment

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o sap-connector cmd/service/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/sap-connector .
COPY --from=builder /app/config ./config

EXPOSE 8080 9090 8081

CMD ["./sap-connector"]
```

### 9.3 Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sap-connector
  namespace: products
spec:
  replicas: 2
  selector:
    matchLabels:
      app: sap-connector
  template:
    metadata:
      labels:
        app: sap-connector
    spec:
      containers:
      - name: sap-connector
        image: myregistry/sap-connector:latest
        ports:
        - containerPort: 8080
          name: webhooks
        - containerPort: 9090
          name: metrics
        - containerPort: 8081
          name: health
        env:
        - name: SAP_CLIENT_ID
          valueFrom:
            secretKeyRef:
              name: sap-credentials
              key: clientId
        - name: SAP_CLIENT_SECRET
          valueFrom:
            secretKeyRef:
              name: sap-credentials
              key: clientSecret
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

## Phase 10: Testing Strategy

### 10.1 Unit Tests

```go
// internal/transformer/sap_to_canonical_test.go
package transformer

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestTransformMATMAS(t *testing.T) {
    transformer := NewSAPToCanonicalTransformer("SAP-TEST")
    
    idoc := &MATMASIDoc{
        E1MARAM: E1MARAM{
            MATNR: "TEST-SKU-001",
            MAKTX: "Test Product",
            MAKTX_LONG: "Test Product Long Description",
            MTART: "FERT",
            BRGEW: "1.5",
        },
    }
    
    event, err := transformer.TransformMATMAS(idoc)
    
    assert.NoError(t, err)
    assert.Equal(t, "ProductMasterUpdated", string(event.EventType))
    assert.Equal(t, "SAP-TEST", event.Source)
    
    payload := event.Payload.(events.ProductMasterPayload)
    assert.Equal(t, "TEST-SKU-001", payload.SKU)
    assert.Equal(t, "Test Product", payload.Name)
    assert.Equal(t, 1.5, payload.Weight)
}
```

### 10.2 Integration Tests with SAP Sandbox

```go
// test/integration/sap_integration_test.go
package integration

import (
    "context"
    "testing"
    "time"
)

func TestSAPProductSync(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Setup test environment
    ctx := context.Background()
    sapClient := setupSAPSandboxClient()
    connector := setupConnector(sapClient)
    
    // Create test product in SAP
    testProduct := &SAPProduct{
        MATNR: "INT-TEST-001",
        MAKTX: "Integration Test Product",
        BRGEW: 2.5,
    }
    
    err := sapClient.CreateProduct(ctx, testProduct)
    assert.NoError(t, err)
    
    // Wait for sync
    time.Sleep(5 * time.Second)
    
    // Verify product exists in products service
    product, err := productsClient.GetProductBySKU(ctx, "INT-TEST-001")
    assert.NoError(t, err)
    assert.Equal(t, "Integration Test Product", product.Name)
    assert.Equal(t, 2.5, product.Weight)
    
    // Cleanup
    sapClient.DeleteProduct(ctx, "INT-TEST-001")
}
```

## Implementation Timeline

1. **Week 1-2**: Core SAP Connector Architecture
   - Set up project structure
   - Implement basic webhook/polling infrastructure
   - Create canonical event models

2. **Week 3-4**: Transformation Engine
   - Build SAP to canonical transformers
   - Implement event publishing to NATS
   - Add error handling and retries

3. **Week 5-6**: Products Service Integration
   - Create event consumers in products service
   - Map canonical events to product commands
   - Test end-to-end flow

4. **Week 7-8**: Production Readiness
   - Add comprehensive monitoring
   - Implement circuit breakers
   - Create deployment configurations
   - Performance testing and optimization

## Key Success Factors

1. **Idempotency**: Ensure all operations are idempotent to handle duplicate events
2. **Monitoring**: Comprehensive observability from day one
3. **Error Handling**: Robust retry and circuit breaker patterns
4. **Testing**: Thorough unit and integration tests
5. **Documentation**: Clear API documentation and runbooks
6. **Security**: Secure credential management and encrypted communication

## Conclusion

This implementation plan provides a robust, scalable, and maintainable solution for SAP integration with your products service. The event-driven architecture ensures loose coupling, while the canonical event model provides flexibility for future integrations with other systems.