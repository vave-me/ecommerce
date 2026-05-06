#!/bin/bash

# ERP Service REST API Test Script
# This script tests all ERP service endpoints

# Configuration
BASE_URL="http://localhost:8080/api/erp"
AUTH_TOKEN=REDACTED

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to make authenticated requests
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    echo -e "\n${BLUE}=== Testing: $method $endpoint ===${NC}"
    
    if [ -z "$data" ]; then
        curl -X "$method" "$BASE_URL$endpoint" \
            -H "Authorization: Bearer $AUTH_TOKEN" \
            -H "Content-Type: application/json" \
            -w "\nHTTP Status: %{http_code}\n" \
            -s | jq '.' 2>/dev/null || echo "Response parsing failed"
    else
        curl -X "$method" "$BASE_URL$endpoint" \
            -H "Authorization: Bearer $AUTH_TOKEN" \
            -H "Content-Type: application/json" \
            -d "$data" \
            -w "\nHTTP Status: %{http_code}\n" \
            -s | jq '.' 2>/dev/null || echo "Response parsing failed"
    fi
}

echo -e "${GREEN}Starting ERP Service API Tests${NC}"

# 1. Connector Management Tests
echo -e "\n${GREEN}1. Connector Management${NC}"

# List all connectors
make_request GET "/connectors"

# Create a test NetSuite connector
echo -e "\n${BLUE}Creating NetSuite test connector...${NC}"
NETSUITE_DATA='{
  "name": "Test NetSuite Connector",
  "type": "netsuite",
  "environment": "sandbox",
  "baseUrl": "https://test.netsuite.com",
  "authType": "oauth1",
  "authConfig": {
    "consumer_key": "test_consumer_key",
    "consumer_secret": "test_consumer_secret",
    "token_id": "test_token_id",
    "token_secret": "test_token_secret",
    "account_id": "test_account"
  },
  "webhookEnabled": true,
  "syncEnabled": true,
  "syncIntervalSeconds": 300,
  "batchSize": 100
}'
RESPONSE=$(make_request POST "/connectors/v2" "$NETSUITE_DATA")
CONNECTOR_ID=$(echo "$RESPONSE" | jq -r '.connectorId' 2>/dev/null)

if [ "$CONNECTOR_ID" != "null" ] && [ -n "$CONNECTOR_ID" ]; then
    echo -e "${GREEN}Created connector with ID: $CONNECTOR_ID${NC}"
    
    # Get connector status
    make_request GET "/connectors/$CONNECTOR_ID/status"
    
    # Get connector health
    make_request GET "/connectors/$CONNECTOR_ID/health"
    
    # Toggle connector (deactivate)
    make_request POST "/connectors/$CONNECTOR_ID/toggle" '{"activate": false, "reason": "Testing toggle"}'
    
    # Toggle connector (activate)
    make_request POST "/connectors/$CONNECTOR_ID/toggle" '{"activate": true, "reason": "Re-enabling after test"}'
fi

# Create an Odoo test connector
echo -e "\n${BLUE}Creating Odoo test connector...${NC}"
ODOO_DATA='{
  "name": "Test Odoo Connector",
  "type": "odoo",
  "environment": "development",
  "baseUrl": "https://demo.odoo.com",
  "authType": "basic",
  "authConfig": {
    "username": "demo",
    "password": "demo",
    "database": "demo_db"
  },
  "webhookEnabled": false,
  "syncEnabled": true,
  "syncIntervalSeconds": 600,
  "batchSize": 50
}'
make_request POST "/connectors/v2" "$ODOO_DATA"

# List connectors with filters
make_request GET "/connectors?type=netsuite"
make_request GET "/connectors?status=active"

# 2. Synchronization Tests
echo -e "\n${GREEN}2. Synchronization Operations${NC}"

if [ "$CONNECTOR_ID" != "null" ] && [ -n "$CONNECTOR_ID" ]; then
    # Get sync history
    make_request GET "/connectors/$CONNECTOR_ID/sync-history"
    
    # Sync products
    make_request POST "/connectors/$CONNECTOR_ID/sync/products" '{"batchSize": 50}'
    
    # Sync stock
    make_request POST "/connectors/$CONNECTOR_ID/sync/stock" '{"productIds": ["SKU001", "SKU002"]}'
    
    # Sync prices
    make_request POST "/connectors/$CONNECTOR_ID/sync/prices" '{}'
    
    # Sync customers
    make_request POST "/connectors/$CONNECTOR_ID/sync/customers" '{}'
fi

# 3. Invoice Management Tests
echo -e "\n${GREEN}3. Invoice Management${NC}"

# Create an invoice
INVOICE_DATA='{
  "invoiceNumber": "INV-TEST-001",
  "orderId": "ORD-123",
  "customerId": "CUST-456",
  "type": "standard",
  "issueDate": "2024-01-24T00:00:00Z",
  "dueDate": "2024-02-24T00:00:00Z",
  "currency": "USD",
  "lines": [
    {
      "description": "Test Product",
      "quantity": 2,
      "unitPrice": 50.00,
      "amount": 100.00
    }
  ],
  "subTotal": 100.00,
  "taxAmount": 10.00,
  "totalAmount": 110.00,
  "paymentTerms": "Net 30",
  "billingAddress": {
    "street": "123 Test St",
    "city": "Test City",
    "state": "TS",
    "postalCode": "12345",
    "country": "US"
  }
}'
INVOICE_RESPONSE=$(make_request POST "/invoices" "$INVOICE_DATA")
INVOICE_ID=$(echo "$INVOICE_RESPONSE" | jq -r '.invoiceId' 2>/dev/null)

if [ "$INVOICE_ID" != "null" ] && [ -n "$INVOICE_ID" ]; then
    echo -e "${GREEN}Created invoice with ID: $INVOICE_ID${NC}"
    
    # Approve invoice
    make_request POST "/invoices/$INVOICE_ID/approve" '{"approvedBy": "admin"}'
    
    # Send invoice
    make_request POST "/invoices/$INVOICE_ID/send" '{"method": "email", "recipientEmail": "redacted-email@example.com"}'
    
    # Record payment
    make_request POST "/invoices/$INVOICE_ID/payments" '{
      "amount": 110,
      "paymentMethod": "credit_card",
      "transactionId": "TXN-123",
      "paymentDate": "2024-01-25T00:00:00Z"
    }'
fi

# 4. Return Management Tests
echo -e "\n${GREEN}4. Return Management${NC}"

# Create a return
RETURN_DATA='{
  "returnNumber": "RET-TEST-001",
  "orderId": "ORD-123",
  "customerId": "CUST-456",
  "reason": "defective",
  "type": "exchange",
  "items": [
    {
      "sku": "SKU001",
      "quantity": 1,
      "reason": "Defective product"
    }
  ],
  "refundAmount": 55.00
}'
RETURN_RESPONSE=$(make_request POST "/returns" "$RETURN_DATA")
RETURN_ID=$(echo "$RETURN_RESPONSE" | jq -r '.returnId' 2>/dev/null)

if [ "$RETURN_ID" != "null" ] && [ -n "$RETURN_ID" ]; then
    echo -e "${GREEN}Created return with ID: $RETURN_ID${NC}"
    
    # Approve return
    make_request POST "/returns/$RETURN_ID/approve" '{"approvedBy": "admin"}'
    
    # Process return
    make_request POST "/returns/$RETURN_ID/process" '{"erpReturnId": "ERP-RET-123"}'
    
    # Restock items
    make_request POST "/returns/$RETURN_ID/restock" '{
      "items": [
        {
          "sku": "SKU001",
          "quantity": 1,
          "locationId": "WAREHOUSE-1"
        }
      ]
    }'
    
    # Complete return
    make_request POST "/returns/$RETURN_ID/complete" '{"refundTransactionId": "REFUND-123"}'
fi

# 5. Inventory Management Tests
echo -e "\n${GREEN}5. Inventory Management${NC}"

# Create inventory reservation
RESERVATION_DATA='{
  "orderId": "ORD-789",
  "sku": "SKU001",
  "quantity": 5,
  "warehouseId": "WAREHOUSE-1",
  "type": "sales_order"
}'
RESERVATION_RESPONSE=$(make_request POST "/inventory/reservations" "$RESERVATION_DATA")
RESERVATION_ID=$(echo "$RESERVATION_RESPONSE" | jq -r '.reservationId' 2>/dev/null)

if [ "$RESERVATION_ID" != "null" ] && [ -n "$RESERVATION_ID" ]; then
    echo -e "${GREEN}Created reservation with ID: $RESERVATION_ID${NC}"
    
    # Transfer reservation
    make_request POST "/inventory/reservations/$RESERVATION_ID/transfer" '{"toWarehouseId": "WAREHOUSE-2"}'
    
    # Fulfill reservation
    make_request POST "/inventory/reservations/$RESERVATION_ID/fulfill" '{}'
    
    # Create another reservation to test release
    TEMP_RESERVATION_RESPONSE=$(make_request POST "/inventory/reservations" "$RESERVATION_DATA")
    TEMP_RESERVATION_ID=$(echo "$TEMP_RESERVATION_RESPONSE" | jq -r '.reservationId' 2>/dev/null)
    
    if [ "$TEMP_RESERVATION_ID" != "null" ] && [ -n "$TEMP_RESERVATION_ID" ]; then
        # Release reservation
        make_request POST "/inventory/reservations/$TEMP_RESERVATION_ID/release" '{"reason": "Order cancelled"}'
    fi
fi

# 6. Webhook Test (public endpoint)
echo -e "\n${GREEN}6. Webhook Endpoint Test${NC}"

if [ "$CONNECTOR_ID" != "null" ] && [ -n "$CONNECTOR_ID" ]; then
    echo -e "\n${BLUE}Testing webhook endpoint (no auth required)...${NC}"
    WEBHOOK_DATA='{
      "event": "product.updated",
      "data": {
        "sku": "SKU001",
        "name": "Updated Product",
        "price": 99.99
      }
    }'
    
    curl -X POST "$BASE_URL/webhook/$CONNECTOR_ID" \
        -H "Content-Type: application/json" \
        -H "X-Webhook-Signature: test-signature" \
        -d "$WEBHOOK_DATA" \
        -w "\nHTTP Status: %{http_code}\n" \
        -s | jq '.' 2>/dev/null || echo "Response parsing failed"
fi

# Cleanup (optional)
echo -e "\n${GREEN}Cleanup (optional)${NC}"
echo "To remove test connectors, uncomment the following lines:"
echo "# make_request DELETE \"/connectors/$CONNECTOR_ID?force=true\""

echo -e "\n${GREEN}ERP Service API Tests Completed!${NC}"