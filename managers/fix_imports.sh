#!/bin/bash

# Script to fix cross-service imports in managers service

echo "Fixing cross-service imports in managers service..."

# Replace assistants imports with managers imports
find . -name "*.go" -type f -exec sed -i \
  -e 's|"middleman/assistants/internal/domain"|"middleman/managers/internal/domain"|g' \
  -e 's|"middleman/assistants/internal/models"|"middleman/managers/internal/models"|g' \
  -e 's|"middleman/assistants/internal/application/services"|"middleman/managers/internal/application/services"|g' \
  -e 's|"middleman/assistants/internal/application/tools"|"middleman/managers/internal/application/tools"|g' \
  {} \;

echo "Import fixes completed!"

# Show any remaining cross-service imports
echo "Checking for any remaining cross-service imports..."
grep -r "middleman/assistants/internal" --include="*.go" || echo "No cross-service imports found!"