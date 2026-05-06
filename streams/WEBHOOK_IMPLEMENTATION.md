# Webhook Implementation for Streams Service

## Overview
The streams service now includes a comprehensive webhook system that follows patterns established in the payments, messages, and comments services. This allows external systems to receive real-time notifications about live streaming events.

## Architecture

### Core Components

1. **Domain Models** (`internal/domain/`)
   - `WebhookSubscription`: Manages webhook endpoints and their configurations
   - `WebhookDelivery`: Tracks delivery attempts and status
   - `WebhookEvent`: Represents events to be delivered

2. **Infrastructure** (`internal/infrastructure/`)
   - `WebhookClient`: Handles HTTP delivery with retry logic
   - Implements exponential backoff with jitter
   - Supports custom headers and HMAC signatures

3. **Event Dispatcher** (`internal/handlers/webhook_dispatcher.go`)
   - Worker pool for concurrent delivery
   - Automatic retry mechanism
   - Failed delivery tracking

4. **Repository Layer** (`internal/postgres/`)
   - PostgreSQL implementations for webhook persistence
   - Efficient queries for pending deliveries
   - Delivery history tracking

## Supported Events

All live streaming domain events are automatically available as webhooks:

- `live_stream.created` - When a new live stream is scheduled
- `live_stream.started` - When streaming begins
- `live_stream.stopped` - When streaming ends
- `viewer.joined` - When a viewer joins the stream
- `viewer.left` - When a viewer leaves
- `stream.quality_changed` - When stream quality changes
- `stream.health_updated` - Stream health metrics updates
- `cdn.endpoint_added` - When CDN endpoints are configured
- `drm.configured` - When DRM is set up

## API Endpoints

### Webhook Management
- `POST /api/streams/webhooks` - Create webhook subscription
- `GET /api/streams/webhooks` - List all subscriptions
- `GET /api/streams/webhooks/{id}` - Get specific subscription
- `PUT /api/streams/webhooks/{id}` - Update subscription
- `DELETE /api/streams/webhooks/{id}` - Delete subscription

### Testing & Monitoring
- `POST /api/streams/webhooks/{id}/test` - Send test webhook
- `GET /api/streams/webhooks/{id}/deliveries` - View delivery history
- `POST /api/streams/webhooks/{id}/retry/{deliveryId}` - Retry failed delivery

## Usage Example

### Creating a Webhook Subscription
```bash
curl -X POST http://localhost:8080/api/streams/webhooks \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/webhooks/streams",
    "events": ["live_stream.started", "live_stream.stopped"],
    "headers": {
      "X-Custom-Header": "value"
    },
    "retry_policy": {
      "max_retries": 5,
      "backoff_factor": 2.0,
      "initial_delay": 1000,
      "max_backoff": 300000
    }
  }'
```

### Webhook Payload Format
```json
{
  "id": "evt_123456",
  "type": "live_stream.started",
  "stream_id": "str_789",
  "timestamp": "2024-01-20T10:00:00Z",
  "data": {
    "stream_id": "str_789",
    "started_at": 1705744800
  }
}
```

## Security

1. **HMAC Signatures**: Each webhook includes an HMAC-SHA256 signature in the `X-Webhook-Signature` header
2. **Secret Generation**: Secrets are automatically generated if not provided
3. **HTTPS Only**: Webhooks are only sent to HTTPS endpoints in production
4. **IP Whitelisting**: Optional IP restrictions can be configured

## Reliability

1. **Retry Logic**: 
   - Exponential backoff with configurable parameters
   - Default: 3 retries with 2x backoff starting at 1 second
   - Maximum backoff: 5 minutes

2. **Delivery Tracking**:
   - All delivery attempts are logged
   - Failed deliveries can be manually retried
   - Automatic cleanup of old delivery records

3. **Concurrent Delivery**:
   - Worker pool processes webhooks concurrently
   - Default: 10 concurrent workers
   - Queue size: 1000 events

## Database Schema

The webhook system uses two main tables:

1. **webhook_subscriptions**: Stores subscription configurations
2. **webhook_deliveries**: Tracks delivery attempts and results

See `migrations/002_add_webhooks.sql` for the complete schema.

## Integration with Domain Events

The webhook system is fully integrated with the domain event system:

1. Domain events are automatically dispatched to the webhook system
2. Event type mapping converts internal events to webhook-friendly names
3. Asynchronous processing ensures domain operations aren't blocked

## Monitoring

- Delivery success/failure rates are tracked
- Failed deliveries are logged with detailed error information
- Webhook statistics available via the delivery history API

## Future Enhancements

1. **Event Filtering**: Allow subscriptions to filter events by criteria
2. **Batch Delivery**: Group multiple events into single webhook calls
3. **Rate Limiting**: Prevent webhook spam
4. **Circuit Breaker**: Temporarily disable failing endpoints
5. **Webhook Templates**: Customizable payload formats