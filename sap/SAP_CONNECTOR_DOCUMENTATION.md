# SAP Connector Service Documentation

## Overview

This service provides a comprehensive SAP integration solution using official SAP libraries:
- **SAP Cloud Security Client** for authentication and authorization
- **SAP HANA Go Driver** for direct database access
- Enhanced webhook handling with security validation

## Architecture

### Core Components

1. **Enhanced SAP Client** (`/internal/sap/enhanced_client.go`)
   - Combines API access, direct HANA queries, and security features
   - Automatically selects the most efficient data access method
   - Supports both cloud and on-premise SAP deployments

2. **HANA Database Client** (`/internal/sap/hana_client.go`)
   - Direct database access using the official SAP HANA driver
   - Connection pooling and transaction support
   - Native support for SAP data types and stored procedures

3. **Security Client** (`/internal/sap/security_client.go`)
   - Integration with SAP Identity Authentication Service (IAS)
   - JWT token validation and claim extraction
   - Webhook authentication with API key validation

4. **Webhook Handler** (`/internal/rest/webhook.go`)
   - Processes both IDoc XML and JSON webhook payloads
   - Transactional processing with automatic rollback on failure
   - Enhanced security validation

## Configuration

### Environment Variables

```bash
# SAP API Configuration
SAP_BASE_URL=https://api.sap.company.com
SAP_API_KEY=REDACTED
SAP_WEBHOOK_SECRET=REDACTED
SAP_WEBHOOK_API_KEY=REDACTED

# OAuth2 Configuration  
SAP_CLIENT_ID=client-id
SAP_CLIENT_SECRET=REDACTED
SAP_TOKEN_URL=https://auth.sap.com/oauth/token

# HANA Database Configuration
SAP_HANA_HOST=hana.company.com
SAP_HANA_PORT=443
SAP_HANA_USER=db-user
SAP_HANA_PASSWORD=REDACTED
SAP_HANA_USE_TLS=true

# Security Configuration
SAP_IAS_INSTANCE_NAME=my-ias-instance
SAP_ISSUER=https://ias.company.com
SAP_AUDIENCE=sap-connector

# Feature Flags
SAP_USE_DIRECT_HANA=true  # Use direct DB queries when possible
SAP_ENABLE_SECURITY=true  # Enable IAS security validation
```

### Kubernetes Configuration

For Kubernetes deployments, mount IAS credentials as secrets:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: sap-ias-credentials
  namespace: sap-connector
data:
  clientid: <base64-encoded-client-id>
  clientsecret: <base64-encoded-client-secret>
  issuer: <base64-encoded-issuer>
  url: <base64-encoded-auth-url>
```

Mount at: `/etc/secrets/sapbtp/identity/<IAS_INSTANCE_NAME>/`

## Usage Examples

### 1. Initialize Enhanced SAP Client

```go
config := &sap.Config{
    BaseURL:       os.Getenv("SAP_BASE_URL"),
    APIKey:        os.Getenv("SAP_API_KEY"),
    WebhookSecret: os.Getenv("SAP_WEBHOOK_SECRET"),
    
    // HANA Configuration
    HANAHost:     os.Getenv("SAP_HANA_HOST"),
    HANAPort:     os.Getenv("SAP_HANA_PORT"),
    HANAUser:     os.Getenv("SAP_HANA_USER"),
    HANAPassword: os.Getenv("SAP_HANA_PASSWORD"),
    HANAUseTLS:   true,
    
    // Security Configuration
    EnableSecurity:  true,
    IASInstanceName: os.Getenv("SAP_IAS_INSTANCE_NAME"),
    
    // Features
    UseDirectHANA: true,
}

client, err := sap.NewEnhancedSAPClient(config)
if err != nil {
    log.Fatal().Err(err).Msg("Failed to create SAP client")
}
defer client.Close()
```

### 2. Query Product Changes

```go
// Uses direct HANA query if enabled, otherwise falls back to API
changes, err := client.GetProductChangesEnhanced(ctx, time.Now().Add(-24*time.Hour))
if err != nil {
    log.Error().Err(err).Msg("Failed to get product changes")
    return err
}

for _, change := range changes {
    log.Info().
        Str("sku", change.SKU).
        Str("name", change.Name).
        Time("changedAt", change.ChangedAt).
        Msg("Product changed")
}
```

### 3. Direct HANA Queries

```go
// For complex queries, use the HANA client directly
if client.HANAClient != nil {
    rows, err := client.HANAClient.ExecuteStoredProcedure(ctx, 
        "GET_PRODUCT_HIERARCHY", 
        "MATERIAL_GROUP_001",
    )
    if err != nil {
        return err
    }
    defer rows.Close()
    
    // Process results...
}
```

### 4. Webhook Processing

The webhook endpoint automatically handles:
- Signature validation
- IAS security authentication (if enabled)
- IDoc XML parsing
- JSON event parsing
- Transactional processing

```bash
# Example webhook call
curl -X POST https://your-service.com/api/sap/webhook \
  -H "X-SAP-Signature: signature-value" \
  -H "X-SAP-API-Key: your-api-key" \
  -H "Content-Type: application/xml" \
  -d @idoc.xml
```

## Security Features

### 1. Multi-Layer Authentication
- API Key validation for webhooks
- OAuth2 client credentials for API calls
- JWT token validation with IAS
- Signature validation for webhook payloads

### 2. Secure Communication
- TLS encryption for all connections
- Certificate validation for HANA Cloud
- Secure credential storage (never in code)

### 3. Audit Trail
- All webhook events stored with status
- Comprehensive logging with correlation IDs
- Failed event tracking for retry

## Performance Optimizations

### 1. Direct Database Access
When `UseDirectHANA` is enabled:
- Queries go directly to HANA instead of through APIs
- Significant performance improvement for bulk operations
- Reduced latency and API rate limit concerns

### 2. Connection Pooling
- HANA client maintains connection pool
- Configurable pool size and connection lifetime
- Automatic connection health checks

### 3. Efficient Data Transfer
- Batch processing for multiple products
- Streaming for large result sets
- Optimized query patterns

## Error Handling

### 1. Webhook Processing
- Automatic transaction rollback on errors
- Event status tracking (received → processing → processed/failed)
- Detailed error logging with context

### 2. Retry Logic
- Configurable retry policies per entity type
- Exponential backoff for transient failures
- Dead letter queue for permanent failures

### 3. Circuit Breaker
- Protects against cascading failures
- Automatic recovery when service is healthy
- Metrics for monitoring circuit state

## Monitoring and Observability

### 1. Metrics
- Event processing duration
- Success/failure rates by event type
- HANA query performance
- Circuit breaker state

### 2. Logging
- Structured JSON logging
- Correlation IDs for request tracing
- Security audit logs

### 3. Health Checks
- Database connectivity
- API endpoint availability
- Authentication token validity

## IDoc Support

### Supported IDoc Types
- **MATMAS**: Material Master data
- **INVCON**: Inventory/Stock levels
- **COND_A**: Pricing conditions
- **ORDERS**: Purchase orders (future)

### IDoc Processing Flow
1. Receive IDoc XML via webhook
2. Parse and validate structure
3. Transform to canonical event format
4. Publish to event bus
5. Update sync status

## Troubleshooting

### Common Issues

1. **HANA Connection Failures**
   - Check TLS configuration for cloud instances
   - Verify firewall rules
   - Confirm credentials and permissions

2. **Webhook Authentication Errors**
   - Ensure API key is set in environment
   - Verify signature calculation method
   - Check IAS configuration in Kubernetes

3. **Performance Issues**
   - Enable direct HANA access if not already
   - Check connection pool settings
   - Monitor query execution plans

### Debug Mode

Enable detailed logging:
```bash
export LOG_LEVEL=debug
export SAP_DEBUG=true
```

## Best Practices

1. **Use Direct HANA for Bulk Operations**
   - Enable `UseDirectHANA` for production
   - Ideal for initial data loads and synchronization

2. **Implement Proper Error Handling**
   - Always check for nil HANA client before direct queries
   - Handle both transient and permanent failures differently

3. **Security First**
   - Enable IAS security for production deployments
   - Rotate API keys regularly
   - Use service accounts with minimal permissions

4. **Monitor Everything**
   - Set up alerts for failed webhooks
   - Track sync lag times
   - Monitor HANA query performance

## Migration Guide

### From Basic SAP Client to Enhanced Client

1. Update configuration to use `EnhancedSAPClient`
2. Add HANA connection details
3. Enable security features as needed
4. Test webhook validation with new security

### Database Schema

Run migrations for sync tracking:
```sql
-- Sync status tracking
CREATE TABLE sap_sync_status (
    id VARCHAR(36) PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(100) NOT NULL,
    last_synced_at TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY idx_entity (entity_type, entity_id)
);

-- Webhook event log
CREATE TABLE sap_webhook_events (
    id VARCHAR(36) PRIMARY KEY,
    event_id VARCHAR(100),
    event_type VARCHAR(50),
    source VARCHAR(50),
    signature VARCHAR(255),
    payload JSON,
    received_at TIMESTAMP NOT NULL,
    processed_at TIMESTAMP,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    INDEX idx_status (status),
    INDEX idx_event_id (event_id)
);
```