#!/bin/bash

# Send a test trace to the collector
curl -X POST http://localhost:4318/v1/traces \
  -H 'Content-Type: application/json' \
  -d '{
  "resourceSpans": [{
    "resource": {
      "attributes": [{
        "key": "service.name",
        "value": {"stringValue": "test-service"}
      }]
    },
    "scopeSpans": [{
      "scope": {
        "name": "test-instrumentation"
      },
      "spans": [{
        "traceId": "5b8aa5a2d2c872e8321cf37308d69df2",
        "spanId": "051581bf3cb55c13",
        "name": "test-operation",
        "kind": 2,
        "startTimeUnixNano": "'$(date +%s%N)'",
        "endTimeUnixNano": "'$(date -d "+1 second" +%s%N)'",
        "attributes": [
          {
            "key": "http.method",
            "value": {"stringValue": "GET"}
          },
          {
            "key": "http.status_code",
            "value": {"intValue": "200"}
          }
        ],
        "status": {}
      }]
    }]
  }]
}'

echo "Test trace sent!"