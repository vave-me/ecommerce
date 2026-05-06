# NetSuite ERP Connector

Production-ready NetSuite connector implementing OAuth 1.0a Token-Based Authentication (TBA) with comprehensive support for products, inventory, customers, and orders.

## Features

- **OAuth 1.0a TBA Authentication**: Secure token-based authentication
- **SuiteQL Support**: Direct SQL-like queries to NetSuite data
- **RESTlet Integration**: Custom business logic execution
- **Webhook Support**: Real-time event notifications
- **Comprehensive Data Sync**: Products, inventory, customers, orders, and prices
- **Rate Limiting**: Built-in rate limiting and retry logic
- **Error Handling**: Robust error handling with detailed logging

## Configuration

```go
config := &erp.ConnectorConfig{
    ID:   "netsuite-prod",
    Type: erp.ERPTypeNetSuite,
    Auth: erp.AuthConfig{
        ConsumerKey:    "your_consumer_key",
        ConsumerSecret: "your_consumer_secret",
        TokenID:        "your_token_id",
        TokenSecret:    "your_token_secret",
    },
    Metadata: map[string]interface{}{
        "account_id":   "123456",                    // Required
        "datacenter":   "sb1",                       // Optional, auto-detected from account_id
        "api_version":  "v2",                        // Optional, defaults to v1
        "restlet_url":  "https://custom.url.com",   // Optional, custom RESTlet URL
        
        // RESTlet paths (optional)
        "suiteql_restlet":      "/restlets/suiteql",
        "create_order_restlet": "/restlets/create_order",
    },
    Webhook: erp.WebhookConfig{
        Secret: "your_webhook_secret",
    },
    Sync: erp.SyncConfig{
        Interval:  time.Minute * 15,
        BatchSize: 100,
    },
}
```

## Supported Operations

### Products
- Fetch products modified since a specific time
- Transform NetSuite items to canonical product format
- Support for custom fields and metadata

### Inventory
- Real-time stock level updates
- Multi-location inventory support
- Available vs on-hand quantity tracking

### Customers
- Customer data synchronization
- Support for company and individual customers
- Address and contact information

### Orders
- Sales order creation and updates
- Line item details with tax calculations
- Order status tracking

### Pricing
- Price level support
- Multi-currency pricing
- Quantity-based pricing

## Webhook Events

The connector supports the following NetSuite webhook events:

- `item.create` → `ProductCreated`
- `item.update` → `ProductMasterUpdated`
- `inventoryadjustment.create` → `StockLevelUpdated`
- `inventoryadjustment.update` → `StockLevelUpdated`
- `itemfulfillment.create` → `StockLevelUpdated`
- `salesorder.create` → `OrderCreated`
- `salesorder.update` → `OrderUpdated`
- `customer.create` → `CustomerCreated`
- `customer.update` → `CustomerUpdated`

## SuiteQL Queries

The connector uses SuiteQL for efficient data retrieval:

```sql
-- Example: Fetch products
SELECT 
    i.id,
    i.itemid,
    i.displayname,
    i.salesdescription,
    i.baseprice,
    i.cost,
    i.weight,
    i.itemtype,
    i.isinactive,
    i.lastmodifieddate
FROM item i
WHERE i.lastmodifieddate > ?
AND i.itemtype IN ('InvtPart', 'NonInvtPart', 'Kit', 'Assembly')
ORDER BY i.lastmodifieddate ASC
LIMIT ?
```

## RESTlet Requirements

The connector expects the following RESTlet endpoints:

### SuiteQL RESTlet
Executes SuiteQL queries and returns results.

```javascript
// Expected request format
{
    "query": "SELECT * FROM item WHERE id = ?",
    "params": [123]
}

// Expected response format
{
    "items": [...],
    "hasMore": false,
    "totalCount": 1
}
```

### Create Order RESTlet
Creates sales orders in NetSuite.

```javascript
// Expected request format
{
    "recordtype": "salesorder",
    "entity": "CUST-001",
    "items": [...],
    // ... other order fields
}

// Expected response format
{
    "success": true,
    "id": "12345",
    "tranid": "SO-001"
}
```

## Testing

Run the test suite:

```bash
go test ./internal/erp/connectors/netsuite/...
```

## Performance Considerations

- Uses connection pooling for HTTP requests
- Implements exponential backoff for retries
- Respects NetSuite API rate limits
- Batches requests where possible
- Caches OAuth nonce values

## Security

- HMAC-SHA256 signature validation for webhooks
- OAuth 1.0a signature generation for API requests
- Secure credential storage
- No logging of sensitive data

## Troubleshooting

### Common Issues

1. **Invalid OAuth Signature**
   - Verify all OAuth credentials are correct
   - Check timestamp synchronization
   - Ensure proper URL encoding

2. **RESTlet Not Found**
   - Verify RESTlet deployment IDs
   - Check script permissions
   - Ensure proper URL paths in configuration

3. **Rate Limiting**
   - Adjust sync intervals
   - Reduce batch sizes
   - Implement request queuing

### Debug Logging

Enable debug logging for detailed request/response information:

```go
logger := zap.NewDevelopment()
connector, _ := NewConnector(config, logger)
```