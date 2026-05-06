# ERP Service REST API Documentation

## Overview

The ERP Service provides a comprehensive REST API for managing ERP connector integrations, synchronizing data, and handling business operations like invoices, returns, and inventory management.

## Authentication

All endpoints (except webhooks) require JWT authentication. Include the token in the Authorization header:

```bash
Authorization: Bearer <your-jwt-token>
```

## Base URL

```
http://localhost:8080/api/erp
```

## API Endpoints

### 1. Connector Management

#### List Connectors
```http
GET /connectors
```

Query Parameters:
- `type` (string): Filter by connector type (odoo, dynamics365, netsuite, etc.)
- `status` (string): Filter by status (active, inactive, error, maintenance)
- `page` (int): Page number (default: 1)
- `pageSize` (int): Items per page (default: 20)
- `sortBy` (string): Sort field
- `sortOrder` (string): Sort direction (asc, desc)

Response:
```json
{
  "connectors": [
    {
      "id": "string",
      "name": "string",
      "type": "string",
      "status": "string",
      "environment": "string",
      "baseUrl": "string",
      "createdAt": "datetime",
      "updatedAt": "datetime"
    }
  ],
  "totalCount": 0,
  "totalPages": 0,
  "currentPage": 0
}
```

#### Add Connector
```http
POST /connectors/v2
```

Request Body:
```json
{
  "name": "string",
  "type": "odoo|dynamics365|netsuite|sap|erpnext|frappe",
  "environment": "production|staging|development|sandbox",
  "baseUrl": "string",
  "authType": "basic|oauth2|oauth1|apikey|token",
  "authConfig": {
    // Type-specific auth fields
  },
  "webhookEnabled": boolean,
  "syncEnabled": boolean,
  "syncIntervalSeconds": number,
  "batchSize": number,
  "rateLimitRequestsPerSecond": number,
  "rateLimitBurst": number
}
```

Response:
```json
{
  "connectorId": "string",
  "message": "Connector created successfully"
}
```

#### Update Connector
```http
PUT /connectors/{connectorId}
```

Request Body: Same as Add Connector (partial updates supported)

#### Remove Connector
```http
DELETE /connectors/{connectorId}?force=true
```

Query Parameters:
- `force` (boolean): Force removal even with active syncs

#### Toggle Connector
```http
POST /connectors/{connectorId}/toggle
```

Request Body:
```json
{
  "activate": boolean,
  "reason": "string"
}
```

#### Get Connector Status
```http
GET /connectors/{connectorId}/status
```

Response:
```json
{
  "connector": {
    "id": "string",
    "name": "string",
    "type": "string",
    "status": "string"
  },
  "healthStatus": "healthy|unhealthy|unknown",
  "lastSync": "datetime",
  "pendingWebhooks": 0,
  "syncStatuses": {},
  "metadata": {}
}
```

### 2. Synchronization Operations

#### Get Sync History
```http
GET /connectors/{connectorId}/sync-history
```

Query Parameters:
- `entityType` (string): Filter by entity type
- `since` (datetime): Start date filter
- `until` (datetime): End date filter
- `page`, `pageSize`, `sortBy`, `sortOrder`

#### Sync Products
```http
POST /connectors/{connectorId}/sync/products
```

Request Body:
```json
{
  "batchSize": number,
  "since": "datetime",
  "filters": {}
}
```

#### Sync Stock
```http
POST /connectors/{connectorId}/sync/stock
```

Request Body:
```json
{
  "batchSize": number,
  "productIds": ["string"],
  "filters": {}
}
```

#### Sync Prices
```http
POST /connectors/{connectorId}/sync/prices
```

#### Sync Customers
```http
POST /connectors/{connectorId}/sync/customers
```

### 3. Order Management

#### Send Order to ERP
```http
POST /connectors/{connectorId}/orders
```

Request Body:
```json
{
  "order": {
    "id": "string",
    "orderNumber": "string",
    "customerId": "string",
    "totalAmount": 0,
    "currency": "string",
    "status": "string",
    "items": [
      {
        "sku": "string",
        "quantity": 0,
        "unitPrice": 0,
        "totalPrice": 0
      }
    ],
    "shippingInfo": {
      "method": "string",
      "carrier": "string",
      "trackingNumber": "string"
    }
  }
}
```

### 4. Invoice Management

#### Create Invoice
```http
POST /invoices
```

Request Body:
```json
{
  "invoiceNumber": "string",
  "orderId": "string",
  "customerId": "string",
  "type": "standard|credit|debit",
  "issueDate": "datetime",
  "dueDate": "datetime",
  "currency": "string",
  "lines": [
    {
      "description": "string",
      "quantity": 0,
      "unitPrice": 0,
      "amount": 0,
      "taxRate": 0
    }
  ],
  "subTotal": 0,
  "taxAmount": 0,
  "discountAmount": 0,
  "shippingAmount": 0,
  "totalAmount": 0,
  "paymentTerms": "string",
  "billingAddress": {
    "street": "string",
    "city": "string",
    "state": "string",
    "postalCode": "string",
    "country": "string"
  }
}
```

#### Approve Invoice
```http
POST /invoices/{invoiceId}/approve
```

Request Body:
```json
{
  "approvedBy": "string"
}
```

#### Send Invoice
```http
POST /invoices/{invoiceId}/send
```

Request Body:
```json
{
  "method": "email|print|api",
  "recipientEmail": "string"
}
```

#### Void Invoice
```http
POST /invoices/{invoiceId}/void
```

Request Body:
```json
{
  "reason": "string",
  "voidedBy": "string"
}
```

#### Record Invoice Payment
```http
POST /invoices/{invoiceId}/payments
```

Request Body:
```json
{
  "amount": 0,
  "paymentMethod": "string",
  "transactionId": "string",
  "paymentDate": "datetime"
}
```

### 5. Return Management

#### Create Return
```http
POST /returns
```

Request Body:
```json
{
  "returnNumber": "string",
  "orderId": "string",
  "customerId": "string",
  "reason": "defective|wrong_item|not_as_described|other",
  "type": "refund|exchange|store_credit",
  "items": [
    {
      "sku": "string",
      "quantity": 0,
      "reason": "string"
    }
  ],
  "refundAmount": 0
}
```

#### Approve Return
```http
POST /returns/{returnId}/approve
```

#### Reject Return
```http
POST /returns/{returnId}/reject
```

#### Process Return
```http
POST /returns/{returnId}/process
```

#### Complete Return
```http
POST /returns/{returnId}/complete
```

#### Restock Return Items
```http
POST /returns/{returnId}/restock
```

### 6. Inventory Management

#### Create Inventory Reservation
```http
POST /inventory/reservations
```

Request Body:
```json
{
  "orderId": "string",
  "sku": "string",
  "quantity": 0,
  "warehouseId": "string",
  "type": "sales_order|transfer|hold"
}
```

#### Release Inventory Reservation
```http
POST /inventory/reservations/{reservationId}/release
```

#### Fulfill Inventory Reservation
```http
POST /inventory/reservations/{reservationId}/fulfill
```

#### Transfer Inventory Reservation
```http
POST /inventory/reservations/{reservationId}/transfer
```

### 7. Webhook Endpoint

#### Receive Webhook
```http
POST /webhook/{connectorId}
```

**Note**: This is a public endpoint (no authentication required)

Headers:
- `X-Webhook-Signature`: Webhook signature for validation
- `X-Hub-Signature`: Alternative signature header
- `Content-Type`: application/json

Request Body: Raw webhook payload from ERP system

## Error Responses

All endpoints use standard HTTP status codes and return errors in the following format:

```json
{
  "code": 3,
  "message": "Invalid argument: field validation failed",
  "details": [
    {
      "type": "BadRequest",
      "field": "name",
      "description": "Name is required"
    }
  ]
}
```

Common Status Codes:
- `200 OK`: Success
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict
- `500 Internal Server Error`: Server error

## Rate Limiting

The API implements rate limiting per connector:
- Default: 10 requests per second with burst of 20
- Configurable per connector during creation/update

## Webhook Security

Webhooks support multiple signature validation methods:
- HMAC-SHA256 signatures
- OAuth 1.0a signatures (NetSuite)
- Custom signature headers per ERP system

## Testing

Use the provided test script to test all endpoints:

```bash
./scripts/test_endpoints.sh
```

Or test individual endpoints with curl:

```bash
curl -X GET http://localhost:8080/api/erp/connectors \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```