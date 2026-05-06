#!/bin/bash
# Script to create Grafana dashboards ConfigMap from Docker directory

DOCKER_DASHBOARDS_DIR="/home/szymon/classified/docker/grafana/provisioning/dashboards"
OUTPUT_FILE="/home/szymon/classified/k8s/16-grafana/grafana-dashboards-configmap.yaml"

# Start the ConfigMap
cat > "$OUTPUT_FILE" << 'EOF'
# ConfigMap containing all Grafana dashboards from Docker directory
# Generated from Docker dashboard files
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-dashboards
  namespace: observability
  labels:
    app: grafana
data:
  # Dashboard provisioning configuration
  dashboards-provider.yaml: |
    apiVersion: 1
    providers:
      - name: 'Imported Dashboards'
        orgId: 1
        folder: 'Docker Import'
        type: file
        disableDeletion: false
        updateIntervalSeconds: 10
        allowUiUpdates: true
        options:
          path: /var/lib/grafana/dashboards
EOF

# Process each JSON dashboard file
for dashboard in "$DOCKER_DASHBOARDS_DIR"/*.json; do
    if [ -f "$dashboard" ]; then
        filename=$(basename "$dashboard")
        echo "  # Dashboard: $filename" >> "$OUTPUT_FILE"
        echo "  $filename: |" >> "$OUTPUT_FILE"
        # Indent each line with 4 spaces for YAML format
        sed 's/^/    /' "$dashboard" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
    fi
done

echo "ConfigMap created at: $OUTPUT_FILE"
echo "Dashboard files included:"
ls -la "$DOCKER_DASHBOARDS_DIR"/*.json | wc -l