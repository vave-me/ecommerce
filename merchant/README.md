# Merchant Service

The Merchant Service is responsible for synchronizing product catalogs with Google Merchant Center, enabling businesses to sell their products through Google Shopping and other Google surfaces.

## Overview

This service provides:
- Automated product synchronization with Google Merchant Center
- Product data validation according to Google's requirements
- Batch operations for efficient catalog management
- Sync status tracking and error handling
- Event-driven updates from the product service

## Architecture

The service follows Domain-Driven Design (DDD) principles with:
- **Domain Layer**: Core business logic and entities
- **Application Layer**: Use cases and orchestration
- **Infrastructure Layer**: External integrations (Google Merchant Center API)
- **Adapters Layer**: Repository implementations and API wrappers

## Prerequisites

1. **Google Merchant Center Account**: You need an active merchant account
2. **Service Account**: Create a service account with Content API access
3. **Service Account JSON**: Download the credentials file
4. **Environment Variables**:
   ```bash
   MERCHANT_ID=your_merchant_id
   SERVICE_ACCOUNT_JSON_PATH=/path/to/service-account.json
   ```

## API Endpoints

### gRPC Service

- `SyncProduct`: Sync a single product to Google Merchant Center
- `BatchSyncProducts`: Sync multiple products in batch
- `RemoveProduct`: Remove a product from Google Merchant Center
- `GetProductStatus`: Get synchronization status for a product
- `ListProducts`: List all synchronized products

### REST API

All gRPC endpoints are also available via REST through the gRPC-Gateway:

- `POST /api/merchant/products/sync` - Sync single product
- `POST /api/merchant/products/batch-sync` - Batch sync products
- `DELETE /api/merchant/products/{product_id}` - Remove product
- `GET /api/merchant/products/{product_id}/status` - Get sync status
- `GET /api/merchant/products` - List products

## Product Data Requirements

Products must include these required fields (validated automatically):
- `offerId`: Unique product identifier
- `title`: Product title (max 150 characters)
- `description`: Product description (max 5000 characters)
- `link`: Product landing page URL
- `imageLink`: Main product image URL
- `price`: Price with currency (ISO 4217)
- `availability`: in_stock, out_of_stock, preorder, or backorder
- `contentLanguage`: Two-letter ISO 639 language code
- `targetCountry`: Two-letter ISO 3166 country code

## Event Integration

The service listens to product events from the product service:
- `ProductAdded`: Creates new product in Google Merchant Center
- `ProductUpdated`: Updates existing product
- `ProductPriceIncreased/Decreased`: Updates pricing
- `ProductStockAdjusted`: Updates availability
- `ProductRemoved`: Removes from Google Merchant Center

## Error Handling

The service implements:
- Automatic retry with exponential backoff for transient errors
- Validation before API calls to prevent invalid submissions
- Comprehensive error logging and monitoring
- Sync status tracking in database

## Rate Limiting

Google Merchant Center API has quotas:
- 2,500 requests per day
- The service implements rate limiting (100ms between requests)
- Batch operations are recommended for bulk updates

## Monitoring

The service provides:
- OpenTelemetry tracing for all operations
- Prometheus metrics for API calls
- Structured logging with operation context
- Sync statistics endpoint

## Database Schema

### products_sync_status
Tracks synchronization status for each product:
- `product_id`: Internal product ID
- `merchant_id`: Google Merchant Center ID
- `sync_status`: PENDING, SYNCED, FAILED, REMOVED
- `last_synced_at`: Last successful sync timestamp
- `last_error`: Error message if sync failed

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o merchant-service ./cmd/service
```

### Running Locally
```bash
./merchant-service
```

## Troubleshooting

### Common Issues

1. **Authentication Errors**: Ensure service account has Content API access
2. **Validation Errors**: Check product data meets Google's requirements
3. **Rate Limit Errors**: Implement batch operations for bulk updates
4. **Network Errors**: Service includes automatic retry logic

### Debug Mode

Enable debug logging:
```bash
LOG_LEVEL=debug ./merchant-service
```

## Best Practices

1. **Use Batch Operations**: For bulk updates, use BatchSyncProducts
2. **Monitor Sync Status**: Check sync status regularly for failed products
3. **Handle Events Properly**: Ensure product service events are processed
4. **Validate Early**: Use the validator before attempting sync
5. **Rate Limit Awareness**: Space out requests to avoid quota issues