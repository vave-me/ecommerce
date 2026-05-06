# Media Service Bulk Import Enhancement Plan

## Executive Summary

This document outlines a comprehensive plan to enhance the media service to handle bulk image imports from enterprise systems (ERP, SAP, CRM). The plan addresses scenarios where only image URLs and SKUs are provided, and scales from 1,000 to 100,000+ images.

## Current State Analysis

### Existing Architecture
- **Storage**: MinIO (S3-compatible)
- **Database**: PostgreSQL with event sourcing
- **API**: gRPC-based single-item operations
- **Pattern**: Direct client uploads via presigned URLs

### Limitations
- No batch operations
- No bulk import capabilities
- No external system integration
- Limited to single-item processing
- No progress tracking for large operations

## Proposed Architecture

### Core Components

#### 1. Import Session Management
Tracks and manages bulk import operations with full lifecycle support.

```go
type ImportSession struct {
    ID               string
    ExternalSystemID string
    TotalImages      int
    ProcessedImages  int
    FailedImages     int
    Status           string
    StartedAt        time.Time
}
```

#### 2. Queue-Based Processing Pipeline
Leverages NATS JetStream for reliable, scalable message processing:

- `media.import.pending` - Initial import requests
- `media.import.fetch` - URL fetching tasks
- `media.import.process` - Image processing
- `media.import.complete` - Completion notifications

#### 3. Worker Pool Architecture
Configurable worker pools for parallel processing:

```go
type WorkerPool struct {
    workers    []*ImageFetchWorker
    jobQueue   chan *ImportJob
    resultChan chan *JobResult
}
```

## Implementation Phases

### Phase 1: Database Schema (Week 1)

```sql
-- Import session tracking
CREATE TABLE import_sessions (
    id UUID PRIMARY KEY,
    external_system_id VARCHAR(50) NOT NULL,
    external_system_type VARCHAR(20) NOT NULL,
    total_images INT NOT NULL,
    processed_images INT DEFAULT 0,
    failed_images INT DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    metadata JSONB
);

-- Import items tracking
CREATE TABLE import_items (
    id UUID PRIMARY KEY,
    session_id UUID REFERENCES import_sessions(id),
    external_id VARCHAR(255) NOT NULL,
    sku VARCHAR(100) NOT NULL,
    image_url TEXT NOT NULL,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    media_id UUID REFERENCES media(id),
    image_id UUID REFERENCES images(id)
);

-- SKU mapping cache
CREATE TABLE sku_mappings (
    sku VARCHAR(100) PRIMARY KEY,
    item_id UUID NOT NULL,
    item_type VARCHAR(20) NOT NULL,
    last_used TIMESTAMP DEFAULT NOW()
);
```

### Phase 2: API Extensions (Week 2)

New gRPC endpoints for bulk operations:

```proto
service MediaService {
    rpc StartBulkImport(StartBulkImportRequest) returns (ImportSession);
    rpc AddImportBatch(AddImportBatchRequest) returns (BatchResult);
    rpc GetImportStatus(GetImportStatusRequest) returns (ImportStatus);
    rpc CancelImport(CancelImportRequest) returns (Empty);
}
```

### Phase 3: Processing Pipeline (Weeks 3-4)

#### Component 1: Import Receiver
- Validates import requests
- Creates import sessions
- Chunks data into manageable batches
- Publishes to processing queue

#### Component 2: URL Fetcher
- Fetches images from external URLs
- Handles retries with exponential backoff
- Streams directly to MinIO (no intermediate storage)
- Updates progress metrics

#### Component 3: Metadata Processor
- Maps SKU to internal ItemID
- Extracts image metadata
- Generates thumbnails
- Stores in PostgreSQL

### Phase 4: External System Adapters (Week 5)

Base adapter interface:

```go
type SystemAdapter interface {
    Name() string
    Connect(ctx context.Context) error
    FetchImageList(ctx context.Context, params FetchParams) (<-chan ImportItem, error)
    ValidateConnection() error
    Close() error
}
```

Specific implementations:
- SAP Adapter (OData/REST API)
- Generic ERP Adapter (REST/SOAP)
- CRM Adapter (Salesforce, Dynamics)
- CSV/Excel file adapter

### Phase 5: Performance Optimizations (Week 6)

#### 1. Batch Processing
- Process images in chunks of 100-500
- Configurable based on system capacity
- Parallel processing with worker pools

#### 2. Database Optimizations
- Bulk inserts using PostgreSQL COPY
- Batch transaction commits
- Connection pooling optimization

#### 3. Network Optimizations
- HTTP connection pooling
- Concurrent downloads (10-20 parallel)
- Rate limiting to respect external systems

#### 4. Storage Optimizations
- Direct streaming from URL to MinIO
- Parallel uploads to object storage
- No intermediate file storage

## Scaling Strategy

### Small Scale (1,000 images)
- 2 workers
- 100 image batches
- Sequential processing
- ~10 minutes completion

### Medium Scale (10,000 images)
- 5 workers
- 200 image batches
- Parallel processing
- ~30 minutes completion

### Large Scale (100,000+ images)
- 10-20 workers
- 500 image batches
- Distributed processing
- Progress checkpoints
- ~2-4 hours completion

## Error Handling

### Retry Strategy
```go
type RetryPolicy struct {
    MaxAttempts:     3
    InitialDelay:    1 * time.Second
    MaxDelay:        30 * time.Second
    BackoffFactor:   2.0
}
```

### Error Categories
1. **Transient Errors** (retry automatically)
   - Network timeouts
   - Temporary DNS failures
   - 503 Service Unavailable

2. **Permanent Errors** (mark as failed)
   - 404 Not Found
   - Invalid image format
   - Authorization failures

3. **System Errors** (pause import)
   - Database connection lost
   - MinIO unavailable
   - Out of memory

## Monitoring & Observability

### Key Metrics
- Import progress (images/second)
- Error rates by type
- External system latency
- Queue depths
- Worker utilization
- Memory usage
- Network bandwidth

### Dashboards
1. **Import Overview**
   - Active imports
   - Success/failure rates
   - Processing speed

2. **System Health**
   - Worker status
   - Queue backlogs
   - Resource usage

3. **Error Analysis**
   - Error types distribution
   - Failed items queue
   - Retry statistics

## Security Considerations

1. **Authentication**
   - Secure storage of external system credentials
   - OAuth2/API key management
   - Certificate validation for HTTPS

2. **Data Protection**
   - Encrypted connections to external systems
   - Secure temporary storage
   - Access control for import operations

3. **Rate Limiting**
   - Prevent DoS on external systems
   - Configurable request limits
   - Circuit breaker implementation

## Deployment Strategy

### Environment Requirements
- Additional MinIO storage capacity
- Increased database connections
- NATS JetStream configuration
- Monitoring infrastructure

### Rollout Plan
1. Deploy database migrations
2. Update media service with new APIs
3. Deploy worker pool infrastructure
4. Configure external system adapters
5. Enable monitoring and alerts
6. Gradual rollout with small test imports

## Testing Strategy

### Unit Tests
- Worker logic
- Retry mechanisms
- Adapter implementations

### Integration Tests
- End-to-end import flow
- Error scenarios
- Performance benchmarks

### Load Tests
- 1,000 image import
- 10,000 image import
- 100,000 image simulation

## Maintenance & Operations

### Regular Tasks
- Monitor import queue health
- Clean up completed imports
- Update SKU mappings
- Review error patterns

### Troubleshooting Guide
- Common error resolutions
- Performance tuning tips
- Debug logging locations
- Manual intervention procedures

## Success Criteria

1. **Performance**
   - Process 1,000 images in < 15 minutes
   - Process 100,000 images in < 4 hours
   - < 1% failure rate for valid images

2. **Reliability**
   - Automatic retry for transient failures
   - Graceful handling of system errors
   - Complete audit trail

3. **Usability**
   - Clear progress reporting
   - Easy error investigation
   - Simple retry mechanisms

## Timeline

- **Week 1**: Database schema and migrations
- **Week 2**: API design and implementation
- **Week 3-4**: Core processing pipeline
- **Week 5**: External system adapters
- **Week 6**: Performance optimization
- **Week 7**: Testing and documentation
- **Week 8**: Deployment and monitoring

## Conclusion

This plan provides a robust foundation for handling enterprise-scale image imports while maintaining system stability and performance. The modular design allows for incremental implementation and easy extension for future requirements.