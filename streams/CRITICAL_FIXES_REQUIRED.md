# Critical Fixes Required - Streams Service

## 🚨 MUST FIX BEFORE PRODUCTION

### 1. Security Fix: WebRTC Origin Validation

**File**: `internal/infrastructure/streaming/webrtc_server.go`

```go
// CURRENT (VULNERABLE):
func checkOrigin(r *http.Request) bool {
    // TODO: Implement proper origin validation
    return true
}

// REQUIRED FIX:
func checkOrigin(r *http.Request) bool {
    origin := r.Header.Get("Origin")
    allowedOrigins := []string{
        "https://yourdomain.com",
        "https://app.yourdomain.com",
        "http://localhost:3000", // Dev only
    }
    
    for _, allowed := range allowedOrigins {
        if origin == allowed {
            return true
        }
    }
    return false
}
```

### 2. Security Fix: Webhook URL Validation

**File**: `internal/infrastructure/webhook_client.go`

Add this validation before making requests:

```go
func (c *WebhookClient) validateURL(urlStr string) error {
    u, err := url.Parse(urlStr)
    if err != nil {
        return err
    }
    
    // Enforce HTTPS in production
    if c.enforceHTTPS && u.Scheme != "https" {
        return errors.New("HTTPS required for webhooks")
    }
    
    // Block private IP ranges
    if ip := net.ParseIP(u.Hostname()); ip != nil {
        if ip.IsLoopback() || ip.IsPrivate() {
            return errors.New("webhook URL points to private IP")
        }
    }
    
    return nil
}
```

### 3. Resource Leak Fix: WebRTC Cleanup

**File**: `internal/infrastructure/streaming/webrtc_server.go`

```go
func (ws *WebRTCServer) handleSignaling(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
    
    // Add cleanup
    defer func() {
        ws.removePeer(peerID)
        if conn != nil {
            conn.Close()
        }
    }()
    
    // Add panic recovery
    defer func() {
        if r := recover(); r != nil {
            ws.logger.Error("Panic in WebRTC handler", 
                zap.Any("panic", r),
                zap.Stack("stack"))
        }
    }()
}
```

### 4. Database Fix: Add Missing Indexes

**File**: `migrations/003_add_production_indexes.sql`

```sql
-- +migrate Up

-- Performance indexes
CREATE INDEX idx_live_streams_status ON live_streams(status) WHERE status = 'active';
CREATE INDEX idx_live_streams_scheduled ON live_streams(scheduled_start_time) WHERE status = 'scheduled';
CREATE INDEX idx_viewer_sessions_active ON viewer_sessions(stream_id, user_id) WHERE left_at IS NULL;

-- Webhook performance
CREATE INDEX idx_webhook_deliveries_pending ON webhook_deliveries(next_retry_at) 
    WHERE status IN ('pending', 'retrying');

-- +migrate Down
DROP INDEX IF EXISTS idx_live_streams_status;
DROP INDEX IF EXISTS idx_live_streams_scheduled;
DROP INDEX IF EXISTS idx_viewer_sessions_active;
DROP INDEX IF EXISTS idx_webhook_deliveries_pending;
```

### 5. Configuration Fix: Environment Variables

**File**: `internal/config/config.go` (create new)

```go
package config

import (
    "time"
    "github.com/kelseyhightower/envconfig"
)

type Config struct {
    // Server
    HTTPPort int `envconfig:"HTTP_PORT" default:"8080"`
    GRPCPort int `envconfig:"GRPC_PORT" default:"9090"`
    
    // Streaming
    RTMPPort          int    `envconfig:"RTMP_PORT" default:"1935"`
    SRTPort           int    `envconfig:"SRT_PORT" default:"4578"`
    HLSSegmentDuration int    `envconfig:"HLS_SEGMENT_DURATION" default:"6"`
    MaxConcurrentViewers int  `envconfig:"MAX_CONCURRENT_VIEWERS" default:"10000"`
    
    // WebRTC
    STUNServers []string `envconfig:"STUN_SERVERS" default:"stun:stun.l.google.com:19302"`
    TURNServer  string   `envconfig:"TURN_SERVER"`
    TURNUser    string   `envconfig:"TURN_USER"`
    TURNPass    string   `envconfig:"TURN_PASS"`
    
    // Security
    JWTSecret       string   `envconfig:"JWT_SECRET" required:"true"`
    AllowedOrigins  []string `envconfig:"ALLOWED_ORIGINS" default:"https://yourdomain.com"`
    WebhookHTTPSOnly bool    `envconfig:"WEBHOOK_HTTPS_ONLY" default:"true"`
    
    // Webhooks
    WebhookMaxRetries     int           `envconfig:"WEBHOOK_MAX_RETRIES" default:"3"`
    WebhookTimeout        time.Duration `envconfig:"WEBHOOK_TIMEOUT" default:"30s"`
    WebhookWorkers        int           `envconfig:"WEBHOOK_WORKERS" default:"10"`
    
    // Database
    DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
    
    // CDN
    CDNProvider string `envconfig:"CDN_PROVIDER" default:"cloudflare"`
    CDNAPIKey   string `envconfig:"CDN_API_KEY"`
    CDNZoneID   string `envconfig:"CDN_ZONE_ID"`
}

func Load() (*Config, error) {
    var cfg Config
    err := envconfig.Process("STREAMS", &cfg)
    return &cfg, err
}
```

### 6. Error Handling Fix: Domain Errors

**File**: `internal/domain/errors.go` (create new)

```go
package domain

import "errors"

var (
    // Stream errors
    ErrStreamNotFound = errors.New("stream not found")
    ErrStreamAlreadyActive = errors.New("stream already active")
    ErrStreamNotActive = errors.New("stream not active")
    ErrMaxViewersReached = errors.New("maximum viewers reached")
    
    // Webhook errors
    ErrWebhookNotFound = errors.New("webhook subscription not found")
    ErrWebhookInvalidURL = errors.New("invalid webhook URL")
    ErrWebhookDeliveryFailed = errors.New("webhook delivery failed")
    
    // Auth errors
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden = errors.New("forbidden")
)
```

### 7. Panic Recovery Middleware

**File**: `internal/middleware/recovery.go` (create new)

```go
package middleware

import (
    "net/http"
    "runtime/debug"
    "go.uber.org/zap"
)

func Recovery(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    logger.Error("Panic recovered",
                        zap.Any("error", err),
                        zap.String("path", r.URL.Path),
                        zap.String("method", r.Method),
                        zap.String("stack", string(debug.Stack())))
                    
                    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                }
            }()
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### 8. Graceful Shutdown

**File**: `cmd/service/main.go` (update)

```go
// Add graceful shutdown
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Handle shutdown signals
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

go func() {
    <-sigCh
    logger.Info("Shutting down gracefully...")
    
    // Stop accepting new connections
    httpServer.SetKeepAlivesEnabled(false)
    
    // Wait for existing connections
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()
    
    if err := httpServer.Shutdown(shutdownCtx); err != nil {
        logger.Error("HTTP server shutdown error", zap.Error(err))
    }
    
    // Stop other services
    webhookDispatcher.Stop()
    streamingServer.Stop()
    
    cancel()
}()
```

## Deployment Checklist

Before deploying to production, ensure:

- [ ] All security fixes implemented and tested
- [ ] Resource leaks fixed and verified
- [ ] Configuration moved to environment variables
- [ ] Database migrations run successfully
- [ ] Panic recovery middleware active
- [ ] Graceful shutdown tested
- [ ] Authentication enabled on all endpoints
- [ ] Rate limiting configured
- [ ] Monitoring and alerting set up
- [ ] Load testing completed

## Testing Commands

```bash
# Test WebRTC origin validation
curl -H "Origin: https://evil.com" http://localhost:8080/ws/signaling

# Test webhook HTTPS enforcement
curl -X POST http://localhost:8080/api/streams/webhooks \
  -d '{"url": "http://insecure.com/webhook"}'

# Test graceful shutdown
kill -SIGTERM $(pgrep streams)

# Check for goroutine leaks
go test -race ./...
```