#!/bin/bash
# Production script to create individual dashboard ConfigMaps

DASHBOARD_DIR="/home/szymon/classified/docker/grafana/provisioning/dashboards"
NAMESPACE="observability"

# Create each dashboard as a separate ConfigMap
for dashboard in "$DASHBOARD_DIR"/*.json; do
    if [ -f "$dashboard" ]; then
        filename=$(basename "$dashboard")
        configmap_name="grafana-dashboard-${filename%.json}"
        
        echo "Creating ConfigMap: $configmap_name"
        kubectl create configmap "$configmap_name" \
            --from-file="$filename=$dashboard" \
            --namespace="$NAMESPACE" \
            --dry-run=client -o yaml | kubectl apply -f -
    fi
done

# Update Grafana deployment to mount all dashboard ConfigMaps
echo "Dashboard ConfigMaps created. Update Grafana deployment to mount them."