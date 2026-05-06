# Consciousness Quick Start Guide

## 1. Prerequisites

- Go 1.21+
- PostgreSQL
- NATS JetStream
- API keys for at least one AI provider (OpenAI, Anthropic, or DeepSeek)

## 2. Configuration

Create a `.env` file in the managers directory:

```bash
# Minimum required configuration
MANAGER_CONSCIOUSNESS_ENABLED=true
AI_PROVIDER_DEFAULT=deepseek
AI_PROVIDER_DEEPSEEK_API_KEY=REDACTED

# Database
DATABASE_URL=postgresql://user:password@localhost:5432/managers

# NATS
NATS_URL=nats://localhost:4222
```

## 3. Build & Run

```bash
# Build the service
make build-managers

# Run with consciousness enabled
MANAGER_CONSCIOUSNESS_ENABLED=true make run-managers
```

## 4. Verify It's Working

Check the logs for consciousness initialization:

```
INFO Consciousness manager initialized
INFO Routing event to consciousness manager
INFO Pattern detected pattern_type=xxx
INFO Decision made decision_id=xxx
```

## 5. Test with Sample Events

Send a test event to trigger consciousness:

```bash
# Example: High-value order that might trigger fraud detection
nats pub orders.events '{
  "type": "OrderCreated",
  "aggregate_id": "order-123",
  "data": {
    "order_id": "order-123",
    "user_id": "user-456",
    "total_amount": 5000,
    "items_count": 1,
    "shipping_address": {
      "country": "US",
      "city": "New York"
    }
  }
}'
```

## 6. Monitor Actions

Watch for autonomous actions in the logs:

```
INFO Executing autonomous decision decision_id=xxx action_count=1
INFO Action executed successfully action_type=flag_order
```

## 7. Common Patterns to Test

### Abandoned Cart Recovery
```bash
nats pub baskets.events '{
  "type": "BasketAbandoned",
  "aggregate_id": "basket-789",
  "data": {
    "basket_id": "basket-789",
    "user_id": "user-123",
    "total_value": 150,
    "items_count": 3,
    "abandoned_at": "2024-01-20T10:00:00Z"
  }
}'
```

### Support Ticket Escalation
```bash
nats pub support.events '{
  "type": "TicketCreated",
  "aggregate_id": "ticket-456",
  "data": {
    "ticket_id": "ticket-456",
    "user_id": "user-789",
    "priority": "urgent",
    "subject": "Order not received",
    "sentiment_score": -0.8
  }
}'
```

## 8. Debugging

Enable debug logging:
```bash
LOG_LEVEL=debug MANAGER_CONSCIOUSNESS_ENABLED=true make run-managers
```

## 9. Disable Consciousness

To run without consciousness (standard mode):
```bash
MANAGER_CONSCIOUSNESS_ENABLED=false make run-managers
```

## 10. Next Steps

- Review decision rules in `internal/consciousness/rules/rules.yaml`
- Adjust confidence thresholds in `.env`
- Monitor success rates in logs
- Enable additional AI providers for fallback