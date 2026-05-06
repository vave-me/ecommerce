#!/bin/bash
# Production deployment script for Grafana with dashboards

set -e

echo "=== Grafana Production Deployment ==="

# 1. Create dashboard ConfigMaps
echo "Creating dashboard ConfigMaps..."
./create-dashboard-configmaps.sh

# 2. Apply updated Grafana configuration
echo "Applying Grafana configuration..."
kubectl apply -f grafana.yaml

# 3. Apply subdomain ingress
echo "Applying subdomain ingress..."
kubectl apply -f grafana-subdomain-ingress.yaml

# 4. Restart Grafana to pick up all changes
echo "Restarting Grafana deployment..."
kubectl rollout restart deployment/grafana -n observability

# 5. Wait for rollout to complete
echo "Waiting for rollout to complete..."
kubectl rollout status deployment/grafana -n observability

echo "=== Deployment Complete ==="
echo "Grafana available at: https://grafana.sfx-markt.de"
echo "Dashboards will be available in the 'Docker Import' folder"