#!/bin/bash

# Update Go import statements from old module name to new module name

set -e

OLD_MODULE="github.com/sst/opencode-sdk-go"
NEW_MODULE="git.j9xym.com/openapi-api-go"

echo "Updating Go import statements..."
echo "Old module: $OLD_MODULE"
echo "New module: $NEW_MODULE"

# Find and update all Go files in packages/tui
find packages/tui -name "*.go" -type f -exec sed -i "s|$OLD_MODULE|$NEW_MODULE|g" {} +

echo "✅ Import statements updated successfully!"
echo "Updated files:"
find packages/tui -name "*.go" -exec grep -l "$NEW_MODULE" {} \;