#!/bin/bash

# Clean Go workspace and prepare for fresh local build

set -e

echo "🧹 Cleaning Go workspace..."
echo

# Remove existing workspace file
if [ -f "go.work" ]; then
	rm -f go.work go.work.sum
	echo "✅ Removed workspace files"
fi

# Clean module caches
echo "🗑️  Cleaning module caches..."
go clean -modcache 2>/dev/null || true

# Clean local module caches
find . -name "go.sum" -exec rm -f {} \; 2>/dev/null || true

echo "✅ Go workspace cleaned"
echo
echo "Ready for fresh local build setup!"