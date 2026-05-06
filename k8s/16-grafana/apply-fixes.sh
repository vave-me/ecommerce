#!/bin/bash
# Apply Grafana fixes

echo "Applying updated Grafana configuration..."
kubectl apply -f grafana.yaml

echo "Ensuring subdomain ingress is applied..."
kubectl apply -f grafana-subdomain-ingress.yaml

echo "Creating dashboard ConfigMaps..."
./create-dashboard-configmaps.sh

echo "Restarting Grafana..."
kubectl rollout restart deployment/grafana -n observability

echo "Waiting for rollout..."
kubectl rollout status deployment/grafana -n observability

echo "Checking Grafana pod status..."
kubectl get pods -n observability | grep grafana

echo "Checking ingress status..."
kubectl get ingress -n observability | grep grafana