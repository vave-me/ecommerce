#!/bin/bash

# Script to add OTEL environment variables to all Kubernetes service deployments

SERVICES_DIR="/home/szymon/classified/k8s/09-services"
OTEL_VARS='            - name: OTEL_SERVICE_NAME
              value: "SERVICE_NAME_PLACEHOLDER"
            - name: OTEL_EXPORTER_OTLP_ENDPOINT
              value: "http://otel-collector.observability.svc.cluster.local:4317"
            - name: OTEL_EXPORTER_OTLP_INSECURE
              value: "true"
            - name: OTEL_TRACES_EXPORTER
              value: "otlp"'

# List of all services
services=(
    "activity" "assistants" "baskets" "categories" "comments" "cosec"
    "following" "geocoding" "mailer" "media" "messages" "metrics"
    "newsletters" "notifications" "offers" "ordering" "payments"
    "posts" "products" "reviews" "search" "shipping" "support"
    "users" "vectors" "wishlists" "services" "scheduler" "erp" "managers" "merchant"
)

for service in "${services[@]}"; do
    file="$SERVICES_DIR/${service}.yaml"
    if [ -f "$file" ]; then
        echo "Processing $service..."
        # Extract the service name for OTEL_SERVICE_NAME
        service_name="$service"
        
        # Replace placeholder with actual service name
        otel_vars_with_name=$(echo "$OTEL_VARS" | sed "s/SERVICE_NAME_PLACEHOLDER/$service_name/")
        
        echo "Would add to $file:"
        echo "$otel_vars_with_name"
        echo "---"
    fi
done

echo "Script complete. This is a dry run - no files were modified."
echo "To apply changes, modify the script to actually edit the files."